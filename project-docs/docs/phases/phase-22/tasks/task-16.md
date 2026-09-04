# Phase 22 - Task 16: Retire escalation state and duplicate recovery truth

## Objective

Remove the exact Phase 22 legacy inventory after its safety replacements pass. Keep activity as event-derived inspection and retain explicitly later-owned historical state.

## Validation Contract Coverage

- VAL-066: Latest schema has no escalation entity or foreign-key references.
- VAL-067: Recovery/activity expose one Event-derived source of truth.
- VAL-068: Retained assertions displaced by removal converge under executing gates.

Full sources, pass/fail criteria, and validator ownership are in
`project-docs/docs/phases/phase-22/qa/validation-contract.md`.

## Binding Constraints

- PROPOSED: execution requires Human acceptance of DEC-180 and a Ready QA plan.
- schema.md section 4; v1-convergence-inventory.md.
- Removal roots and retained exceptions are exact in the inventory; classify every match, including tests and migration-scoped fixtures.
- Do not drop triage_events or its Issue projection (Phase 24), or Concern/Review/telemetry history owned by Phases 23/25.
- Latest-schema tests change to absence; historical migration tests retain their original scope/strength.
- Current-schema tests must reject removed-table Effects; no permissive replay fallback.
- DEC-186/R-13: record one fresh schema-16 reference rebuild sample before
  hand-back. Preserve parity and <30000 ms self-assertion; do not defer a causal
  removal regression to task 17 or replace its three final samples with this one.
- `execution-scope.md` is incorporated into this scope: it assigns only
  subject-preserving convergence of assertions displaced by this task's accepted
  change. All owned ordinary and tagged code executes under the commands below.

## Scope

### Files To Create Or Modify

- `migrations/0016_p22_remove_escalations.up.sql`
- `migrations/0016_p22_remove_escalations.down.sql`
- `internal/maintenance/`
- `internal/activity/`
- `internal/workflow/rules.go`
- `internal/db/`
- `internal/journal/`
- `internal/issues/`
- `internal/export/`
- `cmd/fid/`
- Matching in-package tests and `cmd/fid/integration/` contract fixtures for
  the behavior above, including retained rows invalidated by this change.
- `project-docs/docs/phases/phase-22/tasks/task-16-done.md`.

### Files To Read For Context

- `project-docs/docs/phases/phase-22/contracts.md` - only the sections named above.
- `project-docs/docs/phases/phase-22/public-interface.md` - owned input/result/error types.
- `project-docs/docs/phases/phase-22/schema.md` - owned DDL and Event rows.
- `project-docs/docs/phases/phase-22/execution-scope.md` and `v1-convergence-inventory.md`.
- Canonical source sections cited by this task's validation rows.
- `project-docs/docs/phases/phase-22/tasks/task-15-done.md` - actual predecessor hand-back.

### Out Of Scope

- Workflow control scripts, delivery escalation docs, later-phase product removals, and old-data migration.
- Product policy, protected gate definitions, workflow controls, dependencies,
  and unrelated implementation changes beyond this accepted task.

## Acceptance Criteria

- [ ] Latest schema has no escalation entity or foreign-key references. (VAL-066)
- [ ] Recovery/activity expose one Event-derived source of truth. (VAL-067)
- [ ] Retained assertions displaced by removal converge under executing gates. (VAL-068)
- [ ] Full required commands terminate; preserve red/green and failed-attempt
  evidence with output tokens and observed exit status.
- [ ] Every converged retained row identifies introducing commit/phase/task,
  invalidating change, and the exact command that executes it.

## First Failing Test

TestFinalSchemaHasNoEscalationsOrEscalationReferences fails against schema 15.
Capture the failure before product edits; then implement the smallest behavior
that satisfies it. Build failures alone do not prove the intended rejection.

## Documentation Updates

- `project-docs/docs/phases/phase-22/tasks/task-16-done.md`: implementation,
  red/green, every exact command, observed termination, assertion evidence,
  convergence provenance, dependencies, and deviations.
- Flag spec impact against the named contracts; only Tech Lead accepts a
  contract/schema change. Record no spec change when implementation conforms.

## Exact Commands To Run

Working directory: repository root `/home/mark/projects/fid`, or the isolated
same-task worktree root explicitly recorded by the Runner. No command is a
placeholder. Apply gofmt/goimports to edited Go files before the final gate.

After implementation, every command below is required before task close.
The integration gate is mandatory from task 1, including for retained fixtures.

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
go test -tags integration -race ./cmd/fid/integration/...
go version
go env GOOS GOARCH
uname -srm
nproc
grep -m1 '^model name' /proc/cpuinfo
grep -m1 '^MemTotal:' /proc/meminfo
go test -run '^$' -bench '^BenchmarkRebuildWorkspace100kEvents$' -benchmem -benchtime 1x -timeout 30m ./internal/journal
```

No exception to the full task gate. Commands whose expected pass is nonzero:
none. Focused red/green selections do not replace these commands. Report each
attempt using the command's own output token and observed termination, as
required by AGENTS.md. No interrupted/hung command is PASS; preserve all failed
attempts and rerun the byte-identical command after authorized causal repair.
Report lint is a Runner control after the report commit, not a Developer gate.

## Notes For Developer

- Start after task 15 closes.
- One task per Developer session; recovery may continue the same task.
- This task belongs to QA Review window 6.
- Commit using Conventional Commits with `Refs: phase-22/task-16`.
- Never infer ownership from a historical filename or task label. No-progress
  budgets exclude interruptions and attempts forbidden from repairing the cause.
- Do not read held-out scenario material or activate another agent.
