# NanoDNS Open Source Guidelines

Welcome to the NanoDNS community! This document outlines how you can contribute to NanoDNS and be a part of its development.

## Table of Contents
- [Code of Conduct](#code-of-conduct)
- [How to Contribute](#how-to-contribute)
- [Development Workflow](#development-workflow)
- [Community](#community)
- [Recognition](#recognition)

## Code of Conduct

### Our Pledge
We pledge to make participation in our project and our community a harassment-free experience for everyone, regardless of:
- Age, body size, disability, ethnicity, gender identity and expression
- Level of experience, education, socio-economic status
- Nationality, personal appearance, race, religion
- Sexual identity and orientation

### Expected Behavior
- Use welcoming and inclusive language
- Be respectful of differing viewpoints and experiences
- Accept constructive criticism gracefully
- Focus on what is best for the community
- Show empathy towards other community members

### Unacceptable Behavior
- Use of sexualized language or imagery
- Trolling, insulting/derogatory comments, and personal attacks
- Public or private harassment
- Publishing others' private information without permission
- Other conduct which could reasonably be considered inappropriate

## How to Contribute

### 1. Getting Started
- Star ⭐ and fork the repository
- Join our [GitHub Discussions](https://github.com/mguptahub/nanodns/discussions)
- Look for issues labeled `good first issue` or `help wanted`

### 2. Making Changes
```bash
# Fork and clone
git clone https://github.com/YOUR-USERNAME/nanodns.git

# Create branch
git checkout -b feature/your-feature-name

# Make changes
# Write tests
# Update documentation

# Commit (use conventional commits)
git commit -m "feat: add new feature"
```

### 3. Pull Request Process
1. Update documentation
2. Add tests for new features
3. Ensure all tests pass
4. Update the changelog
5. Submit PR with clear description

### Commit Message Format
We use [Conventional Commits](https://www.conventionalcommits.org/):
```
type(scope): description

[optional body]

[optional footer]
```

Types:
- feat: New feature
- fix: Bug fix
- docs: Documentation
- style: Formatting
- refactor: Code restructuring
- test: Adding tests
- chore: Maintenance

## Development Workflow

### Setting Up Development Environment
```bash
# Install dependencies
go mod download

# Run tests
go test ./...

# Build binary
go build -o nanodns cmd/server/main.go
```

### Code Style
- Follow [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Use `gofmt` for formatting
- Write meaningful comments
- Document public functions

### Testing
- Write unit tests for new features
- Ensure existing tests pass
- Add integration tests when needed
- Test documentation changes

## Community

### Join the Discussion
- GitHub Discussions: Technical discussions and questions
- Issue Tracker: Bug reports and feature requests
- Pull Requests: Code review and contribution

### Communication Guidelines
1. **Be Clear**
   - Write detailed descriptions
   - Include steps to reproduce bugs
   - Explain the context

2. **Be Respectful**
   - Value others' time
   - Acknowledge others' contributions
   - Be patient with new contributors

3. **Be Collaborative**
   - Share knowledge
   - Help others learn
   - Accept different viewpoints

## Recognition

### Contributor Levels
1. **First Time Contributor**
   - First PR merged
   - Added to Contributors list

2. **Regular Contributor**
   - Multiple PRs merged
   - Helping in discussions
   - Bug fixes and features

3. **Core Contributor**
   - Consistent contributions
   - Code review participation
   - Community support

### Hall of Fame
- Contributors are listed in [CONTRIBUTORS.md](../CONTRIBUTORS.md)
- Special mentions in release notes
- Featured in project documentation

### Getting Help
- Ask in GitHub Discussions
- Check existing documentation
- Review closed issues
- Contact maintainers

## License and Legal
- Licensed under [AGPL-2.0](../LICENSE)
- Sign Developer Certificate of Origin (DCO)
- Maintain copyright notices

## Project Governance
- Decisions made through community consensus
- Maintainers guide overall direction
- Community feedback valued and considered

Remember:
- No contribution is too small
- Quality over quantity
- Community and code both matter

Thank you for contributing to NanoDNS! 🎉