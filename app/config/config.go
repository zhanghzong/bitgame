package config

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"os"
	"runtime/debug"
	"strings"
)

var RedisPrefix string

// 初始配置文件
func init() {
	// 添加配置搜索路径
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Println("配置文件加载异常：", err, string(debug.Stack()))
		os.Exit(0)
	}

	gameId := viper.GetString("java.gameId")
	env := strings.ToLower(viper.GetString("app.env"))

	// 构建 Redis 前缀
	RedisPrefix = fmt.Sprintf("%s:%s", gameId, env)
	log.Println("config.init:", RedisPrefix)
}
