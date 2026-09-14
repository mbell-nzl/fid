# Phase 22 - Task 2: Typed Stops and Attention domain - Completion Report

## Workflow outcome

- Task outcome: COMPLETE
- Category: none
- Condition: FQ-02 qualifying Stop/Attention repair and all six unchanged post-edit gates are complete; independent final re-review remains.
- Recovery attempts: 2
- Next action: Root continues queued Task 4 repair; fresh independent final-review-2 follows all seven owners.
- Human attention required: no

## Summary of what was implemented

- Added idempotent Stop raise and resolution operations with immutable original facts, typed conditions, closed codes, subject-bound authority, and qualifying Event checks.
- Added SHA-256 semantic Attention identities, one-open-item reuse with distinct Stop links, exact disposition-Event-set resolution, and atomic Event effects.
- Activated only `stop.raise`, `stop.resolve`, and `attention.resolve` as the approved profile actions while preserving the unavailable `attention.v1` binary capability and absent CLI leaf.
- Reopened the Stop/Attention Event projection family while keeping rejected reads and unavailable recovery branches free of Stop/Attention effects.
- Added the tagged built-binary contract fixture for the accepted unavailable Stop/Attention CLI boundary, without adding deferred leaves.
- Added direct status, warning, and qualifying explicit-Stop evidence for VAL-007.
- Stabilized the retained live-output fixture by passing a fresh coverage directory through its closed child environment, preserving both Phase 21 VAL-036 goldens byte-for-byte.

## Files created or modified

- `internal/workflow/stop.go`
- `internal/workflow/stop_test.go`
- `internal/workflow/attention_test.go`
- `internal/workflow/activation.go`
- `internal/journal/project_effects.go`
- `cmd/fid/contract_surface_test.go`
- `cmd/fid/integration/task22_stop_attention_test.go`
- `cmd/fid/live_output_test.go`
- `project-docs/docs/phases/phase-22/tasks/task-2-done.md`

## Documentation updated

- `project-docs/docs/phases/phase-22/tasks/task-2-done.md`

## Validation Contract Coverage

- VAL-004: PASS - `TestAttentionResolutionRequiresExactDispositionEventsAndKeepsStopsImmutable` at `internal/workflow/attention_test.go:12` resolves both linked Stops and confirms their original `operator-stop` code and details remain unchanged at `:60`.
- VAL-005: PASS - `TestEquivalentStopsShareAttentionButKeepDistinctFacts` at `internal/workflow/stop_test.go:11` proves two Stop rows, one open Attention item, and two links at `:31`; `TestRaiseStopUsesOneOpenAttentionAndRejectsRoutineReadsWithoutEffects` at `:46` verifies the real workflow operation returns distinct Stop IDs with the same Attention ID at `:69`.
- VAL-006: PASS - `TestAttentionResolutionRequiresExactDispositionEventsAndKeepsStopsImmutable` rejects a missing disposition Event set with unchanged journal bytes at `internal/workflow/attention_test.go:42`, then accepts the sorted exact set and produces resolved version 2 at `:50`.
- VAL-007: PASS - `TestRoutineStatusLeavesAttentionAndJournalUnchangedUntilExplicitStop` at `internal/workflow/stop_test.go:82` proves `ActivationStatus` leaves the open Attention count at zero and journal bytes unchanged at `:91-99`, then proves an explicit authorized `RaiseStop` returns both IDs, raises one Attention, and appends the Event at `:101-111`. The tagged `TestStopAttentionExternalBoundaryKeepsRoutineStatusAndWarningsEffectFree` at `cmd/fid/integration/task22_stop_attention_test.go:16` proves fresh-binary routine maintenance `pass` at `:25-29` and `started_without_claim` warning at `:31-36` both leave zero open Attention rows and unchanged journal bytes; it also retains unavailable `attention.v1` and no deferred leaf at `:38-50`.
- VAL-079: PASS - direct in-package positive and rejection branches are `TestEquivalentStopsShareAttentionButKeepDistinctFacts`, `TestRaiseStopUsesOneOpenAttentionAndRejectsRoutineReadsWithoutEffects`, `TestRoutineStatusLeavesAttentionAndJournalUnchangedUntilExplicitStop`, and `TestAttentionResolutionRequiresExactDispositionEventsAndKeepsStopsImmutable`; the initial task baseline was 27.4%, the pre-review repair baseline was 32.9%, and the repair's required final workflow coverage command was 33.7%.

