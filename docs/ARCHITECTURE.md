# OpenAPI Code Generator - Architecture Documentation

## Table of Contents

1. [System Overview](#system-overview)
2. [Core Architectural Principles](#core-architectural-principles)
3. [Component Architecture](#component-architecture)
4. [Plugin System](#plugin-system)
5. [Template System](#template-system)
6. [Data Flow](#data-flow)
7. [Configuration Management](#configuration-management)
8. [Error Handling](#error-handling)

---

## System Overview

The OpenAPI Code Generator is a Go-based tool that generates multi-language client and server code from OpenAPI 3.x specifications. It uses a plugin-based architecture where language-specific logic is separated from the core generation engine.

### Key Design Goals

- **Flexibility**: Support multiple languages through plugins
- **Extensibility**: Easy to add new languages or frameworks
- **Customization**: Users can override templates and configuration
- **Maintainability**: Clear separation between core logic and language-specific code
- **Usability**: Simple CLI with sensible defaults

### High-Level Context Diagram

```
┌─────────────┐
│   User      │
│  (Developer)│
└──────┬──────┘
       │
       │ OpenAPI Spec + Config
       ▼
┌─────────────────────────────────────┐
│   OpenAPI Code Generator (CLI)      │
│                                     │
│  - Parse OpenAPI spec               │
│  - Load configuration               │
│  - Discover & load plugins          │
│  - Orchestrate generation           │
│  - Write output files               │
└─────────────────────────────────────┘
       │
       │ Generated Code Files
       ▼
┌─────────────┐
│  File       │
│  System     │
└─────────────┘
```

---

## Core Architectural Principles

### 1. Plugin-Based Architecture

The system separates **semantic concerns** (what code to generate) from **syntactic concerns** (how to write it).

- **Plugins**: Contain language-specific logic and type mappings
- **Templates**: Define code syntax and structure
- **Core**: Handles spec parsing, orchestration, and file management

### 2. Data-Driven Configuration

Language generators are primarily data-driven, using YAML configuration to define:
- Type mappings (OpenAPI types → language types)
- Naming conventions
- Nullable/optional strategies
- Framework variants

This allows 80% of languages to be added without writing Go code.

### 3. Separation of Concerns

```
┌─────────────────────────────────────────────────────────────┐
│                         CLI Layer                            │
│  - Command parsing                                           │
│  - User interaction                                          │
│  - Error reporting                                           │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│                    Orchestration Layer                       │
│                                                              │
│  ┌─────────────────┐    ┌──────────────┐   ┌─────────────┐│
│  │  Spec Parser    │───▶│  Validator   │──▶│  Generator  ││
│  │  (OpenAPI 3.x)  │    │  (Schema +   │   │  Coordinator││
│  └─────────────────┘    │   Config)    │   └─────────────┘│
│                         └──────────────┘                    │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│                    Plugin System                             │
│                                                              │
│  ┌──────────────────────────────────────────────────────┐  │
│  │            Plugin Registry                            │  │
│  │  - Discover plugins                                   │  │
│  │  - Validate plugin compatibility                      │  │
│  │  - Lifecycle management                               │  │
│  └────────────┬──────────────────────────────────────────┘  │
│               │                                              │
│               ▼                                              │
│  ┌──────────────────────────────────────────────────────┐  │
│  │         Language Generator Interface                  │  │
│  │  (Contract that all generators must implement)        │  │
│  └────────────┬──────────────────────────────────────────┘  │
└───────────────┼──────────────────────────────────────────────┘
                │
                ├────────────────┬────────────────┬──────────────
                ▼                ▼                ▼
    ┌───────────────────┐ ┌──────────────┐ ┌─────────────────┐
    │  Go Generator     │ │ TypeScript   │ │  Python         │
    │  Plugin           │ │ Generator    │ │  Generator      │
    │                   │ │ Plugin       │ │  Plugin         │
    └───────────────────┘ └──────────────┘ └─────────────────┘
```

### 4. Template Composition

Templates are organized by use case and can be composed:

```
generators/
  <language>/
    config.yaml
    templates/
      core/
        models.tmpl          # Data models (always needed)
        types.tmpl           # Type definitions
        errors.tmpl          # Error types
      client/
        base.tmpl            # Client interface/base
        http.tmpl            # HTTP client implementation
      server/
        base.tmpl            # Server interface/base
        handlers.tmpl        # Handler implementations
        routes.tmpl          # Route definitions
      validation/            # Optional
        validators.tmpl
    frameworks/              # Optional framework-specific overrides
      <framework-name>/
        templates/           # Override specific templates
```

---

## Component Architecture

### Core Components

#### 1. CLI Layer

**Responsibilities:**
- Parse command-line arguments
- Display help and usage information
- Coordinate high-level workflow
- Format and display output/errors

**Key Commands:**
```bash
beaverspec generate        # Generate code
beaverspec validate        # Validate spec
beaverspec list            # List available generators
beaverspec version         # Show version
```

#### 2. Spec Parser & Normalizer

**Responsibilities:**
- Parse OpenAPI YAML/JSON files
- Validate against OpenAPI 3.x schema
- Resolve `$ref` references
- Normalize schemas (flatten, resolve allOf/oneOf)
- Build internal enriched representation

**Input:** Raw OpenAPI specification file
**Output:** Enriched Spec Model (internal representation)

#### 3. Configuration Manager

**Responsibilities:**
- Load configuration from multiple sources
- Merge configurations with precedence rules
- Validate configuration values
- Provide configuration to generators

**Configuration Sources (highest to lowest precedence):**
1. CLI flags
2. Environment variables
3. Project config file (`.beaverspec.yaml`)
4. User config file (`~/.beaverspec/config.yaml`)
5. Generator defaults

#### 4. Plugin Registry

**Responsibilities:**
- Discover available generators
- Register generators
- Validate plugin compatibility
- Provide generator instances

**Plugin Loading Strategy:**
- **Phase 1 (MVP)**: Compile-time registration (all generators built-in)
- **Phase 2 (Future)**: Runtime loading from filesystem
- **Phase 3 (Future)**: Download from registry

#### 5. Generation Orchestrator

**Responsibilities:**
- Coordinate the generation process
- Select appropriate generator
- Determine which templates to use
- Invoke generator with configuration
- Write output files

#### 6. Template Engine

**Responsibilities:**
- Load and parse templates
- Provide template helper functions
- Render templates with data
- Handle template composition/includes

**Template Engine:** Go's `text/template` package

---

## Plugin System

### Generator Interface

Every language generator must implement this interface:

```go
type Generator interface {
    // Metadata
    Metadata() GeneratorMetadata
    
    // Lifecycle
    Validate(spec *Spec, config *Config) []ValidationError
    Generate(spec *Spec, config *Config) (*GenerationResult, error)
    
    // Capabilities (optional overrides)
    SupportedFeatures() []Feature
}
```

### Generator Metadata

```go
type GeneratorMetadata struct {
    Name         string
    Language     string
    Version      string
    Description  string
    Author       string
    
    // What this generator can produce
    Capabilities Capabilities
}

type Capabilities struct {
    SupportsClient     bool
    SupportsServer     bool
    SupportsValidation bool
    SupportsAsync      bool
    SupportsStreaming  bool
    
    // Framework support
    ClientFrameworks []string
    ServerFrameworks []string
}
```

### Plugin Configuration Schema

Each generator has a `config.yaml` that defines its behavior:

```yaml
name: "go"
version: "1.0.0"
language: "go"
description: "Go code generator with support for multiple frameworks"

# Type mappings from OpenAPI to Go
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
  
# How to handle nullable/optional
nullable_strategy: "pointer"        # Use *Type for nullable
optional_strategy: "omitempty"      # Use omitempty tag

# Templates
templates:
  required:
    - "core/models"
    - "core/types"
  optional:
    client:
      - "client/base"
      - "client/http"
    server:
      - "server/base"
      - "server/handlers"

# Supported frameworks
frameworks:
  client:
    - "http"
    - "resty"
  server:
    - "chi"
    - "echo"
    - "gin"
    - "http"

# Defaults
defaults:
  generate:
    models: true
    client: true
    server: false
  client:
    framework: "http"
  server:
    framework: "http"
```

### Data-Driven vs Code-Driven Plugins

**80% of generators**: Pure data-driven (config + templates)
**20% of generators**: Require custom Go code for complex logic

This allows most languages to be added without writing code.

---

## Template System

### Template Granularity

Templates are **use-specific** (separate templates for different artifacts):

**Benefits:**
- Generate only what's needed (client-only, server-only, both)
- Easier to maintain smaller templates
- Can mix and match components
- Users can override specific parts

### Template Selection Logic

```
For each template needed:
  1. Check: frameworks/<framework>/templates/<template>
  2. Fallback: templates/<use-case>/<template>
  3. Error if not found
```

### Template Helper Functions

Available in all templates:

**String manipulation:**
- `toCamelCase`, `toPascalCase`, `toSnakeCase`, `toKebabCase`
- `toLower`, `toUpper`, `capitalize`
- `pluralize`, `singularize`

**Type checking:**
- `isArray`, `isObject`, `isString`, `isPrimitive`
- `isNullable`, `isRequired`, `isEnum`

**Code generation:**
- `indent`, `quote`, `escape`
- `joinPath`, `baseName`, `dirName`

**Language-specific (injected by generator):**
- `goType`, `tsType`, `pyType`
- `goPackage`, `tsImport`, `pyImport`

**OpenAPI-specific:**
- `isRef`, `resolveRef`
- `hasDiscriminator`, `discriminatorValue`
- `allOfSchemas`, `oneOfSchemas`

### Template Composition

Templates can include/reference each other:

```go
{{/* client/http.tmpl */}}
package {{.PackageName}}

{{template "imports" .}}

{{template "client-interface" .}}

{{template "http-methods" .}}

{{/* Reuse error handling from core */}}
{{template "core/errors" .}}
```

---

## Data Flow

### Complete Generation Pipeline

```
┌─────────────────┐
│  OpenAPI Spec   │
│  (YAML/JSON)    │
└────────┬────────┘
         │
         ▼
┌─────────────────────────────────────┐
│  Parser & Validator                  │
│  - Parse YAML/JSON                   │
│  - Validate against OpenAPI schema   │
│  - Resolve $refs                     │
└────────┬────────────────────────────┘
         │
         ▼
┌─────────────────────────────────────┐
│  Spec Normalizer                     │
│  - Flatten nested schemas            │
│  - Resolve allOf/oneOf/anyOf         │
│  - Extract reusable components       │
│  - Build dependency graph            │
└────────┬────────────────────────────┘
         │
         ▼
┌─────────────────────────────────────┐
│  Enriched Spec Model                 │
│  (Internal representation)           │
│  - Easier to traverse                │
│  - Pre-computed relationships        │
│  - Metadata attached                 │
└────────┬────────────────────────────┘
         │
         ▼
┌─────────────────────────────────────┐
│  Generator-Specific Transform        │
│  - Apply language type mappings      │
│  - Apply naming conventions          │
│  - Group related schemas             │
│  - Determine file organization       │
└────────┬────────────────────────────┘
         │
         ▼
┌─────────────────────────────────────┐
│  Template Data Models                │
│  (Ready for template rendering)      │
└────────┬────────────────────────────┘
         │
         ▼
┌─────────────────────────────────────┐
│  Template Rendering                  │
└────────┬────────────────────────────┘
         │
         ▼
┌─────────────────────────────────────┐
│  Post-Processing                     │
│  - Format code (gofmt, prettier)     │
│  - Add file headers                  │
│  - Organize imports                  │
└────────┬────────────────────────────┘
         │
         ▼
┌─────────────────┐
│  Output Files   │
└─────────────────┘
```

### Pipeline Responsibility Split

**Shared Core:**
- Parser & Validator
- Spec Normalizer
- Enriched Spec Model

**Generator-Specific:**
- Type mapping
- Naming conventions
- Template data preparation
- Rendering
- Post-processing

---

## Configuration Management

### Configuration Hierarchy

Precedence from highest to lowest:

```
1. CLI Flags
   ↓
2. Environment Variables
   ↓
3. Project Config File (.beaverspec.yaml)
   ↓
4. User Config File (~/.beaverspec/config.yaml)
   ↓
5. Generator Defaults (from generator's config.yaml)
```

### User Configuration File

Location: `~/.beaverspec/config.yaml`

```yaml
# User defaults applied to all projects
defaults:
  output: ./generated
  generate:
    validation: true

generators:
  go:
    package: "api"
  typescript:
    target: "es2020"
```

### Project Configuration File

Location: `.beaverspec.yaml` in project root

```yaml
# Project-specific configuration
generator: go
spec: ./openapi.yaml

output: ./internal/api

generate:
  models: true
  client: true
  server: false
  validation: false

go:
  package: "myservice"
  use_generics: true
  
server:
  framework: "chi"
```

### CLI Overrides

```bash
beaverspec generate \
  --spec ./openapi.yaml \
  --generator typescript \
  --output ./custom-output \
  --no-validation
```

### Generation Control

Users can control what gets generated:

```yaml
generate:
  models: true           # Always recommended
  client: true           # Generate client?
  server: true           # Generate server?
  validation: false      # Optional validators?

client:
  framework: "http"      # "http", "resty", "custom"
  
server:
  framework: "chi"       # "chi", "echo", "gin", "http"
  middleware: true       # Generate middleware?
```

---

## Error Handling

### Error Categories

```go
type ErrorSeverity int

const (
    ErrorSeverityWarning ErrorSeverity = iota
    ErrorSeverityError
    ErrorSeverityFatal
)
```

### Error Types

**Warnings:**
- Unsupported OpenAPI features
- Deprecated patterns
- Best practice violations
- **Behavior:** Collect and display at end, continue generation

**Errors:**
- Invalid spec
- Unsupported required features
- Template rendering failures
- **Behavior:** Collect all, stop generation, display all errors

**Fatal:**
- Config file not found
- Generator not found
- File write permissions
- **Behavior:** Abort immediately

### Error Reporting

```go
type GenerationError struct {
    Severity ErrorSeverity
    Phase    string        // "parsing", "validation", "generation"
    Location *Location     // Where in spec/template
    Message  string
    Hint     string        // Suggestion to fix
    Code     string        // Error code for docs lookup
}

type Location struct {
    File   string
    Line   int
    Column int
    Path   string  // JSON path: "#/paths/users/get"
}
```

### User Experience

```bash
$ beaverspec generate

⚠ Warnings (2):
  • Line 45: Deprecated pattern 'type: file'
    Hint: Use 'type: string, format: binary' instead
  • Path /users/{id}: Missing description

✗ Errors (1):
  • #/components/schemas/User: Discriminator 'type' not found
    Hint: Add 'type' to required properties or remove discriminator
    See: https://docs.beaverspec.io/errors/E1001

Generation failed with 1 error, 2 warnings.
```

---

## Key Architectural Decisions

### 1. Plugin Loading Strategy

**Phase 1 (MVP):** Compile-time plugins
- All generators built into binary
- Templates embedded using `//go:embed`
- Zero external dependencies
- Simple and reliable

**Phase 2 (Future):** Runtime plugins
- Load from `~/.beaverspec/plugins/`
- Load from `./generators/` (project-local)
- Version compatibility checking

**Phase 3 (Future):** Plugin registry
- Download from central registry
- Semantic versioning
- Dependency resolution

### 2. Template vs Plugin Boundary

**Templates handle:** Syntax and code structure
**Plugins handle:** Semantics and language-specific logic

This separation allows maximum flexibility while keeping templates simple.

### 3. Configuration Philosophy

**Data-driven by default:** Most configuration in YAML
**Code when necessary:** Complex logic in Go
**User override everything:** Through config or custom templates

### 4. Output Organization

**Configurable structure:**

```yaml
output:
  structure: "organized"  # "flat" | "organized" | "custom"
  
  # For "custom"
  patterns:
    models: "{{.OutputDir}}/models/{{.Name | toSnakeCase}}.go"
    client: "{{.OutputDir}}/client/{{.Tag | toSnakeCase}}_api.go"
```

---

## Future Enhancements

### Short Term
- Add more built-in language generators
- Improve error messages and validation
- Add watch mode for development
- Better documentation generation

### Medium Term
- Runtime plugin loading
- Plugin marketplace/registry
- IDE integration (Language Server Protocol)
- Interactive mode for configuration

### Long Term
- Code-first approach (generate OpenAPI from code)
- Migration tools between API versions
- Contract testing integration
- Mock server generation

---

## Conclusion

This architecture provides:

1. **Flexibility** through plugins and templates
2. **Simplicity** through data-driven configuration
3. **Extensibility** through clear interfaces
4. **Usability** through sensible defaults
5. **Maintainability** through separation of concerns

The design allows the project to start simple (built-in generators with embedded templates) while providing a clear path to advanced features (runtime plugins, custom generators, framework variants).
