package service

import (
	"context"

	"github.com/ChangSZ/golib/log"

	"github.com/ChangSZ/blog/conf"
	"github.com/ChangSZ/blog/model"
)

func GetUserById(ctx context.Context, userId int) (*model.Users, error) {
	user := new(model.Users)
	err := conf.SqlServer.WithContext(ctx).Select("name, email").Where("id=?", userId).Find(&user).Error
	if err != nil {
		log.WithTrace(ctx).Error(err)
		return user, err
	}
	return user, nil
}

func UserCnt(ctx context.Context) (cnt int64, err error) {
	user := new(model.Users)
	err = conf.SqlServer.WithContext(ctx).Table(user.TableName()).Count(&cnt).Error
	return
}
