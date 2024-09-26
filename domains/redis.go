package domains

import "time"

type RedisRepository interface {
	RedisUserGet(redisKey string) (string, error)
	RedisUserSet(redisKey string, redisValue string, timeToLive time.Duration) error
	RedisUserHMSet(redisKey string, redisValue interface{}, timeToLive time.Duration) error
	RedisUserHMGetAll(redisKey string) (map[string]string, error)

	RedisAuthHMSet(redisKey string, redisValue interface{}, timeToLive time.Duration) error
	RedisAuthHMGetFields(redisKey string, fields []string) (map[string]interface{}, error)
	RedisAuthHSetUpdateField(redisKey string, field string, value interface{}) error
}