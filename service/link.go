package service

import (
	"context"
	"encoding/json"
	"time"

	"github.com/ChangSZ/blog/common"
	"github.com/ChangSZ/blog/conf"
	"github.com/ChangSZ/blog/infra/log"
	"github.com/ChangSZ/blog/model"
	"github.com/go-redis/redis"
)

func LinkList(ctx context.Context, offset int, limit int) (links []model.Links, cnt int64, err error) {
	links = make([]model.Links, 0)
	db := conf.SqlServer.WithContext(ctx).Table((&model.Links{}).TableName()).Order("`order` asc").Offset(offset).Limit(limit)
	if err = db.Count(&cnt).Error; err != nil {
		return
	}
	err = db.Find(&links).Error
	return
}

func LinkSore(ctx context.Context, ls common.LinkStore) (err error) {
	LinkInsert := &model.Links{
		Name:  ls.Name,
		Link:  ls.Link,
		Order: ls.Order,
	}
	err = conf.SqlServer.WithContext(ctx).Create(&LinkInsert).Error
	return
}

func LinkDetail(ctx context.Context, linkId int) (link *model.Links, err error) {
	link = new(model.Links)
	err = conf.SqlServer.WithContext(ctx).Where("id=?", linkId).Find(&link).Error
	return
}

func LinkUpdate(ctx context.Context, ls common.LinkStore, linkId int) (err error) {
	linkUpdate := &model.Links{
		ID:    linkId,
		Link:  ls.Link,
		Name:  ls.Name,
		Order: ls.Order,
	}
	err = conf.SqlServer.WithContext(ctx).Updates(linkUpdate).Error
	return
}

func LinkDestroy(ctx context.Context, linkId int) (err error) {
	link := &model.Links{ID: linkId}
	err = conf.SqlServer.WithContext(ctx).Delete(&link).Error
	return
}

func LinkCnt(ctx context.Context) (cnt int64, err error) {
	link := new(model.Links)
	err = conf.SqlServer.WithContext(ctx).Table(link.TableName()).Count(&cnt).Error
	return
}

func AllLink(ctx context.Context) (links []model.Links, err error) {
	cacheKey := conf.Cnf.LinkIndexKey
	cacheRes, err := conf.CacheClient.Get(cacheKey).Result()
	if err == redis.Nil {
		links, err := doCacheLinkList(ctx, cacheKey)
		if err != nil {
			log.WithTrace(ctx).Error(err)
			return links, err
		}
		return links, nil
	} else if err != nil {
		log.WithTrace(ctx).Error(err)
		return nil, err
	}

	err = json.Unmarshal([]byte(cacheRes), &links)
	if err != nil {
		log.WithTrace(ctx).Error(err)
		links, err = doCacheLinkList(ctx, cacheKey)
		if err != nil {
			log.WithTrace(ctx).Error(err)
			return nil, err
		}
		return links, nil
	}
	return links, nil
}

func doCacheLinkList(ctx context.Context, cacheKey string) (links []model.Links, err error) {
	links = make([]model.Links, 0)
	err = conf.SqlServer.WithContext(ctx).Find(&links).Error
	if err != nil {
		log.WithTrace(ctx).Error(err)
		return links, err
	}
	if len(links) == 0 {
		log.WithTrace(ctx).Info("len(links)==0")
		return links, nil
	}

	jsonRes, err := json.Marshal(&links)
	if err != nil {
		log.WithTrace(ctx).Error(err)
		return nil, err
	}
	err = conf.CacheClient.Set(cacheKey, jsonRes, time.Duration(conf.Cnf.DataCacheTimeDuration)*time.Hour).Err()
	if err != nil {
		log.WithTrace(ctx).Error(err)
		return nil, err
	}
	return links, nil
}
