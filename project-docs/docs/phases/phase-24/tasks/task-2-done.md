# Phase 24 - Task 2: Input-as-written digests and named conflict fields - Completion Report

## Workflow outcome

- Task outcome: COMPLETE
- Category: none
- Condition: None.
- Recovery attempts: 4
- Next action: None.
- Human attention required: no

## Summary of what was implemented

- Added a shared request and per-field SHA-256 helper with sorted field-difference calculation. The helper hashes the family's supplied digest components and derives public field digests from the canonical request object plus `--actor`.
- Adopted the helper for `run.start` through `mutateRun`. Its same-key `claim_version` conflict now carries the stored and supplied field maps and renders an exit-4 `idempotency-conflict` naming `claim_version`.
- At the prior checkpoint, the accepted `differing_fields` member was added to the typed conflict details definition and explicit conflict renderers were moved to one details builder. This continuation adopted the remaining direct request sites and family renderers.
- Continuation attempt 1 adopted the shared helper and field maps in Issues, Labels, Links, Workflow, Evidence, Activation, Claims and Review; removed remaining family-local digest helpers; and added conflict mapping for Gate, Stop and Issue-family renderers. `gate.run` passes real maps through execution-started, Receipt update precondition and committed denial.
- Attempt 1 added process tests for Review optional presence and the DEC-242 exception; seven-leaf replay after HEAD/source changes and through unmerged/deleted-file states with new-key controls; selected family conflicts; and Gate receipt/denial maps. Its live-tree inventory listed 49 keyed leaves before executable conflict/retry rows were added in attempt 2.
- Attempt 2 executed conflict and canonical-retry rows for all 49 planned leaves, and two-field and Actor conflicts in all 20 families. It added raw replay-byte checks, exact input-as-written field-name checks, and per-workspace live/rebuilt map parity. It moved replay ahead of lookup/authority on affected leaves, retained the original Claim result snapshot on replay, and preserved typed Envelope evaluation JSON order.
- Attempt 3 added a singular real workspace/journal fixture. It performs keyed mutations through all 49 live leaves, including recovery and judgment lifecycle transitions and secondary Activation recovery/abort, then checks the exact projection operation inventory and byte-equal field-map JSON before and after `fid rebuild`.
- QA review-1 F3 rework enumerated the five project-file Source leaves in the retained seven-leaf fixture. Each now runs a distinct new-key control after `base.txt` deletion, alongside the original-byte same-key replay and unchanged Event count. The successful observable preflight and unmerged-path control remain in place.

## Files created or modified

- `internal/journal/request_digest.go`; `internal/journal/request_digest_test.go`
- `internal/workflow/run.go`; `internal/cli/codes.go`
- `cmd/fid/idempotency_conflict.go`; `cmd/fid/integration/p24_idempotency_test.go`
- `cmd/fid/claim.go`; `cmd/fid/concern.go`; `cmd/fid/decision.go`; `cmd/fid/deviation.go`; `cmd/fid/envelope.go`; `cmd/fid/evidence.go`; `cmd/fid/finding.go`; `cmd/fid/human_outcome.go`; `cmd/fid/review.go`; `cmd/fid/run.go`
- `project-docs/docs/phases/phase-24/tasks/task-2-done.md`
- Continuation source: `internal/{issues,labels,links,workflow,evidence,claims,review}/` digest call sites and attribution helpers; `cmd/fid/{activation,gate,issue_mutation_response,stop}.go`.
- Continuation tests: `cmd/fid/{activation_test,p24_keyed_inventory_test}.go`; `cmd/fid/integration/{p23_conformance_test,p23_portable_recovery_test,p24_idempotency_test}.go`; `cmd/fid/integration/testdata/p24_keyed_leaves.tsv`.
- `project-docs/specs/06_DESIGN_PLAN_CLI.md`.
- Attempt 2: `cmd/fid/p24_keyed_tree_test.go`; `cmd/fid/integration/p24_{activation,artifact,corpus_completeness,flag,foundation,gate_evidence,judgment,keyed_capture,map_parity,recovery,workflow_digest}_test.go`; `cmd/fid/integration/{e2e_v01_closure_helpers,p22_activation_fault_assertions,p22_event_nested_values,p23_conformance,p23_human_failure,p24_idempotency}_test.go`; and `cmd/fid/integration/testdata/p24_keyed_leaves.tsv`.
- Attempt 2 product: `cmd/fid/{envelope,workflow}.go`; `internal/{claims/generic,evidence/artifact,evidence/evidence,evidence/registration,journal/request_digest,journal/request_digest_test,review/decision,review/human_outcome,workflow/assign,workflow/define,workflow/envelope,workflow/stop}.go`.
- Attempt 3 tests: `cmd/fid/integration/p24_single_journal_test.go` and a reusable adapter entrypoint in `cmd/fid/integration/p23_human_leaf_pairs_test.go`.

