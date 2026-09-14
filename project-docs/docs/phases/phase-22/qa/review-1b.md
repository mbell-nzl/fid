# Phase 22 - Mid-Phase Review 1 (tasks 1-3)

Reviewer: qa-reviewer
Date: 2026-09-06
Tasks covered: phase-22/task-1, phase-22/task-2, phase-22/task-3
Review type: FRESH INDEPENDENT RE-REVIEW
Prior review: project-docs/docs/phases/phase-22/qa/review-1.md

## Verdict

PASS

All three Must-fix findings from review-1 are resolved by freshly inspected
commits, behavior assertions, and independent command evidence. This review
does not assert final QA or Human outcomes for later Phase 22 assertions.

## Scope Reviewed

- Task specs:
  - `project-docs/docs/phases/phase-22/tasks/task-1.md`
  - `project-docs/docs/phases/phase-22/tasks/task-2.md`
  - `project-docs/docs/phases/phase-22/tasks/task-3.md`
- Completion reports:
  - `project-docs/docs/phases/phase-22/tasks/task-1-done.md`
  - `project-docs/docs/phases/phase-22/tasks/task-2-done.md`
  - `project-docs/docs/phases/phase-22/tasks/task-3-done.md`
- Validation contract and QA plan:
  - `project-docs/docs/phases/phase-22/qa/validation-contract.md`
  - `project-docs/docs/phases/phase-22/qa/qa-plan.md`
- Prior review and recovery authority:
  - `project-docs/docs/phases/phase-22/qa/review-1.md`
  - `project-docs/docs/phases/phase-22/decisions.md` (DEC-189)
  - `project-docs/docs/phases/phase-22/triage-log.md` (TRG-22-005)
- Commits and changed files:
  - `daca8392a580ebdd2b6a8d3dab011910e39a779e` — `test(fid): stabilize live output coverage fixture`
  - `b6935f3772f29aa01fc200b1538d029a00dd9aa3` — `docs(task): complete phase-22 task-2 report repair`
  - `d6d140a12bb36e9926f009aba8fc49c1446b39c2` — `test(fid): cover claim renewal run lineage`
  - `8f22fa9c9de86e165214b47070064427bae385b3` — `docs(task): finalize completion report for phase-22 task-3`

## Summary

- Review outcome: PASS
- Must-fix findings: 0
- Improvement findings: 1
- Blocker candidates: 0
- Existing escalation references: none
- Covered validation assertions: VAL-001, VAL-002, VAL-003, VAL-004, VAL-005, VAL-006, VAL-007, VAL-008, VAL-009, VAL-010, VAL-011, VAL-012, VAL-079

All reported implementation and report hashes resolve as commits and their
subjects/`Refs` footers match the applicable task. `report-lint` passed for all
three task reports. The task-2 and task-3 reports are committed and clean at
`b6935f3` and `8f22fa9`, respectively. The sole uncommitted path before this
review artifact is the Runner-owned `project-docs/docs/phases/phase-22/run-ledger.md`.

## Validation Contract Coverage

| VAL ID | Claimed by | Status | Evidence |
|---|---|---|---|
| VAL-001 | task-1 | ON TRACK | `internal/db/recovery_migration_test.go` covers schema-15 tables, constraints, and invalid inserts; the tagged fresh-binary fixture is recorded in task-1's committed report. |
| VAL-002 | task-1 | ON TRACK | The task-1 migration test takes the database down to zero and back up, then verifies the recovery tables are restored empty. |
| VAL-003 | task-1 | ON TRACK | `internal/journal/recovery_schema_test.go` checks matching live-apply and replay rejection for unavailable and corrupt recovery effects. |
| VAL-004 | task-2 | ON TRACK | `internal/workflow/attention_test.go:12` retains original Stop details while recording separate Attention resolution. |
| VAL-005 | task-2 | ON TRACK | `internal/workflow/stop_test.go:11` proves same-key Stops share one open Attention while retaining distinct Stop facts. |
| VAL-006 | task-2 | ON TRACK | `internal/workflow/attention_test.go:12` rejects an incomplete qualifying-event set without effects and accepts the exact sorted set. |
| VAL-007 | task-2 | ON TRACK | `internal/workflow/stop_test.go:82` proves routine status is effect-free before explicit Stop creation; the freshly rerun tagged fixture also proves routine `pass` and a warning leave both journal bytes and Attention count unchanged at the external boundary. |
| VAL-008 | task-3 | ON TRACK | The task-3 report cites clock-controlled before/equal/after expiry and post-expiry/release rejection tests. |
| VAL-009 | task-3 | ON TRACK | The tagged real-process journey records one mutating Claim winner, one typed conflict, and no extra Event. |
| VAL-010 | task-3 | ON TRACK | The observer scenario admits only the activated bound harness and rejects a Builder and forged executor without an observation. |
| VAL-011 | task-3 | ON TRACK | `cmd/fid/generic_claim_test.go:270` starts a Claim-bound Run, renews from version 1 to 2, checks all immutable lineage/binding fields through `run get`, rejects stale version-1 release with unchanged persistence, and accepts version-2 release. The focused independent race run passed. |
| VAL-012 | task-3 | ON TRACK | `cmd/fid/generic_run_test.go:205` covers accepted optional-Obligation binding and mismatched/reused-executor rejections. |
| VAL-079 | task-2 | ON TRACK | Direct positive and rejection paths are present in the Stop/Attention package tests; the independently rerun coverage command reports 33.7% for `internal/workflow`. |

## Gate And Evidence Checks

