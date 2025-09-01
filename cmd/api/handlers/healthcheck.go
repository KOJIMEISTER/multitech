package handlers

import (
	"multitech/pkg/storage"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HealthCheck(ctx *gin.Context) {
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
