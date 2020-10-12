package db

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/spf13/viper"
	"github.com/zhanghuizong/bitgame/app/constants/envConst"
	"github.com/zhanghuizong/bitgame/app/logs"
)

var (
	Db *gorm.DB
)

// MySQL 数据初始化
func init() {
	host := viper.GetString("mysql.host")
	user := viper.GetString("mysql.user")
	passwd := viper.GetString("mysql.passwd")
	database := viper.GetString("mysql.database")
	port := viper.GetString("mysql.port")
	charset := viper.GetString("mysql.charset")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local", user, passwd, host, port, database, charset)

	var err error
	Db, err = gorm.Open("mysql", dsn)
	if err != nil {
		logs.Log.Fatalf("MySQL 连接异常,err:%s", err)
		return
	}

	// 打印 SQL 语句
	env := viper.GetString("app.env")
	if env != envConst.Prod {
		Db.LogMode(true)
	}

	logs.Log.Infof("MySQL 连接成功. dsn:%s", dsn)
}
