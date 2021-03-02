package javaApi

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/wenzhenxi/gorsa"
	"github.com/zhanghuizong/bitgame/service/config"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// rsa 加密
func encode(originData string) string {
	publicKey := config.GetJavaRsaPublic()
	res, err := gorsa.PublicEncrypt(originData, publicKey)
	if err != nil {
		return ""
	}

	return res
}

func post(api string, data map[string]interface{}) (string, error) {
	host := config.GetJavaServerApi()
	url := host + api

	client := &http.Client{}
	client.Timeout = time.Second * 10

	jsonRes, jErr := json.Marshal(data)
	if jErr != nil {
		return "", jErr
	}

	dataStr := string(jsonRes)
	sign := encode(dataStr)
	reqData := map[string]interface{}{
		"sign": sign,
	}
	logrus.Infof("请求数据:%s, url:%s", dataStr, url)

	bodyJson, mErr := json.Marshal(reqData)
	if mErr != nil {
		go SendMessageFormat(url, dataStr, mErr)
		return "", mErr
	}

	req, rErr := http.NewRequest("POST", url, strings.NewReader(string(bodyJson)))
	if rErr != nil {
		go SendMessageFormat(url, dataStr, rErr)
		return "", rErr
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apiKey", config.GetJavaApiKey())

	resp, clientErr := client.Do(req)
	if clientErr != nil {
		logrus.Warnln("POST 请求接口异常", clientErr, url)
		go SendMessageFormat(url, dataStr, clientErr)
		return "", clientErr
	}

	defer resp.Body.Close()

	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		logrus.Warnln("POST 读取响应数据异常", readErr, url)
		go SendMessageFormat(url, dataStr, readErr)
		return "", readErr
	}

	bodyStr := string(body)
	logrus.Infof("请求响应:%s, url:%s", bodyStr, url)

	return bodyStr, nil
}

// 异常报警
func SendMessageFormat(url string, dataStr string, err error) {
	errStr := ""
	if err != nil {
		errStr = err.Error()
	}

	txt := "游戏编号：%s\n" +
		"日期：%s\n" +
		"接口地址：%s\n" +
		"请求参数：%s\n" +
		"异常：接口请求异常"

	// 拼接错误数据
	if errStr != "" {
		txt = txt + ". err:" + errStr
	}

	gameNo := config.GetJavaGameId()
	now := time.Now().Format("2006-01-02 15:04:05")
	txt = fmt.Sprintf(txt, gameNo, now, url, dataStr)
	_, _ = SendMessage(txt)
}
