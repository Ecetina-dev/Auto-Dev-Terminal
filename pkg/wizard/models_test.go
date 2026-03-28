package wizard

import (
	"testing"

	"github.com/auto-dev-terminal/auto-dev-terminal/internal/types"
	"github.com/auto-dev-terminal/auto-dev-terminal/pkg/modules"
)

func TestNewModel(t *testing.T) {
	m := NewModel()

	// Test initial state
	if m.CurrentScreen != ScreenWelcome {
		t.Errorf("CurrentScreen = %v, want %v", m.CurrentScreen, ScreenWelcome)
	}

	if m.PreviousScreen != -1 {
		t.Errorf("PreviousScreen = %v, want -1", m.PreviousScreen)
	}

	if m.SystemInfo != nil {
		t.Error("SystemInfo should be nil for new model")
	}

	if m.SelectedModules == nil {
		t.Error("SelectedModules should be initialized")
	}

	if m.ModuleCursor != 0 {
		t.Errorf("ModuleCursor = %d, want 0", m.ModuleCursor)
	}

	if m.InstallProgress != 0 {
		t.Errorf("InstallProgress = %d, want 0", m.InstallProgress)
	}

	if m.MenuCursor != 0 {
		t.Errorf("MenuCursor = %d, want 0", m.MenuCursor)
	}
}

func TestModelInit(t *testing.T) {
	m := NewModel()

	// Init should load modules
	cmd := m.Init()

	// cmd may be nil or return something - just verify it doesn't panic
	_ = cmd
}

func TestScreenConstants(t *testing.T) {
	// Verify screen constants are sequential
	tests := []struct {
		name     string
		got      Screen
		expected Screen
	}{
		{"ScreenWelcome", ScreenWelcome, 0},
		{"ScreenMainMenu", ScreenMainMenu, 1},
		{"ScreenDetection", ScreenDetection, 2},
		{"ScreenModuleSelection", ScreenModuleSelection, 3},
		{"ScreenPreview", ScreenPreview, 4},
		{"ScreenInstalling", ScreenInstalling, 5},
		{"ScreenResults", ScreenResults, 6},
		{"ScreenExit", ScreenExit, 7},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.expected {
				t.Errorf("%s = %d, want %d", tt.name, tt.got, tt.expected)
			}
		})
	}
}

func TestModelSetScreen(t *testing.T) {
	m := NewModel()

	// Set to a new screen
	m.SetScreen(ScreenMainMenu)

	if m.CurrentScreen != ScreenMainMenu {
		t.Errorf("CurrentScreen = %v, want %v", m.CurrentScreen, ScreenMainMenu)
	}

	if m.PreviousScreen != ScreenWelcome {
		t.Errorf("PreviousScreen = %v, want %v", m.PreviousScreen, ScreenWelcome)
	}

	// Set to another screen
	m.SetScreen(ScreenDetection)

	if m.CurrentScreen != ScreenDetection {
		t.Errorf("CurrentScreen = %v, want %v", m.CurrentScreen, ScreenDetection)
	}

	if m.PreviousScreen != ScreenMainMenu {
		t.Errorf("PreviousScreen = %v, want %v", m.PreviousScreen, ScreenMainMenu)
	}
}

func TestModelGoBack(t *testing.T) {
	m := NewModel()

	// GoBack with no previous should go to main menu
	m.GoBack()
	if m.CurrentScreen != ScreenMainMenu {
		t.Errorf("CurrentScreen = %v, want %v", m.CurrentScreen, ScreenMainMenu)
	}

	// Set up screen history
	m.SetScreen(ScreenModuleSelection)
	m.GoBack()
	if m.CurrentScreen != ScreenMainMenu {
		t.Errorf("CurrentScreen = %v, want %v", m.CurrentScreen, ScreenMainMenu)
	}
}

func TestModelResetSelection(t *testing.T) {
	m := NewModel()

	// Add some selections
	m.SelectedModules["module1"] = true
	m.SelectedModules["module2"] = true
	m.ModuleCursor = 5

	// Reset
	m.ResetSelection()

	if len(m.SelectedModules) != 0 {
		t.Errorf("SelectedModules length = %d, want 0", len(m.SelectedModules))
	}

	if m.ModuleCursor != 0 {
		t.Errorf("ModuleCursor = %d, want 0", m.ModuleCursor)
	}
}

func TestModelHasSelection(t *testing.T) {
	m := NewModel()

	// No selections
	if m.HasSelection() {
		t.Error("HasSelection() should return false when nothing selected")
	}

	// Setup with available modules and selection
	m.AvailableModules = modules.List()
	if len(m.AvailableModules) > 0 {
		// Select first available module
		firstMod := m.AvailableModules[0]
		m.SelectedModules[firstMod.Name()] = true

		if !m.HasSelection() {
			t.Error("HasSelection() should return true when something selected")
		}
	} else {
		// No modules registered - test with manual selection won't work
		// because HasSelection uses GetSelectedModuleList which requires AvailableModules
		t.Skip("No modules in registry to test HasSelection")
	}
}

