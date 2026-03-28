// Package installer provides cross-platform package manager adapters.
package installer

// AllInstallers returns all supported installer types (for registration purposes).
func AllInstallers() map[string]func() Installer {
	return map[string]func() Installer{
		"homebrew":   func() Installer { return NewBrewInstaller() },
		"chocolatey": func() Installer { return NewChocoInstaller() },
		"apt":        func() Installer { return NewAptInstaller() },
		"dnf":        func() Installer { return NewDnfInstaller() },
		"pacman":     func() Installer { return NewPacmanInstaller() },
		"scoop":      func() Installer { return NewScoopInstaller() },
		"winget":     func() Installer { return NewWingetInstaller() },
		"zypper":     func() Installer { return NewZypperInstaller() },
	}
}

// SupportedPackageManagers returns the list of supported package manager names.
func SupportedPackageManagers() []string {
	return []string{
		"homebrew",
		"chocolatey",
		"apt",
		"dnf",
		"pacman",
		"scoop",
		"winget",
		"zypper",
	}
}
