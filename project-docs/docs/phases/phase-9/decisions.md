# Phase 9 Decisions

> In-phase decisions (DEC-NNN / DEC-A..E). Per-doc spec-sync verdicts and the
> full decision history are finalized at phase wrap. Newest at top.

## DEC-024: Concern read ordering is deterministic by journal append order (`rowid DESC`)

Source: ESC-006 / final review benchmark-gate failure
(`TestConcernListFiltersByIssueAndStatus/all_concerns`), Phase 9 final review.
Status: Decided 2026-06-20 (Tech Lead); resolves an ESC-006 closeout blocker.

Decision: `fid concern list` (`internal/issues/concern_list.go`) and the
`fid issue get --section concerns` read model (`internal/issues/detail.go`
`issueConcerns`) order concern rows by **journal append order, newest first**
(`ORDER BY rowid DESC`) instead of `ORDER BY created_at DESC, id DESC`.

Rationale:
- `created_at` has second resolution, so two concerns raised in the same second
  tie; the `id DESC` tiebreak is a content hash, which orders rows
  nondeterministically relative to creation order. The phase benchmark gate
  surfaced this as a flaky `cmd/fid` test (`first concern issue_id` mismatch).
- `rowid DESC` is monotonic, append-ordered, and rebuild-stable: `fid rebuild`
  replays the journal in order, so post-rebuild `rowid` order equals original
  append order. VAL-009 rebuild parity is preserved.
- This mirrors the human-confirmed DEC-021 precedent (AC verification
  latest-per-role resolved by journal append order, no timestamp dependence).

Scope: ordering only. No schema change, no column change, no new public flag.
The newest-first contract is unchanged in intent; only the ordering key becomes
deterministic. `02_TRD.md` / `03_TRD_SCHEMA.md`: no schema change (the `concerns`
table is a normal rowid table). `06_DESIGN_PLAN_CLI.md` §6 gains an explicit
ordering note. Spec-sync verdicts finalized at wrap.

## DEC-023: Concern public read models expose `addressed_in_completion_id`

Source: ESC-006 / VAL-004 / QA-12 public read-back failure, Phase 9 final review.
Status: Decided 2026-06-20 (Tech Lead); resolves an ESC-006 closeout blocker.

Decision: both public concern read surfaces expose the addressing completion:
- `ConcernListItem` (`internal/issues/concern_list.go`) gains
  `addressed_in_completion_id` (nullable `*string`, JSON
  `"addressed_in_completion_id"`), selected and scanned from the existing
  `concerns.addressed_in_completion_id` column.
