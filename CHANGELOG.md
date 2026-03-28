# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Interactive TUI wizard for guided installation
- Module system with pluggable architecture
- Built-in modules: Starship, Oh My Zsh, Git Config, Fonts
- Automatic system detection (OS, Shell, Package Managers)
- Backup and restore functionality for configuration files
- Cross-platform support (Windows, macOS, Linux)
- GitHub Actions CI/CD pipeline with multi-platform testing

### Fixed
- Module registration idempotency (fixed duplicate registration error)
- Proper error handling for missing dependencies

### Changed
- Migrated to Go 1.24
- Updated Cobra to v1.8.1
- Updated Bubble Tea to v0.26.0

## [0.1.0] - 2024-01-15

### Added
- Initial release
- System detection (OS, Shell, Package Manager)
- Basic CLI commands (detect, install, config)
- Starship module
- Configuration backup system

[Unreleased]: https://github.com/auto-dev-terminal/auto-dev-terminal/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/auto-dev-terminal/auto-dev-terminal/releases/tag/v0.1.0