## Documentation updated

- `project-docs/specs/06_DESIGN_PLAN_CLI.md` describes input-as-written hashing, DEC-242, named field digests and `differing_fields`. This report retains the original attempt's failed evidence and records the continuation.
- Attempt 2 made no further design-document change; the prior committed CLI specification update remains the Task 2 documentation update.
- Attempt 3 updated this completion report with the singular-journal evidence and final gate results; the design contract remains the committed CLI specification above.
- QA review-1 F3 rework updated this report with the paired deleted-source control evidence and fresh exact gate results.

## Validation Contract Coverage

- VAL-002: PASS - `TestP24All49KeyedLeavesRebuildInOneJournal` (`cmd/fid/integration/p24_single_journal_test.go:15`) executes all 49 keyed leaves in one real workspace. `p24AssertSingleJournalInventory` (`cmd/fid/integration/p24_single_journal_test.go:196`) queries that workspace's `idempotency_results` projection and requires its exact 49 distinct operation names to equal `testdata/p24_keyed_leaves.tsv`; `p24AssertRebuildMapParity` (`cmd/fid/integration/p24_map_parity_test.go:11`) compares every `request_field_digests_json` string byte-for-byte across live projection and real `fid --json rebuild`. Focused test and full tagged gate passed.
- VAL-005: PASS - The fresh exact `BenchmarkRebuildWorkspace100kEvents-6` printed `22006 ms/rebuild`, below 30000 ms; the command terminated exit 0. The original Task 2 gate printed `22449 ms/rebuild`.
- VAL-006: PASS - `TestP24EveryLiveKeyedLeafHasExecutableCorpusRow` (`cmd/fid/p24_keyed_tree_test.go:14`) and `TestP24ExecutedKeyedCorpusMatches49LeafInventory` (`cmd/fid/integration/p24_corpus_completeness_test.go:12`) bind all 49 live leaves to executable row maps. The Foundation, Recovery, Flag/Claim, Activation, Judgment, Gate/Evidence, Artifact and Workflow corpus tests (`cmd/fid/integration/p24_{foundation,recovery,flag,activation,judgment,gate_evidence,artifact,workflow_digest}*_test.go`) assert each changed-input exit 4, exact sorted `differing_fields`, unchanged journal count, and each family's two-field and Actor cases; the full tagged gate passed.
- VAL-007: PASS - `p24ExerciseCapturedLeaf` (`cmd/fid/integration/p24_foundation_corpus_test.go:52`), `TestP24JudgmentKeyedLeafConflictAndCanonicalRetryCorpus` (`cmd/fid/integration/p24_judgment_corpus_test.go:27`), `TestP24WorkflowInputAsWrittenAndTwoFieldConflict` (`cmd/fid/integration/p24_workflow_digest_test.go:14`), and `TestP24GateEvidenceAndEnvelopeKeyedLeafCorpus` (`cmd/fid/integration/p24_gate_evidence_corpus_test.go:15`) compare original EventRef and response bytes with only `replayed:true`, including reordered input keys/whitespace or flag order. All run in the passing tagged gate.
- VAL-008: PASS - `TestP24ReviewOptionalPresenceAndEmptyException` (`cmd/fid/integration/p24_idempotency_test.go:178`) checks nonempty-versus-omitted `inapplicable_findings` conflicts and DEC-242 `[]`-versus-omitted replay; `p24AssertStoredFieldNames` (`cmd/fid/integration/p24_foundation_corpus_test.go:154`) audits written top-level names for captured leaves, and the Foundation test checks `protected_decisions` presence.
- VAL-009: PASS - `TestP24ObservingJudgmentRetriesAfterSourceAndHEADChange` (`cmd/fid/integration/p24_idempotency_test.go:260`) commits a changed source file and advances HEAD, then compares raw replay bytes and unchanged journal count for all seven named leaves; the full tagged gate passed.
- VAL-052: PASS - Fresh ordinary gate reported `cmd/fid` 221.607 s, `internal/export` 98.415 s as the slowest other uncached package, and cached `internal/workflow`; the fresh tagged `cmd/fid/integration` run reported 288.286 s. The original Task 2 final-tree uncached race run observed `internal/workflow` 67.367 s and `internal/issues` 113.986 s. All recorded ordinary package times are below 600 s and the 480 s warning; tagged integration is below 900 s.
- VAL-066: PASS for Task 2's seven leaves - `TestP24ObservingJudgmentReplayNeverObserves` (`cmd/fid/integration/p24_idempotency_test.go:310`) checks original response bytes and unchanged Event count on all seven same-key unmerged replays. For the five project-file Source leaves enumerated at line 312, it deletes `base.txt`, repeats the same-key replay, then sends a distinct viable `-p24-deleted` new-key request in that same unmerged/deleted state. Each returns a typed JSON error with no stdout and appends no Event (lines 371-394). The existing observable new-key preflight succeeds for each leaf (line 350). Focused race and full tagged gates passed. Task 4 owns isolated deleted-file observation behavior and the `repository-unavailable` code.

