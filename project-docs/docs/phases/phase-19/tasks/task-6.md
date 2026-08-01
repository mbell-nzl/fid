# Phase 19 - Task 6: Closure — integrity, rebuild parity, exit scenario, spec sync

## Objective

Close the v0.2a slice-2 increment: teach `fid maintenance check` to validate
assignment integrity, prove `fid rebuild` reproduces authored workflows and
assignments deterministically, characterize the DEC-079 permissions boundary,
run the end-to-end exit scenario, re-measure the CF-024 performance controls,
and finish the spec and glossary sync.

## Validation Contract Coverage

- VAL-021: rebuild reproduces custom workflows and assignments
  deterministically and non-destructively; journal bytes never change.
- VAL-022: `fid maintenance check` is clean on a healthy workspace and
  reports a deterministic finding for an incompatible seeded assignment.
- VAL-023: `fid workflow define` still rejects `required_gates_current` and
  `has_completion` at the CLI with version-scoped messages.
- VAL-068: no seeded `workflow_transition_requires` row references either
  reserved name.
- VAL-069: the recorded `rg` inventory finds no production evaluator or query
  for either reserved name.
- VAL-070: the recorded `rg` inventory finds no production read of the frozen
  v0.1 `definitions` map (`FrozenDefinitions(`).
- VAL-024: transitioning into a custom-named state is rejected with the
  documented error; a v0.1-named state under a custom workflow succeeds.
- VAL-025: the v0.2a slice-2 exit scenario runs end to end.
- VAL-026 (Human): a human operator can define, assign, and inspect a custom
  workflow and judges the surface coherent — QA Lead fans this into
  per-observable human-checklist rows.

Assertions live in
`project-docs/docs/phases/phase-19/qa/validation-contract.md`.

## Binding Constraints

- **New structural check, existing pattern.** Add
  `checkWorkflowAssignments` to the `checks` slice in
  `internal/maintenance/check.go`, following the `checkWorkflowStates`
  precedent: read-only, deterministic, one `Finding` per violation with a
  stable `Check` name and the offending entity ID. It flags an assignment
  whose workflow `applies_to` is incompatible with the assigned scope
  (kind-scope row whose workflow is bound to a different kind; project-scope
  row whose workflow is not `any`). FK integrity is already covered by
  `checkForeignKeys`; do not duplicate it.
- **Rebuild parity evidence is whole-projection, not workflow-only**
  (VAL-021). Follow the Phase 18 pattern: normalized per-table snapshots of
  **every** projected table before deletion, after rebuild, and after a
  second rebuild, plus `sqlite_schema` equality, clean
  `PRAGMA foreign_key_check`, and per-file journal byte hashes at every
  boundary. Document normalization for any intentionally volatile field.
- **The exit scenario is scripted, not narrated** (VAL-025). Against a real
  `.fid/` workspace using the built binary: define `task-strict`
  (`--applies-to task`, v0.1 state names **including exactly one
  `cancelled`-category state** as DEC-090 requires for issue-applicable
  workflows, plus a `concerns_addressed`
  requirement on `in_progress->dev_complete`), assign it to kind `task`, list
  it in every format, create a task that picks it up instead of
  `task-default`, raise a concern, observe the transition blocked with a
  `BlockedTransitionError` naming the open concern, address the concern,
  observe the transition succeed — and assert throughout that a control issue
  created under `task-default` before the assignment is unaffected.
- **DEC-079 boundary is characterized, not worked around** (VAL-024). Prove
  both states with the built CLI. Do not add matrix rows for custom state
  names; that grant belongs to Phase 20's `fid policy set`.
- **Seam evidence is reproducible** (VAL-023, VAL-068, VAL-069, VAL-070).
  Record the exact commands and their output:

  ```text
  rg -n --glob '*.go' --glob '!**/*_test.go' \
    '\brequired_gates_current\b|\bhas_completion\b' \
    cmd/fid internal/issues internal/dev internal/journal internal/db \
    internal/maintenance internal/policy
  rg -n --glob '*.go' --glob '!**/*_test.go' \
    --glob '!internal/workflow/definitions_fixture.go' \
    'FrozenDefinitions\(' cmd internal
  ```

  Both are expected to produce **no output and exit 1** — `rg` exits non-zero
  when it finds nothing, and that no-match is exactly what VAL-069/VAL-070
  assert. Record each as `PASS (no match) - exit 1 (no output)`; recording
  `exit 0 (no output)` for a search command is a DEC-086 evidence violation and
  a hard stop (DEC-094, TRG-19-001). Reserved declarations and authoring
  rejections inside `internal/workflow` are expected and are classified
  separately in the report. Record this alongside the CLI rejection evidence
  for `required_gates_current` / `has_completion` and the zero-row seeded
  query.
