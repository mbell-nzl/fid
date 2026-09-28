# Phase 23 - Task 13: CLI surface, failure codes, and idempotency - Blocker Report

## Workflow outcome

- Task outcome: BLOCKED
- Category: specification-contradiction
- Condition: The accepted `public-interface.md` section 5.2 and VAL-057 require a `protected-decision-required` denial to return immutable Stop/Attention IDs, null EventRef, and `replayed:false`, but the same document's section 5 declares closed details containing only `envelope_id`, `decision_class`, `required_capabilities`, and `affected_actions`, and section 1 retains the closed `fid-cli.v1` error envelope. The built binary returns that exact closed error without the required IDs or EventRef/replay fields. Adding them changes the protected public interface; omitting them fails VAL-057.
- Recovery attempts: 2
- Next action: Obtain protected authority reconciling `public-interface.md` sections 1, 5, and 5.2 with VAL-057, then resume task 13 under an exact directive. Preserve the red retry assertion and rerun the unchanged tagged gate.
- Human attention required: yes

## Summary of what was implemented

- Added a built-binary request-closure matrix for fourteen Phase 23 mutations. It rejects unknown fields, a null request object, trailing JSON, and duplicate keys without appending an Event.
- Registered six previously omitted Phase 23 failure codes and aligned the exact-schema test. Renamed the unexported failure registry and test-local expectation type per DEC-226.
- Rendered `evidence.link` idempotency conflicts through the accepted failure envelope, verified by a built-binary same-key/different-digest probe.
- At the prior checkpoint, preserved red real-process probes for the missing `review.list` leaf and absent `review-recorded` ActivityEvent. The Runner's TASK-RECOVERY attempt 1 subsequently authorized the bounded Review reader and producer files.
- Added `review.list` with its accepted four filters, default and range limit, complete bounded Review projection, and typed `result-limit-exceeded`. The built-binary list and leaf-inventory probes now pass.
- Tried `review-recorded` and `finding-resolved` effects in the owning Review producer. The first producing call failed at the journal validator with `unknown table`; the trial code was removed so existing Review writes remain available. The failed observation is retained below.
- DEC-262 resolved that scope and task-order gap for attempt 2: the strict journal registry now admits `activity_log` from migration 0002, including the narrow Review cardinality and `details` JSON path. The five authorized Review producers now write the seven accepted issue-subject ActivityEvent names. The four requirement/decision rows remain deferred to task 15 (DEC-262).
- Added real-process failure-code producer checks, two guarded-action precedence probes, replay/conflict probes through the built binary, and repeated noncommitting denial probes. The protected-Decision retry exposes the new specification contradiction above; its failing assertion is retained.

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
- `cmd/fid/claim.go`, `cmd/fid/envelope.go`
- `cmd/fid/integration/p22_fake_test.go`, `p23_concern_deviation_test.go`, `p23_decision_reference_test.go`, `p23_finding_disposition_test.go`, `p23_human_outcome_test.go`, `p23_requirement_reference_test.go`, `p23_review_evidence_budget_test.go`, `p23_failure_contract_test.go`, `p23_task13_fixture_replay_test.go`, `p23_task13_precedence_test.go`
- `internal/journal/current_event_columns.go`, `current_event_cardinality.go`, `current_event_json.go`, `current_event_validation.go`, `current_event_validation_test.go`
- `internal/review/concern.go`, `deviation.go`, `finding.go`, `human_outcome.go`
- `project-docs/specs/06_DESIGN_PLAN_CLI.md`
- `project-docs/docs/phases/phase-23/tasks/task-13-done.md`

## Documentation updated

- `project-docs/specs/06_DESIGN_PLAN_CLI.md` section 5 now lists the accepted Phase 23 failure codes, categories, exits, subjects, and closed detail keys.
- This report records DEC-262's resolved activity scope gap, the task-15 deferral, and the new protected VAL-057 error-envelope contradiction with observed evidence.

## Validation Contract Coverage

