package service_test

import (
	"crypto/rand"
	"errors"
	"fmt"
	serverModels "globe-and-citizen/layer8/server/models"
	"globe-and-citizen/layer8/server/resource_server/dto"
	"globe-and-citizen/layer8/server/resource_server/emails/verification"
	"globe-and-citizen/layer8/server/resource_server/emails/verification/zk"
	"globe-and-citizen/layer8/server/resource_server/models"
	"globe-and-citizen/layer8/server/resource_server/service"
	"globe-and-citizen/layer8/server/resource_server/utils"
	"globe-and-citizen/layer8/server/resource_server/utils/mocks"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/crypto"

	"github.com/stretchr/testify/assert"
)

const userId uint = 1
const adminEmail = "admin@email.com"
const username = "test_user"
const password = "1234"
const userEmail = "user@email.com"
const firstName = "first_name"
const lastName = "second_name"
const displayName = "display_name"
const country = "country"
const verificationCode = "123456"
const publicKey = "0xaaaaa"
const cNonce = "1a7fa8e9dc2a68049358a08349cdde50"
const nonce = "1a7fa8e9dc2a68049358a08349cdde50d6d6f9c5e632f599881d99fe9c62b362"
const clientProof = "96260616beaabaa6d168b9ce7e15b127b6bcbcb88fd0b5de1e6162b948881155"
const storedKey = "222f705167604c99f81c7c6acfa974706fa9dd0b445a8bc34fb5accc9b032558"
const serverKey = "f6e938506893b10799038c2b4225f7e8e72f01c7ab25c3da905803ee93ec5536"
const salt = "c8f720569d7a50d4c812431cf8a242fd608f3a5b4610659b92f6aa553fbe68e0"
const testServerSignature = "138c0fecd0326896fe21398137c1aa0b9866242abc02bd3e9a5a0016013ef5f0"
const iterationCount = 4096
const clientId = "123abc"
const clientSecret = "456def"
const clientName = "testclient"
const redirectUri = "redirect_uri"
const backendUri = "backend_uri"

const zkKeyPairId uint = 2

const verificationCodeValidityDuration = 2 * time.Minute

var emailProof = []byte("proof")
var hashedPassword = utils.SaltAndHashPassword(password, salt)

var zkProof = []byte("zk_proof")
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
	saveProofOfEmailVerification func(userID uint, verificationCode string, proof []byte, zkKeyPairId uint) error
	registerUser                 func(req dto.RegisterUserDTO, hashedPassword string, salt string) error
	registerUserPrecheck         func(req dto.RegisterUserPrecheckDTO, salt string, iterCount int) error
	registerClientPrecheck       func(req dto.RegisterClientPrecheckDTO, salt string, iterCount int) error
	registerUserv2               func(req dto.RegisterUserDTOv2) error
	registerClient               func(client models.Client) error
	registerClientv2             func(req dto.RegisterClientDTOv2, id string, secret string) error
	getUserForUsername           func(username string) (models.User, error)
	profileClient                func(username string) (models.Client, error)
	updateUserPassword           func(username string, password string) error
	updateUserPasswordV2         func(username string, storedKey string, serverKey string) error
}

func (m *mockRepository) FindUser(userId uint) (models.User, error) {
	return m.findUser(userId)
}

func (m *mockRepository) RegisterUser(req dto.RegisterUserDTO, hashedPassword string, salt string) error {
	return m.registerUser(req, hashedPassword, salt)
}

func (m *mockRepository) RegisterPrecheckUser(req dto.RegisterUserPrecheckDTO, salt string, iterCount int) error {
	if m.registerUserPrecheck != nil {
		return m.registerUserPrecheck(req, salt, iterCount)
	}
	return nil
}

func (m *mockRepository) RegisterUserv2(req dto.RegisterUserDTOv2) error {
	return m.registerUserv2(req)
}

func (m *mockRepository) LoginPreCheckUser(req dto.LoginPrecheckDTO) (string, string, error) {
	return username, salt, nil
}

