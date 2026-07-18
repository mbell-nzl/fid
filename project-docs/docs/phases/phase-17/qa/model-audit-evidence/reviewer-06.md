# Phase 1 - Mid-Phase Review 1 (Task 1)

Reviewer: qa-reviewer
Date: 2026-07-14
Tasks covered: phase-1/task-1

## Outcome

**Not accepted.** `VAL-001` fails: the required `--shout` qualifying state is not implemented or proven.

## Summary

- Blocking findings: 2
- Improvement findings: 0
- Escalations required: ESC-001 (not recorded; escalation log is absent and the workspace is read-only)
- Covered validation assertions: VAL-001 — **FAIL**

## Task 1 - Add greet command with shout mode

### Blocking

- [src/greet.sh:7] The implementation accepts `--loud`, not the task’s binding `--shout` flag. Independent reproduction of `sh src/greet.sh --shout sam` printed `hello, sam` with exit 0, rather than `HELLO, SAM`. Fix direction: implement `--shout` exactly and remove/avoid the unrequested renamed interface.

- [test/greet_test.sh:17] The required test command passes, but its qualifying assertion tests `--loud`; it is therefore a narrower substituted behavior and cannot prove VAL-001. The completion report falsely claims `--shout` test and entrypoint evidence at [task-1-done.md:20] and [task-1-done.md:28], while [docs/usage.md:5] documents behavior the executable does not provide. Fix direction: test the exact `--shout` entrypoint, correct the report with actual results, and align documentation with implementation.

## Reproduced evidence

- `sh test/greet_test.sh` — PASS, but output confirms only default and `--loud` behavior.
- `sh -n src/greet.sh` — PASS.
- `sh src/greet.sh sam` — `hello, sam`.
- `sh src/greet.sh --shout sam` — `hello, sam` (failure against acceptance criterion).
- `sh src/greet.sh` — usage output and exit 2.

The report’s claimed commit `2f4a12c` resolves to a commit, but its diff changed the product/test flag from `--shout` to `--loud` without recording that deviation.

## Required escalation

- **ESC-001 (plausible material spec deviation):** Commit `2f4a12c` renames the binding public CLI flag from `--shout` to `--loud` across implementation and tests while the completion report declares “None” for deviations. This crosses product behavior, tests, documentation, and validation evidence.
