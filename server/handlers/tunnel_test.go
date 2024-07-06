package handlers

import (
	"context"
	"encoding/json"
	"globe-and-citizen/layer8/server/models"
	"globe-and-citizen/layer8/server/resource_server/emails/verification"
	"globe-and-citizen/layer8/server/resource_server/repository"
	resourceService "globe-and-citizen/layer8/server/resource_server/service"
	resourceUtils "globe-and-citizen/layer8/server/resource_server/utils"
	"globe-and-citizen/layer8/server/resource_server/utils/mocks"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	utils "github.com/globe-and-citizen/layer8-utils"
)

func prepareInitTunnelRequest(t *testing.T, clientBackendUrl string) *http.Request {
	rmSalt := resourceUtils.GenerateRandomSalt(resourceUtils.SaltSize)
	repo := repository.NewMemoryRepository()
	client := &models.Client{
		ID:          "notanid",
		Secret:      "absolutelynotasecret!",
		Name:        "Ex-C",
		RedirectURI: "http://localhost:5173/oauth2/callback",
		BackendURI:  resourceUtils.RemoveProtocolFromURL(clientBackendUrl),
		Username:    "layer8",
		Password:    resourceUtils.SaltAndHashPassword("12341234", rmSalt),
		Salt:        rmSalt,
	}

	repo.SetClient(client)

	ctrl := gomock.NewController(t)
	payAsYouGoWrapper := mocks.NewMockPayAsYouGoWrapper(ctrl)

	reqToInitTunnel := httptest.NewRequest("GET", "/init-tunnel", nil)
	reqToInitTunnel = reqToInitTunnel.WithContext(
		context.WithValue(
			reqToInitTunnel.Context(),
			"service",
			resourceService.NewService(repo, &verification.EmailVerifier{}, payAsYouGoWrapper),
		),
	)

	queryParams := reqToInitTunnel.URL.Query()
	queryParams.Add("backend", clientBackendUrl)
	reqToInitTunnel.URL.RawQuery = queryParams.Encode()

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

	reqToInitTunnel := prepareInitTunnelRequest(t, mockedServiceProvider.URL)
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
	reqToInitTunnel := prepareInitTunnelRequest(t, "http://localhost:8080")
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

func TestTunnelReturnsExactCodeFromSP(t *testing.T) {
	os.Setenv("MP_123_SECRET_KEY", "mp_secret_key")
	os.Setenv("UP_999_SECRET_KEY", "up_secret_key")

	mpJWT, err := utils.GenerateStandardToken(os.Getenv("MP_123_SECRET_KEY"))
	assert.Nil(t, err)
	upJWT, err := resourceUtils.GenerateUPTokenJWT(os.Getenv("UP_999_SECRET_KEY"), "client_id")
	assert.Nil(t, err)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		paths := []struct {
			Path string
			Code int
			Body string
		}{
			{"/internal-error", http.StatusInternalServerError, `{"error": "internal server error"}`},
			{"/forbidden", http.StatusForbidden, `{"error": "forbidden"}`},
			{"/unauthorized", http.StatusUnauthorized, `{"error": "unauthorized"}`},
		}

		for _, path := range paths {
			if r.URL.Path == path.Path {
				w.Header().Set("Content-Type", "application/json")
				w.Header().Set("mp-jwt", mpJWT)
				w.WriteHeader(path.Code)
				w.Write([]byte(path.Body))
				return
			}
		}

		// default not found
		http.Error(w, "not found", http.StatusNotFound)
	}))

	cases := []struct {
		Path string
		Code int
	}{
		{"/internal-error", http.StatusInternalServerError},
		{"/forbidden", http.StatusForbidden},
		{"/unauthorized", http.StatusUnauthorized},
		{"/not-found", http.StatusNotFound},
	}

	for _, c := range cases {
		req := httptest.NewRequest("GET", c.Path, nil)
		req.Header.Set("X-Forwarded-Proto", "http")
		req.Header.Set("X-Forwarded-Host", strings.Replace(server.URL, "http://", "", 1))
		req.Header.Set("up-jwt", upJWT)
		res := httptest.NewRecorder()

		Tunnel(res, req)

		assert.Equal(t, c.Code, res.Code)
	}
}
