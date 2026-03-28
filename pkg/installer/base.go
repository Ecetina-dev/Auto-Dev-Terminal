// Package installer provides cross-platform package manager adapters.
package installer

import (
	"fmt"
	"os/exec"

	"github.com/auto-dev-terminal/auto-dev-terminal/internal/types"
)

// Installer defines the interface for all package manager adapters.
type Installer interface {
	// Name returns the display name of the package manager
	Name() string

	// IsAvailable checks if the package manager is available on the system
	IsAvailable() bool

	// Install installs a package
	Install(pkg string, opts *types.InstallOptions) *types.InstallResult

	// Uninstall removes a package
	Uninstall(pkg string, opts *types.InstallOptions) *types.InstallResult

	// IsInstalled checks if a package is already installed
	IsInstalled(pkg string) (bool, error)

	// Update updates a package to the latest version
	Update(pkg string, opts *types.InstallOptions) *types.InstallResult

	// GetVersion returns the version of an installed package
	GetVersion(pkg string) (string, error)
}

// BaseInstaller provides common functionality for all installers.
type BaseInstaller struct {
	NameStr    string
	InstallCmd string
	ListCmd    string
	RemoveCmd  string
	UpdateCmd  string
	VersionCmd string
}

// runCommand executes a command and returns the output.
func (b *BaseInstaller) runCommand(cmd string, args ...string) (string, error) {
	execCmd := exec.Command(cmd, args...)
	output, err := execCmd.CombinedOutput()
	if err != nil {
		return string(output), fmt.Errorf("%s: %w", cmd, err)
	}
	return string(output), nil
}

// runCommandWithSudo executes a command with sudo if needed.
func (b *BaseInstaller) runCommandWithSudo(sudo bool, cmd string, args ...string) (string, error) {
	if sudo {
		args = append([]string{cmd}, args...)
		cmd = "sudo"
	}
	return b.runCommand(cmd, args...)
}

// commandExists checks if a command is available in PATH.
func commandExists(cmd string) bool {
	execCmd := exec.Command("command", "-v", cmd)
	return execCmd.Run() == nil
}
