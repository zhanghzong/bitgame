package bitgame

type ClientManager struct {
	// Registered clients.
	clients map[string]*Client

	// Register requests from the clients.
	Register chan *Client

	// Unregister requests from clients.
	Unregister chan *Client

	// 用户与客户ID绑定关系
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
			h.clients[client.Id] = client

		// 客户端退出
		case client := <-h.Unregister:
			if _, ok := h.clients[client.Id]; ok {
				delete(h.clients, client.Id)
				delete(h.UserList, client.Uid)
				close(client.Send)
			}
		}
	}
}

func (h *ClientManager) GetClient(uid string) *Client {
	id, _ := h.UserList[uid]

	return h.clients[id]
}
