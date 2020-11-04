package env

import (
	"github.com/fsnotify/fsnotify"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"runtime/debug"
)

// 初始配置文件
func init() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.WatchConfig()
	err := viper.ReadInConfig()
	if err != nil {
		logrus.Fatalln("配置文件加载异常", err, string(debug.Stack()))
	}

	viper.OnConfigChange(func(e fsnotify.Event) {
		logrus.Infof("Config file changed:%s", e.Name)
	})
}
