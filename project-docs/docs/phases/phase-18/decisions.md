# Phase 18 Decisions

> Purpose: phase-local decision history, spec sync verdicts, workflow notes,
> and wrap-up decisions. DEC numbering is global and never resets.

## Decisions

### DEC-067: Resolve the Phase 18 QA-plan block (QA-BLK-18-01..03 / TRG-001)

- **Date:** 2026-07-19
- **Context:** The QA Lead's first Phase 18 QA plan was committed
  `BLOCKED BEFORE TASK 1` with three findings, all inside Tech Lead authority
  (not product/architecture escalations, per the QA plan's own note). The Phase
  Runner triaged this as TRG-001; TL-Lite returned `VERDICT: BLOCK (escalation
  required)` because rewriting the validation contract exceeds TL-Lite's local
  authority. This decision is the full Tech Lead resolution of that escalation.
  - **QA-BLK-18-01** — the draft contract bundled independent observables in
    VAL-005/007/009/013, left Task 5 read/ready/dev-composite parity and
    production-map isolation uncovered, and VAL-011 was narrower than the
    PRD-OPS-003/004 "reconstruct the query store" claim; VAL-008/012 evidence
    was under-specified.
  - **QA-BLK-18-02** — Task 4 said "recreate the issues indexes," but migrations
    0001–0006 define **no named index directly on `issues`**; the two named
    indexes (`idx_concerns_issue_status`, `idx_triage_issue_created`) are on
    `concerns`/`triage_events` and support `issues_enriched`. The required
    dependent-object set was therefore inexact.
  - **QA-BLK-18-03** — the benchmark command was valid, but scope covered only
    1k/10k-event rebuilds; there was no 100k-event rebuild (TRD-PERF-003) or
    10k-issue registry-backed read (TRD-PERF-001) spot-check for this
    persistence-contract migration.
- **Decision:**
  1. **Contract refined to Ready** (REF-18-01..09). Narrowed VAL-005 (gate JSON
     round-trip), VAL-007 (issues row/value preservation), VAL-009 (started
     entry: missing/wrong/stale/live), VAL-013 (golden replay). Refined VAL-008
     evidence to an exhaustive ordered-pair state matrix, widened VAL-011 to
     full-projection parity across two rebuilds, made VAL-012 both-state, and
     clarified VAL-001 PKs (TRD_SCHEMA §3 keys; `workflow_transition_requires`
     intentionally has no PK/UNIQUE). Added VAL-015 (reserved-name loader
     recognition), VAL-016 (seed exclusion), VAL-017 (phase-end behavior
     absence), VAL-018 (exact dependent-object inventory), VAL-019 (CLI read
     parity), VAL-020 (started auto-release), VAL-021 (journal-byte
     immutability), VAL-022 (get/list/ready parity), VAL-023 (dev-composite
     parity), VAL-024 (no production Go-map read). Contract now holds
     VAL-001..024, all Ready.
  2. **Task 4 dependent-object inventory made exact** (VAL-018). Verified from
     migrations 0001–0006: recreate the five triggers defined `ON issues`
     (`issues_parent_kind_guard_insert/_update`, `issues_kind_immutable`,
     `issues_sequence_scheduling_guard_insert/_update`) and the `issues_enriched`
     view; **preserve unchanged** (do not drop) the two named indexes on
     `concerns`/`triage_events` and the other-table kind/ref triggers that only
     reference `issues`; verify the `issues(id)` PK and `UNIQUE(parent_id,
     sequence)` autoindexes via PRAGMA. Full pre/post `sqlite_schema` fingerprint
     diff, stable across `.down`/`.up`.
  3. **Task 5 performance controls added** (QA-BLK-18-03). Rather than scope them
     out, added a 100k-event rebuild benchmark (<30s, TRD-PERF-003) and a
     10k-issue registry-backed read benchmark (<500ms, TRD-PERF-001), reusing
     `seedScaleReplayWorkspace`; the 10k-event rebuild (<5s, TRD-PERF-002)
     benchmark is retained. Rationale: this phase moves category/edge resolution
     onto per-read DB registry loads and changes the projection schema, so a read
     and a scaled-rebuild spot-check are genuinely warranted, and the fixture
     helper makes them feasible.
  4. Updated task specs 1–5 (coverage + acceptance) and the phase brief VAL
     index to match; contract marked **Ready**.
- **Consequences:** The three QA blocks are resolved on the planning side. On
  Phase Runner resume the QA Lead refreshes the QA plan's IDs/mappings (its
  section 2/3 "proposed VAL" rows become the now-active VAL-015..024) and clears
  the `BLOCKED BEFORE TASK 1` status; the QA plan's own section 10 already
  gates task start on that refresh plus Human approval of the Heavy Mode plan.
  No `ESC-NNN` is opened: contract refinement within Tech Lead authority does
  not meet a standing escalation trigger. Task execution still requires Human
  approval of the Heavy Mode phase plan (phase brief Status).
- **Refs:** `qa/validation-contract.md` (VAL-001..024, revision log);
  `qa/qa-plan.md` §§2–5 (QA-BLK-18-01..03, REF-18-01..09); tasks 1–5;
  `run-ledger.md` TRG-001; migrations 0001–0006; TRD_SCHEMA §3, §2.1, §14;
  02_TRD.md §§4.3, 4.6, TRD-PERF-001/002/003; DEC-066.

