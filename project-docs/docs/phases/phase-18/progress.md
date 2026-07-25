# Phase 18 Progress

> Phase: 18 - Workflow tables migration + engine data-source swap + gate seams
> (roadmap v0.2a, slice 1 of 3)
> Status: COMPLETE / WRAPPED (2026-07-25). All five tasks closed, QA Lead final
> review PASS (rerun), Human checklist PASS, VAL-001..024 all PASS. No Developer
> task is active.

## Plan

Phase 18 converted the v0.1 hard-coded workflow state machines into
data-configurable tables (`workflows`, `workflow_states`, `workflow_transitions`,
`workflow_transition_requires`), repointed the transition engine and every
consumer to resolve from those tables while preserving v0.1 behavior exactly,
migrated `issues.workflow_id` / `issues.status` to real FKs via a table rebuild,
and reserved the DEC-066 gate-obligation seam with no receipt behavior.
Operating mode: **Heavy**.

## Task Index

| Task | Title | Status |
|---|---|---|
| 1 | Workflow tables migration + seed defaults + gate seam | Complete (`127f268`) |
| 2 | DB-backed workflow registry loader + v0.1 parity fixtures | Complete (`799acf0`) |
| 3 | Transition-path data-source swap (behavior-preserving) | Complete (`cf68153`; VAL-020 remediation `77400ce`, DEC-068) |
| 4 | `issues.workflow_id/status` FK migration (table rebuild) | Complete (`0b4b8aa`; VAL-006 remediation `a27fb5f`, DEC-069) |
| 5 | Remaining consumer swap + maintenance parity + closure | Complete (`37b3fb5`; DEC-071 perf fix, DEC-072 loader fix `795a7df`, report `37fb202`) |

## Validation

`qa/validation-contract.md` VAL-001..024 (24 assertions, all machine-checkable,
QA Lead-owned) — **24 PASS / 0 FAIL / 0 BLOCKED / 0 CARRY-FORWARD**, each cited
against fresh evidence in the 2026-07-23 final-review rerun (FRR-E01..E12). The
contract was refined to Ready on 2026-07-19 (DEC-067, VAL-015..024 added). No
assertion, threshold, fixture, or gate command was weakened at any point.

Performance controls (not `VAL-NNN`, per QA-BLK-18-03): 10k-issue registry read
`47.59–49.67 ms` (<500 ms), 10k-event rebuild `2.83–2.92 s` (<5 s), 100k-event
rebuild `26.99–27.66 s` (<30 s, thin headroom — see CF-024). Coverage total
`76.8%` with every enforced package threshold met.

## Progress History

- Planning (2026-07-18/19): Heavy Mode declared; validation contract drafted;
  five-task pack written; QA plan committed `BLOCKED BEFORE TASK 1` with
  QA-BLK-18-01..03.
- Entry unblock (2026-07-19): TRG-001 → **DEC-067** refined the contract to Ready
  (VAL-015..024 added), made the Task 4 dependent-object inventory exact, and
  added the TRD-PERF-001/003 controls. Tasks 1–3 completed.
- Review 1 (2026-07-19): BLOCK on the VAL-020 rebuild-evidence gap. TRG-002 →
  **DEC-068** ordered a do-now Task 3 remediation (evidence supplied, assertion
  unchanged); re-review CONTINUE. Task 4 completed.
- Review 2 (2026-07-19): BLOCK on the VAL-006 valid-pair UPDATE gap. TRG-003 →
  **DEC-069** ordered a do-now Task 4 remediation; re-review CONTINUE.
- Task 5 gate interruption (2026-07-19): TRG-004 → **DEC-070** retained the
  uncommitted work and required the interrupted benchmark/build gate to be
  finished, not relaxed.
