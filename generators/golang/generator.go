package golang

import (
	"fmt"
	"go/format"
	"sort"
	"strings"

	"github.com/imerljak/beaverspec/pkg/codegen"
	"github.com/imerljak/beaverspec/pkg/core"
	"github.com/imerljak/beaverspec/pkg/template"
)

type Generator struct {
	//
}

// ModelData is the data passed to model templates
type ModelData struct {
	Name          string
	Description   string
	Properties    []PropertyData
	EmbeddedTypes []string // types to embed (from allOf with $ref)
	IsArray       bool     // true for array type schemas
	ArrayType     string   // the type for array items
	IsMap         bool     // true for map type schemas
	MapValueType  string   // the type for map values
	HasValidation bool     // true if any property has a validation constraint

	// Union types (oneOf/anyOf)
	IsOneOf       bool
	IsAnyOf       bool
	UnionVariants []string
	Discriminator *DiscriminatorData
}

type DiscriminatorData struct {
	PropertyName string            // e.g., "petType"
	Mapping      map[string]string //e.g., { "cat": "Cat", "dog": "Dog" }
}

type PropertyData struct {
	Name, Type, Description, JsonTag string

	// Metadata
	Default    string // Formatted default value
	ReadOnly   bool
	WriteOnly  bool
	Deprecated bool

	// Validation constraints
	Required          bool
	ZeroCheckExpr     string // Go expression for zero-value check (e.g. `== ""`, `== 0`, `== nil`); empty means skip
	Format            string
	MinLength         *int
	MaxLength         *int
	Minimum           *float64
	Maximum           *float64
	Pattern           string
	EnumValues        []string
	HasEnum           bool
	HasValidation     bool // true if any constraint is present
	IsFormatValidated bool // true if format requires a format check
}

type EnumData struct {
	TypeName    string      // e.g., "OrderStatus"
	BaseType    string      // e.g., "string"
	Description string      // Description of the enum
	Values      []EnumValue // The enum constants
}

type EnumValue struct {
	Name  string      // e.g., "OrderStatusPending"
	Value interface{} // e.g., "pending"
}

type TemplateData struct {
	PackageName string
	Imports     []string
	Enums       []EnumData
	Models      []ModelData
}

// ClientData represents data for client generation
type ClientData struct {
	PackageName     string
	InterfaceName   string
	Operations      []OperationData
	Imports         []codegen.Import
	SecuritySchemes []SecuritySchemeData
}

// ServerData represents data for server generation
type ServerData struct {
	PackageName     string
	Imports         []codegen.Import
	ModelsPackage   string // fully-qualified models package path (e.g. "github.com/example/project/models")
	Handlers        []HandlerGroup
	SecuritySchemes []SecuritySchemeData
}

// SecuritySchemeData is a template-friendly representation of a security scheme
type SecuritySchemeData struct {
	Name        string
	Type        string
	In          string // for apiKey
	KeyName     string // header/query/cookie name for apiKey
	Scheme      string // "bearer" or "basic" for http
	IsAPIKey    bool
	IsBearer    bool
	IsBasicAuth bool
	IsOAuth2    bool
}

// HandlerGroup groups operations by tag for routing
type HandlerGroup struct {
	Name       string
	Operations []OperationData
}

// OperationData represents a single client/server operation
type OperationData struct {
	Name            string // Method name (e.g., "ListPets")
	OperationID     string
	Description     string
	Method          string // HTTP method
	Path            string // URL path
	PathParams      []ParamData
	QueryParams     []ParamData
	HeaderParams    []ParamData
	Security        []string // Security requirement names for this operation
	HasBody         bool
	BodyType        string // Type of request body
	Responses       []ResponseData
	ReturnType      string // Convenience pointer to success return type
	AcceptType      string // Convenience pointer to success accept type
	ErrorReturn     bool   // Whether to return error
	ContentType     string
	IsJSONBody      bool
	IsXMLBody       bool
	IsFormBody      bool
	IsMultipartBody bool
	IsJSONResponse  bool
	IsXMLResponse   bool
}

