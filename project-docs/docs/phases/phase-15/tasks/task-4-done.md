## Summary of what was implemented

- Added `cmd/fid/e2e_v01e_dev_test.go`, one integration-tagged end-to-end smoke that:
  drives `fid dev start -> gate -> finish` and `fid dev start -> block` against a real temp `.fid` workspace,
  proves a `qa_reviewer` composite denial appends zero journal lines,
  and deletes/rebuilds the projection to verify deterministic parity across issues, claims, completions, runs, and gate runs.
- Kept composite product behavior unchanged by consuming the existing CLI surface and asserting journal-line deltas plus projected state only.
- Updated `project-docs/specs/07_DESIGN_PLAN_FLOWS.md` because the documented developer flow had drifted from the implemented CLI (`dev start --model`, `dev gate --exit-code`, and finish output including claim release).

## Files created or modified

- `cmd/fid/e2e_v01e_dev_test.go` - new integration smoke and hermetic snapshot helpers for the composite lifecycle and rebuild parity.
- `project-docs/specs/07_DESIGN_PLAN_FLOWS.md` - reconciled the developer execution example with the implemented composite flags/output.
- `project-docs/docs/phases/phase-15/tasks/task-4-done.md` - this completion report.

## Documentation updated

- Updated `project-docs/specs/07_DESIGN_PLAN_FLOWS.md` to match the implemented `fid dev start`, `fid dev gate`, and `fid dev finish` examples.

## Validation Contract Coverage

- `VAL-011` PASS.
  Evidence: `TestV01EDevCompositeSmokeAndRebuildParity` drives both composite paths over a real temp workspace, asserts one journal-line increment per composite, checks projected claim/issue/run/gate/completion state, proves `qa_reviewer` denial before append, deletes `.fid/fid.sqlite`, runs `fid rebuild`, and confirms the projected snapshot is unchanged.
- `VAL-012` prep only.
  Evidence: the new e2e smoke and the reconciled flow example give QA Lead the deterministic CLI path needed for the phase-end human dogfood run; no human execution was performed in this task.

## Tests run

- Red: `go test -tags integration -race ./cmd/fid -run TestV01EDevCompositeSmokeAndRebuildParity`
  Result: FAIL during initial test-first pass because `assertV01EPlaceholder` was undefined in the new test file.
- Green: `go test -tags integration -race ./cmd/fid -run TestV01EDevCompositeSmokeAndRebuildParity`
  Result: PASS (`ok  	github.com/mbell-nzl/fid/cmd/fid	2.619s`, then cached after formatting).
- Helpful size check: `wc -l cmd/fid/e2e_v01e_dev_test.go project-docs/specs/07_DESIGN_PLAN_FLOWS.md`
  Result: PASS (`746` lines in the new integration test file; `283` lines in the flow doc).
- Conditional benchmark command: not run.
  Reason: no benchmark was added or modified, so `go test -bench . -benchmem ./...` was not applicable for this task.

## Exact quality-gate results

- `gofmt -l cmd internal`
  PASS. No output.
- `golangci-lint run ./...`
  PASS. Output excerpt: `0 issues.`
- `go test -race -coverprofile=coverage.out ./...`
  PASS. Output excerpt: `ok github.com/mbell-nzl/fid/cmd/fid 49.875s coverage: 76.9% of statements`
- `go tool cover -func=coverage.out`
  PASS. Output excerpt: `github.com/mbell-nzl/fid/cmd/fid/dev.go:11: newDevCommand 100.0%` ... `total: (statements) 75.8%`
- `go test -tags integration -race ./...`
  PASS. Output excerpt: `ok github.com/mbell-nzl/fid/cmd/fid 59.956s`
- `go build -o /tmp/fid ./cmd/fid`
  PASS. No output.
- `cmd/fid` package coverage
  PASS. `76.9%` from `go test -race -coverprofile=coverage.out ./...`.

## Deviations from spec

- None.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `64cd9ec`

## Commit subject

- `test(dev): add phase 15 composite smoke`
