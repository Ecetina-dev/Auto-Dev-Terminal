// Package detector provides comprehensive system detection capabilities
// for the Auto-Dev-Terminal CLI tool.
//
// This package detects:
//   - Operating System (Windows, macOS, Linux)
//   - Linux Distribution (Ubuntu, Fedora, Arch, etc.)
//   - Shell (bash, zsh, fish, PowerShell)
//   - Package Managers (brew, apt, dnf, pacman, etc.)
//
// Usage:
//
//	info, err := detector.Detect()
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("OS: %s\n", info.OS)
//	fmt.Printf("Shell: %s\n", info.Shell)
//	fmt.Printf("Package Managers: %v\n", info.PackageManagers)
package detector

import (
	"github.com/auto-dev-terminal/auto-dev-terminal/internal/types"
)

// Detector is the interface for system detection.
type Detector interface {
	Detect() (*types.SystemInfo, error)
}

// SimpleDetector is a basic detector implementation.
type SimpleDetector struct{}

// Detect performs system detection.
func (d *SimpleDetector) Detect() (*types.SystemInfo, error) {
	return Detect()
}

// NewDetector creates a new Detector instance.
func NewDetector() Detector {
	return &SimpleDetector{}
}

// DetectionResult represents the result of a detection operation.
type DetectionResult struct {
	Success bool
	SystemInfo *types.SystemInfo
	Error   error
}

// DetectAsync performs detection asynchronously.
func DetectAsync() chan DetectionResult {
	resultChan := make(chan DetectionResult, 1)

	go func() {
		info, err := Detect()
		resultChan <- DetectionResult{
			Success:    err == nil,
			SystemInfo: info,
			Error:      err,
		}
	}()

	return resultChan
}

// SupportedOS returns the list of supported operating systems.
func SupportedOS() []types.OS {
	return []types.OS{
		types.OSWindows,
		types.OSDarwin,
		types.OSLinux,
	}
}

// SupportedShells returns the list of supported shells.
func SupportedShells() []types.Shell {
	return []types.Shell{
		types.ShellBash,
		types.ShellZsh,
		types.ShellFish,
		types.ShellPowerShell,
		types.ShellPwsh,
		types.ShellCmd,
		types.ShellTcsh,
		types.ShellCsh,
		types.ShellAsh,
	}
}

// SupportedPackageManagers returns the list of supported package managers.
func SupportedPackageManagers() []types.PackageManager {
	return []types.PackageManager{
		types.PkgMgrHomebrew,
		types.PkgMgrMacPorts,
		types.PkgMgrChocolatey,
		types.PkgMgrScoop,
		types.PkgMgrWinget,
		types.PkgMgrAPT,
		types.PkgMgrDNF,
		types.PkgMgrYUM,
		types.PkgMgrPacman,
		types.PkgMgrZypper,
		types.PkgMgrSnap,
		types.PkgMgrFlatpak,
	}
}

// IsSupportedOS checks if the given OS is supported.
func IsSupportedOS(os types.OS) bool {
	for _, supported := range SupportedOS() {
		if os == supported {
			return true
		}
	}
	return false
}

// IsSupportedShell checks if the given shell is supported.
func IsSupportedShell(shell types.Shell) bool {
	for _, supported := range SupportedShells() {
		if shell == supported {
			return true
		}
	}
	return false
}

// IsSupportedPackageManager checks if the given package manager is supported.
func IsSupportedPackageManager(pm types.PackageManager) bool {
	for _, supported := range SupportedPackageManagers() {
		if pm == supported {
			return true
		}
	}
	return false
}
