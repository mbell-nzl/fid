# Phase 21 - Mid-Phase Review 2 (tasks 4-6)

Reviewer: qa-reviewer
Date: 2026-08-23
Tasks covered: phase-21/task-4, phase-21/task-5, phase-21/task-6
Review type: INITIAL
Prior review: none
Re-review of: none

## Summary

- Review outcome: REPAIR
- Must-fix findings: 1
- Improvement findings: 0
- Blocker candidates: 0
- Existing escalation references: none
- Covered validation assertions: VAL-011, VAL-012, VAL-013, VAL-014, VAL-015, VAL-016, VAL-017, VAL-018, VAL-019, VAL-020, VAL-021, VAL-022

## Validation Contract Coverage

| VAL ID | Claimed by | Status | Evidence |
|---|---|---|---|
| VAL-011 | task-4 | NOT ON TRACK | The happy lifecycle is exercised, but the required rejected skipped/extra transition states are not asserted. |
| VAL-012 | task-4 | NOT ON TRACK | Post-commit abort rejection is asserted, but neither qualifying pre-commit abort state nor its no-trust result is exercised. |
| VAL-013 | task-4 | ON TRACK | Task spec source-absence control, closed command inventory, false capability row, and post-commit `activation-recovery-unavailable` test cover the unavailable states. |
| VAL-014 | task-4 | ON TRACK | `TestFinalizedActivationIsTheOnlyTrustedActorSource` rejects staged and committed profiles before resolving the exact finalized profile. |
| VAL-015 | task-4 | NOT ON TRACK | The caller-capability rejection is typed, but no provenance variant or zero-Event delta is asserted. |
| VAL-016 | task-5 | ON TRACK | The test rejects a staged source and registers the exact finalized activated Definition. |
| VAL-017 | task-5 | ON TRACK | The one-field policy matrix asserts `executor-policy-unsupported` and zero persisted Registrations for every listed unsupported field. |
| VAL-018 | task-5 | ON TRACK | Same-binding reuse and changed-binding version 2, including immutable version-1 resolve, are asserted. |
| VAL-019 | task-6 | ON TRACK | One primary Issue, revision 1, no related Issues, one existing Registration-bound Obligation, and wider-input negatives are asserted. |
| VAL-020 | task-6 | ON TRACK | The exact `gate-evidence` dual guard round-trips and computed/type/guard widening variants reject. |
| VAL-021 | task-6 | ON TRACK | Create, get, and activate preserve null immutable budget counters; mutable input rejects without an Event. |
| VAL-022 | task-6 | ON TRACK | Routine, revision, cancellation, protected Decision, Stop, and Attention inputs all have typed no-Event rejection rows. |

## Gate And Evidence Checks

All commands below ran from `/home/mark/projects/fid` in this review session. The task reports contain every command from their respective `## Exact Commands To Run` blocks, including observed termination, preserved failed attempts, and deviation notes. QA-10 through QA-21 add reviewer evidence requirements rather than Developer commands, so no unrecorded QA-plan command applies to the three completion reports.

| Check | Command or evidence | Result | Verbatim token and observed termination |
|---|---|---|---|
| Task 4 representative gate | `go test -race ./internal/workflow ./internal/policy ./internal/cli ./cmd/fid/...` | PASS | `ok  github.com/mbell-nzl/fid/cmd/fid (cached)` (terminated `exit_code: 0`) |
| Task 5 representative gate | `go test -race ./internal/evidence ./internal/workflow ./internal/policy ./internal/cli ./cmd/fid/...` | PASS | `ok  github.com/mbell-nzl/fid/cmd/fid (cached)` (terminated `exit_code: 0`) |
| Task 6 representative gate | `go test -race ./internal/workflow ./internal/validation ./internal/cli ./cmd/fid/...` | PASS | `ok  github.com/mbell-nzl/fid/cmd/fid (cached)` (terminated `exit_code: 0`) |
| Task 4 report control | `project-docs/tools/report-lint.sh project-docs/docs/phases/phase-21/tasks/task-4-done.md project-docs/docs/phases/phase-21/tasks/task-4-activation-trust.md` | PASS | `OUTCOME COMPLETE` (terminated `exit_code: 0`) |
| Task 5 report control | `project-docs/tools/report-lint.sh project-docs/docs/phases/phase-21/tasks/task-5-done.md project-docs/docs/phases/phase-21/tasks/task-5-gate-registration.md` | PASS | `OUTCOME COMPLETE` (terminated `exit_code: 0`) |
| Task 6 report control | `project-docs/tools/report-lint.sh project-docs/docs/phases/phase-21/tasks/task-6-done.md project-docs/docs/phases/phase-21/tasks/task-6-envelope-values.md` | PASS | `OUTCOME COMPLETE` (terminated `exit_code: 0`) |
| Reported task commits | `git rev-parse --verify ad9098fac66f2b530a566d87df087fd72fbb584d^{commit}`; `git rev-parse --verify 0c4c04ea0c7d0786e2e6cf88a7573d8a3105aef5^{commit}`; `git rev-parse --verify 56dad8fa0dcb9b79f7df2a31b356f1b2eb8a694e^{commit}` | PASS | each reported full commit ID (terminated `exit_code: 0`) |
| Reviewed implementation diffs | `git diff --check` for each reported implementation commit against its parent | PASS | `exit 0 (no output)` (terminated `exit_code: 0`) |

## Task 4 - Activation Trust Lifecycle

### Must-fix

