# Phase Brief - Phase 21

> Status: **Active / STOPPED In Task 2 / Task 2 Recovery Authorized (DEC-162)**
> Branch: `phase-21`
> Product implementation: **AUTHORIZED BY DEC-159 UNDER THE TASK PACK**

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

- DEC-152 as refined by DEC-158: Phase 20 retrospective controls task size,
  fresh feasibility, Phase Runner two-argument post-commit lint before the next
  dispatch, three-task review cadence, and checkpoint timing.
- DEC-153: Gate D3.2 DP-1..DP-7 remain accepted and are not reopened.
- DEC-154: exact initial slice/schema contract.
- DEC-155: carry-forward disposition.
- DEC-156: gate evidence guards terminal-success transition before completion;
  both use one shared Obligation evaluator.
- DEC-157: complete wire contract, kebab-case public enums, role-free bootstrap,
  exhaustive V1 disposition, achievable checkpoint, and tasks 19-22 order.
- DEC-158: QA refinements, two-argument Runner report lint, four Human subchecks,
  and three consecutive `<30 s` reference-instance rebuild samples.
- DEC-159: Ready QA plan and Human approval accept the Heavy pack and authorize
  task 1 dispatch through the Phase Runner.
- DEC-160: role-free bootstrap remains exact; task 2 owns a nonpersistent,
  seven-name V1 policy vocabulary bridge for inventoried legacy callers only,
  and task 20 deletes it.
- DEC-161: plain task-prefixed schema-v2 identity values are canonical;
  TRG-21-003 governs exact task-2 RUN-RESUME attempt 3 without a verifier edit.
- DEC-162: retained V1 tests are not protected contracts; the task whose
  accepted change invalidates one converges it in that task without weakening.
  Task 2 gains repo-wide test scope; tasks 8, 9, and 20 inherit the rule.

## Planning Authorities

- Slice semantics/Event fixtures: `../phases/phase-21/contracts.md`
- Complete public interface: `../phases/phase-21/public-interface.md`
- Live V1 inventory: `../phases/phase-21/v1-convergence-inventory.md`
- Validation: `../phases/phase-21/qa/validation-contract.md` (VAL-001..087)
- Checkpoint: `../phases/phase-21/qa/group-1-checkpoint-contract.md`
- QA request: `../phases/phase-21/qa/qa-lead-handoff.md`
- QA plan: `../phases/phase-21/qa/qa-plan.md` (Ready at `e8a636e`)
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

## Dispatch Controls

Before each next dispatch, the Phase Runner runs
`project-docs/tools/report-lint.sh <committed-report> <task-spec>` and records
observed termination. Review artifacts are `qa/review-1.md` through
`qa/review-7.md`; checkpoint is `qa/group-1-checkpoint-review.md`; final
artifacts are `qa/final-review.md` and `qa/human-checklist.md`.

Every external schema-v2 `Originating task`, `Trigger signature`, `Attempt
count`, `Repair category`, and `Repair authority` value is a plain scalar.
Never wrap those values in Markdown backticks. Each trigger begins with its
exact originating task plus `:`.

## Journey Order

Mutate once: init/Issue; discover; activate; register; create/activate Envelope;
resolve Actor; Claim/Run; nonterminal transition; execute Gate; persist/link
Receipt; evaluate Evidence; terminal-success transition; finish/release;
complete; resume; export/validate. Delete the projection and rebuild, then
repeat only named reads, evaluation, completion state, resume, export, and
negative checks. Never create a second Envelope revision.

## Approval Record

1. Complete - Human accepted the Heavy phase shape and Gate D3.2 amendments.
2. Complete - QA Lead marked the plan Ready at `e8a636e` with no Must-fix.
3. Complete - Tech Lead confirmed no unresolved contract refinement.
4. Complete - Human approved execution on 2026-08-23; DEC-159 records it.

## Dispatch Authorization

Task 1 is complete. The ledger is valid and remains `STOPPED` at TRG-21-004;
the Tech Lead must not edit it. DEC-162 authorizes the Phase Runner to append
RUN-RESUME with governing TRG `TRG-21-004`, governing ESC `ESC-21-004`,
resolution authority `DEC-162`, trigger signature
`task-2:specification-contradiction:task-2-done.md:retained-v1-rebuild-role-contract`,
attempt count 6, and this exact next recovery:

`task-2 | in-scope-implementation | task-scope | task-2:specification-contradiction:task-2-done.md:retained-v1-rebuild-role-contract | 6`

Then dispatch only the DEC-162 task-2 continuation. Preserve every failed
required-command result, converge the retained tests under DEC-162 and
`v1-convergence-inventory.md` section 1.2, rerun both failed required commands
byte-for-byte, complete the unchanged remaining gate, commit the truthful
report, and run two-argument report lint. Task 3 remains unauthorized until
task 2 completes mechanically. No protected Human decision remains open.
