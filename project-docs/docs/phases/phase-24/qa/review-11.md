# Phase 24 - Mid-Phase Review 11 (Tasks 12-14)

Reviewer: qa-reviewer
Date: 2026-10-11
Tasks covered: phase-24/task-12, phase-24/task-13, phase-24/task-14
QA window: E, after Task 14
Review type: INITIAL (fresh independent reviewer)
Prior review: none
Re-review of: none

## Verdict

REPAIR

The reviewed behavior passes the reproduced checks. Task 13 omits the
required R13 inventory-command receipt. This is bounded evidence repair
within Task 13; Window E must not advance until a fresh independent review
checks its repair. No protected decision or product repair is required.
Review 10 is Window D2 history, not a re-review target.

## Scope Reviewed

- `AGENTS.md`; current phase brief, active decisions, progress and escalations.
- Phase 24 task specs and completion reports for Tasks 12, 13 and 14;
  QA plan, especially sections 4, 5, 7 and 8; validation contract;
  input inventory; public-interface sections 4 and 8-10; relevant decisions
  DEC-274..278, DEC-280..282, DEC-287, DEC-289 and DEC-290.
- Mid-phase review template; quality plan; workflow sections 2.7, 8 and 9;
  interaction design sections 2-4, 6 and 7; impacted CLI documentation and
  glossary. The requested quality-plan section 5.2 does not exist in the
  current document; its binding sections 3-8 were reviewed. This locator
  issue is already carried under CF-033.
- Task 12 implementation `815a6b66af0381adee77077e826657511ac62292` and
  report finalization `fc4bdff7102b5bcac6a566addf4327d92478c11d`.
- Task 13 implementation checkpoint `595738e780967bb9e9a0a7049d939a4540d011a5`,
  DEC-289 continuation `558c9bc02a911a4c0402465c8585814cb37d9822` and
  report finalization `fd6783f9bb1ece018f1e561189e2e32a2754d449`.
- Task 14 original checkpoint `4b6a3837e52f24a5a8975a1d2abef94635c74823`,
  preserved report `98d45617b35c212683b4346d30a760eb8ed630d3`, DEC-290
  implementation `e97835b9f9983e357b45e80c3dbe4e4b02cedb6d` and report
  finalization/reviewed HEAD `86e0c26ab8c5f9e820259b740438ca90cf7aa93b`.

The runner-owned `run-ledger.md` was dirty on entry. The reviewer did not
edit or stage it. No held-out scenarios or evaluator material was read.
Only this review is a durable artifact of this session.

## Summary

- Review outcome: REPAIR
- Must-fix findings: 1
- Improvement findings: 0
- Blocker candidates: 0
- Existing escalation references: ESC-24-008 (resolved by DEC-290; context only)
- Covered validation assertions: VAL-042, VAL-043, VAL-044, VAL-045,
  VAL-046, VAL-047, VAL-048, VAL-049, VAL-050, VAL-051, VAL-052, VAL-063,
  VAL-068; preparation only for VAL-053, VAL-054, VAL-055, VAL-056, VAL-057

## Validation Contract Coverage

