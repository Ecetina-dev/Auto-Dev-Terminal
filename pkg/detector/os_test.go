package detector

import (
	"testing"

	"github.com/auto-dev-terminal/auto-dev-terminal/internal/types"
)

func TestDetectOS(t *testing.T) {
	// Test the function - it uses runtime.GOOS so we can only verify it returns
	// a valid OS type based on current platform
	result := DetectOS()
	
	// Verify it returns one of the expected values
	validOS := map[types.OS]bool{
		types.OSWindows: true,
		types.OSDarwin:  true,
		types.OSLinux:   true,
	}
	
	// For unknown/custom OS, it returns the raw value
	if !validOS[result] && result != "freebsd" && result != "netbsd" && result != "openbsd" {
		t.Errorf("DetectOS() returned unexpected OS: %v", result)
	}
}

func TestDetectArch(t *testing.T) {
	result := DetectArch()
	
	// Should return non-empty string
	if result == "" {
		t.Error("DetectArch() returned empty string")
	}
	
	// Should be one of common architectures
	validArch := map[string]bool{
		"amd64": true,
		"386":   true,
		"arm64": true,
		"arm":   true,
		"ppc64le": true,
		"s390x": true,
	}
	
	if !validArch[result] {
		t.Logf("DetectArch() returned: %s", result)
	}
}

func TestIsWindows(t *testing.T) {
	// Test current platform - the functions use runtime.GOOS directly
	// so we can only test that the result is consistent with the platform
	result := IsWindows()
	
	// Just verify it returns a boolean without panicking
	_ = result
}

func TestIsDarwin(t *testing.T) {
	result := IsDarwin()
	
	// Just verify it returns a boolean without panicking
	_ = result
}

func TestIsLinux(t *testing.T) {
	result := IsLinux()
	
	// Just verify it returns a boolean without panicking
	_ = result
}

// TestOSConstants verifies the OS type constants
func TestOSConstants(t *testing.T) {
	tests := []struct {
		name     string
		got      types.OS
		expected string
	}{
		{"OSWindows", types.OSWindows, "windows"},
		{"OSDarwin", types.OSDarwin, "darwin"},
		{"OSLinux", types.OSLinux, "linux"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.got) != tt.expected {
				t.Errorf("%s = %q, want %q", tt.name, tt.got, tt.expected)
			}
		})
	}
}
