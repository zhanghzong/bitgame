package redis

import (
	"crypto/tls"
	"encoding/json"
	"github.com/go-redis/redis/v7"
	"github.com/sirupsen/logrus"
	"github.com/zhanghuizong/bitgame/app/constants/redisConst"
	"github.com/zhanghuizong/bitgame/app/definition"
	"github.com/zhanghuizong/bitgame/service/config"
)

var (
	Redis *redis.Client
)

// Redis 非关系型数据初始化
func init() {
	addr := config.GetRedisAddr()
	password := config.GetRedisPassword()
	dbIndex := config.GetRedisDb()
	poolSize := config.GetRedisPoolSize()
	minIdleConns := config.GetRedisMinIdleConns()
	insecureSkipVerify := config.GetInsecureSkipVerify()

	if addr == "" {
		addr = "localhost:6379"
	}

	if dbIndex <= 0 {
		dbIndex = 1
	}

	// Redis 基础配置
	options := &redis.Options{
		Addr:         addr,
		Password:     password,
		DB:           dbIndex,
		PoolSize:     poolSize,
		MinIdleConns: minIdleConns,
	}

	// 启用不安全的TLS认证
	if insecureSkipVerify {
		options.TLSConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	}

	Redis = redis.NewClient(options)

	_, err := Redis.Ping().Result()
	if err != nil {
		logrus.Fatalln("Redis 连接异常", err, addr)
		return
	}

	logrus.Infof("Redis 连接成功. addr:%s", addr)
}

// Subscribe 消息订阅
func Subscribe(clientManger definition.ClientManagerInterface) {
	defer func() {
		err := recover()
		if err != nil {
			logrus.Errorln("Redis 消息订阅异常", err)
		}
	}()

	pubSub := Redis.Subscribe(redisConst.ChannelName)
	msg, err := pubSub.Receive()
	if err != nil {
		logrus.Errorln("Redis 消息订阅异常", err, msg)
		return
	}

	defer pubSub.Close()

	// 用管道来接收消息
	ch := pubSub.Channel()

	// 处理消息
	for msg := range ch {
		channelMsg := new(definition.RedisChannel)
		err := json.Unmarshal([]byte(msg.Payload), channelMsg)
		if err != nil {
			logrus.Errorln("Redis 订阅消息解析异常", err)
			continue
		}

		// 消息分发
		go clientManger.RedisDispatch(channelMsg)
	}
}

// Publish 发布消息
func Publish(message interface{}) {
	res, _ := json.Marshal(message)
	cmd := Redis.Publish(redisConst.ChannelName, res)
	_, err := cmd.Result()
	if err != nil {
		logrus.Errorln("Redis publish 异常", err)
	}
}
