package codegen

import "strings"

// Helper functions for content-type routing/logic in templates

func IsJSON(contentType string) bool {
	return strings.Contains(contentType, MediaTypeJSON)
}

func IsXML(contentType string) bool {
	return strings.Contains(contentType, MediaTypeXML)
}

func IsFormURLEncoded(contentType string) bool {
	return strings.Contains(contentType, MediaTypeFormURLEncoded)
}

func IsMultipartForm(contentType string) bool {
	return strings.Contains(contentType, MediaTypeMultipartForm)
}

func IsTextPlain(contentType string) bool {
	return strings.Contains(contentType, MediaTypeTextPlain)
}
