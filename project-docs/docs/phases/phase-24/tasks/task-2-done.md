# Phase 24 - Task 2: Input-as-written digests and named conflict fields - Completion Report

## Workflow outcome

- Task outcome: RECOVERY_REQUIRED
- Category: session-continuation
- Condition: Shared-helper adoption and seven observing-leaf replay checks progressed, but the 49-leaf behavioral conflict/retry corpus, every-family two-field and Actor cases, all-leaf rebuild parity, optional-presence audit, and final ordered gate remain unfinished.
- Recovery attempts: 1
- Next action: Resume Task 2 from this checkpoint. Execute the 49 rows in `cmd/fid/integration/testdata/p24_keyed_leaves.tsv` as behavioral conflict and canonical-retry cases; add every-family two-field/Actor rows, audit input presence, prove all-leaf rebuild parity, and run the eight exact gates in order.
- Human attention required: no

## Summary of what was implemented

- Added a shared request and per-field SHA-256 helper with sorted field-difference calculation. The helper hashes the family's supplied digest components and derives public field digests from the canonical request object plus `--actor`.
- Adopted the helper for `run.start` through `mutateRun`. Its same-key `claim_version` conflict now carries the stored and supplied field maps and renders an exit-4 `idempotency-conflict` naming `claim_version`.
- At the prior checkpoint, the accepted `differing_fields` member was added to the typed conflict details definition and explicit conflict renderers were moved to one details builder. This continuation adopted the remaining direct request sites and family renderers.
- Continuation attempt 1 adopted the shared helper and field maps in Issues, Labels, Links, Workflow, Evidence, Activation, Claims and Review; removed remaining family-local digest helpers; and added conflict mapping for Gate, Stop and Issue-family renderers. `gate.run` passes real maps through execution-started, Receipt update precondition and committed denial.
- Added process tests for Review optional presence and the DEC-242 exception; seven-leaf replay after HEAD/source changes and through unmerged/deleted-file states with new-key controls; selected family conflicts; and Gate receipt/denial maps. A live-tree inventory lists 49 keyed leaves, but executable conflict/retry rows remain unfinished.

## Files created or modified

- `internal/journal/request_digest.go`; `internal/journal/request_digest_test.go`
- `internal/workflow/run.go`; `internal/cli/codes.go`
- `cmd/fid/idempotency_conflict.go`; `cmd/fid/integration/p24_idempotency_test.go`
- `cmd/fid/claim.go`; `cmd/fid/concern.go`; `cmd/fid/decision.go`; `cmd/fid/deviation.go`; `cmd/fid/envelope.go`; `cmd/fid/evidence.go`; `cmd/fid/finding.go`; `cmd/fid/human_outcome.go`; `cmd/fid/review.go`; `cmd/fid/run.go`
- `project-docs/docs/phases/phase-24/tasks/task-2-done.md`
- Continuation source: `internal/{issues,labels,links,workflow,evidence,claims,review}/` digest call sites and attribution helpers; `cmd/fid/{activation,gate,issue_mutation_response,stop}.go`.
- Continuation tests: `cmd/fid/{activation_test,p24_keyed_inventory_test}.go`; `cmd/fid/integration/{p23_conformance_test,p23_portable_recovery_test,p24_idempotency_test}.go`; `cmd/fid/integration/testdata/p24_keyed_leaves.tsv`.
- `project-docs/specs/06_DESIGN_PLAN_CLI.md`.

## Documentation updated

- `project-docs/specs/06_DESIGN_PLAN_CLI.md` describes input-as-written hashing, DEC-242, named field digests and `differing_fields`. This report retains the original attempt's failed evidence and records the continuation.

## Validation Contract Coverage

- VAL-002: BLOCKED - `TestPortableRoundTripCoversEveryPhase23Table` (`cmd/fid/integration/p23_portable_recovery_test.go:34`) passes with populated maps for its covered families after portable recovery; the required journal with one keyed mutation per all 49 leaves and map parity remains unfinished.
- VAL-005: BLOCKED - The required 100k benchmark has not run on a final tree.
- VAL-006: BLOCKED - `TestP24KeyedLeafCorpusMatchesLiveTree` (`cmd/fid/p24_keyed_inventory_test.go:14`) binds 49 live keyed leaves to planned one-field rows. `TestP24ObservingJudgmentLeavesNameChangedPublicField` (`cmd/fid/integration/p24_idempotency_test.go:209`) and selected Issue, Stop, Claim, Activation and Gate tests assert exact fields and no Event, but executable rows for every leaf and 20 family two-field/Actor cases remain unfinished.
- VAL-007: BLOCKED - `TestP24ObservingJudgmentRetriesAfterSourceAndHEADChange` (`cmd/fid/integration/p24_idempotency_test.go:260`) covers seven replay paths; reordered JSON/flag canonical retries for all 49 keyed leaves remain unfinished.
- VAL-008: Focused PASS - `TestP24ReviewOptionalPresenceAndEmptyException` (`cmd/fid/integration/p24_idempotency_test.go:178`) proves nonempty-versus-omitted conflict naming `inapplicable_findings` and `[]`-versus-omitted replay without an Event. Final tagged gate remains pending.
- VAL-009: Focused PASS - `TestP24ObservingJudgmentRetriesAfterSourceAndHEADChange` (`cmd/fid/integration/p24_idempotency_test.go:260`) retries all seven after committing a changed `base.txt` and advancing HEAD, comparing replay envelopes and Event counts. Final tagged gate remains pending.
- VAL-052: BLOCKED - Repaired touched-package race selection terminated green with `cmd/fid` 251.732 s, `internal/workflow` 111.563 s and slowest other ordinary package `internal/issues` 161.320 s, all below the 480 s warning. Required full ordinary/tagged gate times remain uncaptured.
- VAL-066: Focused PASS - `TestP24ObservingJudgmentReplayNeverObserves` (`cmd/fid/integration/p24_idempotency_test.go:309`) covers all seven in unmerged states, source-file deletion where applicable, no-Event replay, observable new-key preflight success and unobservable new-key refusal. Final tagged gate remains pending.

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

