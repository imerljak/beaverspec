# BeaverSpec API Reference

## CLI Flags

All flags are passed directly to the `beaver` binary. There are no subcommands.

```bash
beaver [flags]
```

| Flag | Type | Default | Description |
|---|---|---|---|
| `-spec` | string | — | Path to an OpenAPI spec file, or an HTTPS URL |
| `-output` | string | — | Directory where generated files are written |
| `-module` | string | — | Go module path (e.g. `github.com/example/myapi`); must match `go.mod` |
| `-package` | string | — | Package name for generated model code |
| `-framework` | string | `net-http` | Server framework: `net-http`, `chi`, `echo`, or `gin` |
| `-generator` | string | `go` | Which generator to use; currently only `go` is available |
| `-config` | string | `beaver.yaml` in the spec directory | Path to a config file |
| `-dry-run` | bool | `false` | Print files that would be generated without writing them |
| `-debug` | bool | `false` | Enable debug-level logging |
| `-version` | bool | `false` | Print the version string and exit |

---

## `beaver.yaml` Config File

The config file is looked up at `beaver.yaml` in the spec's directory unless overridden with `-config`. All fields are optional and correspond to their CLI flag equivalents where both exist. CLI flags take precedence over config file values.

### Top-level fields

| Field | Type | Default | Description |
|---|---|---|---|
| `spec` | string | — | Path or URL to the OpenAPI spec file |
| `output` | string | — | Output directory for generated files |
| `module` | string | — | Go module path; must match `go.mod` |
| `package` | string | — | Package name for generated model code |
| `framework` | string | `net-http` | Server framework: `net-http`, `chi`, `echo`, `gin` |

### `generate` block

Controls which artifact groups are emitted.

| Field | Type | Default | Description |
|---|---|---|---|
| `generate.models` | bool | `true` | Generate `models/models.go` and `models/models_test.go` |
| `generate.client` | bool | `true` | Generate `client/client.go` |
| `generate.server` | bool | `true` | Generate all files under `server/` and `cmd/server/` |
| `generate.validation` | bool | `true` | Include `Validate()` methods on model structs |

### `exclude` block

Suppress generation of specific named items.

| Field | Type | Description |
|---|---|---|
| `exclude.models` | `[]string` | Schema names to skip (e.g. `[InternalError]`) |
| `exclude.tags` | `[]string` | OpenAPI tag names whose operations are omitted entirely |

### `types.overrides` block

A list of format-to-Go-type mappings that override the defaults.

```yaml
types:
  overrides:
    - format: email
      go: string
    - format: uuid
      go: "github.com/google/uuid".UUID
```

| Field | Type | Description |
|---|---|---|
| `format` | string | OpenAPI `format` value to match (e.g. `uuid`, `email`, `date-time`) |
| `go` | string | Replacement Go type; use `"import/path".TypeName` for external types |

### `templates` block

| Field | Type | Description |
|---|---|---|
| `templates.customDir` | string | Directory containing custom templates that replace built-in ones wholesale |
| `templates.overrides` | map[string]string | Per-template overrides: key is the built-in template filename, value is a path to the replacement |

Example:

```yaml
templates:
  customDir: ./my-templates
  overrides:
    models.go.tmpl: ./custom-models.go.tmpl
```

### `options` block

Fine-grained toggles for optional generated files.

| Field | Type | Default | Description |
|---|---|---|---|
| `options.healthCheck` | bool | `true` | Include a `/health` endpoint in the generated server |
| `options.middleware` | bool | `true` | Generate `server/middleware.go` (logging, CORS, rate limiter) |
| `options.exampleMain` | bool | `true` | Generate `cmd/server/main.go` as a runnable entry point |

---

## Generated File Reference

