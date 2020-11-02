package db

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/sirupsen/logrus"
	"github.com/zhanghuizong/bitgame/service/config"
)

var (
	Db *gorm.DB
)

type dbLog struct {
}

func (t dbLog) Print(values ...interface{}) {
	executeFlag := values[0] // 执行类型
	txt := ""
	if executeFlag == "sql" {
		valuesFormat := gorm.LogFormatter(values...)
		executePath := values[1]      // 文件路径
		executeTime := values[2]      // 执行时间
		executeSql := valuesFormat[3] // 执行SQL语句
		txt = fmt.Sprintf("%s; [%s] %s", executeSql, executeTime, executePath)
		logrus.Println(txt)
	} else {
		logrus.Error(values)
	}
}

// MySQL 数据初始化
func init() {
	dsn := config.GetMysqlDsn()

	var err error
	Db, err = gorm.Open("mysql", dsn)
	if err != nil {
		logrus.Fatalf("MySQL 连接异常,err:%s, dsn:%s", err, dsn)
		return
	}

	// 打印 SQL 语句
	Db.SetLogger(new(dbLog))
	Db.LogMode(true)

	logrus.Infof("MySQL 连接成功. dsn:%s", dsn)
}
