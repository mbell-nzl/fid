# Phase 21 - Mid-Phase Review 3 (Tasks 7-9)

Reviewer: qa-reviewer
Date: 2026-08-27
Tasks covered: phase-21/task-7, phase-21/task-8, phase-21/task-9
Review type: INITIAL
Prior review: none (review-1b.md and review-2b.md inspected as prior-window context)

## Verdict

REPAIR

The task-9 public stats surface still accepts the retired `--group-by role`
alias. This contradicts VAL-031's exact closed value set and has a bounded
task-9 implementation and test repair. No protected authority is required.

## Scope Reviewed

- Task specs and completion reports for phase-21 tasks 7, 8, and 9.
- Phase-21 QA plan, validation contract, decisions, current context, and the
  Window 1 and 2 re-review artifacts.
- Task implementation and report history through phase HEAD
  `48e25340420966c1581f8d52a592c2772fbee0c3`, including task commits
  `817ad7c50469c51dc12f88a4bc081ecd7a384d7e`,
  `099906daf6795f79d9e669d562a88702dbc33068`, and
  `bc328f9136ec44d31824d1fe2b1c63d1c8b7c921`.

## Summary

- Review outcome: REPAIR
- Must-fix findings: 1
- Improvement findings: 0
- Blocker candidates: 0
- Existing escalation references: ESC-21-007, ESC-21-008 (both resolved by DEC-165 and DEC-166)
- Covered validation assertions: VAL-023, VAL-024, VAL-025, VAL-026, VAL-027, VAL-028, VAL-029, VAL-030, VAL-031

## Validation Contract Coverage

| VAL ID | Claimed by | Status | Evidence |
|---|---|---|---|
| VAL-023 | task 7 | ON TRACK | `cmd/fid/authority_test.go:50` proves the all-match Grant permits; `:54` varies each Actor/action/capability/Issue/Envelope/revision dimension and `:72` denies no active Activation. |
| VAL-024 | task 7 | ON TRACK | `cmd/fid/authority_test.go:77` invokes each protected operation with project-local `developer` attribution and `:91` requires `EnvelopeAuthorityDeniedError`. |
| VAL-025 | task 8 | ON TRACK | `cmd/fid/task8_contract_test.go:9`, `:43`, and `:132` prove attributed create/update/transition replay; `:177` rejects `--by` for all three operations. |
| VAL-026 | task 8 | ON TRACK | `internal/issues/task8_test.go:14` fixes one clock; `:54` requires Event-v2 aggregate sequences 1-3 and `:61` requires each occurrence time to equal the nonfuture fixed time. |
| VAL-027 | task 8 | ON TRACK | `internal/issues/task8_test.go:111` parses `create.go`, checks its 400-line budget, and checks each function's 60-line budget; `create.go` is currently 173 lines. |
| VAL-028 | task 9 | ON TRACK | `cmd/fid/task9_contract_test.go:10` exercises retained mutation/read families; `retained_actor_request_contract_test.go:79` proves missing and malformed Actor requests use current JSON errors with unchanged journal length, while `:173` proves selected-project actor/idempotency operations. |
| VAL-029 | task 9 | ON TRACK | `cmd/fid/policy_test.go:15` and `:45` prove the public policy command is absent with no journal/activity mutation; the report's final package gate passes. |
| VAL-030 | task 9 | ON TRACK | `cmd/fid/capabilities_test.go:102` checks unavailable AC/Concern surfaces, no journal/projection changes, and unchanged capability advertisement; `cmd/fid/maintenance_vacuum_test.go:13` independently checks vacuum and legacy-flag rejection preserve state. |
| VAL-031 | task 9 | NOT ON TRACK | `internal/stats/stats.go:147` accepts legacy `role` and normalizes it to `actor`. Fresh CLI reproduction returned a successful `stats.tokens` JSON response for `--group-by role` (terminated `exit_code: 0`). Current CLI tests cover actor and epic but no `model` success or closed-value rejection. |

## Gate And Evidence Checks

