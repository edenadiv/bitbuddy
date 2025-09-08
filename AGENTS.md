# Repository Guidelines

This guide helps contributors work efficiently on the BitBuddy TUI written in Go (Bubble Tea + Lip Gloss).

## Project Structure & Module Organization
- Root module: `go.mod` (`module bitbuddy`).
- Entry point: `main.go` (program startup).
- Core domain: `bitbuddy.go` (state and behaviors).
- TUI/UI: `tui.go` (views, update loop, styles).
- Persistence: `storage.go` (JSON read/write to `bitbuddy.json`).
- Binary output: `bitbuddy` (ignored by `.gitignore`).

## Build, Test, and Development Commands
- Build: `go build -o bitbuddy .` — builds the TUI binary.
- Run: `go run .` — runs directly from source.
- Format: `go fmt ./...` — formats code per Go standards.
- Vet: `go vet ./...` — basic static checks.
- Test: `go test ./...` — runs unit tests (none yet; see below).

## Coding Style & Naming Conventions
- Use standard Go style; always run `go fmt` before committing.
- Indentation: tabs (default for Go tools).
- Names: exported identifiers use `CamelCase`; files are lowercase with underscores only when helpful (e.g., `bitbuddy.go`, `storage.go`).
- Keep packages simple; this repo currently uses `package main` across files.

## Testing Guidelines
- Framework: Go’s built‑in `testing` package.
- Location: tests live alongside code as `*_test.go` (e.g., `bitbuddy_test.go`).
- Focus: unit tests for behaviors in `bitbuddy.go` and persistence in `storage.go` (use temp dirs/files).
- Run locally with `go test ./...`. Optional coverage: `go test -cover ./...`.

## Commit & Pull Request Guidelines
- Commit style: follow Conventional Commits where practical (e.g., `feat: ...`, `fix(tui): ...`) consistent with history.
- Pull Requests should include:
  - Clear description of changes and rationale.
  - Screenshots or brief recordings of the TUI when UI changes are made.
  - Linked issue(s) if applicable.
  - Confirmation that `go fmt` and `go vet` pass.

## Security & Configuration Tips
- The app writes state to `bitbuddy.json` in the repo root. Do not commit personal state; consider adding it to your local `.git/info/exclude` if needed.
- Avoid introducing network calls; this is an offline TUI.

## Agent‑Specific Notes
- Keep changes minimal and focused; do not rename files or add new packages without need.
- Respect this guide across the entire repository tree.
