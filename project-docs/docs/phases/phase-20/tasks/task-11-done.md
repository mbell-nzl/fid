# Phase 20 - Task 11: Link contract and dependency invariants - Completion Report

## Workflow outcome

- Task outcome: COMPLETE
- Category: none
- Condition: None.
- Recovery attempts: 0
- Next action: None.
- Human attention required: no

## Summary of what was implemented

- Migrated `link add`, `link remove`, and `link list` to the canonical typed-endpoint request and Response Envelope contract, including resolved project-root use and stdout-only Human successes.
- Added typed Link errors, exact JSON reason mapping, sorted affected Issue revisions, replayable revision Effects, and bounded `blocks`-direction cycle detection.
- Updated impacted in-process and real-process journeys from retired endpoint inference/`--format` calls to the canonical flags and JSON envelope.

## Files created or modified

- `internal/links/add.go`
- `internal/links/remove.go`
- `internal/links/list.go`
- `internal/links/types.go`
- `internal/links/errors.go`
- `internal/links/registry.go`
- `internal/links/cycle.go`
- `internal/links/revision.go`
- `internal/links/cycle_test.go`
- `internal/links/add_activity_test.go`
- `internal/links/add_test.go`
- `internal/links/registry_test.go`
- `cmd/fid/link.go`
- `cmd/fid/link_response.go`
- `cmd/fid/link_test.go`
- `cmd/fid/e2e_v01d_closure_test.go`
- `cmd/fid/e2e_v01d_link_test.go`
- `cmd/fid/e2e_v01e_activity_test.go`
- `cmd/fid/integration/e2e_v01_closure_helpers_test.go`
- `cmd/fid/integration/e2e_v01_closure_test.go`

## Documentation updated

- None required.

## Validation Contract Coverage

- VAL-008: PASS - `cmd/fid/link_test.go:27` asserts Human add success has empty stderr and its confirmation on stdout; the built binary emitted canonical JSON only on stdout during the recorded real-entrypoint probe.
- VAL-010: PASS - `internal/links/add_activity_test.go:48` requires two added Issue revision Effects and `:51` requires both resulting revisions; `/tmp/fid` emitted sorted `issue_revisions` at revision 2 for both endpoints.
- VAL-011: PASS - `cmd/fid/link_test.go:66` requires the five-effect Link Event shape, including the two Issue revision updates at `:78`; the race/coverage suite passed this assertion and replay journey coverage.
- VAL-017: PASS - the real-entrypoint `/tmp/fid --project /tmp/fid-task11.Q4XJP7 --json link add ...` result contained Event `EVT-20260815-000004`; `cmd/fid/link_test.go:465` requires read `link.list` Event `null`.
- VAL-023: PASS - `cmd/fid/link_test.go:27` verifies typed add plus inverse rows, and `cmd/fid/link_test.go:446` verifies the closed typed `link.list` response/count/order.
- VAL-024: PASS - `internal/links/cycle_test.go:8` rejects the three-node closing `blocks` link with the ordered path at `:41` and confirms zero journal append at `:44`.
- VAL-025: PASS - `cmd/fid/link_test.go:343` retains missing-link/self-link rejection-before-append coverage, exercised by the passing `go test -race -coverprofile=coverage.out ./...` gate.
- VAL-041: PASS - `cmd/fid/link_test.go:92` verifies the Link Event activity attribution is `developer`; no protected/trusted effect is introduced by the five-effect Event assertion at `:66`.
- VAL-044: PASS - `cmd/fid/link.go:29`, `:95`, and `:161` require the exact endpoint type flags/list type before handler work; `cmd/fid/link_test.go:27` and `:460` exercise the accepted requests.
- VAL-045: PASS - `internal/links/list.go` orders by direction, relation, entity type, and ID with an explicit 10,000-item limit; `cmd/fid/link_test.go:478` asserts the resulting canonical order.
- VAL-046: PASS - `cmd/fid/link.go:29`/`:33` require add endpoint types, `cmd/fid/link.go:95`/`:99` require remove endpoint types, and `cmd/fid/link.go:161` requires list type; `cmd/fid/link_response.go:56` maps missing endpoints distinctly to `issue-not-found` and `requirement-not-found`.

## Tests run

- RED (first required test): `go test ./internal/links -run '^TestAddRejectsThreeNodeDependencyCycle$' -count=1` - `FAIL` (terminated `exit_code: 1`). Initial compile-red named the missing `DependencyCycleError`; after declaring the type, the behavior-red was `error = <nil> <nil>, want *DependencyCycleError` (terminated `exit_code: 1`).
- GREEN: `go test ./internal/links -run '^TestAddRejectsThreeNodeDependencyCycle$' -count=1` - `ok  	github.com/mbell-nzl/fid/internal/links	0.069s` (terminated `exit_code: 0`).
- GREEN: `go test ./internal/links -count=1` - `ok  	github.com/mbell-nzl/fid/internal/links	0.624s` (terminated `exit_code: 0`).
- GREEN: `go test ./cmd/fid -count=1` - `ok  	github.com/mbell-nzl/fid/cmd/fid	6.393s` (terminated `exit_code: 0`).
- GREEN: `go test -tags integration ./cmd/fid/integration/...` - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	3.236s` (terminated `exit_code: 0`).
- Real entrypoint: `/tmp/fid --project /tmp/fid-task11.Q4XJP7 --json link add TASK11-6dd73e blocks TASK11-268187 --from-type issue --to-type issue --by developer` - `{"contract_version":"fid-cli.v1","ok":true,"operation":"link.add",...}` (terminated `exit_code: 0`).

## Exact quality-gate results

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`).
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`).
- `go test -race -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	(cached)	coverage: 80.5% of statements` (terminated `exit_code: 0`).
- `go tool cover -func=coverage.out` - PASS - `total:										(statements)					77.3%` (terminated `exit_code: 0`).
- `go test -tags integration -race ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	5.296s` (terminated `exit_code: 0`).
- `go build -o /tmp/fid ./cmd/fid` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`).

## Deviations from spec

- None.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `4155acee3037614762d0cfc0d6710f540b3f6441`

## Commit subject

- `feat(links): enforce typed link invariants`
