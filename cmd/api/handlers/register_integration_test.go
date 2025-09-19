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

func TestRegisterHandlerInvalidData(t *testing.T) {
	tx := testutils.TestDB.Begin()
	defer tx.Rollback()

	userRepo := storage.NewGormUserRepository(tx)

	testCases := []struct {
		name          string
		requestBody   string
		expectedError string
	}{
		{
			name:          "Empty Username",
			requestBody:   `{"username":"","email":"test@example.com","password":"password123"}`,
			expectedError: `{"error":"Username must be at least 3 characters"}`,
		},
		{
			name:          "Invalid Email",
			requestBody:   `{"username":"testuser","email":"notanemail","password":"password123"}`,
			expectedError: `{"error":"Invalid email format"}`,
		},
		{
			name:          "Short Password",
			requestBody:   `{"username":"testuser","email":"test@example.com","password":"short"}`,
			expectedError: `{"error":"Password must be at least 8 characters"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, recorder := testutils.NewTestContext()
			testutils.SetJSONBody(ctx, tc.requestBody)

			handler := NewRegisterHandler(userRepo)
			handler.Handler(ctx)

			assert.Equal(t, http.StatusBadRequest, recorder.Code)
			assert.JSONEq(t, tc.expectedError, recorder.Body.String())
		})
	}
}
