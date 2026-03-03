package golang

import (
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/imerljak/beaverspec/pkg/codegen"
	"github.com/imerljak/beaverspec/pkg/parser"
)

func TestConvertEndpointsToClient(t *testing.T) {
	// Load the OpenAPI spec
	loader := openapi3.NewLoader()
	doc, err := loader.LoadFromFile("../../examples/petstore.yaml")
	if err != nil {
		t.Fatalf("Failed to load spec: %v", err)
	}

	// Normalize it
	normalizer := parser.NewNormalizer()
	spec, err := normalizer.Normalize(doc)
	if err != nil {
		t.Fatalf("Failed to normalize: %v", err)
	}

	// Create generator and convert endpoints
	gen := NewGenerator()
	clientData := gen.convertEndpointsToClient(spec.Endpoints, spec.SecurityDef, "github.com/example/project", "generated", nil)

	// Verify basic structure
	if clientData.InterfaceName != "Client" {
		t.Errorf("Expected interface name 'Client', got '%s'", clientData.InterfaceName)
	}

	if clientData.PackageName != "client" {
		t.Errorf("Expected package name 'client', got '%s'", clientData.PackageName)
	}

	// Should have 3 operations from petstore
	if len(clientData.Operations) != 3 {
		t.Fatalf("Expected 3 operations, got %d", len(clientData.Operations))
	}

	// Verify ListPets operation
	listPets := findOperation(clientData.Operations, "ListPets")
	if listPets == nil {
		t.Fatal("ListPets operation not found")
	}
	if listPets.Method != "GET" {
		t.Errorf("Expected ListPets method GET, got %s", listPets.Method)
	}
	if listPets.Path != "/pets" {
		t.Errorf("Expected ListPets path /pets, got %s", listPets.Path)
	}
	if len(listPets.QueryParams) != 1 {
		t.Errorf("Expected ListPets to have 1 query param, got %d", len(listPets.QueryParams))
	}
	if listPets.ReturnType != "*models.Pets" {
		t.Errorf("Expected ListPets return type *models.Pets, got %s", listPets.ReturnType)
	}

	// Verify CreatePets operation
	createPets := findOperation(clientData.Operations, "CreatePets")
	if createPets == nil {
		t.Fatal("CreatePets operation not found")
	}
	if createPets.Method != "POST" {
		t.Errorf("Expected CreatePets method POST, got %s", createPets.Method)
	}
	if !createPets.HasBody {
		t.Error("Expected CreatePets to have body")
	}
	if createPets.BodyType != "*models.Pet" {
		t.Errorf("Expected CreatePets body type *models.Pet, got %s", createPets.BodyType)
	}

	// Verify ShowPetById operation
	showPet := findOperation(clientData.Operations, "ShowPetById")
	if showPet == nil {
		t.Fatal("ShowPetById operation not found")
	}
	if len(showPet.PathParams) != 1 {
		t.Errorf("Expected ShowPetById to have 1 path param, got %d", len(showPet.PathParams))
	}
	if showPet.ReturnType != "*models.Pet" {
		t.Errorf("Expected ShowPetById return type *models.Pet, got %s", showPet.ReturnType)
	}

	// Verify imports
	if len(clientData.Imports) == 0 {
		t.Error("Expected some imports, got none")
	}

	// Should have context and net/http at minimum
	if !containsImport(clientData.Imports, "context") {
		t.Error("Expected 'context' import")
	}
	if !containsImport(clientData.Imports, "net/http") {
		t.Error("Expected 'net/http' import")
	}
}

func findOperation(ops []OperationData, name string) *OperationData {
	for _, op := range ops {
		if op.Name == name {
			return &op
		}
	}
	return nil
}

func containsImport(imports []codegen.Import, target string) bool {
	for _, imp := range imports {
		if imp.Path == target {
			return true
		}
	}
	return false
}
