package detector

import (
	"strings"
	"testing"

	"github.com/auto-dev-terminal/auto-dev-terminal/internal/types"
)

func TestGetPackageManagerBinary(t *testing.T) {
	tests := []struct {
		name          string
		pkgMgr        types.PackageManager
		expectedBin   string
	}{
		{
			name:        "homebrew",
			pkgMgr:      types.PkgMgrHomebrew,
			expectedBin: "brew",
		},
		{
			name:        "macports",
			pkgMgr:      types.PkgMgrMacPorts,
			expectedBin: "port",
		},
		{
			name:        "chocolatey",
			pkgMgr:      types.PkgMgrChocolatey,
			expectedBin: "choco",
		},
		{
			name:        "scoop",
			pkgMgr:      types.PkgMgrScoop,
			expectedBin: "scoop",
		},
		{
			name:        "winget",
			pkgMgr:      types.PkgMgrWinget,
			expectedBin: "winget",
		},
		{
			name:        "apt",
			pkgMgr:      types.PkgMgrAPT,
			expectedBin: "apt",
		},
		{
			name:        "dnf",
			pkgMgr:      types.PkgMgrDNF,
			expectedBin: "dnf",
		},
		{
			name:        "yum",
			pkgMgr:      types.PkgMgrYUM,
			expectedBin: "yum",
		},
		{
			name:        "pacman",
			pkgMgr:      types.PkgMgrPacman,
			expectedBin: "pacman",
		},
		{
			name:        "zypper",
			pkgMgr:      types.PkgMgrZypper,
			expectedBin: "zypper",
		},
		{
			name:        "snap",
			pkgMgr:      types.PkgMgrSnap,
			expectedBin: "snap",
		},
		{
			name:        "flatpak",
			pkgMgr:      types.PkgMgrFlatpak,
			expectedBin: "flatpak",
		},
		{
			name:        "unknown",
			pkgMgr:      types.PkgMgrUnknown,
			expectedBin: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getPackageManagerBinary(tt.pkgMgr)
			if result != tt.expectedBin {
				t.Errorf("getPackageManagerBinary(%v) = %q, want %q", tt.pkgMgr, result, tt.expectedBin)
			}
		})
	}
}

func TestParseVersionOutput(t *testing.T) {
	tests := []struct {
		name     string
		output   string
		expected string
	}{
		{
			name:     "homebrew version",
			output:   "Homebrew 3.6.12\nHomebrew/homebrew-core (git revision 123abc123a; last commit 2022-10-01)",
			expected: "Homebrew",
		},
		{
			name:     "dnf version",
			output:   "dnf 4.14.0\nCached metadata: 1 week, 4 days old",
			expected: "dnf",
		},
		{
			name:     "pacman version",
			output:   " pacman 6.0.2 - libalpm v13.0.2",
			expected: "pacman",
		},
		{
			name:     "v prefix",
			output:   "v1.2.3\nSome additional info",
			expected: "1.2.3",
		},
		{
			name:     "empty output",
			output:   "",
			expected: "",
		},
		{
			name:     "single word with space",
			output:   "myapp 1.0.0",
			expected: "myapp",
		},
		{
			name:     "multi-line with spaces",
			output:   "  v1.2.3  \nline 2",
			expected: "1.2.3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseVersionOutput(tt.output)
			// Check that the result contains expected content or matches exactly
			if tt.expected != "" && !strings.Contains(result, tt.expected) && result != tt.expected {
				t.Errorf("parseVersionOutput(%q) = %q, want %q", tt.output, result, tt.expected)
			}
		})
	}
}

