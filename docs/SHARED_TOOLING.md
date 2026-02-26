# Shared Tooling Architecture

This document describes the cross-cutting concerns and shared utilities that all generators can use. These packages are located in `pkg/codegen/` and `pkg/config/`.

## Philosophy

**80/20 Rule**: Most generator needs (80%) can be handled by shared tooling. Language-specific logic (20%) belongs in individual generators.

**Separation of Concerns**:
- **Core** (`pkg/core/`): OpenAPI spec representation and interfaces
- **Parser** (`pkg/parser/`): OpenAPI spec parsing and normalization
- **Codegen** (`pkg/codegen/`): Cross-language code generation utilities
- **Template** (`pkg/template/`): Template engine and customization
- **Generators** (`generators/*/`): Language-specific implementations

---

## Package: `pkg/codegen/imports`

**Purpose**: Import path resolution and management across languages

### Responsibilities
- Detect module/package systems (go.mod, package.json, setup.py, etc.)
- Construct import paths (relative vs absolute)
- Resolve cross-package references
- Deduplicate and sort imports
- Handle import aliases and naming conflicts

### Interface
```go
type ImportManager interface {
    // Add an import and return the local name to use
    AddImport(packagePath string) string

    // Add an import with specific alias
    AddImportAs(packagePath, alias string) string

    // Get all imports sorted and deduplicated
    GetImports() []Import

    // Resolve import path from current package to target package
    ResolvePath(from, to string) string
}

type Import struct {
    Path  string  // Full import path
    Alias string  // Local alias (empty if no alias)
}
```

### Usage Example
```go
mgr := imports.NewManager("myproject/generated/client")
mgr.AddImport("context")
mgr.AddImport("myproject/generated/models")
mgr.AddImport("net/http")

// Template data
data.Imports = mgr.GetImports()
```

### Language Support
- **Go**: Module-aware, handles replace directives
- **TypeScript**: Relative paths, barrel exports, path mapping
- **Python**: Absolute and relative imports, __init__.py awareness
- **Java**: Package paths, import grouping

---

## Package: `pkg/codegen/naming`

**Purpose**: Name transformation and convention handling

### Responsibilities
- Case conversion (camelCase, PascalCase, snake_case, kebab-case)
- Pluralization and singularization
- Reserved word handling per language
- Keyword collision avoidance
- Special character handling

### Interface
```go
type NamingConverter interface {
    // Convert to language-specific naming convention
    ToTypeName(s string) string      // Go: PascalCase, TS: PascalCase, Python: PascalCase
    ToFieldName(s string) string     // Go: PascalCase, TS: camelCase, Python: snake_case
    ToFunctionName(s string) string  // Go: PascalCase, TS: camelCase, Python: snake_case
    ToVariableName(s string) string  // Go: camelCase, TS: camelCase, Python: snake_case
    ToConstName(s string) string     // Go: PascalCase, TS: UPPER_CASE, Python: UPPER_CASE

    // Utilities
    Pluralize(s string) string
    Singularize(s string) string
    IsReservedWord(s string) bool
    EscapeReserved(s string) string  // Add suffix/prefix if reserved
}
```

### Usage Example
```go
conv := naming.NewGoConverter()
typeName := conv.ToTypeName("user_profile")    // "UserProfile"
fieldName := conv.ToFieldName("first_name")    // "FirstName"
varName := conv.ToVariableName("API_KEY")      // "apiKey"

// Handle reserved words
name := conv.ToVariableName("type")            // "type_" (Go reserved)
```

### Built-in Converters
- `naming.NewGoConverter()`
- `naming.NewTypeScriptConverter()`
- `naming.NewPythonConverter()`
- `naming.NewJavaConverter()`

---

## Package: `pkg/codegen/http`

**Purpose**: HTTP protocol utilities and constants

### Responsibilities
- HTTP method normalization and validation
- Status code categorization
- Content-Type parsing and matching
- Common header definitions
- URL encoding/decoding helpers

### Interface
```go
// Method utilities
func NormalizeMethod(method string) string              // "get" -> "GET"
func IsValidMethod(method string) bool
func IsSafeMethod(method string) bool                   // GET, HEAD, OPTIONS
func IsIdempotent(method string) bool                   // GET, PUT, DELETE, etc.

// Status code utilities
func IsSuccess(code int) bool                           // 2xx
func IsRedirect(code int) bool                          // 3xx
func IsClientError(code int) bool                       // 4xx
func IsServerError(code int) bool                       // 5xx
func StatusCategory(code int) string                    // "success", "error", etc.

// Content-Type utilities
func ParseMediaType(contentType string) (MediaType, error)
func MatchesMediaType(accept, contentType string) bool
func PreferredMediaType(accept []string, available []string) string

// Common constants
const (
    MethodGet    = "GET"
    MethodPost   = "POST"
    MethodPut    = "PUT"
    MethodDelete = "DELETE"
    MethodPatch  = "PATCH"
    // ... etc

    StatusOK                  = 200
    StatusCreated             = 201
    StatusBadRequest          = 400
    StatusUnauthorized        = 401
    StatusNotFound            = 404
    StatusInternalServerError = 500
    // ... etc

    ContentTypeJSON = "application/json"
    ContentTypeXML  = "application/xml"
    ContentTypeForm = "application/x-www-form-urlencoded"
    // ... etc
)
```

