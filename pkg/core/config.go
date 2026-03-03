package core

type Config struct {

	// Source
	SpecPath string

	// Output
	OutputDir string

	// Generator selection
	Generator string

	// What to generate
	Generate GenerateConfig

	// Generator-specific options
	Options map[string]interface{}

	// Template customization
	Templates TemplateConfig

	// Exclusions
	Exclude ExcludeConfig

	// Custom type mappings
	Types TypesConfig
}

type GenerateConfig struct {
	Models, Client, Server, Validation bool
}

type TemplateConfig struct {
	CustomDir string            // Path to custom templates directory (supplements built-ins)
	Overrides map[string]string // template name → file path (replaces individual templates)
}

// ExcludeConfig specifies models and tags to skip during generation.
type ExcludeConfig struct {
	Models []string // model names to skip (exact match)
	Tags   []string // operation tags to skip
}

// TypesConfig holds custom OpenAPI → Go type mappings.
type TypesConfig struct {
	Overrides []TypeOverride
}

// TypeOverride maps an OpenAPI format or type to a custom Go type.
// Format takes precedence over Type when both are set.
type TypeOverride struct {
	Format string // OpenAPI format (e.g. "email", "uuid", "date-time")
	Type   string // OpenAPI type (e.g. "integer") — matched when Format is ""
	Go     string // Target Go type (e.g. "uuid.UUID", "EmailAddress")
}
