package utils

import (
	"fmt"
	"net/http"
)

func IsMethodValid(w http.ResponseWriter, r *http.Request, method string) bool {
	if r.Method != method {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintln(w, http.StatusText(http.StatusMethodNotAllowed))
		return false
	}

	return true
}
