package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNewBackupManager(t *testing.T) {
	bm, err := NewBackupManager()
	if err != nil {
		t.Fatalf("NewBackupManager() error = %v", err)
	}

	if bm == nil {
		t.Fatal("NewBackupManager() returned nil")
	}

	if bm.backupDir == "" {
		t.Error("BackupDir should not be empty")
	}

	// Verify backup directory was created
	if _, err := os.Stat(bm.backupDir); os.IsNotExist(err) {
		t.Errorf("Backup directory was not created: %v", err)
	}
}

func TestBackupManagerBackupDir(t *testing.T) {
	bm, err := NewBackupManager()
	if err != nil {
		t.Fatalf("NewBackupManager() error = %v", err)
	}

	dir := bm.BackupDir()
	if dir == "" {
		t.Error("BackupDir() returned empty string")
	}

	// Verify it returns the correct directory
	if dir != bm.backupDir {
		t.Errorf("BackupDir() = %q, want %q", dir, bm.backupDir)
	}
}

func TestBackupManagerBackupConfig(t *testing.T) {
	bm, err := NewBackupManager()
	if err != nil {
		t.Fatalf("NewBackupManager() error = %v", err)
	}

	// Create a temporary config file to backup
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "testconfig")
	testContent := "# Test Configuration\ntest.key=value\n"
	
	if err := os.WriteFile(configPath, []byte(testContent), 0644); err != nil {
		t.Fatalf("Failed to create test config: %v", err)
	}

	// Test successful backup
	backupPath, err := bm.BackupConfig(configPath)
	if err != nil {
		t.Fatalf("BackupConfig() error = %v", err)
	}

	if backupPath == "" {
		t.Error("BackupConfig() returned empty path")
	}

	// Verify backup file exists
	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		t.Errorf("Backup file was not created: %v", err)
	}

	// Verify backup content matches original
	backupContent, err := os.ReadFile(backupPath)
	if err != nil {
		t.Fatalf("Failed to read backup: %v", err)
	}

	if string(backupContent) != testContent {
		t.Errorf("Backup content = %q, want %q", string(backupContent), testContent)
	}

	// Test with non-existent file
	_, err = bm.BackupConfig("/nonexistent/config")
	if err == nil {
		t.Error("BackupConfig() should fail for non-existent file")
	}
}

func TestBackupManagerListBackups(t *testing.T) {
	bm, err := NewBackupManager()
	if err != nil {
		t.Fatalf("NewBackupManager() error = %v", err)
	}

	// Create a temporary config file to backup
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "testconfig")
	
	if err := os.WriteFile(configPath, []byte("test content"), 0644); err != nil {
		t.Fatalf("Failed to create test config: %v", err)
	}

	// Create multiple backups
	for i := 0; i < 3; i++ {
		// Add small delay to ensure different timestamps
		time.Sleep(10 * time.Millisecond)
		if _, err := bm.BackupConfig(configPath); err != nil {
			t.Fatalf("Failed to create backup %d: %v", i, err)
		}
	}

	// List backups for "testconfig"
	backups, err := bm.ListBackups("testconfig")
	if err != nil {
		t.Fatalf("ListBackups() error = %v", err)
	}

	if len(backups) < 3 {
		t.Errorf("ListBackups() returned %d backups, want at least 3", len(backups))
	}

	// Test with non-existent config name
	backups, err = bm.ListBackups("nonexistent")
	if err != nil {
		t.Fatalf("ListBackups() error = %v", err)
	}

	if len(backups) != 0 {
		t.Errorf("ListBackups() returned %d backups for non-existent config, want 0", len(backups))
	}
}

func TestBackupManagerGetLatestBackup(t *testing.T) {
	bm, err := NewBackupManager()
	if err != nil {
		t.Fatalf("NewBackupManager() error = %v", err)
	}

	// Create a temporary config file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "latesttest")
	
	if err := os.WriteFile(configPath, []byte("test content"), 0644); err != nil {
		t.Fatalf("Failed to create test config: %v", err)
	}

	// Create multiple backups with delay
	var backupTimes []time.Time
	for i := 0; i < 3; i++ {
		time.Sleep(10 * time.Millisecond)
		backupPath, err := bm.BackupConfig(configPath)
		if err != nil {
			t.Fatalf("Failed to create backup %d: %v", i, err)
		}
		info, _ := os.Stat(backupPath)
		backupTimes = append(backupTimes, info.ModTime())
	}

	// Get latest backup
	latest, err := bm.GetLatestBackup("latesttest")
	if err != nil {
		t.Fatalf("GetLatestBackup() error = %v", err)
	}

	if latest == "" {
		t.Error("GetLatestBackup() returned empty path")
	}

	// Test with no backups
	_, err = bm.GetLatestBackup("nonexistentbackup")
	if err == nil {
		t.Error("GetLatestBackup() should fail when no backups exist")
	}
}

func TestBackupConfigTimestampFormat(t *testing.T) {
	bm, err := NewBackupManager()
	if err != nil {
		t.Fatalf("NewBackupManager() error = %v", err)
	}

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "timestamp_test")
	
	if err := os.WriteFile(configPath, []byte("content"), 0644); err != nil {
		t.Fatalf("Failed to create test config: %v", err)
	}

	backupPath, err := bm.BackupConfig(configPath)
	if err != nil {
		t.Fatalf("BackupConfig() error = %v", err)
	}

	// Extract timestamp from backup filename
	filename := filepath.Base(backupPath)
	// Format is: configName_YYYY-MM-DD_HH-MM-SS.bak
	// We expect the timestamp part to be parseable
	parts := filename[len("timestamp_test_"):]
	parts = parts[:len(parts)-len(".bak")] // Remove .bak extension
	
	_, err = time.Parse("2006-01-02_15-04-05", parts)
	if err != nil {
		t.Errorf("Backup filename timestamp is not in expected format: %v", err)
	}
}

func TestBackupConfigAbsolutePath(t *testing.T) {
	bm, err := NewBackupManager()
	if err != nil {
		t.Fatalf("NewBackupManager() error = %v", err)
	}

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "abspath_test")
	
	if err := os.WriteFile(configPath, []byte("content"), 0644); err != nil {
		t.Fatalf("Failed to create test config: %v", err)
	}

	// Test with relative path (should convert to absolute)
	backupPath, err := bm.BackupConfig(configPath)
	if err != nil {
		t.Fatalf("BackupConfig() error = %v", err)
	}

	// Verify backup is in backup directory
	if !filepath.IsAbs(backupPath) {
		t.Error("Backup path should be absolute")
	}

	if !filepath.HasPrefix(backupPath, bm.backupDir) {
		t.Errorf("Backup should be in backup dir %q, got %q", bm.backupDir, backupPath)
	}
}
