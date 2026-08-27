# Phase 21 - Mid-Phase Review 3b (Tasks 7-9)

Reviewer: qa-reviewer
Date: 2026-08-27
Tasks covered: phase-21/task-7, phase-21/task-8, phase-21/task-9
Review type: FRESH INDEPENDENT RE-REVIEW
Prior review: project-docs/docs/phases/phase-21/qa/review-3.md (Re-review of: review-3)

## Verdict

PASS

The bounded Task 9 repair resolves review-3 MF-1. Public stats now accepts
only `actor`, `epic`, and `model`; `role`, another invalid grouping, and
retired `--by` all reject with the current `invalid-argument` behavior and do
not mutate the journal. Task 7 and Task 8 report and commit integrity remain
valid.

## Scope Reviewed

- Task specs and completion reports for phase-21 tasks 7, 8, and 9.
- Phase-21 current context, DEC-166, QA plan, validation contract, review-3,
  and the mid-phase review template.
- Task commits `817ad7c50469c51dc12f88a4bc081ecd7a384d7e`,
  `099906daf6795f79d9e669d562a88702dbc33068`, and the bounded repair
  `98217f4f3d1b20f6f2baa9201c8b16f84d6afeb3`; current report commit
  `62679ff91ea45e65d6e46aeede057a2af5cdf53b`.

## Summary

- Review outcome: PASS
- Must-fix findings: 0
- Improvement findings: 0
- Blocker candidates: 0
- Existing escalation references: ESC-21-007, ESC-21-008 (resolved historical references)
- Covered validation assertions: VAL-023, VAL-024, VAL-025, VAL-026, VAL-027, VAL-028, VAL-029, VAL-030, VAL-031

## Validation Contract Coverage

| VAL ID | Claimed by | Status | Evidence |
|---|---|---|---|
| VAL-023 | task 7 | ON TRACK | Task-7 report maps the all-match and single-dimension failure matrix to `cmd/fid/authority_test.go`; its committed report passes canonical report-lint. |
| VAL-024 | task 7 | ON TRACK | Task-7 report's seeded project-local role attribution adversary rejects all protected operations; report and commit integrity independently pass. |
| VAL-025 | task 8 | ON TRACK | Task-8 report records the actor/idempotency/Event-v2 retained Issue corpus and `--by` rejection; canonical report-lint passes. |
| VAL-026 | task 8 | ON TRACK | Task-8 report records the fixed-clock, nonfuture Event-v2 ordering evidence; canonical report-lint passes. |
| VAL-027 | task 8 | ON TRACK | Task-8 report cites the source-budget inventory and Issue corpus; its committed five-command gate is complete. |
| VAL-028 | task 9 | ON TRACK | Task-9 report records retained-family operation/read matrices and actor/idempotency request negatives; the final task gate passes. |
| VAL-029 | task 9 | ON TRACK | Task-9 report maps public policy-command absence and structural-predicate negatives to the final gate. |
| VAL-030 | task 9 | ON TRACK | Task-9 report maps unavailable AC/Concern/vacuum behavior and no-side-effect checks to the final gate. |
| VAL-031 | task 9 | ON TRACK | `98217f4` removes `role` normalization in `internal/stats/stats.go:142`; `cmd/fid/stats_test.go:37` exercises all qualifying and non-qualifying CLI states with journal-byte checks. Fresh package and real-process evidence below pass. |

## Gate And Evidence Checks

| Check | Command or Evidence | Result | Notes |
|---|---|---|---|
| Task 7 integrity | Canonical two-argument report-lint and `git rev-parse` for `817ad7c` | PASS | `OUTCOME COMPLETE`; hash resolves and subject is `feat(policy): add envelope authority evaluator`. Its five exact spec commands are all recorded PASS with observed `exit_code: 0`. |
| Task 8 integrity | Canonical two-argument report-lint and `git rev-parse` for `099906d` | PASS | `OUTCOME COMPLETE`; hash resolves and subject is `fix(issue): converge retained actor events`. Its five exact spec commands are all recorded PASS with observed `exit_code: 0`. |
| Task 9 bounded diff | `git show 98217f4` and diff through `62679ff` | PASS | Repair commit changes only `internal/stats/stats.go`, `cmd/fid/stats_test.go`, and Task-9 evidence; current commit changes only the Task-9 report evidence. `git diff --check` is clean. |
| Task 9 exact package gate, independently rerun | `go test -race ./internal/policy ./internal/workflow ./internal/issues ./internal/labels ./internal/links ./internal/requirements ./internal/activity ./internal/stats ./internal/maintenance ./internal/cli ./cmd/fid/...` | PASS | Final token: `ok github.com/mbell-nzl/fid/cmd/fid (cached)` (terminated `exit_code: 0`). |
| Task 9 report control, independently rerun | `project-docs/tools/report-lint.sh project-docs/docs/phases/phase-21/tasks/task-9-done.md project-docs/docs/phases/phase-21/tasks/task-9-retained-actor-policy-removal.md` | PASS | `OUTCOME COMPLETE` (terminated `exit_code: 0`). |
| VAL-031 qualifying and non-qualifying real-process probe | Fresh initialized workspace with `/tmp/fid-qa-review-3b --json stats tokens --group-by <value>` | PASS | `actor`, `epic`, and `model` each returned successful `stats.tokens` JSON (terminated `exit_code: 0`). `role`, `team`, and `--by developer` each returned JSON `code:"invalid-argument"` (terminated `exit_code: 1`). The journal SHA-256 was `73089addfb3e10b123fbef6a955dba50ba2aab00674603250b9a7ea758114af8` before and after all six reads. |

## Per-Task Findings

### Task 7 - Generic Authority Evaluator

#### Must-fix

- None.

#### Improvement

- None.

### Task 8 - Issue Actor And Event Convergence

#### Must-fix

- None.

#### Improvement

- None.

### Task 9 - Retained Inventory Attribution And Judgment Freeze

#### Must-fix

- None.

#### Improvement

- None.

## Prior Finding Resolution

| Prior finding | Owning task | Result | Fresh evidence |
|---|---|---|---|
| MF-1 from review-3: public stats accepted retired `--group-by role`; evidence lacked model success and closed-value rejection/no-mutation proof for VAL-031. | phase-21/task-9 | RESOLVED | Repair `98217f4f3d1b20f6f2baa9201c8b16f84d6afeb3` deletes the `role -> actor` normalization and adds `TestStatsTokensAcceptsOnlyClosedGroupByValues`. Fresh exact package-gate PASS, report-lint `OUTCOME COMPLETE`, and the initialized-workspace probe prove actor/epic/model success; role/team/`--by` `invalid-argument` rejection; and identical journal digest. |

## Blocker Candidates

- None.

## Existing Escalation References

- ESC-21-007 and ESC-21-008 are resolved historical Task-9 recovery references under DEC-165 and DEC-166. They do not govern MF-1 or this PASS re-review.

## Final Lines

```text
Review written: project-docs/docs/phases/phase-21/qa/review-3b.md
Must-fix: 0  Improvement: 0  Blocker candidates: 0
VERDICT: PASS
```
