package detector

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/auto-dev-terminal/auto-dev-terminal/internal/types"
)

// DetectShell detects the current shell by checking environment variables.
// Priority: $SHELL > shell-specific env vars > PATH lookup
func DetectShell() (types.Shell, string, error) {
	// FIRST: Check $SHELL - this is the most reliable indicator of current shell
	// (especially important for Git Bash, WSL, etc.)
	shellPath := getEnvOrEmpty("SHELL")
	if shellPath != "" {
		shell := shellFromPath(shellPath)
		if shell != types.ShellUnknown {
			version, _ := getShellVersion(shellPath)
			return shell, version, nil
		}
	}

	// SECOND: Check shell-specific environment variables
	// (these only exist when the shell is actually running)

	// Check for Zsh
	if zshVersion := getEnvOrEmpty("ZSH_VERSION"); zshVersion != "" {
		return types.ShellZsh, zshVersion, nil
	}

	// Check for Fish
	if fishVersion := getEnvOrEmpty("FISH_VERSION"); fishVersion != "" {
		return types.ShellFish, fishVersion, nil
	}

	// Check for Bash (MSYS2/Git Bash on Windows sets this)
	if bashVersion := getEnvOrEmpty("BASH_VERSION"); bashVersion != "" {
		return types.ShellBash, bashVersion, nil
	}

	// Check for Tcsh/Csh
	if tcshVersion := getEnvOrEmpty("TCSH_VERSION"); tcshVersion != "" {
		return types.ShellTcsh, tcshVersion, nil
	}

	// Check for Ash/Dash (usually on minimal Linux)
	if ashVersion := getEnvOrEmpty("ASH_VERSION"); ashVersion != "" {
		return types.ShellAsh, ashVersion, nil
	}

	// THIRD: Only if no env vars found, check if shell binaries exist
	// and are currently running (not just installed)

	// Check for PowerShell - but ONLY as last resort
	// Note: PSModulePath is set on Windows even when not using PowerShell
	// so we only check this if $SHELL and other env vars failed
	if runtime.GOOS == "windows" {
		// On Windows, check parent process to determine actual shell
		// or check if we're actually in a PowerShell session
		parentShell := detectParentShell()
		if parentShell != types.ShellUnknown {
			return parentShell, "", nil
		}
	}

	// Check if we're in cmd.exe on Windows (no special env vars)
	if isWindowsCmd() {
		return types.ShellCmd, "", nil
	}

	return types.ShellUnknown, "", fmt.Errorf("could not detect shell")
}

// shellFromPath converts a shell path to our Shell type.
func shellFromPath(path string) types.Shell {
	base := strings.ToLower(strings.TrimSpace(path))
	switch {
	case strings.Contains(base, "zsh"):
		return types.ShellZsh
	case strings.Contains(base, "fish"):
		return types.ShellFish
	case strings.Contains(base, "bash"):
		return types.ShellBash
	case strings.Contains(base, "pwsh"), strings.Contains(base, "powershell"):
		return types.ShellPwsh
	case strings.Contains(base, "cmd"):
		return types.ShellCmd
	case strings.Contains(base, "tcsh"):
		return types.ShellTcsh
	case strings.Contains(base, "csh"):
		return types.ShellCsh
	case strings.Contains(base, "ash"), strings.Contains(base, "dash"):
		return types.ShellAsh
	default:
		return types.ShellUnknown
	}
}

// getShellVersion attempts to get the version of a shell binary.
func getShellVersion(shellPath string) (string, error) {
	// Try --version first (most common)
	cmd := exec.Command(shellPath, "--version")
	output, err := cmd.Output()
	if err == nil {
		version := strings.TrimSpace(string(output))
		// Take first line and limit length
		lines := strings.Split(version, "\n")
		if len(lines) > 0 {
			firstPart, _, _ := strings.Cut(lines[0], " ")
			return firstPart, nil
		}
	}

	// Try -c "echo $BASH_VERSION" style
	for _, versionFlag := range []string{"-c", "--version"} {
		cmd := exec.Command(shellPath, versionFlag, "echo $VERSION")
		output, err := cmd.Output()
		if err == nil {
			return strings.TrimSpace(string(output)), nil
		}
	}

	return "", fmt.Errorf("could not get version for %s", shellPath)
}

// getEnvOrEmpty returns the value of an environment variable or empty string if not set.
func getEnvOrEmpty(key string) string {
	return os.Getenv(key)
}

// isWindowsCmd checks if running in Windows Command Prompt.
func isWindowsCmd() bool {
	if runtime.GOOS == "windows" {
		// Check if we're not in PowerShell
		if getEnvOrEmpty("PSModulePath") == "" {
			return true
		}
	}
	return false
}

// detectParentShell attempts to detect the parent shell on Windows.
// It checks the parent process name to determine if we're in bash (Git Bash/MSYS2).
func detectParentShell() types.Shell {
	if runtime.GOOS != "windows" {
		return types.ShellUnknown
	}

	// Check if we're running in a Git Bash or MSYS2 environment
	// by checking for MSYSTEM environment variable (set by MSYS2/Git Bash)
	msystem := getEnvOrEmpty("MSYSTEM")
	if msystem != "" {
		// MSYSTEM is set by Git Bash, MSYS2, etc.
		switch msystem {
		case "MINGW64", "MINGW32", "MSYS":
			return types.ShellBash
		}
	}

	// Check for other MSYS2/Unix-like environment indicators
	if getEnvOrEmpty("TERM") == "xterm" {
		// Unix-like terminal on Windows (Git Bash, WSL, etc.)
		if commandExists("/usr/bin/bash") {
			return types.ShellBash
		}
	}

	return types.ShellUnknown
}

// IsShellCompatibleWithZsh returns true if the shell can use Oh-My-Zsh.
func IsShellCompatibleWithZsh(s types.Shell) bool {
	return s == types.ShellZsh || s == types.ShellBash
}

// IsShellCompatibleWithStarship returns true if the shell supports Starship.
func IsShellCompatibleWithStarship(s types.Shell) bool {
	switch s {
	case types.ShellBash, types.ShellZsh, types.ShellFish, types.ShellPwsh, types.ShellPowerShell:
		return true
	default:
		return false
	}
}
