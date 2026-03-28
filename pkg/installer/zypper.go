package installer

import (
	"strings"

	"github.com/auto-dev-terminal/auto-dev-terminal/internal/types"
)

// ZypperInstaller adapts Zypper (openSUSE/SLES) package manager.
type ZypperInstaller struct {
	BaseInstaller
}

// NewZypperInstaller creates a new Zypper installer.
func NewZypperInstaller() *ZypperInstaller {
	return &ZypperInstaller{
		BaseInstaller{
			NameStr:    "Zypper",
			InstallCmd: "zypper",
			ListCmd:    "zypper",
			RemoveCmd:  "zypper",
			UpdateCmd:  "zypper",
			VersionCmd: "rpm",
		},
	}
}

// Name returns the package manager name.
func (z *ZypperInstaller) Name() string {
	return z.NameStr
}

// IsAvailable checks if Zypper is available.
func (z *ZypperInstaller) IsAvailable() bool {
	return commandExists("zypper")
}

// Install installs a package using Zypper.
func (z *ZypperInstaller) Install(pkg string, opts *types.InstallOptions) *types.InstallResult {
	args := []string{"install", "-y"}
	if opts.Verbose {
		args = append(args, "-v")
	}
	args = append(args, pkg)

	output, err := z.runCommandWithSudo(opts.Sudo, "zypper", args...)
	if err != nil {
		return &types.InstallResult{
			Success: false,
			Module:  pkg,
			Error:   err.Error(),
			Output:  output,
		}
	}

	version, _ := z.GetVersion(pkg)
	return &types.InstallResult{
		Success: true,
		Module:  pkg,
		Output:  output,
		Version: version,
	}
}

// Uninstall removes a package using Zypper.
func (z *ZypperInstaller) Uninstall(pkg string, opts *types.InstallOptions) *types.InstallResult {
	args := []string{"remove", "-y", pkg}

	output, err := z.runCommandWithSudo(opts.Sudo, "zypper", args...)
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

// IsInstalled checks if a package is installed via Zypper.
func (z *ZypperInstaller) IsInstalled(pkg string) (bool, error) {
	output, err := z.runCommand("zypper", "search", "-i", pkg)
	if err != nil {
		return false, nil
	}
	return strings.Contains(output, "i+"), nil
}

// Update updates a package using Zypper.
func (z *ZypperInstaller) Update(pkg string, opts *types.InstallOptions) *types.InstallResult {
	args := []string{"update", "-y", pkg}

	output, err := z.runCommandWithSudo(opts.Sudo, "zypper", args...)
	if err != nil {
		return &types.InstallResult{
			Success: false,
			Module:  pkg,
			Error:   err.Error(),
			Output:  output,
		}
	}

	version, _ := z.GetVersion(pkg)
	return &types.InstallResult{
		Success: true,
		Module:  pkg,
		Output:  output,
		Version: version,
	}
}

// GetVersion returns the version of an installed package.
func (z *ZypperInstaller) GetVersion(pkg string) (string, error) {
	output, err := z.runCommand("rpm", "-q", "--queryformat", "%{VERSION}-%{RELEASE}", pkg)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(output), nil
}
