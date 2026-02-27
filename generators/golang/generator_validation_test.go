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
			// Required check for name
			if !strings.Contains(content, "required field is missing") {
				t.Error("Missing required field validation message")
			}
			// MinLength check
			if !strings.Contains(content, "length must be at least") {
				t.Error("Missing minLength validation")
			}
			// Enum check
			if !strings.Contains(content, "is not a valid enum value") {
				t.Error("Missing enum validation")
			}
			// Validator interface
			if !strings.Contains(content, "type Validator interface") {
				t.Error("Missing Validator interface")
			}
		}
	}
}
