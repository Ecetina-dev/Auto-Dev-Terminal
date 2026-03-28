package installer

import (
	"strings"

	"github.com/auto-dev-terminal/auto-dev-terminal/internal/types"
)

// AptInstaller adapts APT (Debian/Ubuntu) package manager.
type AptInstaller struct {
	BaseInstaller
}

// NewAptInstaller creates a new APT installer.
func NewAptInstaller() *AptInstaller {
	return &AptInstaller{
		BaseInstaller{
			NameStr:    "APT",
			InstallCmd: "apt",
			ListCmd:    "dpkg",
			RemoveCmd:  "apt",
			UpdateCmd:  "apt",
			VersionCmd: "dpkg",
		},
	}
}

// Name returns the package manager name.
func (a *AptInstaller) Name() string {
	return a.NameStr
}

// IsAvailable checks if APT is available.
func (a *AptInstaller) IsAvailable() bool {
	return commandExists("apt")
}

// Install installs a package using APT.
func (a *AptInstaller) Install(pkg string, opts *types.InstallOptions) *types.InstallResult {
	// First update package lists
	if !opts.DryRun {
		_, _ = a.runCommandWithSudo(opts.Sudo, "apt", "update")
	}

	args := []string{"install", "-y"}
	if opts.Verbose {
		args = append(args, "-o", "APT::Get::Show-Versions=true")
	}
	args = append(args, pkg)

	output, err := a.runCommandWithSudo(opts.Sudo, "apt", args...)
	if err != nil {
		return &types.InstallResult{
			Success: false,
			Module:  pkg,
			Error:   err.Error(),
			Output:  output,
		}
	}

	version, _ := a.GetVersion(pkg)
	return &types.InstallResult{
		Success: true,
		Module:  pkg,
		Output:  output,
		Version: version,
	}
}

// Uninstall removes a package using APT.
func (a *AptInstaller) Uninstall(pkg string, opts *types.InstallOptions) *types.InstallResult {
	args := []string{"remove", "-y", pkg}

	output, err := a.runCommandWithSudo(opts.Sudo, "apt", args...)
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

// IsInstalled checks if a package is installed via APT.
func (a *AptInstaller) IsInstalled(pkg string) (bool, error) {
	output, err := a.runCommand("dpkg", "-l", pkg)
	if err != nil {
		return false, nil
	}
	return strings.Contains(output, "ii "+pkg), nil
}

// Update updates package lists and upgrades a package.
func (a *AptInstaller) Update(pkg string, opts *types.InstallOptions) *types.InstallResult {
	// Update package lists first
	_, _ = a.runCommandWithSudo(opts.Sudo, "apt", "update")

	args := []string{"install", "--upgrade", "-y", pkg}

	output, err := a.runCommandWithSudo(opts.Sudo, "apt", args...)
	if err != nil {
		return &types.InstallResult{
			Success: false,
			Module:  pkg,
			Error:   err.Error(),
			Output:  output,
		}
	}

	version, _ := a.GetVersion(pkg)
	return &types.InstallResult{
		Success: true,
		Module:  pkg,
		Output:  output,
		Version: version,
	}
}

// GetVersion returns the version of an installed package.
func (a *AptInstaller) GetVersion(pkg string) (string, error) {
	output, err := a.runCommand("dpkg", "-l", pkg)
	if err != nil {
		return "", err
	}

	// Parse version from dpkg output like "ii  vim  2:8.2.3999-1ubuntu2.1"
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "ii ") && strings.Contains(line, " "+pkg+" ") {
			parts := strings.Fields(line)
			if len(parts) >= 3 {
				return parts[2], nil
			}
		}
	}

	return "", nil
}
