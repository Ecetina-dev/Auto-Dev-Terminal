package installer

import (
	"strings"

	"github.com/auto-dev-terminal/auto-dev-terminal/internal/types"
)

// DnfInstaller adapts DNF (Fedora/RHEL) package manager.
type DnfInstaller struct {
	BaseInstaller
}

// NewDnfInstaller creates a new DNF installer.
func NewDnfInstaller() *DnfInstaller {
	return &DnfInstaller{
		BaseInstaller{
			NameStr:    "DNF",
			InstallCmd: "dnf",
			ListCmd:    "dnf",
			RemoveCmd:  "dnf",
			UpdateCmd:  "dnf",
			VersionCmd: "rpm",
		},
	}
}

// Name returns the package manager name.
func (d *DnfInstaller) Name() string {
	return d.NameStr
}

// IsAvailable checks if DNF is available.
func (d *DnfInstaller) IsAvailable() bool {
	return commandExists("dnf")
}

// Install installs a package using DNF.
func (d *DnfInstaller) Install(pkg string, opts *types.InstallOptions) *types.InstallResult {
	args := []string{"install", "-y"}
	if opts.Verbose {
		args = append(args, "-v")
	}
	args = append(args, pkg)

	output, err := d.runCommandWithSudo(opts.Sudo, "dnf", args...)
	if err != nil {
		return &types.InstallResult{
			Success: false,
			Module:  pkg,
			Error:   err.Error(),
			Output:  output,
		}
	}

	version, _ := d.GetVersion(pkg)
	return &types.InstallResult{
		Success: true,
		Module:  pkg,
		Output:  output,
		Version: version,
	}
}

// Uninstall removes a package using DNF.
func (d *DnfInstaller) Uninstall(pkg string, opts *types.InstallOptions) *types.InstallResult {
	args := []string{"remove", "-y", pkg}

	output, err := d.runCommandWithSudo(opts.Sudo, "dnf", args...)
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

// IsInstalled checks if a package is installed via DNF.
func (d *DnfInstaller) IsInstalled(pkg string) (bool, error) {
	output, err := d.runCommand("dnf", "list", "installed", pkg)
	if err != nil {
		return false, nil
	}
	return strings.Contains(output, "."+pkg), nil
}

// Update updates a package using DNF.
func (d *DnfInstaller) Update(pkg string, opts *types.InstallOptions) *types.InstallResult {
	args := []string{"update", "-y", pkg}

	output, err := d.runCommandWithSudo(opts.Sudo, "dnf", args...)
	if err != nil {
		return &types.InstallResult{
			Success: false,
			Module:  pkg,
			Error:   err.Error(),
			Output:  output,
		}
	}

	version, _ := d.GetVersion(pkg)
	return &types.InstallResult{
		Success: true,
		Module:  pkg,
		Output:  output,
		Version: version,
	}
}

// GetVersion returns the version of an installed package.
func (d *DnfInstaller) GetVersion(pkg string) (string, error) {
	output, err := d.runCommand("rpm", "-q", "--queryformat", "%{VERSION}-%{RELEASE}", pkg)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(output), nil
}
