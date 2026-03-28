// Package main is the entry point for the Auto-Dev-Terminal CLI.
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/auto-dev-terminal/auto-dev-terminal/cmd"
	"github.com/auto-dev-terminal/auto-dev-terminal/internal/constants"
	"github.com/spf13/viper"
)

func main() {
	// Initialize Viper configuration
	initializeConfig()

	// Execute root command
	if err := cmd.RootCmd.ExecuteContext(context.Background()); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func initializeConfig() {
	// Set default values
	viper.SetDefault("verbose", false)
	viper.SetDefault("backup_dir", constants.DefaultBackupDir)
	viper.SetDefault("config_dir", constants.DefaultConfigDir)
	viper.SetDefault("modules_dir", constants.DefaultModulesDir)

	// Environment variable support
	viper.SetEnvPrefix("ADT")
	viper.AutomaticEnv()

	// Config file support
	viper.SetConfigName(constants.ConfigFileName)
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME/.auto-dev-terminal/")
	viper.AddConfigPath(".")
}
