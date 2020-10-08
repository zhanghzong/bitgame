package bitgame

import (
	_ "github.com/zhanghuizong/bitgame/app/config"
	"github.com/zhanghuizong/bitgame/app/http"
	_ "github.com/zhanghuizong/bitgame/app/middleware/db"
	"github.com/zhanghuizong/bitgame/app/middleware/redis"
)

func init() {
	// 启动 redis 订阅模式
	go redis.Subscribe(http.WsManager)
}
