package javaApi

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"github.com/zhanghuizong/bitgame/service/config"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func SendMessage(msg string) (*TelegramSendMessageStruct, error) {
	host := config.GetTelegramUrl()
	api := "/telegram/send_message"
	url := host + api

	// 构造请求数据格式
	jsonRes, jErr := json.Marshal(map[string]interface{}{
		"chat_id": config.GetTelegramChatId(),
		"message": msg,
	})
	if jErr != nil {
		return nil, jErr
	}

	dataStr := string(jsonRes)
	logrus.Infof("请求数据:%s, url:%s", dataStr, url)

	client := &http.Client{}
	client.Timeout = time.Second * 10
	req, rErr := http.NewRequest("POST", url, strings.NewReader(dataStr))
	if rErr != nil {
		return nil, rErr
	}

	// 设置请求报文
	req.Header.Set("Content-Type", "application/json")

	resp, clientErr := client.Do(req)
	if clientErr != nil {
		logrus.Warnln("POST 请求接口异常", clientErr, url)
		return nil, clientErr
	}

	defer resp.Body.Close()

	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		logrus.Warnln("POST 读取响应数据异常", readErr, url)
		return nil, readErr
	}

	bodyStr := string(body)
	logrus.Infof("请求响应:%s, url:%s", bodyStr, url)

	responseData := new(TelegramSendMessageStruct)
	err := json.Unmarshal(body, responseData)

	return responseData, err
}
