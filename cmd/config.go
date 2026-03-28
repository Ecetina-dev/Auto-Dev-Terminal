// Package cmd provides the CLI command structure for Auto-Dev-Terminal.
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/auto-dev-terminal/auto-dev-terminal/pkg/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// ConfigCmd represents the config command
var ConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration backups and restores",
	Long: `Manage configuration file backups and restores.

Subcommands:
  backup    Create a backup of configuration files
  restore  Restore configuration from a backup
  list     List available backups
  delete   Delete old backups`,
	RunE: runConfig,
}

var configFlags struct {
	backupDir string
	verbose   bool
}

func init() {
	ConfigCmd.PersistentFlags().StringVarP(&configFlags.backupDir, "backup-dir", "d", "", "Custom backup directory")
	ConfigCmd.PersistentFlags().BoolVarP(&configFlags.verbose, "verbose", "v", false, "Enable verbose output")
}

// BackupCmd represents the config backup subcommand
var BackupCmd = &cobra.Command{
	Use:   "backup [config-file]",
	Short: "Create a backup of configuration files",
	Long: `Creates a timestamped backup of the specified configuration file.

Examples:
  adt config backup ~/.gitconfig
  adt config backup ~/.zshrc`,
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: completeConfigFiles,
	RunE:              runBackup,
}

// RestoreCmd represents the config restore subcommand
var RestoreCmd = &cobra.Command{
	Use:   "restore [backup-path] [target-path]",
	Short: "Restore configuration from a backup",
	Long: `Restores a configuration file from a backup.

Examples:
  adt config restore ~/.auto-dev-terminal/backups/gitconfig_2024-01-15_10-30-00.bak
  adt config restore --latest gitconfig ~/.gitconfig`,
	Args:              cobra.RangeArgs(1, 2),
	ValidArgsFunction: completeBackupFiles,
	RunE:              runRestore,
}

var restoreFlags struct {
	latest bool
}

func init() {
	RestoreCmd.Flags().BoolVarP(&restoreFlags.latest, "latest", "l", false, "Restore from the most recent backup")
}

// ListBackupsCmd represents the config list subcommand
var ListBackupsCmd = &cobra.Command{
	Use:   "list [config-name]",
	Short: "List available backups",
	Long: `Lists all available backups, or only backups for a specific config.

Examples:
  adt config list
  adt config list gitconfig`,
	Args:              cobra.MaximumNArgs(1),
	ValidArgsFunction: completeConfigNames,
	RunE:              runListBackups,
}

// DeleteBackupsCmd represents the config delete subcommand
var DeleteBackupsCmd = &cobra.Command{
	Use:   "delete [config-name]",
	Short: "Delete old backups",
	Long: `Deletes backups older than the specified number of days.

Examples:
  adt config delete gitconfig --days 30`,
	Args: cobra.ExactArgs(1),
	RunE: runDeleteBackups,
}

var deleteFlags struct {
	days int
}

func init() {
	DeleteBackupsCmd.Flags().IntVarP(&deleteFlags.days, "days", "d", 30, "Delete backups older than this many days")
}

func runConfig(cmd *cobra.Command, args []string) error {
	return cmd.Help()
}

func runBackup(cmd *cobra.Command, args []string) error {
	configPath := args[0]
	verbose := configFlags.verbose || viper.GetBool("verbose")

	if verbose {
		fmt.Printf("Backing up: %s\n", configPath)
	}

	// Create backup manager
	bm, err := config.NewBackupManager()
	if err != nil {
		return fmt.Errorf("failed to create backup manager: %w", err)
	}

	// Check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return fmt.Errorf("config file does not exist: %s", configPath)
	}

	// Perform backup
	backupPath, err := bm.BackupConfig(configPath)
	if err != nil {
		return fmt.Errorf("backup failed: %w", err)
	}

	fmt.Printf("Backup created: %s\n", backupPath)

	if verbose {
		fmt.Printf("Backup location: %s\n", bm.BackupDir())
	}

	return nil
}

func runRestore(cmd *cobra.Command, args []string) error {
	verbose := configFlags.verbose || viper.GetBool("verbose")

	// Create restore manager
	rm, err := config.NewRestoreManager()
	if err != nil {
		return fmt.Errorf("failed to create restore manager: %w", err)
	}

	var backupPath string
	var targetPath string

	if restoreFlags.latest {
		// Restore from latest backup
		if len(args) < 1 {
			return fmt.Errorf("config name required when using --latest flag")
		}
		configName := args[0]

		if verbose {
			fmt.Printf("Finding latest backup for: %s\n", configName)
		}

		backup, err := rm.FindLatestBackup(configName)
		if err != nil {
			return fmt.Errorf("no backups found for %s: %w", configName, err)
		}
		backupPath = backup
		targetPath = guessTargetPath(configName)
	} else {
		// Restore from specified backup
		backupPath = args[0]
		if len(args) >= 2 {
			targetPath = args[1]
		} else {
			// Try to guess target path from backup filename
			targetPath = guessTargetPathFromBackup(backupPath)
		}
	}

	if verbose {
		fmt.Printf("Restoring from: %s\n", backupPath)
		fmt.Printf("Restoring to: %s\n", targetPath)
	}

	// Confirm restoration
	fmt.Printf("This will overwrite: %s\n", targetPath)
	fmt.Print("Continue? [y/N]: ")
	var response string
	if _, err := fmt.Scanln(&response); err != nil {
		// If there's an error reading, assume no
		response = "n"
	}
	if response != "y" && response != "Y" {
		fmt.Println("Restore cancelled.")
		return nil
	}

	// Perform restore
	if err := rm.RestoreFromBackup(backupPath, targetPath); err != nil {
		return fmt.Errorf("restore failed: %w", err)
	}

	fmt.Printf("Configuration restored: %s\n", targetPath)

	return nil
}

