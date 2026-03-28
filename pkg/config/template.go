package config

import (
	"bytes"
	"fmt"
	"os"
	"runtime"
	"strings"
	"text/template"
)

// ConfigVariables holds the template variables available for config templating
type ConfigVariables struct {
	HomeDir   string
	Shell     string
	OS        string
	Arch      string
	ConfigDir string
	BinaryDir string
	Username  string
	TempDir   string
}

// NewConfigVariables creates a new ConfigVariables with system values
func NewConfigVariables() (*ConfigVariables, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	shell := os.Getenv("SHELL")
	if shell == "" {
		// Default shells based on OS
		switch runtime.GOOS {
		case "windows":
			shell = "powershell.exe"
		case "darwin":
			shell = "/bin/zsh"
		default:
			shell = "/bin/bash"
		}
	}

	username := os.Getenv("USER")
	if username == "" {
		username = os.Getenv("USERNAME")
	}
	if username == "" {
		username = "unknown"
	}

	configDir := filepathJoin(homeDir, ".auto-dev-terminal")
	binaryDir := filepathJoin(homeDir, ".local", "bin")
	tempDir := os.TempDir()

	return &ConfigVariables{
		HomeDir:   homeDir,
		Shell:     shell,
		OS:        runtime.GOOS,
		Arch:      runtime.GOARCH,
		ConfigDir: configDir,
		BinaryDir: binaryDir,
		Username:  username,
		TempDir:   tempDir,
	}, nil
}

// filepathJoin is a helper to handle path joining across platforms
func filepathJoin(elem ...string) string {
	result := elem[0]
	for _, e := range elem[1:] {
		if !strings.HasSuffix(result, string(os.PathSeparator)) && !strings.HasSuffix(result, "/") && !strings.HasSuffix(result, "\\") {
			result += string(os.PathSeparator)
		}
		result += e
	}
	return result
}

// TemplateEngine processes config templates with variable substitution
type TemplateEngine struct {
	variables *ConfigVariables
	funcMap   template.FuncMap
}

// NewTemplateEngine creates a new TemplateEngine with default variables
func NewTemplateEngine() (*TemplateEngine, error) {
	vars, err := NewConfigVariables()
	if err != nil {
		return nil, err
	}

	funcMap := template.FuncMap{
		"env": os.Getenv,
		"default": func(defaultVal, actual string) string {
			if actual == "" {
				return defaultVal
			}
			return actual
		},
		"upper": strings.ToUpper,
		"lower": strings.ToLower,
	}

	return &TemplateEngine{
		variables: vars,
		funcMap:   funcMap,
	}, nil
}

// ProcessTemplate processes template content and returns the result
func (te *TemplateEngine) ProcessTemplate(content string) (string, error) {
	tmpl, err := template.New("config").Funcs(te.funcMap).Parse(content)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, te.variables); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

// ProcessFile reads a file, processes its template content, and returns the result
func (te *TemplateEngine) ProcessFile(templatePath string) (string, error) {
	content, err := os.ReadFile(templatePath)
	if err != nil {
		return "", fmt.Errorf("failed to read template file: %w", err)
	}

	return te.ProcessTemplate(string(content))
}

// Variables returns the current config variables
func (te *TemplateEngine) Variables() *ConfigVariables {
	return te.variables
}

// UpdateVariable updates a specific template variable
func (te *TemplateEngine) UpdateVariable(key, value string) error {
	switch key {
	case "HomeDir":
		te.variables.HomeDir = value
	case "Shell":
		te.variables.Shell = value
	case "OS":
		te.variables.OS = value
	case "Arch":
		te.variables.Arch = value
	case "ConfigDir":
		te.variables.ConfigDir = value
	case "BinaryDir":
		te.variables.BinaryDir = value
	case "Username":
		te.variables.Username = value
	case "TempDir":
		te.variables.TempDir = value
	default:
		return fmt.Errorf("unknown variable: %s", key)
	}
	return nil
}

// GetTemplateVariables returns available template variable names
func GetTemplateVariables() []string {
	return []string{
		"HomeDir",
		"Shell",
		"OS",
		"Arch",
		"ConfigDir",
		"BinaryDir",
		"Username",
		"TempDir",
	}
}
