# Phase 16 - Task 5: `fid journal inspect <event-id>`

## Objective

Ship read-only journal event inspection (PRD-OPS-013): given an `event_id`,
locate the event in the journal and report its location, identity, actor,
command metadata, effects, canonical hash, and the event-type-keyed revert
verdict. Report-only — no compensating event is appended (PI-5; that is v0.2
`fid event revert`).

- New journal reader in `internal/journal`: find an event by `event_id` across
  journal files, returning the event + `journal_file` + `line` + canonical hash.
- New event-type-keyed revert classifier in `internal/journal` (reused/extended
  by v0.2 `fid event revert`): maps an event type to `{eligible, mode, restores,
  precondition, correction}`.
- New `cmd/fid/journal.go`: `fid journal inspect <event-id>` cobra command, all
  four formats, read-only (PI-1), with the documented JSON shape.

## Requirements

- PRD refs: PRD-OPS-013
- TRD refs: TRD-ARCH-001..004 (journal/event schema); `02_TRD.md` (journal
  package); `03_TRD_SCHEMA.md` §§9-12 (event/effect shape)
- Design refs: `06_DESIGN_PLAN_CLI.md` §16 `fid journal inspect`
- Quality refs: `08_QUALITY_PLAN.md` (deterministic evidence)

## Validation Contract Coverage

Assertions live in
`project-docs/docs/phases/phase-16/qa/validation-contract.md`.

- VAL-010: inspect reports event_id/schema/type, journal file+line, canonical
  hash (== `HashEvent`), actor, command, effects in all four formats; unknown id
  exits non-zero with an actionable error.
- VAL-011: revert verdict is event-type-keyed — mutable-state → inverse with
  restore+precondition; non-reversible → eligible:false with correction.

## Scope

### Blocker Addendum: Internal Activity Fixture Monotonicity (DEC-046)

Task 5 is resumed from the blocker report in `task-5-done.md`. Before closing
Task 5, fix the existing `internal/activity/activity_test.go` failure that
seeds a `2026-07-04` journal event after workspace initialization has created a
later journal file. That failure is out of the `fid journal inspect` feature
scope, but it blocks the required repo-wide gate and violates
`08_QUALITY_PLAN.md` TDD-006.

Allowed remediation scope:

- Edit only `internal/activity/activity_test.go` and directly local test helper
  code needed to keep seeded activity events monotonic with the workspace
  journal.
- Preserve VAL-001..004 behavior: filter include/exclude semantics,
  chronological ordering, structured `details`, date/RFC3339 boundary handling,
  and non-nil empty result behavior.
- Preserve PRD-ERR-014 behavior. Do not weaken the journal out-of-order guard,
  bypass production append paths, or narrow the Task 5 gate.
- Record the remediation in the final Task 5 completion report as a
  DEC-046-backed deviation from the original Task 5 file list.

### Files To Create Or Modify

- `internal/journal/inspect.go` (new) - `FindEvent(journalDir, eventID)`:
  iterate journal files in order (reuse `orderedJournalFiles`), scan lines,
  unmarshal events, return the matching `Event` + `journal_file` (`.fid/journal/
  <name>`) + `line` + `HashEvent` result; a not-found returns a typed
  `EventNotFoundError`. `RevertVerdict(event)`: event-type-keyed classifier
  returning `{Eligible, Mode, Restores, Precondition, Correction, Note}`.
- `internal/journal/inspect_test.go` (new) - find-by-id (existing + missing);
  hash equals `HashEvent`; revert verdict for a transition (inverse), a creation
  (`issue.created` → not eligible + correction), a claim/composite (not
  eligible), and an append-only fact (e.g. `issue.ac_verified` → retraction); a
  safe default for an unknown event type.
- `cmd/fid/journal.go` (new) - `newJournalCommand()` + `inspect` subcommand
  (`Args: cobra.ExactArgs(1)`, `--format`); assemble the documented JSON object
  (`event_id`, `event_schema`, `event_type`, `journal_file`, `line`, `hash`,
  `actor`, `command`, `effects`, `revert`); `renderIssueError` + non-zero on
  unknown id.
- `cmd/fid/journal_test.go` (new) - CLI: inspect a known event (all four
  formats), inspect a missing id (non-zero), and the two revert-verdict shapes.
- `cmd/fid/root.go` - register `newJournalCommand()`.

### Files To Read For Context

- `internal/journal/replay.go` (`orderedJournalFiles`, line scanning, file/line
  tracking, `Event` unmarshal)
