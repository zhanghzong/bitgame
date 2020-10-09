package interfaces

import "github.com/zhanghuizong/bitgame/app/structs"

type ClientManagerInterface interface {
	RedisDispatch(msg *structs.RedisChannel)
}
