# Security Policy

## Supported Versions

We support the latest version of Auto-Dev-Terminal. Older versions may not be supported.

| Version | Supported          |
| ------- | ------------------ |
| 0.1.x   | :white_check_mark: |

## Reporting a Vulnerability

If you discover a security vulnerability within Auto-Dev-Terminal, please send an email to [email protected]. All security vulnerabilities will be promptly addressed.

Please include the following information:
- Type of vulnerability
- Full paths of source file(s) related to the vulnerability
- Location of the affected source code (tag/branch/commit or direct URL)
- Any special configuration required to reproduce the issue
- Step-by-step instructions to reproduce the issue
- Proof-of-concept or exploit code (if possible)
- Impact of the issue, including how an attacker might exploit it

## Security Best Practices

When using Auto-Dev-Terminal:

1. **Review before installing**: Always review what a module will install before confirming
2. **Use dry-run**: Use `--dry-run` to preview changes without executing them
3. **Backup first**: Always backup your configuration files before making changes
4. **Verify sources**: Only install from trusted sources
5. **Keep updated**: Use the latest version to receive security fixes

## Security Considerations

- Auto-Dev-Terminal modifies shell configuration files
- Installation may require elevated privileges (sudo)
- Always review the changes before accepting
- Backup your dotfiles before major changes
