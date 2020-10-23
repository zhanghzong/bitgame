package bitgame

// 加载 appConfig 文件
import _ "github.com/zhanghuizong/bitgame/component/env"

// 加载日志
import _ "github.com/zhanghuizong/bitgame/component/logs"

// 启动 apollo 服务
import _ "github.com/zhanghuizong/bitgame/component/apollo"

// 启动 MySQL 服务
import _ "github.com/zhanghuizong/bitgame/component/db"

import (
	"github.com/zhanghuizong/bitgame/component/redis"
	"github.com/zhanghuizong/bitgame/component/ws"
)

func init() {
	// 启动 redis 订阅模式
	go redis.Subscribe(ws.WsManager)
}
