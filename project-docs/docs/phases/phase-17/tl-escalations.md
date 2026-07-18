# TL Escalations - Phase 17

> **Scope.** Archived Phase 17 escalations, all resolved. Moved out of
> `../../current/tl_escalations_current.md` on 2026-07-17 during the VAL-024
> unblock, together with the WRAP-F1/WRAP-F2 register reconciliation.
>
> Three of the five entries below (ESC-005, ESC-011, and the ESC-004 closure)
> were reconstructed from `run-ledger.md` rather than carried from a live
> register entry. See **Register reconciliation** at the end of this file.

## ESC-004 — Task 1 completion report has uncommitted post-commit edits (RESOLVED)

- **Raised:** 2026-07-14 by TRG-004
- **Phase/task:** Phase 17 / Task 1
- **Trigger:** Quality-gate failure / report-lint evidence invalidation
- **Issue:** `tasks/task-1-done.md` passed its structural and recorded-evidence
  parser checks, but report-lint failed because metadata-wrap changes remained
  uncommitted after commit `65d2f94587310668b252ebfd4141ed9ef26c8ec0`
  (`refactor(format): split markdown renderer`).
- **Tech Lead disposition (2026-07-15):** Triage upheld. TRG-004 was correct to
  BLOCK — an uncommitted completion report has invalid report-lint evidence by
  definition, and the code commit `65d2f94` does not verify the report's current
  modifications. Resolution was mechanical with no spec, interface, contract,
  journal-schema, or architecture impact, so no DEC-NNN was required.
- **Resolution (2026-07-15):** Task 1 completed at `16:50` with report commit
  `2240785f426c3e3360f3f57981807ffa53bd2b7b`, verified subject, and report-lint
  PASS against the committed file. The ledger records "ESC-004 closed this run"
  at `2026-07-15 16:50` (task-1) and again at `17:49` (task-2).
- **Status:** RESOLVED. Archived 2026-07-17 (WRAP-F1); the entry had remained
  stale-open in the current register after its actual closure.

## ESC-005 — Task 2 completion report has bare commit metadata (RESOLVED)

- **Raised:** 2026-07-15 by TRG-005; ledger STOP `2026-07-15 16:59`
- **Phase/task:** Phase 17 / Task 2
- **Trigger:** Quality-gate failure / report-lint evidence invalidation
- **Issue:** `tasks/task-2-done.md` failed report-lint because the commit
  hash/subject section used bare text instead of the backtick-wrapped shape the
  linter parses, with no fresh committed evidence.
- **Resolution (2026-07-15):** Task 2 completed at `17:49` with commit
  `13facdf9616e705db5959634a08c0aa23a23ec0f`, verified subject, and report-lint
  PASS. Benign, same root cause as ESC-004.
- **Status:** RESOLVED. **Backfilled 2026-07-17 (WRAP-F2)** — this escalation was
  never written to the current register when raised, so it was under-reported
  from raise to close.
- **ID collision:** this ID duplicates Phase 7's ESC-005 (AC verification latest
  ordering) in `../phase-7/tl-escalations.md`. The ID was assigned by TL-Lite at
  raise time and is preserved here as recorded in the ledger rather than
  renumbered after the fact. See **Register reconciliation**.

## ESC-011 — Task 7 completion report has uncommitted post-commit edits (RESOLVED)

- **Raised:** 2026-07-15 by TRG-007; ledger TRIAGE and STOP `2026-07-15 19:42`
- **Phase/task:** Phase 17 / Task 7
- **Trigger:** Quality-gate failure / report-lint evidence invalidation
- **Issue:** `tasks/task-7-done.md` failed report-lint because edits remained
  uncommitted after commit `34c6cef754ab34f3d483ccee311dd3b05c3b2337`
  (`feat(cli): guard unshipped governance commands`).
- **Resolution (2026-07-15):** Task 7 completed at `20:37` with verified commit
  `34c6cef754ab34f3d483ccee311dd3b05c3b2337` and report-lint PASS. Benign, same
  root cause as ESC-004/ESC-005.
