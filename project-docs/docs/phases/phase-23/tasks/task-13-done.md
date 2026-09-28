# Phase 23 - Task 13: CLI surface, failure codes, and idempotency - Blocker Report

## Workflow outcome

- Task outcome: BLOCKED
- Category: specification-contradiction
- Condition: VAL-055 requires all ten ActivityEvent names from their producing leaves at task 13, but the current strict journal table registry rejects `activity_log/activity` with `unknown table`, and the current database trigger cannot attach project-level `requirement.record-reference` activity to its new `requirement_refs` row. The accepted trigger change is deferred to migration 0018 in task 15.
- Recovery attempts: 1
- Next action: Phase Runner routes the `activity_log` registry/trigger and task-13-before-task-15 conflict to TL-Lite for protected resolution, then resumes task 13 with an exact bounded directive.
- Human attention required: yes

## Summary of what was implemented

- Added a built-binary request-closure matrix for fourteen Phase 23 mutations. It rejects unknown fields, a null request object, trailing JSON, and duplicate keys without appending an Event.
- Registered six previously omitted Phase 23 failure codes and aligned the exact-schema test. Renamed the unexported failure registry and test-local expectation type per DEC-226.
- Rendered `evidence.link` idempotency conflicts through the accepted failure envelope, verified by a built-binary same-key/different-digest probe.
- At the prior checkpoint, preserved red real-process probes for the missing `review.list` leaf and absent `review-recorded` ActivityEvent. The Runner's TASK-RECOVERY attempt 1 subsequently authorized the bounded Review reader and producer files.
- Added `review.list` with its accepted four filters, default and range limit, complete bounded Review projection, and typed `result-limit-exceeded`. The built-binary list and leaf-inventory probes now pass.
- Tried `review-recorded` and `finding-resolved` effects in the owning Review producer. The first producing call failed at the journal validator with `unknown table`; the trial code was removed so existing Review writes remain available. The failed observation is retained below.

## Files created or modified

- `cmd/fid/evidence.go`
- `cmd/fid/contract_surface_test.go`
- `cmd/fid/review.go`
- `cmd/fid/integration/p23_leaf_surface_test.go`
- `cmd/fid/integration/p23_review_evidence_budget_test.go`
- `cmd/fid/integration/p23_review_record_test.go`
- `internal/cli/codes.go`
- `internal/cli/codes_test.go`
- `internal/cli/envelope.go`
- `internal/cli/phase23_codes_test.go`
- `internal/review/review.go`
- `project-docs/specs/06_DESIGN_PLAN_CLI.md`
- `project-docs/docs/phases/phase-23/tasks/task-13-done.md`

## Documentation updated

- `project-docs/specs/06_DESIGN_PLAN_CLI.md` section 5 now lists the accepted Phase 23 failure codes, categories, exits, subjects, and closed detail keys.
- This report records the resolved list gap and the protected ActivityEvent/schema contradiction with observed evidence.

## Validation Contract Coverage

- VAL-055: BLOCKED - `TestPhase23MutationsRejectMalformedObjectsWithoutEvents` passes the four malformed-object forms across fourteen mutation leaves. `TestPhase23LeavesRejectUnknownRequestFields` now calls the built binary's `review list` successfully (`cmd/fid/integration/p23_leaf_surface_test.go:22`); `TestP23ReviewRecordRealProcessJourney` still finds zero `review-recorded` rows. A scoped `activity_log/activity` trial made the producing call fail with `unknown table` from `internal/journal/current_event_validation.go:161`, because `internal/journal/current_event_columns.go` has no `activity_log` entry. `migrations/0016_p22_remove_escalations.up.sql:9` resolves requirement activity references only against the V1 `requirements` table, while `requirement.record-reference` writes `requirement_refs`; `schema.md` section 6.2 defers the new trigger to 0018. The ten-name inventory cannot pass at this boundary.
- VAL-056: BLOCKED - `TestPhase23FailureCodesAreRegistered` and `TestReasonCodeDetailSchemaIsExactForEveryCode` pass the registered-code inventory and closed shapes; `TestP23ReviewEvidenceRealProcessJourney` proves a built-binary `idempotency-conflict`. The full producer matrix and both accepted precedence probes remain to be recorded.
- VAL-057: BLOCKED - `TestP23ReviewEvidenceRealProcessJourney` proves the conflict emits no Event. Replay of every mutation and both noncommitting Stop denials still needs a complete task-13 corpus.
- VAL-079: PASS for observed DEC-227 timings on this checkpoint - exact ordinary gate reported `cmd/fid 288.696s` and `internal/workflow 116.083s` (both below 600s); exact tagged gate reported `cmd/fid/integration 258.375s` (below 900s). Tagged gate fails VAL-055, not timing. DEC-258/259 already resolved DEC-254 Part 2; no repeated headroom or count=2 step was run. `qa/esc-23-005-repair.md` is prior evidence only.

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

