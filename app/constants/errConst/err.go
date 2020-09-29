package errConst

var BadCmdCode string = "1000"
var NoCmdCode string = "1001"
var BadJwtTokenCode string = "1002"
var AlreadyLoginCode string = "1003"

var BadCmd = map[string]interface{}{
	"code":  BadCmdCode,
	"error": "BAD_CMD",
}

var NoCmd = map[string]interface{}{
	"code":  NoCmdCode,
	"error": "NO_CMD",
}

var BadJwtToken = map[string]interface{}{
	"code":  BadJwtTokenCode,
	"error": "BAD_JWT_TOKEN",
}

var AlreadyLogin = map[string]interface{}{
	"res": AlreadyLoginCode,
	"cmd": "ALREADY_LOGGED",
}
