package apollo

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/zouyx/agollo/v4"
	"github.com/zouyx/agollo/v4/env/config"
	"log"
)

// Config 阿波罗实例对象
var Config *agollo.Client

func Init() {
	log.Println("apollo.init")

	apolloConfig := &config.AppConfig{
		AppID:          viper.GetString("apollo.appId"),
		Cluster:        viper.GetString("apollo.cluster"),
		IP:             viper.GetString("apollo.ip"),
		NamespaceName:  viper.GetString("apollo.namespaceName"),
		IsBackupConfig: false,
		Secret:         viper.GetString("apollo.secret"),
	}

	var err error
	Config, err = agollo.StartWithConfig(func() (*config.AppConfig, error) { return apolloConfig, nil })
	if err != nil {
		logrus.Fatalln("阿波罗连接异常", err)
	}

	Config.AddChangeListener(&changeListener{})
}
