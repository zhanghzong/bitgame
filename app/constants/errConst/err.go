package errConst

var BadCmdCode = "1000"
var NoCmdCode = "1001"
var BadJwtTokenCode = "1002"
var AlreadyLoginCode = "1003"
var TokenExpiredCode = "1004"

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
	"code":  AlreadyLoginCode,
	"error": "ALREADY_LOGGED",
}

var TokenExpired = map[string]interface{}{
	"code":  TokenExpiredCode,
	"error": "TOKEN_EXPIRED",
}
