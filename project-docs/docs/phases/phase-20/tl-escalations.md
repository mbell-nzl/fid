# Phase 20 TL Escalations

> Phase-scoped escalation archive. IDs are `ESC-20-NNN`, allocated from
> `triage-log.md`. Open items also appear in
> `../../current/tl_escalations_current.md` until resolved.

## Escalations

### Resolved

- **ESC-20-001 — Review 1 BLOCK on two cross-task evidence gaps (raised
  2026-08-09 via `TRG-20-001`; resolved 2026-08-09 by DEC-130).** **Neither
  finding was a suspected product defect; both were missing evidence for an
  accepted assertion.** (1) Task 1's registry test compared only flattened
  `DetailFields` names, so no test could detect a changed detail value type,
  nullable rule, or nested object/object-list schema for any of the 23 codes,
  leaving **VAL-002** unproven. (2) Tasks 2-3 proved representative malformed
  config input and one 16-character prefix, but not the required
  just-inside/outside matrix — one-character prefix, 128-byte project ID,
  `math.MaxInt64` executor and maintenance values, overflow rejection — leaving
  **VAL-043** unproven. **DEC-130 sustained both findings and changed no
  acceptance criterion.** Recovery is two bounded test-led tasks rather than
  one, because a single session spanning both matrices repeats the
  oversized-evidence shape that caused the gap: **task-3a** (`internal/cli`
  matrix, with equal-length `duplicate-entity` parallel arrays clarified in
  contracts §2.5 and enforcement pre-authorized) and **task-3b** (config/init
  matrix, plus the do-now dangling-symlink collision regression and the do-now
  `internal/config/config.go` split). `internal/bootstrap/init.go`'s 427-line
  breach is **CF-028**, trigger-based, per QA's own note that splitting now
  would obscure the linear activation protocol. Phase Runner resumes at
  task-3a → task-3b → bounded `qa/review-1b.md` → Task 4; the six accepted
  review windows are unchanged.

## Planning-time escalation candidates considered and resolved without an ESC

| Item | Resolution |
|---|---|
| Phase 20 needs machine contracts adjacent to the Phase 21 backend port | Resolved by DEC-117: `fid capabilities`, Envelopes, Claims/Runs, and idempotency stay in Phase 21. No roadmap deviation. |
| `Envelope` is ambiguous between the JSON wrapper and the Work Envelope entity | Resolved by DEC-109 with a `GLOSSARY.md` amendment before implementation, avoiding the `AGENTS.md` ambiguous-terminology hard stop. |
| `03_TRD_SCHEMA.md` §7 blocks implementation until the slice DDL is accepted | Resolved by DEC-108/129: the Human accepted the §8 slice DDL with the implementation pack. |
| PRD-ISSUE-004 dependency-cycle rejection is unimplemented and unowned | Resolved by Human direction and DEC-114: Phase 20 task 11 owns rejection over existing links; Phase 24 still owns normalized dependency storage/commands. |
| Planning overlapped a mistaken V2 rewrite of binding `09_AI_WORKFLOW.md` | Resolved by commit `f04007d` and DEC-118. The V1 workflow is binding; the full Phase 20 pack was re-audited and re-sliced before Human approval. |
| Independent audit found implicit request/Event/init contracts and oversized tasks | Resolved by DEC-122..128: complete aggregate revision, legacy journal boundary, typed links, atomic init, closed/bounded requests, future-safe workflow authoring, and eighteen-task/six-review reslice. |
