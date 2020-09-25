package bitgame

import (
	"encoding/json"
	"github.com/spf13/viper"
	"github.com/zhanghuizong/bitgame/constants"
	"github.com/zhanghuizong/bitgame/structs"
	"github.com/zhanghuizong/bitgame/utils"
	"github.com/zhanghuizong/bitgame/utils/jwt"
	"log"
)

func parseMsg(c *Client, message []byte) {
	if utils.IsAuth() && c.CommonKey == "" {
		c.Hub.Unregister <- c
		cErr := c.Conn.Close()
		if cErr != nil {
			log.Println(".....", cErr)
		}

		log.Println("客户未进行认证")
		return
	}

	// 解析消息体
	requestMsg := utils.GetRequestMsg(message, c.CommonKey)
	if requestMsg == nil {
		c.SendMsg(constants.Error["1000"])
		return
	}

	// 没有 CMD
	cmd := requestMsg.Cmd
	if cmd == "" {
		c.SendMsg(constants.Error["1001"])
		return
	}

	// 首次解密 JWT
	params := requestMsg.Params
	if c.RequestJwt.Data.Uid == "" {
		jwtRes := jwt.Decode(params["jwt"].(string), viper.GetString("jwt.key"))
		jwtStr, jwtErr := json.Marshal(jwtRes)
		if jwtErr != nil {
			c.SendMsg(constants.Error["1002"])
			return
		}

		jwtData := structs.RequestJwt{}
		jsonErr := json.Unmarshal(jwtStr, &jwtData)
		if jsonErr != nil {
			c.SendMsg(constants.Error["1002"])
			return
		}

		// 记录 JWT 数据
		c.RequestJwt = jwtData

		if c.RequestJwt.Data.Uid == "" {
			c.SendMsg(constants.Error["1002"])
			return
		}

		// 首次连接
		if c.Uid == "" {
			c.Uid = c.RequestJwt.Data.Uid
			c.Hub.UserList[c.Uid] = c.Id
		}
	}

	// 删除 必要 消息体字段
	delete(requestMsg.Params, "jwt")
	delete(requestMsg.Params, "userId")
	delete(requestMsg.Params, "_userInfo")

	msgJson, _ := json.Marshal(map[string]interface{}{
		"message": requestMsg,
		"jwt":     c.RequestJwt,
	})
	log.Printf("接受消息:%s\n", msgJson)

	value, ok := getHandlers(cmd)
	if ok == false {
		c.SendMsg(map[string]interface{}{
			"cmd": "conn::error",
			"res": map[string]interface{}{
				"code":  1001,
				"error": "NO_CMD",
			},
		})
		return
	}

	value(c, requestMsg)
}
