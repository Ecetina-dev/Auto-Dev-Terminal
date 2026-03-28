package installer

import (
	"strings"

	"github.com/auto-dev-terminal/auto-dev-terminal/internal/types"
)

// ChocoInstaller adapts Chocolatey package manager.
type ChocoInstaller struct {
	BaseInstaller
}

// NewChocoInstaller creates a new Chocolatey installer.
func NewChocoInstaller() *ChocoInstaller {
	return &ChocoInstaller{
		BaseInstaller{
			NameStr:    "Chocolatey",
			InstallCmd: "choco",
			ListCmd:    "choco",
			RemoveCmd:  "choco",
			UpdateCmd:  "choco",
			VersionCmd: "choco",
		},
	}
}

// Name returns the package manager name.
func (c *ChocoInstaller) Name() string {
	return c.NameStr
}

// IsAvailable checks if Chocolatey is installed.
func (c *ChocoInstaller) IsAvailable() bool {
	return commandExists("choco")
}

// Install installs a package using Chocolatey.
func (c *ChocoInstaller) Install(pkg string, opts *types.InstallOptions) *types.InstallResult {
	args := []string{"install", pkg, "-y"}
	if opts.Verbose {
		args = append(args, "-v")
	}
	if opts.Sudo {
		// Chocolatey doesn't need sudo on Windows, but we keep the flag for consistency
		_ = opts.Sudo
	}

	output, err := c.runCommand("choco", args...)
	if err != nil {
		return &types.InstallResult{
			Success: false,
			Module:  pkg,
			Error:   err.Error(),
			Output:  output,
		}
	}

	version, _ := c.GetVersion(pkg)
	return &types.InstallResult{
		Success: true,
		Module:  pkg,
		Output:  output,
		Version: version,
	}
}

// Uninstall removes a package using Chocolatey.
func (c *ChocoInstaller) Uninstall(pkg string, opts *types.InstallOptions) *types.InstallResult {
	args := []string{"uninstall", pkg, "-y"}

	output, err := c.runCommand("choco", args...)
	if err != nil {
		return &types.InstallResult{
			Success: false,
			Module:  pkg,
			Error:   err.Error(),
			Output:  output,
		}
	}

	return &types.InstallResult{
		Success: true,
		Module:  pkg,
		Output:  output,
	}
}

// IsInstalled checks if a package is installed via Chocolatey.
func (c *ChocoInstaller) IsInstalled(pkg string) (bool, error) {
	output, err := c.runCommand("choco", "list", "--local-only", pkg)
	if err != nil {
		return false, nil
	}
	return strings.Contains(output, pkg), nil
}

// Update updates a package using Chocolatey.
func (c *ChocoInstaller) Update(pkg string, opts *types.InstallOptions) *types.InstallResult {
	args := []string{"upgrade", pkg, "-y"}

	output, err := c.runCommand("choco", args...)
	if err != nil {
		return &types.InstallResult{
			Success: false,
			Module:  pkg,
			Error:   err.Error(),
			Output:  output,
		}
	}

	version, _ := c.GetVersion(pkg)
	return &types.InstallResult{
		Success: true,
		Module:  pkg,
		Output:  output,
		Version: version,
	}
}

// GetVersion returns the version of an installed package.
func (c *ChocoInstaller) GetVersion(pkg string) (string, error) {
	output, err := c.runCommand("choco", "list", "--local-only", pkg)
	if err != nil {
		return "", err
	}

	// Parse version from output like "node 20.0.0"
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, pkg+" ") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				return parts[1], nil
			}
		}
	}

	return "", nil
}
