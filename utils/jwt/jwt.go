package jwt

import (
	"github.com/dgrijalva/jwt-go"
	"log"
	"time"
)

// 构造 JWT Token
func Encode(data interface{}, secret string) string {
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"data": data,
		"iat":  time.Now().Unix(),                         // 签发时间
		"exp":  time.Now().Add(time.Hour * 24 * 2).Unix(), // 过期时间
	})
	token, err := at.SignedString([]byte(secret))
	if err != nil {
		return ""
	}

	return token
}

// 解密 JWT Token
func Decode(token string, secret string) map[string]interface{} {
	defer func() {
		err := recover()
		if err != nil {
			log.Printf("解密JWT异常：%s\n", err)
		}
	}()

	claim, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		log.Println("解密 JWT Token 异常：" + err.Error())
		return map[string]interface{}{}
	}

	return claim.Claims.(jwt.MapClaims)
}
