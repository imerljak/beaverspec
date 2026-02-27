package golang

import (
	"os"
	"strings"
	"testing"

	"github.com/imerljak/beaverspec/pkg/core"
)

func TestGeneratorErrors(t *testing.T) {
	gen := NewGenerator()

	ep := core.Endpoint{
		OperationID: "auth",
		Path:        "/auth",
		Method:      "POST",
		Responses: []core.Response{
			{
				StatusCode: "200",
				Content: map[string]core.MediaType{
					"application/json": {Schema: &core.Property{RefType: "Session"}},
				},
			},
			{
				StatusCode: "401",
				Content: map[string]core.MediaType{
					"application/json": {Schema: &core.Property{RefType: "AuthError"}},
				},
			},
		},
	}

	spec := &core.Spec{
		Endpoints: []core.Endpoint{ep},
	}
	originalWd, _ := os.Getwd()
	os.Chdir("../../")
	defer os.Chdir(originalWd)

	result, err := gen.Generate(spec, &core.Config{
		OutputDir: "generated",
	})
	if err != nil {
		t.Fatalf("Generation failed: %v", err)
	}

	foundClient := false
	foundHandlers := false
	for _, f := range result.Files {
		if f.Path == "client/client.go" {
			foundClient = true
			if !strings.Contains(string(f.Content), "var apiErr *models.AuthError") {
				t.Error("Client not parsing mapped 401 AuthError")
			}
			if !strings.Contains(string(f.Content), "case 401:") {
				t.Error("Client missing 401 case switch")
			}
		}
		if f.Path == "server/handlers.go" {
			foundHandlers = true
			if !strings.Contains(string(f.Content), "var apiErr models.APIError") {
				t.Error("Server handler missing specific APIError mapping parser")
			}
			if !strings.Contains(string(f.Content), "w.WriteHeader(apiErr.StatusCode())") {
				t.Error("Server not setting custom APIError status codes dynamically")
			}
		}
	}

	if !foundClient {
		t.Error("Did not generate client/client.go")
	}
	if !foundHandlers {
		t.Error("Did not generate server/handlers.go")
	}
}
