# Phase 23 - Mid-Phase Review 14 (Tasks 16-18)

Reviewer: qa-reviewer
Date: 2026-10-01
Tasks covered: phase-23/task-16, phase-23/task-17, phase-23/task-18
Review type: FRESH INDEPENDENT RE-REVIEW
Prior review: `project-docs/docs/phases/phase-23/qa/review-13.md`
Re-review of: review-13

## Verdict

PASS. Fresh source inspection and command runs resolve both review-13 Must-fix findings. The Human judgments and phase-end gate remain assigned to phase-end QA.

## Scope Reviewed

- Task 16-18 specs and completion reports, Phase 23 validation contract and QA plan, the quality plan, current phase context, and both task-18 handoffs (`resume-after-esc-23-010.md`, `resume-after-dec-266.md`).
- Repair commits `4261575c3c84de1acf0089d70cf50fbe310a16c8`, `c2f02906a18b9ba29ace7f4b7e5be3fb44cea581`, `3227197d5f190528825cca3a8c01a480c06712e2`, and `1d2e783a0fa87b97062acc4c7feef248a8436270`. Each resolved separately with `git rev-parse --verify '<hash>^{commit}'`; `git show --format=fuller --no-patch` confirmed the subjects and `Refs: phase-23/task-16` or `Refs: phase-23/task-18` trailers. The two later commits amend only their respective reports.
- Independent commands ran from `/home/mark/projects/fid`. Runner-owned edits in the ledger, triage log, and escalation register were present before this review and were left untouched.

## Summary

- Review outcome: PASS
- Must-fix findings: 0
- Improvement findings: 0
- Blocker candidates: 0
- Existing escalation references: ESC-23-010 (DEC-265/266 dispositions), ESC-23-011 (the earlier capability interruption)
- Covered validation assertions: VAL-004, VAL-027, VAL-047, VAL-064, VAL-065, VAL-066, VAL-067, VAL-068, VAL-069, VAL-070, VAL-071, VAL-072, VAL-073, VAL-076, VAL-079

## Validation Contract Coverage

| VAL ID | Claimed by | Status | Evidence |
|---|---|---|---|
| VAL-004, VAL-027, VAL-064, VAL-065 | task-16 | ON TRACK | The unchanged round-trip fixture and report cover the rebuild pragma leg, one Review-keyed charge after import, every named Phase 23 table and Event, non-mutating validation, accepted schema 18 and rejected schema 17. The earlier migration-17/18 pragma legs belong to tasks 1/15 and final QA. |
| VAL-076 | task-16 | ON TRACK | `p23_portable_recovery_test.go:129-220` takes all 59 projection tables before and after each probe. Fresh verbose output showed nine named Event types × three corruptions plus one unknown Event type: 28 distinct cloned workspace paths, `current Event` errors, and equal before/after SHA-256 pairs. The exact named test passed independently. |
| VAL-047, VAL-066 | task-17 | ON TRACK | The fresh matrix/Attention selection passed. The task-17 report records the earlier false capability state and final true state with the working resolution path; the pinned twelve-row test covers the two unavailable rows and four reporting surfaces. |
| VAL-067 | task-18 | ON TRACK | The shared corpus passed. `TestGrantActionSetIsExactlyThirtyNames` compares the independently pinned 30 members with both closed registries and the registries with each other, then exercises allow and reject cases. The seven Phase 23 introduction commits and owners remain in the report; the repair commits changed none of that inventory. |
| VAL-068 | task-18 | ON TRACK | The new Review ID regression and four-test shared corpus passed independently. `p23NormalizeLeafValues` now follows generated placeholders declared by the independent fake; unmarked caller values compare exactly. Full accepted and rejected cases retain the 23-operation matrix, including authority and stale/binding negatives. |
| VAL-069 | task-18 | ON TRACK | The fresh product-surface selection passed; the test inventories 59 latest-schema tables, 125 command paths, and aliases against pinned expectations. |
| VAL-070, VAL-071, VAL-072, VAL-073 | task-18 fixtures; Human at phase end | DEFERRED | The fresh Human fixture and step selection passed. The report preserves regenerated runnable inputs; comprehension and presentation are explicitly Human judgments at phase end. |
| VAL-079 | tasks 16-18; QA Lead at phase end | ON TRACK | Repaired task-16 times: ordinary `cmd/fid` 196.206 s, uncached `internal/workflow` 57.597 s, tagged integration 232.525 s. Repaired task-18 times: 197.441/56.673/296.429 s. Task 17 records 236.172/90.207/234.507 s. All are below 600/600/900 s; phase-end fresh gates remain with QA Lead. |

