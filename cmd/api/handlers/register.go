package handlers

import (
	"errors"
	"multitech/internal/models"
	"multitech/pkg/storage"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type RegisterHandler struct {
	userRepo storage.UserRepository
}

func NewRegisterHandler(userRepo storage.UserRepository) *RegisterHandler {
	return &RegisterHandler{
		userRepo: userRepo,
	}
}

// @Summary Register new user
// @Description Create a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param user body models.RegisterCredentials true "User registration data"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /register [post]
func (register *RegisterHandler) Handler(ctx *gin.Context) {
	var regCreds models.RegisterCredentials
	if err := ctx.ShouldBindJSON(&regCreds); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	user := models.User{
		Username:  regCreds.Username,
		Email:     regCreds.Email,
		Password:  regCreds.Password,
		CreatedAt: time.Now(),
	}

	if err := user.HashPassword(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error hashing password",
		})
		return
	}

	if err := register.userRepo.CreateUser(ctx.Request.Context(), &user); err != nil {
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
