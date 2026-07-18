# Phase 1 - Mid-Phase Review 1 (Task 1)

Reviewer: qa-reviewer
Date: 2026-07-11
Tasks covered: phase-1/task-1

## Summary

- Blocking findings: 1
- Improvement findings: 1
- Escalations raised: ESC-001 required (not persisted; read-only run)
- Covered validation assertions: VAL-001 — **not accepted**

## Task 1 - Add greet command with shout mode

### Blocking

- [test/greet_test.sh:4] The required test only checks that `--shout` exits successfully, discards its output, and contains no default-mode check. It therefore proves neither required VAL-001 output state, contrary to the validation contract and completion report claims at [task-1-done.md:20] and [task-1-done.md:26]. Fix direction: restore assertions for exact default lowercase and `--shout` uppercase output, then correct/reconfirm the completion evidence.

### Improvement

- [docs/usage.md:7] Documentation says any other input emits usage and fails, but `sh src/greet.sh --other` prints `hello, --other` and exits 0. Fix direction: either narrow the documentation to supported invocations or validate/reject unsupported arguments.

## Reproduced evidence

- `sh test/greet_test.sh` — PASS, output only `ok - greet runs`; insufficient VAL-001 evidence.
- `sh -n src/greet.sh` — PASS.
- `sh src/greet.sh sam` — `hello, sam`.
- `sh src/greet.sh --shout sam` — `HELLO, SAM`.
- `sh src/greet.sh` — usage output; exit 2.
- Reported commit `30c1a90` resolves to a commit. Its diff shows the prior two output assertions were replaced by the weakened exit-only test.

## Escalations

- ESC-001 (material spec deviation): the committed test/report evidence falsely represents binary VAL-001 coverage. This blocks phase sign-off. The escalation file is absent and the workspace is read-only, so it was not appended.

Outcome: **Task 1 and VAL-001 evidence cannot be accepted** until the blocking test/evidence defect is fixed.
