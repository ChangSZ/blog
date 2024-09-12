package service

import (
	"context"
	"encoding/json"
	"html/template"
	"strconv"

	"github.com/ChangSZ/golib/log"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/gorilla/feeds"
	"gorm.io/gorm"

	"github.com/ChangSZ/blog/common"
	"github.com/ChangSZ/blog/conf"
	"github.com/ChangSZ/blog/model"
)

type IndexType string

const (
	IndexTypeTag     IndexType = "tag"
	IndexTypeCate    IndexType = "cate"
	IndexTypeDefault IndexType = "default"
)

func CommonData(ctx context.Context) (h gin.H, err error) {
	h = gin.H{
		"themeJs":          conf.Cnf.ThemeJs,
		"themeCss":         conf.Cnf.ThemeCss,
		"themeImg":         conf.Cnf.ThemeImg,
		"themeFancyboxCss": conf.Cnf.ThemeFancyboxCss,
		"themeFancyboxJs":  conf.Cnf.ThemeFancyboxJs,
		"themeShareCss":    conf.Cnf.ThemeShareCss,
		"themeShareJs":     conf.Cnf.ThemeShareJs,
		"themeArchivesJs":  conf.Cnf.ThemeArchivesJs,
		"themeArchivesCss": conf.Cnf.ThemeArchivesCss,
		"themeNiceImg":     conf.Cnf.ThemeNiceImg,
		"themeAllCss":      conf.Cnf.ThemeAllCss,
		"themeIndexImg":    conf.Cnf.ThemeIndexImg,
		"themeCateImg":     conf.Cnf.ThemeCateImg,
		"themeTagImg":      conf.Cnf.ThemeTagImg,
		"title":            "",

		"tem": "defaultList",
	}
	h["script"] = template.HTML(conf.Cnf.OtherScript)
	cates, err := CateListBySort(ctx)
	if err != nil {
		log.WithTrace(ctx).Error(err)
		return
	}
	var catess []common.IndexCategory
	for _, v := range cates {
		c := common.IndexCategory{
			Cates: v.Cates,
			Html:  template.HTML(v.Html),
		}
		catess = append(catess, c)
	}

	tags, err := AllTags(ctx)
	if err != nil {
		log.WithTrace(ctx).Error(err)
		return
	}

	links, err := AllLink(ctx)
	if err != nil {
		log.WithTrace(ctx).Error(err)
		return
	}

	system, err := IndexSystem(ctx)
	if err != nil {
		log.WithTrace(ctx).Error(err)
		return
	}
	h["cates"] = catess
	h["system"] = system
	h["links"] = links
	h["tags"] = tags
	return
}

func IndexPost(ctx context.Context, page string, limit string, indexType IndexType, name string) (indexPostIndex common.IndexPostList, err error) {
	var postKey string
	switch indexType {
	case IndexTypeTag:
		postKey = conf.Cnf.TagPostIndexKey
	case IndexTypeCate:
		postKey = conf.Cnf.CatePostIndexKey
	case IndexTypeDefault:
		postKey = conf.Cnf.PostIndexKey
		name = "default"
	default:
		postKey = conf.Cnf.PostIndexKey
	}

	field := ":name:" + name + ":page:" + page + ":limit:" + limit
	cacheRes, err := conf.CacheClient.HGet(postKey, field).Result()
	if err == redis.Nil {
		// cache key does not exist
		// set data to the cache what use the cache key
		indexPostIndex, err := doCacheIndexPostList(ctx, postKey, field, indexType, name, page, limit)
		if err != nil {
			log.WithTrace(ctx).Error(err)
			return indexPostIndex, err
		}
		return indexPostIndex, nil
	} else if err != nil {
		log.WithTrace(ctx).Error(err)
		return indexPostIndex, err
	}

	if cacheRes == "" {
		// Data is  null
		// set data to the cache what use the cache key
		indexPostIndex, err := doCacheIndexPostList(ctx, postKey, field, indexType, name, page, limit)
		if err != nil {
			log.WithTrace(ctx).Error(err)
			return indexPostIndex, err
		}
		return indexPostIndex, nil
	}

	err = json.Unmarshal([]byte(cacheRes), &indexPostIndex)
	if err != nil {
		log.WithTrace(ctx).Error(err)
		indexPostIndex, err := doCacheIndexPostList(ctx, postKey, field, indexType, name, page, limit)
		if err != nil {
			log.WithTrace(ctx).Error(err)
			return indexPostIndex, err
		}
		return indexPostIndex, nil
	}
	return
}

