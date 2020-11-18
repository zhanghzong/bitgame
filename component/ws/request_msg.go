package ws

import (
	"encoding/json"
	"github.com/zhanghuizong/bitgame/app/constants/errConst"
	"github.com/zhanghuizong/bitgame/app/definition"
	"github.com/zhanghuizong/bitgame/app/models"
	"github.com/zhanghuizong/bitgame/service/config"
	"github.com/zhanghuizong/bitgame/utils"
	"github.com/zhanghuizong/bitgame/utils/jwt"
	"time"
)

func parseMsg(c *Client, message []byte) {
	// ping/pong
	if string(message) == "ping" {
		c.send <- []byte("pong")
		return
	}

	isAuth := utils.IsAuth()
	if isAuth && c.commonKey == "" {
		closeClient(c)
		c.Warnln("客户未进行认证, common-key 为空")
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
			secret := config.GetJwtKey()
			jwtRes, isJwt := jwt.Decode(reqJwtStr, secret)
			if isJwt != nil {
				c.Errorln("解密JWT失败", isJwt, "token:"+reqJwtStr, "secret:"+secret)
				c.insidePushError(errConst.TokenExpired)
				return
			}

			// jwt 格式异常
			jwtStr, jwtErr := json.Marshal(jwtRes)
			if jwtErr != nil {
				c.Errorln("解析 jwt 执行 json.Marshal 异常", jwtErr)
				c.insidePushError(errConst.BadJwtToken)
				return
			}

			jwtData := definition.ParamJwt{}
			jsonErr := json.Unmarshal(jwtStr, &jwtData)
			if jsonErr != nil {
				c.Errorln("解析 jwt 执行 json.Unmarshal 异常", jsonErr)
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
				c.Jwt.Data.Uid = uid
				c.Uid = uid
				ManagerHub.BindSocketId(c.Uid, c.SocketId)
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

	// call
	value, ok := getHandlers(cmd)
	if ok == false {
		c.insidePushError(errConst.NoCmd)
		return
	}

	c.Entry = c.WithField("pid", time.Now().UnixNano())
	c.Entry = c.WithField("uid", c.Uid)
	c.Msg = requestMsg

	value(c)
}

// 单点登录
func singleLogin(c *Client) {
	uid := c.Uid
	if uid == "" {
		return
	}

	// 判断是否异地登录
	// -1:未登录
	// 0:已登录(通知其它服务器判断)
	// 1:已登录(本服务器已操作)
	isLogin := alreadyLogin(c)
	if isLogin == 0 {
		alreadyLoginNotify(uid) // 通知其余服务器
	}

	model := new(models.LoginModel)
	model.AddSocketId(uid, c.SocketId)

	// 绑定 uid与socketId
	ManagerHub.BindSocketId(c.Uid, c.SocketId)
}

func closeClient(c *Client) {
	if c == nil {
		return
	}

	c.conn.Close()
	ManagerHub.unregister <- c
	c = nil
}

func offline(c *Client) {
	value, ok := getHandlers("offline")
	if ok {
		value(c)
	}

	model := new(models.LoginModel)
	uid := c.Uid
	connSocketId := model.GetSocketId(uid)
	if connSocketId == c.SocketId {
		model.DelSocketId(uid)
	}
}
