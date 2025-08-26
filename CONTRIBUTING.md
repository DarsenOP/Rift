# Contributing to Rift

We use a milestone-based development workflow with strong automation and quality gates.

### Branch Strategy

- Main development branch: dev (always stable and deployable)
- Create feature branches from dev:
```bash
    git checkout dev
    git pull origin dev
    git checkout -b feat/your-feature-name
```
- Use milestone branches (e.g., feat/setup) as integration branches for larger initiatives
- Push your branch and open a Pull Request to the appropriate target branch (dev or a milestone branch)

### Code Quality & Automation

- CI/CD Required: All PRs must pass GitHub Actions checks (linting + tests)
- Pre-merge checks: CI runs on every push and PR, blocking merge until green
- Review required: At least one maintainer must approve before merge
- Squash and merge: Keep history clean with focused commit messages

### Development Setup

```bash
make help        # Show all available targets
make lint        # Run linter (golangci-lint) - skips if no Go files
make test        # Run tests - skips if no test files  
make build       # Build binary - checks for valid project structure
make fmt         # Format code with gofumpt
make clean       # Remove built artifacts
```

### Branch Naming Convention

- feat/* – New features or functionality
- fix/* – Bug fixes and patches
- chore/* – Maintenance, tooling, config changes
- docs/* – Documentation updates
- test/* – Experimental or test code

### Commit Message Convention

Use conventional commit prefixes with milestone scope:

- `feat(m1):` for new features in milestone 1
- `fix(m1):` for bug fixes in milestone 1  
- `docs(m1):` for documentation in milestone 1
- `chore(m1):` for maintenance in milestone 1
- `ci(m1):` for CI/CD changes in milestone 1
- `test(m1):` for test-related changes in milestone 1

**Examples:**

- `feat(m1): initialize Go module and project structure`
- `fix(m2): resolve RESP parser edge case`
- `chore(m1): optimize CI caching strategy`

### Code Standards

- Formatting: gofumpt (enforced via make fmt)
- Linting: golangci-lint (enforced in CI)
- Testing: Write tests for new functionality
- Documentation: Public APIs must be documented
- Imports: Organized with goimports-style grouping

### Pull Request Process

- Ensure your branch is updated with target branch
- Run make lint and make test locally
- Push your changes and open a PR
- Ensure CI passes all checks
- Request review from maintainers
- Address review feedback if needed
- Squash and merge after approval

We value clean, maintainable code and collaborative development!
