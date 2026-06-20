## Summary of what was implemented

Implemented Phase 9 Task 2: `fid concern address` plus the `concerns_addressed`
rule proof.

- Added the `internal/issues.AddressConcern` write path as an in-place
  `concerns` row update from `open` to `addressed`, recording the addressing
  completion and updated timestamp under the journal write lock.
- Added the `fid concern address <concern-id> --in-completion <completion-id> --by <role>`
  CLI command with validate-before-append required-flag and domain-error
  handling.
- Proved the already-shipped derived flag clearing and the DEC-A rule seam:
  partial addressing leaves one open concern visible, addressing the last open
  concern clears the flag, and `evaluateTransitionRule(RuleConcernsAddressed)`
  reports unsatisfied named `CON-` IDs until all concerns are addressed.

## Files created or modified

- `cmd/fid/concern.go`
- `cmd/fid/concern_test.go`
- `internal/issues/concern_address.go`
- `internal/issues/concern_address_test.go`
- `internal/issues/transition_rules_test.go`
- `project-docs/specs/06_DESIGN_PLAN_CLI.md`
- `project-docs/docs/phases/phase-9/tasks/task-2-done.md`

## Documentation updated

- Updated `project-docs/specs/06_DESIGN_PLAN_CLI.md` in the `fid concern address`
  section to add the required argument/flag table and validate-before-append
  rejection rules.
- Added this completion report at
  `project-docs/docs/phases/phase-9/tasks/task-2-done.md`.

## Validation Contract Coverage

- `VAL-004` PASS.
  Evidence:
  `TestAddressConcernUpdatesConcernInPlaceAndClearsDerivedFlags`
  `TestConcernAddressUpdatesConcernAndPrintsConfirmation`
  These prove the same concern row is retained, `status` changes to
  `addressed`, `addressed_in_completion_id` is stored, and the command confirms
  the addressing completion.
- `VAL-005` PASS.
  Evidence:
  `TestAddressConcernRejectsBeforeAppend`
  `TestConcernAddressRequiresCompletionFlag`
  These prove unknown role, nonexistent concern, already-addressed concern,
  nonexistent completion, and cross-issue completion rejections happen before
  any journal append, with required-flag rejection surfaced at the CLI.
- `VAL-006` PASS.
  Evidence:
  `TestAddressConcernUpdatesConcernInPlaceAndClearsDerivedFlags`
  This proves the partial state stays `concerns_raised=true` with one remaining
  open concern on both `Get` and `List`, then clears to `false` and `0` after
  the last concern is addressed.
- `VAL-007` PASS.
  Evidence:
  `TestEvaluateTransitionRuleConcernsAddressedUsesRealConcernWrites`
  This integration-tagged test drives real concern raise/address writes through
  `evaluateTransitionRule` and proves unsatisfied named `CON-` IDs while open,
  then satisfied after all concerns are addressed. No `fid issue transition`
  edge was added.

## Tests run

- `go test -race ./internal/issues -run 'TestAddressConcernUpdatesConcernInPlaceAndClearsDerivedFlags|TestAddressConcernRejectsBeforeAppend'` PASS
- `go test -race ./cmd/fid -run 'TestConcernAddressUpdatesConcernAndPrintsConfirmation|TestConcernAddressRequiresCompletionFlag'` PASS
- `go test -tags integration -race ./internal/issues -run TestEvaluateTransitionRuleConcernsAddressedUsesRealConcernWrites` PASS
- `gofmt -l cmd internal` PASS
- `golangci-lint run ./...` PASS
- `go test -race -coverprofile=coverage.out ./...` PASS
- `go tool cover -func=coverage.out` PASS
- `go build -o /tmp/fid ./cmd/fid` PASS
- `go test -tags integration -race ./...` PASS

## Exact quality-gate results

- `gofmt -l cmd internal` PASS
  Output:
  no files listed.
- `golangci-lint run ./...` PASS
  Output:
  `0 issues.`
- `go test -race -coverprofile=coverage.out ./...` PASS
  Coverage highlights:
  `cmd/fid` 78.2%
  `internal/issues` 81.8%
  `internal/journal` 80.3%
  `internal/policy` 85.0%
  `internal/workflow` 95.5%
- `go tool cover -func=coverage.out` PASS
  Overall statement coverage:
  77.0%
- `go build -o /tmp/fid ./cmd/fid` PASS
- `go test -tags integration -race ./...` PASS
  Package highlights:
  `cmd/fid` PASS
  `internal/issues` PASS

## Deviations from spec

- None.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `3ea2d4f`

## Commit subject

- `feat(concern): add concern address command`
