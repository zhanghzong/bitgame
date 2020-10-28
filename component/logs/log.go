package logs

import (
	"fmt"
	"github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"log"
	"os"
	"time"
)

func init() {
	// 设置日志格式
	logrus.SetFormatter(getTextFormatter())

	// 日志存放路径
	logPath := viper.GetString("log.path")
	if logPath == "" {
		logPath = "logs"
	}

	// 日志分割器
	fullName := logPath + string(os.PathSeparator) + "%Y-%m-%d.log"
	out, err := rotatelogs.New(
		fullName, // 日志文件名称

		// WithRotationTime 设置日志分割的时间，每天分割一次
		rotatelogs.WithRotationTime(time.Hour*24),

		// WithMaxAge 设置文件清理前的最长保存时间
		rotatelogs.WithMaxAge(-1),
	)

	if err != nil {
		log.Println("日志模块启动失败", err, fullName)
		return
	}

	// 设置日志级别
	logLevel := logrus.Level(viper.GetInt("log.level"))
	logrus.SetLevel(logLevel)
	logrus.SetOutput(out)
}

func getFileName() string {
	return fmt.Sprintf("%s.log", time.Now().Format("2006-01-02"))
}
