package detector

import (
	"os/exec"
	"strings"

	"github.com/auto-dev-terminal/auto-dev-terminal/internal/types"
)

// DetectPackageManagers returns a list of available package managers on the system.
func DetectPackageManagers(os types.OS, distro types.Distro) []types.PackageManager {
	var managers []types.PackageManager

	// Always check for these regardless of OS
	allManagers := getAllPackageManagers()

	for _, m := range allManagers {
		if isPackageManagerAvailable(m, os, distro) {
			managers = append(managers, m)
		}
	}

	return managers
}

// getAllPackageManagers returns all supported package managers.
func getAllPackageManagers() []types.PackageManager {
	return []types.PackageManager{
		types.PkgMgrHomebrew,
		types.PkgMgrMacPorts,
		types.PkgMgrChocolatey,
		types.PkgMgrScoop,
		types.PkgMgrWinget,
		types.PkgMgrAPT,
		types.PkgMgrDNF,
		types.PkgMgrYUM,
		types.PkgMgrPacman,
		types.PkgMgrZypper,
		types.PkgMgrSnap,
		types.PkgMgrFlatpak,
	}
}

// isPackageManagerAvailable checks if a specific package manager is available.
func isPackageManagerAvailable(m types.PackageManager, os types.OS, distro types.Distro) bool {
	binary := getPackageManagerBinary(m)
	if binary == "" {
		return false
	}

	if !commandExists(binary) {
		return false
	}

	// Verify the package manager is actually executable (not just present)
	if !verifyPackageManager(m, binary) {
		return false
	}

	return true
}

// getPackageManagerBinary returns the binary name for a package manager.
func getPackageManagerBinary(m types.PackageManager) string {
	switch m {
	case types.PkgMgrHomebrew:
		return "brew"
	case types.PkgMgrMacPorts:
		return "port"
	case types.PkgMgrChocolatey:
		return "choco"
	case types.PkgMgrScoop:
		return "scoop"
	case types.PkgMgrWinget:
		return "winget"
	case types.PkgMgrAPT:
		return "apt"
	case types.PkgMgrDNF:
		return "dnf"
	case types.PkgMgrYUM:
		return "yum"
	case types.PkgMgrPacman:
		return "pacman"
	case types.PkgMgrZypper:
		return "zypper"
	case types.PkgMgrSnap:
		return "snap"
	case types.PkgMgrFlatpak:
		return "flatpak"
	default:
		return ""
	}
}

// verifyPackageManager verifies the package manager is actually working.
func verifyPackageManager(m types.PackageManager, binary string) bool {
	// Try to get version to verify it's working
	var cmd *exec.Cmd

	switch m {
	case types.PkgMgrHomebrew:
		cmd = exec.Command(binary, "--version")
	case types.PkgMgrMacPorts:
		cmd = exec.Command(binary, "--version")
	case types.PkgMgrChocolatey:
		cmd = exec.Command(binary, "--version")
	case types.PkgMgrScoop:
		cmd = exec.Command(binary, "--version")
	case types.PkgMgrWinget:
		cmd = exec.Command(binary, "--version")
	case types.PkgMgrAPT:
		// apt doesn't have --version on all systems, try -v
		cmd = exec.Command(binary, "-v")
	case types.PkgMgrDNF:
		cmd = exec.Command(binary, "--version")
	case types.PkgMgrYUM:
		cmd = exec.Command(binary, "--version")
	case types.PkgMgrPacman:
		cmd = exec.Command(binary, "--version")
	case types.PkgMgrZypper:
		cmd = exec.Command(binary, "--version")
	case types.PkgMgrSnap:
		cmd = exec.Command(binary, "--version")
	case types.PkgMgrFlatpak:
		cmd = exec.Command(binary, "--version")
	default:
		return false
	}

	// Set a short timeout
	if cmd.Process != nil {
		// Just check if command exists and runs
		output, err := cmd.CombinedOutput()
		if err != nil {
			// Some package managers need initialization, but if they don't exist at all, this fails
			// Check if it's just a permission or init issue vs not found
			errStr := strings.ToLower(string(output))
			if strings.Contains(errStr, "not found") || strings.Contains(errStr, "command not found") {
				return false
			}
			// If there's output but also error, it might still be available
			return len(output) > 0
		}
		return true
	}

	return commandExists(binary)
}