- `IssueConcern` (`internal/issues/detail.go`) gains the same nullable field, so
  `fid issue get --section concerns` and `fid concern list` stay parity-aligned
  (design §6: "row shape mirrors `fid issue get --section concerns` plus
  `issue_id`").

Rationale:
- The write path (`concern_address.go`) already persists the column on
  `issue.concern_addressed`; the lifecycle fact existed in the journal/DB but was
  invisible to public read-back, failing VAL-004 / PRD-DEV-006.
- The value is `null` for `open` concerns and the addressing `TC-` ID once
  addressed — the conditional/binary read-back both states VAL-004 requires.

Scope: read-model field addition only. No schema change (column exists since
`migrations/0002_v01_tables.up.sql`), no write-path change, no workflow-edge
change. `06_DESIGN_PLAN_CLI.md` §6 concern list JSON example updated to include
the field. Spec-sync verdicts finalized at wrap.

## Phase 9 planning decisions (DEC-A..E, human-confirmed 2026-06-20)

Recorded in the validation contract §Phase scoping notes. Summarized here for the
wrap audit; each task spec and QA artifact cites them by DEC-A..E.

- **DEC-A — PRD-WF-009 proven at the rule-evaluation seam, no default edge.**
- **DEC-B — concern domain lives in `internal/issues`; no new boundary.**
- **DEC-C — validate-before-append linkage rules (CF-007).**
- **DEC-D — address is a status mutation; raise→address only.**
- **DEC-E — no `activity_log` for concern writes (DEC-020 / CF-011 precedent).**

## Validation Contract Outcome Audit (at wrap)

The final review (`qa/final-review.md`) recorded VAL-001/002/003/005/006/007/008/
009/010 PASS, VAL-004 FAIL (public read-back omitted `addressed_in_completion_id`),
and VAL-011 BLOCKED (no human dogfood yet). The Task 4 fixup (commit `c927de0`,
QA Review 2 `qa/review-2.md` CONTINUE) resolved VAL-004 under DEC-023/DEC-024 and
cleared the benchmark gate, with VAL-008/009 re-confirmed PASS under the new
`rowid DESC` ordering. The human checklist (`qa/human-checklist.md`,
HUMAN-011a..e, 5/5 PASS, 2026-06-21) resolves VAL-011. **All eleven assertions now
stand PASS — no FAIL / BLOCKED / CARRY-FORWARD remains.**

| VAL | Outcome | Disposition |
|---|---|---|
| VAL-001 | PASS | `CON-7d528c` open concern with issue/completion/role/summary/body/status; CLI read-back (Task 1, QA-09). |
| VAL-002 | PASS | Four raise reject paths exit 1 with zero append; valid raise +1 line/+1 row (Task 1, QA-10, CF-007). |
| VAL-003 | PASS | Issue get/list `[false,0]`→`[true,1]` on real raise (Task 1, QA-09/11). |
| VAL-004 | PASS | **Resolved by Task 4 (DEC-023).** `addressed_in_completion_id` now exposed on `fid concern list` and `fid issue get --section concerns`; null for open, addressing `TC-` once addressed (Review 2; was FAIL at final review). |
| VAL-005 | PASS | Four address reject paths exit 1 with unchanged counts; valid address one mutation (Task 2, QA-13). |
| VAL-006 | PASS | Two open concerns `[true,2]`→`[true,1]`→`[false,0]` on successive address (Task 2, QA-14). |
| VAL-007 | PASS | `concerns_addressed` rule-evaluation seam unsatisfied(named IDs)→satisfied with real writes (Task 2, QA-15, DEC-A). |
| VAL-008 | PASS | Issue/status/combined filters in json/toon/markdown/yaml, empty=`[]` exit 0; re-confirmed under `rowid DESC` (Task 3 + Task 4, QA-16, DEC-024). |
| VAL-009 | PASS | Delete-DB + `fid rebuild` parity for concerns + derived flags; held under `rowid DESC` (Task 3 + Task 4, QA-17, ADR-0001). |
| VAL-010 | PASS | End-to-end v0.1c smoke ac add/verify/list→complete→raise(flag lights)→address(flag clears)→list (Task 3, QA-18). |
| VAL-011 | PASS | Human concern-lifecycle dogfood HUMAN-011a..e 5/5 PASS (was BLOCKED at final review). |

No assertion was weakened in-flight. VAL-004 stood as written and was satisfied by
exposing the read-model field rather than relaxing the assertion (the hard-stop
path: final review BLOCK → ESC-006 → DEC-023/024 → Task 4 fixup → QA re-verify).

## Spec Sync Verdict (at wrap)

Recorded for every spec. `verified-no-change` is a positive review assertion, not
a default skip.

| Spec | Verdict | Evidence / Reference |
|---|---|---|
| `00_PROJECT_CHARTER.md` | verified-no-change | Phase 9 delivered the in-scope concern write/read surface + v0.1c closure; no identity / scope / constraint change. |
| `01_PRD.md` | verified-no-change | PRD-DEV-005/006/007 and PRD-WF-009 MET as written (final review + Task 4). No requirement reworded; `escalated`/`wontfix` (PRD-DEV-006 extra statuses) remain deferred per the roadmap. |
| `02_TRD.md` | verified-no-change | Concern writes/reads use the existing write-path seam + append-only projection in `internal/issues` (DEC-B, no new boundary). DEC-E (no `activity_log`) matches the `create`/AC/completion precedent; CF-011 defers timeline population. DEC-A: no `concerns_addressed` edge added to the §4.6 default graph. |
| `03_TRD_SCHEMA.md` | verified-no-change | `concerns` table, `idx_concerns_issue_status`, `addressed_in_completion_id` column, and `issues_enriched` `concerns_raised`/`open_concerns_count` already shipped (v0.1a). DEC-023/DEC-024 are read-model/ordering changes only — no schema/column/trigger change. |
| `04_TRD_CONFIG.md` | N/A | Phase 9 added no `.fid/` layout or config surface. |
| `05_DESIGN_PLAN.md` | verified-no-change | `fid concern raise/address/list` follow the documented command map and mutation/read flows; no flow or shape change. |
| `06_DESIGN_PLAN_CLI.md` | touched: see DEC-023/DEC-024 | §6 concern commands updated during the phase: JSON shape carries `addressed_in_completion_id` (DEC-023) and an explicit "ordered by journal append order, newest first" note (DEC-024). Verified present at lines 769/778/781-782; Task 4 confirmed implementation matches. |
| `07_DESIGN_PLAN_FLOWS.md` | verified-no-change | Concern writes follow the journal→projection→read flow; VAL-009/010 confirm end-to-end determinism and the v0.1c closure flow. |
| `08_QUALITY_PLAN.md` | verified-no-change | COV-001 (`internal/issues` 81.8%, `internal/journal` 80.3%, `internal/policy` 85.0%, `internal/workflow` 95.5%) and COV-002 (`cmd/fid` 77.8%) held at close. New concern files within §4.4 file/function budgets. No plan reword. |
| `09_AI_WORKFLOW.md` | verified-no-change | Workflow ran as specified (contract → QA plan → dev/QA loop → final review BLOCK → escalation → fixup task → re-review → human checklist). The BLOCK→ESC→fixup→re-verify path is exactly the supported hard-stop behavior. |
| `10_IMPLEMENTATION_ROADMAP.md` | verified-no-change | Phase 9 delivered the §v0.1c "concerns + derived flag + v0.1c closure" cut and closes v0.1c; Phase 10 (§v0.1d claims/permissions/links) unchanged. |

## Triage-Log Audit (at wrap)

Two TRG entries in `triage-log.md`:

- **TRG-001 (RESUME) — Task 3 expanded write set.** Confirmed correct in
  hindsight. The `internal/format/markdown.go` change implemented the VAL-008
  required markdown output surface (task-local), and the test-only ordering
  adjustment removed a brittle assumption. Neither met an escalation trigger; the
  full gate passed. RESUME upheld.
- **TRG-002 (BLOCK → ESC-006) — final-review blocking findings.** Confirmed
  correct in hindsight. The benchmark-gate failure was a real required-gate
  failure and VAL-004/QA-12 was a real acceptance-contract failure; neither was
  misclassified, pre-resolved, nor covered by an accepted decision, so BLOCK (not
  autonomous RESUME) was the right call. It escalated to ESC-006, the Tech Lead
  issued DEC-023/DEC-024 + Task 4, and QA Review 2 verified the fix. No entry
  should have escalated that did not; none over-escalated.

Autonomous TL-Lite resume is permanently enabled (DEC-002); there is no warm-up
tally to maintain.

## Carry-Forward Disposition (at wrap)

- **No carry-forward items closed or opened by Phase 9.** The ESC-006 read-back /
  ordering defects were resolved in-phase by the Task 4 fixup (DEC-023/DEC-024),
  not carried forward.
- **CF-003 — re-deferred** to before v0.1d role/permission work (now imminent:
  Phase 10 is §v0.1d). Tech Lead should schedule the corrupt-tail reconciliation
  during Phase 10 planning.
- **CF-008 — re-deferred** to v0.2 (resolution type/reason product enhancement).
- **CF-011 — re-deferred** to v0.1e (consistent `activity_log` timeline with
  `fid activity list`). DEC-E kept concern writes consistent with the
  `create`/AC/completion precedent.
- **CF-013 — re-deferred** to v0.1e, alongside the `fid dev finish`/`block`
  composites that also touch `complete.go`.

## QA Recommendation Disposition (at wrap)

Final review verdict was **BLOCK** (2 blocking / 0 major / 3 improvement); the two
blocking findings were resolved by Task 4 and re-verified by QA Review 2
(CONTINUE, 0 blocking / 1 improvement). Improvement dispositions:

- **VAL row PASS labels / report top-level title omissions (review-1, 3 items).**
  **Acknowledge, no rework** — the deterministic gates and QA evidence independently
  prove the VAL outcomes. Carry the standing "finalize self-contained report
  evidence" nudge.
- **Full 40-char commit hash in completion reports (review-2, 1 item).**
  **Acknowledge, no rework** — `c927de0` verified uniquely in git for this review.
  Carry a standing process nudge: future completion reports record the full hash.
- **Refresh current docs + close ESC-006 in wrap artifacts (final-review
  follow-up).** **Done now** — `phase_brief.md`, `progress_current.md`,
  `decisions_active.md`, and `tl_escalations_current.md` refreshed in this wrap;
  ESC-006 moved to the phase archive.

## Workflow notes

- **What worked well.** The three-task vertical cut (raise / address + rule proof /
  list + rebuild + smoke) mirrored the proven Phase 7/8 shape and closed the whole
  v0.1c slice. Pre-phase human confirmation of DEC-A..E again removed mid-phase
  decision churn. When the final review caught two real defects, the
  BLOCK→ESC-006→DEC-023/024→Task 4 fixup→QA Review 2 loop worked exactly as the
  workflow intends: the defects were fixed under TDD with the benchmark gate as the
  regression guard, and the contract assertion (VAL-004) was satisfied rather than
  weakened. DEC-024 correctly reused the DEC-021 append-order precedent.
- **What was inefficient.** The two final-review blockers were detectable earlier:
  the missing `addressed_in_completion_id` read-back is a cross-surface parity gap
  that the mid-phase review window did not assert against (the field was persisted
  but never read back), and the nondeterministic concern ordering is the *same*
  class of defect already fixed for AC verifications in Phase 7 (DEC-021). Both
  could have been caught at Task 2/Task 3 review with a parity/ordering checklist
  carried from Phase 7.
- **Were tasks the right size?** Yes for Tasks 1-3. The unplanned Task 4 fixup was
  correctly bounded (read-model field + ordering, no schema/CLI change) and is the
  expected shape of an escalation-driven fixup, not a sizing miss.
- **Did QA cadence feel right?** Mostly. One mid-phase review over Tasks 1-3 plus a
  second over the isolated Task 4 fixup matched the work. The gap was *coverage
  depth*, not cadence: a recurring cross-surface read-back parity check and an
  append-order determinism check (both known from prior phases) should be standing
  QA-plan items so they are exercised before the phase-end gate.
- **Token-efficiency opportunities missed?** The full final-review → escalation →
  fixup → re-review cycle cost a phase-end round-trip that a Phase 7-derived
  ordering/parity checklist in the QA plan would likely have avoided. Settling
  DEC-A..E up front remained the big efficiency win.

### Deep retrospective recommendation

**Recommended — lightweight, cross-phase.** Phase 9 surfaced a *recurring*
determinism-ordering defect (DEC-024 is the third append-order fix after DEC-019
and DEC-021) and a cross-surface read-back parity gap that escaped mid-phase
review. These are now a cross-phase pattern, not one-offs. A short retro should
(1) promote "read ordering must be journal-append-order, not `created_at`" toward a
standing principle in `decisions_principles.md`, and (2) add a cross-surface
read-back parity check and an ordering-determinism check to the default QA plan so
future write-surface phases catch both before the phase-end gate. Defer a *deep*
multi-session retro; a focused note capturing these two standing checks is enough.
