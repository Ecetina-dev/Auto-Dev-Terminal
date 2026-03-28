package wizard

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

// ModuleSelectionView renders the module selection view.
func ModuleSelectionView(m Model) string {
	title := TitleStyle.Render("Select Modules")
	title += "\n\n"

	var content string

	// Check if modules are available
	if len(m.AvailableModules) == 0 {
		content = ErrorBoxStyle.Render(
			Error.Render("No modules available") + "\n\n" +
				Subtle.Render("Please run detection first to load available modules.\n") +
				Subtle.Render("Press Esc to go back to the main menu."))
		return title + content
	}

	// Show warning if no system info
	if m.SystemInfo == nil {
		content += Warning.Render("⚠ System not detected") + "\n"
		content += Subtle.Render("Run detection first to ensure modules are compatible.\n\n")
	}

	// Module list
	content += m.renderModuleList()

	// Selection summary
	selectedCount := len(m.GetSelectedModuleList())
	if selectedCount > 0 {
		content += "\n"
		content += PanelStyle.Render(
			fmt.Sprintf("%s %d module(s) selected",
				Highlight.Render("→"),
				selectedCount))
	}

	// Actions hint
	content += "\n\n"
	content += Subtle.Render("Press Space to toggle selection • Enter to continue • Esc to go back")

	footer := "\n"
	footer += ShortKeyMap()

	return title + content + footer
}

// renderModuleList renders the list of available modules.
func (m Model) renderModuleList() string {
	var items []string

	for i, mod := range m.AvailableModules {
		// Checkbox
		var checkbox string
		if m.SelectedModules[mod.Name()] {
			checkbox = CheckboxChecked.Render(" ")
		} else {
			checkbox = CheckboxUnchecked.Render(" ")
		}

		// Cursor position styling
		var name string
		if i == m.ModuleCursor {
			name = SelectedMenuItemStyle.Render(mod.Name())
		} else {
			name = MenuItemStyle.Render(mod.Name())
		}

		// Description (truncated if needed)
		description := Subtle.Render(Truncate(mod.Description(), 50))

		// Version badge
		version := Subtle.Render("v" + mod.Version())

		// Combine
		items = append(items, checkbox+" "+name+" "+version)
		items = append(items, "   "+description)
		items = append(items, "") // Empty line between items
	}

	// Create a bordered list
	listContent := lipgloss.JoinVertical(lipgloss.Left, items...)

	return PanelStyle.Render(listContent)
}

// ModuleCheckbox returns the checkbox string for a module.
func ModuleCheckbox(selected bool) string {
	if selected {
		return CheckboxChecked.Render("[X]")
	}
	return CheckboxUnchecked.Render("[ ]")
}

// ModuleStatus returns the status text for a module.
func ModuleStatus(modName string, selected bool, installed bool) string {
	if installed {
		return Success.Render("installed")
	}
	if selected {
		return Highlight.Render("selected")
	}
	return Subtle.Render("available")
}