// GetPackageManagerVersion returns the version of a package manager.
func GetPackageManagerVersion(m types.PackageManager) (string, error) {
	binary := getPackageManagerBinary(m)
	if binary == "" {
		return "", nil
	}

	var cmd *exec.Cmd

	switch m {
	case types.PkgMgrHomebrew:
		cmd = exec.Command(binary, "--version")
	case types.PkgMgrMacPorts:
		cmd = exec.Command(binary, "--version")
	case types.PkgMgrChocolatey:
		cmd = exec.Command(binary, "--version")
	case types.PkgMgrScoop:
		cmd = exec.Command(binary, "--version")
	case types.PkgMgrWinget:
		cmd = exec.Command(binary, "--version")
	case types.PkgMgrAPT:
		cmd = exec.Command(binary, "-v")
	case types.PkgMgrDNF:
		cmd = exec.Command(binary, "--version")
	case types.PkgMgrYUM:
		cmd = exec.Command(binary, "--version")
	case types.PkgMgrPacman:
		cmd = exec.Command(binary, "--version")
	case types.PkgMgrZypper:
		cmd = exec.Command(binary, "--version")
	case types.PkgMgrSnap:
		cmd = exec.Command(binary, "--version")
	case types.PkgMgrFlatpak:
		cmd = exec.Command(binary, "--version")
	default:
		return "", nil
	}

	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	// Parse version from output
	return parseVersionOutput(string(output)), nil
}

// parseVersionOutput extracts version string from command output.
func parseVersionOutput(output string) string {
	lines := strings.Split(output, "\n")
	if len(lines) > 0 {
		// First line often contains version info
		// Remove common prefixes
		version := strings.TrimPrefix(lines[0], "v")
		parts := strings.Fields(version)
		if len(parts) > 0 {
			return parts[0]
		}
	}
	return strings.TrimSpace(output)
}

// GetPreferredPackageManager returns the preferred package manager for the system.
func GetPreferredPackageManager(os types.OS, distro types.Distro, managers []types.PackageManager) types.PackageManager {
	if len(managers) == 0 {
		return types.PkgMgrUnknown
	}

	// Priority order based on OS and distro
	var preferredOrder []types.PackageManager

	switch os {
	case types.OSDarwin:
		preferredOrder = []types.PackageManager{types.PkgMgrHomebrew, types.PkgMgrMacPorts}
	case types.OSWindows:
		preferredOrder = []types.PackageManager{types.PkgMgrWinget, types.PkgMgrChocolatey, types.PkgMgrScoop}
	case types.OSLinux:
		preferredOrder = getLinuxPreferredOrder(distro)
	}

	// Find first available in preference order
	for _, p := range preferredOrder {
		for _, m := range managers {
			if m == p {
				return m
			}
		}
	}

	// Return first available if none in preferred order
	return managers[0]
}

// getLinuxPreferredOrder returns preferred package managers for a Linux distro.
func getLinuxPreferredOrder(d types.Distro) []types.PackageManager {
	switch {
	case DistroSupportsAPT(d):
		return []types.PackageManager{types.PkgMgrAPT, types.PkgMgrSnap, types.PkgMgrFlatpak}
	case DistroSupportsDNF(d):
		return []types.PackageManager{types.PkgMgrDNF, types.PkgMgrYUM, types.PkgMgrSnap, types.PkgMgrFlatpak}
	case DistroSupportsPacman(d):
		return []types.PackageManager{types.PkgMgrPacman, types.PkgMgrSnap, types.PkgMgrFlatpak}
	case DistroSupportsZypper(d):
		return []types.PackageManager{types.PkgMgrZypper, types.PkgMgrSnap, types.PkgMgrFlatpak}
	default:
		// Generic order
		return []types.PackageManager{types.PkgMgrAPT, types.PkgMgrDNF, types.PkgMgrPacman, types.PkgMgrZypper}
	}
}
