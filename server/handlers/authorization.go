package handlers

import (
	"globe-and-citizen/layer8/server/constants"
	svc "globe-and-citizen/layer8/server/internals/service"
	"globe-and-citizen/layer8/server/utils"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"golang.org/x/oauth2"
)

type AuthorizationHandler interface {
	Authorize(w http.ResponseWriter, r *http.Request)
	OAuthToken(w http.ResponseWriter, r *http.Request)
	Error(w http.ResponseWriter, r *http.Request)
}

type authorizationHandlerImpl struct {
	service   svc.ServiceInterface
	parseHTML func(w http.ResponseWriter, statusCode int, htmlFile string, params map[string]interface{})
}

func NewAuthorizationHandler(
	service svc.ServiceInterface,
	htmlParserFunc func(w http.ResponseWriter, statusCode int, htmlFile string, params map[string]interface{}),
) AuthorizationHandler {
	return &authorizationHandlerImpl{
		service:   service,
		parseHTML: htmlParserFunc,
	}
}

func (a *authorizationHandlerImpl) Authorize(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		a.getAuthorizeHandler(w, r)
	case "POST":
		a.postAuthorizeHandler(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (a *authorizationHandlerImpl) getAuthorizeHandler(w http.ResponseWriter, r *http.Request) {
	var (
		clientID          = r.URL.Query().Get("client_id")
		scopes            = r.URL.Query().Get("scope")
		redirectURI       = r.URL.Query().Get("redirect_uri")
		scopeDescriptions = []string{}
		next              string
	)

	// use the default scope if none is provided
	if scopes == "" {
		scopes = constants.READ_USER_SCOPE
	}

	// add the scope descriptions
	for _, scope := range strings.Split(scopes, ",") {
		scopeDescriptions = append(scopeDescriptions, constants.ScopeDescriptions[scope])
	}

	// get the client
	client, err := a.service.GetClient(clientID)
	if err != nil {
		log.Println(err)
		// redirect to the redirect_uri with error
		http.Redirect(w, r, "/error?opt=invalid_client", http.StatusSeeOther)
		return
	}

	// generate the next url
	uri, err := url.Parse("/authorize")
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/error?opt=server_error", http.StatusSeeOther)
		return
	}

	uri.RawQuery = url.Values{
		"client_id": {clientID},
		"scope":     {scopes},
	}.Encode()

	next = uri.String()

	// check that the user is logged in
	token, err := r.Cookie("token")
	if token == nil || err != nil {
		http.Redirect(w, r, "/login?next="+next, http.StatusSeeOther)
		return
	}

	user, err := a.service.GetUserByToken(token.Value)
	if err != nil || user == nil {
		http.Redirect(w, r, "/login?next="+next, http.StatusSeeOther)
		return
	}

	// check that the redirect_uri is valid match the client's redirect_uri
	if redirectURI != "" && client.RedirectURI != redirectURI {
		http.Redirect(w, r, "/error?opt=redirect_uri_mismatch", http.StatusSeeOther)
		return
	}

	a.parseHTML(w, http.StatusOK, "assets-v1/templates/src/pages/oauth_portal/authorize.html", map[string]interface{}{
		"ClientName": client.Name,
		"Scopes":     scopeDescriptions,
		"Next":       next,
	})
}

func (a *authorizationHandlerImpl) postAuthorizeHandler(w http.ResponseWriter, r *http.Request) {
	var (
		clientID        = r.URL.Query().Get("client_id")
		scopes          = r.URL.Query().Get("scope")
		returnResult, _ = strconv.ParseBool(r.URL.Query().Get("return_result"))
	)

	// decision := r.FormValue("decision")
	// if decision != "allow" {
	// 	log.Println("User denied access")
	// 	utils.MapResponse(
	// 		returnResult, w,
	// 		&utils.JSONResponseInput{
	// 			StatusCode: http.StatusOK,
	// 			Data:       `{"redr": "/error?opt=access_denied"}`,
	// 		},
	// 		&utils.RedirectResponseInput{
	// 			StatusCode: http.StatusSeeOther,
	// 			Location:   "/error?opt=access_denied",
	// 		},
	// 	)

	// 	return
	// }

	// get the client
	client, err := a.service.GetClient(clientID)
	if err != nil {
		log.Println(err)
		utils.MapResponse(
			returnResult, w,
			&utils.JSONResponseInput{
				StatusCode: http.StatusOK,
				Data:       `{"redr": "/error?opt=invalid_client"}`,
			},
			&utils.RedirectResponseInput{
				StatusCode: http.StatusSeeOther,
				Location:   "/error?opt=invalid_client",
			},
		)
		return
	}

	token, err := r.Cookie("token")
	if err != nil || token == nil {
		utils.MapResponse(
			returnResult, w,
			&utils.JSONResponseInput{
				StatusCode: http.StatusOK,
				Data:       `{"redr": "/login?next=` + r.URL.String() + `"}`,
			},
			&utils.RedirectResponseInput{
				StatusCode: http.StatusSeeOther,
				Location:   "/login?next=" + r.URL.String(),
			},
		)
	}

	user, err := a.service.GetUserByToken(token.Value)
	if err != nil || user == nil {
		utils.MapResponse(
			returnResult, w,
			&utils.JSONResponseInput{
				StatusCode: http.StatusOK,
				Data:       `{"redr": "/login?next=` + r.URL.String() + `"}`,
			},
			&utils.RedirectResponseInput{
				StatusCode: http.StatusSeeOther,
				Location:   "/login?next=" + r.URL.String(),
			},
		)

		return
	}

	if scopes == "" {
		scopes = constants.READ_USER_SCOPE
	}

	if r.FormValue("share_details_with_SPA") == "true" {
		scopes += "," + constants.SHARE_DETAILS_WITH_SPA
	}

	redirectURL, err := a.service.GenerateAuthorizationURL(&oauth2.Config{
		ClientID:    client.ID,
		RedirectURL: client.RedirectURI,
		Scopes:      strings.Split(scopes, ","),
	}, int64(user.ID))
	if err != nil {
		utils.MapResponse(
			returnResult, w,
			&utils.JSONResponseInput{
				StatusCode: http.StatusOK,
				Data:       `{"redr": "/error?opt=server_error"}`,
			},
			&utils.RedirectResponseInput{
				StatusCode: http.StatusSeeOther,
				Location:   "/error?opt=server_error",
			},
		)

		return
	}

	code, err := a.service.GenerateAuthJwtCode(&oauth2.Config{
		ClientID:    client.ID,
		RedirectURL: client.RedirectURI,
		Scopes:      strings.Split(scopes, ","),
	}, int64(user.ID))
	if err != nil {
		utils.MapResponse(
			returnResult, w,
			&utils.JSONResponseInput{
				StatusCode: http.StatusOK,
				Data:       `{"redr": "/error?opt=server_error"}`,
			},
			&utils.RedirectResponseInput{
				StatusCode: http.StatusSeeOther,
				Location:   "/error?opt=server_error",
			},
		)
		return
	}

	utils.MapResponse(
		returnResult, w,
		&utils.JSONResponseInput{
			StatusCode: http.StatusOK,
			// Data:       `{"redr": "` + redirectURL.String() + `"}`,
			Data: `{"redr": "` + redirectURL.String() + `", "code": "` + code + `"}`,

		},
		&utils.RedirectResponseInput{
			StatusCode: http.StatusSeeOther,
			Location:   redirectURL.String(),
			// Cookie:     &http.Cookie{Name: "code", Value: code, Path: "/", MaxAge: 60 * 10, HttpOnly: true},
		},
	)
}

func (a *authorizationHandlerImpl) OAuthToken(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		a.postOAauthToken(w, r)
	default:
		utils.WriteJSONResponse(
			w,
			http.StatusMethodNotAllowed,
			`{"error": "method not allowed"}`,
		)

		return
	}
}