- VAL-055: PASS for task 13's DEC-262 subset - `TestP23ReviewRecordRealProcessJourney`, `TestP23FindingDispositionRealProcessJourney`, `TestP23ConcernDeviationRealProcessJourney`, and `TestP23HumanOutcomeAuthorityAndCompatibilityMatrixRealProcess` inspect the seven issue-subject names and their Issue references through `p23AssertActivity` (`cmd/fid/integration/p23_review_record_test.go:160`); that helper also rejects any name outside the accepted ten or any deferred requirement/decision row. `TestCurrentEventColumnDomainsMatchAcceptedDDL` proves the registry entry against migration 0002. The focused built-binary Review/leaf journey passed. `requirement-referenced`, `requirement-linked`, `decision-referenced`, and requirement-subject Concern/Deviation rows are deferred to task 15 (DEC-262); they are not claimed PASS.
- VAL-056: PASS for the added-code producer and precedence probes - `assertP23AddedFailureContract` (`cmd/fid/integration/p23_failure_contract_test.go`) checks the real-process producing operation, category, retryability, subject, and closed detail keys; `assertP22ProcessFailure` checks the exact exit. `TestP23FailurePrecedenceWithProtectedDecisionAndObligation` and `TestP23FailurePrecedenceWithReadinessStopAndObligation` passed under the tagged focused command. The protected Decision's accepted closed details are precisely what conflict with VAL-057.
- VAL-057: BLOCKED - `TestP23AmendedEnvelopeReplayThroughBuiltBinary` passed same-key replay and changed-digest conflict, and `TestP23ReviewCyclesRealProcessJourney` passed repeated exhaustion denials with no Event or idempotency result. `TestP23DecisionReferenceProtectedLifecycleRealProcess` fails at line 34 because both byte-identical `protected-decision-required` responses omit explicit EventRef/replay fields and Stop/Attention IDs. Adding them would alter the closed accepted failure schema or envelope. The exact tagged gate failed on this assertion.
- VAL-079: PASS for observed attempt-2 DEC-227 times: ordinary `cmd/fid 284.783s` and `internal/workflow 99.417s` (each below 600s); tagged `cmd/fid/integration 217.733s` (below 900s). Tagged failure is VAL-057, not timing. DEC-258/259 already resolved DEC-254 Part 2; no count=2 headroom step was run.

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
- `go test -run '^TestCurrentEventColumnDomainsMatchAcceptedDDL$' -count=1 ./internal/journal` - FAIL - `FAIL` (terminated `exit_code: 1`); red registry proof before the DEC-262 entry.
- `go test -run '^TestCurrentEventColumnDomainsMatchAcceptedDDL$' -count=1 ./internal/journal` - PASS - `ok  \tgithub.com/mbell-nzl/fid/internal/journal\t0.070s` (terminated `exit_code: 0`)
- `go test -tags integration -run '^(TestP23ReviewRecordRealProcessJourney|TestPhase23LeavesRejectUnknownRequestFields)$' -count=1 ./cmd/fid/integration` - FAIL - `FAIL` (terminated `exit_code: 1`); activity count zero before the producer repair.
- `go test -tags integration -run '^(TestP23ReviewRecordRealProcessJourney|TestPhase23LeavesRejectUnknownRequestFields)$' -count=1 ./cmd/fid/integration` - PASS - `ok  \tgithub.com/mbell-nzl/fid/cmd/fid/integration\t3.649s` (terminated `exit_code: 0`)
- `go test -tags integration -run '^TestP23FailurePrecedenceWith' -count=1 ./cmd/fid/integration` - PASS - `ok  \tgithub.com/mbell-nzl/fid/cmd/fid/integration\t4.164s` (terminated `exit_code: 0`)
- `go test -tags integration -run '^TestP23AmendedEnvelopeReplayThroughBuiltBinary$' -count=1 ./cmd/fid/integration` - PASS - `ok  \tgithub.com/mbell-nzl/fid/cmd/fid/integration\t1.566s` (terminated `exit_code: 0`)
- `go test -tags integration -run '^TestP23ReviewCyclesRealProcessJourney$' -count=1 ./cmd/fid/integration` - PASS - `ok  \tgithub.com/mbell-nzl/fid/cmd/fid/integration\t3.004s` (terminated `exit_code: 0`)
- `go test -tags integration -run '^(TestP23ReviewCyclesRealProcessJourney|TestP23DecisionReferenceProtectedLifecycleRealProcess|TestP23AmendedEnvelopeReplayThroughBuiltBinary|TestP23ConcernDeviationRealProcessJourney|TestP23ReviewRecordRealProcessJourney)$' -count=1 ./cmd/fid/integration` - FAIL - `FAIL` (terminated `exit_code: 1`); protected denial's closed error omitted EventRef and replay marker.
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 ./cmd/fid ./internal/review ./internal/journal` - PASS - `ok  \tgithub.com/mbell-nzl/fid/internal/journal\t76.583s` (terminated `exit_code: 0`)
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -run '^(TestP23DecisionReferenceProtectedLifecycleRealProcess|TestP23ReviewCyclesRealProcessJourney)$' -count=1 ./cmd/fid/integration` - FAIL - `FAIL` (terminated `exit_code: 1`); same protected denial mismatch.

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

