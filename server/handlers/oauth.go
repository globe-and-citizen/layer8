package handlers

import (
	"encoding/json"
	"fmt"
	"globe-and-citizen/layer8/server/constants"
	"globe-and-citizen/layer8/server/entities"
	svc "globe-and-citizen/layer8/server/internals/service"
	"globe-and-citizen/layer8/server/resource_server/utils"
	"net/http"
	"strings"
)

func TokenHandler(w http.ResponseWriter, r *http.Request) {
	if !validateHttpMethod(w, r.Method, http.MethodPost) {
		return
	}

	req, err := utils.DecodeJsonFromRequest[entities.OauthTokenRequest](w, r.Body)
	if err != nil {
		return
	}

	service := r.Context().Value("Oauthservice").(svc.ServiceInterface)

	err = service.AuthenticateClient(req.ClientUUID, req.ClientSecret)
	if err != nil {
		utils.HandleError(w, http.StatusUnauthorized, "failed to authenticate client", err)
		return
	}

	err = service.VerifyAuthorizationCode(req.AuthorizationCode)
	if err != nil {
		utils.HandleError(w, http.StatusBadRequest, "the authorization code is invalid", err)
		return
	}

	accessToken, err := service.GenerateAccessToken(req.ClientUUID, req.ClientSecret)
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "internal error when generating the access token", err)
		return
	}

	resp := utils.BuildResponse(
		w,
		http.StatusOK,
		"access token generated successfully",
		entities.OauthTokenResponse{
			AccessToken:      accessToken,
			TokenType:        constants.TokenTypeBearer,
			ExpiresInMinutes: constants.AccessTokenValidityMinutes,
		},
	)

	if err := json.NewEncoder(w).Encode(&resp); err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "failed to encode the response", err)
	}
}

func ZkMetadataHandler(w http.ResponseWriter, r *http.Request) {
	if !validateHttpMethod(w, r.Method, http.MethodPost) {
		return
	}

	req, err := utils.DecodeJsonFromRequest[entities.ZkMetadataRequest](w, r.Body)
	if err != nil {
		return
	}

	authHeader := r.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, constants.TokenTypeBearer) {
		errorMsg := "invalid authorization header"
		utils.HandleError(w, http.StatusUnauthorized, errorMsg, fmt.Errorf(errorMsg))
		return
	}

	accessToken := authHeader[len(constants.TokenTypeBearer)+1:]

	service := r.Context().Value("Oauthservice").(svc.ServiceInterface)

	err = service.AuthenticateClient(req.ClientUUID, req.ClientSecret)
	if err != nil {
		utils.HandleError(w, http.StatusUnauthorized, "Failed to authenticate client", err)
		return
	}

	claims, err := service.ValidateAccessToken(req.ClientSecret, accessToken)
	if err != nil {
		utils.HandleError(w, http.StatusUnauthorized, "Failed to validate client access token", err)
		return
	}

	zkMetadata, err := service.GetZkUserMetadata(claims.Scopes, claims.UserID)
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Failed to get user metadata", err)
		return
	}

	resp := utils.BuildResponse(w, http.StatusOK, "User metadata retrieved successfully", zkMetadata)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Failed to encode the response", err)
	}
}

func validateHttpMethod(w http.ResponseWriter, actual string, expected string) bool {
	if actual != expected {
		errorMsg := fmt.Sprintf("invalid http method, expected: %s, got: %s", expected, actual)
		utils.HandleError(w, http.StatusMethodNotAllowed, errorMsg, fmt.Errorf(errorMsg))
		return false
	}

	return true
}
