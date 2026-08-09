# Phase 20 TL Escalations

> Phase-scoped escalation archive. IDs are `ESC-20-NNN`, allocated from
> `triage-log.md`. Open items also appear in
> `../../current/tl_escalations_current.md` until resolved.

## Escalations

_None yet._

## Planning-time escalation candidates considered and resolved without an ESC

| Item | Resolution |
|---|---|
| Phase 20 needs machine contracts adjacent to the Phase 21 backend port | Resolved by DEC-117: `fid capabilities`, Envelopes, Claims/Runs, and idempotency stay in Phase 21. No roadmap deviation. |
| `Envelope` is ambiguous between the JSON wrapper and the Work Envelope entity | Resolved by DEC-109 with a `GLOSSARY.md` amendment before implementation, avoiding the `AGENTS.md` ambiguous-terminology hard stop. |
| `03_TRD_SCHEMA.md` §7 blocks implementation until the slice DDL is accepted | Resolved by DEC-108/129: the Human accepted the §8 slice DDL with the implementation pack. |
| PRD-ISSUE-004 dependency-cycle rejection is unimplemented and unowned | Resolved by Human direction and DEC-114: Phase 20 task 11 owns rejection over existing links; Phase 24 still owns normalized dependency storage/commands. |
| Planning overlapped a mistaken V2 rewrite of binding `09_AI_WORKFLOW.md` | Resolved by commit `f04007d` and DEC-118. The V1 workflow is binding; the full Phase 20 pack was re-audited and re-sliced before Human approval. |
| Independent audit found implicit request/Event/init contracts and oversized tasks | Resolved by DEC-122..128: complete aggregate revision, legacy journal boundary, typed links, atomic init, closed/bounded requests, future-safe workflow authoring, and eighteen-task/six-review reslice. |
