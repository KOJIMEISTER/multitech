package testutils

import (
	"context"
	"fmt"
	"io"
	"log"
	"multitech/internal/models"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	TestDB     *gorm.DB
	testDBOnce sync.Once

	TestRedis     *redis.Client
	testRedisOnce sync.Once
)

func InitTestDB(dsn string) *gorm.DB {
	var initErr error
	testDBOnce.Do(func() {
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			initErr = fmt.Errorf("Failed to connect to test database: %w", err)
			return
		}

		sqlDB, err := db.DB()
		if err != nil {
			initErr = fmt.Errorf("Failed to get database instance: %w", err)
			return
		}

		sqlDB.SetMaxIdleConns(5)
		sqlDB.SetMaxOpenConns(10)

		TestDB = db
	})

	if initErr != nil {
		log.Fatal(initErr)
	}
	return TestDB
}

func InitTestRedis(url string) *redis.Client {
	var initErr error
	testRedisOnce.Do(func() {
		opt, err := redis.ParseURL(url)
		if err != nil {
			initErr = fmt.Errorf("Failed to parse redis URL: %w", err)
		}

		client := redis.NewClient(opt)
		if err := client.Ping(context.Background()).Err(); err != nil {
			initErr = fmt.Errorf("Failed to connect to redis: %w", err)
			return
		}

		TestRedis = client
	})

	if initErr != nil {
		log.Fatal(initErr)
	}
	return TestRedis
}

func RunMigrations(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
	)
}

func NewTestContext() (*gin.Context, *httptest.ResponseRecorder) {
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = &http.Request{
		Header: make(http.Header),
	}
	return ctx, recorder
}

func SetJSONBody(ctx *gin.Context, body string) {
	ctx.Request.Header.Set("Content-Type", "application/json")
	ctx.Request.Body = io.NopCloser(strings.NewReader(body))
}

func CaptureOriginEnv() *map[string]string {
	envs := make(map[string]string)
	for _, env := range os.Environ() {
		if idx := strings.Index(env, "="); idx >= 0 {
			envs[env[:idx]] = env[idx+1:]
		}
	}
	return &envs
}
