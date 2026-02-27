package codegen

import (
	"testing"
)

func TestToPascalCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"user_name", "UserName"},
		{"userName", "UserName"},
		{"User-Name", "UserName"},
		{"user name", "UserName"},
		{"UserName", "UserName"},
		{"", ""},
	}

	for _, tt := range tests {
		actual := ToPascalCase(tt.input)
		if actual != tt.expected {
			t.Errorf("ToPascalCase(%q): expected %q, got %q", tt.input, tt.expected, actual)
		}
	}
}

func TestToSnakeCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"UserName", "user_name"},
		{"userName", "user_name"},
		{"user-name", "user_name"},
		{"user_name", "user_name"},
		{"", ""},
	}

	for _, tt := range tests {
		actual := ToSnakeCase(tt.input)
		if actual != tt.expected {
			t.Errorf("ToSnakeCase(%q): expected %q, got %q", tt.input, tt.expected, actual)
		}
	}
}

func TestToKebabCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"UserName", "user-name"},
		{"userName", "user-name"},
		{"user_name", "user-name"},
		{"user-name", "user-name"},
		{"", ""},
	}

	for _, tt := range tests {
		actual := ToKebabCase(tt.input)
		if actual != tt.expected {
			t.Errorf("ToKebabCase(%q): expected %q, got %q", tt.input, tt.expected, actual)
		}
	}
}

func TestCapitalize(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello", "Hello"},
		{"Hello", "Hello"},
		{"h", "H"},
		{"", ""},
	}

	for _, tt := range tests {
		actual := Capitalize(tt.input)
		if actual != tt.expected {
			t.Errorf("Capitalize(%q): expected %q, got %q", tt.input, tt.expected, actual)
		}
	}
}

func TestPluralize(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		// Regulars
		{"cat", "cats"},
		{"dog", "dogs"},
		{"apple", "apples"},
		// Ends with s, sh, ch, x, z
		{"bus", "buses"},
		{"bush", "bushes"},
		{"watch", "watches"},
		{"box", "boxes"},
		// Ends with y (consonant + y -> ies)
		{"city", "cities"},
		{"puppy", "puppies"},
		// Ends with f or fe (-> ves)
		{"wolf", "wolves"},
		{"wife", "wives"},
		// Ends with o (consonant + o -> oes or os, simplistically tests "oes")
		{"hero", "heroes"},
		{"potato", "potatoes"},
		// Irregulars
		{"person", "people"},
		{"child", "children"},
		{"man", "men"},
		{"woman", "women"},
		{"foot", "feet"},
		{"tooth", "teeth"},
		{"goose", "geese"},
		{"mouse", "mice"},
		{"index", "indices"},
		{"schema", "schemas"},
	}

	for _, tt := range tests {
		actual := Pluralize(tt.input)
		if actual != tt.expected {
			t.Errorf("Pluralize(%q): expected %q, got %q", tt.input, tt.expected, actual)
		}
	}
}

func TestSingularize(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		// Regulars
		{"cats", "cat"},
		{"dogs", "dog"},
		{"apples", "apple"},
		// Ends with es
		{"buses", "bus"},
		{"bushes", "bush"},
		{"watches", "watch"},
		{"boxes", "box"},
		// Ends with ies
		{"cities", "city"},
		{"puppies", "puppy"},
		// Ends with ves
		{"wolves", "wolf"},
		{"wives", "wife"},
		// Irregulars
		{"people", "person"},
		{"children", "child"},
		{"men", "man"},
		{"women", "woman"},
		{"feet", "foot"},
		{"teeth", "tooth"},
		{"geese", "goose"},
		{"mice", "mouse"},
		{"indices", "index"},
	}

	for _, tt := range tests {
		actual := Singularize(tt.input)
		if actual != tt.expected {
			t.Errorf("Singularize(%q): expected %q, got %q", tt.input, tt.expected, actual)
		}
	}
}
