package wizard

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

// PreviewView renders the installation preview/confirmation view.
func PreviewView(m Model) string {
	title := TitleStyle.Render("Installation Preview")
	title += "\n\n"

	// Get selected modules
	selectedModules := m.GetSelectedModuleList()

	if len(selectedModules) == 0 {
		content := ErrorBoxStyle.Render(
			Error.Render("No modules selected") + "\n\n" +
				Subtle.Render("Please select at least one module to install.\n") +
				Subtle.Render("Press Esc to go back and select modules."))
		return title + content
	}

	// System info summary
	var sysInfo string
	if m.SystemInfo != nil {
		sysInfo = fmt.Sprintf("%s %s  •  %s %s  •  %s %s",
			Bold.Render("OS:"), m.SystemInfo.OS,
			Bold.Render("Shell:"), m.SystemInfo.Shell,
			Bold.Render("Arch:"), m.SystemInfo.Arch)
	} else {
		sysInfo = Warning.Render("System info not available")
	}

	content := PanelStyle.Render(Subtle.Render(sysInfo))
	content += "\n\n"

	// Modules to install
	content += Bold.Render("Modules to install:")
	content += "\n"
	content += m.renderSelectedModules()

	// Summary
	content += "\n\n"
	summaryBox := InfoBoxStyle.Render(
		fmt.Sprintf("Total: %d module(s)\n",
			len(selectedModules)) +
			fmt.Sprintf("Package Manager: %s\n",
				Bold.Render(m.getPackageManagerInfo())) +
			fmt.Sprintf("Destination: %s",
				m.SystemInfo.HomeDir))
	content += summaryBox

	// Warning about sudo
	if m.requiresSudo() {
		content += "\n\n"
		content += Warning.Render("⚠ Some installations may require sudo privileges")
	}

	// Confirmation prompt
	content += "\n\n"
	content += Bold.Render("Ready to install?")
	content += "\n"
	content += Subtle.Render("Press Enter to confirm and start installation")
	content += "\n"
	content += Subtle.Render("Press Esc to go back and modify selection")

	footer := "\n"
	footer += ShortKeyMap()

	return title + content + footer
}

// renderSelectedModules renders the list of selected modules.
func (m Model) renderSelectedModules() string {
	selectedModules := m.GetSelectedModuleList()
	var items []string

	for i, mod := range selectedModules {
		number := ListNumberStyle.Render(fmt.Sprintf("%d.", i+1))
		name := Highlight.Render(mod.Name())
		version := Subtle.Render("v" + mod.Version())

		items = append(items, number+" "+name+" "+version)
		items = append(items, "") // Empty line
	}

	return PanelStyle.Render(lipgloss.JoinVertical(lipgloss.Left, items...))
}

// getPackageManagerInfo returns a string describing the package manager to use.
func (m Model) getPackageManagerInfo() string {
	if m.SystemInfo == nil || len(m.SystemInfo.PackageManagers) == 0 {
		return Warning.Render("Not detected")
	}

	// Return first available package manager
	return string(m.SystemInfo.PackageManagers[0])
}

// requiresSudo returns true if any selected module might require sudo.
func (m Model) requiresSudo() bool {
	// This is a simplified check - in a real implementation,
	// we'd check the specific install commands
	if m.SystemInfo == nil {
		return false
	}

	// Check if running on a system that typically needs sudo
	switch m.SystemInfo.OS {
	case "linux":
		return true
	default:
		return false
	}
}

// ConfirmView shows a simple confirmation prompt.
func ConfirmView(prompt string, confirmed bool) string {
	var options []string

	if confirmed {
		options = append(options, Success.Render("▶ Yes"))
		options = append(options, Subtle.Render("  No"))
	} else {
		options = append(options, Subtle.Render("  Yes"))
		options = append(options, Success.Render("▶ No"))
	}

	content := PanelStyle.Render(
		Bold.Render(prompt) + "\n\n" +
			lipgloss.JoinVertical(lipgloss.Left, options...))

	return content
}
