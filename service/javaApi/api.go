package javaApi

import (
	"encoding/json"
	"github.com/zhanghuizong/bitgame/service/config"
	"time"
)

/**
 * 用户账户扣减
 * @param openId string, 用户唯一标识
 * @param outOrderNo  string, 游戏订单编号
 * @param currency  string，币种信息
 * @param amount float32,订单金额
 * @return orderNo 返回平台订单
 */
func DeductUserAccount(openId string, currency string, outOrderNo string, amount float64) (*DeductUserAccountStruct, error) {
	url := "/game/acct/deduct"

	data := map[string]interface{}{
		"gameNo":      config.GetJavaGameId(),
		"channelId":   config.GetJavaChannelId(),
		"requestTime": time.Now().Unix(),
		"openId":      openId,
		"currency":    currency,
		"outOrderNo":  outOrderNo,
		"amount":      amount,
	}

	res, pErr := post(url, data)
	if pErr != nil {
		return nil, pErr
	}

	responseData := new(DeductUserAccountStruct)
	err := json.Unmarshal([]byte(res), responseData)

	// 接口返回异常
	if responseData.RspCode != "0000" {
		jsonRes, _ := json.Marshal(data)
		go SendMessageFormat(url, string(jsonRes), err)
	}

	return responseData, err
}

/**
 * 用户账户增加
 * @param openId string 用户id
 * @param currency  string 货币名称
 * @param outOrderNo  string 游戏订单编号
 * @param orderNo  string 平台订单编号
 * @param amount float32 增加金额
 */
func AcctIncrease(openId string, currency string, outOrderNo string, orderNo string, amount float64) (*AcctIncreaseStruct, error) {
	url := "/game/acct/increase"

	data := map[string]interface{}{
		"gameNo":      config.GetJavaGameId(),
		"channelId":   config.GetJavaChannelId(),
		"requestTime": time.Now().Unix(),
		"openId":      openId,
		"currency":    currency,
		"outOrderNo":  outOrderNo,
		"orderNo":     orderNo,
		"amount":      amount,
	}

	res, pErr := post(url, data)
	if pErr != nil {
		return nil, pErr
	}

	responseData := new(AcctIncreaseStruct)
	err := json.Unmarshal([]byte(res), responseData)

	// 接口返回异常
	if responseData.RspCode != "0000" {
		jsonRes, _ := json.Marshal(data)
		go SendMessageFormat(url, string(jsonRes), err)
	}

	return responseData, err
}

// 查询货币配置列表
func GetCurrencyList() (*GetCurrencyListStruct, error) {
	url := "/game/currency/getCurrencyList"

	data := map[string]interface{}{
		"gameNo":      config.GetJavaGameId(),
		"channelId":   config.GetJavaChannelId(),
		"domainKey":   config.GetJavaDomainKey(),
		"requestTime": time.Now().Unix(),
	}

	res, pErr := post(url, data)
	if pErr != nil {
		return nil, pErr
	}

	responseData := new(GetCurrencyListStruct)
	err := json.Unmarshal([]byte(res), responseData)

	// 接口返回异常
	if responseData.RspCode != "0000" {
		jsonRes, _ := json.Marshal(data)
		go SendMessageFormat(url, string(jsonRes), err)
	}

	return responseData, err
}

// 批量查询用户账户余额
func GetBalanceList(openId string) (*GetBalanceListStruct, error) {
	url := "/game/acct/getBalanceList"

	data := map[string]interface{}{
		"gameNo":      config.GetJavaGameId(),
		"requestTime": time.Now().Unix(),
		"openId":      openId,
	}

	res, pErr := post(url, data)
	if pErr != nil {
		return nil, pErr
	}

	responseData := new(GetBalanceListStruct)
	err := json.Unmarshal([]byte(res), responseData)

	// 接口返回异常
	if responseData.RspCode != "0000" {
		jsonRes, _ := json.Marshal(data)
		go SendMessageFormat(url, string(jsonRes), err)
	}

	return responseData, err
}