| VAL ID | Claimed by | Status | Evidence |
|---|---|---|---|
| VAL-042 | Task 12 | ON TRACK | `TestP24JudgmentGetRecords` reproduced: all three JSON records equal their list items, sole result member, operation, explicit `event:null`, project selection and unchanged journal events. |
| VAL-043 | Task 12 | ON TRACK | `TestP24JudgmentGetRejections` reproduced: unknown ID exit 3 with exact code/subject/details; missing ID, extra argument and unknown flag exit 2 in both modes, inside and outside a project, without journal changes. |
| VAL-044 | Task 13 | ON TRACK | Zero/one/three-item matrix reproduced across judgment leaves and resume sections; exactly one disclaimer before nonempty items, none on items or empty sections. Actual renderers retain `  None.`. |
| VAL-045 | Task 13 | ON TRACK | Every closed field/view matrix reproduced at 80/81 non-ASCII scalars, multiline and whitespace inputs. Lists/resume collapse and cut; detail/mutation/review-get views retain full text, including each Finding disposition. Retained Concern-list regression reproduced. |
| VAL-046 | Task 13 | ON TRACK | Matrix controls cases reproduced: ESC, CR, BEL and NUL escaped; no raw control bytes. |
| VAL-047 | Task 13 | ON TRACK | Real referenced-file sentinel, all retained Human fixture output/replay scans, three new get reads and resume reproduced. Sentinel regression proves its ignored source file does not change the bound repository. |
| VAL-048 | Task 13 | ON TRACK | `GLOSSARY.md:21` and `:32` distinguish referenced file content from every accepted recorded-reasoning field. Definitions match public-interface section 10. QA Lead retains final sign-off. |
| VAL-049 | Task 14 | ON TRACK | Human journey, original terminal-state diagnostic and exact-source-line regression reproduced. Verbose journey reaches `completed envelope ENV-d268e5`; final target and revision are recorded from Human resume. Independent refusal probe checks every journal file's byte digest. |
| VAL-050 | Task 14 | ON TRACK | `TestP24TakeControlRebuildExportImportParity` reproduced: completed VAL-049 workspace, live/rebuilt/imported portable rows, field maps, observed repositories, three detail reads, integrity and schema-19 archive/import checks. |
| VAL-051 | Task 14 | ON TRACK | Recursive registry-wide forbidden-member test and descriptor/type/live-tree completeness reproduced; real Stop corpus asserts no raw JSON; disclaimer absence search terminated with expected exit 1. |
| VAL-052 | Tasks 12-14 | ON TRACK | Reported full-gate times are within accepted budgets; warnings retained. Cached results remain explicitly nonnumeric rather than invented fresh measurements. See timing table. |
| VAL-063 | Tasks 13-14 | ON TRACK | Exact `PRIVATE_` absence search independently terminated exit 1, no output. |
| VAL-068 | Task 14 | ON TRACK | Loaded-registry edge predicate, cancelled edge, sorted categories, JSON `[]`, empty/nil Human `none.`, line placement and both exact help usages reproduced. Binary pairing covers revision 4/several edges and revision 5/no edges. |
| VAL-053 | Task 14 preparation | ON TRACK; Human outcome DEFERRED | Fresh 59-row manifest contains `take-control-start` and the real protected refusal. Human completion judgment remains phase-end work. |
| VAL-054 | Task 14 preparation | ON TRACK; phase-end Human outcome DEFERRED | All five denial rows retained and rehearsed; D2 judgments are historical, not phase-end acceptance. |
| VAL-055 | Task 14 preparation | ON TRACK; phase-end Human outcome DEFERRED | All clean/dirty, same/different-worktree views retained, including DEC-287's exact plain-copy run-inspect refusal. |
| VAL-056 | Task 14 preparation | ON TRACK; Human outcome DEFERRED | Judgment mutation/list pairs, three get views and resume appear in fresh manifest; Human readability remains with Human. |
| VAL-057 | Task 14 preparation | ON TRACK; Human outcome DEFERRED | Three authoring workspaces rehearsed from pre-action state: run close has expected_repository; Decision Reference omits digest; concern record is the third leaf. No Human first/second-try acceptance is claimed. |

The separate no-VAL command-coverage finding F1 prevents review acceptance
even though the reproduced product assertions track toward their criteria.

## Gate And Evidence Checks

All reviewer commands ran from `/home/mark/projects/fid` unless the probe
explicitly used a fixture workspace. Long-running nested command sessions
were polled until an observed `exit_code` was returned; their complete result
objects were preserved in the execution transcript.

### Completion reports and command applicability

Task 12 records all five exact task-spec commands. Tasks 13 and 14 record
all seven, including both predeclared exit-1 absence commands with expected
and observed statuses. Final receipts have output tokens and termination
labels. Earlier failed attempts, superseded evidence and causal repair notes
remain distinguishable. No narrower test selection substitutes for a listed
full gate. Task 12 also records its owned R14 command. Task 13 does not
record its owned R13 command (F1).

