package parser

import (
	"fmt"
	"log/slog"
	"sort"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/imerljak/beaverspec/pkg/core"
)

type Normalizer struct {
	// TODO
}

func NewNormalizer() *Normalizer {
	return &Normalizer{}
}

func (n *Normalizer) Normalize(doc *openapi3.T) (*core.Spec, error) {

	var models []core.Model
	if doc.Components != nil && doc.Components.Schemas != nil {
		models = n.extractModels(doc.Components.Schemas)
	} else {
		models = []core.Model{} // Empty slice if no components
	}

	spec := &core.Spec{
		Raw:       doc,
		Version:   doc.OpenAPI,
		Info:      n.extractInfo(doc.Info),
		Models:    models,
		Endpoints: n.extractEndpoints(doc.Paths),
		Tags:      n.extractTags(doc.Tags),
	}
	return spec, nil
}

// extractEndpoints converts openapi3.Paths into []core.Endpoint
func (n *Normalizer) extractEndpoints(paths *openapi3.Paths) []core.Endpoint {
	var endpoints []core.Endpoint

	for path, pathItem := range paths.Map() {
		if pathItem == nil {
			continue
		}

		// Extract operations for each HTTP method
		operations := map[string]*openapi3.Operation{
			"GET":     pathItem.Get,
			"POST":    pathItem.Post,
			"PUT":     pathItem.Put,
			"DELETE":  pathItem.Delete,
			"PATCH":   pathItem.Patch,
			"HEAD":    pathItem.Head,
			"OPTIONS": pathItem.Options,
		}

		for method, operation := range operations {
			if operation == nil {
				continue
			}

			endpoint := core.Endpoint{
				Path:         path,
				Method:       method,
				OperationID:  operation.OperationID,
				Summary:      operation.Summary,
				Description:  operation.Description,
				Tags:         operation.Tags,
				Parameters:   n.extractParameters(operation.Parameters, pathItem.Parameters),
				RequestBody:  n.extractRequestBody(operation.RequestBody),
				Responses:    n.extractResponses(operation.Responses),
				Security:     n.extractSecurityRequirements(operation.Security),
				IsDeprecated: operation.Deprecated,
			}

			endpoints = append(endpoints, endpoint)
		}
	}

	return endpoints
}

func (n *Normalizer) extractSecurityRequirements(security *openapi3.SecurityRequirements) []core.SecurityRequirement {
	if security == nil {
		return nil
	}

	var requirements []core.SecurityRequirement

	for _, secReq := range *security {
		for name, scopes := range secReq {
			requirements = append(requirements, core.SecurityRequirement{
				Name:   name,
				Scopes: scopes,
			})
		}
	}

	return requirements
}

func (n *Normalizer) extractResponses(responses *openapi3.Responses) []core.Response {
	if responses == nil {
		return nil
	}

	var result []core.Response

	for statusCode, responseRef := range responses.Map() {
		if responseRef == nil || responseRef.Value == nil {
			continue
		}

		response := responseRef.Value
		content := make(map[string]core.MediaType)

		for contentType, mediaType := range response.Content {
			if mediaType == nil {
				continue
			}

			content[contentType] = core.MediaType{
				Schema:  n.extractParameterSchema(mediaType.Schema),
				Example: mediaType.Example,
			}
		}

		headers := make(map[string]core.Header)
		for headerName, headerRef := range response.Headers {
			if headerRef == nil || headerRef.Value == nil {
				continue
			}

			header := headerRef.Value
			headers[headerName] = core.Header{
				Description: header.Description,
				Required:    header.Required,
				Schema:      n.extractParameterSchema(header.Schema),
			}
		}

		result = append(result, core.Response{
			StatusCode:  statusCode,
			Description: *response.Description,
			Content:     content,
			Headers:     headers,
		})
	}

	return result
}

func (n *Normalizer) extractRequestBody(bodyRef *openapi3.RequestBodyRef) *core.RequestBody {
	if bodyRef == nil || bodyRef.Value == nil {
		return nil
	}

	body := bodyRef.Value
	content := make(map[string]core.MediaType)

	for contentType, mediaType := range body.Content {
		if mediaType == nil {
			continue
		}

		content[contentType] = core.MediaType{
			Schema:  n.extractParameterSchema(mediaType.Schema),
			Example: mediaType.Example,
		}
	}

	return &core.RequestBody{
		Description: body.Description,
		Required:    body.Required,
		Content:     content,
	}
}

