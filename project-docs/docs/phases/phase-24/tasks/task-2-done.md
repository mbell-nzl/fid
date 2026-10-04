# Phase 24 - Task 2: Input-as-written digests and named conflict fields - Completion Report

## Workflow outcome

- Task outcome: RECOVERY_REQUIRED
- Category: in-scope-implementation
- Condition: The 49-leaf executable conflict/retry corpus, every-family two-field and Actor cases, optional-presence checks, seven observing-leaf checks, and final eight-command gate pass. VAL-002 still requires one real journal containing a keyed mutation from each of the 49 leaves. Current live-versus-rebuilt map parity spans several fixture journals, so the singular-journal assertion remains unproved.
- Recovery attempts: 2
- Next action: Build one real workspace/journal that performs all 49 keyed leaves (including the lifecycle-sensitive recovery and judgment leaves), assert its exact 49-operation projection inventory and byte-equal live/rebuilt field maps, then rerun the eight exact commands on that final tree and finalize this report as COMPLETE only if every assertion passes.
- Human attention required: no

## Summary of what was implemented

- Added a shared request and per-field SHA-256 helper with sorted field-difference calculation. The helper hashes the family's supplied digest components and derives public field digests from the canonical request object plus `--actor`.
- Adopted the helper for `run.start` through `mutateRun`. Its same-key `claim_version` conflict now carries the stored and supplied field maps and renders an exit-4 `idempotency-conflict` naming `claim_version`.
- At the prior checkpoint, the accepted `differing_fields` member was added to the typed conflict details definition and explicit conflict renderers were moved to one details builder. This continuation adopted the remaining direct request sites and family renderers.
- Continuation attempt 1 adopted the shared helper and field maps in Issues, Labels, Links, Workflow, Evidence, Activation, Claims and Review; removed remaining family-local digest helpers; and added conflict mapping for Gate, Stop and Issue-family renderers. `gate.run` passes real maps through execution-started, Receipt update precondition and committed denial.
- Attempt 1 added process tests for Review optional presence and the DEC-242 exception; seven-leaf replay after HEAD/source changes and through unmerged/deleted-file states with new-key controls; selected family conflicts; and Gate receipt/denial maps. Its live-tree inventory listed 49 keyed leaves before executable conflict/retry rows were added in attempt 2.
- Attempt 2 executed conflict and canonical-retry rows for all 49 planned leaves, and two-field and Actor conflicts in all 20 families. It added raw replay-byte checks, exact input-as-written field-name checks, and per-workspace live/rebuilt map parity. It moved replay ahead of lookup/authority on affected leaves, retained the original Claim result snapshot on replay, and preserved typed Envelope evaluation JSON order.

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

## Documentation updated

- `project-docs/specs/06_DESIGN_PLAN_CLI.md` describes input-as-written hashing, DEC-242, named field digests and `differing_fields`. This report retains the original attempt's failed evidence and records the continuation.
- Attempt 2 made no further design-document change; the prior committed CLI specification update remains the Task 2 documentation update.

## Validation Contract Coverage

