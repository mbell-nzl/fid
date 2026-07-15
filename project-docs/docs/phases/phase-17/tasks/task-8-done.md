# Phase 17 - Task 8: Prove v0.1 Read-Format Conformance - Completion Report

## Summary of what was implemented

- Added one integration-tagged, hermetic conformance matrix for the exact 15 implemented v0.1 read commands.
- Seeded the matrix workspace exclusively through shipped CLI writes, then compared JSON, TOON, and YAML structures while checking markdown primary identities.
- Confirmed every matrix row already dispatches all four declared formats; no production or canonical CLI-spec change was needed.

## Files created or modified

- `cmd/fid/e2e_v01_format_matrix_test.go`
- `project-docs/docs/phases/phase-17/tasks/task-8-done.md`

## Documentation updated

- No canonical CLI-spec change required: `project-docs/specs/06_DESIGN_PLAN_CLI.md` already declares json/toon/markdown/yaml for every shipped matrix command.

## Validation Contract Coverage

- VAL-017: PASS. `TestV01ReadCommandsSupportAllFormatsWithStructuralParity` runs the exact v0.1 read matrix in a real temporary workspace. It deep-compares JSON, TOON, and YAML output per command-specific shape and verifies markdown contains the same seeded primary IDs, labels, actions, or status identity.

## Tests run

- `go test -tags integration -race ./cmd/fid -run '^TestV01ReadCommandsSupportAllFormatsWithStructuralParity$' -count=1` - expected RED: the initial complete matrix reached `journal inspect` and exposed an inaccurate fixture assumption about the latest journal event; no product format behavior was changed.
- `go test -tags integration -race ./cmd/fid -run '^TestV01ReadCommandsSupportAllFormatsWithStructuralParity$' -count=1` - PASS after the deterministic fixture identities were corrected.
- `/tmp/fid` real-entrypoint smoke - PASS: initialized a disposable workspace, created a task, rendered `issue list --format yaml`, and rendered `ready --format markdown`; both outputs contained the created task identity and expected fields.

## Exact quality-gate results

- `gofmt -l cmd internal` - PASS (no files listed).
- `golangci-lint run ./...` - PASS (`0 issues.`).
- `go test -race -coverprofile=coverage.out ./...` - PASS.
- `go tool cover -func=coverage.out` - PASS (total statements: 76.6%).
- `go test -tags integration -race ./...` - PASS.
- `go build -o /tmp/fid ./cmd/fid` - PASS.

## Deviations from spec

- None.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending commit

## Commit subject

- Pending commit
