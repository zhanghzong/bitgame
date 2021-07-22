package ws

var ManagerHub *ClientManager

func Init() *ClientManager {
	ManagerHub = NewHub()

	go ManagerHub.Run()

	return ManagerHub
}
