package ws

var WsManager *ClientManager

func init() {
	WsManager = NewHub()
	go WsManager.Run()
}
