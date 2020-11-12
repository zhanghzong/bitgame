package db

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/sirupsen/logrus"
	"github.com/zhanghuizong/bitgame/service/config"
	"time"
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
		txt = fmt.Sprintf("SQL:%s; [%s] %s", executeSql, executeTime, executePath)
		logrus.Println(txt)
	} else {
		logrus.Error(values)
	}
}

// MySQL 数据初始化
func init() {
	dsn := config.GetMysqlDsn()
	poolSize := config.GetMysqlPoolSize()

	var err error
	Db, err = gorm.Open("mysql", dsn)
	if err != nil {
		logrus.Fatalln("MySQL 连接异常", err, dsn)
		return
	}

	// 打印 SQL 语句
	Db.SetLogger(new(dbLog))
	Db.LogMode(true)
	Db.DB().SetConnMaxLifetime(time.Minute) // 设置连接过期时间

	if poolSize > 0 {
		Db.DB().SetMaxIdleConns(poolSize) // 设置闲置的连接数
		Db.DB().SetMaxOpenConns(poolSize) // 设置最大打开的连接数
	}

	logrus.Infof("MySQL 连接成功. dsn:%s", dsn)
}
