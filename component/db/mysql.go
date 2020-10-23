package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/sirupsen/logrus"
	"github.com/zhanghuizong/bitgame/service/config"
)

var (
	Db *gorm.DB
)

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
	Db.SetLogger(logrus.StandardLogger())
	Db.LogMode(true)

	logrus.Infof("MySQL 连接成功. dsn:%s", dsn)
}
