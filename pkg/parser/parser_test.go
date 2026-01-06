package parser

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
)

// TestParserLoadSpec tests loading a valid OpenAPI spec
func TestParserLoadSpec(t *testing.T) {
	// Create a temporary directory for test files
	tmpDir := t.TempDir()

	// Creata a simple openAPI spec file
	specContent := `
openapi: 3.0.0
info:
  title: Test API
  version: 1.0.0
  description: A test API
paths:
  /users:
    get:
      summary: List users
      responses:
        '200':
          description: Success
`

	// Write the spec to a temp file
	specPath := filepath.Join(tmpDir, "openapi.yaml")
	if err := os.WriteFile(specPath, []byte(specContent), 0644); err != nil {
		t.Fatalf("failed to write test spec: %v", err)
	}

	// Test loading the spec
	parser := NewParser()
	ctx := context.Background()

	spec, err := parser.LoadSpec(ctx, specPath)
	if err != nil {
		t.Fatalf("LoadSpec failed: %v", err)
	}

	if spec == nil {
		t.Fatal("expected spec to be non-nil")
	}

	if spec.Info.Title != "Test API" {
		t.Errorf("expected title 'Test API', got '%s'", spec.Info.Title)
	}
}

// TestParserLoadSpecInvalidFile tests loading a non-existent file
func TestParserLoadSpecInvalidFile(t *testing.T) {
	parser := NewParser()
	ctx := context.Background()

	_, err := parser.LoadSpec(ctx, "/nonexistent/file.yaml")
	if err == nil {
		t.Error("expected error for non-existent file, got nil")
	}
}

// TestNormalizerExtractInfo tests the info extraction
func TestNormalizerExtractInfo(t *testing.T) {
	// TODO: Test the normalizer's extractInfo method
	// Create a sample openapi3.Info and verify conversion
	oai := openapi3.Info{
		Title:       "Info title",
		Description: "Info description",
		Version:     "1.0.0",
		Contact: &openapi3.Contact{
			Name:  "Contact Name",
			Email: "contact@mail.com",
			URL:   "https://www.contact.com",
		},
		License: &openapi3.License{
			Name: "GPL",
			URL:  "https://www.gpl.com",
		},
	}

	n := NewNormalizer()

	specInfo := n.extractInfo(&oai)

	if specInfo.Title != "Info title" {
		t.Errorf("expected title 'Info title', got '%s'", specInfo.Title)
	}

	if specInfo.Contact == nil {
		t.Error("expected Contact to be non-nil")
	}

	if specInfo.Contact.Name != "Contact Name" {
		t.Errorf("expected Contact.name 'Contact Name', got '%s'", specInfo.Contact.Name)
	}

	if specInfo.License == nil {
		t.Error("expected License to be non-nil")
	}

	if specInfo.License.Name != "GPL" {
		t.Errorf("expected License.Name 'GPL', got '%s'", specInfo.License.Name)
	}
}

// TestNormalizerExtractTags tests the tag extraction
func TestNormalizerExtractTags(t *testing.T) {
	tags := openapi3.Tags{
		&openapi3.Tag{
			Name:        "t1",
			Description: "t1 description",
		},
	}

	n := NewNormalizer()
	specTags := n.extractTags(tags)

	if len(specTags) == 0 {
		t.Fatal("expected tags to be non-empty")
	}

	if specTags[0].Name != "t1" {
		t.Errorf("expected Tag[0].Name to be 't1', got '%s'", specTags[0].Name)
	}

	if specTags[0].Description != "t1 description" {
		t.Errorf("expected Tag[0].Description to be 't1 description', got '%s'", specTags[0].Description)
	}
}
