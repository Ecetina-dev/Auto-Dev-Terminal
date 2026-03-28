package installer

import (
	"strings"

	"github.com/auto-dev-terminal/auto-dev-terminal/internal/types"
)

// WingetInstaller adapts Windows Package Manager (winget).
type WingetInstaller struct {
	BaseInstaller
}

// NewWingetInstaller creates a new Winget installer.
func NewWingetInstaller() *WingetInstaller {
	return &WingetInstaller{
		BaseInstaller{
			NameStr:    "Winget",
			InstallCmd: "winget",
			ListCmd:    "winget",
			RemoveCmd:  "winget",
			UpdateCmd:  "winget",
			VersionCmd: "winget",
		},
	}
}

// Name returns the package manager name.
func (w *WingetInstaller) Name() string {
	return w.NameStr
}

// IsAvailable checks if Winget is installed.
func (w *WingetInstaller) IsAvailable() bool {
	return commandExists("winget")
}

// Install installs a package using Winget.
func (w *WingetInstaller) Install(pkg string, opts *types.InstallOptions) *types.InstallResult {
	args := []string{"install", pkg, "--silent"}
	if opts.Yes {
		args = append(args, "--accept-package-agreements", "--accept-source-agreements")
	}
	if opts.Verbose {
		args = append(args, "--verbose")
	}

	output, err := w.runCommand("winget", args...)
	if err != nil {
		return &types.InstallResult{
			Success: false,
			Module:  pkg,
			Error:   err.Error(),
			Output:  output,
		}
	}

	version, _ := w.GetVersion(pkg)
	return &types.InstallResult{
		Success: true,
		Module:  pkg,
		Output:  output,
		Version: version,
	}
}

// Uninstall removes a package using Winget.
func (w *WingetInstaller) Uninstall(pkg string, opts *types.InstallOptions) *types.InstallResult {
	args := []string{"uninstall", pkg, "--silent"}
	if opts.Yes {
		args = append(args, "--accept-package-agreements", "--accept-source-agreements")
	}

	output, err := w.runCommand("winget", args...)
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

// IsInstalled checks if a package is installed via Winget.
func (w *WingetInstaller) IsInstalled(pkg string) (bool, error) {
	output, err := w.runCommand("winget", "list", pkg)
	if err != nil {
		return false, nil
	}
	return strings.Contains(output, pkg), nil
}

// Update updates a package using Winget.
func (w *WingetInstaller) Update(pkg string, opts *types.InstallOptions) *types.InstallResult {
	args := []string{"upgrade", pkg, "--silent"}
	if opts.Yes {
		args = append(args, "--accept-package-agreements", "--accept-source-agreements")
	}
	if opts.Verbose {
		args = append(args, "--verbose")
	}

	output, err := w.runCommand("winget", args...)
	if err != nil {
		return &types.InstallResult{
			Success: false,
			Module:  pkg,
			Error:   err.Error(),
			Output:  output,
		}
	}

	version, _ := w.GetVersion(pkg)
	return &types.InstallResult{
		Success: true,
		Module:  pkg,
		Output:  output,
		Version: version,
	}
}

// GetVersion returns the version of an installed package.
func (w *WingetInstaller) GetVersion(pkg string) (string, error) {
	output, err := w.runCommand("winget", "show", pkg)
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
