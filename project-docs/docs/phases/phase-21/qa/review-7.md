# Phase 21 - Mid-Phase Review 7 (tasks 19-21)

Reviewer: qa-reviewer
Date: 2026-09-02
Tasks covered: phase-21/task-19, phase-21/task-20, phase-21/task-21
Review type: INITIAL
Prior review: project-docs/docs/phases/phase-21/qa/review-6c.md (preceding Window 6; no Window 7 finding lineage)

## Summary

- Review outcome: REPAIR
- Must-fix findings: 1
- Improvement findings: 0
- Blocker candidates: 0
- Existing escalation references: ESC-21-014, closed by DEC-173; provenance only
- Covered validation assertions: VAL-002, VAL-069, VAL-070, VAL-071, VAL-074, VAL-075, VAL-076, VAL-077, VAL-078, VAL-079, VAL-080, VAL-081, VAL-082, VAL-083

## Validation Contract Coverage

| VAL ID | Claimed by | Status | Evidence |
|---|---|---|---|
| VAL-002 | task-20 | ON TRACK | `TestMigration0014RetiresV1AuthorityExecutionWithoutReconstruction` verifies the exact post-0014 drops, final `claims` shape, retained frozen tables, rule deletion, and empty-state down/up at `internal/db/operational_kernel_migration_test.go:47`; the focused rerun passed. |
| VAL-069 | task-19 | ON TRACK | `TestCreateIsByteIdenticalWithCompleteDeterministicInventory` byte-compares repeated archives and asserts canonical inventory/source coverage at `internal/export/export_test.go:20`; fresh focused export tests passed. |
| VAL-070 | task-19 | NOT ON TRACK | Only one aggregate-sequence collision is exercised at `internal/export/export_test.go:70`; QA-68 requires a non-mutating corruption corpus for every safety and integrity defect. MF-19-01 applies. |
| VAL-071 | task-19 | ON TRACK | `cmd/fid/export_test.go:9` checks unknown `export import` and the exact unavailable capability; the required source absence command returned `exit_code: 1` with no output. |
| VAL-074 | task-20 | ON TRACK | The current broad table and Event inventories returned their required `exit_code: 0` matches, all mapped by the report's per-row dispositions, with `UNCLASSIFIED=0`. |
| VAL-075 | task-20 | ON TRACK | Fresh role-policy and `--by` absence searches returned `exit_code: 1`; the report records all ten predeclared removal controls with their declared observed non-zero terminations. |
| VAL-076 | task-21 | ON TRACK | Fresh real-process integration test passed; it asserts Event-v2 role-free bootstrap and no `roles` table before and after rebuild at `cmd/fid/integration/task21_maintenance_rebuild_test.go:15`. |
| VAL-077 | task-21 | ON TRACK | The same real-process test deletes the projection, rebuilds, and byte-compares capabilities while reasserting role-free bootstrap and absent `roles` table. |
| VAL-078 | task-21 | ON TRACK | `internal/maintenance/check_test.go:150` asserts ordered `event-sequence`, `aggregate-sequence`, and `event-binding` findings. |
| VAL-079 | task-21 | ON TRACK | The same test asserts ordered `attempt-receipt-integrity` and `artifact-integrity` findings. |
| VAL-080 | task-21 | ON TRACK | The same test asserts `frozen-state-present` and `removed-state-residue`, then verifies unchanged journal bytes and SQLite change count. |
| VAL-081 | task-21 | ON TRACK | `cmd/fid/task21_maintenance_rebuild_test.go:10` compares capabilities, Envelope, Evidence evaluation, and resume reads before and after a deleted-projection rebuild. |
| VAL-082 | task-21 | ON TRACK | The same test byte-compares export archives and compares export validation results before and after rebuild. |
| VAL-083 | task-21 | ON TRACK | The report has three distinct fresh-process result lines, 27,804 ms, 27,866 ms, and 27,906 ms, on the DEC-158 reference identity. Independent sample: `BenchmarkRebuildWorkspace100kEvents-6 1 26804006990 ns/op 26804 ms/rebuild`; command terminated `exit_code: 0`. |