- Performance breach (2026-07-20): the completed gate exposed a genuine
  TRD-PERF-003 failure (100k rebuild `33.38 s`). TRG-005 → **DEC-071** kept the
  `<30 s` target and made the rebuild faster instead: a rebuild-only connection
  (`synchronous=OFF`, `journal_mode=MEMORY`, FKs still on) via `OpenForRebuild`.
  Task 5 committed at `37b3fb5`.
- Final review BLOCK (2026-07-20): 24 VAL PASS but two Blocking findings —
  `workflow.Load` deadlocked on a one-connection pool (four required commands
  never terminated) and the QA-30 spec sync was absent. TRG-006 → **ESC-001** →
  **DEC-072**: fix the loader (not the fixture), add a bounded-timeout regression
  test, correct the over-claimed command records, and apply the spec sync under
  Tech Lead authority.
- Remediation and re-review (2026-07-23): loader fix `795a7df`, report `37fb202`;
  QA Lead final-review rerun **PASS** (`77a4560`) with 24 PASS / 0 FAIL and no
  Blocking or Major finding.
- Human checklist (2026-07-25): HUMAN-01a..d all **PASS** (`0edcfc8`); Phase 18
  accepted for continued dogfood.
- Wrap (2026-07-25): **DEC-073** signs off the phase, audits all 24 VAL outcomes
  and all six triages, opens CF-024, re-targets CF-017/CF-021/CF-008, and syncs
  `GLOSSARY.md`. **DEC-074** promotes CF-023 to a required workflow-tooling pass
  before Phase 19 and corrects the roadmap sequencing note. ESC-001 archived.

## Carry-Forward (final)

- **CF-024 NEW (DEC-073)** — TRD-PERF-003 headroom is thin: the 100k-event
  rebuild sits at `26.99–27.66 s` against the `<30 s` control after DEC-071.
  Reopen with the first v0.2b projection/rebuild work (receipts/evidence raise
  event volume on this path); next lever is batching commits per N events. Do
  not weaken the target. Target: v0.2b.
- **CF-023 (DEC-058, promoted by DEC-074)** — workflow-tooling hardening is now
  **required before Phase 19 execution**, with scope expanded from Phase 18
  evidence: mint `ESC-NNN` + triage log on every `BLOCK`; `ESC-NNN`/`TRG-NNN`
  global uniqueness; report-lint/commit-metadata friction; run-ledger
  chronological append and no duplicates; a wrap-time reopen audit for
  trigger-based carry-forwards; evidence-honesty hardening so a non-terminating
  command cannot be recorded `PASS`.
- **CF-021 (re-targeted)** — its "next `internal/maintenance` reopen" trigger
  fired in Phase 18 (Task 5 modified `check.go`/`warnings.go`) and was missed.
  Now an explicit **Phase 19 planning item**, not trigger-based.
- **CF-017 (re-targeted)** — its "next `internal/dev` work" trigger fired in
  Phase 18 (Task 5 modified `start/finish/block/gate/dev.go`) and was correctly
  not taken during a behavior-preserving migration. Target stays **v0.2e**, but
  as a scheduled planning item rather than trigger-based.
- **CF-008 (re-deferred)** — explicit decision point at **Phase 19 planning**
  (Phase 19 delivers PRD-WF-011/012, the pairing CF-008 named).
- **CF-018, CF-020, CF-022** unchanged; their trigger packages
  (`internal/sqlquery`, `internal/format`) were not modified this phase.

## Escalations

One escalation, resolved and archived in `tl-escalations.md`: **ESC-001**
(final-review BLOCK — `workflow.Load` connection deadlock + absent QA-30 spec
sync; resolved by DEC-072, verified by the 2026-07-23 rerun). Its ID collides
with Phase 17's ESC-001 and is left as recorded (CF-023 / DEC-074).

Triage: six entries, all `BLOCK`, all correct in hindsight (see the Triage-Log
Audit in `decisions.md`). Only TRG-006 minted a real escalation identifier and a
triage-log entry; the other five did not — the CF-023 defect class.
