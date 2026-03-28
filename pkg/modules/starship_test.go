package modules

import (
	"testing"

	"github.com/auto-dev-terminal/auto-dev-terminal/internal/types"
)

func TestNewStarshipModule(t *testing.T) {
	m := NewStarshipModule()
	if m == nil {
		t.Fatal("NewStarshipModule() returned nil")
	}

	if m.Name() != "starship" {
		t.Errorf("Name() = %q, want %q", m.Name(), "starship")
	}

	if m.Description() == "" {
		t.Error("Description() should not be empty")
	}

	if m.Version() == "" {
		t.Error("Version() should not be empty")
	}

	if m.GetDependencies() == nil {
		t.Error("GetDependencies() should not be nil")
	}
}

func TestStarshipModuleInstall(t *testing.T) {
	m := NewStarshipModule()

	// Test OS detection logic - we can't actually install in test environment
	// so we just verify the OS switch logic works correctly

	tests := []struct {
		name          string
		os            types.OS
		shouldSucceed bool
	}{
		{
			name:          "windows OS",
			os:            types.OSWindows,
			shouldSucceed: true,
		},
		{
			name:          "darwin OS",
			os:            types.OSDarwin,
			shouldSucceed: true,
		},
		{
			name:          "linux OS",
			os:            types.OSLinux,
			shouldSucceed: true,
		},
		{
			name:          "unsupported OS",
			os:            types.OS("freebsd"),
			shouldSucceed: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := &ModuleOptions{
				OS:      tt.os,
				HomeDir: "/home/user",
				Verbose: false,
				Force:   true, // Force to bypass "already installed" check
			}

			result := m.Install(opts)
			if result == nil {
				t.Fatal("Install() returned nil result")
			}

			// For unsupported OS, should always fail
			if !tt.shouldSucceed {
				if result.Success {
					t.Error("Install() should have failed for unsupported OS")
				}
				return
			}

			// For supported OS, we can't guarantee success in test environment
			// (Homebrew might not be installed, network might be down, etc.)
			// So we just verify it didn't panic and returned a valid result
			if result.Module != "starship" {
				t.Errorf("Install() returned wrong module: %s", result.Module)
			}
		})
	}
}

func TestStarshipModuleUninstall(t *testing.T) {
	m := NewStarshipModule()

	tests := []struct {
		name string
		opts *ModuleOptions
	}{
		{
			name: "windows OS",
			opts: &ModuleOptions{
				OS:      types.OSWindows,
				HomeDir: "/home/user",
				Verbose: false,
			},
		},
		{
			name: "darwin OS",
			opts: &ModuleOptions{
				OS:      types.OSDarwin,
				HomeDir: "/home/user",
				Verbose: false,
			},
		},
		{
			name: "linux OS",
			opts: &ModuleOptions{
				OS:      types.OSLinux,
				HomeDir: "/home/user",
				Verbose: false,
			},
		},
		{
			name: "unsupported OS",
			opts: &ModuleOptions{
				OS:      "freebsd",
				HomeDir: "/home/user",
				Verbose: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := m.Uninstall(tt.opts)
			if result == nil {
				t.Fatal("Uninstall() returned nil result")
			}

			// For unsupported OS, expect failure
			if tt.opts != nil && tt.opts.OS == "freebsd" {
				if result.Success {
					t.Error("Uninstall() should fail for unsupported OS")
				}
			}
		})
	}
}

func TestStarshipModuleIsInstalled(t *testing.T) {
	m := NewStarshipModule()

	// This test checks if starship is installed on the system
	// In a CI environment, it likely won't be installed
	installed, err := m.IsInstalled()

	// The function should not error - it returns false if not found
	if err != nil {
		t.Logf("IsInstalled() returned error (expected if starship not installed): %v", err)
	}

	// Just verify we got a result
	_ = installed
}

func TestStarshipModuleGetInitCommand(t *testing.T) {
	m := NewStarshipModule()

	tests := []struct {
		name     string
		shell    types.Shell
		expected string
	}{
		{
			name:     "zsh",
			shell:    types.ShellZsh,
			expected: `eval "$(starship init zsh)"`,
		},
		{
			name:     "bash",
			shell:    types.ShellBash,
			expected: `eval "$(starship init bash)"`,
		},
		{
			name:     "fish",
			shell:    types.ShellFish,
			expected: `starship init fish | source`,
		},
		{
			name:     "powershell",
			shell:    types.ShellPowerShell,
			expected: `Invoke-Expression (&starship init powershell)`,
		},
		// Note: ShellPwsh is not handled in GetInitCommand - falls through to default
		{
			name:     "cmd - unsupported",
			shell:    types.ShellCmd,
			expected: "# Starship init not supported for this shell",
		},
		{
			name:     "tcsh - unsupported",
			shell:    types.ShellTcsh,
			expected: "# Starship init not supported for this shell",
		},
		{
			name:     "unknown shell",
			shell:    types.ShellUnknown,
			expected: "# Starship init not supported for this shell",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := m.GetInitCommand(tt.shell)
			if result != tt.expected {
				t.Errorf("GetInitCommand(%v) = %q, want %q", tt.shell, result, tt.expected)
			}
		})
	}
}

func TestGetDefaultStarshipConfig(t *testing.T) {
	config := getDefaultStarshipConfig()

	if config == "" {
		t.Error("getDefaultStarshipConfig() returned empty string")
	}

	// Check for expected sections
	expectedSections := []string{
		"format =",
		"[character]",
		"[directory]",
		"[git_branch]",
		"[git_status]",
		"[nodejs]",
		"[python]",
		"[rust]",
		"[golang]",
	}

	for _, section := range expectedSections {
		found := false
		for i := 0; i <= len(config)-len(section); i++ {
			if config[i:i+len(section)] == section {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Config missing expected section: %s", section)
		}
	}
}

func TestStarshipModuleImplementsModuleInterface(t *testing.T) {
	// Verify StarshipModule implements Module interface
	var _ Module = (*StarshipModule)(nil)
}
