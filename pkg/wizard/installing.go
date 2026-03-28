package wizard

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// InstallingView renders the installation progress view.
func InstallingView(m Model) string {
	title := TitleStyle.Render("Installing Modules")
	title += "\n\n"

	// Progress bar
	progress := ProgressBar(m.InstallProgress, m.InstallTotal, 40)

	progressText := fmt.Sprintf("Progress: %d/%d",
		m.InstallProgress,
		m.InstallTotal)

	content := PanelStyle.Render(
		progress + "\n\n" +
			Bold.Render(progressText))

	// Current module being installed
	selectedModules := m.GetSelectedModuleList()
	if m.InstallProgress < len(selectedModules) {
		currentMod := selectedModules[m.InstallProgress]
		content += "\n\n"
		content += Highlight.Render("→ Installing: ") + Bold.Render(currentMod.Name())
	}

	// Results so far
	if len(m.InstallResults) > 0 {
		content += "\n\n"
		content += Bold.Render("Results:")
		content += "\n"
		content += m.renderInstallResults()
	}

	// Cancel hint
	content += "\n\n"
	content += Subtle.Render("Installation in progress... Please wait.")

	footer := "\n"
	footer += Subtle.Render("Ctrl+C to cancel (not recommended)")

	return title + content + footer
}

// renderInstallResults shows the results of completed installations.
func (m Model) renderInstallResults() string {
	var items []string

	for _, result := range m.InstallResults {
		var status string
		if result.Success {
			status = Success.Render("✓")
		} else {
			status = Error.Render("✗")
		}

		moduleName := result.Module
		if result.Version != "" {
			moduleName += fmt.Sprintf(" (v%s)", result.Version)
		}

		items = append(items, fmt.Sprintf("  %s %s", status, moduleName))

		// Show error message if failed
		if !result.Success && result.Error != "" {
			errorMsg := Subtle.Render("    Error: " + Truncate(result.Error, 50))
			items = append(items, errorMsg)
		}
	}

	return lipgloss.JoinVertical(lipgloss.Left, items...)
}

// InstallProgressView shows detailed progress for a single module.
func InstallProgressView(moduleName string, step string, progress int) string {
	content := PanelStyle.Render(
		Bold.Render("Installing: ") + Highlight.Render(moduleName) + "\n\n" +
			Subtle.Render(step) + "\n\n" +
			ProgressBar(progress, 100, 30))

	return content
}

// InstallationStatus returns a status message for the current installation.
func InstallationStatus(m Model) string {
	if m.IsInstalling {
		return Highlight.Render("Installing...")
	}

	// Check results
	successCount := m.GetSuccessCount()
	totalCount := len(m.InstallResults)

	if successCount == totalCount {
		return Success.Render("All installations complete!")
	}

	if successCount > 0 {
		return Warning.Render(fmt.Sprintf("%d/%d successful", successCount, totalCount))
	}

	return Error.Render("Installation failed")
}

// RenderCommandOutput renders command output with ANSI codes stripped.
func RenderCommandOutput(output string) string {
	// Strip ANSI codes for clean display
	clean := stripANSI(output)

	// Truncate if too long
	lines := strings.Split(clean, "\n")
	if len(lines) > 10 {
		lines = lines[:10]
		lines = append(lines, Subtle.Render("... (output truncated)"))
	}

	return strings.Join(lines, "\n")
}

// stripANSI removes ANSI escape codes from a string.
func stripANSI(s string) string {
	// Simple ANSI stripping - in production, use a proper library
	result := ""
	inEscape := false

	for _, r := range s {
		if r == '\033' {
			inEscape = true
			continue
		}
		if inEscape {
			if r == 'm' { // End of escape sequence
				inEscape = false
			}
			continue
		}
		result += string(r)
	}

	return result
}
