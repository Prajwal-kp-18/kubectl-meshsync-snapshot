# Contributing to kubectl-meshsync-snapshot

Thank you for your interest in contributing to kubectl-meshsync-snapshot! This document provides guidelines and instructions for contributing to this project.

## Prerequisites

- Go 1.22+
- Access to a Kubernetes cluster for testing
- kubectl installed

## Development Workflow

1. Fork the repository
2. Create a new branch for your feature or bugfix
3. Implement your changes
4. Add tests for your changes
5. Ensure all tests pass with `make test`
6. Format your code with `make fmt`
7. Submit a pull request

## Building and Testing

The project includes a Makefile with several useful targets:

```bash
# Build the binary
make build

# Run tests
make test

# Format code
make fmt

# Clean build artifacts
make clean

# Create distribution packages
make dist
```

## Code Style

- Follow Go best practices and idiomatic Go patterns
- Use meaningful variable and function names
- Write clear comments for public functions and types
- Ensure proper error handling

## Pull Request Process

1. Ensure your code follows the code style guidelines
2. Update the README.md with details of changes if needed
3. Increase version numbers if applicable
4. The PR should be reviewed by at least one maintainer

## Releasing

Project maintainers can release new versions by:

1. Updating version information in both code and docs
2. Creating a new Git tag with the format `v{major}.{minor}.{patch}`
3. Pushing the tag to trigger the GitHub Actions release workflow
4. Creating a PR to the Krew index (if necessary)

## License

By contributing to this project, you agree that your contributions will be licensed under the same MIT License that covers the project. 