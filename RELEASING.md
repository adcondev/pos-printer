# Release Process

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