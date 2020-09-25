package utils

import (
	"encoding/json"
	"game/http/structs"
	"game/utils/aes"
	"github.com/astaxie/beego"
	"log"
)

func IsAuth() bool {
	auth, err := beego.AppConfig.Bool("auth")
	if err != nil {
		log.Println("获取 auth 变量异常", err)
	}

	return auth
}

func GetRequestMsg(message []byte, commonKey string) *structs.RequestMsg {
	msg := string(message)
	msgData := msg[1:]

	// 解析请求数据
	var msgStruct []byte
	if IsAuth() {
		msgStruct = aes.Decode([]byte(msgData), []byte(commonKey))
	} else {
		msgStruct = []byte(msgData)
	}
	log.Println("接受到消息：", string(msgStruct))

	requestMsgData := &structs.RequestMsg{}
	errMsg := json.Unmarshal(msgStruct, requestMsgData)
	if errMsg != nil {
		log.Println("解析数据异常：", errMsg)
	}

	return requestMsgData
}
