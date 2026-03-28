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
// Priority: shell-specific env vars > $SHELL (with Windows-specific logic)
func DetectShell() (types.Shell, string, error) {
	// FIRST: On Windows, check if we're actually in PowerShell
	// This must come BEFORE $SHELL because Git Bash sets $SHELL globally
	// and it persists even when running in PowerShell
	if runtime.GOOS == "windows" {
		// Check PSHOME - this is set when in PowerShell (most reliable)
		if pshome := getEnvOrEmpty("PSHOME"); pshome != "" {
			if strings.Contains(strings.ToLower(pshome), "pwsh") {
				return types.ShellPwsh, getPowerShellVersion("pwsh"), nil
			}
			return types.ShellPowerShell, getPowerShellVersion("powershell"), nil
		}

		// Check if PATH contains PowerShell directories - this indicates we're in PowerShell
		// PowerShell adds its directories to PATH when active
		path := getEnvOrEmpty("PATH")
		if strings.Contains(path, "PowerShell\\v1.0") ||
			strings.Contains(path, "PowerShell/7") ||
			strings.Contains(path, "WindowsPowerShell") {
			// Check if it's pwsh (PowerShell 7+) or powershell (Windows PowerShell)
			if strings.Contains(path, "PowerShell/7") || strings.Contains(path, "PowerShell 7") {
				return types.ShellPwsh, "", nil
			}
			return types.ShellPowerShell, "", nil
		}

		// Check parent process as fallback for PowerShell
		parentShell := detectParentShell()
		if parentShell == types.ShellPowerShell || parentShell == types.ShellPwsh {
			return parentShell, "", nil
		}

		// Check if we're in cmd.exe (no PSModulePath, no MSYSTEM, no bash)
		msystem := getEnvOrEmpty("MSYSTEM")
		shellPath := getEnvOrEmpty("SHELL")
		if msystem == "" && shellPath == "" && getEnvOrEmpty("PSHOME") == "" {
			return types.ShellCmd, "", nil
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

	// THIRD: Check $SHELL - but only if we're NOT in a Windows terminal
	// that was detected above (PowerShell, cmd)
	// On Windows, $SHELL is often set by Git Bash even when in PowerShell
	shellPath := getEnvOrEmpty("SHELL")
	if shellPath != "" && runtime.GOOS != "windows" {
		// On non-Windows, $SHELL is reliable
		shell := shellFromPath(shellPath)
		if shell != types.ShellUnknown {
			version, _ := getShellVersion(shellPath)
			return shell, version, nil
		}
	}

	// On Windows, check $SHELL only if we're confident we're in Git Bash/MSYS2
	// (not in PowerShell or cmd)
	if runtime.GOOS == "windows" && shellPath != "" {
		// Only trust $SHELL on Windows if MSYSTEM is set AND we're not in PowerShell
		msystem := getEnvOrEmpty("MSYSTEM")
		if msystem != "" && getEnvOrEmpty("PSHOME") == "" {
			shell := shellFromPath(shellPath)
			if shell != types.ShellUnknown {
				version, _ := getShellVersion(shellPath)
				return shell, version, nil
			}
		}
	}

	// FOURTH: Last resort - check for common shells in PATH
	// This handles edge cases where none of the above methods work

	return types.ShellUnknown, "", fmt.Errorf("could not detect shell")
}

// getPowerShellVersion attempts to get the PowerShell version.
func getPowerShellVersion(powershellType string) string {
	var cmd *exec.Cmd
	if powershellType == "pwsh" {
		cmd = exec.Command("pwsh", "-NoProfile", "-Command", "$PSVersionTable.PSVersion.ToString()")
	} else {
		cmd = exec.Command("powershell", "-NoProfile", "-Command", "$PSVersionTable.PSVersion.ToString()")
	}
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(output))
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

	// Check for PowerShell - PSHOME is set when in PowerShell
	// This is more reliable than PSModulePath which can be set system-wide
	if pshome := getEnvOrEmpty("PSHOME"); pshome != "" {
		// Check if it's PowerShell 7+ (pwsh) or Windows PowerShell
		if strings.Contains(strings.ToLower(pshome), "pwsh") {
			return types.ShellPwsh
		}
		return types.ShellPowerShell
	}

	// Check for other MSYS2/Unix-like environment indicators
	if getEnvOrEmpty("TERM") == "xterm" {
		// Unix-like terminal on Windows (Git Bash, WSL, etc.)
		if commandExists("/usr/bin/bash") {
			return types.ShellBash
		}
	}

	// Additional check: look for powershell.exe in the process path
	// This is a fallback for edge cases
	parentProc := getParentProcessName()
	if parentProc != "" {
		lower := strings.ToLower(parentProc)
		if strings.Contains(lower, "pwsh") {
			return types.ShellPwsh
		}
		if strings.Contains(lower, "powershell") {
			return types.ShellPowerShell
		}
	}

	return types.ShellUnknown
}

// getParentProcessName returns the name of the parent process on Windows.
func getParentProcessName() string {
	if runtime.GOOS != "windows" {
		return ""
	}

	// Use wmic to get parent process name
	cmd := exec.Command("wmic", "process", "where", "ProcessId="+getCurrentPID(), "get", "ParentProcessId")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}

	lines := strings.Split(string(output), "\n")
	if len(lines) < 2 {
		return ""
	}

	parentPID := strings.TrimSpace(lines[1])
	if parentPID == "" || parentPID == "0" {
		return ""
	}

	// Get process name by PID
	cmd = exec.Command("wmic", "process", "where", "ProcessId="+parentPID, "get", "Name")
	output, err = cmd.Output()
	if err != nil {
		return ""
	}

	lines = strings.Split(string(output), "\n")
	if len(lines) < 2 {
		return ""
	}

	return strings.TrimSpace(lines[1])
}

// getCurrentPID returns the current process ID.
func getCurrentPID() string {
	return fmt.Sprintf("%d", os.Getpid())
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
