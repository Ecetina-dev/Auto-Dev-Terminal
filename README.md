# Auto-Dev-Terminal

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.24+-00ADD8?style=for-the-badge&logo=go" alt="Go Version">
  <img src="https://img.shields.io/badge/platforms-Windows%20%7C%20macOS%20%7C%20Linux-blue?style=for-the-badge&platforms=Windows%20%7C%20macOS%20%7C%20Linux" alt="Platforms">
  <img src="https://img.shields.io/github/actions/workflow/status/auto-dev-terminal/auto-dev-terminal/ci.yml?style=for-the-badge" alt="CI">
  <img src="https://img.shields.io/github/license/auto-dev-terminal/auto-dev-terminal?style=for-the-badge" alt="License">
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

```bash
git clone https://github.com/auto-dev-terminal/auto-dev-terminal.git
cd auto-dev-terminal
go build -o auto-dev-terminal ./cmd/cli
sudo mv auto-dev-terminal /usr/local/bin/
```

### Pre-built Binaries

Download the latest release from the [releases page](https://github.com/auto-dev-terminal/auto-dev-terminal/releases).

## 🖥️ Usage

### Interactive Wizard

Launch the interactive TUI wizard:

```bash
auto-dev-terminal wizard
```

### CLI Commands

#### Detect System Information

```bash
# Basic detection
auto-dev-terminal detect

# JSON output (useful for scripts)
auto-dev-terminal detect --json

# Verbose output
auto-dev-terminal detect --verbose
```

#### Install Modules

```bash
# Install a specific module
auto-dev-terminal install starship

# Install multiple modules
auto-dev-terminal install starship gitconfig

# Dry run (preview what would happen)
auto-dev-terminal install --dry-run starship

# Skip confirmation prompts
auto-dev-terminal install --yes starship

# Force reinstallation
auto-dev-terminal install --force starship
```

#### List Available Modules

```bash
auto-dev-terminal list-modules
```

#### Configuration Backup

```bash
# Create a backup
auto-dev-terminal backup ~/.gitconfig

# List backups
auto-dev-terminal list

# Restore from backup
auto-dev-terminal restore ~/.gitconfig 2024-01-15_10-30-00

# Delete old backups
auto-dev-terminal delete --older-than 30d
```

## 🎯 Examples

### Quick Start

```bash
# 1. Detect your system
auto-dev-terminal detect

# 2. Launch the wizard
auto-dev-terminal wizard

# 3. Select modules and install
```

### Automated Setup

```bash
# Detect system and install starship non-interactively
auto-dev-terminal detect --json
auto-dev-terminal install --yes starship
```

### Windows PowerShell Setup

```powershell
# Detect system
.\auto-dev-terminal.exe detect

# Install modules via Chocolatey or Winget
.\auto-dev-terminal.exe install starship
```

### Linux/macOS Setup

```bash
# Full setup with Oh My Zsh (requires Zsh)
auto-dev-terminal install ohmyzsh starship gitconfig
```

## 🔧 Configuration

### Config File Location

- Windows: `%APPDATA%\auto-dev-terminal\config.yaml`
- macOS: `~/Library/Application Support/auto-dev-terminal/config.yaml`
- Linux: `~/.config/auto-dev-terminal/config.yaml`

### Backup Directory

All backups are stored in:
- Windows: `%APPDATA%\auto-dev-terminal\backups\`
- macOS/Linux: `~/.auto-dev-terminal/backups/`

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
