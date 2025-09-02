package dto

type RegisterUserDTO struct {
	Username  string `json:"username" validate:"required,min=3,max=50"`
	PublicKey []byte `json:"public_key" validate:"required"`
	StoredKey string `json:"stored_key" validate:"required"`
	ServerKey string `json:"server_key" validate:"required"`
}

type LoginUserDTO struct {
	Username    string `json:"username" validate:"required"`
	Nonce       string `json:"nonce" validate:"required"`
	CNonce      string `json:"c_nonce" validate:"required"`
	ClientProof string `json:"client_proof" validate:"required"`
}

type LoginClientDTO struct {
	Username    string `json:"username" validate:"required"`
	Nonce       string `json:"nonce" validate:"required"`
	CNonce      string `json:"c_nonce" validate:"required"`
	ClientProof string `json:"client_proof" validate:"required"`
}

type LoginPrecheckDTO struct {
	Username string `json:"username" validate:"required"`
	CNonce   string `json:"c_nonce" validate:"required"`
}

type UpdateUserMetadataDTO struct {
	DisplayName string `json:"display_name" validate:"required"`
	Color       string `json:"color" validate:"required"`
	Bio         string `json:"bio" validate:"required"`
}

type RegisterClientPrecheckDTO struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
}

type RegisterClientDTO struct {
	Name        string `json:"name" validate:"required"`
	RedirectURI string `json:"redirect_uri" validate:"required"`
	BackendURI  string `json:"backend_uri" validate:"required"`
	Username    string `json:"username" validate:"required,min=3,max=50"`
	StoredKey   string `json:"stored_key" validate:"required"`
	ServerKey   string `json:"server_key" validate:"required"`
}

type CheckBackendURIDTO struct {
	BackendURI string `json:"backend_uri" validate:"required"`
}

type VerifyEmailDTO struct {
	Email string `json:"email" validate:"required,email"`
}

type CheckEmailVerificationCodeDTO struct {
	Email string `json:"email" validate:"required,email"`
	Code  string `json:"code" validate:"required"`
}

type RegisterUserPrecheckDTO struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
}

type ResetPasswordPrecheckDTO struct {
	Username string `json:"username" validate:"required"`
}

type ResetPasswordDTO struct {
	Username  string `json:"username" validate:"required,min=3,max=50"`
	Signature []byte `json:"signature" validate:"required"`
	StoredKey string `json:"stored_key" validation:"required,min=1"`
	ServerKey string `json:"server_key" validation:"required,min=1"`
}

type ClientUnpaidAmountDTO struct {
	ClientId string `json:"client_id" validate:"required"`
}
