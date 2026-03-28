// Package modules provides the Oh My Zsh module.
package modules

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/auto-dev-terminal/auto-dev-terminal/internal/types"
)

// OhMyZshModule implements the Module interface for Oh My Zsh.
type OhMyZshModule struct {
	*BaseModule
}

// NewOhMyZshModule creates a new Oh My Zsh module.
func NewOhMyZshModule() *OhMyZshModule {
	return &OhMyZshModule{
		BaseModule: NewBaseModule(
			"ohmyzsh",
			"A delightful, open source, community-driven framework for managing your Zsh configuration",
			"latest",
			[]string{},
		),
	}
}

// Install installs Oh My Zsh.
func (m *OhMyZshModule) Install(opts *ModuleOptions) *ModuleResult {
	if opts.Verbose {
		fmt.Println("Installing Oh My Zsh...")
	}

	// Oh My Zsh only supports Unix-like systems
	if opts.OS == types.OSWindows {
		return &ModuleResult{
			Success: false,
			Module:  m.Name(),
			Error:   "Oh My Zsh is not supported on Windows. Use WSL or MSYS2.",
		}
	}

	// Check if already installed
	if installed, _ := m.IsInstalled(); installed && !opts.Force {
		return &ModuleResult{
			Success: true,
			Module:  m.Name(),
			Output:  "Oh My Zsh is already installed",
		}
	}

	// Check if Zsh is installed
	if !commandExists("zsh") {
		return &ModuleResult{
			Success: false,
			Module:  m.Name(),
			Error:   "Zsh is not installed. Please install Zsh first.",
		}
	}

	// Set environment variable to skip interactive prompts
	env := append(os.Environ(), "RUNZSH=no", "CHSH=no")

	// Run the official installer
	installer := exec.Command("sh", "-c", `sh -c "$(curl -fsSL https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh)" "" --unattended`)
	installer.Env = env
	if opts.Verbose {
		installer.Stdout = os.Stdout
		installer.Stderr = os.Stderr
	}

	if err := installer.Run(); err != nil {
		return &ModuleResult{
			Success: false,
			Module:  m.Name(),
			Error:   fmt.Sprintf("failed to install Oh My Zsh: %v", err),
		}
	}

	// Set Zsh as the default shell
	if err := m.setDefaultShell(opts); err != nil {
		return &ModuleResult{
			Success: true,
			Module:  m.Name(),
			Output:  "Oh My Zsh installed, but failed to set as default shell: " + err.Error(),
		}
	}

	return &ModuleResult{
		Success: true,
		Module:  m.Name(),
		Output:  "Oh My Zsh installed successfully",
		Version: m.getInstalledVersion(),
	}
}

