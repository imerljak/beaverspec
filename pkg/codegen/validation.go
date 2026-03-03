package codegen

import (
	"net/mail"
	"net/url"
	"regexp"
	"time"
)

// Formats that require explicit validation logic
var validatedFormats = map[string]bool{
	"email":     true,
	"uuid":      true,
	"date-time": true,
	"date":      true,
	"uri":       true,
	"url":       true,
	"hostname":  true,
}

// IsFormatValidated returns true if the given OpenAPI format requires a validation check.
func IsFormatValidated(format string) bool {
	return validatedFormats[format]
}

// HasConstraints returns true if the property has any validation constraints defined.
func HasConstraints(hasMinLength, hasMaxLength, hasMinimum, hasMaximum, hasPattern, hasEnum, isRequired, hasMinItems, hasMaxItems, uniqueItems bool, multipleOf *float64) bool {
	return hasMinLength || hasMaxLength || hasMinimum || hasMaximum || hasPattern || hasEnum || isRequired || hasMinItems || hasMaxItems || uniqueItems || multipleOf != nil
}

// ValidateEmail returns true if the value is a valid email address (RFC 5322).
func ValidateEmail(v string) bool {
	_, err := mail.ParseAddress(v)
	return err == nil
}

var uuidRE = regexp.MustCompile(`(?i)^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)

// ValidateUUID returns true if the value matches the UUID format (RFC 4122).
func ValidateUUID(v string) bool {
	return uuidRE.MatchString(v)
}

// ValidateDate returns true if the value is a valid full-date (YYYY-MM-DD).
func ValidateDate(v string) bool {
	_, err := time.Parse("2006-01-02", v)
	return err == nil
}

// ValidateDateTime returns true if the value is a valid date-time (RFC 3339).
func ValidateDateTime(v string) bool {
	_, err := time.Parse(time.RFC3339, v)
	if err != nil {
		_, err = time.Parse(time.RFC3339Nano, v)
	}
	return err == nil
}

// ValidateURI returns true if the value is a valid absolute URI.
func ValidateURI(v string) bool {
	u, err := url.ParseRequestURI(v)
	return err == nil && u.Scheme != ""
}

var hostnameRE = regexp.MustCompile(`^([a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?\.)*[a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?$`)

// ValidateHostname returns true if the value is a valid hostname (RFC 1123).
func ValidateHostname(v string) bool {
	return len(v) <= 253 && hostnameRE.MatchString(v)
}
