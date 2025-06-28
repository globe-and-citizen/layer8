package models

type LoginPrecheckResponseOutput struct {
	Username string `json:"username"`
	Salt     string `json:"salt"`
}

type LoginPrecheckResponseOutputv2 struct {
	Salt      string `json:"salt"`
	IterCount int    `json:"iter_count"`
	Nonce     string `json:"nonce"`
}

type LoginUserResponseOutput struct {
	Token string `json:"token"`
}

type LoginUserResponseOutputv2 struct {
	ServerSignature string `json:"server_signature"`
	Token           string `json:"token"`
}

type LoginClientResponseOutputv2 struct {
	ServerSignature string `json:"server_signature"`
	Token           string `json:"token"`
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
	ID              string `json:"id"`
	Secret          string `json:"secret"`
	Name            string `json:"name"`
	RedirectURI     string `json:"redirect_uri"`
	BackendURI      string `json:"backend_uri"`
	X509Certificate string `json:"x509_certificate"`
}

type RegisterUserPrecheckResponseOutput struct {
	Salt           string `json:"salt"`
	IterationCount int    `json:"iterationCount"`
}

type RegisterClientPrecheckResponseOutput struct {
	Salt           string `json:"salt"`
	IterationCount int    `json:"iterationCount"`
}

type ResetPasswordPrecheckResponseOutput struct {
	Salt           string `json:"salt"`
	IterationCount int    `json:"iterationCount"`
}

type ClientUnpaidAmountResponseOutput struct {
	UnpaidAmount int `json:"unpaid_amount"`
}
