package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/imerljak/beaverspec/pkg/core"
)

// configFile is the on-disk representation of beaver.yaml / beaver.json.
// All fields are optional; missing fields fall back to CLI defaults.
type configFile struct {
	Spec      string          `yaml:"spec"      json:"spec"`
	Output    string          `yaml:"output"    json:"output"`
	Module    string          `yaml:"module"    json:"module"`
	Package   string          `yaml:"package"   json:"package"`
	Framework string          `yaml:"framework" json:"framework"`
	Generate  generateBlock   `yaml:"generate"  json:"generate"`
	Exclude   excludeBlock    `yaml:"exclude"   json:"exclude"`
	Types     typesBlock      `yaml:"types"     json:"types"`
	Templates templateBlock   `yaml:"templates" json:"templates"`
}

type generateBlock struct {
	Models     *bool `yaml:"models"     json:"models"`
	Client     *bool `yaml:"client"     json:"client"`
	Server     *bool `yaml:"server"     json:"server"`
	Validation *bool `yaml:"validation" json:"validation"`
}

type excludeBlock struct {
	Models []string `yaml:"models" json:"models"`
	Tags   []string `yaml:"tags"   json:"tags"`
}

type typeOverrideEntry struct {
	Format string `yaml:"format" json:"format"`
	Type   string `yaml:"type"   json:"type"`
	Go     string `yaml:"go"     json:"go"`
}

type typesBlock struct {
	Overrides []typeOverrideEntry `yaml:"overrides" json:"overrides"`
}

type templateBlock struct {
	CustomDir string            `yaml:"customDir"  json:"customDir"`
	Overrides map[string]string `yaml:"overrides"  json:"overrides"`
}

// LoadFile reads and parses a beaver.yaml or beaver.json config file.
// Paths in the config file are resolved relative to the config file's directory.
func LoadFile(path string) (*core.Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file %s: %w", path, err)
	}

	var cf configFile
	switch filepath.Ext(path) {
	case ".json":
		if err := json.Unmarshal(data, &cf); err != nil {
			return nil, fmt.Errorf("parsing JSON config %s: %w", path, err)
		}
	default: // .yaml, .yml, or no extension
		if err := yaml.Unmarshal(data, &cf); err != nil {
			return nil, fmt.Errorf("parsing YAML config %s: %w", path, err)
		}
	}

	baseDir := filepath.Dir(path)
	cfg := toConfig(cf, baseDir)
	return cfg, nil
}

// FindConfigFile searches dir for beaver.yaml or beaver.json.
// Returns the absolute path if found, or "" if neither exists.
func FindConfigFile(dir string) string {
	for _, name := range []string{"beaver.yaml", "beaver.yml", "beaver.json"} {
		candidate := filepath.Join(dir, name)
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
	}
	return ""
}

// Merge returns a new Config where non-zero values in cli override base.
// Use this to let explicit CLI flags win over the config file.
func Merge(base, cli *core.Config) *core.Config {
	if base == nil {
		base = &core.Config{}
	}
	if cli == nil {
		return base
	}

	result := *base // copy

	if cli.SpecPath != "" {
		result.SpecPath = cli.SpecPath
	}
	if cli.OutputDir != "" {
		result.OutputDir = cli.OutputDir
	}
	if cli.Generator != "" {
		result.Generator = cli.Generator
	}

	// Merge Options map — CLI values win
	if len(cli.Options) > 0 {
		if result.Options == nil {
			result.Options = make(map[string]interface{})
		}
		for k, v := range cli.Options {
			// Only override if CLI provided a non-empty string value
			if s, ok := v.(string); ok && s == "" {
				continue
			}
			result.Options[k] = v
		}
	}

	// GenerateConfig: CLI wins if any field is explicitly set
	if cli.Generate.Models || cli.Generate.Client || cli.Generate.Server || cli.Generate.Validation {
		result.Generate = cli.Generate
	}

	// Exclude: append (CLI adds to file config exclusions)
	result.Exclude.Models = append(result.Exclude.Models, cli.Exclude.Models...)
	result.Exclude.Tags = append(result.Exclude.Tags, cli.Exclude.Tags...)

	// Types: CLI overrides prepend (higher priority)
	result.Types.Overrides = append(cli.Types.Overrides, result.Types.Overrides...)

	// Templates: CLI wins
	if cli.Templates.CustomDir != "" {
		result.Templates.CustomDir = cli.Templates.CustomDir
	}
	if len(cli.Templates.Overrides) > 0 {
		if result.Templates.Overrides == nil {
			result.Templates.Overrides = make(map[string]string)
		}
		for k, v := range cli.Templates.Overrides {
			result.Templates.Overrides[k] = v
		}
	}

	return &result
}

// toConfig converts the on-disk struct to a core.Config, resolving relative paths.
func toConfig(cf configFile, baseDir string) *core.Config {
	cfg := &core.Config{
		SpecPath:  resolvePath(cf.Spec, baseDir),
		OutputDir: resolvePath(cf.Output, baseDir),
		Options:   make(map[string]interface{}),
	}

	// Standard options stored in the Options map (for backward compat with generator)
	if cf.Module != "" {
		cfg.Options["modulePath"] = cf.Module
	}
	if cf.Package != "" {
		cfg.Options["package"] = cf.Package
	}
	if cf.Framework != "" {
		cfg.Options["framework"] = cf.Framework
	}

	// GenerateConfig — if any field specified use them; otherwise leave all false
	// (generator treats all-false as "generate everything")
	if cf.Generate.Models != nil {
		cfg.Generate.Models = *cf.Generate.Models
	}
	if cf.Generate.Client != nil {
		cfg.Generate.Client = *cf.Generate.Client
	}
	if cf.Generate.Server != nil {
		cfg.Generate.Server = *cf.Generate.Server
	}
	if cf.Generate.Validation != nil {
		cfg.Generate.Validation = *cf.Generate.Validation
	}

	// Exclusions
	cfg.Exclude = core.ExcludeConfig{
		Models: cf.Exclude.Models,
		Tags:   cf.Exclude.Tags,
	}

	// Type overrides
	for _, o := range cf.Types.Overrides {
		cfg.Types.Overrides = append(cfg.Types.Overrides, core.TypeOverride{
			Format: o.Format,
			Type:   o.Type,
			Go:     o.Go,
		})
	}

	// Templates
	cfg.Templates = core.TemplateConfig{
		CustomDir: resolvePath(cf.Templates.CustomDir, baseDir),
		Overrides: cf.Templates.Overrides,
	}
	// Resolve override paths relative to config file dir
	for k, v := range cfg.Templates.Overrides {
		cfg.Templates.Overrides[k] = resolvePath(v, baseDir)
	}

	return cfg
}

// resolvePath resolves path relative to baseDir if it's not absolute.
// Returns "" unchanged.
func resolvePath(path, baseDir string) string {
	if path == "" || filepath.IsAbs(path) {
		return path
	}
	return filepath.Join(baseDir, path)
}
