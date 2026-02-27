package codegen

// Security scheme type constants
const (
	SecurityTypeAPIKey   = "apiKey"
	SecurityTypeHTTP     = "http"
	SecurityTypeOAuth2   = "oauth2"
	SecurityTypeOpenID   = "openIdConnect"
	SecuritySchemeBearer = "bearer"
	SecuritySchemeBasic  = "basic"
)

// IsAPIKey returns true if the security scheme is an API Key type.
func IsAPIKey(schemeType string) bool {
	return schemeType == SecurityTypeAPIKey
}

// IsHTTP returns true if the security scheme is an HTTP type (Bearer, Basic).
func IsHTTP(schemeType string) bool {
	return schemeType == SecurityTypeHTTP
}

// IsOAuth2 returns true if the security scheme is OAuth2 type.
func IsOAuth2(schemeType string) bool {
	return schemeType == SecurityTypeOAuth2
}

// IsBearer returns true if the HTTP security scheme is Bearer token.
func IsBearer(scheme string) bool {
	return scheme == SecuritySchemeBearer
}

// IsBasicAuth returns true if the HTTP security scheme is Basic authentication.
func IsBasicAuth(scheme string) bool {
	return scheme == SecuritySchemeBasic
}
