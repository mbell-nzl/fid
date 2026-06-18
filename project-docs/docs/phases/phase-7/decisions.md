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

## Validation Contract Outcome Audit (at wrap)

Final review (`qa/final-review.md`) recorded VAL-001..VAL-009 PASS with
reproducible QA evidence and VAL-010 CARRY-FORWARD pending human dogfood. The
human checklist (`qa/human-checklist.md`, HV-010a/b/c, 3/3 PASS, 2026-06-14)
resolves VAL-010. **All ten assertions now stand PASS — no FAIL / BLOCKED /
CARRY-FORWARD remains.**

| VAL | Outcome | Disposition |
|---|---|---|
| VAL-001 | PASS | Ordered `ac add`; read-back ordinal order (Task 1). |
| VAL-002 | PASS | Epic rejection before append, zero journal growth (Task 1, CF-007). |
| VAL-003 | PASS | Multi-role independent latest-per-role (Task 2). |
| VAL-004 | PASS | Append-only fail→pass latest by `rowid` DESC, no `time.Sleep` (Task 2b, DEC-021). |
| VAL-005 | PASS | Field persistence + invalid-status rejection; `completion_id` NULL (Task 2). |
| VAL-006 | PASS | Latest-per-role `ac list` JSON, no raw history (Task 3). |
| VAL-007 | PASS | Four-format dispatch + JSON default (Task 3). |
| VAL-008 | PASS | Real-data `all_acs_verified` reject→verify→pass (Task 4). |
| VAL-009 | PASS | Delete-DB + `fid rebuild` AC/verification parity (Tasks 2b/4). |
| VAL-010 | PASS | Human AC-lifecycle dogfood HV-010a/b/c 3/3 PASS (was CARRY-FORWARD at final review). |

No assertion was weakened in-flight. VAL-004 was **strengthened** mid-phase
(Review 1 / ESC-005 → DEC-021: deterministic append-order proof, no
timestamp-boundary dependence) and proven by Task 2b. DEC-022 corrected a
task-scoping criterion (a quality tripwire, not a VAL); no assertion was
softened — the out-of-scope debt was tracked as CF-012.

## Spec Sync Verdict (at wrap)

Recorded for every spec. `verified-no-change` is a positive review assertion,
not a default skip.

| Spec | Verdict | Evidence / Reference |
|---|---|---|
| `00_PROJECT_CHARTER.md` | verified-no-change | Phase 7 delivered the in-scope AC write surface; no identity / scope / constraint change. |
| `01_PRD.md` | verified-no-change | PRD-AC-001..006, PRD-OUT-001..004, PRD-WF-007 MET as written (final review Requirement Compliance); PRD-DEV-* (completions/concerns) confirmed Phase 8/9, not reworded. |
| `02_TRD.md` | verified-no-change | AC writes use the existing write-path seam + append-only projection; no architecture or dependency change. DEC-020 (no `activity_log` for AC writes) matches the existing `create` precedent; CF-011 defers consistent timeline population. |
| `03_TRD_SCHEMA.md` | touched: see DEC-021 | `ac_verifications` comment now states latest-per-role resolves by journal append order for each `(ac_id, verifier_role)`, with `verified_at` recorded/displayed but not the ordering key (lines 237-244). No schema/column change. |
| `04_TRD_CONFIG.md` | N/A | Phase 7 added no `.fid/` layout or config surface. |
| `05_DESIGN_PLAN.md` | verified-no-change | AC add/verify/list follow the documented mutation/read flows; no flow or shape change beyond the CLI example below. |
| `06_DESIGN_PLAN_CLI.md` | touched: see Task 3 | §3 `ac list` JSON example updated to the implemented latest-per-role shape (Task 3); `ac add` / `ac verify` syntax verified to match implementation (final review QA-13). |
| `07_DESIGN_PLAN_FLOWS.md` | verified-no-change | AC writes follow the journal→projection→read flow; VAL-009 rebuild parity confirms end-to-end determinism. |
| `08_QUALITY_PLAN.md` | verified-no-change | COV-001 (`internal/issues` 82.0%, `internal/journal` 80.3%, `internal/policy` 85.0%, `internal/workflow` 95.5%) and COV-002 (`cmd/fid` 78.9%) held. §4.4 file-size: `ac.go` 190 / `ac_verify.go` 255 satisfy the tripwire; `issue.go` 677 tracked as CF-012 (pre-existing). No plan reword. |
| `09_AI_WORKFLOW.md` | verified-no-change | Workflow ran as specified (contract → QA plan → dev/QA loop with mid-phase reviews → final review → human checklist). The one BLOCK (TRG-001 / Review 1) routed through ESC-005 and cleared via DEC-021 / Task 2b. |
| `10_IMPLEMENTATION_ROADMAP.md` | verified-no-change | Phase 7 delivered the §v0.1c "acceptance criteria surface" cut; next cuts (Phase 8 completions, Phase 9 concerns) unchanged. |

