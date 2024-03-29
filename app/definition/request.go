package definition

import "encoding/json"

// RequestMsg 接受消息结构体
type RequestMsg struct {
	Cmd    string                 `json:"cmd"`
	Params map[string]interface{} `json:"params"`
}

func (t RequestMsg) String() string {
	s, _ := json.Marshal(t)

	return string(s)
}

// ParamJwt {
//	"data": {
//		"domain": "btgame.club",
//		"email": "1019***@qq.com",
//		"openId": "487D96C51E87457A9A5F2B0C41085864",
//		"payPwdFlag": true,
//		"pwdFlag": true,
//		"showName": "1019***@qq.com",
//		"timestamp": 1600691912882,
//		"uid": "11",
//		"userId": 100000037
//	},
//	"exp": 1601296712,
//	"iat": 1600691912
//}
type ParamJwt struct {
	Data ParamJwtData `json:"data"`
	Exp  int          `json:"exp"`
	Iat  int          `json:"iat"`
}

type ParamJwtData struct {
	Domain     string `json:"domain"`
	Email      string `json:"email"`
	OpenId     string `json:"openId"`
	PayPwdFlag bool   `json:"payPwdFlag"`
	PwdFlag    bool   `json:"pwdFlag"`
	ShowName   string `json:"showName"`
	Timestamp  int    `json:"timestamp"`
	Uid        string `json:"uid"`
	UserId     int    `json:"userId"`
}