### DEC-068: Resolve Review 1 / TRG-002 VAL-020 rebuild-evidence gap (do-now remediation, no contract change)

- **Date:** 2026-07-19
- **Context:** Mid-phase Review 1 (tasks 1–3) returned `BLOCK` on one finding:
  Task 3's completion report marked VAL-020 `PASS`, but no test proves the
  caller claim remains absent **after rebuilding the exact released-claim
  fixture**. `TestTransitionOutOfStartedReleasesCallerClaimInSameEvent` proves
  only the same-event delete effect and immediate `Get()` absence; the general
  rebuild test (`internal/issues/transition_test.go:17`) rebuilds an
  *entering-started* fixture and therefore cannot prove persistence of the
  *release*. The Phase Runner triaged this as TRG-002; TL-Lite returned `BLOCK`
  but could not mint an `ESC-NNN` / triage-log entry, so the Phase Runner
  stopped and deferred to the Tech Lead. The QA Reviewer explicitly recorded the
  gap as localized to Task 3's own acceptance criterion, needing no
  cross-task / architecture / product / validation-contract decision.
- **Decision:**
  1. **Finding confirmed genuine.** VAL-020's pass criterion ("claim absent in
     the projection **and after rebuild**") and evidence method
     ("rebuilt-projection claim absence") were not met; the Task 3 report
     over-claimed `PASS`. The "existing rebuild check remains covered by the
     transition suite" line is not sufficient evidence for the released-claim
     scenario.
  2. **VAL-020 stands unchanged — no weakening.** The rebuild sub-assertion is
     correct and required. The remedy is to supply the missing evidence, not to
     adjust the assertion or the contract. No `VAL` edit; the validation
     contract stays **Ready**.
  3. **Do-now remediation within Task 3 scope (no new task, no re-plan).** The
     Developer redo, against the existing Task 3 spec (whose acceptance
     criterion already requires "absent ... after rebuild"), must: add a
     characterization test that enters a started-category state holding a live
     caller claim, transitions out of started, rebuilds the exact
     workspace/projection, and asserts the caller claim remains absent
     post-rebuild; correct the Task 3 completion report's VAL-020 evidence to
     cite that test (removing the vague claim); and rerun the full Task 3
     required gate. A QA Reviewer re-review of the VAL-020 remediation must pass
     before Task 4 starts.
  4. **No `ESC-NNN`.** This is a localized in-window test/evidence fix within
     Task authority; it meets no standing escalation trigger (no
     interface/contract/scope/cross-subsystem impact). TL-Lite's `BLOCK
     (escalation required)` was over-cautious; the inability to mint an ESC ID /
     triage log on `BLOCK` is the recurring process gap already tracked under
     CF-023.
  5. **Improvement finding (non-blocking, carry-forward).** The Task 2
     report-finalization commit (`7edfff6`) lacks its `Refs: phase-18/task-2`
     footer. History is immutable; do not rewrite it. Disposition: reminder to
     include the phase/task `Refs` footer on future documentation commits; folds
     into the existing report-lint commit-metadata friction item (CF-023
     deep-retro scope). It does not block the phase.
- **Consequences:** The Phase Runner resumes from the TRG-002 STOP: re-spawn a
  Developer to remediate Task 3 per (3), then a QA Reviewer re-review; on `PASS`,
  proceed to Task 4. Task specs and contract are unchanged except a clarifying
  remediation note appended to `tasks/task-3.md`. No canonical spec-suite doc is
  affected (VAL-020 semantics unchanged; the TRD §14.11 / §4.6 wording sync
  remains the wrap-time item already noted in the phase brief).
- **Refs:** `qa/review-1.md` (Task 3 Blocking finding; Improvement finding);
  `run-ledger.md` REVIEW review-1 / TRG-002 / STOP; `qa/validation-contract.md`
  VAL-020; `tasks/task-3.md` (acceptance criterion; Remediation note);
  `tasks/task-3-done.md` (VAL-020 evidence to correct); CF-023.

## Triage Dispositions

- **TRG-001 (2026-07-19) — RESOLVED by DEC-067.** Phase Runner triage of the
  blocked QA plan; TL-Lite `VERDICT: BLOCK (escalation required)`. Escalated to
  full Tech Lead because it required rewriting the validation contract. Resolved
  by DEC-067 (contract + task specs refined; no product/architecture escalation
  needed). TL-Lite did not mint an `ESC-NNN` or a `triage-log.md`; because the
  finding was correctly routed and is now resolved, no `ESC-NNN` is
  retroactively assigned. The Phase Runner should append a RESUME note to the
  run ledger referencing DEC-067 when it restarts. The recurring gaps
  (ESC-NNN global uniqueness; missing triage log on BLOCK) are already tracked
  under CF-023 for the deep retrospective.
- **TRG-002 (2026-07-19) — RESOLVED by DEC-068.** Review 1 `BLOCK` on the
  VAL-020 rebuild-evidence gap; TL-Lite `VERDICT: BLOCK` (no `ESC-NNN`/triage
  log minted). Confirmed a genuine, localized evidence gap: a do-now
  remediation within Task 3, VAL-020 unchanged, no `ESC-NNN`. Resume: Developer
  Task 3 redo (add the released-claim rebuild characterization test, correct the
  `task-3-done.md` VAL-020 evidence, rerun the full Task 3 gate) → QA Reviewer
  re-review of the VAL-020 remediation → Task 4. The over-escalation and
  missing-triage-log symptoms are the same CF-023 recurring gaps.