- `internal/journal/canonical.go` (`HashEvent`, canonical JSON)
- `internal/journal/event.go` (`Event`/`Actor`/`Command`/`Effect` shapes)
- the event-type inventory (grep `EventType:` across `internal/`): `issue.*`,
  `issue.claim.*`, `link.*`, `label.added`, `requirement.synced`, `run.recorded`,
  `dev.*`, `project.initialized`
- `06_DESIGN_PLAN_CLI.md` §16 (exact JSON shape incl. `revert` object)
- `cmd/fid/run.go` (single-arg cobra + error-render pattern)

### Out Of Scope

- Appending a compensating event (`fid event revert`, v0.2).
- Activity-log / DB reads (inspect reads the journal directly, not the
  projection).
- Any write path or permissions row (read-only, PI-1).

## Acceptance Criteria

- [ ] `fid journal inspect <event-id>` for an existing event reports
  `event_id`/`event_schema`/`event_type`, `journal_file` + `line`, `hash`
  (equal to `journal.HashEvent`), `actor`, `command`, and ordered `effects` in
  all four formats. (VAL-010)
- [ ] An unknown event id exits non-zero with an actionable error. (VAL-010)
- [ ] A mutable-state event (`issue.transitioned`) reports `eligible:true,
  mode:inverse` with `restores` + `precondition`. (VAL-011)
- [ ] A non-reversible event (`issue.created` / `issue.claim.acquired` / a
  `dev.*` composite) reports `eligible:false` with a forward `correction`. (VAL-011)
- [ ] DEC-046 remediation complete: `internal/activity/activity_test.go` no
  longer seeds backdated journal events after workspace initialization, while
  Task 1 VAL-001..004 behavior remains covered. (quality blocker)

## TDD Plan

- First failing test: `FindEvent` for a known id returns the event + file/line +
  hash — fails red (function absent).
- Additional tests: missing id → `EventNotFoundError`; `hash == HashEvent`;
  revert verdicts (inverse for transition; retraction for ac_verified;
  not-eligible+correction for created/claim/composite; safe default for an
  unknown type); CLI all-format render + missing-id non-zero.
- Regression risks: hash must be computed over the **stored** event (pre-upgrade)
  exactly as replay does, so it matches what a future revert would reference;
  the classifier must be exhaustive over current event types with a documented
  safe default; functions ≤60 lines, files ≤400.

## Documentation Updates

- `06_DESIGN_PLAN_CLI.md` §16 - confirm the documented `revert` verdict buckets
  match the implemented classifier; adjust only if an event type is missing from
  the doc's enumeration (note any change in the completion report).
- Completion report: record the event-type → verdict table (PI-5) for the v0.2
  `fid event revert` follow-on.
- Completion report: record the DEC-046 activity-test fixture remediation and
  the exact file(s) touched for it.

## Scenario / Parity / Audit Notes

- Held-out scenario impact: None.
- Cross-surface parity impact: `hash` must equal `journal.HashEvent` (canonical
  parity); output shape must match `06_DESIGN_PLAN_CLI.md` §16.
- Audit or release-readiness impact: the revert classifier is the seed for v0.2
  `fid event revert`; keep it isolated and tested.

## Exact Commands To Run

Run exactly:

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

The full task gate above is required before task close. Record each command in
the completion report as `PASS`, or `BLOCKED` with reason. Record the
`internal/journal` and `cmd/fid` coverage percentages.

Conditional commands:

```text
go test -tags integration -race ./...
go build -o /tmp/fid ./cmd/fid
```

List any conditional command required by this task and why:

- `go test -tags integration -race ./...` — required if any new test is
  integration-tagged (the package set under test includes `cmd/fid`); record
  PASS/BLOCKED.
- `go build -o /tmp/fid ./cmd/fid` — required: CLI build used by the journal
  inspect tests.

## Notes For Developer

- PI-5: v0.1 only **reports** the verdict — never append a compensating event.
- Compute `hash` over the stored event exactly as `replay.go` does (`HashEvent`
  on the unmarshalled event before upgrade) so it matches replay diagnostics and
  a future revert reference (VAL-010 parity).
- The revert classifier must handle every current event type; map per the §16
  buckets (mutable-state → inverse; append-only fact → retraction;
  creation/completion/claim/telemetry/composite → not eligible + correction) and
  give an unknown type a safe `eligible:false` default.
- This is a read command — no `--by`, no matrix row (PI-1).
