package bitgame

import (
	"github.com/spf13/viper"
	"github.com/zhanghuizong/bitgame/app/middleware/db"
	"github.com/zhanghuizong/bitgame/app/middleware/redis"
	"log"
	"os"
	"runtime/debug"
)

var WsManager *ClientManager

// 初始配置文件
func init() {
	viper.SetConfigName("config/app")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		log.Println("配置文件加载异常：", err, string(debug.Stack()))
		os.Exit(0)
	}

	WsManager = NewHub()
	go WsManager.Run()

	// 启动 MySQL 服务
	db.Init()

	// 启动 Redis 服务
	redis.Init()
}
