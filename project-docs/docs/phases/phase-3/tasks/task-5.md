# Phase 3 - Task 5: Conflict / Corruption Error Fidelity (Rebuild + Read Catch-Up)

## Objective

Ensure both failure modes report with full fidelity on **both** the explicit
`fid rebuild` path and the automatic read catch-up path (`fid issue get`):

1. **Corruption** — a malformed trailing journal line is reported as a structured
   corruption error (`file:line` + parse cause), never silently skipped.
2. **Logical conflict** — a failed effect precondition during replay stops before
   serving reads and reports a structured conflict naming the event, entity,
   expected, and actual values (PRD-OPS-012).

The replayer already returns `CorruptJournalError` and `PreconditionConflictError`
(covered by `internal/journal` tests), and `fid rebuild` renders corruption. The
gaps this task closes: `fid rebuild` does **not** render `PreconditionConflictError`
structurally (it falls through to a generic error), and the read catch-up path
(`renderIssueError`) renders neither corruption nor conflict structurally — so a
read over a corrupt/conflicting journal must stop and report, not serve stale data.

## Requirements

- PRD refs: PRD-OPS-012 (structured logical conflict on failed precondition;
  stop before serving reads), PRD-OPS-003/004 (rebuild before serving stale
  reads; `fid rebuild`)
- TRD refs: TRD-ARCH-002 (self-contained events; replay does not re-run command
  logic), TRD-ARCH-003 (versioned upgraders), TRD-ARCH-005 (adapter boundary)
- Design refs: `06_DESIGN_PLAN_CLI.md` §18 (rebuild error examples),
  `05_DESIGN_PLAN.md` §5.4 (error message standard: constraint + entities +
  corrective action)
- Quality refs: Quality Plan §4.3 (COV-001 `internal/journal` ≥ 80%, COV-002
  `cmd/fid` ≥ 70%)

## Validation Contract Coverage

Assertions live in `project-docs/docs/phases/phase-3/qa/validation-contract.md`.

- VAL-008: malformed trailing line reported as corruption (not skipped) on both
  rebuild and read catch-up; clean journal replays cleanly.
- VAL-009: failed precondition stops before serving reads and reports a
  structured logical conflict (event/entity/expected/actual) on both rebuild and
  read catch-up; satisfiable journal replays cleanly.

## Scope

### Files To Create Or Modify

- `cmd/fid/rebuild.go` (`renderRebuildError`: add structured
  `PreconditionConflictError` rendering naming event id/type, entity table+pk,
  expected, actual)
- `cmd/fid/issue.go` (`renderIssueError` / read path: render `CorruptJournalError`
  and `PreconditionConflictError` structurally so a read over a corrupt/
  conflicting journal stops and reports, not serves stale data) — consider a
  shared error-render helper to avoid duplication between rebuild and read
- `cmd/fid/rebuild_test.go` (extend: rebuild renders structured conflict; clean
  rebuild succeeds)
- `cmd/fid/issue_test.go` (extend: `fid issue get` over corrupt journal → structured
  corruption; over conflicting journal → structured conflict + non-zero exit;
  clean journal → read succeeds)
- `internal/journal/*_test.go` (extend only if a both-states gap exists; the
  corruption + conflict replay behaviour is already covered)
- `project-docs/docs/phases/phase-3/tasks/task-5-completion.md`

### Files To Read For Context

- `cmd/fid/rebuild.go` (`renderRebuildError`, currently corruption-only)
- `cmd/fid/issue.go` (`renderIssueError`, currently no corruption/conflict arm)
- `internal/journal/replay.go` (`CorruptJournalError`, `EnsureProjectionFresh`,
  `inspectJournal`) and `internal/journal/project.go`
  (`PreconditionConflictError` fields)
- `internal/journal/replay_test.go` (`TestReplayJournalStopsOnPreconditionConflict`),
  `cmd/fid/rebuild_test.go` (existing corruption assertion)
- `06_DESIGN_PLAN_CLI.md` §18; `05_DESIGN_PLAN.md` §5.4

