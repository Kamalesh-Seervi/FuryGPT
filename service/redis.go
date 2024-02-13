package service

import (
	"context"
	"fmt"
	"time"

	"github.com/kamalesh-seervi/simpleGPT/utils"
	"github.com/redis/go-redis/v9"
)

var rdb *redis.Client

func InitRedis() {
	config, err := utils.LoadConfig()
	if err != nil {
		panic("Error loading config: " + err.Error())
	}
	opts, err := redis.ParseURL(config.RedisURL)
	if err != nil {
		panic(err)
	}
	opts.Password = config.RedisPassword
	rdb = redis.NewClient(opts)

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		panic("Error connecting to Redis: " + err.Error())
	}

	fmt.Println("Connected to Redis")

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
