package service

import (
	"context"
	"fmt"

	"github.com/ChangSZ/golib/log"
	"golang.org/x/crypto/bcrypt"

	"github.com/ChangSZ/blog/common"
	"github.com/ChangSZ/blog/conf"
	"github.com/ChangSZ/blog/model"
)

func GetUserByEmail(ctx context.Context, email string) (user *model.Users, err error) {
	user = new(model.Users)
	err = conf.SqlServer.WithContext(ctx).Where("email = ?", email).Find(&user).Error
	return
}

func GetUserCnt(ctx context.Context) (cnt int64, err error) {
	user := new(model.Users)
	err = conf.SqlServer.WithContext(ctx).Table(user.TableName()).Count(&cnt).Error
	return
}

func UserStore(ctx context.Context, ar common.AuthRegister) (user *model.Users, err error) {
	password := []byte(ar.Password)
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		log.WithTrace(ctx).Error(err)
		return
	}
	userInsert := &model.Users{
		Name:     ar.UserName,
		Email:    ar.Email,
		Password: string(hashedPassword),
		Status:   1,
	}
	err = conf.SqlServer.WithContext(ctx).Create(userInsert).Error
	if err != nil {
		log.WithTrace(ctx).Error(err)
		return
	}
	fmt.Println(userInsert.ID)
	return
}

func DelAllCache() {
	conf.CacheClient.Del(
		conf.Cnf.TagListKey,
		conf.Cnf.CateListKey,
		conf.Cnf.ArchivesKey,
		conf.Cnf.LinkIndexKey,
		conf.Cnf.PostIndexKey,
		conf.Cnf.SystemIndexKey,
		conf.Cnf.TagPostIndexKey,
		conf.Cnf.CatePostIndexKey,
		conf.Cnf.PostDetailIndexKey,
	)
}
