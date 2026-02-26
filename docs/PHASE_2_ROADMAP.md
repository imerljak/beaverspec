# Phase 2 Detailed Roadmap
## Complete OpenAPI 3.x Support + Go Generator + Shared Tooling

This roadmap implements remaining OpenAPI 3.x features while building out the Go generator, extracting shared tooling as we identify cross-cutting concerns.

---

## Milestone 2.1: Package Structure & Import Management
**Duration**: 3-5 days
**Goal**: Establish proper package organization and cross-package imports

### Core Work
- [ ] **Shared Tooling**: Create `pkg/codegen/` package for cross-cutting concerns
  - [ ] `pkg/codegen/imports.go` - Import path resolution and management
    - Module path detection (go.mod, package.json, setup.py)
    - Relative vs absolute import path construction
    - Import deduplication and sorting
  - [ ] `pkg/codegen/naming.go` - Naming convention utilities
    - Case conversion (camelCase, PascalCase, snake_case, kebab-case)
    - Pluralization/singularization
    - Keyword collision detection
    - Language-specific reserved word handling

### Go Generator Work
- [ ] Implement package separation
  - [ ] Update `generator.go` to write to subdirectories (`models/`, `client/`, `server/`)
  - [ ] Add `ModulePath` to generator config (read from go.mod or CLI flag)
  - [ ] Update `collectClientImports()` to add models package import
  - [ ] Add import path construction helper
- [ ] Template updates
  - [ ] Update `client.go.tmpl` to include models import
  - [ ] Add package name to template data
- [ ] Testing
  - [ ] Test generated code compiles with proper imports
  - [ ] Test with different module paths

### Deliverables
- [ ] Reusable import management for all generators
- [ ] Go generator produces proper package structure
- [ ] Documentation on import resolution strategy

---

## Milestone 2.2: Server Generation (Basic)
**Duration**: 5-7 days
**Goal**: Generate basic HTTP server interfaces and handlers

### Core Work
- [ ] **Shared Tooling**: HTTP method/status code utilities
  - [ ] `pkg/codegen/http.go` - HTTP constants and helpers
    - Method validation and normalization (GET, POST, etc.)
    - Status code categorization (2xx, 4xx, 5xx)
    - Content-type parsing
    - Common header definitions
- [ ] Enhance `core.Endpoint` with server-specific metadata
  - [ ] Group endpoints by tag for route organization
  - [ ] Extract success vs error responses
  - [ ] Identify handler return types

### Go Generator Work
- [ ] Create server template data structures
  - [ ] `HandlerData` - per-endpoint handler info
  - [ ] `RouterData` - route registration info
  - [ ] `MiddlewareData` - for future middleware support
- [ ] Implement `convertEndpointsToServer()` function
  - [ ] Group by tags
  - [ ] Extract path parameters for routing
  - [ ] Map request/response types
- [ ] Create server templates
  - [ ] `server/interface.tmpl` - Server interface definition
  - [ ] `server/handlers.tmpl` - Handler function skeletons
  - [ ] `server/routes.tmpl` - Route registration (net/http)
- [ ] Update `Generate()` to include server generation
- [ ] Testing
  - [ ] Generate server from petstore.yaml
  - [ ] Verify server compiles
  - [ ] Test basic HTTP server startup

### Deliverables
- [ ] Shared HTTP utilities for all generators
- [ ] Basic Go server generation working
- [ ] Server interface + handler stubs + route registration

---

## Milestone 2.3: Content Negotiation & Multiple Media Types
**Duration**: 4-6 days
**Goal**: Support multiple content types beyond application/json

### Core Work
- [ ] **OpenAPI Feature**: Enhanced media type processing
  - [ ] Update `core.MediaType` to include encoding info
  - [ ] Update normalizer to extract all content types from request/response
  - [ ] Add content-type priority/preference handling
- [ ] **Shared Tooling**: Content negotiation utilities
  - [ ] `pkg/codegen/content.go` - Media type helpers
    - MIME type parsing and matching
    - Content negotiation strategies
    - Serialization format detection (JSON, XML, FormData, etc.)