// ResponseData holds information over a mapped model response output
type ResponseData struct {
	StatusCode string
	IsSuccess  bool
	IsError    bool
	Type       string
	AcceptType string
	IsJSON     bool
	IsXML      bool
	HasHeaders bool
	Headers    []ParamData
}

// ParamData represents a parameter
type ParamData struct {
	Name        string
	Type        string
	Description string
	Required    bool
}

func NewGenerator() *Generator {
	return &Generator{}
}

// Metadata returns information about this generator
func (g *Generator) Metadata() core.GeneratorMetadata {
	return core.GeneratorMetadata{
		Name:        "go",
		Language:    "go",
		Version:     "0.1.0",
		Description: "Go code generator for OpenAPI specs",
		Author:      "openapi-gen",
		Capabilities: core.Capabilities{
			SupportsClient:     true,
			SupportsServer:     true,
			SupportsValidation: false, // Not yet
			SupportsAsync:      false,
			SupportsStreaming:  false,
			ClientFrameworks:   []string{"net/http"},
			ServerFrameworks:   []string{"net/http"},
		},
	}
}

// Validate checks if the spec is valid for this generator
func (g *Generator) Validate(spec *core.Spec, config *core.Config) []core.ValidationError {
	// For MVP, we'll just return no errors
	// In the future, add validation logic.
	return []core.ValidationError{}
}

