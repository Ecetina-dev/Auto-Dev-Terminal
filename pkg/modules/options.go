// Package modules provides types for module configuration and operations.
package modules

import "github.com/auto-dev-terminal/auto-dev-terminal/internal/types"

// ModuleConfig represents the configuration for a module loaded from YAML/JSON.
type ModuleConfig struct {
	Name           string          `yaml:"name" json:"name"`
	DisplayName    string          `yaml:"display_name" json:"display_name"`
	Description    string          `yaml:"description" json:"description"`
	Version        string          `yaml:"version" json:"version"`
	Dependencies   []string        `yaml:"dependencies" json:"dependencies"`
	Install        InstallAction   `yaml:"install" json:"install"`
	Uninstall      UninstallAction `yaml:"uninstall" json:"uninstall"`
	CheckInstalled CheckAction     `yaml:"check_installed" json:"check_installed"`
	Requirements   []Requirement   `yaml:"requirements" json:"requirements"`
	Files          []FileConfig    `yaml:"files" json:"files"`
	Env            []EnvVar        `yaml:"env" json:"env"`
}

// InstallAction defines how to install the module.
type InstallAction struct {
	Type        string   `yaml:"type" json:"type"` // "package", "script", "git", "download"
	PackageMgr  string   `yaml:"package_manager" json:"package_manager"`
	Packages    []string `yaml:"packages" json:"packages"`
	ScriptURL   string   `yaml:"script_url" json:"script_url"`
	GitRepo     string   `yaml:"git_repo" json:"git_repo"`
	GitBranch   string   `yaml:"git_branch" json:"git_branch"`
	DownloadURL string   `yaml:"download_url" json:"download_url"`
	Commands    []string `yaml:"commands" json:"commands"`
	Destination string   `yaml:"destination" json:"destination"`
}

// UninstallAction defines how to uninstall the module.
type UninstallAction struct {
	Type       string   `yaml:"type" json:"type"`
	PackageMgr string   `yaml:"package_manager" json:"package_manager"`
	Packages   []string `yaml:"packages" json:"packages"`
	Commands   []string `yaml:"commands" json:"commands"`
	Files      []string `yaml:"files" json:"files"`
}

// CheckAction defines how to check if the module is installed.
type CheckAction struct {
	Type     string   `yaml:"type" json:"type"` // "command", "file", "directory"
	Commands []string `yaml:"commands" json:"commands"`
	Paths    []string `yaml:"paths" json:"paths"`
}

// Requirement defines a prerequisite for the module.
type Requirement struct {
	Type     string `yaml:"type" json:"type"` // "shell", "os", "command"
	Value    string `yaml:"value" json:"value"`
	Optional bool   `yaml:"optional" json:"optional"`
}

// FileConfig defines a file to be created or linked.
type FileConfig struct {
	Source      string            `yaml:"source" json:"source"`
	Destination string            `yaml:"destination" json:"destination"`
	Link        bool              `yaml:"link" json:"link"`
	Variables   map[string]string `yaml:"variables" json:"variables"`
}

// EnvVar defines an environment variable to set.
type EnvVar struct {
	Name  string `yaml:"name" json:"name"`
	Value string `yaml:"value" json:"value"`
}

// ModuleLoaderConfig is the root configuration structure for module definitions.
type ModuleLoaderConfig struct {
	Version string         `yaml:"version" json:"version"`
	Modules []ModuleConfig `yaml:"modules" json:"modules"`
}

// DetectResult holds the system detection results needed by modules.
type DetectResult struct {
	OS             types.OS
	Distro         types.Distro
	Shell          types.Shell
	PackageManager types.PackageManager
	HomeDir        string
	Username       string
}

// NewModuleOptions creates ModuleOptions from DetectResult.
func NewModuleOptions(detect DetectResult) *ModuleOptions {
	return &ModuleOptions{
		HomeDir:        detect.HomeDir,
		Shell:          detect.Shell,
		OS:             detect.OS,
		Distro:         detect.Distro,
		PackageManager: detect.PackageManager,
		Sudo:           false,
		Yes:            false,
		Verbose:        false,
		DryRun:         false,
		Force:          false,
	}
}