### Go Generator Work
- [ ] Update client to support multiple content types
  - [ ] Add content-type selection to operation data
  - [ ] Update template to set appropriate Content-Type header
  - [ ] Support Accept header with preference list
- [ ] Update server to support multiple content types
  - [ ] Add content-type routing/switching
  - [ ] Generate request parsers per content type
  - [ ] Generate response serializers per content type
- [ ] Add support for common types:
  - [ ] `application/json` (already done)
  - [ ] `application/xml` (basic support)
  - [ ] `multipart/form-data` (file uploads)
  - [ ] `application/x-www-form-urlencoded` (forms)
  - [ ] `text/plain`
- [ ] Testing
  - [ ] Test client with multiple Accept types
  - [ ] Test server with different Content-Type requests

### Deliverables
- [ ] Content negotiation available to all generators
- [ ] Go client/server handle multiple media types
- [ ] Support for file uploads

---

## Milestone 2.4: Response Headers & Typed Errors
**Duration**: 4-6 days
**Goal**: Generate typed error responses and handle response headers

### Core Work
- [ ] **OpenAPI Feature**: Response processing enhancements
  - [ ] Update normalizer to extract response headers
  - [ ] Map status codes to response schemas
  - [ ] Identify error responses (4xx, 5xx)
  - [ ] Extract response examples
- [ ] **Shared Tooling**: Error handling patterns
  - [ ] `pkg/codegen/errors.go` - Error modeling helpers
    - Error type naming conventions
    - Status code → error mapping
    - Error hierarchy construction

### Go Generator Work
- [ ] Generate typed error structs
  - [ ] Create error types per response schema
  - [ ] Implement Error() interface
  - [ ] Include status code information
- [ ] Update client error handling
  - [ ] Parse error responses into typed errors
  - [ ] Return appropriate error type based on status code
  - [ ] Include response headers in error context
- [ ] Update server error handling
  - [ ] Generate error serialization functions
  - [ ] Add helper for writing error responses
  - [ ] Support custom error handlers
- [ ] Add response header support
  - [ ] Extract headers from successful responses
  - [ ] Generate header setting code in server
  - [ ] Generate header reading code in client
- [ ] Testing
  - [ ] Test typed error returns
  - [ ] Test response header handling
  - [ ] Test error response parsing

### Deliverables
- [ ] Error modeling utilities for all generators
- [ ] Go client returns typed errors
- [ ] Go server includes error helpers
- [ ] Response headers properly handled

---

## Milestone 2.5: Security Scheme Implementation
**Duration**: 5-7 days
**Goal**: Generate authentication/authorization code

### Core Work
- [ ] **OpenAPI Feature**: Security processing
  - [ ] Update normalizer to apply security requirements to endpoints
  - [ ] Resolve security scheme references
  - [ ] Handle global vs operation-level security
  - [ ] Support multiple security options (OR logic)
- [ ] **Shared Tooling**: Auth pattern helpers
  - [ ] `pkg/codegen/security.go` - Security utilities
    - Auth scheme identification
    - Token/credential injection patterns
    - Scope/permission checking helpers

### Go Generator Work
- [ ] Client authentication
  - [ ] Add auth configuration to client options
  - [ ] Generate auth injection code:
    - [ ] API Key (header, query, cookie)
    - [ ] HTTP Basic Auth
    - [ ] HTTP Bearer Token
    - [ ] OAuth2 (basic token injection)
  - [ ] Add per-request auth override
- [ ] Server authentication middleware
  - [ ] Generate auth middleware interfaces
  - [ ] Generate auth verification stubs
  - [ ] Add auth context to handlers
  - [ ] Generate scope checking helpers
- [ ] Testing
  - [ ] Test client with different auth schemes
  - [ ] Test server auth middleware integration
  - [ ] Test multiple auth options (OR logic)

