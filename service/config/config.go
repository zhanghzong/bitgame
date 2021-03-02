package config

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/zhanghuizong/bitgame/component/apollo"
)

// 获取环境变量
// dev= 开发环境
// test= 测试环境
// pre= 预发环境
// prod= 生产环境
// app.env=dev
func GetAppEnv() string {
	return viper.GetString("app.env")
}

// 是否启用加密通信
// app.auth=true
func GetAppAuth() bool {
	if viper.IsSet("app.auth") {
		return viper.GetBool("app.auth")
	}

	return apollo.Config.GetBoolValue("app.auth", true)
}

// RSA 私钥
// app.rsa.private=-----BEGIN RSA PRIVATE KEY-----\n....\n-----END RSA PRIVATE KEY-----
func GetAppRsaPrivate() string {
	if viper.IsSet("app.rsa.private") {
		return viper.GetString("app.rsa.private")
	}

	return apollo.Config.GetValue("app.rsa.private")
}

// RSA 公钥
// app.rsa.public=-----BEGIN PUBLIC KEY-----\n...\n-----END PUBLIC KEY-----
func GetAppRsaPublic() string {
	if viper.IsSet("app.rsa.public") {
		return viper.GetString("app.rsa.public")
	}

	return apollo.Config.GetValue("app.rsa.public")
}

// jwt 秘钥
// jwt.key=8a3d4b8a3f13bc8c013f13bc8c9c0000
func GetJwtKey() string {
	if viper.IsSet("jwt.key") {
		return viper.GetString("jwt.key")
	}

	return apollo.Config.GetValue("jwt.key")
}

// jwt 过期时间. 单位(小时)
// jwt.expired=48
func GetJwtExpired() int {
	if viper.IsSet("jwt.expired") {
		expired := viper.GetInt("jwt.expired")
		if expired <= 0 {
			expired = 48
		}

		return expired
	}

	return apollo.Config.GetIntValue("jwt.expired", 48)
}

// mysql 数据源
// mysql.host=localhost
// mysql.user=root
// mysql.passwd=123
// mysql.database=game_go_fishing
// mysql.port=3306
// mysql.charset=utf8
func GetMysqlDsn() string {
	var host string
	var user string
	var passwd string
	var database string
	var port string
	var charset string
	if viper.IsSet("mysql.host") && viper.IsSet("mysql.user") &&
		viper.IsSet("mysql.database") && viper.IsSet("mysql.port") &&
		viper.IsSet("mysql.charset") {
		host = viper.GetString("mysql.host")
		user = viper.GetString("mysql.user")
		passwd = viper.GetString("mysql.passwd")
		database = viper.GetString("mysql.database")
		port = viper.GetString("mysql.port")
		charset = viper.GetString("mysql.charset")
	} else {
		host = apollo.Config.GetValue("mysql.host")
		user = apollo.Config.GetValue("mysql.user")
		passwd = apollo.Config.GetValue("mysql.passwd")
		database = apollo.Config.GetValue("mysql.database")
		port = apollo.Config.GetValue("mysql.port")
		charset = apollo.Config.GetValue("mysql.charset")
	}

	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local", user, passwd, host, port, database, charset)
}

// 数据库连接池大小
// mysql.pool_size=10
func GetMysqlPoolSize() int {
	if viper.IsSet("mysql.poolSize") {
		return viper.GetInt("mysql.poolSize")
	}

	return apollo.Config.GetIntValue("mysql.poolSize", 0)
}

/** Redis 配置节点 **/

// host:port address
// redis.addr=localhost:6379
func GetRedisAddr() string {
	if viper.IsSet("redis.addr") {
		return viper.GetString("redis.addr")
	}

	return apollo.Config.GetValue("redis.addr")
}

// Optional password. Must match the password specified in the
// requirepass server configuration option (if connecting to a Redis 5.0 instance, or lower),
// or the User Password when connecting to a Redis 6.0 instance, or greater, that is using the Redis ACL system.
// redis.password=
func GetRedisPassword() string {
	if viper.IsSet("redis.password") {
		return viper.GetString("redis.password")
	}

	return apollo.Config.GetValue("redis.password")
}

