package ws

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/rs/xid"
	"github.com/sirupsen/logrus"
	"github.com/zhanghuizong/bitgame/utils"
	"github.com/zhanghuizong/bitgame/utils/rsa"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 10240,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// websocket 协议升级处理逻辑
func ServeWs(context *gin.Context) {
	conn, err := upgrader.Upgrade(context.Writer, context.Request, nil)
	if err != nil {
		logrus.Errorf("Websocket 协议升级异常. err:%s", err)

		if conn == nil {
			return
		}

		cErr := conn.Close()
		if cErr != nil {
			logrus.Errorf("关闭协议异常. err:%s", cErr)
		}

		return
	}

	// 是否启用密文传输
	var commonKey string
	if utils.IsAuth() {
		auth := context.Query("auth")
		commonKey, _ = rsa.Authorize(auth)
		log.Println("commonKey:", commonKey)

		if commonKey == "" {
			conn.WriteJSON(map[string]interface{}{
				"cmd":  "authorize",
				"code": 1,
				"data": "认证失败",
			})
			conn.Close()
			return
		}
	}

	// 生成客户唯一ID
	guid := xid.New().String()

	// 实例化客户端连接对象
	client := &Client{}
	client.SocketId = guid
	client.commonKey = commonKey
	client.Hub = WsManager
	client.conn = conn
	client.send = make(chan []byte, 1024)
	client.Entry = logrus.WithFields(map[string]interface{}{"rid": guid})

	// 注册客户端
	client.Hub.register <- client

	go client.write()
	go client.read()
}
