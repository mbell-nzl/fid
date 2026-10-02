# Phase 23 - Mid-Phase Review 13 (Tasks 16-18)

Reviewer: qa-reviewer
Date: 2026-10-01
Tasks covered: phase-23/task-16, phase-23/task-17, phase-23/task-18
Review type: INITIAL
Prior review: none (review-12 covers the preceding window)

## Verdict

REPAIR. Two bounded evidence defects remain in tasks 16 and 18. The window does not advance. The passing Developer gates and this review's passing focused tests do not close those defects.

## Scope Reviewed

- Task 16-18 specs and completion reports; Phase 23 validation contract, QA plan sections 4, 5, 7, 8, 10, 11 and readiness amendment 22; accepted contracts, public interface, DEC-265/266, both task-18 resume handoffs, review-12, and the quality plan.
- Implementation commits `d02614e6c26757076a91d8a12e08dfa328fe699d`, `994a8734d15580a23a2a79d3c28f5f33d9dbf1e5`, `5fe02e15c508a9b18721dbcff0e463feb53147c7`; task-18 report commit `0c4fc6510f0c09dcb78a1e18c777968614a1ed67`. The reported hashes resolve and their subjects and `Refs:` trailers match.
- Independent commands ran from `/home/mark/projects/fid`. Runner-owned ledger, triage log, and escalation register were already modified and were not edited.

## Summary

- Review outcome: REPAIR
- Must-fix findings: 2
- Improvement findings: 0
- Blocker candidates: 0
- Existing escalation references: ESC-23-010 (resolved by DEC-265/266), ESC-23-011 (restored-capability recovery recorded in the ledger)
- Covered validation assertions: VAL-004, VAL-027, VAL-047, VAL-064, VAL-065, VAL-066, VAL-067, VAL-068, VAL-069, VAL-070, VAL-071, VAL-072, VAL-073, VAL-076, VAL-079

## Validation Contract Coverage

| VAL ID | Claimed by | Status | Evidence |
|---|---|---|---|
| VAL-004, VAL-027, VAL-064, VAL-065 | 16 | ON TRACK | The independently rerun `TestPortableRoundTripCoversEveryPhase23Table` passed. Its four fixtures populate the 14 named Phase 23 tables and nine Events, compare all application rows and archive entries across live/rebuilt/imported projections, verify non-mutating validate, reject schema 17 while importing schema 18, check rebuilt integrity, and retain one Review-keyed charge after import. The earlier migration-17 pragma leg belongs to tasks 1/15 and final QA. |
| VAL-076 | 16 | NOT ON TRACK | `p23PortableAssertCorruptEventsRejected` exercises nine types × three corruptions plus one unknown type and checks the `current Event` error, but all 28 probes use one workspace and connection. The contract requires a fresh workspace for each probe; see MF-1. |
| VAL-047, VAL-066 | 17 | ON TRACK | The named matrix and Attention tests passed independently. The pinned twelve-row fixture compares discovery and export manifest exactly, including two unavailable rows; resume checks its required capability. Task 11's report records the earlier `human_outcome_resolution:false` binary and working resolution path; task 17's current test checks `true`. |
| VAL-067 | 18 | ON TRACK | The independently rerun grant test passed. `grant-actions.json` has 30 independent members; the test compares both closed registries directly and against that corpus, activates the complete set, and rejects unknown/operation-only names without an Event. Fresh `git log -S` checks confirmed all seven Phase 23 names' introducing commits and `Refs:` owners listed in the report. |
| VAL-068 | 18 | NOT ON TRACK | The shared 23-operation corpus and gate regressions passed independently, and DEC-266's authority shape is implemented. The comparator can nevertheless erase a mismatch in a caller-provided `review_id`; see MF-2. |
| VAL-069 | 18 | ON TRACK | The named product-surface test passed independently, comparing 59 latest-schema tables and 125 built command paths with its pinned inventory and checking aliases. |
| VAL-070, VAL-071, VAL-072, VAL-073 | 18 fixture refresh; Human at phase end | DEFERRED | The independently rerun Human fixture and step tests passed. The regenerated manifest and recorded interactive steps are runnable; comprehension and presentation outcomes remain Human judgments at phase end. |
| VAL-079 | 16-18; QA Lead at phase end | ON TRACK | Reports give ordinary `cmd/fid` 246.504/236.172/196.827 s, `internal/workflow` 100.101/90.207/57.728 s, and tagged 269.692/234.507/249.568 s, respectively. Task 18's uncached workflow and `cmd/fid` measurements supplement its cached ordinary gate. All are below the 600/600/900 s budgets. Phase-end budget verification remains assigned to QA Lead. |

