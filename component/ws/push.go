package ws

import (
	"github.com/zhanghuizong/bitgame/app/definition"
	"github.com/zhanghuizong/bitgame/component/redis"
	"os"
)

// 消息单推
func single(uid string, cmd string, data interface{}) {
	client := WsManager.GetClientByUserId(uid)
	if client == nil {
		hostname, _ := os.Hostname()
		channelMsg := new(definition.RedisChannel)
		channelMsg.Type = "response"
		channelMsg.Hostname = hostname
		channelMsg.Users = []string{uid}
		channelMsg.Data = pushSuccess(cmd, data)
		redis.Publish(channelMsg)

		return
	}

	client.Success(cmd, data)
}

// 消息广播
func broadcast(users []string, cmd string, data interface{}) {
	for _, uid := range users {
		single(uid, cmd, data)
	}
}

func pushSuccess(cmd string, data interface{}) map[string]interface{} {
	return map[string]interface{}{
		"cmd":  cmd,
		"code": 200,
		"data": data,
	}
}

func pushError(cmd string, row definition.ErrMsgStruct) map[string]interface{} {
	return map[string]interface{}{
		"cmd":  "error",
		"code": 200,
		"data": map[string]interface{}{
			"from":    cmd,
			"errCode": row.Code,
			"msg":     row.Msg,
		},
	}
}

// 系统错误消息推送
func insidePushError(c *Client, res map[string]interface{}) {
	data := map[string]interface{}{
		"cmd": "conn::error",
		"res": res,
	}

	c.sendMsg(data)
}
