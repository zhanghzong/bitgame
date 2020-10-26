package ws

import (
	"sync"
)

type functions func(client *Client)

var (
	handlers = make(map[string]functions)
	rwLock   sync.RWMutex
)

// 注册
func Register(key string, value functions) {
	rwLock.Lock()
	defer rwLock.Unlock()
	handlers[key] = value

	return
}

// 获取
func getHandlers(key string) (value functions, ok bool) {
	rwLock.RLock()
	defer rwLock.RUnlock()

	value, ok = handlers[key]

	return
}
