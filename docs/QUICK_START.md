# BeaverSpec Quick Start

## What is BeaverSpec?

BeaverSpec (`beaver`) is a CLI tool that generates Go server, client, and model code from OpenAPI 3.x specifications. Point it at a spec file or URL and it produces idiomatic, ready-to-compile Go code for the framework of your choice.

---

## Installation

### go install (recommended)

```bash
go install github.com/imerljak/beaverspec/cmd/beaver@latest
```

### Build from source

```bash
git clone https://github.com/imerljak/beaverspec.git
cd beaverspec
go build -o beaver ./cmd/beaver
```

---

## 5-Minute Quickstart

### 1. Write a minimal spec

Create `openapi.yaml`:

```yaml
openapi: "3.0.3"
info:
  title: Notes API
  version: "1.0.0"
paths:
  /notes:
    get:
      summary: List notes
      operationId: listNotes
      tags: [notes]
      responses:
        "200":
          description: OK
          content:
            application/json:
              schema:
                type: array
                items:
                  $ref: "#/components/schemas/Note"
components:
  schemas:
    Note:
      type: object
      required: [id, body]
      properties:
        id:
          type: integer
          format: int64
        body:
          type: string
```

### 2. Generate code

```bash
beaver -spec openapi.yaml -output ./generated -module github.com/example/myapi
```

### 3. What gets generated

```
generated/
├── models/
│   ├── models.go          # Go structs with json tags and Validate() methods
│   └── models_test.go     # Validation tests
├── client/
│   └── client.go          # HTTP client with typed methods per endpoint
├── server/
│   ├── interface.go       # Handler interface — implement this
│   ├── mocks.go           # Mock implementations for quick start
│   ├── handlers.go        # HTTP adapter layer
│   ├── routes.go          # Route registration per tag group
│   └── middleware.go      # Logging, CORS, rate limiter (stdlib)
└── cmd/server/
    └── main.go            # Runnable example entry point
```

For specs with no endpoints (schemas only), only `models/models.go` and `models/models_test.go` are generated.

---

## Using the Generated Server

### Step 1 — Add the generated code to your module

The generated code lives under the output directory you specified. Make sure the `-module` value matches the `module` line in your `go.mod`:

```
module github.com/example/myapi

go 1.21
```

Run `go mod tidy` from the root of your project after generation.

### Step 2 — Implement the handler interface

Open `generated/server/interface.go`. It will contain something like:

```go
type NotesHandler interface {
    ListNotes(w http.ResponseWriter, r *http.Request) error
}
```

Create your own implementation:

```go
package main

import (
    "encoding/json"
    "net/http"

    "github.com/example/myapi/generated/models"
)

type MyNotesHandler struct{}

func (h *MyNotesHandler) ListNotes(w http.ResponseWriter, r *http.Request) error {
    notes := []models.Note{{Id: 1, Body: "hello"}}
    w.Header().Set("Content-Type", "application/json")
    return json.NewEncoder(w).Encode(notes)
}
```

### Step 3 — Wire it up and run

`cmd/server/main.go` is a real, runnable entry point. Edit it freely — it is yours to own. It wires your handler implementation to the generated routes. Run it with:

```bash
go run ./generated/cmd/server/main.go
```

---

## Framework Selection

Use the `-framework` flag to select a server framework. The generated `handlers.go`, `routes.go`, and `cmd/server/main.go` will use the framework's idioms.

```bash
# stdlib net/http (default)
beaver -spec openapi.yaml -output ./generated -module github.com/example/myapi -framework net-http

# chi
beaver -spec openapi.yaml -output ./generated -module github.com/example/myapi -framework chi

# echo
beaver -spec openapi.yaml -output ./generated -module github.com/example/myapi -framework echo

# gin
beaver -spec openapi.yaml -output ./generated -module github.com/example/myapi -framework gin
```

---

## Config File (beaver.yaml)

Instead of passing flags every time, place a `beaver.yaml` next to your spec file (or point to one with `-config`). All fields are optional.

```yaml
spec: ./openapi.yaml
output: ./generated
module: github.com/example/myapi
package: models
framework: chi

generate:
  models: true
  client: true
  server: true
  validation: true

exclude:
  models: [InternalError]   # skip generating these schema names
  tags: [internal]          # skip all operations under these tags

types:
  overrides:
    - format: uuid
      go: "github.com/google/uuid".UUID

options:
  healthCheck: true   # include /health endpoint (default: true)
  middleware: true    # generate server/middleware.go (default: true)
  exampleMain: true   # generate cmd/server/main.go (default: true)
```

With a config file in place, just run:

```bash
beaver
```

---

## Remote Specs

`beaver` can fetch a spec directly from a URL:

```bash
beaver \
  -spec https://petstore3.swagger.io/api/v3/openapi.json \
  -output ./generated \
  -module github.com/example/petstore
```

---

## Dry Run

Preview the files that would be generated without writing anything to disk:

```bash
beaver -spec openapi.yaml -output ./generated -module github.com/example/myapi -dry-run
```

This prints each file path that would be created, which is useful for understanding output before committing to a directory layout.
