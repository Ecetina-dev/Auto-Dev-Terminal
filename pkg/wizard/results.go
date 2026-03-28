package wizard

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// ResultsView renders the installation results view.
func ResultsView(m Model) string {
	successCount := m.GetSuccessCount()
	totalCount := len(m.InstallResults)

	var title string
	if successCount == totalCount && totalCount > 0 {
		title = Success.Bold(true).Render("✓ Installation Complete!")
	} else if successCount > 0 {
		title = Warning.Bold(true).Render("⚠ Partial Success")
	} else {
		title = Error.Bold(true).Render("✗ Installation Failed")
	}

	title += "\n\n"

	// Summary
	summary := m.renderResultsSummary(successCount, totalCount)

	// Detailed results
	details := "\n\n"
	details += Bold.Render("Details:")
	details += "\n"
	details += m.renderDetailedResults()

	// Config files section
	configFiles := m.renderConfigFiles()
	if configFiles != "" {
		details += "\n\n"
		details += Bold.Render("Configuration Files:")
		details += "\n"
		details += configFiles
	}

	// Next steps
	nextSteps := m.renderNextSteps()

	footer := "\n\n"
	footer += Subtle.Render("Press Esc to return to the main menu")
	footer += "\n"
	footer += Subtle.Render("Press Ctrl+C to exit")

	return title + summary + details + nextSteps + footer
}

// renderResultsSummary shows a summary of the installation results.
func (m Model) renderResultsSummary(success, total int) string {
	var lines []string

	// Overall status
	if success == total && total > 0 {
		lines = append(lines, Success.Render("All "+fmt.Sprint(total)+" module(s) installed successfully!"))
	} else if success > 0 {
		lines = append(lines, Warning.Render(fmt.Sprintf("%d of %d modules installed successfully", success, total)))
	} else {
		lines = append(lines, Error.Render("No modules were installed successfully"))
	}

	content := lipgloss.JoinVertical(lipgloss.Left, lines...)
	return PanelStyle.Render(content)
}

// renderDetailedResults shows each module's installation result.
func (m Model) renderDetailedResults() string {
	var items []string

	for _, result := range m.InstallResults {
		var statusIcon string
		var statusText string

		if result.Success {
			statusIcon = Success.Render("✓")
			statusText = Success.Render("Success")
		} else {
			statusIcon = Error.Render("✗")
			statusText = Error.Render("Failed")
		}

		line := fmt.Sprintf("%s %s - %s", statusIcon, Bold.Render(result.Module), statusText)
		items = append(items, line)

		// Show version if available
		if result.Success && result.Version != "" {
			items = append(items, Subtle.Render("  Version: "+result.Version))
		}

		// Show error message if failed
		if !result.Success && result.Error != "" {
			errorLines := wrapText(result.Error, 60)
			for _, line := range errorLines {
				items = append(items, Error.Render("  "+line))
			}
		}

		items = append(items, "") // Empty line between items
	}

	content := lipgloss.JoinVertical(lipgloss.Left, items...)
	return PanelStyle.Render(content)
}

// renderConfigFiles shows any configuration files that were created/modified.
func (m Model) renderConfigFiles() string {
	// This is a placeholder - in a real implementation,
	// we'd track which config files were created
	var files []string

	if len(files) == 0 {
		return ""
	}

	var items []string
	for _, f := range files {
		items = append(items, Success.Render("✓")+" "+f)
	}

	return lipgloss.JoinVertical(lipgloss.Left, items...)
}

// renderNextSteps shows suggested next steps based on results.
func (m Model) renderNextSteps() string {
	var steps []string

	// Check what was installed to suggest next steps
	hasStarship := false

	for _, result := range m.InstallResults {
		if result.Success && result.Module == "starship" {
			hasStarship = true
		}
	}

	if hasStarship {
		steps = append(steps, "To enable Starship, add the following to your shell config:")
		steps = append(steps, Subtle.Render("  eval \"$(starship init <shell_name>)\""))
	}

	if len(steps) == 0 {
		return ""
	}

	header := Bold.Render("Next Steps:")
	content := lipgloss.JoinVertical(lipgloss.Left, steps...)

	return "\n\n" + header + "\n" + PanelStyle.Render(content)
}

// wrapText wraps text to fit within a given width.
func wrapText(text string, width int) []string {
	words := strings.Fields(text)
	var lines []string
	var currentLine string

	for _, word := range words {
		testLine := currentLine
		if testLine != "" {
			testLine += " "
		}
		testLine += word

		if len(testLine) > width {
			if currentLine != "" {
				lines = append(lines, currentLine)
			}
			currentLine = word
		} else {
			currentLine = testLine
		}
	}

	if currentLine != "" {
		lines = append(lines, currentLine)
	}

	return lines
}

// ErrorResultsView shows the error details when installation fails.
func ErrorResultsView(err error) string {
	title := Error.Bold(true).Render("Installation Error")
	title += "\n\n"

	content := ErrorBoxStyle.Render(
		Error.Render("An error occurred during installation:\n\n") +
			Error.Render(err.Error()))

	footer := "\n\n"
	footer += Subtle.Render("Press Esc to return to the main menu")

	return title + content + footer
}
