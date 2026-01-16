# BeaverSpec 🦫

**A modern OpenAPI code generator with first-class template support**

BeaverSpec is a flexible, plugin-based OpenAPI 3.x code generator written in Go. It generates production-ready client and server code from OpenAPI specifications, with better template customization than traditional tools.

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![OpenAPI](https://img.shields.io/badge/OpenAPI-3.x-6BA539?style=flat&logo=openapi-initiative)](https://www.openapis.org/)

## Why BeaverSpec?

- **🎯 Template-First**: 80% of language support through configuration + templates
- **🔌 Plugin Architecture**: Clean separation between syntax (templates) and semantics (plugins)
- **⚙️ Highly Configurable**: CLI flags → env vars → config files with intelligent merging
- **🚀 Modern Go**: Built with Go 1.21+, leveraging the latest language features
- **📦 Zero Dependencies**: Single binary with embedded templates
- **🎨 Full OpenAPI 3.x Support**: All schema features including oneOf, allOf, discriminators

## Current Status

**Phase 1**: ✅ Complete - Core architecture and Go models generator
**Phase 2**: 🚧 In Progress - Client/server generation and shared tooling

### Implemented Features

- ✅ Complete OpenAPI 3.x spec parsing and normalization
- ✅ Go models generation with comprehensive type support:
  - Basic types (string, integer, number, boolean, object)
  - Arrays and maps (via additionalProperties)
  - Nested objects and $ref resolution
  - Enums with constants
  - allOf composition (struct embedding)
  - oneOf/anyOf polymorphism with discriminators
  - Field metadata (default, readOnly, writeOnly, deprecated)
- ✅ Basic HTTP client generation
- ✅ Auto-generated tests (compilation and serialization)
- ✅ Dynamic import management

## Quick Start

### Installation

```bash
# Clone the repository
git clone https://github.com/imerljak/beaverspec.git
cd beaverspec

# Build the CLI
go build -o beaver ./cmd/beaver

# Or install directly
go install github.com/imerljak/beaverspec/cmd/beaver@latest
```

### Generate Code

```bash
# Generate Go models and client from OpenAPI spec
beaver -spec ./examples/petstore.yaml -output ./generated

# Output:
# ✓ Successfully generated 3 file(s) in generated
#   - models.go
#   - models_test.go
#   - client.go
```

### Use Generated Code

```go
package main

import (
    "context"
    "fmt"
    "github.com/yourproject/generated/models"
)

func main() {
    // Create a new client
    client := models.NewClient("https://petstore.swagger.io/v2")

    // List pets
    limit := int32(10)
    pets, err := client.ListPets(context.Background(), &limit)
    if err != nil {
        panic(err)
    }

    fmt.Printf("Found %d pets\n", len(*pets))
}
```

## Generated Code Features

### Models
- Type-safe structs from OpenAPI schemas
- JSON serialization tags
- Pointer types for optional fields
- Embedded structs for allOf composition
- Interface types for oneOf/anyOf polymorphism
- Enum constants with validation

### Client
- Interface-based design for easy mocking
- Context support for cancellation
- Typed request/response models
- Automatic JSON serialization
- Configurable HTTP client
- Custom headers support

### Tests
- Compilation verification tests
- JSON serialization round-trip tests
- Generated automatically with models

## Project Structure

```
beaverspec/
├── cmd/beaver/              # CLI entry point
├── pkg/
│   ├── core/                # Core interfaces and types
│   ├── parser/              # OpenAPI spec parsing
│   ├── plugin/              # Generator registry
│   └── template/            # Template engine
├── generators/
│   └── golang/              # Go generator
│       ├── generator.go     # Generator implementation
│       ├── config.yaml      # Type mappings and config
│       └── templates/       # Go templates
│           ├── core/        # Models templates
│           └── client/      # Client templates
└── examples/                # Example OpenAPI specs
```

## Examples

### Simple Pet Store
```bash
beaver -spec ./examples/petstore.yaml -output ./gen
```

Generates models and client for the classic Swagger Petstore API.

## Roadmap

See [ROADMAP.md](./ROADMAP.md) for the complete development plan.

**Phase 2** (In Progress):
- Server generation (handlers, routes, middleware)
- Multiple framework support (Chi, Echo, Gin)
- Authentication/authorization
- Request/response validation
- Shared cross-language tooling

**Phase 3** (Planned):
- TypeScript/JavaScript generator
- Python generator

**Phase 4** (Future):
- Additional languages (Java, Rust, etc.)
- Advanced features (streaming, webhooks, etc.)

## Development

### Prerequisites
- Go 1.21 or higher
- Basic understanding of OpenAPI 3.x

### Build
```bash
go build -o beaver ./cmd/beaver
```

### Run Tests
```bash
go test ./...
```

### Run with Example
```bash
go run ./cmd/beaver -spec examples/petstore.yaml -output generated
```

## Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](./CONTRIBUTING.md) for guidelines.

## Philosophy

BeaverSpec follows these principles:

1. **Data over Code**: Configuration and templates should handle most language differences
2. **Plugin-based**: Clear separation between OpenAPI semantics and language syntax
3. **Template Granularity**: Use-specific templates (not monolithic) for selective generation
4. **Sensible Defaults**: Work out-of-the-box with intelligent defaults
5. **Production Ready**: Generated code should be production-quality, not just prototypes

## Why "BeaverSpec"?

Beavers are nature's engineers, known for building robust, well-structured dams. Like a beaver, this tool helps you build solid, reliable API implementations from specifications. Plus, beavers are adorable 🦫

## License

MIT License - see [LICENSE](./LICENSE) for details

## Acknowledgments

- Built with [kin-openapi](https://github.com/getkin/kin-openapi) for OpenAPI parsing
- Inspired by [openapi-generator](https://github.com/OpenAPITools/openapi-generator) but modernized for the Go ecosystem

---

**Built with ❤️ and Go**