- VAL-002: BLOCKED by unfinished in-scope evidence, not a contract decision. `p24AssertRebuildMapParity` (`cmd/fid/integration/p24_map_parity_test.go:11`) compares every live field-map row byte-for-byte after a real `fid --json rebuild` in each fixture workspace, and the tagged gate passed. The test does not yet put all 49 keyed leaves into one journal as the assertion requires.
- VAL-005: PASS - `BenchmarkRebuildWorkspace100kEvents-6` printed `23163 ms/rebuild` in the exact benchmark gate, below 30000 ms; the command terminated exit 0.
- VAL-006: PASS - `TestP24EveryLiveKeyedLeafHasExecutableCorpusRow` (`cmd/fid/p24_keyed_tree_test.go:14`) and `TestP24ExecutedKeyedCorpusMatches49LeafInventory` (`cmd/fid/integration/p24_corpus_completeness_test.go:12`) bind all 49 live leaves to executable row maps. The Foundation, Recovery, Flag/Claim, Activation, Judgment, Gate/Evidence, Artifact and Workflow corpus tests (`cmd/fid/integration/p24_{foundation,recovery,flag,activation,judgment,gate_evidence,artifact,workflow_digest}*_test.go`) assert each changed-input exit 4, exact sorted `differing_fields`, unchanged journal count, and each family's two-field and Actor cases; the full tagged gate passed.
- VAL-007: PASS - `p24ExerciseCapturedLeaf` (`cmd/fid/integration/p24_foundation_corpus_test.go:52`), `TestP24JudgmentKeyedLeafConflictAndCanonicalRetryCorpus` (`cmd/fid/integration/p24_judgment_corpus_test.go:27`), `TestP24WorkflowInputAsWrittenAndTwoFieldConflict` (`cmd/fid/integration/p24_workflow_digest_test.go:14`), and `TestP24GateEvidenceAndEnvelopeKeyedLeafCorpus` (`cmd/fid/integration/p24_gate_evidence_corpus_test.go:15`) compare original EventRef and response bytes with only `replayed:true`, including reordered input keys/whitespace or flag order. All run in the passing tagged gate.
- VAL-008: PASS - `TestP24ReviewOptionalPresenceAndEmptyException` (`cmd/fid/integration/p24_idempotency_test.go:178`) checks nonempty-versus-omitted `inapplicable_findings` conflicts and DEC-242 `[]`-versus-omitted replay; `p24AssertStoredFieldNames` (`cmd/fid/integration/p24_foundation_corpus_test.go:154`) audits written top-level names for captured leaves, and the Foundation test checks `protected_decisions` presence.
- VAL-009: PASS - `TestP24ObservingJudgmentRetriesAfterSourceAndHEADChange` (`cmd/fid/integration/p24_idempotency_test.go:260`) commits a changed source file and advances HEAD, then compares raw replay bytes and unchanged journal count for all seven named leaves; the full tagged gate passed.
- VAL-052: PASS - Final ordinary gate reported `cmd/fid` 195.085 s, `internal/workflow` cached after a directly observed 103.379 s earlier in this attempt, and slowest other ordinary package cached after a directly observed `internal/issues` 190.913 s earlier in this attempt. Final tagged integration was 280.036 s. All are below 600/900 s and the 480 s warning.
- VAL-066: PASS - `TestP24ObservingJudgmentReplayNeverObserves` (`cmd/fid/integration/p24_idempotency_test.go:310`) checks raw replay bytes and no new Event for all seven in unmerged states and after source deletion where applicable; its observable new-key preflight succeeds and unobservable new-key control fails without an Event. The full tagged gate passed.

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

## Exact quality-gate results

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	(cached)	coverage: 62.6% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total:											(statements)			60.5%` (terminated `exit_code: 0`)
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	280.036s` (terminated `exit_code: 0`)
- `go test -run '^$' -bench '^BenchmarkRebuildWorkspace100kEvents$' -benchtime 1x ./internal/journal/` - PASS - `ok  	github.com/mbell-nzl/fid/internal/journal	25.832s` (terminated `exit_code: 0`)
- `rg -n "RequestDigest:" internal cmd -g '!*_test.go'` - PASS - `internal/workflow/run_budget.go:61:		idem := journal.IdempotencyRequest{OperationScope: "run.observe-budget:" + project.ProjectID, IdempotencyKey: options.IdempotencyKey, RequestDigest: digest, RequestFieldDigests: fields}` (terminated `exit_code: 0`)
- `rg -n "func [a-zA-Z]*(RequestDigest|requestDigest|judgmentDigest)\(" internal -g '!*_test.go'` - PASS - `internal/journal/request_digest.go:14:func ComputeRequestDigest(components map[string]any, actor any, aliases map[string]string) (string, map[string]string, error) {` (terminated `exit_code: 0`)

## Deviations from spec

- No accepted-contract deviation. This is a same-task checkpoint, not task completion.
- The first touched-package race attempt failed with final token `FAIL` (terminated `exit_code: 1`): a retained Phase 21 Task 4 Activation assertion expected the displaced generic `invalid-argument`. The repair requires exit-4 `idempotency-conflict`, exact `content`/`expected_version` fields, identical repeated failure and unchanged Event count. The byte-identical repaired race command terminated `exit_code: 0` with final token `ok  	github.com/mbell-nzl/fid/cmd/fid	251.732s`.
- The first retained Phase 23 conformance selection failed with final token `FAIL` (terminated `exit_code: 1`) because its fake model omitted the new closed details member. Exact nine-field convergence made the byte-identical selection pass with final token `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	10.322s` (terminated `exit_code: 0`).