- **Status:** RESOLVED. **Backfilled 2026-07-17 (WRAP-F2)** with the new ID
  ESC-011. The ledger records `Verdict: BLOCK` with the note "TL-Lite could not
  record an escalation ID in its output", so this BLOCK ran and closed with no
  durable identifier. ESC-011 is assigned retroactively at archive time; it is
  numbered after ESC-010 despite being chronologically earlier, because IDs are
  labels and renumbering the register would invalidate existing ledger
  references.

## ESC-009 — Phase 17 QA route model-qualification (RESOLVED)

- **Resolved 2026-07-14** by TRG-001: the committed model audit
  (`qa/model-audit.md`, commit `2d5e536`) records `PASS - QA routing qualified`
  with no MISS, closing QA-BLK-02 and satisfying DEC-053/055. Task 0 and the
  Entry Review proceeded on that basis.
- **Status:** RESOLVED. Archived 2026-07-17.

## ESC-010 — Task 9: shipped CLI rebuild fails derived-projection parity (RESOLVED)

- **Raised:** 2026-07-15 by TRG-008; source `tasks/task-9-done.md`
- **Status:** Resolved 2026-07-17 by DEC-056 — premise not accepted; Task 9
  resumed audit-only behind a required diagnostic and closed on Branch A.
- **Escalation trigger:** Changes a data model or persistence contract; affects
  more than one subsystem; changes testing/validation strategy.
- **Issue:** After the derived DB is deleted, `fid rebuild --by workflow_runner`
  preserved the journal byte-for-byte but appeared to rebuild a different
  `activity_log`. `fid maintenance check` reported `recent_activity_7d` as 17
  before rebuild and 20 after. The focused journey test failed, VAL-022 was
  BLOCKED, and required Task 9 gates remained BLOCKED.
- **Evidence at raise:** none attached; the report recorded the focused journey
  failure and blocked gate status.
- **Routing:** Full Tech Lead

### Tech Lead disposition (2026-07-17)

**TRG-008 was right to BLOCK.** A reported persistence-parity failure with no
verification evidence is exactly what must stop the phase. The escalation was
correct; its *diagnosis* is what was not accepted.

**1. Canonical rebuild-parity expectation — recorded as DEC-056.** The journal is
the sole source of truth for derived state: after delete-and-rebuild every
projected table must be exactly equal to its pre-delete content. That is not
relaxed. But command output is not derived state. A read whose value is a
function of `(projected state, wall-clock now)` is not a parity observable,
because `now` is an input to the read that the journal does not reproduce.
`maintenance check` carries two such fields — `recent_activity_7d` (a
`now`-anchored window) and `claim_count_live` (TTL liveness evaluated at `now`).

**2. The reported defect is very likely not a defect.** Code review did not
support "rebuild reconstitutes different `activity_log` state":

- Live writes and replay run the identical path — `journal.ProjectEvent` applying
  the same effects from the same journal bytes. No production code writes or
  deletes `activity_log` outside that path, `activity_log` has no FK/cascade, and
  activity IDs and `occurred_at` are carried *in* the journal effect, so replay
  cannot invent or lose them.
- There is a mechanism that fully explains 17 -> 20 with the rows identical.
  `nextIssueUpdateTime` (`internal/issues/write_helpers.go`) preserves timestamp
  monotonicity by returning `prior + 1s` when the wall clock hasn't advanced past
  the previous `updated_at`. The journey fires many issue mutations inside the
  same second, so the last few activity rows are stamped slightly **ahead** of
  the wall clock. `countRecentActivity` (`internal/maintenance/info.go`) excludes
  future rows via `!at.After(now)`. The pre-rebuild check omits them; the
  post-rebuild check runs seconds later — after a DB delete, a rebuild, and 14
  `fid sql` subprocesses — by which time the clock has caught up and the same
  rows count. Same rows, later `now`.

**3. The assertion, not the product, is the more likely fault.** Task 9's own
binding constraint said to compare "structurally without timestamp-string
fragility". The implemented `closureSnapshot` embedded `maintenance check` output
— including a `now`-relative counter — inside the compared value. As written,
VAL-022 asserted equality of a time-varying function across two invocations
separated by real elapsed time; no correct implementation could satisfy it, and
it could not distinguish a rebuild defect from clock advance.

