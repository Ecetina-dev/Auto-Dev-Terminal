package wizard

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// UI styles for the wizard TUI.
var (
	// General styles
	Subtle    = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	Bold      = lipgloss.NewStyle().Bold(true)
	Highlight = lipgloss.NewStyle().Foreground(lipgloss.Color("86"))
	Error     = lipgloss.NewStyle().Foreground(lipgloss.Color("204"))
	Success   = lipgloss.NewStyle().Foreground(lipgloss.Color("76"))
	Warning   = lipgloss.NewStyle().Foreground(lipgloss.Color("226"))

	// Container styles
	WindowStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62")).
			Padding(1, 2)

	PanelStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240")).
			Padding(1, 2)

	// Menu styles
	MenuItemStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("white"))

	SelectedMenuItemStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("86")).
				Bold(true)

	// Title styles
	TitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("99")).
			Bold(true).
			Padding(0, 1)

	LargeTitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("99")).
			Bold(true).
			Padding(0, 1)

	// Status styles
	StatusBarStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			Background(lipgloss.Color("235")).
			Padding(0, 1)

	StatusGoodStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("76"))

	StatusBadStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("204"))

	// Progress bar styles
	ProgressBarEmpty = lipgloss.NewStyle().
			Foreground(lipgloss.Color("236"))

	ProgressBarFull = lipgloss.NewStyle().
			Foreground(lipgloss.Color("86"))

	// List styles
	ListNumberStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("86"))

	CheckboxChecked = lipgloss.NewStyle().
				Foreground(lipgloss.Color("76")).SetString("✓")

	CheckboxUnchecked = lipgloss.NewStyle().
				Foreground(lipgloss.Color("241")).SetString("○")

	// Help text styles
	HelpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("245"))

	// Box styles
	InfoBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("62")).
			BorderBackground(lipgloss.Color("235")).
			Background(lipgloss.Color("236")).
			Padding(1, 2)

	ErrorBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("204")).
			Background(lipgloss.Color("235")).
			Padding(1, 2)
)

// ProgressBar creates a text-based progress bar.
func ProgressBar(current, total int, width int) string {
	if total <= 0 {
		total = 1
	}

	percent := float64(current) / float64(total)
	filled := int(float64(width) * percent)

	bar := ProgressBarFull.Render(RepeatChar(filled, '█'))
	bar += ProgressBarEmpty.Render(RepeatChar(width-filled, '░'))

	return bar
}

// RepeatChar returns a string of n copies of the given character.
func RepeatChar(n int, char rune) string {
	result := make([]rune, n)
	for i := range result {
		result[i] = char
	}
	return string(result)
}

// Centered centers text within a given width.
func Centered(text string, width int) string {
	lines := strings.Split(text, "\n")
	var result []string

	for _, line := range lines {
		padding := (width - lipgloss.Width(line)) / 2
		if padding < 0 {
			padding = 0
		}
		result = append(result, RepeatChar(padding, ' ')+line)
	}

	return lipgloss.JoinVertical(lipgloss.Left, result...)
}

// Truncate truncates text to fit within a given width.
func Truncate(text string, width int) string {
	if lipgloss.Width(text) <= width {
		return text
	}

	// Find a good truncation point
	truncated := text
	for lipgloss.Width(truncated) > width-3 {
		truncated = truncated[:len(truncated)-1]
	}
	return truncated + "..."
}

// FormatList formats a slice of strings as a numbered list.
func FormatList(items []string, startNum int) string {
	var result []string
	for i, item := range items {
		num := ListNumberStyle.Render(formatNumber(i + startNum))
		result = append(result, num+". "+item)
	}
	return lipgloss.JoinVertical(lipgloss.Left, result...)
}

// formatNumber formats a number for display.
func formatNumber(n int) string {
	if n < 10 {
		return " " + string(rune('0'+n))
	}
	return string(rune('0' + n/10%10)) + string(rune('0'+n%10))
}

// KeyMap displays the help text for navigation keys.
func KeyMap() string {
	return HelpStyle.Render("↑/↓: Navigate  •  Enter: Select  •  Esc: Back  •  Ctrl+C: Exit")
}

// ShortKeyMap displays a shorter version of the key help.
func ShortKeyMap() string {
	return HelpStyle.Render("↑↓: Navigate  •  Enter: Select  •  Esc: Back  •  Ctrl+C: Exit")
}
