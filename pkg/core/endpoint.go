package core

type Endpoint struct {
	Path        string
	Method      string // GET, POST, PUT, DELETE, etc.
	OperationID string
	Summary     string
	Description string
	Tags        []string

	Parameters  []Parameter
	RequestBody *RequestBody
	Responses   []Response

	Security     []SecurityRequirement
	IsDeprecated bool
}

type Parameter struct {
	Name         string
	In           string // "path", "query", "header", "cookie"
	Description  string
	Required     bool
	Schema       *Property
	Example      interface{}
	IsDeprecated bool
}

type RequestBody struct {
	Description string
	Required    bool
	Content     map[string]MediaType // media type -> schema
}

type MediaType struct {
	Schema  *Property
	Example interface{}
}

type Response struct {
	StatusCode  string // "200", "404", "default"
	Description string
	Headers     map[string]Header
	Content     map[string]MediaType
}

type Header struct {
	Description string
	Required    bool
	Schema      *Property
}

type SecurityRequirement struct {
	Name   string
	Scopes []string
}
