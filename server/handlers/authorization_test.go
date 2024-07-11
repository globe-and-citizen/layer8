package handlers_test

import (
	"errors"
	"fmt"
	"globe-and-citizen/layer8/server/handlers"
	"globe-and-citizen/layer8/server/models"
	"globe-and-citizen/layer8/server/resource_server/utils"
	"globe-and-citizen/layer8/server/utils/mocks"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
)

func Test_GetAuthorize_OK(t *testing.T) {
	var (
		clientID    = "clientID"
		redirectUri = "redirectUri"
		client      = &models.Client{
			Name:        "clientName",
			ID:          clientID,
			RedirectURI: redirectUri,
		}

		user = &models.User{}

		tokenValue = "jwtToken"

		q = url.Values{
			"client_id":    []string{clientID},
			"redirect_uri": []string{redirectUri},
		}
	)

	ctrl := gomock.NewController(t)
	serviceMock := mocks.NewMockServiceInterface(ctrl)
	serviceMock.EXPECT().GetClient(clientID).Return(client, nil)
	serviceMock.EXPECT().GetUserByToken(tokenValue).Return(user, nil)

	parseHtml := func(w http.ResponseWriter, htmlFile string, params map[string]interface{}) {
		q, _ := url.QueryUnescape(params["Next"].(string))

		assert.Equal(t, "assets-v1/templates/src/pages/oauth_portal/authorize.html", htmlFile)
		assert.Equal(t, client.Name, params["ClientName"])
		assert.Equal(t, "/authorize?client_id=clientID&scope=read:user", q)
		assert.Equal(t, []string{"read anonymized information about your account"}, params["Scopes"])
	}

	authorizationHandler := handlers.NewAuthorizationHandler(serviceMock, parseHtml)

	req := httptest.NewRequest("GET", fmt.Sprintf("/authorize?%s", q.Encode()), nil)
	req.AddCookie(&http.Cookie{Name: "token", Value: tokenValue})
	responseRecorder := httptest.NewRecorder()
	authorizationHandler.Authorize(responseRecorder, req)

	// Verify the results
	assert.Equal(t, http.StatusOK, responseRecorder.Code)
}

func Test_GetAuthorize_NoClient(t *testing.T) {
	var (
		clientID    = "clientID"
		redirectUri = "redirectUri"
		q           = url.Values{
			"client_id":    []string{clientID},
			"redirect_uri": []string{redirectUri},
		}
	)

	ctrl := gomock.NewController(t)
	serviceMock := mocks.NewMockServiceInterface(ctrl)
	serviceMock.EXPECT().GetClient(clientID).Return(nil, errors.New("client not found"))
	authorizationHandler := handlers.NewAuthorizationHandler(serviceMock, utils.ParseHTML)

	req := httptest.NewRequest("GET", fmt.Sprintf("/authorize?%s", q.Encode()), nil)
	responseRecorder := httptest.NewRecorder()
	authorizationHandler.Authorize(responseRecorder, req)

	assert.Equal(t, http.StatusSeeOther, responseRecorder.Code)
	assert.Equal(t, "/error?opt=invalid_client", responseRecorder.Header().Get("Location"))
}

func Test_GetAuthorize_NoToken(t *testing.T) {
	var (
		clientID    = "clientID"
		redirectUri = "redirectUri"
		client      = &models.Client{
			Name:        "clientName",
			ID:          clientID,
			RedirectURI: redirectUri,
		}

		q = url.Values{
			"client_id":    []string{clientID},
			"redirect_uri": []string{redirectUri},
		}
	)

	ctrl := gomock.NewController(t)
	serviceMock := mocks.NewMockServiceInterface(ctrl)
	serviceMock.EXPECT().GetClient(clientID).Return(client, nil)

	authorizationHandler := handlers.NewAuthorizationHandler(serviceMock, utils.ParseHTML)

	req := httptest.NewRequest("GET", fmt.Sprintf("/authorize?%s", q.Encode()), nil)
	responseRecorder := httptest.NewRecorder()
	authorizationHandler.Authorize(responseRecorder, req)

	assert.Equal(t, http.StatusSeeOther, responseRecorder.Code)
	assert.Equal(t, "/login?next=/authorize?client_id=clientID&scope=read%3Auser", responseRecorder.Header().Get("Location"))
}

