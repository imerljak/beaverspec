package plugin

import (
	"fmt"
	"sync"

	"github.com/imerljak/openapi-gen/pkg/core"
)

// Registry manages available generators
type Registry struct {
	generators map[string]core.Generator
	mu         sync.RWMutex // Protects concurrent access
}

// NewRegistry creates a new generator registry
func NewRegistry() *Registry {
	return &Registry{
		generators: make(map[string]core.Generator),
	}
}

// Register adds a generator to the registry
func (r *Registry) Register(name string, gen core.Generator) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.generators[name]; ok {
		return fmt.Errorf("generator already registered: %s", name)
	}

	r.generators[name] = gen
	return nil
}

// Get retrieves a generator by name
func (r *Registry) Get(name string) (core.Generator, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	g, ok := r.generators[name]
	if !ok {
		return nil, fmt.Errorf("generator not found: %s", name)
	}
	return g, nil
}

// List returns all registered generator names
func (r *Registry) List() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var keys []string
	for k := range r.generators {
		keys = append(keys, k)
	}
	return keys
}

// Has checks if a generator is registered
func (r *Registry) Has(name string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, ok := r.generators[name]
	return ok
}