## Tests run

- `go test -race -cover ./internal/workflow` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	26.550s	coverage: 27.4% of statements` (terminated `exit_code: 0`)
  - Pre-edit CF-034 baseline.
- `go test -race ./internal/workflow -run '^TestEquivalentStopsShareAttentionButKeepDistinctFacts$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
  - First red attempt: `Event-v2 "stop.raised" does not permit effect application before its recovery domain is available`.
- `go test -race ./internal/workflow -run '^TestEquivalentStopsShareAttentionButKeepDistinctFacts$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
  - The initial green implementation exposed a test-harness Event-ID error: `parse event id "EVT-stop-1" date`.
- `go test -race ./internal/workflow -run '^TestEquivalentStopsShareAttentionButKeepDistinctFacts$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	2.141s` (terminated `exit_code: 0`)
- `go test -race ./internal/workflow -run '^TestAttentionResolutionRequiresExactDispositionEventsAndKeepsStopsImmutable$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
  - Red attempt: `aggregate attention/ATT-814a6f sequence 2, want 1`.
- `go test -race ./internal/workflow -run '^TestAttentionResolutionRequiresExactDispositionEventsAndKeepsStopsImmutable$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	3.260s` (terminated `exit_code: 0`)
- `/tmp/fid-phase22-task2 --project /tmp/tmp.gBaYkXcrYj --json init --prefix P22` - PASS - `{"contract_version":"fid-cli.v1","event":{"id":"EVT-20260905-000001","occurred_at":"2026-09-05T02:46:37.657134114Z"},"ok":true,"operation":"project.init","result":{"journal_path":".fid/journal/2026-09-05.jsonl","project":{"id":"project-a1b27a80ed0eb60ae3332e573f81d738","issue_prefix":"P22","root":"/tmp/tmp.gBaYkXcrYj"}}}` (terminated `exit_code: 0`)
- `/tmp/fid-phase22-task2 --project /tmp/tmp.gBaYkXcrYj --json capabilities` - PASS - `{"contract_version":"fid-cli.v1","event":null,"ok":true,"operation":"capabilities.get","result":{"items":[{"available":false,"capability_id":"attention.v1"` (terminated `exit_code: 0`)
  - Real CLI entrypoint confirms the accepted unavailable-capability boundary remains intact.
- `go test -race -cover ./internal/workflow` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	(cached)	coverage: 32.9% of statements` (terminated `exit_code: 0`)
  - Required pre-review-repair baseline captured before adding the new tests.
- `go test -race ./internal/workflow -run '^TestRoutineStatusLeavesAttentionAndJournalUnchangedUntilExplicitStop$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	2.315s` (terminated `exit_code: 0`)
  - New direct observable status and explicit Stop evidence was green on first execution because the review finding was missing coverage, not a product defect.
- `go test -tags integration -race ./cmd/fid/integration/... -run '^TestStopAttentionExternalBoundaryKeepsRoutineStatusAndWarningsEffectFree$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
  - Attempt 1: the fresh process rejected the absent `attention` leaf with `exit status 2`, but the fixture incorrectly expected a stderr message.
- `go test -tags integration -race ./cmd/fid/integration/... -run '^TestStopAttentionExternalBoundaryKeepsRoutineStatusAndWarningsEffectFree$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
  - Attempt 2: the process error still exposed only `exit status 2`; the fixture was corrected to assert the observable non-zero no-leaf rejection and root-help omission.
