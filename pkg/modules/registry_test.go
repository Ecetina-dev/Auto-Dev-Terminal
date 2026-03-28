package modules

import (
	"testing"
)

// Mock module for testing
type mockModule struct {
	name         string
	description  string
	version      string
	dependencies []string
	installed    bool
}

func (m *mockModule) Name() string              { return m.name }
func (m *mockModule) Description() string       { return m.description }
func (m *mockModule) Version() string           { return m.version }
func (m *mockModule) GetDependencies() []string { return m.dependencies }
func (m *mockModule) Install(opts *ModuleOptions) *ModuleResult {
	return &ModuleResult{Success: true, Module: m.name}
}
func (m *mockModule) Uninstall(opts *ModuleOptions) *ModuleResult {
	return &ModuleResult{Success: true, Module: m.name}
}
func (m *mockModule) IsInstalled() (bool, error) {
	return m.installed, nil
}

func TestNewRegistry(t *testing.T) {
	r := NewRegistry()
	if r == nil {
		t.Fatal("NewRegistry() returned nil")
	}
	if r.modules == nil {
		t.Error("Registry.modules should be initialized")
	}
}

func TestRegistryRegister(t *testing.T) {
	r := NewRegistry()

	// Test registering a new module
	mod := &mockModule{
		name:        "test-module",
		description: "Test module",
		version:     "1.0.0",
	}

	err := r.Register(mod)
	if err != nil {
		t.Errorf("Register() error = %v", err)
	}

	// Test registering duplicate (now idempotent - should NOT fail)
	err = r.Register(mod)
	if err != nil {
		t.Errorf("Register() should NOT fail for duplicate module (idempotent): %v", err)
	}
}

func TestRegistryUnregister(t *testing.T) {
	r := NewRegistry()

	mod := &mockModule{
		name:        "test-module",
		description: "Test module",
		version:     "1.0.0",
	}

	// Register module
	if err := r.Register(mod); err != nil {
		t.Fatalf("Register() error = %v", err)
	}

	// Unregister
	r.Unregister("test-module")

	// Verify it's gone
	if r.Get("test-module") != nil {
		t.Error("Unregister() did not remove module")
	}
}

func TestRegistryGet(t *testing.T) {
	r := NewRegistry()

	mod := &mockModule{
		name:        "test-module",
		description: "Test module",
		version:     "1.0.0",
	}

	// Test get non-existent
	result := r.Get("nonexistent")
	if result != nil {
		t.Error("Get() should return nil for non-existent module")
	}

	// Register and get
	r.Register(mod)
	result = r.Get("test-module")
	if result == nil {
		t.Error("Get() returned nil for registered module")
	}
	if result.Name() != "test-module" {
		t.Errorf("Get() = %q, want %q", result.Name(), "test-module")
	}
}

func TestRegistryList(t *testing.T) {
	r := NewRegistry()

	// Empty list
	modules := r.List()
	if len(modules) != 0 {
		t.Errorf("List() on empty registry = %d, want 0", len(modules))
	}

	// Add modules
	mod1 := &mockModule{name: "module1", description: "Desc1", version: "1.0"}
	mod2 := &mockModule{name: "module2", description: "Desc2", version: "2.0"}
	r.Register(mod1)
	r.Register(mod2)

	modules = r.List()
	if len(modules) != 2 {
		t.Errorf("List() = %d, want 2", len(modules))
	}
}

func TestRegistryNames(t *testing.T) {
	r := NewRegistry()

	// Empty names
	names := r.Names()
	if len(names) != 0 {
		t.Errorf("Names() on empty registry = %d, want 0", len(names))
	}

	// Add modules
	mod1 := &mockModule{name: "module1", description: "Desc1", version: "1.0"}
	mod2 := &mockModule{name: "module2", description: "Desc2", version: "2.0"}
	r.Register(mod1)
	r.Register(mod2)

	names = r.Names()
	if len(names) != 2 {
		t.Errorf("Names() = %d, want 2", len(names))
	}

	// Check names contain expected values
	found := map[string]bool{}
	for _, n := range names {
		found[n] = true
	}
	if !found["module1"] || !found["module2"] {
		t.Error("Names() doesn't contain expected module names")
	}
}

