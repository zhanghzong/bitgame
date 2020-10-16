package javaApi

import (
	"encoding/json"
	"github.com/wenzhenxi/gorsa"
	"github.com/zhanghuizong/bitgame/service/config"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

// rsa加密
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
	client.Timeout = time.Second * 3

	jsonRes, jErr := json.Marshal(data)
	if jErr != nil {
		return "", jErr
	}

	dataStr := string(jsonRes)
	sign := encode(dataStr)
	reqData := map[string]interface{}{
		"sign": sign,
	}
	log.Println("请求接口数据", dataStr, url)

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
		log.Println("POST 请求接口异常", url, clientErr)
		return "", clientErr
	}

	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		log.Println("POST 读取响应数据异常", url, readErr)
		return "", readErr
	}

	bodyStr := string(body)
	log.Println("请求接口响应数据", url, bodyStr)

	return bodyStr, nil
}
