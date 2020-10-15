package redis

import (
	"encoding/json"
	"github.com/go-redis/redis/v7"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/zhanghuizong/bitgame/app/constants/redisConst"
	"github.com/zhanghuizong/bitgame/app/interfaces"
	"github.com/zhanghuizong/bitgame/app/structs"
	"log"
)

var (
	Redis *redis.Client
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

	Redis = redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     password,
		DB:           dbIndex,
		PoolSize:     poolSize,
		MinIdleConns: minIdleConns,
	})

	_, err := Redis.Ping().Result()
	if err != nil {
		log.Printf("Redis 连接异常, err:%s, addr:%s", err, addr)
		logrus.Fatalf("Redis 连接异常, err:%s, addr:%s", err, addr)
		return
	}

	logrus.Info("Redis 连接成功. addr:%s", addr)
}

// 消息订阅
func Subscribe(clientManger interfaces.ClientManagerInterface) {
	defer func() {
		err := recover()
		if err != nil {
			logrus.Errorf("Redis 消息订阅异常. err:%s", err)
		}
	}()

	pubSub := Redis.Subscribe(redisConst.ChannelName)
	msg, err := pubSub.Receive()
	if err != nil {
		logrus.Errorf("Redis 消息订阅异常. err:%s, msg:%s", err, msg)
		return
	}

	defer pubSub.Close()

	// 用管道来接收消息
	ch := pubSub.Channel()

	// 处理消息
	for msg := range ch {
		logrus.Infof("Redis 订阅通道接收数据. msg:%s", msg)

		channelMsg := new(structs.RedisChannel)
		err := json.Unmarshal([]byte(msg.Payload), channelMsg)
		if err != nil {
			logrus.Errorf("Redis 订阅消息解析异常. err:%s", err)
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
		logrus.Errorf("Redis publish 异常. err:%s", err)
	}
}
