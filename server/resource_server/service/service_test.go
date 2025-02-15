package service_test

import (
	"crypto/rand"
	"errors"
	"fmt"
	"globe-and-citizen/layer8/server/resource_server/dto"
	"globe-and-citizen/layer8/server/resource_server/emails/verification"
	"globe-and-citizen/layer8/server/resource_server/emails/verification/zk"
	"globe-and-citizen/layer8/server/resource_server/models"
	"globe-and-citizen/layer8/server/resource_server/service"
	"globe-and-citizen/layer8/server/resource_server/utils"
	"globe-and-citizen/layer8/server/resource_server/utils/mocks"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/crypto"

	"github.com/stretchr/testify/assert"
)

const userId uint = 1
const adminEmail = "admin@email.com"
const username = "user"
const password = "password"
const userEmail = "user@email.com"
const firstName = "first_name"
const lastName = "second_name"
const displayName = "display_name"
const country = "country"
const verificationCode = "123456"
const userSalt = "salt"
const publicKey = "0xaaaaa"
const serverKey = "0xbbbbb"
const storedKey = "0xccccc"

const redirectUri = "redirect_uri"
const backendUri = "backend_uri"

const zkKeyPairId uint = 2

const verificationCodeValidityDuration = 2 * time.Minute

var emailProof = []byte("proof")
var hashedPassword = utils.SaltAndHashPassword(password, userSalt)

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

