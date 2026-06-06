# Phase 1 - TL-Lite Triage Log

> **Purpose.** Append-only log of mid-phase triage events handled by TL-Lite.
> Phase Runner reads TL-Lite's `VERDICT:` line to decide whether to resume or
> stop. Tech Lead audits this log at phase wrap-up.
>
> **Numbering.** TRG-NNN numbering is per-phase and resets each phase. Newest
> entry at top.

| TRG-NNN | Date | Trigger | Source | Verdict | ESC-NNN (if BLOCK) | Summary |
|---|---|---|---|---|---|---|
| TRG-003 | 2026-06-06 | ESC-001 / QA Review 1 says `internal/db.Open` does not prove pragmas on every pooled SQLite connection for VAL-014. | `project-docs/docs/current/tl_escalations_current.md` ESC-001 | BLOCK | ESC-001 | Existing escalation remains unresolved: VAL-014 requires pragmas on every connection, and loaded context does not disprove the pooled-connection concern or record an accepted decision waiving it. |
| TRG-002 | 2026-06-06 | Task 6 report notes that role validation scans journal events directly and tolerates a corrupt trailing line after seeded roles are discovered, to be revisited later. | `project-docs/docs/phases/phase-1/tasks/task-6-completion.md` | RESUME | - | Not a current phase stopper: Phase 1 only requires seeded-role existence checks, and the fallback is limited to bootstrap role discovery so `fid rebuild` can still report the intended journal corruption through the replay path. Revisit before mutable role events or broader policy semantics ship. |
| TRG-001 | 2026-06-06 | Task 4 report notes a mismatch between `02_TRD.md` §4.3 and the Phase 1 validation contract on last-event-hash mismatch handling during projection freshness. | `project-docs/docs/phases/phase-1/tasks/task-4-completion.md` | RESUME | - | Not a current phase stopper: Task 4, VAL-012, and the QA plan consistently require rebuild/catch-up before serving on last-event-hash mismatch. Reconcile TRD wording before later phases build on this behavior. |

---

## TRG-003 - SQLite pragmas on pooled connections

- **Trigger:** ESC-001 / QA Review 1 says `internal/db.Open` does not prove pragmas on every pooled SQLite connection for VAL-014.
- **Source:** `project-docs/docs/current/tl_escalations_current.md` ESC-001
- **Verdict:** BLOCK
- **Reasoning:** QA blocking triage may resume only when the finding is misclassified, already resolved by verified evidence, or covered by an existing accepted decision. The loaded context records ESC-001 as open and states that `internal/db.Open` only proves the physical connection used for the initial PRAGMA Exec calls, while later `database/sql` pooled connections may be opened without those per-connection settings. VAL-014 requires pragmas ON on every connection path, and this issue affects the persistence contract consumed by later Phase 1 tasks.
- **Local decisions made:** None.
- **Escalation:** Existing ESC-001 remains the controlling escalation; no duplicate escalation was created.
- **References:** `project-docs/docs/current/phase_brief.md` Validation Contract Index VAL-014; `project-docs/docs/current/tl_escalations_current.md` ESC-001; `project-docs/docs/phases/phase-1/qa/validation-contract.md` VAL-014

## TRG-002 - Role validation journal scan corrupt-tail tolerance

- **Trigger:** Task 6 report notes that role validation scans journal events directly and tolerates a corrupt trailing line after seeded roles are discovered, to be revisited later.
- **Source:** `project-docs/docs/phases/phase-1/tasks/task-6-completion.md`
- **Verdict:** RESUME
- **Reasoning:** Phase 1 deliberately limits role policy to presence of `--by` plus seeded-role existence; the full permission matrix and mutable role lifecycle are deferred. The task report states the corrupt-tail tolerance is in the role-validation precheck after seeded roles have already been discovered, so it does not authorize replay to skip corruption and does not weaken VAL-010's requirement that malformed trailing journal lines are reported as corruption. The current Task 6 evidence says the role-validation matrix, create/get path, rebuild parity, full task gate, and build all passed. This is acceptable Phase 1 bootstrap debt, but it should be revisited before later phases add mutable role events or richer policy semantics.
- **Local decisions made:** Current Phase 1 execution may continue; this triage entry does not authorize a public-interface, data-model, persistence-contract, or testing-strategy change.
- **References:** `project-docs/docs/current/phase_brief.md` Scope out; `project-docs/docs/phases/phase-1/tasks/task-6.md` Objective / Out Of Scope / Acceptance Criteria; `project-docs/docs/phases/phase-1/qa/validation-contract.md` VAL-010 and VAL-016; `project-docs/docs/phases/phase-1/qa/qa-plan.md` QA-06 and QA-16; `project-docs/docs/phases/phase-1/tasks/task-6-completion.md` Validation Contract Coverage / Tests run / Concerns, technical debt, or notes for the Tech Lead

## TRG-001 - Projection freshness hash-mismatch wording

- **Trigger:** Task 4 report notes a mismatch between `02_TRD.md` §4.3 and the Phase 1 validation contract on last-event-hash mismatch handling during projection freshness.
- **Source:** `project-docs/docs/phases/phase-1/tasks/task-4-completion.md`
- **Verdict:** RESUME
- **Reasoning:** The durable TRD wording says a recorded event ID/hash mismatch at the recorded position stops before answering with projection inconsistency. However, the current Phase 1 task spec, VAL-012, and QA plan all require a missing/stale/mismatched `projection_meta` freshness case, including last-event-hash mismatch, to trigger rebuild/catch-up before serving. Task 4 implemented and tested the behavior selected by those current Phase 1 acceptance artifacts, so this does not block the current phase. The TRD wording should be reconciled before later phases rely on this persistence behavior as a durable contract.
- **Local decisions made:** Current Phase 1 execution may continue under Task 4 / VAL-012 / QA plan behavior; no public interface, schema, or implementation change is authorized by this triage entry.
- **References:** `project-docs/specs/02_TRD.md:501-511`, `project-docs/docs/phases/phase-1/tasks/task-4.md:73-75`, `project-docs/docs/phases/phase-1/qa/validation-contract.md:42`, `project-docs/docs/phases/phase-1/qa/qa-plan.md:157`, `project-docs/docs/phases/phase-1/tasks/task-4-completion.md:33,58`

<!--
Phase wrap-up checklist:

- [ ] Walked every TRG-NNN above.
- [ ] Confirmed each RESUME verdict in hindsight.
- [ ] Flagged any TRG that should have escalated.
-->
