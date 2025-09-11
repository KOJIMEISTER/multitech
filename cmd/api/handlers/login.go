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

type LoginHandler struct {
	userRepo storage.UserRepository
	sessRepo storage.SessionsRepository
}

func NewLoginHandler(userRepo storage.UserRepository, sessRepo storage.SessionsRepository) *LoginHandler {
	return &LoginHandler{
		userRepo: userRepo,
		sessRepo: sessRepo,
	}
}

// @Summary User login
// @Description Authenticate user and return JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body models.LoginCredentials true "Login credentials"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /login [post]
func (login *LoginHandler) Handler(ctx *gin.Context) {
	var creds models.LoginCredentials
	if err := ctx.ShouldBindJSON(&creds); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	user, err := login.userRepo.GetUserByUsername(ctx.Request.Context(), creds.Username)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": storage.ErrUserNotFound.Error(),
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

	if err := login.sessRepo.StoreSession(ctx.Request.Context(), token, user.ID, 24*time.Hour); err != nil {
		if errors.Is(err, storage.ErrSessionExists) {
			ctx.JSON(http.StatusConflict, gin.H{
				"error": storage.ErrSessionExists.Error(),
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error creating session: " + err.Error(),
		})
		return
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
