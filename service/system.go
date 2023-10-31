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

func GetSystemList(ctx context.Context) (system *model.Systems, err error) {
	system = new(model.Systems)
	err = conf.SqlServer.WithContext(ctx).Find(&system).Error
	if err != nil {
		log.WithTrace(ctx).Error(err)
		return
	}
	if system.ID <= 0 {
		systemInsert := &model.Systems{
			Theme:        conf.Cnf.Theme,
			Title:        conf.Cnf.Title,
			Keywords:     conf.Cnf.Keywords,
			Description:  conf.Cnf.Description,
			RecordNumber: conf.Cnf.RecordNumber,
		}
		err = conf.SqlServer.WithContext(ctx).Create(systemInsert).Error
		if err != nil {
			log.WithTrace(ctx).Error(err)
			return
		}
		err = conf.SqlServer.WithContext(ctx).First(system).Error
		if err != nil {
			log.WithTrace(ctx).Error(err)
			return
		}
	}
	return
}

func SystemUpdate(ctx context.Context, sId int, ss common.ConsoleSystem) error {
	systemUpdate := &model.Systems{
		ID:           sId,
		Title:        ss.Title,
		Keywords:     ss.Keywords,
		Description:  ss.Description,
		RecordNumber: ss.RecordNumber,
		Theme:        ss.Theme,
	}
	err := conf.SqlServer.WithContext(ctx).Updates(systemUpdate).Error
	return err
}

func IndexSystem(ctx context.Context) (system *model.Systems, err error) {
	cacheKey := conf.Cnf.SystemIndexKey
	cacheRes, err := conf.CacheClient.Get(cacheKey).Result()
	if err == redis.Nil {
		system, err := doCacheIndexSystem(ctx, cacheKey)
		if err != nil {
			log.WithTrace(ctx).Error(err)
			return system, err
		}
		return system, nil
	} else if err != nil {
		log.WithTrace(ctx).Error(err)
		return system, err
	}

	err = json.Unmarshal([]byte(cacheRes), &system)
	if err != nil {
		log.WithTrace(ctx).Error(err)
		system, err = doCacheIndexSystem(ctx, cacheKey)
		if err != nil {
			log.WithTrace(ctx).Error(err)
			return nil, err
		}
		return system, nil
	}
	return system, nil
}

func doCacheIndexSystem(ctx context.Context, cacheKey string) (system *model.Systems, err error) {
	system = new(model.Systems)
	err = conf.SqlServer.WithContext(ctx).First(system).Error
	if err != nil {
		log.WithTrace(ctx).Error(err)
		return system, err
	}
	jsonRes, err := json.Marshal(&system)
	if err != nil {
		log.WithTrace(ctx).Error(err)
		return system, err
	}
	err = conf.CacheClient.Set(cacheKey, jsonRes, time.Duration(conf.Cnf.DataCacheTimeDuration)*time.Hour).Err()
	if err != nil {
		log.WithTrace(ctx).Error(err)
		return system, err
	}
	return system, nil
}
