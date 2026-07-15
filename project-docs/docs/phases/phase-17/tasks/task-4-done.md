# Phase 17 - Task 4: Implement Maintenance Warnings and Info - Completion Report

## Summary of what was implemented

- Extended the maintenance domain result with deterministic warning findings and the documented `info` counts.
- Added the four v0.1 warning families: stale claims, started-without-claim, stale verifications, and completion-level telemetry gaps.
- Added injected audit time and configured started-without-claim grace handling; finding classification remains separate from Task 5's status and exit-code promotion.

## Files created or modified

- `internal/maintenance/check.go`
- `internal/maintenance/warnings.go`
- `internal/maintenance/info.go`
- `internal/maintenance/warnings_test.go`
- `internal/maintenance/info_test.go`
- `internal/maintenance/test_helpers_test.go`
- `project-docs/docs/phases/phase-17/tasks/task-4-done.md`

## Documentation updated

- None required; DEC-050/051 and the canonical maintenance output contract already specify the app-layer warning semantics and info fields.

## Validation Contract Coverage

- VAL-007: PASS — injected-time fixtures include a live claim and a claim exactly at TTL expiry; only the expiry-boundary claim is classified as `stale_claims`.
- VAL-008: PASS — fixtures prove live-claim exclusion, before/at/after grace boundaries, creation-time fallback, non-status update preservation, and zero grace.
- VAL-009: PASS — only the verification tied to a `superseded` completion is classified; current and NULL completion references are excluded.
- VAL-010: PASS — missing and issue-only runs remain `telemetry_gaps`; an `agent_runs.completion_id` link suppresses the warning.
- VAL-011: PASS — exact `issue_count`, `claim_count_live`, and `recent_activity_7d` values are asserted, including seven-day boundary handling and deep equality across two identical audits.

## Tests run

- `go test ./internal/maintenance -run '^TestCheckWarningsClassifiesQualifyingAndNonQualifyingStates$' -count=1` - expected RED: warning fixture/options/result fields did not exist.
- `go test ./internal/maintenance -run '^(TestCheckWarningsClassifiesQualifyingAndNonQualifyingStates|TestCheckWarningsIncludesStartedIssueAtZeroGrace|TestCheckInfoCountsAndResultAreDeterministic)$' -count=1` - PASS after implementation.
- `go test -race ./internal/maintenance` - PASS.
- Runtime entrypoint - N/A: `fid maintenance check` is explicitly Task 5 scope; this task exposes and tests the maintenance domain API only.

## Exact quality-gate results

- `gofmt -l cmd internal` - PASS
- `golangci-lint run ./...` - PASS (`0 issues.`)
- `go test -race -coverprofile=coverage.out ./...` - PASS (total statement coverage: 76.5%; `internal/journal`: 81.0%)
- `go tool cover -func=coverage.out` - PASS (total statement coverage: 76.5%)
- Conditional commands - N/A (the task specification lists none).

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