// Generate produces Go code from the spec
func (g *Generator) Generate(spec *core.Spec, config *core.Config) (*core.GenerationResult, error) {

	packageName := "models"
	modulePath := "github.com/example/project" // default
	framework := "net-http"                    // default server framework
	if config.Options != nil {
		if pkgName, ok := config.Options["package"].(string); ok && pkgName != "" {
			packageName = pkgName
		}
		if mp, ok := config.Options["modulePath"].(string); ok && mp != "" {
			modulePath = mp
		}
		if fw, ok := config.Options["framework"].(string); ok && fw != "" {
			framework = fw
		}
	}
	baseDir := config.OutputDir // e.g., "generated"

	models := g.convertModels(spec.Models)

	templateData := TemplateData{
		PackageName: packageName,
		Models:      models,
		Imports:     g.collectImports(models),
		Enums:       g.collectEnums(spec.Models),
	}

	// Create template engine
	engine := template.NewGoTemplateEngine()

	// Load templates from the templates dir
	// For now, we need to know where templates are
	// Afterwards this will come from config/embed
	templatesPath := "./generators/golang/templates"

	if err := engine.Load(templatesPath); err != nil {
		return nil, fmt.Errorf("failed to load templates: %w", err)
	}

	// Render the models template
	content, err := engine.Render("models.go.tmpl", templateData)
	if err != nil {
		return nil, fmt.Errorf("failed to render models template: %w", err)
	}

	formattedContent, err := format.Source([]byte(content))
	if err != nil {
		return nil, fmt.Errorf("failed to format models.go: %w", err)
	}

	// Render the test template
	testContent, err := engine.Render("models_test.go.tmpl", templateData)
	if err != nil {
		return nil, fmt.Errorf("failed to render models template: %w", err)
	}

	formattedTestContent, err := format.Source([]byte(testContent))
	if err != nil {
		return nil, fmt.Errorf("failed to format models_test.go: %w", err)
	}

	files := []core.GeneratedFile{
		{
			Path:    "models/models.go",
			Content: []byte(formattedContent),
		},
		{
			Path:    "models/models_test.go",
			Content: []byte(formattedTestContent),
		},
	}

	// Generate client if there are endpoints
	if len(spec.Endpoints) > 0 {
		clientData := g.convertEndpointsToClient(spec.Endpoints, spec.SecurityDef, modulePath, baseDir)

		clientContent, err := engine.Render("client/client.go.tmpl", clientData)
		if err != nil {
			return nil, fmt.Errorf("failed to render client template: %w", err)
		}

		formattedClientContent, err := format.Source([]byte(clientContent))
		if err != nil {
			return nil, fmt.Errorf("failed to format client.go: %w", err)
		}

		files = append(files, core.GeneratedFile{
			Path:    "client/client.go",
			Content: formattedClientContent,
		})

		// Generate server logic
		serverData := g.convertEndpointsToServer(spec.Endpoints, spec.SecurityDef, modulePath, baseDir)

		// interface
		serverInterfaceContent, err := engine.Render("server/interface.go.tmpl", serverData)
		if err != nil {
			return nil, fmt.Errorf("failed to render server interface: %w", err)
		}
		formattedInterface, err := format.Source([]byte(serverInterfaceContent))
		if err != nil {
			return nil, fmt.Errorf("failed to format server/interface.go: %w\ncontent:\n%s", err, serverInterfaceContent)
		}
		files = append(files, core.GeneratedFile{Path: "server/interface.go", Content: formattedInterface})

		// mocks (shared, framework-agnostic)
		serverMocksContent, err := engine.Render("server/mocks.go.tmpl", serverData)
		if err != nil {
			return nil, fmt.Errorf("failed to render server mocks: %w", err)
		}
		formattedMocks, err := format.Source([]byte(serverMocksContent))
		if err != nil {
			return nil, fmt.Errorf("failed to format server/mocks.go: %w\ncontent:\n%s", err, serverMocksContent)
		}
		files = append(files, core.GeneratedFile{Path: "server/mocks.go", Content: formattedMocks})

		// handlers and routes are framework-specific (template layer)
		serverHandlersContent, err := engine.Render("server/"+framework+"/handlers.go.tmpl", serverData)
		if err != nil {
			return nil, fmt.Errorf("failed to render server handlers: %w", err)
		}
		formattedHandlers, err := format.Source([]byte(serverHandlersContent))
		if err == nil {
			files = append(files, core.GeneratedFile{Path: "server/handlers.go", Content: formattedHandlers})
		}

		serverRoutesContent, err := engine.Render("server/"+framework+"/routes.go.tmpl", serverData)
		if err != nil {
			return nil, fmt.Errorf("failed to render server routes: %w", err)
		}
		formattedRoutes, err := format.Source([]byte(serverRoutesContent))
		if err == nil {
			files = append(files, core.GeneratedFile{Path: "server/routes.go", Content: formattedRoutes})
		}
	}

	// Create the generated file
	result := &core.GenerationResult{
		Files:    files,
		Warnings: []core.GenerationError{},
	}

	return result, nil
}

// SupportedFeatures returns OpenAPI features this generator supports
func (g *Generator) SupportedFeatures() []core.Feature {
	return []core.Feature{
		core.FeatureNullable,
		// TODO: Add more as implemented
	}
}

// convertEndpointsToClient converts core endpoints to client data
func (g *Generator) convertEndpointsToClient(endpoints []core.Endpoint, securityDef []core.SecurityScheme, modulePath, baseDir string) ClientData {

	imports := codegen.NewImportManager(modulePath, baseDir, "client")
	imports.Add("context")
	imports.Add("net/http")
	imports.Add("fmt")
	imports.Add("io")

	operations := make([]OperationData, 0, len(endpoints))

	for _, ep := range endpoints {
		op := g.convertEndpoint(ep)
		operations = append(operations, op)

		if op.IsJSONBody {
			imports.Add("encoding/json")
		}
		if op.IsXMLBody {
			imports.Add("encoding/xml")
		}
		for _, resp := range op.Responses {
			if resp.IsJSON {
				imports.Add("encoding/json")
			}
			if resp.IsXML {
				imports.Add("encoding/xml")
			}
		}
		if op.IsFormBody || len(op.QueryParams) > 0 {
			imports.Add("net/url")
		}
		if op.IsFormBody {
			imports.Add("strings")
		}
	}

	imports.AddSibling("models")

	return ClientData{
		PackageName:     "client",
		InterfaceName:   "Client",
		Operations:      operations,
		Imports:         imports.GetImports(),
		SecuritySchemes: g.convertSecuritySchemes(securityDef),
	}
}

