package installer

import (
	"strings"

	"github.com/auto-dev-terminal/auto-dev-terminal/internal/types"
)

// BrewInstaller adapts Homebrew package manager.
type BrewInstaller struct {
	BaseInstaller
}

// NewBrewInstaller creates a new Homebrew installer.
func NewBrewInstaller() *BrewInstaller {
	return &BrewInstaller{
		BaseInstaller{
			NameStr:    "Homebrew",
			InstallCmd: "brew",
			ListCmd:    "brew",
			RemoveCmd:  "brew",
			UpdateCmd:  "brew",
			VersionCmd: "brew",
		},
	}
}

// Name returns the package manager name.
func (b *BrewInstaller) Name() string {
	return b.NameStr
}

// IsAvailable checks if Homebrew is installed.
func (b *BrewInstaller) IsAvailable() bool {
	return commandExists("brew")
}

// Install installs a package using Homebrew.
func (b *BrewInstaller) Install(pkg string, opts *types.InstallOptions) *types.InstallResult {
	args := []string{"install", pkg}
	if opts.Verbose {
		args = append(args, "-v")
	}

	output, err := b.runCommand("brew", args...)
	if err != nil {
		return &types.InstallResult{
			Success: false,
			Module:  pkg,
			Error:   err.Error(),
			Output:  output,
		}
	}

	version, _ := b.GetVersion(pkg)
	return &types.InstallResult{
		Success: true,
		Module:  pkg,
		Output:  output,
		Version: version,
	}
}

// Uninstall removes a package using Homebrew.
func (b *BrewInstaller) Uninstall(pkg string, opts *types.InstallOptions) *types.InstallResult {
	args := []string{"uninstall", pkg}
	if opts.Yes {
		args = append(args, "-y")
	}

	output, err := b.runCommand("brew", args...)
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

// IsInstalled checks if a package is installed via Homebrew.
func (b *BrewInstaller) IsInstalled(pkg string) (bool, error) {
	output, err := b.runCommand("brew", "list", pkg)
	if err != nil {
		return false, nil
	}
	return strings.Contains(output, pkg), nil
}

// Update updates a package using Homebrew.
func (b *BrewInstaller) Update(pkg string, opts *types.InstallOptions) *types.InstallResult {
	args := []string{"upgrade", pkg}

	output, err := b.runCommand("brew", args...)
	if err != nil {
		return &types.InstallResult{
			Success: false,
			Module:  pkg,
			Error:   err.Error(),
			Output:  output,
		}
	}

	version, _ := b.GetVersion(pkg)
	return &types.InstallResult{
		Success: true,
		Module:  pkg,
		Output:  output,
		Version: version,
	}
}

// GetVersion returns the version of an installed package.
func (b *BrewInstaller) GetVersion(pkg string) (string, error) {
	output, err := b.runCommand("brew", "info", pkg)
	if err != nil {
		return "", err
	}

	// Parse version from output like "node: 20.0.0"
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, pkg+":") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				return parts[1], nil
			}
		}
	}

	return "", nil
}
