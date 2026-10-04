# Phase 24 - Task 14: Human-takes-control journey and phase conformance

## Objective

Prove the Phase 24 acceptance journey deterministically. In a tagged
real-process test, a protected-decision Stop blocks a journey, and it is
resolved to `envelope.complete` using only commands and values that fid
printed. Assemble the rebuild and export/import parity over Phase 24 data,
and the removal-inventory absence proofs.

## Validation Contract Coverage

- VAL-049: the journey can be completed from fid output alone.
- VAL-050: rebuild and export/import parity over Phase 24 data.
- VAL-051: retired surfaces are absent.
- VAL-063 (re-proved): the `PRIVATE_` checks stay removed.
- VAL-052.
- Advances VAL-053, VAL-054, VAL-055, VAL-056 and VAL-057 (Human,
  phase-end). This task's extension of the task 11 fixture generator
  prepares their workspaces. The QA Lead runs it and writes no test code.

Assertions live in `project-docs/docs/phases/phase-24/qa/validation-contract.md`.

## Binding Constraints

- Roadmap Phase 24 journey; contract sections 8-10.
- **The journey** starts from an activated project whose Envelope has an
  unsatisfied protected Decision. A Builder's terminal-success transition is
  refused. Every fid command runs in Human mode. Then, at every step:
  - the command comes from the previous `next:` or `Next:` line, or from
    `--help`;
  - each input file comes from that leaf's `--example`, with every
    placeholder substituted by a value of a permitted provenance class.

  The journey ends at `envelope.complete` success.
- **Provenance classes (DEC-278).** The test records, for every command and
  every substituted value, its class and its source line:
  - `fid-output`: a line printed by an earlier Human-mode fid command in this
    journey (stdout or stderr). JSON output is never a source.
  - `fid-help`: the leaf's `--help` or `--example` text (commands, field
    names, enum members).
  - `caller-authored`: values the acting caller chooses or writes: new record
    IDs (DEC-272), idempotency keys, reasoning text, display labels, and the
    project files (and their paths) the caller writes.
  - `git`: the output of a git command that the leaf's `--help` names for
    that field (`git rev-parse HEAD` and `git status --porcelain` for
    `expected_repository`).
  - `setup`: only the acting Actor IDs (`--actor`) and the project
    directory.

  A value with no permitted class fails the test. If the journey cannot be
  completed within these classes, fid has a product gap. Stop with a
  `BLOCKED` escalation candidate naming the step and the value. Never widen
  the classes or read JSON to pass.
- **Fixture generator extension (DEC-278, R-14).** Extend
  `TestP24HumanFixtures` (task 11) with:
  - the journey's stuck start state (VAL-053);
  - workspaces for every judgment view after task 13 (VAL-056);
  - authoring workspaces where `run close` (with `expected_repository`),
    `decision record-reference` (digest omitted) and a third leaf of the
    Human's choice can each be run validly (VAL-057).

  Keep the task 11 denials and views (VAL-054, VAL-055).
- **Parity** covers field-digest maps, observed repositories, the new `get`
  reads, and manifest `projection_schema` 19.
- **Absence.** No descriptor anywhere contains a member named `worktree_id`,
  `capture_source`, `repository_scope` or `target_repository`, nor a
  `repository` object member (a test over the registry). The per-item
  disclaimer strings and the `PRIVATE_` checks are absent (the `rg`
  commands below). Stop Human renderers emit no raw JSON (the task 11
  test).

## Scope

### Files To Create Or Modify

- `cmd/fid/integration/p24_take_control_test.go` (new)
- `cmd/fid/integration/` parity test for Phase 24 data (new or extended)
- `cmd/fid/integration/p24_human_fixtures_test.go` (extension; the tagged
  gate executes it)
- `cmd/fid/input_descriptor_test.go` (forbidden-member scan)

### Files To Read For Context

- `project-docs/docs/phases/phase-24/public-interface.md` sections 4, 8-10
- `project-docs/docs/phases/phase-24/input-inventory.md` all
- `project-docs/docs/phases/phase-24/decisions.md` DEC-273..277
- `project-docs/specs/10_IMPLEMENTATION_ROADMAP.md` Phase 24 row

### Out Of Scope

- Human checklist content (QA Lead). New product behavior of any kind: a gap
  the journey finds is escalated, not fixed here.

## Acceptance Criteria

- [ ] The journey passes in Human mode. Every command and substituted value
  has a recorded permitted class and source line. (VAL-049)
- [ ] The live and rebuilt projections are semantically equal, export
  validates, and import reproduces the same state with manifest 19.
  (VAL-050)
- [ ] The forbidden-member scan passes and the two absence commands exit 1.
  (VAL-051, VAL-063)
- [ ] `TestP24HumanFixtures` passes with the variable unset, and its manifest
  lists the start state, the judgment views and the three authoring
  workspaces in addition to the task 11 fixtures. (VAL-053..057
  preparation)
- [ ] Wall times inside budget. (VAL-052)

## First Failing Test

The journey's first step asserts that the refused transition's Human stderr contains a `next:` line naming `fid decision record-reference`.

## Documentation Updates

- `project-docs/specs/08_QUALITY_PLAN.md`: the Phase 24 journey as a
  retained real-process scenario.

## Exact Commands To Run

Working directory: repository root `/home/mark/projects/fid`, or the isolated
same-task worktree root the Runner records. No command is a placeholder.
Apply gofmt and goimports to edited Go files before the final gate. Every
command below is required before task close. The tagged integration gate is
mandatory and is never substituted.

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...
rg -n "This (Concern|Deviation) is recorded judgment|This Outcome is authorized Human judgment" cmd/fid -g '!*_test.go'
rg -n "PRIVATE_" cmd/fid/integration
```

The `-timeout 30m` bounds follow DEC-227/259. Record the wall time `go test`
reports for `cmd/fid/integration`, `cmd/fid`, `internal/workflow` and the
slowest other ordinary package. The budgets are tagged integration under
900 s and no ordinary package over 600 s (VAL-052); report any ordinary
package over the 480 s warning. A breach caused by this task's own new journeys is a bounded speed
repair owned by this task; otherwise stop and escalate. Never record a
timed-out gate as `PASS`.

Required-command rules live in `AGENTS.md` (Commands). Every required command
is reported `PASS` or `BLOCKED` with reason, with its own output token and
observed termination. Focused red/green selections never replace these
commands. Report lint is a Runner control, not a Developer gate.

List any conditional command required by this task and why:

- The two `rg` lines are absence proofs for VAL-051.

## Expected Non-Zero Commands

- `rg -n "This (Concern|Deviation) is recorded judgment|This Outcome is authorized Human judgment" cmd/fid -g '!*_test.go'` - expected `exit_code: 1` - per-item disclaimers stay removed
- `rg -n "PRIVATE_" cmd/fid/integration` - expected `exit_code: 1` - the old leak checks stay removed

## Notes For Developer

- Report path: `project-docs/docs/phases/phase-24/tasks/task-14-done.md`.
  Commit with `Refs: phase-24/task-14`.
- Red first: record the first failing test before changing product code.
- Retained-test convergence follows DEC-162/165/166. Record each displaced
  assertion with its introducing commit, phase and task. Never infer
  ownership from a filename or task label. Never delete or weaken an
  assertion to pass a gate.
- Any missing registry member, contract ambiguity, or need to change an
  accepted contract shape or assertion is a `BLOCKED` escalation candidate,
  not a local decision (DEC-195).
