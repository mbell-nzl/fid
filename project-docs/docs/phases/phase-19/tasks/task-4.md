# Phase 19 - Task 4: `fid workflow assign` + assignment-aware issue creation

## Objective

Ship PRD-WF-012: `fid workflow assign` binds a defined workflow to an issue
kind or to the project via a `workflow.assigned` journal event, and
`internal/issues.Create` stops hard-coding `defaultWorkflowID(kind)` /
status `"todo"`, resolving both from assignments instead. This is the task
that makes a custom workflow actually take effect.

## Validation Contract Coverage

- VAL-014: a first kind assignment writes exactly one event and exactly one
  kind-scope row.
- VAL-051: reassignment replaces that row and leaves the superseded workflow
  defined and listable.
- VAL-052: kind-scope effect shape — `insert` alone first, ordered
  `delete`+`insert` on reassignment, both replay-safe.
- VAL-053: the same effect-shape rule for project scope.
- VAL-015: `--project` writes exactly one project-scope row with NULL `kind`.
- VAL-054: a kind without a kind-scope assignment resolves to the project
  assignment.
- VAL-055: a kind with a kind-scope assignment wins over the project
  assignment.
- VAL-016: each of the six enumerated assign rejections has its own
  value-naming error; the compatible case succeeds.
- VAL-056: every assign rejection leaves journal bytes and
  `workflow_assignments` rows unchanged.
- VAL-057: a successful assign writes no `activity_log` row (DEC-081).
- VAL-017: creation under an assignment yields the assigned workflow and the
  DEC-077 initial state, for a workflow with an `unstarted` state and one
  without.
- VAL-058: with no assignment rows, all three kinds yield the v0.1 workflow
  and `status='todo'`.
- VAL-059: with no assignment rows, `fid issue get`/`list` output is
  byte-identical to the immutable pre-change baseline (procedure below).

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
- **The v0.1 parity baseline is an immutable, committed artifact** (VAL-059,
  DEC-091). Before editing `internal/issues/create.go`:
  1. Record the **Task 3 completion commit** hash with `git rev-parse HEAD`
     and put it in `task-4-done.md`. That commit is the baseline.
  2. Check it out in a scratch worktree
     (`git worktree add /tmp/fid-baseline <hash>`) and build
     `/tmp/fid-baseline-bin` from it. Do not reuse a stale `/tmp/fid`.
  3. Commit a capture script at
     `cmd/fid/testdata/v01_parity/capture.sh` that builds a fixture workspace
     with **fixed** ids/titles/kinds (one epic, one task, one subtask, no
     assignment rows) and writes `fid issue get <id> --format json` and
     `fid issue list --format json` output to
     `cmd/fid/testdata/v01_parity/<command>-<kind>.json`.
  4. Run it against the baseline binary and **commit the resulting artifacts**
     in this task. The QA Lead must be able to regenerate them from the
     recorded hash and script alone.
  5. Normalization is a committed, enumerated allow-list of volatile fields
     (creation/update timestamps only), applied identically to both sides by a
     committed helper. Prefer none. Anything outside the allow-list is
     compared byte-for-byte.
  If the baseline cannot be captured reproducibly, stop and escalate; do not
  substitute a hand-written expectation.

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
- `cmd/fid/testdata/v01_parity/` (new) — `capture.sh`, the committed baseline
  output artifacts, and the normalization allow-list helper (VAL-059)
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
  the new workflow) and the superseded workflow remains defined and listable
  (VAL-051)
- [ ] A first kind assignment emits an `insert` alone, a kind reassignment
  emits ordered `delete` then `insert`, and replay from empty reproduces one
  row (VAL-052)
- [ ] The same effect-shape rule holds for project scope, with PK
  `{"scope":"project"}` on the delete (VAL-053)
- [ ] `--project` appends one event and creates exactly one project-scope row
  with NULL `kind` (VAL-015)
- [ ] A kind without a kind-scope assignment resolves to the project
  assignment (VAL-054)
- [ ] A kind with a kind-scope assignment resolves to its own workflow even
  when a project assignment exists (VAL-055)
- [ ] Unknown slug, incompatible `applies_to` for `--kind`, non-`any`
  workflow for `--project`, both scopes, neither scope, and a denied role are
  each rejected with a distinct value-naming error, and a compatible assign
  succeeds (VAL-016)
- [ ] Every one of those rejections leaves journal bytes and
  `workflow_assignments` rows unchanged (VAL-056)
- [ ] A successful assign leaves the `activity_log` row set unchanged
  (VAL-057)
- [ ] `fid issue create --kind <k>` under an assignment yields the assigned
  `workflow_id` and the DEC-077 initial state, proven for a workflow with an
  `unstarted` state and one without (VAL-017)
- [ ] With no assignment rows, `workflow_id` and `status` are identical to
  v0.1 for epic, task, and subtask (VAL-058)
- [ ] With no assignment rows, `fid issue get`/`list` output is byte-identical
  to the committed baseline captured from the recorded Task 3 commit, with the
  baseline hash, `capture.sh`, artifacts, and normalization helper all
  committed (VAL-059)

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
- `project-docs/specs/05_DESIGN_PLAN.md` §3.1 — verify the `workflow assign`
  entry added at planning (DEC-092) matches the shipped flag surface; correct
  it in this task if it diverges.

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
