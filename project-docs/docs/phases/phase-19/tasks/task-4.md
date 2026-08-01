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
  DEC-091; procedure synchronized to QA-plan §6 by **DEC-098**). The baseline
  is the **Task 3 implementation commit**, pinned below. **Ambient `HEAD` is
  prohibited as the baseline selector** — when Review 1 checked, `HEAD` was
  Task 3's *report* commit, a different commit, and it moves again with every
  later commit. Before editing `internal/issues/create.go`:

  1. **Pin and verify the baseline commit.** Confirm the pinned value
     abbreviates to the `## Commit hash` in the committed `task-3-done.md`
     (`50533e2`, subject `feat(workflow): add define command and journal
     event`), then verify it resolves. Record the export line and the exact
     40-hex output in `task-4-done.md`:

     ```bash
     PHASE19_TASK3_IMPLEMENTATION_COMMIT=50533e2b813bf073cd75f75eb70dd09cd15e9040
     git rev-parse --verify "$PHASE19_TASK3_IMPLEMENTATION_COMMIT^{commit}"
     ```

     If the hash fails to resolve or disagrees with `task-3-done.md`, stop and
     escalate. Never substitute `git rev-parse HEAD`, a branch name, or any
     other moving selector.
  2. **Build the baseline binary in a detached worktree** at that verified
     hash, and record the `sha256sum` line in `task-4-done.md`. Use exactly
     these paths; do not reuse a stale `/tmp/fid` or `/tmp/fid-baseline`:

     ```bash
     git worktree add --detach /tmp/fid-phase19-v01-baseline "$PHASE19_TASK3_IMPLEMENTATION_COMMIT"
     go -C /tmp/fid-phase19-v01-baseline build -o /tmp/fid-phase19-v01-baseline.bin ./cmd/fid
     sha256sum /tmp/fid-phase19-v01-baseline.bin
     ```

  3. **Commit the fixture and the capture script.** The fixture is rooted at
     `cmd/fid/testdata/v01_parity/fixture/`, contains **no**
     `workflow_assignments` event or row, and fixes IDs (`FID-a10001` epic,
     `FID-b20002` task, `FID-c30003` subtask), titles, kinds, parent links,
     **and timestamps in its journal**.
     `cmd/fid/testdata/v01_parity/capture.sh` takes exactly two positional
     arguments — binary path and output root — recreates the workspace from
     the committed fixture on every run, and executes exactly these six
     command shapes:

     ```text
     fid issue get FID-a10001 --format json
     fid issue get FID-b20002 --format json
     fid issue get FID-c30003 --format json
     fid issue list --kind epic --format json
     fid issue list --kind task --format json
     fid issue list --kind subtask --format json
     ```

     Baseline outputs are exactly `issue-get-epic.json`, `issue-get-task.json`,
     `issue-get-subtask.json`, `issue-list-epic.json`, `issue-list-task.json`,
     and `issue-list-subtask.json`, written **directly under**
     `cmd/fid/testdata/v01_parity/` and committed in this task. Candidate
     outputs go to `/tmp/fid-phase19-v01-parity-candidate/` and are **not**
     committed.
  4. **Normalization is a committed no-op.** Commit
     `cmd/fid/testdata/v01_parity/normalization.json` containing exactly
     `{"json_pointers":[],"replacement":null}`. All six files are compared
     raw-byte. **If the fixed fixture cannot hold timestamps stable, this task
     is `BLOCKED` — stop and escalate. Do not widen the allow-list, edit
     `normalization.json`, or hand-write an expectation.**
  5. **Commit `cmd/fid/testdata/v01_parity/SHA256SUMS`** with one `sha256sum`
     line per baseline output. Every line carries the real 64-character
     lowercase digest, exactly two spaces, and the relative filename;
     placeholders are forbidden. `task-4-done.md` copies the six lines
     verbatim. After this task is committed, the immutable baseline identity
     is the Task 4 implementation commit plus those six digests.

  The QA Lead must be able to regenerate the baseline from the recorded hash
  and committed script alone.

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
- `cmd/fid/testdata/v01_parity/` (new, VAL-059) — `capture.sh`, `fixture/`,
  the six committed baseline `.json` outputs, `normalization.json`, and
  `SHA256SUMS`
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
- [ ] With no assignment rows, all six candidate captures are byte-identical
  to the committed baselines captured from the verified Task 3 **implementation**
  commit `50533e2b813bf073cd75f75eb70dd09cd15e9040` (not ambient `HEAD`), with
  the verified hash, baseline-binary `sha256sum`, `capture.sh`, `fixture/`, the
  six outputs, the no-op `normalization.json`, and `SHA256SUMS` all committed
  and reproducible (VAL-059)

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
- `go build -o /tmp/fid ./cmd/fid` — VAL-059's byte-identical `issue get`/
  `list` evidence needs the built candidate binary.

VAL-059 evidence commands, run in this order and all reported in
`task-4-done.md` with verbatim terminal tokens. The first three run **before**
any `internal/issues/create.go` edit; the last three after the candidate build:

```bash
PHASE19_TASK3_IMPLEMENTATION_COMMIT=50533e2b813bf073cd75f75eb70dd09cd15e9040
git rev-parse --verify "$PHASE19_TASK3_IMPLEMENTATION_COMMIT^{commit}"
git worktree add --detach /tmp/fid-phase19-v01-baseline "$PHASE19_TASK3_IMPLEMENTATION_COMMIT"
go -C /tmp/fid-phase19-v01-baseline build -o /tmp/fid-phase19-v01-baseline.bin ./cmd/fid
sha256sum /tmp/fid-phase19-v01-baseline.bin
cmd/fid/testdata/v01_parity/capture.sh /tmp/fid-phase19-v01-baseline.bin cmd/fid/testdata/v01_parity
cmd/fid/testdata/v01_parity/capture.sh /tmp/fid /tmp/fid-phase19-v01-parity-candidate
sha256sum -c cmd/fid/testdata/v01_parity/SHA256SUMS
diff -r cmd/fid/testdata/v01_parity /tmp/fid-phase19-v01-parity-candidate
```

The `diff` must be empty except for the non-output files that live only in the
committed directory (`capture.sh`, `fixture/`, `normalization.json`,
`SHA256SUMS`); every one of the six `.json` outputs must match byte-for-byte.
Remove the scratch worktree (`git worktree remove /tmp/fid-phase19-v01-baseline`)
before finishing so `git status` is clean.

## Notes For Developer

- Depends on Tasks 1-3. This is the heaviest task in the phase; if it grows
  past a single session, stop and report rather than partially landing the
  `internal/issues` change.
- VAL-058/VAL-059 are the v0.1 no-assignment regression guard on the
  most-exercised path in the product. Capture the baseline `issue get`/`list`
  bytes from the pinned pre-change commit **before** editing `create.go`; once
  `create.go` changes, the baseline can no longer be captured from the working
  tree and only the pinned commit can reproduce it.
- `applyDelete` fails the event if the delete matches zero rows. Decide the
  effect shape from state read inside the lock; never emit a speculative
  delete.
- Escalation tripwire: if resolution turns out to need a change to the
  projection or replay path, stop and escalate — replay re-resolving
  assignments would make history non-deterministic.
- Completion report at
  `project-docs/docs/phases/phase-19/tasks/task-4-done.md`.
