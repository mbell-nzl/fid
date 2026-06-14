# Phase 7 - Task 1: `fid issue ac add`

## Objective

Implement `fid issue ac add <issue-id> --text <text> --ordinal <n> --by <role>`:
add an ordered acceptance criterion to a **task** or **subtask**, reject the
operation on an **epic**, journal it as a single `issue.ac_added` event, and
project it into `acceptance_criteria`. This is the first v0.1c write command.

## Requirements

- PRD refs: PRD-AC-001 (ordered ACs on tasks/subtasks), PRD-AC-002 (reject on epic)
- TRD refs: TRD_SCHEMA §2.3 `acceptance_criteria`; §1 ID scheme (`AC-<hash>`)
- Design refs: `06_DESIGN_PLAN_CLI.md` §3 `fid issue ac add`
- Quality refs: `08_QUALITY_PLAN.md` validate-before-append (CF-007); DEC-007 coverage

## Validation Contract Coverage

- VAL-001: `ac add` on task/subtask creates an `AC-<hash>` row with text + ordinal; read-back is ordinal-ordered.
- VAL-002: `ac add` on an epic is rejected before any journal append; task/subtask succeeds.

## Scope

### Files To Create Or Modify

- `internal/issues/ac.go` (new) - `AddAcceptanceCriterion` write path (validate → lock → generate `AC-` id → build event → dry-run project → append → project).
- `internal/issues/ac_test.go` (new) - unit/integration tests.
- `cmd/fid/issue.go` - wire `newIssueACCommand()` (parent `ac`) + `ac add` subcommand and stderr confirmation.
- `cmd/fid/issue_ac_test.go` (new) - CLI tests.

### Files To Read For Context

- `internal/issues/issue.go` (`Create` write-path pattern: lock, `GenerateID`, effects, `DryRunProjectEvent`, `ProjectEvent`).
- `internal/issues/detail.go` (`issueAcceptanceCriteria`, read-back shape).
- `internal/journal/project_effects.go` (insert effect contract).
- `internal/db/kind_triggers_test.go` (AC kind-guard backstop behavior).

### Out Of Scope

- AC verification (Task 2), `ac list` command (Task 3), transition gating (Task 4).
- `activity_log` emission for AC writes (DEC-020, confirmed: AC writes do not log to the derived timeline table in Phase 7; consistent with `fid issue create`; deferred to v0.1e via CF-011).
- Ordinal uniqueness enforcement (no schema constraint; ties break by id).

## Acceptance Criteria

- [ ] `fid issue ac add <task-id> --text T --ordinal 1 --by developer` creates one `AC-<hash>` row with text=T, ordinal=1 and emits exactly one `issue.ac_added` event (VAL-001).
- [ ] Multiple ACs on one issue read back ordered by ordinal then id via `fid issue get --section acceptance` (VAL-001).
- [ ] `ac add` on an **epic** exits non-zero with a clear kind-constraint error and appends **zero** journal lines / rows (validate-before-append, VAL-002).
- [ ] `--text`, `--ordinal`, `--by` are required; `--ordinal` must be a positive integer; unknown `--by` role rejected via `policy.RequireKnownRole` before append (VAL-002).
- [ ] Issue-not-found exits non-zero with a clear error and zero append.
- [ ] stderr confirmation matches design §3 (`added AC AC-… on <issue> (ordinal N): "…"`).

## TDD Plan

- First failing test: `ac add` on a task projects one `acceptance_criteria` row with the given text/ordinal and an `AC-` id.
- Additional tests: epic rejection (zero append); missing/invalid ordinal; unknown role; issue-not-found; multi-AC ordinal ordering; one-event-per-write shape.
- Regression risks: kind detection must reuse existing helpers; do not regress `issue get --section acceptance` read-back.

## Documentation Updates

- `06_DESIGN_PLAN_CLI.md` §3 - confirm the documented `ac add` output/flags match the implementation; adjust only if behavior differs.
- Completion report - paste exact command transcripts for the add (task + epic-reject) paths (Phase 6 process nudge: pasted transcripts required for public-output changes).

## Scenario / Parity / Audit Notes

- Held-out scenario impact: None (AC theme only; do not paste hidden scenario content).
- Cross-surface parity impact: read-back parity between `ac add` and `issue get --section acceptance`.
- Audit or release-readiness impact: None.

## Exact Commands To Run

Run exactly:

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

Coverage targets (DEC-007): `internal/issues` ≥ 80%, `cmd/fid` ≥ 70%. The full
task gate above is required before task close. Record each command as `PASS` (or
`BLOCKED` with reason) in the completion report.

Conditional commands:

```text
go build -o /tmp/fid ./cmd/fid
```

List any conditional command required by this task and why:

- `go build -o /tmp/fid ./cmd/fid` if pasting real CLI transcripts for evidence.

## Notes For Developer

- Match the `Create` write-path exactly: acquire `journal.WithLock`, open DB, generate id with `journal.GenerateID("AC", existsFn)`, build the event with one `insert` effect into `acceptance_criteria`, call `journal.DryRunProjectEvent` before `AppendEvent`, then `ProjectEvent`.
- Event type: `issue.ac_added` (matches the past-tense `issue.created`/`issue.updated` convention).
- Reject epics in the engine **before** building/appending the event (the AC kind-guard trigger is only a backstop; rejecting late would poison the journal — CF-007).
- Keep the new file under the 400-line size tripwire.
</content>