- **CF-024 measurement obligation.** Run the benchmark gate and record the
  100k-event rebuild against the `<30 s` TRD-PERF-003 control and the
  10k-issue read against the `<500 ms` TRD-PERF-001 control, with the margin
  compared to the Phase 18 baseline (`26.99–27.66 s`). **Do not weaken either
  target.** If the rebuild regresses past the control, stop and escalate.
- **CF-021 rides this task.** Add Go doc comments to the exported
  `maintenance.VacuumResult` and `maintenance.Vacuum`. This task already
  touches `internal/maintenance` and runs the full Go gate.
- **Spec and glossary sync completes here.** `06_DESIGN_PLAN_CLI.md` records
  the DEC-079 limitation and points at `fid policy set` (Phase 20).
  `GLOSSARY.md` gains the Phase 19 terms. This task also **verifies the two
  planning-time spec edits** against what shipped and records the verdict in
  the completion report: `02_TRD.md` §3 (`internal/workflow` responsibilities
  and must-nots, DEC-089) and `05_DESIGN_PLAN.md` §3.1 (the `workflow`
  command family, DEC-092). Correct either here if it diverges.

- **CF-025 — make VAL-059 a standing test, not a ceremony (DEC-101).** Add
  `cmd/fid/v01_parity_test.go` (`//go:build integration`, `package main`) reusing
  the `buildClosureFIDBinary` / `closureRunFID` helpers in
  `e2e_v01_closure_helpers_test.go`: build the current binary, copy
  `cmd/fid/testdata/v01_parity/fixture/` into `t.TempDir()`, run the six
  committed command shapes, and compare each stdout **byte-for-byte** against
  the six committed `.json` goldens, failing with the filename and a diff.
  - The goldens, `fixture/`, `capture.sh`, `SHA256SUMS`, and
    `normalization.json` **stay committed as provenance** — do not delete them.
    What is retired is their status as *required manual procedure*.
  - **Do not regenerate the goldens.** Their provenance is already established
    by reproduction from the pinned Task 3 implementation commit
    `50533e2b813bf073cd75f75eb70dd09cd15e9040` (DEC-099/DEC-101). If the new
    test fails, that is a **real parity regression** — stop and escalate; never
    re-record the goldens to make it pass.
  - No detached worktree, no pinned-commit build, no `mktemp` scratch state, and
    no cleanup command belongs in this test.

## Scope

### Files To Create Or Modify

- `internal/maintenance/check.go`,
  `internal/maintenance/workflow_assignments.go` (new) — the new check
- `internal/maintenance/check_test.go` /
  `internal/maintenance/workflow_assignments_test.go` (new)
- `internal/maintenance/vacuum.go` — CF-021 doc comments
- `cmd/fid/e2e_v02a_workflow_authoring_test.go` (new) — VAL-024, VAL-025
- `cmd/fid/v01_parity_test.go` (new, **CF-025 / DEC-101**) — converts VAL-059
  from a manual ceremony into a standing golden-file test; see Binding
  Constraints
- `internal/journal/` rebuild/replay test files — VAL-021 whole-projection
  parity for authored workflows and assignments
- `project-docs/specs/06_DESIGN_PLAN_CLI.md` — DEC-079 limitation note
- `project-docs/specs/GLOSSARY.md` — Phase 19 terms (see below)
- `project-docs/specs/02_TRD.md` §3 / `project-docs/specs/05_DESIGN_PLAN.md`
  §3.1 — verification of the planning-time DEC-089 / DEC-092 edits

### Files To Read For Context

- `internal/maintenance/check.go` (`checkWorkflowStates` precedent),
  `internal/maintenance/errors.go`
- `internal/journal/phase2_rebuild_test.go`,
  `internal/journal/rebuild_scale_test.go`,
  `internal/journal/rebuild_bench_test.go`
- `cmd/fid/e2e_v01_closure_test.go` (end-to-end scenario conventions)
- `project-docs/docs/phases/phase-18/qa/final-review.md` (Phase 18 rebuild and
  benchmark evidence format, and the CF-024 baseline numbers)
- `project-docs/docs/phases/phase-19/decisions.md` (DEC-079, DEC-082)

### Out Of Scope

