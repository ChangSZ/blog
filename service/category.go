package service

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/ChangSZ/blog/common"
	"github.com/ChangSZ/blog/conf"
	"github.com/ChangSZ/blog/infra/log"
	"github.com/ChangSZ/blog/model"
	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

func GetCateById(ctx context.Context, cateId int) (cate *model.Categories, err error) {
	cate = new(model.Categories)
	err = conf.SqlServer.WithContext(ctx).Where("id=?", cateId).Find(&cate).Error
	return
}

func GetCateByParentId(ctx context.Context, parentId int) (cate *model.Categories, err error) {
	cate = new(model.Categories)
	err = conf.SqlServer.WithContext(ctx).Where("parent_id = ?", parentId).Find(&cate).Error
	return
}

func DelCateRel(ctx context.Context, cateId int) {
	var err error
	err = conf.SqlServer.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		postCate := new(model.PostCate)
		if err = tx.Where("cate_id = ?", cateId).Delete(&postCate).Error; err != nil {
			return err
		}

		cate := &model.Categories{ID: cateId}
		if err = tx.Delete(&cate).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.WithTrace(ctx).Error(err)
		return
	}
	conf.CacheClient.Del(conf.Cnf.CateListKey)
	return
}

func CateStore(ctx context.Context, cs common.CateStore) (bool, error) {

	defaultCate := new(model.Categories)
	err := conf.SqlServer.WithContext(ctx).Where("name = ?", cs.Name).Find(&defaultCate).Error
	if err != nil {
		log.WithTrace(ctx).Error(err)
		return false, err
	}
	if defaultCate.ID > 0 {
		log.WithTrace(ctx).Errorf("Cate has exists,  ID=%v", defaultCate.ID)
		return false, errors.New("Cate has exists ")
	}

	if cs.ParentId > 0 {
		cate := new(model.Categories)
		err = conf.SqlServer.WithContext(ctx).Find(cate, cs.ParentId).Error
		if err != nil {
			log.WithTrace(ctx).Error(err)
			return false, err
		}
		if cate.ID <= 0 {
			log.WithTrace(ctx).Error("The parent id has not data")
			return false, errors.New("The parent id has not data ")
		}
	}

	cate := &model.Categories{
		Name:        cs.Name,
		DisplayName: cs.DisplayName,
		SeoDesc:     cs.SeoDesc,
		ParentID:    cs.ParentId,
	}
	err = conf.SqlServer.WithContext(ctx).Create(cate).Error
	if err != nil {
		log.WithTrace(ctx).Error(err)
		return false, err
	}
	conf.CacheClient.Del(conf.Cnf.CateListKey)
	return true, nil
}

func CateUpdate(ctx context.Context, cateId int, cs common.CateStore) (bool, error) {
	cate := new(model.Categories)
	if cs.ParentId != 0 {
		err := conf.SqlServer.WithContext(ctx).First(cate, cs.ParentId).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				log.WithTrace(ctx).Error("the parent id is not exists ")
				return false, errors.New("the parent id is not exists ")
			}

			log.WithTrace(ctx).Error(err)
			return false, err
		}

		ids := []int{cateId}
		resIds := []int{0}
		_, res2, _ := GetSimilar(ctx, ids, resIds, 0)
		for _, v := range res2 {
			if v == cs.ParentId {
				return false, errors.New("Can not be you child node ")
			}
		}
	}
	cateUpdate := &model.Categories{
		ID:          cateId,
		Name:        cs.Name,
		DisplayName: cs.DisplayName,
		SeoDesc:     cs.SeoDesc,
		ParentID:    cs.ParentId,
	}
	err := conf.SqlServer.WithContext(ctx).Updates(cateUpdate).Error
	if err != nil {
		log.WithTrace(ctx).Error(err)
		return false, err
	}
	conf.CacheClient.Del(conf.Cnf.CateListKey)
	return true, nil
}

func GetSimilar(ctx context.Context, beginId []int, resIds []int, level int) (beginId2 []int, resIds2 []int, level2 int) {
	if len(beginId) != 0 {
		cates := make([]*model.Categories, 0)
		err := conf.SqlServer.WithContext(ctx).Where("parent_id in ?", beginId).Find(&cates)
		if err != nil {
			log.WithTrace(ctx).Error(err, ", the parent id data is not exists")
			return []int{}, []int{}, 0
		}
		if len(cates) != 0 {
			if level == 0 {
				resIds2 = beginId
			} else {
				resIds2 = resIds
			}
			for _, v := range cates {
				id := v.ID
				beginId2 = append(beginId2, id)
				resIds2 = append(resIds2, id)
			}
			level2 = level + 1
			return GetSimilar(ctx, beginId2, resIds2, level2)
		}
		if level == 0 && len(cates) == 0 {
			return beginId, beginId, level
		} else {
			return beginId, resIds, level
		}
	}
	return beginId, resIds, level
}

