# Phase 19 - Task 4: `fid workflow assign` + assignment-aware issue creation

## Objective

Ship PRD-WF-012: `fid workflow assign` binds a defined workflow to an issue
kind or to the project via a `workflow.assigned` journal event, and
`internal/issues.Create` stops hard-coding `defaultWorkflowID(kind)` /
status `"todo"`, resolving both from assignments instead. This is the task
that makes a custom workflow actually take effect.

## Validation Contract Coverage

- VAL-014: kind assignment writes one event and exactly one kind-scope row;
  reassignment replaces it and leaves the superseded workflow defined.
- VAL-015: project assignment creates exactly one project-scope row, and
  kind-scope assignment takes precedence over it (both directions).
- VAL-016: every invalid assign is rejected before append with a distinct
  error; the valid case succeeds.
- VAL-017: issue creation picks up the assigned workflow and its resolved
  initial state, and is byte-identical to v0.1 when no assignment exists.

Assertions live in
`project-docs/docs/phases/phase-19/qa/validation-contract.md`.

## Binding Constraints

- **Flag surface is exactly as specified** in `06_DESIGN_PLAN_CLI.md`:

  ```
  fid workflow assign --name <slug> [--kind <kind>] [--project] --by <role>
  ```

  Exactly one of `--kind` or `--project` is required; both or neither is an
  error. The `workflow.assign` matrix row already exists from Task 3.
- **Compatibility rules (DEC-080).** `--kind <k>` requires the workflow's
  `applies_to` to equal `<k>` or be `any` (NULL). `--project` requires
  `applies_to` to be `any`, because a project default must be legal for every
  kind. Reuse the Task 2 assignment validator; do not re-implement the rules.
- **Reassignment is last-write-wins, expressed literally** (DEC-076). When a
  row already exists for the target scope, the `workflow.assigned` event
  carries an ordered `delete` effect (PK `{"scope":"kind","kind":"<k>"}` or
  `{"scope":"project"}`) followed by an `insert`. When no row exists, it
  carries the `insert` alone. Read current state inside the write lock and
  emit the correct effect set; `applyDelete` requires exactly one affected
  row, so an unconditional delete would break replay on a first assignment.
- **Reject before append**, same ordering contract as Task 3: authorize →
  parse → validate → lock → append → project. Every rejection leaves journal
  bytes and `workflow_assignments` row counts unchanged (VAL-016).
- **No `activity_log` effect** (DEC-081).
- **Resolution at issue creation (DEC-076).** Replace `defaultWorkflowID(kind)`
  with: kind-scope assignment for the issue's kind → project-scope assignment
  → the built-in kind default (`epic-default` / `task-default` /
  `subtask-default`). Resolve inside the existing write lock, from the
  projection, before building the insert effect.
- **Initial state comes from the workflow, not a literal** (DEC-077). Replace
  the hard-coded `"todo"` with the Task 2 initial-state helper applied to the
  resolved workflow: first declared `unstarted` state, else first declared
  state. For all three seeded defaults this resolves to `todo`, so the
  no-assignment path must stay byte-identical to v0.1 (VAL-017).
- **The resolved workflow and status are written into the event effect**, as
  they already are today. Replay must not re-resolve assignments: an issue
  created before an assignment change keeps the workflow it was created with.
  Do not add resolution logic to the projection path.
- **Do not touch the transition path.** Transitioning into custom-named states
  is out of scope and is a documented Phase 19 boundary (DEC-079).

## Scope

### Files To Create Or Modify

- `cmd/fid/workflow.go` — add the `assign` subcommand
- `cmd/fid/workflow_assign_test.go` (new)
- `internal/workflow/assign.go` (new) — the write path
- `internal/workflow/assign_test.go` (new)
- `internal/workflow/resolve.go` (new) — assignment resolution helper used by
  `internal/issues`
- `internal/issues/create.go` — replace `defaultWorkflowID` and the `"todo"`
  literal with resolved values
- `internal/issues/create_test.go` / `internal/issues/issue_test.go` —
  assigned and unassigned creation coverage
