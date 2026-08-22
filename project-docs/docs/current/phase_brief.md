# Phase Brief - Phase 21

> Status: **Planning / Awaiting Human Execution Approval**
> Branch: `phase-21`
> Product implementation: **NOT AUTHORIZED**

## Goal And Mode

Deliver one Heavy one-Issue operational walking skeleton: activated trust,
immutable Gate Registration, one Work Envelope/revision, one gate-evidence
Obligation, generic Actor/Claim/Run authority, fid-executed Receipt Evidence,
guarded terminal-success Issue transition, later Envelope completion, bounded
resume, export create/validate, rebuild parity, and owned V1 removal.

**Operating mode: Heavy.** Authority, public JSON/Event contracts, destructive
schema convergence, local process execution, credentials/redaction, Artifact
safety, completion, and rebuild cross protected boundaries. One numeric phase
has two internal groups and an independent checkpoint; no 21A/21B exists.

## Binding Decisions

- DEC-152: Phase 20 retrospective controls task size, fresh feasibility,
  task-local report lint, three-task review cadence, and checkpoint timing.
- DEC-153: Gate D3.2 DP-1..DP-7 remain accepted and are not reopened.
- DEC-154: exact initial slice/schema contract.
- DEC-155: carry-forward disposition.
- DEC-156: gate evidence guards terminal-success transition before completion;
  both use one shared Obligation evaluator.
- DEC-157: complete wire contract, kebab-case public enums, role-free bootstrap,
  exhaustive V1 disposition, achievable checkpoint, and tasks 19-22 order.

## Planning Authorities

- Slice semantics/Event fixtures: `../phases/phase-21/contracts.md`
- Complete public interface: `../phases/phase-21/public-interface.md`
- Live V1 inventory: `../phases/phase-21/v1-convergence-inventory.md`
- Validation: `../phases/phase-21/qa/validation-contract.md` (VAL-001..087)
- Checkpoint: `../phases/phase-21/qa/group-1-checkpoint-contract.md`
- QA request: `../phases/phase-21/qa/qa-lead-handoff.md`
- Planning audit: `../phases/phase-21/planning-audit.md`
- Tasks: `../phases/phase-21/tasks/task-<N>-*.md`

## Group 1 - Authority Foundation

| Task | Title | Depends on | VAL ownership |
|---:|---|---|---|
| 1 | 0013 schema/constraints | approval only | 001 |
| 2 | Event v2, idempotency, role-free bootstrap | 1 | 003-008 |
| 3 | capability discovery | 1-2 | 009-010 |
| 4 | Activation trust | 2-3 | 011-015 |
| 5 | Gate Registration | 4 | 016-018 |
| 6 | Envelope values | 5 | 019-022 |
| 7 | generic authority | 4,6 | 023-024 |
| 8 | Issue Actor/Event convergence | 2,7 | 025-027 |
| 9 | retained attribution/judgment freeze | 7-8 | 028-031 |
| 10 | generic Claims | 6-7 | 032-034 |
| 11 | generic Runs | 10 | 035,037 |
| 12 | shared evaluator/nonterminal transition | 7,10-11 | 038-041 |

QA Reviews follow tasks 1-3, 4-6, 7-9, and 10-12. The independent checkpoint
owns VAL-042/043. Temporary replacement-bound V1 state may remain exactly as
inventoried; PASS proves new generic paths do not consult it. Task 13 is
blocked without checkpoint PASS and fresh Group 2 feasibility.

## Group 2 - Trusted Execution And Convergence

| Task | Title | Depends on | VAL ownership |
|---:|---|---|---|
| 13 | portable executor policy | checkpoint | 044-046 |
| 14 | live framing/redaction/controls | 13 | 036,047-050 |
| 15 | Attempt/Receipt | 13-14 | 051-054 |
| 16 | Artifacts | 15 | 055-057 |
| 17 | Evidence evaluation | 15-16 | 058-063 |
| 18 | terminal guard/completion/resume | 12,17 | 064-068 |
| 19 | export create/validate | 18 | 069-071 |
| 20 | V1 removal and migration 0014 | 1-19 | 002,074-075 |
| 21 | final-schema init/maintenance/rebuild/export parity | 20 | 076-083 |
| 22 | final interface/journey conformance | 21 | 072-073,084-086 |

QA Reviews follow tasks 13-15, 16-18, and 19-21. Final QA owns task 22
evidence and phase-end VAL-087. VAL-036 is the sole Human checklist assertion.

## Journey Order

Mutate once: init/Issue; discover; activate; register; create/activate Envelope;
resolve Actor; Claim/Run; nonterminal transition; execute Gate; persist/link
Receipt; evaluate Evidence; terminal-success transition; finish/release;
complete; resume; export/validate. Delete the projection and rebuild, then
repeat only named reads, evaluation, completion state, resume, export, and
negative checks. Never create a second Envelope revision.

## Approval Gate

1. Human reviews and approves or amends this complete Heavy planning pack.
2. QA Lead proposes a **Ready** Heavy QA plan from the handoff.
3. Tech Lead applies any QA-driven contract/task command refinements.
4. Human approves execution. Only then may a Phase Runner dispatch task 1.

No protected Human decision remains open inside the proposed contract. The
pending action is execution approval, not a reopened DP-1..DP-7 decision.
