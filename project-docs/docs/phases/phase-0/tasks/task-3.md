# Phase 0 - Task 3: Spec Audit And Phase Readiness Docs

## Objective

Audit the spec and workflow surface after the scaffold lands, then document
whether Phase 1 can be planned without unresolved Phase 0 stack, command, or
workflow blockers.

## Requirements

- PRD refs: PRD §7 open product questions, PRD-OPS-009, PRD-NFR-005,
  PRD-NFR-006
- TRD refs: TRD §1, TRD §9, TRD §11, TRD-SEC-001, TRD-ARCH-013
- Design refs: Design Plan §§1-3
- Quality refs: Quality Plan §§5.4-5.7, AI Workflow §§2.2, 5.1, 7,
  Roadmap Phase 0

## Validation Contract Coverage

- VAL-004: Dependencies remain accepted and documented.
- VAL-007: Phase planning artifacts map assertions to task specs and brief.
- VAL-008: Phase 1 has no unresolved Phase 0 stack choice or command placeholder.
- VAL-009: Final scaffold contains no daemon, listener, remote sync, or MCP adapter.

## Scope

### Files To Create Or Modify

- `project-docs/docs/phases/phase-0/spec-audit.md`
- `project-docs/docs/phases/phase-0/tasks/task-3-completion.md`
- `project-docs/docs/current/progress_current.md` if task execution changes
  current status
- Canonical specs only if the audit finds a real spec defect that can be fixed
  without a product, architecture, data-model, or public-interface decision

### Files To Read For Context

- `AGENTS.md`
- `project-docs/specs/01_PRD.md` §7
- `project-docs/specs/02_TRD.md` §§1, 9, 11
- `project-docs/specs/05_DESIGN_PLAN.md` §§1-3
- `project-docs/specs/08_QUALITY_PLAN.md` §§5.4-5.7
- `project-docs/specs/09_AI_WORKFLOW.md` §§2.2, 5.1, 7
- `project-docs/specs/10_IMPLEMENTATION_ROADMAP.md` Phase 0
- `project-docs/docs/phases/phase-0/qa/validation-contract.md`
- `project-docs/docs/phases/phase-0/tasks/task-1-completion.md`
- `project-docs/docs/phases/phase-0/tasks/task-2-completion.md`

### Out Of Scope

- Resolving future product questions that PRD §7 explicitly defers.
- Accepting future Proposed dependencies for v0.3 or v0.4.
- Weakening validation assertions, quality gates, or dependency policy.
- Reading held-out scenarios or evaluator prompts.

## Acceptance Criteria

- [ ] `spec-audit.md` classifies every in-scope `<TBD>`, Proposed, placeholder, or open-question match as blocking or non-blocking with source references. (VAL-008)
- [ ] Any blocking Phase 0 stack, command, or workflow issue produces a blocker report instead of silent spec adaptation. (VAL-008)
- [ ] Future PRD §7 questions and future Proposed TRD dependencies are explicitly recorded as out of scope unless the audit proves they block Phase 1 planning. (VAL-008)
- [ ] Phase planning artifacts still map every `VAL-NNN` to a task/window and every task to covered `VAL-NNN` IDs. (VAL-007)
- [ ] Dependency-budget notes from task 1 and task 2 are reviewed for TRD §1/§9 compliance. (VAL-004)
- [ ] Final source review records no daemon, listener, remote sync path, or MCP adapter in the Phase 0 scaffold. (VAL-009)

## TDD Plan

- First failing test: N/A. This is an audit/documentation task.
- Additional tests: N/A.
- Regression risks: Treating future open questions as Phase 0 blockers, or
  silently editing specs to hide a real blocker.

## Documentation Updates

- `project-docs/docs/phases/phase-0/spec-audit.md` - required audit result.
- `project-docs/docs/phases/phase-0/tasks/task-3-completion.md` - required
  completion report.
- `project-docs/docs/current/progress_current.md` - update only if this task
  changes current phase status.

## Scenario / Parity / Audit Notes

- Held-out scenario impact: None. Do not read held-out scenarios.
- Cross-surface parity impact: Verify command references across AGENTS,
  Quality Plan, AI Workflow, Makefile, and CI.
- Audit or release-readiness impact: This is the Phase 0 spec audit. It is not
  release-readiness review.

## Exact Commands To Run

Run exactly:

```text
gofmt -l .
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
go build -o /tmp/fid ./cmd/fid
rg -n "<TBD>|PLACEHOLDER|placeholder|Proposed" AGENTS.md project-docs/specs project-docs/docs/current
rg -n "net.Listen|http.ListenAndServe|grpc|mcp|remote sync|daemon" cmd internal .github Makefile
```

The full task gate is required before task close. The two `rg` commands are
audit commands: a zero exit with matches can still be `PASS` only if
`spec-audit.md` classifies every match and records why it is blocking or
non-blocking. For the no-listener audit, exit 1 with no matches is an expected
`PASS` result if the completion report records that no matches were found.
Every required command must be recorded in the completion report as `PASS`, or
as `BLOCKED` with reason in a blocker report.

Conditional commands:

```text
go test -tags integration -race ./...
go test -bench . -benchmem ./...
```

List any conditional command required by this task and why:

- None unless integration or benchmark files are added unexpectedly.

## Notes For Developer

- If the audit finds a product, architecture, data-model, public-interface, or
  validation-contract decision gap, stop with a blocker report.
- Do not remove future open questions just because they are out of Phase 0
  scope.
- Source reads are acceptable evidence for the no-listener/no-MCP audit because
  the scaffold should contain only minimal code at this point.
