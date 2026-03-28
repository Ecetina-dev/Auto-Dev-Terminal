# Auto-Dev-Terminal

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.24+-00ADD8?style=for-the-badge&logo=go" alt="Go Version">
  <img src="https://img.shields.io/badge/platforms-Windows%20%7C%20macOS%20%7C%20Linux-blue?style=for-the-badge&platforms=Windows%20%7C%20macOS%20%7C%20Linux" alt="Platforms">
  <img src="https://img.shields.io/github/actions/workflow/status/Ecetina-dev/Auto-Dev-Terminal/ci.yml?style=for-the-badge" alt="CI">
  <img src="https://img.shields.io/github/license/Ecetina-dev/Auto-Dev-Terminal?style=for-the-badge" alt="License">
</p>

> **Automate your development environment setup in minutes, not hours.**

Auto-Dev-Terminal is a CLI tool that automatically detects your operating system, shell, and package managers, then helps you install and configure essential development tools through an interactive wizard or CLI commands.

## 🚀 Features

- **Intelligent Detection**: Automatically detects OS (Windows/macOS/Linux), Linux distribution, shell (bash/zsh/fish/PowerShell), and available package managers
- **Interactive Wizard**: User-friendly TUI for selecting and installing modules
- **Module System**: Pluggable architecture with built-in modules
- **Backup & Restore**: Automatic backup before modifying configuration files
- **Cross-Platform**: Works on Windows, macOS, and Linux

### Built-in Modules

| Module | Description |
|--------|-------------|
| **Starship** | Blazing-fast, cross-shell prompt |
| **Oh My Zsh** | Community-driven Zsh framework |
| **Git Config** | Enhanced Git configuration with aliases |
| **Fonts** | Nerd Fonts for iconic glyphs |

## 📦 Installation

### From Source

#### Linux / macOS

```bash
git clone https://github.com/Ecetina-dev/Auto-Dev-Terminal.git
cd Auto-Dev-Terminal
go build -o adt ./cmd/cli
sudo mv adt /usr/local/bin/
# Verify installation
adt --version
```

#### Windows

```powershell
# Using PowerShell
git clone https://github.com/Ecetina-dev/Auto-Dev-Terminal.git
cd Auto-Dev-Terminal
go build -o adt.exe ./cmd/cli

# Add to PATH (run as Administrator)
# Option 1: Move to a folder in PATH
Move-Item adt.exe C:\Windows\System32\

# Option 2: Add to user PATH
$env:PATH += ";$PWD"
# Or permanently:
[Environment]::SetEnvironmentVariable("PATH", $env:PATH + ";$PWD", "User")

# Verify installation
.\adt.exe --version
```

### Pre-built Binaries

Download the latest release from the [releases page](https://github.com/Ecetina-dev/Auto-Dev-Terminal/releases).

#### Windows (using Winget)

```powershell
winget install Ecetina-dev.AutoDevTerminal
```

#### macOS (using Homebrew)

```bash
# Coming soon - build from source for now
brew install adt
```

## 🖥️ Usage

> **Note:** Replace `adt` with `./adt` (Linux/macOS) or `.\adt.exe` (Windows) if the command is not found or not in your PATH.

### Interactive Wizard

Launch the interactive TUI wizard:

```bash
adt wizard
```

### CLI Commands

#### Detect System Information

```bash
# Basic detection
adt detect

# Quick detection (simpler output)
adt quick-detect

# JSON output (useful for scripts)
adt detect --json

# Verbose output
adt detect --verbose
```

#### Install Modules

```bash
# Install a specific module
adt install starship

# Install multiple modules
adt install starship gitconfig

# Dry run (preview what would happen)
adt install --dry-run starship

# Skip confirmation prompts
adt install --yes starship

# Force reinstallation
adt install --force starship
```

#### List Available Modules

```bash
adt list-modules
```

#### Additional Commands

```bash
# Quick OS detection (outputs just OS name, useful for scripts)
adt detect-quick

# Run detection wizard (shows detailed system info in wizard style)
adt wizard-detect

# Manage configuration (shows help for config-related commands)
adt config
```

#### Configuration Backup

```bash
# Create a backup
adt backup ~/.gitconfig

# List backups
adt list-backups

# Restore from backup
adt restore ~/.gitconfig 2024-01-15_10-30-00

# Delete old backups
adt delete-backups --older-than 30d
```

## 🎯 Examples

### Quick Start

```bash
# 1. Detect your system
adt detect

# 2. Launch the wizard
adt wizard

# 3. Select modules and install
```

### Automated Setup

```bash
# Detect system and install starship non-interactively
adt detect --json
adt install --yes starship
```

### Windows PowerShell Setup

```powershell
# NOTE: If 'adt' is not recognized, use '.\adt.exe' instead

# Detect system
.\adt.exe detect

# Install modules (uses Chocolatey or Winget automatically)
.\adt.exe install starship
```

### Linux/macOS Setup

```bash
# Full setup with Oh My Zsh (requires Zsh)
adt install ohmyzsh starship gitconfig
```

## 🔧 Configuration

### Config File Location

- Windows: `%APPDATA%\adt\config.yaml`
- macOS: `~/Library/Application Support/adt/config.yaml`
- Linux: `~/.config/adt/config.yaml`

### Backup Directory

All backups are stored in:
- Windows: `%APPDATA%\adt\backups\`
- macOS/Linux: `~/.adt/backups/`

## 🏗️ Architecture

```
auto-dev-terminal/
├── cmd/               # CLI commands
│   ├── cli/main.go   # Entry point
│   ├── detect.go     # System detection
│   ├── install.go    # Module installation
│   ├── wizard.go     # TUI wizard
│   └── config.go     # Backup/restore
├── pkg/
│   ├── detector/     # OS/shell/package manager detection
│   ├── installer/    # Package manager adapters
│   ├── config/       # Backup, templates, writing
│   ├── modules/      # Module system
│   └── wizard/        # Bubble Tea TUI
└── internal/
    ├── types/        # Type definitions
    └── constants/    # Constants
```

## 🤝 Contributing

Contributions are welcome! Please read our [Contributing Guide](CONTRIBUTING.md).

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🔒 Security

See our [Security Policy](SECURITY.md) for reporting vulnerabilities.

## 🙏 Acknowledgments

- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Starship](https://starship.rs/) - Cross-shell prompt
- [Oh My Zsh](https://ohmyz.sh/) - Zsh framework