func (n *Normalizer) extractParameters(opParams, pathParams openapi3.Parameters) []core.Parameter {
	var params []core.Parameter

	// Combine path-level and operation-level parameters
	allParams := append(pathParams, opParams...)

	for _, paramRef := range allParams {
		if paramRef == nil || paramRef.Value == nil {
			continue
		}

		param := paramRef.Value
		schema := n.extractParameterSchema(param.Schema)

		params = append(params, core.Parameter{
			Name:         param.Name,
			In:           param.In,
			Description:  param.Description,
			Required:     param.Required,
			Schema:       schema,
			Example:      param.Example,
			IsDeprecated: param.Deprecated,
		})
	}
	return params
}

func (n *Normalizer) extractParameterSchema(schemaRef *openapi3.SchemaRef) *core.Property {
	if schemaRef == nil || schemaRef.Value == nil {
		return nil
	}

	schema := schemaRef.Value
	prop := &core.Property{}

	// Handle $ref
	if schemaRef.Ref != "" {
		prop.RefType = extractRefTypeName(schemaRef.Ref)
		return prop
	}

	// Extract type
	if schema.Type != nil && len(*schema.Type) > 0 {
		prop.Type = (*schema.Type)[0]
		prop.Format = schema.Format
	}

	// Handle arrays
	if prop.Type == "array" && schema.Items != nil {
		prop.Items = n.extractParameterSchema(schema.Items)
	}

	return prop
}

// extractInfo converts openapi3.Info to core.SpecInfo
func (n *Normalizer) extractInfo(info *openapi3.Info) core.SpecInfo {
	specInfo := core.SpecInfo{
		Title:       info.Title,
		Description: info.Description,
		Version:     info.Version,
	}

	if info.Contact != nil {
		specInfo.Contact = &core.Contact{
			Name:  info.Contact.Name,
			Email: info.Contact.Email,
			URL:   info.Contact.URL,
		}
	}

	if info.License != nil {
		specInfo.License = &core.License{
			Name: info.License.Name,
			URL:  info.License.URL,
		}
	}

	return specInfo
}

// extractTags converts openapi3.Tags to []core.Tag
func (n *Normalizer) extractTags(tags openapi3.Tags) []core.Tag {
	tagList := make([]core.Tag, 0, len(tags))
	for _, t := range tags {
		tagList = append(tagList, core.Tag{
			Name:        t.Name,
			Description: t.Description,
		})
	}
	return tagList
}

func (n *Normalizer) extractModels(schemas openapi3.Schemas) []core.Model {
	models := make([]core.Model, 0, len(schemas))
	nestedModels := make([]core.Model, 0) // Track nested models to add

	for name, schemaRef := range schemas {
		if schemaRef == nil || schemaRef.Value == nil {
			continue
		}

		schema := schemaRef.Value

		//Check if this is an array type schema
		var schemaType string
		if schema.Type != nil && len(*schema.Type) > 0 {
			schemaType = (*schema.Type)[0]
		}

		// Check if this is a map type (object with additionalProperties but no regular properties)
		if schemaType == "object" && schema.AdditionalProperties.Schema != nil && len(schema.Properties) == 0 {
			// This is a map type
			model := core.Model{
				Name:            name,
				Description:     schema.Description,
				IsMap:           true,
				AdditionalProps: n.extractAdditionalProperties(schema.AdditionalProperties.Schema),
			}
			models = append(models, model)
		} else if len(schema.OneOf) > 0 {
			// handle oneOf - require discriminator
			if schema.Discriminator == nil {
				// TODO: Add to validation errors instead of skipping
				slog.Warn("oneOf schema is missing required discriminator, skipping", "schema", name)
				continue
			}

			oneOfVariants := n.extractOneOf(schema.OneOf, name)
			discriminator := n.extractDiscriminator(schema.Discriminator)

			model := core.Model{
				Name:          name,
				Description:   schema.Description,
				IsOneOf:       true,
				OneOf:         oneOfVariants,
				Discriminator: discriminator,
			}

			models = append(models, model)
		} else if len(schema.AnyOf) > 0 {
			// handle anyOf
			anyOfOptions := n.extractAnyOf(schema.AnyOf, name)

			model := core.Model{
				Name:        name,
				Description: schema.Description,
				IsAnyOf:     true,
				AnyOf:       anyOfOptions,
			}

			models = append(models, model)
		} else if schemaType == "array" {
			// handle array type schema
			model := core.Model{
				Name:        name,
				Description: schema.Description,
				IsArray:     true,
				Items:       n.extractArrayItems(schema.Items),
			}
			models = append(models, model)
		} else {
			// Handle object type schema
			props, nested := n.extractPropertiesWithNested(schema.Properties, schema.Required, name)

			// Extract allOf if present
			var allOfModels []core.Model
			if len(schema.AllOf) > 0 {
				allOfModels = n.extractAllOf(schema.AllOf, name)
			}

			model := core.Model{
				Name:        name,
				Description: schema.Description,
				Properties:  props,
				AllOf:       allOfModels,
				// TODO: extract more ?
			}
			models = append(models, model)
			nestedModels = append(nestedModels, nested...)
		}
	}

	models = append(models, nestedModels...)
	return models
}