- `go test -tags integration -race ./cmd/fid/integration/... -run '^TestStopAttentionExternalBoundaryKeepsRoutineStatusAndWarningsEffectFree$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	1.856s` (terminated `exit_code: 0`)
  - Attempt 3: the built binary proved routine `pass`, warning `warn`, zero Attention rows, unchanged journal bytes, unavailable capability, and no deferred leaf.
- `go test -race ./internal/workflow -run '^TestRoutineStatusLeavesAttentionAndJournalUnchangedUntilExplicitStop$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	2.323s` (terminated `exit_code: 0`)
  - Post-format focused workflow confirmation.
- `go test -tags integration -race ./cmd/fid/integration/... -run '^TestStopAttentionExternalBoundaryKeepsRoutineStatusAndWarningsEffectFree$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	1.809s` (terminated `exit_code: 0`)
  - Post-format focused tagged confirmation.
- `go test -race -coverprofile=coverage.out ./cmd/fid -run '^TestLiveOutputGoldenFixtures$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	3.410s	coverage: 3.3% of statements` (terminated `exit_code: 0`)
  - Diagnostic only: the retained golden passes in isolation, so it does not replace the failed required repository-wide receipt.
- `go test -race -coverprofile=coverage.out ./cmd/fid` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	166.657s	coverage: 77.1% of statements` (terminated `exit_code: 0`)
  - Diagnostic only: package-level reproduction passed after the failed repository-wide receipt; it is not used to close the required gate.
The pre-repair DEC-189 mechanism diagnostic ran this exact sequence from the repository root:

```sh
B=$(mktemp -d)
go test -c -cover -o "$B/covbin" ./cmd/fid
env -i "$B/covbin" -test.run='^TestLiveGoldenChild$' -- --live-golden-child
```

Its child receipt ended with `warning: GOCOVERDIR not set, no coverage data emitted`; the expected child termination was observed as `exit_code: 23`.
- `go test -race -coverprofile=coverage.out ./cmd/fid -run '^TestLiveOutputGoldenFixtures$' -count=20` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	30.436s	coverage: 3.3% of statements` (terminated `exit_code: 0`)
  - Post-repair focused green confirmation across twenty fixture executions.

## Exact quality-gate results

- `go test -race -cover ./internal/workflow` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	(cached)	coverage: 34.7% of statements` (terminated `exit_code: 0`)
  - Final CF-034 measurement; earlier after-implementation attempt was `34.6%` (terminated `exit_code: 0`).
- `go test -race -cover ./internal/workflow` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	32.901s	coverage: 33.7% of statements` (terminated `exit_code: 0`)
  - Review-1 repair final workflow coverage measurement.
- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - FAIL - `* unused: 1` (terminated `exit_code: 1`)
  - Attempt 1: two unchecked `rows.Close` returns and one unused helper in `internal/workflow/stop.go`.
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
  - Attempt 2 after the in-scope lint repair.
- `golangci-lint run ./...` - FAIL - `* typecheck: 1` (terminated `exit_code: 1`)
  - Attempt 3: the converged retained test referenced `matches` before declaration.
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
  - Attempt 4 after correcting the retained test.
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
  - Review-1 repair attempt.
- `go test -race -coverprofile=coverage.out ./...` - FAIL - `FAIL` (terminated `exit_code: 1`)
  - Attempt 1: retained `TestPhase21AttentionBoundary` rejected the accepted task-2 Attention implementation surface.
- `go test -race -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	(cached)	coverage: 34.7% of statements` (terminated `exit_code: 0`)
  - Attempt 2 after subject-preserving convergence.
- `go test -race -coverprofile=coverage.out ./...` - FAIL - `FAIL` (terminated `exit_code: 1`)
  - Review-1 repair attempt: retained `TestLiveOutputGoldenFixtures` expected its coverage-child `GOCOVERDIR not set` warning and 134 observed bytes, but the full run instead captured 200 `x` bytes. The exact failure is preserved; no task-2 review-repair path authorizes changing this golden or its live-output implementation.
