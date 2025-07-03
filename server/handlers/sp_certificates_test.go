package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	utilities "github.com/globe-and-citizen/layer8-utils"
	"github.com/golang-jwt/jwt/v4"
	"globe-and-citizen/layer8/server/entities"
	"globe-and-citizen/layer8/server/internals/repository"
	svc "globe-and-citizen/layer8/server/internals/service"
	"globe-and-citizen/layer8/server/models"
	"globe-and-citizen/layer8/server/resource_server/utils"
	"golang.org/x/oauth2"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	Ctl "globe-and-citizen/layer8/server/handlers"

	"github.com/stretchr/testify/assert"
	resourceModels "globe-and-citizen/layer8/server/resource_server/models"
)

var mockService = &svc.Service{
	Repo: &repository.PostgresRepository{},
}

type MockService struct {
	getUserByToken           func(token string) (*models.User, error)
	loginUser                func(username, password string) (map[string]interface{}, error)
	generateAuthorizationURL func(config *oauth2.Config, userID int64, headerMap map[string]string) (*entities.AuthURL, error)
	exchangeCodeForToken     func(config *oauth2.Config, code string) (*oauth2.Token, error)
	accessResourcesWithToken func(token string) (map[string]interface{}, error)
	getClient                func(id string) (*models.Client, error)
	verifyToken              func(token string) (isvalid bool, err error)
	checkClient              func(backendURL string) (*models.Client, error)
	saveX509Certificate      func(clientID string, certificate string) error
	addTestClient            func() (*models.Client, error)
}

func (m MockService) GetUserByToken(token string) (*models.User, error) {
	return m.getUserByToken(token)
}

func (m MockService) LoginUser(username, password string) (map[string]interface{}, error) {
	return m.loginUser(username, password)
}

func (m MockService) GenerateAuthorizationURL(config *oauth2.Config, userID int64, headerMap map[string]string) (*entities.AuthURL, error) {
	return m.generateAuthorizationURL(config, userID, headerMap)
}

func (m MockService) ExchangeCodeForToken(config *oauth2.Config, code string) (*oauth2.Token, error) {
	return m.exchangeCodeForToken(config, code)
}

func (m MockService) AccessResourcesWithToken(token string) (map[string]interface{}, error) {
	return m.accessResourcesWithToken(token)
}

func (m MockService) GetClient(id string) (*models.Client, error) {
	return m.getClient(id)
}

func (m MockService) VerifyToken(token string) (isvalid bool, err error) {
	return m.verifyToken(token)
}

func (m MockService) CheckClient(backendURL string) (*models.Client, error) {
	return m.checkClient(backendURL)
}

func (m MockService) SaveX509Certificate(clientID string, certificate string) error {
	return m.saveX509Certificate(clientID, certificate)
}

func (m MockService) AddTestClient() (*models.Client, error) {
	return m.addTestClient()
}

const clientId = "clientID"
const backendUrl = "backend URL"
const certificate = "x509 certificate"

const jwtSecretKey = "JwtSecretkey"

func Test_GetSPPubKeyHandler_InvalidHttpRequestMethod(t *testing.T) {
	context := context.WithValue(context.Background(), "Oauthservice", mockService)
	req, err := http.NewRequest("POST", "/sp-pub-key", nil)
	if err != nil {
		t.Fatal(err)
	}
	req = req.WithContext(context)

	rr := httptest.NewRecorder()

	Ctl.GetSPPubKey(rr, req)

	assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)
	assert.Equal(t, `{"error": "method not allowed"}`, rr.Body.String())
}

