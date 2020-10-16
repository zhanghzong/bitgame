package definition

type ClientManagerInterface interface {
	RedisDispatch(msg *RedisChannel)
}
