package codegen

import "github.com/imerljak/beaverspec/pkg/core"

// SecuritySchemeData is a template-friendly representation of a security scheme.
// It is language-agnostic: any generator can use it to inspect scheme properties.
type SecuritySchemeData struct {
	Name        string
	Type        string
	In          string // for apiKey: "header", "query", "cookie"
	KeyName     string // for apiKey: actual parameter name (e.g. "X-API-Key")
	Scheme      string // for http: "bearer" or "basic"
	IsAPIKey    bool
	IsBearer    bool
	IsBasicAuth bool
	IsOAuth2    bool
}

// SecurityRequirementData contains all info needed to enforce one security requirement
// in generated code. It is language-agnostic.
type SecurityRequirementData struct {
	SchemeName  string   // scheme identifier, e.g. "bearerAuth"
	Scopes      []string // OAuth2 scopes (informational)
	IsBearer    bool
	IsBasicAuth bool
	IsAPIKey    bool
	IsOAuth2    bool
	APIKeyIn    string // "header", "query", "cookie"
	APIKeyName  string // parameter name for apiKey (e.g. "X-API-Key")
}

// ConvertSecuritySchemes converts core.SecurityScheme values into template-friendly
// SecuritySchemeData. The classification logic delegates to SecurityScheme methods.
func ConvertSecuritySchemes(schemes []core.SecurityScheme) []SecuritySchemeData {
	result := make([]SecuritySchemeData, 0, len(schemes))
	for _, s := range schemes {
		result = append(result, SecuritySchemeData{
			Name:        s.Name,
			Type:        s.Type,
			In:          s.In,
			KeyName:     s.ParameterName,
			Scheme:      s.Scheme,
			IsAPIKey:    s.IsAPIKey(),
			IsBearer:    s.IsBearer(),
			IsBasicAuth: s.IsBasicAuth(),
			IsOAuth2:    s.IsOAuth2(),
		})
	}
	return result
}

// SchemesByName builds a lookup map from scheme identifier → SecuritySchemeData.
func SchemesByName(schemes []SecuritySchemeData) map[string]SecuritySchemeData {
	m := make(map[string]SecuritySchemeData, len(schemes))
	for _, s := range schemes {
		m[s.Name] = s
	}
	return m
}