func doCacheIndexPostList(ctx context.Context, cacheKey string, field string, indexType IndexType, name string, queryPage string, queryLimit string) (res common.IndexPostList, err error) {
	limit, offset := common.Offset(queryPage, queryLimit)
	queryPageInt, err := strconv.Atoi(queryPage)
	if err != nil {
		log.WithTrace(ctx).Error(err)
		return
	}
	var postList []*common.ConsolePostList
	var postCount int64
	switch indexType {
	case IndexTypeTag:
		tag := new(model.Tags)
		err = conf.SqlServer.WithContext(ctx).Where("Name = ?", name).Find(&tag).Error
		if err != nil {
			log.WithTrace(ctx).Error(err)
			return
		}
		postList, err = PostTagList(ctx, tag.ID, limit, offset)
		if err != nil {
			log.WithTrace(ctx).Error(err)
			return
		}
		postCount, err = PostTagListCount(ctx, tag.ID)
		if err != nil {
			log.WithTrace(ctx).Error(err)
			return
		}
	case IndexTypeCate:
		cate := new(model.Categories)
		err = conf.SqlServer.WithContext(ctx).Where("Name = ?", name).Find(&cate).Error
		if err != nil {
			log.WithTrace(ctx).Error(err)
			return
		}
		postList, err = PostCateList(ctx, cate.ID, limit, offset)
		if err != nil {
			log.WithTrace(ctx).Error(err)
			return
		}
		postCount, err = PostCateListCount(ctx, cate.ID)
		if err != nil {
			log.WithTrace(ctx).Error(err)
			return
		}
	case IndexTypeDefault:
		postList, err = ConsolePostIndex(ctx, limit, offset, false)
		if err != nil {
			log.WithTrace(ctx).Error(err)
			return
		}
		postCount, err = ConsolePostCount(ctx, false)
		if err != nil {
			log.WithTrace(ctx).Error(err)
			return
		}
	default:
		postList, err = ConsolePostIndex(ctx, limit, offset, false)
		if err != nil {
			log.WithTrace(ctx).Error(err)
			return
		}

		postCount, err = ConsolePostCount(ctx, false)
		if err != nil {
			log.WithTrace(ctx).Error(err)
			return
		}
	}

	paginate := common.MyPaginate(postCount, limit, queryPageInt)

	res = common.IndexPostList{
		PostListArr: postList,
		Paginate:    paginate,
	}

	jsonRes, err := json.Marshal(&res)
	if err != nil {
		log.WithTrace(ctx).Error(err)
		return
	}
	err = conf.CacheClient.HSet(cacheKey, field, jsonRes).Err()
	if err != nil {
		log.WithTrace(ctx).Error(err)
		return
	}
	return
}

