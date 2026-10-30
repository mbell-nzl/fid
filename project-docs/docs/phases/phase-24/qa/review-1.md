# Phase 24 - Mid-Phase Review 1 (Tasks 1-3)

Reviewer: qa-reviewer
Date: 2026-10-05
Tasks covered: phase-24/task-1, phase-24/task-2, phase-24/task-3
Review type: INITIAL
Prior review: none

## Verdict

REPAIR. Window A does not advance. Five evidence gaps have bounded owners in the reviewed tasks; none calls for a protected decision.

## Scope Reviewed

- Task specs and completion reports: `project-docs/docs/phases/phase-24/tasks/task-{1,2,3}{,-done}.md`.
- `qa/validation-contract.md`, `qa/qa-plan.md` sections 8 and 12.7 and its in-window evidence rows, `input-inventory.md` R1-R20, the phase run ledger and triage log, current context, `AGENTS.md`, `08_QUALITY_PLAN.md`, and `09_AI_WORKFLOW.md` §§2.7, 8.
- Verified with `git rev-parse --verify '<hash>^{commit}'` (each terminated `exit_code: 0`): task 1 `58352c73edf81ac8e89ffb7d7d824a67def3e3c4`, task 2 `95c738f23c34aef1a9e8c6630f1d88c660648591`, task 3 report `290d8cc2c36df820a2fe4c866c11df12ff09666c`. Reviewed their diffs and the task-3 implementation checkpoints `04b6b0992fe20a02df402a183269a9847ddae5ab`, `2ee33ba30ac698de515ab2c1f6142f2d310615e3`, `c942419596aad3fd198cb2c27cf8d06fc5238745`.
- Held-out scenario and evaluator material was not opened.

## Summary

- Review outcome: REPAIR
- Must-fix findings: 5
- Improvement findings: 0
- Blocker candidates: 0
- Existing escalation references: none
- Covered validation assertions: VAL-001-016, VAL-052, VAL-058-060, and task 2's seven keyed-leaf rows of VAL-066

## Validation Contract Coverage

