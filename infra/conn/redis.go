package conn

import (
	"github.com/ChangSZ/golib/log"
	"github.com/go-redis/redis"

	"github.com/ChangSZ/blog/infra/conf"
)

type RedisClient struct {
	Addr     string
	Password string
	Db       int
}

func (rc *RedisClient) SetRedisAddr(addr string) func(*RedisClient) interface{} {
	return func(rc *RedisClient) interface{} {
		rcA := rc.Addr
		rc.Addr = addr
		return rcA
	}
}

func (rc *RedisClient) SetRedisPwd(pwd string) func(*RedisClient) interface{} {
	return func(rc *RedisClient) interface{} {
		rcp := rc.Password
		rc.Password = pwd
		return rcp
	}
}

func (rc *RedisClient) SetRedisDb(db int) func(*RedisClient) interface{} {
	return func(rc *RedisClient) interface{} {
		rcdb := rc.Db
		rc.Db = db
		return rcdb
	}
}

func (rc *RedisClient) RedisInit(options ...func(*RedisClient) interface{}) (*redis.Client, error) {
	redisC := &RedisClient{
		Addr:     conf.REDISADDR,
		Password: conf.REDISPWD,
		Db:       conf.REDISDB,
	}
	for _, option := range options {
		option(redisC)
	}

	client := redis.NewClient(&redis.Options{
		Addr:     redisC.Addr,
		Password: redisC.Password, // no password set
		DB:       redisC.Db,       // use default DB
	})

	_, err := client.Ping().Result()

	if err != nil {
		log.Errorf("redis client ping is error: %v", err)
		return nil, err
	}
	//RedisC1 = client
	return client, nil
}
