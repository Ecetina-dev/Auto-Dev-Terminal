package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestConfigVariables(t *testing.T) {
	// Test NewConfigVariables creates variables
	vars, err := NewConfigVariables()
	if err != nil {
		t.Fatalf("NewConfigVariables() error = %v", err)
	}

	if vars.HomeDir == "" {
		t.Error("HomeDir should not be empty")
	}

	if vars.OS == "" {
		t.Error("OS should not be empty")
	}

	if vars.Arch == "" {
		t.Error("Arch should not be empty")
	}

	if vars.TempDir == "" {
		t.Error("TempDir should not be empty")
	}
}

func TestFilepathJoin(t *testing.T) {
	// Test basic string concatenation behavior
	tests := []struct {
		name  string
		elems []string
	}{
		{
			name:  "single element",
			elems: []string{"/home/user"},
		},
		{
			name:  "two elements",
			elems: []string{"/home/user", ".config"},
		},
		{
			name:  "three elements",
			elems: []string{"/home/user", ".config", "app"},
		},
		{
			name:  "empty elements",
			elems: []string{""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := filepathJoin(tt.elems...)
			// Just verify it doesn't panic and returns non-empty for valid input
			if len(tt.elems) > 0 && tt.elems[0] != "" && result == "" {
				t.Errorf("filepathJoin(%v) returned empty string", tt.elems)
			}
			// Verify it concatenates all elements
			for i, elem := range tt.elems {
				if elem != "" && !strings.Contains(result, elem) {
					t.Errorf("filepathJoin(%v) missing element %d: %q", tt.elems, i, elem)
				}
			}
		})
	}
}

func TestTemplateEngine(t *testing.T) {
	// Test NewTemplateEngine
	engine, err := NewTemplateEngine()
	if err != nil {
		t.Fatalf("NewTemplateEngine() error = %v", err)
	}

	if engine.variables == nil {
		t.Error("TemplateEngine.variables should not be nil")
	}

	if engine.funcMap == nil {
		t.Error("TemplateEngine.funcMap should not be nil")
	}
}

func TestTemplateEngineProcessTemplate(t *testing.T) {
	engine, err := NewTemplateEngine()
	if err != nil {
		t.Fatalf("NewTemplateEngine() error = %v", err)
	}

	tests := []struct {
		name     string
		template string
		wantErr  bool
		check    func(t *testing.T, result string)
	}{
		{
			name:     "empty template",
			template: "",
			wantErr:  false,
			check: func(t *testing.T, result string) {
				if result != "" {
					t.Errorf("Expected empty result, got %q", result)
				}
			},
		},
		{
			name:     "simple text",
			template: "Hello World",
			wantErr:  false,
			check: func(t *testing.T, result string) {
				if result != "Hello World" {
					t.Errorf("Expected 'Hello World', got %q", result)
				}
			},
		},
		{
			name:     "variable substitution",
			template: "Home: {{.HomeDir}}",
			wantErr:  false,
			check: func(t *testing.T, result string) {
				if result == "" || result == "Home: " {
					t.Errorf("Expected variable substitution, got %q", result)
				}
			},
		},
		{
			name:     "upper function",
			template: "{{ \"hello\" | upper }}",
			wantErr:  false,
			check: func(t *testing.T, result string) {
				if result != "HELLO" {
					t.Errorf("Expected 'HELLO', got %q", result)
				}
			},
		},
		{
			name:     "lower function",
			template: "{{ \"HELLO\" | lower }}",
			wantErr:  false,
			check: func(t *testing.T, result string) {
				if result != "hello" {
					t.Errorf("Expected 'hello', got %q", result)
				}
			},
		},
		{
			name:     "default function with empty",
			template: "{{ default \"fallback\" \"\" }}",
			wantErr:  false,
			check: func(t *testing.T, result string) {
				if result != "fallback" {
					t.Errorf("Expected 'fallback', got %q", result)
				}
			},
		},
		{
			name:     "default function with value",
			template: "{{ default \"fallback\" \"actual\" }}",
			wantErr:  false,
			check: func(t *testing.T, result string) {
				if result != "actual" {
					t.Errorf("Expected 'actual', got %q", result)
				}
			},
		},
		{
			name:     "invalid template",
			template: "{{ .Invalid }}}}",
			wantErr:  true,
			check:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := engine.ProcessTemplate(tt.template)
			if (err != nil) != tt.wantErr {
				t.Errorf("ProcessTemplate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.check != nil {
				tt.check(t, result)
			}
		})
	}
}

