# OpenAPI Code Generator - Quick Start Guide

## What is OpenAPI Code Generator?

OpenAPI Code Generator is a Go-based tool that automatically generates client and server code from your OpenAPI 3.x API specifications. It's designed to be flexible, extensible, and easy to customize.

## Why Use This Generator?

- ✅ **Multiple Languages**: Generate code for Go, TypeScript, Python, and more
- ✅ **Framework Support**: Choose from popular frameworks (Chi, Echo, Gin for Go)
- ✅ **Customizable**: Override templates and configuration
- ✅ **Selective Generation**: Generate only what you need (client, server, models)
- ✅ **Type-Safe**: Leverages language type systems for safety
- ✅ **Go-Native**: Written in Go, fast and reliable

## Installation

### Binary Release (Recommended)

```bash
# Download latest release for your platform
curl -L https://github.com/yourorg/beaverspec/releases/latest/download/beaverspec-linux-amd64 -o beaverspec
chmod +x beaverspec
sudo mv beaverspec /usr/local/bin/
```

### Build from Source

```bash
git clone https://github.com/yourorg/beaverspec.git
cd beaverspec
go build -o beaverspec ./cmd/beaver
```

### Go Install

```bash
go install github.com/yourorg/beaverspec/cmd/beaverspec@latest
```

## Quick Start

### 1. Create an OpenAPI Specification

Create a file `openapi.yaml`:

```yaml
openapi: 3.0.0
info:
  title: My API
  version: 1.0.0
paths:
  /users:
    get:
      summary: List users
      responses:
        '200':
          description: Success
          content:
            application/json:
              schema:
                type: array
                items:
                  $ref: '#/components/schemas/User'
components:
  schemas:
    User:
      type: object
      required:
        - id
        - email
      properties:
        id:
          type: integer
          format: int64
        email:
          type: string
          format: email
        name:
          type: string
```

### 2. Generate Code

```bash
# Generate Go client
beaverspec generate --spec openapi.yaml --generator go --client

# Generate TypeScript client
beaverspec generate --spec openapi.yaml --generator typescript --client

# Generate both client and server
beaverspec generate --spec openapi.yaml --generator go --client --server
```

### 3. Use Generated Code

The generated code will be in `./generated/` directory:

```go
// Go example
package main

import (
    "context"
    "fmt"
    "generated/client"
)

func main() {
    client := client.NewClient("https://api.example.com")
    users, err := client.ListUsers(context.Background())
    if err != nil {
        panic(err)
    }
    fmt.Printf("Found %d users\n", len(users))
}
```

## Configuration

### Project Configuration File

Create `.beaverspec.yaml` in your project root:

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

# Generator-specific options
go:
  package: "myservice"
  use_generics: true

# Client configuration
client:
  framework: "http"
```

Now you can just run:

```bash
beaverspec generate
```

### User Configuration File

For user-wide defaults, create `~/.beaverspec/config.yaml`:

```yaml
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

## Common Use Cases

### Generate Client Only

```yaml
# .beaverspec.yaml
generator: go
spec: ./openapi.yaml

generate:
  models: true
  client: true
  server: false
```

### Generate Server Only

```yaml
# .beaverspec.yaml
generator: go
spec: ./openapi.yaml

generate:
  models: true
  client: false
  server: true

server:
  framework: "chi"
```

### Generate with Specific Framework

```yaml
# .beaverspec.yaml
generator: go
spec: ./openapi.yaml

generate:
  models: true
  client: true
  server: true

client:
  framework: "http"

server:
  framework: "echo"
```

### Models Only (for sharing types)

```yaml
# .beaverspec.yaml
generator: typescript
spec: ./openapi.yaml

generate:
  models: true
  client: false
  server: false
```

## Customizing Templates

### Override Specific Templates

```yaml
# .beaverspec.yaml
generator: go
spec: ./openapi.yaml

templates:
  overrides:
    "client/http": "./my-templates/custom-client.go.tmpl"
```

### Use Custom Template Directory

```yaml
# .beaverspec.yaml
generator: go
spec: ./openapi.yaml

templates:
  custom_dir: ./my-templates
```

Your custom templates will be used instead of built-in ones.

## Available Generators

### Go

**Features**:
- Multiple server frameworks (Chi, Echo, Gin, standard http)
- Multiple client frameworks (standard http, resty)
- Generics support (Go 1.18+)
- Validation generation

**Configuration**:
```yaml
generator: go

go:
  package: "api"
  use_generics: true
  json_tag_case: "camelCase"

client:
  framework: "http"  # or "resty"

server:
  framework: "chi"   # or "echo", "gin", "http"
```

### TypeScript

**Features**:
- ES modules or CommonJS
- Type-safe clients
- Validation support

**Configuration**:
```yaml
generator: typescript

typescript:
  target: "es2020"
  module_type: "esm"
  strict_null_checks: true
```

### Python

**Features**:
- Dataclasses or Pydantic models
- Type hints
- Async/await support

**Configuration**:
```yaml
generator: python

python:
  version: "3.10"
  use_dataclasses: true
  use_pydantic: false
```

## CLI Reference

### Generate Command

