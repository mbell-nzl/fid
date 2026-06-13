# Phase 4 Decisions

> Purpose: phase-local decision history, spec sync verdicts, triage audit,
> workflow notes, and wrap-up decisions. DEC numbering is global and never
> resets (Phase 3 ended at DEC-013).

## Decisions

### DEC-014: Rejected `fid issue create` poisons the journal tail — fix scheduled as CF-007

- **Date:** 2026-06-13
- **Context:** The Phase 4 human dogfood checklist (VAL-010, HV-01e) executed a
  deliberate invalid create — a task with `--sequence 1` under prioritised epic
  `FID-e36ee9`. The CLI returned the expected scheduling constraint error, but
  the rejected command was still appended to the canonical journal as
  `EVT-20260610-000007` in `.fid/journal/2026-06-10.jsonl`. Afterwards, unrelated
  reads such as `fid issue get FID-a816f0` failed with
  `project event "EVT-20260610-000007": insert into issues: constraint failed:
  issues: task under prioritised epic must not have sequence (1811)`.
- **Diagnosis (from the checklist):** `issues.Create` appends the event before
  projecting it. When projection rejects the event (the Task 1 sequence
  backstop trigger fires), the append has already committed, poisoning the
  journal tail. Later reads call `EnsureProjectionFresh`, replay the invalid
  tail event, and fail before serving the requested (valid) issue. The valid
  issue `FID-a816f0` is collateral damage.
- **Spec conformance:** This is an **implementation conformance gap, not a spec
  change.** TRD §4.3 already orders the write path as validate → event
  construction → journal append, with "journal append is the commit point" and
  "one logical write = one event." A create whose projection will fail must be
  rejected **before** the append. The current order violates that contract.
- **Decision:** Do **not** hotfix inside this wrap (Tech Lead does not write
  product code, and Phase 4 is already PASS on its contract). Record the defect
  as carry-forward **CF-007** targeting **Phase 5** (mutation/workflow slice,
  where write-path discipline and the validate-before-append seam are already in
  focus). Required remediation: validate scheduling/sequence (and any other
  cross-table trigger outcome) before journal append, or add an equivalent
  pre-append projection-dry-run path (project into a rolled-back transaction),
  so the append remains the commit point only for events that will project
  cleanly. Add a regression test asserting a rejected create leaves the journal
  line count unchanged, plus a poisoned-tail read-recovery test.
- **Human attention:** Flagged in the Phase 4 human checklist wrap note and in
  `progress_current.md`. Severity is journal-integrity / data-availability (a
  single rejected create blocks all subsequent reads until the tail is removed),
  so the human should confirm CF-007 priority before Phase 5 planning closes.
- **Refs:** `qa/human-checklist.md` HV-01e + TL Follow-Up Notes; VAL-005
  (`qa/validation-contract.md`); CF-007 (`decisions_active.md`); TRD §4.3,
  PRD-OPS-001/002; `internal/issues.Create`.

## QA Recommendation Disposition

QA final review (`qa/final-review.md`, VERDICT: PASS) recorded **0 blocking /
0 major / 0 improvement** findings. The single mid-phase blocker (`review-3.md`,
ESC-002) was resolved by Task 4 due-date read-back evidence. No QA recommendation
remains open; nothing to do-now, defer, or cancel from QA.

The only open item out of Phase 4 is human-discovered, not QA-raised: the
journal-poisoning defect (DEC-014 → CF-007).

## Validation Contract Audit

- **VAL-001..VAL-009: PASS** with cited deterministic evidence
  (`qa/final-review.md` Validation Contract Outcomes). VAL-003 due-date public
  read-back was closed in Task 4 per ESC-002 without schema/journal change.
- **VAL-010: PASS.** The human dogfood checklist recorded 5/5 PASS
  (`qa/human-checklist.md`). At final review VAL-010 was `CARRY-FORWARD`
  (human-only judgment not executable by QA Lead); the human has since executed
  it, so the assertion is now satisfied. No VAL outcome remains FAIL, BLOCKED, or
  CARRY-FORWARD.
- **Both-state coverage** held for the binary/conditional assertions: VAL-001
  (accept/reject scheduling states + CLI errors), VAL-002 (valid/invalid
  hierarchy), VAL-004 (epic metadata present / non-epic rejected), VAL-007
  (supported sections + unknown-section rejection), VAL-008 (include/exclude per
  filter + empty result), VAL-009 (format parity + `--format xml` rejection).
- **Coverage gap noted (not a contract failure):** VAL-005 asserted that each
  *successful* create is exactly one rebuildable event. It did not assert
  journal hygiene for *rejected* creates — the exact adjacent state the human
  surfaced (DEC-014). VAL-005 PASS stands as written for the states it covers; a
  Phase 5 assertion should explicitly require rejected writes to leave the
  journal unchanged. Folded into CF-007.

## Spec Sync Verdict

Record a verdict for every spec at phase wrap-up. `verified-no-change` is a
positive assertion after review, not a default skip.

