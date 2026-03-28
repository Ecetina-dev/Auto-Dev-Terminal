package installer

import (
	"strings"

	"github.com/auto-dev-terminal/auto-dev-terminal/internal/types"
)

// PacmanInstaller adapts Pacman (Arch Linux) package manager.
type PacmanInstaller struct {
	BaseInstaller
}

// NewPacmanInstaller creates a new Pacman installer.
func NewPacmanInstaller() *PacmanInstaller {
	return &PacmanInstaller{
		BaseInstaller{
			NameStr:    "Pacman",
			InstallCmd: "pacman",
			ListCmd:    "pacman",
			RemoveCmd:  "pacman",
			UpdateCmd:  "pacman",
			VersionCmd: "pacman",
		},
	}
}

// Name returns the package manager name.
func (p *PacmanInstaller) Name() string {
	return p.NameStr
}

// IsAvailable checks if Pacman is available.
func (p *PacmanInstaller) IsAvailable() bool {
	return commandExists("pacman")
}

// Install installs a package using Pacman.
func (p *PacmanInstaller) Install(pkg string, opts *types.InstallOptions) *types.InstallResult {
	args := []string{"-S", "--noconfirm"}
	if opts.Verbose {
		args = append(args, "-v")
	}
	args = append(args, pkg)

	output, err := p.runCommandWithSudo(opts.Sudo, "pacman", args...)
	if err != nil {
		return &types.InstallResult{
			Success: false,
			Module:  pkg,
			Error:   err.Error(),
			Output:  output,
		}
	}

	version, _ := p.GetVersion(pkg)
	return &types.InstallResult{
		Success: true,
		Module:  pkg,
		Output:  output,
		Version: version,
	}
}

// Uninstall removes a package using Pacman.
func (p *PacmanInstaller) Uninstall(pkg string, opts *types.InstallOptions) *types.InstallResult {
	args := []string{"-R", "--noconfirm", pkg}

	output, err := p.runCommandWithSudo(opts.Sudo, "pacman", args...)
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

// IsInstalled checks if a package is installed via Pacman.
func (p *PacmanInstaller) IsInstalled(pkg string) (bool, error) {
	output, err := p.runCommand("pacman", "-Q", pkg)
	if err != nil {
		return false, nil
	}
	return strings.HasPrefix(output, pkg), nil
}

// Update updates a package using Pacman.
func (p *PacmanInstaller) Update(pkg string, opts *types.InstallOptions) *types.InstallResult {
	args := []string{"-S", "--noconfirm", "--overwrite", "*", pkg}

	output, err := p.runCommandWithSudo(opts.Sudo, "pacman", args...)
	if err != nil {
		return &types.InstallResult{
			Success: false,
			Module:  pkg,
			Error:   err.Error(),
			Output:  output,
		}
	}

	version, _ := p.GetVersion(pkg)
	return &types.InstallResult{
		Success: true,
		Module:  pkg,
		Output:  output,
		Version: version,
	}
}

// GetVersion returns the version of an installed package.
func (p *PacmanInstaller) GetVersion(pkg string) (string, error) {
	output, err := p.runCommand("pacman", "-Q", pkg)
	if err != nil {
		return "", err
	}

	// Parse version from output like "vim 9.0.1000-1"
	parts := strings.Fields(output)
	if len(parts) >= 2 {
		return parts[1], nil
	}

	return "", nil
}
