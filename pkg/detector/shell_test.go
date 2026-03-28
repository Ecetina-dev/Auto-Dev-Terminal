package detector

import (
	"os"
	"testing"

	"github.com/auto-dev-terminal/auto-dev-terminal/internal/types"
)

func TestShellFromPath(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected types.Shell
	}{
		{
			name:     "zsh path",
			path:     "/bin/zsh",
			expected: types.ShellZsh,
		},
		{
			name:     "zsh with version",
			path:     "/bin/zsh-5.9",
			expected: types.ShellZsh,
		},
		{
			name:     "fish path",
			path:     "/usr/bin/fish",
			expected: types.ShellFish,
		},
		{
			name:     "bash path",
			path:     "/bin/bash",
			expected: types.ShellBash,
		},
		{
			name:     "bash with version",
			path:     "/usr/bin/bash-5.2",
			expected: types.ShellBash,
		},
		{
			name:     "pwsh path",
			path:     "/usr/bin/pwsh",
			expected: types.ShellPwsh,
		},
		{
			name:     "powershell path",
			path:     "C:\\Program Files\\PowerShell\\7\\powershell.exe",
			expected: types.ShellPwsh,
		},
		{
			name:     "cmd path",
			path:     "C:\\Windows\\System32\\cmd.exe",
			expected: types.ShellCmd,
		},
		{
			name:     "tcsh path",
			path:     "/bin/tcsh",
			expected: types.ShellTcsh,
		},
		{
			name:     "csh path",
			path:     "/bin/csh",
			expected: types.ShellCsh,
		},
		{
			name:     "ash path",
			path:     "/bin/ash",
			expected: types.ShellAsh,
		},
		{
			name:     "dash path",
			path:     "/usr/bin/dash",
			expected: types.ShellAsh,
		},
		{
			name:     "unknown path",
			path:     "/bin/unknown",
			expected: types.ShellUnknown,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := shellFromPath(tt.path)
			if result != tt.expected {
				t.Errorf("shellFromPath(%q) = %v, want %v", tt.path, result, tt.expected)
			}
		})
	}
}

func TestIsShellCompatibleWithZsh(t *testing.T) {
	tests := []struct {
		name     string
		shell    types.Shell
		expected bool
	}{
		{
			name:     "zsh is compatible",
			shell:    types.ShellZsh,
			expected: true,
		},
		{
			name:     "bash is compatible",
			shell:    types.ShellBash,
			expected: true,
		},
		{
			name:     "fish is not compatible",
			shell:    types.ShellFish,
			expected: false,
		},
		{
			name:     "powershell is not compatible",
			shell:    types.ShellPowerShell,
			expected: false,
		},
		{
			name:     "pwsh is not compatible",
			shell:    types.ShellPwsh,
			expected: false,
		},
		{
			name:     "cmd is not compatible",
			shell:    types.ShellCmd,
			expected: false,
		},
		{
			name:     "unknown is not compatible",
			shell:    types.ShellUnknown,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsShellCompatibleWithZsh(tt.shell)
			if result != tt.expected {
				t.Errorf("IsShellCompatibleWithZsh(%v) = %v, want %v", tt.shell, result, tt.expected)
			}
		})
	}
}

func TestIsShellCompatibleWithStarship(t *testing.T) {
	tests := []struct {
		name     string
		shell    types.Shell
		expected bool
	}{
		{
			name:     "bash is compatible",
			shell:    types.ShellBash,
			expected: true,
		},
		{
			name:     "zsh is compatible",
			shell:    types.ShellZsh,
			expected: true,
		},
		{
			name:     "fish is compatible",
			shell:    types.ShellFish,
			expected: true,
		},
		{
			name:     "powershell is compatible",
			shell:    types.ShellPowerShell,
			expected: true,
		},
		{
			name:     "pwsh is compatible",
			shell:    types.ShellPwsh,
			expected: true,
		},
		{
			name:     "cmd is not compatible",
			shell:    types.ShellCmd,
			expected: false,
		},
		{
			name:     "tcsh is not compatible",
			shell:    types.ShellTcsh,
			expected: false,
		},
		{
			name:     "ash is not compatible",
			shell:    types.ShellAsh,
			expected: false,
		},
		{
			name:     "unknown is not compatible",
			shell:    types.ShellUnknown,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsShellCompatibleWithStarship(tt.shell)
			if result != tt.expected {
				t.Errorf("IsShellCompatibleWithStarship(%v) = %v, want %v", tt.shell, result, tt.expected)
			}
		})
	}
}

func TestGetEnvOrEmpty(t *testing.T) {
	// Set a test environment variable
	key := "AUTO_DEV_TERMINAL_TEST_VAR"
	testValue := "test_value"

	// Clean up after test
	originalValue := os.Getenv(key)
	defer func() { os.Setenv(key, originalValue) }()

	os.Setenv(key, testValue)

	result := getEnvOrEmpty(key)
	if result != testValue {
		t.Errorf("getEnvOrEmpty(%q) = %q, want %q", key, result, testValue)
	}

	// Test non-existent variable
	nonexistent := getEnvOrEmpty("AUTO_DEV_TERMINAL_NONEXISTENT_VAR_12345")
	if nonexistent != "" {
		t.Errorf("getEnvOrEmpty(non-existent) = %q, want empty string", nonexistent)
	}
}