QA-plan section 5 assigns the shared five gates to every task, the absence
searches to their owners, and inventory regeneration to each review. Its
benchmark, reference-instance, archive-scale, full-checkptr, complete phase
gate and task-14-close fixture-retention command are explicitly phase-end
QA work. Those are not missing Developer gates. This review ran focused
tests and a fresh scratch fixture generator; it did not execute the full
phase-level gate or claim its results. The inventory header separately
requires a Developer receipt for every owned registry; the seven-command
Task 13 block does not waive R13.

| Check | Result and observed token |
|---|---|
| `project-docs/tools/report-lint.sh project-docs/docs/phases/phase-24/tasks/task-12-done.md project-docs/docs/phases/phase-24/tasks/task-12.md` | PASS: `PASS report-lint completed with no failures` (terminated `exit_code: 0`). |
| Same exact two-argument report-lint form for Task 13 | PASS: `PASS report-lint completed with no failures` (terminated `exit_code: 0`). Mechanical task-spec checks do not cover F1's inventory requirement. |
| Same exact two-argument report-lint form for Task 14 | PASS: `PASS report-lint completed with no failures` (terminated `exit_code: 0`). |
| Commit identity, subjects and provenance | `git rev-parse --verify <hash>^{commit}` resolves all 17 distinct commit references in the reports; script token `PASS report commit references: 17 distinct commits resolve; numeric scratch-directory suffixes are not commit claims` (terminated `exit_code: 0`). Git subjects match each report; Refs footers match owning tasks. All three report-finalization commits change only the report hash/subject placeholders. |
| `gofmt -l cmd internal` | PASS: `exit 0 (no output)` (terminated `exit_code: 0`); exact required command independently reproduced. |
| `rg -n "This (Concern\|Deviation) is recorded judgment\|This Outcome is authorized Human judgment" cmd/fid -g '!*_test.go'` | PASS: `exit 1 (no output)` (expected `exit_code: 1`; terminated `exit_code: 1`). |
| `rg -n "PRIVATE_" cmd/fid/integration` | PASS: `exit 1 (no output)` (expected `exit_code: 1`; terminated `exit_code: 1`). |

The initial broad hash-extraction diagnostic also tried three numeric
scratch-directory suffixes, each returning exit 128 (`fatal: Needed a
single revision`). They were not reported commit identities. The corrected
commit-reference audit above excludes these path components; no claimed
commit failed verification.

### Reproduced named tests

Each selection used `-count=1`, so these are fresh executions. The exact
commands are retained below, followed by their own observed package tokens.

```text
go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 -run '^(TestRegistryNextStates|TestP24ResumeIssuePosition|TestP24ResumeIssueNextStates|TestP24IssueTransitionPositionHelp|TestP24ConcernListSummary|TestP24JudgmentSectionDisclaimers|TestP24RecordedReasoningEveryFieldAndView|TestPhase23HumanLeafTextMatchesRecordedFacts|TestP24DescriptorForbiddenMembers|TestLiveLeafCommandsMatchContractInventories)$' ./internal/workflow ./cmd/fid
```

PASS: `ok github.com/mbell-nzl/fid/internal/workflow 3.342s` and
`ok github.com/mbell-nzl/fid/cmd/fid 2.660s` (terminated `exit_code: 0`).

```text
go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 -run '^(TestP24JudgmentGetRecords|TestP24JudgmentGetRejections|TestP23ConcernListHumanProjection|TestP24JudgmentDetailAndResumeSentinel|TestP24SentinelPreservesRepositoryBinding|TestP24HumanFixtures|TestP24StopHumanRealCorpus|TestP24HumanTakeControlJourney|TestP24TakeControlTerminalStateProvenance|TestP24TakeControlResumePair|TestP24JourneyProvenanceKeepsExactSourceLine|TestP24TakeControlRebuildExportImportParity|TestP23HumanFixtures|TestP23HumanFixtureStepsRemainRunnable|TestP23SharedLeafConformance)$' ./cmd/fid/integration
```

