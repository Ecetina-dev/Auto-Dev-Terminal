// Package cmd provides the CLI command structure for Auto-Dev-Terminal.
package cmd

import (
	"fmt"

	"github.com/auto-dev-terminal/auto-dev-terminal/pkg/wizard"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// WizardCmd represents the wizard command
var WizardCmd = &cobra.Command{
	Use:   "wizard",
	Short: "Launch interactive TUI wizard",
	Long: `Starts an interactive terminal user interface (TUI) wizard that
guides you through system detection and module installation.

The wizard provides:
- Automatic system detection
- Module selection
- Installation progress
- Configuration backup`,
	RunE: runWizard,
}

var wizardFlags struct {
	configPath string
	verbose    bool
}

func init() {
	WizardCmd.Flags().StringVarP(&wizardFlags.configPath, "config", "c", "", "Path to custom configuration file")
	WizardCmd.Flags().BoolVarP(&wizardFlags.verbose, "verbose", "v", false, "Enable verbose output")
}

func runWizard(cmd *cobra.Command, args []string) error {
	verbose := wizardFlags.verbose || viper.GetBool("verbose")

	if verbose {
		fmt.Println("Starting wizard in verbose mode...")
	}

	// Run the wizard
	var err error
	if wizardFlags.configPath != "" {
		err = wizard.RunWizardWithConfig(wizardFlags.configPath, verbose)
	} else {
		err = wizard.RunWizard()
	}

	if err != nil {
		return fmt.Errorf("wizard error: %w", err)
	}

	return nil
}

// WizardDetectCmd runs only the detection step of the wizard
var WizardDetectCmd = &cobra.Command{
	Use:   "wizard-detect",
	Short: "Run detection wizard",
	Long:  `Runs just the detection portion of the wizard and displays results.`,
	RunE:  runWizardDetect,
}

func runWizardDetect(cmd *cobra.Command, args []string) error {
	fmt.Println("Detecting system configuration...")

	// Use wizard's detection
	detection, err := wizard.DetectSystem()
	if err != nil {
		return fmt.Errorf("detection failed: %w", err)
	}

	fmt.Println("\n" + detection)

	return nil
}