func (m *mockRepository) LoginUser(req dto.LoginUserDTO) (models.User, error) {
	return models.User{
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
	userID uint, verificationCode string, proof []byte, zkKeyPairId uint,
) error {
	return m.saveProofOfEmailVerification(userID, verificationCode, proof, zkKeyPairId)
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

func (m *mockRepository) RegisterClient(client models.Client) error {
	return m.registerClient(client)
}

func (m *mockRepository) RegisterPrecheckClient(req dto.RegisterClientPrecheckDTO, salt string, iterCount int) error {
	return m.registerClientPrecheck(req, salt, iterCount)
}

func (m *mockRepository) RegisterClientv2(req dto.RegisterClientDTOv2, id string, secret string) error {
	return m.registerClientv2(req, id, secret)
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

func (m *mockRepository) ProfileClient(username string) (models.Client, error) {
	return m.profileClient(username)
}

func (m *mockRepository) GetClientDataByBackendURL(backendURL string) (models.Client, error) {
	return models.Client{}, nil
}

func (m *mockRepository) SaveZkSnarksKeyPair(keyPair models.ZkSnarksKeyPair) (uint, error) {
	return 0, nil
}

func (m *mockRepository) GetLatestZkSnarksKeys() (models.ZkSnarksKeyPair, error) {
	return models.ZkSnarksKeyPair{}, nil
}

func (m *mockRepository) GetUserForUsername(username string) (models.User, error) {
	return m.getUserForUsername(username)
}

func (m *mockRepository) UpdateUserPassword(username string, password string) error {
	return m.updateUserPassword(username, password)
}

func (m *mockRepository) UpdateUserPasswordV2(username string, storedKey string, serverKey string) error {
	return m.updateUserPasswordV2(username, storedKey, serverKey)
}

func (m *mockRepository) AddClientTrafficUsage(string, int, time.Time) error {
	return nil
}

func (m *mockRepository) CreateClientTrafficStatisticsEntry(string, int) error {
	return nil
}

func (m *mockRepository) GetAllClientStatistics() ([]models.ClientTrafficStatistics, error) {
	return nil, nil
}

func (m *mockRepository) GetClientTrafficStatistics(string) (*models.ClientTrafficStatistics, error) {
	return &models.ClientTrafficStatistics{}, nil
}

func (m *mockRepository) PayClientTrafficUsage(string, int) error {
	return nil
}

func TestRegisterUser_RepositoryFailedToStoreUserData(t *testing.T) {
	mockRepo := &mockRepository{
		registerUser: func(req dto.RegisterUserDTO, hashedPassword string, salt string) error {
			assert.Equal(t, username, req.Username)
			assert.Equal(t, firstName, req.FirstName)
			assert.Equal(t, lastName, req.LastName)
			assert.Equal(t, displayName, req.DisplayName)
			assert.Equal(t, country, req.Country)

			return fmt.Errorf("failed to store a user")
		},
	}
	currService := service.NewService(mockRepo, &verification.EmailVerifier{}, &mocks.MockProofGenerator{})

	err := currService.RegisterUser(
		dto.RegisterUserDTO{
			Username:    username,
			FirstName:   firstName,
			LastName:    lastName,
			DisplayName: displayName,
			Country:     country,
			Password:    password,
		},
	)

	assert.NotNil(t, err)
}

func TestRegisterUser_Success(t *testing.T) {
	mockRepo := &mockRepository{
		registerUser: func(req dto.RegisterUserDTO, hashedPassword string, salt string) error {
			assert.Equal(t, username, req.Username)
			assert.Equal(t, firstName, req.FirstName)
			assert.Equal(t, lastName, req.LastName)
			assert.Equal(t, displayName, req.DisplayName)
			assert.Equal(t, country, req.Country)

			return nil
		},
	}
	currService := service.NewService(mockRepo, &verification.EmailVerifier{}, &mocks.MockProofGenerator{})

	err := currService.RegisterUser(
		dto.RegisterUserDTO{
			Username:    username,
			FirstName:   firstName,
			LastName:    lastName,
			DisplayName: displayName,
			Country:     country,
			Password:    password,
		},
	)

	assert.Nil(t, err)
}

func TestRegisterClient_RepositoryFailedToStoreClientData(t *testing.T) {
	mockRepo := &mockRepository{
		registerClient: func(client models.Client) error {
			assert.Equal(t, firstName, client.Name)
			assert.Equal(t, username, client.Username)
			assert.Equal(t, redirectUri, client.RedirectURI)
			assert.Equal(t, backendUri, client.BackendURI)

			return fmt.Errorf("failed to store a client")
		},
	}
	currService := service.NewService(mockRepo, &verification.EmailVerifier{}, &mocks.MockProofGenerator{})

	err := currService.RegisterClient(
		dto.RegisterClientDTO{
			Name:        firstName,
			RedirectURI: redirectUri,
			BackendURI:  backendUri,
			Username:    username,
			Password:    password,
		},
	)

	assert.NotNil(t, err)
}

func TestRegisterClient_Success(t *testing.T) {
	mockRepo := &mockRepository{
		registerClient: func(client models.Client) error {
			assert.Equal(t, firstName, client.Name)
			assert.Equal(t, username, client.Username)
			assert.Equal(t, redirectUri, client.RedirectURI)
			assert.Equal(t, backendUri, client.BackendURI)

			return nil
		},
	}
	currService := service.NewService(mockRepo, &verification.EmailVerifier{}, &mocks.MockProofGenerator{})

	os.Setenv("CLIENT_TRAFFIC_RATE_PER_BYTE", "5")
	err := currService.RegisterClient(
		dto.RegisterClientDTO{
			Name:        firstName,
			RedirectURI: redirectUri,
			BackendURI:  backendUri,
			Username:    username,
			Password:    password,
		},
	)

	assert.Nil(t, err)
}

func TestLoginPreCheckUser(t *testing.T) {
	// Create a new mock repository
	mockRepo := new(mockRepository)

	// Create a new service by passing the mock repository
	mockService := service.NewService(mockRepo, &verification.EmailVerifier{}, &mocks.MockProofGenerator{})

	// Create a new mock request
	req := dto.LoginPrecheckDTO{
		Username: username,
	}

	// Call the LoginPreCheckUser method of the mock service
	loginPrecheckResp, err := mockService.LoginPreCheckUser(req)
	if err != nil {
		t.Error("Expected nil, got", err)
	}

	// Use assert to check if the error is nil
	assert.Nil(t, err)
	assert.Equal(t, loginPrecheckResp.Username, username)
	assert.Equal(t, loginPrecheckResp.Salt, salt)
}

func TestLoginPreCheckUserV2_RepositoryError(t *testing.T) {
	mockRepo := &mockRepository{
		getUserForUsername: func(username string) (models.User, error) {
			return models.User{}, errors.New("repository error")
		},
	}
	currService := service.NewService(mockRepo, &verification.EmailVerifier{}, &mocks.MockProofGenerator{})

	req := dto.LoginPrecheckDTOv2{
		Username: username,
		CNonce:   cNonce,
	}

	loginPrecheckResp, err := currService.LoginPrecheckUserv2(req)

	assert.NotNil(t, err)
	assert.Equal(t, "repository error", err.Error())
	assert.Empty(t, loginPrecheckResp)
}

func TestLoginPreCheckUserV2_Success(t *testing.T) {
	mockRepo := &mockRepository{
		getUserForUsername: func(username string) (models.User, error) {
			return models.User{
				Username:       username,
				Salt:           salt,
				IterationCount: iterationCount,
			}, nil
		},
	}

	currService := service.NewService(mockRepo, &verification.EmailVerifier{}, &mocks.MockProofGenerator{})

	req := dto.LoginPrecheckDTOv2{
		Username: username,
		CNonce:   cNonce,
	}

	loginPrecheckResp, err := currService.LoginPrecheckUserv2(req)

	assert.Nil(t, err)
	assert.Equal(t, loginPrecheckResp.Salt, salt)
	assert.Equal(t, strings.HasPrefix(loginPrecheckResp.Nonce, cNonce), true)
	assert.Equal(t, loginPrecheckResp.IterCount, iterationCount)
}

func TestLoginUser(t *testing.T) {
	// Create a new mock repository
	mockRepo := new(mockRepository)

	// Create a new service by passing the mock repository
	mockService := service.NewService(mockRepo, &verification.EmailVerifier{}, &mocks.MockProofGenerator{})

	// Create a new mock request
	req := dto.LoginUserDTO{
		Username: username,
		Password: "12345",
		Salt:     salt,
	}

	// Call the LoginUser method of the mock service
	_, err := mockService.LoginUser(req)
	if err != nil {
		t.Error("Expected nil, got", err)
	}

	// Use assert to check if the error is nil
	assert.Nil(t, err)
}

func TestLoginUserV2_RepositoryError(t *testing.T) {
	mockRepo := &mockRepository{
		getUserForUsername: func(username string) (models.User, error) {
			return models.User{}, errors.New("repository error")
		},
	}
	currService := service.NewService(mockRepo, &verification.EmailVerifier{}, &mocks.MockProofGenerator{})

	req := dto.LoginUserDTOv2{
		Username:    username,
		CNonce:      cNonce,
		Nonce:       nonce,
		ClientProof: clientProof,
	}

	loginUserResp, err := currService.LoginUserv2(req)

	assert.NotNil(t, err)
	assert.Equal(t, "repository error", err.Error())
	assert.Empty(t, loginUserResp)
}

func TestLoginUserV2_DecodingStoredKeyError(t *testing.T) {
	mockRepo := &mockRepository{
		getUserForUsername: func(username string) (models.User, error) {
			return models.User{
				StoredKey: "TEST_STORED_KEY_FOR_DECODE_ERROR_!@#", // Invalid stored key
			}, nil
		},
	}
	currService := service.NewService(mockRepo, &verification.EmailVerifier{}, &mocks.MockProofGenerator{})

	req := dto.LoginUserDTOv2{
		Username:    username,
		CNonce:      cNonce,
		Nonce:       nonce,
		ClientProof: clientProof,
	}

	loginUserResp, err := currService.LoginUserv2(req)

	assert.NotNil(t, err)
	assert.Equal(t, "error decoding stored key: encoding/hex: invalid byte: U+0054 'T'", err.Error())
	assert.Empty(t, loginUserResp)
}

func TestLoginUserV2_DecodingClientProofError(t *testing.T) {
	mockRepo := &mockRepository{
		getUserForUsername: func(username string) (models.User, error) {
			return models.User{
				StoredKey:      storedKey,
				Salt:           salt,
				IterationCount: iterationCount,
			}, nil
		},
	}
	currService := service.NewService(mockRepo, &verification.EmailVerifier{}, &mocks.MockProofGenerator{})

	req := dto.LoginUserDTOv2{
		Username:    username,
		CNonce:      cNonce,
		Nonce:       nonce,
		ClientProof: "TEST_CLIENT_PROOF_FOR_DECODE_ERROR_!@#", // Invalid client proof
	}

	loginUserResp, err := currService.LoginUserv2(req)

	assert.NotNil(t, err)
	assert.Equal(t, "error decoding client proof: encoding/hex: invalid byte: U+0054 'T'", err.Error())
	assert.Empty(t, loginUserResp)
}

func TestLoginUserV2_XorOperationError(t *testing.T) {
	mockRepo := &mockRepository{
		getUserForUsername: func(username string) (models.User, error) {
			return models.User{
				StoredKey:      storedKey,
				Salt:           salt,
				IterationCount: iterationCount,
			}, nil
		},
	}
	currService := service.NewService(mockRepo, &verification.EmailVerifier{}, &mocks.MockProofGenerator{})

	req := dto.LoginUserDTOv2{
		Username:    username,
		CNonce:      cNonce,
		Nonce:       nonce,
		ClientProof: "", // Sending empty client proof to fail the XOR operation, since it requires equal length slices
	}

	loginUserResp, err := currService.LoginUserv2(req)

	assert.NotNil(t, err)
	assert.Equal(t, "error performing XOR operation: slices must have the same length", err.Error())
	assert.Empty(t, loginUserResp)
}

func TestLoginUserV2_AuthFailedKeyMismatchError(t *testing.T) {
	mockRepo := &mockRepository{
		getUserForUsername: func(username string) (models.User, error) {
			return models.User{
				StoredKey:      "", // Empty stored key
				Salt:           salt,
				IterationCount: iterationCount,
			}, nil
		},
	}
	currService := service.NewService(mockRepo, &verification.EmailVerifier{}, &mocks.MockProofGenerator{})

	req := dto.LoginUserDTOv2{
		Username:    username,
		CNonce:      cNonce,
		Nonce:       nonce,
		ClientProof: clientProof,
	}

	loginUserResp, err := currService.LoginUserv2(req)

	assert.NotNil(t, err)
	assert.Equal(t, "server failed to authenticate the user", err.Error())
	assert.Empty(t, loginUserResp)
}

func TestLoginUserV2_DecodingServerKeyError(t *testing.T) {
	mockRepo := &mockRepository{
		getUserForUsername: func(username string) (models.User, error) {
			return models.User{
				StoredKey:      storedKey,
				ServerKey:      "TEST_SERVER_KEY_FOR_DECODE_ERROR_!@#", // Invalid server key
				Salt:           salt,
				IterationCount: iterationCount,
			}, nil
		},
	}
	currService := service.NewService(mockRepo, &verification.EmailVerifier{}, &mocks.MockProofGenerator{})

	req := dto.LoginUserDTOv2{
		Username:    username,
		CNonce:      cNonce,
		Nonce:       nonce,
		ClientProof: clientProof,
	}

	loginUserResp, err := currService.LoginUserv2(req)

	assert.NotNil(t, err)
	assert.Equal(t, "error decoding server key: encoding/hex: invalid byte: U+0054 'T'", err.Error())
	assert.Empty(t, loginUserResp)
}

func TestLoginUserV2_Success(t *testing.T) {
	mockRepo := &mockRepository{
		getUserForUsername: func(username string) (models.User, error) {
			return models.User{
				StoredKey:      storedKey,
				ServerKey:      serverKey,
				Salt:           salt,
				IterationCount: iterationCount,
			}, nil
		},
	}
	currService := service.NewService(mockRepo, &verification.EmailVerifier{}, &mocks.MockProofGenerator{})

	req := dto.LoginUserDTOv2{
		Username:    username,
		CNonce:      cNonce,
		Nonce:       nonce,
		ClientProof: clientProof,
	}

	loginUserResp, err := currService.LoginUserv2(req)

	assert.Nil(t, err)
	assert.Equal(t, testServerSignature, loginUserResp.ServerSignature)
}

func TestProfileUser(t *testing.T) {
	// Create a new mock repository
	mockRepo := new(mockRepository)

	// Create a new service by passing the mock repository
	mockService := service.NewService(mockRepo, &verification.EmailVerifier{}, &mocks.MockProofGenerator{})

	// Call the ProfileUser method of the mock service
	userDetails, err := mockService.ProfileUser(1)
	if err != nil {
		t.Error("Expected nil, got", err)
	}

	// Use assert to check if the error is nil
	assert.Nil(t, err)
	assert.Equal(t, userDetails.Username, username)
}

func TestUpdateDisplayName(t *testing.T) {
	// Create a new mock repository
	mockRepo := new(mockRepository)

	// Create a new service by passing the mock repository
	mockService := service.NewService(mockRepo, &verification.EmailVerifier{}, &mocks.MockProofGenerator{})

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

func TestGetClientData(t *testing.T) {
	// Create a new mock repository
	mockRepo := new(mockRepository)

	// Create a new service by passing the mock repository
	mockService := service.NewService(mockRepo, &verification.EmailVerifier{}, &mocks.MockProofGenerator{})

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

	mockService := service.NewService(mockRepo, &verification.EmailVerifier{}, &zk.ProofProcessor{})

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

	currService := service.NewService(mockRepo, emailVerifier, &mocks.MockProofGenerator{})
	e := currService.VerifyEmail(userId, userEmail)

	assert.NotNil(t, e)
}

func TestVerifyEmail_UserExists_EmailFailedToBeSent(t *testing.T) {
	mockRepo := &mockRepository{
		findUser: func(userId uint) (models.User, error) {
			return models.User{
				ID:               userId,
				Username:         username,
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
	currService := service.NewService(mockRepo, emailVerifier, &mocks.MockProofGenerator{})

	e := currService.VerifyEmail(userId, userEmail)

	assert.NotNil(t, e)
}

func TestVerifyEmail_UserExists_EmailSent_VerificationDataNotSaved(t *testing.T) {
	mockRepo := &mockRepository{
		findUser: func(userId uint) (models.User, error) {
			return models.User{
				ID:               userId,
				Username:         username,
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
	currService := service.NewService(mockRepo, emailVerifier, &mocks.MockProofGenerator{})

	e := currService.VerifyEmail(userId, userEmail)

	assert.NotNil(t, e)
}

func TestVerifyEmail_Success(t *testing.T) {
	mockRepo := &mockRepository{
		findUser: func(userId uint) (models.User, error) {
			return models.User{
				ID:               userId,
				Username:         username,
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
	currService := service.NewService(mockRepo, emailVerifier, &mocks.MockProofGenerator{})

	e := currService.VerifyEmail(userId, userEmail)

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
	currService := service.NewService(mockRepo, emailVerifier, &mocks.MockProofGenerator{})

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
	currService := service.NewService(mockRepo, emailVerifier, &mocks.MockProofGenerator{})

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
	currService := service.NewService(mockRepo, emailVerifier, &mocks.MockProofGenerator{})

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
	currService := service.NewService(mockRepo, emailVerifier, &mocks.MockProofGenerator{})

	e := currService.CheckEmailVerificationCode(userId, verificationCode)

	assert.Nil(t, e)
}

func TestSaveProofOfEmailVerification_ProofFailedToBeSaved(t *testing.T) {
	mockRepo := &mockRepository{
		saveProofOfEmailVerification: func(userID uint, verificationCode string, proof []byte, zkKeyPairId uint) error {
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
	currService := service.NewService(mockRepo, emailVerifier, &mocks.MockProofGenerator{})

	e := currService.SaveProofOfEmailVerification(userId, verificationCode, emailProof, zkKeyPairId)

	assert.NotNil(t, e)
}

func TestSaveProofOfEmailVerification_Success(t *testing.T) {
	mockRepo := &mockRepository{
		saveProofOfEmailVerification: func(userID uint, verificationCode string, proof []byte, zkKeyPairId uint) error {
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
	currService := service.NewService(mockRepo, emailVerifier, &mocks.MockProofGenerator{})

	e := currService.SaveProofOfEmailVerification(userId, verificationCode, emailProof, zkKeyPairId)

	assert.Nil(t, e)
}

func TestFindUser_UserNotFound(t *testing.T) {
	mockRepo := &mockRepository{
		findUser: func(userId uint) (models.User, error) {
			return models.User{}, fmt.Errorf("user not found for id %d", userId)
		},
	}

	currService := service.NewService(mockRepo, &verification.EmailVerifier{}, &mocks.MockProofGenerator{})
	_, err := currService.FindUser(userId)

	assert.NotNil(t, err)
}

func TestFindUser_Success(t *testing.T) {
	mockRepo := &mockRepository{
		findUser: func(userId uint) (models.User, error) {
			return models.User{ID: userId}, nil
		},
	}

	currService := service.NewService(mockRepo, &verification.EmailVerifier{}, &mocks.MockProofGenerator{})
	user, err := currService.FindUser(userId)

	assert.Nil(t, err)
	assert.Equal(t, userId, user.ID)
}

func TestGenerateZkProofOfEmailVerification_FailedToGenerateZkProof(t *testing.T) {
	user := models.User{
		ID:   userId,
		Salt: salt,
	}
	request := dto.CheckEmailVerificationCodeDTO{
		Email: userEmail,
		Code:  verificationCode,
	}

	mockRepo := &mockRepository{}
	mockProofGenerator := &mocks.MockProofGenerator{
		GenerateProofFunc: func(
			emailAddress string, userSalt string, code string,
		) ([]byte, uint, error) {
			if emailAddress != userEmail {
				t.Fatalf("User's email mismatch: expected %s, got %s", userEmail, emailAddress)
			}
			if userSalt != salt {
				t.Fatalf("User's salt mimatch: expected %s, got %s", salt, userSalt)
			}
			if code != verificationCode {
				t.Fatalf("Verification code mismatch: expected %s, got %s", verificationCode, code)
			}

			return nil, zkKeyPairId, fmt.Errorf("failed to generate a zk proof")
		},
	}

	currService := service.NewService(mockRepo, &verification.EmailVerifier{}, mockProofGenerator)
	_, actualZkKeyPairId, err := currService.GenerateZkProofOfEmailVerification(user, request)

	assert.NotNil(t, err)
	assert.Equal(t, zkKeyPairId, actualZkKeyPairId)
}

func TestGenerateZkProofOfEmailVerification_Success(t *testing.T) {
	user := models.User{
		ID:   userId,
		Salt: salt,
	}
	request := dto.CheckEmailVerificationCodeDTO{
		Email: userEmail,
		Code:  verificationCode,
	}

	mockRepo := &mockRepository{}
	mockProofGenerator := &mocks.MockProofGenerator{
		GenerateProofFunc: func(
			emailAddress string, userSalt string, code string,
		) ([]byte, uint, error) {
			if emailAddress != userEmail {
				t.Fatalf("User's email mismatch: expected %s, got %s", userEmail, emailAddress)
			}
			if userSalt != salt {
				t.Fatalf("User's salt mimatch: expected %s, got %s", salt, userSalt)
			}
			if code != verificationCode {
				t.Fatalf("Verification code mismatch: expected %s, got %s", verificationCode, code)
			}

			return zkProof, zkKeyPairId, nil
		},
	}

	currService := service.NewService(mockRepo, &verification.EmailVerifier{}, mockProofGenerator)
	proof, actualZkKeyPairId, err := currService.GenerateZkProofOfEmailVerification(user, request)

	assert.Nil(t, err)

	assert.True(t, utils.Equal(zkProof, proof))
	assert.Equal(t, zkKeyPairId, actualZkKeyPairId)
}

func TestGetUserForUsername_UserNotFound(t *testing.T) {
	mockRepo := &mockRepository{
		getUserForUsername: func(username string) (models.User, error) {
			return models.User{}, fmt.Errorf("user not found")
		},
	}

	currService := service.NewService(mockRepo, &verification.EmailVerifier{}, &mocks.MockProofGenerator{})
	_, err := currService.GetUserForUsername(username)

	assert.NotNil(t, err)
}

func TestGetUserForUsername_Success(t *testing.T) {
	mockRepo := &mockRepository{
		getUserForUsername: func(currUsername string) (models.User, error) {
			if currUsername != username {
				t.Fatalf("Username mismatch: expected %s, got %s", username, currUsername)
			}

			return models.User{
				ID:        userId,
				Username:  username,
				Password:  password,
				FirstName: firstName,
				LastName:  lastName,
			}, nil
		},
	}

	currService := service.NewService(mockRepo, &verification.EmailVerifier{}, &mocks.MockProofGenerator{})
	user, err := currService.GetUserForUsername(username)

	assert.Nil(t, err)
	assert.Equal(t, userId, user.ID)
	assert.Equal(t, username, user.Username)
	assert.Equal(t, password, user.Password)
	assert.Equal(t, firstName, user.FirstName)
	assert.Equal(t, lastName, user.LastName)
}

func TestUpdateUserPassword_FailedToUpdatePasswordInDB(t *testing.T) {
	mockRepo := &mockRepository{
		updateUserPassword: func(currUsername string, currPassword string) error {
			if currUsername != username {
				t.Fatalf("Username mismatch: expected %s, got %s", username, currUsername)
			}
			if currPassword != hashedPassword {
				t.Fatalf("User password mismatch: expected %s, got %s", hashedPassword, currPassword)
			}

			return fmt.Errorf("failed to update user password")
		},
	}

	currService := service.NewService(mockRepo, &verification.EmailVerifier{}, &mocks.MockProofGenerator{})
	err := currService.UpdateUserPassword(username, password, salt)

	assert.NotNil(t, err)
}

func TestUpdateUserPassword_Success(t *testing.T) {
	mockRepo := &mockRepository{
		updateUserPassword: func(currUsername string, currPassword string) error {
			if currUsername != username {
				t.Fatalf("Username mismatch: expected %s, got %s", username, currUsername)
			}
			if currPassword != hashedPassword {
				t.Fatalf("User password mismatch: expected %s, got %s", hashedPassword, currPassword)
			}

			return nil
		},
	}

	currService := service.NewService(mockRepo, &verification.EmailVerifier{}, &mocks.MockProofGenerator{})
	err := currService.UpdateUserPassword(username, password, salt)

	assert.Nil(t, err)
}

func TestValidateSignature_SignatureIsInvalid(t *testing.T) {
	currService := service.NewService(
		&mockRepository{},
		&verification.EmailVerifier{},
		&mocks.MockProofGenerator{},
	)

	message := "sign-in with layer8"
	privateKey, _ := crypto.GenerateKey()

	signature := make([]byte, 64)
	rand.Read(signature)

	err := currService.ValidateSignature(
		message,
		signature,
		crypto.FromECDSAPub(&privateKey.PublicKey),
	)

	assert.NotNil(t, err)
}

func TestValidateSignature_Success(t *testing.T) {
	currService := service.NewService(
		&mockRepository{},
		&verification.EmailVerifier{},
		&mocks.MockProofGenerator{},
	)

	message := "sign-in with layer8"
	privateKey, _ := crypto.GenerateKey()
	signature, _ := crypto.Sign(crypto.Keccak256([]byte(message)), privateKey)

	err := currService.ValidateSignature(
		message,
		signature[:64],
		crypto.FromECDSAPub(&privateKey.PublicKey),
	)

	assert.Nil(t, err)
}

func TestRegisterUserPrecheck_Success(t *testing.T) {
	mockRepo := &mockRepository{
		registerUserPrecheck: func(req dto.RegisterUserPrecheckDTO, rmSalt string, iterCount int) error {
			assert.Equal(t, username, req.Username, "Username should match")
			assert.NotEmpty(t, rmSalt, "Salt should not be empty")
			assert.Equal(t, 4096, iterCount, "Iteration count should match")
			return nil
		},
	}

	currService := service.NewService(mockRepo, &verification.EmailVerifier{}, &mocks.MockProofGenerator{})

	req := dto.RegisterUserPrecheckDTO{
		Username: username,
	}
	iterCount := 4096

	salt, err := currService.RegisterUserPrecheck(req, iterCount)

	assert.Nil(t, err, "Expected no error during RegisterUserPrecheck")
	assert.NotEmpty(t, salt, "Salt should not be empty in the response")
}

func TestRegisterUserPrecheck_RepositoryError(t *testing.T) {
	mockRepo := &mockRepository{
		registerUserPrecheck: func(req dto.RegisterUserPrecheckDTO, rmSalt string, iterCount int) error {
			return errors.New("repository error")
		},
	}

	currService := service.NewService(mockRepo, &verification.EmailVerifier{}, &mocks.MockProofGenerator{})

	req := dto.RegisterUserPrecheckDTO{
		Username: username,
	}
	iterCount := 4096

	salt, err := currService.RegisterUserPrecheck(req, iterCount)

	assert.NotNil(t, err, "Expected an error during RegisterUserPrecheck")
	assert.Equal(t, "repository error", err.Error(), "Error message should match")
	assert.Empty(t, salt, "Salt should be empty in the response")
}

func TestRegisterUserPrecheck_InvalidIterationCount(t *testing.T) {
	mockRepo := &mockRepository{
		registerUserPrecheck: func(req dto.RegisterUserPrecheckDTO, rmSalt string, iterCount int) error {
			assert.Equal(t, username, req.Username, "Username should match")
			assert.Equal(t, 0, iterCount, "Iteration count should match")
			return nil
		},
	}

	currService := service.NewService(mockRepo, &verification.EmailVerifier{}, &mocks.MockProofGenerator{})

	req := dto.RegisterUserPrecheckDTO{
		Username: username,
	}
	iterCount := 0

	salt, err := currService.RegisterUserPrecheck(req, iterCount)

	assert.Nil(t, err, "Expected no error during RegisterUserPrecheck")
	assert.NotEmpty(t, salt, "Salt should not be empty in the response")
}

func TestRegisterUserv2_RepositoryFailedToStoreUserData(t *testing.T) {
	mockRepo := &mockRepository{
		registerUserv2: func(req dto.RegisterUserDTOv2) error {
			assert.Equal(t, username, req.Username)
			assert.Equal(t, firstName, req.FirstName)
			assert.Equal(t, lastName, req.LastName)
			assert.Equal(t, displayName, req.DisplayName)
			assert.Equal(t, country, req.Country)
			assert.Equal(t, []byte(publicKey), req.PublicKey)
			assert.Equal(t, storedKey, req.StoredKey)
			assert.Equal(t, serverKey, req.ServerKey)

			return fmt.Errorf("failed to store a user")
		},
	}
	currService := service.NewService(mockRepo, &verification.EmailVerifier{}, &mocks.MockProofGenerator{})

	err := currService.RegisterUserv2(
		dto.RegisterUserDTOv2{
			Username:    username,
			FirstName:   firstName,
			LastName:    lastName,
			DisplayName: displayName,
			Country:     country,
			PublicKey:   []byte(publicKey),
			StoredKey:   storedKey,
			ServerKey:   serverKey,
		},
	)

	assert.NotNil(t, err)
}

func TestRegisterUserv2_Success(t *testing.T) {
	mockRepo := &mockRepository{
		registerUserv2: func(req dto.RegisterUserDTOv2) error {
			assert.Equal(t, username, req.Username)
			assert.Equal(t, firstName, req.FirstName)
			assert.Equal(t, lastName, req.LastName)
			assert.Equal(t, displayName, req.DisplayName)
			assert.Equal(t, country, req.Country)

			return nil
		},
	}
	currService := service.NewService(mockRepo, &verification.EmailVerifier{}, &mocks.MockProofGenerator{})

	err := currService.RegisterUserv2(
		dto.RegisterUserDTOv2{
			Username:    username,
			FirstName:   firstName,
			LastName:    lastName,
			DisplayName: displayName,
			Country:     country,
			PublicKey:   []byte(publicKey),
			StoredKey:   storedKey,
			ServerKey:   serverKey,
		},
	)

	assert.Nil(t, err)
}

func TestUpdateUserPasswordV2_Success(t *testing.T) {
	mockRepo := &mockRepository{
		updateUserPasswordV2: func(currUsername, currStoredKey, currServerKey string) error {
			assert.Equal(t, username, currUsername)
			assert.Equal(t, storedKey, currStoredKey)
			assert.Equal(t, serverKey, currServerKey)
			return nil
		},
	}

	currService := service.NewService(mockRepo, nil, nil)

	err := currService.UpdateUserPasswordV2(username, storedKey, serverKey)
	assert.NoError(t, err, "Expected no error for successful password update")
}

func TestUpdateUserPasswordV2_RepositoryError(t *testing.T) {
	mockRepo := &mockRepository{
		updateUserPasswordV2: func(currUsername, currStoredKey, currServerKey string) error {
			assert.Equal(t, username, currUsername)
			assert.Equal(t, storedKey, currStoredKey)
			assert.Equal(t, serverKey, currServerKey)
			return fmt.Errorf("database error")
		},
	}

	currService := service.NewService(mockRepo, nil, nil)

	err := currService.UpdateUserPasswordV2(username, storedKey, serverKey)
	assert.Error(t, err, "Expected an error when repository returns an error")
	assert.Equal(t, "database error", err.Error())
}

func TestLoginPreCheckClientV2_RepositoryError(t *testing.T) {
	mockRepo := &mockRepository{
		profileClient: func(username string) (models.Client, error) {
			return models.Client{}, errors.New("repository error")
		},
	}
	currService := service.NewService(mockRepo, &verification.EmailVerifier{}, &mocks.MockProofGenerator{})

	req := dto.LoginPrecheckDTOv2{
		Username: username,
		CNonce:   cNonce,
	}

	loginClientPrecheckResp, err := currService.LoginPrecheckClientv2(req)

	assert.NotNil(t, err)
	assert.Equal(t, "repository error", err.Error())
	assert.Empty(t, loginClientPrecheckResp)
}

func TestLoginPreCheckClientV2_Success(t *testing.T) {
	mockRepo := &mockRepository{
		profileClient: func(username string) (models.Client, error) {
			return models.Client{
				Username:       username,
				Salt:           salt,
				IterationCount: iterationCount,
			}, nil
		},
	}

	currService := service.NewService(mockRepo, &verification.EmailVerifier{}, &mocks.MockProofGenerator{})

	req := dto.LoginPrecheckDTOv2{
		Username: username,
		CNonce:   cNonce,
	}

	loginPrecheckResp, err := currService.LoginPrecheckClientv2(req)

	assert.Nil(t, err)
	assert.Equal(t, loginPrecheckResp.Salt, salt)
	assert.True(t, strings.HasPrefix(loginPrecheckResp.Nonce, cNonce))
	assert.Equal(t, loginPrecheckResp.IterCount, iterationCount)
}

func TestLoginClientV2_RepositoryError(t *testing.T) {
	mockRepo := &mockRepository{
		profileClient: func(username string) (models.Client, error) {
			return models.Client{}, errors.New("repository error")
		},
	}
	currService := service.NewService(mockRepo, &verification.EmailVerifier{}, &mocks.MockProofGenerator{})

	req := dto.LoginClientDTOv2{
		Username:    username,
		CNonce:      cNonce,
		Nonce:       nonce,
		ClientProof: clientProof,
	}

	loginClientResp, err := currService.LoginClientv2(req)

	assert.NotNil(t, err)
	assert.Equal(t, "repository error", err.Error())
	assert.Empty(t, loginClientResp)
}

func TestLoginClientV2_DecodingStoredKeyError(t *testing.T) {
	mockRepo := &mockRepository{
		profileClient: func(username string) (models.Client, error) {
			return models.Client{
				StoredKey: "TEST_STORED_KEY_FOR_DECODE_ERROR_!@#", // Invalid stored key
			}, nil
		},
	}
	currService := service.NewService(mockRepo, &verification.EmailVerifier{}, &mocks.MockProofGenerator{})

	req := dto.LoginClientDTOv2{
		Username:    username,
		CNonce:      cNonce,
		Nonce:       nonce,
		ClientProof: clientProof,
	}

	loginClientResp, err := currService.LoginClientv2(req)

	assert.NotNil(t, err)
	assert.Equal(t, "error decoding stored key: encoding/hex: invalid byte: U+0054 'T'", err.Error())
	assert.Empty(t, loginClientResp)
}

func TestLoginClientV2_DecodingClientProofError(t *testing.T) {
	mockRepo := &mockRepository{
		profileClient: func(username string) (models.Client, error) {
			return models.Client{
				StoredKey:      storedKey,
				Salt:           salt,
				IterationCount: iterationCount,
			}, nil
		},
	}
	currService := service.NewService(mockRepo, &verification.EmailVerifier{}, &mocks.MockProofGenerator{})

	req := dto.LoginClientDTOv2{
		Username:    username,
		CNonce:      cNonce,
		Nonce:       nonce,
		ClientProof: "TEST_CLIENT_PROOF_FOR_DECODE_ERROR_!@#", // Invalid client proof
	}

	loginClientResp, err := currService.LoginClientv2(req)

	assert.NotNil(t, err)
	assert.Equal(t, "error decoding client proof: encoding/hex: invalid byte: U+0054 'T'", err.Error())
	assert.Empty(t, loginClientResp)
}

func TestLoginClientV2_XorOperationError(t *testing.T) {
	mockRepo := &mockRepository{
		profileClient: func(username string) (models.Client, error) {
			return models.Client{
				StoredKey:      storedKey,
				Salt:           salt,
				IterationCount: iterationCount,
			}, nil
		},
	}
	currService := service.NewService(mockRepo, &verification.EmailVerifier{}, &mocks.MockProofGenerator{})

	req := dto.LoginClientDTOv2{
		Username:    username,
		CNonce:      cNonce,
		Nonce:       nonce,
		ClientProof: "", // Sending empty client proof to fail the XOR operation, since it requires equal length slices
	}

	loginClientResp, err := currService.LoginClientv2(req)

	assert.NotNil(t, err)
	assert.Equal(t, "error performing XOR operation: slices must have the same length", err.Error())
	assert.Empty(t, loginClientResp)
}

func TestLoginClientV2_AuthFailedKeyMismatchError(t *testing.T) {
	mockRepo := &mockRepository{
		profileClient: func(username string) (models.Client, error) {
			return models.Client{
				StoredKey:      "", // Empty stored key
				Salt:           salt,
				IterationCount: iterationCount,
			}, nil
		},
	}
	currService := service.NewService(mockRepo, &verification.EmailVerifier{}, &mocks.MockProofGenerator{})

	req := dto.LoginClientDTOv2{
		Username:    username,
		CNonce:      cNonce,
		Nonce:       nonce,
		ClientProof: clientProof,
	}

	loginClientResp, err := currService.LoginClientv2(req)

	assert.NotNil(t, err)
	assert.Equal(t, "server failed to authenticate the user", err.Error())
	assert.Empty(t, loginClientResp)
}

func TestLoginClientV2_DecodingServerKeyError(t *testing.T) {
	mockRepo := &mockRepository{
		profileClient: func(username string) (models.Client, error) {
			return models.Client{
				StoredKey:      storedKey,
				ServerKey:      "TEST_SERVER_KEY_FOR_DECODE_ERROR_!@#", // Invalid server key
				Salt:           salt,
				IterationCount: iterationCount,
			}, nil
		},
	}
	currService := service.NewService(mockRepo, &verification.EmailVerifier{}, &mocks.MockProofGenerator{})

	req := dto.LoginClientDTOv2{
		Username:    username,
		CNonce:      cNonce,
		Nonce:       nonce,
		ClientProof: clientProof,
	}

	loginClientResp, err := currService.LoginClientv2(req)

	assert.NotNil(t, err)
	assert.Equal(t, "error decoding server key: encoding/hex: invalid byte: U+0054 'T'", err.Error())
	assert.Empty(t, loginClientResp)
}

func TestLoginClientV2_Success(t *testing.T) {
	mockRepo := &mockRepository{
		profileClient: func(username string) (models.Client, error) {
			return models.Client{
				StoredKey:      storedKey,
				ServerKey:      serverKey,
				Salt:           salt,
				IterationCount: iterationCount,
			}, nil
		},
	}
	currService := service.NewService(mockRepo, &verification.EmailVerifier{}, &mocks.MockProofGenerator{})

	req := dto.LoginClientDTOv2{
		Username:    username,
		CNonce:      cNonce,
		Nonce:       nonce,
		ClientProof: clientProof,
	}

	loginClientResp, err := currService.LoginClientv2(req)

	assert.Nil(t, err)
	assert.Equal(t, testServerSignature, loginClientResp.ServerSignature)
}

func TestRegisterClientPrecheck_Success(t *testing.T) {
	mockRepo := &mockRepository{
		registerClientPrecheck: func(req dto.RegisterClientPrecheckDTO, rmSalt string, iterCount int) error {
			assert.Equal(t, username, req.Username, "Username should match")
			assert.NotEmpty(t, rmSalt, "Salt should not be empty")
			assert.Equal(t, 4096, iterCount, "Iteration count should match")
			return nil
		},
	}

	currService := service.NewService(mockRepo, &verification.EmailVerifier{}, &mocks.MockProofGenerator{})

	req := dto.RegisterClientPrecheckDTO{
		Username: username,
	}

	salt, err := currService.RegisterClientPrecheck(req, iterationCount)

	assert.Nil(t, err, "Expected no error during RegisterUserPrecheck")
	assert.NotEmpty(t, salt, "Salt should not be empty in the response")
}

func TestRegisterClientPrecheck_RepositoryError(t *testing.T) {
	mockRepo := &mockRepository{
		registerClientPrecheck: func(req dto.RegisterClientPrecheckDTO, rmSalt string, iterCount int) error {
			return errors.New("repository error")
		},
	}

	currService := service.NewService(mockRepo, &verification.EmailVerifier{}, &mocks.MockProofGenerator{})

	req := dto.RegisterClientPrecheckDTO{
		Username: username,
	}

	salt, err := currService.RegisterClientPrecheck(req, iterationCount)

	assert.NotNil(t, err, "Expected an error during RegisterClientPrecheck")
	assert.Equal(t, "repository error", err.Error(), "Error message should match")
	assert.Empty(t, salt, "Salt should be empty in the response")
}

func TestRegisterClientPrecheck_InvalidIterationCount(t *testing.T) {
	mockRepo := &mockRepository{
		registerClientPrecheck: func(req dto.RegisterClientPrecheckDTO, rmSalt string, iterCount int) error {
			assert.Equal(t, username, req.Username, "Username should match")
			assert.Equal(t, 0, iterCount, "Iteration count should match")
			return nil
		},
	}

	currService := service.NewService(mockRepo, &verification.EmailVerifier{}, &mocks.MockProofGenerator{})

	req := dto.RegisterClientPrecheckDTO{
		Username: username,
	}
	iterCount := 0

	salt, err := currService.RegisterClientPrecheck(req, iterCount)

	assert.Nil(t, err, "Expected no error during RegisterClientPrecheck")
	assert.NotEmpty(t, salt, "Salt should not be empty in the response")
}

func TestRegisterClientv2_RepositoryFailedToStoreUserData(t *testing.T) {
	mockRepo := &mockRepository{
		registerClientv2: func(req dto.RegisterClientDTOv2, id string, secret string) error {
			assert.Equal(t, username, req.Username)
			assert.Equal(t, clientName, req.Name)
			assert.Equal(t, redirectUri, req.RedirectURI)
			assert.Equal(t, backendUri, req.BackendURI)
			assert.Equal(t, storedKey, req.StoredKey)
			assert.Equal(t, serverKey, req.ServerKey)

			return fmt.Errorf("failed to store a client")
		},
	}
	currService := service.NewService(mockRepo, &verification.EmailVerifier{}, &mocks.MockProofGenerator{})

	err := currService.RegisterClientv2(
		dto.RegisterClientDTOv2{
			Username:    username,
			Name:        clientName,
			RedirectURI: redirectUri,
			BackendURI:  backendUri,
			StoredKey:   storedKey,
			ServerKey:   serverKey,
		},
	)

	assert.NotNil(t, err)
}

func TestRegisterClientv2_Success(t *testing.T) {
	mockRepo := &mockRepository{
		registerClientv2: func(req dto.RegisterClientDTOv2, id string, secret string) error {
			assert.Equal(t, username, req.Username)
			assert.Equal(t, clientName, req.Name)
			assert.Equal(t, redirectUri, req.RedirectURI)
			assert.Equal(t, backendUri, req.BackendURI)
			assert.Equal(t, storedKey, req.StoredKey)
			assert.Equal(t, serverKey, req.ServerKey)

			return nil
		},
	}
	currService := service.NewService(mockRepo, &verification.EmailVerifier{}, &mocks.MockProofGenerator{})

	err := currService.RegisterClientv2(
		dto.RegisterClientDTOv2{
			Username:    username,
			Name:        clientName,
			RedirectURI: redirectUri,
			BackendURI:  backendUri,
			StoredKey:   storedKey,
			ServerKey:   serverKey,
		},
	)

	assert.Nil(t, err)
}
