# Phase 4 - Task 2: Issue Authoring Domain Service

## Objective

Expand the `internal/issues` authoring service so it can build and project
epic, task, and subtask creation events with v0.1b fields. This is the domain
and journal layer for `fid issue create`; CLI flag wiring lands in Task 3.

## Requirements

- PRD refs: PRD-ISSUE-001, PRD-ISSUE-002, PRD-ISSUE-003, PRD-ISSUE-004,
  PRD-ISSUE-013, PRD-ISSUE-015, PRD-ISSUE-016, PRD-ISSUE-019,
  PRD-ISSUE-020
- TRD refs: `02_TRD.md` section 4.2, section 4.3, section 4.6;
  `03_TRD_SCHEMA.md` section 2.1; TRD-ARCH-001, TRD-ARCH-002, TRD-ARCH-005,
  TRD-ARCH-009
- Design refs: `06_DESIGN_PLAN_CLI.md` `fid issue create`;
  `07_DESIGN_PLAN_FLOWS.md` section 1.1
- Quality refs: `08_QUALITY_PLAN.md` section 5.4, DEC-007 COV-001 carry-in

## Validation Contract Coverage

Assertions live in
`project-docs/docs/phases/phase-4/qa/validation-contract.md`.

- VAL-002: authoring supports epic/task/subtask hierarchy, generated IDs,
  default workflows/statuses, and status categories.
- VAL-003: authoring fields persist through journal/projection/read/rebuild.
- VAL-004: epic metadata is created only for epics.
- VAL-005: create operations remain one logical journal event and rebuildable.

## Scope

### Files To Create Or Modify

- `internal/issues/issue.go` or new focused files under `internal/issues/`
- `internal/issues/issue_test.go` or new package-level tests
- `internal/workflow/` (only if needed for default workflow IDs/status
  categories; do not implement transition edges yet)
- `internal/journal/*test.go` if journal-level rebuild parity needs a focused
  fixture
- `project-docs/docs/phases/phase-4/tasks/task-2-completion.md`

### Files To Read For Context

- `internal/issues/issue.go`
- `internal/journal/event.go`
- `internal/journal/project.go`
- `internal/policy/policy.go`
- `project-docs/specs/02_TRD.md` section 4.6
- `project-docs/specs/03_TRD_SCHEMA.md` section 2.1

### Out Of Scope

- Cobra flag wiring and user-facing CLI error wording (Task 3).
- `fid issue update`, `fid issue transition`, and `fid ready`.
- Label CRUD or label auto-create semantics. `--label`/label assignment support
  must reference existing `labels` rows.
- `epic_metadata.brief_md_path`, `target_start`, and `target_end`; no CLI flag
  exists for those fields in Phase 4.

## Acceptance Criteria

- [ ] `internal/issues.Create` or its replacement accepts an option shape for
  epic/task/subtask creation without hard-coding task-only support. (VAL-002)
- [ ] The service creates valid journal events/effects for epic, task, and
  subtask, with default workflow IDs, initial `todo` status, and owner-role
  defaults (`tl` for epics, `developer` for tasks/subtasks). (VAL-002, VAL-003)
- [ ] Valid hierarchy states succeed and invalid hierarchy states return
  domain errors that identify the violated constraint. (VAL-002)
- [ ] Supported optional fields persist and rebuild correctly: priority,
  sequence, spec markdown path, estimate tokens, due date, existing-label
  assignments, inbox flag, owner role override, and owner actor. (VAL-003,
  VAL-005)
- [ ] Epic creation creates exactly one `epic_metadata` row with scheduling;
  non-epic metadata is rejected by domain validation or DB backstop. (VAL-004)
- [ ] Unknown labels are rejected and no missing label row is auto-created.
  (VAL-003)
- [ ] Each create remains one logical journal event. (VAL-005)

## TDD Plan

- First failing test: package-level create/read round trip for an epic with
  sequential metadata and owner default, red while only task creation is
  supported.
- Additional tests: task under epic, unattached task, subtask under task,
  invalid parent/kind combinations, populated optional fields, omitted-field
  empty/null states, rebuild parity.
- Regression risks: event effects missing preconditions where needed, owner-role
  validation drift, and old task-only error wording surviving.

## Documentation Updates

- `project-docs/docs/phases/phase-4/tasks/task-2-completion.md` - record
  VAL-002/003/004/005 evidence and exact command outcomes.
- Flag any spec ambiguity in the completion report; do not edit product specs
  from this task unless the task proves a durable spec is wrong.

## Scenario / Parity / Audit Notes

- Held-out scenario impact: None; do not read held-out scenarios.
- Cross-surface parity impact: service read-back must remain stable for JSON,
  TOON, markdown, and YAML users once Task 3 exposes it.
- Audit or release-readiness impact: journal event shape is a compatibility
  surface; record any event-type/effect-shape decision in the completion report.

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
go test -tags integration -race ./...
go test -bench . -benchmem ./...
go build -o /tmp/fid ./cmd/fid
```

List any conditional command required by this task and why:

- None. This task is the domain/service layer; Task 3 owns CLI build evidence.

## Output Requirements

- Write `project-docs/docs/phases/phase-4/tasks/task-2-completion.md`.
- Completion report must cite VAL-002, VAL-003, VAL-004, and VAL-005.

## Notes For Developer

- Keep business validation in `internal/issues`, not `cmd/fid`.
- Use `internal/workflow` only for default workflow/status/category constants if
  it reduces duplication and matches the TRD boundary. Do not implement
  transition validation in this task.
- Do not silently invent label auto-create semantics.
- Do not add undocumented epic metadata flags; scheduling is the Phase 4 epic
  metadata surface.
