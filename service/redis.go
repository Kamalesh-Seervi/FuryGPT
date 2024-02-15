package service

import (
	"context"
	"fmt"

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

// func GetUserSessionKey(userID string) string {
// 	return "user:" + userID
// }
