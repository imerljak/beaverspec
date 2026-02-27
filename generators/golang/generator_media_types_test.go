package golang

import (
	"os"
	"strings"
	"testing"

	"github.com/imerljak/beaverspec/pkg/core"
)

func TestGeneratorMediaTypes(t *testing.T) {
	gen := NewGenerator()

	ep := core.Endpoint{
		OperationID: "postXmlData",
		Path:        "/data",
		Method:      "POST",
		RequestBody: &core.RequestBody{
			Content: map[string]core.MediaType{
				"application/xml": {
					Schema: &core.Property{RefType: "MyXmlModel"},
				},
			},
		},
		Responses: []core.Response{
			{
				StatusCode: "200",
				Content: map[string]core.MediaType{
					"application/xml": {
						Schema: &core.Property{RefType: "MyXmlModel"},
					},
				},
			},
		},
	}

	spec := &core.Spec{
		Models: []core.Model{
			{Name: "MyXmlModel", Properties: []core.Property{{Name: "id", Type: "string"}}},
		},
		Endpoints: []core.Endpoint{ep},
	}

	// Change working directory to project root so templates load correctly
	originalWd, _ := os.Getwd()
	os.Chdir("../../")
	defer os.Chdir(originalWd)

	result, err := gen.Generate(spec, &core.Config{
		OutputDir: "generated",
		Options: map[string]interface{}{
			"package":    "models",
			"modulePath": "github.com/example/api",
		},
	})
	if err != nil {
		t.Fatalf("Generation failed: %v", err)
	}

	foundClient := false
	foundHandlers := false
	for _, f := range result.Files {
		if f.Path == "client/client.go" {
			foundClient = true
			if !strings.Contains(string(f.Content), "\"encoding/xml\"") {
				t.Error("Client missing encoding/xml import")
			}
			if !strings.Contains(string(f.Content), "xml.Marshal(body)") {
				t.Error("Client not using xml.Marshal")
			}
			if !strings.Contains(string(f.Content), "xml.NewDecoder(resp.Body).Decode") {
				t.Error("Client not using xml.NewDecoder")
			}
			if !strings.Contains(string(f.Content), "application/xml") {
				t.Error("Client missing application/xml Content-Type / Accept values")
			}
		}
		if f.Path == "server/handlers.go" {
			foundHandlers = true
			if !strings.Contains(string(f.Content), "\"encoding/xml\"") {
				t.Error("Server handlers missing encoding/xml import")
			}
			if !strings.Contains(string(f.Content), "xml.NewDecoder(r.Body).Decode") {
				t.Error("Server handler not decoding XML")
			}
			if !strings.Contains(string(f.Content), "xml.NewEncoder(w).Encode") {
				t.Error("Server handler not encoding XML")
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
