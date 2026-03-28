package wizard

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

// DetectionView renders the system detection view.
func DetectionView(m Model) string {
	title := TitleStyle.Render("System Detection")
	title += "\n\n"
	
	var content string
	
	// Check if detection is running (SystemInfo is nil)
	if m.SystemInfo == nil {
		if m.DetectionErr != nil {
			// Error state
			content = ErrorBoxStyle.Render(
				fmt.Sprintf("Detection Error: %s\n\n%s",
					Error.Render(m.DetectionErr.Error()),
					Subtle.Render("Press Esc to go back and try again")))
		} else {
			// Running state
			content = PanelStyle.Render(
				Highlight.Render("Detecting system configuration..."))
			content += "\n\n"
			content += Subtle.Render("Please wait while we analyze your environment.\n")
			content += Subtle.Render("This may take a few seconds...")
		}
	} else {
		// Show results
		content = m.renderDetectionResults()
	}
	
	footer := "\n"
	footer += ShortKeyMap()
	
	return title + content + footer
}

// renderDetectionResults formats and displays the detection results.
func (m Model) renderDetectionResults() string {
	var sections []string
	
	// OS Section
	osInfo := fmt.Sprintf("%s: %s", Bold.Render("Operating System"), string(m.SystemInfo.OS))
	if m.SystemInfo.Distro != "" {
		osInfo += fmt.Sprintf(" (%s", m.SystemInfo.Distro)
		if m.SystemInfo.DistroVersion != "" {
			osInfo += fmt.Sprintf(" %s", m.SystemInfo.DistroVersion)
		}
		osInfo += ")"
	}
	sections = append(sections, Success.Render("✓")+" "+osInfo)
	
	// Shell Section
	shellInfo := fmt.Sprintf("%s: %s", Bold.Render("Shell"), string(m.SystemInfo.Shell))
	if m.SystemInfo.ShellVersion != "" {
		shellInfo += fmt.Sprintf(" (%s)", m.SystemInfo.ShellVersion)
	}
	sections = append(sections, Success.Render("✓")+" "+shellInfo)
	
	// Architecture
	archInfo := fmt.Sprintf("%s: %s", Bold.Render("Architecture"), m.SystemInfo.Arch)
	sections = append(sections, Success.Render("✓")+" "+archInfo)
	
	// Home Directory
	homeInfo := fmt.Sprintf("%s: %s", Bold.Render("Home Directory"), m.SystemInfo.HomeDir)
	sections = append(sections, Success.Render("✓")+" "+homeInfo)
	
	// Package Managers Section
	pkgMgrsHeader := Bold.Render("Package Managers")
	if len(m.SystemInfo.PackageManagers) > 0 {
		var pkgMgrList []string
		for _, pm := range m.SystemInfo.PackageManagers {
			pkgMgrList = append(pkgMgrList, string(pm))
		}
		pkgMgrSection := fmt.Sprintf("%s: %s", pkgMgrsHeader, Highlight.Render("Available"))
		pkgMgrSection += "\n    "
		pkgMgrSection += Subtle.Render("Found: ")
		pkgMgrSection += FormatList(pkgMgrList, 1)
		sections = append(sections, pkgMgrSection)
	} else {
		sections = append(sections, fmt.Sprintf("%s: %s", pkgMgrsHeader, Warning.Render("None detected")))
	}
	
	// Combine sections
	result := PanelStyle.Render(lipgloss.JoinVertical(lipgloss.Left, sections...))
	result += "\n\n"
	result += Subtle.Render("Press Esc to return to the main menu.")
	result += "\n"
	result += Subtle.Render("Press Enter to continue to module selection.")
	
	return result
}

// renderRunningState shows the detection in progress.
func (m Model) renderRunningState() string {
	return PanelStyle.Render(
		Highlight.Render("Detecting system...") + "\n\n" +
		Subtle.Render("Analyzing your development environment..."))
}
