package service

import (
	"context"

	"github.com/ChangSZ/golib/log"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"

	"github.com/ChangSZ/blog/conf"
	"github.com/ChangSZ/blog/infra/alarm"
)

// 自定义返回值结构体
type MyPutRet struct {
	Key    string
	Hash   string
	Fsize  int
	Bucket string
	Name   string
}

// Upload file to Qiniu
// LocalFile is the local file, such as "./static/images/uploads/2.jpeg"
// FileName is the name what  qiniu name as
// The storage Zone is default
func Qiniu(ctx context.Context, localFile string, fileName string) {
	accessKey := conf.Cnf.QiNiuAccessKey
	secretKey := conf.Cnf.QiNiuSecretKey
	//localFile := "./static/images/uploads/2.jpeg"
	bucket := conf.Cnf.QiNiuBucket
	key := fileName
	putPolicy := storage.PutPolicy{
		Scope: bucket,
	}
	mac := qbox.NewMac(accessKey, secretKey)
	upToken := putPolicy.UploadToken(mac)

	cfg := storage.Config{}
	switch conf.Cnf.QiNiuZone {
	case "HUABEI":
		cfg.Zone = &storage.ZoneHuabei
	case "HUADONG":
		cfg.Zone = &storage.ZoneHuadong
	case "BEIMEI":
		cfg.Zone = &storage.ZoneBeimei
	case "HUANAN":
		cfg.Zone = &storage.ZoneHuanan
	case "XINJIAPO":
		cfg.Zone = &storage.ZoneXinjiapo
	default:
		cfg.Zone = &storage.ZoneHuabei
	}
	cfg.UseHTTPS = false
	cfg.UseCdnDomains = false

	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}

	putExtra := storage.PutExtra{
		Params: map[string]string{
			//"x:name": "github logo",
		},
	}
	err := formUploader.PutFile(context.Background(), &ret, upToken, key, localFile, &putExtra)
	if err != nil {
		log.WithTrace(ctx).Errorf("文件上传七牛云失败, 文件名是：%v", fileName)
		alarm.Alarm("文件上传七牛云失败, 文件名是" + fileName)
		return
	}

	log.WithTrace(ctx).Infof("文件上传七牛云成功, 文件名是：%v", fileName)
	return
}
