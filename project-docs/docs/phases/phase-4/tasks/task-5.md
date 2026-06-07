# Phase 4 - Task 5: `fid issue list` Filters, Ordering, And Formats

## Objective

Add `fid issue list` over the projected read model with the v0.1b filters,
deterministic ordering, and full read-format dispatch. This completes the Phase
4 authoring/discovery loop: create issues, inspect one issue, and list the
working set in agent-friendly formats.

## Requirements

- PRD refs: PRD-ISSUE-007, PRD-OUT-004
- TRD refs: TRD-ARCH-005, TRD-ARCH-008, TRD-PERF-004
- Design refs: `05_DESIGN_PLAN.md` section 6.2; `06_DESIGN_PLAN_CLI.md`
  `fid issue list`
- Quality refs: `08_QUALITY_PLAN.md` section 5.4, COV-002 `cmd/fid` >= 70%
- Dependencies: Tasks 1-4

## Validation Contract Coverage

Assertions live in
`project-docs/docs/phases/phase-4/qa/validation-contract.md`.

- VAL-008: `fid issue list` filters by kind, status, status category, epic,
  existing label, priority, inbox state, and open concerns.
- VAL-009: `fid issue list` ordering is deterministic and supports JSON, TOON,
  markdown, YAML, and unknown-format rejection.

## Scope

### Files To Create Or Modify

- `internal/issues/issue.go` or new focused list files under `internal/issues/`
- `internal/issues/*test.go`
- `internal/format/markdown.go` and tests if a typed markdown list renderer is
  needed
- `cmd/fid/issue.go`
- `cmd/fid/issue_test.go`
- `project-docs/specs/06_DESIGN_PLAN_CLI.md` if markdown list output needs a
  documented shape
- `project-docs/docs/phases/phase-4/tasks/task-5-completion.md`

### Files To Read For Context

- `internal/issues/issue.go`
- `internal/db/issues_enriched_test.go`
- `internal/format/format.go`
- `internal/format/markdown.go`
- `cmd/fid/issue.go`
- `project-docs/specs/05_DESIGN_PLAN.md` section 6.2
- `project-docs/specs/06_DESIGN_PLAN_CLI.md` `fid issue list`

### Out Of Scope

- `fid ready` and ready ordering semantics (Phase 6).
- `fid issue update` and transition behavior (Phase 5).
- Label CRUD. Label filter tests may seed existing label rows; they must not
  require `fid label add`.

## Acceptance Criteria

- [ ] `internal/issues` exposes a list query over `issues_enriched` with filters
  for kind, exact status, status category, epic, existing label, priority, inbox
  state, and open-concern flag. (VAL-008)
- [ ] Each filter has include and exclude evidence in package-level tests and at
  least representative CLI tests. (VAL-008)
- [ ] `--order sequence|priority|created|updated` is deterministic and has
  stable tie-breaks. (VAL-009)
- [ ] `fid issue list --format json|toon|markdown|yaml` succeeds; structured
  formats carry equivalent fields, markdown has a deterministic readable shape,
  and `--format xml` is rejected with known formats. (VAL-009)
- [ ] Empty result returns an empty list at exit 0. (VAL-008)

## TDD Plan

- First failing test (`internal/issues`): mixed issue fixture filtered by
  `--kind task` includes only tasks and excludes epics/subtasks.
- Additional tests: each filter, each order, empty set, JSON/TOON/YAML parity,
  markdown list golden/structural assertion, unknown format rejection.
- Regression risks: base-table scan instead of `issues_enriched`, unstable
  ordering, list markdown shape undocumented, and label filter semantics
  accidentally depending on label CRUD.

## Documentation Updates

- `project-docs/docs/phases/phase-4/tasks/task-5-completion.md` - record
  VAL-008/009 evidence and exact command outcomes.
- `project-docs/specs/06_DESIGN_PLAN_CLI.md` - document markdown list output if
  the implementation adds a new public markdown shape not already shown.

## Scenario / Parity / Audit Notes

- Held-out scenario impact: None; do not read held-out scenarios.
- Cross-surface parity impact: list rows must preserve field names across JSON,
  TOON, and YAML; markdown may be compact but must not contradict the same
  fields.
- Audit or release-readiness impact: public CLI output shape is a compatibility
  surface.

## Exact Commands To Run

Run exactly:

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

The full task gate above is required before task close. Every required command
must be recorded in the completion report as `PASS`, or as `BLOCKED` with reason
in a blocker report. Missing, silently skipped, or narrower substituted commands
are failed gates.

Conditional commands:

```text
go build -o /tmp/fid ./cmd/fid
```

List any conditional command required by this task and why:

- `go build -o /tmp/fid ./cmd/fid` is required because this task adds the public
  `fid issue list` CLI surface.

## Output Requirements

- Write `project-docs/docs/phases/phase-4/tasks/task-5-completion.md`.
- Completion report must cite VAL-008 and VAL-009.

## Notes For Developer

- Reuse the format engine from Phase 3; do not create a separate list
  serialization path.
- Keep `cmd/fid` thin. Query/filter/order logic belongs in `internal/issues`.
- If markdown list shape needs a product decision beyond a compact documented
  list, stop and escalate before implementation.
