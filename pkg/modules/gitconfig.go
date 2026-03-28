// Package modules provides the Git configuration module.
package modules

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/auto-dev-terminal/auto-dev-terminal/internal/types"
)

// GitConfigModule implements the Module interface for Git configuration.
type GitConfigModule struct {
	*BaseModule
}

// NewGitConfigModule creates a new Git configuration module.
func NewGitConfigModule() *GitConfigModule {
	return &GitConfigModule{
		BaseModule: NewBaseModule(
			"gitconfig",
			"Enhanced Git configuration with aliases, prompts, and useful settings",
			"1.0.0",
			[]string{},
		),
	}
}

// Install installs the Git configuration.
func (m *GitConfigModule) Install(opts *ModuleOptions) *ModuleResult {
	if opts.Verbose {
		fmt.Println("Installing Git configuration...")
	}

	// Check if Git is installed
	if !commandExists("git") {
		return &ModuleResult{
			Success: false,
			Module:  m.Name(),
			Error:   "Git is not installed",
		}
	}

	// Get Git version
	version, err := m.getGitVersion()
	if err != nil {
		return &ModuleResult{
			Success: false,
			Module:  m.Name(),
			Error:   fmt.Sprintf("failed to get Git version: %v", err),
		}
	}

	// Back up existing .gitconfig if it exists
	gitconfigPath := filepath.Join(opts.HomeDir, ".gitconfig")
	backupPath := gitconfigPath + ".bak"
	if _, err := os.Stat(gitconfigPath); err == nil {
		if err := copyFile(gitconfigPath, backupPath); err != nil {
			return &ModuleResult{
				Success: false,
				Module:  m.Name(),
				Error:   fmt.Sprintf("failed to backup .gitconfig: %v", err),
			}
		}
		if opts.Verbose {
			fmt.Printf("Backed up existing .gitconfig to %s\n", backupPath)
		}
	}

	// Write new configuration
	if err := m.writeGitConfig(opts); err != nil {
		return &ModuleResult{
			Success: false,
			Module:  m.Name(),
			Error:   fmt.Sprintf("failed to write Git config: %v", err),
		}
	}

	// Configure user name and email if provided
	if os.Getenv("GIT_USER_NAME") != "" {
		if err := m.runGitConfig("user.name", os.Getenv("GIT_USER_NAME")); err != nil {
			return &ModuleResult{
				Success: false,
				Module:  m.Name(),
				Error:   fmt.Sprintf("failed to set git user.name: %v", err),
			}
		}
	}
	if os.Getenv("GIT_USER_EMAIL") != "" {
		if err := m.runGitConfig("user.email", os.Getenv("GIT_USER_EMAIL")); err != nil {
			return &ModuleResult{
				Success: false,
				Module:  m.Name(),
				Error:   fmt.Sprintf("failed to set git user.email: %v", err),
			}
		}
	}

	// Install git-aware prompt if on a Unix-like system
	if opts.OS != types.OSWindows {
		_ = m.installGitPrompt(opts)
	}

	return &ModuleResult{
		Success: true,
		Module:  m.Name(),
		Output:  "Git configuration installed successfully",
		Version: version,
	}
}

// writeGitConfig writes the enhanced Git configuration.
func (m *GitConfigModule) writeGitConfig(opts *ModuleOptions) error {
	configPath := filepath.Join(opts.HomeDir, ".gitconfig")

	config := getGitConfigTemplate()

	if err := os.WriteFile(configPath, []byte(config), 0644); err != nil {
		return fmt.Errorf("writing config: %w", err)
	}

	return nil
}

// runGitConfig runs a git config command.
func (m *GitConfigModule) runGitConfig(key, value string) error {
	cmd := exec.Command("git", "config", "--global", key, value)
	return cmd.Run()
}

