package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
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
	// Redis Configuration
	opts := &redis.Options{
		Addr:     "redis:6379",
		Password: config.RedisPassword,
		DB:       0,
	}

	rdb = redis.NewClient(opts)

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		panic("Error connecting to Redis: " + err.Error())
	}
	fmt.Println("Connected to Redis")
}

func HashAPIKey(apiKey string) string {
	hasher := sha256.New()
	hasher.Write([]byte(apiKey))
	return hex.EncodeToString(hasher.Sum(nil))
}

func StoreHistory(userID, input, response string) {
	hashedUserID := HashAPIKey(userID)
	historyKey := fmt.Sprintf("history:%s", hashedUserID)
	log.Printf("Storing in Redis: Key=%s, Input=%s, Response=%s", userID, input, response)
	maxHistoryLength := 10

	// Create a JSON-like structure to store both input and response
	entry := fmt.Sprintf(`{"input": "%s", "response": "%s"}`, input, response)

	rdb.LPush(context.Background(), historyKey, entry)
	rdb.LTrim(context.Background(), historyKey, 0, int64(maxHistoryLength-1))
	expiration := time.Hour * 24
	rdb.Expire(context.Background(), historyKey, expiration)
}

func GetHistory(userID string) []string {
	hashedUserID := HashAPIKey(userID)
	historyKey := fmt.Sprintf("history:%s", hashedUserID)
	result, err := rdb.LRange(context.Background(), historyKey, 0, -1).Result()
	if err != nil {
		log.Printf("Error retrieving history from Redis: %v", err)
		return []string{}
	}
	return result
}
