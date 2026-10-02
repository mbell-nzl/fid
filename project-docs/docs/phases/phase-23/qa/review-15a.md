# Phase 23 - Mid-Phase Review 15a (tasks 19-21)

Reviewer: qa-reviewer
Date: 2026-10-03
Tasks covered: phase-23/task-19, phase-23/task-20, phase-23/task-21
Review type: FRESH INDEPENDENT RE-REVIEW
Prior review: `project-docs/docs/phases/phase-23/qa/review-15.md`

## Verdict

PASS. Both report-only Must-fix findings from review-15 are resolved by committed, byte-identical retry evidence. This review found no new Must-fix or protected blocker candidate.

## Scope Reviewed

- Task specs and completion reports for tasks 19-21; the phase QA plan, validation contract, current handoff, DEC-268..270, `08_QUALITY_PLAN.md`, and review-15.
- Product commits `180f5168e3f1338cb758f8388f1c26fd6f8428d3` (task 19), `0ba005d7c511d7f5acd1517d1c39dd6e58d3e3e0` (task 20), `c6e070b` and `343b458fc4c1b7a726ae9218315edf319c5e40ba` (task 21). The full reported hashes were verified as commits; report-lint also verified their subjects and committed reports.
- Report-only repair commits `9ea1f5533000d7cb12a48ba0ce92d0729d3d268b` and `025aa410c9efb5c344398d4f45cce6c5798c7529`, each with its owning task's `Refs` trailer. `git diff e7c656f..HEAD --name-only` contains only the task-19 and task-21 completion reports. Runner-owned `run-ledger.md` and `triage-log.md` were already dirty and were not edited here.

## Summary

- Review outcome: PASS
- Must-fix findings: 0
- Improvement findings: 0
- Blocker candidates: 0
- Existing escalation references: ESC-23-012 (separate phase-end Human re-check)
- Covered validation assertions: VAL-085, VAL-086, VAL-087, VAL-088, VAL-089

## Validation Contract Coverage

| VAL ID | Claimed by | Status | Evidence and boundary |
|---|---|---|---|
| VAL-085 | task-19 | ON TRACK | `TestP23StopFailureCodeMatrix` runs 15 real-binary Human/JSON pairs and checks exact stderr, empty stdout, exit parity, stable JSON, and unchanged Event counts. The renderer's seven-code inventory is checked. The original three red selections now have exact green retries; the matrix reran independently below. |
| VAL-086 | task-21 | ON TRACK | `TestP23ReviewCyclesRealProcessJourney` checks an open Stop and denied Claim/Run after wrong class, rejected result, other Issue, and non-human authority; an approved same-Issue Outcome resolves it and admits rework. `TestP23ReviewReworkManualResolutionWithPriorOutcome` covers manual resolution. The journey reran independently below. |
| VAL-087 | task-21 | ON TRACK | `internal/evidence/execution_test.go:86` and `internal/workflow/budget_test.go:168,181` exercise the specified isolated active-Run call sites: both rework states are excluded while attempts and token controls remain. The task spec permits this isolated route because the simultaneous CLI state is unreachable under accepted Claim/Run rules. The shared fake covers the same exclusion and controls; task-21's tagged gate includes it. |
| VAL-088 | task-21 | ON TRACK | The real-process journey checks later consumed totals 2 and 3, remaining 0, and no new rework Stop. `TestP23StopFailureHumanAndJSONParity` checks the two exact Human lines; the matrix checks its JSON twin. The previously failed guidance selection now has an exact retry and an independent rerun below. |
| VAL-089 | task-20 | ON TRACK | `TestP23RequirementReferenceAndLinkRealProcess` checks stable and reference IDs, old-reference lookup, supersession, and links. `TestP23RequirementNewShapeJournalRoundTrip` checks five byte-identical JSON reads, unchanged journal digest, and idempotent retry after rebuild; it reran below. The report retains the DEC-270(b) final-binary `journal-replay-stopped` at line 18, observed exit 6, and identical before/after SHA-256. Review-15 independently reproduced that preserved-journal refusal. |

VAL-072 and VAL-073 fixture evidence remains preparatory. The Human owns the later comprehension and presentation judgments under the phase contract.

## Evidence

### Required commands and reports

Each task spec names the same five exact commands. QA-plan section 22 adopts the `-timeout 30m` forms; the additional full phase, benchmark, build, and full-checkptr controls belong to phase end. I compared every spec command with each report's `## Exact quality-gate results`: all five have a PASS row with an output token and observed `exit_code: 0` for each task. Task 21 retains its earlier lint failure and byte-identical passing retry. Task 20 retains its attempt-1 BLOCKED rows under the former protected replay contradiction and all five later PASS rows after DEC-270. These are Developer report results, not fresh QA runs.