func Test_GetAuthorize_InvalidToken(t *testing.T) {
	var (
		clientID    = "clientID"
		redirectUri = "redirectUri"
		client      = &models.Client{
			Name:        "clientName",
			ID:          clientID,
			RedirectURI: redirectUri,
		}

		tokenValue = "jwtToken"

		q = url.Values{
			"client_id":    []string{clientID},
			"redirect_uri": []string{redirectUri},
		}
	)

	ctrl := gomock.NewController(t)
	serviceMock := mocks.NewMockServiceInterface(ctrl)
	serviceMock.EXPECT().GetClient(clientID).Return(client, nil)
	serviceMock.EXPECT().GetUserByToken(tokenValue).Return(nil, errors.New("invalid token"))

	authorizationHandler := handlers.NewAuthorizationHandler(serviceMock, utils.ParseHTML)
	req := httptest.NewRequest("GET", fmt.Sprintf("/authorize?%s", q.Encode()), nil)
	req.AddCookie(&http.Cookie{Name: "token", Value: tokenValue})
	responseRecorder := httptest.NewRecorder()
	authorizationHandler.Authorize(responseRecorder, req)

	// Verify the results
	assert.Equal(t, http.StatusSeeOther, responseRecorder.Code)
	assert.Equal(t, "/login?next=/authorize?client_id=clientID&scope=read%3Auser", responseRecorder.Header().Get("Location"))
}

func Test_GetAuthorize_InvalidRedirectURI(t *testing.T) {
	var (
		clientID    = "clientID"
		redirectUri = "redirectUri"
		client      = &models.Client{
			Name:        "clientName",
			ID:          clientID,
			RedirectURI: "differentRedirectUri",
		}

		user = &models.User{}

		tokenValue = "jwtToken"

		q = url.Values{
			"client_id":    []string{clientID},
			"redirect_uri": []string{redirectUri},
		}
	)

	ctrl := gomock.NewController(t)
	serviceMock := mocks.NewMockServiceInterface(ctrl)
	serviceMock.EXPECT().GetClient(clientID).Return(client, nil)
	serviceMock.EXPECT().GetUserByToken(tokenValue).Return(user, nil)

	authorizationHandler := handlers.NewAuthorizationHandler(serviceMock, utils.ParseHTML)

	req := httptest.NewRequest("GET", fmt.Sprintf("/authorize?%s", q.Encode()), nil)
	req.AddCookie(&http.Cookie{Name: "token", Value: tokenValue})
	responseRecorder := httptest.NewRecorder()
	authorizationHandler.Authorize(responseRecorder, req)

	// Verify the results
	assert.Equal(t, http.StatusSeeOther, responseRecorder.Code)
	assert.Equal(t, "/error?opt=redirect_uri_mismatch", responseRecorder.Header().Get("Location"))
}

func Test_PostAuthorize_OK(t *testing.T) {
	var (
		clientID             = "clientID"
		scopes               = "read:user"
		returnResult         = "true"
		decision             = "allow"
		shareDisplayName     = "true"
		shareCountry         = "true"
		shareTopFiveMetadata = "true"
		JWTToken             = "jwtToken"

		formData = url.Values{
			"decision":                []string{decision},
			"share_display_name":      []string{shareDisplayName},
			"share_country":           []string{shareCountry},
			"share_top_five_metadata": []string{shareTopFiveMetadata},
		}

		q = url.Values{
			"client_id":     []string{clientID},
			"scope":         []string{scopes},
			"return_result": []string{returnResult},
		}

		client = &models.Client{
			ID:          clientID,
			RedirectURI: "redirectUri",
		}

		user = &models.User{
			ID: 1,
		}
	)

	ctrl := gomock.NewController(t)
	serviceMock := mocks.NewMockServiceInterface(ctrl)
	serviceMock.EXPECT().GetClient(clientID).Return(client, nil)
	serviceMock.EXPECT().GetUserByToken(JWTToken).Return(user, nil)
	serviceMock.EXPECT().GenerateAuthorizationURL(&oauth2.Config{
		ClientID:    clientID,
		RedirectURL: client.RedirectURI,
		Scopes:      []string{"read:user,read:user:display_name,read:user:country,read:user:top_five_metadata"},
	}, user.ID, gomock.Any())
}
