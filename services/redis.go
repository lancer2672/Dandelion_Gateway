package services

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client

type RedisConfig struct {
	RedisURL      string
	RedisUsername string
	RedisPassword string
}

func SetValue(key string, value string, exp time.Duration) {
	redisClient.Set(context.Background(), key, value, exp)
}

func GetValue(key string) (string, error) {
	return redisClient.Get(context.Background(), key).Result()
}

func RemoveValue(key string) {
	redisClient.Del(context.Background(), key)
}