**4. Required next action — diagnose before remediating.** No product change was
authorized on a hypothesis, including the Tech Lead's own. Task 9 carried a
required diagnostic to decide the branch: run `maintenance check` twice at the
end of the journey, seconds apart, **with no rebuild and no DB deletion**.

- **Branch A — the control reproduces the movement.** Rebuild is exonerated. Adopt
  the amended VAL-022 method (exact per-table equality; exact equality on
  `status`/`errors`/`warnings`/`issue_count`; clock-relative counters excluded
  from equality). Task 9 stays audit-only and closes. No product change.
- **Branch B — the control does not reproduce it** and `activity_log` rows
  genuinely differ before-vs-after. Real persistence defect, outranks the closure
  audit. Stop, do **not** amend VAL-022, re-escalate with the row diff attached;
  a scoped remediation task gets scheduled and VAL-022 stays BLOCKED.

**5. Dispositions.**
- VAL-022: BLOCKED -> **amended and Ready**, pending the diagnostic. Recorded as a
  validation-contract deviation in the contract's revision log with reasoning; it
  is a narrowing and was not recorded as anything else. The parity intent is
  preserved and in fact strengthened — per-table exact equality is stricter than
  the previous aggregated-snapshot comparison.
- VAL-023: the rebuild-parity portion unblocks with VAL-022 under Branch A.
- Task 9 gates: remained BLOCKED until the task closed normally. The report's
  BLOCKED gate entries were honest and correctly recorded.
- **CF-020 (new, v0.2):** future-stamped `occurred_at` under burst writes, and its
  interaction with the future-excluding `recent_activity_7d` window. Real but
  minor, bounded, and self-correcting; it is a deliberate monotonicity trade-off
  and does not block v0.1 closure. Not fixed inside an audit task.

**Spec impact:** none pending the diagnostic. Under Branch A no shipped behavior
changes, so no PRD/TRD/design edit is required; DEC-056 plus the contract
revision log carry the record. Under Branch B, spec impact is reassessed with the
remediation task.

### Outcome (2026-07-17)

**Branch A confirmed.** The ledger records Task 9 complete at `19:47`, commit
`7119d5b38d79663f544f093d55f676c375746419`, report-lint PASS, with the
diagnostic result: "the no-rebuild control reproduced `recent_activity_7d`
movement (`17` then `20`), so the amended parity method was used and no product
code changed". QA Lead's final review records VAL-022 and VAL-023 PASS on that
evidence. Rebuild was exonerated; the assertion was the fault.

## Register reconciliation (WRAP-F2, 2026-07-17)

The register was reconciled against `run-ledger.md` for the whole phase. Findings
carried into the phase-end triage audit and the retrospective:

1. **ESC-005 was raised and closed without ever being written to the register.**
   Recovered above from the ledger STOP at `2026-07-15 16:59`.
2. **The TRG-007 escalation ran with no ID at all**, because TL-Lite could not
   record one in its output. Backfilled above as ESC-011.
3. **ESC-004 was stale-open** in the current register long after the ledger
   recorded it closed twice.
4. **ESC-NNN numbering is not globally unique, contrary to
   `templates/tl-escalations-template.md` ("Use global ESC-NNN numbering that
   never resets").** Discovered during this reconciliation, beyond WRAP-F2's
   original scope. Actual collisions across the project: ESC-001 (phase-1 and
   phase-14), ESC-004 (phase-6 and phase-17), ESC-005 (phase-7 and phase-17).
   The register cannot be indexed by ID alone without a phase qualifier. Not
   remediated here — renumbering would invalidate ledger and review references.
   **Retrospective/backlog item for the Tech Lead.**

Both under-reported escalations turned out benign, so nothing was lost this
phase. The concern is structural: the register could not be trusted as the
complete list of what blocked the phase, and the phase-end triage audit walks the
register.

## Retrospective candidate

Five of eight Phase 17 triages (TRG-002/003/004/005/007) were the *same*
report-lint commit-metadata defect, and two of them escalated (ESC-004, ESC-005;
ESC-011 makes three). That is a recurring mechanical failure worth a
deep-retrospective line: the completion-report template or report-lint should
make the backtick-wrapped commit metadata impossible to get wrong, rather than
catching it five times. QA Lead's final review records the same item as an
Improvement finding.
