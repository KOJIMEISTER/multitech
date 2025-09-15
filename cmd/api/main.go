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
	"context"
	"log"
	"multitech/cmd/api/handlers"
	_ "multitech/docs"
	"multitech/internal/config"
	"multitech/middleware"
	"multitech/pkg/storage"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	// postgresClient, err := storage.InitPostgres()
	// if err != nil {
	// 	log.Fatalf("PostgreSQL init error: %v", err)
	// }

	userRepo := storage.NewRedisUserRepository(redisClient)
	sessRepo := storage.NewRedisSessionRepository(redisClient)

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

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Println("Server forced to shudown:", err)
	}

	if err := redisClient.Close(); err != nil {
		log.Println("Error closing Redis:", err)
	}

	log.Println("Server exiting")

}
