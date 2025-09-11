package handlers

import (
	"context"
	"encoding/json"
	"multitech/internal/models"
	"multitech/pkg/testutils"
	"multitech/pkg/testutils/mocks"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegisterHandler(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    string
		mockUserSetup  func(*mocks.MockUserRepository)
		envSetup       func(*mocks.EnvMock)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:        "Success",
			requestBody: `{"username": "user", "email": "test@mail.com", "password": "testpass", "user_id":"*"}`,
			mockUserSetup: func(mur *mocks.MockUserRepository) {
				mur.CreateUserFunc = func(ctx context.Context, user *models.User) error {
					return nil
				}
			},
			envSetup: func(em *mocks.EnvMock) {
				em.Set("JWT_SECRET", "testsecret")
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   `{"message":"User created successfully","user_id":"*"}`,
		},
	}

	originEnv := testutils.CaptureOriginEnv()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserRepo := mocks.NewDefaultUserMock()
			mockEnv := mocks.NewEnvMock()

			if tt.mockUserSetup != nil {
				tt.mockUserSetup(mockUserRepo)
			}
			if tt.envSetup != nil {
				tt.envSetup(mockEnv)
				mockEnv.Apply()
				defer mockEnv.Restore(originEnv)
			}

			ctx, recorder := testutils.NewTestContext()
			testutils.SetJSONBody(ctx, tt.requestBody)

			registerHandler := NewRegisterHandler(mockUserRepo)
			registerHandler.Handler(ctx)

			assert.Equal(t, tt.expectedStatus, recorder.Code)
			if tt.expectedBody != "" {
				if strings.Contains(tt.expectedBody, "*") {
					var expected, actual map[string]interface{}
					if err := json.Unmarshal([]byte(tt.expectedBody), &expected); err != nil {
						t.Fatalf("Failed to parse expected JSON: %v", err)
					}
					if err := json.Unmarshal(recorder.Body.Bytes(), &actual); err != nil {
						t.Fatalf("Failed to parse expected JSON: %v", err)
					}
					for key, value := range expected {
						if strVal, ok := value.(string); ok && strVal == "*" {
							delete(expected, key)
							delete(actual, key)
						}
					}
					assert.Equal(t, expected, actual)
				} else {
					assert.Equal(t, tt.expectedBody, recorder.Body.String())
				}
			}
		})
	}
}
