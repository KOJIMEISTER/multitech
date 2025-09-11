package main

// @title Mutitech API
// @version 1.0
// @description API for Multitech application
// @host localhost:8080
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token
// @BasePath /

import (
	"log"
	"multitech/cmd/api/handlers"
	_ "multitech/docs"
	"multitech/internal/config"
	"multitech/middleware"
	"multitech/pkg/storage"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET("/health", healthCheck.Handler)
	router.GET("/protected", authMiddleware.Middleware(), protectedHandler.Handler)

	router.POST("/login", loginHandler.Handler)
	router.POST("/register", registerHandler.Handler)

	router.Run(":8080")
}
