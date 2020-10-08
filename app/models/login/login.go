package login

import (
	"github.com/zhanghuizong/bitgame/app/constants/redisConst"
	"github.com/zhanghuizong/bitgame/app/middleware/redis"
	"log"
)

type Model struct {
}

func (t Model) AddSocketId(uid string, socketId string) {
	key := redisConst.SingleLogin

	cmd := redis.Redis.HSet(key, uid, socketId)
	_, err := cmd.Result()
	if err != nil {
		log.Println("设置 用户与websocket绑定关系 Redis 异常", err, key, uid, socketId)
	}
}

func (t Model) GetSocketId(uid string) string {
	key := redisConst.SingleLogin
	return redis.Redis.HGet(key, uid).Val()
}

func (t Model) DelSocketId(uid string) int {
	key := redisConst.SingleLogin
	val := redis.Redis.HDel(key, uid).Val()

	return int(val)
}
