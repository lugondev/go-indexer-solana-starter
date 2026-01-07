# Contributing to Go Indexer Solana Starter

Thank you for your interest in contributing! This document provides guidelines for contributing to the project.

## Development Setup

1. Fork the repository
2. Clone your fork: `git clone https://github.com/your-username/go-indexer-solana-starter.git`
3. Create a branch: `git checkout -b feature/your-feature`
4. Install dependencies: `go mod download`
5. Install dev tools: `make install-tools`

## Code Standards

### Go Style Guide
- Follow official Go style guidelines
- Use `gofmt` and `goimports` for formatting
- Write idiomatic Go code
- Keep functions small and focused
- Use meaningful variable names

### Naming Conventions
- **Packages**: lowercase, single-word
- **Files**: lowercase with underscores
- **Exported**: PascalCase
- **Unexported**: camelCase
- **Interfaces**: -er suffix

### Testing
- Write tests for all new code
- Maintain >80% test coverage
- Use table-driven tests
- Test edge cases and error conditions
- Run tests with race detector: `go test -race ./...`

### Error Handling
- Always check and handle errors
- Use `fmt.Errorf` with `%w` for error wrapping
- Create custom error types when needed
- Provide meaningful error messages

## Commit Guidelines

### Commit Message Format
```
<type>(<scope>): <subject>

<body>

<footer>
```

### Types
- **feat**: New feature
- **fix**: Bug fix
- **docs**: Documentation changes
- **style**: Code style changes (formatting)
- **refactor**: Code refactoring
- **test**: Test changes
- **chore**: Build/tool changes

### Examples
```
feat(indexer): add WebSocket support for real-time updates

Implement WebSocket connection to Solana RPC for receiving
real-time block updates instead of polling.

Closes #123
```

```
fix(config): handle missing environment variables gracefully

Previously the app would crash if optional env vars were missing.
Now it uses sensible defaults.
```

## Pull Request Process

1. Update documentation if needed
2. Add tests for new functionality
3. Ensure all tests pass: `make test`
4. Run linters: `make lint`
5. Format code: `make fmt`
6. Update CHANGELOG.md
7. Create a pull request with clear description

### PR Checklist
- [ ] Tests added/updated
- [ ] Documentation updated
- [ ] Code formatted with `gofmt`
- [ ] No linting errors
- [ ] CHANGELOG.md updated
- [ ] Commit messages follow guidelines

## Code Review

- Be respectful and constructive
- Focus on code, not the person
- Explain your suggestions
- Be open to feedback
- Respond promptly to review comments

## Questions?

Feel free to open an issue for:
- Bug reports
- Feature requests
- Questions about the codebase
- Suggestions for improvements

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
