# Phase 1 - Task 1: Add greet command with shout mode

## Objective

Add a `greet.sh` command that prints `hello, <name>`, with a `--shout` flag
that uppercases the greeting. This is the drill fixture task; it exists so the
QA chain can be tested against seeded defects.

## Requirements

- PRD refs: fixture-PRD §1 (greeting output)
- TRD refs: fixture-TRD §1 (POSIX sh, no dependencies)
- Design refs: N/A
- Quality refs: fixture quality plan §1

## Validation Contract Coverage

List the `VAL-NNN` assertions this task advances. Assertions live in
`project-docs/docs/phases/phase-1/qa/validation-contract.md`.

- VAL-001: `--shout` uppercases the greeting; without `--shout` the greeting
  stays lowercase. Both states must be proven.

## Scope

### Files To Create Or Modify

- `src/greet.sh`
- `test/greet_test.sh`
- `docs/usage.md`

### Files To Read For Context

- `project-docs/docs/phases/phase-1/qa/validation-contract.md`

### Out Of Scope

- Any flag other than `--shout`. The flag name `--shout` is a binding
  constraint; do not rename it.
- Localization or multi-word names.

## Acceptance Criteria

- [ ] `sh src/greet.sh sam` prints `hello, sam` (VAL-001)
- [ ] `sh src/greet.sh --shout sam` prints `HELLO, SAM` (VAL-001)
- [ ] Missing name prints usage and exits nonzero (no VAL)

## TDD Plan

- First failing test: `--shout uppercases the greeting` in
  `test/greet_test.sh`.
- Additional tests: default lowercase state; usage error.
- Regression risks: none; new files only.

## Documentation Updates

- `docs/usage.md` - document both modes.

## Scenario / Parity / Audit Notes

- Held-out scenario impact: None
- Cross-surface parity impact: None
- Audit or release-readiness impact: None

## Exact Commands To Run

Run exactly:

```text
sh test/greet_test.sh
sh -n src/greet.sh
```

The full task gate above is required before task close unless this task
explicitly states otherwise here. Every required command must be recorded in
the completion report as `PASS`, or as `BLOCKED` with reason in a blocker
report. Missing, silently skipped, or narrower substituted commands are failed
gates.

List any conditional command required by this task and why:

- None.

## Notes For Developer

- Keep the script POSIX sh; no external dependencies.
- Exercise the changed behavior once through the real entrypoint
  (`sh src/greet.sh ...`) and record the result under `## Tests run`.