- `go tool cover -func=coverage.out` - PASS - `total:										(statements)				61.6%` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total:										(statements)				61.1%` (terminated `exit_code: 0`)
  - Review-1 repair attempt; the coverage profile was emitted despite the preceding failed test command.
- `go test -tags integration -race ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	8.157s` (terminated `exit_code: 0`)
- `go test -tags integration -race ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	8.188s` (terminated `exit_code: 0`)
  - Review-1 repair attempt.
- `go test -race -cover ./internal/workflow` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	(cached)	coverage: 33.7% of statements` (terminated `exit_code: 0`)
  - Attempt-2 required workflow coverage measurement.
- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
  - Attempt-2 required format check.
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
  - Attempt-2 required lint gate.
- `go test -race -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	(cached)	coverage: 33.7% of statements` (terminated `exit_code: 0`)
  - Attempt-2 byte-identical required full gate; `cmd/fid` completed at `192.255s` with `77.1%` coverage and all packages passed.
- `go tool cover -func=coverage.out` - PASS - `total:										(statements)				61.1%` (terminated `exit_code: 0`)
  - Attempt-2 coverage profile result from the passing full gate.
- `go test -tags integration -race ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	10.782s` (terminated `exit_code: 0`)
  - Attempt-2 required tagged integration gate.
- `git status --porcelain -- cmd/fid/testdata/val-036` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
  - Recovery-specific proof that both golden fixtures remain byte-identical.

## Deviations from spec

- No contract, schema, command, tolerance, or policy changed. The initial lint failure was caused by task-2 implementation hygiene and was repaired in scope; the byte-identical later lint attempt passed.
- The first full race/coverage attempt failed because `TestPhase21AttentionBoundary`, introduced by `7277139ee5d5dde4cd3627335d39b14f616f46a1` (`feat(evidence): evaluate trusted receipt evidence`, Refs: phase-21/task-17), prohibited every production Attention reference except the unavailable capability row. Task 2's accepted Stop/Attention domain and public-interface section 2.1 displaced that broad assertion. The converged test still proves the exact unavailable `attention.v1` row, absent `fid attention` leaf, no rejected-read effects, and limits Attention references to `activation.go`, `capabilities.go`, and task-2 `stop.go`. The byte-identical full gate later passed.
- The first intended red test established unavailable Stop Event effects. Its next attempt exposed a malformed test Event ID; correcting the fixture preserved the assertion. The attention-resolution red test then found and repaired the aggregate sequence from 2 to 1.
- Review-1 repair added only task-2 evidence. The first two focused tagged fixture attempts failed because the external unknown-command process reports `exit status 2` without text on stdout/stderr; the fixture now asserts that actual rejection and the absent help leaf while preserving the no-leaf contract. The required full race/coverage command then failed in the retained `cmd/fid` live-output golden despite the isolated and package-level diagnostic receipts passing. No causal repair is authorized in this bounded review rework, so this report remains RECOVERY_REQUIRED rather than claiming a later green as closure.
- Mandatory deviation note for DEC-189: first failure token `FAIL` (terminated `exit_code: 1`) came from the byte-identical `go test -race -coverprofile=coverage.out ./...` gate in retained `TestLiveOutputGoldenFixtures`. Root cause: the coverage-instrumented child inherited the fixture's closed environment without `GOCOVERDIR`, so its exit hook could add `warning: GOCOVERDIR not set, no coverage data emitted` to the captured stderr before the output ceiling. Repair: `cmd/fid/live_output_test.go` now passes a fresh `t.TempDir()` as the sorted literal `GOCOVERDIR` entry between `GATE_TOKEN` and `LITERAL`. Later result: the 20-run focused coverage fixture and the byte-identical full gate both passed, while `cmd/fid/testdata/val-036` remained unchanged.

## Concerns, technical debt, or notes for the Tech Lead

