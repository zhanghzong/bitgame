package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/sirupsen/logrus"
	"github.com/zhanghuizong/bitgame/app/constants/envConst"
	"github.com/zhanghuizong/bitgame/service/config"
	"log"
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
		log.Printf("MySQL 连接异常,err:%s", err)
		logrus.Fatalf("MySQL 连接异常,err:%s", err)
		return
	}

	// 打印 SQL 语句
	env := config.GetAppEnv()
	if env != envConst.Prod {
		Db.LogMode(true)
	}

	logrus.Infof("MySQL 连接成功. dsn:%s", dsn)
}
