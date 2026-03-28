package config

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// ConfigWriter handles safe atomic config file writing
type ConfigWriter struct {
	tempDir string
}

// NewConfigWriter creates a new ConfigWriter
func NewConfigWriter() (*ConfigWriter, error) {
	tempDir := os.TempDir()
	return &ConfigWriter{
		tempDir: tempDir,
	}, nil
}

// WriteConfig atomically writes content to a config file
// Uses temp file + rename for atomicity
func (cw *ConfigWriter) WriteConfig(configPath string, content string) error {
	absPath, err := filepath.Abs(configPath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	// Ensure parent directory exists
	parentDir := filepath.Dir(absPath)
	if err := os.MkdirAll(parentDir, 0755); err != nil {
		return fmt.Errorf("failed to create parent directory: %w", err)
	}

	// Create temp file in the same directory (ensures same filesystem for atomic rename)
	tempFile, err := os.CreateTemp(parentDir, ".config_*.tmp")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	tempPath := tempFile.Name()

	// Write content to temp file
	_, err = io.WriteString(tempFile, content)
	if err != nil {
		tempFile.Close()
		os.Remove(tempPath)
		return fmt.Errorf("failed to write to temp file: %w", err)
	}

	// Close the temp file before renaming
	if err := tempFile.Close(); err != nil {
		os.Remove(tempPath)
		return fmt.Errorf("failed to close temp file: %w", err)
	}

	// Set permissions (use 0644 as default)
	if err := os.Chmod(tempPath, 0644); err != nil {
		os.Remove(tempPath)
		return fmt.Errorf("failed to set permissions: %w", err)
	}

	// Atomic rename
	if err := os.Rename(tempPath, absPath); err != nil {
		os.Remove(tempPath)
		return fmt.Errorf("failed to rename temp file: %w", err)
	}

	return nil
}

// WriteConfigWithBackup writes config and creates a backup first
func (cw *ConfigWriter) WriteConfigWithBackup(configPath string, content string, backup *BackupManager) error {
	// Create backup if file exists
	if _, err := os.Stat(configPath); err == nil {
		if backup != nil {
			if _, err := backup.BackupConfig(configPath); err != nil {
				return fmt.Errorf("failed to create backup: %w", err)
			}
		}
	}

	return cw.WriteConfig(configPath, content)
}

// AppendConfig appends content to an existing config file
func (cw *ConfigWriter) AppendConfig(configPath string, content string) error {
	absPath, err := filepath.Abs(configPath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	// Open file in append mode, create if doesn't exist
	f, err := os.OpenFile(absPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open config file: %w", err)
	}
	defer f.Close()

	_, err = io.WriteString(f, content)
	if err != nil {
		return fmt.Errorf("failed to append to config file: %w", err)
	}

	return nil
}

// EnsureDir ensures a directory exists, creating it if necessary
func (cw *ConfigWriter) EnsureDir(dirPath string) error {
	absPath, err := filepath.Abs(dirPath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	if err := os.MkdirAll(absPath, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	return nil
}

// ReadConfig reads content from a config file
func (cw *ConfigWriter) ReadConfig(configPath string) (string, error) {
	absPath, err := filepath.Abs(configPath)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path: %w", err)
	}

	content, err := os.ReadFile(absPath)
	if err != nil {
		return "", fmt.Errorf("failed to read config file: %w", err)
	}

	return string(content), nil
}