func Test_GetSPPubKeyHandler_MissingBackendURL(t *testing.T) {
	context := context.WithValue(context.Background(), "Oauthservice", mockService)
	req, err := http.NewRequest("GET", "/sp-pub-key?backend_url=", nil)
	if err != nil {
		t.Fatal(err)
	}
	req = req.WithContext(context)

	rr := httptest.NewRecorder()

	Ctl.GetSPPubKey(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Equal(t, `{"error": "missing backend_url parameter"}`, rr.Body.String())
}

func Test_GetSPPubKeyHandler_MissingToken(t *testing.T) {
	context := context.WithValue(context.Background(), "Oauthservice", mockService)
	req, err := http.NewRequest("GET", "/sp-pub-key?backend_url=http://example.com", nil)
	if err != nil {
		t.Fatal(err)
	}
	req = req.WithContext(context)

	rr := httptest.NewRecorder()

	Ctl.GetSPPubKey(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	assert.Equal(t, `{"error": "missing token"}`, rr.Body.String())
}

func Test_GetSPPubKeyHandler_InvalidToken(t *testing.T) {
	context := context.WithValue(context.Background(), "Oauthservice", mockService)
	req, err := http.NewRequest("GET", "/sp-pub-key?backend_url=http://example.com", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer invalid_token")
	req = req.WithContext(context)

	rr := httptest.NewRecorder()

	Ctl.GetSPPubKey(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	assert.Equal(t, `{"error": "invalid token"}`, rr.Body.String())
}

func Test_GetSPPubKeyHandler_CertificateServedSuccessfully(t *testing.T) {
	mockService := MockService{
		checkClient: func(backendURL string) (*models.Client, error) {
			if backendURL != backendUrl {
				t.Fatalf("invalid backend url, expected: %s, got: %s", backendUrl, backendURL)
			}

			return &models.Client{
				ID:                   clientId,
				BackendURI:           backendUrl,
				X509CertificateBytes: []byte(certificate),
			}, nil
		},
		verifyToken: func(token string) (isvalid bool, err error) {
			return true, nil
		},
	}

	os.Setenv("JWT_SECRET_KEY", jwtSecretKey)
	jwtToken, _ := utilities.GenerateStandardToken(jwtSecretKey)

	ctx := context.WithValue(context.Background(), "Oauthservice", mockService)
	req, err := http.NewRequest("GET", fmt.Sprintf("/sp-pub-key?backend_url=%s", backendUrl), nil)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Authorization", "Bearer "+jwtToken)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	Ctl.GetSPPubKey(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response entities.X509CertificateResponse
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response body: %e", err)
	}

	assert.Equal(t, certificate, response.X509Certificate)
}

func TestUploadSPCertificate_InvalidHTTPMethod(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/api/upload-certificate", nil)
	if err != nil {
		t.Fatal(err)
	}

	req = req.WithContext(context.WithValue(context.Background(), "Oauthservice", MockService{}))
	rr := httptest.NewRecorder()

	Ctl.UploadSPCertificate(rr, req)

	assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)
}

func TestUploadSPCertificate_AuthorizationTokenIsMissing(t *testing.T) {
	request := []byte(`{"certificate": "x509 certificate"}`)
	req, err := http.NewRequest(http.MethodPost, "/api/upload-certificate", bytes.NewBuffer(request))
	if err != nil {
		t.Fatal(err)
	}

	req = req.WithContext(context.WithValue(context.Background(), "Oauthservice", MockService{}))
	rr := httptest.NewRecorder()

	Ctl.UploadSPCertificate(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestUploadSPCertificate_InvalidAuthorizationToken(t *testing.T) {
	request := []byte(`{"certificate": "x509 certificate"}`)
	req, err := http.NewRequest(http.MethodPost, "/api/upload-certificate", bytes.NewBuffer(request))
	if err != nil {
		t.Fatal(err)
	}

	req = req.WithContext(context.WithValue(context.Background(), "Oauthservice", MockService{}))
	req.Header.Set("Authorization", "Bearer invalid token")

	rr := httptest.NewRecorder()

	Ctl.UploadSPCertificate(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestUploadSPCertificate_InvalidRequestSchema(t *testing.T) {
	request := []byte(`{somethingelse}`)
	req, err := http.NewRequest(http.MethodPost, "/api/upload-certificate", bytes.NewBuffer(request))
	if err != nil {
		t.Fatal(err)
	}

	req = req.WithContext(context.WithValue(context.Background(), "Oauthservice", MockService{}))

	os.Setenv("JWT_SECRET_KEY", jwtSecretKey)
	req.Header.Set("Authorization", "Bearer "+generateJwtToken())

	rr := httptest.NewRecorder()

	Ctl.UploadSPCertificate(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	response := decodeResponseBodyForResponse(t, rr)

	assert.False(t, response.IsSuccess)
	assert.Equal(t, "Request malformed: error while parsing json", response.Message)
	assert.NotNil(t, response.Error)
}

func TestUploadSPCertificate_RequiredRequestParameterIsMissing(t *testing.T) {
	os.Setenv("JWT_SECRET_KEY", jwtSecretKey)

	request := []byte(`{"something_else": "some_value"}`)
	req, err := http.NewRequest(http.MethodPost, "/api/upload-certificate", bytes.NewBuffer(request))
	if err != nil {
		t.Fatal(err)
	}

	req = req.WithContext(context.WithValue(context.Background(), "Oauthservice", MockService{}))
	req.Header.Set("Authorization", "Bearer "+generateJwtToken())

	rr := httptest.NewRecorder()

	Ctl.UploadSPCertificate(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	response := decodeResponseBodyForResponse(t, rr)

	assert.False(t, response.IsSuccess)
	assert.Equal(t, "Input json is invalid", response.Message)
	assert.NotNil(t, response.Error)
}

func TestUploadSPCertificate_FailedToSaveClientCertificate(t *testing.T) {
	os.Setenv("JWT_SECRET_KEY", jwtSecretKey)

	request := []byte(`{"certificate": "x509 certificate"}`)
	req, err := http.NewRequest(http.MethodPost, "/api/upload-certificate", bytes.NewBuffer(request))
	if err != nil {
		t.Fatal(err)
	}

	service := MockService{
		saveX509Certificate: func(clientID string, x509Certificate string) error {
			if clientID != clientId {
				t.Fatalf("Invalid clientID, expected: %s, got: %s", clientId, clientID)
			}
			if x509Certificate != certificate {
				t.Fatalf("Invalid certificate, expected: %s, got: %s", certificate, x509Certificate)
			}
			return fmt.Errorf("failed to save certificate")
		},
	}

	req = req.WithContext(context.WithValue(context.Background(), "Oauthservice", service))
	req.Header.Set("Authorization", "Bearer "+generateJwtToken())

	rr := httptest.NewRecorder()

	Ctl.UploadSPCertificate(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	response := decodeResponseBodyForResponse(t, rr)

	assert.False(t, response.IsSuccess)
	assert.Equal(t, "failed to save the SP x.509 certificate", response.Message)
	assert.NotNil(t, response.Error)
}

func TestUploadSPCertificate_ClientCertificateSavedSuccessfully(t *testing.T) {
	os.Setenv("JWT_SECRET_KEY", jwtSecretKey)

	request := []byte(`{"certificate": "x509 certificate"}`)
	req, err := http.NewRequest(http.MethodPost, "/api/upload-certificate", bytes.NewBuffer(request))
	if err != nil {
		t.Fatal(err)
	}

	service := MockService{
		saveX509Certificate: func(clientID string, x509Certificate string) error {
			if clientID != clientId {
				t.Fatalf("Invalid clientID, expected: %s, got: %s", clientId, clientID)
			}
			if x509Certificate != certificate {
				t.Fatalf("Invalid certificate, expected: %s, got: %s", certificate, x509Certificate)
			}
			return nil
		},
	}

	req = req.WithContext(context.WithValue(context.Background(), "Oauthservice", service))
	req.Header.Set("Authorization", "Bearer "+generateJwtToken())

	rr := httptest.NewRecorder()

	Ctl.UploadSPCertificate(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)

	response := decodeResponseBodyForResponse(t, rr)

	assert.True(t, response.IsSuccess)
	assert.Equal(t, "x.509 certificate was saved successfully", response.Message)
	assert.Nil(t, response.Error)
}

func generateJwtToken() string {
	claims := &resourceModels.ClientClaims{
		UserName: "username",
		ClientID: clientId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(60 * time.Minute)),
			Issuer:    "GlobeAndCitizen",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtSecretKey))
	if err != nil {
		log.Fatal(err)
	}

	return tokenString
}

func decodeResponseBodyForResponse(t *testing.T, rr *httptest.ResponseRecorder) utils.Response {
	var response utils.Response
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}
	return response
}
