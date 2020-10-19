package ws

import (
	"github.com/zhanghuizong/bitgame/app/definition"
	"sync"
)

type ClientManager struct {
	// 客户端
	// socketId<=>client
	clients map[string]*Client

	// 监听客户注册请求
	register chan *Client

	// 监听客户端退出
	unregister chan *Client

	// 用户与客户ID绑定关系
	// userID<=>socketId
	userList map[string]string

	sync.RWMutex
}

func NewHub() *ClientManager {
	return &ClientManager{
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[string]*Client),
		userList:   make(map[string]string),
	}
}

func (h *ClientManager) Run() {
	defer h.Unlock()

	for {
		select {
		// 客户登录
		case client := <-h.register:
			h.Lock()
			h.clients[client.SocketId] = client
			h.Unlock()

		// 客户端退出
		case client := <-h.unregister:
			h.Lock()
			if _, ok := h.clients[client.SocketId]; ok {
				delete(h.clients, client.SocketId)
				delete(h.userList, client.Uid)
				close(client.send)
			}
			h.Unlock()
		}
	}
}

func (h *ClientManager) GetClientByUserId(uid string) *Client {
	h.RLock()
	defer h.RUnlock()
	id, _ := h.userList[uid]

	return h.clients[id]
}

func (h *ClientManager) GetClientBySocketId(socketId string) *Client {
	h.RLock()
	defer h.RUnlock()
	client, _ := h.clients[socketId]

	return client
}

// Redis channel 消息分发
func (h *ClientManager) RedisDispatch(msg *definition.RedisChannel) {
	users := msg.Users
	for _, uid := range users {
		client := h.GetClientByUserId(uid)
		if client == nil {
			continue
		}

		client.sendMsg(msg.Data)
	}
}
