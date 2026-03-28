package wizard

import (
	"fmt"
	"strings"

	"github.com/auto-dev-terminal/auto-dev-terminal/pkg/detector"
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// WelcomeView renders the welcome screen.
func WelcomeView(m Model) string {
	// Title
	title := LargeTitleStyle.Render("Auto Dev Terminal")
	title += "\n\n"

	// Tagline
	tagline := Subtle.Render("Automated Development Environment Setup")
	tagline += "\n\n"

	// Description
	description := `
This wizard will help you set up your development environment
by detecting your system configuration and installing useful
development tools and configurations.
`
	desc := Subtle.Render(description)

	// Features
	features := PanelStyle.Render(
		Bold.Render("Features:") + "\n" +
			Success.Render("✓") + " Automatic system detection\n" +
			Success.Render("✓") + " Multi-platform support (Linux, macOS, Windows)\n" +
			Success.Render("✓") + " Modular installation\n" +
			Success.Render("✓") + " Configuration backup\n" +
			Success.Render("✓") + " Interactive TUI")

	// Prompt
	prompt := "\n\n"
	prompt += Highlight.Render("Press Enter to continue...")
	prompt += "\n"
	prompt += Subtle.Render("Press Ctrl+C to exit")

	return title + tagline + desc + features + prompt
}

// ExitView renders the exit screen.
func ExitView() string {
	content := TitleStyle.Render("Goodbye!")
	content += "\n\n"
	content += Subtle.Render("Thank you for using Auto Dev Terminal.")
	content += "\n\n"
	content += Success.Render("✓ Session complete")
	content += "\n\n"
	content += Subtle.Render("Press any key to exit...")

	return lipgloss.NewStyle().
		Width(60).
		Height(15).
		Render(content)
}

// NewProgram creates a new Bubble Tea program with the wizard model.
func NewProgram() *tea.Program {
	return tea.NewProgram(NewModel())
}

// NewProgramWithOptions creates a new Bubble Tea program with custom options.
func NewProgramWithOptions(verbose bool) *tea.Program {
	model := NewModel()
	model.Verbose = verbose
	return tea.NewProgram(model)
}

// RunWizard runs the wizard interactively.
func RunWizard() error {
	program := tea.NewProgram(NewModel(), tea.WithAltScreen())

	if err := program.Start(); err != nil {
		return fmt.Errorf("failed to start wizard: %w", err)
	}

	return nil
}

// RunWizardWithConfig runs the wizard with custom configuration.
func RunWizardWithConfig(configPath string, verbose bool) error {
	model := NewModel()
	model.ConfigPath = configPath
	model.Verbose = verbose

	program := tea.NewProgram(model, tea.WithAltScreen())

	if err := program.Start(); err != nil {
		return fmt.Errorf("failed to start wizard: %w", err)
	}

	return nil
}

// StartWizard starts the wizard as a Bubble Tea program.
// This is the main entry point for the wizard subcommand.
func StartWizard() error {
	return RunWizard()
}

// DetectSystem performs system detection and returns a human-readable summary.
func DetectSystem() (string, error) {
	info, err := detector.Detect()
	if err != nil {
		return "", fmt.Errorf("detection failed: %w", err)
	}

	var b strings.Builder
	b.WriteString("System Detection Results:\n")
	b.WriteString("========================\n")
	b.WriteString(fmt.Sprintf("OS: %s\n", info.OS))

	if info.OS == "linux" && info.Distro != "" {
		b.WriteString(fmt.Sprintf("Distribution: %s", info.Distro))
		if info.DistroVersion != "" {
			b.WriteString(fmt.Sprintf(" (%s)", info.DistroVersion))
		}
		b.WriteString("\n")
	}

	b.WriteString(fmt.Sprintf("Architecture: %s\n", info.Arch))
	b.WriteString(fmt.Sprintf("Shell: %s", info.Shell))
	if info.ShellVersion != "" {
		b.WriteString(fmt.Sprintf(" (%s)", info.ShellVersion))
	}
	b.WriteString("\n")

	if len(info.PackageManagers) > 0 {
		b.WriteString("Package Managers: ")
		for i, pm := range info.PackageManagers {
			if i > 0 {
				b.WriteString(", ")
			}
			b.WriteString(string(pm))
		}
		b.WriteString("\n")
	}

	b.WriteString(fmt.Sprintf("User: %s\n", info.Username))
	b.WriteString(fmt.Sprintf("Home: %s\n", info.HomeDir))
	b.WriteString(fmt.Sprintf("Hostname: %s\n", info.Hostname))

	return b.String(), nil
}
