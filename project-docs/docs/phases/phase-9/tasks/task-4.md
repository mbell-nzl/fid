# Phase 9 - Task 4 (fixup): concern read-back `addressed_in_completion_id` + deterministic ordering

## Objective

Clear the two ESC-006 phase-closeout blockers from the Phase 9 final review:

1. Expose `addressed_in_completion_id` on both public concern read surfaces so
   addressing a concern is visible in read-back (VAL-004 / QA-12, DEC-023).
2. Make concern read ordering deterministic by journal append order so the phase
   benchmark gate is green (DEC-024).

This is a bounded fixup task on top of Tasks 1-3; no new commands or schema.

## Requirements

- PRD refs: PRD-DEV-006 (concern lifecycle to `addressed` visible in read-back).
- TRD refs: §5.1 `concerns` schema (`addressed_in_completion_id` already exists,
  `migrations/0002_v01_tables.up.sql`); ADR-0001 (journal-canonical, DB derived).
- Design refs: `06_DESIGN_PLAN_CLI.md` §6 `fid concern list` / `fid concern
  address` (updated JSON shape + ordering note).
- Decision refs: **DEC-023** (read-model field), **DEC-024** (deterministic
  `rowid DESC` ordering, mirrors DEC-021).
- Quality refs: `08_QUALITY_PLAN.md` (TDD, coverage gate DEC-007).

## Validation Contract Coverage

Assertions live in `project-docs/docs/phases/phase-9/qa/validation-contract.md`.

- VAL-004 (currently FAIL): after `fid concern address`, read-back via
  `fid concern list --format json` and `fid issue get --section concerns
  --format json` must show `status='addressed'` **and** the addressing completion
  ID in `addressed_in_completion_id`; an `open` concern shows the field `null`.
- VAL-008 (currently PASS, must stay PASS): filters/formats unchanged; the
  deterministic ordering must not regress the empty-match or filter behavior.
- VAL-009 (currently PASS, must stay PASS): delete-DB + `fid rebuild` parity must
  still hold under `rowid DESC` ordering (append order is rebuild-stable).

## Scope

### Files To Create Or Modify

- `internal/issues/concern_list.go`:
  - Add `AddressedInCompletionID *string` (JSON `addressed_in_completion_id`) to
    `ConcernListItem`.
  - Add `addressed_in_completion_id` to the SELECT column list and scan it as a
    nullable string (via the existing `nullStringPointer` helper).
  - Change `ORDER BY created_at DESC, id DESC` to `ORDER BY rowid DESC` (DEC-024).
- `internal/issues/detail.go`:
  - Add `AddressedInCompletionID *string` (JSON `addressed_in_completion_id`) to
    `IssueConcern`, keeping field order parity with `ConcernListItem` intent.
  - Add the column to the `issueConcerns` SELECT and scan it as nullable.
  - Change the `issueConcerns` `ORDER BY created_at DESC, id DESC` to
    `ORDER BY rowid DESC` (DEC-024).
- `internal/issues/concern_list_test.go` / `internal/issues/issue_test.go` (or the
  existing concern detail test): assert the new field is `null` for `open` and the
  addressing `TC-` ID after address; assert deterministic newest-first ordering.
- `cmd/fid/concern_test.go`: keep `TestConcernListFiltersByIssueAndStatus`
  asserting newest-first ordering, now deterministic without timing dependence;
  add/extend a case asserting `addressed_in_completion_id` round-trips on the
  address path.
- `cmd/fid/e2e_v01c_concern_test.go`: extend the address step to assert
  `addressed_in_completion_id` is exposed on read-back if not already covered.

Keep `detail.go` near its size budget — add only the field + column + scan; do
not add new helpers there. New functions ≤60 lines, files ≤400 lines.

### Files To Read For Context

- `internal/issues/concern_address.go` (`issueConcernAddressEvent` already writes
  `addressed_in_completion_id`; confirm the column name/value).
