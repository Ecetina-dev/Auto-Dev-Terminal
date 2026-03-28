// Package modules provides a modular system for managing shell enhancements,
// prompt customizations, and development tools.
package modules

import "github.com/auto-dev-terminal/auto-dev-terminal/internal/types"

// ModuleOptions contains configuration options for module operations.
type ModuleOptions struct {
	Sudo    bool
	Yes     bool
	Verbose bool
	DryRun  bool
	Force   bool
	// HomeDir is the user's home directory
	HomeDir string
	// Shell is the detected shell type
	Shell types.Shell
	// OS is the detected operating system
	OS types.OS
	// Distro is the detected Linux distribution (if applicable)
	Distro types.Distro
	// PackageManager is the preferred package manager
	PackageManager types.PackageManager
}

// ModuleResult contains the outcome of a module operation.
type ModuleResult struct {
	Success bool
	Module  string
	Error   string
	Output  string
	Version string
}

// Module defines the interface that all modules must implement.
type Module interface {
	// Name returns the unique identifier of the module.
	Name() string

	// Description returns a human-readable description of the module.
	Description() string

	// Version returns the current version of the module.
	Version() string

	// Install performs the installation of the module.
	// Returns a ModuleResult indicating success or failure.
	Install(opts *ModuleOptions) *ModuleResult

	// Uninstall removes the module from the system.
	// Returns a ModuleResult indicating success or failure.
	Uninstall(opts *ModuleOptions) *ModuleResult

	// IsInstalled checks if the module is currently installed.
	// Returns installation status and any error encountered.
	IsInstalled() (bool, error)

	// GetDependencies returns a list of module names that must be
	// installed before this module can be installed.
	GetDependencies() []string
}

// BaseModule provides common functionality for all modules.
type BaseModule struct {
	name        string
	description string
	version     string
	dependencies []string
}

// NewBaseModule creates a new base module with the given properties.
func NewBaseModule(name, description, version string, dependencies []string) *BaseModule {
	return &BaseModule{
		name:         name,
		description:  description,
		version:      version,
		dependencies: dependencies,
	}
}

// Name returns the module name.
func (m *BaseModule) Name() string {
	return m.name
}

// Description returns the module description.
func (m *BaseModule) Description() string {
	return m.description
}

// Version returns the module version.
func (m *BaseModule) Version() string {
	return m.version
}

// GetDependencies returns the module dependencies.
func (m *BaseModule) GetDependencies() []string {
	return m.dependencies
}
