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
	authGroup := router.Group("/auth")
	{
		authGroup.POST("/signup", handlers.Signup)
		authGroup.POST("/login", handlers.Login)
	}
	router.GET("/", handlers.Index)
	router.POST("/run", handlers.Run)
	log.Println("\033[93mFuryAI started. Press CTRL+C to quit.\033[0m")
	router.Run()
}
