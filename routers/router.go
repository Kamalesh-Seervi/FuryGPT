package routers

import (
	"log"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/kamalesh-seervi/simpleGPT/handlers"
)

func Router() {
	router := gin.Default()
	router.Use(gin.Logger())

	staticPath := "static"
	router.Static("/static", staticPath)
	templatesPath := filepath.Join(staticPath, "*.html")
	router.LoadHTMLGlob(templatesPath)

	router.GET("/", handlers.Index)

	router.POST("/run", handlers.Run)
	router.POST("/fetchHistory", handlers.HandleFetchHistory)

	log.Println("\033[93mFuryAI started. Press CTRL+C to quit.\033[0m")
	router.Run()
}
