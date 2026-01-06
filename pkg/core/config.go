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
}

type GenerateConfig struct {
	Models, Client, Server, Validation bool
}

type TemplateConfig struct {
	CustomDir string            // Path to custom templates
	Overrides map[string]string // template name -> path
}
