# Phase 20 - Task 3a: Reason-code schema and redaction evidence - Completion Report

## Summary of what was implemented

- Replaced the flattened reason-code check with one exhaustive table-driven matrix for all 23 closed codes, including categories, retryability, subjects, scalar/null/nested/object-list details, serialized failure shape, and closed-detail redaction.
- Enforced the DEC-130-authorized equal-length rule for `duplicate-entity.identity.fields` and `identity.values` after the new matrix exposed its acceptance of mismatched arrays.

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
- `go test -race ./internal/cli -count=1` - PASS - `ok`
- Real entrypoint check: not required. This recovery task adds no command migration or new runtime surface; the changed registry is exercised through its public `NewFailure`/`WriteFailure` boundary by the exhaustive matrix.

## Exact quality-gate results

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)`
- `golangci-lint run ./...` - PASS - `0 issues.`
- `go test -race -coverprofile=coverage.out ./...` - PASS - `exit 0 (no output)`
- `go tool cover -func=coverage.out` - PASS - `77.4%`
- `go test -tags integration -race ./...` - PASS - `exit 0 (no output)`
- `go build -o /tmp/fid ./cmd/fid` - PASS - `exit 0 (no output)`

## Deviations from spec

- None.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `51763124db2cd6f6c6d00d15548de37532bbc273`

## Commit subject

- `fix(cli): enforce closed reason-code details`
