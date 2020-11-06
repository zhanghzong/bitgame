package jwt

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
	"github.com/zhanghuizong/bitgame/service/config"
	"time"
)

// 构造 JWT Token
func Encode(data interface{}, secret string) string {
	expire := config.GetJwtExpired()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"data": data,
		"iat":  time.Now().Unix(),                                        // 签发时间
		"exp":  time.Now().Add(time.Hour * time.Duration(expire)).Unix(), // 过期时间
	})

	token, err := at.SignedString([]byte(secret))
	if err != nil {
		return ""
	}

	return token
}

// 解密 JWT Token
func Decode(token string, secret string) (map[string]interface{}, error) {
	defer func() {
		err := recover()
		if err != nil {
			logrus.Errorln("解密JWT异常", err, "token:"+token, "secret:"+secret)
		}
	}()

	claim, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return map[string]interface{}{}, err
	}

	return claim.Claims.(jwt.MapClaims), nil
}
