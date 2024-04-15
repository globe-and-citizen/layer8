package handlers

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"globe-and-citizen/layer8/server/constants"
	svc "globe-and-citizen/layer8/server/internals/service"
	"globe-and-citizen/layer8/server/models"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"golang.org/x/oauth2"
)

func Authorize(w http.ResponseWriter, r *http.Request) {
	service := r.Context().Value("Oauthservice").(*svc.Service)

	switch r.Method {
	case "GET":
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
		client, err := service.GetClient(clientID)
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
		q := uri.Query()
		q.Set("client_id", clientID)
		q.Set("scope", scopes)
		uri.RawQuery = q.Encode()
		next = uri.String()

		// check that the user is logged in
		token, err := r.Cookie("token")
		if token != nil && err == nil {
			user, err := service.GetUserByToken(token.Value)
			if err != nil || user == nil {
				http.Redirect(w, r, "/login?next="+next, http.StatusSeeOther)
				return
			}
		} else {
			http.Redirect(w, r, "/login?next="+next, http.StatusSeeOther)

			return
		}

		// check that the redirect_uri is valid match the client's redirect_uri
		if redirectURI != "" && client.RedirectURI != redirectURI {
			http.Redirect(w, r, "/error?opt=redirect_uri_mismatch", http.StatusSeeOther)
			return
		}
		// load the authorize page
		t, err := template.ParseFiles("assets-v1/templates/src/pages/oauth_portal/authorize.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		t.Execute(w, map[string]interface{}{
			"ClientName": client.Name,
			"Scopes":     scopeDescriptions,
			"Next":       next,
		})
		return
	case "POST":
		var (
			clientID        = r.URL.Query().Get("client_id")
			scopes          = r.URL.Query().Get("scope")
			returnResult, _ = strconv.ParseBool(r.URL.Query().Get("return_result"))
		)
		// get authorization decision
		decision := r.FormValue("decision")
		if decision != "allow" {
			log.Println("User denied access")
			if returnResult {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"redr": "/error?opt=access_denied"}`))
				return
			}
			http.Redirect(w, r, "/error?opt=access_denied", http.StatusSeeOther)
			return
		}
		// use the default scope if none is provided
		if scopes == "" {
			scopes = constants.READ_USER_SCOPE
		}
		// get the client
		client, err := service.GetClient(clientID)
		if err != nil {
			log.Println(err)
			if returnResult {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"redr": "/error?opt=invalid_client"}`))
				return
			}
			http.Redirect(w, r, "/error?opt=invalid_client", http.StatusSeeOther)
			return
		}
		// get user
		var user *models.User
		token, err := r.Cookie("token") // Ravi you may need to renag....
		if token != nil && err == nil {
			user, err = service.GetUserByToken(token.Value)
			if err != nil || user == nil {
				if returnResult {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(`{"redr": "/login?next=` + r.URL.String() + `"}`))
					return
				}
				http.Redirect(w, r, "/login?next="+r.URL.String(), http.StatusSeeOther)
				return
			}
		} else {
			if returnResult {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"redr": "/login?next=` + r.URL.String() + `"}`))
				return
			}
			http.Redirect(w, r, "/login?next="+r.URL.String(), http.StatusSeeOther)
			return
		}
		// fmt.Println("scopes before user decision: ", scopes)
		if r.FormValue("share_display_name") == "true" {
			scopes += "," + constants.READ_USER_DISPLAY_NAME_SCOPE
		}
		if r.FormValue("share_country") == "true" {
			scopes += "," + constants.READ_USER_COUNTRY_SCOPE
		}

		// fmt.Println("scopes after user decision: ", scopes)
		// generate authorization url
		authURL, err := service.GenerateAuthorizationURL(&oauth2.Config{
			ClientID:    client.ID,
			RedirectURL: client.RedirectURI,
			Scopes:      strings.Split(scopes, ","),
		}, int64(user.ID))
		if err != nil {
			log.Println("Server error: ", err)
			if returnResult {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"redr": "/error?opt=server_error"}`))
				return
			}
			http.Redirect(w, r, "/error?opt=server_error", http.StatusSeeOther)
			return
		}
		// redirect to the authorization url
		if returnResult {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"redr": "` + authURL.String() + `"}`))
			fmt.Println("Normal exit...")
			return
		}
		http.Redirect(w, r, authURL.String(), http.StatusSeeOther)
		return
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func OAuthToken(w http.ResponseWriter, r *http.Request) {
	service := r.Context().Value("Oauthservice").(*svc.Service)

	// exchange code for token
	switch r.Method {
	case "POST":
		var (
			code        = r.FormValue("code")
			redirectURI = r.FormValue("redirect_uri")
		)

		// decode the basic auth header
		fromBasicAuth := func(t string) (string, string, error) {
			t = strings.TrimPrefix(t, "Basic ")
			b, err := base64.StdEncoding.DecodeString(t)
			if err != nil {
				return "", "", err
			}
			// first remove the "Basic " prefix
			s := strings.SplitN(string(b), ":", 2)
			if len(s) != 2 {
				return "", "", errors.New("invalid basic auth header")
			}
			return s[0], s[1], nil
		}
		clientID, clientSecret, err := fromBasicAuth(r.Header.Get("Authorization"))
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error": "` + err.Error() + `"}`))
			return
		}

		// get the client
		client, err := service.GetClient(clientID)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error": "invalid client"}`))
			return
		}
		// check that the client secret is correct
		if client.Secret != clientSecret {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error": "invalid client secret"}`))
			return
		}
		// exchange code for token
		token, err := service.ExchangeCodeForToken(&oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURI,
		}, code)
		if err != nil {
			res := map[string]string{"error": err.Error()}
			resJSON, _ := json.Marshal(res)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(resJSON)
			return
		}
		// return token
		bToken, err := json.Marshal(token)
		if err != nil {
			res := map[string]string{"error": err.Error()}
			resJSON, _ := json.Marshal(res)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(resJSON)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(bToken)
		return
	default:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(`{"error": "method not allowed"}`))
		return
	}
}

func Error(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
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
		// load the error page
		t, err := template.ParseFiles("assets-v1/templates/src/pages/oauth_portal/error.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		t.Execute(w, map[string]interface{}{
			"Errors": opts,
		})
		return
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}
