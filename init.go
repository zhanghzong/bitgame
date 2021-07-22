package bitgame

import (
	"github.com/zhanghuizong/bitgame/component/apollo"
	"github.com/zhanghuizong/bitgame/component/db"
	"github.com/zhanghuizong/bitgame/component/env"
	"github.com/zhanghuizong/bitgame/component/logs"
	"github.com/zhanghuizong/bitgame/component/redis"
	"github.com/zhanghuizong/bitgame/component/ws"
	"github.com/zhanghuizong/bitgame/utils/console"
)

func init() {
	if console.IsConsoleVersion() {
		return
	}

	// 配置文件初始化
	env.Init()

	// Apollo 初始化
	apollo.Init()

	// 日志初始化
	logs.Init()

	// MySQL 初始化
	db.Init()

	// Redis 初始化
	redis.Init()

	// 启动 redis 订阅模式
	go redis.Subscribe(ws.Init())
}
