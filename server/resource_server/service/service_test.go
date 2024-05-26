package service_test

import (
	"fmt"
	serverModels "globe-and-citizen/layer8/server/models"
	"globe-and-citizen/layer8/server/resource_server/dto"
	"globe-and-citizen/layer8/server/resource_server/models"
	"globe-and-citizen/layer8/server/resource_server/service"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type mockRepository struct {
}

func (m *mockRepository) RegisterUser(req dto.RegisterUserDTO) error {
	return nil
}

func (m *mockRepository) LoginPreCheckUser(req dto.LoginPrecheckDTO) (string, string, error) {
	return "test_user", "ThisIsARandomSalt123!@#", nil
}

func (m *mockRepository) LoginUser(req dto.LoginUserDTO) (models.User, error) {
	return models.User{
		Email:     "test@gcitizen.com",
		Username:  "test_user",
		FirstName: "Test",
		LastName:  "User",
		Password:  "34efcb97e704298f3d64159ee858c6c1826755b37523cfac8a79c2130ea7b16f",
		Salt:      "312c4a2c46405ba4f70f7be070f4d4f7cdede09d4b218bf77c01f9706d7505c9",
	}, nil
}

func (m *mockRepository) ProfileUser(userID uint) (models.User, []models.UserMetadata, error) {
	if userID == 1 {
		return models.User{
				Email:     "test@gcitizen.com",
				Username:  "test_user",
				FirstName: "Test",
				LastName:  "User",
				Password:  "34efcb97e704298f3d64159ee858c6c1826755b37523cfac8a79c2130ea7b16f",
				Salt:      "312c4a2c46405ba4f70f7be070f4d4f7cdede09d4b218bf77c01f9706d7505c9",
			}, []models.UserMetadata{
				{
					UserID: 1,
					Key:    "email_verified",
					Value:  "true",
				},
				{
					UserID: 1,
					Key:    "display_name",
					Value:  "user",
				},
				{
					UserID: 1,
					Key:    "country",
					Value:  "Unknown",
				},
			}, nil
	}
	return models.User{}, []models.UserMetadata{}, fmt.Errorf("User not found")
}

func (m *mockRepository) VerifyEmail(userID uint) error {
	return nil
}

func (m *mockRepository) UpdateDisplayName(userID uint, req dto.UpdateDisplayNameDTO) error {
	return nil
}

func (m *mockRepository) RegisterClient(req dto.RegisterClientDTO) error {
	return nil
}

func (m *mockRepository) GetClientData(clientName string) (models.Client, error) {
	if clientName == "testclient" {
		return models.Client{
			ID:          "1",
			Secret:      "testsecret",
			Name:        "testclient",
			RedirectURI: "https://gcitizen.com/callback",
		}, nil
	}
	return models.Client{}, fmt.Errorf("Client not found")
}

func (m *mockRepository) LoginUserPrecheck(username string) (string, error) {
	return "", nil
}

func (m *mockRepository) GetUser(username string) (*serverModels.User, error) {
	return &serverModels.User{}, nil
}

func (m *mockRepository) GetUserByID(id int64) (*serverModels.User, error) {
	return &serverModels.User{}, nil
}

func (m *mockRepository) GetUserMetadata(userID int64, key string) (*serverModels.UserMetadata, error) {
	return &serverModels.UserMetadata{}, nil
}

func (m *mockRepository) SetClient(client *serverModels.Client) error {
	return nil
}

func (m *mockRepository) GetClient(clientName string) (*serverModels.Client, error) {
	return &serverModels.Client{}, nil
}

func (m *mockRepository) SetTTL(key string, value []byte, time time.Duration) error {
	return nil
}

func (m *mockRepository) GetTTL(key string) ([]byte, error) {
	return []byte{}, nil
}

func (m *mockRepository) LoginClient(req dto.LoginClientDTO) (models.Client, error) {
	return models.Client{}, nil
}

func (m *mockRepository) LoginPreCheckClient(req dto.LoginPrecheckDTO) (string, string, error) {
	return "", "", nil
}

func (m *mockRepository) ProfileClient(userID string) (models.Client, error) {
	return models.Client{}, nil
}

func (m *mockRepository) GetClientDataByBackendURL(backendURL string) (models.Client, error) {
	return models.Client{}, nil
}

func TestRegisterUser(t *testing.T) {
	// Create a new mock repository
	mockRepo := new(mockRepository)

	// Create a new service by passing the mock repository
	mockService := service.NewService(mockRepo)

	// Create a new mock request
	req := dto.RegisterUserDTO{
		Email:       "test@gcitizen.com",
		Username:    "test_user",
		FirstName:   "Test",
		LastName:    "User",
		DisplayName: "user",
		Country:     "Unknown",
		Password:    "12345",
	}

	// Call the RegisterUser method of the mock service
	err := mockService.RegisterUser(req)
	if err != nil {
		t.Error("Expected nil, got", err)
	}

	// Use assert to check if the error is nil
	assert.Nil(t, err)
}

func TestLoginPreCheckUser(t *testing.T) {
	// Create a new mock repository
	mockRepo := new(mockRepository)

	// Create a new service by passing the mock repository
	mockService := service.NewService(mockRepo)

	// Create a new mock request
	req := dto.LoginPrecheckDTO{
		Username: "test_user",
	}

	// Call the LoginPreCheckUser method of the mock service
	loginPrecheckResp, err := mockService.LoginPreCheckUser(req)
	if err != nil {
		t.Error("Expected nil, got", err)
	}

	// Use assert to check if the error is nil
	assert.Nil(t, err)
	assert.Equal(t, loginPrecheckResp.Username, "test_user")
	assert.Equal(t, loginPrecheckResp.Salt, "ThisIsARandomSalt123!@#")
}

func TestLoginUser(t *testing.T) {
	// Create a new mock repository
	mockRepo := new(mockRepository)

	// Create a new service by passing the mock repository
	mockService := service.NewService(mockRepo)

	// Create a new mock request
	req := dto.LoginUserDTO{
		Username: "test_user",
		Password: "12345",
		Salt:     "312c4a2c46405ba4f70f7be070f4d4f7cdede09d4b218bf77c01f9706d7505c9",
	}

	// Call the LoginUser method of the mock service
	_, err := mockService.LoginUser(req)
	if err != nil {
		t.Error("Expected nil, got", err)
	}

	// Use assert to check if the error is nil
	assert.Nil(t, err)
}

func TestProfileUser(t *testing.T) {
	// Create a new mock repository
	mockRepo := new(mockRepository)

	// Create a new service by passing the mock repository
	mockService := service.NewService(mockRepo)

	// Call the ProfileUser method of the mock service
	userDetails, err := mockService.ProfileUser(1)
	if err != nil {
		t.Error("Expected nil, got", err)
	}

	// Use assert to check if the error is nil
	assert.Nil(t, err)
	assert.Equal(t, userDetails.Email, "test@gcitizen.com")
}

func TestVerifyEmail(t *testing.T) {
	// Create a new mock repository
	mockRepo := new(mockRepository)

	// Create a new service by passing the mock repository
	mockService := service.NewService(mockRepo)

	// Call the VerifyEmail method of the mock service
	err := mockService.VerifyEmail(1)
	if err != nil {
		t.Error("Expected nil, got", err)
	}

	// Use assert to check if the error is nil
	assert.Nil(t, err)
}

func TestUpdateDisplayName(t *testing.T) {
	// Create a new mock repository
	mockRepo := new(mockRepository)

	// Create a new service by passing the mock repository
	mockService := service.NewService(mockRepo)

	// Create a new mock request
	req := dto.UpdateDisplayNameDTO{
		DisplayName: "user",
	}

	// Call the UpdateDisplayName method of the mock service
	err := mockService.UpdateDisplayName(1, req)
	if err != nil {
		t.Error("Expected nil, got", err)
	}

	// Use assert to check if the error is nil
	assert.Nil(t, err)
}

func TestRegisterClient(t *testing.T) {
	// Create a new mock repository
	mockRepo := new(mockRepository)

	// Create a new service by passing the mock repository
	mockService := service.NewService(mockRepo)

	// Create a new mock request
	req := dto.RegisterClientDTO{
		Name:        "testclient",
		RedirectURI: "https://gcitizen.com/callback",
		BackendURI:  "https://gcitizen_backend.com/callback",
		Username:    "test_user",
		Password:    "12345",
	}

	// Call the RegisterClient method of the mock service
	err := mockService.RegisterClient(req)
	if err != nil {
		t.Error("Expected nil, got", err)
	}

	// Use assert to check if the error is nil
	assert.Nil(t, err)
}

func TestGetClientData(t *testing.T) {
	// Create a new mock repository
	mockRepo := new(mockRepository)

	// Create a new service by passing the mock repository
	mockService := service.NewService(mockRepo)

	// Call the GetClientData method of the mock service
	clientData, err := mockService.GetClientData("testclient")
	if err != nil {
		t.Error("Expected nil, got", err)
	}

	// Use assert to check if the error is nil
	assert.Nil(t, err)
	assert.Equal(t, clientData.Secret, "testsecret")
	assert.Equal(t, clientData.RedirectURI, "https://gcitizen.com/callback")
}

func (m *mockRepository) CheckBackendURI(backendURL string) (bool, error) {
	return true, nil
}