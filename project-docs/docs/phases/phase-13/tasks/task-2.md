# Phase 13 - Task 2: Per-Issue Label Assignment (`fid issue label add` / `remove`)

## Objective

Ship single-label assignment that adds or removes **one** label on an issue
without replacing the full set, projecting each change to `activity_log`. Append
the `issue.label.add` and `issue.label.remove` rows to the in-code permissions
matrix. `fid issue update --labels` stays replace-all and is unchanged.

## Requirements

- PRD refs: PRD-ISSUE-023
- TRD refs: TRD-ARCH-010; `03_TRD_SCHEMA.md` sections 9.6, 13
- Design refs: `06_DESIGN_PLAN_CLI.md` `fid issue label add`,
  `fid issue label remove`
- Quality refs: `08_QUALITY_PLAN.md` COV-001/COV-002; CF-007 validate-before-append
- Decisions: DEC-028 (held-lock projection freshness); DEC-031 (Phase 12
  per-command `activity_log` projection pattern); planning interpretations PI-2
  (add/remove semantics) and PI-3 (file placement) in the validation contract

## Validation Contract Coverage

Assertions live in `project-docs/docs/phases/phase-13/qa/validation-contract.md`.

- VAL-001 (completes): the `issue.label.add` and `issue.label.remove` matrix
  rows exist and match `03_TRD_SCHEMA.md` section 13; `fid policy show` reflects
  them.
- VAL-005: `fid issue label add` adds one label in one event without replacing
  the others; nonexistent-label and already-present add reject before append.
- VAL-006: `fid issue label remove` removes one label in one event without
  disturbing the rest; removing an absent label rejects before append.
- VAL-007: each add/remove projects exactly one `activity_log` row recording
  issue + label + direction; `fid label add` projects none.
- VAL-008: `issue.label.add`/`issue.label.remove` authorized only for
  developer/tl/workflow_runner/human; denied/unknown role rejects before append.

## Scope

### Files To Create Or Modify

- `internal/issues/` - add single-label add/remove services (new file, e.g.
  `label.go`, plus tests) that build one `issue_labels` insert or delete effect
  and one `activity_log` effect; reuse the existing label-existence and
  issue-label lookups (`requireExistingLabels`, `issueLabels`).
- `internal/policy/matrix.go` and `internal/policy/matrix_test.go` - append the
  `issue.label.add` and `issue.label.remove` rows from `03_TRD_SCHEMA.md`
  section 13 and update implemented-action completeness tests.
- `cmd/fid/issue_label.go` (new) and `cmd/fid/issue_label_test.go` - wire
  `fid issue label add` / `fid issue label remove` subcommands under the existing
  `fid issue` parent. **Do not append these to `cmd/fid/issue.go` (PI-3).**
- `cmd/fid/issue.go` - registration only (add the `label` subcommand to the
  existing issue command), no other growth.

### Files To Read For Context

- `internal/links/remove.go` for the single-event delete + `activity_log`
  effect pattern (`linkRemoveActivityEffect`, `GenerateID("ACT", ...)`).
- `internal/issues/update_changes.go` (`buildLabelEffects`) for the existing
  `issue_labels` insert/delete effect shapes and the `labels` FieldChange.
- `internal/issues/create.go` (`requireExistingLabels`) and
  `internal/issues/detail.go` (`issueLabels`) for label lookups.
- `internal/issues/update.go` for the activity-effect projection discipline on
  issue field changes.
- `internal/policy/matrix.go`, `internal/policy/authorize.go`.
- `cmd/fid/issue.go` for how the `fid issue` parent registers subcommands and
  uses `renderIssueError`.

### Out Of Scope

- `fid issue update --labels` replace-all behaviour (unchanged).
- Label table CRUD (Task 1) and the closure smoke (Task 3).
- The broad CF-011 `activity_log` reconciliation across other write commands
  (v0.1e); only issue label add/remove project here.
- The full CF-014 split of `issue.go` (v0.1e); this task must not grow `issue.go`
  beyond subcommand registration.

## Acceptance Criteria

- [ ] The in-code matrix contains exactly one `issue.label.add` and one
  `issue.label.remove` row matching section 13 (developer + tl +
  workflow_runner + human); `fid policy show` includes them (VAL-001).
