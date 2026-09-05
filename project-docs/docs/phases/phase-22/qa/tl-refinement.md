# Phase 22 Tech Lead Refinement

> Tech Lead: Codex, 2026-09-05. Decisions: DEC-184..187.
> Status: **ACCEPTED** by Human under DEC-188, 2026-09-05; execution authorized.
> Read with qa-plan.md and consistency-review.md. This is a planning disposition,
> not a new QA test result or a Phase 22 completion review.

## Readiness Disposition

All fourteen QA refinement requests are applied now, including the four marked
improvements. Both DECIDE items and every RECORD item have explicit dispositions
in `../kit-contract-mapping.md`. Heavy mode, 18 sequential tasks, 80 assertions,
six review windows, and the independent task-12 checkpoint are unchanged.

The QA plan's readiness condition was disposition of these findings. This
addendum closes that planning condition; Human accepted DEC-180/182 and
DEC-184..187 under DEC-188. Runner performs its required readiness refresh. No Developer has
started, no agent is activated, and no product/test file is edited by this TL.
The existing QA assessment and failed attempts remain historical evidence.

## Request Dispositions

| Request | Disposition and exact owning artifacts |
|---|---|
| R-01 | Apply: public-interface §1 retains --project, no rename; task 6/VAL-023 and final wire corpus. DEC-184. |
| R-02 | Apply: public-interface §2.1 enumerates 23 final profile/grant actions and incremental domain owners. Bootstrap activation operations are explicitly outside the action set. Task 9 removes run.finish action as well as leaf/Event; task 18 independently pins and exercises membership. DEC-184. |
| R-03 | Apply: contracts §§3.1/3.2 specify canonical input images, hashing, ordering, exclusions, failure states and separate reconciliation digest. Task 4 owns vectors/observation; task 12 owns real dirty unchanged/changed evidence. DEC-185. |
| R-04 | Apply: non-null Recovery and recover input IDs plus DDL NOT NULL; apply the same already-required liveness rule to abandon input. Tasks 1/5/6/8/9 own matching negative cases. Error/StopDetails IDs remain nullable to describe absent observations. DEC-184. |
| R-05 | Apply: duplicate-entity / conflict / exit 4 with exact executor_id identity and Run subject rules. Tasks 3/5 plus task 6 wire cases; same-key replay precedes uniqueness rejection. DEC-184. |
| R-06 | Apply: attention.resolve supplies the exact set of Stop disposition-containing Event IDs. Underlying facts alone, missing, extra, duplicate or wrong-subject references reject; task 2/9 and VAL-006. DEC-184. |
| R-07 | Apply: rejected recovery produces errors without Event/Stop/Attention. The three same-named Stop codes require explicit stop.raise; automatic budget/activation effects remain separately specified. Tasks 2/5/6/9. DEC-184. |
| R-08 | Apply enumerated-matrix option: VAL-056 now names 056-a..k and both states of every boundary. Actual cumulative decompression is measured; no impossible forged-ZIP bypass fixture is required. Task 13. Ceilings unchanged. DEC-186. |
| R-09 | Apply now: attention_items.code has the same seven-code CHECK as stops; task 1 exact DDL/live/replay evidence. DEC-184. |
| R-10 | Apply now: activation-version-conflict separates version failures from the retained competing-prepare activation-conflict shape. Tasks 10/11 and final corpus. DEC-184. |
| R-11 | Apply now: shared grammar explicitly excepts read leaves and positional import with no Actor/input flags. Existing owners task 6/15. DEC-184. |
| R-12 | Apply now: owner_phase is opaque scheduling metadata, not control flow. public-interface §2 and kit mapping C-13. DEC-184/187. |
| R-13 | Apply: exact one-sample rebuild command and instance identity commands added to tasks 1/16. Three-sample final rule and <30 s self-assert unchanged. DEC-186. |
| R-14 | Apply with corrected provenance: task 8 explicitly owns gate_run_test.go fixture synchronization and cancellation convergence, bounded load cases, 15-repeat focused command and two full benchmark-gate runs. Earlier failures use bounded same-task amendment, never a waiver. DEC-186. |

