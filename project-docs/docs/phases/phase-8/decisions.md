# Phase 8 Decisions

> In-phase decisions (DEC-NNN / DEC-A..E). Per-doc spec-sync verdicts and the
> full decision history are finalized at phase wrap. Newest at top.

## Phase 8 planning decisions (DEC-A..E, human-confirmed 2026-06-20)

These five decisions were proposed at phase planning and confirmed by the human
before execution (phase brief, validation contract §Phase scoping notes). They
are recorded here as the binding Phase 8 decisions; every Phase 8 task spec,
completion report, and QA artifact cites them by the DEC-A..E labels used
throughout the phase.

- **DEC-A — new-attempt completion status.** A standalone `fid issue complete`
  records the new attempt as `status = 'dev_complete'`; prior attempts on the
  same issue move to `'superseded'`. Gate pass/fail is recorded independently via
  `gates_passed`/`gates_failed`; the `failed_gate` status and the paired workflow
  transition are the v0.1e composites' job (`fid dev finish`/`block`), not this
  primitive. **Status: Confirmed; honoured** (VAL-001/002).
- **DEC-B — attempt numbering.** `attempt_number` is derived `MAX(attempt_number)
  + 1` (1 when none) over the issue's `task_completions`, computed under
  `.fid/write.lock`; not user-supplied. **Status: Confirmed; honoured** (VAL-002,
  `nextCompletionAttemptNumber`).
- **DEC-C — gate snapshot is fixture-seeded in Phase 8.** `fid dev gate` (the
  `gate_runs` writer) is v0.1e, so VAL-005 exercises the snapshot mechanism by
  seeding `gate_runs` via journal fixtures. A real Phase 8 `fid issue complete`
  with no gate writer snapshots an empty set (`gate_results = []`, `0`/`0`) —
  correct. Real-write linkage re-verified in v0.1e. **Status: Confirmed;
  honoured** (VAL-005/007; fixture boundary held — Review 1 confirmed).
- **DEC-D — no `activity_log` for completions (DEC-020 precedent).** Completions
  are append-only fact creations and follow the `create`/AC-write precedent: no
  `activity_log` row. Consistent timeline population deferred to v0.1e via
  **CF-011**. No activity assertion in the contract by design. **Status:
  Confirmed; honoured** (VAL-003 test asserts no `activity_log` rows).
- **DEC-E — completions reject on epics before append (CF-007).** `fid issue
  complete` targets tasks/subtasks; an epic target is rejected before any journal
  append (non-zero exit, clear error, zero new event/row). The kind-guard trigger
  is a backstop; the CLI rejects first. Permissions are `RequireKnownRole` only
  (matrix is v0.1d). **Status: Confirmed; honoured** (VAL-004).

No new in-phase DEC-NNN decisions were required: there were zero Developer
deviations, zero validation-contract deviations, and zero escalations.

## Validation Contract Outcome Audit (at wrap)

Final review (`qa/final-review.md`) recorded VAL-001..VAL-007 and VAL-009 PASS
with reproducible QA evidence, and VAL-008 CARRY-FORWARD pending human dogfood.
The human checklist (`qa/human-checklist.md`, HV-08a..e, 5/5 PASS, 2026-06-20)
resolves VAL-008. **All nine assertions now stand PASS — no FAIL / BLOCKED /
CARRY-FORWARD remains.**

| VAL | Outcome | Disposition |
|---|---|---|
| VAL-001 | PASS | `TC-` row with role/tokens/commits/files/notes; CLI JSON read-back (Task 1). |
| VAL-002 | PASS | Attempt 1→2 numbering; prior attempt retained as `superseded` (Task 1). |
| VAL-003 | PASS | `--human-intervention`/`--intervention-note` both states; no `activity_log` (Task 1, DEC-D). |
| VAL-004 | PASS | Epic rejected before append (zero journal growth); task/subtask succeeds (Task 1, DEC-E/CF-007). |
| VAL-005 | PASS | Fixture-seeded `gate_runs`→`gate_results` + derived 1/1; no-gate `[]`/0/0; source rows preserved (Task 2, DEC-C). |
| VAL-006 | PASS | Real `has_completion_report` reject→complete→pass, no fixtures (Task 3). |
| VAL-007 | PASS | Delete-DB + `fid rebuild` parity for completions + gate snapshot (Task 3, ADR-0001). |
| VAL-008 | PASS | Human completion-lifecycle dogfood HV-08a..e 5/5 PASS (was CARRY-FORWARD at final review). |
| VAL-009 | PASS | CF-012 split: `create.go` 387 / `get.go` 120 / `issue.go` 192, no behavior/API change (Task 4). |

