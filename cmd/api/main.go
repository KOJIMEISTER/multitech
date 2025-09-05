package main

import (
	"log"
	"multitech/cmd/api/handlers"
	"multitech/internal/config"
	"multitech/middleware"
	"multitech/pkg/storage"

	"github.com/gin-gonic/gin"
)

func main() {

	config.LoadEnv()

	redisClient, err := storage.InitRedis()
	if err != nil {
		log.Fatalf("Error init redis: %v", err)
		return
	}

	userRepo := storage.NewUserRepository(redisClient)
	sessRepo := storage.NewSessionRepository(redisClient)

	healthCheck := handlers.NewHealthCheck(redisClient)
	loginHandler := handlers.NewLoginHandler(userRepo, sessRepo)
	registerHandler := handlers.NewRegisterHandler(userRepo)
	protectedHandler := handlers.NewProtectedHandler()

	authMiddleware := middleware.NewAuthMiddleware(sessRepo)

	router := gin.Default()

	router.GET("/health", healthCheck.Handler)
	router.GET("/protected", authMiddleware.Middleware(), protectedHandler.Handler)

	router.POST("/login", loginHandler.Handler)
	router.POST("/register", registerHandler.Handler)

	router.Run(":8080")
}