func TestGetLinuxPreferredOrder(t *testing.T) {
	tests := []struct {
		name       string
		distro     types.Distro
		wantFirst  types.PackageManager
	}{
		{
			name:      "ubuntu",
			distro:    types.DistroUbuntu,
			wantFirst: types.PkgMgrAPT,
		},
		{
			name:      "debian",
			distro:    types.DistroDebian,
			wantFirst: types.PkgMgrAPT,
		},
		{
			name:      "linuxmint",
			distro:    types.DistroLinuxMint,
			wantFirst: types.PkgMgrAPT,
		},
		{
			name:      "fedora",
			distro:    types.DistroFedora,
			wantFirst: types.PkgMgrDNF,
		},
		{
			name:      "rhel",
			distro:    types.DistroRHEL,
			wantFirst: types.PkgMgrDNF,
		},
		{
			name:      "centos",
			distro:    types.DistroCentOS,
			wantFirst: types.PkgMgrDNF,
		},
		{
			name:      "rocky",
			distro:    types.DistroRocky,
			wantFirst: types.PkgMgrDNF,
		},
		{
			name:      "alma",
			distro:    types.DistroAlma,
			wantFirst: types.PkgMgrDNF,
		},
		{
			name:      "arch",
			distro:    types.DistroArch,
			wantFirst: types.PkgMgrPacman,
		},
		{
			name:      "manjaro",
			distro:    types.DistroManjaro,
			wantFirst: types.PkgMgrPacman,
		},
		{
			name:      "opensuse",
			distro:    types.DistroOpenSUSE,
			wantFirst: types.PkgMgrZypper,
		},
		{
			name:      "unknown",
			distro:    types.DistroUnknown,
			wantFirst: types.PkgMgrAPT,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getLinuxPreferredOrder(tt.distro)
			if len(result) == 0 {
				t.Errorf("getLinuxPreferredOrder(%v) returned empty slice", tt.distro)
				return
			}
			if result[0] != tt.wantFirst {
				t.Errorf("getLinuxPreferredOrder(%v)[0] = %v, want %v", tt.distro, result[0], tt.wantFirst)
			}
		})
	}
}

func TestGetPreferredPackageManager(t *testing.T) {
	tests := []struct {
		name      string
		os        types.OS
		distro    types.Distro
		managers  []types.PackageManager
		expected  types.PackageManager
	}{
		{
			name:      "darwin with homebrew",
			os:        types.OSDarwin,
			distro:    types.DistroUnknown,
			managers:  []types.PackageManager{types.PkgMgrHomebrew, types.PkgMgrMacPorts},
			expected:  types.PkgMgrHomebrew,
		},
		{
			name:      "windows with winget",
			os:        types.OSWindows,
			distro:    types.DistroUnknown,
			managers:  []types.PackageManager{types.PkgMgrWinget, types.PkgMgrChocolatey},
			expected:  types.PkgMgrWinget,
		},
		{
			name:      "ubuntu with apt",
			os:        types.OSLinux,
			distro:    types.DistroUbuntu,
			managers:  []types.PackageManager{types.PkgMgrAPT, types.PkgMgrSnap},
			expected:  types.PkgMgrAPT,
		},
		{
			name:      "empty managers returns unknown",
			os:        types.OSLinux,
			distro:    types.DistroUbuntu,
			managers:  []types.PackageManager{},
			expected:  types.PkgMgrUnknown,
		},
		{
			name:      "fallback to first available",
			os:        types.OSLinux,
			distro:    types.DistroUnknown,
			managers:  []types.PackageManager{types.PkgMgrPacman},
			expected:  types.PkgMgrPacman,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetPreferredPackageManager(tt.os, tt.distro, tt.managers)
			if result != tt.expected {
				t.Errorf("GetPreferredPackageManager(%v, %v, %v) = %v, want %v", tt.os, tt.distro, tt.managers, result, tt.expected)
			}
		})
	}
}

func TestGetAllPackageManagers(t *testing.T) {
	result := getAllPackageManagers()

	expectedCount := 12 // All supported package managers

	if len(result) != expectedCount {
		t.Errorf("getAllPackageManagers() returned %d managers, want %d", len(result), expectedCount)
	}

	// Verify expected managers are present
	expectedManagers := []types.PackageManager{
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

	for _, expected := range expectedManagers {
		found := false
		for _, m := range result {
			if m == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected manager %v not found in result", expected)
		}
	}
}
