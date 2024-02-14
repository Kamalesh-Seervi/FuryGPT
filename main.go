package main

import (
	"github.com/kamalesh-seervi/simpleGPT/routers"
	"github.com/kamalesh-seervi/simpleGPT/service"
)

func main() {
	service.InitRedis()
	service.DBSetup()
	routers.Router()
}