// convertEndpointsToServer converts core endpoints into server structures
func (g *Generator) convertEndpointsToServer(endpoints []core.Endpoint, securityDef []core.SecurityScheme, modulePath, baseDir string) ServerData {
	imports := codegen.NewImportManager(modulePath, baseDir, "server")
	imports.Add("context")
	imports.Add("net/http")
	imports.Add("errors")

	groupsMap := make(map[string][]OperationData)

	for _, ep := range endpoints {
		op := g.convertEndpoint(ep)

		// Group by first tag if present, else "Default"
		groupName := "Default"
		if len(ep.Tags) > 0 {
			groupName = template.ToPascalCase(ep.Tags[0])
		}

		groupsMap[groupName] = append(groupsMap[groupName], op)

		if op.IsJSONBody {
			imports.Add("encoding/json")
		}
		if op.IsXMLBody {
			imports.Add("encoding/xml")
		}
		for _, resp := range op.Responses {
			if resp.IsJSON {
				imports.Add("encoding/json")
			}
			if resp.IsXML {
				imports.Add("encoding/xml")
			}
		}
	}

	var handlers []HandlerGroup
	for name, ops := range groupsMap {
		handlers = append(handlers, HandlerGroup{
			Name:       name,
			Operations: ops,
		})
	}

	// Make output consistent
	sort.Slice(handlers, func(i, j int) bool {
		return handlers[i].Name < handlers[j].Name
	})

	imports.AddSibling("models")

	return ServerData{
		PackageName:     "server",
		Imports:         imports.GetImports(),
		ModelsPackage:   modulePath + "/models",
		Handlers:        handlers,
		SecuritySchemes: g.convertSecuritySchemes(securityDef),
	}
}

// convertSecuritySchemes converts core.SecurityScheme slice into template-friendly SecuritySchemeData
func (g *Generator) convertSecuritySchemes(schemes []core.SecurityScheme) []SecuritySchemeData {
	result := make([]SecuritySchemeData, 0, len(schemes))
	for _, s := range schemes {
		result = append(result, SecuritySchemeData{
			Name:        s.Name,
			Type:        s.Type,
			In:          s.In,
			KeyName:     s.Name, // For apiKey, Name doubles as parameter name
			Scheme:      s.Scheme,
			IsAPIKey:    codegen.IsAPIKey(s.Type),
			IsBearer:    codegen.IsHTTP(s.Type) && codegen.IsBearer(s.Scheme),
			IsBasicAuth: codegen.IsHTTP(s.Type) && codegen.IsBasicAuth(s.Scheme),
			IsOAuth2:    codegen.IsOAuth2(s.Type),
		})
	}
	return result
}

