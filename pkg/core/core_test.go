package core

import "testing"

func TestSpec(t *testing.T) {
	spec := &Spec{
		Info: SpecInfo{
			Title:       "Test API",
			Description: "A test API",
			Version:     "1.0.0",
		},
		Models:    []Model{},
		Endpoints: []Endpoint{},
	}

	if spec.Info.Title != "Test API" {
		t.Errorf("expected title 'Test API', got '%s'", spec.Info.Title)
	}
}

func TestConfigCreation(t *testing.T) {
	config := &Config{
		Generator: "go",
		SpecPath:  "./openapi.yaml",
		OutputDir: "./generated",
		Generate: GenerateConfig{
			Models: true,
			Client: true,
			Server: false,
		},
	}

	if !config.Generate.Models {
		t.Error("expected Models to be true")
	}
}
