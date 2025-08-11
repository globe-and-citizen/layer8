package service

import (
	"errors"
	"fmt"
	"globe-and-citizen/layer8/server/constants"
	"globe-and-citizen/layer8/server/models"
	"os"
	"strings"
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

const clientID = "clientID"
const certificate = "certificate"
const clientSecret = "client_secret"
const userID int64 = 10
const country = "some_country"
const displayName = "some_display_name"

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
	returnValues := m.Called(userID, key)
	if returnValues.Get(0) == nil {
		return nil, returnValues.Error(1)
	}

	return returnValues.Get(0).(*models.UserMetadata), returnValues.Error(1)
}

func (m *MockRepository) SaveX509Certificate(clientID string, certificate string) error {
	args := m.Called(clientID, certificate)
	return args.Error(0)
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

func TestSaveX509Certificate_RepositoryFailedToSaveCertificate(t *testing.T) {
	mockRepo := &MockRepository{}
	mockRepo.On(
		"SaveX509Certificate", clientID, certificate,
	).Return(
		fmt.Errorf("repository error"),
	)

	service := NewService(mockRepo)

	err := service.SaveX509Certificate(clientID, certificate)

	assert.NotNil(t, err)
}

func TestSaveX509Certificate_Success(t *testing.T) {
	mockRepo := &MockRepository{}
	mockRepo.On(
		"SaveX509Certificate", clientID, certificate,
	).Return(nil)

	service := NewService(mockRepo)

	err := service.SaveX509Certificate(clientID, certificate)

	assert.Nil(t, err)
}

func TestAuthenticateClient_FailedToGetClient(t *testing.T) {
	mockRepo := &MockRepository{}
	mockRepo.On("GetClient", clientID).Return(
		&models.Client{}, fmt.Errorf("failed to get client"),
	)
	service := NewService(mockRepo)

	err := service.AuthenticateClient(clientID, clientSecret)

	assert.NotNil(t, err)
}

func TestAuthenticateClient_ClientSecretDoesNotMatch(t *testing.T) {
	mockRepo := &MockRepository{}
	mockRepo.On("GetClient", clientID).Return(
		&models.Client{
			ID:     clientID,
			Secret: "other_secret",
		}, nil,
	)
	service := NewService(mockRepo)

	err := service.AuthenticateClient(clientID, clientSecret)

	assert.NotNil(t, err)
	assert.Equal(t, "failed to authenticate client: provided secret value is invalid", err.Error())
}

func TestAuthenticateClient_ClientAuthenticatedSuccessfully(t *testing.T) {
	mockRepo := &MockRepository{}
	mockRepo.On("GetClient", clientID).Return(
		&models.Client{
			ID:     clientID,
			Secret: clientSecret,
		}, nil,
	)
	service := NewService(mockRepo)

	err := service.AuthenticateClient(clientID, clientSecret)

	assert.Nil(t, err)
}

func TestGenerateAccessToken_Success(t *testing.T) {
	mockRepo := &MockRepository{}
	service := NewService(mockRepo)

	accessToken, err := service.GenerateAccessToken(
		&utilities.AuthCodeClaims{UserID: userID, ClientID: clientID},
		clientID,
		clientSecret)
	assert.Nil(t, err)

	claims, err := service.ValidateAccessToken(clientSecret, accessToken)
	assert.Nil(t, err)

	assert.Equal(t, clientID, claims.Subject)
	assert.Equal(t, "Globe and Citizen", claims.Issuer)
}

func TestGetZkUserMetadata_NoScopesProvided(t *testing.T) {
	mockRepo := &MockRepository{}
	service := NewService(mockRepo)

	_, err := service.GetZkUserMetadata("", userID)

	assert.NotNil(t, err)
}

func TestGetZkUserMetadata_FailedToGetCountryMetadata(t *testing.T) {
	scopes := "country"

	mockRepo := &MockRepository{}
	mockRepo.On(
		"GetUserMetadata", userID, constants.USER_COUNTRY_METADATA_KEY,
	).Return(nil, fmt.Errorf("error"))
	service := NewService(mockRepo)

	_, err := service.GetZkUserMetadata(scopes, userID)

	assert.NotNil(t, err)
	assert.True(t, strings.HasPrefix(err.Error(), "failed to get country metadata:"))
}

func TestGetZkUserMetadata_FailedToGetEmailZkMetadata(t *testing.T) {
	scopes := "country,email_verified"
	mockRepo := &MockRepository{}

	mockRepo.On(
		"GetUserMetadata", userID, constants.USER_COUNTRY_METADATA_KEY,
	).Return(&models.UserMetadata{UserID: userID}, nil)

	mockRepo.On(
		"GetUserMetadata", userID, constants.USER_EMAIL_VERIFIED_METADATA_KEY,
	).Return(nil, fmt.Errorf("error"))

	service := NewService(mockRepo)

	_, err := service.GetZkUserMetadata(scopes, userID)

	assert.NotNil(t, err)
	assert.True(t, strings.HasPrefix(err.Error(), "failed to get email metadata:"))
}

func TestGetZkUserMetadata_FailedToGetDisplayNameMetadata(t *testing.T) {
	scopes := "country,email_verified,display_name"
	mockRepo := &MockRepository{}

	mockRepo.On(
		"GetUserMetadata", userID, constants.USER_COUNTRY_METADATA_KEY,
	).Return(
		&models.UserMetadata{
			UserID: userID,
			Key:    constants.USER_COUNTRY_METADATA_KEY,
			Value:  country,
		},
		nil,
	)

	mockRepo.On(
		"GetUserMetadata", userID, constants.USER_EMAIL_VERIFIED_METADATA_KEY,
	).Return(
		&models.UserMetadata{
			UserID: userID,
			Key:    constants.USER_EMAIL_VERIFIED_METADATA_KEY,
			Value:  "true",
		},
		nil,
	)

	mockRepo.On(
		"GetUserMetadata", userID, constants.USER_DISPLAY_NAME_METADATA_KEY,
	).Return(nil, fmt.Errorf("error"))

	service := NewService(mockRepo)

	_, err := service.GetZkUserMetadata(scopes, userID)

	assert.NotNil(t, err)
	assert.True(t, strings.HasPrefix(err.Error(), "failed to get display name metadata:"))
}

func TestGetZkUserMetadata_Success(t *testing.T) {
	scopes := "country,email_verified,display_name,color"
	mockRepo := &MockRepository{}

	mockRepo.On(
		"GetUserMetadata", userID, constants.USER_COUNTRY_METADATA_KEY,
	).Return(
		&models.UserMetadata{
			UserID: userID,
			Key:    constants.USER_COUNTRY_METADATA_KEY,
			Value:  country,
		},
		nil,
	)

	mockRepo.On(
		"GetUserMetadata", userID, constants.USER_EMAIL_VERIFIED_METADATA_KEY,
	).Return(
		&models.UserMetadata{
			UserID: userID,
			Key:    constants.USER_EMAIL_VERIFIED_METADATA_KEY,
			Value:  "true",
		},
		nil,
	)

	mockRepo.On(
		"GetUserMetadata", userID, constants.USER_DISPLAY_NAME_METADATA_KEY,
	).Return(
		&models.UserMetadata{
			UserID: userID,
			Key:    constants.USER_DISPLAY_NAME_METADATA_KEY,
			Value:  displayName,
		},
		nil,
	)

	service := NewService(mockRepo)

	zkMetadata, err := service.GetZkUserMetadata(scopes, userID)

	assert.Nil(t, err)
	assert.Equal(t, country, zkMetadata.Country)
	assert.Equal(t, displayName, zkMetadata.DisplayName)
	assert.Equal(t, "red", zkMetadata.Color)
	assert.True(t, zkMetadata.IsEmailVerified)
}
