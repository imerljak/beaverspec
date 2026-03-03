# BeaverSpec 🦫

**A modern OpenAPI code generator for Go**

BeaverSpec generates production-ready Go server and client code from OpenAPI 3.x specifications — no JVM required, no YAML-hell, just a single Go binary.

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![OpenAPI](https://img.shields.io/badge/OpenAPI-3.x-6BA539?style=flat&logo=openapi-initiative)](https://www.openapis.org/)

## Why BeaverSpec?

- **No JVM** — single Go binary, install in one command
- **Go-idiomatic output** — structs, interfaces, `error` returns, `context.Context`
- **Four frameworks** — `net/http`, Chi, Echo, Gin — switch with one flag
- **Template customization** — override any template with your own `text/template` file
- **Remote specs** — point at a URL, no need to download specs manually
- **Selective generation** — models only, server only, or everything

## Status

**Phase 1**: ✅ Complete — Core architecture, OpenAPI parsing, Go models generation
**Phase 2**: ✅ Complete — Server/client generation, framework support, security, validation, config system, integration polish

## Features

- ✅ Complete OpenAPI 3.x parsing (local files and remote URLs)
- ✅ Go models with JSON tags and `Validate()` methods
- ✅ allOf (struct embedding), oneOf/anyOf (discriminator, union types), enums, maps, arrays
- ✅ HTTP client with typed methods per endpoint
- ✅ Server handler interface + mock implementations
- ✅ Four framework targets: `net-http`, `chi`, `echo`, `gin`
- ✅ Request/response validation with typed error models
- ✅ Content negotiation (multiple `produces`/`consumes` types)
- ✅ Security scheme extraction (Bearer, API Key, Basic, OAuth2)
- ✅ Middleware generation (logging, CORS, rate limiter — stdlib only)
- ✅ Graceful shutdown example `main.go` per framework
- ✅ External `$ref` resolution (local files and URLs)
- ✅ Template customization and override system
- ✅ `beaver.yaml` config file with full schema support

## Quick Start

```bash
# Install
go install github.com/imerljak/beaverspec/cmd/beaver@latest

# Generate from a local spec
beaver -spec openapi.yaml -output ./generated -module github.com/example/myapi

# Generate from a remote spec
beaver -spec https://petstore3.swagger.io/api/v3/openapi.json \
  -output ./generated -module github.com/example/petstore

# Choose a framework
beaver -spec openapi.yaml -output ./generated -module github.com/example/myapi -framework chi

# Preview without writing files
beaver -spec openapi.yaml -output ./generated -module github.com/example/myapi -dry-run
```

## Generated Files

For a spec with endpoints and schemas, `beaver` generates up to 9 files:

```
generated/
├── models/
│   ├── models.go          # Typed structs, json tags, Validate() methods
│   └── models_test.go     # Auto-generated validation tests
├── client/
│   └── client.go          # Typed HTTP client
├── server/
│   ├── interface.go       # Handler interface — implement this
│   ├── mocks.go           # Mock implementations for quick start
│   ├── handlers.go        # HTTP adapter (routes to interface)
│   ├── routes.go          # Route registration per tag group
│   └── middleware.go      # Logging, CORS, rate limiter (stdlib)
└── cmd/server/
    └── main.go            # Runnable example entrypoint (edit freely)
```

## Config File

Create `beaver.yaml` next to your spec to avoid repeating flags:

```yaml
spec: ./openapi.yaml
output: ./generated
module: github.com/example/myapi
framework: chi

generate:
  models: true
  client: true
  server: true

exclude:
  tags: [internal]

types:
  overrides:
    - format: uuid
      go: "github.com/google/uuid".UUID
```

Then just run `beaver` with no flags.

## CLI Reference

| Flag | Default | Description |
|------|---------|-------------|
| `-spec` | — | OpenAPI spec path or URL (required) |
| `-output` | `.` | Output directory |
| `-module` | `github.com/example/project` | Go module path |
| `-framework` | `net-http` | Server framework: `net-http`, `chi`, `echo`, `gin` |
| `-config` | auto | Config file path (auto-discovers `beaver.yaml` next to spec) |
| `-package` | `models` | Package name for generated models |
| `-generator` | `go` | Generator to use |
| `-dry-run` | false | Print files that would be generated, don't write |
| `-debug` | false | Enable debug logging |
| `-version` | — | Print version and exit |

## Project Structure

```
beaverspec/
├── cmd/beaver/          # CLI entry point
├── pkg/
│   ├── core/            # Core types and interfaces
│   ├── parser/          # OpenAPI parsing and normalization
│   ├── config/          # Config file loading and merging
│   ├── codegen/         # Validation codegen helpers
│   └── template/        # Template engine
├── generators/
│   └── golang/          # Go generator + templates
│       └── templates/   # Per-framework Go templates
└── examples/            # Example OpenAPI specs
```

## Documentation

- [Quick Start](docs/QUICK_START.md) — 5-minute walkthrough
- [API Reference](docs/API_REFERENCE.md) — all flags, config schema, generated file reference
- [Migration Guide](docs/MIGRATION_GUIDE.md) — migrating from `openapi-generator-cli`

## Development

```bash
go build -o beaver ./cmd/beaver
go test ./...
go run ./cmd/beaver -spec examples/petstore.yaml -output /tmp/gen -module github.com/example/petstore -dry-run
```

## License

MIT License — see [LICENSE](./LICENSE) for details.

## Acknowledgments

- Built with [kin-openapi](https://github.com/getkin/kin-openapi) for OpenAPI parsing
- Inspired by [openapi-generator](https://github.com/OpenAPITools/openapi-generator), redesigned for the Go ecosystem

---

*Beavers are nature's engineers, known for building robust, well-structured dams. 🦫*
