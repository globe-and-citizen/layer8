package handlers

import (
	"bytes"
	"errors"
	"globe-and-citizen/layer8/server/models"
	"globe-and-citizen/layer8/server/utils/mocks"
	"net/http"
	"net/http/httptest"
	"net/url"
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

func Test_GetLoginHandler_NoToken_OK(t *testing.T) {
	// Prepare the test
	var (
		proxyUrl                  = "http://localhost:5001"
		expectedLoginHTMLPath     = "assets-v1/templates/src/pages/oauth_portal/login.html"
		expectedHTMLParsingParams = map[string]interface{}{
			"HasNext":  true,
			"Next":     "/",
			"ProxyURL": proxyUrl,
		}
	)

	ctrl := gomock.NewController(t)

	serviceMock := mocks.NewMockServiceInterface(ctrl)
	htmlParserMock := func(w http.ResponseWriter, htmlFile string, params map[string]interface{}) {
		assert.Equal(t, expectedLoginHTMLPath, htmlFile)
		assert.Equal(t, expectedHTMLParsingParams, params)
	}

	handler := NewAuthenticationHandler(serviceMock, htmlParserMock)

	os.Setenv("PROXY_URL", proxyUrl)

	// Execute the test
	req := httptest.NewRequest("GET", "/login", nil)
	responseRecorder := httptest.NewRecorder()
	handler.Login(responseRecorder, req)

	// Verify the results
	assert.Equal(t, http.StatusOK, responseRecorder.Code)
}

func Test_GetLoginHandler_TokenExists_OK(t *testing.T) {

	// Prepare the test
	var (
		fakeJwtToken = "fakeJwtToken"
	)

	ctrl := gomock.NewController(t)

	serviceMock := mocks.NewMockServiceInterface(ctrl)
	serviceMock.EXPECT().GetUserByToken(fakeJwtToken).Return(&models.User{}, nil)

	handler := NewAuthenticationHandler(serviceMock, nil)

	// Execute the test
	req := httptest.NewRequest("GET", "/login", nil)
	req.AddCookie(&http.Cookie{Name: "token", Value: fakeJwtToken})
	responseRecorder := httptest.NewRecorder()
	handler.Login(responseRecorder, req)

	// Verify the results
	assert.Equal(t, http.StatusSeeOther, responseRecorder.Code)
	assert.Equal(t, "/", responseRecorder.Header().Get("Location"))
}

func Test_PostLoginHandler_ValidCredentials_OK(t *testing.T) {
	// Prepare the test
	var (
		username     = "username"
		password     = "password"
		nextUrl      = "/next"
		fakeJwtToken = "fakeJwt"

		loginResult = map[string]interface{}{
			"token":    fakeJwtToken,
			"username": username,
		}
	)

	ctrl := gomock.NewController(t)

	serviceMock := mocks.NewMockServiceInterface(ctrl)
	serviceMock.EXPECT().LoginUser(username, password).Return(loginResult, nil)

	handler := NewAuthenticationHandler(serviceMock, nil)

	// Execute the test
	params := url.Values{}
	params.Add("username", username)
	params.Add("password", password)

	req := httptest.NewRequest("POST", "/login", bytes.NewBufferString(params.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.URL.RawQuery = "next=" + nextUrl

	responseRecorder := httptest.NewRecorder()
	handler.Login(responseRecorder, req)

	// Verify the results
	assert.Equal(t, http.StatusSeeOther, responseRecorder.Code)
	assert.Equal(t, nextUrl, responseRecorder.Header().Get("Location"))

	var isCookieStored bool
	for _, cookie := range responseRecorder.Result().Cookies() {
		if cookie.Name == "token" {
			isCookieStored = true
			assert.Equal(t, fakeJwtToken, cookie.Value)
		}
	}

	assert.True(t, isCookieStored)
}

func Test_PostLoginHandler_TokenNotExists_OK(t *testing.T) {
	// Prepare the test
	var (
		username = "username"
		password = "password"
		nextUrl  = "/next"

		loginResult = map[string]interface{}{
			"username": username,
		}

		expectedLoginHTMLPath     = "assets-v1/templates/src/pages/oauth_portal/login.html"
		expectedHTMLParsingParams = map[string]interface{}{
			"HasNext": true,
			"Next":    nextUrl,
			"Error":   "could not get token",
		}
	)

	ctrl := gomock.NewController(t)

	htmlParserMock := func(w http.ResponseWriter, htmlFile string, params map[string]interface{}) {
		assert.Equal(t, expectedLoginHTMLPath, htmlFile)
		assert.Equal(t, expectedHTMLParsingParams, params)
	}
	serviceMock := mocks.NewMockServiceInterface(ctrl)
	serviceMock.EXPECT().LoginUser(username, password).Return(loginResult, nil)

	handler := NewAuthenticationHandler(serviceMock, htmlParserMock)

	// Execute the test
	params := url.Values{}
	params.Add("username", username)
	params.Add("password", password)

	req := httptest.NewRequest("POST", "/login", bytes.NewBufferString(params.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.URL.RawQuery = "next=" + nextUrl

	responseRecorder := httptest.NewRecorder()
	handler.Login(responseRecorder, req)

	// Verify the results
	assert.Equal(t, http.StatusOK, responseRecorder.Code)

}

func Test_PostLoginHandler_InvalidCredentials_OK(t *testing.T) {
	// Prepare the test
	var (
		username = "username"
		password = "password"
		nextUrl  = "/next"

		loginError = errors.New("invalid credentials")

		expectedLoginHTMLPath     = "assets-v1/templates/src/pages/oauth_portal/login.html"
		expectedHTMLParsingParams = map[string]interface{}{
			"HasNext": true,
			"Next":    nextUrl,
			"Error":   loginError.Error(),
		}
	)

	ctrl := gomock.NewController(t)

	htmlParserMock := func(w http.ResponseWriter, htmlFile string, params map[string]interface{}) {
		assert.Equal(t, expectedLoginHTMLPath, htmlFile)
		assert.Equal(t, expectedHTMLParsingParams, params)
	}
	serviceMock := mocks.NewMockServiceInterface(ctrl)
	serviceMock.EXPECT().LoginUser(username, password).Return(nil, loginError)

	handler := NewAuthenticationHandler(serviceMock, htmlParserMock)

	// Execute the test
	params := url.Values{}
	params.Add("username", username)
	params.Add("password", password)

	req := httptest.NewRequest("POST", "/login", bytes.NewBufferString(params.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.URL.RawQuery = "next=" + nextUrl

	responseRecorder := httptest.NewRecorder()
	handler.Login(responseRecorder, req)

	// Verify the results
	assert.Equal(t, http.StatusOK, responseRecorder.Code)
	assert.Len(t, responseRecorder.Result().Cookies(), 0)
}
