# Contributing to Meilisearch MCP

Thanks for wanting to improve this project! This guide explains the expectations, workflow, and conventions. Please read it fully before opening an issue or pull request (PR).

---
## 1. Ground Rules
- Be respectful & inclusive. We follow the [Contributor Covenant v2.1](https://www.contributor-covenant.org/version/2/1/code_of_conduct/) (implicit until a CODE_OF_CONDUCT.md is added).
- Favor small, focused PRs over large rewrites.
- Discuss substantial architectural changes in an issue before implementation.
- Keep the main branch green: add/keep tests passing & follow formatting.
- Security issues: **do not** open a public issue; email the maintainer instead (add GPG/signature if sensitive).

---
## 2. How to Contribute
1. Fork the repo & create a feature branch from `main`.
2. Make changes (see Development below).
3. Ensure formatting, lint, and tests pass.
4. Write helpful commit messages (see Conventional Commits below).
5. Open a PR with the required title format (see PR Naming below).
6. Respond to review feedback promptly.

---
## 3. Development
### Prerequisites
- Go 1.24+
- (Optional) Running Meilisearch instance for integration testing

### Useful Commands
```sh
# Build binary
make build

# Run unit tests
make tests

# Format code (gofumpt)
make fmt

# Lint (golangci-lint – must be installed / available as a go tool)
make check

# Run all (format + lint)
make check-all
```

### Running the Server (examples)
```sh
./build/meilisearch-mcp serve http  --addr :8080 --meili-host http://localhost:7700 --meili-api-key masterKey
./build/meilisearch-mcp serve stdio --meili-host http://localhost:7700 --meili-api-key masterKey
```

---
## 4. Code Style
- Formatting: enforced by `gofumpt` (via `make fmt`).
- Linting: enforced by `golangci-lint` (disable rules only with clear justification).
- Keep functions small & purposeful. Avoid premature abstraction.
- Prefer returning errors with context (`fmt.Errorf("context: %w", err)`).
- Public exports require concise GoDoc comments.

---
## 5. Testing
- Add/adjust tests for all behavioral changes.
- Aim for table-driven tests in Go where logical.
- Name test files `*_test.go`; keep helpers internal to the package unless shared widely.
- Fast unit tests are preferred; long-running or integration-like flows should be skipped by default behind a build tag if added.

---
## 6. Versioning & Releases
This project follows [Semantic Versioning](https://semver.org/). Version data lives in `version/version.go`.
- `fix:` or `perf:` usually increments PATCH.
- `feat:` increments MINOR.
- Breaking changes (`!` or documented in commit body) increment MAJOR.
Release process (maintainer): update `version/version.go`, tag (`vX.Y.Z`), draft GitHub Release with changelog generated from commit history.

---
## 7. Conventional Commits (Commit Message Format)
Format:
```
<type>[optional(scope)][!]: <short summary>

[optional body]
[optional footer(s)]
```
Allowed `type` values here:
- feat:     New user-facing feature
- fix:      Bug fix
- chore:    Non-production code changes (tooling, deps)
- refactor: Code change that neither fixes a bug nor adds a feature
- test:     Add or adjust tests only
- ci:       Continuous integration / build pipeline changes
- docs:     Documentation only changes
- perf:     Performance improvement
- style:    Formatting / stylistic (no logic change; rarely needed if gofumpt is used)
- build:    Build system or dependency changes
- revert:   Reverts a previous commit

Examples:
```
feat: add HTTP rate limiting to pool layer
fix(search): handle empty query causing panic
refactor: simplify transport middleware chain
chore(deps): bump golang.org/x/net to v0.25.0
ci: add lint job with caching
perf: use sync.Pool for request buffers
feat!: remove deprecated stdio flag --mcp-pool (breaking)
```
Footers (when needed):
```
BREAKING CHANGE: <explanation>
Refs: #123
Closes: #456
```

---
## 8. Pull Request Naming & Content
PR Titles MUST start with a bracketed category (capitalized) followed by a concise description:
```
[Feat] Add streaming search response support
[Fix] Prevent nil pointer in key rotation
[Chore] Update gofumpt and lint configs
[Refactor] Consolidate HTTP handlers
[Docs] Add configuration examples for docker-compose
[Test] Improve coverage for pool timeouts
[CI] Add matrix build for Go versions
[Perf] Reduce allocations in index serialization
```
Mapping (PR prefix -> commit type):
- [Feat] -> feat
- [Fix] -> fix
- [Chore] -> chore / build / style
- [Refactor] -> refactor
- [Docs] -> docs
- [Test] -> test
- [CI] -> ci
- [Perf] -> perf

PR Checklist (include in description):
- [ ] Linked issue (or `N/A`)
- [ ] Tests added/updated
- [ ] `make fmt` run clean
- [ ] `make check` passes
- [ ] `make tests` passes
- [ ] Docs updated (README / examples) if user-facing
- [ ] No breaking changes (or documented with BREAKING CHANGE footer)

---
## 9. Branching
- `main` is always releasable.
- Use short-lived feature branches: `feat/rate-limiter`, `fix/panic-empty-query`.
- Rebase (preferred) or squash locally to keep history clean; avoid merge commits in PR branches.

---
## 10. Issues
Before opening a new issue:
- Search existing issues / PRs.
- Provide reproduction steps, expected vs actual, environment (Go version, OS, Meilisearch version).
- For feature requests: explain problem, not just solution; add motivation & potential impact.

---
## 11. Security
Report vulnerabilities privately (do not open a public issue). Provide reproduction details and impact assessment if possible.

---
## 12. Style Quick Reference
| Concern        | Rule / Tool         |
|----------------|---------------------|
| Formatting     | gofumpt (`make fmt`) |
| Linting        | golangci-lint (`make check`) |
| Testing        | `go test -v ./...` or `make tests` |
| Versioning     | Semantic Versioning 2.0.0 |
| Commit format  | Conventional Commits |
| PR title       | [Type] Description   |

---
## 13. FAQ
**Q: My lint fails locally. What now?**  
Update or install `golangci-lint` (`go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest`) and re-run.

**Q: Do I bump the version in my PR?**  
Only if the change is release-ready and you're coordinating a release; otherwise leave it.

**Q: Can I open a draft PR early?**  
Yes—use Draft to gather feedback; ensure checklist passes before marking Ready.

---
Thank you for contributing! 🚀

