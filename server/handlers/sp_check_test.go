package handlers_test

import (
	"context"
	"globe-and-citizen/layer8/server/internals/repository"
	svc "globe-and-citizen/layer8/server/internals/service"
	"net/http"
	"net/http/httptest"
	"testing"

	Ctl "globe-and-citizen/layer8/server/handlers"

	"github.com/stretchr/testify/assert"
)

var mockService = &svc.Service{
	Repo: &repository.PostgresRepository{},
}

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
