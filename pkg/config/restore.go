package config

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// RestoreManager handles restoring config from backups
type RestoreManager struct {
	backupDir string
}

// BackupDir returns the backup directory path
func (rm *RestoreManager) BackupDir() string {
	return rm.backupDir
}

// NewRestoreManager creates a new RestoreManager
func NewRestoreManager() (*RestoreManager, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	backupDir := filepath.Join(homeDir, ".auto-dev-terminal", "backups")

	return &RestoreManager{
		backupDir: backupDir,
	}, nil
}

// RestoreFromBackup restores a config file from a specific backup
func (rm *RestoreManager) RestoreFromBackup(backupPath string, targetPath string) error {
	// Verify backup exists
	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		return fmt.Errorf("backup file does not exist: %s", backupPath)
	}

	// Read backup content
	content, err := os.ReadFile(backupPath)
	if err != nil {
		return fmt.Errorf("failed to read backup file: %w", err)
	}

	// Get absolute target path
	absTarget, err := filepath.Abs(targetPath)
	if err != nil {
		return fmt.Errorf("failed to get absolute target path: %w", err)
	}

	// Ensure target directory exists
	targetDir := filepath.Dir(absTarget)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return fmt.Errorf("failed to create target directory: %w", err)
	}

	// Write to target (atomic write)
	writer, err := NewConfigWriter()
	if err != nil {
		return fmt.Errorf("failed to create config writer: %w", err)
	}

	if err := writer.WriteConfig(absTarget, string(content)); err != nil {
		return fmt.Errorf("failed to write restored config: %w", err)
	}

	return nil
}

// RestoreFromLatestBackup restores from the most recent backup for a config
func (rm *RestoreManager) RestoreFromLatestBackup(configName string, targetPath string) error {
	backup, err := rm.FindLatestBackup(configName)
	if err != nil {
		return err
	}

	return rm.RestoreFromBackup(backup, targetPath)
}

// FindLatestBackup finds the most recent backup for a given config name
func (rm *RestoreManager) FindLatestBackup(configName string) (string, error) {
	backups, err := rm.ListAvailableBackups(configName)
	if err != nil {
		return "", err
	}

	if len(backups) == 0 {
		return "", fmt.Errorf("no backups found for %s", configName)
	}

	// backups is sorted oldest first, so take the last one
	return backups[len(backups)-1].path, nil
}

// BackupEntry represents a single backup file
type BackupEntry struct {
	path       string
	filename   string
	configName string
	timestamp  time.Time
}

// Path returns the backup file path
func (be BackupEntry) Path() string {
	return be.path
}

// Filename returns the backup file name
func (be BackupEntry) Filename() string {
	return be.filename
}

// ConfigName returns the config name this backup belongs to
func (be BackupEntry) ConfigName() string {
	return be.configName
}

// Timestamp returns the backup timestamp
func (be BackupEntry) Timestamp() time.Time {
	return be.timestamp
}

// ListAvailableBackups lists all available backups for a config name
func (rm *RestoreManager) ListAvailableBackups(configName string) ([]BackupEntry, error) {
	entries, err := os.ReadDir(rm.backupDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read backup directory: %w", err)
	}

	var backups []BackupEntry
	prefix := configName + "_"

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		filename := entry.Name()
		if !strings.HasPrefix(filename, prefix) || !strings.HasSuffix(filename, ".bak") {
			continue
		}

		// Extract timestamp from filename: configName_YYYY-MM-DD_HH-MM-SS.bak
		timestampStr := strings.TrimSuffix(filename[len(prefix):], ".bak")
		timestamp, err := time.Parse("2006-01-02_15-04-05", timestampStr)
		if err != nil {
			// Skip files with invalid timestamp format
			continue
		}

		backups = append(backups, BackupEntry{
			path:       filepath.Join(rm.backupDir, filename),
			filename:   filename,
			configName: configName,
			timestamp:  timestamp,
		})
	}

	// Sort by timestamp (oldest first)
	sort.Slice(backups, func(i, j int) bool {
		return backups[i].timestamp.Before(backups[j].timestamp)
	})

	return backups, nil
}

// ListAllBackups lists all backups in the backup directory
func (rm *RestoreManager) ListAllBackups() (map[string][]BackupEntry, error) {
	entries, err := os.ReadDir(rm.backupDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read backup directory: %w", err)
	}

	backupsByConfig := make(map[string][]BackupEntry)

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		filename := entry.Name()
		if !strings.HasSuffix(filename, ".bak") {
			continue
		}

		// Parse config name and timestamp from filename
		// Format: configName_YYYY-MM-DD_HH-MM-SS.bak
		// Note: Find the FIRST underscore that starts the timestamp (date_time)
		base := strings.TrimSuffix(filename, ".bak")

		// Find the underscore that separates configName from timestamp
		// The timestamp format is YYYY-MM-DD_HH-MM-SS (date_HH-MM-SS)
		// We need to find where the date ends (first underscore after the date)
		dateTimeIdx := strings.Index(base, "_20") // Timestamp always starts with year 20xx
		if dateTimeIdx == -1 {
			continue
		}

		configName := base[:dateTimeIdx]
		timestampStr := base[dateTimeIdx+1:]

		timestamp, err := time.Parse("2006-01-02_15-04-05", timestampStr)
		if err != nil {
			continue
		}

		backup := BackupEntry{
			path:       filepath.Join(rm.backupDir, filename),
			filename:   filename,
			configName: configName,
			timestamp:  timestamp,
		}

		backupsByConfig[configName] = append(backupsByConfig[configName], backup)
	}

	// Sort each config's backups
	for configName := range backupsByConfig {
		sort.Slice(backupsByConfig[configName], func(i, j int) bool {
			return backupsByConfig[configName][i].timestamp.Before(backupsByConfig[configName][j].timestamp)
		})
	}

	return backupsByConfig, nil
}

// DeleteOldBackups removes backups older than the specified number of days
func (rm *RestoreManager) DeleteOldBackups(configName string, days int) (int, error) {
	backups, err := rm.ListAvailableBackups(configName)
	if err != nil {
		return 0, err
	}

	cutoff := time.Now().AddDate(0, 0, -days)
	deleted := 0

	for _, backup := range backups {
		if backup.timestamp.Before(cutoff) {
			if err := os.Remove(backup.path); err != nil {
				return deleted, fmt.Errorf("failed to delete backup %s: %w", backup.path, err)
			}
			deleted++
		}
	}

	return deleted, nil
}
