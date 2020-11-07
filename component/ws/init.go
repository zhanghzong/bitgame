package ws

var ManagerHub *ClientManager

func init() {
	ManagerHub = NewHub()
	go ManagerHub.Run()
}