| VAL ID | Claimed by | Status | Evidence and limit |
|---|---|---|---|
| VAL-001 | task-1 | ON TRACK | `internal/db/migration0019_test.go:9` compares schema-18 and schema-19 objects; focused rerun passed. |
| VAL-002 | tasks 1, 2 | ON TRACK | `p24_single_journal_test.go:15,196` requires all 49 keyed operation scopes in one workspace; `p24_map_parity_test.go:11` compares live/rebuilt map bytes. Focused rerun passed. |
| VAL-003 | task-1 | NOT ON TRACK | Insert-map fixtures cover valid, missing, array and non-digest states; receipt-update guard mismatch and per-invalid-state CLI code/exit/line evidence are absent (F1). |
| VAL-004 | task-1 | NOT ON TRACK | Schema-19 import and schema-18 rejection preserve the target, but the test does not assert the required `manifest.json` check path (F2). |
| VAL-005 | tasks 1, 2 | ON TRACK | Reports give 22,373 and 22,449 ms/rebuild. Independent serial sample: 21,531 ms/rebuild; phase-end three fresh samples remain due. |
| VAL-006 | task-2 | ON TRACK | Live-tree and executable-row checks cover 49 keyed leaves; `TestP24` integration rerun passed. Later input-shape convergence remains assigned to tasks 4-9. |
| VAL-007 | task-2 | ON TRACK | `p24ExerciseCapturedLeaf` and judgment corpus compare original response bytes with only `replayed:true`; `TestP24` rerun passed. |
| VAL-008 | task-2 | ON TRACK | `TestP24ReviewOptionalPresenceAndEmptyException` covers nonempty-versus-omitted conflict and empty-versus-omitted replay; Source-digest case belongs to task 4. |
| VAL-009 | task-2 | ON TRACK | `TestP24ObservingJudgmentRetriesAfterSourceAndHEADChange` covers the seven current review leaves; later observing leaves belong to tasks 4-8. |
| VAL-010 | task-3 | ON TRACK | Recursive Short test passed; rebuilt live inventory has 92 leaves and zero without Short. |
| VAL-011 | task-3 | ON TRACK | Three descriptors plus a 36-leaf fence equal the 39 live input leaves; seeded type mismatch is rejected. Completion belongs to task 9. |
| VAL-012 | task-3 | ON TRACK | Three adopted help surfaces and the pinned full `stop raise` help are exercised; later leaves belong to tasks 4-9. |
| VAL-013 | task-3 | ON TRACK | Exact three template strings passed; independent built CLI emitted one JSON object per leaf. Later leaves belong to tasks 4-9. |
| VAL-014 | task-3 | ON TRACK | Three exact templates decode into their reflected request types; later leaves belong to tasks 4-9. |
| VAL-015 | task-3 | NOT ON TRACK | The Developer test checks only `attention resolve` and does not assert complete details or absence of effects for all combinations (F4). An independent CLI spot check passed but does not supply the required task-owned corpus. |
| VAL-016 | task-3 | NOT ON TRACK | Per-leaf structural tests assert paths and nonempty allowed lists, not the exact allowed/value details and no-Event result required by EVID-14 (F5). |
| VAL-052 | tasks 1-3 | ON TRACK | Reported ordinary and tagged package times are below 600/900 s; table below. |
| VAL-058 | task-1 | ON TRACK | Same migration test compares the down-migrated `sqlite_master` objects to schema 18; focused rerun passed. |
| VAL-059 | task-3 | ON TRACK | `attention resolve --help` names Attention resolution; neutral shared flag text and the three adopted helps contain no Run-input wording. |
| VAL-060 | task-3 | ON TRACK | Three template runs from an empty non-project directory leave it empty; independent built-binary spot check reproduced this. Later leaves belong to tasks 4-9. |
| VAL-066 | task-2 rows | NOT ON TRACK | Seven unmerged-path replay/control pairs pass, and source deletion retries pass; the deleted-source branch lacks a paired new-key control (F3). Later leaves belong to tasks 4-8. |

## Gate And Evidence Checks

All commands below use working directory `/home/mark/projects/fid`. The reports' final PASS rows match every byte-identical command in each task's `## Exact Commands To Run`: task 1 has seven, task 2 eight, task 3 five. QA-plan section 8 and the run ledger declare no additional task-level command. Each report row has a command-output token and an observed `exit_code`; no required command was silently skipped or narrowed. `project-docs/tools/report-lint.sh <report> <spec>` was rerun for each task: `PASS report-lint completed with no failures`, each terminated `exit_code: 0`.

| Command (task) | Final report token and observed termination |
|---|---|
| `gofmt -l cmd internal` (1, 2, 3) | `exit 0 (no output)`; each terminated `exit_code: 0`. |
| `golangci-lint run ./...` (1, 2, 3) | `0 issues.`; each terminated `exit_code: 0`. Task 3 preserves earlier `* staticcheck: 1` (terminated `exit_code: 1`). |
| `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` (1, 2, 3) | Final last lines: task 1 `ok .../internal/workflow 72.492s coverage: 62.7%`; tasks 2/3 `ok .../internal/workflow (cached) coverage: 62.6%`. Each terminated `exit_code: 0`. Task 1 preserves a prior `FAIL .../internal/workflow 103.620s` (`exit_code: 1`); task 2 preserves its prior `FAIL` (`exit_code: 1`). Full unabridged tokens are in each report's Exact quality-gate results. |
| `go tool cover -func=coverage.out` (1, 2, 3) | Final `total:` statement coverage 60.4%, 60.5%, 60.2% respectively; each terminated `exit_code: 0`. |
| `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...` (1, 2, 3) | Final `ok .../cmd/fid/integration` at 253.353s, 248.512s, 302.158s respectively; each terminated `exit_code: 0`. Task 2 retains two earlier `FAIL` attempts; task 3 retains two `BLOCKED - FAIL` attempts, all terminated `exit_code: 1`. |
| `go test -run '^$' -bench '^BenchmarkRebuildWorkspace100kEvents$' -benchtime 1x ./internal/journal/` (1, 2) | Reported `22373 ms/rebuild` / `22449 ms/rebuild` and final `ok .../internal/journal 25.020s` / `25.093s`; each terminated `exit_code: 0`. |
| `rg -n "idempotency_results" migrations internal -g '!*_test.go'` (1) | Last reported line `internal/workflow/activation.go:845: ... idempotency_results ...`; terminated `exit_code: 0`. |
| `rg -n "RequestDigest:" internal cmd -g '!*_test.go'` (2) | Last reported line `internal/workflow/run_budget.go:61: ... RequestDigest: digest ...`; terminated `exit_code: 0`. |
| `rg -n "func [a-zA-Z]*(RequestDigest|requestDigest|judgmentDigest)\\(" internal -g '!*_test.go'` (2) | Last reported line `internal/journal/request_digest.go:14:func ComputeRequestDigest(...)`; terminated `exit_code: 0`. |