```bash
beaverspec generate [flags]

Flags:
  -s, --spec string          OpenAPI spec file
  -g, --generator string     Generator to use (go, typescript, python)
  -o, --output string        Output directory
  -c, --config string        Config file
      --models               Generate models
      --client               Generate client
      --server               Generate server
      --validation           Generate validation
      --framework string     Framework to use
      --dry-run             Show what would be generated
  -v, --verbose             Verbose output
```

### Validate Command

```bash
beaverspec validate [flags]

Flags:
  -s, --spec string     OpenAPI spec file
      --strict          Enable strict validation
  -v, --verbose         Verbose output
```

### List Generators

```bash
beaverspec list [--detailed]
```

Shows available generators and their capabilities.

### Version

```bash
beaverspec version
```

## Directory Structure Examples

### Flat Output

```
generated/
├── models.go
├── client.go
├── server.go
└── types.go
```

### Organized Output (Recommended)

```
generated/
├── models/
│   ├── user.go
│   ├── product.go
│   └── order.go
├── client/
│   ├── client.go
│   ├── user_api.go
│   └── product_api.go
└── server/
    ├── handlers.go
    └── routes.go
```

Configure with:

```yaml
output:
  structure: "organized"  # or "flat"
```

## Environment Variables

You can use environment variables for configuration:

```bash
export OPENAPI_GEN_SPEC=./openapi.yaml
export OPENAPI_GEN_GENERATOR=go
export OPENAPI_GEN_OUTPUT=./generated

beaverspec generate
```

## Troubleshooting

### Validation Errors

```bash
# Validate your spec first
beaverspec validate --spec openapi.yaml
```

### Verbose Output

```bash
# See detailed generation process
beaverspec generate --verbose
```

### Dry Run

```bash
# See what would be generated without writing files
beaverspec generate --dry-run
```

### Common Issues

**Issue**: "Generator not found"
```bash
# List available generators
beaverspec list
```

**Issue**: "Invalid OpenAPI spec"
```bash
# Validate your spec
beaverspec validate --spec openapi.yaml --strict
```

**Issue**: "Template rendering failed"
```bash
# Check template syntax and data
beaverspec generate --verbose
```

## Best Practices

### 1. Version Control

Add to `.gitignore`:
```
generated/
.beaverspec.cache
```

Commit your `.beaverspec.yaml` config file.

### 2. CI/CD Integration

```yaml
# .github/workflows/generate.yml
name: Generate API Code

on:
  push:
    paths:
      - 'openapi.yaml'

jobs:
  generate:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - name: Install beaverspec
        run: |
          curl -L https://github.com/yourorg/beaverspec/releases/latest/download/beaverspec-linux-amd64 -o beaverspec
          chmod +x beaverspec
      - name: Generate code
        run: ./beaverspec generate
      - name: Commit changes
        run: |
          git config --local user.email "action@github.com"
          git config --local user.name "GitHub Action"
          git add generated/
          git commit -m "Auto-generate API code" || true
          git push
```

### 3. Project Structure

Recommended structure:

```
my-project/
├── api/
│   └── openapi.yaml          # API specification
├── internal/
│   └── generated/            # Generated code
│       ├── models/
│       ├── client/
│       └── server/
├── cmd/
│   └── server/
│       └── main.go           # Uses generated server
├── .beaverspec.yaml         # Generator config
└── go.mod
```

### 4. Template Customization

Start with built-in templates, then customize:

1. Generate with built-in templates
2. Copy the template you want to customize
3. Modify it for your needs
4. Reference it in config

```bash
# Find built-in templates
beaverspec list --detailed

# Copy and customize
cp ~/.beaverspec/templates/go/client/http.go.tmpl ./templates/
# Edit ./templates/http.go.tmpl

# Use in config
templates:
  overrides:
    "client/http": "./templates/http.go.tmpl"
```

## Examples

Check the `examples/` directory for complete examples:

- `examples/petstore/` - Classic Petstore API
- `examples/todo-api/` - Simple TODO API
- `examples/ecommerce/` - E-commerce API with complex schemas

## Getting Help

- **Documentation**: Read the full docs in `docs/`
- **Issues**: Report bugs or request features on GitHub
- **Discussions**: Ask questions in GitHub Discussions
- **Contributing**: See `CONTRIBUTING.md` for contribution guidelines

## What's Next?

1. **Read the Architecture**: `docs/ARCHITECTURE.md` - Understand the design
2. **Explore Templates**: Look at built-in templates for examples
3. **Customize**: Create your own templates and configurations
4. **Contribute**: Add support for new languages or frameworks

## Quick Reference Card

```bash
# Basic generation
beaverspec generate --spec api.yaml --generator go --client

# With config file
beaverspec generate

# Validate spec
beaverspec validate --spec api.yaml

# List generators
beaverspec list

# Dry run
beaverspec generate --dry-run

# Verbose output
beaverspec generate --verbose
```

## Configuration Quick Reference

```yaml
# Minimal config
generator: go
spec: ./openapi.yaml

# Full config
generator: go
spec: ./openapi.yaml
output: ./internal/api

generate:
  models: true
  client: true
  server: true
  validation: false

go:
  package: "myapi"
  use_generics: true

client:
  framework: "http"

server:
  framework: "chi"

templates:
  custom_dir: ./templates
  overrides:
    "client/http": "./templates/custom.go.tmpl"

output:
  structure: "organized"
```

---

**Happy Code Generating! 🚀**
