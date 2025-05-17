package handlers

import (
	svc "globe-and-citizen/layer8/server/internals/service"
	"log"
	"net/http"
	"strings"
)

// GetSPPubKey handles requests to get the public key for the service provider
func GetSPPubKey(w http.ResponseWriter, r *http.Request) {
	service := r.Context().Value("Oauthservice").(*svc.Service)
	switch r.Method {
	case "GET":
		backendURL := r.URL.Query().Get("backend_url")
		if backendURL == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error": "missing backend_url parameter"}`))
			return
		}
		token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		if token == "" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error": "missing token"}`))
			return
		}
		// Uncomment the following lines to generate a token for testing purposes, to be removed in production
		// token, err := utilities.GenerateStandardToken("ThisIsASecret")
		// if err != nil {
		// 	log.Println(err)
		// 	w.WriteHeader(http.StatusInternalServerError)
		// 	w.Write([]byte(`{"error": "internal server error"}`))
		// 	return
		// }
		isValid, err := service.VerifyToken(token)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error": "invalid token"}`))
			return
		}
		if !isValid {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error": "invalid token"}`))
			return
		}
		_, err = service.CheckClient(backendURL)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error": "internal server error", "registered": false}`))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"registered": true}`))

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(`{"error": "method not allowed"}`))
		return
	}
}
