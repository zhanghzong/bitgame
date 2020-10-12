package logs

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

var Log *logrus.Logger

func init() {
	Log = logrus.New()
	Log.SetReportCaller(true)
	Log.SetFormatter(getTextFormatter())

	name := getFileName()
	file, err := os.OpenFile("logs/"+name, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		Log.Out = file
	} else {
		Log.Out = os.Stdout
	}
}

func getFileName() string {
	return fmt.Sprintf("%s.log", time.Now().Format("2006-01-02"))
}
