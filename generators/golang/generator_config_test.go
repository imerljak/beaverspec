package golang

import (
	"os"
	"strings"
	"testing"

	"github.com/imerljak/beaverspec/pkg/core"
)

func configTestSpec() *core.Spec {
	return &core.Spec{
		Models: []core.Model{
			{Name: "User", Properties: []core.Property{{Name: "id", Type: "string", Required: true}}},
			{Name: "InternalError", Properties: []core.Property{{Name: "code", Type: "integer"}}},
			{Name: "Product", Properties: []core.Property{{Name: "name", Type: "string", Required: true}}},
		},
		Endpoints: []core.Endpoint{
			{
				OperationID: "getUser",
				Method:      "GET",
				Path:        "/users/{id}",
				Tags:        []string{"users"},
				Parameters: []core.Parameter{{Name: "id", In: "path", Schema: &core.Property{Type: "string"}}},
			},
			{
				OperationID: "getInternalDebug",
				Method:      "GET",
				Path:        "/debug",
				Tags:        []string{"internal"},
			},
		},
	}
}

func TestExcludeModels(t *testing.T) {
	originalWd, _ := os.Getwd()
	os.Chdir("../../")
	defer os.Chdir(originalWd)

	gen := NewGenerator()
	spec := configTestSpec()

	result, err := gen.Generate(spec, &core.Config{
		OutputDir: "generated",
		Options:   map[string]interface{}{"modulePath": "github.com/example/project"},
		Exclude: core.ExcludeConfig{
			Models: []string{"InternalError"},
		},
	})
	if err != nil {
		t.Fatalf("generation failed: %v", err)
	}

	modelsFile := findFile(result.Files, "models/models.go")
	if modelsFile == nil {
		t.Fatal("models/models.go not found")
	}
	content := string(modelsFile.Content)

	if strings.Contains(content, "InternalError") {
		t.Error("InternalError should be excluded from models.go")
	}
	if !strings.Contains(content, "type User struct") {
		t.Error("User model should still be present")
	}
	if !strings.Contains(content, "type Product struct") {
		t.Error("Product model should still be present")
	}
}

func TestExcludeTags(t *testing.T) {
	originalWd, _ := os.Getwd()
	os.Chdir("../../")
	defer os.Chdir(originalWd)

	gen := NewGenerator()
	spec := configTestSpec()

	result, err := gen.Generate(spec, &core.Config{
		OutputDir: "generated",
		Options:   map[string]interface{}{"modulePath": "github.com/example/project"},
		Exclude: core.ExcludeConfig{
			Tags: []string{"internal"},
		},
	})
	if err != nil {
		t.Fatalf("generation failed: %v", err)
	}

	// Server interface should not contain the internal operation
	interfaceFile := findFile(result.Files, "server/interface.go")
	if interfaceFile == nil {
		t.Fatal("server/interface.go not found")
	}
	content := string(interfaceFile.Content)

	if strings.Contains(content, "GetInternalDebug") {
		t.Error("internal tag operation should be excluded from server interface")
	}
	if !strings.Contains(content, "GetUser") {
		t.Error("users tag operation should still be present")
	}
}

func TestSelectiveGeneration_ModelsOnly(t *testing.T) {
	originalWd, _ := os.Getwd()
	os.Chdir("../../")
	defer os.Chdir(originalWd)

	gen := NewGenerator()
	spec := configTestSpec()

	result, err := gen.Generate(spec, &core.Config{
		OutputDir: "generated",
		Options:   map[string]interface{}{"modulePath": "github.com/example/project"},
		Generate: core.GenerateConfig{
			Models: true,
			Client: false,
			Server: false,
		},
	})
	if err != nil {
		t.Fatalf("generation failed: %v", err)
	}

	for _, f := range result.Files {
		if strings.HasPrefix(f.Path, "client/") {
			t.Errorf("client file should not be generated when Generate.Client=false: %s", f.Path)
		}
		if strings.HasPrefix(f.Path, "server/") {
			t.Errorf("server file should not be generated when Generate.Server=false: %s", f.Path)
		}
	}

	if findFile(result.Files, "models/models.go") == nil {
		t.Error("models/models.go should be generated when Generate.Models=true")
	}
}

func TestTypeOverrides(t *testing.T) {
	originalWd, _ := os.Getwd()
	os.Chdir("../../")
	defer os.Chdir(originalWd)

	gen := NewGenerator()
	spec := &core.Spec{
		Models: []core.Model{
			{
				Name: "Contact",
				Properties: []core.Property{
					{Name: "email", Type: "string", Format: "email"},
					{Name: "age", Type: "integer"},
				},
			},
		},
	}

	result, err := gen.Generate(spec, &core.Config{
		OutputDir: "generated",
		Options:   map[string]interface{}{"modulePath": "github.com/example/project"},
		Types: core.TypesConfig{
			Overrides: []core.TypeOverride{
				{Format: "email", Go: "EmailAddress"},
				{Type: "integer", Go: "int64"},
			},
		},
	})
	if err != nil {
		t.Fatalf("generation failed: %v", err)
	}

	modelsFile := findFile(result.Files, "models/models.go")
	if modelsFile == nil {
		t.Fatal("models/models.go not found")
	}
	content := string(modelsFile.Content)

	if !strings.Contains(content, "EmailAddress") {
		t.Error("email format should be mapped to EmailAddress via type override")
	}
	if !strings.Contains(content, "int64") {
		t.Error("integer type should be mapped to int64 via type override")
	}
}

func findFile(files []core.GeneratedFile, path string) *core.GeneratedFile {
	for i := range files {
		if files[i].Path == path {
			return &files[i]
		}
	}
	return nil
}
