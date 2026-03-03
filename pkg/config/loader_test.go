package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/imerljak/beaverspec/pkg/config"
	"github.com/imerljak/beaverspec/pkg/core"
)

func TestLoadFile_YAML(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "beaver.yaml")

	content := `
spec: ./openapi.yaml
output: ./generated
module: github.com/example/myapi
package: api
framework: echo

generate:
  models: true
  client: false
  server: true

exclude:
  models: [InternalError, DebugInfo]
  tags: [internal]

types:
  overrides:
    - format: email
      go: string
    - format: uuid
      go: string

templates:
  customDir: ./my-templates
  overrides:
    models.go.tmpl: ./custom-models.go.tmpl
`
	if err := os.WriteFile(cfgPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := config.LoadFile(cfgPath)
	if err != nil {
		t.Fatalf("LoadFile failed: %v", err)
	}

	if want := filepath.Join(dir, "openapi.yaml"); cfg.SpecPath != want {
		t.Errorf("SpecPath: got %q, want %q", cfg.SpecPath, want)
	}
	if want := filepath.Join(dir, "generated"); cfg.OutputDir != want {
		t.Errorf("OutputDir: got %q, want %q", cfg.OutputDir, want)
	}
	if got, _ := cfg.Options["modulePath"].(string); got != "github.com/example/myapi" {
		t.Errorf("modulePath option: got %q", got)
	}
	if got, _ := cfg.Options["package"].(string); got != "api" {
		t.Errorf("package option: got %q", got)
	}
	if got, _ := cfg.Options["framework"].(string); got != "echo" {
		t.Errorf("framework option: got %q", got)
	}

	if !cfg.Generate.Models {
		t.Error("Generate.Models should be true")
	}
	if cfg.Generate.Client {
		t.Error("Generate.Client should be false")
	}
	if !cfg.Generate.Server {
		t.Error("Generate.Server should be true")
	}

	if len(cfg.Exclude.Models) != 2 || cfg.Exclude.Models[0] != "InternalError" {
		t.Errorf("Exclude.Models: got %v", cfg.Exclude.Models)
	}
	if len(cfg.Exclude.Tags) != 1 || cfg.Exclude.Tags[0] != "internal" {
		t.Errorf("Exclude.Tags: got %v", cfg.Exclude.Tags)
	}

	if len(cfg.Types.Overrides) != 2 {
		t.Fatalf("Types.Overrides: expected 2, got %d", len(cfg.Types.Overrides))
	}
	if cfg.Types.Overrides[0].Format != "email" || cfg.Types.Overrides[0].Go != "string" {
		t.Errorf("Types.Overrides[0]: got %+v", cfg.Types.Overrides[0])
	}

	if want := filepath.Join(dir, "my-templates"); cfg.Templates.CustomDir != want {
		t.Errorf("Templates.CustomDir: got %q, want %q", cfg.Templates.CustomDir, want)
	}
	if v := cfg.Templates.Overrides["models.go.tmpl"]; v != filepath.Join(dir, "custom-models.go.tmpl") {
		t.Errorf("Templates.Overrides: got %q", v)
	}
}

func TestFindConfigFile(t *testing.T) {
	dir := t.TempDir()

	// No config file yet
	if got := config.FindConfigFile(dir); got != "" {
		t.Errorf("expected empty, got %q", got)
	}

	// Create beaver.yaml
	cfgPath := filepath.Join(dir, "beaver.yaml")
	os.WriteFile(cfgPath, []byte("spec: openapi.yaml\n"), 0644)

	if got := config.FindConfigFile(dir); got != cfgPath {
		t.Errorf("expected %q, got %q", cfgPath, got)
	}
}

func TestMerge(t *testing.T) {
	base := &core.Config{
		SpecPath:  "base.yaml",
		OutputDir: "base-out",
		Options: map[string]interface{}{
			"package":   "basepkg",
			"framework": "echo",
		},
		Exclude: core.ExcludeConfig{
			Models: []string{"BaseModel"},
		},
	}

	cli := &core.Config{
		SpecPath:  "cli.yaml",
		OutputDir: "",
		Options: map[string]interface{}{
			"package": "clipkg",
		},
		Exclude: core.ExcludeConfig{
			Models: []string{"CLIModel"},
		},
	}

	result := config.Merge(base, cli)

	if result.SpecPath != "cli.yaml" {
		t.Errorf("SpecPath: expected cli wins, got %q", result.SpecPath)
	}
	if result.OutputDir != "base-out" {
		t.Errorf("OutputDir: expected base preserved, got %q", result.OutputDir)
	}
	if got, _ := result.Options["package"].(string); got != "clipkg" {
		t.Errorf("package option: CLI should win, got %q", got)
	}
	if got, _ := result.Options["framework"].(string); got != "echo" {
		t.Errorf("framework option: base should be preserved, got %q", got)
	}

	// Both exclude lists should be combined
	if len(result.Exclude.Models) != 2 {
		t.Errorf("Exclude.Models: expected 2, got %d: %v", len(result.Exclude.Models), result.Exclude.Models)
	}
}
