package ws

import (
	"github.com/zhanghuizong/bitgame/app/constants/errConst"
	"github.com/zhanghuizong/bitgame/app/definition"
	"github.com/zhanghuizong/bitgame/app/models"
	"github.com/zhanghuizong/bitgame/component/redis"
	"os"
	"time"
)

func pushClient(c *Client, data interface{}) {
	uid := c.Uid
	if uid == "" {
		c.sendMsg(data)
		return
	}

	single(uid, data)
}

// 消息单推
func single(uid string, data interface{}) {
	if uid == "" {
		return
	}

	client := ManagerHub.GetClientByUserId(uid)
	if client == nil {
		redisPublish(uid, data)
		return
	}

	client.sendMsg(data)
}

func redisPublish(uid string, data interface{}) {
	hostname, _ := os.Hostname()
	channelMsg := new(definition.RedisChannel)
	channelMsg.Type = "response"
	channelMsg.Hostname = hostname
	channelMsg.Users = []string{uid}
	channelMsg.Data = data
	redis.Publish(channelMsg)
}

// 消息广播
func broadcast(users []string, cmd string, data interface{}) {
	for _, uid := range users {
		single(uid, pushSuccess(cmd, data))
	}
}

// 正常消息格式
func pushSuccess(cmd string, data interface{}) map[string]interface{} {
	return map[string]interface{}{
		"cmd":  cmd,
		"code": 200,
		"data": data,
	}
}

// 错误消息格式
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

// 系统内部错误消息格式
func insideDataDesc(res interface{}) interface{} {
	return map[string]interface{}{
		"cmd": "conn::error",
		"res": res,
	}
}

// 系统错误消息推送
func insidePushError(c *Client, res map[string]interface{}) {
	data := insideDataDesc(res)
	c.sendMsg(data)
}

// 异地登录
func alreadyLogin(c *Client) int {
	if c == nil {
		return -1
	}

	uid := c.Uid
	model := new(models.LoginModel)
	oldSocketId := model.GetSocketId(uid)
	if oldSocketId == "" {
		return -1
	}

	// 本服务器搜索到
	oldClient := ManagerHub.GetClientBySocketId(oldSocketId)
	if oldClient != nil {
		c.Warnln("触发异地登录，即将关闭 socket", oldSocketId)

		// 本服务器推送
		insidePushError(c, errConst.AlreadyLogin)

		// 关闭客户端
		time.AfterFunc(time.Second*3, func() {
			closeClient(oldClient)
		})

		return 1
	}

	return 0
}

// 通知其余服务器检测异地登录
func alreadyLoginNotify(uid string) {
	hostname, _ := os.Hostname()
	channelMsg := new(definition.RedisChannel)
	channelMsg.Type = "alreadyLogin"
	channelMsg.Hostname = hostname
	channelMsg.Users = []string{uid}
	redis.Publish(channelMsg)
}
