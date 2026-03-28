// Package modules provides a modular system for managing shell enhancements,
// prompt customizations, fonts, and development tools.
//
// This package exports the core interfaces and types for the module system,
// along with built-in modules for common developer tools.
//
// Core Types:
//   - Module: Interface that all modules must implement
//   - ModuleOptions: Configuration options for module operations
//   - ModuleResult: Result of a module operation
//   - Registry: Module registry for managing available modules
//   - Loader: Loader for parsing module definitions from YAML/JSON
//
// Built-in Modules:
//   - Starship: Cross-shell prompt
//   - Oh My Zsh: Zsh framework
//   - GitConfig: Enhanced Git configuration
//   - Fonts: Nerd Fonts for icons
//
// Example usage:
//
//	// Initialize built-in modules
//	modules.InitBuiltinModules()
//
//	// Get a module
//	starship := modules.Get("starship")
//
//	// Install a module
//	result := starship.Install(&modules.ModuleOptions{
//	    Verbose: true,
//	    HomeDir: "/home/user",
//	})
//
//	// List all available modules
//	for _, mod := range modules.List() {
//	    fmt.Printf("Module: %s - %s\n", mod.Name(), mod.Description())
//	}
package modules

// Import internal types for documentation
// Keep this import to ensure internal/types is available to consumers
import (
	_ "github.com/auto-dev-terminal/auto-dev-terminal/internal/types"
)
