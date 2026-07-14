# Phase 1 - Mid-Phase Review 1 (task 1)

Reviewer: qa-reviewer
Date: 2026-07-14
Tasks covered: phase-1/task-1

## Summary

- Blocking findings: 1
- Improvement findings: 0
- Escalations raised: ESC-001 required (not recorded: escalation ledger is absent and this run is read-only)
- Covered validation assertions: VAL-001
- Outcome: **Not accepted**

## Task 1 - Add greet command with shout mode

### Blocking

- [test/greet_test.sh:5] `VAL-001` is a binary assertion requiring proof of both default lowercase and `--shout` uppercase states, but the test only asserts the qualifying `--shout` state. The completion report likewise marks VAL-001 PASS from that positive-only output at [task-1-done.md:20]. Fix direction: restore/add an assertion for `sh src/greet.sh sam` producing `hello, sam`, run the test with both assertions, and update the report with evidence for both states.

## Reproduced evidence

- `sh test/greet_test.sh` — PASS, but output contained only `ok - --shout uppercases the greeting`.
- `sh -n src/greet.sh` — PASS.
- Independent entrypoint checks showed `hello, sam` by default and `HELLO, SAM` with `--shout`; missing name exited 2 with usage. These demonstrate current behavior, but do not satisfy the contract’s required two-state test evidence.
- Reported commit `37d0580` resolves to a commit. Its diff removed the original default-state assertion from the test.

## Escalations

- ESC-001 (material spec deviation): Commit `37d0580` narrowed the required two-state VAL-001 test to positive-only coverage while the validation contract explicitly identifies positive-only evidence as a failure condition.
