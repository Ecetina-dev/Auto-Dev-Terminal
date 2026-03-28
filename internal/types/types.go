// Package types provides core type definitions for the Auto-Dev-Terminal system.
package types

import "time"

// OS represents the detected operating system
type OS string

const (
	OSWindows OS = "windows"
	OSDarwin  OS = "darwin"
	OSLinux   OS = "linux"
)

// Distro represents Linux distribution
type Distro string

const (
	DistroUbuntu    Distro = "ubuntu"
	DistroDebian    Distro = "debian"
	DistroFedora    Distro = "fedora"
	DistroRHEL      Distro = "rhel"
	DistroCentOS    Distro = "centos"
	DistroRocky     Distro = "rocky"
	DistroAlma      Distro = "alma"
	DistroArch      Distro = "arch"
	DistroManjaro   Distro = "manjaro"
	DistroEndeavour Distro = "endeavouros"
	DistroOpenSUSE  Distro = "opensuse"
	DistroSLES      Distro = "sles"
	DistroLinuxMint Distro = "linuxmint"
	DistroPop       Distro = "pop"
	DistroUnknown   Distro = "unknown"
)

// Shell represents the detected shell
type Shell string

const (
	ShellBash       Shell = "bash"
	ShellZsh        Shell = "zsh"
	ShellFish       Shell = "fish"
	ShellPowerShell Shell = "powershell"
	ShellPwsh       Shell = "pwsh"
	ShellCmd        Shell = "cmd"
	ShellTcsh       Shell = "tcsh"
	ShellCsh        Shell = "csh"
	ShellAsh        Shell = "ash"
	ShellUnknown    Shell = "unknown"
)

// PackageManager represents available package managers
type PackageManager string

const (
	PkgMgrHomebrew   PackageManager = "homebrew"
	PkgMgrMacPorts   PackageManager = "macports"
	PkgMgrChocolatey PackageManager = "chocolatey"
	PkgMgrScoop      PackageManager = "scoop"
	PkgMgrWinget     PackageManager = "winget"
	PkgMgrAPT        PackageManager = "apt"
	PkgMgrDNF        PackageManager = "dnf"
	PkgMgrYUM        PackageManager = "yum"
	PkgMgrPacman     PackageManager = "pacman"
	PkgMgrZypper     PackageManager = "zypper"
	PkgMgrSnap       PackageManager = "snap"
	PkgMgrFlatpak    PackageManager = "flatpak"
	PkgMgrUnknown    PackageManager = "unknown"
)

// SystemInfo holds all detection results
type SystemInfo struct {
	OS              OS               `json:"os" yaml:"os"`
	Distro          Distro           `json:"distro" yaml:"distro"`
	DistroVersion   string           `json:"distro_version" yaml:"distro_version"`
	Shell           Shell            `json:"shell" yaml:"shell"`
	ShellVersion    string           `json:"shell_version" yaml:"shell_version"`
	PackageManagers []PackageManager `json:"package_managers" yaml:"package_managers"`
	HomeDir         string           `json:"home_dir" yaml:"home_dir"`
	Username        string           `json:"username" yaml:"username"`
	Hostname        string           `json:"hostname" yaml:"hostname"`
	Arch            string           `json:"arch" yaml:"arch"`
}

// Module represents an installable module
type Module struct {
	Name         string        `yaml:"name" json:"name"`
	DisplayName  string        `yaml:"display_name" json:"display_name"`
	Description  string        `yaml:"description" json:"description"`
	Version      string        `yaml:"version" json:"version"`
	Dependencies []string      `yaml:"dependencies" json:"dependencies"`
	Install      InstallConfig `yaml:"install" json:"install"`
	Config       ConfigConfig  `yaml:"config" json:"config"`
	Requirements []Requirement `yaml:"requirements" json:"requirements"`
}

// InstallConfig defines how to install the module
type InstallConfig struct {
	PackageManager PackageManager `yaml:"package_manager" json:"package_manager"`
	Commands       []string       `yaml:"commands" json:"commands"`
	ScriptURL      string         `yaml:"script_url,omitempty" json:"script_url,omitempty"`
}

// ConfigConfig defines configuration file handling
type ConfigConfig struct {
	Source      string            `yaml:"source" json:"source"`
	Destination string            `yaml:"destination" json:"destination"`
	Variables   map[string]string `yaml:"variables" json:"variables"`
}

// Requirement defines a prerequisite for the module
type Requirement struct {
	Type     string `yaml:"type" json:"type"` // "shell", "os", "command"
	Value    string `yaml:"value" json:"value"`
	Optional bool   `yaml:"optional" json:"optional"`
}

// InstallResult represents the outcome of an installation attempt
type InstallResult struct {
	Success bool   `json:"success" yaml:"success"`
	Module  string `json:"module" yaml:"module"`
	Error   string `json:"error,omitempty" yaml:"error,omitempty"`
	Output  string `json:"output,omitempty" yaml:"output,omitempty"`
	Version string `json:"version,omitempty" yaml:"version,omitempty"`
}

// Backup represents a configuration file backup
type Backup struct {
	ID        string    `json:"id" yaml:"id"`
	Original  string    `json:"original" yaml:"original"`
	Backup    string    `json:"backup" yaml:"backup"`
	Timestamp time.Time `json:"timestamp" yaml:"timestamp"`
	Size      int64     `json:"size" yaml:"size"`
}

// BackupManifest tracks all backups
type BackupManifest struct {
	Backups []Backup `json:"backups" yaml:"backups"`
}

// TemplateVariables holds the variables available for template rendering
type TemplateVariables struct {
	OS       string
	Shell    string
	HomeDir  string
	Username string
	Hostname string
	Distro   string
	Arch     string
}

// InstallOptions configures an installation
type InstallOptions struct {
	Sudo    bool
	Yes     bool
	Verbose bool
	DryRun  bool
}
