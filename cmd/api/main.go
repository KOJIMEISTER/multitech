package main

import (
	"errors"
	"log"
	"multitech/internal/config"
	"multitech/internal/models"
	"multitech/middleware"
	"multitech/pkg/storage"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config.LoadEnv()
	storage.InitRedis()

	router := gin.Default()
	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "Hello World"})
	})

	router.GET("/health", healthCheck)

	router.POST("/register", registerHandler)

	router.POST("/login", loginHandler)

	router.GET("/protected", middleware.AuthMiddleware(), protectedHandler)

	router.Run(":8080")
}

func healthCheck(ctx *gin.Context) {
	err := storage.RedisClient.Ping(ctx.Request.Context()).Err()
	if err != nil {
		ctx.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "unhealthy",
			"redis":  err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "healthy",
		"redis":  "connected",
	})
}

func registerHandler(ctx *gin.Context) {
	var user models.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	user.CreatedAt = time.Now()

	if err := user.HashPassword(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error hashing password",
		})
		return
	}

	if err := storage.CreateUser(ctx.Request.Context(), &user); err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			ctx.JSON(http.StatusConflict, gin.H{
				"error": "User already exists",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error creating user: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"user_id": user.ID,
	})
}

func loginHandler(ctx *gin.Context) {
	var creds models.LoginCredentials
	if err := ctx.ShouldBindJSON(&creds); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	user, err := storage.GetUserByUsername(ctx.Request.Context(), creds.Username)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid credentials",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error retrieving user",
		})
		return
	}

	if err := user.CheckPassword(creds.Password); err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid credentials",
		})
		return
	}

	token, err := middleware.GenerateToken(user.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error generating token",
		})
		return
	}

	if err := storage.StoreSession(ctx.Request.Context(), token, user.ID, 24*time.Hour); err != nil {
		if errors.Is(err, storage.ErrSessionExists) {
			ctx.JSON(http.StatusConflict, gin.H{
				"error": "Session already exists",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error creating session: " + err.Error(),
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
	})
}

func protectedHandler(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")
	ctx.JSON(http.StatusOK, gin.H{
		"user_id": userID,
		"message": "Protected content",
	})
}
