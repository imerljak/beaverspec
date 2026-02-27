package codegen

// Formats that require explicit validation logic
var validatedFormats = map[string]bool{
	"email":     true,
	"uuid":      true,
	"date-time": true,
	"date":      true,
	"uri":       true,
	"url":       true,
}

// IsFormatValidated returns true if the given OpenAPI format requires a validation check.
func IsFormatValidated(format string) bool {
	return validatedFormats[format]
}

// HasConstraints returns true if the property has any validation constraints defined.
func HasConstraints(hasMinLength, hasMaxLength, hasMinimum, hasMaximum, hasPattern, hasEnum, isRequired bool) bool {
	return hasMinLength || hasMaxLength || hasMinimum || hasMaximum || hasPattern || hasEnum || isRequired
}
