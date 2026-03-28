// Package detector provides system detection capabilities for Auto-Dev-Terminal.
package detector

import (
	"runtime"

	"github.com/auto-dev-terminal/auto-dev-terminal/internal/types"
)

// DetectOS returns the current operating system using runtime.GOOS.
func DetectOS() types.OS {
	os := runtime.GOOS
	switch os {
	case "windows":
		return types.OSWindows
	case "darwin":
		return types.OSDarwin
	case "linux":
		return types.OSLinux
	default:
		return types.OS(os)
	}
}

// DetectArch returns the system architecture.
func DetectArch() string {
	return runtime.GOARCH
}

// IsWindows returns true if the OS is Windows.
func IsWindows() bool {
	return runtime.GOOS == "windows"
}

// IsDarwin returns true if the OS is macOS.
func IsDarwin() bool {
	return runtime.GOOS == "darwin"
}

// IsLinux returns true if the OS is Linux.
func IsLinux() bool {
	return runtime.GOOS == "linux"
}
