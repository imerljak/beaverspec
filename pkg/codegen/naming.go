package codegen

import (
	"log"
	"regexp"
	"strings"
)

type NamingConverter interface {
	TypeName(string) string
	FieldName(string) string
	FunctionName(string) string
	VariableName(string) string
	FileName(string) string
}

// FilterEmptyStrings Filters out empty strings from input array
func filterEmptyStrings(input []string) []string {
	var filtered []string
	for _, w := range input {
		if w != "" {
			filtered = append(filtered, w)
		}
	}
	return filtered
}

// ToPascalCase converts a string to PascalCase
// Example: "user_name" -> "UserName", "userName" -> "UserName"
func ToPascalCase(s string) string {
	snake := ToSnakeCase(s)
	words := regexp.MustCompile(`[\s_-]`).Split(snake, -1)

	filtered := filterEmptyStrings(words)

	if len(filtered) == 0 {
		return ""
	}

	var result strings.Builder
	result.WriteString(Capitalize(filtered[0]))
	for _, word := range filtered[1:] {
		result.WriteString(Capitalize(word))
	}

	return result.String()
}

// ToSnakeCase converts a string to snake_case
// Example: "UserName" -> "user_name", "userName" -> "user_name"
func ToSnakeCase(s string) string {
	s = strings.ReplaceAll(s, "-", "_")
	replaced := regexp.MustCompile(`[A-Z]`).ReplaceAllFunc([]byte(s), func(b []byte) []byte {
		return append([]byte("_"), b...)
	})

	return strings.TrimPrefix(strings.ToLower(string(replaced)), "_")
}

// ToKebabCase converts a string to kebab-case
// Example: "UserName" -> "user-name"
func ToKebabCase(s string) string {
	s = strings.ReplaceAll(s, "_", "-")
	replaced := regexp.MustCompile(`[A-Z]`).ReplaceAllFunc([]byte(s), func(b []byte) []byte {
		return append([]byte("-"), b...)
	})

	return strings.TrimPrefix(strings.ToLower(string(replaced)), "-")
}

// Capitalize capitalizes the first letter of a string
func Capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

// Add ES to the end of the word
var esPattern = compileRegex(`(i?)(s|sh|ch|x|z)$`)

// Remove the y and add ies
var iesPattern = compileRegex(`(i?)[qwrtpsdfghjklzxcvbnm]y$`)

// Add es to the end of the word
var oesPattern = compileRegex(`(i?)[qwrtpsdfghjklzxcvbnm]o$`)

var pluralIrregulars = map[string]string{
	"person": "people",
	"child":  "children",
	"man":    "men",
	"woman":  "women",
	"foot":   "feet",
	"tooth":  "teeth",
	"goose":  "geese",
	"mouse":  "mice",
	"index":  "indices", // or "indexes"
	"schema": "schemas", // technical term, keep simple
}

// Pluralize converts a string to its plural form
func Pluralize(singular string) string {
	if plural := pluralIrregulars[singular]; plural != "" {
		return plural
	}

	if esPattern.MatchString(singular) {
		return singular + "es"
	}

	if iesPattern.MatchString(singular) {
		return singular[:len(singular)-1] + "ies"
	}

	if strings.HasSuffix(singular, "fe") {
		return singular[:len(singular)-2] + "ves"
	}

	if strings.HasSuffix(singular, "f") {
		return singular[:len(singular)-1] + "ves"
	}

	if oesPattern.MatchString(singular) {
		return singular + "s" // TODO: enhance to handle exceptions
	}

	return singular + "s"
}

var singularIrregulars = map[string]string{
	"people":   "person",
	"children": "child",
	"men":      "man",
	"women":    "woman",
	"feet":     "foot",
	"teeth":    "tooth",
	"geese":    "goose",
	"mice":     "mouse",
	"indices":  "index",
}

var esSingularPattern = compileRegex(`(i?)(ses|shes|ches|xes|zes)$`)

// Singularize converts a plural string to singular
func Singularize(plural string) string {

	if s := singularIrregulars[plural]; s != "" {
		return s
	}

	if strings.HasSuffix(plural, "ies") {
		return plural[:len(plural)-3] + "y"
	}

	if strings.HasSuffix(plural, "ves") {
		return plural[:len(plural)-3] + "f" // TODO: enhance this
	}

	if esSingularPattern.MatchString(plural) {
		return plural[:len(plural)-2]
	}

	if strings.HasSuffix(plural, "s") && !strings.HasSuffix(plural, "ss") {
		return plural[:len(plural)-1]
	}

	return plural
}

func compileRegex(pattern string) *regexp.Regexp {
	re, err := regexp.Compile(pattern)
	if err != nil {
		log.Fatalf("invalid regex %s: %v", pattern, err)
	}
	return re
}
