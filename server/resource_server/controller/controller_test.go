package controller_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"globe-and-citizen/layer8/server/resource_server/dto"
	"globe-and-citizen/layer8/server/resource_server/models"
	"globe-and-citizen/layer8/server/resource_server/utils"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	Ctl "globe-and-citizen/layer8/server/resource_server/controller"
)

const userId = 1
const username = "test_user"
const firstName = "first name"
const lastName = "last name"
const displayName = "display name"
const country = "country"
const verificationCode = "123467"
const userEmail = "user@email.com"
const userPassword = "test_password"
const newUserPassword = "new_password"
const userSalt = "ThisIsARandomSalt123!@#"
const clientSalt = "TestSaltForClient123!@#$%"
const iterationCount = 4096
const nonce = "Test_Nonce"
const serverSignature = "Test_Server_Signature"
const testToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImhtayIsInVzZXJfaWQiOjIsImlzcyI6Ikdsb2JlQW5kQ2l0aXplbiIsImV4cCI6MTcwNjUyNzY0NH0.AeQk23OPvlvauDEf45IlxxJ8ViSM5BlC6OlNkhXTomw"
const serverKey = "user_server_key"
const storedKey = "user_stored_key"

const zkKeyPairId uint = 2

var authenticationToken, _ = utils.GenerateToken(
	models.User{
		ID:       userId,
		Username: username,
	},
)
var emailProof = []byte("email_proof")

func decodeResponseBodyForResponse(t *testing.T, rr *httptest.ResponseRecorder) utils.Response {
	var response utils.Response
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}
	return response
}
func decodeResponseBodyForErrorResponse(t *testing.T, rr *httptest.ResponseRecorder) utils.Response {
	var response utils.Response

	body, err := io.ReadAll(rr.Body)
	if err != nil {
		t.Fatal(err)
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		t.Fatal(err)
	}

	return response
}

// MockService implements interfaces.IService for testing purposes.
type MockService struct {
	verifyEmail                        func(userID uint, userEmail string) error
	checkEmailVerificationCode         func(userID uint, code string) error
	findUser                           func(userID uint) (models.User, error)
	generateZkProofOfEmailVerification func(user models.User, request dto.CheckEmailVerificationCodeDTO) ([]byte, uint, error)
	saveProofOfEmailVerification       func(userID uint, verificationCode string, zkProof []byte, zkKeyPairId uint) error
	profileUser                        func(userID uint) (models.ProfileResponseOutput, error)
	getUserForUsername                 func(username string) (models.User, error)
	validateSignature                  func(message string, signature []byte, publicKey []byte) error
	updateUserPassword                 func(username string, newPassword string, salt string) error
	updateUserPasswordV2               func(username string, storedKey string, serverKey string) error
	registerUserPrecheck               func(req dto.RegisterUserPrecheckDTO, iterCount int) (string, error)
	registerClientPrecheck             func(req dto.RegisterClientPrecheckDTO, iterCount int) (string, error)
	registerUserv2                     func(req dto.RegisterUserDTOv2) error
	registerClientv2                   func(req dto.RegisterClientDTOv2) error
	loginPrecheckUserv2                func(req dto.LoginPrecheckDTOv2) (models.LoginPrecheckResponseOutputv2, error)
	loginPrecheckClientv2              func(req dto.LoginPrecheckDTOv2) (models.LoginPrecheckResponseOutputv2, error)
	loginUserv2                        func(req dto.LoginUserDTOv2) (models.LoginUserResponseOutputv2, error)
	loginPrecheckClientv2              func(req dto.LoginPrecheckDTOv2) (models.LoginPrecheckResponseOutputv2, error)
	loginClientv2                      func(req dto.LoginClientDTOv2) (models.LoginClientResponseOutputv2, error)
}

func (ms *MockService) RegisterUser(req dto.RegisterUserDTO) error {
	// Mock implementation for testing purposes.
	return nil
}

func (ms *MockService) LoginPreCheckUser(req dto.LoginPrecheckDTO) (models.LoginPrecheckResponseOutput, error) {
	// Mock implementation for testing purposes.
	return models.LoginPrecheckResponseOutput{
		Username: "test_user",
		Salt:     "ThisIsARandomSalt123!@#",
	}, nil
}

func (ms *MockService) LoginPrecheckUserv2(req dto.LoginPrecheckDTOv2) (response models.LoginPrecheckResponseOutputv2, err error) {
	return ms.loginPrecheckUserv2(req)
}

func (ms *MockService) LoginUser(req dto.LoginUserDTO) (models.LoginUserResponseOutput, error) {
	return models.LoginUserResponseOutput{
		Token: testToken,
	}, nil
}

func (ms *MockService) LoginUserv2(req dto.LoginUserDTOv2) (response models.LoginUserResponseOutputv2, err error) {
	return ms.loginUserv2(req)
}

func (ms *MockService) ProfileUser(userID uint) (models.ProfileResponseOutput, error) {
	return ms.profileUser(userID)
}

func (ms *MockService) FindUser(userID uint) (models.User, error) {
	return ms.findUser(userID)
}

func (ms *MockService) VerifyEmail(userID uint, userEmail string) error {
	return ms.verifyEmail(userID, userEmail)
}

func (ms *MockService) CheckEmailVerificationCode(userID uint, code string) error {
	return ms.checkEmailVerificationCode(userID, code)
}

func (ms *MockService) GenerateZkProofOfEmailVerification(
	user models.User,
	request dto.CheckEmailVerificationCodeDTO,
) ([]byte, uint, error) {
	return ms.generateZkProofOfEmailVerification(user, request)
}

func (ms *MockService) SaveProofOfEmailVerification(
	userID uint, verificationCode string, zkProof []byte, zkKeyPairId uint,
) error {
	return ms.saveProofOfEmailVerification(userID, verificationCode, zkProof, zkKeyPairId)
}

func (ms *MockService) UpdateDisplayName(userID uint, req dto.UpdateDisplayNameDTO) error {
	// Mock implementation for testing purposes.
	return nil
}

func (ms *MockService) RegisterClient(req dto.RegisterClientDTO) error {
	// Mock implementation for testing purposes.
	return nil
}

func (ms *MockService) RegisterClientv2(req dto.RegisterClientDTOv2) error {
	return ms.registerClientv2(req)
}

func (ms *MockService) RegisterClientPrecheck(req dto.RegisterClientPrecheckDTO, iterCount int) (string, error) {
	return ms.registerClientPrecheck(req, iterCount)
}

func (ms *MockService) GetClientData(clientName string) (models.ClientResponseOutput, error) {
	// Mock implementation for testing purposes.
	return models.ClientResponseOutput{
		ID:          "0",
		Secret:      "",
		Name:        "testclient",
		RedirectURI: "https://gcitizen.com/callback",
	}, nil
}

func (ms *MockService) LoginPreCheckClient(req dto.LoginPrecheckDTO) (models.LoginPrecheckResponseOutput, error) {
	// Mock implementation for testing purposes.
	return models.LoginPrecheckResponseOutput{}, nil
}

func (ms *MockService) ProfileClient(userID string) (models.ClientResponseOutput, error) {
	// Mock implementation for testing purposes.
	return models.ClientResponseOutput{}, nil
}

func (ms *MockService) GetClientDataByBackendURL(backendURL string) (models.ClientResponseOutput, error) {
	return models.ClientResponseOutput{}, nil
}

func (ms *MockService) CheckBackendURI(backendURL string) (bool, error) {
	// Mock implementation for testing purposes.
	return true, nil
}

func (m *MockService) LoginClient(req dto.LoginClientDTO) (models.LoginUserResponseOutput, error) {
	// Mock implementation for LoginClient method
	return models.LoginUserResponseOutput{
		Token: testToken,
	}, nil
}

func (m *MockService) LoginClientv2(req dto.LoginClientDTOv2) (models.LoginClientResponseOutputv2, error) {
	return m.loginClientv2(req)
}

func (m *MockService) LoginPrecheckClientv2(req dto.LoginPrecheckDTOv2) (models.LoginPrecheckResponseOutputv2, error) {
	return m.loginPrecheckClientv2(req)
}

func (m *MockService) LoginClientv2(req dto.LoginClientDTOv2) (models.LoginClientResponseOutputv2, error) {
	return m.loginClientv2(req)
}

func (m *MockService) LoginPrecheckClientv2(req dto.LoginPrecheckDTOv2) (models.LoginPrecheckResponseOutputv2, error) {
	return m.loginPrecheckClientv2(req)
}

func (m *MockService) GetUserForUsername(username string) (models.User, error) {
	return m.getUserForUsername(username)
}

func (m *MockService) ValidateSignature(message string, signature []byte, publicKey []byte) error {
	return m.validateSignature(message, signature, publicKey)
}

func (m *MockService) UpdateUserPassword(username string, newPassword string, salt string) error {
	return m.updateUserPassword(username, newPassword, salt)
}

func (m *MockService) UpdateUserPasswordV2(username string, storedKey string, serverKey string) error {
	return m.updateUserPasswordV2(username, storedKey, serverKey)
}

func (m *MockService) RegisterUserPrecheck(req dto.RegisterUserPrecheckDTO, iterCount int) (string, error) {
	return m.registerUserPrecheck(req, iterCount)
}

// Mock RegisterUser method for unit tests
func (m *MockService) RegisterUserv2(req dto.RegisterUserDTOv2) error {
	return m.registerUserv2(req)
}

func TestRegisterUserHandler_InvalidHttpRequestMethod(t *testing.T) {
	requestBody := []byte(`{
		"email": "test@gcitizen.com",
		"username": "test_user",
		"first_name": "Test",
		"last_name": "User",
		"display_name": "user",
		"country": "Unknown",
		"password": "12345"
	}`)

	req, err := http.NewRequest("GET", "/api/v1/register-user", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	setMockServiceInContext(req)

	rr := httptest.NewRecorder()

	Ctl.RegisterUserHandler(rr, req)

	assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)

	response := decodeResponseBodyForErrorResponse(t, rr)

	assert.False(t, response.IsSuccess)
	assert.Equal(t, "Invalid http method. Expected POST", response.Message)
	assert.NotNil(t, response.Error)
}

