package redis

import (
	"encoding/json"
	v7 "github.com/go-redis/redis/v7"
	"github.com/spf13/viper"
	"github.com/zhanghuizong/bitgame/app/constants/redisConst"
	"github.com/zhanghuizong/bitgame/app/interfaces"
	"github.com/zhanghuizong/bitgame/app/logs"
	"github.com/zhanghuizong/bitgame/app/structs"
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

	Redis = v7.NewClient(&v7.Options{
		Addr:         addr,
		Password:     password,
		DB:           dbIndex,
		PoolSize:     poolSize,
		MinIdleConns: minIdleConns,
	})

	_, err := Redis.Ping().Result()
	if err != nil {
		logs.Log.WithFields(map[string]interface{}{"err": err, "addr": addr}).Error("Redis 连接异常")
	}

	logs.Log.WithFields(map[string]interface{}{"addr": addr, "dbIndex": dbIndex}).Info("Redis 连接成功")
}

// 消息订阅
func Subscribe(clientManger interfaces.ClientManagerInterface) {
	defer func() {
		err := recover()
		if err != nil {
			logs.Log.WithFields(map[string]interface{}{"err": err}).Error("Redis 消息订阅异常")
		}
	}()

	pubSub := Redis.Subscribe(redisConst.ChannelName)
	msg, err := pubSub.Receive()
	if err != nil {
		logs.Log.WithFields(map[string]interface{}{"err": err, "msg": msg}).Error("Redis 消息订阅异常")
		return
	}

	defer pubSub.Close()

	logs.Log.WithFields(map[string]interface{}{"msg": msg}).Info("Redis 订阅通道")

	// 用管道来接收消息
	ch := pubSub.Channel()

	// 处理消息
	for msg := range ch {
		logs.Log.WithFields(map[string]interface{}{"msg": msg.String()}).Info("Redis 订阅通道接收数据")

		channelMsg := new(structs.RedisChannel)
		err := json.Unmarshal([]byte(msg.Payload), channelMsg)
		if err != nil {
			logs.Log.WithFields(map[string]interface{}{"err": err}).Error("Redis 订阅消息解析异常")
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
		logs.Log.WithFields(map[string]interface{}{"err": err}).Error("Redis publish 异常")
	}
}
