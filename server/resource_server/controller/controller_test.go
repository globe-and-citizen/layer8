package controller_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"globe-and-citizen/layer8/server/resource_server/dto"
	"globe-and-citizen/layer8/server/resource_server/models"
	"globe-and-citizen/layer8/server/resource_server/utils"
	l8utils "globe-and-citizen/layer8/server/utils"
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
const userSalt = "salt"

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
	generateZkProofOfEmailVerification func(user models.User, request dto.CheckEmailVerificationCodeDTO) ([]byte, error)
	saveProofOfEmailVerification       func(userID uint, verificationCode string, zkProof []byte) error
	profileUser                        func(userID uint) (models.ProfileResponseOutput, error)
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

func (ms *MockService) LoginUser(req dto.LoginUserDTO) (models.LoginUserResponseOutput, error) {
	// Mock implementation for testing purposes.
	return models.LoginUserResponseOutput{
		Token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImhtayIsInVzZXJfaWQiOjIsImlzcyI6Ikdsb2JlQW5kQ2l0aXplbiIsImV4cCI6MTcwNjUyNzY0NH0.AeQk23OPvlvauDEf45IlxxJ8ViSM5BlC6OlNkhXTomw",
	}, nil
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
) ([]byte, error) {
	return ms.generateZkProofOfEmailVerification(user, request)
}

func (ms *MockService) SaveProofOfEmailVerification(
	userID uint, verificationCode string, zkProof []byte,
) error {
	return ms.saveProofOfEmailVerification(userID, verificationCode, zkProof)
}

func (ms *MockService) UpdateDisplayName(userID uint, req dto.UpdateDisplayNameDTO) error {
	// Mock implementation for testing purposes.
	return nil
}

func (ms *MockService) RegisterClient(req dto.RegisterClientDTO) error {
	// Mock implementation for testing purposes.
	return nil
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
		Token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImhtayIsInVzZXJfaWQiOjIsImlzcyI6Ikdsb2JlQW5kQ2l0aXplbiIsImV4cCI6MTcwNjUyNzY0NH0.AeQk23OPvlvauDEf45IlxxJ8ViSM5BlC6OlNkhXTomw",
	}, nil
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

	rr := httptest.NewRecorder()

	Ctl.RegisterUserHandler(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	response := decodeResponseBody(t, rr)

	assert.False(t, response.Status)
	assert.Equal(t, "Request malformed: error while parsing json", response.Message)
	assert.NotNil(t, response.Error)
}

func TestRegisterUserHandler_Success(t *testing.T) {
	// Mock request body
	requestBody := []byte(`{
		"email": "test@gcitizen.com",
		"username": "test_user",
		"first_name": "Test",
		"last_name": "User",
		"display_name": "user",
		"country": "Unknown",
		"password": "12345"
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
	response := decodeResponseBody(t, rr)

	// Now assert the fields directly
	assert.True(t, response.IsSuccess)
	assert.Equal(t, "User registered successfully", response.Message)
	assert.Equal(t, map[string]interface{}{}, response.Data)
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

	response := decodeResponseBody(t, rr)

	assert.False(t, response.Status)
	assert.Equal(t, "Request malformed: error while parsing json", response.Message)
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
	assert.Equal(t, map[string]interface{}{}, response.Data)
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

	response := decodeResponseBody(t, rr)

	assert.False(t, response.Status)
	assert.Equal(t, "Request malformed: error while parsing json", response.Message)
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

	response := decodeResponseBody(t, rr)

	assert.False(t, response.Status)
	assert.Equal(t, "Request malformed: error while parsing json", response.Message)
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
	assert.Equal(t, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImhtayIsInVzZXJfaWQiOjIsImlzcyI6Ikdsb2JlQW5kQ2l0aXplbiIsImV4cCI6MTcwNjUyNzY0NH0.AeQk23OPvlvauDEf45IlxxJ8ViSM5BlC6OlNkhXTomw", response.Token)
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

	response := decodeResponseBody(t, rr)

	assert.False(t, response.Status)
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

	response := decodeResponseBody(t, rr)

	assert.False(t, response.Status)
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

func TestGetClientData(t *testing.T) {
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
	assert.Equal(t, map[string]interface{}{}, response.Data)
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

	response := decodeResponseBody(t, rr)

	assert.False(t, response.Status)
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
		) ([]byte, error) {
			return []byte{}, fmt.Errorf("failed to generate the zk email proof")
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
		) ([]byte, error) {
			return emailProof, nil
		},
		saveProofOfEmailVerification: func(
			userID uint, verificationCode string, zkProof []byte,
		) error {
			if !l8utils.Equal(zkProof, emailProof) {
				t.Fatalf("Email proof mismatch: expected %s, got %s", emailProof, zkProof)
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
		) ([]byte, error) {
			return emailProof, nil
		},
		saveProofOfEmailVerification: func(
			userID uint, verificationCode string, zkProof []byte,
		) error {
			if !l8utils.Equal(zkProof, emailProof) {
				t.Fatalf("Email proof mismatch: expected %s, got %s", emailProof, zkProof)
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
	assert.Equal(t, map[string]interface{}{}, response.Data)
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

	response := decodeResponseBody(t, rr)

	assert.False(t, response.Status)
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

	response := decodeResponseBody(t, rr)

	assert.False(t, response.Status)
	assert.Equal(t, "Request malformed: error while parsing json", response.Message)
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

	var response = decodeResponseBody(t, rr)

	assert.True(t, response.IsSuccess)
	assert.Equal(t, "Display name updated successfully", response.Message)
	assert.Equal(t, map[string]interface{}{}, response.Data)
  assert.Nil(t, response.Error)
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
	response := decodeResponseBody(t, rr)

	assert.False(t, response.Status)
	assert.Equal(t, "Request malformed: error while parsing json", response.Message)
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
	assert.Equal(t, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImhtayIsInVzZXJfaWQiOjIsImlzcyI6Ikdsb2JlQW5kQ2l0aXplbiIsImV4cCI6MTcwNjUyNzY0NH0.AeQk23OPvlvauDEf45IlxxJ8ViSM5BlC6OlNkhXTomw", tokenResp.Token)
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

	response := decodeResponseBody(t, rr)

	assert.False(t, response.Status)
	assert.Equal(t, "Request malformed: error while parsing json", response.Message)
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

func setMockServiceInContext(req *http.Request) *http.Request {
	mockSvc := &MockService{}
	ctx := context.WithValue(req.Context(), "service", mockSvc)
	return req.WithContext(ctx)
}
