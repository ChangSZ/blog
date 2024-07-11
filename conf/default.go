package conf

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ChangSZ/golib/log"
	"github.com/ChangSZ/golib/mail"
	"github.com/go-redis/redis"
	"github.com/speps/go-hashids"
	"gopkg.in/yaml.v2"
	"gorm.io/gorm"

	"github.com/ChangSZ/blog/infra/alarm"
	"github.com/ChangSZ/blog/infra/backup"
	"github.com/ChangSZ/blog/infra/conn"
	"github.com/ChangSZ/blog/infra/hashid"
	"github.com/ChangSZ/blog/infra/jwt"
)

var (
	SqlServer   *gorm.DB
	ZHashId     *hashids.HashID
	CacheClient *redis.Client
	MailClient  *mail.EmailParam
	Cnf         *Conf
	Env         string
)

func DefaultInit() {
	CnfInit()
	LoggerInit()
	DbInit()
	AlarmInit()
	MailInit()
	ZHashIdInit()
	RedisInit()
	JwtInit()
	// the customer error code init
	SetMsg(Msg)
}

func CnfInit() {
	cf := &Conf{
		AppUrl:                "http://localhost:8081",
		AppImgUrl:             "http://localhost:8081/static/uploads/images/",
		DefaultLimit:          "20",
		DefaultIndexLimit:     "3",
		LogFilePath:           "./log/blog",
		LogRotateMaxDays:      7,
		DbUser:                "root",
		DbPassword:            "",
		DbPort:                "3306",
		DbDataBase:            "blog",
		DbHost:                "127.0.0.1",
		AlarmType:             "mail",
		MailUser:              "test@test.com",
		MailPwd:               "",
		MailHost:              "smtp.qq.com",
		MailPort:              587,
		HashIdSalt:            "i must add a salt what is only for me",
		HashIdLength:          8,
		JwtIss:                "blog",
		JwtAudience:           "blog",
		JwtJti:                "blog",
		JwtSecretKey:          "blog",
		JwtTokenLife:          3,
		RedisAddr:             "localhost:6379",
		RedisPwd:              "",
		RedisDb:               0,
		BackUpFilePath:        "./backup/",
		BackUpDuration:        "* * */1 * *",
		BackUpSentTo:          "792386472@qq.com",
		DataCacheTimeDuration: 720,
		ImgUploadUrl:          "http://localhost:8081/console/post/imgUpload",
		ImgUploadDst:          "./static/uploads/images/",
		MinioUploadImg:        true,
		MinioEndpoint:         "",
		MinioBucketName:       "",
		MinioAccessKey:        "",
		MinioSecretKey:        "",
		QiNiuUploadImg:        false,
		QiNiuHostName:         "",
		QiNiuAccessKey:        "",
		QiNiuSecretKey:        "",
		QiNiuBucket:           "blog",
		QiNiuZone:             "HUABEI",
		CateListKey:           "all:cate:sort",
		TagListKey:            "all:tag",
		Theme:                 0,
		Title:                 "默认Title",
		Keywords:              "默认关键词,搬运工",
		Description:           "个人网站,https://github.com/ChangSZ/blog",
		RecordNumber:          "000-0000",
		UserCnt:               2,
		Author:                "搬运工",
		Email:                 "",
		PostIndexKey:          "index:all:post:list",
		TagPostIndexKey:       "index:all:tag:post:list",
		CatePostIndexKey:      "index:all:cate:post:list",
		LinkIndexKey:          "index:all:link:list",
		SystemIndexKey:        "index:all:system:list",
		PostDetailIndexKey:    "index:post:detail",
		ArchivesKey:           "index:archives:list",
		GithubName:            "",
		GithubRepo:            "",
		GithubClientId:        "",
		GithubClientSecret:    "",
		GithubLabels:          "Gitalk",
		ThemeJs:               "/static/assets/js",
		ThemeCss:              "/static/assets/css",
		ThemeImg:              "/static/assets/img",
		ThemeFancyboxCss:      "/static/assets/fancybox",
		ThemeFancyboxJs:       "/static/assets/fancybox",
		ThemeHLightCss:        "/static/assets/highlightjs",
		ThemeHLightJs:         "/static/assets/highlightjs",
		ThemeShareCss:         "/static/assets/css",
		ThemeShareJs:          "/static/assets/js",
		ThemeArchivesJs:       "/static/assets/js",
		ThemeArchivesCss:      "/static/assets/css",
		ThemeNiceImg:          "/static/assets/img",
		ThemeAllCss:           "/static/assets/css",
		ThemeIndexImg:         "/static/assets/img",
		ThemeCateImg:          "/static/assets/img",
		ThemeTagImg:           "/static/assets/img",
	}

	files, _ := filepath.Glob("./env.*.yaml")
	dev := false
	prod := false
	for _, v := range files {
		switch v {
		case "env.dev.yaml":
			dev = true
		case "env.prod.yaml":
			prod = true
		default:
			continue
		}
	}

	var fileName string
	var env string
	if dev {
		fileName = "/env.dev.yaml"
		env = "dev"
	} else if prod {
		fileName = "/env.prod.yaml"
		env = "prod"
	} else {
		fileName = "default"
		env = "dev"
	}

	if fileName == "default" {
		Cnf = cf
		Env = env
		return
	}

	res, err := filepath.Abs(filepath.Dir("./main.go"))
	if err != nil {
		log.Error(err)
	}

	//读取yaml配置文件
	yamlFile, err := os.ReadFile(res + fileName)
	if err != nil {
		log.Error(err)
	}

	err = yaml.Unmarshal(yamlFile, &cf)
	if err != nil {
		log.Error(err)
	}

	Cnf = cf
	Env = env
}

