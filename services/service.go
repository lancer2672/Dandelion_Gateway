package services

import (
	"log"

	"github.com/lancer2672/Dandelion_Gateway/internal/utils"
	"github.com/redis/go-redis/v9"
)

func ConfigServices() {
	redisConfig := RedisConfig{
		RedisURL:      utils.ConfigIns.RedisURL,
		RedisUsername: utils.ConfigIns.RedisUsername,
		RedisPassword: utils.ConfigIns.RedisPassword,
	}

	connectRedis("default", redisConfig)
}

func connectRedis(clientName string, redisConfig RedisConfig) {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     redisConfig.RedisURL,
		Password: redisConfig.RedisPassword,
		Username: redisConfig.RedisUsername,
		DB:       0, // use default DB
	})
	log.Println("Connected to Redis")
}
