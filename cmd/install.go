// Package cmd provides the CLI command structure for Auto-Dev-Terminal.
package cmd

import (
	"fmt"

	"github.com/auto-dev-terminal/auto-dev-terminal/internal/types"
	"github.com/auto-dev-terminal/auto-dev-terminal/pkg/detector"
	"github.com/auto-dev-terminal/auto-dev-terminal/pkg/modules"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// InstallCmd represents the install command
var InstallCmd = &cobra.Command{
	Use:   "install [module1 module2 ...]",
	Short: "Install development modules",
	Long: `Install one or more development modules. If no modules are specified,
the command will prompt for module selection.

Available modules are loaded from the modules directory.`,
	Args:              cobra.RangeArgs(0, 10),
	ValidArgsFunction: completeModuleNames,
	RunE:              runInstall,
}

var installFlags struct {
	yes     bool
	sudo    bool
	dryRun  bool
	verbose bool
	force   bool
}

func init() {
	InstallCmd.Flags().BoolVarP(&installFlags.yes, "yes", "y", false, "Skip confirmation prompts")
	InstallCmd.Flags().BoolVarP(&installFlags.sudo, "sudo", "s", false, "Use sudo for installations")
	InstallCmd.Flags().BoolVarP(&installFlags.dryRun, "dry-run", "d", false, "Show what would be installed without installing")
	InstallCmd.Flags().BoolVarP(&installFlags.verbose, "verbose", "v", false, "Show detailed installation progress")
	InstallCmd.Flags().BoolVarP(&installFlags.force, "force", "f", false, "Force reinstallation of already installed modules")
}

func runInstall(cmd *cobra.Command, args []string) error {
	verbose := installFlags.verbose || viper.GetBool("verbose")

	if verbose {
		fmt.Println("Starting installation process...")
	}

	// Detect system information
	if verbose {
		fmt.Println("Detecting system configuration...")
	}

	systemInfo, err := detector.Detect()
	if err != nil {
		return fmt.Errorf("failed to detect system: %w", err)
	}

	if verbose {
		fmt.Printf("Detected OS: %s\n", systemInfo.OS)
		if len(systemInfo.PackageManagers) > 0 {
			fmt.Printf("Available package managers: %v\n", systemInfo.PackageManagers)
		}
	}

	// Determine which modules to install
	modulesToInstall, err := determineModulesToInstall(args)
	if err != nil {
		return err
	}

	if len(modulesToInstall) == 0 {
		fmt.Println("No modules selected for installation.")
		return nil
	}

	// Show what will be installed
	fmt.Println("\nModules to install:")
	for _, m := range modulesToInstall {
		fmt.Printf("  - %s: %s\n", m.Name(), m.Description())
	}

	// Confirmation prompt
	if !installFlags.yes {
		fmt.Print("\nDo you want to continue? [y/N]: ")
		var response string
		fmt.Scanln(&response)
		if response != "y" && response != "Y" {
			fmt.Println("Installation cancelled.")
			return nil
		}
	}

	// Dry run mode
	if installFlags.dryRun {
		fmt.Println("\n--- DRY RUN MODE ---")
		fmt.Println("The following actions would be performed:")
		for _, m := range modulesToInstall {
			fmt.Printf("  - Install %s\n", m.Name())
		}
		fmt.Println("--- END DRY RUN ---")
		return nil
	}

	// Perform installation
	return performInstallation(modulesToInstall, systemInfo, verbose)
}

func determineModulesToInstall(args []string) ([]modules.Module, error) {
	// Initialize built-in modules
	if err := modules.InitBuiltinModules(); err != nil {
		return nil, fmt.Errorf("failed to initialize modules: %w", err)
	}

	// If specific modules provided, validate and return them
	if len(args) > 0 {
		var selected []modules.Module
		moduleNames := make(map[string]bool)

		for _, arg := range args {
			// Check if already selected (avoid duplicates)
			if moduleNames[arg] {
				continue
			}
			moduleNames[arg] = true

			mod := modules.GetModuleByName(arg)
			if mod == nil {
				return nil, fmt.Errorf("unknown module: %s", arg)
			}
			selected = append(selected, mod)
		}

		return selected, nil
	}

	// No modules specified - return all built-in modules
	available := modules.BuiltinModules()
	if len(available) == 0 {
		fmt.Println("No modules available.")
		return nil, nil
	}

	return available, nil
}

func performInstallation(modulesToInstall []modules.Module, systemInfo *types.SystemInfo, verbose bool) error {
	// Create module installer
	installer := modules.NewModuleInstaller()
	installer = installer.WithSystemInfo(*systemInfo)
	if installFlags.sudo {
		installer = installer.WithSudo(true)
	}
	if installFlags.verbose {
		installer = installer.WithVerbose(true)
	}
	if installFlags.force {
		installer = installer.WithForce(true)
	}

	// Track results
	successCount := 0
	failCount := 0

	for _, mod := range modulesToInstall {
		fmt.Printf("\nInstalling %s...\n", mod.Name())

		result := installer.Install(mod.Name())
		if result.Success {
			successCount++
			if verbose {
				fmt.Printf("  ✓ %s installed successfully\n", mod.Name())
			}
		} else {
			failCount++
			fmt.Printf("  ✗ Failed to install %s: %s\n", mod.Name(), result.Error)
		}
	}

	// Summary
	fmt.Printf("\n--- Installation Summary ---\n")
	fmt.Printf("Successful: %d\n", successCount)
	fmt.Printf("Failed: %d\n", failCount)

	if failCount > 0 {
		return fmt.Errorf("installation completed with %d errors", failCount)
	}

	return nil
}

func completeModuleNames(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	// Initialize built-in modules
	if err := modules.InitBuiltinModules(); err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	// Get available module names
	builtins := modules.BuiltinModules()
	var completions []string
	for _, mod := range builtins {
		completions = append(completions, mod.Name())
	}

	return completions, cobra.ShellCompDirectiveNoFileComp
}

// ListModulesCmd represents the list-modules command
var ListModulesCmd = &cobra.Command{
	Use:   "list-modules",
	Short: "List available modules",
	Long:  `Lists all modules that can be installed.`,
	RunE:  runListModules,
}

func runListModules(cmd *cobra.Command, args []string) error {
	// Initialize built-in modules
	if err := modules.InitBuiltinModules(); err != nil {
		return fmt.Errorf("failed to initialize modules: %w", err)
	}

	builtins := modules.BuiltinModules()
	if len(builtins) == 0 {
		fmt.Println("No modules available.")
		return nil
	}

	fmt.Println("Available modules:")
	for _, mod := range builtins {
		fmt.Printf("  %s - %s\n", mod.Name(), mod.Description())
	}

	return nil
}
