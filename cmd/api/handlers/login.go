package handlers

import (
	"errors"
	"multitech/internal/models"
	"multitech/middleware"
	"multitech/pkg/storage"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func LoginHandler(ctx *gin.Context) {
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