func (a *authorizationHandlerImpl) postOAauthToken(w http.ResponseWriter, r *http.Request) {
	var (
		code        = r.FormValue("code")
		redirectURI = r.FormValue("redirect_uri")
	)

	clientID, clientSecret, err := utils.GetClientIDAndSecretFromAuthHeader(r.Header.Get("Authorization"))
	if err != nil {
		utils.WriteJSONResponse(
			w,
			http.StatusUnauthorized,
			`{"error": "`+err.Error()+`"}`,
		)

		return
	}

	// get the client
	client, err := a.service.GetClient(clientID)
	if err != nil {
		utils.WriteJSONResponse(
			w,
			http.StatusUnauthorized,
			`{"error": "invalid client"}`,
		)
		return
	}

	if client.Secret != clientSecret {
		utils.WriteJSONResponse(
			w,
			http.StatusUnauthorized,
			`{"error": "invalid client secret"}`,
		)
		return
	}

	token, err := a.service.ExchangeCodeForToken(&oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURI,
	}, code)
	if err != nil {
		utils.WriteJSONResponse(
			w,
			http.StatusInternalServerError,
			`{"error": "`+err.Error()+`"}`,
		)

		return
	}

	utils.WriteJSONResponse(
		w,
		http.StatusOK,
		token,
	)
}

func (a *authorizationHandlerImpl) Error(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		a.getError(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func (a *authorizationHandlerImpl) getError(w http.ResponseWriter, r *http.Request) {
	var (
		opt    = r.URL.Query().Get("opt")
		opts   = []string{}
		errors = map[string]string{
			"invalid_client":        "The client is invalid.",
			"access_denied":         "The user denied the request.",
			"server_error":          "An error occurred on the server.",
			"redirect_uri_mismatch": "The redirect uri does not match the client's redirect uri.",
		}
	)
	// add the error to the list of errors
	for _, v := range strings.Split(opt, ",") {
		opts = append(opts, errors[v])
	}

	a.parseHTML(w, http.StatusBadRequest, "assets-v1/templates/src/pages/oauth_portal/error.html", map[string]interface{}{
		"Errors": opts,
	})
}
