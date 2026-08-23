# Phase 21 - Mid-Phase Review 1 (tasks 1-3)

Reviewer: qa-reviewer
Date: 2026-08-23
Tasks covered: phase-21/task-1, phase-21/task-2, phase-21/task-3
Review type: INITIAL
Prior review: none

## Summary

- Review outcome: REPAIR
- Must-fix findings: 2
- Improvement findings: 0
- Blocker candidates: 0
- Existing escalation references: none
- Covered validation assertions: VAL-001, VAL-003, VAL-004, VAL-005, VAL-006, VAL-007, VAL-008, VAL-009, VAL-010

## Validation Contract Coverage

| VAL ID | Claimed by | Status | Evidence |
|---|---|---|---|
| VAL-001 | task-1 | NOT ON TRACK | The up/down/up exercise passes, but the inventory assertions only require selected members to be present; they do not prove the required schema is exact. |
| VAL-003 | task-2 | ON TRACK | `TestBuildProjectInitializedEventMatchesV2Golden` byte-compares the Event-v2 bootstrap fixture; unknown top-level-key rejection is covered in `event_v2_test.go`. |
| VAL-004 | task-2 | ON TRACK | Append and replay tests reject aggregate sequence gaps and out-of-order facts. |
| VAL-005 | task-2 | ON TRACK | The idempotency replay test deletes the projection, rebuilds it, and asserts original Event identity with only `replayed:true` changing. |
| VAL-006 | task-2 | ON TRACK | The digest-conflict test asserts typed conflict details and a one-Event journal count. |
| VAL-007 | task-2 | ON TRACK | The race-enabled two-caller test asserts one Event identity, one replay, and one journal Event. |
| VAL-008 | task-2 | ON TRACK | Bootstrap golden and CLI tests assert service Actor, sequence 1, null repository/idempotency, and empty Effects. |
| VAL-009 | task-3 | ON TRACK | The complete literal 12-row, sorted JSON matrix and projection-row query are asserted in `capabilities_test.go`. |
| VAL-010 | task-3 | NOT ON TRACK | The false-row test exercises caller and projection override attempts, but does not execute the seven named deferred operation/surface probes required by QA-09. |

## Gate And Evidence Checks

| Check | Command or evidence | Result | Notes |
|---|---|---|---|
| Task 1 exact gates | Five commands in `task-1-schema-projection.md` | PASS | All five are recorded PASS with tokens and observed `exit_code: 0` in `task-1-done.md`; current two-argument report lint returned `OUTCOME COMPLETE` (terminated `exit_code: 0`). |
| Task 2 exact gates | Six commands in `task-2-event-idempotency.md` | PASS | The report preserves earlier failed package and coverage attempts, documents DEC-162 causal repair, and records later byte-identical PASS rows for every required command; current report lint returned `OUTCOME COMPLETE` (terminated `exit_code: 0`). |
| Task 3 exact gates | Five commands in `task-3-capabilities.md` | PASS | The initial package-gate failure and later byte-identical PASS are recorded with a deviation note; current report lint returned `OUTCOME COMPLETE` (terminated `exit_code: 0`). |
| Reproduced evidence | `go test -race ./internal/workflow ./internal/cli ./cmd/fid/...` | PASS | Re-run in this review: `ok  github.com/mbell-nzl/fid/cmd/fid  (cached)` (terminated `exit_code: 0`). This is Task 3's exact required package gate and QA-08's command. |
| Commit and diff audit | `a094e690815c9940ca926bcad68aedacea663c43`, Task 2 commits `94e752d`, `861e90d`, `c0fd2cd8343826fdcd1a597f8f2f3a671e0e15dc`, and `6ff90152fdcdd373c14a434ec736e557144a1617` | PASS | Reported implementation hashes resolve and have task-scoped Conventional Commit subjects with required `Refs` footers. `git diff --check` is clean for each task's committed range. |

## Task 1 - Operational Schema And Projection

### Must-fix

- [internal/db/operational_kernel_migration_test.go:49] VAL-001 requires the exact additive table/index/check/FK inventory, but the test only rejects missing expected tables and selected schema fragments. It accepts extra tables, indexes, checks, and foreign keys, while the post-up snapshot is compared only to another generated post-up snapshot. Owning task: task-1. Repair boundary: make the migration inventory exhaustive and canonical in `internal/db/operational_kernel_migration_test.go` (and correct 0013 SQL only if that test exposes drift); preserve the existing reversible round trip. Required evidence: failing-then-passing exact-inventory test, `go test -race ./internal/db/...`, all remaining task-1 exact gates, and a passing two-argument report lint result. Re-review scope: inspect 0013 up/down SQL and the exhaustive DDL/index/FK assertions, then rerun `go test -race ./internal/db/...`.

### Improvement

- None.

## Task 2 - Event V2, Idempotency, And Role-Free Bootstrap

### Must-fix

- None.

### Improvement

- None.

## Task 3 - Capability Discovery

### Must-fix

- [cmd/fid/capabilities_test.go:72] VAL-010 and QA-09 require an exact false-row/operation conformance matrix. The loop labels the seven deferred surfaces, then only writes `available = 1` to `system_capabilities`, reruns discovery, and checks the row resets to false; it never probes Activation, Claim lifecycle/list, export create/validate, AC/Concern judgment, historical stats, or Attention absence. Owning task: task-3. Repair boundary: add executable table-driven probes for every named deferred surface that prove its associated capability remains `available:false` and cannot be advertised or authorized by that surface; retain explicitly deferred V1 behavior and do not perform V1 removal or task-4 work. Required evidence: qualifying and non-qualifying capability-state matrix, no forbidden Event/projection side effect for each unavailable probe, `go test -race ./internal/workflow ./internal/cli ./cmd/fid/...`, all remaining task-3 exact gates, and a passing two-argument report lint result. Re-review scope: inspect the new surface probes and rerun the exact Task 3 package gate.

### Improvement

- None.

## Repair Plan

| Finding | Owning task | Affected assertion | Exact repair boundary | Required evidence | Fresh re-review scope |
|---|---|---|---|---|---|
| MF-1 | task-1 | VAL-001 | Exhaustive 0013 schema inventory test in `internal/db/operational_kernel_migration_test.go`; migration SQL only if the exhaustive test finds a mismatch. | Red/green exact inventory evidence; all task-1 exact gates; report lint PASS. | Inspect 0013 SQL and exhaustive assertions; rerun the Task 1 DB gate. |
| MF-2 | task-3 | VAL-010 | Executable named-surface probes in Task 3's `cmd/fid`/`internal/{workflow,cli}` scope; no V1 removal or task-4 implementation. | Binary qualifying/non-qualifying matrix, no-effect assertions, all task-3 exact gates, report lint PASS. | Inspect all seven probes; rerun the Task 3 exact package gate. |

## Cross-cutting observations

- Task 2's earlier protected contradiction was resolved by DEC-162. The final report retains failed evidence, names the causal retained-test convergence, and its final package and repository-wide required commands pass; no active escalation remains.
- Review scope ends at task 3. This review neither assesses nor authorizes task 4.

## Blocker Candidates

- None. Both findings have decision-free, bounded repair ownership in reviewed tasks and do not change accepted product, architecture, public-interface, dependency, or validation-contract authority.

## Existing Escalation References

- None. `tl_escalations_current.md` records no open escalation; ESC-21-001 through ESC-21-004 are resolved and archived.

## Final Lines
Review written: project-docs/docs/phases/phase-21/qa/review-1.md
Must-fix: 2  Improvement: 0  Blocker candidates: 0
VERDICT: REPAIR
