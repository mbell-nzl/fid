# Phase 11 - Task 3: `fid claim list`

## Objective

Add the read surface for claim inspection: `fid claim list` with live/stale
filtering, role filtering, all four formats, required fields, and deterministic
append-order output.

## Requirements

- PRD refs: PRD-CLAIM-007, PRD-CLAIM-008
- TRD refs: `03_TRD_SCHEMA.md` section 11.1; TRD-ARCH-005 boundary rules
- Design refs: `06_DESIGN_PLAN_CLI.md` `fid claim list`; `05_DESIGN_PLAN.md`
  claim state field parity
- Quality refs: `08_QUALITY_PLAN.md` COV-001/COV-002/COV-003
- Decisions: `decisions_principles.md` deterministic read ordering rule; Phase
  11 DEC-027

## Validation Contract Coverage

- VAL-008: `fid claim list` fields, `--by` filter, and json/toon/markdown/yaml parity.
- VAL-009: `--stale` returns only expired current claims and ordering is rowid-stable.
- VAL-010: human-held claims display as `human`.

## Scope

### Files To Create Or Modify

- `internal/claims/list.go` and tests - query current claims, compute
  `expires_at` and `is_stale`, filter by stale and claimed-by role, order by
  `claims.rowid`.
- `cmd/fid/claim.go` and tests - root `claim` command with `list` subcommand,
  `--stale`, `--by`, and `--format json|toon|markdown|yaml`.
- `cmd/fid/root.go` - register `newClaimCommand`.
- `internal/format/` only if the existing formatter needs a narrow claim-list
  adapter; prefer existing generic format helpers.

### Files To Read For Context

- Existing list commands in `cmd/fid/issue.go`, `cmd/fid/policy.go`, and
  `internal/issues/list.go`.
- `internal/issues/detail.go` claim section for existing claim row shape.
- `internal/format/` json/toon/markdown/yaml conventions.
- `decisions_principles.md` read-ordering rule.

### Out Of Scope

- Mutating claim commands (Task 2).
- Named stale-claims report (v0.2).
- `fid ready --explain` claimed/stale flags (v0.2).
- Permission checks on `claim list`; it is read-only.

## Acceptance Criteria

- [ ] `fid claim list` defaults to JSON and returns current non-stale claims
  with `ref_type`, `ref_id`, `claimed_by`, `claimed_at`, `ttl_seconds`,
  `expires_at`, `is_stale`, and `note` (VAL-008).
- [ ] `--by <role>` filters by the stored `claimed_by` value and is treated as a
  read filter, not an acting role or authorization parameter (VAL-008).
- [ ] `--format toon`, `--format markdown`, and `--format yaml` render the same
  data as JSON, with structural parity tests (VAL-008).
- [ ] `--stale` returns only past-TTL current claims and excludes live claims
  (VAL-009).
- [ ] Default ordering uses `claims.rowid` and remains stable across rebuild for
  the same journal (VAL-009).
- [ ] A human-held claim appears as `claimed_by: human` in list output (VAL-010).
- [ ] `fid claim list` appends no journal line and writes no `activity_log` row.

## TDD Plan

- First failing test: `fid claim list` is registered and returns a JSON array for
  a workspace with one live claim.
- Additional tests: all formats match JSON reference; `--by` filter; `--stale`
  filter; empty results; rowid ordering; read-only journal/activity assertion;
  human claim display.
- Regression risks: do not order by `claimed_at` because same-second ties can be
  nondeterministic; do not route `--by` through `Authorize` because this command
  is read-only.

## Documentation Updates

- `project-docs/docs/phases/phase-11/tasks/task-3-done.md` - completion report
  with command results, coverage, and commit hash/subject.
- `project-docs/docs/phases/phase-11/decisions.md` - record any output-shape
  deviation from `06_DESIGN_PLAN_CLI.md`.
- `06_DESIGN_PLAN_CLI.md` - edit only if implementation proves the existing
  list contract wrong; otherwise record verified-no-change in the completion
  report.

## Scenario / Parity / Audit Notes

- Held-out scenario impact: None.
- Cross-surface parity impact: list fields should match the existing issue
  claims section where the same field exists.
- Audit or release-readiness impact: public read command requires build
  verification and read-only proof.

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

The full gate above is required before task close because this task adds a root
read command and format surface. A focused test command may be used during the
red/green loop, but it does not substitute for the full gate at close. Every
required command must be recorded in the completion report as `PASS`, or as
`BLOCKED` with reason in a blocker report.

Conditional commands:

```text
go test -bench . -benchmem ./...
```

- `go test -bench . -benchmem ./...` - required only if benchmark fixtures or
  list performance paths are touched; otherwise record N/A with reason.

## Notes For Developer

- Use `ORDER BY rowid` for deterministic default ordering unless a spec is
  updated before implementation starts.
- `--by` is a filter over `claimed_by`; unknown strings should simply filter to
  no rows unless an existing read-filter pattern clearly requires otherwise. If
  that conflicts with local patterns, stop and escalate.
- Keep format rendering in the adapter/format layer. Claim query and stale
  classification belong in `internal/claims`.
- Pair completion-report commit-hash filling with the task commit or immediate
  reconciliation commit.