func TestRegistryLen(t *testing.T) {
	r := NewRegistry()

	if r.Len() != 0 {
		t.Errorf("Len() on empty registry = %d, want 0", r.Len())
	}

	r.Register(&mockModule{name: "mod1", description: "Desc", version: "1.0"})
	if r.Len() != 1 {
		t.Errorf("Len() = %d, want 1", r.Len())
	}

	r.Register(&mockModule{name: "mod2", description: "Desc", version: "1.0"})
	if r.Len() != 2 {
		t.Errorf("Len() = %d, want 2", r.Len())
	}
}

func TestRegistryHas(t *testing.T) {
	r := NewRegistry()

	// Test non-existent
	if r.Has("nonexistent") {
		t.Error("Has() should return false for non-existent module")
	}

	// Register and test
	mod := &mockModule{name: "test-module", description: "Desc", version: "1.0"}
	r.Register(mod)

	if !r.Has("test-module") {
		t.Error("Has() should return true for registered module")
	}
}

func TestRegistryGetDependencies(t *testing.T) {
	r := NewRegistry()

	// Register modules with dependencies
	mod1 := &mockModule{
		name:         "module1",
		description:  "Module 1",
		version:      "1.0",
		dependencies: []string{},
	}
	mod2 := &mockModule{
		name:         "module2",
		description:  "Module 2",
		version:      "1.0",
		dependencies: []string{"module1"},
	}
	mod3 := &mockModule{
		name:         "module3",
		description:  "Module 3",
		version:      "1.0",
		dependencies: []string{"module2"},
	}

	r.Register(mod1)
	r.Register(mod2)
	r.Register(mod3)

	// Get dependencies for module3
	deps, err := r.GetDependencies("module3")
	if err != nil {
		t.Errorf("GetDependencies() error = %v", err)
	}

	// Should have module1 and module2 (and potentially duplicates of module2)
	if len(deps) < 2 {
		t.Errorf("GetDependencies() = %v, want at least 2 deps", deps)
	}

	// Test with non-existent module
	_, err = r.GetDependencies("nonexistent")
	if err == nil {
		t.Error("GetDependencies() should fail for non-existent module")
	}
}

func TestRegistryResolveDependencies(t *testing.T) {
	r := NewRegistry()

	// Register modules in wrong order
	mod3 := &mockModule{
		name:         "module3",
		description:  "Module 3",
		version:      "1.0",
		dependencies: []string{"module2", "module1"},
	}
	mod1 := &mockModule{
		name:         "module1",
		description:  "Module 1",
		version:      "1.0",
		dependencies: []string{},
	}
	mod2 := &mockModule{
		name:         "module2",
		description:  "Module 2",
		version:      "1.0",
		dependencies: []string{"module1"},
	}

	r.Register(mod3)
	r.Register(mod1)
	r.Register(mod2)

	// Resolve dependencies - should return correct order
	order, err := r.ResolveDependencies("module3")
	if err != nil {
		t.Errorf("ResolveDependencies() error = %v", err)
	}

	// module1 should come before module2, module2 before module3
	pos := map[string]int{}
	for i, name := range order {
		pos[name] = i
	}

	if pos["module1"] > pos["module2"] {
		t.Error("module1 should come before module2")
	}
	if pos["module2"] > pos["module3"] {
		t.Error("module2 should come before module3")
	}

	// Test with missing dependency
	mod4 := &mockModule{
		name:         "module4",
		description:  "Module 4",
		version:      "1.0",
		dependencies: []string{"nonexistent"},
	}
	r.Register(mod4)

	_, err = r.ResolveDependencies("module4")
	if err == nil {
		t.Error("ResolveDependencies() should fail for missing dependency")
	}
}

// Test global registry functions
func TestGlobalRegistry(t *testing.T) {
	// Save original global registry
	original := globalRegistry
	defer func() { globalRegistry = original }()

	// Create new registry and set it
	r := NewRegistry()
	SetGlobalRegistry(r)

	if globalRegistry != r {
		t.Error("SetGlobalRegistry() did not set global registry")
	}

	// Test convenience functions
	if GetGlobalRegistry() != r {
		t.Error("GetGlobalRegistry() did not return set registry")
	}
}
