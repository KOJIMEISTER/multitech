package handlers

import (
	"errors"
	"fmt"
	"multitech/internal/models"
	"multitech/pkg/storage"
	"net/http"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	minPasswordLength = 8
	maxPasswordLength = 254
	minUsernameLength = 3
	maxUsernameLength = 254
	minEmailLength    = 3
	maxEmailLength    = 254
	emailRegexPattern = `^[a-zA-Z0-9.!#$%&'*+/=?^_` + "`" + `{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$`
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

	if len(regCreds.Password) < minPasswordLength {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Password must be at least %d characters", minPasswordLength),
		})
		return
	}

	if len(regCreds.Username) < minUsernameLength {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Username must be at least %d characters", minUsernameLength),
		})
		return
	}

	if err := register.validateEmail(regCreds.Email); err != nil {
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

func (register *RegisterHandler) validateEmail(email string) error {
	if len(email) < minEmailLength || len(email) > maxEmailLength {
		return fmt.Errorf("Email must be between %d-%d characters", minEmailLength, maxEmailLength)
	}

	matched, err := regexp.MatchString(emailRegexPattern, email)
	if err != nil || !matched {
		return errors.New("Invalid email format")
	}
	return nil
}