// convertEndpoint converts a single endpoint to operation data
func (g *Generator) convertEndpoint(ep core.Endpoint) OperationData {
	// Generate method name from operation ID
	methodName := g.toMethodName(ep.OperationID)

	// Extract parameters by location
	paramMap := g.extractParametersByLocation(ep.Parameters)

	// Determine body type and return type
	bodyType, contentType := g.getRequestBodyType(ep.RequestBody)
	responsesData := g.extractResponsesData(ep.Responses)

	var returnType, acceptType string
	var isJSONResponse, isXMLResponse bool
	for _, r := range responsesData {
		if r.IsSuccess && r.Type != "" {
			returnType = r.Type
			acceptType = r.AcceptType
			isJSONResponse = r.IsJSON
			isXMLResponse = r.IsXML
			break
		}
	}

	return OperationData{
		Name:            methodName,
		OperationID:     ep.OperationID,
		Description:     ep.Description,
		Method:          ep.Method,
		Path:            ep.Path,
		PathParams:      paramMap["path"],
		QueryParams:     paramMap["query"],
		HeaderParams:    paramMap["header"],
		HasBody:         bodyType != "",
		BodyType:        bodyType,
		Responses:       responsesData,
		ReturnType:      returnType,
		AcceptType:      acceptType,
		ErrorReturn:     true, // Always return error
		ContentType:     contentType,
		IsJSONBody:      codegen.IsJSON(contentType),
		IsXMLBody:       codegen.IsXML(contentType),
		IsFormBody:      codegen.IsFormURLEncoded(contentType),
		IsMultipartBody: codegen.IsMultipartForm(contentType),
		IsJSONResponse:  isJSONResponse,
		IsXMLResponse:   isXMLResponse,
	}
}

// extractResponsesData compiles parsed successful and error-typed responses and headers
func (g *Generator) extractResponsesData(responses []core.Response) []ResponseData {
	var respData []ResponseData

	for _, resp := range responses {
		isSuccess := codegen.IsSuccessStatus(resp.StatusCode)
		isError := codegen.IsErrorStatus(resp.StatusCode)

		var returnType string
		var acceptType string

		// Simplify by picking the first mapped type implementation matching JSON or XML
		for cType, mediaType := range resp.Content {
			if mediaType.Schema != nil {
				acceptType = cType
				if mediaType.Schema.RefType != "" {
					returnType = "*models." + mediaType.Schema.RefType // TODO: parameterize
				} else {
					returnType = g.mapParameterType(mediaType.Schema)
				}
				break // Only map the first matched type payload per status code for now MVP
			}
		}

		var headers []ParamData
		for name, header := range resp.Headers {
			headers = append(headers, ParamData{
				Name:        name,
				Type:        g.mapParameterType(header.Schema),
				Description: header.Description,
				Required:    header.Required,
			})
		}

		respData = append(respData, ResponseData{
			StatusCode: resp.StatusCode,
			IsSuccess:  isSuccess,
			IsError:    isError,
			Type:       returnType,
			AcceptType: acceptType,
			IsJSON:     codegen.IsJSON(acceptType),
			IsXML:      codegen.IsXML(acceptType),
			HasHeaders: len(headers) > 0,
			Headers:    headers,
		})
	}

	return respData
}

// getRequestBodyType determines the request body type
func (g *Generator) getRequestBodyType(body *core.RequestBody) (string, string) {
	if body == nil {
		return "", ""
	}

	for contentType, mediaType := range body.Content {
		if mediaType.Schema != nil {
			if mediaType.Schema.RefType != "" {
				return "*models." + mediaType.Schema.RefType, contentType // TODO: parameterize the package name
			}
			return g.mapParameterType(mediaType.Schema), contentType
		}
	}
	return "", ""
}

// extractParamsByLocation filters parameters by location
func (g *Generator) extractParametersByLocation(parameters []core.Parameter) map[string][]ParamData {
	result := make(map[string][]ParamData, 3)

	for _, p := range parameters {
		paramType := g.mapParameterType(p.Schema)

		if !p.Required && (p.In == "query" || p.In == "header") {
			paramType = "*" + paramType
		}

		if result[p.In] == nil {
			result[p.In] = make([]ParamData, 0, len(parameters))
		}

		result[p.In] = append(result[p.In], ParamData{
			Name:        p.Name,
			Type:        paramType,
			Description: p.Description,
			Required:    p.Required,
		})
	}

	return result
}

// mapParameterType maps a parameter schema to Go type
func (g *Generator) mapParameterType(schema *core.Property) string {
	if schema == nil {
		return "string"
	}

	// Handle references — qualify with models package
	if schema.RefType != "" {
		return "models." + schema.RefType
	}

	// Handle arrays
	if schema.Type == "array" && schema.Items != nil {
		itemType := g.mapParameterType(schema.Items)
		return "[]" + itemType
	}

	// Map primitive types
	return g.mapPrimitiveType(schema.Type, schema.Format)
}

