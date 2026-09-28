# Phase 23 - Task 13: CLI surface, failure codes, and idempotency - Recovery Report

## Workflow outcome

- Task outcome: RECOVERY_REQUIRED
- Category: in-scope-implementation
- Condition: Accepted `review.list` and ten ActivityEvent names cannot be completed within task 13's named production file scope. The focused real-process probes remain red.
- Recovery attempts: 1
- Next action: Phase Runner obtains a bounded task-13 scope amendment for the owning Review list reader and Phase 23 Event producers, then dispatches a same-task continuation to finish VAL-055..057 and the exact gate.
- Human attention required: no

## Summary of what was implemented

- Added a built-binary request-closure matrix for fourteen Phase 23 mutations. It rejects unknown fields, a null request object, trailing JSON, and duplicate keys without appending an Event.
- Registered six previously omitted Phase 23 failure codes and aligned the exact-schema test. Renamed the unexported failure registry and test-local expectation type per DEC-226.
- Rendered `evidence.link` idempotency conflicts through the accepted failure envelope, verified by a built-binary same-key/different-digest probe.
- Preserved red real-process probes for the missing `review.list` leaf and the absent `review-recorded` ActivityEvent. The existing CLI and domain Event producer have no accepted in-scope route to add the required list query and activity effect.

## Files created or modified

- `cmd/fid/evidence.go`
- `cmd/fid/integration/p23_leaf_surface_test.go`
- `cmd/fid/integration/p23_review_evidence_budget_test.go`
- `cmd/fid/integration/p23_review_record_test.go`
- `internal/cli/codes.go`
- `internal/cli/codes_test.go`
- `internal/cli/envelope.go`
- `internal/cli/phase23_codes_test.go`
- `project-docs/specs/06_DESIGN_PLAN_CLI.md`
- `project-docs/docs/phases/phase-23/tasks/task-13-done.md`

## Documentation updated

- `project-docs/specs/06_DESIGN_PLAN_CLI.md` section 5 now lists the accepted Phase 23 failure codes, categories, exits, subjects, and closed detail keys.
- This recovery report records the remaining scope gap and observed evidence.

## Validation Contract Coverage

- VAL-055: BLOCKED - `TestPhase23MutationsRejectMalformedObjectsWithoutEvents` passes the four malformed-object forms across fourteen mutation leaves. `TestPhase23LeavesRejectUnknownRequestFields` still reaches CLI help for missing `review.list`; `TestP23ReviewRecordRealProcessJourney` finds zero `review-recorded` activity rows after a successful producing leaf. The ten-name ActivityEvent inventory remains incomplete.
- VAL-056: BLOCKED - `TestPhase23FailureCodesAreRegistered` and `TestReasonCodeDetailSchemaIsExactForEveryCode` pass the registered-code inventory and closed shapes; `TestP23ReviewEvidenceRealProcessJourney` proves a built-binary `idempotency-conflict`. The full producer matrix and both accepted precedence probes remain to be recorded.
- VAL-057: BLOCKED - `TestP23ReviewEvidenceRealProcessJourney` proves the conflict emits no Event. Replay of every mutation and both noncommitting Stop denials still needs a complete task-13 corpus.
- VAL-079: BLOCKED - the three DEC-227 package times cannot be claimed until the exact ordinary and tagged gates terminate on a completed tree. DEC-258/259 already resolved DEC-254 Part 2; no repeated headroom or count=2 step was run.

## Tests run