QA-81 through QA-85 are assigned to task 22 and final QA by `qa-plan.md` section 7.9. They are not claimed by this window and remain deferred, not a finding against tasks 19-21.

## Gate And Evidence Checks

| Check | Command or evidence | Result | Notes |
|---|---|---|---|
| Task 19 exact command audit | `task-19-done.md` against `task-19-export-create-validate.md` | PASS | All six required commands are recorded PASS with a token and observed termination. Historical lint failure `1 issues:` (`exit_code: 1`) and byte-identical successful rerun `0 issues.` (`exit_code: 0`) are preserved at `task-19-done.md:59`. |
| Task 20 exact command audit | `task-20-done.md` against amended `task-20-v1-removal.md` | PASS | The package gate, DEC-173 tagged integration gate, five `ls` absence controls, five `rg` absence controls, two broad classified inventories, format, lint, race/coverage, and coverage report are all recorded with their required statuses and observed terminations. Historical failed attempts and later reruns are retained. |
| Task 21 exact command audit | `task-21-done.md` against `task-21-maintenance-rebuild.md` | PASS | Every required command is recorded PASS with observed termination, including all six identity commands, three separately recorded benchmark invocations, benchmark suite, and task gate. Historical focused race and coverage failures plus byte-identical passing reruns are recorded in the deviation section. |
| Task 19 report control | `project-docs/tools/report-lint.sh project-docs/docs/phases/phase-21/tasks/task-19-done.md project-docs/docs/phases/phase-21/tasks/task-19-export-create-validate.md` | PASS | `PASS report-lint completed with no failures` (terminated `exit_code: 0`). |
| Task 20 report control | `project-docs/tools/report-lint.sh project-docs/docs/phases/phase-21/tasks/task-20-done.md project-docs/docs/phases/phase-21/tasks/task-20-v1-removal.md` | PASS | `PASS report-lint completed with no failures` (terminated `exit_code: 0`). This independently confirms the current report repair despite its preserved historical lint failures. |
| Task 21 report control | `project-docs/tools/report-lint.sh project-docs/docs/phases/phase-21/tasks/task-21-done.md project-docs/docs/phases/phase-21/tasks/task-21-maintenance-rebuild.md` | PASS | `PASS report-lint completed with no failures` (terminated `exit_code: 0`). |
| Reproduced task 19 command | `rg -n --glob '*.go' --glob '!**/*_test.go' 'Use:\\s*\"import|newExportImport|export\\.import' cmd/fid internal` | PASS | No match, expected `exit_code: 1`; terminated `exit_code: 1`. Fresh focused test of deterministic create, one corruption rejection, and import/capability behavior passed: `ok github.com/mbell-nzl/fid/cmd/fid 0.119s` (terminated `exit_code: 0`). |
| Reproduced task 20 commands | Exact role-policy absence and broad table/Event inventory recipes | PASS | Role-policy absence returned no output (expected and terminated `exit_code: 1`). Broad table inventory terminated `exit_code: 0` and contained only the report-classified final `claims` and frozen diagnostic readers; Event inventory terminated `exit_code: 0` and contained only report-classified Event-v2 families. Focused migration/removal tests passed: `ok github.com/mbell-nzl/fid/cmd/fid 0.091s` (terminated `exit_code: 0`). |
| Reproduced task 21 command | `go test -run '^$' -bench '^BenchmarkRebuildWorkspace100kEvents$' -benchmem -benchtime 1x -timeout 30m ./internal/journal` | PASS | `BenchmarkRebuildWorkspace100kEvents-6 1 26804006990 ns/op 26804 ms/rebuild`; `ok github.com/mbell-nzl/fid/internal/journal 29.424s` (terminated `exit_code: 0`). The benchmark's measured rebuild time is below 30 s. Focused maintenance/parity and fresh real-process role-free init/rebuild tests also passed, terminating `exit_code: 0`. |
| Public/frozen/deferred boundary | Current sources, contracts, decisions, and triage history | PASS | Migration 0014 remains the specified drop/rename/delete sequence. DEC-173's bounded task-20 amendment explicitly preserved commands, VALs, migration, frozen tables, public interface, Event schema, and task-21 rebuild threshold. The only production Attention match is the required unavailable `attention.v1` capability row; no Phase 22 Attention mutation is present. |
| Commit and worktree provenance | `git rev-parse --verify`, `git show -s`, `git diff --check`, `git status --short` | PASS | Reported commits `2518931`, `2d34575`, and `76ad781` resolve and each has the matching subject and `Refs: phase-21/task-N`. `git diff --check 22ad81d..053ab699` terminated `exit_code: 0`. HEAD is `053ab699708cebf38fef9d600dd5e4ef853e58cb`; only parent-owned `run-ledger.md` and `triage-log.md` are dirty. |

