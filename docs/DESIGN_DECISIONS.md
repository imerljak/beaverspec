# Design Decisions

This document tracks important architectural and design decisions made during the development of the OpenAPI Code Generator.

## Format

Each decision follows this template:

```
## [ID] Decision Title

**Date**: YYYY-MM-DD
**Status**: Proposed | Accepted | Deprecated | Superseded
**Context**: Why this decision was needed
**Decision**: What we decided
**Consequences**: Positive and negative outcomes
**Alternatives Considered**: What else we looked at
```

---

## [DD-001] Plugin-Based Architecture

**Date**: 2026-01-03
**Status**: Accepted

### Context

The project needs to support multiple programming languages for code generation. We needed to decide how to structure the codebase to allow adding new languages without modifying core logic.

### Decision

Implement a plugin-based architecture where:
- Core system handles OpenAPI parsing, orchestration, and file I/O
- Each language is implemented as a plugin with a standardized interface
- Plugins can be compiled into the binary (Phase 1) or loaded dynamically (future)

### Consequences

**Positive**:
- Clear separation of concerns
- Easy to add new languages
- Core logic remains stable
- Plugins can be developed independently
- Testable in isolation

**Negative**:
- More complex initial setup
- Need to design a good plugin interface upfront
- Potential performance overhead (minimal for compile-time plugins)

### Alternatives Considered

1. **Monolithic per-language approach**: Separate tool per language
   - Rejected: Too much code duplication, inconsistent UX

2. **Template-only approach**: Everything in templates
   - Rejected: Templates would become too complex for type mapping logic

3. **Code generation framework**: Use existing framework
   - Rejected: Want full control and Go-native solution

---

## [DD-002] Templates vs Plugins Boundary

**Date**: 2026-01-03
**Status**: Accepted

### Context

We needed to decide what logic goes in templates vs what goes in plugin code.

### Decision

- **Templates**: Handle syntax and code structure (the "how to write")
- **Plugins**: Handle semantics and language-specific logic (the "what it means")

Specifically:
- Templates: Code formatting, imports, class/function structures
- Plugins: Type mappings, nullable strategies, validation rules

### Consequences

**Positive**:
- Templates remain simple and readable
- Complex logic is type-safe Go code
- Clear responsibility boundaries
- Templates can be customized without touching Go code

**Negative**:
- Need both templates AND plugin code for complex generators
- Learning curve for contributors (need to understand both)

### Alternatives Considered

1. **Everything in templates**: Use template helpers for all logic
   - Rejected: Templates become unmaintainable

2. **Everything in code**: Generate strings programmatically
   - Rejected: Hard to customize, no separation of concerns

3. **Template inheritance**: Complex template system with OOP-like features
   - Rejected: Over-engineered for our needs

---

## [DD-003] Data-Driven Plugin Configuration

**Date**: 2026-01-03
**Status**: Accepted

### Context

We want to make it easy to add new language support without requiring deep Go knowledge. Need to decide how much should be configurable vs coded.

### Decision

Implement a hybrid approach where:
- 80% of generator behavior is defined in YAML configuration
- 20% complex cases can use custom Go code
- Config includes: type mappings, naming conventions, template references, framework support

### Consequences

**Positive**:
- Low barrier to add new languages
- Most languages require no Go code
- Configuration is self-documenting
- Easy to understand and modify

**Negative**:
- Config format needs to be comprehensive
- Need to design flexible mapping system
- Edge cases might push against config limitations

### Alternatives Considered

1. **Code-only plugins**: Everything in Go
   - Rejected: High barrier for contributors

2. **DSL for generators**: Custom language for generator logic
   - Rejected: Over-engineered, learning curve too high

3. **Pure template-driven**: Everything in template logic
   - Rejected: Type mappings don't fit well in templates

---

## [DD-004] Template Granularity (Use-Specific)

**Date**: 2026-01-03
**Status**: Accepted

### Context

Need to decide the level of granularity for templates - one big template vs many small templates.

### Decision

Use **use-specific templates**:
- Separate templates for models, client, server, validation
- Further subdivision within each category (e.g., client/base, client/http)
- Framework-specific templates can override base templates

### Consequences

**Positive**:
- Generate only what's needed (client-only, server-only)
- Easier to maintain smaller templates
- Users can override specific parts
- Can mix and match components

**Negative**:
- More files to manage
- Need template composition/includes
- Coordination between related templates

### Alternatives Considered

1. **Monolithic template**: One template per language
   - Rejected: Hard to maintain, all-or-nothing generation

2. **Per-file templates**: One template per output file
   - Rejected: Too granular, coordination nightmare

3. **Per-schema templates**: Template per schema type
   - Rejected: Doesn't align with user needs (want client/server selection)

---

