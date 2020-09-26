package base64

import "encoding/base64"

// base64 加密
func Encode(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}

// base64 解密
func Decode(s string) string {
	res, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return ""
	}

	return string(res)
}