- Any new command or flag.
- Any permissions-matrix change (DEC-079; Phase 20).
- Broader markdown characterization backfill (CF-022 remainder).
- Any receipt, gate, or evidence behavior (v0.2b).

## Acceptance Criteria

- [ ] `fid maintenance check` reports no assignment finding on a healthy
  workspace with valid kind and project assignments (VAL-022)
- [ ] A test-seeded incompatible assignment produces a deterministic finding
  and a non-zero `--warn-as-error` exit (VAL-022)
- [ ] Normalized snapshots of every projected table are identical before
  deletion, after rebuild, and after a second rebuild; `sqlite_schema` equal;
  `foreign_key_check` clean; every journal file byte-identical at every
  boundary (VAL-021)
- [ ] `fid workflow define` rejects `required_gates_current` and
  `has_completion` at the CLI with version-scoped messages and no append
  (VAL-023)
- [ ] Zero seeded `workflow_transition_requires` rows reference either
  reserved name (VAL-068)
- [ ] The recorded `rg` inventory shows no production evaluator or query for
  either reserved name (VAL-069)
- [ ] The recorded `rg` inventory shows no production read of
  `FrozenDefinitions(` outside the fixture file itself (VAL-070)
- [ ] Transitioning into a custom-named state fails with the documented
  no-matrix-entry error and appends nothing; a v0.1-named state under a
  custom workflow succeeds (VAL-024)
- [ ] The scripted exit scenario passes every step, including the unaffected
  `task-default` control issue (VAL-025)
- [ ] Benchmark gate run; 100k-event rebuild and 10k-issue read recorded
  against their controls with the margin vs the Phase 18 baseline (CF-024)
- [ ] `maintenance.VacuumResult` and `maintenance.Vacuum` carry Go doc
  comments (CF-021)
- [ ] `06_DESIGN_PLAN_CLI.md` DEC-079 note and `GLOSSARY.md` terms added
- [ ] `02_TRD.md` §3 (DEC-089) and `05_DESIGN_PLAN.md` §3.1 (DEC-092) verified
  against the shipped surface, with the verdict recorded

## First Failing Test

`TestCheckWorkflowAssignmentsFlagsIncompatibleAppliesTo`: seed a kind-scope
assignment whose workflow is bound to another kind and assert the finding —
fails until the check exists.

## Documentation Updates

- `project-docs/specs/06_DESIGN_PLAN_CLI.md` — the DEC-079 limitation:
  transitions into custom state names require a permissions entry, delivered
  by `fid policy set` in Phase 20.
- `project-docs/specs/02_TRD.md` §3 and `project-docs/specs/05_DESIGN_PLAN.md`
  §3.1 — verify the planning-time DEC-089 / DEC-092 edits against the shipped
  package responsibilities and command surface.
- `project-docs/specs/GLOSSARY.md` — add **Custom workflow** (a user-authored
  workflow created by `fid workflow define`, stored as journal-derived rows in
  the projected workflow tables, as opposed to the seeded defaults),
  **Workflow assignment** (a `workflow_assignments` row binding a workflow to
  an issue kind or to the project, with kind scope taking precedence), and
  **Initial state** (the state an issue is created in: the first declared
  `unstarted` state, else the first declared state).

## Exact Commands To Run

Run exactly:

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
go test -tags integration -race ./...
go test -bench . -benchmem ./...
go build -o /tmp/fid ./cmd/fid
```

Required-command rules live in `AGENTS.md` (Commands): every required command
is reported `PASS`, or `BLOCKED` with reason; missing, silently skipped, or
narrower substituted commands are failed gates. A non-terminating command is
never `PASS`; record wall-clock evidence that each command terminated.

Conditional commands required by this task:

- `go test -bench . -benchmem ./...` — the CF-024 measurement obligation.
- `go test -tags integration -race ./...` and `go build -o /tmp/fid ./cmd/fid`
  — the rebuild-parity and end-to-end scenario evidence.

## Notes For Developer

- Depends on Tasks 1-5. This is the phase-closure task; VAL-021, VAL-023, and
  VAL-026 are phase-end assertions.
- Phase 18's ESC-001 was a `workflow.Load` connection deadlock that left four
  required gate commands non-terminating and was recorded `PASS` anyway. When
  running this task's gate, confirm and record that **every** command
  terminated before claiming a result.
- Do not weaken the TRD-PERF-003 or TRD-PERF-001 controls. CF-024 records
  roughly 8-10% headroom on the rebuild; a regression past the control is an
  escalation, not a note.
- Completion report at
  `project-docs/docs/phases/phase-19/tasks/task-6-done.md`.