// toMethodName converts an operation ID to a Go method name
func (g *Generator) toMethodName(operationID string) string {
	if operationID == "" {
		return "Execute"
	}

	// PascalCase the operationID
	return template.ToPascalCase(operationID)
}

func (g *Generator) convertModels(models []core.Model) []ModelData {
	var result []ModelData
	for _, m := range models {
		// Top-level scalar enum schemas are emitted in the Enums section, not as structs.
		if isTopLevelEnum(m) {
			continue
		}
		result = append(result, g.convertModel(m))
	}
	return result
}

// isTopLevelEnum returns true for models that represent a named scalar enum
// type (e.g. type: string + enum values) rather than a struct.
func isTopLevelEnum(m core.Model) bool {
	return len(m.Enum) > 0 && m.Type != "" && len(m.Properties) == 0 && !m.IsArray && !m.IsMap && !m.IsOneOf && !m.IsAnyOf
}

func (g *Generator) convertModel(m core.Model) ModelData {

	// Extract embedded types from allOf (only $ref, not inline schemas)
	var embeddedTypes []string
	for _, allOfModel := range m.AllOf {
		// Only embed if it's a reference (has a name but no properties)
		if allOfModel.Name != "" && len(allOfModel.Properties) == 0 {
			embeddedTypes = append(embeddedTypes, allOfModel.Name)
		} else if len(allOfModel.Properties) > 0 {
			// Inline schema properties should be merged into the main properties
			for _, prop := range allOfModel.Properties {
				m.Properties = append(m.Properties, prop)
			}
		}
	}

	// Extract union variants from oneOf/anyOf
	var unionVariants []string
	if m.IsOneOf {
		for _, variant := range m.OneOf {
			if variant.Name != "" {
				unionVariants = append(unionVariants, variant.Name)
			}
		}
	} else if m.IsAnyOf {
		for _, variant := range m.AnyOf {
			if variant.Name != "" {
				unionVariants = append(unionVariants, variant.Name)
			}
		}
	}

	var discriminator *DiscriminatorData
	if m.Discriminator != nil {
		// Extract type names from mapping (strip OpenAPI $ref paths)
		mapping := make(map[string]string)

		if len(m.Discriminator.Mapping) > 0 {
			// Explicit mapping provided
			for key, value := range m.Discriminator.Mapping {
				// value might be '#/components/schemas/Cat' or just 'Cat'
				typeName := value
				if strings.Contains(value, "/") {
					// Extract last part after final '/'
					typeName = value[strings.LastIndex(value, "/")+1:]
				}
				mapping[key] = typeName
			}
		} else {
			// Implicit mapping - use variant names (lowercase as discriminator value)
			for _, variant := range unionVariants {
				// Use lowercase variant name as discriminator value
				// e.g., "Circle" -> "circle"
				key := strings.ToLower(variant)
				mapping[key] = variant
			}
		}

		discriminator = &DiscriminatorData{
			PropertyName: m.Discriminator.PropertyName,
			Mapping:      mapping,
		}
	}

	props := g.convertProperties(m.Properties, m.Name)
	hasValidation := false
	for _, p := range props {
		if p.HasValidation {
			hasValidation = true
			break
		}
	}

	return ModelData{
		Name:          m.Name,
		Description:   m.Description,
		Properties:    props,
		EmbeddedTypes: embeddedTypes,
		IsArray:       m.IsArray,
		ArrayType:     g.getArrayItemsType(m.Items),
		IsMap:         m.IsMap,
		MapValueType:  g.getArrayItemsType(m.AdditionalProps),
		HasValidation: hasValidation,

		IsOneOf:       m.IsOneOf,
		IsAnyOf:       m.IsAnyOf,
		UnionVariants: unionVariants,
		Discriminator: discriminator,
	}
}

