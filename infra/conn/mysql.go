package conn

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/ChangSZ/golib/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/ChangSZ/golib/file"

	"github.com/ChangSZ/blog/infra/conf"

	_ "github.com/go-sql-driver/mysql"
)

var mysqlDB *gorm.DB
var sqlParam *SqlParam

type SqlParam struct {
	Host     string
	Port     string
	DataBase string
	UserName string
	Password string
	LogPath  string
}

type Sp func(*SqlParam) interface{}

func (p *Sp) SetDbHost(host string) Sp {
	return func(p *SqlParam) interface{} {
		h := p.Host
		p.Host = host
		return h
	}
}

func (p *Sp) SetDbPort(port string) Sp {
	return func(p *SqlParam) interface{} {
		pt := p.Port
		p.Port = port
		return pt
	}
}

func (p *Sp) SetDbDataBase(dataBase string) Sp {
	return func(p *SqlParam) interface{} {
		db := p.DataBase
		p.DataBase = dataBase
		return db
	}
}

func (p *Sp) SetDbPassword(pwd string) Sp {
	return func(p *SqlParam) interface{} {
		password := p.Password
		p.Password = pwd
		return password
	}
}

func (p *Sp) SetDbUserName(u string) Sp {
	return func(p *SqlParam) interface{} {
		name := p.UserName
		p.UserName = u
		return name
	}
}

func InitMysql(options ...Sp) (*gorm.DB, error) {
	q := &SqlParam{
		Host:     conf.DBHOST,
		Port:     conf.DBPORT,
		Password: conf.DBPASSWORD,
		DataBase: conf.DBDATABASE,
		UserName: conf.DBUSERNAME,
	}
	for _, option := range options {
		option(q)
	}

	dsn := ""
	var db *gorm.DB
	var err error

	dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", q.UserName, q.Password, q.Host, q.Port, q.DataBase)
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: setLogger()})
	if err != nil {
		panic(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}

	sqlDB.SetMaxIdleConns(3)
	sqlDB.SetMaxOpenConns(20)
	sqlDB.SetConnMaxLifetime(25)
	sqlDB.SetConnMaxIdleTime(10)

	mysqlDB = db
	sqlParam = q
	return mysqlDB, err
}

func setLogger() logger.Interface {
	return log.NewSQLLogger(logger.Config{
		SlowThreshold:             time.Second, // 慢 SQL 阈值
		LogLevel:                  logger.Info, // 日志级别
		IgnoreRecordNotFoundError: true,        // 忽略ErrRecordNotFound（记录未找到）错误
		Colorful:                  false,       // 禁用彩色打印
	})
}

func MySQLDump(fileName string, filePath string) error {
	log.Infof("sql dump file, name: %v, path: %v", fileName, filePath)

	err := file.MkdirAll(filePath)
	if err != nil {
		log.Errorf("mkdir -p has error: %v", err)
	}

	dumpFile := filePath + fileName
	if file.FileOrDirExists(dumpFile) {
		err = os.Remove(dumpFile)
		if err != nil {
			log.Errorf("sql dump has error: %v", err)
		}
	}

	// 使用 mysqldump 命令来导出数据和结构
	cmd := exec.Command("mysqldump",
		"-h", sqlParam.Host,
		"-P", sqlParam.Port,
		fmt.Sprintf("-u%s", sqlParam.UserName),
		fmt.Sprintf("-p%s", sqlParam.Password),
		sqlParam.DataBase)
	output, err := cmd.Output()
	if err != nil {
		log.Errorf("mysqldump error: %v", err)
		return err
	}

	// 将导出的内容写入文件
	f, err := os.Create(dumpFile)
	if err != nil {
		log.Errorf("create mysqldump file error: %v", err)
		return err
	}
	defer f.Close()

	_, err = f.Write(output)
	if err != nil {
		log.Errorf("write mysqldump file error: %v", err)
		return err
	}

	return nil
}
