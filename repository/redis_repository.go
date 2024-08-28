package repository

import (
	"auth-service/bootstrap"
	"auth-service/domains"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

type redisRepository struct {
	redisUser *redis.Client
	redisAuth *redis.Client
	ctx       context.Context
}

func NewRedisRepository() domains.RedisRepository {
	return &redisRepository{
		bootstrap.RedisAuth,
		bootstrap.RedisUser,
		context.Background(),
	}
}

func (r *redisRepository) RedisUserGet(redisKey string) (string, error) {
	val, err := r.redisUser.Get(r.ctx, redisKey).Result()

	return val, err
}

func (r *redisRepository) RedisUserSet(redisKey, redisValue string, timeToLive time.Duration) error {
	if timeToLive > 0 {
		err := r.redisUser.Set(r.ctx, redisKey, redisValue, timeToLive).Err()
		if err != nil {
			// Handle the error if needed
			log.Printf("Failed to set key with TTL: %v", err)
			return err
		}
	} else {
		err := r.redisUser.Set(r.ctx, redisKey, redisValue, -1).Err()
		if err != nil {
			// Handle the error if needed
			log.Printf("Failed to set key without TTL: %v", err)
			return err
		}
	}

	return nil
}

func (r *redisRepository) RedisUserHMSet(redisKey string, redisValue interface{}, timeToLive time.Duration) error {

	errSaveRedis := r.redisUser.HMSet(r.ctx, redisKey, redisValue).Err()

	if errSaveRedis != nil {
		fmt.Println(errSaveRedis)
		return errSaveRedis
	}

	if timeToLive > 0 {
		errSaveTTL := r.redisUser.Expire(r.ctx, redisKey, timeToLive).Err()

		if errSaveTTL != nil {
			fmt.Println(errSaveTTL)
			return errSaveTTL
		}
	}

	return nil
}

func (r *redisRepository) RedisAuthHMSet(redisKey string, redisValue interface{}, timeToLive time.Duration) error {
	errSaveRedis := r.redisAuth.HMSet(r.ctx, redisKey, redisValue).Err()

	if errSaveRedis != nil {
		fmt.Println(errSaveRedis)
		return errSaveRedis
	}

	if timeToLive > 0 {
		errSaveTTL := r.redisAuth.Expire(r.ctx, redisKey, timeToLive).Err()

		if errSaveTTL != nil {
			fmt.Println(errSaveTTL)
			return errSaveTTL
		}
	}

	return nil
}