PASS: `ok github.com/mbell-nzl/fid/cmd/fid/integration 53.585s`
(terminated `exit_code: 0`). `FID_P24_HUMAN_FIXTURE_ROOT` was unset for
this selection, exercising the ordinary temporary-generator branch.

```text
go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 -run '^(TestInputDescriptorRegistryMatchesLiveTreeAndTypes|TestGrammarWholeTreeClasses|TestGrammarWholeTreeFlagDomains|TestEveryLiveLeafHasOneLineShort)$' ./cmd/fid
go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 -run '^(TestP24KeyedLeafCorpusMatchesLiveTree|TestP24EveryLiveKeyedLeafHasExecutableCorpusRow)$' ./cmd/fid
go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 -run '^TestP24GenericFallbackCorpusCompleteness$' ./cmd/fid/integration
```

PASS respectively: `ok github.com/mbell-nzl/fid/cmd/fid 47.922s`,
`ok github.com/mbell-nzl/fid/cmd/fid 2.466s`, and
`ok github.com/mbell-nzl/fid/cmd/fid/integration 2.594s`
(each terminated `exit_code: 0`). These verify live-tree, type, typed-flag,
keyed-leaf and fallback completeness beyond the in-window behavior tests.

```text
FID_P24_HUMAN_FIXTURE_ROOT=/tmp/fid-p24-qa11-92i2l_en/human go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -run '^TestP24HumanFixtures$' -count=1 -v ./cmd/fid/integration
go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 -run '^TestP24HumanTakeControlJourney$' -v ./cmd/fid/integration
```

PASS respectively: `ok github.com/mbell-nzl/fid/cmd/fid/integration 28.781s`
and `ok github.com/mbell-nzl/fid/cmd/fid/integration 7.532s`
(each terminated `exit_code: 0`). The second command was captured by a
Python subprocess wrapper that also returned observed exit 0; its complete
transcript is scratch `/tmp/fid-p24-qa11-92i2l_en/journey.log`.

Fresh manifest:
`/tmp/fid-p24-qa11-92i2l_en/human/fixtures-698117270/manifest.json`.
Its binary hash, commands and 59 fixture rows were checked by the generator
and an independent probe. Probe tokens include `PASS fixture preparation:
59 rows`, `PASS authoring inputs: run.close expectation present; Decision
digest omitted; concern.record third leaf`, and `PASS real help: both exact
DEC-290 usages; child exit_code=0` (probe terminated `exit_code: 0`).

### Timing and budget evidence

| Task | cmd/fid | internal/workflow | Slowest other reported ordinary package | Tagged integration |
|---|---:|---|---|---:|
| 12 final gate | 376.609 s | Cached; first exact attempt 142.144 s | internal/export 122.143 s; first attempt 183.669 s | 381.138 s |
| 13 final continuation | 340.352 s | Cached; original exact attempt 136.519 s, explicitly historical | internal/export 108.394 s; original internal/issues 238.583 s | 342.605 s |
| 14 final tree | 488.166 s | 120.646 s | internal/issues 228.100 s | 431.049 s |

Every numeric ordinary result is below 600 s; tagged results are below
900 s. Task 14's cmd/fid is 81.36% of budget, 8.166 s above the 480 s
warning, with 111.834 s headroom. This is the contract's reported warning,
not a failed budget. Task 12's earlier 484.599 s, Task 13's historical
530.445 s and Task 14's superseded 530.567 s warnings remain recorded.
The cached rows provide no new wall-time measurement; the reports state
that limitation and do not invent durations. Narrow reviewer selections
do not remeasure full-suite budgets.

## Per-Task Findings

### Task 12 - Concern, Deviation and Human Outcome get leaves

#### Must-fix

- None.

#### Improvement

- None.

Observation (interaction design sections 2-4 and 6): the adapters delegate
record loading to `internal/review`, reuse the list loaders, and preserve
typed not-found handling. The real-process tests assert record equality,
transport metadata, grammar before project access and absence of journal
changes. The three additions enter the live inventory without changing
capability grants. DEC-274 explicitly requires the Phase 23 producer-table
sync; that documentation change is authorized and matches the behavior.

