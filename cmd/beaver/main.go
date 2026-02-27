package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/imerljak/beaverspec/generators/golang"
	"github.com/imerljak/beaverspec/pkg/core"
	"github.com/imerljak/beaverspec/pkg/parser"
	"github.com/imerljak/beaverspec/pkg/plugin"
)

const version = "0.1.0-alpha"

func main() {
	// 1. Parse flags
	specFilePath := flag.String("spec", ".", "OpenAPI spec file path")
	outputFilePath := flag.String("output", ".", "output directory path")
	generatorFlag := flag.String("generator", "go", "which generator to use, default 'go'")
	packageNameFlag := flag.String("package", "models", "package name for generated code")
	modulePathFlag := flag.String("module", "github.com/example/project", "base module path for the generated code")
	versionFlag := flag.Bool("version", false, "print version and exit")
	debugFlag := flag.Bool("debug", false, "enable debug logging")
	flag.Parse()

	// Setup structured logging
	logLevel := slog.LevelInfo
	if *debugFlag {
		logLevel = slog.LevelDebug
	}
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: logLevel,
	}))
	slog.SetDefault(logger)

	// Show version
	if *versionFlag {
		fmt.Printf("BeaverSpec v%s\n", version)
		fmt.Println("A modern OpenAPI code generator")
		os.Exit(0)
	}

	// 2. Initialize registry and register generators
	registry := plugin.NewRegistry()
	if err := registry.Register("go", golang.NewGenerator()); err != nil {
		slog.Error("failed to register generator", "error", err)
		os.Exit(1)
	}

	// 3. Load and parse spec
	ctx := context.Background()
	p := parser.NewParser()
	t, err := p.LoadSpec(ctx, *specFilePath)
	if err != nil {
		slog.Error("failed to load spec file", "path", *specFilePath, "error", err)
		os.Exit(1)
	}

	n := parser.NewNormalizer()
	spec, err := n.Normalize(t)
	if err != nil {
		slog.Error("failed to normalize OpenAPI specification", "error", err)
		os.Exit(1)
	}

	// 4. Get generator and generate code
	gen, err := registry.Get(*generatorFlag)
	if err != nil {
		slog.Error("failed to load generator", "generator", *generatorFlag, "error", err)
		os.Exit(1)
	}

	result, err := gen.Generate(spec, &core.Config{
		SpecPath:  *specFilePath,
		OutputDir: *outputFilePath,
		Generator: *generatorFlag,
		Options: map[string]interface{}{
			"package":    *packageNameFlag,
			"modulePath": *modulePathFlag,
		},
	})
	if err != nil {
		slog.Error("failed to generate files", "error", err)
		os.Exit(1)
	}

	for _, f := range result.Files {
		path := filepath.Join(*outputFilePath, f.Path)

		// 5. Create output directory if it doesn't exist
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			slog.Error("failed to create output directory", "path", filepath.Dir(path), "error", err)
			os.Exit(1)
		}

		// 5.1 Write files to output directory
		err := os.WriteFile(path, f.Content, 0644)
		if err != nil {
			slog.Error("failed to write file", "path", f.Path, "error", err)
			os.Exit(1)
		}
	}

	// 6. Print success message
	slog.Info("Code generation successful", "files_created", len(result.Files), "output_dir", *outputFilePath)
	for _, f := range result.Files {
		fmt.Printf("  - %s\n", f.Path)
	}
}
