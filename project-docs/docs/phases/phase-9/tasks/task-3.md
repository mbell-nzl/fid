# Phase 9 - Task 3: `fid concern list`, rebuild parity, v0.1c closure smoke

## Objective

Add the top-level `fid concern list` read command (filter by issue and status,
all formats), prove delete-DB + `fid rebuild` parity for concerns and derived
flags, and run the end-to-end v0.1c closure smoke across the whole slice. This
closes the v0.1c roadmap version.

## Requirements

- PRD refs: PRD-DEV-006 / PRD-DEV-007 (`fid concern list` filters and read).
- TRD refs: §5.1 `concerns` schema; ADR-0001 (journal-canonical, DB derived).
- Design refs: `06_DESIGN_PLAN_CLI.md` §6 `fid concern list` (JSON shape incl.
  `issue_id`).
- Quality refs: `08_QUALITY_PLAN.md` (TDD, coverage gate; DEC-007); roadmap
  §v0.1c exit criteria.

## Validation Contract Coverage

Assertions live in `project-docs/docs/phases/phase-9/qa/validation-contract.md`.

- VAL-008: `fid concern list` filters by `--issue` and `--status` (composing),
  across json/toon/markdown/yaml; empty match returns an empty collection, exit 0.
- VAL-009: delete DB + `fid rebuild` reproduces identical `concerns` state and
  derived `concerns_raised`/`open_concerns_count`.
- VAL-010: end-to-end v0.1c smoke ac add → ac verify → ac list → issue complete →
  concern raise (flag lights) → concern address (flag clears) → concern list.

## Scope

### Files To Create Or Modify

- `internal/issues/concern_list.go` (new) - `ListConcerns` returning concern
  rows (incl. `issue_id`) with optional issue/status filters; new
  `ConcernListItem` struct matching the design §6 JSON shape. Do **not** extend
  `detail.go` (already near the 400-line budget).
- `internal/issues/concern_list_test.go` (new).
- `cmd/fid/concern.go` - add the `list` subcommand (`--issue`, `--status`,
  `--format`) using the existing `internal/format` engine.
- `cmd/fid/concern_test.go` - extend with list/format/filter tests.
- `cmd/fid/e2e_v01c_concern_test.go` (new) - v0.1c closure smoke (VAL-010).
- `cmd/fid/rebuild_test.go` or a concern rebuild test - VAL-009 parity.

### Files To Read For Context

- `internal/issues/list.go` (filter/query + format dispatch precedent).
- `internal/issues/detail.go` (`IssueConcern` fields; concern query shape).
- `internal/format/format.go` (format dispatch for read commands).
- `cmd/fid/e2e_v01c_ac_test.go`, `cmd/fid/e2e_v01c_complete_test.go` (existing
  v0.1c e2e structure to extend).
- `cmd/fid/rebuild.go` / `rebuild_test.go` (rebuild parity pattern).

### Out Of Scope

- `escalated`/`wontfix` producers (listable if present; no Phase 9 writer).
- Any workflow-edge change (DEC-A).
- `activity_log` rows (DEC-E).

## Acceptance Criteria

- [ ] `fid concern list` returns concerns; `--issue` and `--status` filter and
  compose; output renders in json (default), toon, markdown, yaml (VAL-008)
- [ ] A filter matching nothing returns an empty collection, exit 0 (VAL-008)
- [ ] Delete DB + `fid rebuild` reproduces identical `concerns` rows (status,
  `raised_in`/`addressed_in_completion_id`) and identical derived flags (VAL-009)
- [ ] End-to-end v0.1c smoke runs the full slice with the derived flag lighting
  up on raise and clearing on address, every step asserted (VAL-010)

## TDD Plan

- First failing test: `ListConcerns` returns seeded concerns with `issue_id`;
  `--issue` and `--status` filters narrow correctly; combined filter; empty match
  returns empty slice.
- Additional tests: format dispatch per format; rebuild parity (snapshot →
  delete DB → rebuild → compare); e2e smoke driving the real command path.
- Regression risks: format engine field shape drift vs design §6 JSON; empty
  result must be `[]`/empty, not null/error.

## Documentation Updates

- `06_DESIGN_PLAN_CLI.md` §6 `fid concern list` - confirm JSON shape (incl.
  `issue_id`) and filter behavior.
- `project-docs/docs/current/progress_current.md` - mark v0.1c closure (handled
  at phase wrap by Tech Lead; note readiness in the completion report).
- Completion report - record the v0.1c slice as closed end-to-end.

## Scenario / Parity / Audit Notes

- Held-out scenario impact: None.
- Cross-surface parity impact: confirm `fid concern list` concern fields match
  the nested `fid issue get --section concerns` shape plus `issue_id`.
- Audit or release-readiness impact: v0.1c version closure — note in the
  completion report for the Tech Lead phase wrap.

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

This task ships a read command, a rebuild-parity test, and an end-to-end smoke,
so the integration suite and build verification are **required**, not optional.
Coverage must show `internal/issues` ≥ 80% and `cmd/fid` ≥ 70% (DEC-007); record
percentages in the completion report.

Conditional commands:

```text
go test -bench . -benchmem ./...
```

- Run the bench suite only if a concern read path is added to a hot query (not
  expected); otherwise N/A.

## Notes For Developer

- `fid concern list` is a top-level command (`fid concern list`), not under
  `fid issue`; add it to the `concern` parent command from Task 1.
- The list JSON shape must include `issue_id` (design §6) — the nested
  `IssueConcern` struct does not carry it, so use a dedicated `ConcernListItem`.
- Empty results must serialize as an empty collection with exit 0 (mirror
  `fid issue list`).
- Keep new functions ≤60 lines and new files ≤400 lines.