No assertion was weakened in-flight. No assertion required strengthening. DEC-C's
fixture-seeded VAL-005 evidence was planned and human-confirmed (not a silent
softening): the contract proves the snapshot mechanism, and real-write linkage is
re-verified in v0.1e.

## Spec Sync Verdict (at wrap)

Recorded for every spec. `verified-no-change` is a positive review assertion,
not a default skip.

| Spec | Verdict | Evidence / Reference |
|---|---|---|
| `00_PROJECT_CHARTER.md` | verified-no-change | Phase 8 delivered the in-scope task-completion write surface; no identity / scope / constraint change. |
| `01_PRD.md` | verified-no-change | PRD-DEV-001/002/003/004/009/010 and PRD-WF-008 MET as written (final review Requirement Compliance). No requirement reworded; composites/`fid dev gate` (PRD-DEV-005–008) remain v0.1e. |
| `02_TRD.md` | verified-no-change | Completions use the existing write-path seam + append-only projection; no architecture or dependency change. DEC-D (no `activity_log` for completions) matches the existing `create`/AC precedent; CF-011 defers consistent timeline population. |
| `03_TRD_SCHEMA.md` | verified-no-change | `task_completions` / `gate_runs` schema and the gate-snapshot comment already shipped (v0.1a). Task 2 verified the snapshot shape against the existing schema/comment; no schema/column/trigger change. |
| `04_TRD_CONFIG.md` | N/A | Phase 8 added no `.fid/` layout or config surface. |
| `05_DESIGN_PLAN.md` | verified-no-change | `issue complete` follows the documented command map and mutation/read flows; no flow or shape change. |
| `06_DESIGN_PLAN_CLI.md` | verified-no-change | Task 1 verified the `fid issue complete` section (flags, output) matches the implementation; no spec text change required. |
| `07_DESIGN_PLAN_FLOWS.md` | verified-no-change | Completions follow the journal→projection→read flow; VAL-007 rebuild parity confirms end-to-end determinism. |
| `08_QUALITY_PLAN.md` | verified-no-change | COV-001 (`internal/issues` 81.9%, `internal/journal` 80.3%, `internal/policy` 85.0%, `internal/workflow` 95.5%) and COV-002 (`cmd/fid` 78.7%) held. §4.4 file-size: `create.go` 387 / `get.go` 120 / `issue.go` 192 / `complete.go` 391 all ≤400; CF-012 closed. No plan reword. |
| `09_AI_WORKFLOW.md` | verified-no-change | Workflow ran as specified (contract → QA plan → dev/QA loop with two mid-phase reviews → final review → human checklist). Zero BLOCKs; both reviews CONTINUE. |
| `10_IMPLEMENTATION_ROADMAP.md` | verified-no-change | Phase 8 delivered the §v0.1c "task completions" cut; Phase 9 (concerns + derived flag + v0.1c closure) unchanged. |

## Triage-Log Audit (at wrap)

**No TRG-NNN entries this phase.** No `triage-log.md` was created: there were
zero Developer BLOCKs, zero TL-Lite triage events, and zero escalations. Both
mid-phase reviews returned `VERDICT: CONTINUE` (`review-1.md` Blocking 0 /
Improvement 3; `review-2.md` Blocking 0 / Improvement 1). In hindsight no event
should have escalated that did not; none over-escalated.

Autonomous TL-Lite resume is permanently enabled (DEC-002); there is no warm-up
tally to maintain.

## Carry-Forward Disposition (at wrap)

- **CF-012 — CLOSED.** Delivered by Task 4 (VAL-009 PASS): behavior-preserving
  split of `internal/issues/issue.go` into `create.go` (387), `get.go` (120), and
  residual `issue.go` (192) — all ≤400 lines, no exported-surface or behavior
  change (Review 2 verified commit `6bf1bf3`, empty `*_test.go` diff, package
  test green). Mirrors the Phase 6 CF-010 precedent.
