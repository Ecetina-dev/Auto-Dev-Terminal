// Package modules provides a registry for managing and discovering modules.
package modules

import (
	"fmt"
	"sync"
)

// Registry manages the collection of available modules.
type Registry struct {
	modules map[string]Module
	mu      sync.RWMutex
}

// NewRegistry creates a new empty module registry.
func NewRegistry() *Registry {
	return &Registry{
		modules: make(map[string]Module),
	}
}

// Register adds a module to the registry.
// If a module with the same name already exists, it will be replaced with the new one.
// This makes registration idempotent - safe to call multiple times.
func (r *Registry) Register(m Module) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	name := m.Name()
	// Replace existing module (idempotent behavior)
	r.modules[name] = m
	return nil
}

// Unregister removes a module from the registry by name.
func (r *Registry) Unregister(name string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.modules, name)
}

// Get retrieves a module by name.
// Returns nil if the module is not found.
func (r *Registry) Get(name string) Module {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.modules[name]
}

// List returns all registered modules.
func (r *Registry) List() []Module {
	r.mu.RLock()
	defer r.mu.RUnlock()

	modules := make([]Module, 0, len(r.modules))
	for _, m := range r.modules {
		modules = append(modules, m)
	}
	return modules
}

// Names returns the names of all registered modules.
func (r *Registry) Names() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	names := make([]string, 0, len(r.modules))
	for name := range r.modules {
		names = append(names, name)
	}
	return names
}

// Len returns the number of registered modules.
func (r *Registry) Len() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.modules)
}

// Has checks if a module with the given name is registered.
func (r *Registry) Has(name string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, exists := r.modules[name]
	return exists
}

// GetDependencies returns all dependencies for a given module,
// including transitive dependencies.
func (r *Registry) GetDependencies(name string) ([]string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	visited := make(map[string]bool)
	var deps []string

	var visit func(modName string) error
	visit = func(modName string) error {
		if visited[modName] {
			return nil
		}
		visited[modName] = true

		mod, exists := r.modules[modName]
		if !exists {
			return fmt.Errorf("module %q not found", modName)
		}

		for _, dep := range mod.GetDependencies() {
			deps = append(deps, dep)
			if err := visit(dep); err != nil {
				return err
			}
		}
		return nil
	}

	if err := visit(name); err != nil {
		return nil, err
	}

	return deps, nil
}

// ResolveDependencies resolves the installation order for a module
// based on its dependencies.
func (r *Registry) ResolveDependencies(name string) ([]string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	type node struct {
		name       string
		dependents []string
	}

	// Build dependency graph
	graph := make(map[string]*node)
	for n, m := range r.modules {
		graph[n] = &node{name: n, dependents: m.GetDependencies()}
	}

	visited := make(map[string]bool)
	var order []string

	var visit func(n string) error
	visit = func(n string) error {
		if visited[n] {
			return nil
		}
		visited[n] = true

		// Visit dependencies first
		if nd, exists := graph[n]; exists {
			for _, dep := range nd.dependents {
				if _, depExists := graph[dep]; !depExists {
					return fmt.Errorf("missing dependency: %q for module %q", dep, n)
				}
				if err := visit(dep); err != nil {
					return err
				}
			}
		}

		order = append(order, n)
		return nil
	}

	if err := visit(name); err != nil {
		return nil, err
	}

	return order, nil
}

// GlobalRegistry is the default module registry used throughout the package.
var globalRegistry = NewRegistry()

// Register is a convenience function that registers a module with the global registry.
func Register(m Module) error {
	return globalRegistry.Register(m)
}

// Unregister is a convenience function that removes a module from the global registry.
func Unregister(name string) {
	globalRegistry.Unregister(name)
}

// Get is a convenience function that gets a module from the global registry.
func Get(name string) Module {
	return globalRegistry.Get(name)
}

// List is a convenience function that lists all modules in the global registry.
func List() []Module {
	return globalRegistry.List()
}

// Names is a convenience function that returns all module names from the global registry.
func Names() []string {
	return globalRegistry.Names()
}

// GetGlobalRegistry returns the global module registry.
func GetGlobalRegistry() *Registry {
	return globalRegistry
}

// SetGlobalRegistry replaces the global module registry with a custom one.
func SetGlobalRegistry(reg *Registry) {
	globalRegistry = reg
}