- The first committed-report lint attempt terminated `exit_code: 3` with final token `FAIL report-lint found 3 failure(s)`: three PASS rows had prose after their termination tokens. A docs-only correction moved that prose outside the rows; lint is rerun against the corrected committed report.
- Attempt 2's first tagged gate failed with token `FAIL` (terminated `exit_code: 1`): `TestP23StopFailureCodeMatrix` lacked a conflict row, `p22RejectActivationFaultConflicts` expected pre-Task-2 invalid-argument for a same-key changed request, and `p22RejectNestedCurrentValues` treated field-map names as closed. Added a real Stop conflict row and exact Activation conflict assertion; retained nested JSON type/digest checks while allowing names to follow written optional input. The focused three-test selection passed and the final byte-identical tagged gate passed with `ok  \tgithub.com/mbell-nzl/fid/cmd/fid/integration\t280.036s` (terminated `exit_code: 0`).
- The second tagged attempt failed with token `FAIL` (terminated `exit_code: 1`) because the Phase 23 amended fake conflict model expected `differing_fields:[""]` for Envelope, Claim and Evidence. Added the four changed public names and kept the fake/real comparison; the focused selection and final byte-identical tagged gate passed.
- An intermediate ordinary gate failed with token `FAIL` (terminated `exit_code: 1`) after an attempted `--actor` Event-validator guard rejected historical synthetic `{}` maps. The guard was removed; real CLI field-name assertions remain in the new corpus. The final byte-identical ordinary gate passed with `ok  \tgithub.com/mbell-nzl/fid/internal/workflow\t(cached)\tcoverage: 62.6% of statements` (terminated `exit_code: 0`).
- No accepted contract, command, budget or threshold was changed. This remains an in-scope Task 2 checkpoint because VAL-002's one-journal fixture is unfinished.

## Concerns, technical debt, or notes for the Tech Lead

- The original checkpoint left most write sites storing `{}` and some renderers returning generic failures; attempt 1 adopted their direct request constructors and Gate, Stop and Issue renderers, and attempt 2 ran the all-leaf behavioral corpus.
- The original checkpoint had changed no retained assertion. This continuation converged three displaced retained assertions after checking their introducing commits, as required by DEC-162/165/166.
- The first red exposed a missing `run.start` CLI conflict mapping. The integration test and focused green prove that one path only; they do not prove the 49-leaf or observation contracts.
- Retained-test provenance: `cmd/fid/activation_test.go` replay assertion was introduced by `0444ccc458f3ac5a720724825f7f81589cb41e09` (`phase-21/task-4`); `cmd/fid/integration/p23_conformance_test.go` fake conflict details by `2e7fd2f9fa05bce648510789453c8d2ed8633b1d` (`phase-23/task-18`); `cmd/fid/integration/p23_portable_recovery_test.go` by `d02614e6c26757076a91d8a12e08dfa328fe699d` (`phase-23/task-16`). Each preserves its replay, exact-envelope, or round-trip subject while asserting Task 2's accepted maps or details.
- The 49-row TSV, live-tree test, executable corpus and per-workspace rebuild checks now pass. The remaining gap is specifically the singular journal with all 49 keyed leaves for VAL-002. Final gate package times were `cmd/fid` 195.085 s and `cmd/fid/integration` 280.036 s; `internal/workflow` and other ordinary packages were cached on the last gate, with this attempt's directly observed non-cached times `internal/workflow` 103.379 s and slowest other `internal/issues` 190.913 s. No ordinary package reached the 480 s warning or 600 s limit; tagged integration was under 900 s.
- Attempt 2 retained-test introducing commits: Stop failure inventory `180f5168e3f1338cb758f8388f1c26fd6f8428d3` (`phase-23/task-19`); Activation changed-request assertion `e7d5963105fa3a40b8d0d687902def63595ade2f` (`phase-22/task-10`); nested Event-value assertion `c8e0061a11c1405e98fd5e4f09612b4f9099af51` (`phase-22/task-1`); amended fake conflict assertion `2e7fd2f9fa05bce648510789453c8d2ed8633b1d` (`phase-23/task-18`). Each keeps the original parity, failure, interruption, or nested-value subject while converging to Task 2's open field maps and exact conflict names.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending commit

## Commit subject

- Pending commit
