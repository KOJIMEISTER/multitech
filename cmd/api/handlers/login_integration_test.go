package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"multitech/internal/models"
	"multitech/pkg/storage"
	"multitech/pkg/testutils"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	containers, err := testutils.SetupContainers(ctx)
	if err != nil {
		fmt.Printf("Failed to setup containers: %v", err)
		os.Exit(1)
	}
	defer containers.Terminate(ctx)

	testutils.InitTestDB(containers.PostgresDSN)
	testutils.InitTestRedis(containers.RedisURL)

	testutils.RunMigrations(testutils.TestDB)

	os.Exit(m.Run())
}

func TestLoginHandlerSuccess(t *testing.T) {
	tx := testutils.TestDB.Begin()
	defer tx.Rollback()

	user := &models.User{
		Username: "testuser",
		Password: "testpass",
	}
	user.HashPassword()
	assert.NoError(t, tx.Create(user).Error)

	userRepo := storage.NewGormUserRepository(tx)
	sessRepo := storage.NewRedisSessionRepository(testutils.TestRedis)

	ctx, recorder := testutils.NewTestContext()
	testutils.SetJSONBody(ctx, `{"username":"testuser","password":"testpass"}`)

	os.Setenv("JWT_SECRET", "testsecret")
	defer os.Unsetenv("JWT_SECRET")

	handler := NewLoginHandler(userRepo, sessRepo)
	handler.Handler(ctx)

	assert.Equal(t, http.StatusOK, recorder.Code)

	var response map[string]interface{}
	assert.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &response))
	assert.NotEmpty(t, response["token"])
	assert.NotEmpty(t, response["user"])

	token := response["token"].(string)
	userID, err := sessRepo.GetSession(ctx, token)
	assert.NoError(t, err)
	assert.Equal(t, user.ID, userID)
}

func TestLoginHandlerInvalidPassword(t *testing.T) {
	tx := testutils.TestDB.Begin()
	defer tx.Rollback()

	user := &models.User{
		Username: "testuser",
		Password: "testpass",
	}
	user.HashPassword()
	assert.NoError(t, tx.Create(user).Error)

	userRepo := storage.NewGormUserRepository(tx)
	sessRepo := storage.NewRedisSessionRepository(testutils.TestRedis)

	ctx, recoder := testutils.NewTestContext()
	testutils.SetJSONBody(ctx, `{"username":"testuser","password":"wrongpass"}`)

	os.Setenv("JWT_SECRET", "testsecret")
	defer os.Unsetenv("JWT_SECRET")

	handler := NewLoginHandler(userRepo, sessRepo)
	handler.Handler(ctx)

	assert.Equal(t, http.StatusUnauthorized, recoder.Code)
	assert.JSONEq(t, `{"error":"Invalid credentials"}`, recoder.Body.String())
}

func TestLoginHandlerUserNotFound(t *testing.T) {
	tx := testutils.TestDB.Begin()
	defer tx.Rollback()

	userRepo := storage.NewGormUserRepository(tx)
	sessRepo := storage.NewRedisSessionRepository(testutils.TestRedis)

	ctx, recorder := testutils.NewTestContext()
	testutils.SetJSONBody(ctx, `{"username":"nonexistent","password":"testpass"}`)

	handler := NewLoginHandler(userRepo, sessRepo)
	handler.Handler(ctx)

	assert.Equal(t, http.StatusUnauthorized, recorder.Code)
	assert.JSONEq(t, `{"error":"User not found"}`, recorder.Body.String())
}
