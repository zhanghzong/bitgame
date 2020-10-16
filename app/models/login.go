package models

import (
	"github.com/sirupsen/logrus"
	"github.com/zhanghuizong/bitgame/app/constants/redisConst"
	"github.com/zhanghuizong/bitgame/component/redis"
)

type LoginModel struct {
}

func (t LoginModel) AddSocketId(uid string, socketId string) {
	key := redisConst.SingleLogin

	cmd := redis.Redis.HSet(key, uid, socketId)
	_, err := cmd.Result()
	if err != nil {
		logrus.Warn("设置 用户与websocket绑定关系 Redis 异常", err, key, uid, socketId)
	}
}

func (t LoginModel) GetSocketId(uid string) string {
	key := redisConst.SingleLogin
	return redis.Redis.HGet(key, uid).Val()
}

func (t LoginModel) DelSocketId(uid string) int {
	key := redisConst.SingleLogin
	val := redis.Redis.HDel(key, uid).Val()

	return int(val)
}
