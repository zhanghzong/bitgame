package redis

import (
	v7 "github.com/go-redis/redis/v7"
	"github.com/spf13/viper"
	"log"
)

var (
	Redis *v7.Client
)

func Init() {
	addr := viper.GetString("redis.addr")
	password := viper.GetString("redis.password")
	dbIndex := viper.GetInt("redis.db")
	poolSize := viper.GetInt("redis.poolSize")
	minIdleConns := viper.GetInt("redis.minIdleConns")

	if addr == "" {
		addr = "localhost:6379"
	}

	if dbIndex <= 0 {
		dbIndex = 1
	}

	log.Println("Redis 初始化：", addr, dbIndex)
	Redis = v7.NewClient(&v7.Options{
		Addr:         addr,
		Password:     password,
		DB:           dbIndex,
		PoolSize:     poolSize,
		MinIdleConns: minIdleConns,
	})

	_, err := Redis.Ping().Result()
	if err != nil {
		log.Println("Redis 连接异常：", err)
	}
}
