package main

import (
	"multitech/cmd/api/handlers"
	"multitech/internal/config"
	"multitech/middleware"
	"multitech/pkg/storage"

	"github.com/gin-gonic/gin"
)

func main() {

	config.LoadEnv()

	storage.InitRedis()

	router := gin.Default()

	router.GET("/health", handlers.HealthCheck)

	router.POST("/register", handlers.RegisterHandler)

	router.POST("/login", handlers.LoginHandler)

	router.GET("/protected", middleware.AuthMiddleware(), handlers.ProtectedHandler)

	router.Run(":8080")
}
