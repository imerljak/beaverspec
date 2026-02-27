package codegen

// Method constants for generated HTTP routers if frameworks need them.
const (
	MethodGet     = "GET"
	MethodHead    = "HEAD"
	MethodPost    = "POST"
	MethodPut     = "PUT"
	MethodPatch   = "PATCH"
	MethodDelete  = "DELETE"
	MethodConnect = "CONNECT"
	MethodOptions = "OPTIONS"
	MethodTrace   = "Trace"
)

// SupportedMediaTypes
const (
	MediaTypeJSON           = "application/json"
	MediaTypeXML            = "application/xml"
	MediaTypeFormURLEncoded = "application/x-www-form-urlencoded"
	MediaTypeMultipartForm  = "multipart/form-data"
	MediaTypeTextPlain      = "text/plain"
)
