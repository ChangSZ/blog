package conf

import "time"

const (
	DBHOST     = "127.0.0.1"
	DBPORT     = "3306"
	DBPASSWORD = "123456"
	DBUSERNAME = "root"
	DBDATABASE = "blog"

	MAIlTYPE = "html"

	HASHIDSALT      = "salt"
	HASHIDMINLENGTH = 8

	REDISADDR = ""
	REDISPWD  = ""
	REDISDB   = 0

	JWTISS       = "csz"
	JWTAUDIENCE  = "csz"
	JWTJTI       = "csz"
	JWTSECRETKEY = "csz"
	JWTTOKENKEY  = "login:token:"
	JWTTOKENLIFE = time.Hour * time.Duration(72)
)

const (
	BackUpDest        = "./backup"
	BackUpDuration    = "0 0 0 * * *"
	BackUpSqlFileName = "-sql-backup.sql"
	BackUpFilePath    = "./backup/"
)
