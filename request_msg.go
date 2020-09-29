package bitgame

import (
	"encoding/json"
	"github.com/spf13/viper"
	"github.com/zhanghuizong/bitgame/app/constants/errConst"
	"github.com/zhanghuizong/bitgame/app/models/login"
	"github.com/zhanghuizong/bitgame/app/structs"
	"github.com/zhanghuizong/bitgame/utils"
	"github.com/zhanghuizong/bitgame/utils/jwt"
	"log"
	"time"
)

func parseMsg(c *Client, message []byte) {
	if utils.IsAuth() && c.CommonKey == "" {
		closeClient(c)
		log.Println("客户未进行认证, common-key 为空")
		return
	}

	// 解析消息体
	requestMsg := utils.GetRequestMsg(message, c.CommonKey)
	if requestMsg == nil {
		c.pushError(errConst.BadCmd)
		return
	}

	// 没有 CMD
	cmd := requestMsg.Cmd
	if cmd == "" {
		c.pushError(errConst.NoCmd)
		return
	}

	// 首次解密 JWT
	params := requestMsg.Params
	if c.ParamJwt.Data.Uid == "" {
		jwtRes := jwt.Decode(params["jwt"].(string), viper.GetString("jwt.key"))
		jwtStr, jwtErr := json.Marshal(jwtRes)
		if jwtErr != nil {
			c.pushError(errConst.BadJwtToken)
			return
		}

		jwtData := structs.ParamJwt{}
		jsonErr := json.Unmarshal(jwtStr, &jwtData)
		if jsonErr != nil {
			c.pushError(errConst.BadJwtToken)
			return
		}

		// 记录 JWT 数据
		c.ParamJwt = jwtData

		// 参数异常
		if c.ParamJwt.Data.Uid == "" {
			c.pushError(errConst.BadJwtToken)
			return
		}

		// 首次连接
		if c.Uid == "" {
			c.Uid = c.ParamJwt.Data.Uid

			// websocket hook 上线操作
			value, ok := getHandlers("online")
			if ok {
				value(c, nil)
			}

			// 异地登录检测
			singleLogin(c)
		}
	}

	// 删除 必要 消息体字段
	delete(requestMsg.Params, "jwt")
	delete(requestMsg.Params, "userId")
	delete(requestMsg.Params, "_userInfo")

	msgJson, _ := json.Marshal(map[string]interface{}{
		"message": requestMsg,
		"jwt":     c.ParamJwt,
	})
	log.Printf("接受消息:%s\n", msgJson)

	value, ok := getHandlers(cmd)
	if ok == false {
		c.pushError(errConst.NoCmd)
		return
	}

	value(c, requestMsg)
}

// 单点登录
func singleLogin(c *Client) {
	uid := c.Uid

	model := new(login.Model)
	oldSocketId := model.GetSocketId(uid)

	if oldSocketId != "" {
		oldClient := WsManager.GetClientBySocketId(oldSocketId)
		if oldClient != nil {
			oldClient.pushError(errConst.AlreadyLogin)

			time.AfterFunc(time.Second*3, func() {
				closeOtherClient(oldClient)
			})
		}
	}

	model.AddSocketId(uid, c.SocketId)

	// 绑定 uid与socketId
	c.Hub.UserList[c.Uid] = c.SocketId
}

func closeOtherClient(c *Client) {
	if c == nil {
		return
	}

	// 手动触发离线事件
	method := c.conn.CloseHandler()
	if method != nil {
		method(1, "异地登录关闭")
	}

	closeClient(c)
}
