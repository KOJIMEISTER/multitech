package middleware

import (
	"context"
	"multitech/pkg/storage"
	"multitech/pkg/testutils"
	"multitech/pkg/testutils/mocks"
	"net/http"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware(t *testing.T) {
	validToken, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, &Claims{
		UserID: 1,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
		},
	}).SignedString([]byte("test-secret"))

	expiredToken, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, &Claims{
		UserID: 1,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
		},
	}).SignedString([]byte("test-secret"))

	tests := []struct {
		name           string
		token          string
		mockSessSetup  func(*mocks.MockSessionsRepository)
		envSetup       func(*mocks.EnvMock)
		expectedStatus int
		expectedError  string
	}{
		{
			name:           "Missing authorization header",
			expectedStatus: http.StatusUnauthorized,
			expectedError:  `{"error":"Authorization header required"}`,
		},
		{
			name:           "Invalid token format",
			token:          "Invalid",
			expectedStatus: http.StatusUnauthorized,
			expectedError:  `{"error":"Invalid token"}`,
		},
		{
			name:  "Expired token",
			token: "Bearer " + expiredToken,
			envSetup: func(em *mocks.EnvMock) {
				em.Set("JWT_SECRET", "test-secret")
			},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  `{"error":"Invalid token"}`,
		},
		{
			name:  "Invalid signature",
			token: "Bearer " + validToken,
			envSetup: func(em *mocks.EnvMock) {
				em.Set("JWT_SECRET", "wrong-secret")
			},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  `{"error":"Invalid token"}`,
		},
		{
			name:  "Valid token but missing session",
			token: "Bearer " + validToken,
			envSetup: func(em *mocks.EnvMock) {
				em.Set("JWT_SECRET", "test-secret")
			},
			mockSessSetup: func(msr *mocks.MockSessionsRepository) {
				msr.GetSessionFunc = func(ctx context.Context, token string) (uint, error) {
					return 0, storage.ErrSessionExists
				}
			},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  `{"error":"Invalid or expired session"}`,
		},
		{
			name:  "Valid token with session",
			token: "Bearer " + validToken,
			envSetup: func(em *mocks.EnvMock) {
				em.Set("JWT_SECRET", "test-secret")
			},
			mockSessSetup: func(msr *mocks.MockSessionsRepository) {
				msr.GetSessionFunc = func(ctx context.Context, token string) (uint, error) {
					return 1, nil
				}
			},
			expectedStatus: http.StatusOK,
		},
	}

	originEnv := testutils.CaptureOriginEnv()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSessRepo := mocks.NewDefaultSessionsMock()
			mockEnv := mocks.NewEnvMock()

			if tt.mockSessSetup != nil {
				tt.mockSessSetup(mockSessRepo)
			}
			if tt.envSetup != nil {
				tt.envSetup(mockEnv)
				mockEnv.Apply()
				defer mockEnv.Restore(originEnv)
			}

			ctx, recorder := testutils.NewTestContext()
			if tt.token != "" {
				ctx.Request.Header.Set("Authorization", tt.token)
			}

			middleware := NewAuthMiddleware(mockSessRepo)
			handler := middleware.Middleware()
			handler(ctx)

			assert.Equal(t, tt.expectedStatus, recorder.Code)
			if tt.expectedError != "" {
				assert.Equal(t, tt.expectedError, recorder.Body.String())
			}
		})
	}
}