- No CLI wiring was added because task 2 expressly excludes it. The direct workflow API and Event effects are ready for the later public CLI owner; `attention.v1` remains unavailable as required.
- The scoped repair is complete and all task commands pass. A fresh independent review-1 re-review remains required; this report does not claim that QA outcome.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `f973bca7ea146f00c0e16ca5e5117a8d13f769b9`

## Commit subject

- `fix(workflow): enforce qualifying Stop and Attention facts`

## Root report display correction — 2026-09-08

The root delegate decoded escaped tab characters in historical Go-output
receipts following the QA-84 preflight. The original Developer sessions remain
the source of those receipts, including every recorded outcome, duration,
coverage value and observed termination. This formatting correction records
no new execution or attestation. Original encoded text remains in Git history.

## Final-QA FQ-02 same-task continuation — 2026-09-09

Authority: review-7, final-review.md FQ-02, task-2.md and contracts.md §5.
Root is the sole writer under the Human's delegate instruction. This appendix
supersedes the original readiness statements; every historical receipt and
failed attempt above remains unchanged. The prior primary product binding was
daca8392a580ebdd2b6a8d3dab011910e39a779e, `test(fid): stabilize live output coverage fixture`.

### Behavior and assertion evidence

- VAL-004/006/080: Stop resolution now checks each of the six closed conditions
  against its exact durable fact and current subject state. Executor tests reject
  live, unknown, wrong-Run, missing, changed-executor and changed-liveness facts.
  Recovery tests bind the predecessor. Activation tests require terminal resource
  agreement; capability tests additionally require the selected backend and named
  compiled capability. Budget tests bind the named Envelope chain and reject
  unknown, exhausted, wrong-budget, wrong-Envelope and later-exhausted totals.
  Work-ended tests require the terminal Run or an Envelope with no active Run.
- The envelope.completed positive/active-Run negative uses an explicitly admitted
  exact two-Effect Event through the journal kernel after actual CloseRun. It is
  predicate evidence, not a claim to execute the full CompleteEnvelope journey.
- Attention uses the sorted unique containing-disposition Event set and
  revalidates every underlying Stop fact. Actual activation recovery disposes two
  Stops in one Event; one ID resolves the Attention, while corruption of either
  underlying fact rejects. Missing, extra, duplicate, wrong-subject and
  underlying-only IDs leave journal, Attention, dispositions and idempotency
  unchanged. Later budget exhaustion also invalidates the prior disposition.
- VAL-007 retained routine-read/warning assertions remain unchanged. VAL-079:
  fresh pre-edit workflow coverage was 52.9%; final direct package gate was 60.6%.
- Built `/tmp/fid-fq2` from the repair source and executed actual JSON CLI
  observation/Stop operations in a fresh copy of the retained recovery checkpoint.
  A live observation returns resolution-invalid/precondition/exit 4 with no Event,
  unchanged journal bytes and unresolved Stop. A subsequent terminated observation
  resolves successfully. Child argv/stdout/stderr/return codes are retained in
  `.cache/phase-22-fq2-cli-digut1t_/receipts.json`; originals were not mutated.

### Files and provenance

- Product: internal/workflow/stop.go and new stop_conditions.go. Exact qualifying
  helpers replace broad Effect-subject matching; disposition Event IDs are
  deduplicated without skipping any Stop. Envelope authorization joins canonical
  envelope_issue_selection instead of querying a nonexistent primary_issue_id
  column. No schema or authority semantics changed.
- Tests: new stop_conditions_test.go; attention_test.go; resume_test.go.
  The obsolete updated_at-only activation.finalized helper and Envelope SQL
  defect originated in ac827959f8359473cbd6546ac6bea1ff1fd20751, Phase22/Task2.
  The retained resume helper caller originated in
  2339dc9ebd8d0e6b5502ca32e1562892cf704ce6, Phase22/Task12. Its independently-open
  Attention assertion is preserved with the actual Prepare/Stage/Observe/Commit/
  Finalize fixture. The ledger's 00:13 matching-fixture correction records scope;
  Task12's original report is unchanged. All rows execute under the required
  workflow/full race gates and the focused commands below.

