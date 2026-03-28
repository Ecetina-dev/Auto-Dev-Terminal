package detector

import (
	"testing"

	"github.com/auto-dev-terminal/auto-dev-terminal/internal/types"
)

func TestParseDistroID(t *testing.T) {
	tests := []struct {
		name     string
		id       string
		idLike   string
		expected types.Distro
	}{
		// Primary ID tests
		{
			name:     "ubuntu",
			id:       "ubuntu",
			idLike:   "",
			expected: types.DistroUbuntu,
		},
		{
			name:     "debian",
			id:       "debian",
			idLike:   "",
			expected: types.DistroDebian,
		},
		{
			name:     "fedora",
			id:       "fedora",
			idLike:   "",
			expected: types.DistroFedora,
		},
		{
			name:     "rhel",
			id:       "rhel",
			idLike:   "",
			expected: types.DistroRHEL,
		},
		{
			name:     "centos",
			id:       "centos",
			idLike:   "",
			expected: types.DistroCentOS,
		},
		{
			name:     "rocky",
			id:       "rocky",
			idLike:   "",
			expected: types.DistroRocky,
		},
		{
			name:     "alma",
			id:       "alma",
			idLike:   "",
			expected: types.DistroAlma,
		},
		{
			name:     "almalinux",
			id:       "almalinux",
			idLike:   "",
			expected: types.DistroAlma,
		},
		{
			name:     "arch",
			id:       "arch",
			idLike:   "",
			expected: types.DistroArch,
		},
		{
			name:     "manjaro",
			id:       "manjaro",
			idLike:   "",
			expected: types.DistroManjaro,
		},
		{
			name:     "endeavouros",
			id:       "endeavouros",
			idLike:   "",
			expected: types.DistroEndeavour,
		},
		{
			name:     "opensuse",
			id:       "opensuse",
			idLike:   "",
			expected: types.DistroOpenSUSE,
		},
		{
			name:     "opensuse-leap",
			id:       "opensuse-leap",
			idLike:   "",
			expected: types.DistroOpenSUSE,
		},
		{
			name:     "opensuse-tumbleweed",
			id:       "opensuse-tumbleweed",
			idLike:   "",
			expected: types.DistroOpenSUSE,
		},
		{
			name:     "sles",
			id:       "sles",
			idLike:   "",
			expected: types.DistroSLES,
		},
		{
			name:     "linuxmint",
			id:       "linuxmint",
			idLike:   "",
			expected: types.DistroLinuxMint,
		},
		{
			name:     "pop",
			id:       "pop",
			idLike:   "",
			expected: types.DistroPop,
		},
		// ID_LIKE fallback tests
		{
			name:     "unknown with ubuntu like",
			id:       "neon",
			idLike:   "ubuntu debian",
			expected: types.DistroUbuntu,
		},
		{
			name:     "unknown with debian like",
			id:       "kali",
			idLike:   "debian",
			expected: types.DistroDebian,
		},
		{
			name:     "unknown with fedora like",
			id:       "rawhide",
			idLike:   "fedora",
			expected: types.DistroFedora,
		},
		{
			name:     "unknown with rhel like",
			id:       "oraclelinux",
			idLike:   "rhel centos",
			expected: types.DistroRHEL,
		},
		{
			name:     "unknown with arch like",
			id:       "garuda",
			idLike:   "arch",
			expected: types.DistroArch,
		},
		// Unknown distro returns raw ID
		{
			name:     "unknown distro",
			id:       "unknown Distro",
			idLike:   "",
			expected: "unknown Distro",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseDistroID(tt.id, tt.idLike)
			if result != tt.expected {
				t.Errorf("parseDistroID(%q, %q) = %v, want %v", tt.id, tt.idLike, result, tt.expected)
			}
		})
	}
}

