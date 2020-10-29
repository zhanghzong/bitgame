package ws

import (
	"encoding/json"
	"github.com/rs/xid"
	"github.com/zhanghuizong/bitgame/app/constants/envConst"
	"github.com/zhanghuizong/bitgame/app/constants/errConst"
	"github.com/zhanghuizong/bitgame/app/definition"
	"github.com/zhanghuizong/bitgame/app/models"
	"github.com/zhanghuizong/bitgame/service/config"
	"github.com/zhanghuizong/bitgame/utils"
	"github.com/zhanghuizong/bitgame/utils/jwt"
	"time"
)

func parseMsg(c *Client, message []byte) {
	isAuth := utils.IsAuth()
	if isAuth && c.commonKey == "" {
		closeClient(c)
		c.Warnf("客户未进行认证, common-key 为空")
		return
	}

	// 解析消息体
	requestMsg := utils.GetRequestMsg(message, c.commonKey)
	if requestMsg == nil {
		c.insidePushError(errConst.BadCmd)
		return
	}

	// 没有 CMD
	cmd := requestMsg.Cmd
	if cmd == "" {
		c.insidePushError(errConst.NoCmd)
		return
	}

	// 首次解密 JWT
	params := requestMsg.Params
	if c.Jwt.Data.Uid == "" {
		// 解密 JWT
		reqJwtStr, isOk := params["jwt"].(string)
		if isAuth && isOk && reqJwtStr != "" {
			jwtRes := jwt.Decode(reqJwtStr, config.GetJwtKey())

			// jwt 格式异常
			jwtStr, jwtErr := json.Marshal(jwtRes)
			if jwtErr != nil {
				c.insidePushError(errConst.BadJwtToken)
				return
			}

			jwtData := definition.ParamJwt{}
			jsonErr := json.Unmarshal(jwtStr, &jwtData)
			if jsonErr != nil {
				c.insidePushError(errConst.BadJwtToken)
				return
			}

			// 记录 JWT 数据
			c.Jwt = jwtData

			// 首次连接
			if c.Uid == "" {
				c.Uid = c.Jwt.Data.Uid

				// websocket hook 上线操作
				value, ok := getHandlers("online")
				if ok {
					value(c)
				}

				// 异地登录检测
				singleLogin(c)
			}
		} else {
			uid, isOk := params["uid"].(string)
			if isOk {
				// 绑定 uid与socketId
				env := config.GetAppEnv()
				if env == envConst.Local || env == envConst.Dev || env == envConst.Test {
					c.Jwt.Data.Uid = uid
					c.Uid = uid
					c.Hub.userList[c.Uid] = c.SocketId
				}
			}
		}
	}

	// 删除 必要 消息体字段
	delete(requestMsg.Params, "jwt")
	delete(requestMsg.Params, "userId")

	msgJson, _ := json.Marshal(map[string]interface{}{
		"message": requestMsg,
		"jwt":     c.Jwt,
	})

	c.Infof("接收消息:%s", msgJson)

	value, ok := getHandlers(cmd)
	if ok == false {
		c.insidePushError(errConst.NoCmd)
		return
	}

	_, isOk := c.Data["uid"]
	if isOk {
		c.Data["rid"] = xid.New().String()
	}

	c.Data["uid"] = c.Uid

	c.Msg = requestMsg

	value(c)
}

// 单点登录
func singleLogin(c *Client) {
	uid := c.Uid

	model := new(models.LoginModel)
	oldSocketId := model.GetSocketId(uid)

	if oldSocketId != "" {
		oldClient := c.Hub.GetClientBySocketId(oldSocketId)
		if oldClient != nil {
			oldClient.insidePushError(errConst.AlreadyLogin)

			time.AfterFunc(time.Second*3, func() {
				closeOtherClient(oldClient)
			})
		}
	}

	model.AddSocketId(uid, c.SocketId)

	// 绑定 uid与socketId
	c.Hub.userList[c.Uid] = c.SocketId
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
