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
	PackageName   string
	InterfaceName string
	Operations    []OperationData
	Imports       []codegen.Import
}

// ServerData represents data for server generation
type ServerData struct {
	PackageName string
	Imports     []codegen.Import
	Handlers    []HandlerGroup
}

// HandlerGroup groups operations by tag for routing
type HandlerGroup struct {
	Name       string
	Operations []OperationData
}

// OperationData represents a single client/server operation
type OperationData struct {
	Name         string // Method name (e.g., "ListPets")
	OperationID  string
	Description  string
	Method       string // HTTP method
	Path         string // URL path
	PathParams   []ParamData
	QueryParams  []ParamData
	HeaderParams []ParamData
	HasBody      bool
	BodyType     string // Type of request body
	ReturnType   string // Return type
	ErrorReturn  bool   // Whether to return error
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
	if config.Options != nil {
		if pkgName, ok := config.Options["package"].(string); ok && pkgName != "" {
			packageName = pkgName
		}
		if mp, ok := config.Options["modulePath"].(string); ok && mp != "" {
			modulePath = mp
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
		clientData := g.convertEndpointsToClient(spec.Endpoints, modulePath, baseDir)

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
		serverData := g.convertEndpointsToServer(spec.Endpoints, modulePath, baseDir)

		// interface
		serverInterfaceContent, err := engine.Render("server/interface.go.tmpl", serverData)
		if err != nil {
			return nil, fmt.Errorf("failed to render server interface: %w", err)
		}
		formattedInterface, err := format.Source([]byte(serverInterfaceContent))
		if err == nil {
			files = append(files, core.GeneratedFile{Path: "server/interface.go", Content: formattedInterface})
		}

		// handlers
		serverHandlersContent, err := engine.Render("server/handlers.go.tmpl", serverData)
		if err != nil {
			return nil, fmt.Errorf("failed to render server handlers: %w", err)
		}
		formattedHandlers, err := format.Source([]byte(serverHandlersContent))
		if err == nil {
			files = append(files, core.GeneratedFile{Path: "server/handlers.go", Content: formattedHandlers})
		}

		// routes
		serverRoutesContent, err := engine.Render("server/routes.go.tmpl", serverData)
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
func (g *Generator) convertEndpointsToClient(endpoints []core.Endpoint, modulePath, baseDir string) ClientData {

	imports := codegen.NewImportManager(modulePath, baseDir, "client")
	imports.Add("context")
	imports.Add("net/http")
	imports.Add("fmt")
	imports.Add("io")
	imports.Add("encoding/json")

	operations := make([]OperationData, 0, len(endpoints))

	for _, ep := range endpoints {
		op := g.convertEndpoint(ep)
		operations = append(operations, op)
	}

	for _, op := range operations {
		if len(op.QueryParams) > 0 {
			imports.Add("net/url")
			break
		}
	}

	imports.AddSibling("models")

	return ClientData{
		PackageName:   "client",
		InterfaceName: "Client",
		Operations:    operations,
		Imports:       imports.GetImports(),
	}
}

// convertEndpointsToServer converts core endpoints into server structures
func (g *Generator) convertEndpointsToServer(endpoints []core.Endpoint, modulePath, baseDir string) ServerData {
	imports := codegen.NewImportManager(modulePath, baseDir, "server")
	imports.Add("context")
	imports.Add("net/http")
	imports.Add("encoding/json")

	groupsMap := make(map[string][]OperationData)

	for _, ep := range endpoints {
		op := g.convertEndpoint(ep)

		// Group by first tag if present, else "Default"
		groupName := "Default"
		if len(ep.Tags) > 0 {
			groupName = template.ToPascalCase(ep.Tags[0])
		}

		groupsMap[groupName] = append(groupsMap[groupName], op)
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

	// Add fmt/io imports if operations need them (future-proofing)
	// We add them statically for now as they are very common
	imports.Add("fmt")
	imports.Add("io")

	return ServerData{
		PackageName: "server",
		Imports:     imports.GetImports(),
		Handlers:    handlers,
	}
}

// convertEndpoint converts a single endpoint to operation data
func (g *Generator) convertEndpoint(ep core.Endpoint) OperationData {
	// Generate method name from operation ID
	methodName := g.toMethodName(ep.OperationID)

	// Extract parameters by location
	paramMap := g.extractParametersByLocation(ep.Parameters)

	// Determine body type and return type
	bodyType := g.getRequestBodyType(ep.RequestBody)
	returnType := g.getResponseType(ep.Responses)

	return OperationData{
		Name:         methodName,
		OperationID:  ep.OperationID,
		Description:  ep.Description,
		Method:       ep.Method,
		Path:         ep.Path,
		PathParams:   paramMap["path"],
		QueryParams:  paramMap["query"],
		HeaderParams: paramMap["header"],
		HasBody:      bodyType != "",
		BodyType:     bodyType,
		ReturnType:   returnType,
		ErrorReturn:  true, // Always return error
	}
}

// getResponseType determines the return type from responses
func (g *Generator) getResponseType(responses []core.Response) string {
	// Look for 200, 201, or default success response
	for _, resp := range responses {
		if resp.StatusCode == "200" || resp.StatusCode == "201" {
			// Look for JSON content
			if mediaType, ok := resp.Content["application/json"]; ok {
				if mediaType.Schema != nil {
					if mediaType.Schema.RefType != "" {
						return "*models." + mediaType.Schema.RefType // TODO: parameterize the package name
					}
					return g.mapParameterType(mediaType.Schema)
				}
			}
		}
	}

	// No body response
	return ""
}

// getRequestBodyType determines the request body type
func (g *Generator) getRequestBodyType(body *core.RequestBody) string {
	if body == nil {
		return ""
	}

	// Look for JSON content
	if mediaType, ok := body.Content["application/json"]; ok {
		if mediaType.Schema != nil {
			if mediaType.Schema.RefType != "" {
				return "*models." + mediaType.Schema.RefType // TODO: parameterize the package name
			}
			return g.mapParameterType(mediaType.Schema)
		}
	}
	return ""
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

	// Handle references
	if schema.RefType != "" {
		return schema.RefType
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
		result = append(result, g.convertModel(m))
	}
	return result
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

	return ModelData{
		Name:          m.Name,
		Description:   m.Description,
		Properties:    g.convertProperties(m.Properties, m.Name),
		EmbeddedTypes: embeddedTypes,
		IsArray:       m.IsArray,
		ArrayType:     g.getArrayItemsType(m.Items),
		IsMap:         m.IsMap,
		MapValueType:  g.getArrayItemsType(m.AdditionalProps),

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
		result = append(result, PropertyData{
			Name:        p.Name,
			Type:        g.mapType(p, modelName), // Map OpenAPI to Go types
			Description: p.Description,
			JsonTag:     p.Name, // TODO: apply naming convention

			Default:    g.formatDefault(p.Default),
			ReadOnly:   p.ReadOnly,
			WriteOnly:  p.WriteOnly,
			Deprecated: p.Deprecated,
		})
	}
	return result
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
		{
			if prop.Items != nil && prop.Items.Type != "" {
				itemType := g.mapPrimitiveType(prop.Items.Type, prop.Items.Format)
				baseType = "[]" + itemType
			} else {
				baseType = "[]interface{}" // Fallback for untyped arrays
			}
		}
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

		for _, prop := range model.Properties {
			// Check if the type uses time.Time
			if strings.Contains(prop.Type, "time.Time") {
				importsNeeded["time"] = true
			}
			// Future: add other import detection here
			// e.g., "encoding/json" for custom marshalers
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

	for _, model := range models {
		for _, prop := range model.Properties {
			// Skip if not an enum
			if len(prop.Enum) == 0 {
				continue
			}

			// Generate eunm type name: ModelName + PropertyName
			enumTypeName := model.Name + template.Capitalize(prop.Name)

			// Skip ifwe've already processed this enum type
			if enumMap[enumTypeName] {
				continue
			}
			enumMap[enumTypeName] = true

			// Determine base type from property type
			baseType := g.mapPrimitiveType(prop.Type, prop.Format)

			// Convert enum values to EnumValue structs
			enumValues := make([]EnumValue, 0, len(prop.Enum))
			for _, val := range prop.Enum {
				// Generate consntatn name: EnuMType + Value
				// e.g., "OrderStatus" + "pending" = "OrderStatusPending"
				constName := enumTypeName + template.Capitalize(fmt.Sprintf("%v", val))

				enumValues = append(enumValues, EnumValue{
					Name:  constName,
					Value: val,
				})
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
