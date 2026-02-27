package golang

import (
	"os"
	"strings"
	"testing"

	"github.com/imerljak/beaverspec/pkg/core"
)

func frameworkSpec() *core.Spec {
	return &core.Spec{
		Models: []core.Model{
			{
				Name: "Item",
				Properties: []core.Property{
					{Name: "id", Type: "string"},
					{Name: "name", Type: "string"},
				},
			},
		},
		Endpoints: []core.Endpoint{
			{
				OperationID: "getItem",
				Path:        "/items/{id}",
				Method:      "GET",
				Tags:        []string{"items"},
				Parameters: []core.Parameter{
					{Name: "id", In: "path", Required: true, Schema: &core.Property{Type: "string"}},
				},
				Responses: []core.Response{
					{
						StatusCode: "200",
						Content: map[string]core.MediaType{
							"application/json": {Schema: &core.Property{RefType: "Item"}},
						},
					},
				},
			},
		},
	}
}

func generateForFramework(t *testing.T, framework string) map[string]string {
	t.Helper()
	gen := NewGenerator()

	originalWd, _ := os.Getwd()
	os.Chdir("../../")
	defer os.Chdir(originalWd)

	result, err := gen.Generate(frameworkSpec(), &core.Config{
		OutputDir: "generated",
		Options: map[string]interface{}{
			"framework": framework,
		},
	})
	if err != nil {
		t.Fatalf("Generation failed for framework %s: %v", framework, err)
	}

	files := make(map[string]string)
	for _, f := range result.Files {
		files[f.Path] = string(f.Content)
	}
	return files
}

func TestFrameworkNetHTTP(t *testing.T) {
	files := generateForFramework(t, "net-http")

	handlers := files["server/handlers.go"]
	routes := files["server/routes.go"]

	if !strings.Contains(handlers, "http.ResponseWriter") {
		t.Error("net-http handlers should use http.ResponseWriter")
	}
	if !strings.Contains(handlers, "r.PathValue(") {
		t.Error("net-http handlers should use r.PathValue for path params")
	}
	if !strings.Contains(routes, "*http.ServeMux") {
		t.Error("net-http routes should use *http.ServeMux")
	}
	if !strings.Contains(routes, "mux.HandleFunc(") {
		t.Error("net-http routes should use mux.HandleFunc")
	}
	// net/http uses {param} syntax
	if !strings.Contains(routes, "{id}") {
		t.Error("net-http routes should use {id} path param syntax")
	}
}

func TestFrameworkChi(t *testing.T) {
	files := generateForFramework(t, "chi")

	handlers := files["server/handlers.go"]
	routes := files["server/routes.go"]

	if !strings.Contains(handlers, "chi.URLParam(") {
		t.Error("chi handlers should use chi.URLParam for path params")
	}
	if !strings.Contains(routes, "chi.Router") {
		t.Error("chi routes should accept chi.Router")
	}
	// chi uses {param} syntax
	if !strings.Contains(routes, "{id}") {
		t.Error("chi routes should use {id} path param syntax")
	}
	if !strings.Contains(routes, "go-chi/chi") {
		t.Error("chi routes should import go-chi/chi")
	}
}

func TestFrameworkEcho(t *testing.T) {
	files := generateForFramework(t, "echo")

	handlers := files["server/handlers.go"]
	routes := files["server/routes.go"]

	if !strings.Contains(handlers, "echo.Context") {
		t.Error("echo handlers should use echo.Context")
	}
	if !strings.Contains(handlers, "c.Param(") {
		t.Error("echo handlers should use c.Param for path params")
	}
	if !strings.Contains(handlers, "c.JSON(") {
		t.Error("echo handlers should use c.JSON for responses")
	}
	if !strings.Contains(routes, "*echo.Echo") {
		t.Error("echo routes should accept *echo.Echo")
	}
	// echo uses :param syntax
	if !strings.Contains(routes, ":id") {
		t.Error("echo routes should use :id path param syntax")
	}
}

func TestFrameworkGin(t *testing.T) {
	files := generateForFramework(t, "gin")

	handlers := files["server/handlers.go"]
	routes := files["server/routes.go"]

	if !strings.Contains(handlers, "*gin.Context") {
		t.Error("gin handlers should use *gin.Context")
	}
	if !strings.Contains(handlers, "c.Param(") {
		t.Error("gin handlers should use c.Param for path params")
	}
	if !strings.Contains(handlers, "c.JSON(") {
		t.Error("gin handlers should use c.JSON for responses")
	}
	if !strings.Contains(routes, "*gin.Engine") {
		t.Error("gin routes should accept *gin.Engine")
	}
	// gin uses :param syntax
	if !strings.Contains(routes, ":id") {
		t.Error("gin routes should use :id path param syntax")
	}
}
