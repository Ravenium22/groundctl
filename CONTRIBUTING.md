# Contributing to groundctl

Thanks for your interest in contributing! groundctl is an open-source project and we welcome contributions of all kinds.

## Quick Start

```bash
# Clone the repo
git clone https://github.com/groundctl/groundctl.git
cd groundctl

# Install Go 1.22+ (https://go.dev/dl/)

# Build
go build -o ground .

# Run tests
go test ./...

# Run the CLI
./ground --help
```

## Project Structure

```
groundctl/
  cmd/              # Cobra CLI commands
  internal/
    model/          # Core data types
    detector/       # Tool detection engine
    config/         # .ground.yaml load/save/validate
    version/        # Semver constraint checking
    drift/          # Diff engine (snapshot vs config)
    fixer/          # Fix strategy resolver + executor
    pkgmanager/     # Package manager detection
    profile/        # Profile management
    team/           # Team sharing (git/URL)
    shell/          # Shell hooks and prompt integration
    secrets/        # Secret reference system
    snapshot/       # Machine state capture
    report/         # Report formatters (JSON/HTML/MD)
  docs/             # Docusaurus docs site
  .github/          # CI workflows and GitHub Action
```

## Development Workflow

1. **Fork** the repository
2. **Create a branch** from `main`: `git checkout -b feat/my-feature`
3. **Make your changes** with tests
4. **Run tests**: `go test ./...`
5. **Run lint**: `golangci-lint run`
6. **Commit** with a clear message
7. **Open a PR** against `main`

## Code Style

- Follow standard Go conventions (`gofmt`, `go vet`)
- Write table-driven tests using `testify`
- Keep functions focused and under 50 lines where possible
- Use `lipgloss` styles for terminal output
- Error messages should be lowercase, without trailing punctuation

## Adding a New Tool to the Detector

The detector supports 18 tools. To add a new one:

1. Edit `internal/detector/detector.go`
2. Add an entry to the `registry` map with the detection command and version regex
3. Add the tool to `internal/fixer/strategy.go` with install commands per package manager
4. Add tests
5. Update the README tool table

## Adding a New Secret Backend

1. Create `internal/secrets/backend_<name>.go` implementing the `Backend` interface
2. Register it in `DefaultRegistry()` in `internal/secrets/secrets.go`
3. Add tests in `internal/secrets/secrets_test.go`
4. Document in `docs/docs/secrets.md`

## Good First Issues

Look for issues labeled [`good first issue`](https://github.com/groundctl/groundctl/labels/good%20first%20issue). These are scoped, well-documented tasks ideal for newcomers:

- Add detection for a new tool (e.g. `helm`, `deno`, `bun`)
- Add a package manager install command for an existing tool
- Improve error messages
- Add shell completion for a specific command
- Write documentation for a use case

## Running Tests

```bash
# All tests
go test ./...

# Specific package
go test ./internal/detector/...

# With verbose output
go test -v ./internal/version/...

# With coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Reporting Bugs

Open a [GitHub issue](https://github.com/groundctl/groundctl/issues/new) with:

- groundctl version (`ground version`)
- OS and architecture
- Steps to reproduce
- Expected vs actual behavior
- `ground doctor` output (if relevant)

## Code of Conduct

Be kind, be constructive, be inclusive. We follow the [Contributor Covenant](https://www.contributor-covenant.org/version/2/1/code_of_conduct/).

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
