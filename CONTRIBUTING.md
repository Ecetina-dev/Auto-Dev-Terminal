# Contributing to Auto-Dev-Terminal

Thank you for your interest in contributing to Auto-Dev-Terminal!

## Code of Conduct

By participating in this project, you are expected to uphold our [Code of Conduct](https://github.com/auto-dev-terminal/auto-dev-terminal/blob/main/CODE_OF_CONDUCT.md). Please report unacceptable behavior to [email protected].

## How Can I Contribute?

### Reporting Bugs

Before creating bug reports, please check the [issue tracker](https://github.com/auto-dev-terminal/auto-dev-terminal/issues) to see if the problem has already been reported. If it has, add a comment to the existing issue instead of creating a new one.

When creating a bug report, please include:
- A quick summary and background
- Steps to reproduce
- What you expected vs what happened
- Notes (possibly including why you think this might be happening)

### Suggesting Enhancements

Enhancement suggestions are tracked as GitHub issues. When creating an enhancement suggestion:
- Use a clear and descriptive title
- Provide a step-by-step description of the suggested enhancement
- Explain why this enhancement would be useful

### Pull Requests

1. Fork the repo and create your branch from `main`
2. If you've added code that should be tested, add tests
3. If you've changed APIs, update the documentation
4. Ensure the test suite passes
5. Make sure your code lints
6. Use meaningful commit messages
7. Write a clear description of your PR

## Development Setup

### Prerequisites

- Go 1.24 or later
- Make

### Local Development

```bash
# Clone the repository
git clone https://github.com/auto-dev-terminal/auto-dev-terminal.git
cd auto-dev-terminal

# Install dependencies
go mod download

# Run tests
go test ./...

# Build
go build -o bin/auto-dev-terminal ./cmd/cli

# Run linter
make lint
```

### Coding Standards

- Follow Go standard coding conventions
- Use meaningful variable and function names
- Add comments to exported functions
- Write tests for new functionality
- Ensure code passes `golangci-lint`

### Commit Messages

We use [Conventional Commits](https://www.conventionalcommits.org/) for commit messages:

```
<type>(<scope>): <description>

[optional body]

[optional footer]
```

Types:
- `feat`: A new feature
- `fix`: A bug fix
- `docs`: Documentation only changes
- `style`: Changes that do not affect the meaning of the code
- `refactor`: A code change that neither fixes a bug nor adds a feature
- `test`: Adding missing tests or correcting existing tests
- `chore`: Changes to the build process or auxiliary tools

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
