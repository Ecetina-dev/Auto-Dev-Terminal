// Package wizard provides the interactive Bubble Tea TUI for the auto-dev-terminal application.
package wizard

import (
	"github.com/auto-dev-terminal/auto-dev-terminal/internal/types"
	"github.com/auto-dev-terminal/auto-dev-terminal/pkg/detector"
	"github.com/auto-dev-terminal/auto-dev-terminal/pkg/installer"
	"github.com/auto-dev-terminal/auto-dev-terminal/pkg/modules"

	"github.com/charmbracelet/bubbletea"
)

// Screen represents the current screen state in the wizard.
type Screen int

const (
	ScreenWelcome Screen = iota
	ScreenMainMenu
	ScreenDetection
	ScreenModuleSelection
	ScreenPreview
	ScreenInstalling
	ScreenResults
	ScreenExit
)

// Model is the main Bubble Tea model for the wizard.
type Model struct {
	// Screen state
	CurrentScreen   Screen
	PreviousScreen  Screen

	// System information
	SystemInfo   *types.SystemInfo
	DetectionErr error

	// Module selection
	AvailableModules []modules.Module
	SelectedModules map[string]bool
	ModuleCursor   int

	// Installation state
	InstallProgress int
	InstallTotal   int
	InstallResults []types.InstallResult
	InstallError   error
	IsInstalling  bool

	// Navigation
	MenuCursor    int
	ConfirmChoice bool

	// Configuration
	ConfigPath string
	Verbose    bool
}

// NewModel creates a new wizard model with default values.
func NewModel() Model {
	return Model{
		CurrentScreen:   ScreenWelcome,
		PreviousScreen:  -1,
		SystemInfo:      nil,
		DetectionErr:    nil,
		AvailableModules: nil,
		SelectedModules: make(map[string]bool),
		ModuleCursor:    0,
		InstallProgress: 0,
		InstallTotal:    0,
		InstallResults:  nil,
		InstallError:    nil,
		IsInstalling:    false,
		MenuCursor:      0,
		ConfirmChoice:   false,
		ConfigPath:      "",
		Verbose:         false,
	}
}

// Init implements tea.Model Init method - runs any necessary initialization.
func (m Model) Init() tea.Cmd {
	// Load available modules when starting
	loadedModules := modules.List()
	if len(loadedModules) > 0 {
		m.AvailableModules = loadedModules
	}
	return nil
}

// Update implements tea.Model Update method.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return handleKeyMsg(&m, msg)

	case tea.WindowSizeMsg:
		return handleWindowSize(&m, msg)

	case detectionRequestedMsg:
		// Run detection
		err := m.StartDetection()
		if err != nil {
			m.DetectionErr = err
		}
		return m, nil

	case installCompleteMsg:
		return handleInstallComplete(&m)

	default:
		return m, nil
	}
}

// handleKeyMsg handles keyboard input for the wizard.
func handleKeyMsg(m *Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyCtrlC:
		m.CurrentScreen = ScreenExit
		return *m, tea.Quit

	case tea.KeyEsc:
		return handleEscape(m)

	case tea.KeyEnter:
		return handleEnter(m)

	case tea.KeyUp:
		return handleUp(m)

	case tea.KeyDown:
		return handleDown(m)

	case tea.KeySpace:
		return handleSpace(m)

	default:
		return *m, nil
	}
}

// handleEscape handles the Escape key based on current screen.
func handleEscape(m *Model) (Model, tea.Cmd) {
	switch m.CurrentScreen {
	case ScreenWelcome:
		m.CurrentScreen = ScreenExit
		return *m, tea.Quit

	case ScreenMainMenu:
		m.CurrentScreen = ScreenExit
		return *m, tea.Quit

	case ScreenDetection:
		m.SetScreen(ScreenMainMenu)

	case ScreenModuleSelection:
		m.SetScreen(ScreenMainMenu)

	case ScreenPreview:
		m.SetScreen(ScreenModuleSelection)

	case ScreenInstalling:
		// Can't go back during installation
		return *m, nil

	case ScreenResults:
		m.ResetSelection()
		m.SetScreen(ScreenMainMenu)

	default:
		m.SetScreen(ScreenMainMenu)
	}

	return *m, nil
}

