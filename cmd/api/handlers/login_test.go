package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"multitech/internal/models"
	"multitech/pkg/storage"
	"multitech/pkg/testutils"
	"multitech/pkg/testutils/mocks"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoginHandler(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    string
		mockUserSetup  func(*mocks.MockUserRepository)
		mockSessSetup  func(*mocks.MockSessionsRepository)
		envSetup       func(*mocks.EnvMock)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:        "Success",
			requestBody: `{"username": "testuser", "password": "testpass"}`,
			mockUserSetup: func(mur *mocks.MockUserRepository) {
				mur.GetUserByUsernameFunc = func(ctx context.Context, username string) (*models.User, error) {
					return &models.User{
						ID:       1,
						Username: username,
						Password: "$2a$10$V2ezVCk4gXWAQhkCHV4wfOq0b/LtD0PHnx.t1GcALDYAZ96NUChrm",
					}, nil
				}
			},
			envSetup: func(em *mocks.EnvMock) {
				em.Set("JWT_SECRET", "testsecret")
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"token":"*", "user":{"id":1,"username":"testuser","email":""}}`,
		},
		{
			name:        "Invalid Password",
			requestBody: `{"username": "testuser", "password": "wrongpass"}`,
			mockUserSetup: func(mur *mocks.MockUserRepository) {
				mur.GetUserByUsernameFunc = func(ctx context.Context, username string) (*models.User, error) {
					return &models.User{
						ID:       1,
						Username: username,
						Password: "$2a$10$V2ezVCk4gXWAQhkCHV4wfOq0b/LtD0PHnx.t1GcALDYAZ96NUChrm",
					}, nil
				}
			},
			envSetup: func(em *mocks.EnvMock) {
				em.Set("JWT_SECRET", "testsecret")
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"error": "Invalid credentials"}`,
		},
		{
			name:        "User Not Found",
			requestBody: `{"username": "nonexistent", "password": "testpass"}`,
			mockUserSetup: func(mur *mocks.MockUserRepository) {
				mur.GetUserByUsernameFunc = func(ctx context.Context, username string) (*models.User, error) {
					return nil, storage.ErrUserNotFound
				}
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   fmt.Sprintf(`{"error": "%s"}`, storage.ErrUserNotFound.Error()),
		},
	}

	originEnv := testutils.CaptureOriginEnv()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserRepo := mocks.NewDefaultUserMock()
			mockSessRepo := mocks.NewDefaultSessionsMock()
			mockEnv := mocks.NewEnvMock()

			if tt.mockUserSetup != nil {
				tt.mockUserSetup(mockUserRepo)
			}
			if tt.mockSessSetup != nil {
				tt.mockSessSetup(mockSessRepo)
			}
			if tt.envSetup != nil {
				tt.envSetup(mockEnv)
				mockEnv.Apply()
				defer mockEnv.Restore(originEnv)
			}

			ctx, recorder := testutils.NewTestContext()
			testutils.SetJSONBody(ctx, tt.requestBody)

			loginHandler := NewLoginHandler(mockUserRepo, mockSessRepo)
			loginHandler.Handler(ctx)

			assert.Equal(t, tt.expectedStatus, recorder.Code)
			if tt.expectedBody != "" {
				if strings.Contains(tt.expectedBody, "*") {
					var expected, actual map[string]interface{}
					if err := json.Unmarshal([]byte(tt.expectedBody), &expected); err != nil {
						t.Fatalf("failed to parse expected JSON: %v", err)
					}
					if err := json.Unmarshal(recorder.Body.Bytes(), &actual); err != nil {
						t.Fatalf("failed to parse expected JSON: %v", err)
					}
					for key, value := range expected {
						if strVal, ok := value.(string); ok && strVal == "*" {
							delete(expected, key)
							delete(actual, key)
						}
					}
					assert.Equal(t, expected, actual)
				} else {
					assert.JSONEq(t, tt.expectedBody, recorder.Body.String())
				}
			}
		})
	}
}
