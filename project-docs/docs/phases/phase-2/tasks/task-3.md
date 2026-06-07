# Phase 2 - Task 3: Kind-Invariant And Kind-Immutability Backstop Triggers (Migration 0004)

## Objective

Add the v0.1 kind-invariant and kind-immutability SQLite triggers as a backstop
for non-CLI write paths (future dashboard, `fid sql`, direct DB edits). These
enforce structural rules that static CHECKs cannot, because they require
cross-row reads of the referenced issue's `kind`. The application layer remains
the primary validator; these triggers are the defense-in-depth backstop.

## Requirements

- PRD refs: N/A (structural backstop; no user-facing command in this task)
- TRD refs: TRD §4.5 (application-primary, triggers as backstop), TRD SCHEMA
  §14.1 (kind-invariant rules), §14.7 (kind immutability); TRD-ARCH-006,
  TRD-SEC-004
- Design refs: N/A
- Quality refs: Quality Plan §5.1; TRD-SEC-004 verification (raw-SQL RAISE)

## Validation Contract Coverage

Assertions live in `project-docs/docs/phases/phase-2/qa/validation-contract.md`.

- VAL-009: `issues.parent_id` kind-hierarchy trigger (both states).
- VAL-010: kind-guard triggers on `acceptance_criteria`, `task_completions`
  (task or subtask) and `epic_metadata` (epic).
- VAL-011: `epic_id` kind-guard triggers on `triage_events` and `decisions`.
- VAL-012: kind-immutability trigger on `issues` (BEFORE UPDATE OF kind).

## Scope

### Files To Create Or Modify

- `migrations/0004_v01_kind_triggers.up.sql` (kind-hierarchy on `issues`;
  kind guards on `acceptance_criteria`, `task_completions`, `epic_metadata`,
  `triage_events`, `decisions`; kind immutability on `issues`)
- `migrations/0004_v01_kind_triggers.down.sql` (drop all triggers added)
- `internal/db/*_test.go` (raw-SQL backstop tests: invalid RAISE + valid insert,
  both states for each trigger; confirm valid journal-applied issue still
  projects)
- `internal/db/projection_meta.go` (bump `ProjectionSchemaVersion` if a fresh
  generation is warranted; see Notes)
- `project-docs/docs/phases/phase-2/tasks/task-3-completion.md`

### Files To Read For Context

- `project-docs/specs/03_TRD_SCHEMA.md` §14.1, §14.7
- `project-docs/specs/02_TRD.md` §4.5
- `migrations/0001_v01a_base.up.sql` (issues kind values),
  `migrations/0002_v01_tables.up.sql` (child tables shipped in Task 1)
- `project-docs/docs/phases/phase-2/qa/validation-contract.md`

### Out Of Scope

- Polymorphic ref integrity, link scope, claim scope (Task 4).
- Sequence-vs-scheduling backstop trigger (deferred — see contract Scope
  Deferrals).
- App-layer hierarchy validation in `fid issue create` lifecycle (v0.1b).
- `carry_forward.scheduled_epic_id` kind guard (carry_forward is v0.4).

## Acceptance Criteria

- [ ] BEFORE INSERT/UPDATE trigger(s) on `issues` reject: an epic with a
  non-NULL `parent_id`; a subtask whose `parent_id` is NULL or does not
  reference a `task`; a task whose `parent_id` is non-NULL and does not
  reference an `epic`. Valid epic(no parent)/task(epic parent or NULL)/
  subtask(task parent) insert successfully. (VAL-009)
- [ ] BEFORE INSERT trigger(s) reject `acceptance_criteria` and
  `task_completions` whose `issue_id` references an issue that is not `task` or
  `subtask`; reject `epic_metadata` whose `issue_id` is not an `epic`. Valid-kind
  rows insert. (VAL-010)
- [ ] BEFORE INSERT trigger(s) reject `triage_events` whose `epic_id` does not
  reference an `epic`, and `decisions` whose `epic_id` is set but does not
  reference an `epic`. Valid references (and NULL `decisions.epic_id`) insert.
  (VAL-011)
- [ ] BEFORE UPDATE OF kind trigger on `issues` rejects any change to `kind`;
  an update to a non-kind column on the same row succeeds. (VAL-012)
- [ ] All triggers raise a clear `RAISE(ABORT, ...)` message; a valid issue
  created via the journal/projection path still projects without trigger error.

## TDD Plan

- First failing test: raw-SQL `INSERT INTO issues` with kind `epic` and a
  non-NULL `parent_id` should RAISE — red before the trigger exists.
- Additional tests: one invalid + one valid case per rule above (parent_id
  hierarchy across all three kinds; ac/completion/epic_metadata kind guards;
  triage/decision epic_id guards; kind immutability vs non-kind update); a
  control test that a valid `issue.created` event still projects under the new
  triggers.
- Regression risks: `recursive_triggers = ON` interaction; trigger WHEN clauses
  for NULL `parent_id`; matching `RAISE` on both INSERT and UPDATE where §14.1
  specifies it; not breaking the existing Phase-1 issue-create projection.

## Documentation Updates

- `project-docs/docs/phases/phase-2/tasks/task-3-completion.md` - completion
  report citing VAL-009..012 evidence, gate outcomes, and a note recording the
  §14.1-vs-§4.1 `task_completions` kind resolution (task or subtask) actually
  implemented.

## Scenario / Parity / Audit Notes

- Held-out scenario impact: None. Do not read held-out scenarios.
- Cross-surface parity impact: None.
- Audit or release-readiness impact: Triggers satisfy TRD-SEC-004 backstop
  verification; ensure raw-SQL RAISE evidence is captured for the security review.

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
in a blocker report.

Conditional commands:

```text
go test -tags integration -race ./...
go build -o /tmp/fid ./cmd/fid
```

List any conditional command required by this task and why:

- `go test -tags integration -race ./...` if trigger backstop tests run with
  temp DBs behind the `integration` build tag.

## Notes For Developer

- Triggers are a **backstop**, so test them with raw SQL inserts (the non-CLI
  write path), per TRD-SEC-004 / TRD-ARCH-006. Also keep one control test that a
  valid journal-applied row still projects, so a trigger never false-positives on
  the happy path.
- `task_completions` kind guard: implement per §14.1 ("same" as
  `acceptance_criteria` = **task or subtask**). TRD §4.1's "(task only)" is the
  indicative entity note (§4.1 says column lists are indicative; §14.1 is the
  trigger reference). Record this resolution in the completion report so QA can
  confirm it at spec sync.
- Use `RAISE(ABORT, '<table>: <reason>')` style messages naming the rule, so the
  backstop is debuggable.
- Only add the `ProjectionSchemaVersion` bump if you judge a fresh projection
  generation is warranted; triggers do not change projected row *data*, so a bump
  is optional here. Coordinate with Tasks 1-2 to keep the constant monotonic and
  document the choice.
- Keep all trigger DDL in migration `0004`; do not edit `0001`-`0003`.
