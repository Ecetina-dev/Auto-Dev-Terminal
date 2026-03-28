package detector

import (
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
)

// FileExists checks if a file exists and is accessible.
func FileExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

// DirExists checks if a directory exists and is accessible.
func DirExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// PathExists checks if a path (file or directory) exists.
func PathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// commandExists checks if a command is available in PATH.
func commandExists(cmd string) bool {
	// Check if the command contains a path separator
	if filepath.Base(cmd) != cmd {
		// Has path separator, check directly
		return FileExists(cmd)
	}

	// Look in PATH
	_, err := exec.LookPath(cmd)
	return err == nil
}

// getHomeDir returns the user's home directory.
func getHomeDir() string {
	user, err := user.Current()
	if err != nil {
		// Fallback to environment variable
		home := os.Getenv("HOME")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return user.HomeDir
}

// expandPath expands ~ to home directory and resolves environment variables.
func expandPath(path string) string {
	if strings.HasPrefix(path, "~") {
		home := getHomeDir()
		if home != "" {
			path = filepath.Join(home, strings.TrimPrefix(path, "~"))
		}
	}

	// Expand environment variables
	path = os.ExpandEnv(path)

	// Resolve to absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		return path
	}
	return absPath
}

// getEnvWithDefault returns the environment variable or a default value.
func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvOrDefault is an alias for getEnvWithDefault.
var getEnvOrDefault = getEnvWithDefault

// isExecutable checks if a file is executable.
func isExecutable(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.Mode()&0111 != 0
}

// readFileFirstLine reads the first line of a file.
func readFileFirstLine(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return strings.Split(string(content), "\n")[0], nil
}

// getShellFromEnv returns the shell from SHELL environment variable.
func getShellFromEnv() string {
	return os.Getenv("SHELL")
}

// getPWD returns the current working directory.
func getPWD() string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}
	return dir
}

// joinPath joins path elements with the appropriate separator.
func joinPath(elem ...string) string {
	return filepath.Join(elem...)
}

// normalizePath normalizes a path (resolve dots, handle separators).
func normalizePath(path string) string {
	return filepath.Clean(path)
}
