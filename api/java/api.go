package java

import (
	"github.com/spf13/viper"
	"time"
)

/**
 * 用户账户扣减
 * @param openId string, 用户唯一标识
 * @param outOrderNo  string, 游戏订单编号
 * @param amount float32,订单金额
 * @param currency  string，币种信息
 * @return orderNo 返回平台订单
 */
func DeductUserAccount(openId string, outOrderNo string, amount float32, currency string) string {
	url := "/game/acct/deduct"

	data := map[string]interface{}{
		"gameNo":      viper.GetString("java.gameId"),
		"channelId":   viper.GetString("java.channelId"),
		"requestTime": time.Now().Unix(),
		"openId":      openId,
		"currency":    currency,
		"outOrderNo":  outOrderNo,
		"amount":      amount,
	}

	return post(url, data)
}

/**
 * 用户账户兑换（游戏充值）
 * @param {string}  openId       用户id
 * @param {string}  outOrderNo   游戏订单编号
 * @param {number}  orderAmount  扣减的金额
 * @param {string}  currency     币种
 * @param {number}  converAmount 兑换的金额
 * @param {array}   langueList   语言配置
 * @param {boolean} decrFlag     true不扣减账号 (捕鱼专用)
 */
func DeductUserAccountNew(openId string, outOrderNo string, amount float32, currency string, converAmount float32, langueList []interface{}) string {
	url := "/game/acct/convert"

	data := map[string]interface{}{
		"gameNo":       viper.GetString("java.gameId"),
		"channelId":    viper.GetString("java.channelId"),
		"requestTime":  time.Now().Unix(),
		"openId":       openId,
		"currency":     currency,
		"outOrderNo":   outOrderNo,
		"amount":       amount,
		"converAmount": converAmount,
		"langueList":   langueList,
	}

	return post(url, data)
}

// 查询货币配置列表
func GetCurrencyList() string {
	url := "/game/currency/getCurrencyList"

	data := map[string]interface{}{
		"gameNo":      viper.GetString("java.gameId"),
		"channelId":   viper.GetString("java.channelId"),
		"domainKey":   viper.GetString("java.domainKey"),
		"requestTime": time.Now().Unix(),
	}

	return post(url, data)
}