func IndexPostDetail(ctx context.Context, postIdStr string) (postDetail common.IndexPostDetail, err error) {
	cacheKey := conf.Cnf.PostDetailIndexKey
	field := ":post:id:" + postIdStr

	postIdInt, err := strconv.Atoi(postIdStr)
	if err != nil {
		log.WithTrace(ctx).Error(err)
		return
	}
	cacheRes, err := conf.CacheClient.HGet(cacheKey, field).Result()
	if err == redis.Nil {
		// cache key does not exist
		// set data to the cache what use the cache key
		postDetail, err := doCacheIndexPostDetail(ctx, cacheKey, field, postIdInt)
		if err != nil {
			log.WithTrace(ctx).Error(err)
			return postDetail, err
		}
		return postDetail, nil
	} else if err != nil {
		log.WithTrace(ctx).Error(err)
		return postDetail, err
	}

	if cacheRes == "" {
		// Data is  null
		// set data to the cache what use the cache key
		postDetail, err = doCacheIndexPostDetail(ctx, cacheKey, field, postIdInt)
		if err != nil {
			log.WithTrace(ctx).Error(err)
			return postDetail, err
		}
		return postDetail, nil
	}

	err = json.Unmarshal([]byte(cacheRes), &postDetail)
	if err != nil {
		log.WithTrace(ctx).Error(err)
		postDetail, err = doCacheIndexPostDetail(ctx, cacheKey, field, postIdInt)
		if err != nil {
			log.WithTrace(ctx).Error(err)
			return postDetail, err
		}
		return postDetail, nil
	}
	return

}

func doCacheIndexPostDetail(ctx context.Context, cacheKey string, field string, postIdInt int) (postDetail common.IndexPostDetail, err error) {
	postDetail, err = IndexPostDetailDao(ctx, postIdInt)
	if err != nil {
		log.WithTrace(ctx).Error(err)
		return
	}
	jsonRes, err := json.Marshal(&postDetail)
	if err != nil {
		log.WithTrace(ctx).Error(err)
		return
	}
	err = conf.CacheClient.HSet(cacheKey, field, jsonRes).Err()
	if err != nil {
		log.WithTrace(ctx).Error(err)
		return
	}
	return
}

func PostViewAdd(ctx context.Context, postIdStr string) {
	postIdInt, err := strconv.Atoi(postIdStr)
	if err != nil {
		log.WithTrace(ctx).Error(err)
		return
	}

	err = conf.SqlServer.WithContext(ctx).Table((&model.PostViews{}).TableName()).
		Where("post_id = ?", postIdInt).
		UpdateColumn("num", gorm.Expr("num + ?", 1)).
		Error
	if err != nil {
		log.WithTrace(ctx).Error(err)
	}
}

func PostArchives(ctx context.Context) (archivesList map[string][]*model.Posts, err error) {
	cacheKey := conf.Cnf.ArchivesKey
	field := ":all:"

	cacheRes, err := conf.CacheClient.HGet(cacheKey, field).Result()
	if err == redis.Nil {
		// cache key does not exist
		// set data to the cache what use the cache key
		archivesList, err := doCacheArchives(ctx, cacheKey, field)
		if err != nil {
			log.WithTrace(ctx).Error(err)
			return archivesList, err
		}
		return archivesList, nil
	} else if err != nil {
		log.WithTrace(ctx).Error(err)
		return archivesList, err
	}

	if cacheRes == "" {
		// Data is  null
		// set data to the cache what use the cache key
		archivesList, err := doCacheArchives(ctx, cacheKey, field)
		if err != nil {
			log.WithTrace(ctx).Error(err)
			return archivesList, err
		}
		return archivesList, nil
	}

	archivesList = make(map[string][]*model.Posts)
	err = json.Unmarshal([]byte(cacheRes), &archivesList)
	if err != nil {
		log.WithTrace(ctx).Error(err)
		archivesList, err := doCacheArchives(ctx, cacheKey, field)
		if err != nil {
			log.WithTrace(ctx).Error(err)
			return archivesList, err
		}
		return archivesList, nil
	}
	return
}

func doCacheArchives(ctx context.Context, cacheKey string, field string) (archivesList map[string][]*model.Posts, err error) {
	posts := make([]*model.Posts, 0)
	err = conf.SqlServer.WithContext(ctx).Where("deleted_at IS NULL OR deleted_at = ?", "0001-01-01 00:00:00").Order("created_at desc").Find(&posts).Error
	if err != nil {
		log.WithTrace(ctx).Error(err)
		return
	}
	archivesList = make(map[string][]*model.Posts)
	for _, v := range posts {
		date := v.CreatedAt.Format("2006-01")
		archivesList[date] = append(archivesList[date], v)
	}

	jsonRes, err := json.Marshal(&archivesList)
	if err != nil {
		log.WithTrace(ctx).Error(err)
		return
	}
	err = conf.CacheClient.HSet(cacheKey, field, jsonRes).Err()
	if err != nil {
		log.WithTrace(ctx).Error(err)
		return
	}
	return
}