| Path | Package | Description | Always Generated? |
|---|---|---|---|
| `models/models.go` | `models` | Go structs with `json` tags; includes `Validate()` methods when validation is enabled | Yes (if `generate.models: true`) |
| `models/models_test.go` | `models` | Unit tests for `Validate()` methods | Yes (if `generate.models: true`) |
| `client/client.go` | `client` | HTTP client with one typed method per operation | Only when spec has endpoints and `generate.client: true` |
| `server/interface.go` | `server` | Handler interface; one method per operation grouped by tag | Only when spec has endpoints and `generate.server: true` |
| `server/mocks.go` | `server` | Stub implementations of the handler interface for quick start | Only when spec has endpoints and `generate.server: true` |
| `server/handlers.go` | `server` | HTTP adapter layer that calls the handler interface methods | Only when spec has endpoints and `generate.server: true` |
| `server/routes.go` | `server` | Route registration functions, one per tag group | Only when spec has endpoints and `generate.server: true` |
| `server/middleware.go` | `server` | Logging, CORS, and rate limiter middleware (stdlib only) | Only when spec has endpoints, `generate.server: true`, and `options.middleware: true` |
| `cmd/server/main.go` | `main` | Runnable example entry point; intended to be edited freely | Only when spec has endpoints, `generate.server: true`, and `options.exampleMain: true` |

---

## Supported OpenAPI Features

| Feature | Status | Notes |
|---|---|---|
| OpenAPI 3.x | Full | Versions 3.0.x and 3.1.x |
| `components/schemas` — object | Full | Generates a named Go struct |
| `components/schemas` — array | Full | Generates a Go slice type |
| `components/schemas` — map (`additionalProperties`) | Full | Generates `map[string]T` |
| `components/schemas` — enum | Full | Generates typed string/int constants |
| `components/schemas` — `allOf` | Full | Merged into a single struct |
| `components/schemas` — `oneOf` with discriminator | Full | Discriminator field used to select the concrete type |
| `components/schemas` — `anyOf` | Full | Generates a union-style struct |
| Path parameters | Full | Extracted and passed to handler methods |
| Query parameters | Full | Extracted and passed to handler methods |
| Header parameters | Full | Extracted and passed to handler methods |
| Request bodies (JSON) | Full | Decoded into the appropriate model struct |
| Multiple response content types | Full | Content negotiation via `Accept` header |
| Security schemes — Bearer token | Partial | Extracted and exposed; enforcement is left to the user |
| Security schemes — API Key | Partial | Extracted and exposed; enforcement is left to the user |
| Security schemes — Basic auth | Partial | Extracted and exposed; enforcement is left to the user |
| Security schemes — OAuth2 | Partial | Extracted and exposed; enforcement is left to the user |
| External `$ref` (local files) | Full | Resolved relative to the spec file |
| External `$ref` (URLs) | Full | Fetched over HTTPS at generation time |
| Remote spec (`-spec https://...`) | Full | Spec fetched over HTTPS before parsing |
| XML bodies | Not Supported | — |
| `application/x-www-form-urlencoded` bodies | Not Supported | — |
| Webhooks | Not Supported | — |
| Callbacks | Not Supported | — |
| Links | Not Supported | — |

---

## Type Mapping

The table below shows how OpenAPI type + format combinations are mapped to Go types by default. Any of these can be overridden with `types.overrides` in `beaver.yaml`.

| OpenAPI Type | OpenAPI Format | Go Type |
|---|---|---|
| `string` | _(none)_ | `string` |
| `string` | `date` | `string` |
| `string` | `date-time` | `time.Time` |
| `string` | `email` | `string` |
| `string` | `uuid` | `string` |
| `string` | `uri` | `string` |
| `string` | `byte` | `[]byte` |
| `string` | `binary` | `[]byte` |
| `string` | `password` | `string` |
| `integer` | _(none)_ | `int` |
| `integer` | `int32` | `int32` |
| `integer` | `int64` | `int64` |
| `number` | _(none)_ | `float64` |
| `number` | `float` | `float32` |
| `number` | `double` | `float64` |
| `boolean` | _(none)_ | `bool` |
| `array` | _(none)_ | `[]T` |
| `object` | _(none)_ | named struct or `map[string]interface{}` |
| `object` with `additionalProperties` | _(none)_ | `map[string]T` |
