package golang

import (
	"os"
	"strings"
	"testing"

	"github.com/imerljak/beaverspec/pkg/core"
)

func TestMiddlewareGenerated(t *testing.T) {
	files := generateForFramework(t, "net-http")

	middleware := files["server/middleware.go"]
	if middleware == "" {
		t.Fatal("server/middleware.go not generated")
	}
	for _, want := range []string{
		"LoggingMiddleware",
		"CORSMiddleware",
		"RateLimiter",
		"MetricsStore",
		"log/slog",
	} {
		if !strings.Contains(middleware, want) {
			t.Errorf("server/middleware.go missing %q", want)
		}
	}
	// Should not contain any framework-specific identifiers
	for _, bad := range []string{"echo.Context", "*gin.Context", "chi.URLParam"} {
		if strings.Contains(middleware, bad) {
			t.Errorf("server/middleware.go should be framework-agnostic, found %q", bad)
		}
	}
}

func TestExampleMainNetHTTP(t *testing.T) {
	files := generateForFramework(t, "net-http")

	main := files["cmd/server/main.go"]
	if main == "" {
		t.Fatal("cmd/server/main.go not generated for net-http")
	}
	for _, want := range []string{
		"http.NewServeMux()",
		"srv.Shutdown(ctx)",
		"signal.Notify",
		"RegisterHealthRoutes",
		"Register",
	} {
		if !strings.Contains(main, want) {
			t.Errorf("net-http main.go missing %q", want)
		}
	}
}

func TestExampleMainChi(t *testing.T) {
	files := generateForFramework(t, "chi")

	main := files["cmd/server/main.go"]
	if main == "" {
		t.Fatal("cmd/server/main.go not generated for chi")
	}
	for _, want := range []string{
		"chi.NewRouter()",
		"chimiddleware.RequestID",
		"srv.Shutdown(ctx)",
	} {
		if !strings.Contains(main, want) {
			t.Errorf("chi main.go missing %q", want)
		}
	}
}

func TestExampleMainEcho(t *testing.T) {
	files := generateForFramework(t, "echo")

	main := files["cmd/server/main.go"]
	if main == "" {
		t.Fatal("cmd/server/main.go not generated for echo")
	}
	for _, want := range []string{
		"echo.New()",
		"e.Shutdown(ctx)",
		"echomiddleware",
	} {
		if !strings.Contains(main, want) {
			t.Errorf("echo main.go missing %q", want)
		}
	}
}

func TestExampleMainGin(t *testing.T) {
	files := generateForFramework(t, "gin")

	main := files["cmd/server/main.go"]
	if main == "" {
		t.Fatal("cmd/server/main.go not generated for gin")
	}
	for _, want := range []string{
		"gin.New()",
		"srv.Shutdown(ctx)",
	} {
		if !strings.Contains(main, want) {
			t.Errorf("gin main.go missing %q", want)
		}
	}
}

func TestHealthCheckRoute(t *testing.T) {
	for _, fw := range []string{"net-http", "chi", "echo", "gin"} {
		fw := fw
		t.Run(fw, func(t *testing.T) {
			files := generateForFramework(t, fw)
			routes := files["server/routes.go"]
			if !strings.Contains(routes, "RegisterHealthRoutes") {
				t.Errorf("%s routes.go missing RegisterHealthRoutes", fw)
			}
		})
	}
}

func TestMiddlewareOptOut(t *testing.T) {
	originalWd, _ := os.Getwd()
	os.Chdir("../../")
	defer os.Chdir(originalWd)

	gen := NewGenerator()
	result, err := gen.Generate(frameworkSpec(), &core.Config{
		OutputDir: "generated",
		Options: map[string]interface{}{
			"framework":  "net-http",
			"modulePath": "github.com/example/project",
			"middleware": false,
		},
	})
	if err != nil {
		t.Fatalf("generation failed: %v", err)
	}
	files := make(map[string]string)
	for _, f := range result.Files {
		files[f.Path] = string(f.Content)
	}
	if files["server/middleware.go"] != "" {
		t.Error("server/middleware.go should not be generated when middleware=false")
	}
	if files["cmd/server/main.go"] == "" {
		t.Error("cmd/server/main.go should still be generated when only middleware=false")
	}
}

func TestExampleMainOptOut(t *testing.T) {
	originalWd, _ := os.Getwd()
	os.Chdir("../../")
	defer os.Chdir(originalWd)

	gen := NewGenerator()
	result, err := gen.Generate(frameworkSpec(), &core.Config{
		OutputDir: "generated",
		Options: map[string]interface{}{
			"framework":   "net-http",
			"modulePath":  "github.com/example/project",
			"exampleMain": false,
		},
	})
	if err != nil {
		t.Fatalf("generation failed: %v", err)
	}
	files := make(map[string]string)
	for _, f := range result.Files {
		files[f.Path] = string(f.Content)
	}
	if files["cmd/server/main.go"] != "" {
		t.Error("cmd/server/main.go should not be generated when exampleMain=false")
	}
	if files["server/middleware.go"] == "" {
		t.Error("server/middleware.go should still be generated when only exampleMain=false")
	}
}