func TestRegisterUserHandler_RequestJsonIsMalformed(t *testing.T) {
	requestBody := []byte(`{
		"email": "test@gcitizen.com",
		"username": "test_user",
		"first_name": "Test",
		"last_name": "User",
		"display_name": "user",
		"country": "Unknown",
		"password": "12345"
	}something_else`)

	req, err := http.NewRequest("POST", "/api/v1/register-user", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	mockService := &MockService{}
	req = req.WithContext(context.WithValue(req.Context(), "service", mockService))
	// setMockServiceInContext(req)

	rr := httptest.NewRecorder()

	Ctl.RegisterUserHandler(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	response := decodeResponseBodyForErrorResponse(t, rr)

	assert.False(t, response.IsSuccess)
	assert.Equal(t, "Request malformed: error while parsing json", response.Message)
	assert.NotNil(t, response.Error)
}

func TestRegisterUserHandler_RequiredRequestJsonFieldsAreMissing(t *testing.T) {
	requestBody := []byte(`{
		"email": "test@gcitizen.com",
		"first_name": "Test",
		"last_name": "User",
		"display_name": "user",
		"country": "Unknown",
		"password": "12345"
	}`)

	req, err := http.NewRequest("POST", "/api/v1/register-user", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	mockService := &MockService{}
	req = req.WithContext(context.WithValue(req.Context(), "service", mockService))

	rr := httptest.NewRecorder()

	Ctl.RegisterUserHandler(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	response := decodeResponseBodyForErrorResponse(t, rr)

	assert.False(t, response.IsSuccess)
	assert.Equal(t, "Input json is invalid", response.Message)
	assert.NotNil(t, response.Error)
}

func TestRegisterUserHandler_Success(t *testing.T) {
	// Mock request body
	requestBody := []byte(`{
		"username": "test_user",
		"first_name": "Test",
		"last_name": "User",
		"display_name": "user",
		"country": "Unknown",
		"password": "12345",
		"public_key": "0xaaaaaa"
	}`)

	// Create a mock request
	req, err := http.NewRequest("POST", "/api/v1/register-user", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	// Create a mock service and set it in the request context
	mockService := &MockService{}
	req = req.WithContext(context.WithValue(req.Context(), "service", mockService))

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the handler function
	Ctl.RegisterUserHandler(rr, req)

	// Check the status code
	assert.Equal(t, http.StatusCreated, rr.Code)

	// Decode the response body
	response := decodeResponseBodyForResponse(t, rr)

	// Now assert the fields directly
	assert.True(t, response.IsSuccess)
	assert.Equal(t, "User registered successfully", response.Message)
	assert.Equal(t, nil, response.Data)
}

func TestRegisterClientHandler_InvalidHttpRequestMethod(t *testing.T) {
	requestBody := []byte(`{
		"name": "testclient", 
		"redirect_uri": "https://gcitizen.com/callback", 
		"username": "test_user", 
		"password": "12345"
	}`)

	req, err := http.NewRequest("PUT", "/api/v1/register-client", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	setMockServiceInContext(req)

	rr := httptest.NewRecorder()

	Ctl.RegisterClientHandler(rr, req)

	assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)

	response := decodeResponseBodyForErrorResponse(t, rr)

	assert.False(t, response.IsSuccess)
	assert.Equal(t, "Invalid http method. Expected POST", response.Message)
	assert.NotNil(t, response.Error)
}

func TestRegisterClientHandler_RequestJsonIsMalformed(t *testing.T) {
	requestBody := []byte(`{
		"name": "testclient", 
		"redirect_uri": "https://gcitizen.com/callback",
		"backend_uri": "https://backend.com",
		"username": "test_user", 
		"password": "12345"
	}something_else`)

	req, err := http.NewRequest("POST", "/api/v1/register-client", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	mockService := &MockService{}
	req = req.WithContext(context.WithValue(req.Context(), "service", mockService))

	rr := httptest.NewRecorder()

	Ctl.RegisterClientHandler(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	response := decodeResponseBodyForErrorResponse(t, rr)

	assert.False(t, response.IsSuccess)
	assert.Equal(t, "Request malformed: error while parsing json", response.Message)
	assert.NotNil(t, response.Error)
}

func TestRegisterClientHandler_RequiredRequestJsonFieldsAreMissing(t *testing.T) {
	requestBody := []byte(`{
		"name": "testclient", 
		"redirect_uri": "https://gcitizen.com/callback",
		"backend_uri": "https://backend.com",
		"username": "test_user"
	}`)

	req, err := http.NewRequest("POST", "/api/v1/register-client", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	mockService := &MockService{}
	req = req.WithContext(context.WithValue(req.Context(), "service", mockService))

	rr := httptest.NewRecorder()

	Ctl.RegisterClientHandler(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	response := decodeResponseBodyForErrorResponse(t, rr)

	assert.False(t, response.IsSuccess)
	assert.Equal(t, "Input json is invalid", response.Message)
	assert.NotNil(t, response.Error)
}

func TestRegisterClientHandler_Success(t *testing.T) {
	// Mock request body
	requestBody := []byte(`{
		"name": "testclient", 
		"redirect_uri": "https://gcitizen.com/callback",
		"backend_uri": "https://backend.com",
		"username": "test_user", 
		"password": "12345"
	}`)

	// Create a mock request
	req, err := http.NewRequest("POST", "/api/v1/register-client", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	// Create a mock service and set it in the request context
	mockService := &MockService{}
	req = req.WithContext(context.WithValue(req.Context(), "service", mockService))

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the handler function
	Ctl.RegisterClientHandler(rr, req)

	// Check the status code
	assert.Equal(t, http.StatusCreated, rr.Code)

	// Decode the response body
	var response utils.Response
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}

	// Now assert the fields directly
	assert.True(t, response.IsSuccess)
	assert.Equal(t, "Client registered successfully", response.Message)
	assert.Equal(t, nil, response.Data)
}

func TestLoginPrecheckHandler_InvalidHttpRequestMethod(t *testing.T) {
	requestBody := []byte(`{"username": "test_user"}`)

	req, err := http.NewRequest("GET", "/api/v1/login-precheck", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	mockService := &MockService{}
	req = req.WithContext(context.WithValue(req.Context(), "service", mockService))

	rr := httptest.NewRecorder()

	Ctl.LoginPrecheckHandler(rr, req)

	assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)

	response := decodeResponseBodyForErrorResponse(t, rr)

	assert.False(t, response.IsSuccess)
	assert.Equal(t, "Invalid http method. Expected POST", response.Message)
	assert.NotNil(t, response.Error)
}

func TestLoginPrecheckHandler_RequestJsonIsMalformed(t *testing.T) {
	requestBody := []byte(`{"username": "test_user"}something_else`)
	req, err := http.NewRequest("POST", "/api/v1/login-precheck", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	mockService := &MockService{}
	req = req.WithContext(context.WithValue(req.Context(), "service", mockService))

	rr := httptest.NewRecorder()

	Ctl.LoginPrecheckHandler(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	response := decodeResponseBodyForErrorResponse(t, rr)

	assert.False(t, response.IsSuccess)
	assert.Equal(t, "Request malformed: error while parsing json", response.Message)
	assert.NotNil(t, response.Error)
}

func TestLoginPrecheckHandler_RequiredRequestJsonFieldsAreMissing(t *testing.T) {
	requestBody := []byte(`{}`)

	req, err := http.NewRequest("POST", "/api/v1/login-precheck", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	mockService := &MockService{}
	req = req.WithContext(context.WithValue(req.Context(), "service", mockService))

	rr := httptest.NewRecorder()

	Ctl.LoginPrecheckHandler(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	response := decodeResponseBodyForErrorResponse(t, rr)

	assert.False(t, response.IsSuccess)
	assert.Equal(t, "Input json is invalid", response.Message)
	assert.NotNil(t, response.Error)
}

func TestLoginPrecheckHandler_Success(t *testing.T) {
	// Mock request body
	requestBody := []byte(`{"username": "test_user"}`)

	// Create a mock request
	req, err := http.NewRequest("POST", "/api/v1/login-precheck", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	// Create a mock service and set it in the request context
	mockService := &MockService{}
	req = req.WithContext(context.WithValue(req.Context(), "service", mockService))

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the handler function
	Ctl.LoginPrecheckHandler(rr, req)

	// Check the status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Decode the response body
	var response models.LoginPrecheckResponseOutput
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}

	// Now assert the fields directly
	assert.Equal(t, "test_user", response.Username)
	assert.Equal(t, "ThisIsARandomSalt123!@#", response.Salt)
}

func TestLoginPrecheckHandlerv2_InvalidHttpRequestMethod(t *testing.T) {
	requestBody := []byte(`{"username": "test_user", "c_nonce": "Test_Nonce"}`)

	req, err := http.NewRequest("GET", "/api/v2/login-precheck", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	mockService := &MockService{}

	req = req.WithContext(context.WithValue(req.Context(), "service", mockService))

	rr := httptest.NewRecorder()

	Ctl.LoginPrecheckHandlerv2(rr, req)

	assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)

	response := decodeResponseBodyForErrorResponse(t, rr)

	assert.False(t, response.IsSuccess)
	assert.Equal(t, "Invalid http method. Expected POST", response.Message)
	assert.NotNil(t, response.Error)
}

func TestLoginPrecheckHandlerv2_RequestJsonIsMalformed(t *testing.T) {
	requestBody := []byte(`{"username": "test_user", "c_nonce": "Test_Nonce"}something_else`)
	req, err := http.NewRequest("POST", "/api/v2/login-precheck", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	mockService := &MockService{}

	req = req.WithContext(context.WithValue(req.Context(), "service", mockService))

	rr := httptest.NewRecorder()

	Ctl.LoginPrecheckHandlerv2(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	response := decodeResponseBodyForErrorResponse(t, rr)

	assert.False(t, response.IsSuccess)
	assert.Equal(t, "Request malformed: error while parsing json", response.Message)
	assert.NotNil(t, response.Error)
}

func TestLoginPrecheckHandlerv2_RequiredRequestJsonFieldsAreMissing(t *testing.T) {
	requestBody := []byte(`{}`)

	req, err := http.NewRequest("POST", "/api/v2/login-precheck", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	mockService := &MockService{}

	req = req.WithContext(context.WithValue(req.Context(), "service", mockService))

	rr := httptest.NewRecorder()

	Ctl.LoginPrecheckHandlerv2(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	response := decodeResponseBodyForErrorResponse(t, rr)

	assert.False(t, response.IsSuccess)
	assert.Equal(t, "Input json is invalid", response.Message)
	assert.NotNil(t, response.Error)
}

func TestLoginPrecheckHandlerv2_ServiceError(t *testing.T) {
	// Mock request body
	requestBody := []byte(`{"username": "test_user", "c_nonce": "Test_Nonce"}`)

	// Create a mock request
	req, err := http.NewRequest("POST", "/api/v2/login-precheck", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	// Create a mock service and set it in the request context
	mockService := &MockService{
		loginPrecheckUserv2: func(req dto.LoginPrecheckDTOv2) (models.LoginPrecheckResponseOutputv2, error) {
			return models.LoginPrecheckResponseOutputv2{}, fmt.Errorf("mock service error")
		},
	}

	req = req.WithContext(context.WithValue(req.Context(), "service", mockService))

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the handler function
	Ctl.LoginPrecheckHandlerv2(rr, req)

	response := decodeResponseBodyForErrorResponse(t, rr)

	// Now assert the fields directly
	assert.False(t, response.IsSuccess)
	assert.Equal(t, "Failed to perform precheck, service error", response.Message)
	assert.NotNil(t, response.Error)

	// Check the status code
	assert.Equal(t, http.StatusBadRequest, rr.Code, "Expected HTTP 400 Bad Request")
}

func TestLoginPrecheckHandlerv2_Success(t *testing.T) {
	// Mock request body
	requestBody := []byte(`{"username": "test_user", "c_nonce": "Test_Nonce"}`)

	// Create a mock request
	req, err := http.NewRequest("POST", "/api/v2/login-precheck", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	// Create a mock service and set it in the request context
	mockService := &MockService{
		loginPrecheckUserv2: func(req dto.LoginPrecheckDTOv2) (models.LoginPrecheckResponseOutputv2, error) {
			return models.LoginPrecheckResponseOutputv2{
				Salt:      userSalt,
				IterCount: iterationCount,
				Nonce:     nonce,
			}, nil
		},
	}

	req = req.WithContext(context.WithValue(req.Context(), "service", mockService))

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the handler function
	Ctl.LoginPrecheckHandlerv2(rr, req)

	// Check the status code
	assert.Equal(t, http.StatusOK, rr.Code)

	response := decodeResponseBodyForResponse(t, rr)

	// Convert response.Data to JSON bytes for unmarshalling
	dataBytes, err := json.Marshal(response.Data)
	if err != nil {
		t.Fatal(err)
	}

	var loginPrecheckResponse models.LoginPrecheckResponseOutputv2
	if err := json.Unmarshal(dataBytes, &loginPrecheckResponse); err != nil {
		t.Fatal(err)
	}

	assert.True(t, response.IsSuccess)
	assert.Nil(t, response.Error)

	// Now assert the fields directly
	assert.Equal(t, userSalt, loginPrecheckResponse.Salt)
	assert.Equal(t, iterationCount, loginPrecheckResponse.IterCount)
	assert.Equal(t, nonce, loginPrecheckResponse.Nonce)
}

func TestLoginUserHandler_InvalidHttpRequestMethod(t *testing.T) {
	requestBody := []byte(`{
		"username": "test_user",
		"password": "12345",
		"salt": 	"ThisIsARandomSalt123!@#"}`)

	req, err := http.NewRequest("GET", "/api/v1/login-user", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	setMockServiceInContext(req)

	rr := httptest.NewRecorder()

	Ctl.LoginUserHandler(rr, req)

	assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)

	response := decodeResponseBodyForErrorResponse(t, rr)

	assert.False(t, response.IsSuccess)
	assert.Equal(t, "Invalid http method. Expected POST", response.Message)
	assert.NotNil(t, response.Error)
}

func TestLoginUserHandler_RequestJsonIsMalformed(t *testing.T) {
	requestBody := []byte(`{
		"username": "test_user",
		"password": "12345",
		"salt": 	"ThisIsARandomSalt123!@#"}something_else`)
	req, err := http.NewRequest("POST", "/api/v1/login-user", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	mockService := &MockService{}
	req = req.WithContext(context.WithValue(req.Context(), "service", mockService))

	rr := httptest.NewRecorder()

	Ctl.LoginUserHandler(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	response := decodeResponseBodyForErrorResponse(t, rr)

	assert.False(t, response.IsSuccess)
	assert.Equal(t, "Request malformed: error while parsing json", response.Message)
	assert.NotNil(t, response.Error)
}

func TestLoginUserHandler_RequiredRequestJsonFieldsAreMissing(t *testing.T) {
	requestBody := []byte(`{
		"username": "test_user",
		"password": "12345"
	}`)
	req, err := http.NewRequest("POST", "/api/v1/login-user", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	mockService := &MockService{}
	req = req.WithContext(context.WithValue(req.Context(), "service", mockService))

	rr := httptest.NewRecorder()

	Ctl.LoginUserHandler(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	response := decodeResponseBodyForErrorResponse(t, rr)

	assert.False(t, response.IsSuccess)
	assert.Equal(t, "Input json is invalid", response.Message)
	assert.NotNil(t, response.Error)
}

func TestLoginUserHandler_Success(t *testing.T) {
	// Mock request body
	requestBody := []byte(`{
		"username": "test_user",
		"password": "12345",
		"salt": 	"ThisIsARandomSalt123!@#"}`)

	// Create a mock request
	req, err := http.NewRequest("POST", "/api/v1/login-user", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	// Create a mock service and set it in the request context
	mockService := &MockService{}
	req = req.WithContext(context.WithValue(req.Context(), "service", mockService))

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the handler function
	Ctl.LoginUserHandler(rr, req)

	// Check the status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Decode the response body
	var response models.LoginUserResponseOutput
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}

	// Now assert the fields directly
	assert.Equal(t, testToken, response.Token)
}

func TestLoginUserHandlerv2_InvalidHttpRequestMethod(t *testing.T) {
	requestBody := []byte(`{
		"username": 	"test_user",
		"nonce": 		"Test_Nonce",
		"c_nonce": 		"Test_Nonce",
		"client_proof": "Test_Client_Proof"
		}`)

	req, err := http.NewRequest("GET", "/api/v2/login-user", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	mockService := &MockService{}

	req = req.WithContext(context.WithValue(req.Context(), "service", mockService))

	rr := httptest.NewRecorder()

	Ctl.LoginUserHandlerv2(rr, req)

	assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)

	response := decodeResponseBodyForErrorResponse(t, rr)

	assert.False(t, response.IsSuccess)
	assert.Equal(t, "Invalid http method. Expected POST", response.Message)
	assert.NotNil(t, response.Error)
}

func TestLoginUserHandlerv2_RequestJsonIsMalformed(t *testing.T) {
	requestBody := []byte(`{
		"username": 	"test_user",
		"nonce": 		"Test_Nonce",
		"c_nonce": 		"Test_Nonce",
		"client_proof": "Test_Client_Proof"
		}something_else`)

	req, err := http.NewRequest("POST", "/api/v2/login-user", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	mockService := &MockService{}

	req = req.WithContext(context.WithValue(req.Context(), "service", mockService))

	rr := httptest.NewRecorder()

	Ctl.LoginUserHandlerv2(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	response := decodeResponseBodyForErrorResponse(t, rr)

	assert.False(t, response.IsSuccess)
	assert.Equal(t, "Request malformed: error while parsing json", response.Message)
	assert.NotNil(t, response.Error)
}

func TestLoginUserHandlerv2_RequiredRequestJsonFieldsAreMissing(t *testing.T) {
	requestBody := []byte(`{
		"nonce": 		"Test_Nonce",
		"client_proof": "Test_Client_Proof"
		}`)

	req, err := http.NewRequest("POST", "/api/v2/login-user", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	mockService := &MockService{}

	req = req.WithContext(context.WithValue(req.Context(), "service", mockService))

	rr := httptest.NewRecorder()

	Ctl.LoginUserHandlerv2(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	response := decodeResponseBodyForErrorResponse(t, rr)

	assert.False(t, response.IsSuccess)
	assert.Equal(t, "Input json is invalid", response.Message)
	assert.NotNil(t, response.Error)
}

func TestLoginUserHandlerv2_ServiceError(t *testing.T) {
	// Mock request body
	requestBody := []byte(`{
		"username": 	"test_user",
		"nonce": 		"Test_Nonce",
		"c_nonce": 		"Test_Nonce",
		"client_proof": "Test_Client_Proof"
		}`)

	// Create a mock request
	req, err := http.NewRequest("POST", "/api/v2/login-user", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	// Create a mock service and set it in the request context
	mockService := &MockService{
		loginUserv2: func(req dto.LoginUserDTOv2) (models.LoginUserResponseOutputv2, error) {
			return models.LoginUserResponseOutputv2{}, fmt.Errorf("mock service error")
		},
	}

	req = req.WithContext(context.WithValue(req.Context(), "service", mockService))

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the handler function
	Ctl.LoginUserHandlerv2(rr, req)

	// Decode the response body
	response := decodeResponseBodyForErrorResponse(t, rr)

	// Now assert the fields directly
	assert.False(t, response.IsSuccess)
	assert.Equal(t, "Failed to perform login", response.Message)
	assert.NotNil(t, response.Error)

	// Check the status code
	assert.Equal(t, http.StatusBadRequest, rr.Code, "Expected HTTP 400 Bad Request")
}

func TestLoginUserHandlerv2_Success(t *testing.T) {
	// Mock request body
	requestBody := []byte(`{
		"username": 	"test_user",
		"nonce": 		"Test_Nonce",
		"c_nonce": 		"Test_Nonce",
		"client_proof": "Test_Client_Proof"
		}`)

	// Create a mock request
	req, err := http.NewRequest("POST", "/api/v2/login-user", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	// Create a mock service and set it in the request context
	mockService := &MockService{
		loginUserv2: func(req dto.LoginUserDTOv2) (models.LoginUserResponseOutputv2, error) {
			return models.LoginUserResponseOutputv2{
				Token:           testToken,
				ServerSignature: serverSignature,
			}, nil
		},
	}

	req = req.WithContext(context.WithValue(req.Context(), "service", mockService))

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the handler function
	Ctl.LoginUserHandlerv2(rr, req)

	// Check the status code
	assert.Equal(t, http.StatusOK, rr.Code)

	response := decodeResponseBodyForErrorResponse(t, rr)

	// Convert response.Data to JSON bytes for unmarshalling
	dataBytes, err := json.Marshal(response.Data)
	if err != nil {
		t.Fatal(err)
	}

	var loginUserResponse models.LoginUserResponseOutputv2
	if err := json.Unmarshal(dataBytes, &loginUserResponse); err != nil {
		t.Fatal(err)
	}

	assert.True(t, response.IsSuccess)
	assert.Nil(t, response.Error)
	assert.Equal(t, "Login successful", response.Message)

	// Now assert the fields directly
	assert.Equal(t, testToken, loginUserResponse.Token)
	assert.Equal(t, serverSignature, loginUserResponse.ServerSignature)
}

func TestProfileHandler_InvalidHttpRequestMethod(t *testing.T) {
	req, err := http.NewRequest("POST", "/api/v1/profile", nil)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Authorization", "Bearer "+authenticationToken)
	setMockServiceInContext(req)

	rr := httptest.NewRecorder()

	Ctl.ProfileHandler(rr, req)

	assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)

	response := decodeResponseBodyForErrorResponse(t, rr)

	assert.False(t, response.IsSuccess)
	assert.Equal(t, "Invalid http method. Expected GET", response.Message)
	assert.NotNil(t, response.Error)
}

func TestProfileHandler_InvalidAuthenticationToken(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/v1/profile", nil)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Authorization", "Bearer invalid token")

	mockService := &MockService{}
	req = req.WithContext(context.WithValue(req.Context(), "service", mockService))

	rr := httptest.NewRecorder()

	Ctl.ProfileHandler(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)

	response := decodeResponseBodyForErrorResponse(t, rr)

	assert.False(t, response.IsSuccess)
	assert.Equal(t, "Authentication error: invalid token", response.Message)
	assert.NotNil(t, response.Error)
}

func TestProfileHandler_FailedToProfileUser(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/v1/profile", nil)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Authorization", "Bearer "+authenticationToken)

	mockService := &MockService{
		profileUser: func(userID uint) (models.ProfileResponseOutput, error) {
			return models.ProfileResponseOutput{}, fmt.Errorf("could not profile user %d", userID)
		},
	}
	req = req.WithContext(context.WithValue(req.Context(), "service", mockService))

	rr := httptest.NewRecorder()

	Ctl.ProfileHandler(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	response := decodeResponseBodyForErrorResponse(t, rr)

	assert.False(t, response.IsSuccess)
	assert.Equal(t, "Failed to get user profile", response.Message)
	assert.NotNil(t, response.Error)
}

func TestProfileHandler_Success(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/v1/profile", nil)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Authorization", "Bearer "+authenticationToken)

	mockService := &MockService{
		profileUser: func(userID uint) (models.ProfileResponseOutput, error) {
			return models.ProfileResponseOutput{
				Username:    username,
				FirstName:   firstName,
				LastName:    lastName,
				DisplayName: displayName,
				Country:     country,
			}, nil
		},
	}
	req = req.WithContext(context.WithValue(req.Context(), "service", mockService))

	rr := httptest.NewRecorder()

	Ctl.ProfileHandler(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response models.ProfileResponseOutput
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, username, response.Username)
	assert.Equal(t, firstName, response.FirstName)
	assert.Equal(t, lastName, response.LastName)
	assert.Equal(t, displayName, response.DisplayName)
	assert.Equal(t, country, response.Country)
}

func TestGetClientData_Success(t *testing.T) {
	// Create a mock request
	req, err := http.NewRequest("GET", "/api/v1/client", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Set the Authorization header
	req.Header.Set("Name", "testclient")

	// Create a mock service and set it in the request context
	mockService := &MockService{}
	req = req.WithContext(context.WithValue(req.Context(), "service", mockService))

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the handler function
	Ctl.GetClientData(rr, req)

	// Check the status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Decode the response body
	var response models.ClientResponseOutput
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}

	// Now assert the fields directly
	assert.Equal(t, "0", response.ID)
	assert.Equal(t, "", response.Secret)
	assert.Equal(t, "testclient", response.Name)
	assert.Equal(t, "https://gcitizen.com/callback", response.RedirectURI)
}

func TestVerifyEmailHandler_InvalidHttpRequestMethod(t *testing.T) {
	requestBody := []byte(`{"email": "user@email.com"}`)
	req, err := http.NewRequest("GET", "/api/v1/verify-email", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	Ctl.VerifyEmailHandler(rr, req)

	assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)

	response := decodeResponseBodyForErrorResponse(t, rr)

	assert.False(t, response.IsSuccess)
	assert.Equal(t, "Invalid http method. Expected POST", response.Message)
	assert.NotNil(t, response.Error)
}

func TestVerifyEmailHandler_InvalidAuthorizationToken(t *testing.T) {
	requestBody := []byte(`{"email": "user@email.com"}`)
	req, err := http.NewRequest("POST", "/api/v1/verify-email", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer invalid token")

	req = req.WithContext(context.WithValue(req.Context(), "service", &MockService{}))

	rr := httptest.NewRecorder()

	Ctl.VerifyEmailHandler(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)

	response := decodeResponseBodyForErrorResponse(t, rr)

	assert.False(t, response.IsSuccess)
	assert.Equal(t, "Authentication error: invalid token", response.Message)
	assert.NotNil(t, response.Error)
}

func TestVerifyEmailHandler_MalformedRequestBodyJson(t *testing.T) {
	requestBody := []byte(`{"email": "user@email.com";}`)
	req, err := http.NewRequest("POST", "/api/v1/verify-email", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+authenticationToken)

	req = req.WithContext(context.WithValue(req.Context(), "service", &MockService{}))

	rr := httptest.NewRecorder()

	Ctl.VerifyEmailHandler(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	response := decodeResponseBodyForErrorResponse(t, rr)

	assert.False(t, response.IsSuccess)
	assert.Equal(t, "Request malformed: error while parsing json", response.Message)
	assert.NotNil(t, response.Error)
}

func TestVerifyEmailHandler_RequestJsonSchemeIsInvalid(t *testing.T) {
	requestBody := []byte(`{"emal": "user@email.com"}`)
	req, err := http.NewRequest("POST", "/api/v1/verify-email", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+authenticationToken)

	req = req.WithContext(context.WithValue(req.Context(), "service", &MockService{}))

	rr := httptest.NewRecorder()

	Ctl.VerifyEmailHandler(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	response := decodeResponseBodyForErrorResponse(t, rr)

	assert.False(t, response.IsSuccess)
	assert.Equal(t, "Input json is invalid", response.Message)
	assert.NotNil(t, response.Error)
}

func TestVerifyEmailHandler_FailedToVerifyEmail(t *testing.T) {
	requestBody := []byte(`{"email": "user@email.com"}`)
	req, err := http.NewRequest("POST", "/api/v1/verify-email", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+authenticationToken)

	mockService := &MockService{
		verifyEmail: func(userID uint, email string) error {
			if email != userEmail {
				t.Fatalf("User email mismatch: expected %s, got %s", userEmail, email)
			}
			return fmt.Errorf("failed to verify email for user %d", userID)
		},
	}
	req = req.WithContext(context.WithValue(req.Context(), "service", mockService))

	rr := httptest.NewRecorder()

	Ctl.VerifyEmailHandler(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	response := decodeResponseBodyForErrorResponse(t, rr)

	assert.False(t, response.IsSuccess)
	assert.Equal(t, "Failed to verify email", response.Message)
	assert.NotNil(t, response.Error)
}

func TestVerifyEmailHandler_Success(t *testing.T) {
	requestBody := []byte(`{"email": "user@email.com"}`)
	req, err := http.NewRequest("POST", "/api/v1/verify-email", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+authenticationToken)

	mockService := &MockService{
		verifyEmail: func(userID uint, email string) error {
			if email != userEmail {
				t.Fatalf("User email mismatch: expected %s, got %s", userEmail, email)
			}
			return nil
		},
	}
	req = req.WithContext(context.WithValue(req.Context(), "service", mockService))

	rr := httptest.NewRecorder()

	Ctl.VerifyEmailHandler(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	response := decodeResponseBodyForResponse(t, rr)

	assert.True(t, response.IsSuccess)
	assert.Equal(t, "Verification email sent", response.Message)
	assert.Equal(t, nil, response.Data)
}

func TestCheckEmailVerificationCode_InvalidHttpRequestMethod(t *testing.T) {
	requestBody := []byte(`{
		"email": "user@email.com",
		"code": "123467"
	}`)
	req, err := http.NewRequest("PUT", "/api/v1/check-email-verification-code", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	Ctl.CheckEmailVerificationCode(rr, req)

	assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)

	response := decodeResponseBodyForErrorResponse(t, rr)

	assert.False(t, response.IsSuccess)
	assert.Equal(t, "Invalid http method. Expected POST", response.Message)
	assert.NotNil(t, response.Error)
}

func TestCheckEmailVerificationCode_InvalidAuthenticationToken(t *testing.T) {
	requestBody := []byte(`{
		"email": "user@email.com",
		"code": "123467"
	}`)
	req, err := http.NewRequest("POST", "/api/v1/check-email-verification-code", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer invalid token")

	mockService := &MockService{}
	req = req.WithContext(context.WithValue(req.Context(), "service", mockService))

	rr := httptest.NewRecorder()

	Ctl.CheckEmailVerificationCode(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)

	response := decodeResponseBodyForErrorResponse(t, rr)

	assert.False(t, response.IsSuccess)
	assert.Equal(t, "Authentication error: invalid token", response.Message)
	assert.NotNil(t, response.Error)
}

func TestCheckEmailVerificationCode_MalformedRequestBody(t *testing.T) {
	requestBody := []byte(`{
		"email": "user@email.com",
		"code": "123467"
	`)
	req, err := http.NewRequest("POST", "/api/v1/check-email-verification-code", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+authenticationToken)

	mockService := &MockService{}
	req = req.WithContext(context.WithValue(req.Context(), "service", mockService))

	rr := httptest.NewRecorder()

	Ctl.CheckEmailVerificationCode(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	response := decodeResponseBodyForErrorResponse(t, rr)

	assert.False(t, response.IsSuccess)
	assert.Equal(t, "Request malformed: error while parsing json", response.Message)
	assert.NotNil(t, response.Error)
}

func TestCheckEmailVerificationCode_RequestJSONDoesNotMatchTheScheme(t *testing.T) {
	requestBody := []byte(`{
		"email": "user@email.com",
		"cod": "123467"
	}`)
	req, err := http.NewRequest("POST", "/api/v1/check-email-verification-code", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+authenticationToken)

	mockService := &MockService{}
	req = req.WithContext(context.WithValue(req.Context(), "service", mockService))

	rr := httptest.NewRecorder()

	Ctl.CheckEmailVerificationCode(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	response := decodeResponseBodyForErrorResponse(t, rr)

	assert.False(t, response.IsSuccess)
	assert.Equal(t, "Input json is invalid", response.Message)
	assert.NotNil(t, response.Error)
}

func TestCheckEmailVerificationCode_VerificationCodeIsInvalid(t *testing.T) {
	requestBody := []byte(`{
		"email": "user@email.com", 
		"code": "123467"
	}`)
	req, err := http.NewRequest("POST", "/api/v1/check-email-verification-code", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+authenticationToken)

	mockService := &MockService{
		checkEmailVerificationCode: func(userID uint, code string) error {
			if code != verificationCode {
				t.Fatalf("Verification code mismatch, expected %s, got %s", verificationCode, code)
			}
			return fmt.Errorf("failed to verify code for user %d", userID)
		},
	}
	req = req.WithContext(context.WithValue(req.Context(), "service", mockService))

	rr := httptest.NewRecorder()

	Ctl.CheckEmailVerificationCode(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	response := decodeResponseBodyForErrorResponse(t, rr)

	assert.False(t, response.IsSuccess)
	assert.Equal(t, "Failed to verify code", response.Message)
	assert.NotNil(t, response.Error)
}

func TestCheckEmailVerificationCode_UserNotFound(t *testing.T) {
	requestBody := []byte(`{
		"email": "user@email.com", 
		"code": "123467"
	}`)
	req, err := http.NewRequest("POST", "/api/v1/check-email-verification-code", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+authenticationToken)

	mockService := &MockService{
		checkEmailVerificationCode: func(userID uint, code string) error {
			if code != verificationCode {
				t.Fatalf("Verification code mismatch, expected %s, got %s", verificationCode, code)
			}
			return nil
		},
		findUser: func(userId uint) (models.User, error) {
			return models.User{}, fmt.Errorf("user was not found for id %d", userId)
		},
	}
	req = req.WithContext(context.WithValue(req.Context(), "service", mockService))

	rr := httptest.NewRecorder()

	Ctl.CheckEmailVerificationCode(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	response := decodeResponseBodyForErrorResponse(t, rr)

	assert.False(t, response.IsSuccess)
	assert.Equal(t, "User with provided id does not exist", response.Message)
	assert.NotNil(t, response.Error)
}

func TestCheckEmailVerificationCode_ZkEmailProofFailedToBeGenerated(t *testing.T) {
	requestBody := []byte(`{
		"email": "user@email.com", 
		"code": "123467"
	}`)
	req, err := http.NewRequest("POST", "/api/v1/check-email-verification-code", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+authenticationToken)

	mockService := &MockService{
		checkEmailVerificationCode: func(userID uint, code string) error {
			if code != verificationCode {
				t.Fatalf("Verification code mismatch, expected %s, got %s", verificationCode, code)
			}
			return nil
		},
		findUser: func(userID uint) (models.User, error) {
			if userID != userId {
				t.Fatalf("User id mismatch, expected %d, got %d", userId, userID)
			}
			return models.User{
				ID:   userID,
				Salt: userSalt,
			}, nil
		},
		generateZkProofOfEmailVerification: func(
			user models.User,
			request dto.CheckEmailVerificationCodeDTO,
		) ([]byte, uint, error) {
			return []byte{}, zkKeyPairId, fmt.Errorf("failed to generate the zk email proof")
		},
	}
	req = req.WithContext(context.WithValue(req.Context(), "service", mockService))

	rr := httptest.NewRecorder()

	Ctl.CheckEmailVerificationCode(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	response := decodeResponseBodyForErrorResponse(t, rr)

	assert.False(t, response.IsSuccess)
	assert.Equal(t, "Failed to generate zk proof of email verification", response.Message)
	assert.NotNil(t, response.Error)
}

func TestCheckEmailVerificationCode_FailedToSaveProofOfEmailVerification(t *testing.T) {
	requestBody := []byte(`{
		"email": "user@email.com", 
		"code": "123467"
	}`)
	req, err := http.NewRequest("POST", "/api/v1/check-email-verification-code", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+authenticationToken)

	mockService := &MockService{
		checkEmailVerificationCode: func(userID uint, code string) error {
			if code != verificationCode {
				t.Fatalf("Verification code mismatch, expected %s, got %s", verificationCode, code)
			}
			return nil
		},
		findUser: func(userID uint) (models.User, error) {
			if userID != userId {
				t.Fatalf("User id mismatch, expected %d, got %d", userId, userID)
			}
			return models.User{
				ID:   userID,
				Salt: userSalt,
			}, nil
		},
		generateZkProofOfEmailVerification: func(
			user models.User,
			request dto.CheckEmailVerificationCodeDTO,
		) ([]byte, uint, error) {
			return emailProof, zkKeyPairId, nil
		},
		saveProofOfEmailVerification: func(
			userID uint, verificationCode string, zkProof []byte, zkKeyId uint,
		) error {
			if !utils.Equal(zkProof, emailProof) {
				t.Fatalf("Email proof mismatch: expected %s, got %s", emailProof, zkProof)
			}
			if zkKeyId != zkKeyPairId {
				t.Fatalf("Unexpected zk key pair id: expected %d, got %d", zkKeyPairId, zkKeyId)
			}
			return fmt.Errorf("failed to save proof of email verification")
		},
	}
	req = req.WithContext(context.WithValue(req.Context(), "service", mockService))

	rr := httptest.NewRecorder()

	Ctl.CheckEmailVerificationCode(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	response := decodeResponseBodyForErrorResponse(t, rr)

	assert.False(t, response.IsSuccess)
	assert.Equal(t, "Failed to save proof of the email verification procedure", response.Message)
	assert.NotNil(t, response.Error)
}

func TestCheckEmailVerificationCode_Success(t *testing.T) {
	requestBody := []byte(`{
		"email": "user@email.com", 
		"code": "123467"
	}`)
	req, err := http.NewRequest("POST", "/api/v1/check-email-verification-code", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+authenticationToken)

	mockService := &MockService{
		checkEmailVerificationCode: func(userID uint, code string) error {
			if code != verificationCode {
				t.Fatalf("Verification code mismatch, expected %s, got %s", verificationCode, code)
			}
			return nil
		},
		findUser: func(userID uint) (models.User, error) {
			if userID != userId {
				t.Fatalf("User id mismatch, expected %d, got %d", userId, userID)
			}
			return models.User{
				ID:   userID,
				Salt: userSalt,
			}, nil
		},
		generateZkProofOfEmailVerification: func(
			user models.User,
			request dto.CheckEmailVerificationCodeDTO,
		) ([]byte, uint, error) {
			return emailProof, zkKeyPairId, nil
		},
		saveProofOfEmailVerification: func(
			userID uint, verificationCode string, zkProof []byte, zkKeyId uint,
		) error {
			if !utils.Equal(zkProof, emailProof) {
				t.Fatalf("Email proof mismatch: expected %s, got %s", emailProof, zkProof)
			}
			if zkKeyId != zkKeyPairId {
				t.Fatalf("Unexpected zk key pair id: expected %d, got %d", zkKeyPairId, zkKeyId)
			}
			return nil
		},
	}
	req = req.WithContext(context.WithValue(req.Context(), "service", mockService))

	rr := httptest.NewRecorder()

	Ctl.CheckEmailVerificationCode(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	response := decodeResponseBodyForResponse(t, rr)

	assert.True(t, response.IsSuccess)
	assert.Equal(t, "Your email was successfully verified!", response.Message)
	assert.Equal(t, nil, response.Data)
}

func TestUpdateDisplayNameHandler_InvalidHttpRequestMethod(t *testing.T) {
	requestBody := []byte(`{"display_name": "test_user"}`)

	req, err := http.NewRequest("PUT", "/api/v1/update-display-name", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	Ctl.UpdateDisplayNameHandler(rr, req)

	assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)

	response := decodeResponseBodyForErrorResponse(t, rr)

	assert.False(t, response.IsSuccess)
	assert.Equal(t, "Invalid http method. Expected POST", response.Message)
	assert.NotNil(t, response.Error)
}

func TestUpdateDisplayNameHandler_AuthenticationTokenIsInvalid(t *testing.T) {
	requestBody := []byte(`{"display_name": "test_user"}`)
	req, err := http.NewRequest("POST", "/api/v1/update-display-name", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Authorization", "Bearer invalid token")

	mockService := &MockService{}
	req = req.WithContext(context.WithValue(req.Context(), "service", mockService))

	rr := httptest.NewRecorder()

	Ctl.UpdateDisplayNameHandler(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)

	response := decodeResponseBodyForErrorResponse(t, rr)

	assert.False(t, response.IsSuccess)
	assert.Equal(t, "Authentication error: invalid token", response.Message)
	assert.NotNil(t, response.Error)
}

func TestUpdateDisplayNameHandler_RequestJsonIsMalformed(t *testing.T) {
	requestBody := []byte(`{"display_name": "test_user"}something_else`)
	req, err := http.NewRequest("POST", "/api/v1/update-display-name", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Authorization", "Bearer "+authenticationToken)

	mockService := &MockService{}
	req = req.WithContext(context.WithValue(req.Context(), "service", mockService))

	rr := httptest.NewRecorder()

	Ctl.UpdateDisplayNameHandler(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	response := decodeResponseBodyForErrorResponse(t, rr)

	assert.False(t, response.IsSuccess)
	assert.Equal(t, "Request malformed: error while parsing json", response.Message)
	assert.NotNil(t, response.Error)
}

func TestUpdateDisplayNameHandler_RequiredRequestJsonFieldsAreMissing(t *testing.T) {
	requestBody := []byte(`{}`)
	req, err := http.NewRequest("POST", "/api/v1/update-display-name", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Authorization", "Bearer "+authenticationToken)

	mockService := &MockService{}
	req = req.WithContext(context.WithValue(req.Context(), "service", mockService))

	rr := httptest.NewRecorder()

	Ctl.UpdateDisplayNameHandler(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	response := decodeResponseBodyForErrorResponse(t, rr)

	assert.False(t, response.IsSuccess)
	assert.Equal(t, "Input json is invalid", response.Message)
	assert.NotNil(t, response.Error)
}

func TestUpdateDisplayNameHandler_Success(t *testing.T) {
	requestBody := []byte(`{"display_name": "test_user"}`)
	req, err := http.NewRequest("POST", "/api/v1/update-display-name", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Authorization", "Bearer "+authenticationToken)

	mockService := &MockService{}
	req = req.WithContext(context.WithValue(req.Context(), "service", mockService))

	rr := httptest.NewRecorder()

	Ctl.UpdateDisplayNameHandler(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response = decodeResponseBodyForErrorResponse(t, rr)

	assert.True(t, response.IsSuccess)
	assert.Equal(t, "Display name updated successfully", response.Message)
	assert.Equal(t, nil, response.Data)
	assert.Nil(t, response.Error)
}

func TestLoginClientHandler_InvalidHttpRequestMethod(t *testing.T) {
	reqBody := []byte(`{
		"username": "testuser",
		"password": "testpassword"
	}`)

	req := httptest.NewRequest("PUT", "/api/v1/login-client", bytes.NewBuffer(reqBody))

	rr := httptest.NewRecorder()

	Ctl.LoginClientHandler(rr, req)

	assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)

	response := decodeResponseBodyForErrorResponse(t, rr)

	assert.False(t, response.IsSuccess)
	assert.Equal(t, "Invalid http method. Expected POST", response.Message)
	assert.NotNil(t, response.Error)
}

func TestLoginClientHandler_RequestJsonIsMalformed(t *testing.T) {
	loginReq := []byte(`{
		"username": "testuser",
		"password": "testpassword"
	}something_else`)
	req := httptest.NewRequest("POST", "/api/v1/login-client", bytes.NewBuffer(loginReq))

	req = setMockServiceInContext(req)

	rr := httptest.NewRecorder()

	Ctl.LoginClientHandler(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	response := decodeResponseBodyForErrorResponse(t, rr)

	assert.False(t, response.IsSuccess)
	assert.Equal(t, "Request malformed: error while parsing json", response.Message)
	assert.NotNil(t, response.Error)
}

func TestLoginClientHandler_RequiredRequestJsonFieldsAreMissing(t *testing.T) {
	loginReq := []byte(`{
		"username": "testuser"
	}`)

	req := httptest.NewRequest("POST", "/api/v1/login-client", bytes.NewBuffer(loginReq))
	req = setMockServiceInContext(req)

	rr := httptest.NewRecorder()

	Ctl.LoginClientHandler(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	response := decodeResponseBodyForErrorResponse(t, rr)

	assert.False(t, response.IsSuccess)
	assert.Equal(t, "Input json is invalid", response.Message)
	assert.NotNil(t, response.Error)
}

func TestLoginClientHandler_Success(t *testing.T) {
	loginReq := []byte(`{
		"username": "testuser",
		"password": "testpassword"
	}`)

	req := httptest.NewRequest("POST", "/api/v1/login-client", bytes.NewBuffer(loginReq))

	req = setMockServiceInContext(req)

	w := httptest.NewRecorder()

	Ctl.LoginClientHandler(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var tokenResp models.LoginUserResponseOutput
	err := json.NewDecoder(w.Body).Decode(&tokenResp)
	if err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	// Validate the response
	assert.Equal(t, testToken, tokenResp.Token)
}

func TestCheckBackendURIHandler_InvalidHttpRequestMethod(t *testing.T) {
	reqBody := []byte(`{
		"backend_uri": "https://example.com"
	}`)
	req := httptest.NewRequest("GET", "/api/v1/check-backend-uri", bytes.NewBuffer(reqBody))

	rr := httptest.NewRecorder()

	Ctl.CheckBackendURI(rr, req)

	assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)

	response := decodeResponseBodyForErrorResponse(t, rr)

	assert.False(t, response.IsSuccess)
	assert.Equal(t, "Invalid http method. Expected POST", response.Message)
	assert.NotNil(t, response.Error)
}

func TestCheckBackendURIHandler_RequestJsonIsMalformed(t *testing.T) {
	reqBody := []byte(`{
		"backend_uri": "https://example.com"
	}something_else`)
	req := httptest.NewRequest("POST", "/api/v1/check-backend-uri", bytes.NewBuffer(reqBody))
	req = setMockServiceInContext(req)

	rr := httptest.NewRecorder()

	Ctl.CheckBackendURI(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	response := decodeResponseBodyForErrorResponse(t, rr)

	assert.False(t, response.IsSuccess)
	assert.Equal(t, "Request malformed: error while parsing json", response.Message)
	assert.NotNil(t, response.Error)
}

func TestCheckBackendURIHandler_RequiredRequestJsonFieldsAreMissing(t *testing.T) {
	reqBody := []byte(`{}`)
	req := httptest.NewRequest("POST", "/api/v1/check-backend-uri", bytes.NewBuffer(reqBody))
	req = setMockServiceInContext(req)

	rr := httptest.NewRecorder()

	Ctl.CheckBackendURI(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	response := decodeResponseBodyForErrorResponse(t, rr)

	assert.False(t, response.IsSuccess)
	assert.Equal(t, "Input json is invalid", response.Message)
	assert.NotNil(t, response.Error)
}

func TestCheckBackendURIHandler_Success(t *testing.T) {
	checkReq := dto.CheckBackendURIDTO{
		BackendURI: "https://example.com",
	}
	reqBody, err := json.Marshal(checkReq)
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest("POST", "/api/v1/check-backend-uri", bytes.NewBuffer(reqBody))
	req = setMockServiceInContext(req)

	w := httptest.NewRecorder()

	Ctl.CheckBackendURI(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response bool
	err = json.NewDecoder(w.Body).Decode(&response)
	if err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	assert.True(t, response)
}

func TestResetPasswordHandler_InvalidRequestMethod(t *testing.T) {
	reqBody := []byte(`{
		"username": "test_user",
		"signature": "aaaaaab",
		"new_password": "new_password"
	}`)
	req := httptest.NewRequest("GET", "/api/v1/reset-password", bytes.NewBuffer(reqBody))

	rr := httptest.NewRecorder()

	Ctl.ResetPasswordHandler(rr, req)

	assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)

	response := decodeResponseBodyForErrorResponse(t, rr)

	assert.False(t, response.IsSuccess)
	assert.Equal(t, "Invalid http method. Expected POST", response.Message)
	assert.NotNil(t, response.Error)
}

func TestResetPasswordHandler_RequestJsonIsMalformed(t *testing.T) {
	reqBody := []byte(`{
		"username": "test_user",
		"signature": "aaaaaab",
		"new_password": "new_password"
	}somethingElse`)

	req := httptest.NewRequest("POST", "/api/v1/reset-password", bytes.NewBuffer(reqBody))
	req = setMockServiceInContext(req)

	rr := httptest.NewRecorder()

	Ctl.ResetPasswordHandler(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	response := decodeResponseBodyForErrorResponse(t, rr)

	assert.False(t, response.IsSuccess)
	assert.Equal(t, "Request malformed: error while parsing json", response.Message)
	assert.NotNil(t, response.Error)
}

func TestResetPasswordHandler_RequiredRequestJsonFieldsAreMissing(t *testing.T) {
	reqBody := []byte(`{
		"username": "test_user",
		"new_password": "new_password"
	}`)

	req := httptest.NewRequest("POST", "/api/v1/reset-password", bytes.NewBuffer(reqBody))
	req = setMockServiceInContext(req)

	rr := httptest.NewRecorder()

	Ctl.ResetPasswordHandler(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	response := decodeResponseBodyForErrorResponse(t, rr)

	assert.False(t, response.IsSuccess)
	assert.Equal(t, "Input json is invalid", response.Message)
	assert.NotNil(t, response.Error)
}

func TestResetPasswordHandler_UserNotFoundForUsername(t *testing.T) {
	reqBody := []byte(`{
		"username": "test_user",
		"signature": "aaabbbbc",
		"new_password": "new_password"
	}`)

	mockService := &MockService{
		getUserForUsername: func(currUsername string) (models.User, error) {
			if currUsername != username {
				t.Fatalf("Username mismatch: expected %s, got %s", username, currUsername)
			}
			return models.User{}, fmt.Errorf("user not found")
		},
	}

	req := httptest.NewRequest("POST", "/api/v1/reset-password", bytes.NewBuffer(reqBody))
	req = req.WithContext(context.WithValue(req.Context(), "service", mockService))

	rr := httptest.NewRecorder()

	Ctl.ResetPasswordHandler(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	response := decodeResponseBodyForErrorResponse(t, rr)

	assert.False(t, response.IsSuccess)
	assert.Equal(t, "user does not exist", response.Message)
	assert.NotNil(t, response.Error)
}

func TestResetPasswordHandler_InvalidSignature(t *testing.T) {
	reqBody := []byte(`{
		"username": "test_user",
		"signature": "aaabbbbc",
		"new_password": "new_password"
	}`)

	mockService := &MockService{
		getUserForUsername: func(currUsername string) (models.User, error) {
			if currUsername != username {
				t.Fatalf("Username mismatch: expected %s, got %s", username, currUsername)
			}
			return models.User{
				ID:        userId,
				Username:  username,
				FirstName: firstName,
				LastName:  lastName,
				Password:  userPassword,
			}, nil
		},
		validateSignature: func(message string, signature []byte, publicKey []byte) error {
			return fmt.Errorf("failed to validate signature")
		},
	}

	req := httptest.NewRequest("POST", "/api/v1/reset-password", bytes.NewBuffer(reqBody))
	req = req.WithContext(context.WithValue(req.Context(), "service", mockService))

	rr := httptest.NewRecorder()

	Ctl.ResetPasswordHandler(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	response := decodeResponseBodyForErrorResponse(t, rr)

	assert.False(t, response.IsSuccess)
	assert.Equal(t, "signature is invalid", response.Message)
	assert.NotNil(t, response.Error)
}

func TestResetPasswordHandler_UserPasswordUpdateFailed(t *testing.T) {
	reqBody := []byte(`{
		"username": "test_user",
		"signature": "aaabbbbc",
		"new_password": "new_password"
	}`)

	mockService := &MockService{
		getUserForUsername: func(currUsername string) (models.User, error) {
			if currUsername != username {
				t.Fatalf("Username mismatch: expected %s, got %s", username, currUsername)
			}
			return models.User{
				ID:        userId,
				Username:  username,
				FirstName: firstName,
				LastName:  lastName,
				Password:  userPassword,
				Salt:      userSalt,
			}, nil
		},
		validateSignature: func(message string, signature []byte, publicKey []byte) error {
			return nil
		},
		updateUserPassword: func(currUsername string, newPassword string, salt string) error {
			if currUsername != username {
				t.Fatalf("Username mismatch: expected %s, got %s", username, currUsername)
			}
			if newPassword != newUserPassword {
				t.Fatalf("Password mismatch: expected %s, got %s", newUserPassword, newPassword)
			}
			if salt != userSalt {
				t.Fatalf("Salt mismatch: expected %s, got %s", userSalt, salt)
			}
			return fmt.Errorf("failed to update user password")
		},
	}

	req := httptest.NewRequest("POST", "/api/v1/reset-password", bytes.NewBuffer(reqBody))
	req = req.WithContext(context.WithValue(req.Context(), "service", mockService))

	rr := httptest.NewRecorder()

	Ctl.ResetPasswordHandler(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	response := decodeResponseBodyForErrorResponse(t, rr)

	assert.False(t, response.IsSuccess)
	assert.Equal(t, "Internal error: failed to update password", response.Message)
	assert.NotNil(t, response.Error)
}

func TestResetPasswordHandler_Success(t *testing.T) {
	reqBody := []byte(`{
		"username": "test_user",
		"signature": "aaabbbbc",
		"new_password": "new_password"
	}`)

	mockService := &MockService{
		getUserForUsername: func(currUsername string) (models.User, error) {
			if currUsername != username {
				t.Fatalf("Username mismatch: expected %s, got %s", username, currUsername)
			}
			return models.User{
				ID:        userId,
				Username:  username,
				FirstName: firstName,
				LastName:  lastName,
				Password:  userPassword,
				Salt:      userSalt,
			}, nil
		},
		validateSignature: func(message string, signature []byte, publicKey []byte) error {
			return nil
		},
		updateUserPassword: func(currUsername string, newPassword string, salt string) error {
			if currUsername != username {
				t.Fatalf("Username mismatch: expected %s, got %s", username, currUsername)
			}
			if newPassword != newUserPassword {
				t.Fatalf("Password mismatch: expected %s, got %s", newUserPassword, newPassword)
			}
			if salt != userSalt {
				t.Fatalf("Salt mismatch: expected %s, got %s", userSalt, salt)
			}
			return nil
		},
	}

	req := httptest.NewRequest("POST", "/api/v1/reset-password", bytes.NewBuffer(reqBody))
	req = req.WithContext(context.WithValue(req.Context(), "service", mockService))

	rr := httptest.NewRecorder()

	Ctl.ResetPasswordHandler(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)

	response := decodeResponseBodyForResponse(t, rr)

	assert.Equal(t, true, response.IsSuccess)
	assert.Equal(t, "Your password was updated successfully!", response.Message)
	assert.Nil(t, response.Error)
}

func setMockServiceInContext(req *http.Request) *http.Request {
	mockSvc := &MockService{}
	ctx := context.WithValue(req.Context(), "service", mockSvc)
	return req.WithContext(ctx)
}

func TestRegisterUserPrecheck_Success(t *testing.T) {
	requestBody := []byte(`{"username": "test_user"}`)

	req, err := http.NewRequest("POST", "/api/v1/register-user-precheck", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	mockService := &MockService{
		registerUserPrecheck: func(req dto.RegisterUserPrecheckDTO, iterCount int) (string, error) {
			assert.Equal(t, "test_user", req.Username, "Username should match")
			assert.Equal(t, 4096, iterCount, "Iteration count should match")

			return userSalt, nil
		},
	}

	req = req.WithContext(context.WithValue(req.Context(), "service", mockService))

	t.Setenv("SCRAM_ITERATION_COUNT", "4096")

	rr := httptest.NewRecorder()

	Ctl.RegisterUserPrecheck(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code, "Response code should be 201 Created")

	response := decodeResponseBodyForResponse(t, rr)

	assert.Equal(t, true, response.IsSuccess)
	assert.Equal(t, "User is successfully registered", response.Message, "Response message should match")
	assert.Nil(t, response.Error)
}

func TestRegisterUserPrecheck_InvalidMethod(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/v1/register-user-precheck", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	Ctl.RegisterUserPrecheck(rr, req)

	assert.Equal(t, http.StatusMethodNotAllowed, rr.Code, "Expected HTTP 405 Method Not Allowed")
}

func TestRegisterUserPrecheck_InvalidIterationCount(t *testing.T) {
	requestBody := []byte(`{
		"username": "test_user"
	}`)

	req, err := http.NewRequest("POST", "/api/v1/register-user-precheck", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	mockService := &MockService{}
	req = req.WithContext(context.WithValue(req.Context(), "service", mockService))

	t.Setenv("SCRAM_ITERATION_COUNT", "invalid_value")

	rr := httptest.NewRecorder()

	Ctl.RegisterUserPrecheck(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code, "Expected HTTP 500 Internal Server Error")
}

func TestRegisterUserPrecheck_InvalidJSON(t *testing.T) {
	requestBody := []byte(`invalid_json`)

	req, err := http.NewRequest("POST", "/api/v1/register-user-precheck", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	mockService := &MockService{}
	req = req.WithContext(context.WithValue(req.Context(), "service", mockService))

	t.Setenv("SCRAM_ITERATION_COUNT", "4096")

	rr := httptest.NewRecorder()

	Ctl.RegisterUserPrecheck(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code, "Expected HTTP 400 Bad Request")
}

func TestRegisterUserPrecheck_MissingRequiredFields(t *testing.T) {
	reqBody := []byte(`{
		"other_field": "some value"
	}`)

	mockService := &MockService{
		registerUserPrecheck: func(request dto.RegisterUserPrecheckDTO, iterCount int) (string, error) {
			return "", nil
		},
	}

	req := httptest.NewRequest("POST", "/api/register-user-precheck", bytes.NewBuffer(reqBody))
	req = req.WithContext(context.WithValue(req.Context(), "service", mockService))
	req.Header.Set("Content-Type", "application/json")

	t.Setenv("SCRAM_ITERATION_COUNT", "4096")
	rr := httptest.NewRecorder()

	Ctl.RegisterUserPrecheck(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	response := decodeResponseBodyForResponse(t, rr)

	assert.Equal(t, false, response.IsSuccess)
	assert.Equal(t, "Input json is invalid", response.Message)
}

func TestRegisterUserPrecheck_ServiceError(t *testing.T) {
	requestBody := []byte(`{
		"username": "test_user"
	}`)

	req, err := http.NewRequest("POST", "/api/v1/register-user-precheck", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	mockService := &MockService{
		registerUserPrecheck: func(req dto.RegisterUserPrecheckDTO, iterCount int) (string, error) {
			return "", fmt.Errorf("mock service error")
		},
	}

	req = req.WithContext(context.WithValue(req.Context(), "service", mockService))

	t.Setenv("SCRAM_ITERATION_COUNT", "4096")

	rr := httptest.NewRecorder()

	Ctl.RegisterUserPrecheck(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code, "Expected HTTP 400 Bad Request")
}

func TestRegisterUserHandlerv2_InvalidHttpRequestMethod(t *testing.T) {
	requestBody := []byte(`{
		"username": "test_user",
		"first_name": "Test",
		"last_name": "User",
		"display_name": "user",
		"country": "Unknown",
		"public_key": "0xaaaaaa",
		"server_key": "0xbbbbbb",
		"stored_key": "0xcccccc"
	}`)

	req, err := http.NewRequest("GET", "/api/v2/register-user", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	setMockServiceInContext(req)

	rr := httptest.NewRecorder()

	Ctl.RegisterUserHandlerv2(rr, req)

	assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)

	response := decodeResponseBodyForErrorResponse(t, rr)

	assert.False(t, response.IsSuccess)
	assert.Equal(t, "Invalid http method. Expected POST", response.Message)
	assert.NotNil(t, response.Error)
}

func TestRegisterUserHandlerv2_RequestJsonIsMalformed(t *testing.T) {
	requestBody := []byte(`{
		"username": "test_user",
		"first_name": "Test",
		"last_name": "User",
		"display_name": "user",
		"country": "Unknown",
		"public_key": "0xaaaaaa",
		"server_key": "0xbbbbbb",
		"stored_key": "0xcccccc"
	}something_else`)

	req, err := http.NewRequest("POST", "/api/v2/register-user", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	mockService := &MockService{}
	req = req.WithContext(context.WithValue(req.Context(), "service", mockService))
	// setMockServiceInContext(req)

	rr := httptest.NewRecorder()

	Ctl.RegisterUserHandlerv2(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	response := decodeResponseBodyForErrorResponse(t, rr)

	assert.False(t, response.IsSuccess)
	assert.Equal(t, "Request malformed: error while parsing json", response.Message)
	assert.NotNil(t, response.Error)
}

func TestRegisterUserHandlerv2_RequiredRequestJsonFieldsAreMissing(t *testing.T) {
	requestBody := []byte(`{
		"first_name": "Test",
		"last_name": "User",
		"display_name": "user",
		"country": "Unknown",
		"public_key": "0xaaaaaa",
		"server_key": "0xbbbbbb",
		"stored_key": "0xcccccc"
	}`)

	req, err := http.NewRequest("POST", "/api/v2/register-user", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	mockService := &MockService{}
	req = req.WithContext(context.WithValue(req.Context(), "service", mockService))

	rr := httptest.NewRecorder()

	Ctl.RegisterUserHandlerv2(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	response := decodeResponseBodyForErrorResponse(t, rr)

	assert.False(t, response.IsSuccess)
	assert.Equal(t, "Input json is invalid", response.Message)
	assert.NotNil(t, response.Error)
}

func TestRegisterUserHandlerv2_ServiceError(t *testing.T) {
	// Mock request body
	requestBody := []byte(`{
		"username": "test_user",
		"first_name": "Test",
		"last_name": "User",
		"display_name": "user",
		"country": "Unknown",
		"public_key": "0xaaaaaa",
		"server_key": "0xbbbbbb",
		"stored_key": "0xcccccc"
	}`)

	// Create a mock request
	req, err := http.NewRequest("POST", "/api/v2/register-user", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	// Create a mock service and set it in the request context
	mockService := &MockService{
		registerUserv2: func(req dto.RegisterUserDTOv2) error {
			return fmt.Errorf("failed to register user")
		},
	}

	req = req.WithContext(context.WithValue(req.Context(), "service", mockService))

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the handler function
	Ctl.RegisterUserHandlerv2(rr, req)

	// Check the status code
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	// Decode the response body
	response := decodeResponseBodyForResponse(t, rr)

	// Now assert the fields directly
	assert.False(t, response.IsSuccess)
	assert.Equal(t, "Failed to register user", response.Message)
	assert.NotNil(t, response.Error)
}

func TestRegisterUserHandlerv2_Success(t *testing.T) {
	// Mock request body
	requestBody := []byte(`{
		"username": "test_user",
		"first_name": "Test",
		"last_name": "User",
		"display_name": "user",
		"country": "Unknown",
		"public_key": "0xaaaaaa",
		"server_key": "0xbbbbbb",
		"stored_key": "0xcccccc"
	}`)

	// Create a mock request
	req, err := http.NewRequest("POST", "/api/v2/register-user", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	// Create a mock service and set it in the request context
	mockService := &MockService{
		registerUserv2: func(req dto.RegisterUserDTOv2) error {
			return nil
		},
	}

	req = req.WithContext(context.WithValue(req.Context(), "service", mockService))

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the handler function
	Ctl.RegisterUserHandlerv2(rr, req)

	// Check the status code
	assert.Equal(t, http.StatusCreated, rr.Code)

	// Decode the response body
	response := decodeResponseBodyForResponse(t, rr)

	// Now assert the fields directly
	assert.True(t, response.IsSuccess)
	assert.Equal(t, "User registered successfully", response.Message)
	assert.Equal(t, nil, response.Data)
}

func TestResetPasswordPrecheckHandler_Success(t *testing.T) {
	requestBody := []byte(`{"username": "test_user"}`)

	req, err := http.NewRequest("POST", "/api/v2/reset-password-precheck", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	mockService := &MockService{
		getUserForUsername: func(username string) (models.User, error) {
			return models.User{
				ID:        userId,
				Username:  username,
				PublicKey: []byte("mock_public_key"),
			}, nil
		},
	}

	req = req.WithContext(context.WithValue(req.Context(), "service", mockService))

	rr := httptest.NewRecorder()

	Ctl.ResetPasswordPrecheckHandler(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "Expected response status code to be 200 OK")

	response := decodeResponseBodyForResponse(t, rr)

	assert.Equal(t, true, response.IsSuccess)
	assert.Equal(t, "User does exist!", response.Message, "Response message should match")
	assert.Nil(t, response.Error)
}

func TestResetPasswordPrecheckHandler_RequiredRequestJsonFieldIsMissing(t *testing.T) {
	requestBody := []byte(`{}`)

	req, err := http.NewRequest("POST", "/api/v2/reset-password-precheck", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	mockService := &MockService{}

	req = req.WithContext(context.WithValue(req.Context(), "service", mockService))

	rr := httptest.NewRecorder()

	Ctl.ResetPasswordPrecheckHandler(rr, req)

	response := decodeResponseBodyForErrorResponse(t, rr)

	assert.Equal(t, false, response.IsSuccess)
	assert.Equal(t, http.StatusBadRequest, rr.Code, "Response should be 400 Bad Request")
	assert.NotNil(t, response.Error)
}

func TestResetPasswordPrecheckHandler_InvalidHttpMethod(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/v2/reset-password-precheck", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	Ctl.ResetPasswordPrecheckHandler(rr, req)

	assert.Equal(t, http.StatusMethodNotAllowed, rr.Code, "Response should be 405 Method Not Allowed")

	response := decodeResponseBodyForErrorResponse(t, rr)

	assert.Equal(t, false, response.IsSuccess)
	assert.Equal(t, "Invalid http method. Expected POST", response.Message)
	assert.NotNil(t, response.Error)
}

func TestResetPasswordPrecheckHandler_UserNotFound(t *testing.T) {
	requestBody := []byte(`{"username": "nonexistent_user"}`)

	req, err := http.NewRequest("POST", "/api/v2/reset-password-precheck", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	mockService := &MockService{
		getUserForUsername: func(username string) (models.User, error) {
			return models.User{}, errors.New("User not found")
		},
	}

	req = req.WithContext(context.WithValue(req.Context(), "service", mockService))

	rr := httptest.NewRecorder()

	Ctl.ResetPasswordPrecheckHandler(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code, "Response should be 400 Bad Request")

	response := decodeResponseBodyForErrorResponse(t, rr)

	assert.Equal(t, false, response.IsSuccess)
	assert.Equal(t, "User does not exist!", response.Message, "Response message should match")
	assert.NotNil(t, response.Error)
}

func TestResetPasswordPrecheckHandler_InvalidJSON(t *testing.T) {
	requestBody := []byte(`invalid_json`)

	req, err := http.NewRequest("POST", "/api/v2/reset-password-precheck", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	mockService := &MockService{}
	req = req.WithContext(context.WithValue(req.Context(), "service", mockService))

	rr := httptest.NewRecorder()

	Ctl.ResetPasswordPrecheckHandler(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code, "Expected HTTP 400 Bad Request")

	response := decodeResponseBodyForErrorResponse(t, rr)

	assert.Equal(t, false, response.IsSuccess)
	assert.Equal(t, "Request malformed: error while parsing json", response.Message)
	assert.NotNil(t, response.Error)
}

func TestResetPasswordHandlerV2_Success(t *testing.T) {
	reqBody := []byte(`{
		"username": "test_user",
		"signature": "aaabbbbc",
		"stored_key": "user_stored_key",
		"server_key": "user_server_key"
	}`)

	mockService := &MockService{
		getUserForUsername: func(currUsername string) (models.User, error) {
			if currUsername != username {
				t.Fatalf("Username mismatch: expected %s, got %s", "test_user", currUsername)
			}
			return models.User{
				ID:        userId,
				Username:  username,
				PublicKey: []byte("mock_public_key"),
			}, nil
		},
		validateSignature: func(message string, signature []byte, publicKey []byte) error {
			if message != "Sign-in with Layer8" {
				t.Fatalf("Message mismatch: expected %s, got %s", "Sign-in with Layer8", message)
			}
			return nil
		},
		updateUserPasswordV2: func(currUsername, currStoredKey, currServerKey string) error {
			if currUsername != username {
				t.Fatalf("Username mismatch: expected %s, got %s", username, currUsername)
			}
			if currStoredKey != storedKey {
				t.Fatalf("currStoredKey mismatch: expected %s, got %s", storedKey, currStoredKey)
			}
			if currServerKey != serverKey {
				t.Fatalf("currServerKey mismatch: expected %s, got %s", serverKey, currServerKey)
			}
			return nil
		},
	}
	req := httptest.NewRequest("POST", "/api/v2/reset-password", bytes.NewBuffer(reqBody))
	req = req.WithContext(context.WithValue(req.Context(), "service", mockService))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	Ctl.ResetPasswordHandlerV2(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)

	response := decodeResponseBodyForResponse(t, rr)

	assert.Equal(t, true, response.IsSuccess)
	assert.Equal(t, "Your password was updated successfully!", response.Message)
	assert.Nil(t, response.Error)
}

func TestResetPasswordHandlerV2_InvalidJSON(t *testing.T) {
	reqBody := []byte(`{"invalid-request"}`)

	mockService := &MockService{}

	req := httptest.NewRequest("POST", "/api/v2/reset-password", bytes.NewBuffer(reqBody))
	req = req.WithContext(context.WithValue(req.Context(), "service", mockService))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	Ctl.ResetPasswordHandlerV2(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	response := decodeResponseBodyForErrorResponse(t, rr)

	assert.Equal(t, false, response.IsSuccess)
	assert.Equal(t, "Request malformed: error while parsing json", response.Message)
	assert.NotNil(t, response.Error)
}

func TestResetPasswordHandlerV2_UserNotFound(t *testing.T) {
	reqBody := []byte(`{
		"username": "non_existent_user",
		"signature": "aaabbbbc",
		"stored_key": "test_stored_key",
		"server_key": "test_server_key"
	}`)

	mockService := &MockService{
		getUserForUsername: func(currUsername string) (models.User, error) {
			if currUsername != "non_existent_user" {
				t.Fatalf("Username mismatch: expected %s, got %s", "non_existent_user", currUsername)
			}
			return models.User{}, errors.New("User not found")
		},
	}

	req := httptest.NewRequest("POST", "/api/v2/reset-password", bytes.NewBuffer(reqBody))
	req = req.WithContext(context.WithValue(req.Context(), "service", mockService))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	Ctl.ResetPasswordHandlerV2(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)

	response := decodeResponseBodyForErrorResponse(t, rr)

	assert.Equal(t, false, response.IsSuccess)
	assert.Contains(t, response.Error, "User not found")
	assert.NotNil(t, response.Error)
}

func TestResetPasswordHandlerV2_InvalidSignature(t *testing.T) {
	reqBody := []byte(`{
		"username": "test_user",
		"signature": "aaabbbbc",
		"stored_key": "user_stored_key",
		"server_key": "user_server_key"
	}`)

	mockService := &MockService{
		getUserForUsername: func(currUsername string) (models.User, error) {
			if currUsername != username {
				t.Fatalf("Username mismatch: expected %s, got %s", username, currUsername)
			}
			return models.User{
				ID:        userId,
				Username:  username,
				PublicKey: []byte("mock_public_key"),
			}, nil
		},
		validateSignature: func(message string, signature []byte, publicKey []byte) error {
			return fmt.Errorf("invalid signature")
		},
		updateUserPasswordV2: func(currUsername, storedKey, serverKey string) error {
			return nil
		},
	}

	req := httptest.NewRequest("POST", "/api/v2/reset-password", bytes.NewBuffer(reqBody))
	req = req.WithContext(context.WithValue(req.Context(), "service", mockService))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	Ctl.ResetPasswordHandlerV2(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	response := decodeResponseBodyForErrorResponse(t, rr)

	assert.Equal(t, false, response.IsSuccess)
	assert.Equal(t, "Signature is invalid!", response.Message)
	assert.NotNil(t, response.Error)
}

func TestResetPasswordHandlerV2_UpdatePasswordFailure(t *testing.T) {
	reqBody := []byte(`{
		"username": "test_user",
		"signature": "aaabbbbc",
		"stored_key": "test_stored_key",
		"server_key": "test_server_key"
	}`)

	mockService := &MockService{
		getUserForUsername: func(currUsername string) (models.User, error) {
			if currUsername != username {
				t.Fatalf("Username mismatch: expected %s, got %s", username, currUsername)
			}
			return models.User{
				ID:        userId,
				Username:  username,
				PublicKey: []byte("mock_public_key"),
			}, nil
		},
		validateSignature: func(message string, signature []byte, publicKey []byte) error {
			return nil
		},
		updateUserPasswordV2: func(currUsername, storedKey, serverKey string) error {
			return errors.New("failed to update password")
		},
	}

	req := httptest.NewRequest("POST", "/api/v2/reset-password", bytes.NewBuffer(reqBody))
	req = req.WithContext(context.WithValue(req.Context(), "service", mockService))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	Ctl.ResetPasswordHandlerV2(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	response := decodeResponseBodyForErrorResponse(t, rr)

	assert.Equal(t, false, response.IsSuccess)
	assert.Equal(t, "Internal error: failed to update user", response.Message)
	assert.NotNil(t, response.Error)
}

func TestResetPasswordHandlerV2_MissingRequiredField(t *testing.T) {
	reqBody := []byte(`{
		"other_field": "missing_field",
	}`)

	mockService := &MockService{
		getUserForUsername: func(currUsername string) (models.User, error) {
			return models.User{}, errors.New("User not found")
		},
	}

	req := httptest.NewRequest("POST", "/api/v2/reset-password", bytes.NewBuffer(reqBody))
	req = req.WithContext(context.WithValue(req.Context(), "service", mockService))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	Ctl.ResetPasswordHandlerV2(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	response := decodeResponseBodyForResponse(t, rr)

	assert.Equal(t, false, response.IsSuccess)
	assert.Equal(t, "Request malformed: error while parsing json", response.Message)
	assert.NotNil(t, response.Message)
}

func TestLoginClientPrecheckHandlerv2_InvalidHttpRequestMethod(t *testing.T) {
	requestBody := []byte(`{"username": "test_client", "c_nonce": "Test_Nonce"}`)

	req, err := http.NewRequest("GET", "/api/v2/login-client-precheck", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	mockService := &MockService{}

	req = req.WithContext(context.WithValue(req.Context(), "service", mockService))

	rr := httptest.NewRecorder()

	Ctl.LoginClientPrecheckHandlerv2(rr, req)

	assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)

	response := decodeResponseBodyForErrorResponse(t, rr)

	assert.False(t, response.IsSuccess)
	assert.Equal(t, "Invalid http method. Expected POST", response.Message)
	assert.NotNil(t, response.Error)
}

func TestLoginClientPrecheckHandlerv2_RequestJsonIsMalformed(t *testing.T) {
	requestBody := []byte(`{"username": "test_client", "c_nonce": "Test_Nonce"}something_else`)

	req, err := http.NewRequest("POST", "/api/v2/login-client-precheck", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	mockService := &MockService{}

	req = req.WithContext(context.WithValue(req.Context(), "service", mockService))

	rr := httptest.NewRecorder()

	Ctl.LoginClientPrecheckHandlerv2(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	response := decodeResponseBodyForErrorResponse(t, rr)

	assert.False(t, response.IsSuccess)
	assert.Equal(t, "Request malformed: error while parsing json", response.Message)
	assert.NotNil(t, response.Error)
}

func TestLoginClientPrecheckHandlerv2_RequiredRequestJsonFieldsAreMissing(t *testing.T) {
	requestBody := []byte(`{}`)

	req, err := http.NewRequest("POST", "/api/v2/login-client-precheck", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	mockService := &MockService{}

	req = req.WithContext(context.WithValue(req.Context(), "service", mockService))

	rr := httptest.NewRecorder()

	Ctl.LoginClientPrecheckHandlerv2(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	response := decodeResponseBodyForErrorResponse(t, rr)

	assert.False(t, response.IsSuccess)
	assert.Equal(t, "Input json is invalid", response.Message)
	assert.NotNil(t, response.Error)
}

func TestLoginClientPrecheckHandlerv2_ServiceError(t *testing.T) {
	// Mock request body
	requestBody := []byte(`{"username": "test_client", "c_nonce": "Test_Nonce"}`)

	// Create a mock request
	req, err := http.NewRequest("POST", "/api/v2/login-client-precheck", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	// Create a mock service and set it in the request context
	mockService := &MockService{
		loginPrecheckClientv2: func(req dto.LoginPrecheckDTOv2) (models.LoginPrecheckResponseOutputv2, error) {
			return models.LoginPrecheckResponseOutputv2{}, fmt.Errorf("mock service error")
		},
	}

	req = req.WithContext(context.WithValue(req.Context(), "service", mockService))

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the handler function
	Ctl.LoginClientPrecheckHandlerv2(rr, req)

	response := decodeResponseBodyForErrorResponse(t, rr)

	// Now assert the fields directly
	assert.False(t, response.IsSuccess)
	assert.Equal(t, "Failed to perform precheck, service error", response.Message)
	assert.NotNil(t, response.Error)

	// Check the status code
	assert.Equal(t, http.StatusBadRequest, rr.Code, "Expected HTTP 400 Bad Request")
}

func TestLoginClientPrecheckHandlerv2_Success(t *testing.T) {
	// Mock request body
	requestBody := []byte(`{"username": "test_client", "c_nonce": "Test_Nonce"}`)

	// Create a mock request
	req, err := http.NewRequest("POST", "/api/v2/login-client-precheck", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	// Create a mock service and set it in the request context
	mockService := &MockService{
		loginPrecheckClientv2: func(req dto.LoginPrecheckDTOv2) (models.LoginPrecheckResponseOutputv2, error) {
			return models.LoginPrecheckResponseOutputv2{
				Salt:      userSalt,
				IterCount: iterCount,
				Nonce:     nonce,
			}, nil
		},
	}

	req = req.WithContext(context.WithValue(req.Context(), "service", mockService))

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the handler function
	Ctl.LoginClientPrecheckHandlerv2(rr, req)

	// Check the status code
	assert.Equal(t, http.StatusOK, rr.Code)

	response := decodeResponseBodyForResponse(t, rr)

	// Convert response.Data to JSON bytes for unmarshalling
	dataBytes, err := json.Marshal(response.Data)
	if err != nil {
		t.Fatal(err)
	}

	var loginPrecheckResponse models.LoginPrecheckResponseOutputv2
	if err := json.Unmarshal(dataBytes, &loginPrecheckResponse); err != nil {
		t.Fatal(err)
	}

	assert.True(t, response.IsSuccess)
	assert.Nil(t, response.Error)

	// Now assert the fields directly
	assert.Equal(t, userSalt, loginPrecheckResponse.Salt)
	assert.Equal(t, iterCount, loginPrecheckResponse.IterCount)
	assert.Equal(t, nonce, loginPrecheckResponse.Nonce)
}

func TestLoginClientHandlerv2_InvalidHttpRequestMethod(t *testing.T) {
	requestBody := []byte(`{
		"username": 	"test_client",
		"nonce": 		"Test_Nonce",
		"c_nonce": 		"Test_Nonce",
		"client_proof": "Test_Client_Proof"
		}`)

	req, err := http.NewRequest("GET", "/api/v2/login-client", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	mockService := &MockService{}

	req = req.WithContext(context.WithValue(req.Context(), "service", mockService))

	rr := httptest.NewRecorder()

	Ctl.LoginClientHandlerv2(rr, req)

	assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)

	response := decodeResponseBodyForErrorResponse(t, rr)

	assert.False(t, response.IsSuccess)
	assert.Equal(t, "Invalid http method. Expected POST", response.Message)
	assert.NotNil(t, response.Error)
}

func TestLoginClientHandlerv2_RequestJsonIsMalformed(t *testing.T) {
	requestBody := []byte(`{
		"username": 	"test_client",
		"nonce": 		"Test_Nonce",
		"c_nonce": 		"Test_Nonce",
		"client_proof": "Test_Client_Proof"
		}something_else`)

	req, err := http.NewRequest("POST", "/api/v2/login-client", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	mockService := &MockService{}

	req = req.WithContext(context.WithValue(req.Context(), "service", mockService))

	rr := httptest.NewRecorder()

	Ctl.LoginClientHandlerv2(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	response := decodeResponseBodyForErrorResponse(t, rr)

	assert.False(t, response.IsSuccess)
	assert.Equal(t, "Request malformed: error while parsing json", response.Message)
	assert.NotNil(t, response.Error)
}

func TestLoginClientHandlerv2_RequiredRequestJsonFieldsAreMissing(t *testing.T) {
	requestBody := []byte(`{
		"nonce": 		"Test_Nonce",
		"client_proof": "Test_Client_Proof"
		}`)

	req, err := http.NewRequest("POST", "/api/v2/login-client", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	mockService := &MockService{}

	req = req.WithContext(context.WithValue(req.Context(), "service", mockService))

	rr := httptest.NewRecorder()

	Ctl.LoginClientHandlerv2(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	response := decodeResponseBodyForErrorResponse(t, rr)

	assert.False(t, response.IsSuccess)
	assert.Equal(t, "Input json is invalid", response.Message)
	assert.NotNil(t, response.Error)
}

func TestLoginClientHandlerv2_ServiceError(t *testing.T) {
	// Mock request body
	requestBody := []byte(`{
		"username": 	"test_client",
		"nonce": 		"Test_Nonce",
		"c_nonce": 		"Test_Nonce",
		"client_proof": "Test_Client_Proof"
		}`)

	// Create a mock request
	req, err := http.NewRequest("POST", "/api/v2/login-client", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	// Create a mock service and set it in the request context
	mockService := &MockService{
		loginClientv2: func(req dto.LoginClientDTOv2) (models.LoginClientResponseOutputv2, error) {
			return models.LoginClientResponseOutputv2{}, fmt.Errorf("mock service error")
		},
	}

	req = req.WithContext(context.WithValue(req.Context(), "service", mockService))

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the handler function
	Ctl.LoginClientHandlerv2(rr, req)

	// Decode the response body
	response := decodeResponseBodyForErrorResponse(t, rr)

	// Now assert the fields directly
	assert.False(t, response.IsSuccess)
	assert.Equal(t, "Failed to perform login", response.Message)
	assert.NotNil(t, response.Error)

	// Check the status code
	assert.Equal(t, http.StatusBadRequest, rr.Code, "Expected HTTP 400 Bad Request")
}

func TestLoginClientHandlerv2_Success(t *testing.T) {
	// Mock request body
	requestBody := []byte(`{
		"username": 	"test_client",
		"nonce": 		"Test_Nonce",
		"c_nonce": 		"Test_Nonce",
		"client_proof": "Test_Client_Proof"
		}`)

	// Create a mock request
	req, err := http.NewRequest("POST", "/api/v2/login-client", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	// Create a mock service and set it in the request context
	mockService := &MockService{
		loginClientv2: func(req dto.LoginClientDTOv2) (models.LoginClientResponseOutputv2, error) {
			return models.LoginClientResponseOutputv2{
				Token:           testToken,
				ServerSignature: serverSignature,
			}, nil
		},
	}

	req = req.WithContext(context.WithValue(req.Context(), "service", mockService))

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the handler function
	Ctl.LoginClientHandlerv2(rr, req)

	// Check the status code
	assert.Equal(t, http.StatusOK, rr.Code)

	response := decodeResponseBodyForResponse(t, rr)

	// Convert response.Data to JSON bytes for unmarshalling
	dataBytes, err := json.Marshal(response.Data)
	if err != nil {
		t.Fatal(err)
	}

	var loginUserResponse models.LoginClientResponseOutputv2
	if err := json.Unmarshal(dataBytes, &loginUserResponse); err != nil {
		t.Fatal(err)
	}

	assert.True(t, response.IsSuccess)
	assert.Nil(t, response.Error)
	assert.Equal(t, "Login successful", response.Message)

	// Now assert the fields directly
	assert.Equal(t, testToken, loginUserResponse.Token)
	assert.Equal(t, serverSignature, loginUserResponse.ServerSignature)
}

func TestRegisterClientPrecheck_Success(t *testing.T) {
	requestBody := []byte(`{"username": "test_user"}`)

	req, err := http.NewRequest("POST", "/api/v2/register-client-precheck", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	mockService := &MockService{
		registerClientPrecheck: func(req dto.RegisterClientPrecheckDTO, iterCount int) (string, error) {
			assert.Equal(t, username, req.Username, "Username should match")
			assert.Equal(t, iterationCount, iterCount, "Iteration count should match")

			return clientSalt, nil
		},
	}

	req = req.WithContext(context.WithValue(req.Context(), "service", mockService))

	t.Setenv("SCRAM_ITERATION_COUNT", "4096")

	rr := httptest.NewRecorder()

	Ctl.RegisterClientPrecheckHandler(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code, "Response code should be 201 Created")

	response := decodeResponseBodyForResponse(t, rr)

	assert.Equal(t, true, response.IsSuccess)
	assert.Equal(t, "Client is successfully registered", response.Message, "Response message should match")
	assert.Nil(t, response.Error)
}

func TestRegisterClientPrecheck_InvalidMethod(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/v2/register-client-precheck", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	Ctl.RegisterClientPrecheckHandler(rr, req)

	assert.Equal(t, http.StatusMethodNotAllowed, rr.Code, "Expected HTTP 405 Method Not Allowed")
}

func TestRegisterClientPrecheck_InvalidIterationCount(t *testing.T) {
	requestBody := []byte(`{
		"username": "test_user"
	}`)

	req, err := http.NewRequest("POST", "/api/v2/register-client-precheck", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	mockService := &MockService{}
	req = req.WithContext(context.WithValue(req.Context(), "service", mockService))

	t.Setenv("SCRAM_ITERATION_COUNT", "invalid_value")

	rr := httptest.NewRecorder()

	Ctl.RegisterClientPrecheckHandler(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code, "Expected HTTP 500 Internal Server Error")
}

func TestRegisterClientPrecheck_InvalidJSON(t *testing.T) {
	requestBody := []byte(`invalid_json`)

	req, err := http.NewRequest("POST", "/api/v2/register-client-precheck", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	mockService := &MockService{}
	req = req.WithContext(context.WithValue(req.Context(), "service", mockService))

	t.Setenv("SCRAM_ITERATION_COUNT", "4096")

	rr := httptest.NewRecorder()

	Ctl.RegisterClientPrecheckHandler(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code, "Expected HTTP 400 Bad Request")
}

func TestRegisterClientPrecheck_MissingRequiredFields(t *testing.T) {
	reqBody := []byte(`{
		"other_field": "some value"
	}`)

	mockService := &MockService{
		registerClientPrecheck: func(req dto.RegisterClientPrecheckDTO, iterCount int) (string, error) {
			return "", nil
		},
	}

	req := httptest.NewRequest("POST", "/api/v2/register-client-precheck", bytes.NewBuffer(reqBody))
	req = req.WithContext(context.WithValue(req.Context(), "service", mockService))
	req.Header.Set("Content-Type", "application/json")

	t.Setenv("SCRAM_ITERATION_COUNT", "4096")
	rr := httptest.NewRecorder()

	Ctl.RegisterClientPrecheckHandler(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	response := decodeResponseBodyForResponse(t, rr)

	assert.Equal(t, false, response.IsSuccess)
	assert.Equal(t, "Input json is invalid", response.Message)
}

func TestRegisterClientPrecheck_ServiceError(t *testing.T) {
	requestBody := []byte(`{
		"username": "test_user"
	}`)

	req, err := http.NewRequest("POST", "/api/v2/register-client-precheck", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	mockService := &MockService{
		registerClientPrecheck: func(req dto.RegisterClientPrecheckDTO, iterCount int) (string, error) {
			return "", fmt.Errorf("mock service error")
		},
	}

	req = req.WithContext(context.WithValue(req.Context(), "service", mockService))

	t.Setenv("SCRAM_ITERATION_COUNT", "4096")

	rr := httptest.NewRecorder()

	Ctl.RegisterClientPrecheckHandler(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code, "Expected HTTP 400 Bad Request")
}

func TestRegisterClientHandlerv2_InvalidHttpRequestMethod(t *testing.T) {
	requestBody := []byte(`{
		"name": "test_client",
		"redirect_uri": "https://localhost:3000/callback",
		"backend_uri": "https://localhost:8080",
		"username": "test_user",
		"server_key": "0xbbbbbb",
		"stored_key": "0xcccccc"
	}`)

	req, err := http.NewRequest("GET", "/api/v2/register-client", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	setMockServiceInContext(req)

	rr := httptest.NewRecorder()

	Ctl.RegisterClientHandlerv2(rr, req)

	assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)

	response := decodeResponseBodyForErrorResponse(t, rr)

	assert.False(t, response.IsSuccess)
	assert.Equal(t, "Invalid http method. Expected POST", response.Message)
	assert.NotNil(t, response.Error)
}

func TestRegisterClientHandlerv2_RequestJsonIsMalformed(t *testing.T) {
	requestBody := []byte(`{
		"name": "test_client",
		"redirect_uri": "https://localhost:3000/callback",
		"backend_uri": "https://localhost:8080",
		"username": "test_user",
		"server_key": "0xbbbbbb",
		"stored_key": "0xcccccc"
	}something_else`)

	req, err := http.NewRequest("POST", "/api/v2/register-client", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	mockService := &MockService{}
	req = req.WithContext(context.WithValue(req.Context(), "service", mockService))

	rr := httptest.NewRecorder()

	Ctl.RegisterClientHandlerv2(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	response := decodeResponseBodyForErrorResponse(t, rr)

	assert.False(t, response.IsSuccess)
	assert.Equal(t, "Request malformed: error while parsing json", response.Message)
	assert.NotNil(t, response.Error)
}

func TestRegisterClientHandlerv2_RequiredRequestJsonFieldsAreMissing(t *testing.T) {
	requestBody := []byte(`{
		"name": "test_client",
		"redirect_uri": "https://localhost:3000/callback",
		"backend_uri": "https://localhost:8080",
		"server_key": "0xbbbbbb",
		"stored_key": "0xcccccc"
	}`)

	req, err := http.NewRequest("POST", "/api/v2/register-client", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	mockService := &MockService{}
	req = req.WithContext(context.WithValue(req.Context(), "service", mockService))

	rr := httptest.NewRecorder()

	Ctl.RegisterClientHandlerv2(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	response := decodeResponseBodyForErrorResponse(t, rr)

	assert.False(t, response.IsSuccess)
	assert.Equal(t, "Input json is invalid", response.Message)
	assert.NotNil(t, response.Error)
}

func TestRegisterClientHandlerv2_ServiceError(t *testing.T) {
	// Mock request body
	requestBody := []byte(`{
		"name": "test_client",
		"redirect_uri": "https://localhost:3000/callback",
		"backend_uri": "https://localhost:8080",
		"username": "test_user",
		"server_key": "0xbbbbbb",
		"stored_key": "0xcccccc"
	}`)

	// Create a mock request
	req, err := http.NewRequest("POST", "/api/v2/register-client", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	// Create a mock service and set it in the request context
	mockService := &MockService{
		registerClientv2: func(req dto.RegisterClientDTOv2) error {
			return fmt.Errorf("failed to register client")
		},
	}

	req = req.WithContext(context.WithValue(req.Context(), "service", mockService))

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the handler function
	Ctl.RegisterClientHandlerv2(rr, req)

	// Check the status code
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	// Decode the response body
	response := decodeResponseBodyForResponse(t, rr)

	// Now assert the fields directly
	assert.False(t, response.IsSuccess)
	assert.Equal(t, "Failed to register client", response.Message)
	assert.NotNil(t, response.Error)
}

func TestRegisterClientHandlerv2_Success(t *testing.T) {
	// Mock request body
	requestBody := []byte(`{
		"name": "test_client",
		"redirect_uri": "https://localhost:3000/callback",
		"backend_uri": "https://localhost:8080",
		"username": "test_user",
		"server_key": "0xbbbbbb",
		"stored_key": "0xcccccc"
	}`)

	// Create a mock request
	req, err := http.NewRequest("POST", "/api/v2/register-client", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	// Create a mock service and set it in the request context
	mockService := &MockService{
		registerClientv2: func(req dto.RegisterClientDTOv2) error {
			return nil
		},
	}

	req = req.WithContext(context.WithValue(req.Context(), "service", mockService))

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the handler function
	Ctl.RegisterClientHandlerv2(rr, req)

	// Check the status code
	assert.Equal(t, http.StatusCreated, rr.Code)

	// Decode the response body
	response := decodeResponseBodyForResponse(t, rr)

	// Now assert the fields directly
	assert.True(t, response.IsSuccess)
	assert.Equal(t, "Client registered successfully", response.Message)
	assert.Equal(t, nil, response.Data)
}
