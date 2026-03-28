package wizard

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

// MainMenuOption represents an option in the main menu.
type MainMenuOption struct {
	Title       string
	Description string
	Screen      Screen
}

// MainMenuOptions returns the available main menu options.
func MainMenuOptions() []MainMenuOption {
	return []MainMenuOption{
		{
			Title:       "Run Detection",
			Description: "Detect your system configuration (OS, Shell, Package Managers)",
			Screen:      ScreenDetection,
		},
		{
			Title:       "Setup Modules",
			Description: "Select and install development modules (Starship, Oh-My-Zsh, etc.)",
			Screen:      ScreenModuleSelection,
		},
		{
			Title:       "View Configuration",
			Description: "View and manage your configuration files and backups",
			Screen:      ScreenResults, // Placeholder - would be ScreenConfig
		},
		{
			Title:       "Exit",
			Description: "Exit the wizard",
			Screen:      ScreenExit,
		},
	}
}

// MainMenuView renders the main menu view.
func MainMenuView(m Model) string {
	options := MainMenuOptions()
	
	// Header
	header := LargeTitleStyle.Render("Auto Dev Terminal")
	header += "\n\n"
	header += Subtle.Render("Automated Development Environment Setup\n")
	header += "\n"
	
	// Menu
	var menuItems []string
	for i, opt := range options {
		if i == m.MenuCursor {
			menuItems = append(menuItems, SelectedMenuItemStyle.Render("▶ "+opt.Title))
			menuItems = append(menuItems, Subtle.Render("  "+opt.Description))
		} else {
			menuItems = append(menuItems, MenuItemStyle.Render("  "+opt.Title))
			menuItems = append(menuItems, Subtle.Render("  "+opt.Description))
		}
		menuItems = append(menuItems, "") // Empty line between options
	}
	
	menu := lipgloss.JoinVertical(lipgloss.Left, menuItems...)
	
	// Footer
	footer := "\n\n"
	footer += KeyMap()
	
	// Combine
	content := header + menu + footer
	
	// Apply container style
	return lipgloss.NewStyle().
		Width(70).
		Height(20).
		Render(content)
}

// MainMenuHelp returns additional help text for the main menu.
func MainMenuHelp() string {
	return fmt.Sprintf("\n%s\n",
		Subtle.Render("Press Enter to select an option • Use arrow keys to navigate"))
}