func runListBackups(cmd *cobra.Command, args []string) error {
	verbose := configFlags.verbose || viper.GetBool("verbose")

	rm, err := config.NewRestoreManager()
	if err != nil {
		return fmt.Errorf("failed to create restore manager: %w", err)
	}

	if len(args) == 1 {
		// List backups for specific config
		configName := args[0]
		backups, err := rm.ListAvailableBackups(configName)
		if err != nil {
			return fmt.Errorf("failed to list backups: %w", err)
		}

		if len(backups) == 0 {
			fmt.Printf("No backups found for: %s\n", configName)
			return nil
		}

		fmt.Printf("Backups for %s:\n", configName)
		for _, backup := range backups {
			fmt.Printf("  %s - %s\n", formatBackupTime(backup.Timestamp()), backup.Path())
		}
	} else {
		// List all backups
		backups, err := rm.ListAllBackups()
		if err != nil {
			return fmt.Errorf("failed to list backups: %w", err)
		}

		if len(backups) == 0 {
			fmt.Println("No backups found.")
			return nil
		}

		fmt.Println("All backups:")
		for configName, backupList := range backups {
			fmt.Printf("\n%s:\n", configName)
			for _, backup := range backupList {
				fmt.Printf("  %s - %s\n", formatBackupTime(backup.Timestamp()), backup.Path())
			}
		}
	}

	if verbose {
		bm, err := config.NewBackupManager()
		if err == nil {
			fmt.Printf("\nBackup directory: %s\n", bm.BackupDir())
		}
	}

	return nil
}

func runDeleteBackups(cmd *cobra.Command, args []string) error {
	configName := args[0]
	verbose := configFlags.verbose || viper.GetBool("verbose")

	if configName == "all" {
		return fmt.Errorf("'all' is not supported yet. Please specify a config name.")
	}

	rm, err := config.NewRestoreManager()
	if err != nil {
		return fmt.Errorf("failed to create restore manager: %w", err)
	}

	if verbose {
		fmt.Printf("Deleting backups older than %d days for: %s\n", deleteFlags.days, configName)
	}

	deleted, err := rm.DeleteOldBackups(configName, deleteFlags.days)
	if err != nil {
		return fmt.Errorf("failed to delete backups: %w", err)
	}

	fmt.Printf("Deleted %d old backup(s)\n", deleted)

	return nil
}

// Helper functions

func formatBackupTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

func guessTargetPath(configName string) string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, configName)
}

func guessTargetPathFromBackup(backupPath string) string {
	// Extract config name from backup filename
	// Format: configName_YYYY-MM-DD_HH-MM-SS.bak
	filename := ""

	// Get just the filename from path
	for i := len(backupPath) - 1; i >= 0; i-- {
		if backupPath[i] == '/' || backupPath[i] == '\\' {
			filename = backupPath[i+1:]
			break
		}
	}

	if filename == "" {
		filename = backupPath
	}

	// Remove timestamp suffix
	for i := len(filename) - 1; i >= 0; i-- {
		if filename[i] == '_' {
			configName := filename[:i]
			return guessTargetPath(configName)
		}
	}

	return guessTargetPath(filename)
}

func completeConfigFiles(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	// Common config files
	completions := []string{
		".gitconfig",
		".zshrc",
		".bashrc",
		".bash_profile",
		".profile",
		".vimrc",
		".config/nvim/init.vim",
		".config/fish/config.fish",
	}

	return completions, cobra.ShellCompDirectiveNoFileComp
}

func completeBackupFiles(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	// Provide completion for backup files
	homeDir, _ := os.UserHomeDir()
	backupDir := filepath.Join(homeDir, ".auto-dev-terminal", "backups")

	entries, err := os.ReadDir(backupDir)
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	var completions []string
	for _, entry := range entries {
		if !entry.IsDir() {
			completions = append(completions, entry.Name())
		}
	}

	return completions, cobra.ShellCompDirectiveNoFileComp
}

func completeConfigNames(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	// Common config names
	completions := []string{
		"gitconfig",
		"zshrc",
		"bashrc",
		"profile",
		"vimrc",
		"nvim",
		"fish",
	}

	return completions, cobra.ShellCompDirectiveNoFileComp
}
