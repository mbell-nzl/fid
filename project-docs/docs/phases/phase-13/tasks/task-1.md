# Phase 13 - Task 1: Label CRUD (`fid label add` / `fid label list`)

## Objective

Ship the first label-creation surface: a `fid label` command group with `add`
(create one label through a journal event) and `list` (read all labels in four
formats). Append the `label.add` row to the in-code permissions matrix. This
unblocks per-issue label assignment (Task 2), because before this command a
label cannot exist and `requireExistingLabels` on issue create/update can never
be satisfied.

## Requirements

- PRD refs: PRD-ISSUE-015, PRD-ISSUE-023 (prerequisite), PRD-OUT-001..004
- TRD refs: TRD-ARCH-010; `03_TRD_SCHEMA.md` sections 9.6, 13
- Design refs: `06_DESIGN_PLAN_CLI.md` section 9 (`fid label add`, `fid label list`)
- Quality refs: `08_QUALITY_PLAN.md` COV-001/COV-002; CF-007 validate-before-append
- Decisions: DEC-028 (held-lock projection freshness); planning interpretation
  PI-1 (label add semantics) in the validation contract

## Validation Contract Coverage

Assertions live in `project-docs/docs/phases/phase-13/qa/validation-contract.md`.

- VAL-001 (partial): the `label.add` matrix row exists and matches
  `03_TRD_SCHEMA.md` section 13 (developer + tl + workflow_runner + human);
  `fid policy show` reflects it. (The `issue.label.*` rows are added in Task 2.)
- VAL-002: `fid label add` creates one `labels` row in one event; a duplicate
  name rejects before append with zero new journal lines.
- VAL-003: `fid label list` returns name/colour/description across json, toon,
  markdown, and yaml with the same data.
- VAL-004: `label.add` authorized only for developer/tl/workflow_runner/human;
  denied/unknown role rejects before append naming role + action.

## Scope

### Files To Create Or Modify

- `internal/labels/` (new package) - add service, list service, types, errors,
  and tests for `fid label add` / `fid label list`.
- `internal/policy/matrix.go` and `internal/policy/matrix_test.go` - append the
  `label.add` row from `03_TRD_SCHEMA.md` section 13 and update
  implemented-action completeness tests.
- `cmd/fid/label.go` (new) and `cmd/fid/label_test.go` - wire `fid label` with
  `add` and `list` subcommands and flags.
- `cmd/fid/root.go` - register `newLabelCommand()`.

### Files To Read For Context

- `internal/links/add.go` and `internal/links/list.go` for the
  lock -> ensure-fresh-with-lock -> build event -> dry-run -> append -> project
  write pattern and the four-format list pattern.
- `internal/journal/project_effects.go` for generic insert/delete effect
  semantics (no projector change is needed for the `labels` table).
- `internal/issues/create.go` (`requireExistingLabels`) for the existing label
  lookup against the `labels` table.
- `internal/policy/matrix.go`, `internal/policy/authorize.go`.
- `cmd/fid/link.go`, `cmd/fid/policy.go` for new-top-level-command wiring and
  `fidformat.Write` / `renderIssueError` usage.
- `migrations/0002_v01_tables.up.sql` (`labels` DDL).

### Out Of Scope

- `fid issue label add` / `fid issue label remove` (Task 2).
- The end-to-end closure smoke (Task 3).
- Label update/delete/rename or colour-format validation beyond what the
  completion report records.
- Schema migrations. The `labels` table already exists.

## Acceptance Criteria

- [ ] `internal/labels` owns the add + list business rules; `cmd/fid` contains
  only flag wiring and output (VAL-002, VAL-003).
- [ ] The in-code matrix contains exactly one `label.add` row matching section 13
  (developer + tl + workflow_runner + human true, others false);
  `fid policy show` includes it (VAL-001 partial).
- [ ] `fid label add <name> [--colour <hex>] [--description <text>] --by <role>`
  creates one `labels` row with the given fields via one `label.added` journal
  event with a single `labels` insert effect (VAL-002).
- [ ] A second `fid label add` for an existing name rejects before append with a
  non-zero exit, an actionable error naming the existing label, and zero new
  journal lines/rows (VAL-002).
- [ ] `fid label list` returns every label with `name`, `colour`, and
  `description` in json (default), toon, markdown, and yaml, with null
  colour/description preserved (VAL-003).
- [ ] `label.add` authorizes via the matrix: developer/tl/workflow_runner/human
  succeed; a denied role and an unknown role get a permission error before
  append naming role + `label.add` (VAL-004).
- [ ] Label creation projects no `activity_log` row (PI-1).

## TDD Plan

- First failing test: `fid label add core --by developer` creates a `labels` row
  `{name: core}` and emits exactly one journal event.
- Additional tests: add with `--colour`/`--description`; duplicate-name add
  rejects before append (zero new lines); list returns the seeded labels in all
  four formats with matching data and null colour/description; denied-role and
  unknown-role rejection with zero append; matrix completeness for `label.add`.
- Regression risks: write via journal effect + dry-run before append, not direct
  DB mutation; the duplicate path must append nothing; use
  `EnsureProjectionFreshWithLock` inside `WithLock` (no nested lock, DEC-028);
  no `activity_log` row for creation.

## Documentation Updates

- `project-docs/docs/phases/phase-13/decisions.md` - record DEC for PI-1 (label
  add semantics: duplicate rejection, colour/description handling, no activity
  row), plus any section-13 vs in-code mismatch found.
- `06_DESIGN_PLAN_CLI.md` section 9 - reconcile only if the implemented
  flags/output differ from the spec; otherwise record verified-no-change in the
  completion report.
- `project-docs/docs/phases/phase-13/tasks/task-1-done.md` - completion report
  with command results, coverage, label-add semantics summary, and commit hash.

## Scenario / Parity / Audit Notes

- Held-out scenario impact: None.
- Cross-surface parity impact: policy show vs `Authorize` parity for the
  `label.add` row begins here and is finalized with the `issue.label.*` rows in
  Task 2; four-format list parity for labels.
- Audit or release-readiness impact: None.

## Exact Commands To Run

Run exactly:

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

The full task gate above is required before task close because this task
introduces a new domain package, a journaled write, and a new command. A focused
test command may be used during the red/green loop, but it does not substitute
for the full gate at close. Every required command must be recorded in the
completion report as `PASS`, or as `BLOCKED` with reason in a blocker report.

Conditional commands:

```text
go build -o /tmp/fid ./cmd/fid
go test -tags integration -race ./...
```

- `go build -o /tmp/fid ./cmd/fid` - required because this task registers a new
  CLI command; record PASS.
- `go test -tags integration -race ./...` - run only if an integration fixture is
  touched; otherwise record N/A with reason (the closure smoke is Task 3).

## Notes For Developer

- Use red/green TDD. Do not add a schema migration; the `labels` table exists.
- One logical write = one event with one insert effect.
- Keep the add + list logic in `internal/labels` so Task 2 can reuse a
  `LabelExists`-style helper for its precheck.
- Completion report commit fields must be filled from `git rev-parse HEAD`.
