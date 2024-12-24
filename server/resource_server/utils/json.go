package utils

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"io"
	"net/http"
)

func DecodeJsonFromRequest[T any](w http.ResponseWriter, byteBuffer io.ReadCloser) (T, error) {
	var to T

	body, e := io.ReadAll(byteBuffer)
	if e != nil {
		HandleError(w, http.StatusBadRequest, "Error while reading request body", e)
		return to, e
	}

	e = json.Unmarshal(body, &to)
	if e != nil {
		HandleError(w, http.StatusBadRequest, "Request malformed: error while parsing json", e)
		return to, e
	}

	e = validator.New().Struct(to)
	if e != nil {
		HandleError(w, http.StatusBadRequest, "Input json is invalid", e)
		return to, e
	}

	return to, nil
}
