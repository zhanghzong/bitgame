package structs

//{
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
type RequestJwtData struct {
	Domain     string
	Email      string
	OpenId     string
	PayPwdFlag bool
	PwdFlag    bool
	ShowName   string
	Timestamp  int
	Uid        string
	UserId     int
}

type RequestJwt struct {
	Data RequestJwtData
	Exp  int
	Iat  int
}

// 接受消息结构体
type RequestMsg struct {
	Cmd    string                 `json:"cmd"`
	Params map[string]interface{} `json:"params"`
}
