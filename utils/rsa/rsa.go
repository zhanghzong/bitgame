package rsa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"github.com/zhanghuizong/bitgame/constants"
	"strings"
)

// 解密
func Decode(ciphertext []byte) ([]byte, error) {
	// 将密钥解析成私钥实例
	block, _ := pem.Decode(constants.PrivateKey)
	if block == nil {
		return nil, errors.New("private key error")
	}

	// 解析pem.Decode（）返回的Block指针实例
	pri, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return rsa.DecryptPKCS1v15(rand.Reader, pri, ciphertext) //RSA算法解密
}

// rsa 认证
func Authorize(auth string) string {
	auth = replaceSpace(auth)
	desc, err1 := base64.StdEncoding.DecodeString(auth)
	if err1 != nil {
		return ""
	}

	commonKey, err2 := Decode(desc)
	if err2 != nil {
		return ""
	}

	return string(commonKey)
}

func replaceSpace(s string) string {
	return strings.Join(strings.Split(s, " "), "+")
}
