## Summary of what was implemented

- Added the in-code Phase 10 permissions matrix in `internal/policy/matrix.go` for the implemented v0.1a-v0.1c write actions only, matching `project-docs/specs/03_TRD_SCHEMA.md` section 13.
- Added `policy.Authorize(workspaceRoot, role, action string) error` and `PermissionDeniedError` in `internal/policy/authorize.go`, preserving `MissingRoleError` and `UnknownRoleError` while denying disallowed role/action combinations from the matrix.
- Added the read-only `fid policy show` root command in `cmd/fid/policy.go`, registered it in `cmd/fid/root.go`, and reused the shared format engine for `json`, `toon`, `markdown`, and `yaml`.
- Added Task 2 tests covering matrix completeness, `Authorize` permit/deny behavior, preserved role-validation errors, `fid policy show` format dispatch, action coverage, and read-only/no-`--by` behavior.

## Files created or modified

- `internal/policy/matrix.go`
- `internal/policy/authorize.go`
- `internal/policy/matrix_test.go`
- `internal/policy/authorize_test.go`
- `cmd/fid/policy.go`
- `cmd/fid/policy_test.go`
- `cmd/fid/root.go`
- `internal/format/markdown.go`
- `project-docs/docs/phases/phase-10/tasks/task-2-done.md`

## Documentation updated

- `project-docs/specs/06_DESIGN_PLAN_CLI.md` section 11 — verified-no-change; the documented JSON row shape already matched the implemented `fid policy show` output.
- `project-docs/docs/phases/phase-10/decisions.md` — verified-no-change; DEC-A, DEC-B, and DEC-C already matched the implemented seam, shared matrix, and read-only command behavior.

## Validation Contract Coverage

- `VAL-002` — covered by `TestPolicyShowDefaultsToJSONAndDispatchesFormats` and `TestPolicyShowIsReadOnlyAndDoesNotAcceptByFlag`, proving `fid policy show` renders the shared matrix in `json`/`toon`/`markdown`/`yaml`, includes every implemented write action from the shared matrix, stays read-only, and rejects `--by`.
- `VAL-003` — covered by `TestMatrixContainsExactlyTheImplementedV01WriteActions`, `TestAuthorizePermitsAllowedCombinationAndRejectsDeniedCombination`, and `TestAuthorizePreservesMissingAndUnknownRoleErrors`, proving every implemented action has exactly one matrix entry and `Authorize` preserves role validation while denying disallowed matrix cells with `PermissionDeniedError`.

## Tests run

- `go test ./internal/policy -run TestAuthorizePermitsAllowedCombinationAndRejectsDeniedCombination` — FAIL (expected red: `Authorize` and `PermissionDeniedError` undefined)
- `go test ./internal/policy -run TestAuthorizePermitsAllowedCombinationAndRejectsDeniedCombination` — PASS
- `go test ./internal/policy` — PASS
- `go test ./cmd/fid -run TestPolicyShowDefaultsToJSONAndDispatchesFormats` — FAIL (expected red: unknown command `policy`)
- `go test ./cmd/fid -run TestPolicyShowDefaultsToJSONAndDispatchesFormats` — PASS
- `go test ./cmd/fid -run 'TestPolicyShow'` — FAIL (expected red: `--by` rejection assertion needed to match Cobra unknown-flag behavior)
- `go test ./cmd/fid -run 'TestPolicyShow'` — PASS

## Exact quality-gate results

- `gofmt -l cmd internal` — PASS
- `golangci-lint run ./...` — PASS
- `go test -race -coverprofile=coverage.out ./...` — PASS
- `go tool cover -func=coverage.out` — PASS; COV-001: `internal/issues` 81.8%, `internal/journal` 80.3%, `internal/policy` 81.2%, `internal/workflow` 95.5%. COV-002: `cmd/fid` 77.8%.
- `go build -o /tmp/fid ./cmd/fid` — PASS
- `go test -tags integration -race ./...` — N/A; Task 2 added no integration-tagged tests.

## Deviations from spec

- None.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `6f4ce6b03b823b40f7f3b4ea7d6c0f3a59167df3`

## Commit subject

- `feat(policy): add permissions matrix and show command`
