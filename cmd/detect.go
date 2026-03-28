// Package cmd provides the CLI command structure for Auto-Dev-Terminal.
package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/auto-dev-terminal/auto-dev-terminal/internal/types"
	"github.com/auto-dev-terminal/auto-dev-terminal/pkg/detector"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// DetectCmd represents the detect command
var DetectCmd = &cobra.Command{
	Use:   "detect",
	Short: "Detect system information",
	Long: `Detects and displays system information including OS, distribution,
shell, architecture, and available package managers.`,
	RunE: runDetect,
}

var detectFlags struct {
	json    bool
	verbose bool
}

func init() {
	DetectCmd.Flags().BoolVarP(&detectFlags.json, "json", "j", false, "Output in JSON format")
	DetectCmd.Flags().BoolVarP(&detectFlags.verbose, "verbose", "v", false, "Show detailed detection information")
}

func runDetect(cmd *cobra.Command, args []string) error {
	verbose := detectFlags.verbose || viper.GetBool("verbose")

	if verbose {
		fmt.Println("Detecting system information...")
	}

	// Perform system detection
	info, err := detector.Detect()
	if err != nil {
		return fmt.Errorf("failed to detect system: %w", err)
	}

	// Output based on format flag
	if detectFlags.json {
		return outputJSON(info)
	}

	// Human-readable output
	return outputText(info, verbose)
}

func outputJSON(info *types.SystemInfo) error {
	output, err := json.MarshalIndent(info, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	fmt.Println(string(output))
	return nil
}

func outputText(info *types.SystemInfo, verbose bool) error {
	fmt.Println("System Detection Results:")
	fmt.Println("========================")

	fmt.Printf("OS: %s\n", info.OS)

	if info.OS == types.OSLinux && info.Distro != "" {
		fmt.Printf("Distribution: %s", info.Distro)
		if info.DistroVersion != "" {
			fmt.Printf(" (%s)", info.DistroVersion)
		}
		fmt.Println()
	}

	fmt.Printf("Architecture: %s\n", info.Arch)

	fmt.Printf("Shell: %s", info.Shell)
	if info.ShellVersion != "" {
		fmt.Printf(" (%s)", info.ShellVersion)
	}
	fmt.Println()

	if len(info.PackageManagers) > 0 {
		fmt.Print("Package Managers: ")
		for i, pm := range info.PackageManagers {
			if i > 0 {
				fmt.Print(", ")
			}
			fmt.Print(pm)
		}
		fmt.Println()
	}

	if verbose {
		fmt.Printf("\nUser: %s\n", info.Username)
		fmt.Printf("Home Directory: %s\n", info.HomeDir)
		fmt.Printf("Hostname: %s\n", info.Hostname)
	}

	return nil
}

// QuickDetectCmd represents a quick detect command for scripting
var QuickDetectCmd = &cobra.Command{
	Use:   "detect-quick",
	Short: "Quick OS detection for scripts",
	Long:  `Performs a quick OS detection and outputs the OS name. Useful for scripting.`,
	RunE:  runQuickDetect,
}

func runQuickDetect(cmd *cobra.Command, args []string) error {
	os := detector.DetectSimple()
	fmt.Println(os)
	return nil
}