### Out Of Scope

- CF-003 (policy bootstrap corrupt-tail tolerance, `internal/policy`) — a
  separate code path with a later target (v0.1d); do not change it here.
- `fid merge --journal` and `fid event revert` (v0.2 conflict tooling).
- Any change to *when* corruption/conflict is detected in the replayer; this task
  is about faithful CLI reporting on both paths.
- Automatic intent reconciliation — explicitly not in v0.1 (PRD-OPS-012).

## Acceptance Criteria

- [ ] `fid rebuild` over a journal with a malformed trailing line emits
  `error: malformed journal entry at <file>:<line>` with the parse cause and
  manual-repair guidance (existing behaviour — keep covered). (VAL-008)
- [ ] `fid issue get` over a journal with a malformed trailing line stops and
  emits the same structured corruption error (no stale/partial read, non-zero
  exit). (VAL-008)
- [ ] `fid rebuild` over a journal with a failed effect precondition emits a
  structured conflict naming the event id/type, entity (table + pk), expected,
  and actual values. (VAL-009)
- [ ] `fid issue get` over a journal with a failed effect precondition stops
  before serving the read and emits the same structured conflict (non-zero exit).
  (VAL-009)
- [ ] A clean / satisfiable journal: both `fid rebuild` and `fid issue get`
  succeed (both states proven). (VAL-008, VAL-009)

## TDD Plan

- First failing test (`cmd/fid`): `fid rebuild` over a precondition-conflicting
  journal asserts the structured conflict message — red while
  `renderRebuildError` falls through to generic `error: %v`.
- Additional tests: read-path corruption and conflict rendering + non-zero exit;
  clean-path success for both commands; assert the conflict message includes
  event, entity, expected, actual.
- Regression risks: the read path must not serve stale data when catch-up fails
  (PRD-OPS-012 "stop before serving reads"); keep the structured fields out of
  stdout (errors go to stderr per 05_DESIGN_PLAN §5.3); hold COV-001
  `internal/journal` ≥ 80% and COV-002 `cmd/fid` ≥ 70%; keep `cmd/fid` thin
  (VAL-010) — formatting of the typed errors is presentation, but detection logic
  stays in `internal/journal`.

## Documentation Updates

- `project-docs/docs/phases/phase-3/tasks/task-5-completion.md` - completion
  report citing VAL-008/009 evidence (both rebuild and read-catch-up paths, both
  states) and gate outcomes.

## Scenario / Parity / Audit Notes

- Held-out scenario impact: the "journal corruption recovery: malformed trailing
  line" and "cross-branch journal merge with conflicting writes" scenario themes
  (05_DESIGN_PLAN §10) are exercised structurally here — do not read any held-out
  scenario content; implement to the spec only.
- Cross-surface parity impact: error message standard (05_DESIGN_PLAN §5.4) —
  errors name the constraint, entities, and corrective action.
- Audit or release-readiness impact: None.

## Exact Commands To Run

Run exactly:

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

The full task gate above is required before task close. Every required command
must be recorded in the completion report as `PASS`, or as `BLOCKED` with reason
in a blocker report. Missing, silently skipped, or narrower substituted commands
are failed gates.

Conditional commands:

```text
go build -o /tmp/fid ./cmd/fid
```

List any conditional command required by this task and why:

- `go build -o /tmp/fid ./cmd/fid` is required: this task changes user-facing
  error output for `fid rebuild` and `fid issue get`; confirm the binary builds
  and smoke both error paths.

## Notes For Developer

- Detection already lives in `internal/journal` (`CorruptJournalError`,
  `PreconditionConflictError`); this task adds faithful CLI *rendering* on both
  paths. Use `errors.As` to match the typed errors and render their fields.
- The read path must refuse to serve a read when catch-up fails — propagate the
  error and exit non-zero; never fall back to a stale projection.
- A shared render helper for the two typed errors avoids drift between
  `rebuild.go` and `issue.go`.
- Keep detection logic in `internal/journal`; only the message formatting belongs
  in `cmd/fid` (VAL-010).
