# Migration Guide: openapi-generator-cli to BeaverSpec

This guide is for Go teams currently using `openapi-generator-cli` who want to switch to BeaverSpec (`beaver`).

---

## Why Migrate

- **No JVM required.** `openapi-generator-cli` requires a Java runtime. BeaverSpec is a single Go binary with zero runtime dependencies.
- **Native Go tooling.** Install with `go install`, run in any Go CI pipeline without container tricks or version pinning Java.
- **Go-idiomatic output.** Generated code follows Go conventions: interfaces, typed errors, standard `net/http` patterns, and `go vet`-clean output.
- **Template customization in Go.** Override or replace any template using Go `text/template`. No Mustache, no Java classpath manipulation.
- **Simple YAML config.** One `beaver.yaml` next to your spec. No generator-specific flags spread across shell scripts.
- **MIT licensed.** No contributor license agreements, no premium plugins, no hosted service required.

---

## Feature Comparison

| Feature | openapi-generator-cli | BeaverSpec |
|---|---|---|
| Installation | `npm install` or JAR download | `go install github.com/imerljak/beaverspec/cmd/beaver@latest` |
| Runtime dependency | Java 11+ | None |
| Config file | `config.yaml` (generator-specific keys) | `beaver.yaml` (single schema) |
| Go frameworks | `go`, `go-gin-server`, `go-echo-server`, etc. | `net-http`, `chi`, `echo`, `gin` |
| Languages | 50+ | Go only |
| Template engine | Mustache | Go `text/template` |
| Validation | External library (`oapi-codegen`-style) | Built-in `Validate()` on generated structs |
| Discriminator (`oneOf`) | Partial | Supported |
| Remote specs | Yes | Yes (`-spec https://...`) |
| Dry run | No built-in flag | `-dry-run` flag |
| CI integration | Requires JVM in image | Single binary, no extra deps |

---

## Command Mapping

### Basic generation

```bash
# openapi-generator-cli
openapi-generator-cli generate \
  -g go \
  -i spec.yaml \
  -o ./generated \
  --additional-properties=packageName=myapi

# beaver
beaver -spec spec.yaml -output ./generated -module github.com/example/myapi
```

### Framework-specific server

```bash
# openapi-generator-cli
openapi-generator-cli generate \
  -g go-gin-server \
  -i spec.yaml \
  -o ./generated

# beaver
beaver -spec spec.yaml -output ./generated -module github.com/example/myapi -framework gin
```

### Dry run

```bash
# openapi-generator-cli
# (no built-in dry run; workaround: generate to a temp dir)

# beaver
beaver -spec spec.yaml -output ./generated -module github.com/example/myapi -dry-run
```

### Remote spec

```bash
# openapi-generator-cli
openapi-generator-cli generate -g go -i https://petstore3.swagger.io/api/v3/openapi.json -o ./generated

# beaver
beaver -spec https://petstore3.swagger.io/api/v3/openapi.json -output ./generated -module github.com/example/petstore
```

### Using a config file

```bash
# openapi-generator-cli
openapi-generator-cli generate -g go -i spec.yaml -c config.yaml -o ./generated

# beaver
beaver -config beaver.yaml -output ./generated
```

---

## Config File Migration

A typical `openapi-generator-cli` config file:

```yaml
# openapi-generator config.yaml
packageName: myapi
packageVersion: 1.0.0
generateInterfaces: true
enumClassPrefix: true
structPrefix: true
```

The equivalent `beaver.yaml`:

```yaml
# beaver.yaml
spec: ./openapi.yaml
output: ./generated
module: github.com/example/myapi
package: models
framework: net-http  # net-http | chi | echo | gin

generate:
  models: true
  client: true
  server: true
  validation: true

exclude:
  models: [InternalError]
  tags: [internal]

types:
  overrides:
    - format: uuid
      go: "github.com/google/uuid".UUID

templates:
  customDir: ./my-templates
  overrides:
    models.go.tmpl: ./custom-models.go.tmpl
```

All fields in `beaver.yaml` are optional. The `spec` field can point to a local file or a remote URL.

---

## Generated Code Differences

### File layout

`openapi-generator-cli` produces one file per model and one file per API operation by default. BeaverSpec produces a smaller, grouped set of files:

```
# openapi-generator-cli (go generator)
model_pet.go
model_error.go
api_pets.go
api_store.go
...

# BeaverSpec
models/models.go          - all structs in one file, with Validate() methods
models/models_test.go     - generated validation tests
client/client.go          - typed HTTP client
server/interface.go       - handler interface (one method per operation)
server/mocks.go           - mock implementations for testing
server/handlers.go        - HTTP adapter (decodes request, calls interface, encodes response)
server/routes.go          - route registration
server/middleware.go      - logging, CORS, rate limiting (net-http only)
cmd/server/main.go        - runnable example entrypoint
```

### Handler pattern

`openapi-generator-cli` generates function-based handlers or a loosely typed `ApiService` interface. BeaverSpec generates a strict interface where every operation is a method:

```go
// BeaverSpec server/interface.go (example)
type PetsHandler interface {
    ListPets(ctx context.Context, req ListPetsRequest) (ListPetsResponse, error)
    CreatePet(ctx context.Context, req CreatePetRequest) (CreatePetResponse, error)
}
```

You implement the interface and pass it to the generated router. The HTTP layer (decoding, encoding, content negotiation, error mapping) is handled by `server/handlers.go`, which you do not edit.

### Entrypoint

The generated `cmd/server/main.go` is a working, compilable example that wires your implementation to the HTTP server. It is intended to be copied into your own `cmd/` directory and edited freely. It is not regenerated on subsequent runs if it already exists.

### Validation

BeaverSpec generates a `Validate() error` method on each model struct derived from the OpenAPI schema constraints (`minLength`, `minimum`, `required`, `pattern`, etc.). There is no separate validation library to import.

---

## Known Limitations

The following features supported by `openapi-generator-cli` are not yet available in BeaverSpec:

- **XML bodies.** Only `application/json` request and response bodies are supported. `application/xml` is ignored.
- **Form-urlencoded and multipart.** `application/x-www-form-urlencoded` and `multipart/form-data` are not supported.
- **Non-Go output.** BeaverSpec only generates Go code. Use `openapi-generator-cli` if you need TypeScript, Java, Python, or other languages.
- **Webhooks and callbacks.** OpenAPI 3.1 webhooks and 3.0 callbacks are parsed but not included in generated output.
- **Links.** OpenAPI `links` objects are ignored.
- **Security enforcement.** Security schemes are extracted and surfaced in the generated interface, but no authentication middleware is generated. You wire your own auth middleware.
