# Phase 20 - Mid-Phase Review 3 (Tasks 7-9)

Reviewer: qa-reviewer
Date: 2026-08-12
Tasks covered: phase-20/task-7, phase-20/task-8, phase-20/task-9

## Verdict

BLOCK

## Scope Reviewed

- Task specs and completion reports: `task-7.md` / `task-7-done.md`, `task-8.md` / `task-8-done.md`, and `task-9.md` / `task-9-done.md`.
- Phase contract and QA artifacts: `contracts.md` §§3, 6, and 7; `qa-plan.md`; and `validation-contract.md`.
- Verified commits: `4950e946a31183cd5386a1a721fa8a2ca2f2988c` (`feat(issue): enforce hierarchy invariant`), `4df4d6690a0562873c570ddaeba2140ad3adef81` (`feat(issue): migrate canonical issue read responses`), and `86d927da88f0fc070709a02ce39a3d68e560e2a6` (`feat(issue): migrate canonical mutation envelopes`). Each resolves with `git rev-parse --verify <hash>^{commit}` and has the required task reference.
- Current context: `AGENTS.md`, active Phase 20 brief/decisions/progress/escalations, and `08_QUALITY_PLAN.md`.

## Summary

- Blocking findings: 3
- Improvement findings: 0
- Escalations raised: none
- Covered validation assertions: VAL-002, VAL-008, VAL-010, VAL-011, VAL-013 through VAL-020, VAL-040, VAL-041, VAL-044, VAL-045

## Validation Contract Coverage

| VAL ID | Claimed by | Status | Evidence |
|---|---|---|---|
| VAL-002, VAL-008 | task-9 | ON TRACK | Canonical stale-update failure and Human mutation channel tests; broader close remains Task 18. |
| VAL-010, VAL-011 | tasks 6-9 | ON TRACK | Task 6 revision/Event fixtures remain in place; Task 9 returns revision/Event fields, with later families owned by Tasks 10-11. |
| VAL-013, VAL-014 | task-8 | ON TRACK | Reproduced closed get/list schema, retained filters/order, and collection-limit tests. |
| VAL-015 | task-8 | ON TRACK | Reproduced exclusive no-ready and sequential-blocker read-state tests; both result forms are asserted. |
| VAL-016 | task-8 | ON TRACK | Retired read flag rejection is covered by the Task 8 command tests. |
| VAL-017 | tasks 8-9 | DEFERRED | Read `event:null` is covered. The phase-wide mutation identity comparison closes at Task 18; Task 9's missing transition identity evidence is recorded below. |
| VAL-018 | task-9 | NOT ON TRACK | The required absent/matching/stale matrix is not evidenced for `issue transition`; see Blocking 1. |
| VAL-019 | task-9 | NOT ON TRACK | Only one canonical unsatisfied-precondition rejection is asserted; see Blocking 2. |
| VAL-020 | task-7 | ON TRACK | Reproduced legal/illegal/missing hierarchy tests and the request/flag absence proof. |
| VAL-040, VAL-041 | task-9 | NOT ON TRACK | Owner-input and attribution negative-state evidence is incomplete; see Blocking 3. |
| VAL-044 | tasks 8-9 | NOT ON TRACK | The Task 9 closed mutation-request matrix is incomplete for the missing owner and transition cases; full leaf closure remains Task 18. |
| VAL-045 | tasks 8-9 | ON TRACK | Task 8 read limit/order and Task 9 lexicographic `changes` ordering are asserted; later collection families close in Tasks 10-13. |

## Gate And Evidence Checks

