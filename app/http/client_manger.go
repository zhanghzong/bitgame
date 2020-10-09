package http

import "github.com/zhanghuizong/bitgame/app/structs"

type ClientManager struct {
	// 客户端
	// socketId<=>client
	clients map[string]*Client

	// 监听客户注册请求
	Register chan *Client

	// 监听客户端退出
	Unregister chan *Client

	// 用户与客户ID绑定关系
	// userID<=>socketId
	UserList map[string]string
}

func NewHub() *ClientManager {
	return &ClientManager{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		clients:    make(map[string]*Client),
		UserList:   make(map[string]string),
	}
}

func (h *ClientManager) Run() {
	for {
		select {
		// 客户登录
		case client := <-h.Register:
			h.clients[client.SocketId] = client

		// 客户端退出
		case client := <-h.Unregister:
			if _, ok := h.clients[client.SocketId]; ok {
				delete(h.clients, client.SocketId)
				delete(h.UserList, client.Uid)
				close(client.send)
			}
		}
	}
}

func (h *ClientManager) GetClientByUserId(uid string) *Client {
	id, _ := h.UserList[uid]

	return h.clients[id]
}

func (h *ClientManager) GetClientBySocketId(socketId string) *Client {
	client, _ := h.clients[socketId]

	return client
}

// Redis channel 消息分发
func (h *ClientManager) RedisDispatch(msg *structs.RedisChannel) {
	users := msg.Users
	for _, uid := range users {
		client := h.GetClientByUserId(uid)
		if client == nil {
			continue
		}

		client.sendMsg(msg.Data)
	}
}
