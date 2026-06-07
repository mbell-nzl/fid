# Phase 4 - Task 3: `fid issue create` CLI Surface

## Objective

Expose the expanded authoring service through `fid issue create` with the
v0.1b flags from the CLI spec. This task is the user-visible vertical slice for
creating epics, tasks, and subtasks from a real workspace.

## Requirements

- PRD refs: PRD-ISSUE-001, PRD-ISSUE-002, PRD-ISSUE-003, PRD-ISSUE-004,
  PRD-ISSUE-013, PRD-ISSUE-015, PRD-ISSUE-016, PRD-ISSUE-019,
  PRD-ISSUE-020
- TRD refs: `02_TRD.md` section 4.3, section 4.6; TRD-ARCH-001,
  TRD-ARCH-002, TRD-ARCH-005
- Design refs: `06_DESIGN_PLAN_CLI.md` `fid issue create`;
  `05_DESIGN_PLAN.md` section 5.6
- Quality refs: `08_QUALITY_PLAN.md` section 5.4, COV-002 `cmd/fid` >= 70%
- Dependencies: Task 1 and Task 2

## Validation Contract Coverage

Assertions live in
`project-docs/docs/phases/phase-4/qa/validation-contract.md`.

- VAL-002: CLI creates epic/task/subtask with hierarchy/defaults.
- VAL-003: CLI persists supported authoring fields.
- VAL-004: CLI creates epic metadata for epics only.
- VAL-005: CLI create remains one journal event and rebuildable.
- VAL-001: CLI invalid scheduling/sequence states surface the persistence
  backstop as actionable create errors.

## Scope

### Files To Create Or Modify

- `cmd/fid/issue.go`
- `cmd/fid/issue_test.go`
- `internal/issues/issue.go` only for small service adjustments required by CLI
  wiring
- `project-docs/docs/phases/phase-4/tasks/task-3-completion.md`

### Files To Read For Context

- `cmd/fid/issue.go`
- `cmd/fid/issue_test.go`
- `internal/issues/issue.go`
- `project-docs/specs/06_DESIGN_PLAN_CLI.md` `fid issue create`
- `project-docs/specs/05_DESIGN_PLAN.md` section 5.6

### Out Of Scope

- `fid issue update`, `fid issue transition`, `fid ready`.
- `fid label add/list` and `fid issue label add/remove`.
- Label auto-create semantics. If `--label` is wired in this phase, unknown
  labels must produce an explicit error unless a recorded spec decision says
  otherwise.
- Epic metadata target dates and brief path flags; Phase 4 exposes scheduling.

## Acceptance Criteria

- [ ] `fid issue create` accepts `--kind`, `--title`, `--by`, `--epic`,
  `--parent`, `--priority`, `--sequence`, `--spec-md`, `--est-tokens`,
  `--due-date`, repeatable `--label`, `--inbox`, `--owner-role`,
  `--owner-actor`, and epic-only `--scheduling` as specified. (VAL-002,
  VAL-003, VAL-004)
- [ ] CLI smoke creates an epic, a task under that epic, an unattached task, and
  a subtask; `fid issue get --format json` reads each back with expected
  defaults and fields. (VAL-002, VAL-003)
- [ ] Invalid hierarchy and invalid scheduling/sequence combinations exit
  non-zero with actionable `error:` output naming the violated constraint.
  (VAL-001, VAL-002)
- [ ] Unknown `--label` values are rejected without creating new label rows.
  (VAL-003)
- [ ] One CLI create appends exactly one journal event and survives delete-DB +
  `fid rebuild`. (VAL-005)
- [ ] Write confirmations go to stderr and remain brief per
  `05_DESIGN_PLAN.md` section 5.6. (VAL-002)

## TDD Plan

- First failing test (`cmd/fid`): create an epic with `--scheduling sequential`
  and assert the JSON read-back has `kind=epic`, `workflow_id=epic-default`,
  owner role `tl`, and metadata evidence.
- Additional tests: task under epic, subtask under task, invalid parent flags,
  populated authoring fields, one-event journal count, rebuild parity, role
  errors unchanged.
- Regression risks: Cobra flag aliases drifting from CLI spec, command metadata
  leaking unsanitized paths/arguments, `cmd/fid` absorbing domain validation.

## Documentation Updates

- `project-docs/docs/phases/phase-4/tasks/task-3-completion.md` - record
  VAL-001/002/003/004/005 evidence and exact command outcomes.
- If implementation clarifies `--label` behavior, update the completion report
  and flag spec sync; do not silently pick auto-create semantics.

## Scenario / Parity / Audit Notes

- Held-out scenario impact: None; do not read held-out scenarios.
- Cross-surface parity impact: create/read parity across JSON, TOON, markdown,
  and YAML is exercised by later read/list tasks; do not change output field
  names without spec sync.
- Audit or release-readiness impact: public CLI flags and journal event shape
  are compatibility surfaces.

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

- `go build -o /tmp/fid ./cmd/fid` is required because this task changes the
  public `fid issue create` CLI surface.

## Output Requirements

- Write `project-docs/docs/phases/phase-4/tasks/task-3-completion.md`.
- Completion report must cite VAL-001, VAL-002, VAL-003, VAL-004, and VAL-005.

## Notes For Developer

- Preserve the existing `--by` role validation behavior and error rendering.
- Keep `cmd/fid` thin: parse flags, call `internal/issues`, render success or
  structured errors.
- Do not implement update, transition, ready, or label CRUD in this task.
