package detector

import (
	"os"
	"os/user"
	"strings"

	"github.com/auto-dev-terminal/auto-dev-terminal/internal/types"
)

// Detect performs a full system detection and returns a SystemInfo struct.
func Detect() (*types.SystemInfo, error) {
	info := &types.SystemInfo{}

	// Detect OS
	info.OS = DetectOS()
	info.Arch = DetectArch()

	// Detect distro if Linux
	if info.OS == types.OSLinux {
		distro, version, err := DetectDistro()
		if err == nil {
			info.Distro = distro
			info.DistroVersion = version
		}
	}

	// Detect shell
	shell, shellVersion, err := DetectShell()
	if err == nil {
		info.Shell = shell
		info.ShellVersion = shellVersion
	}

	// Detect package managers
	info.PackageManagers = DetectPackageManagers(info.OS, info.Distro)

	// Get user info
	if user, err := user.Current(); err == nil {
		info.Username = user.Username
		info.HomeDir = user.HomeDir
	}

	// Get hostname
	if hostname, err := os.Hostname(); err == nil {
		info.Hostname = hostname
	}

	return info, nil
}

// DetectSimple performs a quick detection with minimal system calls.
func DetectSimple() types.OS {
	return DetectOS()
}

// DetectWithOptions performs detection with additional options.
func DetectWithOptions(detectShell, detectPkgMgrs bool) (*types.SystemInfo, error) {
	info := &types.SystemInfo{}

	// Detect OS
	info.OS = DetectOS()
	info.Arch = DetectArch()

	// Detect distro if Linux
	if info.OS == types.OSLinux {
		distro, version, err := DetectDistro()
		if err == nil {
			info.Distro = distro
			info.DistroVersion = version
		}
	}

	// Optionally detect shell
	if detectShell {
		shell, shellVersion, err := DetectShell()
		if err == nil {
			info.Shell = shell
			info.ShellVersion = shellVersion
		}
	}

	// Optionally detect package managers
	if detectPkgMgrs {
		info.PackageManagers = DetectPackageManagers(info.OS, info.Distro)
	}

	// Get user info
	if user, err := user.Current(); err == nil {
		info.Username = user.Username
		info.HomeDir = user.HomeDir
	}

	// Get hostname
	if hostname, err := os.Hostname(); err == nil {
		info.Hostname = hostname
	}

	return info, nil
}

// GetSystemInfo is an alias for Detect for clarity.
var GetSystemInfo = Detect

// DetectPrintable returns a human-readable summary of the detected system.
func DetectPrintable() string {
	info, err := Detect()
	if err != nil {
		return "Error detecting system: " + err.Error()
	}

	var b strings.Builder
	b.WriteString("System Detection Results:\n")
	b.WriteString("========================\n")
	b.WriteString("OS: " + string(info.OS) + "\n")

	if info.OS == types.OSLinux && info.Distro != "" {
		b.WriteString("Distro: " + string(info.Distro))
		if info.DistroVersion != "" {
			b.WriteString(" (" + info.DistroVersion + ")")
		}
		b.WriteString("\n")
	}

	b.WriteString("Arch: " + info.Arch + "\n")
	b.WriteString("Shell: " + string(info.Shell))
	if info.ShellVersion != "" {
		b.WriteString(" (" + info.ShellVersion + ")")
	}
	b.WriteString("\n")

	if len(info.PackageManagers) > 0 {
		b.WriteString("Package Managers: ")
		for i, pm := range info.PackageManagers {
			if i > 0 {
				b.WriteString(", ")
			}
			b.WriteString(string(pm))
		}
		b.WriteString("\n")
	}

	b.WriteString("User: " + info.Username + "\n")
	b.WriteString("Home: " + info.HomeDir + "\n")
	b.WriteString("Hostname: " + info.Hostname + "\n")

	return b.String()
}