## Tests run

- First red: `go test -tags integration -run '^TestP24RunStartConflictNamesChangedClaimVersion$' -count=1 ./cmd/fid/integration/` - FAIL - `FAIL` (terminated `exit_code: 1`); `run.start` returned exit 2 `invalid-argument`, rather than the required exit-4 conflict.
- Helper red: `go test -run '^TestComputeRequestDigest' -count=1 ./internal/journal` - FAIL - `FAIL` (terminated `exit_code: 1`); `ComputeRequestDigest` and `DifferingRequestFields` were undefined.
- `go test -run '^TestComputeRequestDigest' -count=1 ./internal/journal` - PASS - `ok  	github.com/mbell-nzl/fid/internal/journal	0.004s` (terminated `exit_code: 0`)
- `go test -tags integration -run '^TestP24RunStartConflictNamesChangedClaimVersion$' -count=1 ./cmd/fid/integration/` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	2.213s` (terminated `exit_code: 0`)
- `go test -run 'Conflict|Idempotency' -count=1 ./cmd/fid` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	0.742s` (terminated `exit_code: 0`)
- After formatting, `go test -tags integration -run '^TestP24RunStartConflictNamesChangedClaimVersion$' -count=1 ./cmd/fid/integration/` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	1.560s` (terminated `exit_code: 0`)

- Continuation portable-recovery red: `go test -tags integration -run '^TestPortableRoundTripCoversEveryPhase23Table$' -count=1 ./cmd/fid/integration/` - FAIL - `FAIL` (terminated `exit_code: 1`); Issue, Gate registration, Envelope, Run observation and Evidence link stored `{}`. After helper adoption, the byte-identical command passed: `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	10.834s` (terminated `exit_code: 0`).
- Gate red/green: `go test -tags integration -run '^TestP24GateRunPersistsNamedFieldsThroughReceipt$' -count=1 ./cmd/fid/integration/` - FAIL - `FAIL` (terminated `exit_code: 1`, field map empty after admitted execution); after real map propagation and conflict rendering, PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	1.704s` (terminated `exit_code: 0`). `TestP24GateRunCommittedDenialPersistsNamedFields` also passed in the combined Gate selection.
- Issue renderer red/green: `go test -tags integration -run '^TestP24IssueCreateConflictNamesChangedTitle$' -count=1 ./cmd/fid/integration/` - FAIL - `FAIL` (terminated `exit_code: 1`, generic invalid argument); after renderer repair, PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	1.597s` (terminated `exit_code: 0`).
- Stop renderer red/green: `go test -tags integration -run '^TestP24StopRaiseConflictNamesChangedDetails$' -count=1 ./cmd/fid/integration/` - FAIL - `FAIL` (terminated `exit_code: 1`, generic invalid argument); after renderer repair, PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	1.261s` (terminated `exit_code: 0`).
- `go test -tags integration -run '^TestP24' -count=1 ./cmd/fid/integration/` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	13.420s` (terminated `exit_code: 0`)
  - Real entrypoint: the Issue fixture invoked `fid issue create --kind task --title "different title" --actor agent:developer --idempotency-key p22-corpus-issue --json` and observed exit 4, `differing_fields:["--title"]`, with no Event.
- Retained conformance red: `go test -tags integration -run '^TestP23SharedLeafConformance$' -count=1 ./cmd/fid/integration/` - FAIL - `FAIL` (terminated `exit_code: 1`); nine fake conflict envelopes omitted `differing_fields`. After exact-field convergence the byte-identical command passed: `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	10.322s` (terminated `exit_code: 0`).
- Touched-package race attempt: `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 ./internal/journal ./internal/issues ./internal/labels ./internal/links ./internal/workflow ./internal/evidence ./internal/claims ./internal/review ./cmd/fid` - FAIL - `FAIL` (terminated `exit_code: 1`); retained Activation assertion expected generic `invalid-argument`. After exact conflict convergence the byte-identical command passed: `ok  	github.com/mbell-nzl/fid/cmd/fid	251.732s` (terminated `exit_code: 0`).
- `go test -run '^TestP24KeyedLeafCorpusMatchesLiveTree$' -count=1 ./cmd/fid` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	0.152s` (terminated `exit_code: 0`)

