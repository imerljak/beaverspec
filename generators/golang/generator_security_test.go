package golang

import (
	"os"
	"strings"
	"testing"

	"github.com/imerljak/beaverspec/pkg/core"
)

func TestGeneratorSecuritySchemes(t *testing.T) {
	gen := NewGenerator()

	ep := core.Endpoint{
		OperationID: "getProtected",
		Path:        "/protected",
		Method:      "GET",
		Responses:   []core.Response{{StatusCode: "200"}},
	}

	spec := &core.Spec{
		Endpoints: []core.Endpoint{ep},
		SecurityDef: []core.SecurityScheme{
			{
				Name:   "bearerAuth",
				Type:   "http",
				Scheme: "bearer",
			},
			{
				Name: "apiKeyAuth",
				Type: "apiKey",
				In:   "header",
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
			// Bearer token helper
			if !strings.Contains(string(f.Content), "WithBearerAuthToken") {
				t.Error("Client missing WithBearerAuthToken helper")
			}
			// API key helper
			if !strings.Contains(string(f.Content), "WithApiKeyAuthAPIKey") {
				t.Error("Client missing WithApiKeyAuthAPIKey helper")
			}
		}
		if f.Path == "server/interface.go" {
			if !strings.Contains(string(f.Content), "BearerAuthAuthMiddleware") {
				t.Error("Server interface missing BearerAuth middleware interface")
			}
			if !strings.Contains(string(f.Content), "ValidateBearerToken") {
				t.Error("Server interface missing ValidateBearerToken method stub")
			}
		}
	}
}
