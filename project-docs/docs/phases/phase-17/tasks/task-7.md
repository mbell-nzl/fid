# Phase 17 - Task 7: Add Governance Stub-Table CLI Guards

## Objective

Expose deterministic version-gate errors for commands whose governance tables
exist in the v0.1 schema but whose behavior ships in v0.2. Prevent users and
agents from mistaking Cobra's unknown-command response or flag validation for a
supported governance surface.

## Validation Contract Coverage

- VAL-016: every approved decision/escalation/triage command gates first,
  non-zero and non-mutating, while `requirement sync` remains shipped.

## Binding Constraints

- DEC-052 is binding. Exact guarded commands:
  `decision create/get/list`; `escalation create/resolve/get/list`; `triage
  record/start/verdict`.
- Exact first-line error: `error: commands for <entity> ship in v0.2`.
- The version gate runs before argument, flag, actor, workspace, DB, or journal
  validation, even when callers pass future flags or omit future required flags.
- These are guards, not partial implementations. Add no domain package, matrix
  row, schema write, or placeholder success response.
- `requirement sync` is already shipped and must remain unchanged.

## Scope

### Files To Create Or Modify

- `cmd/fid/governance_stubs.go` - gated command families/subcommands and shared
  typed/version error helper.
- `cmd/fid/governance_stubs_test.go` - exact command matrix, error precedence,
  non-mutation, help visibility, and requirement-sync control.
- `cmd/fid/root.go` - register the three guarded families.

### Files To Read For Context

- `cmd/fid/requirement.go`
- `cmd/fid/root.go`
- `migrations/0002_v01_tables.up.sql`
- `project-docs/specs/01_PRD.md` PRD-ERR-012
- `project-docs/specs/06_DESIGN_PLAN_CLI.md` section 24
- `project-docs/specs/10_IMPLEMENTATION_ROADMAP.md` Phase 17

### Out Of Scope

- Decision, escalation, or triage CRUD/read behavior.
- Requirement get/list/drift.
- Workflow, role, policy-set, comments, or other later-version command guards.
- `maintenance migrate` and upgrade-path guidance.

## Acceptance Criteria

- [ ] Each exact command returns the right entity-specific v0.2 error and a
  non-zero result, including invocations with missing/future flags. (VAL-016)
- [ ] Every guard leaves journal bytes/line count and governance table counts
  unchanged. (VAL-016)
- [ ] Root help exposes the guarded families without claiming CRUD is shipped.
  (VAL-016)
- [ ] `fid requirement sync` still writes successfully in the control case.
  (VAL-016)

## First Failing Test

`TestGovernanceStubCommandsReturnVersionGateBeforeValidation` fails because the
root command currently reports the governance families as unknown.

## Documentation Updates

- None expected. DEC-052 synchronized the PRD, CLI design, and roadmap at
  kickoff; stop if implementation proves that contract wrong.

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

Required-command rules live in `AGENTS.md`. Every command is reported `PASS`,
or `BLOCKED` with reason; the full gate, integration suite, and build are
required before task close.

List any conditional command required by this task and why:

- None.

## Notes For Developer

- Independent of tasks 1-6; DEC-052 supplies the approved public contract.
- Do not implement future flags merely to parse them; the version error must win.
- Completion report: `project-docs/docs/phases/phase-17/tasks/task-7-done.md`.