| Check | Command or Evidence | Result | Notes |
|---|---|---|---|
| Task 7 exact gate coverage | `gofmt -l cmd internal`; `golangci-lint run ./...`; `go test -race -coverprofile=coverage.out ./...`; `go tool cover -func=coverage.out`; `go test -tags integration -race ./cmd/fid/integration/...`; `go build -o /tmp/fid ./cmd/fid` | PASS | All six match Task 7 and QA-plan requirements exactly, with PASS tokens and observed `exit_code: 0` in the report. |
| Task 8 exact gate coverage | Same six commands | PASS | All six required commands and valid output tokens are recorded. The reports predate the forward-effective termination annotation rule; DEC-139 explicitly closes that record without reopening it. |
| Task 9 exact gate coverage | Same six commands | PASS | All six match Task 9 and QA-plan requirements exactly, with PASS tokens and observed `exit_code: 0`. |
| Mechanical report verification | `project-docs/tools/report-lint.sh <report> <task-spec>` for Tasks 7, 8, and 9 | PASS | Each report-lint invocation completed `exit_code: 0` with `PASS report-lint completed with no failures`. |
| Required independent reproduction | `go test -tags integration -race ./cmd/fid/integration/...` | PASS | From `/home/mark/projects/fid`: `ok  \tgithub.com/mbell-nzl/fid/cmd/fid/integration\t8.949s` (terminated `exit_code: 0`). |
| Additional reproduced behavior | Focused Task 7 hierarchy, Task 8 read-state, and transition domain tests | PASS | `go test` runs completed with `ok` for `internal/issues` (0.155s), `cmd/fid` (0.129s), and `internal/issues` (0.119s), each terminated `exit_code: 0`. |

## Per-Task Findings

### Task 7 - Issue hierarchy invariant

#### Blocking

- None.

#### Improvement

- None.

Observation: `internal/issues/create_test.go:38` exercises legal, illegal, and missing-parent states with unchanged journal counts; `cmd/fid/issue_parent_surface_test.go:9` confirms update and transition do not register parent mutation flags.

### Task 8 - Issue read contract

#### Blocking

- None.

#### Improvement

- None.

Observation: `cmd/fid/issue_get_list_test.go:144` and `:194` assert both exclusive ready-result forms and their code-specific nullability; the focused reproduction passed.

### Task 9 - Issue mutation and concurrency contract

#### Blocking

- [cmd/fid/issue_update_test.go:132; cmd/fid/issue_transition_test.go:15] The required three-state optimistic-concurrency matrix is missing for `issue transition`: the only `expected-revision` assertion is a stale **update**, while the transition test file has no `expected-revision` case. The report's matching real-entrypoint case is likewise update-only (`task-9-done.md:30,50`). This leaves absent, matching, and stale transition behavior—including zero append on stale—without deterministic evidence. Fix direction: add table-driven update and transition tests for absent, matching, and stale revisions, asserting exact `state-conflict` details and journal/projection invariance for each rejection.

- [cmd/fid/issue_transition_test.go:15; task-9-done.md:31] The canonical transition evidence asserts only an unsatisfied live-claim precondition. It does not prove the required legal success, unknown target, absent edge, and unsupported/unsatisfied precondition matrix with exact `transition-not-permitted` details and zero append on every rejection. The report substitutes a source reference for missing-edge behavior rather than an observable assertion. It also supplies no canonical transition envelope-to-journal Event-identity comparison. Fix direction: add a canonical JSON transition table covering all four outcomes, state/journal counts, exact details, and success `event` identity against the journal tail.

- [cmd/fid/issue_update_test.go:229; internal/issues/update_test.go:244; task-9-done.md:32-34] Owner-boundary and attribution negative-state evidence is incomplete. The owner-flag test exercises only `issue create`, not `issue update`, despite both retained Issue mutations needing both removed flags rejected. The sole attribution test is update-only and checks a projected activity actor, but contains no forbidden-table/Event-effect assertion proving that no trusted Evidence, Human Outcome, or protected-authority entity/effect was created. Fix direction: cover both owner flags on create and update, assert closed results omit owner changes, and add create/update/transition attribution tests that inspect journal effects and relevant tables for the required absence.

#### Improvement

- None.

## Cross-cutting observations

- The three completion reports record the exact six task/QA-plan gates. Task 8's missing termination annotations are not a finding: the active DEC-139 instruction makes the stronger annotation requirement forward-effective and expressly prevents reopening Task 8.
- The independent integration suite passes, but it cannot substitute for the contract-specific negative matrices above; a broad green command does not prove a missing conditional state.

## Escalations

- None. All findings are bounded to Task 9's required evidence and test coverage; no cross-task decision or contract ambiguity was identified.

## Final Lines

```text
Review written: project-docs/docs/phases/phase-20/qa/review-3.md
Blocking: 3  Improvement: 0  Escalations: none
VERDICT: BLOCK
```
