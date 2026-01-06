package template

import (
	"testing"
)

type TestCase struct {
	Input       string
	Expectation string
}

func TestCapitalize(t *testing.T) {
	tests := []TestCase{
		{Input: "banana", Expectation: "Banana"},
		{Input: "hello world", Expectation: "Hello world"},
		{Input: "", Expectation: ""},
		{Input: "A", Expectation: "A"},
	}

	exec(t, tests, NamedFn{Name: "capitalize", Fn: capitalize})
}

func TestToCamelCase(t *testing.T) {
	tests := []TestCase{
		{Input: "user_name", Expectation: "userName"},
		{Input: "User_Name", Expectation: "userName"},
		{Input: "User Name", Expectation: "userName"},
		{Input: "User name", Expectation: "userName"},
		{Input: "user name", Expectation: "userName"},
		{Input: "UserName", Expectation: "userName"},
		{Input: "user-name", Expectation: "userName"},
		{Input: "", Expectation: ""},
		{Input: "userName", Expectation: "userName"},
		{Input: "API_KEY", Expectation: "aPIKEY"}, // TODO: fix issue with acronyms in the future.
		{Input: "get_user_by_id", Expectation: "getUserById"},
	}

	exec(t, tests, NamedFn{Name: "toCamelCase", Fn: toCamelCase})
}

func TestToPascalCase(t *testing.T) {
	tests := []TestCase{
		{Input: "user_name", Expectation: "UserName"},
		{Input: "User_Name", Expectation: "UserName"},
		{Input: "User Name", Expectation: "UserName"},
		{Input: "User name", Expectation: "UserName"},
		{Input: "user name", Expectation: "UserName"},
		{Input: "UserName", Expectation: "UserName"},
		{Input: "userName", Expectation: "UserName"},
		{Input: "user-name-test", Expectation: "UserNameTest"},
	}

	exec(t, tests, NamedFn{Name: "toPascalCase", Fn: toPascalCase})
}

func TestToSnakeCase(t *testing.T) {
	rules := []TestCase{
		{Input: "UserName", Expectation: "user_name"},
		{Input: "userName", Expectation: "user_name"},
		{Input: "user_name", Expectation: "user_name"},
		{Input: "user-name-test", Expectation: "user_name_test"},
		{Input: "HTTPSConnection", Expectation: "h_t_t_p_s_connection"},
		{Input: "getUserByID", Expectation: "get_user_by_i_d"},
	}

	exec(t, rules, NamedFn{Name: "toSnakeCase", Fn: toSnakeCase})
}

func TestToKebabCase(t *testing.T) {
	rules := []TestCase{
		{Input: "UserName", Expectation: "user-name"},
		{Input: "userName", Expectation: "user-name"},
		{Input: "user-name", Expectation: "user-name"},
		{Input: "user_name", Expectation: "user-name"},
	}

	exec(t, rules, NamedFn{Name: "toKebabCase", Fn: toKebabCase})
}

type Fn func(string) string
type NamedFn struct {
	Name string
	Fn   Fn
}

func exec(t *testing.T, rules []TestCase, fn NamedFn) {
	for _, rule := range rules {
		if result := fn.Fn(rule.Input); result != rule.Expectation {
			t.Fatalf("%s(%s) -> expected %s, got %s", fn.Name, rule.Input, rule.Expectation, result)
		}
	}
}