func TestDistroSupportsAPT(t *testing.T) {
	tests := []struct {
		name     string
		distro   types.Distro
		expected bool
	}{
		{
			name:     "ubuntu supports APT",
			distro:   types.DistroUbuntu,
			expected: true,
		},
		{
			name:     "debian supports APT",
			distro:   types.DistroDebian,
			expected: true,
		},
		{
			name:     "linuxmint supports APT",
			distro:   types.DistroLinuxMint,
			expected: true,
		},
		{
			name:     "pop supports APT",
			distro:   types.DistroPop,
			expected: true,
		},
		{
			name:     "fedora does not support APT",
			distro:   types.DistroFedora,
			expected: false,
		},
		{
			name:     "arch does not support APT",
			distro:   types.DistroArch,
			expected: false,
		},
		{
			name:     "unknown does not support APT",
			distro:   types.DistroUnknown,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DistroSupportsAPT(tt.distro)
			if result != tt.expected {
				t.Errorf("DistroSupportsAPT(%v) = %v, want %v", tt.distro, result, tt.expected)
			}
		})
	}
}

func TestDistroSupportsDNF(t *testing.T) {
	tests := []struct {
		name     string
		distro   types.Distro
		expected bool
	}{
		{
			name:     "fedora supports DNF",
			distro:   types.DistroFedora,
			expected: true,
		},
		{
			name:     "rhel supports DNF",
			distro:   types.DistroRHEL,
			expected: true,
		},
		{
			name:     "centos supports DNF",
			distro:   types.DistroCentOS,
			expected: true,
		},
		{
			name:     "rocky supports DNF",
			distro:   types.DistroRocky,
			expected: true,
		},
		{
			name:     "alma supports DNF",
			distro:   types.DistroAlma,
			expected: true,
		},
		{
			name:     "ubuntu does not support DNF",
			distro:   types.DistroUbuntu,
			expected: false,
		},
		{
			name:     "arch does not support DNF",
			distro:   types.DistroArch,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DistroSupportsDNF(tt.distro)
			if result != tt.expected {
				t.Errorf("DistroSupportsDNF(%v) = %v, want %v", tt.distro, result, tt.expected)
			}
		})
	}
}

func TestDistroSupportsPacman(t *testing.T) {
	tests := []struct {
		name     string
		distro   types.Distro
		expected bool
	}{
		{
			name:     "arch supports Pacman",
			distro:   types.DistroArch,
			expected: true,
		},
		{
			name:     "manjaro supports Pacman",
			distro:   types.DistroManjaro,
			expected: true,
		},
		{
			name:     "endeavour supports Pacman",
			distro:   types.DistroEndeavour,
			expected: true,
		},
		{
			name:     "ubuntu does not support Pacman",
			distro:   types.DistroUbuntu,
			expected: false,
		},
		{
			name:     "fedora does not support Pacman",
			distro:   types.DistroFedora,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DistroSupportsPacman(tt.distro)
			if result != tt.expected {
				t.Errorf("DistroSupportsPacman(%v) = %v, want %v", tt.distro, result, tt.expected)
			}
		})
	}
}

func TestDistroSupportsZypper(t *testing.T) {
	tests := []struct {
		name     string
		distro   types.Distro
		expected bool
	}{
		{
			name:     "opensuse supports Zypper",
			distro:   types.DistroOpenSUSE,
			expected: true,
		},
		{
			name:     "sles supports Zypper",
			distro:   types.DistroSLES,
			expected: true,
		},
		{
			name:     "ubuntu does not support Zypper",
			distro:   types.DistroUbuntu,
			expected: false,
		},
		{
			name:     "fedora does not support Zypper",
			distro:   types.DistroFedora,
			expected: false,
		},
		{
			name:     "arch does not support Zypper",
			distro:   types.DistroArch,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DistroSupportsZypper(tt.distro)
			if result != tt.expected {
				t.Errorf("DistroSupportsZypper(%v) = %v, want %v", tt.distro, result, tt.expected)
			}
		})
	}
}