// setDefaultShell sets Zsh as the default shell.
func (m *OhMyZshModule) setDefaultShell(opts *ModuleOptions) error {
	chshPath, err := exec.LookPath("chsh")
	if err != nil {
		return fmt.Errorf("chsh not found: %w", err)
	}

	// Find Zsh path
	zshPath, err := exec.LookPath("zsh")
	if err != nil {
		return fmt.Errorf("zsh not found: %w", err)
	}

	var cmd *exec.Cmd
	if opts.Sudo {
		cmd = exec.Command("sudo", chshPath, "-s", zshPath)
	} else {
		cmd = exec.Command(chshPath, "-s", zshPath)
	}

	if opts.Verbose {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	return cmd.Run()
}

// Uninstall removes Oh My Zsh.
func (m *OhMyZshModule) Uninstall(opts *ModuleOptions) *ModuleResult {
	if opts.Verbose {
		fmt.Println("Uninstalling Oh My Zsh...")
	}

	if opts.OS == types.OSWindows {
		return &ModuleResult{
			Success: false,
			Module:  m.Name(),
			Error:   "Oh My Zsh is not supported on Windows",
		}
	}

	// Remove the Oh My Zsh directory
	ohMyZshDir := filepath.Join(opts.HomeDir, ".oh-my-zsh")
	if _, err := os.Stat(ohMyZshDir); err == nil {
		if err := os.RemoveAll(ohMyZshDir); err != nil {
			return &ModuleResult{
				Success: false,
				Module:  m.Name(),
				Error:   fmt.Sprintf("failed to remove .oh-my-zsh: %v", err),
			}
		}
	}

	// Restore the original .zshrc if there's a backup
	backupPath := filepath.Join(opts.HomeDir, ".zshrc.bak")
	if _, err := os.Stat(backupPath); err == nil {
		originalPath := filepath.Join(opts.HomeDir, ".zshrc")
		if err := os.Rename(backupPath, originalPath); err != nil {
			return &ModuleResult{
				Success: false,
				Module:  m.Name(),
				Error:   fmt.Sprintf("failed to restore .zshrc: %v", err),
			}
		}
	}

	return &ModuleResult{
		Success: true,
		Module:  m.Name(),
		Output:  "Oh My Zsh uninstalled successfully",
	}
}

// IsInstalled checks if Oh My Zsh is installed.
func (m *OhMyZshModule) IsInstalled() (bool, error) {
	ohMyZshDir := filepath.Join(os.Getenv("HOME"), ".oh-my-zsh")
	if _, err := os.Stat(ohMyZshDir); err == nil {
		return true, nil
	}
	return false, nil
}

// getInstalledVersion returns the installed version of Oh My Zsh.
func (m *OhMyZshModule) getInstalledVersion() string {
	versionFile := filepath.Join(os.Getenv("HOME"), ".oh-my-zsh", "VERSION")
	data, err := os.ReadFile(versionFile)
	if err != nil {
		return "unknown"
	}
	return strings.TrimSpace(string(data))
}

// InstallPlugin installs an Oh My Zsh plugin.
func (m *OhMyZshModule) InstallPlugin(pluginName string, opts *ModuleOptions) *ModuleResult {
	pluginDir := filepath.Join(opts.HomeDir, ".oh-my-zsh", "custom", "plugins", pluginName)

	if _, err := os.Stat(pluginDir); err == nil {
		return &ModuleResult{
			Success: true,
			Module:  m.Name(),
			Output:  fmt.Sprintf("Plugin %s is already installed", pluginName),
		}
	}

	// Try to clone the plugin repository
	pluginURL := fmt.Sprintf("https://github.com/%s.git", pluginName)
	cmd := exec.Command("git", "clone", "--depth", "1", pluginURL, pluginDir)
	if opts.Verbose {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	if err := cmd.Run(); err != nil {
		return &ModuleResult{
			Success: false,
			Module:  m.Name(),
			Error:   fmt.Sprintf("failed to install plugin %s: %v", pluginName, err),
		}
	}

	return &ModuleResult{
		Success: true,
		Module:  m.Name(),
		Output:  fmt.Sprintf("Plugin %s installed successfully", pluginName),
	}
}

// InstallTheme installs an Oh My Zsh theme.
func (m *OhMyZshModule) InstallTheme(themeName string, opts *ModuleOptions) *ModuleResult {
	themesDir := filepath.Join(opts.HomeDir, ".oh-my-zsh", "themes")

	// Check if theme file exists
	themePath := filepath.Join(themesDir, themeName+".zsh-theme")
	if _, err := os.Stat(themePath); err == nil {
		return &ModuleResult{
			Success: true,
			Module:  m.Name(),
			Output:  fmt.Sprintf("Theme %s is already installed", themeName),
		}
	}

	// Try to clone the theme repository
	themeURL := fmt.Sprintf("https://github.com/%s.git", themeName)
	cmd := exec.Command("git", "clone", "--depth", "1", themeURL, themesDir)
	if opts.Verbose {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	if err := cmd.Run(); err != nil {
		return &ModuleResult{
			Success: false,
			Module:  m.Name(),
			Error:   fmt.Sprintf("failed to install theme %s: %v", themeName, err),
		}
	}

	return &ModuleResult{
		Success: true,
		Module:  m.Name(),
		Output:  fmt.Sprintf("Theme %s installed successfully", themeName),
	}
}

// GetPlugins returns a list of recommended plugins.
func (m *OhMyZshModule) GetPlugins() []string {
	return []string{
		"zsh-autosuggestions",
		"zsh-syntax-highlighting",
		"zsh-completions",
		"z",
		"git",
		"docker",
		"kubectl",
	}
}

// Ensure OhMyZshModule implements Module interface
var _ Module = (*OhMyZshModule)(nil)

// init registers the Oh My Zsh module with the global registry.
func init() {
	_ = Register(NewOhMyZshModule())
}
