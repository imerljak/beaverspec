package core

// Generator is the interface that all language
// generators must implement
type Generator interface {

	// Metadata returs informatio about this generator
	Metadata() GeneratorMetadata

	// Validate checks if the spec and config are valid for this generator
	// Returns validation errors (does not stop on first error)
	Validate(spec *Spec, config *Config) []ValidationError

	// Generate produces code from the spec
	Generate(spec *Spec, config *Config) (*GenerationResult, error)

	// SupportedFeatures returns OpenAPI features this generator supports
	SupportedFeatures() []Feature
}

// GeneratorMetadata describes the generator's capabilities
type GeneratorMetadata struct {
	Name         string       // Unique identifier (e.g., "go-generator")
	Language     string       // Target language (e.g., "go", "typescript")
	Version      string       // Semantic version
	Description  string       // Human-readable description
	Author       string       // Author/maintainer
	Capabilities Capabilities // What this generator can produce
}

// Capabilities describes what a generator can produce
type Capabilities struct {
	SupportsClient     bool     // Can generate client code
	SupportsServer     bool     // Can generate server code
	SupportsValidation bool     // Can generate validators
	SupportsAsync      bool     // Supports async/await patterns
	SupportsStreaming  bool     // Supports streaming (SSE, Websockets)
	ClientFrameworks   []string // Supported client frameworks
	ServerFrameworks   []string // Supported server frameworks
}

// Feature represents an OpeAPI feature
type Feature string

const (
	FeatureNullable      Feature = "nullable"
	FeatureDiscriminator Feature = "discriminator"
	FeatureOneOf         Feature = "oneOf"
	FeatureAllOf         Feature = "allOf"
	FeatureCallbacks     Feature = "callbacks"
	FeatureWebhooks      Feature = "webhooks"
	FeatureDeprecated    Feature = "deprecated"
	FeatureExamples      Feature = "examples"
	FeatureReadOnly      Feature = "readOnly"
	FeatureWriteOnly     Feature = "writeOnly"
)
