# Phase 20 - Task 3f: Dedicated integration package and gate migration

## Objective

Stop the integration gate from rerunning the ordinary suite. Classify every
tagged test, move the four real-binary journeys and helpers to
`cmd/fid/integration`, make tagged in-process tests ordinary `cmd/fid` tests,
and migrate the accepted integration command with parity evidence.

## Validation Contract Coverage

- No VAL-NNN applies. Existing scenarios retain their original assertions;
  this task changes test classification and delivery strategy.

## Binding Constraints

- Account for every test name before and after; delete or weaken none.
- Dedicated tests execute the real binary in isolated workspaces. In-process
  tests remain unchanged except removal of their tag.
- Run old and new integration commands before active docs retire the old one.
- Update `AGENTS.md`, Quality Plan, AI Workflow, templates, QA plan, and
  remaining Tasks 4-18. Completed historical specs/reports are immutable.
- Record the test-only package boundary in `02_TRD.md`; add no dependency.

## Scope

### Files To Create Or Modify

- `cmd/fid/integration/` test files and helpers
- Existing integration-tagged `cmd/fid/*_test.go` files
- `AGENTS.md`
- `project-docs/specs/02_TRD.md`
- `project-docs/specs/08_QUALITY_PLAN.md`
- `project-docs/specs/09_AI_WORKFLOW.md`
- `project-docs/templates/phase-runner-prompt.md`
- `project-docs/docs/phases/phase-20/qa/qa-plan.md`
- `project-docs/docs/phases/phase-20/tasks/task-4.md` through `task-18.md`
- `project-docs/docs/phases/phase-20/infra-improvement-pack.md`
- `project-docs/docs/phases/phase-20/tasks/task-3f-done.md`

### Files To Read For Context

- All `cmd/fid` files containing `//go:build integration`
- `cmd/fid/e2e_v01_closure_helpers_test.go`
- DEC-133 and the gate diagnosis

### Out Of Scope

- Building the binary once; Task 3g owns that optimization.
- Rewriting scenarios or product behavior.

## Acceptance Criteria

- [ ] A committed before/after test-name inventory accounts for every tagged
      top-level test exactly once. (no VAL)
- [ ] Four real-process journeys live in `cmd/fid/integration`; tagged
      in-process tests are ordinary `cmd/fid` tests. (no VAL)
- [ ] Default race/coverage executes reclassified ordinary tests and the
      dedicated gate executes only real-process journeys.
- [ ] Old and new integration commands both PASS during migration.
- [ ] Every active gate surface uses the new exact command; workflow-lint
      passes and historical artifacts are unchanged.

## First Failing Test

Add an inventory regression expecting the dedicated package to contain the
four accepted real-process names and the old tagged `cmd/fid` inventory to be
empty; confirm it fails before moving/reclassifying files.

## Documentation Updates

- Canonical gate/spec/template/QA/task surfaces listed above.
- `02_TRD.md` - record the test-only integration adapter boundary.

## Exact Commands To Run

Run exactly:

```text
project-docs/tools/workflow-lint.sh
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
go test -tags integration -race ./...
go test -tags integration -race ./cmd/fid/integration/...
go build -o /tmp/fid ./cmd/fid
```

The repository-wide integration command is transitional parity evidence, not a
narrowing substitute. The dedicated command becomes canonical after both pass.

List any conditional command required by this task and why:

- None.

## Notes For Developer

- Preserve Git history with `git mv` where files move.
- Do not export product internals solely for tests.
- Completion report: `tasks/task-3f-done.md`.