## Gate And Evidence Checks

Each completion report's `## Tests run` and final `## Exact quality-gate results` was compared with its task spec and QA-plan section 22. Tasks 17 and 18 retain the five exact commands; task 16 retains those five plus the required fresh-process 100k rebuild benchmark. Each final result has its command's output token and observed `exit_code: 0`; no narrower command replaces a required gate. Task 16 preserves its named-test red (`corruption probe reused a projection connection`, `FAIL`, `exit_code: 1`) and first report-lint failure (`exit_code: 2`). Task 18 preserves the attempt-6 tagged `claim.acquire` authority failure (`FAIL`, `exit_code: 1`), the new regression's red (`normalized Evidence review_id = <generated:REV>`), the first shared-corpus marker failure, and report-lint's earlier failure. Each has a later byte-identical passing command and a causal deviation note. The QA plan assigns its full phase gate, three fresh benchmark samples, archive-scale check, and full-checkptr control to final QA, so they were not rerun in this bounded review.

| Check | Command or evidence | Result | Reproduced evidence |
|---|---|---|---|
| Task 16 report control | `project-docs/tools/report-lint.sh project-docs/docs/phases/phase-23/tasks/task-16-done.md project-docs/docs/phases/phase-23/tasks/task-16.md` | PASS | `PASS report-lint completed with no failures` (terminated `exit_code: 0`); all six commands and reported commit resolved. |
| Task 17 report control | Same two-argument command for task 17 | PASS | `PASS report-lint completed with no failures` (terminated `exit_code: 0`); all five commands and reported commit resolved. |
| Task 18 report control | Same two-argument command for task 18 | PASS | `PASS report-lint completed with no failures` (terminated `exit_code: 0`); all five commands and reported commit resolved. |
| Task 16 MF-1 exact test | `go test -tags integration -run '^TestPortableAssertCorruptEventsRejected$' -count=1 ./cmd/fid/integration/` | PASS | `ok  \tgithub.com/mbell-nzl/fid/cmd/fid/integration\t7.742s` (terminated `exit_code: 0`). A `-v` rerun showed all 28 subtests and matching 59-table digests; `ok  \tgithub.com/mbell-nzl/fid/cmd/fid/integration\t7.751s` (terminated `exit_code: 0`). |
| Task 16 unchanged benchmark gate | `go test -run '^$' -bench '^BenchmarkRebuildWorkspace100kEvents$' -benchtime 1x ./internal/journal/` | PASS | `BenchmarkRebuildWorkspace100kEvents-6` reported `22482 ms/rebuild`, below 30000; last line `ok  \tgithub.com/mbell-nzl/fid/internal/journal\t25.713s` (terminated `exit_code: 0`). |
| Task 18 MF-2 regression | `go test -tags integration -run '^TestP23CallerProvidedReviewIDsRemainDistinct$' -count=1 ./cmd/fid/integration/` | PASS | `ok  \tgithub.com/mbell-nzl/fid/cmd/fid/integration\t0.643s` (terminated `exit_code: 0`). |
| Task 18 shared corpus | `go test -tags integration -run '^(TestP23FailureCorpusHasEveryAcceptedCode|TestP23SharedLeafConformance|TestP23SharedAmendedLeafConformance|TestGrantActionSetIsExactlyThirtyNames)$' -count=1 ./cmd/fid/integration/` | PASS | `ok  \tgithub.com/mbell-nzl/fid/cmd/fid/integration\t10.487s` (terminated `exit_code: 0`). |
| Task 18 unchanged lint gate | `golangci-lint run ./...` | PASS | `0 issues.` (terminated `exit_code: 0`). |
| DEC-266 gate regressions | `go test -run '^(TestGateRegisterUnknownActorReturnsContractAuthorityFailure|TestGateRegisterDeniedActorReturnsContractAuthorityFailure)$' -count=1 ./cmd/fid` | PASS | `ok  \tgithub.com/mbell-nzl/fid/cmd/fid\t0.328s` (terminated `exit_code: 0`). |
| Task 17 and remaining handoff selection | `go test -tags integration -run '^(TestFinalCapabilityMatrixMatchesAcceptedTable|TestP23AttentionResolutionThroughReviewRecordedAndFindingResolvedRealProcess|TestP23CancelledFamiliesAbsentFromProduct|TestP23HumanFixtures|TestP23HumanFixtureStepsRemainRunnable)$' -count=1 ./cmd/fid/integration/` | PASS | `ok  \tgithub.com/mbell-nzl/fid/cmd/fid/integration\t11.430s` (terminated `exit_code: 0`). |
| Shared format gate | `gofmt -l cmd internal` | PASS | `exit 0 (no output)` (terminated `exit_code: 0`). |