R-14 provenance correction is based on this session's git log -S, git blame and
git show reads: b694991 introduced timeout_ms=1 and the signal helper;
76ad781 changed only the signal case's timeout. Both predate the phase base.
The QA report's original wording and command results are retained, not rewritten.

## QA Plan Integration

The following refinements supplement existing QA evidence rows and checklist
steps; they introduce no dropped checks or replacement pass/fail semantics.

| Existing QA reference | Required refined evidence |
|---|---|
| EVID-01..03 / QA schema checks | Exact amended DDL, including required observation FK and equal code CHECKs; earlier draft DDL feasibility is not a rerun of this amendment. |
| EVID-06 / QA-06 | Exact disposition Event set with underlying-fact revalidation, positive and all negative reference cases. |
| EVID-12/17/18/21/22 | Required recovery/abandon observation references, executor reuse details and no-Stop rejection cases. |
| EVID-13/42/43 | Fixed canonical digest vectors plus real same-target dirty and changed-byte cases; no matching arbitrary hashes. |
| EVID-28 / QA-78 | Synchronized signal/timeout helpers, cancellation and bounded load, repeated focused command and two consecutive full benchmark runs; preserve the known failed baseline. |
| EVID-49 / QA-57 | Every VAL-056 case 056-a..k with permitted/rejected boundaries, measured total cap and generated-per-run resource fixtures. |
| EVID-62 / QA reference rebuild | One early sample at tasks 1/16 supplements three fresh samples at task 17/final review. |
| EVID-65 / wire corpus | Independently pinned 23-member action set, bootstrap exclusions, retired/unknown rejection and all amended error/input variants. |
| EVID-66 / QA-83 | Every C-01..13 disposition plus pin currency and as-built comparison; preserve kit qualification exclusions. |

Each task's full gate remains required. Commands needed by Developers are in
their task specs; this addendum is not a hidden extra command list. Task 17 now
also includes both benchmark runs required by QA-78. Final QA retains all
84 deterministic checklist items and all six Human fan-out checks.

## Evidence And Remaining Limits

QA's measurements remain attributed to its session. This TL did not rerun the
DDL chain, benchmark, coverage, archive probes, product tests, or quality gates,
and records no new PASS receipt. The runtime behavior and the amended DDL/digest
fixtures are UNVERIFIED in this session and must earn task/final evidence.
No required command is changed into a weaker substitute. Current QA baseline
failures remain failures; assigning a repair does not claim it has landed.

VAL coverage and command reach were inspected as planning documents: all 80
IDs retain one owning task and one final validator (77 QA Lead, 3 Human), with
no added/retired IDs. The explicit matrix/evidence sections refine those same
assertions. All rows are Ready after DEC-188; no outcome is promoted to PASS.

Kit coordinator export binding, terminalState mapping, scope/asOf, dedup-key
alignment and missing code representations remain CF-037/Phase 25 qualification
work. Phase 23/24 own their named Review/retry extensions. Three fid capability
rows becoming true do not qualify a production kit adapter.

QA's incidental documentation notes are dispositioned too: the obsolete
AGENTS.md quality §5.2 pointer is queued with CF-033's separately authorized
AGENTS maintenance; quality §8 now identifies Phase 21 as historical and names
Phase 22's review/checkpoint rules. No workflow-control or prompt edit is made.

## Approval Package

Human approval covers Heavy mode, the 18-task index in the current phase brief,
the refined validation contract and exact slice/wire/DDL proposals, including
CF-035 and DEC-184..187. Approval authorizes the normal Runner workflow; this
Tech Lead session does not dispatch it. Phase wrap-up is not applicable yet.