Independent named-test reruns and process evidence:

| Command or probe | Reproduced result |
|---|---|
| `go test -count=1 -run '^(TestMigration0019AddsOnlyFieldDigestColumn|TestCurrentIdempotencyFieldDigestsReplayFailClosed|TestSchema19ArchiveImportsAndSchema18ManifestFailsClosed|TestInputValidationOrderAndSensitiveValues)$' ./internal/db ./internal/journal ./internal/export ./internal/cli` | Last `ok  github.com/mbell-nzl/fid/internal/cli 0.008s`; all four packages printed `ok`; terminated `exit_code: 0`. Covers valid/invalid map and archive states, but not F1/F2. |
| `go test -count=1 -run '^(TestAttentionResolveIntegerBoundWritesJSONFailure|TestInputDescriptorRegistryMatchesLiveTreeAndTypes|TestEveryLiveLeafHasOneLineShort|TestAdoptedInputHelpAndExamples|TestExampleRejectsOtherInputsInFixedOrder|TestAdoptedStructuralValidationPaths|TestStructuralFailuresCoverEachAdoptedLeaf|TestP24EveryLiveKeyedLeafHasExecutableCorpusRow)$' ./cmd/fid` | `ok  github.com/mbell-nzl/fid/cmd/fid 0.189s`; terminated `exit_code: 0`. |
| `go test -tags integration -count=1 -run '^(TestP24All49KeyedLeavesRebuildInOneJournal|TestP24ExecutedKeyedCorpusMatches49LeafInventory|TestP24ReviewOptionalPresenceAndEmptyException|TestP24ObservingJudgmentRetriesAfterSourceAndHEADChange|TestP24ObservingJudgmentReplayNeverObserves|TestP23StopFailureCodeMatrix/attention/invalid-argument|TestP24StopRaiseConflictNamesChangedDetails)$' ./cmd/fid/integration/` | `ok  github.com/mbell-nzl/fid/cmd/fid/integration 14.777s`; terminated `exit_code: 0`. Both retained repair targets pass. |
| `go test -tags integration -count=1 -run '^TestP24' ./cmd/fid/integration/` | `ok  github.com/mbell-nzl/fid/cmd/fid/integration 24.624s`; terminated `exit_code: 0`. Covers the full in-window Phase-24 integration corpus. |
| `go test -count=1 -run '^(TestCommitIdempotentReplaysStoredResultAfterRebuild|TestTask9PrechangeBudgetObservedJournalFailsClosedWithoutMutation)$' ./internal/journal` | `ok  github.com/mbell-nzl/fid/internal/journal 0.247s`; terminated `exit_code: 0`. |
| `go test -count=1 -run '^TestPhase22AttentionReadBoundary$' ./cmd/fid` | `ok  github.com/mbell-nzl/fid/cmd/fid 0.234s`; terminated `exit_code: 0`. The retained production-path inventory accepts the new descriptor file. |
| `go test -tags integration -count=1 -run '^(TestPortableRoundTripCoversEveryPhase23Table|TestP23StopFailureCodeMatrix|TestP22DualAdapterRecoveryConformance)$' ./cmd/fid/integration/` | `ok  github.com/mbell-nzl/fid/cmd/fid/integration 10.199s`; terminated `exit_code: 0`. Covers Task 1 portable parity and Task 3's retained failure/recovery journeys. |
| Exact 100k benchmark, first independent attempt | Ran concurrently with the preceding integration selection: `100k-event rebuild = 30.122364792s, want <30s`, then `FAIL`; terminated `exit_code: 1`. This failed attempt is retained. `08_QUALITY_PLAN.md` §7 says to run heavy gates serially. |
| Exact 100k benchmark, serial rerun after integration ended | `BenchmarkRebuildWorkspace100kEvents-6 ... 21531 ms/rebuild`, `PASS`, `ok  github.com/mbell-nzl/fid/internal/journal 24.138s`; terminated `exit_code: 0`. This is the representative independent VAL-005 sample. |
| Built CLI in fresh non-project `/tmp` directory | `PASS built CLI: 3 standalone templates, 18 conflict refusals, exact JSON --json precedence, no files created`; probe terminated `exit_code: 0`. It supports current behavior but does not replace F4's persistent task-owned corpus. |

