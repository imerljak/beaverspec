package golang

import (
	"os"
	"strings"
	"testing"

	"github.com/imerljak/beaverspec/pkg/core"
)

func TestGeneratorValidation(t *testing.T) {
	gen := NewGenerator()

	minLen := 3
	maxLen := 100
	minimum := 0.0
	maximum := 999.0
	pattern := `^\d{4}-\d{2}-\d{2}$`

	spec := &core.Spec{
		Models: []core.Model{
			{
				Name: "CreateUserRequest",
				Properties: []core.Property{
					{
						Name:      "name",
						Type:      "string",
						Required:  true,
						MinLength: &minLen,
						MaxLength: &maxLen,
					},
					{
						Name:     "email",
						Type:     "string",
						Format:   "email",
						Required: true,
					},
					{
						Name:    "age",
						Type:    "integer",
						Minimum: &minimum,
						Maximum: &maximum,
					},
					{
						Name: "role",
						Type: "string",
						Enum: []interface{}{"admin", "user", "guest"},
					},
					{
						Name:    "birthDate",
						Type:    "string",
						Pattern: pattern,
					},
				},
			},
		},
	}

	originalWd, _ := os.Getwd()
	os.Chdir("../../")
	defer os.Chdir(originalWd)

	result, err := gen.Generate(spec, &core.Config{OutputDir: "generated"})
	if err != nil {
		t.Fatalf("Generation failed: %v", err)
	}

	for _, f := range result.Files {
		if f.Path == "models/models.go" {
			content := string(f.Content)
			// Should have Validate method
			if !strings.Contains(content, "func (m *CreateUserRequest) Validate() error") {
				t.Error("Missing Validate() method for CreateUserRequest")
			}
			// Required check
			if !strings.Contains(content, "required field is missing") {
				t.Error("Missing required field validation message")
			}
			// MinLength check
			if !strings.Contains(content, "length must be at least") {
				t.Error("Missing minLength validation")
			}
			// MaxLength check
			if !strings.Contains(content, "length must be at most") {
				t.Error("Missing maxLength validation")
			}
			// Minimum check
			if !strings.Contains(content, "value must be >=") {
				t.Error("Missing minimum validation")
			}
			// Maximum check
			if !strings.Contains(content, "value must be <=") {
				t.Error("Missing maximum validation")
			}
			// Pattern check
			if !strings.Contains(content, "does not match required pattern") {
				t.Error("Missing pattern validation")
			}
			// Enum check
			if !strings.Contains(content, "is not a valid enum value") {
				t.Error("Missing enum validation")
			}
			// Format validation (email uses regexp)
			if !strings.Contains(content, "regexp.MatchString") {
				t.Error("Missing format validation (regexp)")
			}
			// Validator interface
			if !strings.Contains(content, "type Validator interface") {
				t.Error("Missing Validator interface")
			}
		}
	}
}

func TestGeneratorClientValidation(t *testing.T) {
	gen := NewGenerator()

	minLen := 1
	spec := &core.Spec{
		Models: []core.Model{
			{
				Name: "CreateItemRequest",
				Properties: []core.Property{
					{Name: "title", Type: "string", Required: true, MinLength: &minLen},
				},
			},
		},
		Endpoints: []core.Endpoint{
			{
				OperationID: "createItem",
				Path:        "/items",
				Method:      "POST",
				RequestBody: &core.RequestBody{
					Content: map[string]core.MediaType{
						"application/json": {Schema: &core.Property{RefType: "CreateItemRequest"}},
					},
				},
				Responses: []core.Response{
					{StatusCode: "201"},
				},
			},
		},
	}

	originalWd, _ := os.Getwd()
	os.Chdir("../../")
	defer os.Chdir(originalWd)

	result, err := gen.Generate(spec, &core.Config{OutputDir: "generated"})
	if err != nil {
		t.Fatalf("Generation failed: %v", err)
	}

	for _, f := range result.Files {
		if f.Path == "client/client.go" {
			content := string(f.Content)
			if !strings.Contains(content, "validator") {
				t.Error("Client missing validator duck-type interface")
			}
			if !strings.Contains(content, "v.Validate()") {
				t.Error("Client missing v.Validate() call")
			}
			if !strings.Contains(content, "validation:") {
				t.Error("Client missing validation error wrapping")
			}
		}
	}
}