func TestModelGetSelectedModuleList(t *testing.T) {
	m := NewModel()

	// Use the global registry's registered modules
	// The Starship module is registered in init()
	m.AvailableModules = modules.List()

	// If no modules are registered, skip this test
	if len(m.AvailableModules) == 0 {
		t.Skip("No modules registered in global registry")
	}

	// Select first available module
	firstMod := m.AvailableModules[0]
	m.SelectedModules[firstMod.Name()] = true

	selected := m.GetSelectedModuleList()

	if len(selected) != 1 {
		t.Errorf("GetSelectedModuleList() = %d, want 1", len(selected))
	}

	if len(selected) > 0 && selected[0].Name() != firstMod.Name() {
		t.Errorf("Selected module = %q, want %q", selected[0].Name(), firstMod.Name())
	}
}

func TestModelGetManager(t *testing.T) {
	m := NewModel()

	// No system info
	manager := m.GetManager()
	if manager != nil {
		t.Error("GetManager() should return nil when SystemInfo is nil")
	}

	// With system info
	m.SystemInfo = &types.SystemInfo{
		OS:     types.OSLinux,
		Distro: types.DistroUbuntu,
	}

	manager = m.GetManager()
	if manager == nil {
		t.Error("GetManager() should return manager when SystemInfo is set")
	}
}

func TestModelCanInstall(t *testing.T) {
	m := NewModel()

	// No selection, no system info
	if m.CanInstall() {
		t.Error("CanInstall() should return false with no selection and no system info")
	}

	// Setup with available modules
	m.AvailableModules = modules.List()

	// Only add selection (no system info)
	if len(m.AvailableModules) > 0 {
		firstMod := m.AvailableModules[0]
		m.SelectedModules[firstMod.Name()] = true

		if m.CanInstall() {
			t.Error("CanInstall() should return false with no system info")
		}

		// Add system info but still installing
		m.SystemInfo = &types.SystemInfo{OS: types.OSLinux}
		m.IsInstalling = true
		if m.CanInstall() {
			t.Error("CanInstall() should return false when already installing")
		}

		// Ready to install
		m.IsInstalling = false
		if !m.CanInstall() {
			t.Error("CanInstall() should return true when ready")
		}
	} else {
		// Can't fully test without modules
		t.Skip("No modules available for testing")
	}
}

func TestModelStartDetection(t *testing.T) {
	m := NewModel()

	// StartDetection should set screen and attempt detection
	// Note: This will likely fail in test environment without proper system detection
	// but we just verify it doesn't panic
	err := m.StartDetection()

	// Screen should be set to Detection regardless of detection success
	if m.CurrentScreen != ScreenDetection {
		t.Errorf("CurrentScreen = %v, want %v", m.CurrentScreen, ScreenDetection)
	}

	// If detection failed, DetectionErr should be set
	if err != nil && m.DetectionErr == nil {
		t.Error("DetectionErr should be set when detection fails")
	}
}

func TestModelStartInstallation(t *testing.T) {
	m := NewModel()

	// Can't install without prerequisites
	m.StartInstallation()
	if m.IsInstalling {
		t.Error("StartInstallation() should not start without prerequisites")
	}

	// Setup with available modules and system info
	m.AvailableModules = modules.List()

	if len(m.AvailableModules) > 0 {
		firstMod := m.AvailableModules[0]
		m.SelectedModules[firstMod.Name()] = true
		m.SystemInfo = &types.SystemInfo{OS: types.OSLinux}

		m.StartInstallation()

		if !m.IsInstalling {
			t.Error("StartInstallation() should start when ready")
		}

		if m.InstallTotal != 1 {
			t.Errorf("InstallTotal = %d, want 1", m.InstallTotal)
		}

		if m.CurrentScreen != ScreenInstalling {
			t.Errorf("CurrentScreen = %v, want %v", m.CurrentScreen, ScreenInstalling)
		}
	} else {
		t.Skip("No modules available for testing")
	}
}

func TestModelAddInstallResult(t *testing.T) {
	m := NewModel()
	m.InstallTotal = 3
	m.InstallProgress = 0

	// Add first result
	m.AddInstallResult(types.InstallResult{
		Success: true,
		Module:  "mod1",
	})

	if m.InstallProgress != 1 {
		t.Errorf("InstallProgress = %d, want 1", m.InstallProgress)
	}

	// Add more results
	m.AddInstallResult(types.InstallResult{Success: true, Module: "mod2"})
	m.AddInstallResult(types.InstallResult{Success: true, Module: "mod3"})

	// Should be complete now
	if m.IsInstalling {
		t.Error("IsInstalling should be false after all results added")
	}

	if m.CurrentScreen != ScreenResults {
		t.Errorf("CurrentScreen = %v, want %v", m.CurrentScreen, ScreenResults)
	}
}

func TestModelGetSuccessCount(t *testing.T) {
	m := NewModel()

	m.InstallResults = []types.InstallResult{
		{Success: true, Module: "mod1"},
		{Success: false, Module: "mod2"},
		{Success: true, Module: "mod3"},
		{Success: true, Module: "mod4"},
	}

	count := m.GetSuccessCount()
	if count != 3 {
		t.Errorf("GetSuccessCount() = %d, want 3", count)
	}
}

func TestModelView(t *testing.T) {
	m := NewModel()

	// Each screen should return a non-empty view
	testScreens := []Screen{
		ScreenWelcome,
		ScreenMainMenu,
		ScreenDetection,
		ScreenModuleSelection,
		ScreenPreview,
		ScreenInstalling,
		ScreenResults,
		ScreenExit,
	}

	for _, screen := range testScreens {
		m.CurrentScreen = screen
		view := m.View()
		if view == "" {
			t.Errorf("View() returned empty string for screen %v", screen)
		}
	}
}
