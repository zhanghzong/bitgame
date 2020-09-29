package login

import (
	"github.com/zhanghuizong/bitgame/app/middleware/redis"
	"github.com/zhanghuizong/bitgame/utils"
	"log"
)

type Model struct {
}

var singleLogin = ":bitgame:ws:single_login"

func (t Model) AddSocketId(uid string, socketId string) {
	key := utils.GetRedisPrefix() + singleLogin

	cmd := redis.Redis.HSet(key, uid, socketId)
	_, err := cmd.Result()
	if err != nil {
		log.Println("设置 用户与websocket绑定关系 Redis 异常", err, key, uid, socketId)
	}
}

func (t Model) GetSocketId(uid string) string {
	key := utils.GetRedisPrefix() + singleLogin
	return redis.Redis.HGet(key, uid).Val()
}

func (t Model) DelSocketId(uid string) int {
	key := utils.GetRedisPrefix() + singleLogin
	val := redis.Redis.HDel(key, uid).Val()

	return int(val)
}
