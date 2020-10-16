package rsa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"github.com/zhanghuizong/bitgame/service/config"
	"github.com/zhanghuizong/bitgame/utils/base64"
	"strings"
)

// 解密
func Decode(originData string, privateKey string) (string, error) {
	originData = base64.Decode(originData)

	// 将密钥解析成私钥实例
	block, _ := pem.Decode([]byte(privateKey))
	if block == nil {
		return "", errors.New("private key error")
	}

	// 解析pem.Decode（）返回的Block指针实例
	pri, pErr := x509.ParsePKCS1PrivateKey(block.Bytes)
	if pErr != nil {
		return "", pErr
	}

	res, dErr := rsa.DecryptPKCS1v15(rand.Reader, pri, []byte(originData))
	if dErr != nil {
		return "", dErr
	}

	return string(res), nil
}

// Rsa 认证
func Authorize(originData string) (string, error) {
	key := config.GetAppRsaPrivate()
	originData = strings.Join(strings.Split(originData, " "), "+")

	return Decode(originData, key)
}
