package ws

// 根据 Uid 消息单播
func Single(uid string, cmd string, data interface{}) {
	single(uid, pushSuccess(cmd, data))
}
