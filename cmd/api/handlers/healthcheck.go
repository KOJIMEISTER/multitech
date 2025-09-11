package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type HealthCheck struct {
	redisClient *redis.Client
}

func NewHealthCheck(redisClient *redis.Client) *HealthCheck {
	return &HealthCheck{
		redisClient: redisClient,
	}
}

// @Summary Health check
// @Description Check if the service is running
// @Tags system
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /health [get]
func (health *HealthCheck) Handler(ctx *gin.Context) {
	err := health.redisClient.Ping(ctx.Request.Context()).Err()
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
