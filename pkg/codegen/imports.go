package codegen

import (
	"path"
	"sort"
)

// Import represents a single import statement
type Import struct {
	Path  string // "github.com/user/project/models"
	Alias string
}

// ImportManager manages imports for a generated file
type ImportManager struct {
	modulePath string            // "github.com/user/proj"
	baseDir    string            // "generated" (output dir relative to module)
	currentPkg string            // "client" (package we're generating)
	imports    map[string]string // path -> alias
}

// NewImportManager creates a new import manager
func NewImportManager(modulePath, baseDir, currentPkg string) *ImportManager {
	return &ImportManager{
		modulePath: modulePath,
		baseDir:    baseDir,
		currentPkg: currentPkg,
		imports:    make(map[string]string),
	}
}

// Add adds a standad import and returns the local name to use
func (m *ImportManager) Add(pkgPath string) string {
	if _, exists := m.imports[pkgPath]; !exists {
		m.imports[pkgPath] = ""
	}
	return path.Base(pkgPath)
}

// AddWithAlias adds an import with an explicit alias
func (m *ImportManager) AddWithAlias(pkgPath, alias string) string {
	m.imports[pkgPath] = alias
	return alias
}

// AddSibling adds an import for a sibling package in the same output directory
// e.g., from "client" package, AddSibling("models") adds "github.com/.../generated/models"
func (m *ImportManager) AddSibling(pkgName string) string {
	if pkgName == m.currentPkg {
		return "" // don't import self
	}

	fullPath := path.Join(m.modulePath, pkgName)
	return m.Add(fullPath)
}

// GetImports returns all imports sorted alphabetically
func (m *ImportManager) GetImports() []Import {
	var imports []Import

	for pkgPath, alias := range m.imports {
		imports = append(imports, Import{Path: pkgPath, Alias: alias})
	}

	sortImports(imports)
	return imports
}

// HasImports return true if there are any imports
func (m *ImportManager) HasImports() bool {
	return len(m.imports) > 0
}

func sortImports(imports []Import) {
	sort.Slice(imports, func(i, j int) bool {
		return imports[i].Path < imports[j].Path
	})
}
