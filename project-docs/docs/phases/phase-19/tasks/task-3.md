# Phase 19 - Task 3: `fid workflow define` command + `workflow.defined` event

## Objective

Ship the first user-facing half of PRD-WF-011: register the `fid workflow`
command family, implement `fid workflow define` on top of the Task 2 core,
append a `workflow.defined` journal event, and project it into the workflow
tables. Add the `workflow.define` and `workflow.assign` permission rows so
the matrix change lands once.

## Validation Contract Coverage

- VAL-009: define appends exactly one `workflow.defined` event whose
  projected rows match the declaration and resolve through `workflow.Load`.
- VAL-010: every invalid input is rejected before any append or DB write, with
  the error naming the offending value.
- VAL-011: `workflow.define` / `workflow.assign` permit exactly `tl` and
  `human`, appear in `fid policy show` in all four formats, and a denied
  define appends nothing.
- VAL-012: Unicode arrows in `--transitions` normalize to the same projected
  edges as ASCII `->`.
- VAL-013: the stderr confirmation matches the documented shape and its counts
  equal the projected row counts.

Assertions live in
`project-docs/docs/phases/phase-19/qa/validation-contract.md`.

## Binding Constraints

- **Flag surface is exactly as specified** in `06_DESIGN_PLAN_CLI.md`:

  ```
  fid workflow define --name <slug> --applies-to <kind>
      --states <csv> --categories <csv> --transitions <csv>
      [--require <edge:rule-json>]... --by <role>
  ```

  Do not invent flags. `--name` sets **both** `workflows.id` and
  `workflows.name` (DEC-076); there is no display-name flag.
- **`--applies-to` accepts** `epic`, `task`, `subtask`, `qa_review`,
  `escalation`, `any` (DEC-080). `any` projects as
  `workflows.applies_to IS NULL`.
- **`--transitions` parsing.** Entries are `from->to`. Accept the Unicode
  arrows `→`, `⟶`, and `-->` as equivalent to `->` and normalize before
  handing state names to the core (VAL-012). Trim surrounding whitespace.
- **`--require` parsing.** Format is `<from>-><to>:<rule>[<json>]`, e.g.
  `in_progress->dev_complete:all_acs_verified{"roles":["developer"]}`. Split
  on the **first** `:` after the normalized arrow so JSON colons are safe.
  Absent JSON means NULL `params_json`.
- **Order of operations is a contract** (VAL-010): resolve workspace →
  authorize `workflow.define` → parse flags → validate through the Task 2
  core → only then take the write lock, append, and project. Every rejection
  must happen with zero journal bytes written and zero rows mutated. Follow
  the `internal/issues.Create` pattern: `journal.WithLock`, `PrepareEvent`,
  `DryRunProjectEvent`, `AppendEvent`, `HashEvent`, `ProjectEvent`.
- **One event per define.** Event type `workflow.defined`, `EventSchema` =
  `journal.CurrentEventSchema` (new event types are additive; do not bump the
  event schema). All effects from the Task 2 builder go in that single event
  so definition is atomic on replay.
- **No `activity_log` effect** (DEC-081). The `activity_log` ref-type
  whitelist does not include workflows and is not extended in this phase.
- **Permissions.** Add `{Action: "workflow.define", TL: true, Human: true}`
  and `{Action: "workflow.assign", TL: true, Human: true}` to
  `defaultMatrixRows`, matching `03_TRD_SCHEMA.md` §1. Add **both** rows in
  this task so `fid policy show` output changes once; Task 4 consumes the
  assign row. Update the `fid policy show` tests and any golden output.
  `fid policy show` output is a compatibility surface — additive rows are
  non-breaking, but the change must be visible in the completion report.
- **Confirmation output goes to stderr** and matches the documented block
  exactly (VAL-013):

  ```text
  defined workflow dev-flow-strict for task
    states: 7
    transitions: 8
    requirements: 2
  ```

  Counts are the **projected** row counts, so the transition count includes
  materialized cancel edges. Use the existing `writeIssueLine` /
  `renderIssueError` helpers for consistency with every other command.
- **Register `newWorkflowCommand()` in `cmd/fid/root.go`.** `fid workflow` is
  not a governance stub; it ships now.

## Scope

### Files To Create Or Modify

- `cmd/fid/workflow.go` (new) — `fid workflow` parent + `define` subcommand
- `cmd/fid/workflow_define_test.go` (new)
- `cmd/fid/root.go` — register the command family
- `internal/workflow/define.go` (new) — the write path
  (lock/append/project) wrapping the Task 2 core
