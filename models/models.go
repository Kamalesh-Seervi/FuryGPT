package models

type Config struct {
	APIKey        string
	RedisPassword string
	RedisDB       string
	RedisURL      string
	SecretKey     string
}

type User struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}
