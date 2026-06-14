# Phase 7 - Task 2b: Deterministic AC verification latest-per-role (ESC-005 / DEC-021 fixup)

## Objective

Resolve the Review 1 blockers on the `fid issue ac verify` slice:

1. Make "latest verification per role" deterministic and rebuild-stable by
   resolving it in **journal append order** (DEC-021), removing the dependence
   on second-resolution `verified_at` and on the random `ACV-` id that left
   VAL-004 unproven.
2. Replace the `time.Sleep`-based fail→pass tests with deterministic
   append-order coverage.
3. Commit the Task 1 completion-report finalization so the committed artifact
   contains the verified hash/subject (Review 1 Task 1 finalization FAIL).
4. Split `internal/issues/ac.go` (now 431 lines, over the 400-line tripwire)
   into cohesive files before Task 3 adds the `ac list` read path — behavior
   preserving (Review 1 improvement finding).

## Requirements

- PRD refs: PRD-AC-005 (append-only history; latest reflects most recent)
- TRD refs: TRD_SCHEMA §2.3 `ac_verifications` (append-only; latest-per-role)
- Design refs: `06_DESIGN_PLAN_CLI.md` §3 (no CLI shape change)
- Quality refs: DEC-021 (canonical ordering); DEC-007 coverage; 400-line non-test tripwire
- Decision: DEC-021 (`phase-7/decisions.md`); resolves ESC-005

## Validation Contract Coverage

- VAL-004: Append-only fail→pass yields two retained rows and latest-per-role =
  pass, proven deterministically (no `time.Sleep` boundary dependence).
- VAL-009: Delete-DB + `fid rebuild` reproduces identical latest-per-role after
  the ordering change (rebuild stability of the new tie-break).

## Scope

### Files To Create Or Modify

- `internal/issues/detail.go` - change `latestVerificationsByRole` latest
  selection to journal append order (`rowid DESC` per role, per DEC-021).
- `internal/issues/transition_rules.go` - audit `loadAcceptanceCriterionVerifications`
  / any other latest-per-role read so it uses the same append-order resolution.
- `internal/issues/ac_test.go` - replace sleep-based fail→pass test with
  deterministic coverage; add/extend rebuild latest-per-role assertion.
- `cmd/fid/issue_ac_test.go` - replace the sleep-based CLI fail→pass test.
- `internal/issues/ac.go` + new `internal/issues/<verify>.go` (or similar) -
  behavior-preserving split to bring each file ≤ 400 lines.

### Files To Read For Context

- `internal/issues/ac.go` (`VerifyAcceptanceCriterion`, `verified_at` write).
- `migrations/0002_v01_tables.up.sql` (`ac_verifications` has no append-order
  column; `id TEXT PRIMARY KEY` ⇒ implicit monotonic `rowid`).
- `project-docs/docs/phases/phase-7/decisions.md` (DEC-021).

### Out Of Scope

- Any change to `ac add`, the CLI flag/output shape, or stored `verified_at`
  format/value (it stays as recorded; only ordering semantics change).
- A schema migration or explicit sequence column, unless `rowid` proves
  unsuitable — if so, stop and escalate before adding a migration (DEC-021
  treats this as the persistence contract).
- `ac list` (Task 3) and transition gating (Task 4) behavior.

## Acceptance Criteria

- [ ] Latest-per-role is resolved by append/journal order (`rowid DESC`), not by
      `verified_at`/`id`, in every latest-per-role read path (VAL-004).
- [ ] A same-AC same-role `fail` then `pass` written in immediate succession
      (no `time.Sleep`) yields two rows and latest-per-role = pass; after the
      fail only, latest = fail (VAL-004).
- [ ] No remaining `time.Sleep`/timestamp-boundary dependence in AC
      verification ordering tests (`ac_test.go`, `issue_ac_test.go`).
- [ ] Delete-DB + `fid rebuild` reproduces identical latest-per-role for a
      multi-write fixture (VAL-009 stays green under the new ordering).
- [ ] Multi-role independence (VAL-003) and field persistence (VAL-005) remain
      green (no regression).
- [ ] The committed `task-1-done.md` contains the verified Task 1 commit hash
      and subject (no `Pending git commit`); confirm via
      `git show HEAD:project-docs/docs/phases/phase-7/tasks/task-1-done.md`.
- [ ] Every non-test file under `internal/issues` is ≤ 400 lines; the ac.go
      split is behavior-preserving (full gate green, coverage holds).

## TDD Plan

- First failing test: same-role fail→pass written with no sleep asserts
  latest-per-role = pass (fails today because the random-id tie-break can pick
  the fail row).
- Additional tests: rebuild latest-per-role parity for a multi-write fixture;
  keep multi-role independence and field-persistence coverage.
- Regression risks: the `all_acs_verified` rule loader and
  `issue get --section acceptance` must read latest through the same
  append-order resolution; do not change the loader's external contract.

## Documentation Updates

- `03_TRD_SCHEMA.md` `ac_verifications` - state that latest-per-role resolves in
  journal/append order (DEC-021). Adjust only this wording.
- `06_DESIGN_PLAN_CLI.md` §3 - confirm unchanged (no CLI shape change); note
  only if divergence is found.
- Completion report - paste the deterministic fail→pass transcript and the
  rebuild latest-per-role evidence; cite DEC-021 and ESC-005.

## Scenario / Parity / Audit Notes

- Held-out scenario impact: None (AC theme only).
- Cross-surface parity impact: `issue get --section acceptance`, `ac list`
  (Task 3), and the `all_acs_verified` loader must all agree on latest-per-role.
- Audit or release-readiness impact: None.

## Exact Commands To Run

Run exactly:

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

Coverage targets (DEC-007): `internal/issues` ≥ 80%, `cmd/fid` ≥ 70%.

Conditional commands:

```text
go build -o /tmp/fid ./cmd/fid
```

List any conditional command required by this task and why:

- `go build -o /tmp/fid ./cmd/fid` - required to regenerate the deterministic
  fail→pass and rebuild transcripts pasted in the completion report.

## Notes For Developer

- DEC-021 (`phase-7/decisions.md`) is the governing decision: latest = last
  appended in journal order. Resolve via SQLite implicit `rowid` DESC per role.
  `ac_verifications` is append-only and `fid rebuild` replays in journal order,
  so `rowid` order is deterministic and rebuild-stable. Keep `verified_at`
  exactly as recorded; it is display-only, not a sort key.
- If `rowid` is unsuitable for any reason, **stop and escalate** before adding a
  migration or sequence column — do not silently change the persistence
  contract.
- Do the Task 1 report finalization commit as its own commit (docs), separate
  from the code-behavior commit, to keep the history clean.
- Append-only is unchanged: two verifications still mean two rows; never update
  a prior row.
