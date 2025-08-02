package entities

import "github.com/dgrijalva/jwt-go"

type Client struct {
	ID          string `json:"id"`
	Secret      string `json:"secret"`
	Name        string `json:"name"`
	RedirectURI string `json:"redirect_uri"`
}

type ClientClaims struct {
	jwt.StandardClaims
	UserID int64
	Scopes string
}
