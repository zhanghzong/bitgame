package javaApi

import (
	"encoding/json"
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
		return "", mErr
	}

	req, rErr := http.NewRequest("POST", url, strings.NewReader(string(bodyJson)))
	if rErr != nil {
		return "", rErr
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apiKey", config.GetJavaApiKey())

	resp, clientErr := client.Do(req)
	if clientErr != nil {
		logrus.Warnln("POST 请求接口异常", clientErr, url)
		return "", clientErr
	}

	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		logrus.Warnln("POST 读取响应数据异常", readErr, url)
		return "", readErr
	}

	bodyStr := string(body)
	logrus.Infof("请求响应:%s, url:%s", bodyStr, url)

	return bodyStr, nil
}