// handleEnter handles the Enter key based on current screen.
func handleEnter(m *Model) (Model, tea.Cmd) {
	switch m.CurrentScreen {
	case ScreenWelcome:
		m.SetScreen(ScreenMainMenu)
		return *m, nil

	case ScreenMainMenu:
		return handleMainMenuSelection(m)

	case ScreenDetection:
		// If we have system info, go to module selection
		if m.SystemInfo != nil {
			m.SetScreen(ScreenModuleSelection)
		}
		return *m, nil

	case ScreenModuleSelection:
		// Go to preview if we have selections
		if m.HasSelection() {
			m.SetScreen(ScreenPreview)
		}
		return *m, nil

	case ScreenPreview:
		// Start installation
		m.StartInstallation()
		return *m, startInstallationWorker(m)

	case ScreenInstalling:
		// Can't interact during installation
		return *m, nil

	case ScreenResults:
		m.ResetSelection()
		m.SetScreen(ScreenMainMenu)
		return *m, nil

	default:
		return *m, nil
	}
}

// handleMainMenuSelection handles menu item selection.
func handleMainMenuSelection(m *Model) (Model, tea.Cmd) {
	options := MainMenuOptions()

	if m.MenuCursor < 0 || m.MenuCursor >= len(options) {
		return *m, nil
	}

	selectedOption := options[m.MenuCursor]

	switch selectedOption.Screen {
	case ScreenDetection:
		m.SetScreen(ScreenDetection)
		// Start detection in background
		return *m, runDetection

	case ScreenModuleSelection:
		// Need system info first
		if m.SystemInfo == nil {
			m.SetScreen(ScreenDetection)
			return *m, runDetection
		}
		m.SetScreen(ScreenModuleSelection)

	case ScreenExit:
		m.CurrentScreen = ScreenExit
		return *m, tea.Quit

	default:
		m.SetScreen(selectedOption.Screen)
	}

	return *m, nil
}

// handleUp handles upward navigation.
func handleUp(m *Model) (Model, tea.Cmd) {
	switch m.CurrentScreen {
	case ScreenMainMenu:
		m.MenuCursor--
		if m.MenuCursor < 0 {
			m.MenuCursor = len(MainMenuOptions()) - 1
		}

	case ScreenModuleSelection:
		m.ModuleCursor--
		if m.ModuleCursor < 0 {
			m.ModuleCursor = len(m.AvailableModules) - 1
		}

	case ScreenPreview:
		// Toggle confirmation
		m.ConfirmChoice = !m.ConfirmChoice

	case ScreenInstalling:
		// Can't navigate during installation
		return *m, nil

	default:
		// No navigation needed
	}

	return *m, nil
}

// handleDown handles downward navigation.
func handleDown(m *Model) (Model, tea.Cmd) {
	switch m.CurrentScreen {
	case ScreenMainMenu:
		m.MenuCursor++
		if m.MenuCursor >= len(MainMenuOptions()) {
			m.MenuCursor = 0
		}

	case ScreenModuleSelection:
		m.ModuleCursor++
		if m.ModuleCursor >= len(m.AvailableModules) {
			m.ModuleCursor = 0
		}

	case ScreenPreview:
		// Toggle confirmation
		m.ConfirmChoice = !m.ConfirmChoice

	case ScreenInstalling:
		// Can't navigate during installation
		return *m, nil

	default:
		// No navigation needed
	}

	return *m, nil
}

// handleSpace handles space bar for module selection.
func handleSpace(m *Model) (Model, tea.Cmd) {
	if m.CurrentScreen == ScreenModuleSelection && len(m.AvailableModules) > 0 {
		if m.ModuleCursor >= 0 && m.ModuleCursor < len(m.AvailableModules) {
			mod := m.AvailableModules[m.ModuleCursor]
			if m.SelectedModules[mod.Name()] {
				delete(m.SelectedModules, mod.Name())
			} else {
				m.SelectedModules[mod.Name()] = true
			}
		}
	}

	return *m, nil
}

// handleWindowSize handles window resize events.
func handleWindowSize(m *Model, msg tea.WindowSizeMsg) (Model, tea.Cmd) {
	// Update any width/height dependent state if needed
	// Currently, our views are responsive, so no action needed
	return *m, nil
}

// detectionRequestedMsg signals that detection should be run.
type detectionRequestedMsg struct{}

// runDetection is a command that runs system detection.
func runDetection() tea.Msg {
	return detectionRequestedMsg{}
}

