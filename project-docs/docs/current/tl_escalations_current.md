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

- **ESC-19-003 — VAL-059 baseline instruction conflicted with the required
  immutable-baseline procedure (raised 2026-07-30, resolved 2026-07-30 by
  DEC-098).** Task 4 directed the Developer to record `git rev-parse HEAD` as
  the Task 3 baseline. At Review 1 ambient `HEAD` was the Task 3 *report*
  commit `697c673`, while QA-plan §6 requires the Task 3 *implementation*
  commit `50533e2`, explicitly not ambient `HEAD` — so the immutable v0.1
  parity evidence could have been unreproducible or contaminated with later
  artifacts. **QA was right and the task spec was the artifact that drifted.**
  Resolution (**DEC-098**): the baseline is now **pinned** to the literal
  verified hash `50533e2b813bf073cd75f75eb70dd09cd15e9040` in both `task-4.md`
  and `validation-contract.md`, with a required `git rev-parse --verify` plus
  `task-3-done.md` cross-check and an explicit prohibition on `HEAD`/branch
  selectors; Task 4's whole procedure is synchronized to QA-plan §6 (worktree
  paths, `capture.sh` signature, fixture, six output names, `SHA256SUMS`);
  normalization is fixed to the committed **no-op**, with `BLOCKED`+escalate as
  the only response to unstable timestamps. `qa-plan.md` needed no edit.
  **Task 4 is cleared to start.**
- **ESC-19-002 — Task 2 lint gate blocked by a test-only import cycle (raised
  2026-07-29, resolved 2026-07-29 by DEC-096/DEC-097).** Task 2's binding
  `[]journal.Effect` return type makes `internal/workflow` import
  `internal/journal`, which imports `internal/db`; the in-package test file
  `internal/db/workflow_tables_test.go` imports `internal/workflow`, so the
  `db` test binary closes a cycle and `golangci-lint run ./...` fails
  `typecheck`. Resolution: **DEC-096** confirms `workflow → journal` as the
  intended, minimal (`Effect`-only) dependency and authorizes the narrow repair
  — move the two workflow-importing tests into the existing external `db_test`
  package. No production code changes; no change to the Task 2 API. **DEC-097**
  root-causes the missing triage record as a **Codex provider-prompt defect**
  (shared body asserted a "cannot run commands" limit that only Claude's
  frontmatter enforces; on Codex that reads as "cannot write files") and fixes
  `project-docs/agents/tl-lite.md` — **`sync-agents.sh` must be run.**
  Backfilled by the Tech Lead as TRG-19-002.
- **ESC-19-001 — Task 1 VAL-031 evidence token (raised 2026-07-26, resolved
  2026-07-29 by DEC-094/DEC-095).** `task-1-done.md` records the no-match `rg`
  inventory as `PASS - exit 0 (no output)`; `rg` exits 1 on no match, so the
  token was not copied from the command (DEC-086 violation). The VAL-031
  observable itself holds. Resolution: report correction + re-execution
  (DEC-094); the missing triage record is backfilled as TRG-19-001 and the
  TL-Lite write failure is dispositioned by DEC-095. Backfilled by the Tech
  Lead — TL-Lite raised the BLOCK without minting this ID.

## Recent Archive Notes

- Phase 18 wrapped 2026-07-25. It raised one escalation, **ESC-001**
  (final-review BLOCK: the `workflow.Load` connection deadlock that left four
  required gate commands non-terminating, plus the then-absent QA-30 spec sync).
  Resolved by **DEC-072**, verified by the 2026-07-23 final-review rerun, and
  archived to `../phases/phase-18/tl-escalations.md`. Its ID collides with Phase
  17's ESC-001 and is left as recorded — the known CF-023 defect.
- Phase 18's other five triages (TRG-001..005) returned `BLOCK` and were
  adjudicated as DEC-067..DEC-071 without an `ESC-NNN`; that routing was correct
  in substance, but in five of six cases TL-Lite reported an escalation it never
  minted and wrote no triage log. **CF-023 is now a required pre-Phase-19
  workflow-tooling pass (DEC-074)** covering exactly that: mint `ESC-NNN` **and**
  the triage-log entry on every `BLOCK`, and enforce `ESC-NNN`/`TRG-NNN` global
  uniqueness.
- Phase 17 wrapped 2026-07-18. It raised five escalations, all resolved and
  archived to `../phases/phase-17/tl-escalations.md`: ESC-004 (Task 1
  report-lint), ESC-005 (Task 2 report-lint, backfilled), ESC-009 (QA model
  qualification), ESC-010 (rebuild parity; resolved by DEC-056 Branch A), and
  ESC-011 (Task 7 report-lint, backfilled with a new ID). That file also records
  the WRAP-F1/F2 register reconciliation. The two open retrospective items -
  recurring report-lint commit-metadata friction and **ESC-NNN IDs not being
  globally unique** (ESC-001/004/005 collide across phases) - are now tracked as
  **CF-023** (DEC-058) and flagged for the deep retrospective before v0.2a.
- Phase 16 raised no `ESC-NNN` items. Local unblock and workflow decisions
  DEC-043..048 plus TRG-046 are archived in
  `../phases/phase-16/tl-escalations.md` and `../phases/phase-16/decisions.md`.
- Phase 16 VAL-013 is not an escalation; the completed human checklist records
  PASS for HUMAN-16a..d.

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
