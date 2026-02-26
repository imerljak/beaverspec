# Phase 2 Overview: Complete OpenAPI Support + Go Generator + Shared Tooling

**Timeline**: 10-14 weeks
**Goal**: Production-ready Go generator with all OpenAPI 3.x features and reusable cross-language utilities

---

## What We're Building

### 1. Complete OpenAPI 3.x Support (Core)
All OpenAPI specification features properly parsed and available to generators:
- ✅ Basic types, refs, nested objects, enums (DONE)
- ✅ allOf, oneOf, anyOf, discriminators (DONE)
- ✅ additionalProperties, arrays (DONE)
- 🔲 Content negotiation (multiple media types)
- 🔲 Response headers
- 🔲 Security schemes (applied to operations)
- 🔲 Callbacks (webhooks)
- 🔲 Links
- 🔲 Cookie parameters
- 🔲 Server variables
- 🔲 Examples (plural)
- 🔲 Encoding objects

### 2. Production-Ready Go Generator
Full-featured Go code generation for real-world use:
- ✅ Models with all OpenAPI features (DONE)
- ✅ Basic HTTP client (DONE)
- 🔲 Advanced client (retries, timeouts, interceptors)
- 🔲 Server interfaces and handlers
- 🔲 Multiple frameworks (net/http, Chi, Echo, Gin)
- 🔲 Authentication/authorization
- 🔲 Request/response validation
- 🔲 Typed errors
- 🔲 Documentation generation
- 🔲 Test generation
- 🔲 Mock generation

### 3. Shared Tooling (Reusable Across Languages)
Cross-cutting utilities for all generators:
- 🔲 Import management (`pkg/codegen/imports`)
- 🔲 Naming conventions (`pkg/codegen/naming`)
- 🔲 HTTP utilities (`pkg/codegen/http`)
- 🔲 Content negotiation (`pkg/codegen/content`)
- 🔲 Error patterns (`pkg/codegen/errors`)
- 🔲 Security patterns (`pkg/codegen/security`)
- 🔲 Validation (`pkg/codegen/validation`)
- 🔲 Framework abstraction (`pkg/codegen/framework`)
- 🔲 Code formatting (`pkg/codegen/format`)
- 🔲 Documentation (`pkg/codegen/docs`)
- 🔲 Testing utilities (`pkg/codegen/testing`)
- 🔲 Resilience patterns (`pkg/codegen/resilience`)
- 🔲 Middleware patterns (`pkg/codegen/middleware`)
- 🔲 Config management (`pkg/config/loader`)
- 🔲 Template overrides (`pkg/template/override`)

---

## Implementation Strategy

### Parallel Development Approach

We'll work on **three tracks simultaneously**:

```
Track 1: OpenAPI Features (Core)
└─ Add features to parser/normalizer as needed

Track 2: Go Generator (Language-Specific)
└─ Implement features using Go idioms

Track 3: Shared Tooling (Cross-Language)
└─ Extract patterns as we identify them
```

**Why this works:**
- Implement feature in Go first (learn the problem)
- Extract reusable pattern (solve it generally)
- Document for other languages (enable future generators)

### Example Flow

**Milestone 2.1: Package Structure**
1. Implement import management for Go → Learn the problem
2. Extract generic import resolution logic → Create `pkg/codegen/imports`
3. Document how TypeScript/Python imports differ → Enable future work

**Milestone 2.4: Typed Errors**
1. Implement typed errors for Go → Learn the problem
2. Extract error naming and categorization → Create `pkg/codegen/errors`
3. Document how other languages handle errors → Enable future work

---

## Milestones Quick Reference

| # | Milestone | Duration | Deliverables |
|---|-----------|----------|--------------|
| 2.1 | Package Structure & Imports | 3-5 days | Import management, proper package organization |
| 2.2 | Server Generation (Basic) | 5-7 days | Server interface, handlers, routes |
| 2.3 | Content Negotiation | 4-6 days | Multiple media types, file uploads |
| 2.4 | Response Headers & Typed Errors | 4-6 days | Typed errors, header handling |
| 2.5 | Security Scheme Implementation | 5-7 days | Auth in client/server, all auth types |
| 2.6 | Request Validation | 5-7 days | Validation generation, all constraints |
| 2.7 | Framework Support | 6-8 days | Chi, Echo, Gin support |
| 2.8 | Advanced OpenAPI Features | 5-7 days | Callbacks, links, cookies, etc. |
| 2.9 | Code Quality & DX | 5-7 days | Formatting, docs, tests, mocks |
| 2.10 | Advanced Client Features | 4-6 days | Retries, circuit breakers, interceptors |
| 2.11 | Advanced Server Features | 4-6 days | Middleware, observability, health checks |
| 2.12 | Configuration & Customization | 3-5 days | Flexible config, template overrides |
| 2.13 | Integration & Polish | 5-7 days | E2E testing, examples, optimization |

