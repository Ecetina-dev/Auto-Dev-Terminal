// Package constants provides constants and error messages for the Auto-Dev-Terminal system.
package constants

import "path/filepath"

// Application constants
const (
	AppName        = "auto-dev-terminal"
	AppDescription = "CLI tool for automated development environment setup"
	Version        = "1.0.0"
)

// Default paths
var (
	// DefaultBackupDir is the default directory for backups
	DefaultBackupDir = filepath.Join("~", ".auto-dev-terminal", "backups")

	// DefaultConfigDir is the default configuration directory
	DefaultConfigDir = filepath.Join("~", ".auto-dev-terminal", "config")

	// DefaultModulesDir is the default directory for module definitions
	DefaultModulesDir = filepath.Join("~", ".auto-dev-terminal", "modules")

	// ConfigFileName is the name of the main config file
	ConfigFileName = "config.yaml"

	// ManifestFileName is the name of the backup manifest
	ManifestFileName = "manifest.json"
)

// Error messages
var (
	// Detection errors
	ErrOSDetectionFailed     = "failed to detect operating system"
	ErrShellDetectionFailed  = "failed to detect shell"
	ErrPkgMgrDetectionFailed = "failed to detect package managers"
	ErrDistroDetectionFailed = "failed to detect Linux distribution"

	// Installation errors
	ErrInstallerNotAvailable = "package manager not available"
	ErrInstallationFailed    = "installation failed"
	ErrVerificationFailed    = "installation verification failed"
	ErrPackageNotFound       = "package not found"
	ErrAlreadyInstalled      = "package already installed"
	ErrRequiresSudo          = "installation requires elevated privileges"

	// Configuration errors
	ErrBackupFailed      = "failed to create backup"
	ErrRestoreFailed     = "failed to restore backup"
	ErrTemplateFailed    = "failed to render template"
	ErrWriteFailed       = "failed to write configuration file"
	ErrManifestCorrupted = "backup manifest is corrupted"

	// Module errors
	ErrModuleNotFound    = "module not found"
	ErrModuleInvalid     = "invalid module definition"
	ErrModuleUnsupported = "module not supported on current platform"
	ErrDependencyFailed  = "module dependency not satisfied"

	// General errors
	ErrFileNotFound     = "file not found"
	ErrPermissionDenied = "permission denied"
	ErrInvalidInput     = "invalid input"
	ErrCancelled        = "operation cancelled"
)

// Backup filename format
const (
	// BackupTimestampFormat is the format used for backup timestamps
	BackupTimestampFormat = "20060102_150405"

	// BackupFileExtension is the extension used for backup files
	BackupFileExtension = ".bak"
)

// Supported shells for module requirements
var SupportedShells = []string{
	"bash",
	"zsh",
	"fish",
	"powershell",
	"pwsh",
}

// Supported operating systems
var SupportedOS = []string{
	"windows",
	"darwin",
	"linux",
}

// Package manager binary names
var PackageManagerBinaries = map[string][]string{
	"homebrew":   {"brew"},
	"macports":   {"port"},
	"chocolatey": {"choco"},
	"scoop":      {"scoop"},
	"winget":     {"winget"},
	"apt":        {"apt"},
	"dnf":        {"dnf"},
	"yum":        {"yum"},
	"pacman":     {"pacman"},
	"zypper":     {"zypper"},
	"snap":       {"snap"},
	"flatpak":    {"flatpak"},
}
