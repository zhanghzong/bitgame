package config

import (
	"fmt"
	"github.com/zhanghuizong/bitgame/component/apollo"
)

// 获取环境变量
// dev= 开发环境
// test= 测试环境
// pre= 预发环境
// prod= 生产环境
// app.env=dev
func GetAppEnv() string {
	return apollo.Config.GetValue("app.env")
}

// 获取服务协议
// websocket 启动协议
// http
// https
// 空，则默认启动两种
// app.protocol=http
func GetAppProtocol() string {
	return apollo.Config.GetValue("app.protocol")
}

// 是否启用加密通信
// app.auth=true
func GetAppAuth() bool {
	return apollo.Config.GetBoolValue("app.auth", true)
}

// RSA 私钥
// app.rsa.private=-----BEGIN RSA PRIVATE KEY-----\n....\n-----END RSA PRIVATE KEY-----
func GetAppRsaPrivate() string {
	return apollo.Config.GetValue("app.rsa.private")
}

// RSA 公钥
// app.rsa.public=-----BEGIN PUBLIC KEY-----\n...\n-----END PUBLIC KEY-----
func GetAppRsaPublic() string {
	return apollo.Config.GetValue("app.rsa.public")
}

// jwt 秘钥
// jwt.key=8a3d4b8a3f13bc8c013f13bc8c9c0000
func GetJwtKey() string {
	return apollo.Config.GetValue("jwt.key")
}

// mysql 数据源
// mysql.host=localhost
// mysql.user=root
// mysql.passwd=123
// mysql.database=game_go_fishing
// mysql.port=3306
// mysql.charset=utf8
func GetMysqlDsn() string {
	host := apollo.Config.GetValue("mysql.host")
	user := apollo.Config.GetValue("mysql.user")
	passwd := apollo.Config.GetValue("mysql.passwd")
	database := apollo.Config.GetValue("mysql.database")
	port := apollo.Config.GetValue("mysql.port")
	charset := apollo.Config.GetValue("mysql.charset")

	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local", user, passwd, host, port, database, charset)
}

/** Redis 配置节点 **/

// host:port address
// redis.addr=localhost:6379
func GetRedisAddr() string {
	return apollo.Config.GetValue("redis.addr")
}

// Optional password. Must match the password specified in the
// requirepass server configuration option (if connecting to a Redis 5.0 instance, or lower),
// or the User Password when connecting to a Redis 6.0 instance, or greater, that is using the Redis ACL system.
// redis.password=
func GetRedisPassword() string {
	return apollo.Config.GetValue("redis.password")
}

// Database to be selected after connecting to the server.
// redis.db=1
func GetRedisDb() int {
	return apollo.Config.GetIntValue("redis.db", 1)
}

// Maximum number of socket connections.
// Default is 10 connections per every CPU as reported by runtime.NumCPU.
// redis.poolSize=30
func GetRedisPoolSize() int {
	return apollo.Config.GetIntValue("redis.poolSize", 10)
}

// Minimum number of idle connections which is useful when establishing
// new connection is slow.
// redis.minIdleConns=30
func GetRedisMinIdleConns() int {
	return apollo.Config.GetIntValue("redis.minIdleConns", 30)
}

/** java 配置节点 */

// 游戏ID
//java.gameId=10008
func GetJavaGameId() string {
	return apollo.Config.GetValue("java.gameId")
}

// 客户端ID
//java.clientId=game-fishing
func GetJavaClientId() string {
	return apollo.Config.GetValue("java.clientId")
}

// 客户端授权秘钥
// java.clientSecret=f857f55b86f04b78824ad3a94948a584
func GetJavaClientSecret() string {
	return apollo.Config.GetValue("java.clientSecret")
}

// 接口请求密钥
// java.apiKey=009093eb938e4f0e97579132d29e235d
func GetJavaApiKey() string {
	return apollo.Config.GetValue("java.apiKey")
}

// 接口域名地址
// java.serverApi=http://api.btgame.club
func GetJavaServerApi() string {
	return apollo.Config.GetValue("java.serverApi")
}

// 渠道ID
// java.channelId=BITGAME
func GetJavaChannelId() string {
	return apollo.Config.GetValue("java.channelId")
}

// 游戏一级域名
// java.domainKey=btgame.club
func GetJavaDomainKey() string {
	return apollo.Config.GetValue("java.domainKey")
}

// RSA
//java.rsa.public=-----BEGIN PUBLIC KEY-----\n....\n-----END PUBLIC KEY-----
func GetJavaRsaPublic() string {
	return apollo.Config.GetValue("java.rsa.public")
}
