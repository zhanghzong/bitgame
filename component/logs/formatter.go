package logs

import (
	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/sirupsen/logrus"
	"runtime"
)

func getTextFormatter() *nested.Formatter {
	return &nested.Formatter{
		TimestampFormat: "2006-01-02 15:04:05",
		CallerFirst:     true,
		NoColors:        true,
		CustomCallerFormatter: func(frame *runtime.Frame) string {
			return " [" + frame.Function + "]"
		},
		FieldsOrder: []string{"rid", "tid", "uid"},
	}
}

func getJsonFormatter() *logrus.JSONFormatter {
	return &logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			return frame.Function, ""
		},
	}
}
