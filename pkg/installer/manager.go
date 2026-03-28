package installer

import (
	"fmt"

	"github.com/auto-dev-terminal/auto-dev-terminal/internal/types"
)

// Manager provides a unified interface for all package managers.
type Manager struct {
	installers map[types.PackageManager]Installer
	systemInfo *types.SystemInfo
}

// NewManager creates a new installer manager.
func NewManager(sysInfo *types.SystemInfo) *Manager {
	m := &Manager{
		installers: make(map[types.PackageManager]Installer),
		systemInfo: sysInfo,
	}

	// Register all available installers
	m.registerInstallers()

	return m
}

// registerInstallers registers all supported package managers.
func (m *Manager) registerInstallers() {
	// Unix-like systems
	if m.systemInfo != nil {
		switch m.systemInfo.Distro {
		case types.DistroUbuntu, types.DistroDebian, types.DistroLinuxMint, types.DistroPop:
			m.installers[types.PkgMgrAPT] = NewAptInstaller()
		case types.DistroFedora, types.DistroRHEL, types.DistroCentOS, types.DistroRocky, types.DistroAlma:
			m.installers[types.PkgMgrDNF] = NewDnfInstaller()
		case types.DistroArch, types.DistroManjaro, types.DistroEndeavour:
			m.installers[types.PkgMgrPacman] = NewPacmanInstaller()
		case types.DistroOpenSUSE, types.DistroSLES:
			m.installers[types.PkgMgrZypper] = NewZypperInstaller()
		}
	}

	// Cross-platform installers (available on multiple OS)
	m.installers[types.PkgMgrHomebrew] = NewBrewInstaller()
	m.installers[types.PkgMgrChocolatey] = NewChocoInstaller()
	m.installers[types.PkgMgrScoop] = NewScoopInstaller()
	m.installers[types.PkgMgrWinget] = NewWingetInstaller()
}

// GetInstaller returns the installer for a specific package manager.
func (m *Manager) GetInstaller(pkgMgr types.PackageManager) (Installer, error) {
	installer, ok := m.installers[pkgMgr]
	if !ok {
		return nil, fmt.Errorf("unsupported package manager: %s", pkgMgr)
	}

	if !installer.IsAvailable() {
		return nil, fmt.Errorf("package manager %s is not available", pkgMgr)
	}

	return installer, nil
}

// GetAvailableInstallers returns all available installers for the current system.
func (m *Manager) GetAvailableInstallers() []Installer {
	var available []Installer
	for _, installer := range m.installers {
		if installer.IsAvailable() {
			available = append(available, installer)
		}
	}
	return available
}

// Install installs a package using the specified package manager.
func (m *Manager) Install(pkgMgr types.PackageManager, pkg string, opts *types.InstallOptions) *types.InstallResult {
	installer, err := m.GetInstaller(pkgMgr)
	if err != nil {
		return &types.InstallResult{
			Success: false,
			Module:  pkg,
			Error:   err.Error(),
		}
	}

	return installer.Install(pkg, opts)
}

// Uninstall removes a package using the specified package manager.
func (m *Manager) Uninstall(pkgMgr types.PackageManager, pkg string, opts *types.InstallOptions) *types.InstallResult {
	installer, err := m.GetInstaller(pkgMgr)
	if err != nil {
		return &types.InstallResult{
			Success: false,
			Module:  pkg,
			Error:   err.Error(),
		}
	}

	return installer.Uninstall(pkg, opts)
}

// IsInstalled checks if a package is installed using the specified package manager.
func (m *Manager) IsInstalled(pkgMgr types.PackageManager, pkg string) (bool, error) {
	installer, err := m.GetInstaller(pkgMgr)
	if err != nil {
		return false, err
	}

	return installer.IsInstalled(pkg)
}

// Update updates a package using the specified package manager.
func (m *Manager) Update(pkgMgr types.PackageManager, pkg string, opts *types.InstallOptions) *types.InstallResult {
	installer, err := m.GetInstaller(pkgMgr)
	if err != nil {
		return &types.InstallResult{
			Success: false,
			Module:  pkg,
			Error:   err.Error(),
		}
	}

	return installer.Update(pkg, opts)
}

// GetVersion returns the version of an installed package.
func (m *Manager) GetVersion(pkgMgr types.PackageManager, pkg string) (string, error) {
	installer, err := m.GetInstaller(pkgMgr)
	if err != nil {
		return "", err
	}

	return installer.GetVersion(pkg)
}

// AutoDetectInstaller returns the best available installer for the current system.
func (m *Manager) AutoDetectInstaller() (Installer, types.PackageManager) {
	// Priority order based on package manager detection from system info
	if m.systemInfo != nil {
		for _, pm := range m.systemInfo.PackageManagers {
			if installer, ok := m.installers[pm]; ok && installer.IsAvailable() {
				return installer, pm
			}
		}
	}

	// Fallback: try to find any available installer
	priorityOrder := []types.PackageManager{
		types.PkgMgrAPT,
		types.PkgMgrDNF,
		types.PkgMgrPacman,
		types.PkgMgrHomebrew,
		types.PkgMgrWinget,
		types.PkgMgrScoop,
		types.PkgMgrChocolatey,
	}

	for _, pm := range priorityOrder {
		if installer, ok := m.installers[pm]; ok && installer.IsAvailable() {
			return installer, pm
		}
	}

	return nil, types.PkgMgrUnknown
}
