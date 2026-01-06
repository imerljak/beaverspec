package parser

import (
	"sort"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/imerljak/openapi-gen/pkg/core"
)

type Normalizer struct {
	// TODO
}

func NewNormalizer() *Normalizer {
	return &Normalizer{}
}

func (n *Normalizer) Normalize(doc *openapi3.T) (*core.Spec, error) {
	spec := &core.Spec{
		Raw:       doc,
		Version:   doc.OpenAPI,
		Info:      n.extractInfo(doc.Info),
		Models:    n.extractModels(doc.Components.Schemas),
		Endpoints: []core.Endpoint{}, //TODO
		Tags:      n.extractTags(doc.Tags),
	}
	return spec, nil
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

	for name, schemaRef := range schemas {
		if schemaRef == nil || schemaRef.Value == nil {
			continue
		}

		schema := schemaRef.Value
		model := core.Model{
			Name:        name,
			Description: schema.Description,
			Properties:  n.extractProperties(schema.Properties, schema.Required),
			// TODO: extract more ?
		}
		models = append(models, model)
	}
	return models
}

func (n *Normalizer) extractProperties(props openapi3.Schemas, required []string) []core.Property {
	names := make([]string, 0, len(props))
	for name := range props {
		names = append(names, name)
	}
	sort.Strings(names)

	properties := make([]core.Property, 0, len(props))
	for _, name := range names {
		propRef := props[name]

		if propRef == nil || propRef.Value == nil {
			continue
		}

		// Extract type - openapi3.Types is a slice, take the first type
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

		isRequired := false
		for _, reqField := range required {
			if reqField == name {
				isRequired = true
				break
			}
		}

		property := core.Property{
			Name:        name,
			Type:        propType,
			Format:      prop.Format,
			Description: prop.Description,
			Required:    isRequired,
			Nullable:    prop.Nullable,
			Items: &core.Property{
				Type: itemsType,
			},
		}

		// Handle nested objects
		if propType == "object" {
			property.Properties = n.extractProperties(prop.Properties, prop.Required)
		}

		properties = append(properties, property)
	}
	return properties
}