## Per-Task Findings

### Task 16 - Portable recovery for the review and judgment domain

#### Must-fix

- None.

#### Improvement

- None.

Observation: `TestPortableAssertCorruptEventsRejected` at `cmd/fid/integration/p23_portable_recovery_test.go:88` builds the nine Event examples from four real fixtures. Its helper clones the same clean `.fid` per subtest, opens a fresh DB handle, confirms the clone matches the baseline, rejects through `journal.ProjectEvent`, and compares every table's rows before/after. The 28 verbose paths and digests independently confirm MF-1's isolation requirement. The separate round-trip test retains positive import/rebuild evidence and schema-17 rejection.

### Task 17 - Final capability matrix and availability audit

#### Must-fix

- None.

#### Improvement

- None.

Observation: the fresh matrix and Attention tests passed; `task-17-done.md` retains both the prequalification `human_outcome_resolution:false` evidence and the final `true` row with a working path. Its five required commands still pass report-lint. No repair commit touched task 17.

### Task 18 - Conformance assembly and wire-surface audit

#### Must-fix

- None.

#### Improvement

- None.

Observation: `cmd/fid/integration/p23_conformance_test.go:799-852` uses fake-authored generated markers as the only normalization positions. `p23_amended_evidence_oracle_test.go:9-25` copies `review_id` from the request, and the new regression at `p23_conformance_test.go:104-139` checks two distinct valid IDs against their inputs and keeps their normalized results different; it also guards caller `id` and `envelope_id`. The seven oracle files read fixture requests and contract values, not fid result maps. The full shared corpus still compares complete accepted values and rejection envelopes, with Event-count checks; its 23-row report matrix cites `public-interface.md` or `contracts.md` for each inapplicable rejection class. The independently marked `failures.json` partition remains 19 declared = 18 added + sole retained `requirement-not-found`, and the seven grant-action introduction rows remain in the report.

Authority cross-check: `rg -n 'CodeAuthority(Missing|UnknownActor|Denied)' cmd internal --glob '*.go'` found the expected ten non-test producer lines: eight converged paths in `claim.go`, `run.go`, `finding.go`, `concern.go`, `deviation.go`, `decision.go`, and `human_outcome.go` twice, plus unchanged `gate.go:325/327`. `internal/cli/codes.go:286-290,310` gives all four `authority-*` codes `operation`, nullable `actor_id`, `required_actions`, `required_capabilities`, with no subject. The gate tests reach the two formerly latent paths and assert exact details, null subject, exit 5, and no Event. `p23FakeLeafAdapter.UntrustedActor` and `.MissingCapability` contain literal four-field contract expectations, including the `human-outcome.record` row; no fake authority expectation is copied from fid output.

## Prior Finding Resolution

| Prior finding | Owning task | Result | Fresh evidence |
|---|---|---|---|
| MF-1: 28 corrupt Event cases reused one workspace and connection (VAL-076) | task-16 | RESOLVED | Commit `4261575` moves workspace cloning and DB opening inside each subtest; commit `c2f0290` corrects report gate formatting. Source lines 129-220, the exact named test, 28 verbose subtests with distinct paths and equal 59-table digests, benchmark gate, and report-lint all passed. |
| MF-2: field-name normalization erased caller Review IDs (VAL-068) | task-18 | RESOLVED | Commit `3227197` changes the comparator to fake-declared generated positions and adds `TestP23CallerProvidedReviewIDsRemainDistinct`; commit `1d2e783` records committed-report lint. Source lines 104-139 and 787-852, the new regression, shared corpus, unchanged lint gate, and report-lint all passed. |

## Blocker Candidates

- None. Both original findings have bounded, independently verified repairs; no new protected decision is required.

## Existing Escalation References

- ESC-23-010 has the accepted DEC-265/266 dispositions. ESC-23-011 records the earlier Developer-capability interruption; review-13 records the restored-capability resume. Neither supplies a current finding in this review.

## Final Lines

```text
Review written: project-docs/docs/phases/phase-23/qa/review-14.md
Must-fix: 0  Improvement: 0  Blocker candidates: 0
VERDICT: PASS
```

OUTCOME PASS
