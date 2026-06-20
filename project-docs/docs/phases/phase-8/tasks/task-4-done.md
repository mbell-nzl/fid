# Task 4 Completion Report

## Summary of what was implemented

- Performed the CF-012 behavior-preserving split of `internal/issues/issue.go`.
- Moved the create path into `internal/issues/create.go`.
- Moved `Get` into `internal/issues/get.go`.
- Reduced `internal/issues/issue.go` to shared exported types/errors and shared package helpers.
- Preserved all exported symbols, function signatures, package boundaries, and behavior.

## Files created or modified

- `internal/issues/create.go` (new, 387 lines)
- `internal/issues/get.go` (new, 120 lines)
- `internal/issues/issue.go` (residual, 192 lines)
- `project-docs/docs/phases/phase-8/tasks/task-4-done.md`

## Documentation updated

- `project-docs/docs/phases/phase-8/tasks/task-4-done.md`

## Validation Contract Coverage

- `VAL-009` PASS.
- Evidence:
  - Same-package split completed with all three in-scope non-test files under 400 lines: `internal/issues/create.go` 387, `internal/issues/get.go` 120, `internal/issues/issue.go` 192.
  - Diff is function moves plus import cleanup only; no exported symbol was added, removed, or renamed.
  - Regression oracle stayed green before and after the split, and the full task gate plus conditional build passed.
  - `internal/issues` coverage remained above the DEC-007 threshold at 81.9%.

## Tests run

- `go test -race ./internal/issues` — PASS (pre-split baseline)
- `go test -race ./internal/issues` — PASS (post-split regression check)
- `go test -race -coverprofile=coverage.out ./...` — PASS

## Exact quality-gate results

- `gofmt -l cmd internal` — PASS
- `golangci-lint run ./...` — PASS
- `go test -race -coverprofile=coverage.out ./...` — PASS
- `go tool cover -func=coverage.out` — PASS
  - `internal/issues` coverage: 81.9%
- `go build -o /tmp/fid ./cmd/fid` — PASS

## Deviations from spec

- None.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `6bf1bf3`

## Commit subject

- `refactor: split issues issue file`
