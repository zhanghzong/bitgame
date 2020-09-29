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
	cmd := redis.Redis.HGet(key, uid)
	val, err := cmd.Result()
	if err != nil {
		log.Println("获取 用户与websocket绑定关系 Redis 异常", err, key, uid)
	}

	return val
}

func (t Model) DelSocketId(uid string) int {
	key := utils.GetRedisPrefix() + singleLogin
	cmd := redis.Redis.HDel(key, uid)
	val, err := cmd.Result()
	if err != nil {
		log.Println("删除 用户与websocket绑定关系 Redis 异常", err, key, uid)
	}

	return int(val)
}
