// Package modules provides built-in module definitions and initialization.
package modules

import (
	"fmt"

	"github.com/auto-dev-terminal/auto-dev-terminal/internal/types"
)

// BuiltinModules returns a list of all built-in modules.
func BuiltinModules() []Module {
	return []Module{
		NewStarshipModule(),
		NewOhMyZshModule(),
		NewGitConfigModule(),
		NewFontsModule(),
	}
}

// GetModuleByName returns a built-in module by name.
func GetModuleByName(name string) Module {
	for _, mod := range BuiltinModules() {
		if mod.Name() == name {
			return mod
		}
	}
	return nil
}

// InitBuiltinModules registers all built-in modules with the global registry.
func InitBuiltinModules() error {
	var errors []string

	for _, mod := range BuiltinModules() {
		if err := Register(mod); err != nil {
			errors = append(errors, fmt.Sprintf("failed to register %s: %v", mod.Name(), err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("errors initializing modules: %s", errors)
	}

	return nil
}

// ModuleInfo contains metadata about a module.
type ModuleInfo struct {
	Name         string
	Description  string
	Version      string
	Dependencies []string
	Installed    bool
}

// GetModuleInfo returns detailed information about a module.
func GetModuleInfo(name string) (*ModuleInfo, error) {
	mod := Get(name)
	if mod == nil {
		// Try to find in built-in modules
		mod = GetModuleByName(name)
		if mod == nil {
			return nil, fmt.Errorf("module %q not found", name)
		}
	}

	installed, err := mod.IsInstalled()
	if err != nil {
		installed = false
	}

	return &ModuleInfo{
		Name:         mod.Name(),
		Description:  mod.Description(),
		Version:      mod.Version(),
		Dependencies: mod.GetDependencies(),
		Installed:    installed,
	}, nil
}

// InstallWithDeps installs a module and all its dependencies.
func InstallWithDeps(name string, opts *ModuleOptions) ([]*ModuleResult, error) {
	mod := Get(name)
	if mod == nil {
		return nil, fmt.Errorf("module %q not found", name)
	}

	// Get dependency order
	registry := GetGlobalRegistry()
	deps, err := registry.ResolveDependencies(name)
	if err != nil {
		// If resolve fails, just install the module directly
		deps = []string{name}
	}

	// If the requested module isn't in the resolved list, add it
	found := false
	for _, d := range deps {
		if d == name {
			found = true
			break
		}
	}
	if !found {
		deps = append(deps, name)
	}

	results := make([]*ModuleResult, 0, len(deps))

	for _, depName := range deps {
		depMod := Get(depName)
		if depMod == nil {
			results = append(results, &ModuleResult{
				Success: false,
				Module:  depName,
				Error:   "module not found in registry",
			})
			continue
		}

		// Check if already installed
		installed, err := depMod.IsInstalled()
		if err == nil && installed && !opts.Force {
			results = append(results, &ModuleResult{
				Success: true,
				Module:  depName,
				Output:  "already installed",
			})
			continue
		}

		// Install the dependency
		result := depMod.Install(opts)
		results = append(results, result)

		if !result.Success {
			// Stop on first failure
			return results, fmt.Errorf("failed to install %s: %s", depName, result.Error)
		}
	}

	return results, nil
}

// UninstallWithDeps uninstalls a module (optionally including dependents).
func UninstallWithDeps(name string, opts *ModuleOptions) ([]*ModuleResult, error) {
	mod := Get(name)
	if mod == nil {
		return nil, fmt.Errorf("module %q not found", name)
	}

	result := mod.Uninstall(opts)
	return []*ModuleResult{result}, nil
}

// ValidateModuleConfig validates a module configuration.
func ValidateModuleConfig(config ModuleConfig) error {
	if config.Name == "" {
		return fmt.Errorf("module name is required")
	}

	// Validate install type
	validInstallTypes := map[string]bool{
		"package":  true,
		"script":   true,
		"git":      true,
		"download": true,
	}

	if config.Install.Type != "" && !validInstallTypes[config.Install.Type] {
		return fmt.Errorf("invalid install type: %s", config.Install.Type)
	}

	// Validate requirements
	for _, req := range config.Requirements {
		validTypes := map[string]bool{
			"shell":    true,
			"os":       true,
			"command": true,
		}
		if !validTypes[req.Type] {
			return fmt.Errorf("invalid requirement type: %s", req.Type)
		}
	}

	return nil
}

// CreateModuleConfig creates a ModuleConfig from user input.
func CreateModuleConfig(name, description, version string, deps []string) ModuleConfig {
	return ModuleConfig{
		Name:         name,
		DisplayName:  name,
		Description:  description,
		Version:      version,
		Dependencies: deps,
	}
}

// ModuleStatus represents the installation status of a module.
type ModuleStatus struct {
	Name       string
	Version    string
	Installed  bool
	InstallErr string
}

// CheckAllModulesStatus checks the installation status of all registered modules.
func CheckAllModulesStatus() []ModuleStatus {
	modules := List()
	statuses := make([]ModuleStatus, 0, len(modules))

	for _, mod := range modules {
		installed, err := mod.IsInstalled()
		status := ModuleStatus{
			Name:      mod.Name(),
			Version:   mod.Version(),
			Installed: installed,
		}
		if err != nil {
			status.InstallErr = err.Error()
		}
		statuses = append(statuses, status)
	}

	return statuses
}

// GetModulesByCategory returns modules filtered by category (if supported).
// Currently returns all modules as we don't have categories defined.
func GetModulesByCategory(category string) []Module {
	// For now, return all modules
	// In the future, we could add category support to the Module interface
	return List()
}

// CreateModuleFromConfig creates a new module from a ModuleConfig.
func CreateModuleFromConfig(config ModuleConfig) (Module, error) {
	if err := ValidateModuleConfig(config); err != nil {
		return nil, err
	}

	return NewConfigurableModule(config), nil
}

// ModuleInstaller provides a convenient way to install modules with common options.
type ModuleInstaller struct {
	opts *ModuleOptions
}

// NewModuleInstaller creates a new module installer with default options.
func NewModuleInstaller() *ModuleInstaller {
	return &ModuleInstaller{
		opts: &ModuleOptions{
			Sudo:    false,
			Yes:     false,
			Verbose: false,
			DryRun:  false,
			Force:   false,
		},
	}
}

// WithSudo sets the sudo option.
func (mi *ModuleInstaller) WithSudo(sudo bool) *ModuleInstaller {
	mi.opts.Sudo = sudo
	return mi
}

// WithVerbose sets the verbose option.
func (mi *ModuleInstaller) WithVerbose(verbose bool) *ModuleInstaller {
	mi.opts.Verbose = verbose
	return mi
}

// WithForce sets the force option.
func (mi *ModuleInstaller) WithForce(force bool) *ModuleInstaller {
	mi.opts.Force = force
	return mi
}

// WithSystemInfo sets system information for the installer.
func (mi *ModuleInstaller) WithSystemInfo(info types.SystemInfo) *ModuleInstaller {
	mi.opts.HomeDir = info.HomeDir
	mi.opts.Shell = info.Shell
	mi.opts.OS = info.OS
	mi.opts.Distro = info.Distro
	mi.opts.PackageManager = info.PackageManagers[0]
	return mi
}

// Install installs a module by name.
func (mi *ModuleInstaller) Install(name string) *ModuleResult {
	mod := Get(name)
	if mod == nil {
		return &ModuleResult{
			Success: false,
			Module:  name,
			Error:   "module not found",
		}
	}

	return mod.Install(mi.opts)
}

// Uninstall uninstalls a module by name.
func (mi *ModuleInstaller) Uninstall(name string) *ModuleResult {
	mod := Get(name)
	if mod == nil {
		return &ModuleResult{
			Success: false,
			Module:  name,
			Error:   "module not found",
		}
	}

	return mod.Uninstall(mi.opts)
}

// GetOptions returns the current module options.
func (mi *ModuleInstaller) GetOptions() *ModuleOptions {
	return mi.opts
}
