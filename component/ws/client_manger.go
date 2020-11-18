package ws

import (
	"github.com/zhanghuizong/bitgame/app/definition"
	"sync"
)

type ClientManager struct {
	// 客户端
	// socketId<=>client
	clients sync.Map

	// 监听客户注册请求
	register chan *Client

	// 监听客户端退出
	unregister chan *Client

	// 用户与客户ID绑定关系
	// userID<=>socketId
	userList sync.Map
}

func NewHub() *ClientManager {
	return &ClientManager{
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *ClientManager) Run() {
	for {
		select {
		// 客户登录
		case client := <-h.register:
			h.clients.Store(client.SocketId, client)

		// 客户端退出
		case client := <-h.unregister:
			if _, ok := h.clients.Load(client.SocketId); ok {
				h.clients.Delete(client.SocketId)
				h.userList.Delete(client.Uid)
				close(client.send)
			}
		}
	}
}

// 用户与客户ID绑定关系
// userID<=>socketId
func (h *ClientManager) BindSocketId(uid string, socketId string) {
	h.userList.Store(uid, socketId)
}

// 根据 uid 换取 socketId
func (h *ClientManager) GetClientByUserId(uid string) *Client {
	socketId, ok1 := h.userList.Load(uid)
	if !ok1 {
		return nil
	}

	client, ok2 := h.clients.Load(socketId)
	if !ok2 {
		return nil
	}

	c, ok := client.(*Client)
	if !ok {
		return nil
	}

	return c
}

func (h *ClientManager) GetClientBySocketId(socketId string) *Client {
	client, _ := h.clients.Load(socketId)

	c, ok := client.(*Client)
	if !ok {
		return nil
	}

	return c
}

// Redis channel 消息分发
func (h *ClientManager) RedisDispatch(msg *definition.RedisChannel) {
	switch msg.Type {
	// 正常消息推送
	case "response":
		users := msg.Users
		for _, uid := range users {
			client := h.GetClientByUserId(uid)
			if client == nil {
				continue
			}

			client.sendMsg(msg.Data)
		}

	// 异常登录
	case "alreadyLogin":
		users := msg.Users
		for _, uid := range users {
			// -1:未登录
			// 0:已登录(通知其它服务器判断)
			// 1:已登录(本服务器已操作)
			alreadyLogin(h.GetClientByUserId(uid))
		}
	}
}
