package handlers

import (
	"errors"
	svc "globe-and-citizen/layer8/server/internals/service"
	"globe-and-citizen/layer8/server/resource_server/utils"
	"html/template"
	"net/http"
	"os"
)

type AuthenticationHandler interface {
	Login(w http.ResponseWriter, r *http.Request)
	Register(w http.ResponseWriter, r *http.Request)
}

type authenticationHandlerImpl struct {
	service   svc.ServiceInterface
	parseHTML func(w http.ResponseWriter, htmlFile string, params map[string]interface{})
}

func NewAuthenticationHandler(
	service svc.ServiceInterface,
	htmlParserFunc func(w http.ResponseWriter, htmlFile string, params map[string]interface{}),
) AuthenticationHandler {
	return &authenticationHandlerImpl{
		service:   service,
		parseHTML: htmlParserFunc,
	}
}

func (a *authenticationHandlerImpl) Login(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		a.getLoginHandler(w, r)
	case "POST":
		a.postLoginHandler(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (a *authenticationHandlerImpl) getLoginHandler(w http.ResponseWriter, r *http.Request) {
	next := r.URL.Query().Get("next")
	if next == "" {
		next = "/"
	}
	// check if the user is already logged in
	token, err := r.Cookie("token")
	if token != nil && err == nil {
		user, err := a.service.GetUserByToken(token.Value)
		if err == nil && user != nil {
			http.Redirect(w, r, next, http.StatusSeeOther)
			return
		}
	}

	a.parseHTML(w, "assets-v1/templates/src/pages/oauth_portal/login.html",
		map[string]interface{}{
			"HasNext":  next != "",
			"Next":     next,
			"ProxyURL": os.Getenv("PROXY_URL"),
		},
	)
}

func (a *authenticationHandlerImpl) postLoginHandler(w http.ResponseWriter, r *http.Request) {
	next := r.URL.Query().Get("next")
	username := r.FormValue("username")
	password := r.FormValue("password")

	rUser, err := a.service.LoginUser(username, password)
	if err != nil {
		a.parseLoginWithErr(w, r, err)
		return
	}

	token, ok := rUser["token"].(string)
	if !ok {
		a.parseLoginWithErr(w, r, errors.New("could not get token"))

		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "token",
		Value: token,
		Path:  "/",
	})

	http.Redirect(w, r, next, http.StatusSeeOther)
}

func (a *authenticationHandlerImpl) parseLoginWithErr(w http.ResponseWriter, r *http.Request, err error) {
	t, errT := template.ParseFiles("assets-v1/templates/src/pages/oauth_portal/login.html")
	if errT != nil {
		http.Error(w, errT.Error(), http.StatusInternalServerError)
		return
	}

	t.Execute(w, map[string]interface{}{
		"HasNext": true,
		"Next":    r.URL.Query().Get("next"),
		"Error":   err.Error(),
	})
}

func (a *authenticationHandlerImpl) Register(w http.ResponseWriter, r *http.Request) {
	service := r.Context().Value("OauthService").(*svc.Service)

	switch r.Method {
	case "GET":
		next := r.URL.Query().Get("next")
		if next == "" {
			next = "/"
		}
		// check if the user is already logged in
		token, err := r.Cookie("token")
		if token != nil && err == nil {
			user, err := service.GetUserByToken(token.Value)
			if err == nil && user != nil {
				http.Redirect(w, r, next, http.StatusSeeOther)
				return
			}
		}

		utils.ParseHTML(w, "registerClient.html",
			map[string]interface{}{
				"HasNext":  next != "",
				"Next":     next,
				"ProxyURL": os.Getenv("PROXY_URL"),
			},
		)
		return
	case "POST":
		next := r.URL.Query().Get("next")
		username := r.FormValue("username")
		password := r.FormValue("password")
		// login the user
		rUser, err := service.LoginUser(username, password)
		if err != nil {
			utils.ParseHTML(w, "registerClient.html",
				map[string]interface{}{
					"HasNext":  next != "",
					"Next":     next,
					"Error":    err.Error(),
					"ProxyURL": os.Getenv("PROXY_URL"),
				},
			)
			return
		}

		// set the token cookie
		token, ok := rUser["token"].(string)
		if !ok {
			utils.ParseHTML(w, "registerClient.html",
				map[string]interface{}{
					"HasNext":  next != "",
					"Next":     next,
					"Error":    "could not get token",
					"ProxyURL": os.Getenv("PROXY_URL"),
				},
			)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:  "token",
			Value: token,
			Path:  "/",
		})
		// redirect to next page - here the user already knows their pseudo profile
		// when they registered
		http.Redirect(w, r, next, http.StatusSeeOther)
		return
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}
