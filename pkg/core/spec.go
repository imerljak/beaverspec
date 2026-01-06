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

	// Complex type handling
	AllOf []Model
	OneOf []Model
	AnyOf []Model

	// Metadata
	Discriminator *Discriminator
	IsDeprecated  bool
	Example       interface{}

	// Source location
	SourcePath string // JSON Path in spec
}

type Property struct {
	Name        string
	Type        string
	Description string
	Required    bool
	Nullable    bool
	ReadOnly    bool
	WriteOnly   bool
	Format      string
	Pattern     string
	Default     interface{}
	Example     interface{}

	// Constraints
	MinLength *int
	MaxLength *int
	Minimum   *float64
	Maximum   *float64
	MinItems  *int
	MaxItems  *int

	// For arrays
	Items *Property

	// For objects
	Properties []Property
}

type Discriminator struct {
	PropertyName string
	Mapping      map[string]string
}
