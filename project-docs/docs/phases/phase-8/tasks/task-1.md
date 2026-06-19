# Phase 8 - Task 1: `fid issue complete` core write (attempt numbering, supersession, human intervention)

## Objective

Implement the `fid issue complete` command and its `internal/issues` write path:
an append-only `task_completions` row with derived attempt numbering,
prior-attempt supersession, the basic completion fields (tokens, commits, files,
notes), and the human-intervention flag/note. Reject completions on epics before
any journal append. No gate snapshot in this task (Task 2) — `gate_results` is
empty and `gates_passed`/`gates_failed` default to 0.

## Requirements

- PRD refs: PRD-DEV-001 (append-only completion record), PRD-DEV-002 (multiple
  numbered attempts), PRD-DEV-003 (mark prior attempt superseded), PRD-DEV-009
  (human-intervention flag/note)
- TRD refs: TRD_SCHEMA §2 `task_completions`; TRD §4.3 write-path seam
  (validate/build before append); ADR-0001
- Design refs: `06_DESIGN_PLAN_CLI.md` §2 `fid issue complete`
- Quality refs: `08_QUALITY_PLAN.md` (TDD, validate-before-append, DEC-007
  coverage); CF-007

## Validation Contract Coverage

- VAL-001: `fid issue complete` creates a `TC-<hash>` row storing developer_role,
  tokens, commits (JSON array), files (JSON array), notes path; read-back shows them.
- VAL-002: First attempt is 1; second attempt is 2 and prior attempt is retained
  with `status = superseded`.
- VAL-003: `--human-intervention`/`--intervention-note` set the flag + note;
  omitting them yields `human_intervention = 0` and no note.
- VAL-004: `fid issue complete` on an epic is rejected before any journal append;
  a task/subtask succeeds.

## Scope

### Files To Create Or Modify

- `internal/issues/complete.go` (new) - `Complete(options)` write path mirroring
  `VerifyAcceptanceCriterion` (EnsureProjectionFresh → WithLock → validate →
  GenerateID `TC` → derive attempt → build event with insert + supersede update
  effects → DryRunProjectEvent → AppendEvent → ProjectEvent).
- `internal/issues/complete_test.go` (new) - unit/integration tests.
- `cmd/fid/issue.go` - register `issue complete <id>` with flags `--by`,
  `--notes-md`, `--tokens`, `--commits`, `--files` (repeatable), and
  `--human-intervention`, `--intervention-note`; stderr confirmation per design.
- `cmd/fid/issue_complete_test.go` (new) - command-level tests.

### Files To Read For Context

- `internal/issues/ac_verify.go` (write-path / effect / projection pattern to mirror).
- `internal/issues/detail.go` `issueCompletions` (read-back shape, JSON fields).
- `internal/issues/transition_rules.go` `loadCompletionReportStatus` (the count
  this completion will satisfy — but do not change the rule here; Task 3 wires it).
- `migrations/0004_v01_kind_triggers.up.sql` `task_completions_issue_kind_guard_insert`
  (the epic-rejection backstop the CLI must pre-empt).
- `06_DESIGN_PLAN_CLI.md` §2 `fid issue complete` (flags + stderr output).

### Out Of Scope

- Gate snapshot (`gate_runs` → `gate_results`, `gates_passed`/`gates_failed`):
  Task 2. Leave `gate_results` empty and counts at 0 here.
- `has_completion_report` transition wiring and rebuild parity: Task 3.
- `internal/issues/issue.go` size split: Task 4.
- `model` capture, `expected_files_touched`, concerns, and any `fid dev *`
  composite: later phases.

## Acceptance Criteria

- [ ] `fid issue complete <task> --by developer` creates one `TC-<hash>`
  `task_completions` row; `fid issue get --format json` reads back
  `developer_role`, `actual_tokens`, `commit_hashes`, `files_touched`, and the
  notes path with stored values (VAL-001).
- [ ] `--commits` (comma-separated) and `--files` (repeatable) persist as JSON
  arrays; stderr confirmation matches the design block (`completed <id> (attempt N)`,
  `completion: TC-...`, tokens, commits) (VAL-001).
- [ ] `attempt_number` is derived `MAX(attempt_number)+1` (1 when none) under the
  write lock; a second completion on the same task is attempt 2 and the first
  row is retained with `status = superseded`; the new row is `dev_complete`
  (DEC-A) (VAL-002).
- [ ] `--human-intervention --intervention-note "<text>"` stores
  `human_intervention = 1` + note; omitting the flag stores `0` + no note (VAL-003).
- [ ] `fid issue complete` on an **epic** exits non-zero with a clear error
  naming the kind constraint and appends **zero** journal lines/rows; on a
  task/subtask it succeeds (VAL-004).
- [ ] Unknown `--by` role is rejected via `policy.RequireKnownRole` before append;
  a missing issue id returns a `NotFoundError`-style error before append.
- [ ] No `activity_log` row is emitted by `fid issue complete` (DEC-D / DEC-020
  precedent).

## TDD Plan

- First failing test: `fid issue complete` on a fresh task creates attempt 1 with
  the supplied fields and a `TC-` id.
- Additional tests: second completion → attempt 2 + prior `superseded` + both rows
  retained; human-intervention both states; epic rejection with zero new journal
  lines; unknown role rejection; JSON read-back of all fields.
- Regression risks: attempt-number derivation must run under the lock against the
  fresh projection; supersession update effects must target each prior completion
  PK individually (the projection's `applyUpdate` asserts single-row-affected).

## Documentation Updates

- `06_DESIGN_PLAN_CLI.md` §2 `fid issue complete` - only if the implemented flags
  or stderr output differ from the spec; otherwise note "verified, no change" in
  the completion report (spec-sync input for wrap).
- `project-docs/docs/phases/phase-8/` completion report - paste full real CLI
  transcripts for create→complete→complete (attempt 2)→epic-reject (process nudge:
  pasted exact command transcripts for public-output changes).

## Scenario / Parity / Audit Notes

- Held-out scenario impact: None (do not read hidden scenario content).
- Cross-surface parity impact: completion fields must read back consistently via
  `fid issue get` JSON and any markdown/section view that surfaces completions.
- Audit or release-readiness impact: None.

## Exact Commands To Run

Run exactly:

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

Coverage targets (DEC-007): `cmd/fid` ≥ 70%, `internal/issues` ≥ 80%.

Conditional commands:

```text
go build -o /tmp/fid ./cmd/fid
```

List any conditional command required by this task and why:

- `go build -o /tmp/fid ./cmd/fid` to produce the binary for the pasted CLI
  transcripts in the completion report.

## Notes For Developer

- Mirror `VerifyAcceptanceCriterion` for the lock/effect/projection seam and
  `issueCreateArgv`/`issueACVerifyArgv` for the sanitized command argv.
- The new completion's `status` is `dev_complete` and prior attempts become
  `superseded` (DEC-A); do not derive `failed_gate` from gates — that is the
  v0.1e `fid dev block` composite's job.
- File-size acceptance is scoped to the files this task writes (`complete.go`,
  `complete_test.go`, the `issue.go` command additions) — keep each new non-test
  file ≤400 lines. The pre-existing `internal/issues/issue.go` 677-line tripwire
  (CF-012) is Task 4; do not let it block this task.
- Reject the epic case in the CLI/write path **before** the journal append; the
  `task_completions_issue_kind_guard_insert` trigger is only a backstop (CF-007).
