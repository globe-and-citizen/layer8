package constants

// Scopes
const (
	READ_USER_SCOPE              = "read:user"
	READ_USER_DISPLAY_NAME_SCOPE = "read:user:display_name"
	READ_USER_COUNTRY_SCOPE      = "read:user:country"
)

const (
	USER_DISPLAY_NAME_METADATA_KEY   = "display_name"
	USER_COUNTRY_METADATA_KEY        = "country"
	USER_EMAIL_VERIFIED_METADATA_KEY = "email_verified"
)

// Scope descriptions
var ScopeDescriptions = map[string]string{
	READ_USER_SCOPE: "read anonymized information about your account",
}
