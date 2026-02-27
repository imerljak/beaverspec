package template

import (
	"regexp"
	"strings"
	"text/template"
)

// DefaultHelpers returns a map of default template helper functions
func DefaultHelpers() template.FuncMap {
	return template.FuncMap{
		"toCamelCase":  ToCamelCase,
		"toPascalCase": ToPascalCase,
		"toSnakeCase":  ToSnakeCase,
		"toKebabCase":  ToKebabCase,
		"toLower":      strings.ToLower,
		"toUpper":      strings.ToUpper,
		"capitalize":   Capitalize,
		"deref": func(p *int) int {
			if p != nil {
				return *p
			}
			return 0
		},
		"derefF": func(p *float64) float64 {
			if p != nil {
				return *p
			}
			return 0
		},
		"zeroCheck": ZeroCheck,
	}
}

// ZeroCheck returns the Go zero-value comparison snippet for the given type.
// Used to generate required-field validation in Validate() methods.
func ZeroCheck(goType string) string {
	switch goType {
	case "string":
		return `== ""`
	case "int", "int8", "int16", "int32", "int64",
		"uint", "uint8", "uint16", "uint32", "uint64",
		"float32", "float64":
		return "== 0"
	case "bool":
		return "== false"
	default:
		// Pointer types, interfaces, slices, maps
		return "== nil"
	}
}

// FilterEmptyStrings Filters out empty strings from input array
func FilterEmptyStrings(input []string) []string {
	var filtered []string
	for _, w := range input {
		if w != "" {
			filtered = append(filtered, w)
		}
	}
	return filtered
}

// ToCamelCase converts a string to camelCase
// Example: "user_name" -> "userName", "UserName" -> "userName"
func ToCamelCase(s string) string {
	snake := ToSnakeCase(s)
	words := regexp.MustCompile(`[\s_-]`).Split(snake, -1)

	// Filter out empty strings
	filtered := FilterEmptyStrings(words)

	if len(filtered) == 0 {
		return ""
	}

	// First word lowecase, rest capitalized
	var result strings.Builder
	result.WriteString(strings.ToLower(filtered[0]))
	for _, word := range filtered[1:] {
		result.WriteString(Capitalize(word))
	}
	return result.String()
}

// ToPascalCase converts a string to PascalCase
// Example: "user_name" -> "UserName", "userName" -> "UserName"
func ToPascalCase(s string) string {
	snake := ToSnakeCase(s)
	words := regexp.MustCompile(`[\s_-]`).Split(snake, -1)

	filtered := FilterEmptyStrings(words)

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