func (g *Generator) getArrayItemsType(items *core.Property) string {
	if items == nil {
		return "interface{}"
	}

	if items.RefType != "" {
		return items.RefType
	}

	return g.mapPrimitiveType(items.Type, items.Format)
}

func (g *Generator) convertProperties(props []core.Property, modelName string) []PropertyData {
	var result []PropertyData
	for _, p := range props {
		// Build enum value strings
		enumVals := make([]string, 0, len(p.Enum))
		for _, e := range p.Enum {
			enumVals = append(enumVals, fmt.Sprintf("%v", e))
		}
		isFormatValidated := codegen.IsFormatValidated(p.Format)
		goType := g.mapType(p, modelName)
		zeroExpr := zeroCheckExpr(p, goType)
		hasValidation := codegen.HasConstraints(
			p.MinLength != nil, p.MaxLength != nil,
			p.Minimum != nil, p.Maximum != nil,
			p.Pattern != "", len(p.Enum) > 0, p.Required && zeroExpr != "",
		) || isFormatValidated

		result = append(result, PropertyData{
			Name:        p.Name,
			Type:        goType,
			Description: p.Description,
			JsonTag:     p.Name,

			Default:    g.formatDefault(p.Default),
			ReadOnly:   p.ReadOnly,
			WriteOnly:  p.WriteOnly,
			Deprecated: p.Deprecated,

			Required:          p.Required,
			ZeroCheckExpr:     zeroExpr,
			Format:            p.Format,
			MinLength:         p.MinLength,
			MaxLength:         p.MaxLength,
			Minimum:           p.Minimum,
			Maximum:           p.Maximum,
			Pattern:           p.Pattern,
			EnumValues:        enumVals,
			HasEnum:           len(p.Enum) > 0,
			HasValidation:     hasValidation,
			IsFormatValidated: isFormatValidated,
		})
	}
	return result
}

// zeroCheckExpr returns the Go comparison expression used to detect a zero/missing
// required field. It uses the underlying OpenAPI scalar type (p.Type) rather than the
// final Go type name, so that named enum types are handled correctly.
// Returns "" for types where a required check cannot be generated (e.g. required struct refs).
func zeroCheckExpr(p core.Property, goType string) string {
	// Pointer types, slices, maps, and interface{} are always nil-checked.
	if strings.HasPrefix(goType, "*") ||
		strings.HasPrefix(goType, "[]") ||
		strings.HasPrefix(goType, "map[") ||
		goType == "interface{}" {
		return "== nil"
	}

	// Use the underlying OpenAPI type to determine the zero value.
	// This handles named types (enums) that map to scalar Go types.
	switch p.Type {
	case "string":
		return `== ""`
	case "integer", "number":
		return "== 0"
	case "boolean":
		return "== false"
	default:
		// Object refs (structs): no meaningful scalar zero check — skip.
		return ""
	}
}

func (g *Generator) formatDefault(val interface{}) string {
	if val == nil {
		return ""
	}

	return fmt.Sprintf("%v", val)
}

// mapType converts OpenAPI types to Golang types.
func (g *Generator) mapType(prop core.Property, modelName string) string {

	if prop.RefType != "" {
		if !prop.Required {
			return "*" + prop.RefType
		}

		return prop.RefType
	}

	// Handle enum types
	if len(prop.Enum) > 0 {
		enumTypeName := modelName + template.Capitalize(prop.Name)
		if !prop.Required {
			return "*" + enumTypeName
		}
		return enumTypeName
	}

	// Handle array types
	var baseType string
	switch prop.Type {
	case "array":
		baseType = "[]" + g.getArrayItemsType(prop.Items)
	case "object":
		{
			// For now, use map[string]interface{} for nested objects
			// TODO: Support inline struct or type references
			baseType = "map[string]interface{}"
		}
	default:
		{
			baseType = g.mapPrimitiveType(prop.Type, prop.Format)
		}
	}

	if !prop.Required && prop.Type != "array" && prop.Type != "object" {
		return "*" + baseType
	}

	return baseType
}

