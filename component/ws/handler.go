package ws

import (
	"github.com/zhanghuizong/bitgame/app/definition"
	"sync"
)

type functions func(client *Client, message *definition.RequestMsg)

var (
	handlers        = make(map[string]functions)
	handlersRWMutex sync.RWMutex
)

// 注册
func Register(key string, value functions) {
	handlersRWMutex.Lock()
	defer handlersRWMutex.Unlock()
	handlers[key] = value

	return
}

// 获取
func getHandlers(key string) (value functions, ok bool) {
	handlersRWMutex.RLock()
	defer handlersRWMutex.RUnlock()

	value, ok = handlers[key]

	return
}
