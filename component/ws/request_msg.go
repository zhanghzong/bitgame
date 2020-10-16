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

	// (预发|生产)禁止使用未加密数据传输
	if !isAuth {
		env := config.GetAppEnv()
		if env == envConst.Pre || env == envConst.Prod {
			c.insidePushError(errConst.BadJwtToken)
			return
		}
	}

	if isAuth && c.commonKey == "" {
		closeClient(c)
		c.Log.Warnf("客户未进行认证, common-key 为空")
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
	if c.ParamJwt.Data.Uid == "" {
		var jwtStr []byte
		jwtData := definition.ParamJwt{}

		// 解密 JWT
		if isAuth {
			var jwtErr error
			jwtRes := jwt.Decode(params["jwt"].(string), config.GetJwtKey())
			jwtStr, jwtErr = json.Marshal(jwtRes)
			if jwtErr != nil {
				c.insidePushError(errConst.BadJwtToken)
				return
			}
		} else {
			requestJwt, ok := params["jwt"].(map[string]interface{})

			// 不存在 jwt 参数
			if !ok {
				uid, uidOk := params["uid"].(string)
				if uidOk {
					jwtData.Data.Uid = uid
					jwtData.Data.UserId = int(time.Now().Unix())
					jwtData.Data.ShowName = uid
				}
			}
			jwtStr, _ = json.Marshal(requestJwt)
		}

		jsonErr := json.Unmarshal(jwtStr, &jwtData)
		if jsonErr != nil {
			c.insidePushError(errConst.BadJwtToken)
			return
		}

		// 记录 JWT 数据
		c.ParamJwt = jwtData

		// 参数异常
		if c.ParamJwt.Data.Uid == "" {
			c.insidePushError(errConst.BadJwtToken)
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

	msgJson, _ := json.Marshal(map[string]interface{}{
		"message": requestMsg,
		"jwt":     c.ParamJwt,
	})

	c.Log.Infof("接收消息:%s", msgJson)

	value, ok := getHandlers(cmd)
	if ok == false {
		c.insidePushError(errConst.NoCmd)
		return
	}

	_, isOk := c.Log.Data["uid"]
	if isOk {
		c.Log.Data["rid"] = xid.New().String() // 请求ID
	}

	c.Log.Data["uid"] = c.Uid

	value(c, requestMsg)
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