## Gate And Evidence Checks

Each report's `## Tests run` and final `## Exact quality-gate results` was compared with its spec and QA plan. Tasks 17 and 18 list all five exact commands; task 16 lists those five plus its required fresh-process 100k rebuild benchmark. Each final result has a command output token and observed `exit_code: 0`. Task 18 preserves the attempt-5/6 focused and tagged `FAIL`/`BLOCKED` records with `exit_code: 1`, the DEC-266 causal repair, and the later byte-identical passing commands. Report-attempt-7's first report-lint failure and passing retry are also preserved. No task in this window has an expected non-zero command. The QA plan assigns the full seven-command, archive-scale, three-sample benchmark, and full-checkptr phase gates to final QA; running them here would duplicate that later gate without resolving these two evidence defects.

| Check | Command or evidence | Result | Reproduced evidence |
|---|---|---|---|
| Task 16 report control | `project-docs/tools/report-lint.sh project-docs/docs/phases/phase-23/tasks/task-16-done.md project-docs/docs/phases/phase-23/tasks/task-16.md` | PASS | `PASS report-lint completed with no failures` (terminated `exit_code: 0`). |
| Task 17 report control | Same two-argument command for task 17 | PASS | `PASS report-lint completed with no failures` (terminated `exit_code: 0`). |
| Task 18 report control | Same two-argument command for task 18 | PASS | `PASS report-lint completed with no failures` (terminated `exit_code: 0`); historical failed control retained in report. |
| Required format command | `gofmt -l cmd internal` | PASS | `exit 0 (no output)` (terminated `exit_code: 0`). |
| Independent in-window named tests | `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 -run '^(TestPortableRoundTripCoversEveryPhase23Table|TestFinalCapabilityMatrixMatchesAcceptedTable|TestP23AttentionResolutionThroughReviewRecordedAndFindingResolvedRealProcess|TestGrantActionSetIsExactlyThirtyNames|TestP23ConditionWireGrammarMatchesJournalAndFixtures|TestP23SharedLeafConformance|TestP23SharedAmendedLeafConformance|TestP23ConcernDispositionReferenceRejections|TestP23FailureCorpusHasEveryAcceptedCode|TestP23CancelledFamiliesAbsentFromProduct|TestP23HumanFixtures|TestP23HumanFixtureStepsRemainRunnable)$' ./cmd/fid/integration/` | PASS | `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	20.524s` (terminated `exit_code: 0`). This is fresh focused QA evidence, not a substitute for the Developer's full tagged gate. |
| Independent DEC-266 regression and CLI matrix tests | `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 -run '^(TestGateRegisterUnknownActorReturnsContractAuthorityFailure|TestGateRegisterDeniedActorReturnsContractAuthorityFailure|TestCapabilitiesJSONGoldenAndProjectionRows)$' ./cmd/fid` | PASS | `ok  	github.com/mbell-nzl/fid/cmd/fid	3.780s` (terminated `exit_code: 0`). |

## Per-Task Findings

### Task 16 - Portable recovery for the review and judgment domain

#### Must-fix

- **MF-1** [`cmd/fid/integration/p23_portable_recovery_test.go:104`]: `p23PortableAssertCorruptEventsRejected` opens the projection once and reuses it for every malformed/unknown Event. VAL-076 explicitly requires each payload to be attempted against a fresh workspace with before/after projection evidence. Passing errors and row equality on one shared connection leave that isolation requirement unproven. Owning task: task-16. Repair boundary: this task's integration corruption fixture/test only; retain all 28 cases, their validator-error assertion, and full projection comparison. Required evidence: each case uses a separately initialized or cloned clean workspace and proves the complete projection unchanged; rerun the named test and all six unchanged task-16 gates with observed termination, then update the report. Re-review scope: inspect case isolation and before/after evidence, rerun the affected named test and representative gate, and recheck report-lint. Affected assertion: VAL-076.

#### Improvement

- None.

Observation: `p23PortableRoundTrip` compares all application tables and canonical archive entries, while `p23PortableAssertPriorSchemaRejected` supplies both accepted schema-18 and rejected schema-17 states. `ProjectEvent` calls `validateCurrentEvent` before applying Effects, so the malformed probes do reach the intended validator.

