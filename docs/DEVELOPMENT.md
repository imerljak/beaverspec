# Development Guide

This guide covers everything needed to work on the codebase itself — building, testing, tooling, and conventions.

## About This Project

This is a **learning project** designed to practice Go concepts while solving a real problem. The goals are:

- Interfaces and type embedding
- Template processing (`text/template`)
- File I/O and path handling
- CLI application building with Cobra
- Testing patterns in Go

If you're new to the codebase, start by reading [ARCHITECTURE.md](ARCHITECTURE.md) before writing any code.

---

## Technical Stack

| Concern | Library | Version |
|---|---|---|
| Language | Go | 1.21+ |
| OpenAPI parsing | `github.com/getkin/kin-openapi/openapi3` | latest |
| CLI framework | `github.com/spf13/cobra` | latest |
| Template engine | Go stdlib `text/template` | — |
| Testing | `github.com/stretchr/testify` | latest |

---

## Development Commands

### Build

```bash
go mod tidy
go build -o beaver ./cmd/beaver
```

### Run

```bash
./beaver generate --spec ./examples/petstore/openapi.yaml --generator go
```

### Test

```bash
# All tests
go test ./...

# With coverage
go test -cover ./...

# Specific package
go test ./pkg/parser/

# Integration tests (future)
go test -tags=integration ./...
```

### Code Quality

```bash
# Format
go fmt ./...
gofmt -s -w .

# Vet
go vet ./...

# Lint
golangci-lint run
```

---

## Implementation Guidelines

Follow these when adding or modifying features:

1. **Start with interfaces** — Define the contract in `pkg/core/` before writing any implementation.
2. **Write tests first** — Especially for core logic, transformations, and type mappings.
3. **Keep templates simple** — Move complex logic into Go code, not into `.tmpl` files.
4. **Validate early, collect all errors** — Never fail on the first error; gather all validation issues and report them together with context (`file`, `line`, JSON path).
5. **Godoc all exported types** — Every exported symbol in `pkg/` needs a godoc comment.
6. **Follow Go conventions** — `gofmt`, standard project layout, idiomatic error handling.
7. **Test with real specs** — Use the examples in `examples/` to validate end-to-end behaviour.

---

## Commit Message Convention

```
feat:     New feature
fix:      Bug fix
docs:     Documentation only
refactor: Code change without feature or fix
test:     Adding or updating tests
chore:    Tooling, dependencies, maintenance
```

Examples:
```
feat: add oneOf discriminator support for Go generator
fix: resolve $ref cycles in schema normalizer
test: add roundtrip tests for allOf composition
```

---

## Documentation Map

| File | Purpose | Read when... |
|---|---|---|
| [ARCHITECTURE.md](ARCHITECTURE.md) | System design, component boundaries, data flow | Before writing any code |
| [TECHNICAL_SPEC.md](TECHNICAL_SPEC.md) | Interface contracts, data models, CLI spec | Implementing a new generator or core feature |
| [DESIGN_DECISIONS.md](DESIGN_DECISIONS.md) | Rationale for key choices | Wondering _why_ something is the way it is |
| [ROADMAP.md](ROADMAP.md) | All phases and milestones | Planning new work |
| [PHASE_2_ROADMAP.md](PHASE_2_ROADMAP.md) | Detailed Phase 2 breakdown | Working on Phase 2 features |
| [SHARED_TOOLING.md](SHARED_TOOLING.md) | Shared helpers, APIs across generators | Adding a new generator |
| [QUICK_START.md](QUICK_START.md) | User-facing usage guide | Understanding the end-user experience |