### Deliverables
- [ ] Security pattern utilities for all generators
- [ ] Go client supports all auth schemes
- [ ] Go server includes auth middleware
- [ ] OAuth2 basic support

---

## Milestone 2.6: Request Validation
**Duration**: 5-7 days
**Goal**: Generate validation code for requests and responses

### Core Work
- [ ] **OpenAPI Feature**: Validation metadata extraction
  - [ ] Already captured in `core.Property` (pattern, min/max, etc.)
  - [ ] Add validation for enum values
  - [ ] Add format validation (email, uri, date-time, uuid, etc.)
  - [ ] Add required field validation
- [ ] **Shared Tooling**: Validation utilities
  - [ ] `pkg/codegen/validation.go` - Validation helpers
    - Validation expression generation
    - Common regex patterns (email, uri, etc.)
    - Validation error formatting
    - Format validator registry

### Go Generator Work
- [ ] Generate validation functions
  - [ ] Create `Validate()` method for each model
  - [ ] Check required fields
  - [ ] Check constraints (min/max, pattern, etc.)
  - [ ] Check format (email, uri, etc.)
  - [ ] Check enum values
  - [ ] Validate nested objects and arrays
- [ ] Client-side validation
  - [ ] Add optional pre-request validation
  - [ ] Validate request bodies before sending
  - [ ] Validate responses (for testing/debugging)
- [ ] Server-side validation
  - [ ] Auto-validate incoming requests
  - [ ] Return 400 with validation errors
  - [ ] Add validation middleware option
- [ ] Testing
  - [ ] Test validation with valid data
  - [ ] Test validation with invalid data
  - [ ] Test all constraint types
  - [ ] Test format validators

### Deliverables
- [ ] Validation utilities for all generators
- [ ] Go models include Validate() methods
- [ ] Client/server validate requests
- [ ] Clear validation error messages

---

## Milestone 2.7: Framework Support
**Duration**: 6-8 days
**Goal**: Support multiple Go HTTP frameworks

### Core Work
- [ ] **Shared Tooling**: Framework abstraction
  - [ ] `pkg/codegen/framework.go` - Framework utilities
    - Framework detection and selection
    - Framework-specific template overrides
    - Router pattern translation
    - Middleware pattern translation

### Go Generator Work
- [ ] Framework configuration
  - [ ] Add framework selection to config.yaml
  - [ ] Define framework capabilities
  - [ ] Set up template override system
- [ ] Implement framework support:
  - [ ] **net/http** (standard library - already done)
  - [ ] **Chi** router
    - [ ] Chi-specific route registration
    - [ ] Chi middleware integration
    - [ ] Path parameter extraction
  - [ ] **Echo** framework
    - [ ] Echo route registration
    - [ ] Echo middleware
    - [ ] Echo context handling
  - [ ] **Gin** framework
    - [ ] Gin route registration
    - [ ] Gin middleware
    - [ ] Gin context handling
- [ ] Template organization
  - [ ] `frameworks/net-http/` - standard library templates
  - [ ] `frameworks/chi/` - Chi templates
  - [ ] `frameworks/echo/` - Echo templates
  - [ ] `frameworks/gin/` - Gin templates
- [ ] Update template loading to check framework overrides
- [ ] Testing
  - [ ] Generate server for each framework
  - [ ] Verify compilation
  - [ ] Test basic server operations

### Deliverables
- [ ] Framework abstraction for all generators
- [ ] Go generator supports 4 frameworks
- [ ] Template override system working
- [ ] Framework selection documented

---

## Milestone 2.8: Advanced OpenAPI Features
**Duration**: 5-7 days
**Goal**: Implement remaining OpenAPI 3.x features

### Core Work
- [ ] **OpenAPI Feature**: Callbacks (webhooks)
  - [ ] Add `Callback` to core types
  - [ ] Parse callback definitions
  - [ ] Treat callbacks as reverse endpoints
- [ ] **OpenAPI Feature**: Links
  - [ ] Add `Link` to core types
  - [ ] Parse link definitions
  - [ ] Extract parameter mappings
