package handlers

import (
	"context"
	"encoding/json"
	"globe-and-citizen/layer8/server/resource_server/dto"
	"globe-and-citizen/layer8/server/resource_server/repository"
	resourceService "globe-and-citizen/layer8/server/resource_server/service"
	resourceUtils "globe-and-citizen/layer8/server/resource_server/utils"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	utils "github.com/globe-and-citizen/layer8-utils"
)

func makeInitTunnelRequest(clientBackendUrl string) *http.Request {
	repo := repository.NewMemoryRepository()
	repo.RegisterClient(dto.RegisterClientDTO{
		Name:        "name",
		RedirectURI: "redirect_uri",
		BackendURI:  resourceUtils.RemoveProtocolFromURL(clientBackendUrl),
		Username:    "username",
		Password:    "password",
	})

	reqToInitTunnel := httptest.NewRequest("GET", "/api/tunnel", nil)
	reqToInitTunnel = reqToInitTunnel.WithContext(
		context.WithValue(reqToInitTunnel.Context(), "service", resourceService.NewService(repo)),
	)

	return reqToInitTunnel
}

func Test_InitTunnel_OK(t *testing.T) {
	var (
		MockedMp123SecretKey = "MOCKED_MP_123_SECRET_KEY"
		MockedUp999SecretKey = "MOCKED_UP_999_SECRET_KEY"
	)

	os.Setenv("MP_123_SECRET_KEY", MockedMp123SecretKey)
	os.Setenv("UP_999_SECRET_KEY", MockedUp999SecretKey)

	mockedServiceProvider := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, pubJWK_ecdh, err := utils.GenerateKeyPair(utils.ECDH)
		if err != nil {
			return
		}

		b64PubJWK, err := pubJWK_ecdh.ExportAsBase64()
		if err != nil {
			return
		}

		w.Write([]byte(b64PubJWK))
	}))

	clientBackendUrl := mockedServiceProvider.URL

	reqToInitTunnel := makeInitTunnelRequest(clientBackendUrl)

	queryParams := reqToInitTunnel.URL.Query()
	queryParams.Add("backend", clientBackendUrl)
	reqToInitTunnel.URL.RawQuery = queryParams.Encode()

	responseRecorder := httptest.NewRecorder()

	InitTunnel(responseRecorder, reqToInitTunnel)

	assert.Equal(t, http.StatusOK, responseRecorder.Code)

	resBody, err := io.ReadAll(responseRecorder.Body)
	if err != nil {
		t.Fatal(err)
	}

	var res map[string]interface{}
	if err := json.Unmarshal(resBody, &res); err != nil {
		t.Fatal(err)
	}

	assert.NotEmpty(t, res["up-JWT"])
	assert.NotEmpty(t, res["server_pubKeyECDH"])

	upJWT, ok := res["up-JWT"].(string)
	if !ok {
		t.Fatal("up-JWT is not a string")
	}

	_, err = resourceUtils.ValidateUPTokenJWT(upJWT, MockedUp999SecretKey)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_InitTunnel_InvalidBackendURL(t *testing.T) {
	reqToInitTunnel := httptest.NewRequest("GET", "/api/tunnel", nil)
	responseRecorder := httptest.NewRecorder()

	InitTunnel(responseRecorder, reqToInitTunnel)
	resBody, err := io.ReadAll(responseRecorder.Body)
	if err != nil {
		t.Fatal(err)
	}

	var res map[string]interface{}
	if err := json.Unmarshal(resBody, &res); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "Failed to get User. Malformed query string.", res["message"])
}

func Test_InitTunnel_UnavailableBackend(t *testing.T) {
	clientBackendUrl := "http://localhost:8080"

	reqToInitTunnel := makeInitTunnelRequest(clientBackendUrl)
	queryParams := reqToInitTunnel.URL.Query()
	queryParams.Add("backend", clientBackendUrl)
	reqToInitTunnel.URL.RawQuery = queryParams.Encode()

	responseRecorder := httptest.NewRecorder()

	InitTunnel(responseRecorder, reqToInitTunnel)

	assert.Equal(t, http.StatusInternalServerError, responseRecorder.Code)
}

func Test_TunnelAPI_OK(t *testing.T) {
	var (
		MOCKED_MP_123_SECRET_KEY = "MOCKED_MP_123_SECRET_KEY"
		MOCKED_UP_999_SECRET_KEY = "MOCKED_UP_999_SECRET_KEY"
	)

	os.Setenv("MP_123_SECRET_KEY", MOCKED_MP_123_SECRET_KEY)
	os.Setenv("UP_999_SECRET_KEY", MOCKED_UP_999_SECRET_KEY)

	mpJWT, _ := utils.GenerateStandardToken(os.Getenv("MP_123_SECRET_KEY"))
	upJWT, _ := resourceUtils.GenerateUPTokenJWT(os.Getenv("UP_999_SECRET_KEY"), "client_id")

	mockedServiceProvider := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/ping" {
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("mp-jwt", mpJWT)
			w.Header().Set("trace-id", "trace-id-mock")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"message": "pong"}`))
		}
	}))

	reqToInitTunnel := httptest.NewRequest("GET", "/ping", nil)
	reqToInitTunnel.Header.Set("X-Forwarded-Proto", "http")
	reqToInitTunnel.Header.Set("X-Forwarded-Host", strings.Replace(mockedServiceProvider.URL, "http://", "", 1))
	reqToInitTunnel.Header.Set("up-jwt", upJWT)
	responseRecorder := httptest.NewRecorder()

	Tunnel(responseRecorder, reqToInitTunnel)

	assert.Equal(t, http.StatusOK, responseRecorder.Code)
	assert.Equal(t, "trace-id-mock", responseRecorder.Header().Get("trace-id"))
	assert.Equal(t, `{"message": "pong"}`, responseRecorder.Body.String())
}

func Test_TunnelAPI_OK_BadRequest(t *testing.T) {
	var (
		MOCKED_MP_123_SECRET_KEY = "MOCKED_MP_123_SECRET_KEY"
		MOCKED_UP_999_SECRET_KEY = "MOCKED_UP_999_SECRET_KEY"
	)

	os.Setenv("MP_123_SECRET_KEY", MOCKED_MP_123_SECRET_KEY)
	os.Setenv("UP_999_SECRET_KEY", MOCKED_UP_999_SECRET_KEY)

	mpJWT, _ := utils.GenerateStandardToken(os.Getenv("MP_123_SECRET_KEY"))
	upJWT, _ := resourceUtils.GenerateUPTokenJWT(os.Getenv("UP_999_SECRET_KEY"), "client_id")

	mockedServiceProvider := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/ping" {
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("mp-jwt", mpJWT)
			w.Header().Set("trace-id", "trace-id-mock")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error": "invalid value"}`))
		}
	}))

	reqToInitTunnel := httptest.NewRequest("GET", "/ping", nil)
	reqToInitTunnel.Header.Set("X-Forwarded-Proto", "http")
	reqToInitTunnel.Header.Set("X-Forwarded-Host", strings.Replace(mockedServiceProvider.URL, "http://", "", 1))
	reqToInitTunnel.Header.Set("up-jwt", upJWT)
	responseRecorder := httptest.NewRecorder()

	Tunnel(responseRecorder, reqToInitTunnel)

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
	assert.Equal(t, "trace-id-mock", responseRecorder.Header().Get("trace-id"))
	assert.Equal(t, `{"error": "invalid value"}`, responseRecorder.Body.String())
}
