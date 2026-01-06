package template

import (
	"regexp"
	"strings"
	"text/template"
)

// DefaultHelpers returns a map of default template helper functions
func DefaultHelpers() template.FuncMap {
	return template.FuncMap{
		"toCamelCase":  toCamelCase,
		"toPascalCase": toPascalCase,
		"toSnakeCase":  toSnakeCase,
		"toKebabCase":  toKebabCase,
		"toLower":      strings.ToLower,
		"toUpper":      strings.ToUpper,
		"capitalize":   capitalize,
	}
}

// filterEmptyStrings Filters out empty strings from input array
func filterEmptyStrings(input []string) []string {
	var filtered []string
	for _, w := range input {
		if w != "" {
			filtered = append(filtered, w)
		}
	}
	return filtered
}

// toCamelCase converts a string to camelCase
// Example: "user_name" -> "userName", "UserName" -> "userName"
func toCamelCase(s string) string {
	snake := toSnakeCase(s)
	words := regexp.MustCompile(`[\s_-]`).Split(snake, -1)

	// Filter out empty strings
	filtered := filterEmptyStrings(words)

	if len(filtered) == 0 {
		return ""
	}

	// First word lowecase, rest capitalized
	var result strings.Builder
	result.WriteString(strings.ToLower(filtered[0]))
	for _, word := range filtered[1:] {
		result.WriteString(capitalize(word))
	}
	return result.String()
}

// toPascalCase converts a string to PascalCase
// Example: "user_name" -> "UserName", "userName" -> "UserName"
func toPascalCase(s string) string {
	snake := toSnakeCase(s)
	words := regexp.MustCompile(`[\s_-]`).Split(snake, -1)

	filtered := filterEmptyStrings(words)

	if len(filtered) == 0 {
		return ""
	}

	var result strings.Builder
	result.WriteString(capitalize(filtered[0]))
	for _, word := range filtered[1:] {
		result.WriteString(capitalize(word))
	}

	return result.String()
}

// toSnakeCase converts a string to snake_case
// Example: "UserName" -> "user_name", "userName" -> "user_name"
func toSnakeCase(s string) string {
	s = strings.ReplaceAll(s, "-", "_")
	replaced := regexp.MustCompile(`[A-Z]`).ReplaceAllFunc([]byte(s), func(b []byte) []byte {
		return append([]byte("_"), b...)
	})

	return strings.TrimPrefix(strings.ToLower(string(replaced)), "_")
}

// toKebabCase converts a string to kebab-case
// Example: "UserName" -> "user-name"
func toKebabCase(s string) string {
	s = strings.ReplaceAll(s, "_", "-")
	replaced := regexp.MustCompile(`[A-Z]`).ReplaceAllFunc([]byte(s), func(b []byte) []byte {
		return append([]byte("-"), b...)
	})

	return strings.TrimPrefix(strings.ToLower(string(replaced)), "-")
}

// capitalize capitalizes the first letter of a string
func capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
