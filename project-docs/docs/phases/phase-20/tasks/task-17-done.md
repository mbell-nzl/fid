# Phase 20 - Task 17: Delete promised-command stubs - Completion Report

## Workflow outcome

- Task outcome: COMPLETE
- Category: none
- Condition: None.
- Recovery attempts: 0
- Next action: None.
- Human attention required: no

## Summary of what was implemented

- Deleted the decision, escalation, and triage governance-stub adapter and its obsolete version-gate test suite.
- Removed the three root registrations, so Cobra now treats each former family root as an unknown command; no tombstone or replacement family was added.
- Added a root-tree guard for all three unknown command roots and updated the real-process closure journey to retain its maintenance and journal-immutability assertions under the new absence behavior.

## Files created or modified

- Deleted `cmd/fid/governance_stubs.go`.
- Deleted `cmd/fid/governance_stubs_test.go`.
- `cmd/fid/root.go`
- `cmd/fid/root_test.go`
- `cmd/fid/integration/e2e_v01_closure_helpers_test.go`
- `project-docs/docs/phases/phase-20/tasks/task-17-done.md`
- Scope check: the pre-existing modified `project-docs/docs/phases/phase-20/run-ledger.md` was left unstaged and unmodified by this task.

## Documentation updated

- None required.

## Validation Contract Coverage

- VAL-034: PASS - `TestPromisedGovernanceCommandsAreAbsent` at `cmd/fid/root_test.go:9` proves `fid decision`, `fid escalation`, and `fid triage` are Cobra unknown commands with no output. The final declared absence command returned `exit 1 (no output)` for every stub constructor, `versionGateError`, and `ships in v0.2` promise. The real-process closure journey's `assertClosureMaintenanceAndReadCommands` at `cmd/fid/integration/e2e_v01_closure_helpers_test.go:278` still runs maintenance check/vacuum, verifies journal bytes after each, and verifies the unknown `decision` command leaves that journal unchanged.

## Tests run

- `go test ./cmd/fid -run '^TestPromisedGovernanceCommandsAreAbsent$' -count=1` - EXPECTED RED - `FAIL` (terminated `exit_code: 1`); before removal, all three registered roots returned nil instead of Cobra unknown-command errors.
- `go test ./cmd/fid -run '^TestPromisedGovernanceCommandsAreAbsent$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	0.047s` (terminated `exit_code: 0`)
- `go test -tags integration -race ./cmd/fid/integration/... -run '^TestV01ClosureJourneyAndRebuildParity$'` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	3.569s` (terminated `exit_code: 0`)
- `/tmp/fid decision` - EXPECTED REJECTION - `exit 1 (no output)` (terminated `exit_code: 1`)
  Real-process entrypoint evidence: the built CLI rejects the removed root with no tombstone output.

## Exact quality-gate results

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	(cached)	coverage: 80.6% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total:										(statements)					78.7%` (terminated `exit_code: 0`)
- `go test -tags integration -race ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	4.435s` (terminated `exit_code: 0`)
- `go build -o /tmp/fid ./cmd/fid` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `rg -n "newDecisionCommand|newEscalationCommand|newTriageCommand|versionGateError|ships in v0\.2" cmd/fid` - PASS - `exit 1 (no output)` (expected `exit_code: 1`; terminated `exit_code: 1`)

## Deviations from spec

- Required integration-gate first attempt: `go test -tags integration -race ./cmd/fid/integration/...` - FAIL - `FAIL` (terminated `exit_code: 1`). Root cause: `cmd/fid/integration/e2e_v01_closure_helpers_test.go` still asserted the removed decision version-gate stderr. Causal repair: preserved its real-process maintenance and journal-immutability assertions while changing only the now-invalidated expectation to a non-zero unknown command with empty streams. The byte-identical command later passed with `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	4.435s` (terminated `exit_code: 0`).
- No product, architecture, data-model, public-interface, dependency, command, threshold, or assertion-strength deviation.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending commit

## Commit subject

- Pending commit
