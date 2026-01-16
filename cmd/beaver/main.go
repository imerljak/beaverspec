package main

import (
	"context"
	"flag"
	"fmt"
	"log"
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
	versionFlag := flag.Bool("version", false, "print version and exit")
	flag.Parse()

	// Show version
	if *versionFlag {
		fmt.Printf("BeaverSpec v%s\n", version)
		fmt.Println("A modern OpenAPI code generator")
		os.Exit(0)
	}

	// 2. Initialize registry and register generators
	registry := plugin.NewRegistry()
	if err := registry.Register("go", golang.NewGenerator()); err != nil {
		log.Fatalf("failed to register generator: %v", err)
	}

	// 3. Load and parse spec
	ctx := context.Background()
	p := parser.NewParser()
	t, err := p.LoadSpec(ctx, *specFilePath)
	if err != nil {
		log.Fatalf("failed to load spec file: %v", err)
	}

	n := parser.NewNormalizer()
	spec, err := n.Normalize(t)
	if err != nil {
		log.Fatalf("failed to normalize OpenAPI specification: %v", err)
	}

	// 4. Get generator and generate code
	gen, err := registry.Get(*generatorFlag)
	if err != nil {
		log.Fatalf("failed to load generator: %v", err)
	}

	result, err := gen.Generate(spec, &core.Config{
		SpecPath:  *specFilePath,
		OutputDir: *outputFilePath,
		Generator: *generatorFlag,
		Options: map[string]interface{}{
			"package": *packageNameFlag,
		},
	})
	if err != nil {
		log.Fatalf("failed to generate files: %v", err)
	}

	for _, f := range result.Files {
		path := filepath.Join(*outputFilePath, f.Path)

		// 5. Create output directory if it doesn't exist
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			log.Fatalf("failed to create output directory: %v", err)
		}

		// 5.1 Write files to output directory
		err := os.WriteFile(path, f.Content, 0644)
		if err != nil {
			log.Fatalf("failed to write file - %s : %v", f.Path, err)
		}
	}

	// 6. Print success message
	fmt.Printf("✓ Successfully generated %d file(s) in %s\n", len(result.Files), *outputFilePath)
	for _, f := range result.Files {
		fmt.Printf("  - %s\n", f.Path)
	}
}