func LoggerInit() {
	logCfg := log.Config{
		FilePath: Cnf.LogFilePath,
		MaxDays:  Cnf.LogRotateMaxDays,
		LogLevel: "INFO",
	}
	log.Init(logCfg)
}

func DbInit() {
	sp := new(conn.Sp)
	dbUser := sp.SetDbUserName(Cnf.DbUser)
	dbPwd := sp.SetDbPassword(Cnf.DbPassword)
	dbPort := sp.SetDbPort(Cnf.DbPort)
	dbHost := sp.SetDbHost(Cnf.DbHost)
	dbDb := sp.SetDbDataBase(Cnf.DbDataBase)
	sqlServer, err := conn.InitMysql(dbUser, dbPwd, dbPort, dbHost, dbDb)
	SqlServer = sqlServer
	if err != nil {
		log.Errorf("some errors: %v", err)
		panic(err.Error())
	}
}

func BackUpInit() {
	bp := new(backup.BackUpParam)
	dest := "./zip/" + time.Now().Format("2006-01-02") + ".zip"
	bp.SetFilePath(Cnf.BackUpFilePath).
		SetFiles("./backup", "./static/uploads/images").
		SetDest(dest).
		SetCronSpec(Cnf.BackUpDuration)
	data := make(map[string]string)
	data[time.Now().Format("2006-01-02")+".zip"] = dest
	bp.Ep = MailClient
	subject := time.Now().Format("2006-01-02") + "备份邮件"
	bp.Ep.SetSubject(subject).
		SetAttaches(data).
		SetBody(
			`<html><body>
			<p><img src="https://golang.org/doc/gopher/doc.png"></p><br/>
			<h1>日常备份</h1>
			</body></html>`).
		SetTo(strings.Split(Cnf.BackUpSentTo, ","))
	bp.Backup()
}

func AlarmInit() {
	a := new(alarm.AlarmParam)
	alarmT := a.SetType(alarm.AlarmType(Cnf.AlarmType))
	mailTo := a.SetMailTo("792386472@qq.com")
	err := a.AlarmInit(alarmT, mailTo)
	if err != nil {
		log.Error(err)
	}
}

func MailInit() {
	client, err := mail.Init(
		mail.WithUser(Cnf.MailUser),
		mail.WithPwd(Cnf.MailPwd),
		mail.WithHost(Cnf.MailHost),
		mail.WithPort(Cnf.MailPort))
	if err != nil {
		log.Error(err)
		return
	}
	MailClient = client
	log.Info("begin to backup")
	BackUpInit()
}

func ZHashIdInit() {
	hd := new(hashid.HashIdParams)
	salt := hd.SetHashIdSalt(Cnf.HashIdSalt)
	hdLength := hd.SetHashIdLength(Cnf.HashIdLength)
	zHashId, err := hd.HashIdInit(hdLength, salt)
	if err != nil {
		log.Error(err)
	}
	ZHashId = zHashId
}

func RedisInit() {
	rc := new(conn.RedisClient)
	addr := rc.SetRedisAddr(Cnf.RedisAddr)
	pwd := rc.SetRedisPwd(Cnf.RedisPwd)
	db := rc.SetRedisDb(Cnf.RedisDb)
	client, err := rc.RedisInit(addr, db, pwd)
	if err != nil {
		log.Error(err)
		panic(err)
	}
	CacheClient = client
}

func JwtInit() {
	jt := new(jwt.JwtParam)
	ad := jt.SetDefaultAudience(Cnf.JwtAudience)
	jti := jt.SetDefaultJti(Cnf.JwtJti)
	iss := jt.SetDefaultIss(Cnf.JwtIss)
	sk := jt.SetDefaultSecretKey(Cnf.JwtSecretKey)
	rc := jt.SetRedisCache(CacheClient)
	tl := jt.SetTokenLife(time.Hour * time.Duration(Cnf.JwtTokenLife))
	_ = jt.JwtInit(ad, jti, iss, sk, rc, tl)
}
