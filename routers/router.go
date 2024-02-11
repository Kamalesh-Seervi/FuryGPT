package routers

import (
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/", index)
	router.POST("/run", run)
	log.Println("\033[93mFuryAI started. Press CTRL+C to quit.\033[0m")
	router.Run()
}
