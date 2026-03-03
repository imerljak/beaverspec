package template

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// Engine handles template loading and rendering
type Engine interface {
	// Load templates from a directory or embed
	Load(templateDir string) error

	// Render a template with data
	Render(templateName string, data interface{}) (string, error)

	// Register a custom helper function
	RegisterFunction(name string, fn interface{}) error

	// Check if a template exists
	HasTemplate(name string) bool
}

// GoTemplateEngine implements Engine using Go's text/template
type GoTemplateEngine struct {
	templates *template.Template
	funcMap   template.FuncMap
}

// NewGoTemplateEngine creates a new template engine
func NewGoTemplateEngine() *GoTemplateEngine {
	return &GoTemplateEngine{
		funcMap: DefaultHelpers(), // Load default helper functions
	}
}

// RegisterFunction registers a custom template function
func (e *GoTemplateEngine) RegisterFunction(name string, fn interface{}) error {
	if _, ok := e.funcMap[name]; !ok {
		e.funcMap[name] = fn
		return nil
	}
	return fmt.Errorf("function '%s' already exists", name)
}

// Load loads templates from a directory (recursively)
func (e *GoTemplateEngine) Load(templateDir string) error {
	templates := template.New("root").Funcs(e.funcMap)

	// Walk the directory tree to find all .tmpl files
	err := filepath.Walk(templateDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Only process .tmpl files
		if !strings.HasSuffix(path, ".tmpl") {
			return nil
		}

		// Get relative path from templateDir for the template name
		relPath, err := filepath.Rel(templateDir, path)
		if err != nil {
			return err
		}

		// Use forward slashes for template names (cross-platform)
		templateName := filepath.ToSlash(relPath)

		// Read template content
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		// Parse the template with the correct name
		_, err = templates.New(templateName).Parse(string(content))
		return err
	})

	if err != nil {
		return fmt.Errorf("failed to load templates from %s: %w", templateDir, err)
	}

	e.templates = templates
	return nil
}

// Render executes a template with the given data
func (e *GoTemplateEngine) Render(templateName string, data interface{}) (string, error) {
	t := e.templates.Lookup(templateName)
	if t == nil {
		return "", fmt.Errorf("no template found with name '%s'", templateName)
	}

	var buff bytes.Buffer
	if err := t.Execute(&buff, data); err != nil {
		return "", fmt.Errorf("failed to execute template '%s': %w", templateName, err)
	}

	return buff.String(), nil
}

// Override registers or replaces a single template by name with the given content.
// Call this after Load() to replace individual built-in templates with custom ones.
func (e *GoTemplateEngine) Override(name, content string) error {
	if e.templates == nil {
		return fmt.Errorf("engine not loaded: call Load() before Override()")
	}
	_, err := e.templates.New(name).Parse(content)
	return err
}

// HasTemplate checks if a template exists
func (e *GoTemplateEngine) HasTemplate(name string) bool {
	if e.templates == nil {
		return false
	}

	return e.templates.Lookup(name) != nil
}
