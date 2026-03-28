package installer

import (
	"strings"

	"github.com/auto-dev-terminal/auto-dev-terminal/internal/types"
)

// ScoopInstaller adapts Scoop Windows package manager.
type ScoopInstaller struct {
	BaseInstaller
}

// NewScoopInstaller creates a new Scoop installer.
func NewScoopInstaller() *ScoopInstaller {
	return &ScoopInstaller{
		BaseInstaller{
			NameStr:    "Scoop",
			InstallCmd: "scoop",
			ListCmd:    "scoop",
			RemoveCmd:  "scoop",
			UpdateCmd:  "scoop",
			VersionCmd: "scoop",
		},
	}
}

// Name returns the package manager name.
func (s *ScoopInstaller) Name() string {
	return s.NameStr
}

// IsAvailable checks if Scoop is installed.
func (s *ScoopInstaller) IsAvailable() bool {
	return commandExists("scoop")
}

// Install installs a package using Scoop.
func (s *ScoopInstaller) Install(pkg string, opts *types.InstallOptions) *types.InstallResult {
	args := []string{"install", pkg}
	if opts.Yes {
		args = append(args, "-y")
	}
	if opts.Verbose {
		args = append(args, "-v")
	}

	output, err := s.runCommand("scoop", args...)
	if err != nil {
		return &types.InstallResult{
			Success: false,
			Module:  pkg,
			Error:   err.Error(),
			Output:  output,
		}
	}

	version, _ := s.GetVersion(pkg)
	return &types.InstallResult{
		Success: true,
		Module:  pkg,
		Output:  output,
		Version: version,
	}
}

// Uninstall removes a package using Scoop.
func (s *ScoopInstaller) Uninstall(pkg string, opts *types.InstallOptions) *types.InstallResult {
	args := []string{"uninstall", pkg}
	if opts.Yes {
		args = append(args, "-y")
	}

	output, err := s.runCommand("scoop", args...)
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

// IsInstalled checks if a package is installed via Scoop.
func (s *ScoopInstaller) IsInstalled(pkg string) (bool, error) {
	output, err := s.runCommand("scoop", "list", pkg)
	if err != nil {
		return false, nil
	}
	// Look for the package in the list output
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.Contains(line, pkg) && !strings.HasPrefix(strings.TrimSpace(line), "---") {
			return true, nil
		}
	}
	return false, nil
}

// Update updates a package using Scoop.
func (s *ScoopInstaller) Update(pkg string, opts *types.InstallOptions) *types.InstallResult {
	args := []string{"update", pkg}
	if opts.Verbose {
		args = append(args, "-v")
	}

	output, err := s.runCommand("scoop", args...)
	if err != nil {
		return &types.InstallResult{
			Success: false,
			Module:  pkg,
			Error:   err.Error(),
			Output:  output,
		}
	}

	version, _ := s.GetVersion(pkg)
	return &types.InstallResult{
		Success: true,
		Module:  pkg,
		Output:  output,
		Version: version,
	}
}

// GetVersion returns the version of an installed package.
func (s *ScoopInstaller) GetVersion(pkg string) (string, error) {
	output, err := s.runCommand("scoop", "info", pkg)
	if err != nil {
		return "", err
	}

	// Parse version from output like "Version: 20.0.0"
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "Version:") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				return parts[1], nil
			}
		}
	}

	return "", nil
}
