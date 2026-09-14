# Phase 23 - Task 16: Portable recovery for the review and judgment domain

## Objective

Extend portable export, non-mutating validation, and import to cover every
Phase 23 table and Event, and prove live, rebuilt, and imported projections are
semantically equivalent.

## Validation Contract Coverage

- VAL-004: Pragma checks are clean after the final schema and after rebuild.
- VAL-064: Export, validate, and import round-trip every Phase 23 table and
  Event with semantically equivalent current state.
- VAL-065: Import accepts only the binary's exact current accepted projection
  schema.

## Binding Constraints

- `public-interface.md` section 6 and `schema.md` section 7 are the exact
  authority. Manifest schema stays `fid.export-manifest.v2` with its Phase 22
  fields; only `projection_schema` moves, and it is 18 at phase end.
- `export validate` is always non-mutating (PRD-REC-002). Import preserves
  stable identities and history (PRD-REC-003) and stays empty-target.
- Equivalence is proven across live, rebuilt, and imported projections for every
  Phase 23 table, not sampled. Name each table in the evidence.
- A prior projection schema rejects `archive-invalid` with check code
  `schema-unsupported`. No old-journal conversion or dual-write is introduced.
- The kit coordinator's recovery-export binding remains an accepted Phase 25
  obligation under CF-037. Do not claim or imply that this task satisfies it.
- `execution-scope.md` is incorporated on the usual convergence terms.

## Scope

### Files To Create Or Modify

- `internal/export/export.go`, `archive_profiles.go`, import binding.
- `internal/journal/` rebuild parity for the Phase 23 tables.
- `cmd/fid/export.go` - result shapes only.
- Matching in-package tests and `cmd/fid/integration/` portable fixtures.
- `project-docs/docs/phases/phase-23/tasks/task-16-done.md`.

### Files To Read For Context

- `public-interface.md` section 6; `schema.md` sections 2, 4, 7;
  `contracts.md` section 15.
- `project-docs/docs/phases/phase-22/contracts.md` section 10 for the retained
  portable recovery rules.
- `task-15-done.md`; `execution-scope.md`; `contract-coverage.md`.

### Out Of Scope

- Capability advertisement and the conformance audit.
- Any coordinator-side activation or export binding.

## Acceptance Criteria

- [ ] Live, rebuilt, and imported projections are semantically equivalent for
  every named Phase 23 table; validate mutates nothing; identities and history
  are preserved. (VAL-064)
- [ ] A prior projection schema rejects `archive-invalid` with
  `schema-unsupported` while the current one imports. (VAL-065)
- [ ] Pragma checks are clean at the final schema and after rebuild. (VAL-004)
- [ ] Full required commands terminate with preserved red/green evidence.
- [ ] Every converged retained row records its provenance and executing command.

## First Failing Test

`TestPortableRoundTripCoversEveryPhase23Table` fails against the incomplete
archive profile.

## Documentation Updates

- `project-docs/docs/phases/phase-23/tasks/task-16-done.md` in the full format,
  naming every table proven equivalent.
## Exact Commands To Run

Working directory: repository root `/home/mark/projects/fid`, or the isolated
same-task worktree root explicitly recorded by the Runner. No command is a
placeholder. Apply gofmt and goimports to edited Go files before the final gate.

After implementation, every command below is required before task close.
The tagged integration gate is mandatory and is never substituted.

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
go test -tags integration -race ./cmd/fid/integration/...
go test -run '^$' -bench '^BenchmarkRebuildWorkspace100kEvents$' -benchtime 1x ./internal/journal/
```

Required-command rules live in `AGENTS.md` (Commands): every required command
is reported `PASS`, or `BLOCKED` with reason; missing, silently skipped, or
narrower substituted commands are failed gates. No exception to the full task
gate. Focused red/green selections never replace these commands. Report each
attempt with the command's own output token and observed termination. An
interrupted or hung command is never `PASS`; preserve all failed attempts and
rerun the byte-identical command after an authorized causal repair. Report lint
is a Runner control after the report commit, not a Developer gate.

List any conditional command required by this task and why:

- the named fresh-process rebuild benchmark, because this task changes rebuild and import coverage over the full corpus.

## Notes For Developer

- This task belongs to QA Review window 6.
- One task per Developer session; recovery may continue the same task.
- Commit using Conventional Commits with `Refs: phase-23/task-16`.
- Never infer ownership from a historical filename or task label. No-progress
  budgets exclude interruptions and attempts forbidden from repairing the cause.
- Do not read held-out scenario material or activate another agent.