- `internal/issues/detail.go` `nullStringPointer` usage precedent.
- `migrations/0002_v01_tables.up.sql` (`concerns` columns).
- Phase 7 DEC-021 (`project-docs/docs/phases/phase-7/decisions.md`) for the
  `rowid DESC` append-order precedent and rebuild-stability reasoning.

### Out Of Scope

- Any new CLI command, flag, or status (`escalated`/`wontfix` still have no
  producer; DEC-D).
- Schema/migration changes (the column already exists).
- Workflow-edge changes (DEC-A) and `activity_log` rows (DEC-E).
- Re-running or completing the VAL-011 human dogfood (tracked on the human
  checklist, not this task).

## Acceptance Criteria

- [ ] `fid concern list --format json` includes `addressed_in_completion_id`:
  `null` for an `open` concern, the addressing `TC-` ID after `fid concern
  address` (VAL-004, DEC-023)
- [ ] `fid issue get <issue> --section concerns --format json` includes
  `addressed_in_completion_id` with the same semantics (VAL-004, DEC-023)
- [ ] Both concern read queries order by `rowid DESC` (newest-appended first) and
  the ordering is deterministic regardless of `created_at` ties (DEC-024)
- [ ] `go test -bench . -benchmem ./...` exits 0 with
  `TestConcernListFiltersByIssueAndStatus` passing (ESC-006 benchmark gate)
- [ ] VAL-008 and VAL-009 still pass (filters/formats/empty-match unchanged;
  delete-DB + `fid rebuild` parity holds under the new ordering)

## TDD Plan

- First failing test: read-back assertion that `addressed_in_completion_id` is
  present and equals the addressing completion ID after `fid concern address`
  (fails today because the field is omitted).
- Second failing test: deterministic newest-first ordering for two concerns that
  share a `created_at` second (drive the same path the benchmark gate exercised).
- Additional tests: `addressed_in_completion_id` is `null` for `open`; rebuild
  parity snapshot includes the new field and survives delete-DB + rebuild.
- Regression risks: JSON null vs empty-string for `open`; field-shape parity drift
  between `ConcernListItem` and `IssueConcern`; ordering change must not break the
  VAL-008 filter expectations or the VAL-009 parity snapshot.

## Documentation Updates

- `06_DESIGN_PLAN_CLI.md` §6 — JSON example already updated by the Tech Lead to
  carry `addressed_in_completion_id` and the ordering note; confirm the
  implementation matches and adjust only if the field placement differs.
- Completion report — cite VAL-004/008/009, DEC-023, DEC-024, and record the
  benchmark gate now green; record coverage percentages.

## Scenario / Parity / Audit Notes

- Held-out scenario impact: None.
- Cross-surface parity impact: `fid concern list` and `fid issue get --section
  concerns` must expose the same `addressed_in_completion_id` field/semantics.
- Audit or release-readiness impact: this unblocks the Phase 9 final-review
  benchmark gate and VAL-004; note in the completion report for the TL wrap.

## Exact Commands To Run

Run exactly:

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
go test -tags integration -race ./...
go test -bench . -benchmem ./...
go build -o /tmp/fid ./cmd/fid
```

The benchmark command is **required** for this task (it is the failing phase
gate from ESC-006) — it must exit 0 with `TestConcernListFiltersByIssueAndStatus`
passing. Coverage must show `internal/issues` ≥ 80% and `cmd/fid` ≥ 70%
(DEC-007); record percentages in the completion report.

## Notes For Developer

- The write path already records `addressed_in_completion_id` on
  `issue.concern_addressed`; this task only adds it to the **read** path — do not
  change `concern_address.go`.
- Use the existing `nullStringPointer` helper for the nullable scan; an `open`
  concern must serialize the field as JSON `null`, not `""`.
- `ORDER BY rowid DESC` is rebuild-stable because `fid rebuild` replays the
  journal in order; do not introduce a `time.Sleep` to disambiguate ordering.
- Keep the two structs' field shapes parity-aligned so the cross-surface parity
  assertion in VAL-008/VAL-004 holds.