// installCompleteMsg signals that installation is complete.
type installCompleteMsg struct{}

// startInstallationWorker starts the background installation process.
func startInstallationWorker(m *Model) tea.Cmd {
	return func() tea.Msg {
		selectedModules := m.GetSelectedModuleList()

		for i, mod := range selectedModules {
			// Simulate installation delay
			// In production, this would run actual installation

			// Create a result (simulated)
			result := types.InstallResult{
				Success: true,
				Module:  mod.Name(),
				Version: mod.Version(),
				Output:  "Successfully installed " + mod.Name(),
			}

			m.InstallResults = append(m.InstallResults, result)
			m.InstallProgress = i + 1
		}

		return installCompleteMsg{}
	}
}

// handleInstallComplete handles the installation complete message.
func handleInstallComplete(m *Model) (Model, tea.Cmd) {
	m.IsInstalling = false
	m.CurrentScreen = ScreenResults
	return *m, nil
}

// View implements tea.Model View method.
func (m Model) View() string {
	switch m.CurrentScreen {
	case ScreenWelcome:
		return WelcomeView(m)

	case ScreenMainMenu:
		return MainMenuView(m)

	case ScreenDetection:
		return DetectionView(m)

	case ScreenModuleSelection:
		return ModuleSelectionView(m)

	case ScreenPreview:
		return PreviewView(m)

	case ScreenInstalling:
		return InstallingView(m)

	case ScreenResults:
		return ResultsView(m)

	case ScreenExit:
		return ExitView()

	default:
		return WelcomeView(m)
	}
}

// GetSelectedModuleList returns a slice of selected modules.
func (m Model) GetSelectedModuleList() []modules.Module {
	var selected []modules.Module
	for _, mod := range m.AvailableModules {
		if m.SelectedModules[mod.Name()] {
			selected = append(selected, mod)
		}
	}
	return selected
}

// GetManager returns the installer manager based on detected package managers.
func (m Model) GetManager() *installer.Manager {
	if m.SystemInfo == nil {
		return nil
	}
	return installer.NewManager(m.SystemInfo)
}

// HasSelection returns true if any module is selected.
func (m Model) HasSelection() bool {
	return len(m.GetSelectedModuleList()) > 0
}

// CanInstall returns true if the wizard is ready to install.
func (m Model) CanInstall() bool {
	return m.HasSelection() && m.SystemInfo != nil && !m.IsInstalling
}

// SetScreen sets the current screen and tracks previous.
func (m *Model) SetScreen(screen Screen) {
	m.PreviousScreen = m.CurrentScreen
	m.CurrentScreen = screen
}

// GoBack returns to the previous screen.
func (m *Model) GoBack() {
	if m.PreviousScreen >= 0 {
		m.CurrentScreen = m.PreviousScreen
		m.PreviousScreen = -1
	} else {
		m.CurrentScreen = ScreenMainMenu
	}
}

// ResetSelection clears all module selections.
func (m *Model) ResetSelection() {
	m.SelectedModules = make(map[string]bool)
	m.ModuleCursor = 0
}

// StartDetection runs system detection.
func (m *Model) StartDetection() error {
	m.CurrentScreen = ScreenDetection

	info, err := detector.Detect()
	if err != nil {
		m.DetectionErr = err
		return err
	}

	m.SystemInfo = info
	m.DetectionErr = nil
	return nil
}

// StartInstallation begins the installation process.
func (m *Model) StartInstallation() {
	if !m.CanInstall() {
		return
	}

	m.IsInstalling = true
	m.InstallProgress = 0
	m.InstallTotal = len(m.GetSelectedModuleList())
	m.InstallResults = nil
	m.CurrentScreen = ScreenInstalling
}

// AddInstallResult adds a result to the installation results.
func (m *Model) AddInstallResult(result types.InstallResult) {
	m.InstallResults = append(m.InstallResults, result)
	m.InstallProgress++

	// Check if all installations are complete
	if m.InstallProgress >= m.InstallTotal {
		m.IsInstalling = false
		m.CurrentScreen = ScreenResults
	}
}

// GetSuccessCount returns the number of successful installations.
func (m Model) GetSuccessCount() int {
	count := 0
	for _, r := range m.InstallResults {
		if r.Success {
			count++
		}
	}
	return count
}
