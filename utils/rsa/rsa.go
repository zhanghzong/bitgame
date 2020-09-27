package rsa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"github.com/spf13/viper"
	"github.com/zhanghuizong/bitgame/utils/base64"
	"strings"
)

// 加密
func Encode(originData string, publicKey string) (string, error) {
	block, _ := pem.Decode([]byte(publicKey))
	if block == nil {
		return "", errors.New("public key error")
	}

	pubInterface, pErr := x509.ParsePKIXPublicKey(block.Bytes)
	if pErr != nil {
		return "", pErr
	}

	pub := pubInterface.(*rsa.PublicKey)
	res, dErr := rsa.EncryptPKCS1v15(rand.Reader, pub, []byte(originData))
	if dErr != nil {
		return "", dErr
	}

	return base64.Encode(string(res)), nil
}

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
	key := viper.GetString("app.rsa.private")
	originData = strings.Join(strings.Split(originData, " "), "+")

	return Decode(originData, key)
}
