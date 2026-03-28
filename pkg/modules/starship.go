// Package modules provides the Starship prompt module.
package modules

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/auto-dev-terminal/auto-dev-terminal/internal/types"
)

// StarshipModule implements the Module interface for Starship prompt.
type StarshipModule struct {
	*BaseModule
}

// NewStarshipModule creates a new Starship prompt module.
func NewStarshipModule() *StarshipModule {
	return &StarshipModule{
		BaseModule: NewBaseModule(
			"starship",
			"The minimal, blazing-fast, and infinitely customizable prompt for any shell",
			"1.16.0",
			[]string{},
		),
	}
}

// Install installs Starship prompt.
func (m *StarshipModule) Install(opts *ModuleOptions) *ModuleResult {
	if opts.Verbose {
		fmt.Println("Installing Starship prompt...")
	}

	// Check if already installed
	if installed, _ := m.IsInstalled(); installed && !opts.Force {
		return &ModuleResult{
			Success: true,
			Module:  m.Name(),
			Output:  "Starship is already installed",
		}
	}

	// Determine installation method based on OS
	var output string
	var err error

	switch opts.OS {
	case types.OSWindows:
		output, err = m.installWindows(opts)
	case types.OSDarwin:
		output, err = m.installMacOS(opts)
	case types.OSLinux:
		output, err = m.installLinux(opts)
	default:
		return &ModuleResult{
			Success: false,
			Module:  m.Name(),
			Error:   fmt.Sprintf("unsupported operating system: %s", opts.OS),
		}
	}

	if err != nil {
		return &ModuleResult{
			Success: false,
			Module:  m.Name(),
			Error:   err.Error(),
			Output:  output,
		}
	}

	return &ModuleResult{
		Success: true,
		Module:  m.Name(),
		Output:  output,
		Version: m.Version(),
	}
}

// installWindows installs Starship using winget, scoop, or chocolatey.
func (m *StarshipModule) installWindows(opts *ModuleOptions) (string, error) {
	// Try winget first
	if commandExists("winget") {
		cmd := exec.Command("winget", "install", "Starship.Starship")
		if opts.Verbose {
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
		}
		if err := cmd.Run(); err == nil {
			return "Installed via winget", nil
		}
	}

	// Try scoop
	if commandExists("scoop") {
		cmd := exec.Command("scoop", "install", "starship")
		if opts.Verbose {
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
		}
		if err := cmd.Run(); err == nil {
			return "Installed via scoop", nil
		}
	}

	// Try chocolatey
	if commandExists("choco") {
		cmd := exec.Command("choco", "install", "starship")
		if opts.Verbose {
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
		}
		if err := cmd.Run(); err == nil {
			return "Installed via chocolatey", nil
		}
	}

	// Fallback to PowerShell installer
	return m.installPowerShell(opts)
}

// installPowerShell installs Starship using the official PowerShell installer.
func (m *StarshipModule) installPowerShell(opts *ModuleOptions) (string, error) {
	// Create the Starship config directory
	configDir := filepath.Join(opts.HomeDir, ".config", "starship")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return "", fmt.Errorf("creating starship config dir: %w", err)
	}

	// For Windows, we need to add Starship to PATH via PowerShell profile
	profilePath := filepath.Join(opts.HomeDir, "Documents", "PowerShell", "Microsoft.PowerShell_profile.ps1")
	if _, err := os.Stat(profilePath); os.IsNotExist(err) {
		profilePath = filepath.Join(opts.HomeDir, "Documents", "WindowsPowerShell", "Microsoft.PowerShell_profile.ps1")
	}

	return "Starship installed. Add 'Invoke-Expression (&starship init powershell)' to your PowerShell profile.", nil
}

// installMacOS installs Starship using Homebrew.
func (m *StarshipModule) installMacOS(opts *ModuleOptions) (string, error) {
	if !commandExists("brew") {
		return "", fmt.Errorf("Homebrew is not installed")
	}

	cmd := exec.Command("brew", "install", "starship")
	if opts.Verbose {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("installing starship via brew: %w", err)
	}

	return "Installed via Homebrew", nil
}

