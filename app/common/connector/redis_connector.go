package connector

import (
	"github.com/abyss414/house/app/common/config"
	"github.com/piaohao/godis"
)

var RedisPool *godis.Pool

func InitGlobalRedisClient() {
	pool := godis.NewPool(&godis.PoolConfig{}, &godis.Option{
		Host: config.GlobalConfig().Redis.Host,
		Port: config.GlobalConfig().Redis.Port,
		Db:   0,
	})
	RedisPool = pool
}

func GetRedisClient() (*godis.Redis, error) {
	return RedisPool.GetResource()
}

func InitDependentRedisClient() *godis.Redis {
	redis := godis.NewRedis(&godis.Option{
		Host: config.GlobalConfig().Redis.Host,
		Port: config.GlobalConfig().Redis.Port,
		Db:   0,
	})
	return redis
}
