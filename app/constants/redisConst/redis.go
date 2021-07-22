package redisConst

import (
	"github.com/zhanghuizong/bitgame/utils"
	"log"
)

// SingleLogin 单点登录
var SingleLogin = ""

// ChannelName 消息通道名称
var ChannelName = ""

func Init() {
	log.Println("redis.constants.init")

	SingleLogin = utils.GetRedisPrefix() + ":bitgame:single_login:ws"

	ChannelName = utils.GetRedisPrefix() + ":bitgame:channel:communication"
}