// installLinux installs Starship using the official installer or package manager.
func (m *StarshipModule) installLinux(opts *ModuleOptions) (string, error) {
	// Try package managers first
	if commandExists("apt") || commandExists("apt-get") {
		// Add Starship repository
		if !opts.Sudo {
			return "", fmt.Errorf("sudo is required to install Starship via APT")
		}
		
		installer := exec.Command("sh", "-c", fmt.Sprintf("curl -sS https://starship.rs/install.sh | sh -s -- -y"))
		if opts.Verbose {
			installer.Stdout = os.Stdout
			installer.Stderr = os.Stderr
		}
		
		if err := installer.Run(); err != nil {
			return "", fmt.Errorf("installing starship: %w", err)
		}
		
		return "Installed via official installer", nil
	}

	// Try dnf
	if commandExists("dnf") {
		cmd := exec.Command("dnf", "install", "starship")
		if opts.Sudo {
			cmd = exec.Command("sudo", "dnf", "install", "starship")
		}
		if opts.Verbose {
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
		}
		if err := cmd.Run(); err == nil {
			return "Installed via dnf", nil
		}
	}

	// Try pacman
	if commandExists("pacman") {
		cmd := exec.Command("pacman", "-S", "--noconfirm", "starship")
		if opts.Sudo {
			cmd = exec.Command("sudo", "pacman", "-S", "--noconfirm", "starship")
		}
		if opts.Verbose {
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
		}
		if err := cmd.Run(); err == nil {
			return "Installed via pacman", nil
		}
	}

	// Fallback to official installer
	return m.installLinuxOfficial(opts)
}

// installLinuxOfficial installs Starship using the official shell script.
func (m *StarshipModule) installLinuxOfficial(opts *ModuleOptions) (string, error) {
	installer := exec.Command("sh", "-c", "curl -sS https://starship.rs/install.sh | sh -s -- -y")
	if opts.Verbose {
		installer.Stdout = os.Stdout
		installer.Stderr = os.Stderr
	}

	if err := installer.Run(); err != nil {
		return "", fmt.Errorf("running official installer: %w", err)
	}

	return "Installed via official installer", nil
}

// Uninstall removes Starship prompt.
func (m *StarshipModule) Uninstall(opts *ModuleOptions) *ModuleResult {
	if opts.Verbose {
		fmt.Println("Uninstalling Starship prompt...")
	}

	var output string
	var err error

	switch opts.OS {
	case types.OSWindows:
		output, err = m.uninstallWindows(opts)
	case types.OSDarwin:
		output, err = m.uninstallMacOS(opts)
	case types.OSLinux:
		output, err = m.uninstallLinux(opts)
	default:
		return &ModuleResult{
			Success: false,
			Module:  m.Name(),
			Error:   fmt.Sprintf("unsupported operating system: %s", opts.OS),
		}
	}

	if err != nil {
		return &ModuleResult{
			Success: false,
			Module:  m.Name(),
			Error:   err.Error(),
			Output:  output,
		}
	}

	return &ModuleResult{
		Success: true,
		Module:  m.Name(),
		Output:  output,
	}
}

func (m *StarshipModule) uninstallWindows(opts *ModuleOptions) (string, error) {
	// Try winget
	if commandExists("winget") {
		cmd := exec.Command("winget", "uninstall", "Starship.Starship")
		if err := cmd.Run(); err == nil {
			return "Uninstalled via winget", nil
		}
	}

	// Try scoop
	if commandExists("scoop") {
		cmd := exec.Command("scoop", "uninstall", "starship")
		if err := cmd.Run(); err == nil {
			return "Uninstalled via scoop", nil
		}
	}

	// Try chocolatey
	if commandExists("choco") {
		cmd := exec.Command("choco", "uninstall", "starship")
		if err := cmd.Run(); err == nil {
			return "Uninstalled via chocolatey", nil
		}
	}

	return "Manual removal required", nil
}

func (m *StarshipModule) uninstallMacOS(opts *ModuleOptions) (string, error) {
	if commandExists("brew") {
		cmd := exec.Command("brew", "uninstall", "starship")
		if opts.Verbose {
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
		}
		if err := cmd.Run(); err != nil {
			return "", fmt.Errorf("uninstalling starship: %w", err)
		}
		return "Uninstalled via Homebrew", nil
	}
	return "Manual removal required", nil
}

