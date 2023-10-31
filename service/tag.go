package service

import (
	"context"
	"encoding/json"
	"time"

	"github.com/ChangSZ/blog/common"
	"github.com/ChangSZ/blog/conf"
	"github.com/ChangSZ/blog/infra/log"
	"github.com/ChangSZ/blog/model"
	"github.com/go-errors/errors"
	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

func TagStore(ctx context.Context, ts common.TagStore) (err error) {
	tag := new(model.Tags)
	err = conf.SqlServer.WithContext(ctx).Where("name = ?", ts.Name).Find(&tag).Error
	if err != nil {
		log.WithTrace(ctx).Error(err)
		return err
	}

	if tag.ID > 0 {
		log.WithTrace(ctx).Errorf("Tag has exists, ID=%v", tag.ID)
		return errors.New("Tag has exists")
	}

	tagInsert := &model.Tags{
		Name:        ts.Name,
		DisplayName: ts.DisplayName,
		SeoDesc:     ts.SeoDesc,
		Num:         0,
	}
	err = conf.SqlServer.WithContext(ctx).Create(tagInsert).Error
	conf.CacheClient.Del(conf.Cnf.TagListKey)
	return
}

func GetPostTagsByPostId(ctx context.Context, postId int) (tagsArr []int, err error) {
	postTags := make([]model.PostTag, 0)
	err = conf.SqlServer.WithContext(ctx).Where("post_id = ?", postId).Find(&postTags).Error
	if err != nil {
		log.WithTrace(ctx).Error(err)
		return nil, nil
	}

	for _, postTag := range postTags {
		tagsArr = append(tagsArr, postTag.TagID)
	}
	return
}

func GetTagById(ctx context.Context, tagId int) (tag *model.Tags, err error) {
	tag = new(model.Tags)
	err = conf.SqlServer.WithContext(ctx).Where("id=?", tagId).Find(&tag).Error
	return
}

func TagUpdate(ctx context.Context, tagId int, ts common.TagStore) error {
	tagUpdate := &model.Tags{
		ID:          tagId,
		Name:        ts.Name,
		DisplayName: ts.DisplayName,
		SeoDesc:     ts.SeoDesc,
	}
	err := conf.SqlServer.WithContext(ctx).Updates(tagUpdate).Error
	return err
}

func GetTagsByIds(ctx context.Context, tagIds []int) ([]*model.Tags, error) {
	tags := make([]*model.Tags, 0)
	err := conf.SqlServer.WithContext(ctx).Where("id in ?", tagIds).Find(&tags).Error
	if err != nil {
		return nil, err
	}
	return tags, nil
}

func TagsIndex(ctx context.Context, limit int, offset int) (num int64, tags []*model.Tags, err error) {
	tags = make([]*model.Tags, 0)
	db := conf.SqlServer.WithContext(ctx).Table((&model.Tags{}).TableName()).Order("num desc").Offset(offset).Limit(limit)
	if err = db.Count(&num).Error; err != nil {
		return
	}
	err = db.Find(&tags).Error
	return
}

func DelTagRel(ctx context.Context, tagId int) {
	var err error
	err = conf.SqlServer.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		postTag := new(model.PostTag)
		if err = tx.Where("tag_id = ?", tagId).Delete(&postTag).Error; err != nil {
			return err
		}

		tag := &model.Tags{ID: tagId}
		if err = tx.Delete(&tag).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.WithTrace(ctx).Error(err)
		return
	}

	conf.CacheClient.Del(conf.Cnf.TagListKey)
	return
}

func AllTags(ctx context.Context) ([]model.Tags, error) {
	cacheKey := conf.Cnf.TagListKey
	cacheRes, err := conf.CacheClient.Get(cacheKey).Result()
	if err == redis.Nil {
		tags, err := doCacheTagList(ctx, cacheKey)
		if err != nil {
			log.WithTrace(ctx).Error(err)
			return tags, err
		}
		return tags, nil
	} else if err != nil {
		log.WithTrace(ctx).Error(err)
		return nil, err
	}

	var cacheTag []model.Tags
	err = json.Unmarshal([]byte(cacheRes), &cacheTag)
	if err != nil {
		log.WithTrace(ctx).Error(err)
		tags, err := doCacheTagList(ctx, cacheKey)
		if err != nil {
			log.WithTrace(ctx).Error(err)
			return nil, err
		}
		return tags, nil
	}
	return cacheTag, nil
}

func doCacheTagList(ctx context.Context, cacheKey string) ([]model.Tags, error) {
	tags, err := tags(ctx)
	if err != nil {
		log.WithTrace(ctx).Error(err)
		return tags, err
	}
	jsonRes, err := json.Marshal(&tags)
	if err != nil {
		log.WithTrace(ctx).Error(err)
		return nil, err
	}
	err = conf.CacheClient.Set(cacheKey, jsonRes, time.Duration(conf.Cnf.DataCacheTimeDuration)*time.Hour).Err()
	if err != nil {
		log.WithTrace(ctx).Error(err)
		return nil, err
	}
	return tags, nil
}

func tags(ctx context.Context) ([]model.Tags, error) {
	tags := make([]model.Tags, 0)
	err := conf.SqlServer.WithContext(ctx).Find(&tags).Error
	if err != nil {
		log.WithTrace(ctx).Error(err)
		return tags, err
	}

	return tags, nil
}

func TagCnt(ctx context.Context) (cnt int64, err error) {
	tag := new(model.Tags)
	err = conf.SqlServer.WithContext(ctx).Table(tag.TableName()).Count(&cnt).Error
	return
}
