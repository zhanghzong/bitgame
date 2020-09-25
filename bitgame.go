package bitgame

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/rs/xid"
	"github.com/spf13/viper"
	"github.com/zhanghuizong/bitgame/utils"
	"github.com/zhanghuizong/bitgame/utils/rsa"
	"log"
	"net/http"
	"os"
	"runtime/debug"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 10240,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// websocket 协议升级处理逻辑
func ServeWs(hub *ClientManager, context *gin.Context) {
	conn, err := upgrader.Upgrade(context.Writer, context.Request, nil)
	if err != nil {
		log.Println("Websocket 协议升级异常：", err)

		if conn == nil {
			return
		}

		cErr := conn.Close()
		if cErr != nil {
			log.Println("关闭协议异常", cErr)
		}

		return
	}

	// 是否启用密文传输
	var commonKey string
	if utils.IsAuth() {
		auth := context.Query("auth")
		commonKey = rsa.Authorize(auth)

		log.Println("commonKey:", commonKey)
		if commonKey == "" {
			conn.WriteJSON(map[string]interface{}{
				"cmd":  "authorize",
				"code": 1,
				"msg":  "认证失败",
			})
			conn.Close()
			return
		}
	}

	// 生成客户唯一ID
	guid := xid.New().String()

	// 实例化客户端连接对象
	client := &Client{}
	client.Id = guid
	client.CommonKey = commonKey
	client.Hub = hub
	client.Conn = conn
	client.Send = make(chan []byte, 1024)
	client.Context = context

	client.Hub.Register <- client

	go client.write()
	go client.read()
}

// 初始配置文件
func InitConfig() {
	viper.SetConfigName("config/app")
	viper.AddConfigPath(".") // 添加搜索路径

	err := viper.ReadInConfig()
	if err != nil {
		log.Println("配置文件加载异常：", err, string(debug.Stack()))
		os.Exit(0)
	}
}
