package golang

import (
	"fmt"
	"go/format"
	"sort"
	"strings"

	"github.com/imerljak/openapi-gen/pkg/core"
	"github.com/imerljak/openapi-gen/pkg/template"
)

type Generator struct {
	//
}

// ModelData is the data passed to model templates
type ModelData struct {
	Name        string
	Description string
	Properties  []PropertyData
}

type PropertyData struct {
	Name, Type, Description, JsonTag string
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
			SupportsClient:     false, // Not yet
			SupportsServer:     false, // Not yet
			SupportsValidation: false, // Not yet
			SupportsAsync:      false,
			SupportsStreaming:  false,
			ClientFrameworks:   []string{},
			ServerFrameworks:   []string{},
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
	if config.Options != nil {
		if pkgName, ok := config.Options["package"].(string); ok && pkgName != "" {
			packageName = pkgName
		}
	}

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

	// Create the generated file
	result := &core.GenerationResult{
		Files: []core.GeneratedFile{
			{
				Path:    "models.go",
				Content: []byte(formattedContent),
			},
			{
				Path:    "models_test.go",
				Content: []byte(formattedTestContent),
			},
		},
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

func (g *Generator) convertModels(models []core.Model) []ModelData {
	var result []ModelData
	for _, m := range models {
		result = append(result, g.convertModel(m))
	}
	return result
}

func (g *Generator) convertModel(m core.Model) ModelData {
	return ModelData{
		Name:        m.Name,
		Description: m.Description,
		Properties:  g.convertProperties(m.Properties, m.Name),
	}
}

func (g *Generator) convertProperties(props []core.Property, modelName string) []PropertyData {
	var result []PropertyData
	for _, p := range props {
		result = append(result, PropertyData{
			Name:        p.Name,
			Type:        g.mapType(p, modelName), // Map OpenAPI to Go types
			Description: p.Description,
			JsonTag:     p.Name, // TODO: apply naming convention
		})
	}
	return result
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
