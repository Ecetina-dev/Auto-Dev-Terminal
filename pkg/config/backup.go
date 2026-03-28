package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// BackupManager handles creating backups before config modifications
type BackupManager struct {
	backupDir string
}

// NewBackupManager creates a new BackupManager
func NewBackupManager() (*BackupManager, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	backupDir := filepath.Join(homeDir, ".auto-dev-terminal", "backups")
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create backup directory: %w", err)
	}

	return &BackupManager{
		backupDir: backupDir,
	}, nil
}

// BackupConfig creates a timestamped backup of a config file
// Returns the path to the backup file
func (bm *BackupManager) BackupConfig(configPath string) (string, error) {
	absPath, err := filepath.Abs(configPath)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path: %w", err)
	}

	// Check if config file exists
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return "", fmt.Errorf("config file does not exist: %s", absPath)
	}

	// Generate timestamped backup filename
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	filename := filepath.Base(absPath)
	backupFilename := fmt.Sprintf("%s_%s.bak", filename, timestamp)
	backupPath := filepath.Join(bm.backupDir, backupFilename)

	// Read original config
	content, err := os.ReadFile(absPath)
	if err != nil {
		return "", fmt.Errorf("failed to read config file: %w", err)
	}

	// Write to backup
	if err := os.WriteFile(backupPath, content, 0644); err != nil {
		return "", fmt.Errorf("failed to write backup file: %w", err)
	}

	return backupPath, nil
}

// BackupDir returns the backup directory path
func (bm *BackupManager) BackupDir() string {
	return bm.backupDir
}

// ListBackups returns all backup files for a given config name
func (bm *BackupManager) ListBackups(configName string) ([]string, error) {
	entries, err := os.ReadDir(bm.backupDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read backup directory: %w", err)
	}

	var backups []string
	prefix := configName + "_"
	for _, entry := range entries {
		if !entry.IsDir() && len(entry.Name()) > len(prefix) && entry.Name()[:len(prefix)] == prefix {
			backups = append(backups, filepath.Join(bm.backupDir, entry.Name()))
		}
	}

	return backups, nil
}

// GetLatestBackup returns the most recent backup for a config file
func (bm *BackupManager) GetLatestBackup(configName string) (string, error) {
	backups, err := bm.ListBackups(configName)
	if err != nil {
		return "", err
	}

	if len(backups) == 0 {
		return "", fmt.Errorf("no backups found for %s", configName)
	}

	// Return the most recent (last in sorted list)
	return backups[len(backups)-1], nil
}