func TestTemplateEngineProcessFile(t *testing.T) {
	engine, err := NewTemplateEngine()
	if err != nil {
		t.Fatalf("NewTemplateEngine() error = %v", err)
	}

	// Test with non-existent file
	_, err = engine.ProcessFile("/nonexistent/path/file.txt")
	if err == nil {
		t.Error("ProcessFile() should fail for non-existent file")
	}
}

func TestTemplateEngineVariables(t *testing.T) {
	engine, err := NewTemplateEngine()
	if err != nil {
		t.Fatalf("NewTemplateEngine() error = %v", err)
	}

	vars := engine.Variables()
	if vars == nil {
		t.Error("Variables() returned nil")
	}
}

func TestTemplateEngineUpdateVariable(t *testing.T) {
	engine, err := NewTemplateEngine()
	if err != nil {
		t.Fatalf("NewTemplateEngine() error = %v", err)
	}

	tests := []struct {
		name      string
		key       string
		value     string
		wantErr   bool
		checkFunc func(t *testing.T, engine *TemplateEngine)
	}{
		{
			name:    "update HomeDir",
			key:     "HomeDir",
			value:   "/test/home",
			wantErr: false,
			checkFunc: func(t *testing.T, engine *TemplateEngine) {
				if engine.variables.HomeDir != "/test/home" {
					t.Errorf("HomeDir = %q, want %q", engine.variables.HomeDir, "/test/home")
				}
			},
		},
		{
			name:    "update Username",
			key:     "Username",
			value:   "testuser",
			wantErr: false,
			checkFunc: func(t *testing.T, engine *TemplateEngine) {
				if engine.variables.Username != "testuser" {
					t.Errorf("Username = %q, want %q", engine.variables.Username, "testuser")
				}
			},
		},
		{
			name:    "update OS",
			key:     "OS",
			value:   "freebsd",
			wantErr: false,
			checkFunc: func(t *testing.T, engine *TemplateEngine) {
				if engine.variables.OS != "freebsd" {
					t.Errorf("OS = %q, want %q", engine.variables.OS, "freebsd")
				}
			},
		},
		{
			name:    "unknown variable",
			key:     "UnknownVar",
			value:   "value",
			wantErr: true,
			checkFunc: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := engine.UpdateVariable(tt.key, tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateVariable() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.checkFunc != nil {
				tt.checkFunc(t, engine)
			}
		})
	}
}

func TestGetTemplateVariables(t *testing.T) {
	vars := GetTemplateVariables()

	expectedVars := []string{
		"HomeDir",
		"Shell",
		"OS",
		"Arch",
		"ConfigDir",
		"BinaryDir",
		"Username",
		"TempDir",
	}

	if len(vars) != len(expectedVars) {
		t.Errorf("GetTemplateVariables() returned %d vars, want %d", len(vars), len(expectedVars))
	}

	for _, expected := range expectedVars {
		found := false
		for _, v := range vars {
			if v == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected variable %q not found", expected)
		}
	}
}

// Helper to create temp directory for tests
func createTempDir(t *testing.T) string {
	dir, err := os.MkdirTemp("", "template-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	return dir
}

// Helper to create temp file
func createTempFile(t *testing.T, dir, name, content string) string {
	path := filepath.Join(dir, name)
	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	return path
}
