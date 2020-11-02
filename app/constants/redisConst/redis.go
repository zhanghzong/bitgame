package redisConst

import (
	"github.com/zhanghuizong/bitgame/utils"
)

// 单点登录
var SingleLogin = utils.GetRedisPrefix() + ":bitgame:single_login:ws"

// 消息通道名称
var ChannelName = utils.GetRedisPrefix() + ":bitgame:channel:communication"
