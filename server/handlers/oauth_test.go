package handlers_test

import (
	"bytes"
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"globe-and-citizen/layer8/server/handlers"
	"net/http"
	"net/http/httptest"
	"testing"
)

const clientUUID = "test_uuid"
const clientSecret = "test_secret"
const authorizationCode = "test_authorization_code"

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

func TestTokenHandler_FailedToVerifyAuthorizationCode(t *testing.T) {
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

			return nil
		},
		verifyAuthorizationCode: func(code string) error {
			if code != authorizationCode {
				t.Fatalf("Invalid authorization code, expected: %s, got: %s", authorizationCode, code)
			}
			return fmt.Errorf("invalid code")
		},
	}

	req = req.WithContext(context.WithValue(context.Background(), "Oauthservice", mockService))
	rr := httptest.NewRecorder()

	handlers.TokenHandler(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	response := decodeResponseBodyForResponse(t, rr)

	assert.False(t, response.IsSuccess)
	assert.NotNil(t, response.Error)
	assert.Equal(t, "the authorization code is invalid", response.Message)
}

func TestTokenHandler_FailedToGenerateAccessToken(t *testing.T) {

}

func TestTokenHandler_AccessTokenServedSuccessfully(t *testing.T) {

}
