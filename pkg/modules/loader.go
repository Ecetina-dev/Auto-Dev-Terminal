// Package modules provides functionality for loading module definitions from files.
package modules

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Loader handles loading module definitions from YAML or JSON files.
type Loader struct {
	registry *Registry
}

// NewLoader creates a new module loader with the given registry.
func NewLoader(reg *Registry) *Loader {
	return &Loader{registry: reg}
}

// LoadFromFile loads module definitions from a YAML or JSON file.
func (l *Loader) LoadFromFile(path string) ([]ModuleConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading file %q: %w", path, err)
	}

	return l.Parse(data)
}

// Parse parses module definitions from YAML or JSON data.
func (l *Loader) Parse(data []byte) ([]ModuleConfig, error) {
	var config ModuleLoaderConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("parsing module config: %w", err)
	}

	return config.Modules, nil
}

// LoadAndRegister loads module definitions from a file and registers them
// with the loader's registry.
func (l *Loader) LoadAndRegister(path string) error {
	configs, err := l.LoadFromFile(path)
	if err != nil {
		return err
	}

	for _, cfg := range configs {
		mod := NewConfigurableModule(cfg)
		if err := l.registry.Register(mod); err != nil {
			return fmt.Errorf("registering module %q: %w", cfg.Name, err)
		}
	}

	return nil
}

// LoadFromDirectory loads all module definition files from a directory.
func (l *Loader) LoadFromDirectory(dirPath string) error {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return fmt.Errorf("reading directory %q: %w", dirPath, err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		// Only process .yaml and .yml files
		if len(name) < 5 {
			continue
		}
		ext := name[len(name)-5:]
		if ext != ".yaml" && ext != ".yml" && ext != ".json" {
			continue
		}

		path := dirPath + string(os.PathSeparator) + name
		if err := l.LoadAndRegister(path); err != nil {
			return fmt.Errorf("loading modules from %q: %w", path, err)
		}
	}

	return nil
}

// ConfigurableModule is a module that is configured from a ModuleConfig.
type ConfigurableModule struct {
	config ModuleConfig
}

// NewConfigurableModule creates a new configurable module from a ModuleConfig.
func NewConfigurableModule(config ModuleConfig) *ConfigurableModule {
	return &ConfigurableModule{config: config}
}

// Name returns the module name.
func (m *ConfigurableModule) Name() string {
	return m.config.Name
}

// Description returns the module description.
func (m *ConfigurableModule) Description() string {
	return m.config.Description
}

// Version returns the module version.
func (m *ConfigurableModule) Version() string {
	return m.config.Version
}

// GetDependencies returns the module dependencies.
func (m *ConfigurableModule) GetDependencies() []string {
	return m.config.Dependencies
}

// Install performs the installation of the module.
func (m *ConfigurableModule) Install(opts *ModuleOptions) *ModuleResult {
	// Check requirements first
	for _, req := range m.config.Requirements {
		if !req.Optional {
			if !m.checkRequirement(req, opts) {
				return &ModuleResult{
					Success: false,
					Module:  m.Name(),
					Error:   fmt.Sprintf("requirement not met: %s %s", req.Type, req.Value),
				}
			}
		}
	}

	// TODO: Implement actual installation logic based on config.Install
	return &ModuleResult{
		Success: true,
		Module:  m.Name(),
		Output:  "Installation not yet implemented for configurable modules",
	}
}

// Uninstall performs the uninstallation of the module.
func (m *ConfigurableModule) Uninstall(opts *ModuleOptions) *ModuleResult {
	// TODO: Implement actual uninstallation logic based on config.Uninstall
	return &ModuleResult{
		Success: true,
		Module:  m.Name(),
		Output:  "Uninstallation not yet implemented for configurable modules",
	}
}

// IsInstalled checks if the module is installed.
func (m *ConfigurableModule) IsInstalled() (bool, error) {
	// TODO: Implement actual check based on config.CheckInstalled
	return false, nil
}

// checkRequirement checks if a requirement is met.
func (m *ConfigurableModule) checkRequirement(req Requirement, opts *ModuleOptions) bool {
	switch req.Type {
	case "shell":
		return string(opts.Shell) == req.Value
	case "os":
		return string(opts.OS) == req.Value
	case "command":
		return commandExists(req.Value)
	default:
		return false
	}
}

// commandExists checks if a command is available in PATH.
func commandExists(cmd string) bool {
	// This is a placeholder - actual implementation would check PATH
	// For now, return false as we need proper OS-specific implementation
	return false
}
