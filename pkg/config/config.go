// Package config provides configuration management for auto-dev-terminal
//
// This package handles:
//
//   - Backup creation before config modifications
//   - Template variable substitution for config files
//   - Atomic file writing for safe config updates
//   - Restore functionality from backups
//
// # Usage
//
// Backup and modify config:
//
//	backupMgr, _ := config.NewBackupManager()
//	backupPath, _ := backupMgr.BackupConfig("/path/to/config")
//
//	engine, _ := config.NewTemplateEngine()
//	processed, _ := engine.ProcessTemplate(configContent)
//
//	writer, _ := config.NewConfigWriter()
//	writer.WriteConfig("/path/to/config", processed)
//
// Restore from backup:
//
//	restoreMgr, _ := config.NewRestoreManager()
//	restoreMgr.RestoreFromLatestBackup("myconfig", "/path/to/config")
package config

// Package version
const Version = "0.1.0"

// DefaultConfigDir is the default config directory
const DefaultConfigDir = ".auto-dev-terminal"

// DefaultBackupDir is the default backup subdirectory
const DefaultBackupDir = "backups"