func PostRss(ctx context.Context) (rssList []*common.IndexRss, err error) {
	cacheKey := conf.Cnf.ArchivesKey
	field := ":rss:"

	cacheRes, err := conf.CacheClient.HGet(cacheKey, field).Result()
	if err == redis.Nil {
		// cache key does not exist
		// set data to the cache what use the cache key
		rssList, err := doCacheRss(ctx, cacheKey, field)
		if err != nil {
			log.WithTrace(ctx).Error(err)
			return rssList, err
		}
		return rssList, nil
	} else if err != nil {
		log.WithTrace(ctx).Error(err)
		return rssList, err
	}

	if cacheRes == "" {
		// Data is  null
		// set data to the cache what use the cache key
		rssList, err := doCacheRss(ctx, cacheKey, field)
		if err != nil {
			log.WithTrace(ctx).Error(err)
			return rssList, err
		}
		return rssList, nil
	}

	err = json.Unmarshal([]byte(cacheRes), &rssList)
	if err != nil {
		log.WithTrace(ctx).Error(err)
		rssList, err := doCacheRss(ctx, cacheKey, field)
		if err != nil {
			log.WithTrace(ctx).Error(err)
			return rssList, err
		}
		return rssList, nil
	}
	return
}

func CommonRss(ctx context.Context) (feed *feeds.Feed, err error) {
	res, err := PostRss(ctx)
	if err != nil {
		log.WithTrace(ctx).Error(err)
		return
	}

	system, err := IndexSystem(ctx)
	if err != nil {
		log.WithTrace(ctx).Error(err)
		return
	}

	feed = &feeds.Feed{
		Title:       system.Title,
		Link:        &feeds.Link{Href: conf.Cnf.AppUrl},
		Description: system.Description,
		Author:      &feeds.Author{Name: conf.Cnf.Author, Email: conf.Cnf.Email},
	}

	for _, v := range res {
		idStr := strconv.Itoa(v.Id)
		feedItem := &feeds.Item{
			Title:       v.Title,
			Link:        &feeds.Link{Href: conf.Cnf.AppUrl + "/detail/" + idStr},
			Author:      &feeds.Author{Name: conf.Cnf.Author, Email: conf.Cnf.Email},
			Description: v.Summary,
			Id:          idStr,
			Updated:     v.UpdatedAt,
			Created:     v.CreatedAt,
		}
		feed.Items = append(feed.Items, feedItem)
	}
	return feed, nil
}

func doCacheRss(ctx context.Context, cacheKey string, field string) (rssList []*common.IndexRss, err error) {
	posts := make([]*model.Posts, 0)
	err = conf.SqlServer.WithContext(ctx).Where("deleted_at IS NULL OR deleted_at = ?", "0001-01-01 00:00:00").Order("created_at desc").Find(&posts).Error
	if err != nil {
		log.WithTrace(ctx).Error(err)
		return
	}
	for _, v := range posts {
		user := new(model.Users)
		user, err = GetUserById(ctx, v.UserID)
		if err != nil {
			log.WithTrace(ctx).Error(err)
			return
		}
		rl := &common.IndexRss{
			Id:        v.ID,
			Uid:       v.UID,
			Author:    user.Name,
			Title:     v.Title,
			Summary:   v.Summary,
			CreatedAt: v.CreatedAt,
			UpdatedAt: v.UpdatedAt,
		}
		rssList = append(rssList, rl)
	}

	jsonRes, err := json.Marshal(&rssList)
	if err != nil {
		log.WithTrace(ctx).Error(err)
		return
	}
	err = conf.CacheClient.HSet(cacheKey, field, jsonRes).Err()
	if err != nil {
		log.WithTrace(ctx).Error(err)
		return
	}
	return
}
