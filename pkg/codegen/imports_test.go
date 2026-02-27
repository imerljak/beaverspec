package codegen

import (
	"testing"
)

func TestImportManager(t *testing.T) {
	manager := NewImportManager("github.com/example/project", "generated", "client")

	// Test Add
	alias1 := manager.Add("fmt")
	if alias1 != "fmt" {
		t.Errorf("Expected alias 'fmt', got '%s'", alias1)
	}

	// Test AddWithAlias
	alias2 := manager.AddWithAlias("encoding/json", "js")
	if alias2 != "js" {
		t.Errorf("Expected alias 'js', got '%s'", alias2)
	}

	// Test AddSibling
	// currentPkg is "client", so it should add "github.com/example/project/models"
	alias3 := manager.AddSibling("models")
	if alias3 != "models" {
		t.Errorf("Expected alias 'models', got '%s'", alias3)
	}

	// Test AddSibling to self
	aliasSelf := manager.AddSibling("client")
	if aliasSelf != "" {
		t.Errorf("Expected empty alias for self import, got '%s'", aliasSelf)
	}

	// Test HasImports
	if !manager.HasImports() {
		t.Error("Expected HasImports to be true")
	}

	// Test GetImports (should be sorted alphabetically by path)
	imports := manager.GetImports()
	if len(imports) != 3 {
		t.Fatalf("Expected 3 imports, got %d", len(imports))
	}

	if imports[0].Path != "encoding/json" || imports[0].Alias != "js" {
		t.Errorf("Expected first import to be encoding/json with alias js, got %v", imports[0])
	}
	if imports[1].Path != "fmt" || imports[1].Alias != "" {
		t.Errorf("Expected second import to be fmt, got %v", imports[1])
	}
	if imports[2].Path != "github.com/example/project/models" || imports[2].Alias != "" {
		t.Errorf("Expected third import to be models, got %v", imports[2])
	}
}