### Task 13 - Judgment disclaimers once per section; recorded reasoning shown

#### Must-fix

- **F1 — no VAL, required evidence:**
  [task-13-done.md:179](../tasks/task-13-done.md#current-final-required-gate-dec-289-continuation)
  records the seven task-spec gates, but neither its Tests run nor its
  quality-gate results records R13's `rg -n "p23CheckJudgmentQualifiers|PRIVATE_" cmd`.
  [input-inventory.md:7](../input-inventory.md) requires each Developer to
  rerun and record owned registry commands before close; its R13 row at
  line 94 assigns this command to Task 13. The narrower `PRIVATE_` absence
  scan and passing judgment tests do not replace the required registry
  command. The reviewer reproduced three remaining qualifier references,
  exit 0, but reviewer evidence cannot be transcribed into a Developer PASS.
  Owning task: **phase-24/task-13**. Repair boundary: execute the existing
  R13 command and append its observed receipt to Task 13's completion
  report; preserve historical results and change no product, spec, command
  or assertion. Required evidence: same-session command, output token,
  observed exit 0, working directory, accounted three-reference result,
  committed clean report and two-argument report-lint. Fresh re-review
  scope: F1, report delta/history, R13 regeneration and representative
  judgment/sentinel verification within Window E.

#### Improvement

- None.

Observation (interaction design sections 2-4 and 7): one shared disclaimer
function serves leaf and resume sections, and the rendering distinguishes
list from detail views over the accepted closed field list. The glossary
and CLI examples agree. Sentinel fixtures use already ignored project-file
content and preserve repository identity; they prove file-content
non-disclosure while recorded reasoning remains visible. DEC-289's
continuation changes only the authorized Concern-list predicate; the rest
of that file is unchanged. Earlier assertion provenance resolves to the
reported introducing commits. F1 is a report-coverage omission, not a
judgment-rendering failure.

### Task 14 - Human-takes-control journey and phase conformance

#### Must-fix

- None.

#### Improvement

- None.

Observation (interaction design sections 2-4 and 6-7): the registry owns
edge enumeration; the Human adapter only formats the resulting ordered
states. The product diff is exactly DEC-290's member, two Human lines and
two transition usages. It uses `RulesForEdge`, including cancellation,
rather than the deliberately filtered `Definition.Transitions` view.
The help strings match DEC-290 verbatim; issue-update help is unchanged.
Empty judgment sections use `None.`; empty/nil Next states use `none.`.

The actual journey rediscovers the Issue, Envelope, revision, Decision
class/capability and Obligation from earlier Human output. It decodes each
leaf's `--example`, records substitutions under the closed five provenance
classes, permits setup only for Actors/project directory, and writes input
files outside the worktree. Its read guard retains the JSON/SQL/journal
prohibition. Setup JSON and post-journey parity inspection do not supply
acting-caller values. The final source receipts are:

```text
transition --to provenance class=fid-output source="  Next states: cancelled (cancelled), qa_block (started), qa_pass (completed)"; --expected-revision provenance class=fid-output source="  Revision: 4"
envelope.complete success: completed envelope ENV-d268e5
```

The final resume follows Evidence linking and precedes terminal transition;
the source recorder retains the exact observed source line even when the
same token occurs elsewhere. The approved Decision Reference/Outcome pair
resolves the Stops, qualifying Evidence permits transition, and the exact
completion confirmation is asserted. Independent replay of the retained
protected refusal exits 5 with empty stdout and unchanged SHA-256 digests
for every journal file, beyond the test's event-count assertion.

Parity enumerates every SQLite table, then excludes only projection metadata
and migration bookkeeping. Sorted complete row comparisons therefore
include observed repository data. Separate field-map comparisons reject
empty maps; all three judgment get results match live/rebuilt/imported
states. Export validation and both manifest/import schema-19 assertions
execute. Task 14's preserved blocked attempt and superseded run are not
presented as the final passing evidence.

## Prior Finding Resolution

Not applicable. This is the initial Window E review by a fresh reviewer,
with `Re-review of: none`. Review 10's Window D2 findings are not reopened.

## Repair Plan

| Finding | Owning task | Affected assertion | Exact repair boundary | Required evidence | Fresh re-review scope |
|---|---|---|---|---|---|
| F1 | phase-24/task-13 | no VAL; required owned-registry evidence | Run `rg -n "p23CheckJudgmentQualifiers\|PRIVATE_" cmd` from repository root and append the current-session receipt to task-13-done.md. Report-only durable change; preserve earlier failures, history and valid gate receipts. Do not change code, specs, accepted commands or validation semantics. | Command's own token and observed exit 0; accounted qualifier matches and no PRIVATE_ references; finalized committed clean report; `project-docs/tools/report-lint.sh project-docs/docs/phases/phase-24/tasks/task-13-done.md project-docs/docs/phases/phase-24/tasks/task-13.md` on that committed report. | Fresh independent Window E reviewer cites review-11, verifies F1 against the new receipt/report delta, independently regenerates R13 and reruns representative affected judgment/sentinel evidence. Product evidence here may be traced, not relabelled as fresh execution. |

## Cross-cutting observations

### R1-R20 completeness regeneration

All registries were regenerated in this session. Markdown-escaped table
alternations were executed as regex alternation. Query commands below
returned observed exit 0 except R5, whose expected absence returned exit 1.
The Python inventory runner captured each child stdout/stderr and exit;
its enclosing result terminated exit 0. R8/R12 are file registries, so
their listed artifacts were inspected through concrete queries/counts.

| Registry | Command/evidence and observed token | Current result / termination |
|---|---|---|
| R1 | `rg -n "RequestDigest:" internal cmd -g '!*_test.go'`; token `RequestFieldDigests: fields}` | 45 sites, 33 files; `exit_code: 0`. |
| R2 | `rg -n "func [a-zA-Z]*(RequestDigest\|requestDigest\|judgmentDigest)\(" internal -g '!*_test.go'`; token `func ComputeRequestDigest(` | One shared helper; `exit_code: 0`. |
| R3 | `rg -n "idempotency_results" migrations internal -g '!*_test.go'`; token `ALTER TABLE idempotency_results DROP COLUMN request_field_digests_json;` | 126 matches in 18 files; `exit_code: 0`. |
| R4 | `ls migrations`; token `0019_p24_idempotency_fields.up.sql`; `rg -n "ProjectionSchemaVersion\|projection_schema" internal -g '!*_test.go'`; token `CurrentProjectionSchema:` | 38 up/down files, 17 schema references; each `exit_code: 0`. Current archive/import schema 19 reproduced separately. |
| R5 | `rg -n "func decodeClosed[A-Za-z]*Input" cmd/fid`; token `exit 1 (no output)` | Zero old decoders; expected/observed `exit_code: 1`. |
| R6 | Inventory repository-helper regex; token `func sameRepository(left, right Repository) bool {` | Four matches in three files; `exit_code: 0`. |
| R7 | Inventory caller-repository-writer regex; token `repositoryJSON := mustEnvelopeJSON(request.RepositoryScope)` | 16 matches in two files; `exit_code: 0`; remaining references are resolved internal facts, not public input revival. |
| R8 | `rg -n "definition\(" internal/cli/codes.go`; token `CodeImportTargetConflict: definition(`; registered definitions inspected | 85 declaration/definition lines; `exit_code: 0`; unavailable codes and three get not-found producers accounted. |
| R9 | Inventory failure-renderer regex; token `func renderActivationError(` | 43 functions in 31 files; `exit_code: 0`; shared failure mapping retained. |
| R10 | `rg -n '"error: ' cmd internal -g '*_test.go'`; token `out.String() != want {` | 63 matches in 27 files; `exit_code: 0`; expanded Human/grammar assertions. |
| R11 | `rg -n -i "worktree\|dirty digest\|HeadOID\|head_oid" cmd/fid -g '!*_test.go'`; token `return "different worktree", nil` | 21 matches in 11 files; `exit_code: 0`; renderer/help/fixture references inspected. |
| R12 | `wc -l cmd/fid/stop_human.go cmd/fid/resume_human.go cmd/fid/resume_human_guidance.go`; token `440 total` | All three listed renderer files exist and were inspected; `exit_code: 0`. |
| R13 | `rg -n "p23CheckJudgmentQualifiers\|PRIVATE_" cmd`; token `p23CheckJudgmentQualifiers(t, pair.Operation, human)` | Three references in two test files, no PRIVATE_ values; `exit_code: 0`. Reviewer PASS does not close F1. |
| R14 | `rg -n "func TestLiveLeafCommandsMatchContractInventories" cmd/fid`; token `func TestLiveLeafCommandsMatchContractInventories(t *testing.T) {` | One completeness test; `exit_code: 0`; named test reproduced. |
| R15-R18 | Verbatim inventory section-1 Bash block, executed with `bash /tmp/fid-p24-qa11-92i2l_en/inventory.sh`; terminal token `274 total` | 95 leaves, 29 groups, 39 input leaves, 50 flag-keyed leaves, zero no-Short leaves, 61 typed flags; observed `exit_code: 0`. |
| R19 | `rg -n '"argument": *"request"' cmd/fid -g '!*_test.go'`; token `"accepted issue read input"` | Exactly four family fallbacks in issue_get_list, issue_mutation_response, link_response and receipt; `exit_code: 0`. |
| R20 | `rg -n 'CodeInvalidArgument, "invalid [^"]*"' cmd/fid -g '!*_test.go'`; token `"invalid gate request"` | 23 accounted sites: 18 generic fallbacks, four typed failures, one centralized canonical-request construction; `exit_code: 0`. AST corpus completeness also reproduced. |

The exact inventory block generated files under `/tmp/tmp.IEmQr6RsPY`.
An independent set comparison printed `PASS registry list sets: 95 leaves
(93 contract + 2 diagnostics), 39 input, 49 canonical keyed in 20 families,
29 groups, 61 typed flags (17 int + 44 bool), zero missing Short`
(terminated `exit_code: 0`). The 50 flag-keyed set minus export import
equals the 49-row corpus exactly.

Planning counts are historical. The three DEC-274 get leaves explain
92 → 95 total leaves; the grammar corpus excludes only the two retained
diagnostics, leaving 93. The 39 new example booleans explain 22 → 61 typed
flags. Shared digest-helper convergence, migration 0019, descriptor adoption
and centralized grammar/failure handling explain the older registry deltas;
the current counts agree with the prior D1 baseline except the in-window
get additions and judgment/fixture assertions. R13's former reasoning
markers are removed while its three qualifier references remain. No
unaccounted fallback or missing live-tree/corpus row was found.

### Wrap-up notes

- Carry Task 14's 488.166 s warning to wrap-up; its budget passes and no
  speed repair is required here. Preserve the existing CF-043 performance
  planning boundary.
- The empty-judgment matrix currently checks case-insensitive `none`
  rather than the exact `  None.` marker. Actual output is correct and
  DEC-289 is satisfied; tightening that test on a future causal edit would
  improve regression specificity without changing the accepted wording.
- Human VAL-053..057 outcomes, complete phase gates and final QA Lead
  acceptance remain pending at their assigned phase-end boundaries.
- The reviewed diff adds no dependency or architecture boundary. The
  new registry method remains in the domain layer, and updated production
  Human/registry files stay within the 400-line target.

## Blocker Candidates

- None.

## Existing Escalation References

- ESC-24-008: resolved by DEC-290. Its protected gap is preserved in Task
  14 history and its accepted resolution is verified here. It does not
  cover F1 and is not reopened.

## Final Lines

```text
Review written: project-docs/docs/phases/phase-24/qa/review-11.md
Must-fix: 1  Improvement: 0  Blocker candidates: 0
VERDICT: REPAIR
```
