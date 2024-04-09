package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	utils "github.com/globe-and-citizen/layer8-utils"
)

func Test_InitTunnel_OK(t *testing.T) {
	var (
		MOCKED_MP_123_SECRET_KEY = "MOCKED_MP_123_SECRET_KEY"
		MOCKED_UP_999_SECRET_KEY = "MOCKED_UP_999_SECRET_KEY"
	)

	os.Setenv("MP_123_SECRET_KEY", MOCKED_MP_123_SECRET_KEY)
	os.Setenv("UP_999_SECRET_KEY", MOCKED_UP_999_SECRET_KEY)

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

	reqToInitTunnel := httptest.NewRequest("GET", "/api/tunnel", nil)

	queryParams := reqToInitTunnel.URL.Query()
	queryParams.Add("backend", mockedServiceProvider.URL)
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
	assert.NotEmpty(t, res["server_pubKeyECDHs"])

	upJWT, ok := res["up-JWT"].(string)
	if !ok {
		t.Fatal("up-JWT is not a string")
	}

	_, err = utils.VerifyStandardToken(upJWT, MOCKED_UP_999_SECRET_KEY)
	if err != nil {
		t.Fatal(err)
	}

	// TODO: validate server_pubKeyECDHs

}