## [DD-005] Configuration Hierarchy

**Date**: 2026-01-03
**Status**: Accepted

### Context

Users need to configure the generator at different levels (global, project, per-run). Need to define precedence and sources.

### Decision

Five-level hierarchy (highest to lowest precedence):
1. CLI flags (per-run overrides)
2. Environment variables (deployment-specific)
3. Project config file (`.beaverspec.yaml`)
4. User config file (`~/.beaverspec/config.yaml`)
5. Generator defaults (from generator's `config.yaml`)

### Consequences

**Positive**:
- Flexible for different use cases
- Reasonable defaults at each level
- Clear override mechanism
- Supports both local and CI/CD workflows

**Negative**:
- Can be confusing to debug which config is active
- Need good tooling to show effective configuration

### Alternatives Considered

1. **Single config file only**: No hierarchy
   - Rejected: Not flexible enough for different environments

2. **Flat precedence**: CLI > File > Defaults
   - Rejected: Doesn't support user-wide preferences

3. **Inheritance model**: Config files inherit from each other
   - Rejected: Too complex, hard to understand

---

## [DD-006] Compile-Time Plugin Loading (Phase 1)

**Date**: 2026-01-03
**Status**: Accepted

### Context

Need to decide plugin loading strategy for MVP. Balance between simplicity and flexibility.

### Decision

**Phase 1**: Compile-time registration
- All generators built into binary
- Templates embedded using `//go:embed`
- Simple registry with compile-time registration

**Future phases**: Runtime loading, marketplace

### Consequences

**Positive**:
- Simple implementation
- Zero runtime dependencies
- Fast startup
- Easy to distribute (single binary)
- Type-safe

**Negative**:
- Can't add generators without recompiling
- Binary size grows with generators
- Less flexible for users

### Alternatives Considered

1. **Runtime plugins from day 1**: Go plugin system
   - Rejected: Complex, platform-limited, harder to debug

2. **External process plugins**: Fork/exec style
   - Rejected: Over-engineered for MVP, slower

3. **No plugin system**: Hardcoded generators
   - Rejected: Not extensible

---

## [DD-007] Go Templates for Rendering

**Date**: 2026-01-03
**Status**: Accepted

### Context

Need to choose a template engine for code generation.

### Decision

Use Go's built-in `text/template` package with custom helper functions.

### Consequences

**Positive**:
- No external dependencies
- Well-documented and stable
- Powerful enough for our needs
- Go developers already familiar
- Good error messages

**Negative**:
- Less features than Mustache/Handlebars
- Go's template syntax has quirks
- No template inheritance (need to implement)

### Alternatives Considered

1. **Mustache**: Language-agnostic templates
   - Rejected: Too limited (same problem as openapi-generator)

2. **Handlebars**: More powerful Mustache
   - Rejected: Need Go library, adds dependency

3. **Custom DSL**: Build our own
   - Rejected: Reinventing the wheel, maintenance burden

---

## [DD-008] Error Collection Strategy

**Date**: 2026-01-03
**Status**: Accepted

### Context

Need to decide error handling philosophy - fail fast vs collect all errors.

### Decision

Three-tier error handling:
- **Warnings**: Collect all, display at end, continue generation
- **Errors**: Collect all, stop generation, display all errors
- **Fatal**: Stop immediately

### Consequences

**Positive**:
- Users see all validation errors at once
- Better UX (fix multiple issues in one go)
- Clear severity levels

**Negative**:
- More complex error collection
- Need to handle partial state
- Some errors might cascade

### Alternatives Considered

1. **Fail fast**: Stop on first error
   - Rejected: Poor UX, user has to fix one at a time

2. **Best effort**: Try to continue despite errors
   - Rejected: Can produce invalid output

3. **Warnings only**: Never stop generation
   - Rejected: Dangerous, could generate bad code

---

## [DD-009] OpenAPI Parser Library

**Date**: 2026-01-03
**Status**: Accepted

### Context

Need to parse and validate OpenAPI specifications. Decide between building custom parser or using existing library.

### Decision

Use `github.com/getkin/kin-openapi` as the OpenAPI parser.

### Consequences

**Positive**:
- Battle-tested, widely used
- Full OpenAPI 3.x support
- Good validation
- Active maintenance
- Go-native

**Negative**:
- External dependency
- API might change in future
- Need to wrap/adapt for our needs

### Alternatives Considered

1. **Build custom parser**: Parse YAML ourselves
   - Rejected: Huge amount of work, error-prone

2. **go-swagger**: Alternative parser
   - Rejected: Focuses on Swagger 2.0, older codebase

3. **Use JSON Schema validator**: Lower-level approach
   - Rejected: Still need to parse OpenAPI structure

---

## [DD-010] Output Organization

**Date**: 2026-01-03
**Status**: Accepted

### Context

Need to decide how generated files should be organized in output directory.

### Decision

Support three modes (configurable):
1. **Flat**: All files in one directory
2. **Organized**: Grouped by concern (models/, client/, server/)
3. **Custom**: User-defined pattern

Default to "organized" structure.

### Consequences

**Positive**:
- Flexibility for different project structures
- Organized mode works well for most cases
- Custom mode allows full control

**Negative**:
- Need to implement all three modes
- Path handling gets more complex

### Alternatives Considered

1. **Fixed structure only**: No choice
   - Rejected: Different projects have different conventions

2. **Organized only**: Remove flat mode
   - Rejected: Some projects prefer flat structure

3. **Custom only**: Pattern-based from the start
   - Rejected: Too complex for simple use cases

---

## [DD-011] Framework Variants

**Date**: 2026-01-03
**Status**: Accepted

### Context

Different languages have popular frameworks (Chi, Echo, Gin for Go). Need to decide how to support them.

### Decision

Support framework variants through:
- Base templates that work generically
- Framework-specific template directories that override base
- Template selection based on config

Example: `frameworks/chi/templates/` overrides `templates/server/`

### Consequences

**Positive**:
- Can generate idiomatic code for each framework
- Framework-specific features can be leveraged
- Base templates provide fallback

**Negative**:
- Need to maintain multiple template sets
- Potential divergence between frameworks
- More testing required

### Alternatives Considered

1. **Generic only**: Framework-agnostic code
   - Rejected: Loses framework benefits, users customize anyway

2. **Framework per generator**: Separate generator per framework
   - Rejected: Too much duplication

3. **Post-generation adapters**: Generate generic, transform after
   - Rejected: Complex, brittle

---

## [DD-012] Validation Generation

**Date**: 2026-01-03
**Status**: Accepted

### Context

OpenAPI specs contain validation rules (min/max, pattern, etc.). Need to decide if/how to generate validators.

### Decision

Make validation generation **optional and separate**:
- Not generated by default
- Can be enabled via config flag
- Generated as separate files/functions
- Use language-specific validation libraries where available

### Consequences

**Positive**:
- Keeps generated code cleaner
- Users opt-in to validation
- Can leverage existing libraries
- Reduces generated code size

**Negative**:
- More work to implement
- Another thing to configure
- Validation might be inconsistent across languages

### Alternatives Considered

1. **Always generate validation**: Built into models
   - Rejected: Makes models bloated, not always needed

2. **Never generate validation**: Users do it themselves
   - Rejected: Loses valuable spec information

3. **Validation via annotations only**: Just comments
   - Rejected: Not executable, easy to ignore

---

## [DD-013] Spec Normalization Strategy

**Date**: 2026-01-03
**Status**: Accepted

### Context

OpenAPI specs can be complex (allOf, oneOf, $refs). Need to decide how much to normalize/simplify before generation.

### Decision

Create "enriched spec model" that:
- Resolves all `$ref` references
- Flattens `allOf` where possible
- Preserves `oneOf`/`anyOf` for generator to handle
- Pre-computes commonly needed information
- Maintains link to original spec

### Consequences

**Positive**:
- Generators work with simpler model
- Shared normalization logic (DRY)
- Can add caching/memoization
- Easier to reason about

**Negative**:
- Additional processing step
- Memory overhead for enriched model
- Need to maintain both representations

### Alternatives Considered

1. **No normalization**: Generators handle raw spec
   - Rejected: Every generator reimplements same logic

2. **Full normalization**: Resolve everything
   - Rejected: Loses important distinctions (oneOf semantics)

3. **Lazy normalization**: Resolve on-demand
   - Rejected: Complex to implement correctly

---

## Future Decisions

Items that need decisions but are deferred:

- **[DD-014]** Runtime plugin loading mechanism
- **[DD-015]** Plugin marketplace/registry design
- **[DD-016]** Caching strategy for parsed specs
- **[DD-017]** Watch mode implementation
- **[DD-018]** IDE integration approach
- **[DD-019]** Code formatting integration
- **[DD-020]** Migration tools from openapi-generator

---

## Decision Review Process

Decisions should be reviewed:
- When implementation reveals unforeseen issues
- When user feedback suggests different approach
- After 6 months of use
- When new OpenAPI features emerge

## Changing Decisions

To change a decision:
1. Propose new decision with "Supersedes DD-XXX"
2. Mark old decision as "Superseded by DD-YYY"
3. Document migration path if needed
4. Update affected documentation

---

## References

- [Architectural Decision Records](https://adr.github.io/)
- [ARCHITECTURE.md](./ARCHITECTURE.md)
- [TECHNICAL_SPEC.md](./TECHNICAL_SPEC.md)
