// Package modules provides the Nerd Fonts module.
package modules

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/auto-dev-terminal/auto-dev-terminal/internal/types"
)

// FontsModule implements the Module interface for Nerd Fonts.
type FontsModule struct {
	*BaseModule
}

// NewFontsModule creates a new Nerd Fonts module.
func NewFontsModule() *FontsModule {
	return &FontsModule{
		BaseModule: NewBaseModule(
			"fonts",
			"Nerd Fonts - Iconic font aggregator, collection, and patcher",
			"3.2.1",
			[]string{},
		),
	}
}

// Install installs Nerd Fonts.
func (m *FontsModule) Install(opts *ModuleOptions) *ModuleResult {
	if opts.Verbose {
		fmt.Println("Installing Nerd Fonts...")
	}

	// Check if already installed (check for a common font)
	if installed, _ := m.IsInstalled(); installed && !opts.Force {
		return &ModuleResult{
			Success: true,
			Module:  m.Name(),
			Output:  "Nerd Fonts appear to be already installed",
		}
	}

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

// installWindows installs Nerd Fonts on Windows.
func (m *FontsModule) installWindows(opts *ModuleOptions) (string, error) {
	// Create fonts directory
	fontsDir := filepath.Join(os.Getenv("LOCALAPPDATA"), "Microsoft", "Windows", "Fonts")
	if err := os.MkdirAll(fontsDir, 0755); err != nil {
		return "", fmt.Errorf("creating fonts directory: %w", err)
	}

	// Download and install JetBrains Mono Nerd Font
	fontName := "JetBrainsMono"
	fontURL := "https://github.com/ryanoasis/nerd-fonts/releases/download/v3.2.1/JetBrainsMono.zip"

	if opts.Verbose {
		fmt.Printf("Downloading %s from %s\n", fontName, fontURL)
	}

	// Download the font zip
	tmpDir, err := os.MkdirTemp("", "nerd-fonts-")
	if err != nil {
		return "", fmt.Errorf("creating temp directory: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	zipPath := filepath.Join(tmpDir, fontName+".zip")
	if err := m.downloadFile(fontURL, zipPath); err != nil {
		return "", fmt.Errorf("downloading font: %w", err)
	}

	// Extract and install fonts
	if err := m.extractAndInstallFonts(zipPath, fontsDir); err != nil {
		return "", fmt.Errorf("installing fonts: %w", err)
	}

	// Register fonts in Windows
	if err := m.registerWindowsFonts(); err != nil {
		return "", fmt.Errorf("registering fonts: %w", err)
	}

	return fmt.Sprintf("Installed %s Nerd Fonts to %s", fontName, fontsDir), nil
}

// installMacOS installs Nerd Fonts on macOS.
func (m *FontsModule) installMacOS(opts *ModuleOptions) (string, error) {
	// Check if Homebrew is available
	if !commandExists("brew") {
		return "", fmt.Errorf("Homebrew is required to install Nerd Fonts on macOS")
	}

	// Install font cask
	cmd := exec.Command("brew", "install", "--cask", "font-jetbrains-mono-nerd-font")
	if opts.Verbose {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("installing font via Homebrew: %w", err)
	}

	return "Installed JetBrains Mono Nerd Font via Homebrew", nil
}

// installLinux installs Nerd Fonts on Linux.
func (m *FontsModule) installLinux(opts *ModuleOptions) (string, error) {
	// Try various package managers
	if commandExists("apt") || commandExists("apt-get") {
		return m.installLinuxAPT(opts)
	}

	if commandExists("dnf") {
		return m.installLinuxDNF(opts)
	}

	if commandExists("pacman") {
		return m.installLinuxPacman(opts)
	}

	// Fallback: manual installation
	return m.installLinuxManual(opts)
}

// installLinuxAPT installs Nerd Fonts using APT on Debian/Ubuntu.
func (m *FontsModule) installLinuxAPT(opts *ModuleOptions) (string, error) {
	// Try to install from package if available
	cmd := exec.Command("apt", "install", "-y", "fonts-jetbrains-mono")
	if opts.Sudo {
		cmd = exec.Command("sudo", "apt", "install", "-y", "fonts-jetbrains-mono")
	}
	if opts.Verbose {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	if err := cmd.Run(); err == nil {
		return "Installed fonts via APT", nil
	}

	// Fallback to manual installation
	return m.installLinuxManual(opts)
}

// installLinuxDNF installs Nerd Fonts using DNF on Fedora/RHEL.
func (m *FontsModule) installLinuxDNF(opts *ModuleOptions) (string, error) {
	cmd := exec.Command("dnf", "install", "-y", "fontconfig", "mkfontscale")
	if opts.Sudo {
		cmd = exec.Command("sudo", "dnf", "install", "-y", "fontconfig", "mkfontscale")
	}
	if opts.Verbose {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("installing font dependencies: %w", err)
	}

	return m.installLinuxManual(opts)
}

// installLinuxPacman installs Nerd Fonts using Pacman on Arch Linux.
func (m *FontsModule) installLinuxPacman(opts *ModuleOptions) (string, error) {
	// Try to find and install from AUR or official repos
	// First check if there's a package
	cmd := exec.Command("pacman", "-Ss", "nerd-fonts")
	output, err := cmd.Output()
	if err == nil && strings.Contains(string(output), "nerd-fonts") {
		installCmd := exec.Command("pacman", "-S", "--noconfirm", "nerd-fonts-jetbrains-mono")
		if opts.Sudo {
			installCmd = exec.Command("sudo", "pacman", "-S", "--noconfirm", "nerd-fonts-jetbrains-mono")
		}
		if opts.Verbose {
			installCmd.Stdout = os.Stdout
			installCmd.Stderr = os.Stderr
		}
		if err := installCmd.Run(); err == nil {
			return "Installed via Pacman", nil
		}
	}

	return m.installLinuxManual(opts)
}

// installLinuxManual manually installs Nerd Fonts on Linux.
func (m *FontsModule) installLinuxManual(opts *ModuleOptions) (string, error) {
	// Create user font directory
	fontDir := filepath.Join(opts.HomeDir, ".local", "share", "fonts")
	if err := os.MkdirAll(fontDir, 0755); err != nil {
		return "", fmt.Errorf("creating font directory: %w", err)
	}

	// Download the font
	fontURL := "https://github.com/ryanoasis/nerd-fonts/releases/download/v3.2.1/JetBrainsMono.tar.xz"
	fontFile := filepath.Join(fontDir, "JetBrainsMono.tar.xz")

	if opts.Verbose {
		fmt.Printf("Downloading %s\n", fontURL)
	}

	if err := m.downloadFile(fontURL, fontFile); err != nil {
		return "", fmt.Errorf("downloading font: %w", err)
	}

	// Extract
	extractCmd := exec.Command("tar", "-xf", fontFile, "-C", fontDir)
	if err := extractCmd.Run(); err != nil {
		return "", fmt.Errorf("extracting font: %w", err)
	}

	// Clean up
	os.Remove(fontFile)

	// Refresh font cache
	if err := m.refreshFontCache(); err != nil {
		return "", fmt.Errorf("refreshing font cache: %w", err)
	}

	return fmt.Sprintf("Installed JetBrains Mono Nerd Font to %s", fontDir), nil
}

// downloadFile downloads a file from URL to path.
func (m *FontsModule) downloadFile(url, path string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	out, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("creating file: %w", err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

// extractAndInstallFonts extracts font zip and installs to the specified directory.
func (m *FontsModule) extractAndInstallFonts(zipPath, destDir string) error {
	// Determine the archive format and extract accordingly
	var cmd *exec.Cmd

	if strings.HasSuffix(zipPath, ".zip") {
		cmd = exec.Command("powershell", "-Command",
			fmt.Sprintf("Expand-Archive -Path '%s' -DestinationPath '%s' -Force", zipPath, destDir))
	} else if strings.HasSuffix(zipPath, ".tar.xz") {
		tmpDir := filepath.Dir(zipPath)
		cmd = exec.Command("tar", "-xf", zipPath, "-C", tmpDir)
	} else {
		return fmt.Errorf("unsupported archive format")
	}

	return cmd.Run()
}

// registerWindowsFonts registers the installed fonts with Windows.
func (m *FontsModule) registerWindowsFonts() error {
	fontsDir := filepath.Join(os.Getenv("LOCALAPPDATA"), "Microsoft", "Windows", "Fonts")

	// Add fonts to registry using PowerShell
	psScript := fmt.Sprintf(`
$fontsDir = '%s'
Get-ChildItem -Path $fontsDir -Filter '*.ttf' | ForEach-Object {
    $fontName = $_.Name
    $fontPath = $_.FullName
    $regPath = "HKCU:\Software\Microsoft\Windows NT\CurrentVersion\Fonts"
    Set-ItemProperty -Path $regPath -Name "$fontName (TrueType)" -Value $fontPath
}
`, fontsDir)

	cmd := exec.Command("powershell", "-Command", psScript)
	return cmd.Run()
}

// refreshFontCache refreshes the font cache on Linux.
func (m *FontsModule) refreshFontCache() error {
	// Run fc-cache to refresh font cache
	cmd := exec.Command("fc-cache", "-f", "-v")
	if runtime.GOOS != "windows" {
		// Try to run as user first, then with fc-cache in user directory
		cmd = exec.Command("fc-cache", "-f", "-v", filepath.Join(os.Getenv("HOME"), ".local", "share", "fonts"))
	}
	return cmd.Run()
}

// Uninstall removes Nerd Fonts.
func (m *FontsModule) Uninstall(opts *ModuleOptions) *ModuleResult {
	if opts.Verbose {
		fmt.Println("Uninstalling Nerd Fonts...")
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

func (m *FontsModule) uninstallWindows(opts *ModuleOptions) (string, error) {
	fontsDir := filepath.Join(os.Getenv("LOCALAPPDATA"), "Microsoft", "Windows", "Fonts")

	// Remove font files
	if err := os.RemoveAll(fontsDir); err != nil {
		return "", fmt.Errorf("removing fonts directory: %w", err)
	}

	// Remove registry entries
	psScript := `
$regPath = "HKCU:\Software\Microsoft\Windows NT\CurrentVersion\Fonts"
Get-ItemProperty -Path $regPath | Get-Member -MemberType NoteProperty | 
    Where-Object { $_.Name -match "Nerd Font|JetBrains Mono" } | 
    ForEach-Object { Remove-ItemProperty -Path $regPath -Name $_.Name }
`
	cmd := exec.Command("powershell", "-Command", psScript)
	cmd.Run()

	return "Removed Nerd Fonts from Windows", nil
}

func (m *FontsModule) uninstallMacOS(opts *ModuleOptions) (string, error) {
	if !commandExists("brew") {
		return "", fmt.Errorf("Homebrew not found")
	}

	cmd := exec.Command("brew", "uninstall", "--cask", "font-jetbrains-mono-nerd-font")
	if opts.Verbose {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("uninstalling font: %w", err)
	}

	return "Uninstalled JetBrains Mono Nerd Font via Homebrew", nil
}

func (m *FontsModule) uninstallLinux(opts *ModuleOptions) (string, error) {
	fontDir := filepath.Join(opts.HomeDir, ".local", "share", "fonts")

	if err := os.RemoveAll(fontDir); err != nil {
		return "", fmt.Errorf("removing font directory: %w", err)
	}

	if err := m.refreshFontCache(); err != nil {
		return "", fmt.Errorf("refreshing font cache: %w", err)
	}

	return "Removed Nerd Fonts from user font directory", nil
}

// IsInstalled checks if Nerd Fonts are installed.
func (m *FontsModule) IsInstalled() (bool, error) {
	// Check for common Nerd Font files
	fontsToCheck := []string{
		"JetBrains Mono Regular Nerd Font Complete.ttf",
		"JetBrains Mono Regular Nerd Font Complete Mono.ttf",
		"JetBrainsMonoNerdFont-Regular.ttf",
	}

	// Check user font directories
	searchDirs := []string{
		filepath.Join(os.Getenv("HOME"), ".local", "share", "fonts"),
		filepath.Join(os.Getenv("HOME"), ".fonts"),
		filepath.Join(os.Getenv("LOCALAPPDATA"), "Microsoft", "Windows", "Fonts"),
		"/usr/share/fonts",
		"/usr/local/share/fonts",
	}

	for _, dir := range searchDirs {
		if _, err := os.Stat(dir); err != nil {
			continue
		}

		for _, font := range fontsToCheck {
			fontPath := filepath.Join(dir, font)
			if _, err := os.Stat(fontPath); err == nil {
				return true, nil
			}
		}
	}

	return false, nil
}

// GetRecommendedFonts returns a list of recommended Nerd Fonts.
func (m *FontsModule) GetRecommendedFonts() []string {
	return []string{
		"JetBrains Mono",
		"FiraCode",
		"Source Code Pro",
		"Cascadia Code",
		"Hack",
		"DejaVu Sans Mono",
	}
}

// Ensure FontsModule implements Module interface
var _ Module = (*FontsModule)(nil)

// init registers the Fonts module with the global registry.
func init() {
	Register(NewFontsModule())
}