func GetPostCateByPostId(ctx context.Context, postId int) (cates *model.Categories, err error) {
	postCate := new(model.PostCate)
	err = conf.SqlServer.WithContext(ctx).Where("post_id = ?", postId).Find(&postCate).Error
	if err != nil {
		log.WithTrace(ctx).Error(err)
		return cates, err
	}
	if postCate.ID != 0 {
		cates = new(model.Categories)
		err = conf.SqlServer.WithContext(ctx).Where("id = ?", postCate.CateID).Find(&cates).Error
		if err != nil {
			log.WithTrace(ctx).Error(err)
			return cates, err
		}
		if cates.ID == 0 {
			log.WithTrace(ctx).Error("there has not data")
			return cates, errors.New("can not get the post cate")
		}
	} else {
		log.WithTrace(ctx).Error("there has not data")
		return cates, errors.New("can not get the post cate")
	}

	return cates, nil

}

// Get the cate list what by parent sort
func CateListBySort(ctx context.Context) ([]common.Category, error) {
	cacheKey := conf.Cnf.CateListKey
	cacheRes, err := conf.CacheClient.Get(cacheKey).Result()
	if err == redis.Nil {
		// cache key does not exist
		// set data to the cache what use the cache key
		cates, err := doCacheCateList(ctx, cacheKey)
		if err != nil {
			log.WithTrace(ctx).Error(err)
			return nil, err
		}
		return cates, nil
	} else if err != nil {
		log.WithTrace(ctx).Error(err)
		return nil, err
	}

	if cacheRes == "" {
		// Data is  null
		// set data to the cache what use the cache key
		cates, err := doCacheCateList(ctx, cacheKey)
		if err != nil {
			log.WithTrace(ctx).Error(err)
			return nil, err
		}
		return cates, nil
	}

	var comCate []common.Category
	err = json.Unmarshal([]byte(cacheRes), &comCate)
	if err != nil {
		log.WithTrace(ctx).Error(err)
		cates, err := doCacheCateList(ctx, cacheKey)
		if err != nil {
			log.WithTrace(ctx).Error(err)
			return nil, err
		}
		return cates, nil
	}
	return comCate, nil
}

// Get the all cate
// then set it to cache
func doCacheCateList(ctx context.Context, cacheKey string) ([]common.Category, error) {
	allCates, err := allCates(ctx)
	if err != nil {
		log.WithTrace(ctx).Error(err)
		return nil, err
	}
	var cate common.Category
	var cates []common.Category
	for _, v := range allCates {
		cate.Cates = v
		cates = append(cates, cate)
	}
	res := tree(ctx, cates, 0, 0, 0)
	jsonRes, err := json.Marshal(&res)
	if err != nil {
		log.WithTrace(ctx).Error(err)
		return nil, err
	}
	err = conf.CacheClient.Set(cacheKey, jsonRes, time.Duration(conf.Cnf.DataCacheTimeDuration)*time.Hour).Err()
	if err != nil {
		log.WithTrace(ctx).Error(err)
		return nil, err
	}
	return res, nil
}

// data recursion
func tree(ctx context.Context, cate []common.Category, parent int, level int, key int) []common.Category {
	html := "-"
	var data []common.Category
	for _, v := range cate {
		var ParentId = v.Cates.ParentID
		var Id = v.Cates.ID
		if ParentId == parent {
			var newHtml string
			if level != 0 {
				newHtml = common.GoRepeat("&nbsp;&nbsp;&nbsp;&nbsp;", level) + "|"
			}
			v.Html = newHtml + common.GoRepeat(html, level)
			data = append(data, v)
			data = merge(data, tree(ctx, cate, Id, level+1, key+1))
		}
	}
	return data
}

// merge two arr
func merge(arr1 []common.Category, arr2 []common.Category) []common.Category {
	for _, val := range arr2 {
		arr1 = append(arr1, val)
	}
	return arr1
}

// Get all cate
// if not exists
// create the default one
func allCates(ctx context.Context) ([]model.Categories, error) {
	cates := make([]model.Categories, 0)
	err := conf.SqlServer.WithContext(ctx).Find(&cates).Error
	if err != nil {
		log.WithTrace(ctx).Error(err)
		return cates, err
	}

	if len(cates) == 0 {
		cateCreate := &model.Categories{
			Name:        "default",
			DisplayName: "默认分类",
			SeoDesc:     "默认的分类",
			ParentID:    0,
		}
		err := conf.SqlServer.WithContext(ctx).Create(&cateCreate).Error
		if err != nil {
			log.WithTrace(ctx).Error(err)
			return cates, err
		}

		err = conf.SqlServer.WithContext(ctx).Find(&cates).Error
		if err != nil {
			log.WithTrace(ctx).Error(err)
			return cates, err
		}

		return cates, nil
	}

	return cates, nil
}

func CateCnt(ctx context.Context) (cnt int64, err error) {
	cate := new(model.Categories)
	err = conf.SqlServer.WithContext(ctx).Table(cate.TableName()).Count(&cnt).Error
	return
}