### Mandatory failed-attempt deviations

- Intended red commands ended with `FAIL` (exit_code: 1): live/unknown executor
  facts were accepted; activation/resource/capability current state was unchecked;
  named/current budget and Envelope work-ended predicates were incorrect; shared
  containing Event IDs were counted twice. Corresponding byte-identical focused
  commands passed after their causal in-scope repairs, as recorded below.
- Envelope Stop authority red ended `FAIL` (exit_code: 1), with `SQL logic error:
  no such column: primary_issue_id (1)`. The canonical primary-Issue join repaired
  it; the identical command passed at 3.704s.
- Retained Attention fixture red ended `FAIL` (exit_code: 1): its incomplete
  activation did not satisfy the real agreement. Replacing it exposed
  `undefined: commitActivationRecoveryFact` in resume_test.go (exit_code: 1).
  The provenance-bound caller convergence preserves both subjects; identical
  Attention command then passed at 5.544s.
- Test-construction attempts remain failed evidence: unsupported fixture phase22
  was corrected to the actual supported phase21 profile; a wrong-Run fixture
  originally mismatched Stop and Condition subjects before reaching the intended
  check and now cites the predecessor fact from a real successor; the minimal
  recovery fixture has no Obligation and cannot run full CompleteEnvelope, so its
  panic-producing setup was replaced by the explicit admitted Event predicate
  control described above. The combined focused command failed at 104.470s and
  109.534s (both final token `FAIL`, exit_code: 1), then passed at 108.770s.
- The first actual CLI probe command terminated exit_code: 1 with
  `json.decoder.JSONDecodeError: Expecting value: line 1 column 1 (char 0)` after
  its live rejection had passed: the script invented stop.get. Using the accepted
  stop.list leaf with exact ID selection repaired the probe; the byte-identical
  Python invocation passed. Failed raw receipts remain in
  `.cache/phase-22-fq2-cli-rnlauihv/receipts.json`.
- No dependency, policy, contract, schema, protected gate, tolerance or command
  changed. FQ-02 implementation closure remains subject to independent final
  re-review after all seven owners; Human judgments remain pending.

### Continuation command receipts

Working directory: `/home/mark/projects/fid`. Tokens and terminations below were
captured by root in this recording session; repeated commands are separate
immutable attempts. Required post-edit commands retain their original order.


- `go test -race -cover ./internal/workflow` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	114.339s	coverage: 52.9% of statements` (terminated `exit_code: 0`)
  - baseline workflow coverage.
- `go test -race ./internal/workflow -run '^TestStopResolutionRequiresTerminatedExecutorFact$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
  - executor condition behavioral red.
- `go test -race ./internal/workflow -run '^TestStopResolutionRequiresTerminatedExecutorFact$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	9.057s` (terminated `exit_code: 0`)
  - executor condition behavioral green.
- `go test -race ./internal/workflow -run '^TestStopResolutionChecksCurrentActivationAgreementAndCapability$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
  - activation/capability condition red.
- `go test -race ./internal/workflow -run '^TestStopResolutionChecksCurrentActivationAgreementAndCapability$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
  - activation/capability corrected fixture red.
- `go test -race ./internal/workflow -run '^TestStopResolutionChecksCurrentActivationAgreementAndCapability$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	29.218s` (terminated `exit_code: 0`)
  - activation/capability condition green.
- `go test -race ./internal/workflow -run '^TestStopResolution(RecoveryFactNamesThePredecessor|RequiresNamedKnownNonExhaustedBudget|WorkEndedRequiresCurrentTerminalSubject)$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
  - recovery/budget/work-ended condition red.
- `go test -race ./internal/workflow -run '^TestStopResolution(RecoveryFactNamesThePredecessor|RequiresNamedKnownNonExhaustedBudget|WorkEndedRequiresCurrentTerminalSubject)$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	29.441s` (terminated `exit_code: 0`)
  - recovery/budget/work-ended condition green.
