package logs

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"log"
	"os"
	"time"
)

func init() {
	// 设置日志格式
	logrus.SetFormatter(getTextFormatter())

	// 设置日志输出流
	logPath := viper.GetString("log.path")
	if logPath == "" {
		logPath = "logs/"
	}

	fullName := logPath + string(os.PathSeparator) + getFileName()
	file, err := os.OpenFile(fullName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Printf("打开日志文件异常. err:%s", err)
		return
	}

	// 设置日志级别
	logLevel := logrus.Level(viper.GetInt("log.level"))
	logrus.SetLevel(logLevel)
	logrus.SetOutput(file)
}

func getFileName() string {
	return fmt.Sprintf("%s.log", time.Now().Format("2006-01-02"))
}
