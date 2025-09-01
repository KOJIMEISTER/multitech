package handlers

import (
	"errors"
	"multitech/internal/models"
	"multitech/pkg/storage"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func RegisterHandler(ctx *gin.Context) {
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
