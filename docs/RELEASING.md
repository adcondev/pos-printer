# Release Process

## 🚀 DevOps and CI/CD

This project is built with a strong focus on automation, code quality, and a streamlined development process. The
following is a summary of the DevOps and CI/CD infrastructure:

### CI/CD Pipeline

The CI/CD pipeline is built on **GitHub Actions** and is designed to ensure code quality, security, and compatibility.
The pipeline consists of the following jobs:

* **Validate Commits:** Enforces conventional commit messages and semantic pull request titles to maintain a clean and
  readable project history.
* **Test and Coverage:** Runs the full suite of Go tests on Ubuntu, Windows, and macOS to ensure cross-platform
  compatibility. Code coverage reports are uploaded to **Codecov** to monitor code quality.
* **Benchmarks:** Runs and reports benchmarks on every pull request to monitor performance and prevent regressions.
* **Security Scan:** Performs a security scan using **Trivy** to identify and report vulnerabilities. The results are
  uploaded to the GitHub Security tab for easy tracking.
* **Go Linters:** Runs `golangci-lint` to enforce coding standards and best practices.

### Automated Release Management

The project uses an automated release process that triggers on merges to the `main` branch. The process is as follows:

1. **Analyze Commits:** The pipeline analyzes the commit messages since the last release to determine the next semantic
   version (major, minor, or patch).
2. **Generate Release Notes:** The `CHANGELOG.md` is updated based on the commit messages.
3. **Create GitHub Release:** A new GitHub release is created with the new version and release notes.
4. **Warm up Go Proxy:** The Go module proxy is notified of the new version to ensure it is available to developers as
   quickly as possible.

### Code Quality and Governance

The project uses a combination of tools to enforce code quality and governance:

* **`golangci-lint`:** A fast and comprehensive Go linter that checks for a wide range of issues.
* **`pre-commit` hooks:** A set of hooks that run before each commit to format code, run tests, and check for common
  issues.
* **`commitlint`:** Enforces conventional commit messages to ensure a clean and readable project history.

## Automated Release Flow

This project uses automated semantic versioning based on conventional commits.

### Version Bumping Rules

| Commit Type     | Example                          | Version Bump          |
|-----------------|----------------------------------|-----------------------|
| Breaking Change | `feat(escpos)!: refactor API`    | Major (1.0.0 → 2.0.0) |
| Feature         | `feat(print): add color support` | Minor (1.0.0 → 1.1.0) |
| Fix             | `fix(barcode): correct encoding` | Patch (1.0.0 → 1.0.1) |

### Current Version: v2.0.0

The v2.0.0 release includes major architectural changes:

- Modular ESC/POS architecture
- Enhanced graphics processing
- QR code support
- Improved print position control

### Integration Methods

| Method           | Description                   | Best For                       |
|------------------|-------------------------------|--------------------------------|
| **Squash Merge** | Combines all commits into one | Feature branches (recommended) |
| **Merge Commit** | Preserves all commit history  | When history is important      |
| **Rebase**       | Linear history                | Small fixes                    |

**Important**: For breaking changes, ensure the `!` indicator is preserved in the final commit message.

### Pre-release Checklist

- [ ] All tests passing
- [ ] Documentation updated
- [ ] CHANGELOG.md reviewed
- [ ] No uncommitted changes

### Post-release Tasks

1. Verify GitHub release was created
2. Check that Docker image was built and pushed
3. Update any dependent projects
4. Announce release if major version