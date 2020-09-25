package bitgame

import (
	"github.com/zhanghuizong/bitgame/structs"
	"sync"
)

type disposeFunc func(client *Client, message *structs.RequestMsg)

var (
	handlers        = make(map[string]disposeFunc)
	handlersRWMutex sync.RWMutex
)

// 注册
func Register(key string, value disposeFunc) {
	handlersRWMutex.Lock()
	defer handlersRWMutex.Unlock()
	handlers[key] = value

	return
}

// 获取
func getHandlers(key string) (value disposeFunc, ok bool) {
	handlersRWMutex.RLock()
	defer handlersRWMutex.RUnlock()

	value, ok = handlers[key]

	return
}