- [ ] `fid issue label add <issue-id> <label> --by <role>` adds exactly that
  label in one journal event and leaves the issue's other labels in place
  (VAL-005).
- [ ] Adding a label that does not exist in the `labels` table, and adding a
  label already present on the issue, each reject before append with a non-zero
  exit, an actionable error, and zero new journal lines (VAL-005).
- [ ] `fid issue label remove <issue-id> <label> --by <role>` removes exactly
  that label in one journal event and leaves the issue's other labels in place;
  removing a label not on the issue rejects before append (VAL-006).
- [ ] Each successful add and remove projects exactly one `activity_log` row
  recording the issue, the label, and the direction (added/removed); removing a
  label does not delete it from the `labels` table (VAL-007).
- [ ] `issue.label.add` / `issue.label.remove` authorize via the matrix:
  developer/tl/workflow_runner/human succeed; a denied role and an unknown role
  reject before append naming role + action (VAL-008).
- [ ] New subcommands live in `cmd/fid/issue_label.go`; `cmd/fid/issue.go` grows
  by registration only (PI-3).

## TDD Plan

- First failing test: an issue with labels `[a, b]`; `fid issue label add c`
  yields `[a, b, c]` (sorted), one journal event, and one `activity_log` add row.
- Additional tests: add already-present label rejects (zero append); add
  nonexistent label rejects (zero append); remove `b` yields `[a]` with one event
  and one `activity_log` remove row; remove absent label rejects (zero append);
  remove leaves the `labels` table row intact; denied-role and unknown-role
  rejection for both actions; matrix completeness for both rows.
- Regression risks: do not reuse the replace-all `--labels` path; the add/remove
  delta must touch exactly one `issue_labels` row; use
  `EnsureProjectionFreshWithLock` inside `WithLock` (no nested lock, DEC-028);
  the `activity_log` row id must be generated with `GenerateID("ACT", ...)` and
  collision-checked like Phase 12 link removal.

## Documentation Updates

- `project-docs/docs/phases/phase-13/decisions.md` - record DEC for PI-2 (add/
  remove reject semantics and the `activity_log` row shape) and PI-3 (file
  placement), plus any section-13 vs in-code mismatch found.
- `06_DESIGN_PLAN_CLI.md` `fid issue label add`/`remove` - reconcile only if the
  implemented flags/output differ from the spec; otherwise record
  verified-no-change in the completion report.
- `project-docs/docs/phases/phase-13/tasks/task-2-done.md` - completion report
  with command results, coverage, the `activity_log` row shape, and commit hash.

## Scenario / Parity / Audit Notes

- Held-out scenario impact: None.
- Cross-surface parity impact: completes policy show vs `Authorize` parity for
  the three Phase 13 matrix rows; `activity_log` projection parity with the
  Phase 12 link-removal pattern (DEC-031).
- Audit or release-readiness impact: None.

## Exact Commands To Run

Run exactly:

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

The full task gate above is required before task close because this task adds
journaled writes with `activity_log` projection and two new subcommands. A
focused test command may be used during the red/green loop, but it does not
substitute for the full gate at close. Every required command must be recorded in
the completion report as `PASS`, or as `BLOCKED` with reason in a blocker report.

Conditional commands:

```text
go build -o /tmp/fid ./cmd/fid
go test -tags integration -race ./...
```

- `go build -o /tmp/fid ./cmd/fid` - required because this task registers new
  CLI subcommands; record PASS.
- `go test -tags integration -race ./...` - run only if an integration fixture is
  touched; otherwise record N/A with reason (the closure smoke is Task 3).

## Notes For Developer

- Use red/green TDD. Do not add a schema migration; the `issue_labels` table
  exists.
- One logical write = one event: an add/remove changes exactly one `issue_labels`
  row plus one `activity_log` row, in a single event.
- Mirror Phase 12 link removal for the `activity_log` effect (op `activity`,
  generated `ACT-...` id, details recording label + direction).
- Depends on Task 1: a label must exist (`fid label add`) before it can be added
  to an issue.
- Completion report commit fields must be filled from `git rev-parse HEAD`.