| Exact task-spec command | Task 19 report | Task 20 report | Task 21 final report |
|---|---|---|---|
| `gofmt -l cmd internal` | `exit 0 (no output)`; exit 0 | same; exit 0 | same; exit 0 |
| `golangci-lint run ./...` | `0 issues.`; exit 0 | same; exit 0 | same; exit 0 |
| `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` | workflow `62.1%`; exit 0 | workflow `62.1%`; exit 0 | workflow `62.7%`; exit 0 |
| `go tool cover -func=coverage.out` | `60.5%`; exit 0 | `60.5%`; exit 0 | `60.7%`; exit 0 |
| `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...` | integration `215.530s`; exit 0 | integration `253.239s`; exit 0 | integration `291.867s`; exit 0 |

The reports record ordinary `cmd/fid` times of 215.387s, 344.249s, and 215.093s respectively; task 21's final workflow result was cached and its prior uncached run was 113.517s. Reported ordinary package and tagged integration times remain below the accepted 600s and 900s budgets. No report-only retry was counted as a new product gate.

### Fresh independent checks

All commands below ran from `/home/mark/projects/fid`. Each nested command reached its own `exit_code`; none was inferred from wrapper completion.

| Check | Fresh result |
|---|---|
| `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/... -run '^TestP23StopFailureCodeMatrix$' -count=1` | `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	6.638s` (terminated `exit_code: 0`); representative exact task-19 failed selection. |
| `go test -tags integration -run '^TestP23StopFailureHumanAndJSONParity$' -count=1 ./cmd/fid/integration` | `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	2.173s` (terminated `exit_code: 0`); exact task-21 guidance selection. |
| `go test -tags integration -run '^(TestP23ReviewCyclesRealProcessJourney|TestP23RequirementNewShapeJournalRoundTrip)$' -count=1 ./cmd/fid/integration` | `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	6.059s` (terminated `exit_code: 0`); affected rework journey and DEC-270(a) new-shape control. |
| `project-docs/tools/report-lint.sh <task-N-done.md> <task-N.md>` for N=19, 20, 21 | Each printed `OUTCOME COMPLETE` and `PASS report-lint completed with no failures` (each terminated `exit_code: 0`); actual paths were the corresponding phase-23 task files. |

## Per-Task Findings

### Task 19 - Human-mode failures from the Stop-error family

#### Must-fix

- None.

#### Improvement

- None.

Observation: `cmd/fid/integration/p23_human_failure_test.go:41,168` asserts externally visible Human/JSON parity and the seven-code inventory. The retry commit changes only its completion report; `cmd/fid/stop.go` and the tests remain as reviewed in review-15.

### Task 20 - Stable requirement identity on every reference

#### Must-fix

- None.

#### Improvement

- None.

Observation: the task-20 report and DEC-270 evidence remain committed and unchanged after review-15. The committed round-trip test passed independently in this review. Its pre-change refusal is the accepted DEC-270(b) state, with `journal-replay-stopped`, exit 6, and unchanged source digest; no new contradiction or task-20 repair evidence appeared.

### Task 21 - Resolvable review-rework exhaustion

#### Must-fix

- None.

#### Improvement

- None.

Observation: `cmd/fid/integration/p23_review_evidence_budget_test.go:77` asserts both qualifying and non-qualifying Outcome behavior and later single charges. The generic-admission helper is shared by gate execution and resume reasons. The retry commit changes only the completion report; the guidance test and affected journey passed independently.

## Prior Finding Resolution

| Prior finding | Owning task | Result | Fresh evidence |
|---|---|---|---|
| Review-15 MF-1: task-19's three expected-red Stop parity, code-matrix, and `gate/open-stop` selections lacked byte-identical passing retries | task-19; VAL-085 | RESOLVED | `task-19-done.md:34-63` retains all three `FAIL` rows (`FAIL`, terminated `exit_code: 1`). An exact string comparison of the backticked command bytes finds one later PASS for each original command, with verbatim integration `ok` tokens (`4.811s`, `5.502s`, `4.665s`) and terminated `exit_code: 0`. The causal note names the silent Human path and the gate fixture's missing budget baseline. Commit `9ea1f55` changes only the report; independent code-matrix rerun and committed-report lint pass above. |
| Review-15 MF-2: task-21's expected-red guidance selection lacked a byte-identical passing retry | task-21; VAL-088 | RESOLVED | `task-21-done.md:47-93` retains the `FAIL` row (`FAIL`, terminated `exit_code: 1`). The later PASS row has identical backticked command bytes, verbatim `ok` token (`2.270s`), observed `exit_code: 0`, and separately recorded 2.560s command wall time. The causal note names the missing guidance and implementation repair. Commit `025aa41` changes only the report; independent exact guidance and affected journey reruns and committed-report lint pass above. |

## Cross-cutting observations

- The repairs add evidence to task reports without changing product, tests, task commands, validation criteria, or the DEC-270 decision. Each validation assertion has both its relevant qualifying and non-qualifying evidence route. The later Human re-check under ESC-23-012 remains separate from this bounded deterministic review.

## Blocker Candidates

- None. No new protected decision or TRG/ESC identity is requested.

## Existing Escalation References

- ESC-23-012: existing phase-end Human re-check and final-review work only.

## Final Lines

Review written: project-docs/docs/phases/phase-23/qa/review-15a.md
Must-fix: 0  Improvement: 0  Blocker candidates: 0
VERDICT: PASS
OUTCOME PASS
