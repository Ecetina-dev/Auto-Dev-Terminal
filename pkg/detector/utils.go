package detector

import (
	"os"
	"os/exec"
	"path/filepath"
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