- [ ] **OpenAPI Feature**: Cookie parameters
  - [ ] Already in `Parameter.In`, enhance processing
  - [ ] Add cookie handling to endpoint extraction
- [ ] **OpenAPI Feature**: Server variables
  - [ ] Parse server variable definitions
  - [ ] Support templated server URLs
- [ ] **OpenAPI Feature**: External documentation
  - [ ] Parse externalDocs
  - [ ] Include in comments/documentation
- [ ] **OpenAPI Feature**: Examples (plural)
  - [ ] Parse multiple examples
  - [ ] Include in generated tests
- [ ] **OpenAPI Feature**: Encoding objects
  - [ ] Parse encoding for multipart/form-data
  - [ ] Support style and explode for arrays/objects

### Go Generator Work
- [ ] **Callbacks**: Generate webhook handlers
  - [ ] Generate callback interface
  - [ ] Generate registration functions
  - [ ] Generate webhook signature verification helpers
- [ ] **Links**: Generate navigation helpers
  - [ ] Generate link builder functions
  - [ ] Parameter substitution helpers
- [ ] **Cookie parameters**: Add cookie support
  - [ ] Client: Set-Cookie header handling
  - [ ] Server: Cookie reading/writing
- [ ] **Server variables**: Support URL templating
  - [ ] Client: Allow server variable configuration
  - [ ] Generate variable validation
- [ ] **Examples**: Generate example tests
  - [ ] Use examples in test generation
  - [ ] Create example request/response fixtures
- [ ] Testing
  - [ ] Test specs with each new feature
  - [ ] Verify generated code compiles
  - [ ] Test runtime behavior

### Deliverables
- [ ] All major OpenAPI 3.x features supported in core
- [ ] Go generator implements all features
- [ ] Comprehensive test coverage
- [ ] Updated documentation

---

## Milestone 2.9: Code Quality & Developer Experience
**Duration**: 5-7 days
**Goal**: Improve generated code quality and usability

### Shared Tooling
- [ ] **Code formatting integration**
  - [ ] `pkg/codegen/format.go` - Post-generation formatting
    - Language-specific formatter registry (gofmt, prettier, black)
    - Format detection and execution
    - Error handling for formatter failures
- [ ] **Documentation generation**
  - [ ] `pkg/codegen/docs.go` - Documentation helpers
    - Extract OpenAPI descriptions
    - Generate language-specific doc comments
    - Link to OpenAPI spec sections
    - Generate markdown documentation
- [ ] **Testing utilities**
  - [ ] `pkg/codegen/testing.go` - Test generation helpers
    - Mock generation patterns
    - Test fixture creation
    - Example-based test generation

### Go Generator Work
- [ ] Code formatting
  - [ ] Auto-run gofmt on generated code
  - [ ] Auto-run goimports for import organization
  - [ ] Handle formatting errors gracefully
- [ ] Documentation improvements
  - [ ] Better godoc comments from OpenAPI descriptions
  - [ ] Add examples to comments
  - [ ] Link to OpenAPI spec in comments
  - [ ] Generate README for generated code
- [ ] Test generation
  - [ ] Generate compilation tests (already done)
  - [ ] Generate example-based tests
  - [ ] Generate mock implementations for testing
  - [ ] Generate test fixtures
- [ ] Code organization
  - [ ] Add file headers with generation info
  - [ ] Group related code sections
  - [ ] Add section comments
  - [ ] Optimize imports
- [ ] Developer ergonomics
  - [ ] Add builder patterns for complex types
  - [ ] Add convenience constructors
  - [ ] Add fluent interfaces where appropriate
  - [ ] Generate example usage code

### Deliverables
- [ ] Formatting utilities for all generators
- [ ] Documentation generation utilities
- [ ] Test generation utilities
- [ ] Go generator produces high-quality code
- [ ] Generated code is well-documented
- [ ] Generated code includes tests and examples

---

