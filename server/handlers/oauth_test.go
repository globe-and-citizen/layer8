package handlers_test

import (
	"bytes"
	"context"
	"fmt"
	"globe-and-citizen/layer8/server/entities"
	"globe-and-citizen/layer8/server/handlers"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

const clientUUID = "test_uuid"
const clientSecret = "test_secret"
const authorizationCode = "test_authorization_code"
const accessToken = "test_access_token"
const scopes = "country,email_verified"
const userID = 25
const userCountry = "test_country"

func TestTokenHandler_InvalidHttpMethod(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/api/token", nil)
	if err != nil {
		t.Fatal(err)
	}

	req = req.WithContext(context.WithValue(context.Background(), "Oauthservice", MockService{}))
	rr := httptest.NewRecorder()

	handlers.TokenHandler(rr, req)

	assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)

	response := decodeResponseBodyForResponse(t, rr)

	assert.False(t, response.IsSuccess)
	assert.NotNil(t, response.Error)
	assert.Equal(t, "invalid http method, expected: POST, got: GET", response.Message)
}

func TestTokenHandler_InvalidRequestJsonSchema(t *testing.T) {
	request := []byte(`{invalid_json}`)
	req, err := http.NewRequest(http.MethodPost, "/api/token", bytes.NewBuffer(request))
	if err != nil {
		t.Fatal(err)
	}

	req = req.WithContext(context.WithValue(context.Background(), "Oauthservice", MockService{}))
	rr := httptest.NewRecorder()

	handlers.TokenHandler(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	response := decodeResponseBodyForResponse(t, rr)

	assert.False(t, response.IsSuccess)
	assert.NotNil(t, response.Error)
	assert.Equal(t, "Request malformed: error while parsing json", response.Message)
}

func TestTokenHandler_RequiredRequestParametersMissing(t *testing.T) {
	request := []byte(`{"something": "else"}`)
	req, err := http.NewRequest(http.MethodPost, "/api/token", bytes.NewBuffer(request))
	if err != nil {
		t.Fatal(err)
	}

	req = req.WithContext(context.WithValue(context.Background(), "Oauthservice", MockService{}))
	rr := httptest.NewRecorder()

	handlers.TokenHandler(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	response := decodeResponseBodyForResponse(t, rr)

	assert.False(t, response.IsSuccess)
	assert.NotNil(t, response.Error)
	assert.Equal(t, "Input json is invalid", response.Message)
}

func TestTokenHandler_FailedToAuthenticateClient(t *testing.T) {
	request := []byte(`{
		"client_oauth_uuid": "test_uuid", 
		"client_oauth_secret": "test_secret",
		"authorization_code": "test_authorization_code"
	}`)
	req, err := http.NewRequest(http.MethodPost, "/api/token", bytes.NewBuffer(request))
	if err != nil {
		t.Fatal(err)
	}

	mockService := MockService{
		authenticateClient: func(uuid string, secret string) error {
			if uuid != clientUUID {
				t.Fatalf("Invalid uuid, expected: %s, got: %s", clientUUID, uuid)
			}
			if secret != clientSecret {
				t.Fatalf("Invalid secret, expected: %s, got: %s", clientSecret, secret)
			}

			return fmt.Errorf("failed to authenticate client")
		},
	}

	req = req.WithContext(context.WithValue(context.Background(), "Oauthservice", mockService))
	rr := httptest.NewRecorder()

	handlers.TokenHandler(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)

	response := decodeResponseBodyForResponse(t, rr)

	assert.False(t, response.IsSuccess)
	assert.NotNil(t, response.Error)
	assert.Equal(t, "failed to authenticate client", response.Message)
}

// func TestTokenHandler_FailedToVerifyAuthorizationCode(t *testing.T) {
// 	request := []byte(`{
// 		"client_oauth_uuid": "test_uuid",
// 		"client_oauth_secret": "test_secret",
// 		"authorization_code": "test_authorization_code"
// 	}`)
// 	req, err := http.NewRequest(http.MethodPost, "/api/token", bytes.NewBuffer(request))
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	mockService := MockService{
// 		authenticateClient: func(uuid string, secret string) error {
// 			if uuid != clientUUID {
// 				t.Fatalf("Invalid uuid, expected: %s, got: %s", clientUUID, uuid)
// 			}
// 			if secret != clientSecret {
// 				t.Fatalf("Invalid secret, expected: %s, got: %s", clientSecret, secret)
// 			}

// 			return nil
// 		},
// 		verifyAuthorizationCode: func(code string) error {
// 			if code != authorizationCode {
// 				t.Fatalf("Invalid authorization code, expected: %s, got: %s", authorizationCode, code)
// 			}
// 			return fmt.Errorf("invalid code")
// 		},
// 	}

// 	req = req.WithContext(context.WithValue(context.Background(), "Oauthservice", mockService))
// 	rr := httptest.NewRecorder()

// 	handlers.TokenHandler(rr, req)

// 	assert.Equal(t, http.StatusBadRequest, rr.Code)

// 	response := decodeResponseBodyForResponse(t, rr)

// 	assert.False(t, response.IsSuccess)
// 	assert.NotNil(t, response.Error)
// 	assert.Equal(t, "the authorization code is invalid", response.Message)
// }

// func TestTokenHandler_FailedToGenerateAccessToken(t *testing.T) {
// 	request := []byte(`{
// 		"client_oauth_uuid": "test_uuid",
// 		"client_oauth_secret": "test_secret",
// 		"authorization_code": "test_authorization_code"
// 	}`)
// 	req, err := http.NewRequest(http.MethodPost, "/api/token", bytes.NewBuffer(request))
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	mockService := MockService{
// 		authenticateClient: func(uuid string, secret string) error {
// 			if uuid != clientUUID {
// 				t.Fatalf("Invalid uuid, expected: %s, got: %s", clientUUID, uuid)
// 			}
// 			if secret != clientSecret {
// 				t.Fatalf("Invalid secret, expected: %s, got: %s", clientSecret, secret)
// 			}

// 			return nil
// 		},
// 		verifyAuthorizationCode: func(code string) error {
// 			if code != authorizationCode {
// 				t.Fatalf("Invalid authorization code, expected: %s, got: %s", authorizationCode, code)
// 			}
// 			return nil
// 		},
// 		generateAccessToken: func(uuid string, secret string) (string, error) {
// 			if uuid != clientUUID {
// 				t.Fatalf("Invalid uuid, expected: %s, got: %s", clientUUID, uuid)
// 			}
// 			if secret != clientSecret {
// 				t.Fatalf("Invalid secret, expected: %s, got: %s", clientSecret, secret)
// 			}

// 			return "", fmt.Errorf("failed to generate access token")
// 		},
// 	}

// 	req = req.WithContext(context.WithValue(context.Background(), "Oauthservice", mockService))
// 	rr := httptest.NewRecorder()

// 	handlers.TokenHandler(rr, req)

// 	assert.Equal(t, http.StatusInternalServerError, rr.Code)

// 	response := decodeResponseBodyForResponse(t, rr)

// 	assert.False(t, response.IsSuccess)
// 	assert.NotNil(t, response.Error)
// 	assert.Equal(t, "internal error when generating the access token", response.Message)
// }

// func TestTokenHandler_AccessTokenServedSuccessfully(t *testing.T) {
// 	request := []byte(`{
// 		"client_oauth_uuid": "test_uuid",
// 		"client_oauth_secret": "test_secret",
// 		"authorization_code": "test_authorization_code"
// 	}`)
// 	req, err := http.NewRequest(http.MethodPost, "/api/token", bytes.NewBuffer(request))
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	mockService := MockService{
// 		authenticateClient: func(uuid string, secret string) error {
// 			if uuid != clientUUID {
// 				t.Fatalf("Invalid uuid, expected: %s, got: %s", clientUUID, uuid)
// 			}
// 			if secret != clientSecret {
// 				t.Fatalf("Invalid secret, expected: %s, got: %s", clientSecret, secret)
// 			}

// 			return nil
// 		},
// 		verifyAuthorizationCode: func(code string) error {
// 			if code != authorizationCode {
// 				t.Fatalf("Invalid authorization code, expected: %s, got: %s", authorizationCode, code)
// 			}
// 			return nil
// 		},
// 		generateAccessToken: func(uuid string, secret string) (string, error) {
// 			if uuid != clientUUID {
// 				t.Fatalf("Invalid uuid, expected: %s, got: %s", clientUUID, uuid)
// 			}
// 			if secret != clientSecret {
// 				t.Fatalf("Invalid secret, expected: %s, got: %s", clientSecret, secret)
// 			}

// 			return accessToken, nil
// 		},
// 	}

// 	req = req.WithContext(context.WithValue(context.Background(), "Oauthservice", mockService))
// 	rr := httptest.NewRecorder()

// 	handlers.TokenHandler(rr, req)

// 	assert.Equal(t, http.StatusOK, rr.Code)

// 	response := decodeResponseBodyForResponse(t, rr)

// 	assert.True(t, response.IsSuccess)
// 	assert.Nil(t, response.Error)
// 	assert.Equal(t, "access token generated successfully", response.Message)

// 	resp := response.Data.(map[string]interface{})

// 	assert.Equal(t, accessToken, resp["access_token"])
// 	assert.Equal(t, "bearer", resp["token_type"])
// }

func TestZkMetadataHandler_InvalidHttpMethod(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/api/zk-metadata", nil)
	if err != nil {
		t.Fatal(err)
	}

	req = req.WithContext(context.WithValue(context.Background(), "Oauthservice", MockService{}))
	rr := httptest.NewRecorder()

	handlers.ZkMetadataHandler(rr, req)

	assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)

	response := decodeResponseBodyForResponse(t, rr)

	assert.False(t, response.IsSuccess)
	assert.NotNil(t, response.Error)
	assert.Equal(t, "invalid http method, expected: POST, got: GET", response.Message)
}

func TestZkMetadataHandler_InvalidRequestJsonSchema(t *testing.T) {
	request := []byte(`{invalid_json}`)
	req, err := http.NewRequest(http.MethodPost, "/api/zk-metadata", bytes.NewBuffer(request))
	if err != nil {
		t.Fatal(err)
	}

	req = req.WithContext(context.WithValue(context.Background(), "Oauthservice", MockService{}))
	rr := httptest.NewRecorder()

	handlers.ZkMetadataHandler(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	response := decodeResponseBodyForResponse(t, rr)

	assert.False(t, response.IsSuccess)
	assert.NotNil(t, response.Error)
	assert.Equal(t, "Request malformed: error while parsing json", response.Message)
}

func TestZkMetadataHandler_RequiredRequestParametersMissing(t *testing.T) {
	request := []byte(`{"something": "else"}`)
	req, err := http.NewRequest(http.MethodPost, "/api/zk-metadata", bytes.NewBuffer(request))
	if err != nil {
		t.Fatal(err)
	}

	req = req.WithContext(context.WithValue(context.Background(), "Oauthservice", MockService{}))
	rr := httptest.NewRecorder()

	handlers.ZkMetadataHandler(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	response := decodeResponseBodyForResponse(t, rr)

	assert.False(t, response.IsSuccess)
	assert.NotNil(t, response.Error)
	assert.Equal(t, "Input json is invalid", response.Message)
}

func TestZkMetadataHandler_NoAuthorizationHeader(t *testing.T) {
	request := []byte(`{
		"client_oauth_uuid": "test_uuid",
		"client_oauth_secret": "test_secret"
	}`)
	req, err := http.NewRequest(http.MethodPost, "/api/zk-metadata", bytes.NewBuffer(request))
	if err != nil {
		t.Fatal(err)
	}

	req = req.WithContext(context.WithValue(context.Background(), "Oauthservice", MockService{}))
	rr := httptest.NewRecorder()

	handlers.ZkMetadataHandler(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)

	response := decodeResponseBodyForResponse(t, rr)

	assert.False(t, response.IsSuccess)
	assert.NotNil(t, response.Error)
	assert.Equal(t, "invalid authorization header", response.Message)
}

func TestZkMetadataHandler_InvalidAuthorizationHeader(t *testing.T) {
	request := []byte(`{
		"client_oauth_uuid": "test_uuid",
		"client_oauth_secret": "test_secret"
	}`)
	req, err := http.NewRequest(http.MethodPost, "/api/zk-metadata", bytes.NewBuffer(request))
	if err != nil {
		t.Fatal(err)
	}

	req = req.WithContext(context.WithValue(context.Background(), "Oauthservice", MockService{}))
	req.Header.Set("Authorization", "some value")

	rr := httptest.NewRecorder()

	handlers.ZkMetadataHandler(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)

	response := decodeResponseBodyForResponse(t, rr)

	assert.False(t, response.IsSuccess)
	assert.NotNil(t, response.Error)
	assert.Equal(t, "invalid authorization header", response.Message)
}

func TestZkMetadataHandler_FailedToAuthenticateClient(t *testing.T) {
	request := []byte(`{
		"client_oauth_uuid": "test_uuid",
		"client_oauth_secret": "test_secret"
	}`)
	req, err := http.NewRequest(http.MethodPost, "/api/zk-metadata", bytes.NewBuffer(request))
	if err != nil {
		t.Fatal(err)
	}

	mockService := MockService{
		authenticateClient: func(uuid string, secret string) error {
			if uuid != clientUUID {
				t.Fatalf("Invalid uuid, expected: %s, got: %s", clientUUID, uuid)
			}
			if secret != clientSecret {
				t.Fatalf("Invalid secret, expected: %s, got: %s", clientSecret, secret)
			}

			return fmt.Errorf("error")
		},
	}

	req = req.WithContext(context.WithValue(context.Background(), "Oauthservice", mockService))
	req.Header.Set("Authorization", "Bearer "+accessToken)

	rr := httptest.NewRecorder()

	handlers.ZkMetadataHandler(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)

	response := decodeResponseBodyForResponse(t, rr)

	assert.False(t, response.IsSuccess)
	assert.NotNil(t, response.Error)
	assert.Equal(t, "Failed to authenticate client", response.Message)
}

func TestZkMetadataHandler_FailedToValidateAccessToken(t *testing.T) {
	request := []byte(`{
		"client_oauth_uuid": "test_uuid",
		"client_oauth_secret": "test_secret"
	}`)
	req, err := http.NewRequest(http.MethodPost, "/api/zk-metadata", bytes.NewBuffer(request))
	if err != nil {
		t.Fatal(err)
	}

	mockService := MockService{
		authenticateClient: func(uuid string, secret string) error {
			if uuid != clientUUID {
				t.Fatalf("Invalid uuid, expected: %s, got: %s", clientUUID, uuid)
			}
			if secret != clientSecret {
				t.Fatalf("Invalid secret, expected: %s, got: %s", clientSecret, secret)
			}

			return nil
		},
		validateAccessToken: func(secret string, token string) (*entities.ClientClaims, error) {
			if secret != clientSecret {
				t.Fatalf("Invalid secret, expected: %s, got: %s", clientSecret, secret)
			}
			if token != accessToken {
				t.Fatalf("Invalid access token, expected: %s, got: %s", accessToken, token)
			}
			return nil, fmt.Errorf("error")
		},
	}

	req = req.WithContext(context.WithValue(context.Background(), "Oauthservice", mockService))
	req.Header.Set("Authorization", "Bearer "+accessToken)

	rr := httptest.NewRecorder()

	handlers.ZkMetadataHandler(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)

	response := decodeResponseBodyForResponse(t, rr)

	assert.False(t, response.IsSuccess)
	assert.NotNil(t, response.Error)
	assert.Equal(t, "Failed to validate client access token", response.Message)
}

func TestZkMetadataHandler_FailedToGetZkUserMetadata(t *testing.T) {
	request := []byte(`{
		"client_oauth_uuid": "test_uuid",
		"client_oauth_secret": "test_secret"
	}`)
	req, err := http.NewRequest(http.MethodPost, "/api/zk-metadata", bytes.NewBuffer(request))
	if err != nil {
		t.Fatal(err)
	}

	mockService := MockService{
		authenticateClient: func(uuid string, secret string) error {
			if uuid != clientUUID {
				t.Fatalf("Invalid uuid, expected: %s, got: %s", clientUUID, uuid)
			}
			if secret != clientSecret {
				t.Fatalf("Invalid secret, expected: %s, got: %s", clientSecret, secret)
			}

			return nil
		},
		validateAccessToken: func(secret string, token string) (*entities.ClientClaims, error) {
			if secret != clientSecret {
				t.Fatalf("Invalid secret, expected: %s, got: %s", clientSecret, secret)
			}
			if token != accessToken {
				t.Fatalf("Invalid access token, expected: %s, got: %s", accessToken, token)
			}

			return &entities.ClientClaims{
				Scopes: scopes,
				UserID: userID,
			}, nil
		},
		getZkUserMetadata: func(scopesStr string, userId int64) (*entities.ZkMetadataResponse, error) {
			if scopesStr != scopes {
				t.Fatalf("Invalid scopes, expected: %s, got: %s", scopes, scopesStr)
			}
			if userId != userID {
				t.Fatalf("Invalid access token, expected: %d, got: %d", userID, userId)
			}

			return &entities.ZkMetadataResponse{}, fmt.Errorf("error")
		},
	}

	req = req.WithContext(context.WithValue(context.Background(), "Oauthservice", mockService))
	req.Header.Set("Authorization", "Bearer "+accessToken)

	rr := httptest.NewRecorder()

	handlers.ZkMetadataHandler(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	response := decodeResponseBodyForResponse(t, rr)

	assert.False(t, response.IsSuccess)
	assert.NotNil(t, response.Error)
	assert.Equal(t, "Failed to get user metadata", response.Message)
}

func TestZkMetadataHandler_UserZkMetadataServedSuccessfully(t *testing.T) {
	request := []byte(`{
		"client_oauth_uuid": "test_uuid",
		"client_oauth_secret": "test_secret"
	}`)
	req, err := http.NewRequest(http.MethodPost, "/api/zk-metadata", bytes.NewBuffer(request))
	if err != nil {
		t.Fatal(err)
	}

	mockService := MockService{
		authenticateClient: func(uuid string, secret string) error {
			if uuid != clientUUID {
				t.Fatalf("Invalid uuid, expected: %s, got: %s", clientUUID, uuid)
			}
			if secret != clientSecret {
				t.Fatalf("Invalid secret, expected: %s, got: %s", clientSecret, secret)
			}

			return nil
		},
		validateAccessToken: func(secret string, token string) (*entities.ClientClaims, error) {
			if secret != clientSecret {
				t.Fatalf("Invalid secret, expected: %s, got: %s", clientSecret, secret)
			}
			if token != accessToken {
				t.Fatalf("Invalid access token, expected: %s, got: %s", accessToken, token)
			}

			return &entities.ClientClaims{
				Scopes: scopes,
				UserID: userID,
			}, nil
		},
		getZkUserMetadata: func(scopesStr string, userId int64) (*entities.ZkMetadataResponse, error) {
			if scopesStr != scopes {
				t.Fatalf("Invalid scopes, expected: %s, got: %s", scopes, scopesStr)
			}
			if userId != userID {
				t.Fatalf("Invalid access token, expected: %d, got: %d", userID, userId)
			}

			return &entities.ZkMetadataResponse{
				Country:         userCountry,
				IsEmailVerified: true,
			}, nil
		},
	}

	req = req.WithContext(context.WithValue(context.Background(), "Oauthservice", mockService))
	req.Header.Set("Authorization", "Bearer "+accessToken)

	rr := httptest.NewRecorder()

	handlers.ZkMetadataHandler(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	response := decodeResponseBodyForResponse(t, rr)

	assert.True(t, response.IsSuccess)
	assert.Nil(t, response.Error)
	assert.Equal(t, "User metadata retrieved successfully", response.Message)

	resp := response.Data.(map[string]interface{})

	assert.Equal(t, userCountry, resp["country"])
	assert.Equal(t, true, resp["is_email_verified"])
}
