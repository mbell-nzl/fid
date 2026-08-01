# TL Escalations - Current

> **Purpose.** Open items that exceed local-decision authority and need full
> Tech Lead or human input.
>
> **Scope.** Open items only. Resolved phase items move to the phase archive.
>
> **Size target:** under 80 lines.

## Open Escalations

None.

## Resolved This Phase (archive at wrap)

Full reasoning lives in `../phases/phase-19/decisions.md` and
`../phases/phase-19/triage-log.md`; these are pointers.

- **ESC-19-004 — Task 4's first VAL-059 baseline capture failed (2026-08-01 →
  DEC-099/DEC-100).** `capture.sh` resolved a relative output root after
  `cd`-ing into its temp workspace. The task stopped with no commit, six gates
  unrun, and a registered scratch worktree. **The BLOCK rightly stopped the run
  and wrongly reasoned about why:** the failure was a path bug in a *Task 4
  deliverable being authored*, it wrote nothing so could not contaminate
  anything, and DEC-098's guard is **structural** (detached worktree at the
  pinned commit), not temporal. **DEC-099:** preserve the partial work, settle
  baseline provenance by **reproduction against `SHA256SUMS`**, remove the
  worktree via `git worktree remove`, finish normally — no restart. **DEC-100:**
  standing rule — a task-authored harness failing while being authored is a
  **deviation, not a failed gate**. Task 4 cleared to resume.
- **ESC-19-003 — VAL-059 baseline instruction conflicted with the required
  immutable-baseline procedure (2026-07-30 → DEC-098).** Task 4 said `git
  rev-parse HEAD`; QA-plan §6 requires the Task 3 *implementation* commit, and
  at Review 1 the two differed. **QA was right; the task spec had drifted.** The
  baseline is now pinned to `50533e2b813bf073cd75f75eb70dd09cd15e9040`, moving
  selectors are prohibited, Task 4 is synchronized to QA-plan §6, and
  normalization is the committed no-op. `qa-plan.md` needed no edit.
- **ESC-19-002 — Task 2 lint gate blocked by a test-only import cycle
  (2026-07-29 → DEC-096/DEC-097).** The binding `[]journal.Effect` return type
  makes `workflow → journal → db`, so the in-package
  `internal/db/workflow_tables_test.go` closed a cycle and failed `typecheck`.
  **DEC-096** confirms the `Effect`-only edge and authorizes moving the two
  tests to the external `db_test` package; no production or API change.
  **DEC-097** root-causes the missing triage record as a **Codex
  provider-prompt defect** and fixes `project-docs/agents/tl-lite.md`;
  `sync-agents.sh` has been run and verified. Backfilled as TRG-19-002.
- **ESC-19-001 — Task 1 VAL-031 evidence token (2026-07-26 → DEC-094/DEC-095).**
  The report recorded a no-match `rg` as `PASS - exit 0 (no output)`; `rg` exits
  1 on no match. The observable itself held. Report corrected; TRG-19-001
  backfilled and the TL-Lite write failure dispositioned.

## Recent Archive Notes

- **Phase 18** (wrapped 2026-07-25) raised **ESC-001** — a final-review BLOCK on
  the `workflow.Load` connection deadlock. Resolved by **DEC-072**, archived to
  `../phases/phase-18/tl-escalations.md`. Its ID collides with Phase 17's
  ESC-001; left as recorded (the CF-023 defect that motivated DEC-084).
- **Phase 18's** other five triages (TRG-001..005) returned `BLOCK` and were
  adjudicated as DEC-067..071 without an `ESC-NNN`. Correct in substance, but in
  five of six cases TL-Lite reported an escalation it never minted and wrote no
  triage log — now fixed by DEC-085 and tracked as **CF-023**.
- **Phase 17** (wrapped 2026-07-18) raised five escalations, all resolved and
  archived to `../phases/phase-17/tl-escalations.md`, plus the WRAP-F1/F2
  register reconciliation. Its two open retrospective items became **CF-023**.
- **Phase 16** raised no `ESC-NNN`; DEC-043..048 and TRG-046 are archived in
  `../phases/phase-16/`. Phase 16 VAL-013 is not an escalation — the completed
  human checklist records PASS for HUMAN-16a..d.

## Standing Escalation Triggers

TL-Lite escalates instead of deciding locally when an issue:

1. Changes a public interface, CLI contract, output shape, journal event schema,
   data model, persistence contract, dependency, tool, quality gate, or workflow.
2. Contradicts `decisions_active.md`, `decisions_principles.md`, `AGENTS.md`, or
   a recorded phase decision.
3. Affects more than one task/subsystem, exposes product-intent ambiguity, or
   changes testing strategy / coverage policy.
4. Resumes context-exhausted or partial work spanning schema/data model,
   dependencies, public interfaces, canonical docs, or multiple Developer
   sessions.

QA findings escalate here if they have architecture implications, ambiguous
acceptance criteria, unclear security/privacy severity, contradict a standing
decision, span tasks/subsystems beyond the review window, plausibly indicate a
material spec deviation, or expose a quality-gate anomaly not anticipated by the
QA plan.