---

## Package: `pkg/codegen/content`

**Purpose**: Content negotiation and media type handling

### Responsibilities
- Media type parsing and comparison
- Content negotiation (Accept header processing)
- Serialization format detection
- Charset and encoding handling
- Quality value parsing (q=0.8)

### Interface
```go
type MediaType struct {
    Type       string            // "application"
    Subtype    string            // "json"
    Parameters map[string]string // "charset": "utf-8"
    Quality    float64           // From q parameter
}

// Content negotiation
func ParseAccept(accept string) []MediaType
func NegotiateContent(accept string, available []string) (string, bool)
func QualitySort(mediaTypes []MediaType) []MediaType

// Format detection
func IsJSON(contentType string) bool
func IsXML(contentType string) bool
func IsFormData(contentType string) bool
func IsMultipart(contentType string) bool

// Serialization helpers
func NeedsBodySerialization(contentType string) bool
func GetDefaultSerializer(contentType string) string  // "json", "xml", etc.
```

---

## Package: `pkg/codegen/errors`

**Purpose**: Error modeling and handling patterns

### Responsibilities
- Error type naming conventions
- Status code to error type mapping
- Error hierarchy construction
- Error message formatting
- Error context preservation

### Interface
```go
// Error naming
func ErrorTypeName(statusCode int, operation string) string
func ErrorConstName(statusCode int) string

// Error categorization
func IsTransientError(statusCode int) bool  // Should retry
func IsPermanentError(statusCode int) bool  // Don't retry
func IsAuthError(statusCode int) bool       // 401, 403

// Error context
type ErrorContext struct {
    StatusCode int
    Message    string
    Details    map[string]interface{}
    Cause      error
}

func FormatErrorContext(ctx ErrorContext) string
```

### Language-Specific Patterns
- **Go**: Error interface, wrapped errors, typed errors
- **TypeScript**: Error subclasses, discriminated unions
- **Python**: Exception hierarchy
- **Java**: Checked vs unchecked exceptions

---

## Package: `pkg/codegen/security`

**Purpose**: Authentication and authorization patterns

### Responsibilities
- Security scheme identification
- Auth token injection patterns
- Credential management patterns
- Scope/permission checking
- Common auth middleware patterns

### Interface
```go
type SecurityScheme struct {
    Type   string  // "apiKey", "http", "oauth2", "openIdConnect"
    In     string  // "header", "query", "cookie"
    Name   string  // Header/query/cookie name
    Scheme string  // "bearer", "basic"
}

// Pattern generation
func GetAuthInjectionPattern(scheme SecurityScheme, lang string) string
func GetAuthVerificationPattern(scheme SecurityScheme, lang string) string
func RequiresMiddleware(scheme SecurityScheme) bool

// Common patterns
const (
    PatternAPIKey       = "api_key"
    PatternBearer       = "bearer"
    PatternBasic        = "basic"
    PatternOAuth2       = "oauth2"
    PatternOpenIDConnect = "openid_connect"
)
```

---

## Package: `pkg/codegen/validation`

**Purpose**: Validation code generation

### Responsibilities
- Generate validation expressions
- Common regex patterns
- Format validators
- Constraint checking
- Validation error formatting

### Interface
```go
type Validator interface {
    // Generate validation code
    ValidateRequired(fieldName string) string
    ValidatePattern(fieldName, pattern string) string
    ValidateMinMax(fieldName string, min, max *float64) string
    ValidateLength(fieldName string, minLen, maxLen *int) string
    ValidateEnum(fieldName string, values []interface{}) string
    ValidateFormat(fieldName, format string) string
}

// Common formats
const (
    FormatEmail    = "email"
    FormatURI      = "uri"
    FormatUUID     = "uuid"
    FormatDateTime = "date-time"
    FormatDate     = "date"
    FormatIPv4     = "ipv4"
    FormatIPv6     = "ipv6"
)

// Regex patterns
var CommonPatterns = map[string]string{
    FormatEmail: `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`,
    FormatUUID:  `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`,
    // ... etc
}
```

