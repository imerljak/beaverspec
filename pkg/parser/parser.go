package parser

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

type Parser struct {
	// Validation settings
	// Caching options
	// Custom loaders
	// etc...
}

func NewParser() *Parser {
	return &Parser{}
}

// LoadSpec loads an OpenAPI specification from a file path or HTTP/HTTPS URL.
// Returns the raw openapi3 document.
func (p *Parser) LoadSpec(ctx context.Context, specPath string) (*openapi3.T, error) {
	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true

	if strings.HasPrefix(specPath, "http://") || strings.HasPrefix(specPath, "https://") {
		u, err := url.ParseRequestURI(specPath)
		if err != nil {
			return nil, fmt.Errorf("invalid spec URL %s: %w", specPath, err)
		}
		t, err := loader.LoadFromURI(u)
		if err != nil {
			return nil, fmt.Errorf("failed to load spec from %s: %w", specPath, err)
		}
		return t, nil
	}

	t, err := loader.LoadFromFile(specPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load spec from %s: %w", specPath, err)
	}
	return t, nil
}

// ValidateSpec validates an OpenAPI specification
func (p *Parser) ValidateSpec(ctx context.Context, spec *openapi3.T) error {
	if err := spec.Validate(ctx); err != nil {
		return fmt.Errorf("spec validation failed: %w", err)
	}
	return nil
}
