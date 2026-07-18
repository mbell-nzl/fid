# Phase 18 - Task 5: Remaining consumer swap + maintenance parity + Phase 18 closure

## Objective

Repoint the remaining consumers of the package-global workflow Go map to the
DB-backed registry — `fid maintenance check` (§14.11 `(workflow_id, status)`
audit), issue get/list status-category display, the `fid dev *` composites, and
`fid ready` — then prove the fully-migrated system end-to-end: v0.1 behavior
parity, deterministic rebuild, non-destructive delete-and-rebuild, and
unchanged v0.1 journal history.

## Validation Contract Coverage

- VAL-010: `fid maintenance check` validates `(workflow_id, status)` against the workflow tables — clean on healthy, flags a seeded invalid pair, non-zero under `--warn-as-error`.
- VAL-011: `fid rebuild` reproduces the migrated projection deterministically; delete-and-rebuild is non-destructive (per-table parity).
- VAL-013: v0.1 journal history is never rewritten; existing events replay unchanged into the migrated schema (golden fixture); journal files byte-unchanged.

## Binding Constraints

- **Swap every remaining `workflow.*` package-function consumer to the registry.**
  From the current call sites: `internal/maintenance/check.go` +
  `internal/maintenance/warnings.go` (§14.11 workflow-pair audit and any
  category-derived warnings), `internal/issues/get.go` / `issue.go` /
  `detail.go` / `list.go` / `ready.go` (status-category display and ordering),
  and `internal/dev/{start,finish,block,dev}.go` (category checks in the
  composites). After this task, **no production code path reads the Go
  `definitions` map**; it remains only as frozen parity fixtures for tests.
- **Maintenance §14.11 now validates against the workflow tables.** The audit
  must resolve each issue's `(workflow_id, status)` against `workflow_states`
  (the migrated source of truth) rather than the Go map. Behavior is otherwise
  unchanged: deterministic findings ordered by check, entity ID, then message;
  `--warn-as-error` exits non-zero on a violation.
- **Behavior parity everywhere.** Category display, ready-queue ordering, and
  dev-composite category gating must be identical to v0.1. Because the seed and
  the Go map are behavior-identical, this is a source swap, not a semantic change.
- **Closure evidence.** Provide a delete-and-rebuild per-table parity check over
  the migrated schema and a golden replay fixture proving a v0.1 journal
  reprojects unchanged, and assert journal files are byte-identical across
  migrate/rebuild (DEC-066: history never rewritten).
- The gate-obligation seam remains behavior-free; re-confirm no receipt/gate
  code is reachable from any consumer (VAL-005 stays true after the swap).

## Scope

### Files To Create Or Modify

- `internal/maintenance/check.go`, `internal/maintenance/warnings.go`
- `internal/issues/get.go`, `internal/issues/issue.go`, `internal/issues/detail.go`,
  `internal/issues/list.go`, `internal/issues/ready.go`
- `internal/dev/start.go`, `internal/dev/finish.go`, `internal/dev/block.go`, `internal/dev/dev.go`
- Their tests, plus a phase-closure test/fixture (delete-and-rebuild parity;
  golden v0.1 journal replay). Existing e2e closure tests under `cmd/fid/`
  updated/extended as needed.

### Files To Read For Context

- `internal/workflow/registry.go` (loader)
- `project-docs/specs/03_TRD_SCHEMA.md` §14.11 (maintenance parity)
- Current consumers listed above (call sites of `workflow.CategoryForState` /
  `RulesForEdge` / `Workflow`)
- `internal/journal/` rebuild + golden-fixture tests

### Out Of Scope

- Deleting the Go `definitions` map (kept as parity fixtures).
- `fid workflow define/assign/list`, escalation-default (Phase 19).
- `role_permissions` migration (Phase 20); any gate/receipt behavior (v0.2b).

## Acceptance Criteria

- [ ] No production code path reads the Go `definitions` map; all consumers use
  the DB registry
- [ ] `fid maintenance check` resolves `(workflow_id, status)` from the workflow
  tables: clean on healthy, flags a seeded invalid pair, non-zero under
  `--warn-as-error` (VAL-010)
- [ ] Issue get/list category display and ready-queue ordering are identical to v0.1
- [ ] Delete-and-rebuild is non-destructive with per-table parity over the
  migrated schema (VAL-011)
- [ ] A golden v0.1 journal replays unchanged into the migrated projection;
  journal files are byte-identical across migrate/rebuild (VAL-013)
- [ ] Full phase-level gate green

## First Failing Test

A maintenance test that seeds an issue whose `(workflow_id, status)` pair is not
in `workflow_states` and asserts `fid maintenance check` reports a structural
finding sourced from the tables — fails until the audit reads the registry.

## Documentation Updates

- None authored by the Developer. The Tech Lead syncs `03_TRD_SCHEMA.md`
  §14.11 / §4.6 and `02_TRD.md` wording (maintenance now audits against the
  workflow tables; defaults are migrated) at phase wrap. Flag any behavior
  nuance discovered here in the completion report.

## Exact Commands To Run

Run exactly the full phase-level gate:

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
go test -tags integration -race ./...
go test -bench . -benchmem ./...
go build -o /tmp/fid ./cmd/fid
```

Required-command rules live in `AGENTS.md` (Commands): every required command
is reported `PASS`, or `BLOCKED` with reason; missing, silently skipped, or
narrower substituted commands are failed gates.

## Notes For Developer

- Closure task: depends on Tasks 1–4. Run the full phase-level gate (this task
  uses the phase-level gate, not the shorter task gate) because it proves the
  v0.2a slice-1 exit criteria.
- If any consumer's behavior would change under the registry, stop and escalate
  rather than adjusting the assertion.
- Completion report at `project-docs/docs/phases/phase-18/tasks/task-5-done.md`.