- `internal/workflow/define_test.go` (new)
- `internal/policy/matrix.go`, `internal/policy/matrix_test.go` — the two new
  action rows
- `cmd/fid/policy_test.go` — updated `policy show` expectations
- `project-docs/specs/06_DESIGN_PLAN_CLI.md` — `fid workflow define`:
  `--applies-to` vocabulary incl. `escalation`, accepted rule vocabulary and
  the three version-scoped rejection messages, cancel-edge and initial-state
  rules, `--name` setting id and name
- `project-docs/specs/01_PRD.md` — PRD-WF-013 scoping note (DEC-077)

### Files To Read For Context

- `internal/issues/create.go` (the canonical write-path pattern)
- `internal/journal/writer.go`, `internal/journal/lock.go`,
  `internal/journal/project.go`
- `cmd/fid/issue_create.go`, `cmd/fid/issue_helpers.go`, `cmd/fid/policy.go`
- `internal/workflow/authoring.go` (Task 2)
- `project-docs/specs/06_DESIGN_PLAN_CLI.md` `fid workflow define`
- `project-docs/docs/phases/phase-19/decisions.md` (DEC-076..DEC-081)

### Out Of Scope

- `fid workflow assign` (Task 4) and `fid workflow list` (Task 5).
- Any change to `internal/issues.Create` workflow resolution (Task 4).
- Workflow editing, deactivation, or deletion commands (unscheduled).
- Any `activity_log` schema or trigger change (v0.2d).

## Acceptance Criteria

- [ ] `fid workflow define` creates a workflow whose projected rows in all
  four workflow tables match the declaration tuple-for-tuple, including
  declared order as `ordinal` and materialized cancel edges (VAL-009)
- [ ] Exactly one `workflow.defined` event is appended per successful define,
  and the workflow resolves through `workflow.Load` (VAL-009)
- [ ] Every VAL-007 rejection case fails the CLI with an error naming the
  offending value, leaving journal bytes and all four table row counts
  unchanged (VAL-010)
- [ ] `workflow.define` and `workflow.assign` permit exactly `tl` and
  `human`; a denied role is rejected with the standard permission error and
  appends nothing; both rows appear in `fid policy show` in json, toon,
  markdown, and yaml (VAL-011)
- [ ] A definition declared with Unicode arrows projects transitions
  identical to the ASCII `->` form (VAL-012)
- [ ] The stderr confirmation is byte-identical to the documented block and
  its three counts equal the projected row counts (VAL-013)
- [ ] `06_DESIGN_PLAN_CLI.md` and the PRD-WF-013 scoping note are updated in
  this task

## First Failing Test

`TestWorkflowDefineProjectsDeclaredStatesTransitionsAndRules`: run
`fid workflow define` against a temp workspace and assert the four projected
tables match the declaration — fails until the command exists.

## Documentation Updates

- `project-docs/specs/06_DESIGN_PLAN_CLI.md` — `fid workflow define` flags,
  vocabulary, rules, errors, and confirmation output (DEC-077, DEC-078,
  DEC-080).
- `project-docs/specs/01_PRD.md` — PRD-WF-013 gains the DEC-077 scoping note
  that cancel edges exist for workflows declaring a `cancelled`-category
  state, and that the seeded escalation default deliberately declares none.

## Exact Commands To Run

Run exactly:

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
go test -tags integration -race ./...
go build -o /tmp/fid ./cmd/fid
```

Required-command rules live in `AGENTS.md` (Commands): every required command
is reported `PASS`, or `BLOCKED` with reason; missing, silently skipped, or
narrower substituted commands are failed gates. A non-terminating command is
never `PASS`.

Conditional commands required by this task:

- `go test -tags integration -race ./...` — a new journal event type touches
  replay and rebuild paths.
- `go build -o /tmp/fid ./cmd/fid` — VAL-012/VAL-013 need built-CLI evidence,
  and the command tree changed.

## Notes For Developer

- Depends on Tasks 1 and 2.
- The "reject before append" ordering is the single most important property
  in this task. Capture journal byte size before and after each failing case
  in the tests; do not settle for asserting the error text alone.
- `cmd/fid` is a COV-002 package with a 70% threshold; `internal/workflow` is
  COV-001 at 80%. Keep the new write path well covered.
- Escalation tripwire: if the documented confirmation block or flag surface
  cannot be produced as specified, stop and escalate rather than adjusting
  the documented CLI contract on your own.
- Completion report at
  `project-docs/docs/phases/phase-19/tasks/task-3-done.md`.
