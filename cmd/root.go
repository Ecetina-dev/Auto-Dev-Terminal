// Package cmd provides the CLI command structure for Auto-Dev-Terminal.
package cmd

import (
	"fmt"
	"os"

	"github.com/auto-dev-terminal/auto-dev-terminal/internal/constants"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "auto-dev-terminal",
	Short: "Automated development environment setup tool",
	Long: `Auto-Dev-Terminal is a CLI tool that automates development environment
setup through intelligent system detection, package manager integration,
and an interactive wizard interface.`,
	Version: constants.Version,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Set verbose mode if flag is provided
		if verbose, _ := cmd.Flags().GetBool("verbose"); verbose {
			fmt.Println("Verbose mode enabled")
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately
func Execute() error {
	return RootCmd.Execute()
}

func init() {
	// Add global persistent flags
	RootCmd.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose output")
	RootCmd.PersistentFlags().StringP("config", "c", "", "Custom config file path")

	// Bind viper to flags
	viper.BindPFlag("verbose", RootCmd.PersistentFlags().Lookup("verbose"))
	viper.BindPFlag("config", RootCmd.PersistentFlags().Lookup("config"))

	// Add subcommands
	RootCmd.AddCommand(DetectCmd)
	RootCmd.AddCommand(QuickDetectCmd)
	RootCmd.AddCommand(InstallCmd)
	RootCmd.AddCommand(ListModulesCmd)
	RootCmd.AddCommand(WizardCmd)
	RootCmd.AddCommand(WizardDetectCmd)
	RootCmd.AddCommand(ConfigCmd)
	RootCmd.AddCommand(BackupCmd)
	RootCmd.AddCommand(RestoreCmd)
	RootCmd.AddCommand(ListBackupsCmd)
	RootCmd.AddCommand(DeleteBackupsCmd)
}

// GetVerbose returns the verbose flag value
func GetVerbose() bool {
	return viper.GetBool("verbose")
}

// GetConfigPath returns the custom config path if provided
func GetConfigPath() string {
	return viper.GetString("config")
}

// setupBackupDir ensures the backup directory exists
func setupBackupDir() error {
	backupDir := viper.GetString("backup_dir")
	if _, err := os.Stat(backupDir); os.IsNotExist(err) {
		if err := os.MkdirAll(backupDir, 0755); err != nil {
			return fmt.Errorf("failed to create backup directory: %w", err)
		}
	}
	return nil
}
