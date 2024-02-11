package service

import (
	"context"
	"time"

	"github.com/kamalesh-seervi/simpleGPT/utils"
	"github.com/redis/go-redis/v9"
)

var rdb *redis.Client

func InitRedis() {
	config, _ := utils.LoadConfig()
	opts, err := redis.ParseURL(config.RedisURL)
	if err != nil {
		panic(err)
	}
	rdb = redis.NewClient(opts)
}

func GetPromptFromCache(input string) (string, error) {
	ctx := context.Background()
	response, err := rdb.Get(ctx, input).Result()
	if err == redis.Nil {
		// Key not found in cache
		return "", nil
	} else if err != nil {
		return "", err
	}
	return response, nil
}

func SetPromptInCache(input, response string) error {
	ctx := context.Background()
	return rdb.Set(ctx, input, response, time.Hour).Err()
}