### Language Implementations
- **Go**: `naming.NewGoValidator()` - generates Go validation code
- **TypeScript**: `naming.NewTSValidator()` - generates TS validation code
- **Python**: `naming.NewPyValidator()` - generates Python validation code

---

## Package: `pkg/codegen/framework`

**Purpose**: Framework abstraction and template routing

### Responsibilities
- Framework detection and selection
- Framework capability queries
- Template override resolution
- Router pattern translation
- Middleware pattern translation

### Interface
```go
type Framework interface {
    Name() string
    Language() string

    // Capabilities
    SupportsMiddleware() bool
    SupportsWebSockets() bool
    SupportsStreaming() bool

    // Template resolution
    GetTemplatePath(baseTemplate string) string  // Returns override or base

    // Pattern translation
    TranslatePathPattern(openAPIPath string) string  // "/users/{id}" -> framework format
    TranslateMiddlewarePattern(middleware string) string
}

// Framework registry
func RegisterFramework(f Framework)
func GetFramework(language, name string) (Framework, error)
func ListFrameworks(language string) []Framework
```

### Example Frameworks
- **Go**: `net/http`, `chi`, `echo`, `gin`, `fiber`
- **TypeScript**: `express`, `fastify`, `koa`, `nest`
- **Python**: `flask`, `fastapi`, `django`, `starlette`
- **Java**: `spring`, `micronaut`, `quarkus`

---

## Package: `pkg/codegen/format`

**Purpose**: Post-generation code formatting

### Responsibilities
- Detect available formatters
- Execute language-specific formatters
- Handle formatter errors gracefully
- Support custom formatter configuration

### Interface
```go
type Formatter interface {
    Name() string
    Language() string

    // Format code
    Format(code []byte) ([]byte, error)
    FormatFile(path string) error

    // Availability
    IsAvailable() bool
    Install() error  // Installation instructions
}

// Built-in formatters
func NewGoFormatter() Formatter       // gofmt + goimports
func NewTSFormatter() Formatter       // prettier
func NewPyFormatter() Formatter       // black
func NewJavaFormatter() Formatter     // google-java-format

// Formatter registry
func RegisterFormatter(f Formatter)
func GetFormatter(language string) (Formatter, error)
func FormatDirectory(dir, language string) error
```

---

## Package: `pkg/codegen/docs`

**Purpose**: Documentation generation

### Responsibilities
- Extract documentation from OpenAPI
- Generate language-specific doc comments
- Link to OpenAPI spec sections
- Generate markdown documentation
- Generate API reference

### Interface
```go
type DocGenerator interface {
    // Generate doc comments for code
    TypeDoc(model core.Model) string
    FieldDoc(prop core.Property) string
    FunctionDoc(endpoint core.Endpoint) string

    // Generate documentation files
    GenerateMarkdown(spec *core.Spec) ([]byte, error)
    GenerateHTML(spec *core.Spec) ([]byte, error)
    GenerateReference(spec *core.Spec) ([]byte, error)
}

// Language-specific doc generators
func NewGoDocGenerator() DocGenerator      // Godoc format
func NewTSDocGenerator() DocGenerator      // TSDoc/JSDoc format
func NewPyDocGenerator() DocGenerator      // Docstring format
func NewJavaDocGenerator() DocGenerator    // Javadoc format
```

---

## Package: `pkg/codegen/testing`

**Purpose**: Test generation utilities

### Responsibilities
- Generate test fixtures
- Generate mock implementations
- Generate example-based tests
- Generate integration tests
- Test data generation

### Interface
```go
type TestGenerator interface {
    // Generate test code
    GenerateModelTests(model core.Model) string
    GenerateClientTests(endpoint core.Endpoint) string
    GenerateServerTests(endpoint core.Endpoint) string

    // Generate mocks
    GenerateMockClient(spec *core.Spec) string
    GenerateMockServer(spec *core.Spec) string

    // Generate fixtures
    GenerateFixtures(spec *core.Spec) map[string]interface{}
}

// Test data generation
func GenerateTestValue(prop core.Property) interface{}
func GenerateExampleRequest(endpoint core.Endpoint) interface{}
func GenerateExampleResponse(endpoint core.Endpoint) interface{}
```

---

## Package: `pkg/codegen/resilience`

**Purpose**: Retry, timeout, and circuit breaker patterns

### Responsibilities
- Retry strategy definitions
- Circuit breaker patterns
- Timeout configuration
- Rate limiting patterns
- Backoff strategies