// getGitVersion returns the installed Git version.
func (m *GitConfigModule) getGitVersion() (string, error) {
	cmd := exec.Command("git", "--version")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	// Parse "git version 2.40.0" -> "2.40.0"
	parts := strings.Fields(string(output))
	if len(parts) >= 3 {
		return parts[2], nil
	}
	return string(output), nil
}

// installGitPrompt installs the git-aware prompt for bash/zsh.
func (m *GitConfigModule) installGitPrompt(opts *ModuleOptions) error {
	// Create the git-prompt directory
	promptDir := filepath.Join(opts.HomeDir, ".git-prompt")
	if err := os.MkdirAll(promptDir, 0755); err != nil {
		return fmt.Errorf("creating prompt directory: %w", err)
	}

	promptScript := filepath.Join(promptDir, "git-prompt.sh")

	if _, err := os.Stat(promptScript); os.IsNotExist(err) {
		// We don't actually need the full git repo, just the prompt script
		// For now, create a simple prompt script
		promptContent := getGitPromptScript()
		if err := os.WriteFile(promptScript, []byte(promptContent), 0644); err != nil {
			return fmt.Errorf("writing prompt script: %w", err)
		}
	}

	return nil
}

// Uninstall removes the Git configuration.
func (m *GitConfigModule) Uninstall(opts *ModuleOptions) *ModuleResult {
	if opts.Verbose {
		fmt.Println("Uninstalling Git configuration...")
	}

	gitconfigPath := filepath.Join(opts.HomeDir, ".gitconfig")
	backupPath := gitconfigPath + ".bak"

	// Check if our config is installed (look for our marker)
	data, err := os.ReadFile(gitconfigPath)
	if err != nil {
		return &ModuleResult{
			Success: false,
			Module:  m.Name(),
			Error:   "Git config file not found",
		}
	}

	// Check for our marker
	if !strings.Contains(string(data), "# Auto-Dev-Terminal Git Config") {
		return &ModuleResult{
			Success: false,
			Module:  m.Name(),
			Error:   "Git config was not installed by Auto-Dev-Terminal",
		}
	}

	// Restore backup if it exists
	if _, err := os.Stat(backupPath); err == nil {
		if err := os.Rename(backupPath, gitconfigPath); err != nil {
			return &ModuleResult{
				Success: false,
				Module:  m.Name(),
				Error:   fmt.Sprintf("failed to restore backup: %v", err),
			}
		}
		return &ModuleResult{
			Success: true,
			Module:  m.Name(),
			Output:  "Restored original .gitconfig from backup",
		}
	}

	// No backup, just remove
	if err := os.Remove(gitconfigPath); err != nil {
		return &ModuleResult{
			Success: false,
			Module:  m.Name(),
			Error:   fmt.Sprintf("failed to remove config: %v", err),
		}
	}

	return &ModuleResult{
		Success: true,
		Module:  m.Name(),
		Output:  "Git configuration removed",
	}
}

// IsInstalled checks if the enhanced Git configuration is installed.
func (m *GitConfigModule) IsInstalled() (bool, error) {
	gitconfigPath := filepath.Join(os.Getenv("HOME"), ".gitconfig")
	data, err := os.ReadFile(gitconfigPath)
	if err != nil {
		return false, nil
	}
	return strings.Contains(string(data), "# Auto-Dev-Terminal Git Config"), nil
}

// copyFile copies a file from src to dst.
func copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, data, 0644)
}

