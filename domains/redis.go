package domains

import "time"

type RedisRepository interface {
	RedisUserGet(redisKey string) (string, error)
	RedisUserSet(redisKey string, redisValue string, timeToLive time.Duration) error
	RedisUserHMSet(redisKey string, redisValue interface{}, timeToLive time.Duration) error
	RedisAuthHMSet(redisKey string, redisValue interface{}, timeToLive time.Duration) error
}