// mapType converts OpenAPI primitive types to Golang types.
func (g *Generator) mapPrimitiveType(specType string, format string) string {
	switch specType {
	case "string":
		if format == "date-time" {
			return "time.Time"
		}
		return "string"
	case "integer":
		if format == "int64" {
			return "int64"
		}
		if format == "int32" {
			return "int32"
		}
		return "int"
	case "number":
		if format == "float" {
			return "float32"
		}
		return "float64"
	case "boolean":
		return "bool"
	default:
		return "interface{}" // fallback
	}
}

func (g *Generator) collectImports(models []ModelData) []string {
	importsNeeded := make(map[string]bool)

	for _, model := range models {
		// Check for oneOf types (need fmt for error messages)
		if model.IsOneOf {
			importsNeeded["fmt"] = true
			importsNeeded["encoding/json"] = true
		}

		if model.HasValidation {
			importsNeeded["fmt"] = true
		}

		for _, prop := range model.Properties {
			// Check if the type uses time.Time
			if strings.Contains(prop.Type, "time.Time") {
				importsNeeded["time"] = true
			}
			if prop.IsFormatValidated && (prop.Format == "date-time" || prop.Format == "date") {
				importsNeeded["time"] = true
			}
			if prop.Pattern != "" {
				importsNeeded["regexp"] = true
			}
		}
	}

	// Convert map to sorted slice for consistent output
	imports := make([]string, 0, len(importsNeeded))
	for imp := range importsNeeded {
		imports = append(imports, imp)
	}
	sort.Strings(imports)

	return imports
}

func (g *Generator) collectEnums(models []core.Model) []EnumData {
	var enums []EnumData
	enumMap := make(map[string]bool) // Track unique enum type names

	// First pass: top-level scalar enum schemas (e.g. CustomerOrderStatus, PickValidationType).
	// These must be registered before scanning properties so that property-level
	// inline enums with the same derived name are correctly deduplicated.
	for _, model := range models {
		if !isTopLevelEnum(model) {
			continue
		}
		if enumMap[model.Name] {
			continue
		}
		enumMap[model.Name] = true

		baseType := g.mapPrimitiveType(model.Type, "")
		enumValues := make([]EnumValue, 0, len(model.Enum))
		for _, val := range model.Enum {
			constName := model.Name + template.Capitalize(fmt.Sprintf("%v", val))
			enumValues = append(enumValues, EnumValue{Name: constName, Value: val})
		}
		enums = append(enums, EnumData{
			TypeName:    model.Name,
			BaseType:    baseType,
			Description: model.Description,
			Values:      enumValues,
		})
	}

	// Second pass: inline property-level enum values (e.g. CustomField.type → CustomFieldType).
	for _, model := range models {
		for _, prop := range model.Properties {
			if len(prop.Enum) == 0 {
				continue
			}

			// Generate enum type name: ModelName + PropertyName
			enumTypeName := model.Name + template.Capitalize(prop.Name)

			if enumMap[enumTypeName] {
				continue
			}
			enumMap[enumTypeName] = true

			baseType := g.mapPrimitiveType(prop.Type, prop.Format)
			enumValues := make([]EnumValue, 0, len(prop.Enum))
			for _, val := range prop.Enum {
				constName := enumTypeName + template.Capitalize(fmt.Sprintf("%v", val))
				enumValues = append(enumValues, EnumValue{Name: constName, Value: val})
			}

			enums = append(enums, EnumData{
				TypeName:    enumTypeName,
				BaseType:    baseType,
				Description: prop.Description,
				Values:      enumValues,
			})
		}
	}

	return enums
}
