package utils

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"

	"github.com/dgrijalva/jwt-go"
	"github.com/xdg-go/pbkdf2"
)

type CustomClaims struct {
	UserAgent       string `json:"user_agent"`
	SecChUaPlatform string `json:"sec-ch-ua-platform"`
	SecChUaMobile   string `json:"sec-ch-ua-mobile"`
	SecFetchSite    string `json:"sec-fetch-site"`
	SecFetchDest    string `json:"sec-fetch-dest"`
	SecFetchMode    string `json:"sec-fetch-mode"`
	Referer         string `json:"referer"`
	Origin          string `json:"origin"`
	ContentLength   string `json:"content-length"`
	jwt.StandardClaims
}

func GenerateToken(secret string, claims *CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("could not generate user token: %s", err)
	}
	return tokenString, nil
}

func SaltAndHashPassword(password string, salt string) string {
	dk := pbkdf2.Key([]byte(password), []byte(salt), 4096, 32, sha1.New)
	return hex.EncodeToString(dk[:])
}
