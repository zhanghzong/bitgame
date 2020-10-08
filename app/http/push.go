package http

import (
	"github.com/zhanghuizong/bitgame/app/middleware/redis"
	"github.com/zhanghuizong/bitgame/app/structs"
	"os"
)

// 消息单推
func Single(uid string, data interface{}) {
	client := WsManager.GetClientByUserId(uid)
	if client == nil {
		hostname, _ := os.Hostname()
		channelMsg := new(structs.RedisChannel)
		channelMsg.Type = "response"
		channelMsg.Hostname = hostname
		channelMsg.Users = []string{uid}
		channelMsg.Data = data
		redis.Publish(channelMsg)

		return
	}

	client.SendMsg(data)
}

// 消息广播
func Broadcast(users []string, data interface{}) {
	for _, uid := range users {
		Single(uid, data)
	}
}
