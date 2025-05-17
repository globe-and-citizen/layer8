package service

import (
	"errors"
	"globe-and-citizen/layer8/server/models"
	"os"
	"testing"
	"time"

	utilities "github.com/globe-and-citizen/layer8-utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRepository is a mock implementation of the Repository interface
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) GetClient(key string) (*models.Client, error) {
	args := m.Called(key)
	return args.Get(0).(*models.Client), args.Error(1)
}

func (m *MockRepository) GetClientByURL(backendURL string) (*models.Client, error) {
	args := m.Called(backendURL)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Client), args.Error(1)
}

// Implement other required repository methods with empty implementations
func (m *MockRepository) SetClient(client *models.Client) error                           { return nil }
func (m *MockRepository) GetUserByID(id int64) (*models.User, error)                      { return nil, nil }
func (m *MockRepository) LoginUserPrecheck(username string) (string, error)               { return "", nil }
func (m *MockRepository) GetUser(username string) (*models.User, error)                   { return nil, nil }
func (m *MockRepository) SetTTL(key string, value []byte, expiration time.Duration) error { return nil }
func (m *MockRepository) GetTTL(key string) ([]byte, error)                               { return nil, nil }
func (m *MockRepository) GetUserMetadata(userID int64, key string) (*models.UserMetadata, error) {
	return &models.UserMetadata{}, nil
}
func TestService_VerifyToken(t *testing.T) {
	// Set up JWT secret for testing
	os.Setenv("JWT_SECRET_KEY", "test-secret-key")
	defer os.Unsetenv("JWT_SECRET_KEY")

	// Generate test tokens first
	validToken, err := utilities.GenerateStandardToken("test-secret-key")
	if err != nil {
		t.Fatalf("Failed to generate valid test token: %v", err)
	}

	tests := []struct {
		name        string
		token       string
		wantIsValid bool
		wantErr     bool
		errMsg      string
	}{
		{
			name:        "valid token",
			token:       validToken,
			wantIsValid: true,
			wantErr:     false,
		},
		{
			name:        "invalid token",
			token:       "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.l7W7V46W3NlY9tY80y660y0z022b2960201920200222",
			wantIsValid: false,
			wantErr:     true,
			errMsg:      "signature is invalid",
		},
	}

	// Initialize service with mock repository
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isValid, err := service.VerifyToken(tt.token)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.wantIsValid, isValid)
		})
	}
}

func TestService_CheckClient(t *testing.T) {
	tests := []struct {
		name       string
		backendURL string
		wantClient *models.Client
		wantErr    bool
		errMsg     string
		mockSetup  func(*MockRepository)
	}{
		{
			name:       "client found",
			backendURL: "http://valid-client.com",
			wantClient: &models.Client{
				ID:         "test-client",
				Secret:     "test-secret",
				Name:       "Test Client",
				BackendURI: "http://valid-client.com",
			},
			wantErr: false,
			mockSetup: func(m *MockRepository) {
				m.On("GetClientByURL", "http://valid-client.com").
					Return(&models.Client{
						ID:         "test-client",
						Secret:     "test-secret",
						Name:       "Test Client",
						BackendURI: "http://valid-client.com",
					}, nil)
			},
		},
		{
			name:       "client not found",
			backendURL: "http://invalid-client.com",
			wantClient: nil,
			wantErr:    true,
			errMsg:     "client not found",
			mockSetup: func(m *MockRepository) {
				m.On("GetClientByURL", "http://invalid-client.com").
					Return(nil, errors.New("client not found"))
			},
		},
		{
			name:       "repository error",
			backendURL: "http://error-client.com",
			wantClient: nil,
			wantErr:    true,
			errMsg:     "could not get client",
			mockSetup: func(m *MockRepository) {
				m.On("GetClientByURL", "http://error-client.com").
					Return(nil, errors.New("database error"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			if tt.mockSetup != nil {
				tt.mockSetup(mockRepo)
			}

			service := NewService(mockRepo)
			client, err := service.CheckClient(tt.backendURL)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.wantClient, client)
			mockRepo.AssertExpectations(t)
		})
	}
}
