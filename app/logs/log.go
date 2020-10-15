package logs

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

func init() {
	// 设置日志格式
	logrus.SetFormatter(getTextFormatter())

	// 设置日志输出流
	name := getFileName()
	file, err := os.OpenFile("logs/"+name, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		logrus.SetOutput(file)
	}

	// 设置日志级别
	logrus.SetLevel(logrus.DebugLevel)

	// TODO Hooks
}

func getFileName() string {
	return fmt.Sprintf("%s.log", time.Now().Format("2006-01-02"))
}