func (n *Normalizer) extractOneOf(oneOfSchemas openapi3.SchemaRefs, parentName string) []core.Model {
	var models []core.Model

	for _, schemaRef := range oneOfSchemas {
		if schemaRef == nil || schemaRef.Value == nil {
			continue
		}

		var model core.Model

		// oneOf should typically be $refs to other schemas
		if schemaRef.Ref != "" {
			model = core.Model{
				Name: extractRefTypeName(schemaRef.Ref),
			}
		} else {
			// Inline schema - less common for oneOf
			schema := schemaRef.Value
			props, _ := n.extractPropertiesWithNested(schema.Properties, schema.Required, parentName+"Variant")

			model = core.Model{
				Name:       parentName + "Variant",
				Properties: props,
			}
		}
		models = append(models, model)
	}
	return models
}

func (n *Normalizer) extractAnyOf(anyOfSchemas openapi3.SchemaRefs, parentName string) []core.Model {
	// Similar to extractOneOf
	var models []core.Model

	for _, schemaRef := range anyOfSchemas {
		if schemaRef == nil || schemaRef.Value == nil {
			continue
		}

		var model core.Model

		if schemaRef.Ref != "" {
			model = core.Model{
				Name: extractRefTypeName(schemaRef.Ref),
			}
		} else {
			schema := schemaRef.Value
			props, _ := n.extractPropertiesWithNested(schema.Properties, schema.Required, parentName+"Option")

			model = core.Model{
				Name:       parentName + "Option",
				Properties: props,
			}
		}

		models = append(models, model)
	}

	return models
}

func (n *Normalizer) extractDiscriminator(disc *openapi3.Discriminator) *core.Discriminator {
	if disc == nil {
		return nil
	}

	return &core.Discriminator{
		PropertyName: disc.PropertyName,
		Mapping:      disc.Mapping,
	}
}

func (n *Normalizer) extractAllOf(allOfSchemas openapi3.SchemaRefs, parentName string) []core.Model {
	var models []core.Model

	for i, schemaRef := range allOfSchemas {
		if schemaRef == nil || schemaRef.Value == nil {
			continue
		}

		var model core.Model

		// Check if it's a $ref
		if schemaRef.Ref != "" {
			model = core.Model{
				Name: extractRefTypeName(schemaRef.Ref),
			}
		} else {
			// Inline schema - extract properties
			schema := schemaRef.Value
			props, _ := n.extractPropertiesWithNested(schema.Properties, schema.Required, fmt.Sprintf("%s_AllOf%d", parentName, i))

			model = core.Model{
				Name:       fmt.Sprintf("%sAllOf%d", parentName, i), // Generate inline name
				Properties: props,
			}
		}

		models = append(models, model)
	}

	return models
}