- **CF-013 — opened (new).** `internal/issues/complete.go` `Complete` is a large
  orchestration function (~194 lines, file is 391 ≤ the 400-line file tripwire).
  Not a Phase 8 blocker — behavior and gates passed — but a candidate for a
  function-level structure cleanup (extract narrow helpers for event construction
  and append/project result assembly). Low priority; schedule alongside the v0.1e
  `fid dev finish`/`block` composites that will also touch this path.
- **CF-011 — re-deferred** to v0.1e (consistent `activity_log` timeline
  population with `fid activity list`, PRD-OUT-017). DEC-D keeps completions
  consistent with the `create`/AC precedent until then.
- **CF-003 — re-deferred** to before v0.1d role-permission / mutable-role work.
- **CF-008 — re-deferred** to v0.2 (resolution type/reason product enhancement).

## QA Recommendation Disposition (at wrap)

Final review verdict: **PASS WITH IMPROVEMENTS** (0 blocking / 0 major / 4
improvement). VAL-008 carried to the human checklist and resolved 5/5 PASS.

- **`complete.go` orchestration size (improvement).** **Defer** — tracked as
  **CF-013**; under the file tripwire, behavior/gates green.
- **Duplicate fixture-only `gate_runs` helpers (improvement).** **Defer to
  v0.1e** — acceptable under DEC-C; centralize/replace when the real `fid dev
  gate` writer lands and the fixture boundary is no longer needed.
- **Task 3 summarized transcript with `<task>` placeholders (improvement).**
  **Acknowledge, no rework** — the deterministic integration test
  (`TestV01CHasCompletionReportGateAndRebuildParity`) is sufficient for
  VAL-006/007. Carry the standing process nudge: future public-output reports
  prefer literal transcripts with concrete IDs.
- **Task 4 compact report omits literal command excerpts (improvement).**
  **Acknowledge, no rework** — Review 2 and final review independently verified
  the commit, file sizes, no-test-change state, package test, format, and build.
  Carry the same self-contained-evidence nudge.
- **Refresh current docs + close CF-012 in wrap artifacts (final-review
  follow-up).** **Done now** — `phase_brief.md`, `progress_current.md`,
  `decisions_active.md`, and `tl_escalations_current.md` refreshed in this wrap;
  CF-012 closed in the carry-forward register.

## Workflow notes

- **What worked well.** A clean four-task vertical cut landed the whole
  completion write surface with **zero BLOCKs and zero escalations**. Pre-phase
  human confirmation of DEC-A..E removed the usual mid-phase decision churn — the
  contentious calls (new-attempt status, fixture-seeded gates, no `activity_log`,
  epic rejection) were settled before Task 1, so the Developer never had to stop
  to escalate. The earlier-phase nudge to scope file-size criteria to the task's
  allowed write set (DEC-022 lesson) meant Task 4's CF-012 split went through
  without the over-broad-criterion BLOCK that bit Phase 7's Task 2b.
- **What was inefficient.** Minor evidence-hygiene drift only: Task 3 used a
  summarized `<task>`-placeholder transcript and Task 4's committed report still
  said "Pending commit creation" at review time. Neither blocked (QA verified
  independently), but both are the same recurring "finalize self-contained report
  evidence before the docs commit" nudge surfacing again.
- **Were tasks the right size?** Yes. Three feature tasks (core write / gate
  snapshot / real-data gating + rebuild) plus the deferred-to-last CF-012 split
  was a clean dependency order. No task needed splitting; none was too thin.
- **Did QA cadence feel right?** Yes. One mid-phase review over Tasks 1–3 and a
  second over the isolated Task 4 refactor matched the work shape; the human
  dogfood fanned cleanly into five per-observable rows (HV-08a..e).
- **Token-efficiency opportunities missed?** Few. Settling DEC-A..E up front was
  the big efficiency win. The only repeatable saving is finalizing report
  transcripts/metadata in the code-behavior commit so phase-end QA does not have
  to reconcile a "pending" committed copy against an updated worktree copy.

### Deep retrospective recommendation

Not required. Phase 8 closed with zero blocking/major findings, all nine VALs
PASS, a 5/5 human dogfood, zero escalations, and zero BLOCKs — the smoothest cut
in the v0.1c slice. The only mildly recurring signal (finalize self-contained
report evidence before the docs commit) is an evidence-hygiene nudge already
carried in standing commitments, not a retro-worthy cross-phase theme. Defer a
deep retro until a recurring or cross-phase problem surfaces.
</content>
</invoke>