## Exact quality-gate results

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- Formatting result applies to the checkpoint tree; rerun on the final Task 2 tree.
- `golangci-lint run ./...` - BLOCKED - Final-tree gate deferred until the remaining Task 2 implementation and corpus are complete.
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` - BLOCKED - Final-tree gate deferred until the remaining Task 2 implementation and corpus are complete.
- `go tool cover -func=coverage.out` - BLOCKED - Coverage gate depends on the unfinished ordinary test gate.
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...` - BLOCKED - Final-tree gate deferred until the remaining Task 2 implementation and corpus are complete.
- `go test -run '^$' -bench '^BenchmarkRebuildWorkspace100kEvents$' -benchtime 1x ./internal/journal/` - BLOCKED - Final-tree benchmark deferred until the remaining Task 2 implementation and corpus are complete.
- `rg -n "RequestDigest:" internal cmd -g '!*_test.go'` - BLOCKED - Required byte-identical final registry command has not run on a final Task 2 tree; direct idempotency request constructors now carry maps, and Gate Event effects carry the map.
- `rg -n "func [a-zA-Z]*(RequestDigest|requestDigest|judgmentDigest)\(" internal -g '!*_test.go'` - BLOCKED - Required byte-identical final registry command has not run on a final Task 2 tree; a diagnostic currently lists only `ComputeRequestDigest`.

## Deviations from spec

- No accepted-contract deviation. This is a same-task checkpoint, not task completion.
- The first touched-package race attempt failed with final token `FAIL` (terminated `exit_code: 1`): a retained Phase 21 Task 4 Activation assertion expected the displaced generic `invalid-argument`. The repair requires exit-4 `idempotency-conflict`, exact `content`/`expected_version` fields, identical repeated failure and unchanged Event count. The byte-identical repaired race command terminated `exit_code: 0` with final token `ok  	github.com/mbell-nzl/fid/cmd/fid	251.732s`.
- The first retained Phase 23 conformance selection failed with final token `FAIL` (terminated `exit_code: 1`) because its fake model omitted the new closed details member. Exact nine-field convergence made the byte-identical selection pass with final token `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	10.322s` (terminated `exit_code: 0`).

- The first committed-report lint attempt terminated `exit_code: 3` with final token `FAIL report-lint found 3 failure(s)`: three PASS rows had prose after their termination tokens. A docs-only correction moved that prose outside the rows; lint is rerun against the corrected committed report.

## Concerns, technical debt, or notes for the Tech Lead

- The original checkpoint left most write sites storing `{}` and some renderers returning generic failures; this continuation adopted their direct request constructors and Gate, Stop and Issue renderers. The all-leaf behavioral corpus must still verify these paths.
- The original checkpoint had changed no retained assertion. This continuation converged three displaced retained assertions after checking their introducing commits, as required by DEC-162/165/166.
- The first red exposed a missing `run.start` CLI conflict mapping. The integration test and focused green prove that one path only; they do not prove the 49-leaf or observation contracts.
- Retained-test provenance: `cmd/fid/activation_test.go` replay assertion was introduced by `0444ccc458f3ac5a720724825f7f81589cb41e09` (`phase-21/task-4`); `cmd/fid/integration/p23_conformance_test.go` fake conflict details by `2e7fd2f9fa05bce648510789453c8d2ed8633b1d` (`phase-23/task-18`); `cmd/fid/integration/p23_portable_recovery_test.go` by `d02614e6c26757076a91d8a12e08dfa328fe699d` (`phase-23/task-16`). Each preserves its replay, exact-envelope, or round-trip subject while asserting Task 2's accepted maps or details.
- The 49-row TSV and live-tree test prove planning completeness only. The continuation must execute each row's changed-input conflict and canonical retry and prove every-family two-field/Actor cases. Several flag/positional and typed `--input` paths still need an explicit input-as-written optional-presence audit.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `1c94956a5c0b79ad0ab1ff325b7097ac621ed2b2`

## Commit subject

- `feat(idempotency): checkpoint shared digest adoption across families`
