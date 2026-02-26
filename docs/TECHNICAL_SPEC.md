# OpenAPI Code Generator - Technical Specification

## Table of Contents

1. [Core Interfaces](#core-interfaces)
2. [Data Models](#data-models)
3. [Plugin Contract](#plugin-contract)
4. [Configuration Schema](#configuration-schema)
5. [Template Data Models](#template-data-models)
6. [File Organization](#file-organization)
7. [CLI Specification](#cli-specification)

---

## Core Interfaces

### Generator Interface

All language generators must implement this interface:

```go
package core

import (
    "github.com/getkin/kin-openapi/openapi3"
)

// Generator is the interface that all language generators must implement
type Generator interface {
    // Metadata returns information about this generator
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
    SupportsStreaming  bool     // Supports streaming (SSE, WebSockets)
    
    ClientFrameworks   []string // Supported client frameworks
    ServerFrameworks   []string // Supported server frameworks
}

// Feature represents an OpenAPI feature
type Feature string

const (
    FeatureNullable       Feature = "nullable"
    FeatureDiscriminator  Feature = "discriminator"
    FeatureOneOf          Feature = "oneOf"
    FeatureAllOf          Feature = "allOf"
    FeatureAnyOf          Feature = "anyOf"
    FeatureCallbacks      Feature = "callbacks"
    FeatureWebhooks       Feature = "webhooks"
    FeatureDeprecated     Feature = "deprecated"
    FeatureExamples       Feature = "examples"
    FeatureReadOnly       Feature = "readOnly"
    FeatureWriteOnly      Feature = "writeOnly"
)
```

### Spec Model

```go
package core

// Spec represents a processed OpenAPI specification
type Spec struct {
    // Original OpenAPI document
    Raw *openapi3.T
    
    // Enriched/normalized data
    Info        SpecInfo
    Models      []Model
    Endpoints   []Endpoint
    Tags        []Tag
    SecurityDef []SecurityScheme
    
    // Metadata
    Version     string // OpenAPI version (e.g., "3.0.0")
    
    // Internal
    refResolver *RefResolver
}

type SpecInfo struct {
    Title       string
    Description string
    Version     string
    Contact     *Contact
    License     *License
}

type Contact struct {
    Name  string
    Email string
    URL   string
}

type License struct {
    Name string
    URL  string
}
```

### Model Definitions

```go
package core

// Model represents a schema/data type
type Model struct {
    Name        string
    Description string
    Type        string // "object", "array", "string", etc.
    Properties  []Property
    Required    []string
    Enum        []interface{}
    
    // Complex type handling
    AllOf       []Model
    OneOf       []Model
    AnyOf       []Model
    
    // Metadata
    Discriminator *Discriminator
    IsDeprecated  bool
    Example       interface{}
    
    // Source location
    SourcePath string // JSON path in spec
}

type Property struct {
    Name        string
    Type        string
    Description string
    Required    bool
    Nullable    bool
    ReadOnly    bool
    WriteOnly   bool
    Format      string
    Pattern     string
    Default     interface{}
    Example     interface{}
    
    // Constraints
    MinLength   *int
    MaxLength   *int
    Minimum     *float64
    Maximum     *float64
    MinItems    *int
    MaxItems    *int
    
    // For arrays
    Items       *Property
    
    // For objects
    Properties  []Property
}

type Discriminator struct {
    PropertyName string
    Mapping      map[string]string
}
```

### Endpoint Definitions

```go
package core

// Endpoint represents an API operation
type Endpoint struct {
    Path        string
    Method      string // GET, POST, PUT, DELETE, etc.
    OperationID string
    Summary     string
    Description string
    Tags        []string
    
    Parameters  []Parameter
    RequestBody *RequestBody
    Responses   []Response
    
    Security    []SecurityRequirement
    IsDeprecated bool
}

type Parameter struct {
    Name        string
    In          string // "path", "query", "header", "cookie"
    Description string
    Required    bool
    Schema      *Property
    Example     interface{}
    IsDeprecated bool
}

type RequestBody struct {
    Description string
    Required    bool
    Content     map[string]MediaType // media type -> schema
}

type MediaType struct {
    Schema  *Property
    Example interface{}
}

type Response struct {
    StatusCode  string // "200", "404", "default"
    Description string
    Headers     map[string]Header
    Content     map[string]MediaType
}

type Header struct {
    Description string
    Required    bool
    Schema      *Property
}

type SecurityRequirement struct {
    Name   string
    Scopes []string
}
```

### Configuration

```go
package core

// Config represents the complete configuration for generation
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
    Models     bool
    Client     bool
    Server     bool
    Validation bool
}

type TemplateConfig struct {
    CustomDir string                 // Path to custom templates
    Overrides map[string]string      // template name -> path
}
```

### Generation Result

```go
package core

// GenerationResult contains the output of code generation
type GenerationResult struct {
    Files      []GeneratedFile
    Warnings   []GenerationError
}

type GeneratedFile struct {
    Path     string                 // Relative path from output dir
    Content  []byte                 // File content
    Metadata map[string]interface{} // Optional metadata
}
```

### Error Handling

```go
package core

type ErrorSeverity int

const (
    ErrorSeverityWarning ErrorSeverity = iota
    ErrorSeverityError
    ErrorSeverityFatal
)

type GenerationError struct {
    Severity ErrorSeverity
    Phase    string    // "parsing", "validation", "generation"
    Location *Location
    Message  string
    Hint     string    // Suggestion to fix
    Code     string    // Error code (e.g., "E1001")
}

type Location struct {
    File   string // Spec file or template file
    Line   int
    Column int
    Path   string // JSON path (e.g., "#/paths/users/get")
}

func (e GenerationError) Error() string {
    return e.Message
}

type ValidationError struct {
    GenerationError
}
```

---

## Plugin Contract

### Base Generator

Provides common functionality that all generators can use:

```go
package generator

import (
    "text/template"
    "github.com/yourorg/beaverspec/core"
)

// BaseGenerator provides common functionality for all generators
type BaseGenerator struct {
    metadata core.GeneratorMetadata
    engine   TemplateEngine
    config   *PluginConfig
}

// PluginConfig is loaded from generator's config.yaml
type PluginConfig struct {
    Name         string
    Language     string
    Version      string
    Description  string
    
    TypeMappings map[string]TypeMapping
    Naming       NamingConventions
    
    NullableStrategy string // "pointer", "optional", "union"
    OptionalStrategy string // "pointer", "optional", "omitempty"
    
    Templates    TemplateManifest
    Frameworks   FrameworkSupport
    Defaults     map[string]interface{}
}

type TypeMapping struct {
    Default string
    Formats map[string]string // format -> type
}

type NamingConventions struct {
    Type      string // "PascalCase", "camelCase", etc.
    Field     string
    Parameter string
    Const     string
}

type TemplateManifest struct {
    Required []string
    Optional map[string][]string // category -> templates
}

type FrameworkSupport struct {
    Client []string
    Server []string
}
```

### Template Engine Interface

```go
package generator

// TemplateEngine handles template loading and rendering
type TemplateEngine interface {
    // Load templates from directory
    Load(templateDir string) error
    
    // Render a template with data
    Render(templateName string, data interface{}) (string, error)
    
    // Register custom function
    RegisterFunction(name string, fn interface{}) error
    
    // Check if template exists
    HasTemplate(name string) bool
}

// GoTemplateEngine implements TemplateEngine using Go's text/template
type GoTemplateEngine struct {
    templates *template.Template
    funcMap   template.FuncMap
}
```

---

## Configuration Schema

### Generator Config File (`config.yaml`)

Each generator has a `config.yaml` in its directory:

```yaml
# Metadata
name: "go"
language: "go"
version: "1.0.0"
description: "Go code generator with framework support"
author: "OpenAPI Gen Team"

# Capabilities
capabilities:
  client: true
  server: true
  validation: true
  async: true
  streaming: false

# Type mappings from OpenAPI to target language
type_mappings:
  integer:
    default: "int"
    int32: "int32"
    int64: "int64"
  string:
    default: "string"
    date: "time.Time"
    date-time: "time.Time"
    byte: "[]byte"
    binary: "[]byte"
    email: "string"
    uuid: "string"
  number:
    default: "float64"
    float: "float32"
    double: "float64"
  boolean:
    default: "bool"
  array:
    default: "[]{{.ItemType}}"
  object:
    default: "map[string]interface{}"

# Naming conventions
naming:
  type: "PascalCase"      # User, ProductInfo
  field: "PascalCase"     # ID, FirstName
  parameter: "camelCase"  # userID, productName
  const: "PascalCase"     # MaxRetries
  file: "snake_case"      # user_model.go

# Null/optional handling
nullable_strategy: "pointer"        # Use *Type for nullable
optional_strategy: "omitempty"      # Use omitempty tag for optional

# Template configuration
templates:
  required:
    - "core/models"
    - "core/types"
    - "core/errors"
  
  optional:
    client:
      - "client/base"
      - "client/http"
    server:
      - "server/base"
      - "server/handlers"
      - "server/routes"
    validation:
      - "validation/validators"

# Framework support
frameworks:
  client:
    - name: "http"
      description: "Standard net/http client"
      default: true
    - name: "resty"
      description: "Resty HTTP client"
      default: false
  
  server:
    - name: "http"
      description: "Standard net/http server"
      default: true
    - name: "chi"
      description: "Chi router"
      default: false
    - name: "echo"
      description: "Echo framework"
      default: false
    - name: "gin"
      description: "Gin framework"
      default: false

# Default generation settings
defaults:
  generate:
    models: true
    client: true
    server: false
    validation: false
  
  client:
    framework: "http"
  
  server:
    framework: "http"
  
  output:
    structure: "organized"  # "flat" | "organized"

# Language-specific options
options:
  use_generics: false
  json_tag_case: "camelCase"
  generate_constructors: true
  generate_getters: false
```

### User Config File (`~/.beaverspec/config.yaml`)

User-level defaults:

```yaml
# Global defaults for all projects
defaults:
  output: ./generated
  
  generate:
    validation: true
    models: true

# Per-generator defaults
generators:
  go:
    package: "api"
    use_generics: true
  
  typescript:
    target: "es2020"
    module_type: "esm"
  
  python:
    version: "3.10"
    use_dataclasses: true
```

### Project Config File (`.beaverspec.yaml`)

Project-specific configuration:

```yaml
# Which generator to use
generator: go

# OpenAPI spec location
spec: ./openapi.yaml

# Output directory
output: ./internal/api

# What to generate
generate:
  models: true
  client: true
  server: false
  validation: false

# Generator-specific configuration
go:
  package: "myservice"
  use_generics: true
  json_tag_case: "snake_case"

# Client configuration
client:
  framework: "resty"

# Server configuration  
server:
  framework: "chi"
  middleware: true

# Template customization
templates:
  custom_dir: ./custom-templates
  overrides:
    "client/http": "./templates/custom-client.go.tmpl"

# Output organization
output:
  structure: "organized"  # "flat" | "organized" | "custom"
  
  # Custom patterns (if structure: custom)
  patterns:
    models: "{{.OutputDir}}/models/{{.Name | toSnakeCase}}.go"
    client: "{{.OutputDir}}/client/{{.Tag | toSnakeCase}}_api.go"
    server: "{{.OutputDir}}/server/{{.Tag | toSnakeCase}}.go"
```

---

## Template Data Models

Data structures passed to templates during rendering:

### Root Template Data

```go
type TemplateData struct {
    // Spec metadata
    SpecInfo    SpecInfo
    PackageName string
    
    // Generated content
    Models      []ModelData
    Endpoints   []EndpointData
    Tags        []TagData
    
    // Configuration
    Config      *Config
    Options     map[string]interface{}
    
    // Utilities
    Imports     []string
}
```

### Model Template Data

```go
type ModelData struct {
    // Basic info
    Name        string  // Transformed according to naming convention
    RawName     string  // Original name from spec
    Description string
    
    // Type information
    Type        string  // Mapped to target language
    IsArray     bool
    IsObject    bool
    IsEnum      bool
    IsPrimitive bool
    
    // Properties
    Properties  []PropertyData
    Required    []string
    
    // Complex types
    AllOf       []ModelData
    OneOf       []ModelData
    AnyOf       []ModelData
    
    // Metadata
    HasDiscriminator bool
    Discriminator    *DiscriminatorData
    IsDeprecated     bool
    
    // Source
    SourcePath  string
}

type PropertyData struct {
    Name         string  // Transformed name
    RawName      string  // Original JSON name
    Type         string  // Mapped type
    Description  string
    
    // Flags
    IsRequired   bool
    IsNullable   bool
    IsReadOnly   bool
    IsWriteOnly  bool
    IsArray      bool
    IsObject     bool
    
    // Validation
    MinLength    *int
    MaxLength    *int
    Pattern      string
    Minimum      *float64
    Maximum      *float64
    
    // Default/Example
    Default      interface{}
    Example      interface{}
    
    // For arrays
    ItemType     string
    
    // JSON tag
    JsonTag      string
}
```

### Endpoint Template Data

```go
type EndpointData struct {
    // Basic info
    Path         string
    Method       string
    OperationID  string  // Transformed for method name
    RawOperationID string
    Summary      string
    Description  string
    Tags         []string
    
    // Parameters
    PathParams   []ParameterData
    QueryParams  []ParameterData
    HeaderParams []ParameterData
    
    // Request/Response
    HasRequestBody bool
    RequestBody    *RequestBodyData
    Responses      []ResponseData
    
    // Metadata
    IsDeprecated bool
    RequiresSecurity bool
}

type ParameterData struct {
    Name         string  // Transformed
    RawName      string  // Original
    Type         string  // Mapped type
    Description  string
    Required     bool
    Example      interface{}
}

type RequestBodyData struct {
    Description  string
    Required     bool
    ContentTypes []string
    Model        *ModelData  // Primary content type
}

type ResponseData struct {
    StatusCode   string
    Description  string
    Model        *ModelData  // Response body model
    IsDefault    bool
}
```

---

## File Organization

### Project Structure

```
beaverspec/
├── cmd/
│   └── beaverspec/
│       └── main.go              # CLI entry point
├── pkg/
│   ├── core/
│   │   ├── spec.go              # Spec model
│   │   ├── generator.go         # Generator interface
│   │   ├── config.go            # Configuration
│   │   └── errors.go            # Error types
│   ├── parser/
│   │   ├── parser.go            # OpenAPI parser
│   │   ├── normalizer.go        # Spec normalization
│   │   └── resolver.go          # $ref resolution
│   ├── plugin/
│   │   ├── registry.go          # Plugin registry
│   │   └── loader.go            # Plugin loading
│   ├── template/
│   │   ├── engine.go            # Template engine interface
│   │   ├── go_engine.go         # Go template implementation
│   │   └── helpers.go           # Template helper functions
│   └── cli/
│       ├── cli.go               # CLI application
│       ├── generate.go          # Generate command
│       ├── validate.go          # Validate command
│       └── list.go              # List generators command
├── generators/
│   ├── base/
│   │   └── generator.go         # Base generator utilities
│   ├── golang/
│   │   ├── generator.go         # Go generator implementation
│   │   ├── config.yaml          # Go generator config
│   │   ├── templates/
│   │   │   ├── core/
│   │   │   │   ├── models.go.tmpl
│   │   │   │   ├── types.go.tmpl
│   │   │   │   └── errors.go.tmpl
│   │   │   ├── client/
│   │   │   │   ├── base.go.tmpl
│   │   │   │   └── http.go.tmpl
│   │   │   └── server/
│   │   │       ├── base.go.tmpl
│   │   │       ├── handlers.go.tmpl
│   │   │       └── routes.go.tmpl
│   │   └── frameworks/
│   │       ├── chi/
│   │       │   └── templates/
│   │       └── echo/
│   │           └── templates/
│   ├── typescript/
│   │   ├── generator.go
│   │   ├── config.yaml
│   │   └── templates/
│   └── python/
│       ├── generator.go
│       ├── config.yaml
│       └── templates/
├── examples/
│   ├── petstore/
│   │   ├── openapi.yaml
│   │   └── .beaverspec.yaml
│   └── simple-api/
│       ├── openapi.yaml
│       └── .beaverspec.yaml
├── docs/
│   ├── ARCHITECTURE.md
│   ├── TECHNICAL_SPEC.md
│   └── USER_GUIDE.md
├── go.mod
├── go.sum
├── README.md
└── LICENSE
```

### Generated Output Structure

**Flat structure:**
```
generated/
├── models.go
├── client.go
├── server.go
└── types.go
```

**Organized structure:**
```
generated/
├── models/
│   ├── user.go
│   ├── product.go
│   └── order.go
├── client/
│   ├── client.go
│   ├── user_api.go
│   ├── product_api.go
│   └── order_api.go
└── server/
    ├── handlers.go
    ├── routes.go
    └── middleware.go
```

---

## CLI Specification

### Commands

```bash
beaverspec [command] [flags]
```

#### Generate Command

```bash
beaverspec generate [flags]

Flags:
  -s, --spec string           OpenAPI spec file (default: openapi.yaml)
  -g, --generator string      Generator to use (default: from config)
  -o, --output string         Output directory (default: ./generated)
  -c, --config string         Config file (default: .beaverspec.yaml)
      --models                Generate models (default: true)
      --client                Generate client (default: true)
      --server                Generate server (default: false)
      --validation            Generate validation (default: false)
      --framework string      Framework to use
      --dry-run               Show what would be generated
  -v, --verbose               Verbose output
  -h, --help                  Help for generate
```

#### Validate Command

```bash
beaverspec validate [flags]

Validates an OpenAPI specification file.

Flags:
  -s, --spec string      OpenAPI spec file (default: openapi.yaml)
      --strict           Enable strict validation
  -v, --verbose          Verbose output
  -h, --help             Help for validate
```

#### List Command

```bash
beaverspec list [flags]

Lists available generators and their capabilities.

Flags:
      --detailed         Show detailed information
  -h, --help             Help for list
```

#### Version Command

```bash
beaverspec version

Shows version information.
```

### Exit Codes

- `0`: Success
- `1`: General error
- `2`: Configuration error
- `3`: Validation error
- `4`: Generation error

### Environment Variables

```bash
OPENAPI_GEN_CONFIG    # Path to config file
OPENAPI_GEN_SPEC      # Path to spec file
OPENAPI_GEN_OUTPUT    # Output directory
OPENAPI_GEN_GENERATOR # Generator to use
```

---

## Implementation Guidelines

### Error Handling

1. **Collect validation errors** - Don't fail on first error
2. **Provide context** - Include file, line, column, JSON path
3. **Suggest fixes** - Include hints when possible
4. **Error codes** - Use codes for documentation lookup

### Testing Strategy

1. **Unit tests** - Test individual components
2. **Integration tests** - Test complete generation pipeline
3. **Snapshot tests** - Compare generated output
4. **Example validation** - Validate example specs regularly

### Template Best Practices

1. Keep templates simple and focused
2. Use helper functions for complex logic
3. Document template variables
4. Test templates with edge cases
5. Provide clear error messages

### Performance Considerations

1. Cache parsed specs
2. Lazy load templates
3. Parallel file writing
4. Incremental generation (future)

---

## Versioning and Compatibility

### Semantic Versioning

- **Major**: Breaking changes to interfaces or CLI
- **Minor**: New features, backward compatible
- **Patch**: Bug fixes

### Plugin Compatibility

```yaml
# Generator specifies required core version
requires:
  core: ">=1.0.0 <2.0.0"
```

### OpenAPI Support

- **v3.0.x**: Full support
- **v3.1.x**: Full support (goal)
- **v2.0** (Swagger): Not supported (use converter)

---

## Security Considerations

1. **Template injection** - Validate/sanitize template content
2. **Path traversal** - Validate output paths
3. **Resource limits** - Limit spec size and complexity
4. **Dependency scanning** - Regular security audits

---

## Migration Path

### From openapi-generator (Java)

Provide migration guide showing:
1. Configuration mapping
2. Template conversion guide
3. Customization migration
4. Command equivalents

### Between versions

1. Deprecation warnings
2. Migration scripts
3. Changelog with breaking changes
4. Compatibility mode (if needed)