| Check | Command or Evidence | Result | Notes |
|---|---|---|---|
| Task 7 required gates | task-7 report versus task-7 spec | PASS | All five exact commands are recorded PASS with verbatim tokens and observed `exit_code: 0`; task-specific QA-plan additions: none. |
| Task 7 representative rerun | `go test -race ./internal/policy ./internal/workflow ./internal/validation ./internal/cli ./cmd/fid/...` | PASS | Exact final token: `ok  github.com/mbell-nzl/fid/cmd/fid (cached)` (terminated `exit_code: 0`). |
| Task 7 report control | `project-docs/tools/report-lint.sh project-docs/docs/phases/phase-21/tasks/task-7-done.md project-docs/docs/phases/phase-21/tasks/task-7-authority-evaluator.md` | PASS | `OUTCOME COMPLETE` (terminated `exit_code: 0`). Commit hash and subject resolve/match. |
| Task 8 required gates | task-8 report versus task-8 spec | PASS | All five exact commands are recorded PASS with verbatim tokens and observed `exit_code: 0`; preserved failed attempts record real failures and later byte-identical passes. |
| Task 8 representative rerun | `go test -race ./internal/issues ./internal/activity ./internal/journal ./internal/cli ./cmd/fid/...` | PASS | Exact final token: `ok  github.com/mbell-nzl/fid/cmd/fid (cached)` (terminated `exit_code: 0`). |
| Task 8 report control | `project-docs/tools/report-lint.sh project-docs/docs/phases/phase-21/tasks/task-8-done.md project-docs/docs/phases/phase-21/tasks/task-8-issue-actor-convergence.md` | PASS | `OUTCOME COMPLETE` (terminated `exit_code: 0`). Commit hash and subject resolve/match. |
| Task 9 required gates | task-9 report versus task-9 spec | PASS | All five exact commands are recorded PASS with verbatim tokens and observed `exit_code: 0` in the final Attempt 12 block. The previous failures remain recorded with their real tokens and terminations. |
| Task 9 representative rerun | `go test -race ./internal/policy ./internal/workflow ./internal/issues ./internal/labels ./internal/links ./internal/requirements ./internal/activity ./internal/stats ./internal/maintenance ./internal/cli ./cmd/fid/...` | PASS | Exact final token: `ok  github.com/mbell-nzl/fid/cmd/fid (cached)` (terminated `exit_code: 0`). |
| Task 9 report-lint closure | `project-docs/tools/report-lint.sh project-docs/docs/phases/phase-21/tasks/task-9-done.md project-docs/docs/phases/phase-21/tasks/task-9-retained-actor-policy-removal.md` | PASS | `OUTCOME COMPLETE` (terminated `exit_code: 0`). Ledger lines 1336-1432 show the report-only recoveries preserved historical failures and ended at committed report `48e2534`. |
| VAL-031 independent reproduction | `go run ./cmd/fid --project /tmp/fid-qa-stats.9ijxEy --json stats tokens --group-by role` | FAIL | Returned `{"contract_version":"fid-cli.v1","event":null,"ok":true,"operation":"stats.tokens","result":{"items":[]}}` (terminated `exit_code: 0`); this non-qualifying state must be rejected. `--group-by model` separately returned the same successful envelope with `exit_code: 0`; arbitrary `invalid` rejected with `exit_code: 1`. |

## Task 7 - Generic Authority Evaluator

### Must-fix

- None.

### Improvement

- None.

## Task 8 - Issue Actor And Event Convergence

### Must-fix

- None.

### Improvement

- None.

## Task 9 - Retained Inventory Attribution And Judgment Freeze

### Must-fix

- [internal/stats/stats.go:147] The public `--group-by role` alias remains accepted and is converted to `actor`, so Stats does not accept exactly `actor`, `epic`, or `model` as VAL-031 requires. The retained report also overclaims VAL-031 from tests that exercise only actor and epic. Owning task: phase-21/task-9. Repair boundary: remove public `role` acceptance without changing accepted actor/epic/model behavior, then add observable CLI coverage for all three accepted values and for rejected `role`, another invalid value, and retired `--by`; update only the task-9 completion report evidence. Required evidence: the exact task-9 five-command gate, a focused CLI test that asserts successful actor/epic/model envelopes and state-preserving rejected values, and a passing two-argument task-9 report-lint result with token and observed termination. Re-review scope: inspect the bounded stats implementation/test/report diff, verify both qualifying and non-qualifying VAL-031 states, rerun the affected exact task-9 package command and the report-lint control.

### Improvement

- None.

## Repair Plan

| Finding | Owning task | Affected assertion | Exact repair boundary | Required evidence | Fresh re-review scope |
|---|---|---|---|---|---|
| MF-1 | phase-21/task-9 | VAL-031 | `internal/stats/stats.go`, CLI stats tests, and `task-9-done.md` only; remove `role` as a public accepted grouping and preserve actor/epic/model. | Exact five-command task gate; focused actor/epic/model success plus role/invalid/`--by` rejection and no-side-effect proof; two-argument report-lint PASS. | Bounded diff and current code/tests/report; rerun exact task-9 package command and report lint; reproduce both sides of VAL-031. |

## Cross-cutting observations

- Task 7's evaluator boundary remains policy-pure and its workflow adapter resolves the active Actor before evaluation; its matrix gives both allowed and fail-closed coverage.
- Task 8's recovered same-clock test explicitly proves nonfuture Event-v2 ordering, and the task-owned source-budget test guards the reopened creation split.
- Task 9's DEC-166 recovery stayed within the authorized retained-fixture convergence: `git log --follow` traces the renamed request fixture to Phase 20 task 10, while `bc328f9` changes the named Task-9 request rows and later report-only commits `033a586`, `f50673e`, and `48e2534` modify only the completion report. The report-lint closure is mechanically valid, but it cannot substitute for VAL-031's closed-value evidence.

## Blocker Candidates

- None.

## Existing Escalation References

- ESC-21-007 and ESC-21-008 are resolved historical task-9 recovery references; neither covers MF-1 or creates a current protected boundary.

## Final Lines

```text
Review written: project-docs/docs/phases/phase-21/qa/review-3.md
Must-fix: 1  Improvement: 0  Blocker candidates: 0
VERDICT: REPAIR
```
