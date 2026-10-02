# Phase <N> - Task <M>: <Title>

## Objective

<2-3 sentences: what this task changes and why.>

## Validation Contract Coverage

- VAL-NNN: <one-line restatement of behavior or invariant>

Assertions live in `project-docs/docs/phases/phase-<N>/qa/validation-contract.md`;
contract rows carry the source requirement references. If no assertion
applies, state why explicitly (for example, `documentation-only task`).

## Binding Constraints

The decisions and invariants that bind this implementation - interpretation
decisions, ordering rules, contract semantics. Highest-value section; omit
only for trivial tasks.

- <constraint>

## Scope

### Files To Create Or Modify

- `<path>`

### Files To Read For Context

- `<path>`

### Out Of Scope

- <explicit non-goal>

## Acceptance Criteria

- [ ] <observable/testable criterion> (VAL-NNN or no VAL)

## First Failing Test

<One line naming the first failing test; add regression risks only when
non-obvious.>

## Documentation Updates

- `<doc path>` - <why> (or "None expected")

## Exact Commands To Run

Run exactly:

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -gcflags=modernc.org/...=-d=checkptr=0 -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

Required-command rules live in `AGENTS.md` (Commands): every required command
is reported `PASS`, or `BLOCKED` with reason; missing, silently skipped, or
narrower substituted commands are failed gates. The full gate above is
required before task close unless this task records an exception here.

List any conditional command required by this task and why:

- None.

## Expected Non-Zero Commands

Optional. Include this section only when a required command's pass condition is
a non-zero exit - an `rg`/`grep` absence proof, an `ls` on a path that must not
exist. Declare the expected status here so the Developer never infers one and
the report can record the real observed status:

- `<command>` - expected `exit_code: <N>` - <what its absence proves>

One declaration per line: the exact command in backticks, then
`- expected \`exit_code: <N>\``, then optional prose that is ignored. Never wrap
before the expected field. Every declared command must also appear in
`## Exact Commands To Run`. Omit the section entirely when every required
command must exit 0.

A declaration is a line starting at column 0 with a backtick bullet (`` - ` ``);
every other line here - this prose, an indented continuation of a declaration's
trailing explanation, a bullet that does not open with a backtick - is
explanatory and is not parsed (DEC-144). A backtick bullet that does not match
the grammar is a failure, not prose.

## Notes For Developer

- <constraints, ordering, pitfalls, escalation tripwires; completion-report
  path when it deviates from the `task-M-done.md` convention>