// Database to be selected after connecting to the server.
// redis.db=1
func GetRedisDb() int {
	if viper.IsSet("redis.db") {
		return viper.GetInt("redis.db")
	}

	return apollo.Config.GetIntValue("redis.db", 1)
}

// Maximum number of socket connections.
// Default is 10 connections per every CPU as reported by runtime.NumCPU.
// redis.poolSize=30
func GetRedisPoolSize() int {
	if viper.IsSet("redis.poolSize") {
		return viper.GetInt("redis.poolSize")
	}

	return apollo.Config.GetIntValue("redis.poolSize", 10)
}

// Minimum number of idle connections which is useful when establishing
// new connection is slow.
// redis.minIdleConns=30
func GetRedisMinIdleConns() int {
	if viper.IsSet("redis.minIdleConns") {
		return viper.GetInt("redis.minIdleConns")
	}

	return apollo.Config.GetIntValue("redis.minIdleConns", 30)
}

// InsecureSkipVerify controls whether a client verifies the server's
// certificate chain and host name. If InsecureSkipVerify is true, crypto/tls
// accepts any certificate presented by the server and any host name in that
// certificate. In this mode, TLS is susceptible to machine-in-the-middle
// attacks unless custom verification is used. This should be used only for
// testing or in combination with VerifyConnection or VerifyPeerCertificate.
func GetInsecureSkipVerify() bool {
	if viper.IsSet("redis.insecureSkipVerify") {
		return viper.GetBool("redis.insecureSkipVerify")
	}

	return apollo.Config.GetBoolValue("redis.insecureSkipVerify", false)
}

/** java 配置节点 */

// 游戏ID
//java.gameId=10008
func GetJavaGameId() string {
	if viper.IsSet("java.gameId") {
		return viper.GetString("java.gameId")
	}

	return apollo.Config.GetValue("java.gameId")
}

// 客户端ID
//java.clientId=game-fishing
func GetJavaClientId() string {
	if viper.IsSet("java.clientId") {
		return viper.GetString("java.clientId")
	}

	return apollo.Config.GetValue("java.clientId")
}

// 客户端授权秘钥
// java.clientSecret=f857f55b86f04b78824ad3a94948a584
func GetJavaClientSecret() string {
	if viper.IsSet("java.clientSecret") {
		return viper.GetString("java.clientSecret")
	}

	return apollo.Config.GetValue("java.clientSecret")
}

// 接口请求密钥
// java.apiKey=009093eb938e4f0e97579132d29e235d
func GetJavaApiKey() string {
	if viper.IsSet("java.apiKey") {
		return viper.GetString("java.apiKey")
	}

	return apollo.Config.GetValue("java.apiKey")
}

// 接口域名地址
// java.serverApi=http://api.btgame.club
func GetJavaServerApi() string {
	if viper.IsSet("java.serverApi") {
		return viper.GetString("java.serverApi")
	}

	return apollo.Config.GetValue("java.serverApi")
}

// 渠道ID
// java.channelId=BITGAME
func GetJavaChannelId() string {
	if viper.IsSet("java.channelId") {
		return viper.GetString("java.channelId")
	}

	return apollo.Config.GetValue("java.channelId")
}

// 游戏一级域名
// java.domainKey=btgame.club
func GetJavaDomainKey() string {
	if viper.IsSet("java.domainKey") {
		return viper.GetString("java.domainKey")
	}

	return apollo.Config.GetValue("java.domainKey")
}

// RSA
//java.rsa.public=-----BEGIN PUBLIC KEY-----\n....\n-----END PUBLIC KEY-----
func GetJavaRsaPublic() string {
	if viper.IsSet("java.rsa.public") {
		return viper.GetString("java.rsa.public")
	}

	return apollo.Config.GetValue("java.rsa.public")
}

// # 监控 URL 接口地址
// telegram.url=http://fat.monitor.testbitgame.com/
func GetTelegramUrl() string {
	if viper.IsSet("telegram.url") {
		return viper.GetString("telegram.url")
	}

	return apollo.Config.GetValue("telegram.url")
}

// telegram 聊天群ID
// telegram.chatId=-429713498
func GetTelegramChatId() string {
	if viper.IsSet("telegram.chatId") {
		return viper.GetString("telegram.chatId")
	}

	return apollo.Config.GetValue("telegram.chatId")
}
