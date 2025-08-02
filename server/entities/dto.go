package entities

type X509CertificateRequest struct {
	Certificate string `json:"certificate" validate:"required"`
}

type X509CertificateResponse struct {
	X509Certificate string `json:"x509_certificate"`
}

type OauthTokenRequest struct {
	ClientUUID        string `json:"client_oauth_uuid" validate:"required"`
	ClientSecret      string `json:"client_oauth_secret" validate:"required"`
	AuthorizationCode string `json:"authorization_code" validate:"required"`
}

type OauthTokenResponse struct {
	AccessToken      string `json:"access_token"`
	TokenType        string `json:"token_type"`
	ExpiresInMinutes int    `json:"expires_in_minutes"`
}

type ZkMetadataRequest struct {
	ClientUUID   string `json:"client_oauth_uuid" validate:"required"`
	ClientSecret string `json:"client_oauth_secret" validate:"required"`
}

type ZkMetadataResponse struct {
	Country         string `json:"country"`
	IsEmailVerified bool   `json:"is_email_verified"`
	DisplayName     string `json:"display_name"`
	Color           string `json:"color"`
}