### Task 17 - Final capability matrix and availability audit

#### Must-fix

- None.

#### Improvement

- None.

Observation: commit `994a8734` changes only the three accepted availability rows in workflow/export, converges the retained tests, and updates `02_TRD.md`. The independent matrix test compares the twelve exact discovery and manifest rows, and the real-process Attention test exercises the resolution path after the final `true` limit.

### Task 18 - Conformance assembly and wire-surface audit

#### Must-fix

- **MF-2** [`cmd/fid/integration/p23_conformance_test.go:829`]: the shared normalizer replaces any syntactically valid `review_id` with `<generated:REV>` based only on the field name. `p23EvidenceLinkModel` copies `review_id` from the caller's request. Two different valid Review IDs in that returned field therefore compare equal after normalization, contrary to the handoff's exact equality for caller-provided values and its rule that normalization touch only fid-generated values. The same generic traversal also handles `id` and `envelope_id` without checking their origin. Owning task: task-18. Repair boundary: the shared conformance comparator/oracle tests in `cmd/fid/integration/`; retain all 23 operations, the independent fake, rejection classes, and DEC-266 authority expectations. Required evidence: restrict placeholder conversion to declared generated positions while applying the rule identically to both adapters; add a regression showing distinct well-formed caller-provided IDs stay distinct and compare the Evidence `review_id` with its input; rerun the named shared corpus, five unchanged task-18 gates, and report-lint with observed termination. Re-review scope: inspect normalization and fake independence, rerun the new regression plus the affected shared corpus and representative gate, and recheck the report. Affected assertion: VAL-068.

#### Improvement

- None.

Observation: the DEC-266 amendment's non-test producer inventory has ten uses: eight changed Phase 20-shaped producer lines and the unchanged `gate.go:325/327` pair. Registry definitions for all four `authority-*` codes have the Phase 21 four fields and null subject; the gate tests exercise both latent paths with exit 5 and no Event. The two corrected fake authority rows are literal contract expectations; the models build expected values from contract rules and fixture requests, not fid response maps. The `failures.json` markers and test enforce DEC-265's 19 declared = 18 added + sole retained `requirement-not-found`. The report's VAL-068 matrix cites public-interface §§1, 4, 5.1 or `contracts.md` for each N/A rejection class, and its eight changed-producer rows identify the introducing commits and `Refs:` trailers. The historical attempt-6 failure remains separate from the report-attempt-7 pass. The ledger records restored capability, task-18 REPAIR `RESOLVED`, and task COMPLETE at report commit `0c4fc651`; this review has not edited it.

## Repair Plan

| Finding | Owning task | Affected assertion | Exact repair boundary | Required evidence | Fresh re-review scope |
|---|---|---|---|---|---|
| MF-1 | task-16 | VAL-076 | `cmd/fid/integration/p23_portable_recovery_test.go` corruption fixture/test and task-16 report; no product or contract change. | Fresh workspace per nine × three corrupted Event cases and unknown type; full projection before/after; named test and six unchanged task gates with tokens and observed termination; committed report and report-lint PASS. | Inspect isolation and snapshots; rerun named corruption test and representative required gate; verify report evidence. |
| MF-2 | task-18 | VAL-068 | `cmd/fid/integration/` conformance comparator/oracle tests and task-18 report; no product, fake-derived-from-fid, contract, or command change. | Exact caller-value comparison despite generated-looking syntax; generated-only normalization regression; 23-operation shared corpus and five unchanged task gates with tokens and observed termination; committed report and report-lint PASS. | Inspect normalization paths and independent fake; rerun new regression, shared corpus, and representative required gate; verify report evidence. |

Repair order: task-16, then task-18. A fresh independent QA Reviewer must write the next review artifact for this same window and assess both findings; neither is resolved by this review.

## Blocker Candidates

- None. Both repairs have decision-free ownership within reviewed task files. The historical DEC-265/266 authority questions have accepted dispositions, and the ledger records restored Developer capability.

## Existing Escalation References

- ESC-23-010 is resolved by DEC-265 and DEC-266. ESC-23-011 remains marked open in the current register, while the Runner ledger records the restored-capability resume and task-18 completion. No new escalation ID is allocated here.

## Final Lines

```text
Review written: project-docs/docs/phases/phase-23/qa/review-13.md
Must-fix: 2  Improvement: 0  Blocker candidates: 0
VERDICT: REPAIR
```