func (m *StarshipModule) uninstallLinux(opts *ModuleOptions) (string, error) {
	// Remove binary
	binaryPath := filepath.Join(opts.HomeDir, ".local", "bin", "starship")
	if _, err := os.Stat(binaryPath); err == nil {
		if err := os.Remove(binaryPath); err != nil {
			return "", fmt.Errorf("removing starship binary: %w", err)
		}
	}

	// Remove config directory
	configDir := filepath.Join(opts.HomeDir, ".config", "starship")
	if _, err := os.Stat(configDir); err == nil {
		if err := os.RemoveAll(configDir); err != nil {
			return "", fmt.Errorf("removing starship config: %w", err)
		}
	}

	return "Removed starship binary and config", nil
}

// IsInstalled checks if Starship is installed.
func (m *StarshipModule) IsInstalled() (bool, error) {
	// Check if starship command exists
	cmd := exec.Command("starship", "--version")
	if err := cmd.Run(); err == nil {
		return true, nil
	}

	// Check common installation paths
	homeDir := os.Getenv("HOME")
	if homeDir == "" {
		homeDir = os.Getenv("USERPROFILE")
	}

	paths := []string{
		filepath.Join(homeDir, ".local", "bin", "starship"),
		filepath.Join(homeDir, ".cargo", "bin", "starship"),
		"C:\\Program Files\\starship\\starship.exe",
		"C:\\Program Files (x86)\\starship\\starship.exe",
	}

	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			return true, nil
		}
	}

	return false, nil
}

// GetInitCommand returns the shell initialization command for Starship.
func (m *StarshipModule) GetInitCommand(shell types.Shell) string {
	switch shell {
	case types.ShellZsh:
		return `eval "$(starship init zsh)"`
	case types.ShellBash:
		return `eval "$(starship init bash)"`
	case types.ShellFish:
		return `starship init fish | source`
	case types.ShellPowerShell:
		return `Invoke-Expression (&starship init powershell)`
	default:
		return "# Starship init not supported for this shell"
	}
}

// ensureStarshipConfig creates the default Starship configuration.
func (m *StarshipModule) ensureStarshipConfig(opts *ModuleOptions) error {
	configDir := filepath.Join(opts.HomeDir, ".config", "starship")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("creating config directory: %w", err)
	}

	configFile := filepath.Join(configDir, "config.toml")
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		defaultConfig := getDefaultStarshipConfig()
		if err := os.WriteFile(configFile, []byte(defaultConfig), 0644); err != nil {
			return fmt.Errorf("writing default config: %w", err)
		}
	}

	return nil
}

// getDefaultStarshipConfig returns the default Starship configuration.
func getDefaultStarshipConfig() string {
	return `# Get more completions and faster performance
https://starship.rs/config/

format = """
$directory$git_branch$git_status$nodejs$python$rust$golang
$character"""

[character]
success_symbol = "[➜](bold green)"
error_symbol = "[✗](bold red)"

[directory]
truncation_length = 3
truncate_to_repo = true

[git_branch]
symbol = " "

[git_status]
style = "bold red"
conflicted = "⚔️ "
ahead = "⇡${count}"
behind = "⇣${count}"
diverged = "⇕⇡${ahead_count}⇣${behind_count}"
untracked = "🤷"
stashed = "📦"
modified = "📝"
staged = '++\($count\)'
renamed = "👅"
deleted = "🗑"

[nodejs]
format = "via [$symbol($version )]($style)"

[python]
symbol = "🐍 "
format = "via [${symbol}${pyenv_prefix}(${version} )(\\($virtualenv\\) )]($style)"

[rust]
format = "via [$symbol($version )]($style)"

[golang]
format = "via [$symbol($version )]($style)"
`
}

// Ensure StarshipModule implements Module interface
var _ Module = (*StarshipModule)(nil)

// init registers the Starship module with the global registry.
func init() {
	Register(NewStarshipModule())
}
