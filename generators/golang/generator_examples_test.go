package golang

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/imerljak/beaverspec/pkg/core"
	"github.com/imerljak/beaverspec/pkg/parser"
)

// loadExampleSpec parses an example spec file and returns a normalized core.Spec.
// Caller is responsible for setting the working directory to the project root first.
func loadExampleSpec(t *testing.T, specPath string) *core.Spec {
	t.Helper()
	ctx := context.Background()
	p := parser.NewParser()
	raw, err := p.LoadSpec(ctx, specPath)
	if err != nil {
		t.Fatalf("LoadSpec(%q) failed: %v", specPath, err)
	}
	n := parser.NewNormalizer()
	spec, err := n.Normalize(raw)
	if err != nil {
		t.Fatalf("Normalize(%q) failed: %v", specPath, err)
	}
	return spec
}

// TestAllExampleSpecsGenerate loops over every *.yaml in examples/ and asserts
// that generation succeeds and produces at least one non-empty file.
func TestAllExampleSpecsGenerate(t *testing.T) {
	originalWd, _ := os.Getwd()
	os.Chdir("../../")
	defer os.Chdir(originalWd)

	entries, err := filepath.Glob("examples/*.yaml")
	if err != nil || len(entries) == 0 {
		t.Fatal("no example specs found in examples/")
	}

	gen := NewGenerator()

	for _, specPath := range entries {
		specPath := specPath
		name := strings.TrimSuffix(filepath.Base(specPath), ".yaml")
		t.Run(name, func(t *testing.T) {
			spec := loadExampleSpec(t, specPath)
			result, err := gen.Generate(spec, &core.Config{
				OutputDir: "generated",
				Options: map[string]interface{}{
					"framework":  "net-http",
					"modulePath": "github.com/example/test",
				},
			})
			if err != nil {
				t.Fatalf("Generate failed: %v", err)
			}
			if len(result.Files) == 0 {
				t.Error("expected at least one generated file")
			}
			for _, f := range result.Files {
				if len(f.Content) == 0 {
					t.Errorf("file %q is empty", f.Path)
				}
			}
		})
	}
}

// TestPetstoreFileRefIntegration is an end-to-end test using the petstore-file-ref
// example which exercises external $ref resolution across multiple files.
func TestPetstoreFileRefIntegration(t *testing.T) {
	originalWd, _ := os.Getwd()
	os.Chdir("../../")
	defer os.Chdir(originalWd)

	spec := loadExampleSpec(t, "examples/petstore-file-ref/spec.yaml")
	gen := NewGenerator()
	result, err := gen.Generate(spec, &core.Config{
		OutputDir: "generated",
		Options: map[string]interface{}{
			"framework":  "net-http",
			"modulePath": "github.com/example/petstore",
		},
	})
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	files := make(map[string]string)
	for _, f := range result.Files {
		files[f.Path] = string(f.Content)
	}

	required := []string{
		"models/models.go",
		"server/interface.go",
	}
	for _, path := range required {
		if files[path] == "" {
			t.Errorf("expected file %q to be generated", path)
		}
	}

	models := files["models/models.go"]
	for _, want := range []string{"type Pet struct", "Pets"} {
		if !strings.Contains(models, want) {
			t.Errorf("models/models.go missing %q", want)
		}
	}

	iface := files["server/interface.go"]
	for _, want := range []string{"ListPets", "CreatePets", "ShowPetById"} {
		if !strings.Contains(iface, want) {
			t.Errorf("server/interface.go missing %q", want)
		}
	}
}
