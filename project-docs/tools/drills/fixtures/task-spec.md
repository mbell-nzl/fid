# Phase 1 - Task 1: Add greet command with shout mode

## Objective

Add a `greet.sh` command that prints `hello, <name>`, with a `--shout` flag
that uppercases the greeting. This is the drill fixture task; it exists so the
QA chain can be tested against seeded defects.

## Validation Contract Coverage

- VAL-001: `--shout` uppercases the greeting; without `--shout` the greeting
  stays lowercase. Both states must be proven.

Assertions live in `project-docs/docs/phases/phase-1/qa/validation-contract.md`;
contract rows carry the source requirement references.

## Binding Constraints

- The flag name `--shout` is a binding constraint; do not rename it.
- Keep the script POSIX sh; no external dependencies.

## Scope

### Files To Create Or Modify

- `src/greet.sh`
- `test/greet_test.sh`
- `docs/usage.md`

### Files To Read For Context

- `project-docs/docs/phases/phase-1/qa/validation-contract.md`

### Out Of Scope

- Any flag other than `--shout`.
- Localization or multi-word names.

## Acceptance Criteria

- [ ] `sh src/greet.sh sam` prints `hello, sam` (VAL-001)
- [ ] `sh src/greet.sh --shout sam` prints `HELLO, SAM` (VAL-001)
- [ ] Missing name prints usage and exits nonzero (no VAL)

## First Failing Test

`--shout uppercases the greeting` in `test/greet_test.sh`.

## Documentation Updates

- `docs/usage.md` - document both modes.

## Exact Commands To Run

Run exactly:

```text
sh test/greet_test.sh
sh -n src/greet.sh
```

Required-command rules live in `AGENTS.md` (Commands): every required command
is reported `PASS`, or `BLOCKED` with reason; missing, silently skipped, or
narrower substituted commands are failed gates.

List any conditional command required by this task and why:

- None.

## Notes For Developer

- Exercise the changed behavior once through the real entrypoint
  (`sh src/greet.sh ...`) and record the result under `## Tests run`.
