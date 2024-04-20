package handlers

import (
	"encoding/json"
	svc "globe-and-citizen/layer8/server/internals/service"
	"log"
	"net/http"
	"strings"
)

// UserInfo handles requests to get a user's anonymized data
// The last step of the oauth flow
func UserInfo(w http.ResponseWriter, r *http.Request) {
	service := r.Context().Value("Oauthservice").(*svc.Service)

	// Print headers
	// fmt.Println("*** Headers ***")
	// for name, values := range r.Header {
	// 	for _, value := range values {
	// 		fmt.Println("Header: ", name, value)
	// 	}
	// }

	switch r.Method {
	case "GET":
		token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		data, err := service.AccessResourcesWithToken(token)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error": "invalid token"}`))
			return
		}
		b, err := json.Marshal(data)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error": "server error"}`))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(b)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(`{"error": "method not allowed"}`))
		return
	}
}
