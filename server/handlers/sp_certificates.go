package handlers

import (
	"encoding/json"
	"fmt"
	"globe-and-citizen/layer8/server/entities"
	svc "globe-and-citizen/layer8/server/internals/service"
	"globe-and-citizen/layer8/server/resource_server/utils"
	"log"
	"net/http"
	"strings"
)

func UploadSPCertificate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errorMessage := fmt.Sprintf("Invalid http method, expected POST, actual: %s", r.Method)

		utils.HandleError(w, http.StatusMethodNotAllowed, errorMessage, fmt.Errorf(errorMessage))
		return
	}

	token := r.Header.Get("Authorization")
	if token == "" {
		utils.HandleError(
			w,
			http.StatusUnauthorized,
			"authorization token is empty",
			fmt.Errorf("token is invalid"),
		)
		return
	}
	token = token[7:]

	clientClaims, err := utils.ValidateClientToken(token)
	if err != nil {
		utils.HandleError(
			w,
			http.StatusUnauthorized,
			"failed to verify client auth token",
			err,
		)
		return
	}

	req, err := utils.DecodeJsonFromRequest[entities.X509CertificateRequest](w, r.Body)
	if err != nil {
		return
	}

	service := r.Context().Value("Oauthservice").(svc.ServiceInterface)

	err = service.SaveX509Certificate(clientClaims.ClientID, req.Certificate)
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "failed to save the SP x.509 certificate", err)
		return
	}

	response := utils.BuildResponseWithNoBody(
		w,
		http.StatusCreated,
		"x.509 certificate was saved successfully",
	)

	if err := json.NewEncoder(w).Encode(&response); err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "failed to encode the response", err)
	}
}

// GetSPPubKey handles requests to get the public key for the service provider
func GetSPPubKey(w http.ResponseWriter, r *http.Request) {
	service := r.Context().Value("Oauthservice").(svc.ServiceInterface)
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
		client, err := service.CheckClient(backendURL)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error": "internal server error", "registered": false}`))
			return
		}

		response := entities.X509CertificateResponse{
			X509Certificate: string(client.X509CertificateBytes),
		}

		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			w.Write([]byte(`{"error": "internal server error", "registered": false}`))
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(`{"error": "method not allowed"}`))
		return
	}
}
