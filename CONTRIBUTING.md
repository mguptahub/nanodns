# Contributing to NanoDNS

Thank you for your interest in contributing to NanoDNS! This document will guide you through the contribution process.

## Getting Started

1. Fork the repository on GitHub
2. Clone your fork:
```bash
git clone https://github.com/YOUR_USERNAME/nanodns.git
cd nanodns
```

3. Add the upstream repository:
```bash
git remote add upstream https://github.com/mguptahub/nanodns.git
```

4. Create a new branch for your changes:
```bash
git checkout -b feature/your-feature-name
```

## Before Creating a Pull Request

1. Sync your fork with upstream:
```bash
git fetch upstream
git rebase upstream/main
```

2. Test your changes locally:
```bash
# Run tests
go test ./...

# Build and test the binary
go build -o nanodns cmd/server/main.go
```

3. Ensure your code follows the project's style:
- Follow standard Go coding conventions
- Use meaningful variable and function names
- Add comments for complex logic
- Keep functions small and focused

4. Update documentation if needed:
- Update README.md if you've added new features
- Add/update code comments
- Update examples if relevant

## Creating a Pull Request

1. Push your changes to your fork:
```bash
git push origin feature/your-feature-name
```

2. Go to GitHub and create a Pull Request

3. In your PR description:
- Clearly describe the changes
- Link any related issues
- Include examples if applicable
- Mention breaking changes if any

## After Creating a Pull Request

1. Wait for PR review
   - Maintainers will review your code
   - Address any requested changes
   - Keep the PR updated with upstream changes

2. Wait for approval
   - At least one maintainer must approve
   - All discussions must be resolved
   - All checks must pass

3. Wait for release
   - Once merged, your changes will be included in the next release
   - Releases are created by maintainers
   - Your name will be included in the changelog

## Release Process

Releases are managed by maintainers following this process:

1. Code is merged into main
2. Maintainers prepare release
3. New version tag is created
4. GitHub Actions automatically:
   - Run tests
   - Build binaries
   - Create GitHub release
   - Push Docker images
   - Update documentation

## Need Help?

Feel free to:
- Join our discussions on GitHub
- Ask questions in issues
- Reach out to maintainers

Thank you for contributing!