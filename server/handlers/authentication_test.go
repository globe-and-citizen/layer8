package handlers

import (
	"globe-and-citizen/layer8/server/mocks"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

/* TESTS THAT I WILL NEED TO WRITE
* 1) Query param next == ""
* 2) Query param next == "<?>"
* 3) request.Cookie("token") != nil
* 4) etc...
 */

func Test_GetLogin_NoToken_OK(t *testing.T) {
	ctrl := gomock.NewController(t)

	var (
		proxyUrl                  = "http://localhost:5001"
		expectedLoginHTMLPath     = "assets-v1/templates/src/pages/oauth_portal/login.html"
		expectedHTMLParsingParams = map[string]interface{}{
			"HasNext":  true,
			"Next":     "/",
			"ProxyURL": proxyUrl,
		}
	)

	os.Setenv("PROXY_URL", proxyUrl)

	serviceMock := mocks.NewMockServiceInterface(ctrl)
	htmlParserMock := func(w http.ResponseWriter, htmlFile string, params map[string]interface{}) {
		assert.Equal(t, expectedLoginHTMLPath, htmlFile)
		assert.Equal(t, expectedHTMLParsingParams, params)
	}

	handler := NewAuthenticationHandler(serviceMock, htmlParserMock)

	req := httptest.NewRequest("GET", "/login", nil)
	responseRecorder := httptest.NewRecorder()

	handler.Login(responseRecorder, req)
	assert.Equal(t, http.StatusOK, responseRecorder.Code)
}