- `go test -tags integration -run '^TestPhase23LeavesRejectUnknownRequestFields$' -count=1 ./cmd/fid/integration` - FAIL - `FAIL` (terminated `exit_code: 1`)
- `go test -tags integration -run '^TestPhase23LeavesRejectUnknownRequestFields$' -count=1 ./cmd/fid/integration` - FAIL - `FAIL` (terminated `exit_code: 1`)
- `go test -tags integration -run '^TestPhase23LeavesRejectUnknownRequestFields$' -count=1 ./cmd/fid/integration` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	1.677s` (terminated `exit_code: 0`)
- `go test -run '^TestPhase23FailureCodesAreRegistered$' -count=1 ./internal/cli` - FAIL - `FAIL` (terminated `exit_code: 1`)
- `go test -run '^(TestPhase23FailureCodesAreRegistered|TestReasonCodeDetailSchemaIsExactForEveryCode)$' -count=1 ./internal/cli` - PASS - `ok  	github.com/mbell-nzl/fid/internal/cli	0.026s` (terminated `exit_code: 0`)
- `go test -tags integration -run '^TestP23ReviewEvidenceRealProcessJourney$' -count=1 ./cmd/fid/integration` - FAIL - `FAIL` (terminated `exit_code: 1`)
- `go test -tags integration -run '^TestP23ReviewEvidenceRealProcessJourney$' -count=1 ./cmd/fid/integration` - FAIL - `FAIL` (terminated `exit_code: 1`)
- `go test -tags integration -run '^TestP23ReviewEvidenceRealProcessJourney$' -count=1 ./cmd/fid/integration` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	4.522s` (terminated `exit_code: 0`)
- `go test -tags integration -run '^TestPhase23MutationsRejectMalformedObjectsWithoutEvents$' -count=1 ./cmd/fid/integration` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	2.123s` (terminated `exit_code: 0`)
- `go test -tags integration -run '^TestP23ReviewRecordRealProcessJourney$' -count=1 ./cmd/fid/integration` - FAIL - `FAIL` (terminated `exit_code: 1`)
- `go test -tags integration -run '^(TestPhase23LeavesRejectUnknownRequestFields|TestP23ReviewRecordRealProcessJourney)$' -count=1 ./cmd/fid/integration` - FAIL - `FAIL` (terminated `exit_code: 1`)
- `go test -tags integration -run '^(TestPhase23MutationsRejectMalformedObjectsWithoutEvents|TestP23ReviewEvidenceRealProcessJourney)$' -count=1 ./cmd/fid/integration` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	3.796s` (terminated `exit_code: 0`)

## Exact quality-gate results

Working directory: `/home/mark/projects/fid`. The two short exact commands ran in task-spec order. The remaining exact commands await the bounded source-scope amendment and a completed tree; none has a termination result in this attempt.

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` - BLOCKED - scope gap leaves VAL-055 probes red; command not run (no termination observed).
- `go tool cover -func=coverage.out` - BLOCKED - no task-13 coverage profile from a completed ordinary gate (no termination observed).
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...` - BLOCKED - scope gap leaves `review.list` and ActivityEvent probes red; command not run (no termination observed).

## Deviations from spec

- Task 13's accepted VAL-055 behavior needs `internal/review` list and Event-producer edits beyond its named CLI/adapter and idempotency-snapshot production scope. No ActivityEvent was invented and no gate or assertion was changed. The temporarily green list implementation was removed; its red test remains. This is a scope gap for the Phase Runner to resolve before same-task continuation.
- The first evidence conflict probe used an unchanged digest and correctly replayed. Changing a request field exposed the generic error adapter, which was repaired and rerun green. The first failing outputs remain recorded above.

## Concerns, technical debt, or notes for the Tech Lead

- Exact paths needing scoped authority include `internal/review/review.go` (Review list projection and `review.recorded` activity effect) and the owning producer files `internal/review/finding.go`, `concern.go`, `deviation.go`, `decision.go`, `human_outcome.go`, and `requirement.go` for the remaining accepted activity names. Their current Event effects omit `activity_log`; `internal/journal/current_event_catalog.go` already permits the optional effect. The accepted contract fixes the ten names and producers.
- DEC-227 package times remain unmeasured for this checkpoint; the full gates are pending.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending commit

## Commit subject

- Pending commit
