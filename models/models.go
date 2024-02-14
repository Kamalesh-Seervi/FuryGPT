package models

type Config struct {
	APIKey        string
	RedisPassword string
	RedisDB       string
	RedisURL      string
}


type User struct {
    ID       uint   `json:"id"`
    Username string `json:"username"`
    Password string `json:"-"`
    APIKey   string `json:"apiKey"`
}