| Spec | Verdict | Evidence / Reference |
|---|---|---|
| `00_PROJECT_CHARTER.md` | verified-no-change | Phase 4 delivered the in-scope v0.1b authoring/discovery slice; no identity/scope/constraint change. |
| `01_PRD.md` | verified-no-change | PRD-ISSUE-001..007/013..016/018..020, PRD-WF-005/006, PRD-OUT-004 MET as written (final review Requirement Compliance). The DEC-014 defect is an implementation conformance gap against PRD-OPS-001/002, not a requirement reword. |
| `02_TRD.md` | verified-no-change | TRD §4.3 already mandates validate → construct → append with append as the commit point; DEC-014 is an implementation deviation from that order (CF-007), not a spec change. TRD-ARCH-005 (`cmd/fid` thin) MET via CF-006 closure (VAL-006). No dependency change this phase. |
| `03_TRD_SCHEMA.md` | verified-no-change | Task 1 implemented the §14.1 sequence-vs-scheduling backstop trigger (CF-004) exactly as already specified; no schema/trigger/view spec reword. |
| `04_TRD_CONFIG.md` | N/A | Phase 4 added no `.fid/` layout or config surface. |
| `05_DESIGN_PLAN.md` | verified-no-change | §6.1 enriched read-model centralisation (VAL-006) and §6.2 list filters/order (VAL-008/009) realized as described. |
| `06_DESIGN_PLAN_CLI.md` | updated | Task 5 updated the `fid issue list` filters/ordering/format documentation in-phase (see `qa/final-review.md` Documentation Review). `fid issue get --section` and create surfaces match the documented contract. |
| `07_DESIGN_PLAN_FLOWS.md` | verified-no-change | create/get/list follow the documented journal→projection→read flow. DEC-014 is a conformance defect in the create write-path implementation (CF-007), not a flow redesign. |
| `08_QUALITY_PLAN.md` | verified-no-change | COV-001 (`internal/issues` 85.5%, `internal/journal` 80.7%, `internal/policy` 85.0%) and COV-002 (`cmd/fid` 79.5%) held; DEC-007 carry-in applied. The rejected-create journal-hygiene coverage gap is recorded for CF-007, not a quality-plan reword. |
| `09_AI_WORKFLOW.md` | verified-no-change | Workflow ran as specified (contract → QA plan → dev/QA loop → one mid-phase review → final review → human checklist). The human-dogfood step (VAL-010) caught a defect the automated gates passed over — the workflow performed as designed. |
| `10_IMPLEMENTATION_ROADMAP.md` | verified-no-change | Phase 4 delivered the v0.1b "Issue authoring + discovery" slice per the 2026-06-07 slicing note; Phase 5/6 cuts unchanged. |

## Triage-Log Audit

- **TRG-001** (Task 1 post-commit completion-report repo-state note) — verdict
  **RESUME**. Confirmed correct in hindsight: the completion report met all
  required headings/evidence; the note was a repo-history hygiene observation,
  not an acceptance blocker. Matched no escalation trigger.
- **TRG-002** (due-date public read-back gap) — verdict **BLOCK**. Confirmed
  correct in hindsight: it correctly stopped on a real public read-path gap and
  routed to the existing ESC-002 rather than passing a persistence-only proof as
  read-back. ESC-002 was dispositioned DO-NOW and closed by Task 4. No duplicate
  escalation; appropriate handling.
- No TRG should have escalated that did not. Autonomous TL-Lite resume remains
  permanently ENABLED (DEC-002); no warm-up tally to maintain.

## ESC-002 Disposition

**Resolved and archived.** ESC-002 (due-date create field lacked a public
read-back surface) was dispositioned DO-NOW on 2026-06-09 (keep VAL-003 intact,
fold the fix into Task 4). Task 4 evidence proves `due_date` reads back through
`internal/issues.Get` (`2026-06-30` / `nil`) and `fid issue get --format json`
(`"due_date":"2026-06-30"` / `null`). Moved to
`project-docs/docs/phases/phase-4/tl-escalations.md`.

## Workflow Notes

Answer in 1-3 sentences each. Completed at phase wrap-up.

- **What worked well?** The human dogfood step (VAL-010) earned its keep: every
  automated gate and all nine deterministic VALs passed, yet a realistic invalid
  create surfaced a genuine journal-integrity defect that test coverage had not
  exercised. Resolving ESC-002 as a folded-in Task 4 closeout (rather than a
  contract bounce) kept Tasks 2-3 valid while still enforcing public read-back.
- **What was inefficient?** Task 5 shipped an inaccurate COV statement that had
  to be corrected by a dedicated Task 6 closeout; tighter coverage reporting in
  the Task 5 report would have avoided the extra task. Otherwise the phase ran
  clean — zero deviations, zero dependency changes across six tasks.
- **Were tasks the right size?** Mostly. Tasks 1-5 were well-sized one-session
  slices; Task 6 was a small remedial coverage-close task that ideally would not
  have existed. The original five-task cut matched the authoring/discovery scope.
- **Did the QA review cadence feel right?** Yes. One mid-phase review after Task 3
  caught the due-date read-back gap (ESC-002) at the right moment, and the
  phase-end final review re-ran the full gate. The defect that slipped through is
  a test-coverage gap, not a cadence problem.
- **Were there token-efficiency opportunities missed?** Minor. The Task 5 → Task 6
  COV correction loop spent tokens on a reporting fix; accurate first-pass
  coverage evidence would have folded that into Task 5.

### Deep Retrospective Recommendation

**Lightweight retro warranted on one theme: validation-contract coverage of
rejected/abort write paths.** Phase 4 proved the happy-path and rejected-CLI-error
states but did not assert that a *rejected* mutation leaves the canonical journal
unchanged — and that exact gap became a real, read-blocking defect (DEC-014 /
CF-007). The durable lesson for Phase 5 and every later write-command phase: when
a write path can be rejected by a projection-time trigger, the validation contract
must assert journal hygiene (line-count-unchanged + read-recovery) for the
rejected state, not only an actionable CLI error. No broader process change is
needed; capture this as a planning checklist item when slicing write commands.