- `project-docs/specs/06_DESIGN_PLAN_CLI.md` — `fid workflow assign`
  compatibility rules and errors
- `project-docs/specs/03_TRD_SCHEMA.md` §2.1 — note that kind-appropriate
  workflow assignment at creation now resolves through `workflow_assignments`

### Files To Read For Context

- `internal/issues/create.go` (`defaultWorkflowID`, the `"todo"` literal, the
  effect builder)
- `internal/workflow/authoring.go`, `internal/workflow/assignment.go` (Task 2)
- `internal/workflow/define.go` (Task 3 write-path pattern)
- `internal/journal/project_effects.go` (`applyDelete` single-row rule)
- `migrations/0009_v02a_workflow_assignments.up.sql` (Task 1)
- `project-docs/docs/phases/phase-19/decisions.md` (DEC-076, DEC-077, DEC-080)

### Out Of Scope

- `fid workflow list` (Task 5).
- Unassign / clear-assignment commands (unscheduled).
- Any change to `fid issue transition` or the permissions matrix (DEC-079;
  Phase 20).
- Any change to `defaultIssueOwnerRole` or other creation defaults.

## Acceptance Criteria

- [ ] `fid workflow assign --kind <k>` appends exactly one
  `workflow.assigned` event and leaves exactly one kind-scope row (VAL-014)
- [ ] Re-assigning the same kind replaces that row (still exactly one, naming
  the new workflow) and the superseded workflow remains defined (VAL-014)
- [ ] `--project` creates exactly one project-scope row; a kind without a
  kind-scope assignment resolves to it, and a kind with one resolves to its
  own workflow (VAL-015)
- [ ] Unknown slug, incompatible `applies_to` for `--kind`, non-`any`
  workflow for `--project`, both scopes, neither scope, and a denied role are
  each rejected with a distinct error, leaving journal bytes and
  `workflow_assignments` row counts unchanged (VAL-016)
- [ ] `fid issue create --kind <k>` under an assignment yields the assigned
  `workflow_id` and the DEC-077 initial state, proven for a workflow with an
  `unstarted` state and one without (VAL-017)
- [ ] With no assignment rows, `workflow_id` and `status` are identical to
  v0.1 for epic, task, and subtask, and `fid issue get`/`list` output is
  byte-identical to the pre-change build (VAL-017)
- [ ] A first assignment emits an `insert` alone; a reassignment emits
  `delete` then `insert`, and both replay cleanly

## First Failing Test

`TestIssueCreateUsesAssignedWorkflowAndItsInitialState`: assign a custom
workflow to kind `task`, create a task, and assert `workflow_id` and `status`
— fails until resolution replaces the hard-coded default.

## Documentation Updates

- `project-docs/specs/06_DESIGN_PLAN_CLI.md` — `fid workflow assign`
  compatibility rules, error list, and reassignment semantics (DEC-076,
  DEC-080).
- `project-docs/specs/03_TRD_SCHEMA.md` §2.1 — the note that kind-appropriate
  defaults are assigned in application code is updated to name
  `workflow_assignments` as the resolution source.

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

- `go test -tags integration -race ./...` — issue creation and a new event
  type both cross replay/rebuild paths.
- `go build -o /tmp/fid ./cmd/fid` — VAL-017's byte-identical `issue get`/
  `list` evidence needs the built binary.

## Notes For Developer

- Depends on Tasks 1-3. This is the heaviest task in the phase; if it grows
  past a single session, stop and report rather than partially landing the
  `internal/issues` change.
- The v0.1 no-assignment parity half of VAL-017 is a regression guard on the
  most-exercised path in the product. Capture the `issue get`/`list` bytes
  from the pre-change build **before** editing `create.go`.
- `applyDelete` fails the event if the delete matches zero rows. Decide the
  effect shape from state read inside the lock; never emit a speculative
  delete.
- Escalation tripwire: if resolution turns out to need a change to the
  projection or replay path, stop and escalate — replay re-resolving
  assignments would make history non-deterministic.
- Completion report at
  `project-docs/docs/phases/phase-19/tasks/task-4-done.md`.
