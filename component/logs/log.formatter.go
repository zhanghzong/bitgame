package logs

import (
	nested "github.com/antonfisher/nested-logrus-formatter"
	logstash "github.com/bshuster-repo/logrus-logstash-hook"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/zhanghuizong/bitgame/utils/ip"
	"runtime"
)

func getTextFormatter() *nested.Formatter {
	return &nested.Formatter{
		TimestampFormat: "2006-01-02 15:04:05.000",
		CallerFirst:     true,
		NoColors:        true,
		CustomCallerFormatter: func(frame *runtime.Frame) string {
			return " [" + frame.Function + "]"
		},
		FieldsOrder: []string{"pid", "rid", "tid", "uid"},
	}
}

func getJsonFormatter() *logrus.JSONFormatter {
	return &logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05.000",
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			return frame.Function, ""
		},
	}
}

// 输出 logstash 日志格式
func getLogstashFormatter() logrus.Formatter {
	return logstash.DefaultFormatter(logrus.Fields{
		"appName": viper.GetString("apollo.appId"),
		"host":    ip.GetIp(),
	})
}
