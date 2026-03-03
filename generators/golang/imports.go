package golang

import (
	"path"
	"sort"
)

// Import represents a single Go import statement.
type Import struct {
	Path  string // e.g. "github.com/user/project/models"
	Alias string // optional alias (e.g. "echomiddleware")
}

// ImportManager manages imports for a generated Go file.
type ImportManager struct {
	modulePath string
	baseDir    string
	currentPkg string
	imports    map[string]string // path -> alias
}

// NewImportManager creates a new ImportManager.
func NewImportManager(modulePath, baseDir, currentPkg string) *ImportManager {
	return &ImportManager{
		modulePath: modulePath,
		baseDir:    baseDir,
		currentPkg: currentPkg,
		imports:    make(map[string]string),
	}
}

// Add registers a standard import and returns the local package name.
func (m *ImportManager) Add(pkgPath string) string {
	if _, exists := m.imports[pkgPath]; !exists {
		m.imports[pkgPath] = ""
	}
	return path.Base(pkgPath)
}

// AddWithAlias registers an import with an explicit alias.
func (m *ImportManager) AddWithAlias(pkgPath, alias string) string {
	m.imports[pkgPath] = alias
	return alias
}

// AddSibling registers a sibling package in the same output directory.
// e.g. from the "client" package, AddSibling("models") adds "modulePath/models".
func (m *ImportManager) AddSibling(pkgName string) string {
	if pkgName == m.currentPkg {
		return ""
	}
	fullPath := path.Join(m.modulePath, pkgName)
	return m.Add(fullPath)
}

// GetImports returns all registered imports sorted alphabetically.
func (m *ImportManager) GetImports() []Import {
	imports := make([]Import, 0, len(m.imports))
	for pkgPath, alias := range m.imports {
		imports = append(imports, Import{Path: pkgPath, Alias: alias})
	}
	sort.Slice(imports, func(i, j int) bool {
		return imports[i].Path < imports[j].Path
	})
	return imports
}

// HasImports returns true if any imports are registered.
func (m *ImportManager) HasImports() bool {
	return len(m.imports) > 0
}