func TestRegisterUser_RepositoryFailedToStoreUserData(t *testing.T) {
	mockRepo := &mocks.MockRepository{
		RegisterUserMock: func(req dto.RegisterUserDTO, hashedPassword string, salt string) error {
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
	mockRepo := &mocks.MockRepository{
		RegisterUserMock: func(req dto.RegisterUserDTO, hashedPassword string, salt string) error {
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
	mockRepo := &mocks.MockRepository{
		RegisterClientMock: func(client models.Client) error {
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
	mockRepo := &mocks.MockRepository{
		RegisterClientMock: func(client models.Client) error {
			assert.Equal(t, firstName, client.Name)
			assert.Equal(t, username, client.Username)
			assert.Equal(t, redirectUri, client.RedirectURI)
			assert.Equal(t, backendUri, client.BackendURI)

			return nil
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

	assert.Nil(t, err)
}

func TestLoginPreCheckUser(t *testing.T) {
	// Create a new mock repository
	mockRepo := new(mocks.MockRepository)

	// Create a new service by passing the mock repository
	mockService := service.NewService(mockRepo, &verification.EmailVerifier{}, &mocks.MockProofGenerator{})

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
	mockRepo := new(mocks.MockRepository)

	// Create a new service by passing the mock repository
	mockService := service.NewService(mockRepo, &verification.EmailVerifier{}, &mocks.MockProofGenerator{})

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
	mockRepo := new(mocks.MockRepository)

	// Create a new service by passing the mock repository
	mockService := service.NewService(mockRepo, &verification.EmailVerifier{}, &mocks.MockProofGenerator{})

	// Call the ProfileUser method of the mock service
	userDetails, err := mockService.ProfileUser(1)
	if err != nil {
		t.Error("Expected nil, got", err)
	}

	// Use assert to check if the error is nil
	assert.Nil(t, err)
	assert.Equal(t, userDetails.Username, "test_user")
}

func TestUpdateDisplayName(t *testing.T) {
	// Create a new mock repository
	mockRepo := new(mocks.MockRepository)

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
	mockRepo := new(mocks.MockRepository)

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
	mockRepo := new(mocks.MockRepository)

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
	mockRepo := &mocks.MockRepository{
		FindUserMock: func(userId uint) (models.User, error) {
			return models.User{}, fmt.Errorf("user %d does not exist", userId)
		},
	}
	emailVerifier := &verification.EmailVerifier{}

	currService := service.NewService(mockRepo, emailVerifier, &mocks.MockProofGenerator{})
	e := currService.VerifyEmail(userId, userEmail)

	assert.NotNil(t, e)
}

func TestVerifyEmail_UserExists_EmailFailedToBeSent(t *testing.T) {
	mockRepo := &mocks.MockRepository{
		FindUserMock: func(userId uint) (models.User, error) {
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
	mockRepo := &mocks.MockRepository{
		FindUserMock: func(userId uint) (models.User, error) {
			return models.User{
				ID:               userId,
				Username:         username,
				VerificationCode: "",
			}, nil
		},
		SaveEmailVerificationDataMock: func(data models.EmailVerificationData) error {
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
	mockRepo := &mocks.MockRepository{
		FindUserMock: func(userId uint) (models.User, error) {
			return models.User{
				ID:               userId,
				Username:         username,
				VerificationCode: "",
			}, nil
		},
		SaveEmailVerificationDataMock: func(data models.EmailVerificationData) error {
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
	mockRepo := &mocks.MockRepository{
		GetEmailVerificationDataMock: func(userId uint) (models.EmailVerificationData, error) {
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
	mockRepo := &mocks.MockRepository{
		GetEmailVerificationDataMock: func(userId uint) (models.EmailVerificationData, error) {
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
	mockRepo := &mocks.MockRepository{
		GetEmailVerificationDataMock: func(userId uint) (models.EmailVerificationData, error) {
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
	mockRepo := &mocks.MockRepository{
		GetEmailVerificationDataMock: func(userId uint) (models.EmailVerificationData, error) {
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
	mockRepo := &mocks.MockRepository{
		SaveProofOfEmailVerificationMock: func(userID uint, verificationCode string, proof []byte, zkKeyPairId uint) error {
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
	mockRepo := &mocks.MockRepository{
		SaveProofOfEmailVerificationMock: func(userID uint, verificationCode string, proof []byte, zkKeyPairId uint) error {
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
	mockRepo := &mocks.MockRepository{
		FindUserMock: func(userId uint) (models.User, error) {
			return models.User{}, fmt.Errorf("user not found for id %d", userId)
		},
	}

	currService := service.NewService(mockRepo, &verification.EmailVerifier{}, &mocks.MockProofGenerator{})
	_, err := currService.FindUser(userId)

	assert.NotNil(t, err)
}

func TestFindUser_Success(t *testing.T) {
	mockRepo := &mocks.MockRepository{
		FindUserMock: func(userId uint) (models.User, error) {
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
		Salt: userSalt,
	}
	request := dto.CheckEmailVerificationCodeDTO{
		Email: userEmail,
		Code:  verificationCode,
	}

	mockRepo := &mocks.MockRepository{}
	mockProofGenerator := &mocks.MockProofGenerator{
		GenerateProofFunc: func(
			emailAddress string, salt string, code string,
		) ([]byte, uint, error) {
			if emailAddress != userEmail {
				t.Fatalf("User's email mismatch: expected %s, got %s", userEmail, emailAddress)
			}
			if salt != userSalt {
				t.Fatalf("User's salt mimatch: expected %s, got %s", userSalt, salt)
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
		Salt: userSalt,
	}
	request := dto.CheckEmailVerificationCodeDTO{
		Email: userEmail,
		Code:  verificationCode,
	}

	mockRepo := &mocks.MockRepository{}
	mockProofGenerator := &mocks.MockProofGenerator{
		GenerateProofFunc: func(
			emailAddress string, salt string, code string,
		) ([]byte, uint, error) {
			if emailAddress != userEmail {
				t.Fatalf("User's email mismatch: expected %s, got %s", userEmail, emailAddress)
			}
			if salt != userSalt {
				t.Fatalf("User's salt mimatch: expected %s, got %s", userSalt, salt)
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
	mockRepo := &mocks.MockRepository{
		GetUserForUsernameMock: func(username string) (models.User, error) {
			return models.User{}, fmt.Errorf("user not found")
		},
	}

	currService := service.NewService(mockRepo, &verification.EmailVerifier{}, &mocks.MockProofGenerator{})
	_, err := currService.GetUserForUsername(username)

	assert.NotNil(t, err)
}

func TestGetUserForUsername_Success(t *testing.T) {
	mockRepo := &mocks.MockRepository{
		GetUserForUsernameMock: func(currUsername string) (models.User, error) {
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
	mockRepo := &mocks.MockRepository{
		UpdateUserPasswordMock: func(currUsername string, currPassword string) error {
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
	err := currService.UpdateUserPassword(username, password, userSalt)

	assert.NotNil(t, err)
}

func TestUpdateUserPassword_Success(t *testing.T) {
	mockRepo := &mocks.MockRepository{
		UpdateUserPasswordMock: func(currUsername string, currPassword string) error {
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
	err := currService.UpdateUserPassword(username, password, userSalt)

	assert.Nil(t, err)
}

func TestValidateSignature_SignatureIsInvalid(t *testing.T) {
	currService := service.NewService(
		&mocks.MockRepository{},
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
		&mocks.MockRepository{},
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
	mockRepo := &mocks.MockRepository{
		RegisterUserPrecheckMock: func(req dto.RegisterUserPrecheckDTO, rmSalt string, iterCount int) error {
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
	mockRepo := &mocks.MockRepository{
		RegisterUserPrecheckMock: func(req dto.RegisterUserPrecheckDTO, rmSalt string, iterCount int) error {
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
	mockRepo := &mocks.MockRepository{
		RegisterUserPrecheckMock: func(req dto.RegisterUserPrecheckDTO, rmSalt string, iterCount int) error {
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
	mockRepo := &mocks.MockRepository{
		RegisterUserv2Mock: func(req dto.RegisterUserDTOv2) error {
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
	mockRepo := &mocks.MockRepository{
		RegisterUserv2Mock: func(req dto.RegisterUserDTOv2) error {
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
	mockRepo := &mocks.MockRepository{
		UpdateUserPasswordV2Mock: func(currUsername, currStoredKey, currServerKey string) error {
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
	mockRepo := &mocks.MockRepository{
		UpdateUserPasswordV2Mock: func(currUsername, currStoredKey, currServerKey string) error {
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
