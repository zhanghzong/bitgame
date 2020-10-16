package utils

import (
	"encoding/json"
	"fmt"
	"github.com/zhanghuizong/bitgame/app/definition"
	"github.com/zhanghuizong/bitgame/service/config"
	"github.com/zhanghuizong/bitgame/utils/aes"
	"log"
	"strings"
)

func GetRedisPrefix() string {
	gameId := config.GetJavaGameId()
	env := strings.ToLower(config.GetAppEnv())

	// 构建 Redis 前缀
	return fmt.Sprintf("%s:%s", gameId, env)
}

func IsAuth() bool {
	return config.GetAppAuth()
}

func GetRequestMsg(message []byte, commonKey string) *definition.RequestMsg {
	if IsAuth() {
		msgData := message[1:]
		var err error
		message, err = aes.Decode(msgData, []byte(commonKey))
		if err != nil {
			log.Println("消息体解密异常", err)
			return nil
		}
	}

	requestMsgData := &definition.RequestMsg{}
	errMsg := json.Unmarshal(message, requestMsgData)
	if errMsg != nil {
		log.Println("解析数据异常：", errMsg)
		return nil
	}

	return requestMsgData
}