**Total: 10-14 weeks**

---

## Architecture After Phase 2

```
beaverspec/
├── cmd/beaverspec/              # CLI
├── pkg/
│   ├── core/                     # OpenAPI spec representation
│   ├── parser/                   # OpenAPI parsing & normalization
│   ├── plugin/                   # Generator registry
│   ├── template/                 # Template engine
│   │   └── override.go          # NEW: Template override system
│   ├── config/                   # NEW: Enhanced configuration
│   │   └── loader.go
│   └── codegen/                  # NEW: Shared tooling
│       ├── imports.go           # Import management
│       ├── naming.go            # Naming conventions
│       ├── http.go              # HTTP utilities
│       ├── content.go           # Content negotiation
│       ├── errors.go            # Error patterns
│       ├── security.go          # Security patterns
│       ├── validation.go        # Validation
│       ├── framework.go         # Framework abstraction
│       ├── format.go            # Code formatting
│       ├── docs.go              # Documentation
│       ├── testing.go           # Test generation
│       ├── resilience.go        # Retry/timeout/circuit breaker
│       └── middleware.go        # Middleware patterns
├── generators/
│   ├── base/                    # Base generator utilities
│   └── golang/                  # Go generator
│       ├── generator.go         # Generator implementation
│       ├── config.yaml          # Type mappings, naming, etc.
│       ├── templates/
│       │   ├── core/            # Models, types, errors
│       │   │   ├── models.tmpl
│       │   │   ├── types.tmpl
│       │   │   └── errors.tmpl
│       │   ├── client/          # Client generation
│       │   │   ├── interface.tmpl
│       │   │   ├── http.tmpl
│       │   │   └── mock.tmpl
│       │   └── server/          # Server generation
│       │       ├── interface.tmpl
│       │       ├── handlers.tmpl
│       │       └── routes.tmpl
│       └── frameworks/          # Framework-specific templates
│           ├── net-http/
│           ├── chi/
│           ├── echo/
│           └── gin/
└── examples/
    ├── petstore/                # Simple example
    ├── github/                  # Real-world example
    └── fullstack/               # Complete app example
```

---

## Key Decisions & Rationale

### 1. Implement in Go First, Extract Later
**Decision**: Build features in Go generator, then extract shared tooling
**Rationale**:
- Learn by doing (real problems > theoretical solutions)
- Faster initial development
- Better shared APIs (informed by real use)

### 2. Language-Specific Package Organization
**Decision**: Let each generator decide its package structure
**Rationale**:
- Go prefers separate packages, TS prefers barrel exports, Python prefers modules
- Forcing universal structure leads to unidiomatic code
- Shared import management utilities can handle different patterns

### 3. Framework Support via Template Overrides
**Decision**: Use template inheritance for framework-specific code
**Rationale**:
- Most framework differences are syntactic (how to register routes, etc.)
- Template overrides are simpler than code generation branches
- Easy for users to customize or add frameworks

### 4. Shared Tooling in `pkg/codegen`
**Decision**: Create dedicated package for cross-language utilities
**Rationale**:
- Clear separation from OpenAPI-specific logic (`pkg/core`, `pkg/parser`)
- Reusable by all generators
- Can be tested independently

### 5. Security Features in Phase 2
**Decision**: Implement auth/security before Phase 3 (other languages)
**Rationale**:
- Security is critical for production use
- Patterns will inform shared tooling
- TypeScript/Python generators will benefit from lessons learned

---

## Success Metrics

Phase 2 is complete when:

### Functional Completeness
- [ ] All OpenAPI 3.x features supported (except rare edge cases)
- [ ] Can generate production-ready Go client
- [ ] Can generate production-ready Go server
- [ ] Support 4+ Go frameworks
- [ ] Full authentication/authorization support
- [ ] Request/response validation

### Code Quality
- [ ] Generated code passes `go vet` and `golangci-lint`
- [ ] Generated code includes comprehensive tests
- [ ] Generated code is well-documented
- [ ] Generated code handles errors properly

### Developer Experience
- [ ] Clear error messages with actionable hints
- [ ] Configuration is intuitive and well-documented
- [ ] Template customization is straightforward
- [ ] CLI provides helpful feedback

### Shared Tooling
- [ ] All 15 shared packages implemented
- [ ] Each package has comprehensive tests
- [ ] Each package is documented with examples
- [ ] Packages are used by Go generator
- [ ] Packages designed for TypeScript/Python use

### Documentation
- [ ] Complete API reference
- [ ] Tutorial for Go generator
- [ ] Template customization guide
- [ ] Migration guide from openapi-generator
- [ ] Troubleshooting guide

### Examples
- [ ] Petstore (simple example)
- [ ] Real-world API (GitHub/Stripe/etc.)
- [ ] Microservices example
- [ ] Full-stack application

### Performance
- [ ] Generates petstore in < 100ms
- [ ] Generates large specs (100+ endpoints) in < 1s
- [ ] Memory usage < 50MB for typical specs
- [ ] Concurrent generation support

---

## Risk Mitigation

### Risk: Scope Creep
**Mitigation**:
- Strict milestone boundaries
- Only implement features needed for Go generator
- Defer "nice to have" features to Phase 3
- Regular scope reviews

### Risk: Over-Engineering Shared Tooling
**Mitigation**:
- Start simple, extract as needed
- Validate with real use cases
- Don't try to predict all language needs
- Iterate based on feedback

### Risk: Go-Specific Bias in Shared Code
**Mitigation**:
- Document how other languages differ
- Consider TypeScript/Python patterns during design
- Keep shared code abstract (interfaces)
- Review shared APIs with multi-language lens

### Risk: Template Complexity
**Mitigation**:
- Keep templates simple (logic in Go, not templates)
- Break into small, focused templates
- Comprehensive template testing
- Template documentation and examples

---

## Next Steps

### Immediate Actions (This Week)
1. Start Milestone 2.1 (Package Structure & Imports)
2. Create `pkg/codegen/imports.go` package
3. Create `pkg/codegen/naming.go` package
4. Update Go generator to use separate packages
5. Test with petstore example

### Short-Term (Next 2-4 Weeks)
1. Complete Milestones 2.1-2.3
2. Basic server generation working
3. Content negotiation support
4. Multiple examples tested

### Mid-Term (Next 4-8 Weeks)
1. Complete Milestones 2.4-2.8
2. All OpenAPI features implemented
3. Security and validation complete
4. Framework support complete

### Long-Term (Next 8-14 Weeks)
1. Complete Milestones 2.9-2.13
2. Production-ready Go generator
3. All shared tooling extracted and documented
4. Ready for Phase 3 (TypeScript generator)

---

## Questions to Answer During Phase 2

1. **Import Management**: Can we create a truly generic import resolver?
2. **Framework Abstraction**: Is template override sufficient or do we need plugin system?
3. **Validation**: Should validation be runtime or compile-time (where possible)?
4. **Error Handling**: What error patterns work across all languages?
5. **Testing**: Can we generate meaningful tests or just stubs?
6. **Documentation**: What documentation formats are most useful?
7. **Performance**: Where are the bottlenecks? Template rendering? Normalization?
8. **Configuration**: What needs to be configurable vs opinionated?

These questions will inform Phase 3 (additional languages) and Phase 4 (advanced features).

---

## Related Documents

- [PHASE_2_ROADMAP.md](./PHASE_2_ROADMAP.md) - Detailed milestone breakdown
- [SHARED_TOOLING.md](./SHARED_TOOLING.md) - Shared tooling architecture and APIs
- [ARCHITECTURE.md](./ARCHITECTURE.md) - Overall project architecture
- [ROADMAP.md](./ROADMAP.md) - Complete project roadmap (all phases)
- [CLAUDE.md](./CLAUDE.md) - Project context for Claude Code
