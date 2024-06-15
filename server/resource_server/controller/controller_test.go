package controller_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"globe-and-citizen/layer8/server/resource_server/dto"
	"globe-and-citizen/layer8/server/resource_server/models"
	"globe-and-citizen/layer8/server/resource_server/utils"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	Ctl "globe-and-citizen/layer8/server/resource_server/controller"
)

var authenticationToken, _ = utils.GenerateToken(
	models.User{
		ID:       1,
		Username: "test_user",
	},
)

const verificationCode = "123467"
const emailProof = "email_proof"
const userEmail = "user@email.com"

func decodeResponseBody(t *testing.T, rr *httptest.ResponseRecorder) utils.Response {
	var response utils.Response
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}
	return response
}

// MockService implements interfaces.IService for testing purposes.
type MockService struct {
	verifyEmail                        func(userID uint, userEmail string) error
	checkEmailVerificationCode         func(userID uint, code string) error
	generateZkProofOfEmailVerification func(userID uint) (string, error)
	saveProofOfEmailVerification       func(userID uint, verificationCode string, zkProof string) error
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
	if userID == 1 {
		return models.ProfileResponseOutput{
			Username:    "test_user",
			FirstName:   "Test",
			LastName:    "User",
			DisplayName: "user",
			Country:     "Unknown",
		}, nil
	}
	return models.ProfileResponseOutput{}, fmt.Errorf("user not found")
}

func (ms *MockService) VerifyEmail(userID uint, userEmail string) error {
	return ms.verifyEmail(userID, userEmail)
}

func (ms *MockService) CheckEmailVerificationCode(userID uint, code string) error {
	return ms.checkEmailVerificationCode(userID, code)
}

func (ms *MockService) GenerateZkProofOfEmailVerification(userID uint) (string, error) {
	return ms.generateZkProofOfEmailVerification(userID)
}

func (ms *MockService) SaveProofOfEmailVerification(
	userID uint, verificationCode string, zkProof string,
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

// func (ms *MockService) LoginClient(req dto.LoginClientDTO) (models.LoginUserResponseOutput, error) {
// 	// Mock implementation for testing purposes.
// 	return models.LoginUserResponseOutput{}, nil
// }

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

func TestRegisterUserHandler(t *testing.T) {
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
	assert.Equal(t, http.StatusOK, rr.Code)

	// Decode the response body
	var response utils.Response
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}

	// Now assert the fields directly
	assert.True(t, response.Status)
	assert.Equal(t, "OK!", response.Message)
	assert.Nil(t, response.Error)
	assert.Equal(t, "User registered successfully", response.Data.(string))
}

func TestRegisterClientHandler(t *testing.T) {
	// Mock request body
	requestBody := []byte(`{"name": "testclient", "redirect_uri": "https://gcitizen.com/callback", "username": "test_user", "password": "12345"}`)

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
	assert.Equal(t, http.StatusOK, rr.Code)

	// Decode the response body
	var response utils.Response
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}

	// Now assert the fields directly
	assert.True(t, response.Status)
	assert.Equal(t, "OK!", response.Message)
	assert.Nil(t, response.Error)
	assert.Equal(t, "Client registered successfully", response.Data.(string))
}

