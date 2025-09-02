package testutils

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

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

func CaptureOrigingEnv() *map[string]string {
	envs := make(map[string]string)
	for _, env := range os.Environ() {
		if idx := strings.Index(env, "="); idx >= 0 {
			envs[env[:idx]] = env[idx+1:]
		}
	}
	return &envs
}