- `go test -tags integration -run '^(TestPhase23LeavesRejectUnknownRequestFields|TestP23ReviewRecordRealProcessJourney)$' -count=1 ./cmd/fid/integration` - FAIL - `FAIL` (terminated `exit_code: 1`); initial recovery red: list absent and activity count zero.
- `go test -tags integration -run '^(TestPhase23LeavesRejectUnknownRequestFields|TestP23ReviewRecordRealProcessJourney)$' -count=1 ./cmd/fid/integration` - FAIL - `FAIL` (terminated `exit_code: 1`); list green, Review write failed while trial activity effect was present.
- `go test -tags integration -run '^TestP23ReviewRecordRealProcessJourney$' -count=1 ./cmd/fid/integration` - FAIL - `FAIL` (terminated `exit_code: 1`); temporary local diagnostic exposed `DEBUG REVIEW: unknown table` from the journal validator. Diagnostic was removed.
- `go test -tags integration -run '^TestPhase23LeavesRejectUnknownRequestFields$' -count=1 ./cmd/fid/integration` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	2.028s` (terminated `exit_code: 0`)
- `go test -tags integration -run '^TestP23ReviewRecordRealProcessJourney$' -count=1 ./cmd/fid/integration` - FAIL - `FAIL` (terminated `exit_code: 1`); preserved activity count zero after trial effect removal.
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 ./cmd/fid ./internal/review` - FAIL - `FAIL` (terminated `exit_code: 1`); retained leaf inventory omitted newly accepted `review list`.
- `go test -run '^TestLiveLeafCommandsMatchContractInventories$' -count=1 ./cmd/fid` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	0.147s` (terminated `exit_code: 0`)
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...` - FAIL - `FAIL` (terminated `exit_code: 1`)
- `project-docs/tools/report-lint.sh project-docs/docs/phases/phase-23/tasks/task-13-done.md project-docs/docs/phases/phase-23/tasks/task-13.md` - FAIL - `FAIL report-lint found 3 failure(s)` (terminated `exit_code: 3`)

## Exact quality-gate results

Working directory for both attempts: `/home/mark/projects/fid`. The following is the preserved first attempt; its last three commands were not run at that checkpoint.

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` - BLOCKED - scope gap leaves VAL-055 probes red; command not run (no termination observed).
- `go tool cover -func=coverage.out` - BLOCKED - no task-13 coverage profile from a completed ordinary gate (no termination observed).
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...` - BLOCKED - scope gap leaves `review.list` and ActivityEvent probes red; command not run (no termination observed).

Recovery attempt 1, on the list-converged checkpoint, ran all five byte-identical commands in task-spec order. Every command terminated:

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	116.083s	coverage: 64.6% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total:											(statements)					61.3%` (terminated `exit_code: 0`)
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...` - BLOCKED - `FAIL` (terminated `exit_code: 1`); `TestP23ReviewRecordRealProcessJourney` observed zero `review-recorded` rows. The accepted activity contract cannot be made green under the current journal registry and migration order.

## Deviations from spec

- Prior checkpoint: task 13's accepted VAL-055 behavior needed `internal/review` list and Event-producer edits beyond its then-named production scope. No ActivityEvent was invented and no gate or assertion was changed. The temporarily green list implementation was removed at that checkpoint; its red test was retained. The Runner resolved this scope gap in TASK-RECOVERY attempt 1, and the list is now implemented.
- The first evidence conflict probe used an unchanged digest and correctly replayed. Changing a request field exposed the generic error adapter, which was repaired and rerun green. The first failing outputs remain recorded above.
- Recovery attempt 1 added the missing `review.list` leaf and converged the retained exact leaf inventory. The first touched-package race command failed on that inventory (`FAIL`, exit 1); after adding the accepted leaf to the inventory, the focused test passed (`ok  	github.com/mbell-nzl/fid/cmd/fid	0.147s`, exit 0). The original failure remains in Tests run.
- The scoped Review activity trial failed at `internal/journal/current_event_validation.go:161` with `unknown table`, before projection or ref-integrity checks. The trial and its temporary diagnostic were removed; the producing-leaf probe remained red (`FAIL`, exit 1) and the fifth exact gate remained red (`FAIL`, exit 1). No event catalog, journal table registry, migration, threshold, assertion, or command was changed.
- The first report-lint attempt failed with `FAIL report-lint found 3 failure(s)` (exit 3): two PASS rows had explanatory prose after the termination field, and the failed tagged required command was marked FAIL instead of BLOCKED in a blocker report. The PASS rows were reduced to their observed results; the tagged gate's original FAIL result is retained in Tests run and recorded as BLOCKED with its reason in the exact-gate section.

## Concerns, technical debt, or notes for the Tech Lead

- The Phase 23 catalog permits `activity_log/activity` on each producing Event, but `internal/journal/current_event_columns.go` omits the table, so strict current-Event validation rejects every such effect. That file is outside the Runner's bounded producer-file authorization. `requirement.record-reference` has no Issue ID in its accepted request; the current trigger references only V1 `requirements`, while task 15 owns migration 0018's new `requirement_refs` trigger. Choosing an arbitrary Issue or adding a temporary V1 row would change the model. This is the protected task-order/data-model contradiction for TL-Lite to resolve. The six remaining owning producer files were left unchanged.
- DEC-227 package times were measured from the exact commands: `cmd/fid 288.696s`, `internal/workflow 116.083s`, `cmd/fid/integration 258.375s`. No budget was breached, and no 80% headroom warning applies.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `838fc32e37ce1f3a5352f0e371fc09f65590007a`

## Commit subject

- `feat(review): add bounded review list projection`
