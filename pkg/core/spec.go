package core

import "github.com/getkin/kin-openapi/openapi3"

// Spec represents a processed OpenAPI specification
type Spec struct {
	Raw *openapi3.T

	Info        SpecInfo
	Models      []Model
	Endpoints   []Endpoint
	Tags        []Tag
	SecurityDef []SecurityScheme

	Version string // OpenAPI Version (e.g., "3.0.0")

	refResolver *RefResolver
}

// Tag represennts an API tag for grouping operations
type Tag struct {
	Name        string
	Description string
}

// SecuritySchema represents an authentication/authorization schema
type SecurityScheme struct {
	Type             string // "apiKey", "http", "oauth2", "openIdConnect"
	Description      string
	Name             string // For apiKey: header/query/cookie name
	In               string // For apiKey: "header", "query", "cookie"
	Scheme           string // For http: "bearer", "basic", etc.
	BearerFormat     string // For http bearer
	Flows            *OAuthFlows
	OpenIDConnectURL string
}

// OAuthFlows represents OAuth 2.0 flow configurations
type OAuthFlows struct {
	Implicit          *OAuthFlow
	Password          *OAuthFlow
	ClientCredentials *OAuthFlow
	AuthorizationCode *OAuthFlow
}

// OAuthFlow represetns a single OAuth 2.0 flow
type OAuthFlow struct {
	AuthorizationURL string
	TokenURL         string
	RefreshURL       string
	Scopes           map[string]string
}

// RefResolver handles $ref resolution in OpenAPI specs
type RefResolver struct {
	// TODO
}

type SpecInfo struct {
	Title       string
	Description string
	Version     string
	Contact     *Contact
	License     *License
}

type Contact struct {
	Name  string
	Email string
	URL   string
}

type License struct {
	Name string
	URL  string
}

type Model struct {
	Name        string
	Description string
	Type        string // "object", "array", "string", etc.
	Properties  []Property
	Required    []string
	Enum        []interface{}

	// For array type schemas
	IsArray bool      // true if this is an array schema
	Items   *Property // full items definition

	// For map type schemas (additionalProperties)
	IsMap           bool      // true if this is a map schema
	AdditionalProps *Property // the type of map values

	// Complex type handling
	AllOf []Model
	OneOf []Model
	AnyOf []Model

	// Metadata
	Discriminator *Discriminator
	IsDeprecated  bool
	Example       interface{}

	// Union type flags
	IsOneOf bool
	IsAnyOf bool

	// Source location
	SourcePath string // JSON Path in spec
}

type Property struct {
	Name        string
	Type        string
	Format      string
	Description string
	Required    bool
	Nullable    bool
	Enum        []interface{} // Enum values if this is an enum

	// Field metadata
	Default    interface{} // Default value
	ReadOnly   bool        // Read-only field
	WriteOnly  bool        // Write-only field
	Deprecated bool        // Deprecated field

	// Constraints
	Pattern   string
	MinLength *int
	MaxLength *int
	Minimum   *float64
	Maximum   *float64
	MinItems  *int
	MaxItems  *int

	// For nested types
	RefType string // Referenced type name (for $ref or nested objects)
	Items   *Property

	// For objects
	Properties []Property

	// Additional properties for maps
	AdditionalProperties *Property

	Example interface{}
}

type Discriminator struct {
	PropertyName string
	Mapping      map[string]string
}
