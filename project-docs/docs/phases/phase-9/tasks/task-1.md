# Phase 9 - Task 1: `fid concern raise`

## Objective

Add the `fid concern raise` write command: create a structured `concerns` row
linked to an issue and a task completion, with validate-before-append rejection
and read-back. This is the first write into the existing `concerns` schema and
proves the already-shipped derived concern flag (`concerns_raised` /
`open_concerns_count`) with a real writer.

## Requirements

- PRD refs: PRD-DEV-005 (structured concern linked to issue + completion),
  PRD-DEV-007 (derived concern flag on issue reads).
- TRD refs: §5.1 `concerns` schema; §2.2 `issues_enriched` derived columns;
  §4.3 write-path seam (validate/build before append).
- Design refs: `06_DESIGN_PLAN_CLI.md` §6 `fid concern raise`.
- Quality refs: `08_QUALITY_PLAN.md` (TDD, coverage gate; DEC-007 COV-001/002).

## Validation Contract Coverage

Assertions live in `project-docs/docs/phases/phase-9/qa/validation-contract.md`.

- VAL-001: `fid concern raise` creates a `CON-<hash>` concern row with all named
  fields and `status='open'`, read back via `fid concern list` /
  `fid issue get --section concerns`.
- VAL-002: raise is rejected before any journal append for unknown role,
  nonexistent issue, nonexistent completion, and completion-not-on-issue; valid
  raise appends exactly one event/row (CF-007 / DEC-C).
- VAL-003: raising an open concern flips `concerns_raised` false→true and
  increments `open_concerns_count` on `fid issue get` and `fid issue list`.

## Scope

### Files To Create Or Modify

- `internal/issues/concern_raise.go` (new) - `RaiseConcern` domain function:
  `RequireKnownRole`, validate issue exists, validate `--in-completion` exists
  and belongs to the issue, generate `CON-<hash>`, build `insert` effect under
  `.fid/write.lock`, dry-run project, append, project. Follow
  `internal/issues/ac_verify.go` / `complete.go` structure verbatim.
- `internal/issues/concern_raise_test.go` (new) - unit/integration tests.
- `cmd/fid/concern.go` (new) - `fid concern` parent command + `raise` subcommand
  wiring flags `--on`, `--by`, `--in-completion`, `--summary`, `--body-md`, and
  the stderr confirmation line from design §6.
- `cmd/fid/concern_test.go` (new) - command-level tests.
- `cmd/fid/root.go` - register the `concern` command.

### Files To Read For Context

- `internal/issues/ac_verify.go` (closest write precedent: validate-before-append
  + lock + dry-run + append + project).
- `internal/issues/complete.go` (completion lookups; `requireCompletionIssue`).
- `internal/issues/detail.go` (`issueConcerns`, `IssueConcern` struct).
- `internal/issues/get.go` / `list.go` (derived `concerns_raised` read path).
- `internal/issues/transition_rules.go` (`loadConcernStatus`).
- `cmd/fid/issue.go` (command wiring + stderr confirmation patterns).

### Out Of Scope

- `fid concern address` (Task 2) and `fid concern list` top-level command
  (Task 3 — Task 1 reads back via `fid issue get --section concerns` and may
  add a minimal `list` only if Task 3 is not yet present; prefer reading back
  via `fid issue get` to keep scope tight).
- `escalated`/`wontfix` statuses; `triaged_in_event_id` linkage.
- `activity_log` rows for concern writes (DEC-E: none, by design).
- Any workflow-edge change (DEC-A: no edge gates on `concerns_addressed`).

## Acceptance Criteria

- [ ] `fid concern raise --on <task> --by <role> --in-completion <TC> --summary`
  creates a `CON-<hash>` row with `issue_id`, `raised_in_completion_id`,
  `raised_by_role`, `summary`, optional `body_md_path`, `status='open'` (VAL-001)
- [ ] Read-back via `fid issue get --section concerns --format json` shows the
  concern with those fields (VAL-001)
- [ ] Rejected before append (non-zero exit, clear error, zero new journal
  lines/rows) for: unknown role; nonexistent `--on`; nonexistent
  `--in-completion`; completion not belonging to `--on` (VAL-002)
- [ ] Valid raise appends exactly one journal event and one concern row (VAL-002)
- [ ] `concerns_raised` flips false→true and `open_concerns_count` increments on
  `fid issue get` and `fid issue list` after a real raise (VAL-003)

## TDD Plan

- First failing test: `RaiseConcern` creates a concern row with `status='open'`
  and all fields, read back from the DB.
- Additional tests: each validate-before-append rejection asserts zero new
  journal lines (count journal file) + error; completion-belongs-to-issue guard;
  derived-flag round-trip via `Get`/`List`; command-level `cmd/fid` test for flag
  parsing + stderr confirmation + exit codes.
- Regression risks: journal-poisoning on rejected writes (CF-007) — assert zero
  appends on every reject path; ID-collision retry via `journal.GenerateID`.

## Documentation Updates

- `06_DESIGN_PLAN_CLI.md` §6 `fid concern raise` - confirm flag/behavior match;
  update if the validate-before-append linkage rules (DEC-C) need stating.
- Completion report - record DEC-C linkage rules as implemented and any deviation.

## Scenario / Parity / Audit Notes

- Held-out scenario impact: None.
- Cross-surface parity impact: confirm `concerns_raised`/`open_concerns_count`
  read identically from `fid issue get` and `fid issue list` (already-shipped
  read path; do not re-derive).
- Audit or release-readiness impact: None.

## Exact Commands To Run

Run exactly:

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

The full task gate above is required before task close. The coverage output must
show `internal/issues` ≥ 80% and `cmd/fid` ≥ 70% (DEC-007 COV-001/002); record
the percentages in the completion report.

Conditional commands:

```text
go test -tags integration -race ./...
go build -o /tmp/fid ./cmd/fid
```

- Run `go build -o /tmp/fid ./cmd/fid` (new command surface).
- Run the integration suite if any test is added under the `integration` tag.

## Notes For Developer

- Mirror `ac_verify.go` precisely for the lock/dry-run/append/project sequence
  and the `issueACVerifyArgv` argv helper; name the helper `issueConcernRaiseArgv`.
- Validate-before-append: do all existence/ownership/role checks inside
  `WithLock` **before** `AppendEvent`, exactly as `complete.go` /`ac_verify.go` do.
- The completion lookup should reuse / parallel `completionIDExists` and verify
  the completion's `issue_id` equals `--on` (DEC-C). A completion only exists for
  task/subtask issues, so this implicitly rejects epics.
- Keep new functions ≤60 lines and new files ≤400 lines. Do **not** add concern
  write logic to `detail.go` (already near the 400-line budget).
- `summary` is required; `body_md_path` is optional. Emit the design §6 stderr
  confirmation: `raised concern CON-... on FID-...: "<summary>"`.