Recovery attempt 2, working directory `/home/mark/projects/fid`, ran the required focused journey first and then all five byte-identical task commands in order. The focused journey passed with `ok  \tgithub.com/mbell-nzl/fid/cmd/fid/integration\t3.649s` (terminated `exit_code: 0`). The exact commands terminated as follows:

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` - PASS - `ok  \tgithub.com/mbell-nzl/fid/internal/workflow\t99.417s\tcoverage: 64.6% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total:\t\t\t\t\t\t\t\t\t\t\t(statements)\t\t\t\t\t61.2%` (terminated `exit_code: 0`)
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...` - BLOCKED - `FAIL` (terminated `exit_code: 1`); `TestP23DecisionReferenceProtectedLifecycleRealProcess` found that both protected-Decision denials omitted the explicit null EventRef, `replayed:false`, and Stop/Attention IDs required by VAL-057. The package reported 217.733s, below its budget.

## Deviations from spec

- Prior checkpoint: task 13's accepted VAL-055 behavior needed `internal/review` list and Event-producer edits beyond its then-named production scope. No ActivityEvent was invented and no gate or assertion was changed. The temporarily green list implementation was removed at that checkpoint; its red test was retained. The Runner resolved this scope gap in TASK-RECOVERY attempt 1, and the list is now implemented.
- The first evidence conflict probe used an unchanged digest and correctly replayed. Changing a request field exposed the generic error adapter, which was repaired and rerun green. The first failing outputs remain recorded above.
- Recovery attempt 1 added the missing `review.list` leaf and converged the retained exact leaf inventory. The first touched-package race command failed on that inventory (`FAIL`, exit 1); after adding the accepted leaf to the inventory, the focused test passed (`ok  	github.com/mbell-nzl/fid/cmd/fid	0.147s`, exit 0). The original failure remains in Tests run.
- The scoped Review activity trial failed at `internal/journal/current_event_validation.go:161` with `unknown table`, before projection or ref-integrity checks. The trial and its temporary diagnostic were removed; the producing-leaf probe remained red (`FAIL`, exit 1) and the fifth exact gate remained red (`FAIL`, exit 1). No event catalog, journal table registry, migration, threshold, assertion, or command was changed.
- The first report-lint attempt failed with `FAIL report-lint found 3 failure(s)` (exit 3): two PASS rows had explanatory prose after the termination field, and the failed tagged required command was marked FAIL instead of BLOCKED in a blocker report. The PASS rows were reduced to their observed results; the tagged gate's original FAIL result is retained in Tests run and recorded as BLOCKED with its reason in the exact-gate section.
- DEC-262 authorizes the task-13 `activity_log` registry and issue-subject producers; it resolves the prior `TRG-23-007`/`ESC-23-007` scope gap after the 20:46 STOP. The earlier BLOCKED report, zero-row activity probes, failed activity trial, failed tagged gate, and failed report-lint remain above as immutable attempt-1 evidence. Migration 0018 and four requirement/decision activity rows remain task 15 work.
- Attempt 2's failing tagged gate token was `FAIL` (terminated `exit_code: 1`). The cause is the protected-Decision denial schema contradiction: the actual response obeys the accepted section 5 closed details and retained `fid-cli.v1` error envelope, while section 5.2 and VAL-057 require fields absent from both. No gate, test assertion, detail schema, or interface was weakened or changed; the failing probe remains red pending protected authority.

## Concerns, technical debt, or notes for the Tech Lead

- DEC-262 resolved the prior activity table scope gap. The four requirement/decision rows are deliberately deferred to task 15 under migration 0018. Attempt 2 now needs a protected interface/validation decision: `public-interface.md` section 5.2 requires fields and IDs that section 5's protected-Decision detail schema and section 1's retained error envelope cannot carry. The built binary returns the closed error exactly, and the retry test preserves the unsatisfied assertion.
- DEC-227 times from attempt 2's exact commands are `cmd/fid 284.783s`, `internal/workflow 99.417s`, and `cmd/fid/integration 217.733s`; no budget or 80% headroom warning applies. Attempt 1's times were `288.696s`, `116.083s`, and `258.375s`.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending commit

## Commit subject

- Pending commit
