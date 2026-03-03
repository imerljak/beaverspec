package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/imerljak/beaverspec/generators/golang"
	"github.com/imerljak/beaverspec/pkg/config"
	"github.com/imerljak/beaverspec/pkg/core"
	"github.com/imerljak/beaverspec/pkg/parser"
	"github.com/imerljak/beaverspec/pkg/plugin"
)

const version = "0.1.0-alpha"

func main() {
	// 1. Parse flags
	configFlag := flag.String("config", "", "config file path (default: beaver.yaml in spec directory)")
	specFilePath := flag.String("spec", "", "OpenAPI spec file path")
	outputFilePath := flag.String("output", "", "output directory path")
	generatorFlag := flag.String("generator", "go", "which generator to use, default 'go'")
	packageNameFlag := flag.String("package", "", "package name for generated code")
	modulePathFlag := flag.String("module", "", "base module path for the generated code")
	frameworkFlag := flag.String("framework", "", "server framework: net-http, chi, echo, gin")
	dryRunFlag := flag.Bool("dry-run", false, "print files that would be generated without writing them")
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

	// 2. Build config: start from file (if any), then apply CLI overrides
	cfg := &core.Config{}

	// Determine config file path: explicit flag or auto-discover next to spec
	cfgPath := *configFlag
	if cfgPath == "" && *specFilePath != "" {
		cfgPath = config.FindConfigFile(filepath.Dir(*specFilePath))
	}
	if cfgPath != "" {
		loaded, err := config.LoadFile(cfgPath)
		if err != nil {
			slog.Error("failed to load config file", "path", cfgPath, "error", err)
			os.Exit(1)
		}
		cfg = loaded
		slog.Debug("loaded config file", "path", cfgPath)
	}

	// CLI flags override file config (only when explicitly provided)
	cliOpts := map[string]interface{}{}
	if *packageNameFlag != "" {
		cliOpts["package"] = *packageNameFlag
	}
	if *modulePathFlag != "" {
		cliOpts["modulePath"] = *modulePathFlag
	}
	if *frameworkFlag != "" {
		cliOpts["framework"] = *frameworkFlag
	}
	cliConfig := &core.Config{
		SpecPath:  *specFilePath,
		OutputDir: *outputFilePath,
		Generator: *generatorFlag,
		Options:   cliOpts,
	}
	cfg = config.Merge(cfg, cliConfig)

	// Apply defaults for required fields that were not set anywhere
	if cfg.Options == nil {
		cfg.Options = map[string]interface{}{}
	}
	if _, ok := cfg.Options["package"]; !ok {
		cfg.Options["package"] = "models"
	}
	if _, ok := cfg.Options["modulePath"]; !ok {
		cfg.Options["modulePath"] = "github.com/example/project"
	}
	if _, ok := cfg.Options["framework"]; !ok {
		cfg.Options["framework"] = "net-http"
	}
	if cfg.Generator == "" {
		cfg.Generator = "go"
	}

	if cfg.SpecPath == "" {
		slog.Error("spec file path is required (use -spec flag or set 'spec' in beaver.yaml)")
		os.Exit(1)
	}
	if cfg.OutputDir == "" {
		cfg.OutputDir = "."
	}

	// 3. Initialize registry and register generators
	registry := plugin.NewRegistry()
	if err := registry.Register("go", golang.NewGenerator()); err != nil {
		slog.Error("failed to register generator", "error", err)
		os.Exit(1)
	}

	// 4. Load and parse spec
	ctx := context.Background()
	p := parser.NewParser()
	t, err := p.LoadSpec(ctx, cfg.SpecPath)
	if err != nil {
		slog.Error("failed to load spec file", "path", cfg.SpecPath, "error", err)
		os.Exit(1)
	}

	n := parser.NewNormalizer()
	spec, err := n.Normalize(t)
	if err != nil {
		slog.Error("failed to normalize OpenAPI specification", "error", err)
		os.Exit(1)
	}

	// 5. Get generator and generate code
	gen, err := registry.Get(cfg.Generator)
	if err != nil {
		slog.Error("failed to load generator", "generator", cfg.Generator, "error", err)
		os.Exit(1)
	}

	result, err := gen.Generate(spec, cfg)
	if err != nil {
		slog.Error("failed to generate files", "error", err)
		os.Exit(1)
	}

	// 6. Dry-run: just list files, don't write
	if *dryRunFlag {
		fmt.Printf("Dry run — %d file(s) would be generated to %s:\n", len(result.Files), cfg.OutputDir)
		for _, f := range result.Files {
			fmt.Printf("  %s\n", f.Path)
		}
		return
	}

	// 7. Write files to output directory
	for _, f := range result.Files {
		path := filepath.Join(cfg.OutputDir, f.Path)
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			slog.Error("failed to create output directory", "path", filepath.Dir(path), "error", err)
			os.Exit(1)
		}
		if err := os.WriteFile(path, f.Content, 0644); err != nil {
			slog.Error("failed to write file", "path", f.Path, "error", err)
			os.Exit(1)
		}
	}

	slog.Info("Code generation successful", "files_created", len(result.Files), "output_dir", cfg.OutputDir)
	for _, f := range result.Files {
		fmt.Printf("  - %s\n", f.Path)
	}
}
