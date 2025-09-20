package handlers

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"globe-and-citizen/layer8/server/resource_server/dto"
	"globe-and-citizen/layer8/server/resource_server/emails/verification"
	"globe-and-citizen/layer8/server/resource_server/emails/verification/code"
	"globe-and-citizen/layer8/server/resource_server/emails/verification/zk"
	"globe-and-citizen/layer8/server/resource_server/service"
	resourceUtils "globe-and-citizen/layer8/server/resource_server/utils"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/coder/websocket"
	"github.com/stretchr/testify/assert"

	"globe-and-citizen/layer8/server/resource_server/utils/mocks"

	"globe-and-citizen/layer8/server/resource_server/models"

	utils "github.com/globe-and-citizen/layer8-utils"
)

type mockWsResponseRecorder struct {
	http.ResponseWriter
	hijack func() (net.Conn, *bufio.ReadWriter, error)
}

func (mj mockWsResponseRecorder) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return mj.hijack()
}

// Asserting we implement the http.Hijacker interface
var _ http.Hijacker = &mockWsResponseRecorder{}

const name = "name"
const redirectUri = "redirect_uri"
const username = "username"
const password = "password"

func TestTunnel_WebSocketImpl(t *testing.T) {
	// create the ws Server connection
	wsMockClient := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := websocket.Accept(w, r, nil)
		if err != nil {
			t.Fatal(err)
		}
		defer c.CloseNow()

		// Set the context as needed. Use of r.Context() is not recommended
		// to avoid surprising behavior (see http.Hijacker).
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		_, msg, err := c.Read(ctx)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, "Hello, Server!", string(msg))
		c.Write(ctx, websocket.MessageText, []byte("Hello, Client!"))
	}))

	req, err := http.NewRequest("GET", wsMockClient.URL, nil)
	if err != nil {
		t.Fatal(err)
	}

	url_parts := strings.Split(wsMockClient.URL, "://")
	req.Header.Add("X-Forwarded-Proto", "ws")
	req.Header.Add("X-Forwarded-Host", url_parts[1])

	// this bit here is to make the request a websocket request
	server, _ := net.Pipe()

	rw := bufio.NewReadWriter(bufio.NewReader(server), bufio.NewWriter(server))
	responseRecorder := mockWsResponseRecorder{
		ResponseWriter: httptest.NewRecorder(),
		hijack: func() (net.Conn, *bufio.ReadWriter, error) {
			return server, rw, nil
		},
	}

	req.Header.Add("upgrade", "websocket")
	req.Header.Add("connection", "upgrade")
	req.Header.Add("sec-websocket-version", "13")
	req.Header.Add("sec-websocket-key", "dGhlIHNhbXBsZSBub25jZQ==")

	Tunnel(responseRecorder, req)
}

func prepareInitTunnelRequest(clientBackendUrl string, mockRepo *mocks.MockRepository) *http.Request {
	resourceService := service.NewService(
		mockRepo,
		&verification.EmailVerifier{},
		&zk.ProofProcessor{},
		code.NewMIMCCodeGenerator(),
	)
	resourceService.RegisterClient(dto.RegisterClientDTO{
		Name:        name,
		RedirectURI: redirectUri,
		BackendURI:  resourceUtils.RemoveProtocolFromURL(clientBackendUrl),
		Username:    username,
		StoredKey:   "storedKey",
		ServerKey:   "serverKey",
	})

	reqToInitTunnel := httptest.NewRequest("GET", "/init-tunnel", nil)
	reqToInitTunnel = reqToInitTunnel.WithContext(
		context.WithValue(
			reqToInitTunnel.Context(),
			"service",
			resourceService,
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

	mockRepo := &mocks.MockRepository{
		RegisterClientMock: func(client models.Client) error {
			assert.Equal(t, name, client.Name)
			assert.Equal(t, username, client.Username)
			assert.Equal(t, redirectUri, client.RedirectURI)
			assert.Equal(t, resourceUtils.RemoveProtocolFromURL(mockedServiceProvider.URL), client.BackendURI)

			return fmt.Errorf("failed to store a client")
		},
	}

	reqToInitTunnel := prepareInitTunnelRequest(mockedServiceProvider.URL, mockRepo)
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
	mockRepo := &mocks.MockRepository{
		RegisterClientMock: func(client models.Client) error {
			return fmt.Errorf("failed to store a client")
		},
	}
	reqToInitTunnel := prepareInitTunnelRequest("http://localhost:8080", mockRepo)
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