## Milestone 2.10: Advanced Client Features
**Duration**: 4-6 days
**Goal**: Add production-ready client features

### Shared Tooling
- [ ] **Retry and resilience patterns**
  - [ ] `pkg/codegen/resilience.go` - Resilience helpers
    - Retry strategy definitions
    - Circuit breaker patterns
    - Timeout configuration
    - Rate limiting patterns

### Go Generator Work
- [ ] Client configuration enhancements
  - [ ] Add timeout configuration
  - [ ] Add retry policies
  - [ ] Add circuit breaker support
  - [ ] Add request/response interceptors
- [ ] Advanced client options
  - [ ] Custom HTTP client injection
  - [ ] Request middleware chain
  - [ ] Response middleware chain
  - [ ] Logging/tracing hooks
- [ ] Error handling improvements
  - [ ] Automatic retry on transient errors
  - [ ] Exponential backoff
  - [ ] Retry budget
  - [ ] Detailed error context
- [ ] Performance optimizations
  - [ ] Connection pooling configuration
  - [ ] Keep-alive settings
  - [ ] Request/response streaming for large payloads
  - [ ] Compression support
- [ ] Testing helpers
  - [ ] Generate mock server for testing
  - [ ] Generate test client with recording
  - [ ] Generate fixture loader
- [ ] Testing
  - [ ] Test retry behavior
  - [ ] Test timeout handling
  - [ ] Test interceptor chains
  - [ ] Test performance features

### Deliverables
- [ ] Resilience utilities for all generators
- [ ] Production-ready Go client
- [ ] Comprehensive client configuration
- [ ] Testing helpers for client users

---

## Milestone 2.11: Advanced Server Features
**Duration**: 4-6 days
**Goal**: Add production-ready server features

### Shared Tooling
- [ ] **Middleware patterns**
  - [ ] `pkg/codegen/middleware.go` - Middleware helpers
    - Middleware chain patterns
    - Common middleware types (logging, metrics, cors)
    - Middleware composition

### Go Generator Work
- [ ] Server configuration enhancements
  - [ ] Add server startup/shutdown
  - [ ] Add graceful shutdown support
  - [ ] Add health check endpoints
  - [ ] Add metrics endpoints
- [ ] Middleware generation
  - [ ] Logging middleware
  - [ ] Request ID middleware
  - [ ] CORS middleware
  - [ ] Rate limiting middleware
  - [ ] Metrics middleware (Prometheus)
  - [ ] Recovery middleware (panic handling)
- [ ] Advanced routing
  - [ ] Route grouping by tag
  - [ ] Route prefix support
  - [ ] Versioned API support
  - [ ] Conditional route registration
- [ ] Server implementations
  - [ ] Generate main.go example
  - [ ] Configuration loading
  - [ ] Dependency injection setup
  - [ ] Example handler implementations
- [ ] Observability
  - [ ] Structured logging integration
  - [ ] OpenTelemetry tracing support
  - [ ] Prometheus metrics
  - [ ] Health check implementations
- [ ] Testing
  - [ ] Test server startup/shutdown
  - [ ] Test middleware chain
  - [ ] Test graceful shutdown
  - [ ] Integration tests

### Deliverables
- [ ] Middleware utilities for all generators
- [ ] Production-ready Go server
- [ ] Observability built-in
- [ ] Complete server examples

---

## Milestone 2.12: Configuration & Customization
**Duration**: 3-5 days
**Goal**: Flexible generator configuration

### Shared Tooling
- [ ] **Configuration management**
  - [ ] `pkg/config/loader.go` - Enhanced config loading
    - Multiple config sources (file, env, CLI)
    - Config merging and precedence
    - Config validation
    - Schema definition for configs
- [ ] **Template customization**
  - [ ] `pkg/template/override.go` - Template override system
    - User template directories
    - Template inheritance
    - Partial template includes
    - Template validation

### Go Generator Work
- [ ] Enhanced configuration
  - [ ] Add more type mapping options
  - [ ] Add naming customization per-model
  - [ ] Add generation exclusions/inclusions
  - [ ] Add template customization paths
