package utils

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/xdg-go/pbkdf2"
)

func SaltAndHashPassword(password string, salt string) string {
	dk := pbkdf2.Key([]byte(password), []byte(salt), 4096, 32, sha1.New)
	return hex.EncodeToString(dk[:])
}

func WriteJSONResponse(
	w http.ResponseWriter,
	status int,
	data interface{},
) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	byteData, ok := data.(string)
	if ok {
		w.Write([]byte(byteData))
		return
	}

	resBody, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "Error marshalling JSON", http.StatusInternalServerError)
		return
	}

	w.Write(resBody)
}

type JSONResponseInput struct {
	StatusCode int
	Data       interface{}
}

type RedirectResponseInput struct {
	StatusCode int
	Location   string
}

func MapResponse(
	isJson bool,
	w http.ResponseWriter,
	j *JSONResponseInput,
	red *RedirectResponseInput,
) {
	if isJson {
		WriteJSONResponse(w, j.StatusCode, j.Data)
	} else {
		http.Redirect(w, &http.Request{}, red.Location, red.StatusCode)
	}
}

func GetClientIDAndSecretFromAuthHeader(t string) (string, string, error) {
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
