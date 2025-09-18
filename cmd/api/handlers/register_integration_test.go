package handlers

import (
	"encoding/json"
	"multitech/internal/models"
	"multitech/pkg/storage"
	"multitech/pkg/testutils"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegisterHandlerSuccess(t *testing.T) {
	tx := testutils.TestDB.Begin()
	defer tx.Rollback()

	userRepo := storage.NewGormUserRepository(tx)

	ctx, recorder := testutils.NewTestContext()
	testutils.SetJSONBody(ctx, `{"username":"newuser","email":"newuser@example.com","password":"securepassword123"}`)

	handler := NewRegisterHandler(userRepo)
	handler.Handler(ctx)

	assert.Equal(t, http.StatusCreated, recorder.Code)

	var response map[string]interface{}
	assert.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &response))
	assert.Equal(t, "User created successfully", response["message"])
	assert.NotEmpty(t, response["user_id"])

	var user models.User
	err := tx.Where("username = ?", "newuser").First(&user).Error
	assert.NoError(t, err)
	assert.Equal(t, "newuser@example.com", user.Email)
	assert.NotEmpty(t, user.Password)
}

func TestRegisterHandlerDuplicateUsername(t *testing.T) {
	tx := testutils.TestDB.Begin()
	defer tx.Rollback()

	exisingUser := &models.User{
		Username: "existinguser",
		Email:    "existing@example.com",
		Password: "hashedpassword",
	}
	assert.NoError(t, tx.Create(exisingUser).Error)

	userRepo := storage.NewGormUserRepository(tx)

	ctx, recorder := testutils.NewTestContext()
	testutils.SetJSONBody(ctx, `{"username":"existinguser","email":"new@example.com","password":"password123"}`)

	handler := NewRegisterHandler(userRepo)
	handler.Handler(ctx)

	assert.Equal(t, http.StatusConflict, recorder.Code)
	assert.JSONEq(t, `{"error":"User already exists"}`, recorder.Body.String())
}
