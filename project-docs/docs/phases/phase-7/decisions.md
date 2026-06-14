# Phase 7 Decisions

> In-phase decisions (DEC-NNN). Per-doc spec-sync verdicts and the full
> decision history are finalized at phase wrap. Newest at top.

## DEC-022: Task 2b file-size criterion scoped to the ac.go split; issue.go tripwire deferred to CF-012

- **Date:** 2026-06-14
- **By:** Tech Lead (resolves the Task 2b BLOCK).
- **Status:** Confirmed.

### Context

Task 2b's last acceptance criterion read "Every non-test file under
`internal/issues` is ≤ 400 lines," but the task **objective (#4)** and **Scope**
only ever targeted splitting `ac.go` (431 lines) into cohesive files. The
Developer completed that split — `ac.go` = 190, new `ac_verify.go` = 255, both
≤ 400 — but the broadly-worded criterion swept in `internal/issues/issue.go`
at **677 lines**: a pre-existing tripwire violation untouched by Phase 7 and
outside Task 2b's allowed write set. The Developer correctly stopped at the
line-count check rather than make an out-of-scope edit (blocker report:
`tasks/task-2b-done.md`).

### Decision

The Task 2b file-size acceptance criterion is **scoped to the AC write-path
files produced by the split** (`ac.go`, `ac_verify.go`), matching the task
objective. `issue.go`'s 677-line tripwire violation is **pre-existing debt**,
tracked as **CF-012** for a dedicated behavior-preserving split (mirroring the
Phase 6 CF-010 transition/update/project split). No VAL assertion is affected:
VAL-004/VAL-009 are unchanged; the 400-line item is a quality **tripwire**
(DEC-007), not a validation assertion — so this is a task-scoping correction,
not a silent weakening of an assertion.

### Consequences

- Task 2b is unblocked: the Developer runs the full task gate for the in-scope
  split and creates the code-behavior commit plus the separate Task 1 report
  finalization commit.
- The out-of-scope debt is **explicitly tracked** (CF-012), not dropped.
- No schema/interface/persistence-contract change.

### Spec impact

- No spec touched by DEC-022 (task-scoping correction only). The seven-spec
  sync verdict for DEC-022 is `verified-no-change` at wrap.

## DEC-021: AC verification "latest per role" is resolved by journal append order

- **Date:** 2026-06-14
- **By:** Tech Lead (resolves ESC-005); **human-confirmed 2026-06-14** (persistence-contract semantics).
- **Status:** Confirmed — implementation scheduled via Task 2b.

### Context

ESC-005 / Review 1 found that `latestVerificationsByRole`
(`internal/issues/detail.go`) selects the latest verification per role with
`ORDER BY verified_at DESC, id DESC`. But `VerifyAcceptanceCriterion`
(`internal/issues/ac.go`) stores `verified_at` at `time.RFC3339` **second**
precision, and `ACV-` ids are content hashes, not append-order keys. Two
same-role verifications inside the same wall-clock second therefore tie on
`verified_at` and break the tie by **random** id, so a fast fail→pass can
report the wrong latest status. The Task 2 tests only pass because they
`time.Sleep` past a one-second boundary (`ac_test.go:654`,
`issue_ac_test.go:310`), so VAL-004 is not reliably proven and downstream
consumers (VAL-006 `ac list`, VAL-008 `all_acs_verified`) inherit the
ambiguity.

### Decision

**The canonical "latest verification per role" is the one appended last in
journal order** for that `(ac_id, verifier_role)`, not the one with the largest
timestamp or id. This follows ADR-0001 (journal is canonical; the DB is a
derived projection replayed in append order), so "most recent" must mean "last
journaled."

- Resolve latest-per-role by **SQLite implicit `rowid` DESC** (monotonic with
  insert/replay order): `... ORDER BY rowid DESC LIMIT 1` per role. `rowid` is
  strictly increasing in journal order because `ac_verifications` is an
  append-only derived table that is never deleted mid-life, and `fid rebuild`
  re-inserts rows in journal order — so the ordering is **deterministic and
  rebuild-stable**, independent of timestamp resolution and of the random
  `ACV-` id.
- `verified_at` remains a recorded/display field; it is **not** the ordering
  key. (Keeping it as a sort key would reintroduce the clock-skew/second-tie
  ambiguity.)
- An explicit monotonic sequence column is an acceptable alternative *within
  this contract* if the Developer finds `rowid` unsuitable, but `rowid` is the
  recommended minimal change (no migration).

### Consequences

- VAL-004 evidence must prove fail→pass latest-per-role **without** sleeping
  past a timestamp boundary (deterministic append-order proof).
- VAL-009 rebuild parity must continue to hold for latest-per-role after the
  ordering change.
- No schema migration; no public CLI shape change. Display of `verified_at` is
  unchanged.

### Spec impact (to confirm at wrap)

- `03_TRD_SCHEMA.md` `ac_verifications` "latest-per-role" wording: state
  append/journal-order resolution explicitly — **touched: see DEC-021** at wrap.
- Other six specs: expected `verified-no-change`; verify at wrap.