### Interface
```go
type RetryStrategy struct {
    MaxRetries  int
    BackoffType string  // "exponential", "linear", "constant"
    InitialWait time.Duration
    MaxWait     time.Duration
    Jitter      bool
}

type CircuitBreakerConfig struct {
    FailureThreshold int
    SuccessThreshold int
    Timeout          time.Duration
    HalfOpenRequests int
}

// Pattern generation
func GenerateRetryCode(strategy RetryStrategy, lang string) string
func GenerateCircuitBreakerCode(config CircuitBreakerConfig, lang string) string
func GenerateTimeoutCode(timeout time.Duration, lang string) string
func GenerateRateLimitCode(rps int, lang string) string
```

---

## Package: `pkg/codegen/middleware`

**Purpose**: Middleware composition patterns

### Responsibilities
- Middleware chain patterns
- Common middleware implementations
- Middleware ordering
- Framework-specific translation

### Interface
```go
type Middleware struct {
    Name     string
    Priority int  // Lower = earlier in chain
    Config   map[string]interface{}
}

// Common middleware types
const (
    MiddlewareLogging     = "logging"
    MiddlewareRequestID   = "request_id"
    MiddlewareCORS        = "cors"
    MiddlewareRateLimit   = "rate_limit"
    MiddlewareAuth        = "auth"
    MiddlewareMetrics     = "metrics"
    MiddlewareRecovery    = "recovery"
    MiddlewareValidation  = "validation"
)

// Pattern generation
func GenerateMiddlewareChain(middlewares []Middleware, framework string) string
func GenerateMiddlewareImpl(mw Middleware, framework string) string
```

---

## Package: `pkg/config/loader`

**Purpose**: Enhanced configuration management

### Responsibilities
- Load from multiple sources (file, env, CLI)
- Merge configurations with precedence
- Validate configurations
- Schema definition and validation
- Default value handling

### Interface
```go
type Config struct {
    // Multiple sources
    File map[string]interface{}
    Env  map[string]interface{}
    CLI  map[string]interface{}

    // Merged result
    Merged map[string]interface{}
}

type Loader interface {
    // Load from sources
    LoadFile(path string) error
    LoadEnv(prefix string) error
    LoadCLI(flags map[string]interface{}) error

    // Merge and validate
    Merge() error
    Validate(schema Schema) error

    // Access
    Get(key string) interface{}
    GetString(key string) string
    GetInt(key string) int
    GetBool(key string) bool
}

// Precedence: CLI > Env > File > Defaults
```

---

## Package: `pkg/template/override`

**Purpose**: Template customization and override system

### Responsibilities
- User template directory support
- Template inheritance
- Partial template includes
- Template validation
- Template documentation

### Interface
```go
type TemplateResolver interface {
    // Resolve template path with overrides
    Resolve(templateName string) (string, error)

    // Load template with overrides
    Load(templateName string) (*template.Template, error)

    // Add override directories
    AddOverrideDir(dir string)

    // List available templates
    ListTemplates() []string
}

// Template includes
func Include(templateName string, data interface{}) string
func Partial(partialName string, data interface{}) string
```

---

## Usage in Generators

Generators use shared tooling like this:

```go
// generators/golang/generator.go

import (
    "github.com/yourorg/beaverspec/pkg/codegen/imports"
    "github.com/yourorg/beaverspec/pkg/codegen/naming"
    "github.com/yourorg/beaverspec/pkg/codegen/http"
)

type Generator struct {
    naming    naming.NamingConverter
    imports   *imports.Manager
    config    *Config
}

func (g *Generator) Generate(spec *core.Spec, config *core.Config) error {
    // Use shared utilities
    typeName := g.naming.ToTypeName(model.Name)

    // Manage imports
    g.imports.AddImport("context")
    g.imports.AddImport("net/http")

    // HTTP utilities
    method := http.NormalizeMethod(endpoint.Method)
    isSuccess := http.IsSuccess(response.StatusCode)

    // ... rest of generation
}
```

---

## Benefits of Shared Tooling

1. **Consistency**: All generators handle common concerns the same way
2. **Maintainability**: Fix a bug once, all generators benefit
3. **Extensibility**: New generators can focus on language-specific logic
4. **Testing**: Test shared utilities once, comprehensively
5. **Documentation**: Document patterns once, reference everywhere
6. **Best Practices**: Encode best practices in shared code

---

## Adding New Shared Utilities

When adding language-specific logic to a generator, ask:

1. **Is this pattern common across languages?** → Consider shared tooling
2. **Would other generators need this?** → Consider shared tooling
3. **Is this purely language syntax?** → Keep in generator/templates
4. **Is this OpenAPI-related?** → Consider if it belongs in core/parser

**Decision Tree**:
```
Is it language-specific syntax/idiom?
├─ Yes → Keep in generator
└─ No → Is it common across languages?
    ├─ Yes → Add to pkg/codegen
    └─ No → Is it OpenAPI spec related?
        ├─ Yes → Add to pkg/core or pkg/parser
        └─ No → Keep in generator (may become shared later)
```
