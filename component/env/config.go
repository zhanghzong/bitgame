package env

import (
	"github.com/spf13/viper"
	"log"
	"os"
	"runtime/debug"
)

// 初始配置文件
func init() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Println("配置文件加载异常：", err, string(debug.Stack()))
		os.Exit(0)
	}

	log.Println("应用配置文件加载成功")
}