VAL-052 wall-time comparison from the three reports:

| Task | cmd/fid | internal/workflow | Slowest other ordinary package | Tagged integration | Assessment |
|---|---:|---:|---:|---:|---|
| 1 | 215.296 s | 72.492 s | internal/issues 194.345 s (first full attempt; final cached) | 253.353 s | Under 480 s warning and 600/900 s limits. |
| 2 | 196.324 s | 67.367 s (uncached final tree) | internal/issues 113.986 s (uncached final tree) | 248.512 s | Under limits. |
| 3 | 200.072 s | 93.837 s (uncached) | internal/issues 166.391 s (uncached) | 302.158 s | Under limits. |

Inventory regeneration from `input-inventory.md` §1 and exact R1-R20 searches terminated `exit_code: 0`. The section-1 `wc -l` token was `235 total`: 92 leaves, 29 groups, 39 input leaves, 50 keyed leaves, zero without Short, 25 typed flags. R8 and R12 are file registries and were checked by `rg`/file existence, each terminating `exit_code: 0`.

| Registry | Current regeneration result | Delta from accepted planning inventory |
|---|---|---|
| R1 | 45 `RequestDigest:` rows in 33 files | 43/32 → 45/33: one new central `internal/journal/idempotency.go` reference and one extra direct Gate reference in `internal/evidence/execution.go`; Task 2 adoption retains the indexed sites. |
| R2 | One helper: `internal/journal/request_digest.go:14` | 11 family helpers → one shared helper. |
| R3 | 126 matches in 18 files | 16 files → 18: migration 0019 up/down; column registries and direct Gate effect are present. |
| R4 | 38 migration files, latest `0019_p24_idempotency_fields.up.sql`; 17 schema/projection references | Latest 0018/schema 18 → 0019/schema 19; projection meta format remains 5. |
| R5 | Seven original closed decoders | Unchanged; no new descriptor decoder matches the excluded pattern. |
| R6, R7 | 11 repository comparison/validator matches in eight files; 50 caller-repository matches in eight files | Still present for later observation tasks 4-8, as scheduled. |
| R8 | `invalid-argument` and `idempotency-conflict` definitions present; `repository-unavailable` and `operation-unavailable` absent | Conflict details gained `differing_fields`; the unavailable codes belong to later tasks. |
| R9, R10 | 41 Human failure functions; 32 Human `error:` assertions in 13 files | Unchanged pending task 10. |
| R11, R12 | 11 repository-view references; all three Stop/Attention renderer files exist | Unchanged pending task 11. |
| R13, R14 | 15 qualifier/`PRIVATE_` references; one live-leaf inventory test | Unchanged pending tasks 12-13. |
| R15, R16 | Zero no-Short leaves; 50 keyed leaves, 49 after excluding `export import` | R15 20 → 0 under task 3; R16 unchanged. |
| R17, R18 | 29 groups; 25 typed leaf flags | R17 unchanged; R18 22 → 25 because the three adopted leaves add `--example`. |
| R19, R20 | Seven `argument:"request"` sites; 25 generic-invalid-argument search rows | Unchanged; later tasks 9a/10 own convergence. |