| Check | Command or Evidence | Result | Notes |
|---|---|---|---|
| Task-1 exact and QA-plan commands | Task-1 report versus task-1 specification and QA-plan section 5 | PASS | Format, lint, race/coverage, coverage report, tagged integration, reference-instance inventory, and schema-15 rebuild sample are recorded with terminating receipts; `report-lint` passed. |
| Task-2 exact and QA-plan commands | Task-2 report versus task-2 specification and QA-plan section 5 | PASS | Required baseline/final workflow coverage, format, lint, byte-identical race/coverage gate, coverage report, and tagged integration receipts are present; `report-lint` passed. |
| Task-3 exact and QA-plan commands | Task-3 report versus task-3 specification and QA-plan section 5 | PASS | Format, lint, race/coverage, coverage report, and tagged integration receipts are present for the original task and its MF-3 rework; `report-lint` passed. |
| MF-3 reproduced behavior | `go test -race ./cmd/fid -run '^TestClaimRenewalAfterRunPreservesLineageAndFencesClaimVersions$' -count=1` | PASS | `ok  github.com/mbell-nzl/fid/cmd/fid  6.301s` (terminated `exit_code: 0`). |
| MF-1/MF-2 reproduced boundary | `go test -tags integration -race ./cmd/fid/integration/... -run '^TestStopAttentionExternalBoundaryKeepsRoutineStatusAndWarningsEffectFree$' -count=1` | PASS | `ok  github.com/mbell-nzl/fid/cmd/fid/integration  2.191s` (terminated `exit_code: 0`). |
| Task-2 representative coverage | `go test -race -cover ./internal/workflow` | PASS | `ok  github.com/mbell-nzl/fid/internal/workflow (cached) coverage: 33.7% of statements` (terminated `exit_code: 0`). |
| Task-2/task-3 tagged suite | `go test -tags integration -race ./cmd/fid/integration/...` | PASS | `ok  github.com/mbell-nzl/fid/cmd/fid/integration  8.875s` (terminated `exit_code: 0`). |
| DEC-189 mechanism | Coverage-built `TestLiveGoldenChild` under `env -i` | PASS | Fresh diagnostic emitted `warning: GOCOVERDIR not set, no coverage data emitted` and the child terminated at expected `exit_code: 23`; `daca839` adds only the fresh literal `GOCOVERDIR` fixture entry. |
| DEC-189 repair scope and goldens | Actual diffs and blob comparison from `4dd41c9` through `HEAD` | PASS | Harness change is limited to `cmd/fid/live_output_test.go`; `live-frames.jsonl` and `human-session.txt` have identical pre-repair, repair-report, and current blob IDs. |
| Worktree and escalation state | `git status --short`; `tl_escalations_current.md` | PASS | Before writing this review, only the Runner ledger was modified; the current escalation register lists no open ESC. |

## Per-Task Findings

### Task 1 - Recovery schema and replay foundation

#### Must-fix

- None.

#### Improvement

- None. The committed migration/replay evidence continues to align with VAL-001 through VAL-003 and the task's complete command record.

### Task 2 - Typed Stops and Attention domain

#### Must-fix

- None.

#### Improvement

- [internal/workflow/stop.go:1] The file remains 668 lines, over the repository's 400-line non-test-file target. Fix direction: extract Attention-specific behavior to `internal/workflow/attention.go` when a later task changes this area.

### Task 3 - Lease lifecycle and executor observations

#### Must-fix

- None.

#### Improvement

- None. The new test asserts observable CLI contracts and persistence stability rather than package layout, and its report correctly attributes VAL-011 to that test.

## Prior Finding Resolution

| Prior finding | Owning task | Result | Fresh evidence |
|---|---|---|---|
| MF-1 — tagged Stop/Attention external-boundary fixture | task-2 | RESOLVED | `cmd/fid/integration/task22_stop_attention_test.go:16` launches a fresh binary, proves routine pass and warning paths leave journal bytes and open Attention count unchanged, and preserves the unavailable capability/no-leaf boundary. The focused tagged command passed independently. |
| MF-2 — routine status/warning no-op and explicit Stop evidence | task-2 | RESOLVED | `internal/workflow/stop_test.go:82` proves status is effect-free and explicit `RaiseStop` creates the effect. The tagged fixture provides the warning non-qualifying state. The independent workflow coverage and full tagged-suite commands passed. |
| MF-3 — Claim renewal after Run creation, immutable Run lineage/binding, and stale/current Claim-version admission | task-3 | RESOLVED | `d6d140a` adds the single task-owned scenario at `cmd/fid/generic_claim_test.go:270`; it captures original Run lineage, renews, rejects stale version 1 without persistence change, accepts current version 2, and rechecks the immutable lineage. The focused race command passed independently. |

## Repair Plan

Not required: no Must-fix finding remains.

## Blocker Candidates

- None.

## Existing Escalation References

- None. `project-docs/docs/current/tl_escalations_current.md` lists no open ESC; its archived ESC-22-001 resolution is not an active reference for this review.

## Cross-cutting observations

- Task-2 preserves the original failed repository-wide receipt, reports the later byte-identical green, and changes neither the protected output goldens nor their assertion semantics. The observed DEC-189 child diagnostic validates the stated causal mechanism.
- Task-2 and task-3 rework reports remain committed evidence artifacts after their test commits. `git diff --check` is clean and no dependency change appears in the reviewed commits.

## Final Lines

```text
Review written: project-docs/docs/phases/phase-22/qa/review-1b.md
Must-fix: 0  Improvement: 1  Blocker candidates: 0
VERDICT: PASS
```