// getGitConfigTemplate returns the enhanced Git configuration template.
func getGitConfigTemplate() string {
	return `# Auto-Dev-Terminal Git Config
# Generated by Auto-Dev-Terminal

[user]
	name = 
	email = 

[core]
	editor = vim
	autocrlf = input
	whitespace = fix
	excludesfile = ~/.gitignore_global

[init]
	defaultBranch = main

[color]
	ui = auto

[color "branch"]
	current = yellow reverse
	local = yellow
	remote = green

[color "diff"]
	meta = yellow bold
	frag = magenta bold
	old = red bold
	new = green bold

[color "status"]
	added = green
	changed = yellow
	untracked = red

[alias]
	# Shortcuts
	s = status
	ss = status -sb
	a = add
	aa = add --all
	c = commit
	cm = commit -m
	ca = commit --amend
	co = checkout
	cob = checkout -b
	l = log
	ll = log --oneline --graph --decorate
	llg = log --graph --abbrev-commit --decorate --format=format:'%C(bold blue)%h%C(reset) - %C(bold green)(%ar)%C(reset) %C(white)%s%C(reset) %C(dim white)- %an%C(reset)%C(auto)%d%C(reset)'
	d = diff
	ds = diff --stat
	dc = diff --cached
	dr = diff --name-only
	r = remote
	ra = remote add
	rr = remote rm
	f = fetch
	p = pull
	m = merge
	b = branch
	ba = branch -a
	bd = branch -d
	bdd = branch -D
	
	# Undo
	undo = reset HEAD~1 --mixed
	undocommit = reset --soft HEAD~1
	
	# Stash
	ssave = stash save
	slist = stash list
	spop = stash pop
	sdrop = stash drop
	
	# Tags
	taglist = tag -l
	
	# Search
	grep = grep --color=auto
	find = log --all --full-history --oneline --author=
	
	# Work in progress
	wip = !git add -A && git commit -m 'WIP'
	unwip = reset HEAD~1

[push]
	default = simple
	autoSetupRemote = true
	followTags = true

[pull]
	rebase = true

[fetch]
	prune = true
	pruneTags = true

[rebase]
	autoStash = true
	autosquash = true

[merge]
	tool = vimdiff
	conflictstyle = diff3

[diff]
	tool = vimdiff
	algorithm = histogram

[difftool "vimdiff"]
	cmd = vim -d \"$LOCAL\" \"$REMOTE\"

[mergetool "vimdiff"]
	cmd = vim -d \"$LOCAL\" \"$MERGED\" \"$REMOTE\"

[help]
	autocorrect = 1

[credential]
	helper = store

[url "git@github.com:"]
	insteadOf = https://github.com/

[url "git@gitlab.com:"]
	insteadOf = https://gitlab.com/

[url "git@bitbucket.org:"]
	insteadOf = https://bitbucket.org/

[include]
	path = ~/.gitconfig_local

[filter "lfs"]
	clean = git-lfs clean -- %f
	smudge = git-lfs smudge -- %f
	process = git-lfs filter-process
	required = true
`
}

// getGitPromptScript returns a simple git-aware prompt script.
func getGitPromptScript() string {
	return `#!/bin/bash
# git-prompt.sh: bash/zsh support for showing git branch and status
# Simplified version for Auto-Dev-Terminal

export GIT_PS1_SHOWDIRTYSTATE=1
export GIT_PS1_SHOWSTASHSTATE=1
export GIT_PS1_SHOWUNTRACKEDFILES=1
export GIT_PS1_SHOWUPSTREAM=auto
export GIT_PS1_HIDE_IF_PWD_IGNORED=1

__git_ps1() {
    local g="$(git rev-parse --git-dir 2>/dev/null)"
    if [ -n "$g" ]; then
        local r=""
        local b="$(git symbolic-ref HEAD 2>/dev/null | sed 's|refs/heads/||')"
        
        if [ -d "$g/../.dotest" ]; then
            r="|REBASE-i"
        elif [ -d "$g/.dotest-merge" ]; then
            r="|REBASE-m"
        else
            git rev-parse --verify HEAD >/dev/null 2>&1 || return
            git diff-index --cached --quiet HEAD -- 2>/dev/null || r="|COMMIT"
        fi
        
        [ -n "$b" ] && printf "(%s)" "$b"
        [ -n "$r" ] && printf "%s" "$r"
    fi
}
`
}

// Ensure GitConfigModule implements Module interface
var _ Module = (*GitConfigModule)(nil)

// init registers the Git config module with the global registry.
func init() {
	_ = Register(NewGitConfigModule())
}
