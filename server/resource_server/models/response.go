package models

type LoginPrecheckResponseOutput struct {
	Username string `json:"username"`
	Salt     string `json:"salt"`
}

type LoginUserResponseOutput struct {
	Token string `json:"token"`
}

type ProfileResponseOutput struct {
	Username      string `json:"username"`
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	DisplayName   string `json:"display_name"`
	Country       string `json:"country"`
	EmailVerified bool   `json:"email_verified"`
}
type ClientResponseOutput struct {
	ID          string `json:"id"`
	Secret      string `json:"secret"`
	Name        string `json:"name"`
	RedirectURI string `json:"redirect_uri"`
	BackendURI  string `json:"backend_uri"`
}

type RegisterUserPrecheckResponseOutput struct {
	Salt           string `json:"salt"`
	IterationCount int    `json:"iterationCount"`
}

type LoginPrecheckResponseOutputv2 struct {
	Salt      string `json:"salt"`
	IterCount int    `json:"iter_count"`
	Nonce     string `json:"nonce"`
}

type LoginUserResponseOutputv2 struct {
	ServerSignature string `json:"server_signature"`
}