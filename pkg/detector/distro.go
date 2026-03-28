package detector

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/auto-dev-terminal/auto-dev-terminal/internal/types"
)

const osReleasePath = "/etc/os-release"

// DetectDistro parses /etc/os-release to determine the Linux distribution.
func DetectDistro() (types.Distro, string, error) {
	file, err := os.Open(osReleasePath)
	if err != nil {
		return types.DistroUnknown, "", fmt.Errorf("failed to open /etc/os-release: %w", err)
	}
	defer file.Close()

	var (
		id        string
		idLike    string
		versionID string
	)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "ID=") {
			id = strings.Trim(strings.TrimPrefix(line, "ID="), `"`)
		}
		if strings.HasPrefix(line, "ID_LIKE=") {
			idLike = strings.Trim(strings.TrimPrefix(line, "ID_LIKE="), `"`)
		}
		if strings.HasPrefix(line, "VERSION_ID=") {
			versionID = strings.Trim(strings.TrimPrefix(line, "VERSION_ID="), `"`)
		}
	}

	if err := scanner.Err(); err != nil {
		return types.DistroUnknown, "", fmt.Errorf("error reading /etc/os-release: %w", err)
	}

	if id == "" {
		return types.DistroUnknown, "", errors.New("could not determine distro: ID not found in /etc/os-release")
	}

	distro := parseDistroID(id, idLike)
	return distro, versionID, nil
}

// parseDistroID converts the os-release ID to our Distro type.
func parseDistroID(id, idLike string) types.Distro {
	// Primary ID check
	switch id {
	case "ubuntu":
		return types.DistroUbuntu
	case "debian":
		return types.DistroDebian
	case "fedora":
		return types.DistroFedora
	case "rhel":
		return types.DistroRHEL
	case "centos":
		return types.DistroCentOS
	case "rocky":
		return types.DistroRocky
	case "alma", "almalinux":
		return types.DistroAlma
	case "arch":
		return types.DistroArch
	case "manjaro":
		return types.DistroManjaro
	case "endeavouros":
		return types.DistroEndeavour
	case "opensuse", "opensuse-leap", "opensuse-tumbleweed":
		return types.DistroOpenSUSE
	case "sles", "sles_latest":
		return types.DistroSLES
	case "linuxmint":
		return types.DistroLinuxMint
	case "pop":
		return types.DistroPop
	}

	// Fallback: check ID_LIKE for derivatives
	if idLike != "" {
		likeParts := strings.Fields(idLike)
		for _, part := range likeParts {
			switch part {
			case "ubuntu":
				return types.DistroUbuntu
			case "debian":
				return types.DistroDebian
			case "fedora":
				return types.DistroFedora
			case "rhel", "centos":
				return types.DistroRHEL
			case "arch":
				return types.DistroArch
			case "opensuse":
				return types.DistroOpenSUSE
			}
		}
	}

	// Return the raw ID if we can't map it
	return types.Distro(id)
}

// DistroSupportsAPT returns true if the distribution uses APT package manager.
func DistroSupportsAPT(d types.Distro) bool {
	switch d {
	case types.DistroUbuntu, types.DistroDebian, types.DistroLinuxMint, types.DistroPop:
		return true
	default:
		return false
	}
}

// DistroSupportsDNF returns true if the distribution uses DNF/YUM package manager.
func DistroSupportsDNF(d types.Distro) bool {
	switch d {
	case types.DistroFedora, types.DistroRHEL, types.DistroCentOS, types.DistroRocky, types.DistroAlma:
		return true
	default:
		return false
	}
}

// DistroSupportsPacman returns true if the distribution uses Pacman package manager.
func DistroSupportsPacman(d types.Distro) bool {
	switch d {
	case types.DistroArch, types.DistroManjaro, types.DistroEndeavour:
		return true
	default:
		return false
	}
}

// DistroSupportsZypper returns true if the distribution uses Zypper package manager.
func DistroSupportsZypper(d types.Distro) bool {
	switch d {
	case types.DistroOpenSUSE, types.DistroSLES:
		return true
	default:
		return false
	}
}