## Triage-Log Audit (at wrap)

| TRG | Verdict | Hindsight |
|---|---|---|
| TRG-001 | BLOCK → RESOLVED (ESC-005 / DEC-021) | **Correct.** Same-second same-role latest-ordering nondeterminism is a persistence-contract semantics question spanning later `ac list` / workflow consumers (QA Reviewer triggers 5/6) and a plausible material spec deviation — escalating rather than resolving locally was right. The secondary Task 1 report-finalization FAIL was a genuine uncommitted-artifact blocker. Both cleared (DEC-021 + Task 2b; finalization commit verified by QA-15). No TRG should have escalated that did not; none over-escalated. |

Autonomous TL-Lite resume is permanently enabled (DEC-002); no warm-up tally to
maintain. Phase 7 had one BLOCK (TRG-001) and no autonomous RESUME entries.

## Carry-Forward Disposition (at wrap)

- **CF-011 — re-deferred** to v0.1e (consistent `activity_log` timeline
  population with `fid activity list`, PRD-OUT-017). Opened this phase by DEC-020.
- **CF-012 — re-deferred** to the next available v0.1c slot (behavior-preserving
  `internal/issues/issue.go` 677-line split, mirroring CF-010). Opened this phase
  by DEC-022; pre-existing debt, **not** Phase 7-introduced.
- **CF-003 — re-deferred** to before v0.1d role-permission / mutable-role work.
- **CF-008 — re-deferred** to v0.2 (resolution type/reason product enhancement).

## QA Recommendation Disposition (at wrap)

- **VAL-010 pending human dogfood (improvement).** **Done** — human ran
  HV-010a/b/c (3/3 PASS); VAL-010 resolved, human-checklist summary reconciled.
- **`issue.go` 677 lines above tripwire (improvement).** **Defer** — explicitly
  CF-012; not a Phase 7 blocker, scheduled behavior-preserving split.
- **Current-context docs need wrap-up sync (improvement).** **Done now** —
  `tl_escalations_current.md` (ESC-005 archived), `progress_current.md`, and
  `phase_brief.md` refreshed as part of this wrap.

## Workflow notes

- **What worked well.** Catching the same-second latest-ordering nondeterminism
  at Review 1 — before `ac list` (Task 3) and the `all_acs_verified` gate
  (Task 4) consumed it — and resolving it with the ADR-0001-aligned `rowid`
  append-order rule (DEC-021) kept the fix to one focused fixup (Task 2b) with no
  schema migration and no downstream rework. Splitting the AC write path into
  `ac.go` / `ac_verify.go` during that fixup kept both files under tripwire.
- **What was inefficient.** The Task 2b BLOCK came from an over-worded acceptance
  criterion ("every non-test file ≤400 lines") that swept in pre-existing
  `issue.go` debt the task never intended to touch. The Developer correctly
  stopped, but a tighter, scope-matched criterion (only the files the split
  produces) would have avoided a BLOCK and a decision (DEC-022). Lesson: scope
  file-size criteria to the task's allowed write set, not the whole package.
- **Were tasks the right size?** Mostly. The original four-task cut was clean;
  the inserted Task 2b fixup was the right granularity for an ESC-driven
  correction. Bundling the latest-ordering fix and the file split into one fixup
  was reasonable since both touched the same AC write-path files.
- **Did QA cadence feel right?** Yes. One BLOCK review + fixup re-review + a
  second mid-phase review across five task units matched the work; the human
  dogfood fanned cleanly into three per-observable rows.
- **Token-efficiency opportunities missed?** Minor. The Review 1 BLOCK plus
  fixup re-review cost an extra cycle, but it was load-bearing (it prevented a
  latent nondeterminism from shipping). The earlier-phase "paste exact
  transcript" nudge was honoured in task specs and reduced repeat findings.

### Deep retrospective recommendation

Not required. Phase 7 closed with zero blocking/major findings at final review,
all ten VALs PASS, a 3/3 human dogfood, and one clean BLOCK→DEC resolution
(DEC-021) plus one scope-correction DEC (DEC-022). The one mildly systemic
signal — over-broad acceptance criteria pulling in out-of-scope debt — is a
task-spec authoring tweak (scope file-size checks to the allowed write set), not
a retro-worthy cross-phase theme. Defer a deep retro until a recurring or
cross-phase problem surfaces.
