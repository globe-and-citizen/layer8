package constants

// Scopes
const (
	ReadUserScope                = "read:user"
	ReadUserDisplayNameScope     = "read:user:display_name"
	ReadUserColorScope           = "read:user:color"
	ReadUserIsEmailVerifiedScope = "read:user:is_email_verified"
)

const (
	UserDisplayNameMetadataKey   = "display_name"
	UserEmailVerifiedMetadataKey = "email_verified"
	UserColorMetadataKey         = "color"
)

// Scope descriptions
var ScopeDescriptions = map[string]string{
	ReadUserScope: "read anonymized information about your account",
}