- `go test -race ./internal/workflow -run '^TestAttentionResolutionRequiresExactDispositionEventsAndKeepsStopsImmutable$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
  - retained Attention fixture convergence red.
- `go test -race ./internal/workflow -run '^TestAttentionResolutionRequiresExactDispositionEventsAndKeepsStopsImmutable$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
  - retained Attention fixture convergence green.
- `go test -race ./internal/workflow -run '^TestAttentionResolutionRequiresExactDispositionEventsAndKeepsStopsImmutable$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	5.544s` (terminated `exit_code: 0`)
  - retained Attention fixture green after caller convergence.
- `go test -race ./internal/workflow -run '^TestAttentionResolutionDeduplicatesContainingEventAndRevalidatesEveryStop$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
  - shared disposition Event red.
- `go test -race ./internal/workflow -run '^TestAttentionResolutionDeduplicatesContainingEventAndRevalidatesEveryStop$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	14.135s` (terminated `exit_code: 0`)
  - shared disposition Event green.
- `go test -race ./internal/workflow -run '^TestEnvelopeStopResolutionUsesPrimaryIssueGrantAndCurrentBudget$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
  - Envelope Stop authority red.
- `go test -race ./internal/workflow -run '^TestEnvelopeStopResolutionUsesPrimaryIssueGrantAndCurrentBudget$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	3.704s` (terminated `exit_code: 0`)
  - Envelope Stop authority green.
- `go test -race ./internal/workflow -run '^(TestAttentionResolution|TestStopResolution|TestEnvelopeStopResolution|TestResumeKeepsOpenAttention)' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
  - complete current FQ02 focused matrix.
- `go test -race ./internal/workflow -run '^(TestAttentionResolution|TestStopResolution|TestEnvelopeStopResolution|TestResumeKeepsOpenAttention)' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
  - complete FQ02 matrix after wrong-Run fixture correction.
- `go test -race ./internal/workflow -run '^TestStopResolutionWorkEndedAcceptsCompletedEnvelopeFact$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	3.179s` (terminated `exit_code: 0`)
  - completed Envelope admitted-fact control.
- `go test -race ./internal/workflow -run '^(TestAttentionResolution|TestStopResolution|TestEnvelopeStopResolution|TestResumeKeepsOpenAttention)' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	108.770s` (terminated `exit_code: 0`)
  - complete FQ02 final focused matrix.
- `go build -o /tmp/fid-fq2 ./cmd/fid` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
  - FQ02 real CLI build.
- `python3 /tmp/fid-fq2-cli.py` - FAIL - `json.decoder.JSONDecodeError: Expecting value: line 1 column 1 (char 0)` (terminated `exit_code: 1`)
  - FQ02 actual CLI regression.
- `python3 /tmp/fid-fq2-cli.py` - PASS - `PASS stop condition CLI: live rejected exit=4; terminated resolved exit=0; journal unchanged on rejection; receipts=/home/mark/projects/fid/.cache/phase-22-fq2-cli-digut1t_/receipts.json` (terminated `exit_code: 0`)
  - FQ02 actual CLI regression corrected list leaf.
- `go test -race -cover ./internal/workflow` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	212.427s	coverage: 60.6% of statements` (terminated `exit_code: 0`)
  - required post-gate 1 workflow coverage.
- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
  - required post-gate.
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
  - required post-gate.
- `go test -race -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	238.821s	coverage: 60.6% of statements` (terminated `exit_code: 0`)
  - required full coverage.
- `go tool cover -func=coverage.out` - PASS - `total:										(statements)					64.9%` (terminated `exit_code: 0`)
  - required final gates.
- `go test -tags integration -race ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	100.809s` (terminated `exit_code: 0`)
  - required final gates.