- [cmd/fid/activation_test.go:60] MF-1: `TestActivationComponentStagesValidatedProfilesAndRejectsCallerAugmentation` asserts only `caller_capabilities` and only a typed/stderr rejection at lines 105-123; it does not assert a zero Event delta or a caller-provenance variant. `TestActivationCommitIsPointOfNoReturn` asserts only the happy commit and post-commit abort at lines 147-185, while the contract and QA-10/QA-11 require every skipped/extra transition to reject and abort to succeed from both `prepared` and `staged` with no active trust source. This leaves the non-qualifying evidence for VAL-011, VAL-012, and VAL-015 incomplete. Owning task: phase-21/task-4. Repair boundary: extend `cmd/fid/activation_test.go` with CLI-level lifecycle/abort and caller-augmentation tables, and update only `task-4-done.md` with the new red/green and exact-gate evidence; do not change accepted lifecycle, error contract, task spec, QA plan, public interface, or ledger. Required evidence: each skipped/extra lifecycle and post-finalization transition rejects with no Event; prepared and staged abort each succeed and leave `ResolveTrustedActor` unavailable; both caller-capability and caller-provenance inputs return the typed error with unchanged journal Event count; the byte-identical Task 4 gate and all Task 4 required commands pass with observed termination; the two-argument report lint returns `OUTCOME COMPLETE`. Re-review scope: a fresh independent reviewer inspects those test rows/report evidence and reruns the exact Task 4 representative command and Task 4 two-argument report-lint control.

### Improvement

- None.

### Evidence

- The implementation enforces the expected status checks in `internal/workflow/activation.go:347`, `internal/workflow/activation.go:384`, and `internal/workflow/activation.go:406`, and the recovery source-absence command is recorded with the declared `exit_code: 1`; these are not substitutes for the missing observable binary-state tests.
- The report preserves its initial focused-gate and lint failures with root causes and byte-identical passing reruns. Its finalized report commit `f95a45a354cbee271219b7a329366d1ecd7b4a4b` and implementation commit `ad9098fac66f2b530a566d87df087fd72fbb584d` are present; the implementation commit resolves to `feat(activation): add trusted activation lifecycle`.

## Task 5 - Activated Gate Registration

### Must-fix

- None.

### Improvement

- None.

### Evidence

- `cmd/fid/gate_test.go:26` rejects registration before finalization, then lines 44-68 assert the finalized Activation-owned Definition and immutable registration fields. The registration implementation resolves actor and Definition from finalized Activation rows before policy validation at `internal/evidence/registration.go:126` and `internal/evidence/registration.go:139`.
- The unsupported-policy table at `cmd/fid/gate_test.go:71` covers argv, shell, sandbox, network, write, environment, limits, exit codes, credentials, redaction, and paths; each row asserts the exact typed error and zero Registration persistence at lines 101-107.
- `cmd/fid/gate_test.go:112` proves identical binding reuse, changed-binding version 2, and version-1 resolution. The report preserves its focused-gate failure/deviation and the report-only recovery needed to restore canonical final PASS notation. Implementation commit `0c4c04ea0c7d0786e2e6cf88a7573d8a3105aef5` resolves to `feat(gate): register activated definitions`.

## Task 6 - One-Revision Envelope Values

### Must-fix

- None.

### Improvement

- None.

### Evidence

- `cmd/fid/envelope_test.go:13` asserts the exact create/get/activate round trip, while its table at lines 43-81 rejects related Issues, a second Obligation, computed status, another type, guard widening, mutable budgets, Routine, revision, cancellation, protected Decision, Stop, and Attention with an unchanged Event count.
- The exact persisted Obligation guards and null counters are checked at `cmd/fid/envelope_test.go:176` and `cmd/fid/envelope_test.go:181`. `internal/workflow/envelope.go:311` validates the closed boundary before constructing any Event effects.
- The report preserves its focused-gate and two lint failures, their causes, and byte-identical passing reruns. Implementation commit `56dad8fa0dcb9b79f7df2a31b356f1b2eb8a694e` resolves to `feat(envelope): add one-revision values`.

## Repair Plan

| Finding | Owning task | Affected assertion | Exact repair boundary | Required evidence | Fresh re-review scope |
|---|---|---|---|---|---|
| MF-1 | phase-21/task-4 | VAL-011, VAL-012, VAL-015 | `cmd/fid/activation_test.go` behavioral tables and truthful `task-4-done.md` evidence only | Both abort qualifiers/no-trust result, skipped/extra transition rejects/no Event, caller capability and provenance rejects/no Event, all Task 4 exact gates, and passing two-argument report lint | Inspect new assertions and report evidence; rerun exact Task 4 representative command and report-lint |

## Cross-cutting observations

- The run ledger records `QA-plan commands beyond the task spec: none required` for tasks 4-6. This is consistent with the plan: QA-10 through QA-21 prescribe reviewer checks, while task-spec command lists remain the Developer gate contracts.
- No out-of-scope files appear in the implementation diffs for tasks 4-6. The only current worktree modification before this review was the runner-owned `project-docs/docs/phases/phase-21/run-ledger.md`.
- Tasks 5 and 6 satisfy the qualifying and non-qualifying assertions called for by their validation rows. Task 4 must meet the same evidence standard before Window 2 advances.

## Blocker Candidates

- None. MF-1 is decision-free, preserves the accepted contract, and fits entirely within phase-21/task-4's existing scope.

## Existing Escalation References

- None. `tl_escalations_current.md` has no open escalation; resolved Phase 21 entries do not govern this bounded repair.

## Final Lines

```text
Review written: project-docs/docs/phases/phase-21/qa/review-2.md
Must-fix: 1  Improvement: 0  Blocker candidates: 0
VERDICT: REPAIR
```