## Task 19 - Deterministic Export Create And Validate

### Must-fix

- [internal/export/export_test.go:70] MF-19-01: VAL-070 evidence covers only an aggregate-sequence collision, while the accepted QA-68 contract requires the full safety and integrity defect corpus, each rejected as `archive-invalid` without mutating the project. Owning task: task-19. Repair boundary: add task-19 export tests only, covering every accepted validation check code and verifying unchanged `.fid` state for each rejection; do not change export semantics, public error codes, archive format, or the QA contract. Required evidence: the complete corpus test passes, each row asserts its expected code and before/after project digest, the task-19 exact focused gate is rerun byte-identically, and the two-argument report lint passes. Re-review scope: fresh independent review of task 19's corpus tests, exact task-19 command, report control, and VAL-070; retain the current task-20/task-21 evidence unless the repair touches it.

### Improvement

- None.

## Task 20 - V1 Removal And Migration 0014

### Must-fix

- None.

### Improvement

- None.

## Task 21 - Final-Schema Maintenance, Init, Rebuild, And Export Parity

### Must-fix

- None.

### Improvement

- None.

## Repair Plan

| Finding | Owning task | Affected assertion | Exact repair boundary | Required evidence | Fresh re-review scope |
|---|---|---|---|---|---|
| MF-19-01 | task-19 | VAL-070 | `internal/export/export_test.go` and, only if necessary for CLI-observable coverage, existing task-19 export test files; add full defect-corpus/no-mutation assertions without changing accepted behavior or contracts. | One corpus row per closed `archive-invalid` check code, each with identical before/after `.fid` state; byte-identical task-19 focused gate and two-argument report-lint PASS with observed terminations. | Inspect the corpus, report, and implementation; rerun the affected test and exact task-19 absence command. Confirm task-20/task-21 commits and evidence remain unchanged. |

## Cross-cutting observations

- Task 20 respects DEC-173: its fresh broad recipes remain broad and classified, the role/`--by` absence controls are zero-match proofs, `UNCLASSIFIED=0` is recorded for both inventories, and the task-21 coordination constraints were not weakened.
- Task 21 respects DEC-158: the report records three distinct fresh-process samples and all six reference identity commands. The independent fourth sample remains below the same protected threshold; it is corroboration, not a replacement for the required three report samples.
- The preceding Window 6 review is PASS and has no Window 7 Must-fix lineage. No QA finding is self-resolved in this artifact.

## Blocker Candidates

- None. MF-19-01 is a bounded, decision-free task-19 test-evidence repair; it does not require a public-interface, archive-schema, security-policy, or validation-contract change.

## Existing Escalation References

- ESC-21-014, closed by DEC-173; provenance only. Its task-20 scope contradiction was resolved with no change to the command, expected-status, VAL, migration, frozen-table, public-interface, Event-schema, or task-21 benchmark contracts.

## Final Lines

Review written: project-docs/docs/phases/phase-21/qa/review-7.md
Must-fix: 1  Improvement: 0  Blocker candidates: 0
VERDICT: REPAIR
