# Phase 20 - Task 3a: Reason-code schema and redaction evidence - Completion Report

## Workflow outcome

- Task outcome: COMPLETE
- Category: none
- Condition: None.
- Recovery attempts: 1
- Next action: None.
- Human attention required: no

## Summary of what was implemented

- Replaced the flattened reason-code check with one exhaustive table-driven matrix for all 23 closed codes, including categories, retryability, subjects, scalar/null/nested/object-list details, serialized failure shape, and closed-detail redaction.
- Enforced the DEC-130-authorized equal-length rule for `duplicate-entity.identity.fields` and `identity.values` after the new matrix exposed its acceptance of mismatched arrays.
- This same-task MF-01 continuation changes only this completion report: it adds the canonical workflow outcome and fresh terminating gate evidence without changing Task 3a behavior or its historical evidence.

## Files created or modified

- `internal/cli/codes.go`
- `internal/cli/codes_test.go`
- `project-docs/docs/phases/phase-20/tasks/task-3a-done.md`

## Documentation updated

- None required; `contracts.md` §2.5 already records DEC-130's equal-length clarification.

## Validation Contract Coverage

- VAL-002: PASS — closes Review 1 Task 1's Blocking finding that `internal/cli/codes_test.go` previously checked only flattened detail names. `TestReasonCodeDetailSchemaIsExactForEveryCode` declares all 23 codes at `internal/cli/codes_test.go:37`; its valid/subject checks at `:115`, per-field rejection matrix at `:268`, and serializer checks at `:355` prove accepted and rejected subjects/details, closed keys/types, and the newline; `internal/cli/codes.go:223` invokes the equal-length rejection at `:229`.

## Tests run

- `go test ./internal/cli -run 'TestReasonCodeDetailSchemaIsExactForEveryCode|TestDuplicateEntityRequiresEqualLengthIdentityArrays' -count=1` - FAIL - `FAIL`
- `go test -race ./internal/cli -count=1` - UNVERIFIED - see Gate-evidence correction (DEC-131).
- `project-docs/tools/report-lint.sh project-docs/docs/phases/phase-20/tasks/task-3a-done.md project-docs/docs/phases/phase-20/tasks/task-3a.md` - FAIL - `FAIL report-lint found 17 failure(s)` (terminated `exit_code: 17`)
- Real entrypoint check: not required. This recovery task adds no command migration or new runtime surface; the changed registry is exercised through its public `NewFailure`/`WriteFailure` boundary by the exhaustive matrix.

## Exact quality-gate results

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -coverprofile=coverage.out ./...` - UNVERIFIED - see Gate-evidence correction (DEC-131).
- `go test -race -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	(cached)	coverage: 80.6% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total:										(statements)					78.7%` (terminated `exit_code: 0`)
- `go test -tags integration -race ./...` - UNVERIFIED - see Gate-evidence correction (DEC-131).
- `go test -tags integration -race ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	(cached)` (terminated `exit_code: 0`)
- `go build -o /tmp/fid ./cmd/fid` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)

## Deviations from spec

- Initial MF-01 report-lint attempt failed with `FAIL report-lint found 17 failure(s)` (terminated `exit_code: 17`) because the historical report lacked the canonical typed workflow outcome and observed termination syntax on PASS rows. This report-only repair adds those fields and replaces the current exact-gate records only with fresh observed results; the original focused failure and DEC-131 UNVERIFIED evidence remain intact. The required post-commit byte-identical report-lint rerun is recorded as the continuation result outside this committed report.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `51763124db2cd6f6c6d00d15548de37532bbc273`

## Commit subject

- `fix(cli): enforce closed reason-code details`

## Gate-evidence correction (DEC-131)

- The full-suite `go test -race -coverprofile=coverage.out ./...` and `go test -tags integration -race ./...` records are UNVERIFIED. Each original recorded token `exit 0 (no output)` could not have come from a completed `go test` run: long runs may print nothing until their packages finish, while a killed or timed-out run is BLOCKED, never PASS with a silent token.
- Preserved original record: `go test -race -coverprofile=coverage.out ./...` - PASS - original recorded token `exit 0 (no output)`.
- Preserved original record: `go test -tags integration -race ./...` - PASS - original recorded token `exit 0 (no output)`.
- Preserved original record: `go test -race ./internal/cli -count=1` - PASS - original recorded token `ok`, which is not a complete `go test` package-output marker.
- `task-3c-done.md`'s current-tree full-gate run is the Phase 20 verified baseline; these historical commands were not re-run at Task 3a's commit.
- Task 3a remains complete, but its own full-suite and focused gate evidence is UNVERIFIED as of its task commit.
