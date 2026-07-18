# Phase 17 - Mid-Phase Review 3 (Tasks 6-7)

Reviewer: qa-reviewer
Date: 2026-07-15
Tasks covered: phase-17/task-6, phase-17/task-7

## Verdict

CONTINUE

## Scope Reviewed

- Task specs:
  - `project-docs/docs/phases/phase-17/tasks/task-6.md`
  - `project-docs/docs/phases/phase-17/tasks/task-7.md`
- Completion reports:
  - `project-docs/docs/phases/phase-17/tasks/task-6-done.md`
  - `project-docs/docs/phases/phase-17/tasks/task-7-done.md`
- Validation contract: `project-docs/docs/phases/phase-17/qa/validation-contract.md`
- QA plan: `project-docs/docs/phases/phase-17/qa/qa-plan.md`
- Commits and finalized report state:
  - `f01a308377a59137fd3b9c509dad0742d6119aee` — `feat(maintenance): add locked vacuum command`
  - `20b9f0c` — finalized Task 6 report commit
  - `34c6cef754ab34f3d483ccee311dd3b05c3b2337` — `feat(cli): guard unshipped governance commands`
  - `4ebb32cab3a413d78413e67095a4f4295232ce8b` — finalized Task 7 report commit

## Summary

- Blocking findings: 0
- Improvement findings: 1
- Escalations raised: none
- Covered validation assertions: VAL-015, VAL-016

## Validation Contract Coverage

| VAL ID | Claimed by | Status | Evidence |
|---|---|---|---|
| VAL-015 | task-6 | ON TRACK | Fresh focused maintenance tests passed uncached for projected-state/activity and byte-identical journal preservation, held write-lock waiting/completion, and missing-projection recovery. Fresh CLI tests passed for `tl`/`human`, every missing/unknown/default-denied role, and policy-show parity. |
| VAL-016 | task-7 | ON TRACK | Fresh uncached CLI tests passed for all ten exact decision/escalation/triage leaves with future/missing-input probes, exact version errors, unchanged journal and governance-table counts, no-workspace precedence, root-help visibility, and the successful `requirement sync` control. |

## Gate And Evidence Checks

| Check | Command or Evidence | Result | Notes |
|---|---|---|---|
| Task 6 report evidence | `project-docs/tools/report-lint.sh project-docs/docs/phases/phase-17/tasks/task-6-done.md project-docs/docs/phases/phase-17/tasks/task-6.md` | PASS | All six exact task commands are recorded PASS; the committed report, implementation hash, and subject are valid. |
| Task 7 report evidence | `project-docs/tools/report-lint.sh project-docs/docs/phases/phase-17/tasks/task-7-done.md project-docs/docs/phases/phase-17/tasks/task-7.md` | PASS | All six exact task commands are recorded PASS; the committed report, implementation hash, and subject are valid. |
| Required integration gate | `go test -tags integration -race ./...` | PASS | Reproduced the exact Task 6/7 gate; all packages passed, including `cmd/fid`, `internal/maintenance`, `internal/db`, and `internal/policy`. |
| Vacuum behavior | `go test -race ./internal/maintenance -run '^TestVacuum(PreservesProjectedStateJournalAndActivity|WaitsForWriteLockThenCompletes|RefreshesMissingProjectionBeforeCompacting)$' -count=1 -v` | PASS | Reproduced both preservation states, held-lock wait/no nested-lock deadlock, and missing-projection refresh without journal mutation. |
| Vacuum CLI policy and denial states | `go test -race ./cmd/fid -run '^(TestMaintenanceVacuumCommandAuthorizesCompactsAndPreservesState|TestMaintenanceVacuumCommandRejectsUnauthorizedRolesBeforeMutation|TestMaintenanceVacuumPolicyShowMatchesEnforcement|TestGovernanceStubCommandsReturnVersionGateBeforeValidation|TestGovernanceStubVersionGateWinsBeforeWorkspaceValidation|TestGovernanceStubFamiliesAppearInRootHelp|TestRequirementSyncRemainsShippedAlongsideGovernanceStubGuards)$' -count=1 -v` | PASS | Reproduced authorized, missing/unknown/denied, policy-parity, guard-precedence, no-mutation, help, and shipped-control cases uncached. |
| Commit and scope verification | `git rev-parse --verify <reported-hash>^{commit}`, `git show`, and `git diff --check` for `f01a308` and `34c6cef` | PASS | Both reported implementation hashes resolve with their reported subjects and task references. Task 6 changes are within its specified maintenance/db/policy/CLI/test/documentation scope; Task 7 contains only the guard tree, root registration, tests, and report state. No dependency change is present. |
| CLI contract and documentation parity | `project-docs/specs/05_DESIGN_PLAN.md:26-27,123-126`; `project-docs/specs/06_DESIGN_PLAN_CLI.md:1485-1510,1928-1947`; `project-docs/specs/02_TRD.md:480-485`; `project-docs/docs/playbook.md:135-147` | PASS | Vacuum stays under the maintenance command surface and emits its confirmation to stderr; its lock, derived-DB-only, and no-journal/activity semantics match. The ten governance leaves, exact v0.2 error, first-before-validation precedence, and requirement-sync exemption match the guarded CLI contract. |
| Phase-end ownership | `project-docs/docs/phases/phase-17/qa/validation-contract.md:70`; `project-docs/docs/phases/phase-17/qa/qa-plan.md:356-370` | PASS | This review records only the in-window evidence for VAL-015/016. Terminal VAL-024 remains assigned to the QA Lead at phase end. |

## Task 6 - `fid maintenance vacuum`

### Blocking

- None.

### Improvement

- [internal/maintenance/vacuum.go:14] Exported `VacuumResult` and `Vacuum` have no Go doc comments, contrary to STRUCT-003 in `08_QUALITY_PLAN.md`. Fix direction: add concise comments that state the locked, derived-projection-only operation and its returned byte sizes.

### Observation

- The authorization check occurs before `journal.WithLock` and projection refresh (`internal/maintenance/vacuum.go:19-29`); the paired CLI test proves the allowed `tl`/`human` path and each missing, unknown, and denied default-role path, while the maintenance test proves preservation and lock contention.

## Task 7 - Governance stub-table guards

### Blocking

- None.

### Improvement

- None.

### Observation

- The leaf command uses arbitrary arguments plus the unknown-flag whitelist before it constructs and renders the version error (`cmd/fid/governance_stubs.go:18-30`). The table-driven test covers all ten leaves, exact errors, future/missing-input precedence, journal/table non-mutation, and the `requirement sync` positive control.

## Cross-cutting observations

- The task reports have complete, committed exact-command evidence; the independent report-lint checks confirm no required Task 6 or Task 7 gate was skipped, narrowed, or left uncommitted.
- No phase-end outcome is inferred from the passing task gates. VAL-024 remains outside this review window and under QA Lead ownership.

## Escalations

- None.

## Final Lines

```text
Review written: project-docs/docs/phases/phase-17/qa/review-3.md
Blocking: 0  Improvement: 1  Escalations: none
VERDICT: CONTINUE
```
