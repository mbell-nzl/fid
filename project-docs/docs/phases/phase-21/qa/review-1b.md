# Phase 21 - Mid-Phase Review 1b (tasks 1-3)

Reviewer: qa-reviewer
Date: 2026-08-23
Tasks covered: phase-21/task-1, phase-21/task-2, phase-21/task-3
Review type: FRESH INDEPENDENT RE-REVIEW
Prior review: project-docs/docs/phases/phase-21/qa/review-1.md
Re-review of: review-1

## Summary

- Review outcome: PASS
- Must-fix findings: 0
- Improvement findings: 0
- Blocker candidates: 0
- Existing escalation references: none; ESC-21-005 was resolved by DEC-163 before the repaired review-owner sequence.
- Covered validation assertions: VAL-001, VAL-003, VAL-004, VAL-005, VAL-006, VAL-007, VAL-008, VAL-009, VAL-010

## Reproduced Command Evidence

All commands ran from `/home/mark/projects/fid` in this re-review session.

| Command | Result | Verbatim token and observed termination |
|---|---|---|
| `go test -race ./internal/db/...` | PASS | `ok  github.com/mbell-nzl/fid/internal/db (cached)` (terminated `exit_code: 0`) |
| `go test -race ./internal/workflow ./internal/cli ./cmd/fid/...` | PASS | `ok  github.com/mbell-nzl/fid/cmd/fid (cached)` (terminated `exit_code: 0`) |
| `project-docs/tools/report-lint.sh project-docs/docs/phases/phase-21/tasks/task-1-done.md project-docs/docs/phases/phase-21/tasks/task-1-schema-projection.md` | PASS | `OUTCOME COMPLETE` (terminated `exit_code: 0`) |
| `project-docs/tools/report-lint.sh project-docs/docs/phases/phase-21/tasks/task-3-done.md project-docs/docs/phases/phase-21/tasks/task-3-capabilities.md` | PASS | `OUTCOME COMPLETE` (terminated `exit_code: 0`) |
| `project-docs/tools/report-lint.sh project-docs/docs/phases/phase-21/tasks/task-2-done.md project-docs/docs/phases/phase-21/tasks/task-2-event-idempotency.md` | PASS | `OUTCOME COMPLETE` (terminated `exit_code: 0`) |

## Task 1 - Operational Schema And Projection

### Must-fix

- None.

### Improvement

- None.

### Evidence

- VAL-001 is resolved. [operational_kernel_migration_test.go](/home/mark/projects/fid/internal/db/operational_kernel_migration_test.go:14) captures the pre-0013 state, applies 0013, compares the complete new SQLite schema-object delta against fixed canonical DDL, verifies representative constraints, then verifies down and second-up state equality. The expected inventory at [operational_kernel_migration_test.go](/home/mark/projects/fid/internal/db/operational_kernel_migration_test.go:145) contains all 23 additive tables and four named indexes, so the normalized DDL comparison includes every CHECK and foreign key.
- The negative comparison table at [operational_kernel_migration_test.go](/home/mark/projects/fid/internal/db/operational_kernel_migration_test.go:47) proves added table, index, CHECK, and foreign-key drift is rejected; this supplies the non-qualifying state missing from review-1. The inspected 0013 [up migration](/home/mark/projects/fid/migrations/0013_p21_operational_kernel.up.sql:1) and [down migration](/home/mark/projects/fid/migrations/0013_p21_operational_kernel.down.sql:1) match that exact additive/reverse inventory.
- `task-1-done.md` preserves the original red/green cycle and all five exact task gates with observed terminations. The fresh two-argument report lint above verifies required-command coverage and the reported implementation commit. Repair commits `553965143d3d61e7b178401caa4be5d46853136a` and `a169b4114dfbd6f3d19625dd577cf6867bc707cc` resolve and have the recorded task-scoped subjects.

## Task 2 - Event V2, Idempotency, And Role-Free Bootstrap

### Must-fix

- None.

### Improvement

- None.

### Evidence

- Prior accepted VAL-003 through VAL-008 evidence remains intact: the Event-v2 bootstrap, aggregate-sequence, idempotency replay/conflict/concurrency, and role-free compatibility tests named in `task-2-done.md` are present in `internal/bootstrap`, `internal/journal`, and `internal/policy`.
- The fresh two-argument report-lint result confirms all six Task 2 exact commands remain recorded with prior failed evidence, later byte-identical passing evidence, tokens, and observed terminations. The repaired review range contains no Task 2 implementation diff, so there is no post-review regression to attribute to this window.

## Task 3 - Capability Discovery

### Must-fix

- None.

### Improvement

- None.

### Evidence

- VAL-010 is resolved. [capabilities_test.go](/home/mark/projects/fid/cmd/fid/capabilities_test.go:102) first checks the qualifying and non-qualifying values of the complete 12-row matrix, including all seven false capability IDs.
- Its seven executable probes at [capabilities_test.go](/home/mark/projects/fid/cmd/fid/capabilities_test.go:140) cover Activation lifecycle, Claim lifecycle/list, export create/validate, AC judgment, Concern judgment, historical stats, and Attention absence. Each probe must fail without stdout, must expose the unavailable-surface failure, preserves byte-identical journal content and all projection-table counts, then re-discovers the corresponding false capability at [capabilities_test.go](/home/mark/projects/fid/cmd/fid/capabilities_test.go:185). This proves the named surfaces neither advertise nor authorize the deferred capability and introduce no journal or projection side effect.
- `task-3-done.md` retains its red/green evidence, required gate rows, and causal deviation. The fresh Task 3 package gate and two-argument report lint both pass. Repair commits `401f609d4d7e24c0ab5ca45d27a16424000ded9f` and `68cc7ab9c028d9588e2455e18ad56fa5a8474578` resolve and have the recorded task-scoped subjects.

## Prior Finding Resolution

- MF-1, task-1 / VAL-001: RESOLVED. Fresh inspection confirms a full fixed canonical inventory with explicit extra-object/CHECK/FK rejection, reversible SQL, preserved red/green and gate evidence, passing report lint, and a fresh exact DB-gate PASS.
- MF-2, task-3 / VAL-010: RESOLVED. Fresh inspection confirms all seven named executable probes, complete binary-state coverage, unavailable/no-authorization assertions, no journal/projection-side-effect assertions, preserved red/green and gate evidence, passing report lint, and a fresh exact package-gate PASS.

## Cross-cutting observations

- The ledger records TASK-RECOVERY, REPAIR, and TASK COMPLETE for task-1 before the corresponding task-3 sequence. This satisfies DEC-163's ordered review-owner repair condition.
- Task 4 remains unauthorized. DEC-163 requires this independent `review-1b` PASS artifact to close review-1; only the Phase Runner may then record the REVIEW event and route the next task. This review does not edit the ledger or authorize task 4.
- QA-plan command applicability is unchanged: the ledger records no commands beyond the task-spec gates for tasks 1-3, and the three completion reports mechanically confirm their respective exact command sets.

## Blocker Candidates

- None. The prior workflow-policy blocker is resolved by DEC-163, and no current evidence requires product, architecture, data-model, public-interface, dependency, security, infrastructure, or specification authority.

## Existing Escalation References

- None. `tl_escalations_current.md` records no open escalation; DEC-163 resolved ESC-21-005.

## Final Lines

Review written: project-docs/docs/phases/phase-21/qa/review-1b.md
Must-fix: 0  Improvement: 0  Blocker candidates: 0
VERDICT: PASS