## Per-Task Findings

### Task 1 - Idempotency field-digest storage

#### Must-fix

- **F1 [Must-fix; VAL-003]** `internal/journal/current_event_validation_test.go:14`: the four insert-map cases assert a replay error and line for invalid input, while no test mutates the `gate.receipt-recorded` update guard to differ from its stored map. The required `journal-replay-stopped` code/exit 6/line evidence for each invalid map is also absent; the report's real CLI refusal covers only an old journal with a missing map. Owning task: task-1. Repair boundary: task-1 journal/CLI test fixtures and report evidence, plus a causal in-scope fix only if a new assertion exposes a defect; preserve accepted Event and gate contracts. Required evidence: valid receipt update passes, mismatched guard fails closed without mutation, and each malformed insert map yields the specified code, exit and line; rerun affected named tests and task-1 exact gates, retain earlier attempts. Re-review scope: inspect new assertions, receipt effect/projection diff and rerun affected focused tests.
- **F2 [Must-fix; VAL-004]** `internal/export/export_test.go:104`: `assertArchiveCheck` checks only the error type and code. It does not assert that the schema-18 refusal's check names `manifest.json`, although the accepted criterion and EVID-04 require that path. Owning task: task-1. Repair boundary: strengthen the archive test/report within the existing schema-19 import fixture, without changing the archive contract. Required evidence: assert `archive-invalid`/`schema-unsupported` on `manifest.json`, preserve the unchanged-target digest check, then rerun the focused test and task-1 exact gates. Re-review scope: inspect the assertion and rerun the focused archive test.

#### Improvement

- None.

### Task 2 - Input-as-written digests and named conflict fields

#### Must-fix

- **F3 [Must-fix; VAL-066]** `cmd/fid/integration/p24_idempotency_test.go:367`: the seven-leaf test sends a successful observable new-key control and a refusing unmerged-path control, then deletes `base.txt` for source-bearing leaves and checks only same-key replay. It does not send a new-key control in the deleted-source state, so the second conditional state lacks its non-qualifying half. Owning task: task-2. Repair boundary: extend the existing seven-leaf integration fixture and report; do not change observation semantics assigned to task 4. Required evidence: for each project-file Source leaf, replay after deletion still has original bytes/no Event and a distinct viable new-key control in that same state refuses/no Event; retain the observable preflight. Rerun the focused test and task-2 exact gates. Re-review scope: inspect the source-leaf enumeration and paired control outcomes, then rerun the focused integration test.

#### Improvement

- None.

### Task 3 - Input descriptor framework, help, templates and field validation

#### Must-fix