- [ ] Template customization
  - [ ] Support user template overrides
  - [ ] Support template includes/partials
  - [ ] Document template variables
  - [ ] Provide template examples
- [ ] Generation options
  - [ ] Select what to generate (models only, client only, etc.)
  - [ ] Tag-based filtering
  - [ ] Path-based filtering
  - [ ] Custom package naming
- [ ] CLI enhancements
  - [ ] Better config file support
  - [ ] Config validation command
  - [ ] Template listing command
  - [ ] Dry-run mode
- [ ] Documentation
  - [ ] Configuration reference
  - [ ] Template customization guide
  - [ ] Examples of common customizations

### Deliverables
- [ ] Flexible config system for all generators
- [ ] Template override system
- [ ] Go generator fully configurable
- [ ] Comprehensive customization docs

---

## Milestone 2.13: Integration & Polish
**Duration**: 5-7 days
**Goal**: End-to-end testing and refinement

### All Areas
- [ ] End-to-end testing
  - [ ] Test complete workflow with multiple specs
  - [ ] Test all OpenAPI features
  - [ ] Test all frameworks
  - [ ] Test all configuration options
  - [ ] Performance testing
- [ ] Bug fixes
  - [ ] Address edge cases
  - [ ] Fix template issues
  - [ ] Improve error messages
  - [ ] Handle malformed specs gracefully
- [ ] Documentation
  - [ ] Complete API documentation
  - [ ] Tutorial for Go generator
  - [ ] Migration guide from openapi-generator
  - [ ] Best practices guide
  - [ ] Troubleshooting guide
- [ ] Examples
  - [ ] Complete petstore example
  - [ ] Real-world API example (GitHub, Stripe, etc.)
  - [ ] Microservices example
  - [ ] Full-stack example app
- [ ] Performance optimization
  - [ ] Profile generation performance
  - [ ] Optimize template rendering
  - [ ] Optimize normalizer
  - [ ] Memory usage optimization
- [ ] Developer experience
  - [ ] Better CLI output
  - [ ] Progress indicators
  - [ ] Colored output
  - [ ] Interactive mode
  - [ ] IDE integration helpers

### Deliverables
- [ ] Fully tested Go generator
- [ ] Complete documentation
- [ ] Example projects
- [ ] Performance optimized
- [ ] Production ready

---

## Summary: Shared Tooling Packages

By the end of Phase 2, we'll have these reusable packages:

```
pkg/codegen/
├── imports.go      # Import path resolution and management
├── naming.go       # Naming conventions and transformations
├── http.go         # HTTP method/status code utilities
├── content.go      # Media type and content negotiation
├── errors.go       # Error modeling and handling patterns
├── security.go     # Authentication/authorization patterns
├── validation.go   # Validation expression generation
├── framework.go    # Framework abstraction and selection
├── format.go       # Code formatting integration
├── docs.go         # Documentation generation
├── testing.go      # Test generation utilities
├── resilience.go   # Retry, timeout, circuit breaker patterns
└── middleware.go   # Middleware composition patterns

pkg/config/
└── loader.go       # Enhanced configuration loading

pkg/template/
└── override.go     # Template override and customization
```

These packages will be usable by TypeScript, Python, and future generators.

---

## Phase 2 Success Criteria

- [ ] All major OpenAPI 3.x features supported
- [ ] Complete Go client generation
- [ ] Complete Go server generation
- [ ] Multiple framework support (net/http, Chi, Echo, Gin)
- [ ] Shared tooling packages extracted and documented
- [ ] Production-ready code generation
- [ ] Comprehensive test coverage
- [ ] Complete documentation
- [ ] Example projects
- [ ] Ready for Phase 3 (additional languages)

---

## Estimated Total Duration: 10-14 weeks

This aggressive timeline assumes:
- Focused development time
- Parallel work on shared tooling and generator features
- Learning and iteration built in
- Regular testing and validation
