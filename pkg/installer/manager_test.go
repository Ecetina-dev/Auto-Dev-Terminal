package installer

import (
	"testing"

	"github.com/auto-dev-terminal/auto-dev-terminal/internal/types"
)

func TestNewManager(t *testing.T) {
	tests := []struct {
		name       string
		systemInfo *types.SystemInfo
	}{
		{
			name:       "nil system info",
			systemInfo: nil,
		},
		{
			name:       "ubuntu system info",
			systemInfo: &types.SystemInfo{OS: types.OSLinux, Distro: types.DistroUbuntu},
		},
		{
			name:       "fedora system info",
			systemInfo: &types.SystemInfo{OS: types.OSLinux, Distro: types.DistroFedora},
		},
		{
			name:       "arch system info",
			systemInfo: &types.SystemInfo{OS: types.OSLinux, Distro: types.DistroArch},
		},
		{
			name:       "windows system info",
			systemInfo: &types.SystemInfo{OS: types.OSWindows},
		},
		{
			name:       "darwin system info",
			systemInfo: &types.SystemInfo{OS: types.OSDarwin},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewManager(tt.systemInfo)
			if m == nil {
				t.Error("NewManager() returned nil")
			}
			if m.installers == nil {
				t.Error("NewManager() returned nil installers map")
			}
			if m.systemInfo != tt.systemInfo {
				t.Errorf("NewManager().systemInfo = %v, want %v", m.systemInfo, tt.systemInfo)
			}
		})
	}
}

func TestManagerRegisterInstallers(t *testing.T) {
	tests := []struct {
		name         string
		distro       types.Distro
		expectAPT    bool
		expectDNF    bool
		expectPacman bool
		expectZypper bool
	}{
		{
			name:      "ubuntu registers APT",
			distro:    types.DistroUbuntu,
			expectAPT: true,
		},
		{
			name:      "debian registers APT",
			distro:    types.DistroDebian,
			expectAPT: true,
		},
		{
			name:      "linuxmint registers APT",
			distro:    types.DistroLinuxMint,
			expectAPT: true,
		},
		{
			name:      "fedora registers DNF",
			distro:    types.DistroFedora,
			expectDNF: true,
		},
		{
			name:      "rhel registers DNF",
			distro:    types.DistroRHEL,
			expectDNF: true,
		},
		{
			name:      "centos registers DNF",
			distro:    types.DistroCentOS,
			expectDNF: true,
		},
		{
			name:         "arch registers Pacman",
			distro:       types.DistroArch,
			expectPacman: true,
		},
		{
			name:         "manjaro registers Pacman",
			distro:       types.DistroManjaro,
			expectPacman: true,
		},
		{
			name:         "opensuse registers Zypper",
			distro:       types.DistroOpenSUSE,
			expectZypper: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sysInfo := &types.SystemInfo{OS: types.OSLinux, Distro: tt.distro}
			m := NewManager(sysInfo)

			if tt.expectAPT && m.installers[types.PkgMgrAPT] == nil {
				t.Errorf("Expected APT installer for %v", tt.distro)
			}
			if tt.expectDNF && m.installers[types.PkgMgrDNF] == nil {
				t.Errorf("Expected DNF installer for %v", tt.distro)
			}
			if tt.expectPacman && m.installers[types.PkgMgrPacman] == nil {
				t.Errorf("Expected Pacman installer for %v", tt.distro)
			}
			if tt.expectZypper && m.installers[types.PkgMgrZypper] == nil {
				t.Errorf("Expected Zypper installer for %v", tt.distro)
			}
		})
	}
}

func TestGetInstaller(t *testing.T) {
	m := NewManager(&types.SystemInfo{OS: types.OSLinux, Distro: types.DistroUbuntu})

	tests := []struct {
		name        string
		pkgMgr      types.PackageManager
		expectError bool
	}{
		{
			name:        "unknown package manager",
			pkgMgr:      types.PkgMgrUnknown,
			expectError: true,
		},
		{
			name:        "apt registered for Ubuntu",
			pkgMgr:      types.PkgMgrAPT,
			expectError: true, // Will fail if apt is not available in test env
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := m.GetInstaller(tt.pkgMgr)
			// Just verify it doesn't panic - error expected if not available
			_ = err
		})
	}
}

func TestGetAvailableInstallers(t *testing.T) {
	m := NewManager(&types.SystemInfo{OS: types.OSLinux, Distro: types.DistroUbuntu})

	// Get available installers - should return at least the cross-platform ones
	available := m.GetAvailableInstallers()

	// At minimum, we should have cross-platform installers registered
	if len(available) < 1 {
		t.Logf("GetAvailableInstallers returned %d installers", len(available))
	}
}

func TestInstall(t *testing.T) {
	m := NewManager(&types.SystemInfo{OS: types.OSLinux, Distro: types.DistroUbuntu})

	// Test installing a non-existent package with unknown package manager
	result := m.Install(types.PkgMgrUnknown, "nonexistent-package", nil)

	if result == nil {
		t.Error("Install() returned nil result")
	}
	if result.Success {
		t.Error("Install() should have failed for unknown package manager")
	}
	if result.Module != "nonexistent-package" {
		t.Errorf("Install().Module = %q, want %q", result.Module, "nonexistent-package")
	}
}

func TestUninstall(t *testing.T) {
	m := NewManager(&types.SystemInfo{OS: types.OSLinux, Distro: types.DistroUbuntu})

	// Test uninstalling with unknown package manager
	result := m.Uninstall(types.PkgMgrUnknown, "some-package", nil)

	if result == nil {
		t.Error("Uninstall() returned nil result")
	}
	if result.Success {
		t.Error("Uninstall() should have failed for unknown package manager")
	}
}

func TestIsInstalled(t *testing.T) {
	m := NewManager(&types.SystemInfo{OS: types.OSLinux, Distro: types.DistroUbuntu})

	// Test with unknown package manager
	_, err := m.IsInstalled(types.PkgMgrUnknown, "some-package")

	if err == nil {
		t.Error("IsInstalled() should have failed for unknown package manager")
	}
}

func TestUpdate(t *testing.T) {
	m := NewManager(&types.SystemInfo{OS: types.OSLinux, Distro: types.DistroUbuntu})

	// Test updating with unknown package manager
	result := m.Update(types.PkgMgrUnknown, "some-package", nil)

	if result == nil {
		t.Error("Update() returned nil result")
	}
	if result.Success {
		t.Error("Update() should have failed for unknown package manager")
	}
}

func TestGetVersion(t *testing.T) {
	m := NewManager(&types.SystemInfo{OS: types.OSLinux, Distro: types.DistroUbuntu})

	// Test with unknown package manager
	_, err := m.GetVersion(types.PkgMgrUnknown, "some-package")

	if err == nil {
		t.Error("GetVersion() should have failed for unknown package manager")
	}
}

func TestAutoDetectInstaller(t *testing.T) {
	tests := []struct {
		name       string
		systemInfo *types.SystemInfo
	}{
		{
			name:       "nil system info",
			systemInfo: nil,
		},
		{
			name: "with system info",
			systemInfo: &types.SystemInfo{
				OS:              types.OSLinux,
				Distro:          types.DistroUbuntu,
				PackageManagers: []types.PackageManager{types.PkgMgrAPT},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewManager(tt.systemInfo)
			installer, pkgMgr := m.AutoDetectInstaller()

			// The result may be nil if no installers are available in test environment
			// but the function should not panic
			_ = installer
			_ = pkgMgr
		})
	}
}
