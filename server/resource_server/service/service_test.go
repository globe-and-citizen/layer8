package service_test

import (
	"fmt"
	serverModels "globe-and-citizen/layer8/server/models"
	"globe-and-citizen/layer8/server/resource_server/dto"
	"globe-and-citizen/layer8/server/resource_server/emails/verification"
	"globe-and-citizen/layer8/server/resource_server/models"
	"globe-and-citizen/layer8/server/resource_server/service"
	"globe-and-citizen/layer8/server/resource_server/utils/mocks"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const userId = 1
const adminEmail = "admin@email.com"
const username = "user"
const userEmail = "user@email.com"
const verificationCode = "123456"
const emailProof = "proof"
const verificationCodeValidityDuration = 2 * time.Minute

var timestamp = time.Date(2024, time.May, 24, 14, 0, 0, 0, time.UTC)
var timestampPlusTwoSeconds = timestamp.Add(2 * time.Second)
var now = func() time.Time {
	return timestamp
}
var mockCodeGenerator = &mocks.MockCodeGenerator{
	VerificationCode: verificationCode,
}
var defaultMockSenderService = &mocks.MockEmailSenderService{
	SendEmailFunc: func(email *models.Email) error {
		return nil
	},
}

type mockRepository struct {
	findUser                     func(userId uint) (models.User, error)
	saveEmailVerificationData    func(data models.EmailVerificationData) error
	getEmailVerificationData     func(userId uint) (models.EmailVerificationData, error)
	deleteEmailVerificationData  func(userId uint) error
	saveProofOfEmailVerification func(userID uint, verificationCode string, proof string) error
	setUserEmailVerified         func(userID uint) error
}

func (m *mockRepository) FindUser(userId uint) (models.User, error) {
	return m.findUser(userId)
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

func (m *mockRepository) SaveProofOfEmailVerification(
	userID uint, verificationCode string, proof string,
) error {
	return m.saveProofOfEmailVerification(userID, verificationCode, proof)
}

func (m *mockRepository) SaveEmailVerificationData(data models.EmailVerificationData) error {
	return m.saveEmailVerificationData(data)
}

func (m *mockRepository) GetEmailVerificationData(userId uint) (models.EmailVerificationData, error) {
	return m.getEmailVerificationData(userId)
}

func (m *mockRepository) UpdateDisplayName(userID uint, req dto.UpdateDisplayNameDTO) error {
	return nil
}

func (m *mockRepository) RegisterClient(req dto.RegisterClientDTO) error {
	return nil
}

func (m *mockRepository) IsBackendURIExists(backendURL string) (bool, error) {
    return true, nil
}

func (m *mockRepository) CheckBackendURI(backendURL string) (bool, error) {
    // Your mock implementation of CheckBackendURI here
    return true, nil
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
	mockService := service.NewService(mockRepo, &verification.EmailVerifier{})

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
	mockService := service.NewService(mockRepo, &verification.EmailVerifier{})

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
	mockService := service.NewService(mockRepo, &verification.EmailVerifier{})

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
	mockService := service.NewService(mockRepo, &verification.EmailVerifier{})

	// Call the ProfileUser method of the mock service
	userDetails, err := mockService.ProfileUser(1)
	if err != nil {
		t.Error("Expected nil, got", err)
	}

	// Use assert to check if the error is nil
	assert.Nil(t, err)
	assert.Equal(t, userDetails.Email, "test@gcitizen.com")
}

func TestUpdateDisplayName(t *testing.T) {
	// Create a new mock repository
	mockRepo := new(mockRepository)

	// Create a new service by passing the mock repository
	mockService := service.NewService(mockRepo, &verification.EmailVerifier{})

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
	mockService := service.NewService(mockRepo, &verification.EmailVerifier{})

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
	mockService := service.NewService(mockRepo, &verification.EmailVerifier{})

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

func TestCheckBackendURI(t *testing.T) {
    mockRepo := new(mockRepository)

    mockService := service.NewService(mockRepo)

    backendURL := "example.com"

    expectedResponse := true

    response, err := mockService.CheckBackendURI(backendURL)
    if err != nil {
        t.Error("Expected nil error, got", err)
    }

    assert.Nil(t, err)

    if response != expectedResponse {
        t.Errorf("Expected response: %v, got: %v", expectedResponse, response)
    }
}

func TestVerifyEmail_UserDoesNotExist(t *testing.T) {
	mockRepo := &mockRepository{
		findUser: func(userId uint) (models.User, error) {
			return models.User{}, fmt.Errorf("user %d does not exist", userId)
		},
	}
	emailVerifier := &verification.EmailVerifier{}

	currService := service.NewService(mockRepo, emailVerifier)
	e := currService.VerifyEmail(userId)

	assert.NotNil(t, e)
}

func TestVerifyEmail_UserExists_EmailFailedToBeSent(t *testing.T) {
	mockRepo := &mockRepository{
		findUser: func(userId uint) (models.User, error) {
			return models.User{
				ID:               userId,
				Username:         username,
				Email:            userEmail,
				VerificationCode: "",
			}, nil
		},
	}
	mockSenderService := &mocks.MockEmailSenderService{
		SendEmailFunc: func(email *models.Email) error {
			return fmt.Errorf("failed to send email")
		},
	}
	emailVerifier := verification.NewEmailVerifier(
		adminEmail,
		mockSenderService,
		mockCodeGenerator,
		verificationCodeValidityDuration,
		now,
	)
	currService := service.NewService(mockRepo, emailVerifier)

	e := currService.VerifyEmail(userId)

	assert.NotNil(t, e)
}

func TestVerifyEmail_UserExists_EmailSent_VerificationDataNotSaved(t *testing.T) {
	mockRepo := &mockRepository{
		findUser: func(userId uint) (models.User, error) {
			return models.User{
				ID:               userId,
				Username:         username,
				Email:            userEmail,
				VerificationCode: "",
			}, nil
		},
		saveEmailVerificationData: func(data models.EmailVerificationData) error {
			return fmt.Errorf("could not save the verification data")
		},
	}
	emailVerifier := verification.NewEmailVerifier(
		adminEmail,
		defaultMockSenderService,
		mockCodeGenerator,
		verificationCodeValidityDuration,
		now,
	)
	currService := service.NewService(mockRepo, emailVerifier)

	e := currService.VerifyEmail(userId)

	assert.NotNil(t, e)
}

func TestVerifyEmail_Success(t *testing.T) {
	mockRepo := &mockRepository{
		findUser: func(userId uint) (models.User, error) {
			return models.User{
				ID:               userId,
				Username:         username,
				Email:            userEmail,
				VerificationCode: "",
			}, nil
		},
		saveEmailVerificationData: func(data models.EmailVerificationData) error {
			return nil
		},
	}
	emailVerifier := verification.NewEmailVerifier(
		adminEmail,
		defaultMockSenderService,
		mockCodeGenerator,
		verificationCodeValidityDuration,
		now,
	)
	currService := service.NewService(mockRepo, emailVerifier)

	e := currService.VerifyEmail(userId)

	assert.Nil(t, e)
}

func TestCheckEmailVerificationCode_VerificationDataDoesNotExist(t *testing.T) {
	mockRepo := &mockRepository{
		getEmailVerificationData: func(userId uint) (models.EmailVerificationData, error) {
			return models.EmailVerificationData{},
				fmt.Errorf("could not get the verification data for user %d", userId)
		},
	}
	emailVerifier := &verification.EmailVerifier{}
	currService := service.NewService(mockRepo, emailVerifier)

	e := currService.CheckEmailVerificationCode(userId, verificationCode)

	assert.NotNil(t, e)
}

func TestCheckEmailVerificationCode_VerificationCodeMismatch(t *testing.T) {
	mockRepo := &mockRepository{
		getEmailVerificationData: func(userId uint) (models.EmailVerificationData, error) {
			return models.EmailVerificationData{
				UserId:           userId,
				VerificationCode: verificationCode,
				ExpiresAt:        timestampPlusTwoSeconds,
			}, nil
		},
	}
	emailVerifier := verification.NewEmailVerifier(
		adminEmail,
		defaultMockSenderService,
		mockCodeGenerator,
		verificationCodeValidityDuration,
		now,
	)
	currService := service.NewService(mockRepo, emailVerifier)

	e := currService.CheckEmailVerificationCode(userId, "567890")

	assert.NotNil(t, e)
}

func TestCheckEmailVerificationCode_VerificationCodeIsExpired(t *testing.T) {
	mockRepo := &mockRepository{
		getEmailVerificationData: func(userId uint) (models.EmailVerificationData, error) {
			return models.EmailVerificationData{
				UserId:           userId,
				VerificationCode: verificationCode,
				ExpiresAt:        timestamp,
			}, nil
		},
	}
	emailVerifier := verification.NewEmailVerifier(
		adminEmail,
		defaultMockSenderService,
		mockCodeGenerator,
		verificationCodeValidityDuration,
		func() time.Time {
			return timestampPlusTwoSeconds
		},
	)
	currService := service.NewService(mockRepo, emailVerifier)

	e := currService.CheckEmailVerificationCode(userId, verificationCode)

	assert.NotNil(t, e)
}

func TestCheckEmailVerificationCode_Success(t *testing.T) {
	mockRepo := &mockRepository{
		getEmailVerificationData: func(userId uint) (models.EmailVerificationData, error) {
			return models.EmailVerificationData{
				UserId:           userId,
				VerificationCode: verificationCode,
				ExpiresAt:        timestampPlusTwoSeconds,
			}, nil
		},
	}
	emailVerifier := verification.NewEmailVerifier(
		adminEmail,
		defaultMockSenderService,
		mockCodeGenerator,
		verificationCodeValidityDuration,
		now,
	)
	currService := service.NewService(mockRepo, emailVerifier)

	e := currService.CheckEmailVerificationCode(userId, verificationCode)

	assert.Nil(t, e)
}

func TestSaveProofOfEmailVerification_ProofFailedToBeSaved(t *testing.T) {
	mockRepo := &mockRepository{
		saveProofOfEmailVerification: func(userID uint, verificationCode string, proof string) error {
			return fmt.Errorf("could not save proof of verification for user %d", userID)
		},
	}
	emailVerifier := verification.NewEmailVerifier(
		adminEmail,
		defaultMockSenderService,
		mockCodeGenerator,
		verificationCodeValidityDuration,
		now,
	)
	currService := service.NewService(mockRepo, emailVerifier)

	e := currService.SaveProofOfEmailVerification(userId, verificationCode, emailProof)

	assert.NotNil(t, e)
}

func TestSaveProofOfEmailVerification_Success(t *testing.T) {
	mockRepo := &mockRepository{
		saveProofOfEmailVerification: func(userID uint, verificationCode string, proof string) error {
			return nil
		},
	}
	emailVerifier := verification.NewEmailVerifier(
		adminEmail,
		defaultMockSenderService,
		mockCodeGenerator,
		verificationCodeValidityDuration,
		now,
	)
	currService := service.NewService(mockRepo, emailVerifier)

	e := currService.SaveProofOfEmailVerification(userId, verificationCode, emailProof)

	assert.Nil(t, e)
}