func (n *Normalizer) extractAdditionalProperties(propsRef *openapi3.SchemaRef) *core.Property {
	if propsRef == nil || propsRef.Value == nil {
		return &core.Property{Type: "interface{}"}
	}

	prop := &core.Property{}

	// Check if it's a $ref
	if propsRef.Ref != "" {
		prop.RefType = extractRefTypeName(propsRef.Ref)
	} else {
		// Extract value type
		if propsRef.Value.Type != nil && len(*propsRef.Value.Type) > 0 {
			prop.Type = (*propsRef.Value.Type)[0]
			prop.Format = propsRef.Value.Format
		}
	}

	return prop
}

func (n *Normalizer) extractArrayItems(itemsRef *openapi3.SchemaRef) *core.Property {
	if itemsRef == nil || itemsRef.Value == nil {
		return &core.Property{Type: "interface{}"}
	}

	items := &core.Property{}

	// Check if items is a $ref
	if itemsRef.Ref != "" {
		items.RefType = extractRefTypeName(itemsRef.Ref)
	} else {
		// Extract item type
		if itemsRef.Value.Type != nil && len(*itemsRef.Value.Type) > 0 {
			items.Type = (*itemsRef.Value.Type)[0]
			items.Format = itemsRef.Value.Format
		}
	}

	return items
}

func (n *Normalizer) extractPropertiesWithNested(props openapi3.Schemas, required []string, parentModelName string) ([]core.Property, []core.Model) {
	properties := make([]core.Property, 0, len(props))
	nestedModels := make([]core.Model, 0)

	names := make([]string, 0, len(props))
	for name := range props {
		names = append(names, name)
	}
	sort.Strings(names)

	requiredMap := make(map[string]bool)
	for _, req := range required {
		requiredMap[req] = true
	}

	for _, name := range names {
		propRef := props[name]

		if propRef == nil || propRef.Value == nil {
			continue
		}

		prop := propRef.Value
		propType := ""
		if prop.Type != nil && len(*prop.Type) > 0 {
			propType = (*prop.Type)[0]
		}

		var itemsType string
		if propType == "array" && prop.Items != nil && prop.Items.Value != nil {
			if prop.Items.Value.Type != nil && len(*prop.Items.Value.Type) > 0 {
				itemsType = (*prop.Items.Value.Type)[0]
			}
		}

		var refType string
		if propRef.Ref != "" {
			refType = extractRefTypeName(propRef.Ref)
		} else if propType == "object" && len(prop.Properties) > 0 {
			// Handle inline nested object - extract as separate model
			nestedTypeName := generateNestedTypeName(parentModelName, name)
			refType = nestedTypeName

			// Recursively extract nested properties
			nestedProps, deeplyNested := n.extractPropertiesWithNested(prop.Properties, prop.Required, nestedTypeName)

			nestedModel := core.Model{
				Name:        nestedTypeName,
				Description: prop.Description,
				Properties:  nestedProps,
			}

			nestedModels = append(nestedModels, nestedModel)
			nestedModels = append(nestedModels, deeplyNested...)
		}

		property := core.Property{
			Name:        name,
			Type:        propType,
			Format:      prop.Format,
			Description: prop.Description,
			Required:    requiredMap[name],
			Nullable:    prop.Nullable,
			Enum:        prop.Enum,
			Items: &core.Property{
				Type: itemsType,
			},
			RefType: refType,

			Default:    prop.Default,
			ReadOnly:   prop.ReadOnly,
			WriteOnly:  prop.WriteOnly,
			Deprecated: prop.Deprecated,
		}

		// Still set Properties for backward compatibility, but RefType takes precedence
		if propType == "object" && refType == "" {
			nestedProps, _ := n.extractPropertiesWithNested(prop.Properties, prop.Required, "")
			property.Properties = nestedProps
		}

		properties = append(properties, property)
	}

	return properties, nestedModels
}

// capitalizeFirst capitalizes the first letter of a string
func capitalizeFirst(s string) string {
	if len(s) == 0 {
		return s
	}

	return strings.ToUpper(s[:1]) + s[1:]
}

// extractRefTypeName extracts the type name from a $ref path
// Example "#/components/schemas/User" -> "User"
func extractRefTypeName(ref string) string {
	if ref == "" {
		return ""
	}

	i := strings.LastIndex(ref, "/")
	if i == -1 {
		return ref
	}

	return ref[i+1:]
}

// generateNestedTypeName generates a type name for nested objects
func generateNestedTypeName(parentModel, propName string) string {
	return parentModel + capitalizeFirst(propName)
}
