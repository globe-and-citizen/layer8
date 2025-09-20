package service_test

import (
	"crypto/rand"
	"errors"
	"fmt"
	serverModels "globe-and-citizen/layer8/server/models"
	"globe-and-citizen/layer8/server/resource_server/dto"
	"globe-and-citizen/layer8/server/resource_server/emails/verification"
	"globe-and-citizen/layer8/server/resource_server/emails/verification/code"
	"globe-and-citizen/layer8/server/resource_server/emails/verification/zk"
	"globe-and-citizen/layer8/server/resource_server/models"
	"globe-and-citizen/layer8/server/resource_server/service"
	"globe-and-citizen/layer8/server/resource_server/utils"
	"globe-and-citizen/layer8/server/resource_server/utils/mocks"
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
const displayName = "display_name"
const color = "color"
const bio = "bio"
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
	profileUser                  func(userID uint) (models.User, models.UserMetadata, error)
	saveEmailVerificationData    func(data models.EmailVerificationData) error
	getEmailVerificationData     func(userId uint) (models.EmailVerificationData, error)
	saveProofOfEmailVerification func(userID uint, verificationCode string, proof []byte, zkKeyPairId uint) error
	registerUserPrecheck         func(req dto.RegisterUserPrecheckDTO, salt string, iterCount int) error
	registerClientPrecheck       func(req dto.RegisterClientPrecheckDTO, salt string, iterCount int) error
	registerUser                 func(req dto.RegisterUserDTO) error
	registerClient               func(req dto.RegisterClientDTO, id string, secret string) error
	getUserForUsername           func(username string) (models.User, error)
	profileClient                func(username string) (models.Client, error)
	updateUserPassword           func(username string, storedKey string, serverKey string) error
	getUserMetadata              func(userID int64, key string) (*serverModels.UserMetadata, error)
	updateUserMetadata           func(userID uint, req dto.UpdateUserMetadataDTO) error
}

func (m *mockRepository) FindUser(userId uint) (models.User, error) {
	return m.findUser(userId)
}

func (m *mockRepository) RegisterPrecheckUser(req dto.RegisterUserPrecheckDTO, salt string, iterCount int) error {
	if m.registerUserPrecheck != nil {
		return m.registerUserPrecheck(req, salt, iterCount)
	}
	return nil
}

func (m *mockRepository) RegisterUser(req dto.RegisterUserDTO) error {
	return m.registerUser(req)
}