- `go test -tags integration -run '^TestP24FoundationKeyedLeafConflictAndCanonicalRetryCorpus$' -count=1 ./cmd/fid/integration/` - FAIL - `FAIL` (terminated `exit_code: 1`)
- `go test -tags integration -run '^TestP24FoundationKeyedLeafConflictAndCanonicalRetryCorpus$' -count=1 ./cmd/fid/integration/` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	2.521s` (terminated `exit_code: 0`)
- `go test -tags integration -run '^TestP24' -count=1 ./cmd/fid/integration/` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	18.378s` (terminated `exit_code: 0`)
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 ./cmd/fid ./internal/claims ./internal/evidence ./internal/journal ./internal/review ./internal/workflow` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	80.599s` (terminated `exit_code: 0`)
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 -run '^TestP24' ./cmd/fid/integration/` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	22.976s` (terminated `exit_code: 0`)
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...` - FAIL - `FAIL` (terminated `exit_code: 1`)
- `go test -tags integration -run '^(TestP23StopFailureCodeMatrix|TestP22ActivationRecoverySurvivesAppendProjectionInterruption|TestP22CurrentEventFamiliesLiveReplayValidation)$' -count=1 ./cmd/fid/integration/` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	32.345s` (terminated `exit_code: 0`)
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...` - FAIL - `FAIL` (terminated `exit_code: 1`)
- `go test -tags integration -run '^TestP23SharedAmendedLeafConformance$' -count=1 ./cmd/fid/integration/` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	2.025s` (terminated `exit_code: 0`)
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` - FAIL - `FAIL` (terminated `exit_code: 1`)
- The actual `fid --json` entrypoint is exercised by `p24ExerciseCapturedLeaf` for all 49 keyed operations and by `TestP24GateEvidenceAndEnvelopeKeyedLeafCorpus` for streamed Gate output; the final tagged gate below passed.
- Attempt 3 singular-journal red: `go test -tags integration -run '^TestP24All49KeyedLeavesRebuildInOneJournal$' -count=1 ./cmd/fid/integration/` - FAIL - `FAIL` (terminated `exit_code: 1`); the new real workspace projected only 11 operation names, versus the required 49. Subsequent red iterations exposed an invalid authority-profile capability, stale Run version, and Issue transition preconditions; sequencing the same-workspace mutations repaired the fixture.
- `go test -tags integration -run '^TestP24All49KeyedLeavesRebuildInOneJournal$' -count=1 ./cmd/fid/integration/` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	10.033s` (terminated `exit_code: 0`)
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 -run '^(TestP24All49KeyedLeavesRebuildInOneJournal|TestP24JudgmentKeyedLeafConflictAndCanonicalRetryCorpus|TestP23SharedLeafConformance)$' ./cmd/fid/integration/` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	14.415s` (terminated `exit_code: 0`)
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 ./internal/workflow ./internal/issues` - PASS - `ok  	github.com/mbell-nzl/fid/internal/issues	113.986s` (terminated `exit_code: 0`)
- The same uncached race command printed `ok  	github.com/mbell-nzl/fid/internal/workflow	67.367s` before its final Issues line.
- Real entrypoint in the singular-journal fixture: `fid --project <one workspace> --json rebuild` exited 0 after 49 real keyed operation names were projected, and every live `request_field_digests_json` byte string equaled its rebuilt value.
- QA review-1 F3 red: `go test -tags integration -run '^TestP24ObservingJudgmentReplayNeverObserves$' -count=1 ./cmd/fid/integration/` - FAIL - `FAIL` (terminated `exit_code: 1`; tool wall 13.287 s). All five Source-bearing subtests reported `deleted-source new-key control was not verified` before the paired call was added.
- Isolated-deletion probe with the same focused command - FAIL - `FAIL` (terminated `exit_code: 1`; tool wall 12.342 s). After clearing the unmerged index, `concern.record`, `deviation.record`, and `human-outcome.record` accepted a new key with a missing Source file. This is Task 4 observation behavior, so the fixture retains the existing unmerged state through deletion.
- QA review-1 F3 green: `go test -tags integration -run '^TestP24ObservingJudgmentReplayNeverObserves$' -count=1 ./cmd/fid/integration/` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	10.931s` (terminated `exit_code: 0`)
- Touched tagged race pre-gate: `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 -run '^TestP24ObservingJudgmentReplayNeverObserves$' ./cmd/fid/integration/` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	12.677s` (terminated `exit_code: 0`)
- Real CLI entrypoint in that fixture: `fid --project <workspace> --json <source leaf> --input <fixture> --actor <actor> --idempotency-key <distinct deleted key>` returned a nonzero status and typed JSON error for each of the five Source leaves after `base.txt` removal, while the paired same-key invocation returned its original stdout bytes. The test checks the journal line count before and after both calls.
- Fresh exact-gate wall evidence: format 0.025 s; lint 3.501 s; ordinary race gate's longest reported package `cmd/fid` 221.607 s (slowest other ordinary package `internal/export` 98.415 s; `internal/workflow` cached); coverage tool 0.095 s; tagged integration package 288.286 s; benchmark tool 24.832 s (22,006 ms/rebuild); each `rg` tool call <0.001 s. The execution surface did not provide a single aggregate elapsed value across the ordinary gate's polling chunks; the reported package times are its reproducible wall evidence.

## Exact quality-gate results

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	(cached)	coverage: 62.6% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total:											(statements)			60.2%` (terminated `exit_code: 0`)
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	288.286s` (terminated `exit_code: 0`)
- `go test -run '^$' -bench '^BenchmarkRebuildWorkspace100kEvents$' -benchtime 1x ./internal/journal/` - PASS - `ok  	github.com/mbell-nzl/fid/internal/journal	24.634s` (terminated `exit_code: 0`)
- `rg -n "RequestDigest:" internal cmd -g '!*_test.go'` - PASS - `internal/workflow/run_budget.go:61:		idem := journal.IdempotencyRequest{OperationScope: "run.observe-budget:" + project.ProjectID, IdempotencyKey: options.IdempotencyKey, RequestDigest: digest, RequestFieldDigests: fields}` (terminated `exit_code: 0`)
- `rg -n "func [a-zA-Z]*(RequestDigest|requestDigest|judgmentDigest)\(" internal -g '!*_test.go'` - PASS - `internal/journal/request_digest.go:14:func ComputeRequestDigest(components map[string]any, actor any, aliases map[string]string) (string, map[string]string, error) {` (terminated `exit_code: 0`)

## Deviations from spec

- QA review-1 F3 first focused attempt failed with `FAIL` (terminated `exit_code: 1`) because the deleted-source branch had no new-key call. Adding the paired call and exact five-leaf enumeration repaired that evidence gap; the byte-identical focused command passed with `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	10.931s` (terminated `exit_code: 0`). A second focused probe failed with `FAIL` (terminated `exit_code: 1`) when it cleared the unmerged index: three reference leaves accepted the missing-file request under today's behavior. Keeping the accepted unmerged/deleted state restored the paired test; Task 4 owns file-only observation semantics. The full byte-identical Task 2 gate then passed.
- No accepted-contract deviation. The final singular-journal fixture uses only real keyed CLI writes in one workspace and changes no public contract.
- The first touched-package race attempt failed with final token `FAIL` (terminated `exit_code: 1`): a retained Phase 21 Task 4 Activation assertion expected the displaced generic `invalid-argument`. The repair requires exit-4 `idempotency-conflict`, exact `content`/`expected_version` fields, identical repeated failure and unchanged Event count. The byte-identical repaired race command terminated `exit_code: 0` with final token `ok  	github.com/mbell-nzl/fid/cmd/fid	251.732s`.
- The first retained Phase 23 conformance selection failed with final token `FAIL` (terminated `exit_code: 1`) because its fake model omitted the new closed details member. Exact nine-field convergence made the byte-identical selection pass with final token `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	10.322s` (terminated `exit_code: 0`).

- The first committed-report lint attempt terminated `exit_code: 3` with final token `FAIL report-lint found 3 failure(s)`: three PASS rows had prose after their termination tokens. A docs-only correction moved that prose outside the rows; lint is rerun against the corrected committed report.
- Attempt 2's first tagged gate failed with token `FAIL` (terminated `exit_code: 1`): `TestP23StopFailureCodeMatrix` lacked a conflict row, `p22RejectActivationFaultConflicts` expected pre-Task-2 invalid-argument for a same-key changed request, and `p22RejectNestedCurrentValues` treated field-map names as closed. Added a real Stop conflict row and exact Activation conflict assertion; retained nested JSON type/digest checks while allowing names to follow written optional input. The focused three-test selection passed and the final byte-identical tagged gate passed with `ok  \tgithub.com/mbell-nzl/fid/cmd/fid/integration\t280.036s` (terminated `exit_code: 0`).
- The second tagged attempt failed with token `FAIL` (terminated `exit_code: 1`) because the Phase 23 amended fake conflict model expected `differing_fields:[""]` for Envelope, Claim and Evidence. Added the four changed public names and kept the fake/real comparison; the focused selection and final byte-identical tagged gate passed.
- An intermediate ordinary gate failed with token `FAIL` (terminated `exit_code: 1`) after an attempted `--actor` Event-validator guard rejected historical synthetic `{}` maps. The guard was removed; real CLI field-name assertions remain in the new corpus. The final byte-identical ordinary gate passed with `ok  \tgithub.com/mbell-nzl/fid/internal/workflow\t(cached)\tcoverage: 62.6% of statements` (terminated `exit_code: 0`).
- Attempt 3's first singular-journal test failed with `FAIL` (terminated `exit_code: 1`): its initial 11-operation projection did not match the required 49-operation inventory. The fixture was extended with real mutations in the same workspace; invalid capabilities, a stale Run version, and post-close Issue transition preconditions were repaired by accepted capability selection and lifecycle ordering. The byte-identical focused command later passed with `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	10.033s` (terminated `exit_code: 0`), and the exact tagged gate passed. No gate definition, assertion, budget or threshold changed.

## Concerns, technical debt, or notes for the Tech Lead

- QA review-1 F3 limit: the deleted-source control runs while the fixture is also unmerged, exactly matching the retained state after the first refusal. The three reference leaves' isolated missing-file acceptance is left for Task 4; this test proves replay bypasses observation and that a fresh viable request is refused in the combined state. Fresh gates reported `cmd/fid` 221.607 s, `internal/export` 98.415 s and tagged integration 288.286 s, below their 600/900 s budgets and 480 s ordinary headroom warning.
- The original checkpoint left most write sites storing `{}` and some renderers returning generic failures; attempt 1 adopted their direct request constructors and Gate, Stop and Issue renderers, and attempt 2 ran the all-leaf behavioral corpus.
- The original checkpoint had changed no retained assertion. This continuation converged three displaced retained assertions after checking their introducing commits, as required by DEC-162/165/166.
- The first red exposed a missing `run.start` CLI conflict mapping. The integration test and focused green prove that one path only; they do not prove the 49-leaf or observation contracts.
- Retained-test provenance: `cmd/fid/activation_test.go` replay assertion was introduced by `0444ccc458f3ac5a720724825f7f81589cb41e09` (`phase-21/task-4`); `cmd/fid/integration/p23_conformance_test.go` fake conflict details by `2e7fd2f9fa05bce648510789453c8d2ed8633b1d` (`phase-23/task-18`); `cmd/fid/integration/p23_portable_recovery_test.go` by `d02614e6c26757076a91d8a12e08dfa328fe699d` (`phase-23/task-16`). Each preserves its replay, exact-envelope, or round-trip subject while asserting Task 2's accepted maps or details.
- The 49-row TSV, live-tree test, executable corpus, per-workspace rebuild checks, and attempt 3's singular-journal rebuild check all pass. The final ordinary gate reported `cmd/fid` 196.324 s; the final tagged gate reported `cmd/fid/integration` 248.512 s. An uncached final-tree race run reported `internal/workflow` 67.367 s and slowest other ordinary package `internal/issues` 113.986 s. No ordinary package reached the 480 s warning or 600 s limit; tagged integration remained under 900 s.
- Attempt 2 retained-test introducing commits: Stop failure inventory `180f5168e3f1338cb758f8388f1c26fd6f8428d3` (`phase-23/task-19`); Activation changed-request assertion `e7d5963105fa3a40b8d0d687902def63595ade2f` (`phase-22/task-10`); nested Event-value assertion `c8e0061a11c1405e98fd5e4f09612b4f9099af51` (`phase-22/task-1`); amended fake conflict assertion `2e7fd2f9fa05bce648510789453c8d2ed8633b1d` (`phase-23/task-18`). Each keeps the original parity, failure, interruption, or nested-value subject while converging to Task 2's open field maps and exact conflict names.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending commit

## Commit subject

- Pending commit
