package redis

import (
	"encoding/json"
	v7 "github.com/go-redis/redis/v7"
	"github.com/spf13/viper"
	"github.com/zhanghuizong/bitgame/app/constants/redisConst"
	"github.com/zhanghuizong/bitgame/app/interfaces"
	"github.com/zhanghuizong/bitgame/app/structs"
	"log"
)

var (
	Redis *v7.Client
)

// Redis 非关系型数据初始化
func init() {
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

// 消息订阅
func Subscribe(clientManger interfaces.ClientManagerInterface) {
	defer func() {
		err := recover()
		if err != nil {
			log.Println("Redis 消息订阅异常", err)
		}
	}()

	pubSub := Redis.Subscribe(redisConst.ChannelName)
	msg, err := pubSub.Receive()
	if err != nil {
		log.Println("redis 订阅失败", err, msg)
		return
	}

	defer pubSub.Close()

	log.Println("redis 订阅通道 ", msg)

	// 用管道来接收消息
	ch := pubSub.Channel()

	// 处理消息
	for msg := range ch {
		log.Println("Redis 消息订阅", msg.String())
		channelMsg := new(structs.RedisChannel)
		err := json.Unmarshal([]byte(msg.Payload), channelMsg)
		if err != nil {
			log.Println("解析 Redis channel 消息异常", err)
			continue
		}

		// 消息分发
		go clientManger.RedisDispatch(channelMsg)
	}
}

// 发布消息
func Publish(message interface{}) {
	res, _ := json.Marshal(message)
	cmd := Redis.Publish(redisConst.ChannelName, res)
	_, err := cmd.Result()
	if err != nil {
		log.Println("Redis publish 异常", err)
	}
}