- **F4 [Must-fix; VAL-015]** `cmd/fid/input_descriptor_test.go:220`: the six-row exclusivity test exercises only `attention resolve`; it checks exit and a message substring, with exact `argument` only for `--json`. It does not assert the other required `value`/`allowed` details, operation, or no read/write for each adopted leaf. The report's one real JSON probe covers only `--json` precedence. Owning task: task-3. Repair boundary: add a built-binary or equivalent process corpus for the three adopted leaves and update the report; preserve the accepted `--example` rule. Required evidence: six conflicts per leaf, ordered first-conflict details including `--json` JSON envelope, exit 2, no effects, and standalone success outside a project; rerun task-3 exact gates. Re-review scope: inspect row completeness and exact assertions, rerun the corpus.
- **F5 [Must-fix; VAL-016]** `cmd/fid/input_descriptor_test.go:288`: the per-adopted-leaf structural corpus checks `Path` and merely `len(Allowed)>0`; it does not prove each class's exact `allowed` and `value` details or an Event-free process refusal. The exact-value unit test at `internal/cli/input_validate_test.go:8` uses one synthetic descriptor, and the one real Attention bound probe covers only one class. Owning task: task-3. Repair boundary: strengthen the existing descriptor/CLI tests and report within the three adopted leaves; keep the accepted detail schema. Required evidence: qualifying well-typed input versus all applicable top-level and nested structural failure classes with exact path, value, allowed, exit 2 and no Event; rerun task-3 exact gates. Re-review scope: inspect class/leaf coverage and rerun the strengthened tests.

#### Improvement

- None.

Task 3 retained-test convergence is otherwise sound: `git blame` attributes the old `stop.raise` fixture line to `25caa7d4` (phase-22/task-9), the Attention matrix row to `180f5168` (phase-23/task-19), and the exact production-path inventory to `59d2653d` (phase-23/task-14). TRG-24-001 authorized only the two explicit nullable fixture members; diff `04b6b099..2ee33ba3` makes exactly that fixture change and retains the journey. The first exact tagged gate remains `FAIL` (`exit_code: 1`). TRG-24-002 then authorized the structural error adapter; diff `2ee33ba3..c942419` changes `InputError.Details()` from numeric to accepted string-or-null details, without editing `internal/cli/codes.go` or the retained matrix. The second exact tagged gate remains `FAIL` (`exit_code: 1`), and the later byte-identical gate passed at 302.158s (`exit_code: 0`). The ledger records MATERIAL-PROGRESS then RESOLVED, no exhausted no-progress budget. Focused retained reruns passed above.

## Repair Plan

| Finding | Owning task | Affected assertion | Exact repair boundary | Required evidence | Fresh re-review scope |
|---|---|---|---|---|---|
| F1 | task-1 | VAL-003 | Existing journal/CLI fixtures and report; causal task-1 fix only if exposed | Valid/invalid insert and receipt-guard states, exact CLI code/exit/line, no mutation; named tests and all seven exact gates | Diff, report, focused replay/receipt rerun |
| F2 | task-1 | VAL-004 | Existing archive test/report | Exact `manifest.json` check path plus rejection and unchanged target; focused test and all seven exact gates | Archive assertion and focused rerun |
| F3 | task-2 | VAL-066 | Existing seven-leaf integration fixture/report | Paired deleted-source replay and new-key refusal/no Event, observable preflight; focused test and all eight exact gates | Source-bearing row enumeration, state isolation, focused rerun |
| F4 | task-3 | VAL-015 | Adopted-leaf exclusivity tests/report | Three leaves × six combinations, exact details/exit/no effects and standalone success; all five exact gates | Row completeness and process rerun |
| F5 | task-3 | VAL-016 | Adopted-leaf structural tests/report | Exact class/path/value/allowed/no-Event outcomes and qualifying controls; all five exact gates | Class coverage and focused rerun |

The Phase Runner should dispatch owners in task-1, task-2, task-3 order, preserving failed attempts and later byte-identical passing gates. A fresh independent QA Reviewer must write the next review artifact, cite this review, check every F1-F5 evidence obligation, and rerun representative affected verification.

## Blocker Candidates

- None. All five findings are bounded evidence or causal test repairs inside the reviewed tasks. No accepted product, architecture, data-model, public-interface, dependency, security or gate definition needs changing.

## Existing Escalation References

- None. TRG-24-001 and TRG-24-002 are ordinary repair-chain references, not ESC identities; the current escalation register has no open Phase 24 entry.

## Final Lines

```text
Review written: project-docs/docs/phases/phase-24/qa/review-1.md
Must-fix: 5  Improvement: 0  Blocker candidates: 0
VERDICT: REPAIR
```
