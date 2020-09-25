package utils

import (
	"encoding/json"
	"github.com/spf13/viper"
	"github.com/zhanghuizong/bitgame/structs"
	"github.com/zhanghuizong/bitgame/utils/aes"
	"log"
)

func IsAuth() bool {
	return viper.GetBool("app.auth")
}

func GetRequestMsg(message []byte, commonKey string) *structs.RequestMsg {
	if IsAuth() {
		msgData := message[1:]
		var err error
		message, err = aes.Decode(msgData, []byte(commonKey))
		if err != nil {
			log.Println("消息体解密异常", err)
			return nil
		}
	}

	requestMsgData := &structs.RequestMsg{}
	errMsg := json.Unmarshal(message, requestMsgData)
	if errMsg != nil {
		log.Println("解析数据异常：", errMsg)
		return nil
	}

	return requestMsgData
}