func (m *mockRepository) ProfileUser(userID uint) (models.User, models.UserMetadata, error) {
	return m.profileUser(userID)
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

func (m *mockRepository) UpdateUserMetadata(userID uint, req dto.UpdateUserMetadataDTO) error {
	return m.updateUserMetadata(userID, req)
}

func (m *mockRepository) RegisterPrecheckClient(req dto.RegisterClientPrecheckDTO, salt string, iterCount int) error {
	return m.registerClientPrecheck(req, salt, iterCount)
}

func (m *mockRepository) RegisterClient(req dto.RegisterClientDTO, id string, secret string) error {
	return m.registerClient(req, id, secret)
}

func (m *mockRepository) IsBackendURIExists(backendURL string) (bool, error) {
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

func (m *mockRepository) GetUser(username string) (*serverModels.User, error) {
	return &serverModels.User{}, nil
}

func (m *mockRepository) GetUserByID(id int64) (*serverModels.User, error) {
	return &serverModels.User{}, nil
}

func (m *mockRepository) GetUserMetadata(userID int64, key string) (*serverModels.UserMetadata, error) {
	return m.getUserMetadata(userID, key)
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

func (m *mockRepository) UpdateUserPassword(username string, storedKey string, serverKey string) error {
	return m.updateUserPassword(username, storedKey, serverKey)
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

func (m *mockRepository) SavePhoneNumberVerificationData(data models.PhoneNumberVerificationData) error {
	return nil
}

func (m *mockRepository) GetPhoneNumberVerificationData(userID uint) (models.PhoneNumberVerificationData, error) {
	return models.PhoneNumberVerificationData{}, nil
}

func (m *mockRepository) SaveProofOfPhoneNumberVerification(userID uint, verificationCode string, zkProof []byte, zkPairID uint) error {
	return nil
}

func (m *mockRepository) SaveTelegramSessionIDHash(userID uint, sessionID []byte) error {
	return nil
}

func TestLoginPreCheckUser_RepositoryError(t *testing.T) {
	mockRepo := &mockRepository{
		getUserForUsername: func(username string) (models.User, error) {
			return models.User{}, errors.New("repository error")
		},
	}
	currService := service.NewService(mockRepo, &verification.EmailVerifier{}, &mocks.MockProofGenerator{}, code.NewMIMCCodeGenerator())

	req := dto.LoginPrecheckDTO{
		Username: username,
		CNonce:   cNonce,
	}

	loginPrecheckResp, err := currService.LoginPrecheckUser(req)

	assert.NotNil(t, err)
	assert.Equal(t, "repository error", err.Error())
	assert.Empty(t, loginPrecheckResp)
}

func TestLoginPreCheckUser_Success(t *testing.T) {
	mockRepo := &mockRepository{
		getUserForUsername: func(username string) (models.User, error) {
			return models.User{
				Username:       username,
				Salt:           salt,
				IterationCount: iterationCount,
			}, nil
		},
	}

	currService := service.NewService(mockRepo, &verification.EmailVerifier{}, &mocks.MockProofGenerator{}, code.NewMIMCCodeGenerator())

	req := dto.LoginPrecheckDTO{
		Username: username,
		CNonce:   cNonce,
	}

	loginPrecheckResp, err := currService.LoginPrecheckUser(req)

	assert.Nil(t, err)
	assert.Equal(t, loginPrecheckResp.Salt, salt)
	assert.Equal(t, strings.HasPrefix(loginPrecheckResp.Nonce, cNonce), true)
	assert.Equal(t, loginPrecheckResp.IterCount, iterationCount)
}

func TestLoginUser_RepositoryError(t *testing.T) {
	mockRepo := &mockRepository{
		getUserForUsername: func(username string) (models.User, error) {
			return models.User{}, errors.New("repository error")
		},
	}
	currService := service.NewService(mockRepo, &verification.EmailVerifier{}, &mocks.MockProofGenerator{}, code.NewMIMCCodeGenerator())

	req := dto.LoginUserDTO{
		Username:    username,
		CNonce:      cNonce,
		Nonce:       nonce,
		ClientProof: clientProof,
	}

	loginUserResp, err := currService.LoginUser(req)

	assert.NotNil(t, err)
	assert.Equal(t, "repository error", err.Error())
	assert.Empty(t, loginUserResp)
}

func TestLoginUser_DecodingStoredKeyError(t *testing.T) {
	mockRepo := &mockRepository{
		getUserForUsername: func(username string) (models.User, error) {
			return models.User{
				StoredKey: "TEST_STORED_KEY_FOR_DECODE_ERROR_!@#", // Invalid stored key
			}, nil
		},
	}
	currService := service.NewService(mockRepo, &verification.EmailVerifier{}, &mocks.MockProofGenerator{}, code.NewMIMCCodeGenerator())

	req := dto.LoginUserDTO{
		Username:    username,
		CNonce:      cNonce,
		Nonce:       nonce,
		ClientProof: clientProof,
	}

	loginUserResp, err := currService.LoginUser(req)

	assert.NotNil(t, err)
	assert.Equal(t, "error decoding stored key: encoding/hex: invalid byte: U+0054 'T'", err.Error())
	assert.Empty(t, loginUserResp)
}

func TestLoginUser_DecodingClientProofError(t *testing.T) {
	mockRepo := &mockRepository{
		getUserForUsername: func(username string) (models.User, error) {
			return models.User{
				StoredKey:      storedKey,
				Salt:           salt,
				IterationCount: iterationCount,
			}, nil
		},
	}
	currService := service.NewService(mockRepo, &verification.EmailVerifier{}, &mocks.MockProofGenerator{}, code.NewMIMCCodeGenerator())

	req := dto.LoginUserDTO{
		Username:    username,
		CNonce:      cNonce,
		Nonce:       nonce,
		ClientProof: "TEST_CLIENT_PROOF_FOR_DECODE_ERROR_!@#", // Invalid client proof
	}

	loginUserResp, err := currService.LoginUser(req)

	assert.NotNil(t, err)
	assert.Equal(t, "error decoding client proof: encoding/hex: invalid byte: U+0054 'T'", err.Error())
	assert.Empty(t, loginUserResp)
}

func TestLoginUser_XorOperationError(t *testing.T) {
	mockRepo := &mockRepository{
		getUserForUsername: func(username string) (models.User, error) {
			return models.User{
				StoredKey:      storedKey,
				Salt:           salt,
				IterationCount: iterationCount,
			}, nil
		},
	}
	currService := service.NewService(mockRepo, &verification.EmailVerifier{}, &mocks.MockProofGenerator{}, code.NewMIMCCodeGenerator())

	req := dto.LoginUserDTO{
		Username:    username,
		CNonce:      cNonce,
		Nonce:       nonce,
		ClientProof: "", // Sending empty client proof to fail the XOR operation, since it requires equal length slices
	}

	loginUserResp, err := currService.LoginUser(req)

	assert.NotNil(t, err)
	assert.Equal(t, "error performing XOR operation: slices must have the same length", err.Error())
	assert.Empty(t, loginUserResp)
}

func TestLoginUser_AuthFailedKeyMismatchError(t *testing.T) {
	mockRepo := &mockRepository{
		getUserForUsername: func(username string) (models.User, error) {
			return models.User{
				StoredKey:      "", // Empty stored key
				Salt:           salt,
				IterationCount: iterationCount,
			}, nil
		},
	}
	currService := service.NewService(mockRepo, &verification.EmailVerifier{}, &mocks.MockProofGenerator{}, code.NewMIMCCodeGenerator())

	req := dto.LoginUserDTO{
		Username:    username,
		CNonce:      cNonce,
		Nonce:       nonce,
		ClientProof: clientProof,
	}

	loginUserResp, err := currService.LoginUser(req)

	assert.NotNil(t, err)
	assert.Equal(t, "server failed to authenticate the user", err.Error())
	assert.Empty(t, loginUserResp)
}

func TestLoginUser_DecodingServerKeyError(t *testing.T) {
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
	currService := service.NewService(mockRepo, &verification.EmailVerifier{}, &mocks.MockProofGenerator{}, code.NewMIMCCodeGenerator())

	req := dto.LoginUserDTO{
		Username:    username,
		CNonce:      cNonce,
		Nonce:       nonce,
		ClientProof: clientProof,
	}

	loginUserResp, err := currService.LoginUser(req)

	assert.NotNil(t, err)
	assert.Equal(t, "error decoding server key: encoding/hex: invalid byte: U+0054 'T'", err.Error())
	assert.Empty(t, loginUserResp)
}

func TestLoginUser_Success(t *testing.T) {
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
	currService := service.NewService(mockRepo, &verification.EmailVerifier{}, &mocks.MockProofGenerator{}, code.NewMIMCCodeGenerator())

	req := dto.LoginUserDTO{
		Username:    username,
		CNonce:      cNonce,
		Nonce:       nonce,
		ClientProof: clientProof,
	}

	loginUserResp, err := currService.LoginUser(req)

	assert.Nil(t, err)
	assert.Equal(t, testServerSignature, loginUserResp.ServerSignature)
}

func TestProfileUser(t *testing.T) {
	// Create a new mock repository
	mockRepo := &mockRepository{
		profileUser: func(userID uint) (models.User, models.UserMetadata, error) {
			if userID != userId {
				t.Fatalf("userID mismatch, expected %d, got %d", userId, userID)
			}

			return models.User{
					ID:       userID,
					Username: username,
				}, models.UserMetadata{
					ID:          userID,
					DisplayName: displayName,
					Color:       color,
					Bio:         bio,
				}, nil
		},
	}

	// Create a new service by passing the mock repository
	mockService := service.NewService(mockRepo, &verification.EmailVerifier{}, &mocks.MockProofGenerator{}, code.NewMIMCCodeGenerator())

	// Call the ProfileUser method of the mock service
	userDetails, err := mockService.ProfileUser(userId)
	if err != nil {
		t.Error("Expected nil, got", err)
	}

	// Use assert to check if the error is nil
	assert.Nil(t, err)
	assert.Equal(t, username, userDetails.Username)
	assert.Equal(t, displayName, userDetails.DisplayName)
	assert.Equal(t, bio, userDetails.Bio)
	assert.Equal(t, color, userDetails.Color)
}

func TestUpdateUserMetadata(t *testing.T) {
	mockRepo := &mockRepository{
		updateUserMetadata: func(userID uint, req dto.UpdateUserMetadataDTO) error {
			if userID != userId {
				t.Fatalf("userID mismatch, expected %d, got %d", userId, userID)
			}
			if req.DisplayName != displayName {
				t.Fatalf("displayName mismatch, expected %s, got %s", displayName, req.DisplayName)
			}
			if req.Color != color {
				t.Fatalf("color mismatch, expected %s, got %s", color, req.Color)
			}
			if req.Bio != bio {
				t.Fatalf("bio mismatch, expected %s, got %s", bio, req.Bio)
			}
			return nil
		},
	}

	// Create a new service by passing the mock repository
	mockService := service.NewService(mockRepo, &verification.EmailVerifier{}, &mocks.MockProofGenerator{}, code.NewMIMCCodeGenerator())

	// Create a new mock request
	req := dto.UpdateUserMetadataDTO{
		DisplayName: displayName,
		Color:       color,
		Bio:         bio,
	}

	// Call the UpdateUserMetadata method of the mock service
	err := mockService.UpdateUserMetadata(userId, req)
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
	mockService := service.NewService(mockRepo, &verification.EmailVerifier{}, &mocks.MockProofGenerator{}, code.NewMIMCCodeGenerator())

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

	mockService := service.NewService(mockRepo, &verification.EmailVerifier{}, &zk.ProofProcessor{}, code.NewMIMCCodeGenerator())

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

	currService := service.NewService(mockRepo, emailVerifier, &mocks.MockProofGenerator{}, code.NewMIMCCodeGenerator())
	e := currService.VerifyEmail(userId, userEmail)

	assert.NotNil(t, e)
}

func TestVerifyEmail_UserExists_EmailFailedToBeSent(t *testing.T) {
	mockRepo := &mockRepository{
		findUser: func(userId uint) (models.User, error) {
			return models.User{
				ID:                    userId,
				Username:              username,
				EmailVerificationCode: "",
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
	currService := service.NewService(mockRepo, emailVerifier, &mocks.MockProofGenerator{}, code.NewMIMCCodeGenerator())

	e := currService.VerifyEmail(userId, userEmail)

	assert.NotNil(t, e)
}

func TestVerifyEmail_UserExists_EmailSent_VerificationDataNotSaved(t *testing.T) {
	mockRepo := &mockRepository{
		findUser: func(userId uint) (models.User, error) {
			return models.User{
				ID:                    userId,
				Username:              username,
				EmailVerificationCode: "",
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
	currService := service.NewService(mockRepo, emailVerifier, &mocks.MockProofGenerator{}, code.NewMIMCCodeGenerator())

	e := currService.VerifyEmail(userId, userEmail)

	assert.NotNil(t, e)
}

func TestVerifyEmail_Success(t *testing.T) {
	mockRepo := &mockRepository{
		findUser: func(userId uint) (models.User, error) {
			return models.User{
				ID:                    userId,
				Username:              username,
				EmailVerificationCode: "",
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
	currService := service.NewService(mockRepo, emailVerifier, &mocks.MockProofGenerator{}, code.NewMIMCCodeGenerator())

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
	currService := service.NewService(mockRepo, emailVerifier, &mocks.MockProofGenerator{}, code.NewMIMCCodeGenerator())

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
	currService := service.NewService(mockRepo, emailVerifier, &mocks.MockProofGenerator{}, code.NewMIMCCodeGenerator())

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
	currService := service.NewService(mockRepo, emailVerifier, &mocks.MockProofGenerator{}, code.NewMIMCCodeGenerator())

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
	currService := service.NewService(mockRepo, emailVerifier, &mocks.MockProofGenerator{}, code.NewMIMCCodeGenerator())

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
	currService := service.NewService(mockRepo, emailVerifier, &mocks.MockProofGenerator{}, code.NewMIMCCodeGenerator())

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
	currService := service.NewService(mockRepo, emailVerifier, &mocks.MockProofGenerator{}, code.NewMIMCCodeGenerator())

	e := currService.SaveProofOfEmailVerification(userId, verificationCode, emailProof, zkKeyPairId)

	assert.Nil(t, e)
}

func TestFindUser_UserNotFound(t *testing.T) {
	mockRepo := &mockRepository{
		findUser: func(userId uint) (models.User, error) {
			return models.User{}, fmt.Errorf("user not found for id %d", userId)
		},
	}

	currService := service.NewService(mockRepo, &verification.EmailVerifier{}, &mocks.MockProofGenerator{}, code.NewMIMCCodeGenerator())
	_, err := currService.FindUser(userId)

	assert.NotNil(t, err)
}

func TestFindUser_Success(t *testing.T) {
	mockRepo := &mockRepository{
		findUser: func(userId uint) (models.User, error) {
			return models.User{ID: userId}, nil
		},
	}

	currService := service.NewService(mockRepo, &verification.EmailVerifier{}, &mocks.MockProofGenerator{}, code.NewMIMCCodeGenerator())
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

	currService := service.NewService(mockRepo, &verification.EmailVerifier{}, mockProofGenerator, code.NewMIMCCodeGenerator())
	_, actualZkKeyPairId, err := currService.GenerateZkProof(user, request.Email, request.Code)

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

	currService := service.NewService(mockRepo, &verification.EmailVerifier{}, mockProofGenerator, code.NewMIMCCodeGenerator())
	proof, actualZkKeyPairId, err := currService.GenerateZkProof(user, request.Email, request.Code)

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

	currService := service.NewService(mockRepo, &verification.EmailVerifier{}, &mocks.MockProofGenerator{}, code.NewMIMCCodeGenerator())
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
				ID:       userId,
				Username: username,
			}, nil
		},
	}

	currService := service.NewService(mockRepo, &verification.EmailVerifier{}, &mocks.MockProofGenerator{}, code.NewMIMCCodeGenerator())
	user, err := currService.GetUserForUsername(username)

	assert.Nil(t, err)
	assert.Equal(t, userId, user.ID)
	assert.Equal(t, username, user.Username)
}

func TestValidateSignature_SignatureIsInvalid(t *testing.T) {
	currService := service.NewService(
		&mockRepository{},
		&verification.EmailVerifier{},
		&mocks.MockProofGenerator{},
		code.NewMIMCCodeGenerator(),
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
		code.NewMIMCCodeGenerator(),
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

	currService := service.NewService(mockRepo, &verification.EmailVerifier{}, &mocks.MockProofGenerator{}, code.NewMIMCCodeGenerator())

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

	currService := service.NewService(mockRepo, &verification.EmailVerifier{}, &mocks.MockProofGenerator{}, code.NewMIMCCodeGenerator())

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

	currService := service.NewService(mockRepo, &verification.EmailVerifier{}, &mocks.MockProofGenerator{}, code.NewMIMCCodeGenerator())

	req := dto.RegisterUserPrecheckDTO{
		Username: username,
	}
	iterCount := 0

	salt, err := currService.RegisterUserPrecheck(req, iterCount)

	assert.Nil(t, err, "Expected no error during RegisterUserPrecheck")
	assert.NotEmpty(t, salt, "Salt should not be empty in the response")
}

func TestRegisterUser_RepositoryFailedToStoreUserData(t *testing.T) {
	mockRepo := &mockRepository{
		registerUser: func(req dto.RegisterUserDTO) error {
			assert.Equal(t, username, req.Username)
			assert.Equal(t, []byte(publicKey), req.PublicKey)
			assert.Equal(t, storedKey, req.StoredKey)
			assert.Equal(t, serverKey, req.ServerKey)

			return fmt.Errorf("failed to store a user")
		},
	}
	currService := service.NewService(mockRepo, &verification.EmailVerifier{}, &mocks.MockProofGenerator{}, code.NewMIMCCodeGenerator())

	err := currService.RegisterUser(
		dto.RegisterUserDTO{
			Username:  username,
			PublicKey: []byte(publicKey),
			StoredKey: storedKey,
			ServerKey: serverKey,
		},
	)

	assert.NotNil(t, err)
}

func TestRegisterUser_Success(t *testing.T) {
	mockRepo := &mockRepository{
		registerUser: func(req dto.RegisterUserDTO) error {
			assert.Equal(t, username, req.Username)
			assert.Equal(t, storedKey, req.StoredKey)
			assert.Equal(t, serverKey, req.ServerKey)
			return nil
		},
	}
	currService := service.NewService(mockRepo, &verification.EmailVerifier{}, &mocks.MockProofGenerator{}, code.NewMIMCCodeGenerator())

	err := currService.RegisterUser(
		dto.RegisterUserDTO{
			Username:  username,
			PublicKey: []byte(publicKey),
			StoredKey: storedKey,
			ServerKey: serverKey,
		},
	)

	assert.Nil(t, err)
}

func TestUpdateUserPassword_Success(t *testing.T) {
	mockRepo := &mockRepository{
		updateUserPassword: func(currUsername, currStoredKey, currServerKey string) error {
			assert.Equal(t, username, currUsername)
			assert.Equal(t, storedKey, currStoredKey)
			assert.Equal(t, serverKey, currServerKey)
			return nil
		},
	}

	currService := service.NewService(mockRepo, nil, nil, code.NewMIMCCodeGenerator())

	err := currService.UpdateUserPassword(username, storedKey, serverKey)
	assert.NoError(t, err, "Expected no error for successful password update")
}

func TestUpdateUserPassword_RepositoryError(t *testing.T) {
	mockRepo := &mockRepository{
		updateUserPassword: func(currUsername, currStoredKey, currServerKey string) error {
			assert.Equal(t, username, currUsername)
			assert.Equal(t, storedKey, currStoredKey)
			assert.Equal(t, serverKey, currServerKey)
			return fmt.Errorf("database error")
		},
	}

	currService := service.NewService(mockRepo, nil, nil, code.NewMIMCCodeGenerator())

	err := currService.UpdateUserPassword(username, storedKey, serverKey)
	assert.Error(t, err, "Expected an error when repository returns an error")
	assert.Equal(t, "database error", err.Error())
}

func TestLoginPreCheckClient_RepositoryError(t *testing.T) {
	mockRepo := &mockRepository{
		profileClient: func(username string) (models.Client, error) {
			return models.Client{}, errors.New("repository error")
		},
	}
	currService := service.NewService(mockRepo, &verification.EmailVerifier{}, &mocks.MockProofGenerator{}, code.NewMIMCCodeGenerator())

	req := dto.LoginPrecheckDTO{
		Username: username,
		CNonce:   cNonce,
	}

	loginClientPrecheckResp, err := currService.LoginPrecheckClient(req)

	assert.NotNil(t, err)
	assert.Equal(t, "repository error", err.Error())
	assert.Empty(t, loginClientPrecheckResp)
}

func TestLoginPreCheckClient_Success(t *testing.T) {
	mockRepo := &mockRepository{
		profileClient: func(username string) (models.Client, error) {
			return models.Client{
				Username:       username,
				Salt:           salt,
				IterationCount: iterationCount,
			}, nil
		},
	}

	currService := service.NewService(mockRepo, &verification.EmailVerifier{}, &mocks.MockProofGenerator{}, code.NewMIMCCodeGenerator())

	req := dto.LoginPrecheckDTO{
		Username: username,
		CNonce:   cNonce,
	}

	loginPrecheckResp, err := currService.LoginPrecheckClient(req)

	assert.Nil(t, err)
	assert.Equal(t, loginPrecheckResp.Salt, salt)
	assert.True(t, strings.HasPrefix(loginPrecheckResp.Nonce, cNonce))
	assert.Equal(t, loginPrecheckResp.IterCount, iterationCount)
}

func TestLoginClient_RepositoryError(t *testing.T) {
	mockRepo := &mockRepository{
		profileClient: func(username string) (models.Client, error) {
			return models.Client{}, errors.New("repository error")
		},
	}
	currService := service.NewService(mockRepo, &verification.EmailVerifier{}, &mocks.MockProofGenerator{}, code.NewMIMCCodeGenerator())

	req := dto.LoginClientDTO{
		Username:    username,
		CNonce:      cNonce,
		Nonce:       nonce,
		ClientProof: clientProof,
	}

	loginClientResp, err := currService.LoginClient(req)

	assert.NotNil(t, err)
	assert.Equal(t, "repository error", err.Error())
	assert.Empty(t, loginClientResp)
}

func TestLoginClient_DecodingStoredKeyError(t *testing.T) {
	mockRepo := &mockRepository{
		profileClient: func(username string) (models.Client, error) {
			return models.Client{
				StoredKey: "TEST_STORED_KEY_FOR_DECODE_ERROR_!@#", // Invalid stored key
			}, nil
		},
	}
	currService := service.NewService(mockRepo, &verification.EmailVerifier{}, &mocks.MockProofGenerator{}, code.NewMIMCCodeGenerator())

	req := dto.LoginClientDTO{
		Username:    username,
		CNonce:      cNonce,
		Nonce:       nonce,
		ClientProof: clientProof,
	}

	loginClientResp, err := currService.LoginClient(req)

	assert.NotNil(t, err)
	assert.Equal(t, "error decoding stored key: encoding/hex: invalid byte: U+0054 'T'", err.Error())
	assert.Empty(t, loginClientResp)
}

func TestLoginClient_DecodingClientProofError(t *testing.T) {
	mockRepo := &mockRepository{
		profileClient: func(username string) (models.Client, error) {
			return models.Client{
				StoredKey:      storedKey,
				Salt:           salt,
				IterationCount: iterationCount,
			}, nil
		},
	}
	currService := service.NewService(mockRepo, &verification.EmailVerifier{}, &mocks.MockProofGenerator{}, code.NewMIMCCodeGenerator())

	req := dto.LoginClientDTO{
		Username:    username,
		CNonce:      cNonce,
		Nonce:       nonce,
		ClientProof: "TEST_CLIENT_PROOF_FOR_DECODE_ERROR_!@#", // Invalid client proof
	}

	loginClientResp, err := currService.LoginClient(req)

	assert.NotNil(t, err)
	assert.Equal(t, "error decoding client proof: encoding/hex: invalid byte: U+0054 'T'", err.Error())
	assert.Empty(t, loginClientResp)
}

func TestLoginClient_XorOperationError(t *testing.T) {
	mockRepo := &mockRepository{
		profileClient: func(username string) (models.Client, error) {
			return models.Client{
				StoredKey:      storedKey,
				Salt:           salt,
				IterationCount: iterationCount,
			}, nil
		},
	}
	currService := service.NewService(mockRepo, &verification.EmailVerifier{}, &mocks.MockProofGenerator{}, code.NewMIMCCodeGenerator())

	req := dto.LoginClientDTO{
		Username:    username,
		CNonce:      cNonce,
		Nonce:       nonce,
		ClientProof: "", // Sending empty client proof to fail the XOR operation, since it requires equal length slices
	}

	loginClientResp, err := currService.LoginClient(req)

	assert.NotNil(t, err)
	assert.Equal(t, "error performing XOR operation: slices must have the same length", err.Error())
	assert.Empty(t, loginClientResp)
}

func TestLoginClient_AuthFailedKeyMismatchError(t *testing.T) {
	mockRepo := &mockRepository{
		profileClient: func(username string) (models.Client, error) {
			return models.Client{
				StoredKey:      "", // Empty stored key
				Salt:           salt,
				IterationCount: iterationCount,
			}, nil
		},
	}
	currService := service.NewService(mockRepo, &verification.EmailVerifier{}, &mocks.MockProofGenerator{}, code.NewMIMCCodeGenerator())

	req := dto.LoginClientDTO{
		Username:    username,
		CNonce:      cNonce,
		Nonce:       nonce,
		ClientProof: clientProof,
	}

	loginClientResp, err := currService.LoginClient(req)

	assert.NotNil(t, err)
	assert.Equal(t, "server failed to authenticate the user", err.Error())
	assert.Empty(t, loginClientResp)
}

func TestLoginClient_DecodingServerKeyError(t *testing.T) {
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
	currService := service.NewService(mockRepo, &verification.EmailVerifier{}, &mocks.MockProofGenerator{}, code.NewMIMCCodeGenerator())

	req := dto.LoginClientDTO{
		Username:    username,
		CNonce:      cNonce,
		Nonce:       nonce,
		ClientProof: clientProof,
	}

	loginClientResp, err := currService.LoginClient(req)

	assert.NotNil(t, err)
	assert.Equal(t, "error decoding server key: encoding/hex: invalid byte: U+0054 'T'", err.Error())
	assert.Empty(t, loginClientResp)
}

func TestLoginClient_Success(t *testing.T) {
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
	currService := service.NewService(mockRepo, &verification.EmailVerifier{}, &mocks.MockProofGenerator{}, code.NewMIMCCodeGenerator())

	req := dto.LoginClientDTO{
		Username:    username,
		CNonce:      cNonce,
		Nonce:       nonce,
		ClientProof: clientProof,
	}

	loginClientResp, err := currService.LoginClient(req)

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

	currService := service.NewService(mockRepo, &verification.EmailVerifier{}, &mocks.MockProofGenerator{}, code.NewMIMCCodeGenerator())

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

	currService := service.NewService(mockRepo, &verification.EmailVerifier{}, &mocks.MockProofGenerator{}, code.NewMIMCCodeGenerator())

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

	currService := service.NewService(mockRepo, &verification.EmailVerifier{}, &mocks.MockProofGenerator{}, code.NewMIMCCodeGenerator())

	req := dto.RegisterClientPrecheckDTO{
		Username: username,
	}
	iterCount := 0

	salt, err := currService.RegisterClientPrecheck(req, iterCount)

	assert.Nil(t, err, "Expected no error during RegisterClientPrecheck")
	assert.NotEmpty(t, salt, "Salt should not be empty in the response")
}

func TestRegisterClient_RepositoryFailedToStoreUserData(t *testing.T) {
	mockRepo := &mockRepository{
		registerClient: func(req dto.RegisterClientDTO, id string, secret string) error {
			assert.Equal(t, username, req.Username)
			assert.Equal(t, clientName, req.Name)
			assert.Equal(t, redirectUri, req.RedirectURI)
			assert.Equal(t, backendUri, req.BackendURI)
			assert.Equal(t, storedKey, req.StoredKey)
			assert.Equal(t, serverKey, req.ServerKey)

			return fmt.Errorf("failed to store a client")
		},
	}
	currService := service.NewService(mockRepo, &verification.EmailVerifier{}, &mocks.MockProofGenerator{}, code.NewMIMCCodeGenerator())

	err := currService.RegisterClient(
		dto.RegisterClientDTO{
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

func TestRegisterClient_Success(t *testing.T) {
	mockRepo := &mockRepository{
		registerClient: func(req dto.RegisterClientDTO, id string, secret string) error {
			assert.Equal(t, username, req.Username)
			assert.Equal(t, clientName, req.Name)
			assert.Equal(t, redirectUri, req.RedirectURI)
			assert.Equal(t, backendUri, req.BackendURI)
			assert.Equal(t, storedKey, req.StoredKey)
			assert.Equal(t, serverKey, req.ServerKey)

			return nil
		},
	}
	currService := service.NewService(mockRepo, &verification.EmailVerifier{}, &mocks.MockProofGenerator{}, code.NewMIMCCodeGenerator())

	err := currService.RegisterClient(
		dto.RegisterClientDTO{
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