func TestLoginPrecheckHandler(t *testing.T) {
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

func TestLoginUserHandler(t *testing.T) {
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

func TestProfileHandler(t *testing.T) {
	// Generate a Mock JWT token
	tokenString, err := utils.GenerateToken(models.User{
		ID:       1,
		Username: "test_user",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Create a mock request
	req, err := http.NewRequest("GET", "/api/v1/profile", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Set the Authorization header
	req.Header.Set("Authorization", "Bearer "+tokenString)

	// Create a mock service and set it in the request context
	mockService := &MockService{}
	req = req.WithContext(context.WithValue(req.Context(), "service", mockService))

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the handler function
	Ctl.ProfileHandler(rr, req)

	// Check the status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Decode the response body
	var response models.ProfileResponseOutput
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}

	// Now assert the fields directly
	assert.Equal(t, "test_user", response.Username)
	assert.Equal(t, "Test", response.FirstName)
	assert.Equal(t, "User", response.LastName)
	assert.Equal(t, "user", response.DisplayName)
	assert.Equal(t, "Unknown", response.Country)
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

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	response := decodeResponseBody(t, rr)

	assert.False(t, response.Status)
	assert.Equal(t, "Request failed: invalid authorization token", response.Message)
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

	response := decodeResponseBody(t, rr)

	assert.False(t, response.Status)
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

	response := decodeResponseBody(t, rr)

	assert.False(t, response.Status)
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

	response := decodeResponseBody(t, rr)

	assert.False(t, response.Status)
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

	response := decodeResponseBody(t, rr)

	assert.True(t, response.Status)
	assert.Equal(t, "OK!", response.Message)
	assert.Equal(t, "Verification email sent", response.Data)
	assert.Nil(t, response.Error)
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

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	response := decodeResponseBody(t, rr)

	assert.False(t, response.Status)
	assert.Equal(t, "Failed to verify user's token", response.Message)
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

	response := decodeResponseBody(t, rr)

	assert.False(t, response.Status)
	assert.Equal(t, "Error while unmarshalling json", response.Message)
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

	response := decodeResponseBody(t, rr)

	assert.False(t, response.Status)
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

	response := decodeResponseBody(t, rr)

	assert.False(t, response.Status)
	assert.Equal(t, "Failed to verify code", response.Message)
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
		generateZkProofOfEmailVerification: func(userID uint) (string, error) {
			return "", fmt.Errorf("failed to generate the zk email proof")
		},
	}
	req = req.WithContext(context.WithValue(req.Context(), "service", mockService))

	rr := httptest.NewRecorder()

	Ctl.CheckEmailVerificationCode(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	response := decodeResponseBody(t, rr)

	assert.False(t, response.Status)
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
		generateZkProofOfEmailVerification: func(userID uint) (string, error) {
			return emailProof, nil
		},
		saveProofOfEmailVerification: func(
			userID uint, verificationCode string, zkProof string,
		) error {
			if zkProof != emailProof {
				t.Fatalf("Email proof mismatch: expected %s, got %s", emailProof, zkProof)
			}
			return fmt.Errorf("failed to save proof of email verification")
		},
	}
	req = req.WithContext(context.WithValue(req.Context(), "service", mockService))

	rr := httptest.NewRecorder()

	Ctl.CheckEmailVerificationCode(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	response := decodeResponseBody(t, rr)

	assert.False(t, response.Status)
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
		generateZkProofOfEmailVerification: func(userID uint) (string, error) {
			return emailProof, nil
		},
		saveProofOfEmailVerification: func(
			userID uint, verificationCode string, zkProof string,
		) error {
			if zkProof != emailProof {
				t.Fatalf("Email proof mismatch: expected %s, got %s", emailProof, zkProof)
			}
			return nil
		},
	}
	req = req.WithContext(context.WithValue(req.Context(), "service", mockService))

	rr := httptest.NewRecorder()

	Ctl.CheckEmailVerificationCode(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	response := decodeResponseBody(t, rr)

	assert.True(t, response.Status)
	assert.Equal(t, "OK!", response.Message)
	assert.Equal(t, "Your email was successfully verified!", response.Data)
	assert.Nil(t, response.Error)
}

func TestUpdateDisplayNameHandler(t *testing.T) {
	// Generate a Mock JWT token
	tokenString, err := utils.GenerateToken(models.User{
		ID:       1,
		Username: "test_user",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Mock request body
	requestBody := []byte(`{"display_name": "test_user"}`)

	// Create a mock request
	req, err := http.NewRequest("POST", "/api/v1/update-display-name", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	// Set the Authorization header
	req.Header.Set("Authorization", "Bearer "+tokenString)

	// Create a mock service and set it in the request context
	mockService := &MockService{}
	req = req.WithContext(context.WithValue(req.Context(), "service", mockService))

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the handler function
	Ctl.UpdateDisplayNameHandler(rr, req)

	// Check the status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Decode the response body
	var response utils.Response
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}

	// Now assert the fields directly
	assert.True(t, response.Status)
	assert.Equal(t, "OK!", response.Message)
	assert.Nil(t, response.Error)
	assert.Equal(t, "Display name updated successfully", response.Data.(string))
}

// Javokhir started the testing
func (m *MockService) LoginClient(req dto.LoginClientDTO) (models.LoginUserResponseOutput, error) {
	// Mock implementation for LoginClient method
	return models.LoginUserResponseOutput{
		Token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImhtayIsInVzZXJfaWQiOjIsImlzcyI6Ikdsb2JlQW5kQ2l0aXplbiIsImV4cCI6MTcwNjUyNzY0NH0.AeQk23OPvlvauDEf45IlxxJ8ViSM5BlC6OlNkhXTomw",
	}, nil
}

func TestLoginClientHandler(t *testing.T) {
	// Prepare request body
	loginReq := dto.LoginClientDTO{
		Username: "testuser",
		Password: "testpassword",
	}
	reqBody, err := json.Marshal(loginReq)
	if err != nil {
		t.Fatal(err)
	}

	// Prepare request with request body
	req := httptest.NewRequest("POST", "/api/v1/login-client", bytes.NewBuffer(reqBody))

	// Set up mock service in request context
	req = setMockServiceInContext(req)

	// Create a response recorder to capture the handler's response
	w := httptest.NewRecorder()

	// Call the handler function
	Ctl.LoginClientHandler(w, req)

	// Check the response status code
	assert.Equal(t, http.StatusOK, w.Code)

	// Decode the response body
	var tokenResp models.LoginUserResponseOutput
	err = json.NewDecoder(w.Body).Decode(&tokenResp)
	if err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	// Validate the response
	assert.Equal(t, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImhtayIsInVzZXJfaWQiOjIsImlzcyI6Ikdsb2JlQW5kQ2l0aXplbiIsImV4cCI6MTcwNjUyNzY0NH0.AeQk23OPvlvauDEf45IlxxJ8ViSM5BlC6OlNkhXTomw", tokenResp.Token)
}

func TestCheckBackendURIHandler(t *testing.T) {
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
// Javokhir finished the testing
