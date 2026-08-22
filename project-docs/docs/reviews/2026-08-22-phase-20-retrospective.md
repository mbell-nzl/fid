# Phase 20 Deep Retrospective - Temporary Delivery Workflow

**Date:** 2026-08-22. **Scope:** the deep retrospective requested by the
Human and recommended by Phase 20 DEC-151 before Phase 21 planning. This
reviews the repository's temporary V1 delivery workflow, not fid product
behavior.

## Evidence Base

- Phase 20 closed with 46/46 validation assertions PASS, deterministic QA
  PASS, and Human verification 3/3 PASS.
- The phase allocated 23 unique `TRG-20-NNN` identifiers and 19 escalations.
- `decisions.md` DEC-130..150, `triage-log.md`, `tl-escalations.md`, the
  Task 3e-3l infrastructure pack, and all final reviews were inspected.
- The earlier stoppage analysis and workflow-reliability audit were treated as
  evidence, not assumed conclusions:
  `2026-08-11-developer-stoppage-analysis.md` and
  `2026-08-15-workflow-reliability-overhaul-audit.md`.

The product boundary was sound. The dominant cost came from verification and
delivery controls being specified later than, more narrowly than, or
differently from the artifacts they judged.

## Findings

### 1. Strict evidence worked; late evidence design did not

Independent review found missing exact-schema matrices, config boundary cases,
transition failure states, and the absent real-process Phase 20 journey. Those
were real evidence gaps. Keeping both-state assertions, independent reruns,
and immutable failed attempts prevented a green full-suite command from being
mistaken for proof of an untested state.

The avoidable part was timing. Report syntax, termination annotations,
expected-nonzero grammar, and full-output rules were corrected after task
reports already existed. Final QA then had to lint 22 historical reports.
Phase 21 must make the owning task the first and normal failure boundary for
report conformance.

### 2. Unexecuted verifier instructions became protected recovery work

DEC-143..150 repeatedly corrected commands or control directives whose claimed
pass state had not been shown reachable before publication. The YAML module
absence command, `go mod why` output shape, Task 3k timeout, and sequence-lint
red fixtures all failed because the prescribed evidence procedure was wrong or
under-budgeted, not because the protected assertion was wrong.

The standing protection around gates remains correct. The planning correction
is to validate every non-default verifier, fixture, expected non-zero result,
and wall-clock assumption before it becomes binding. Heavy-mode feasibility
cannot be a prose review.

### 3. Several task boundaries split one green build across two owners

Task 14 removed a public flag while Task 15 owned callers that Task 14 needed
to compile. Task 9 compressed a wide three-state/four-state evidence matrix
into one-line criteria. The infrastructure pack then required eight recovery
tasks because test ownership, binary construction, workspace setup, process
working directory, parallelism, and time were coupled.

The durable planning rule is causal ownership: a task that removes or changes
a public field, flag, symbol, schema object, or Event owns every caller and
fixture invalidated by that change. Large matrices are sliced by one behavior
or trust boundary, not by arbitrary file groups.

### 4. Typed recovery improved routing and should remain

The `RESUME | REPAIR | BLOCK` overhaul correctly distinguished incomplete
in-scope work from protected decisions. Task 10's ordinary incomplete
migration, bounded QA rework, report correction, and control-conformance
repair no longer need to be interpreted as product decisions.

Phase 20's later control repairs were expensive, but that is not a reason to
return to binary stop/resume routing. The problem was the novelty and
unvalidated complexity of the controls. Preserve strict evidence and typed
recovery; reduce how often recovery is needed through earlier feasibility.

### 5. Review cadence found product gaps; closure cadence found them late

Six product review windows were justified and effective. The inefficiency was
that report/control conformance accumulated outside those windows and was
rediscovered during final QA. Phase 21 is larger and security-sensitive, so
reviews must stay frequent but narrow: every three Developer tasks, plus one
fresh independent review and executor-feasibility checkpoint between the two
internal task groups.

## Decisions Applied To Phase 21

1. **Task size:** one Developer task owns one coherent behavior, schema
   transition, security boundary, or removal closure. Removal tasks own all
   invalidated callers. No task may defer compile- or gate-required fallout to
   the next task.
2. **Verifier feasibility:** the QA plan must freshly validate every command,
   fixture, special environment assumption, expected non-zero result, live
   execution framing probe, and wall-clock budget before execution. A new or
   changed procedure without a demonstrated reachable pass state blocks the QA
   plan.
3. **Report timing:** every task spec includes the canonical task gate and the
   exact `report-lint.sh` invocation for its own completion report. The next
   task does not start until that report is committed, clean, and linted.
4. **Review cadence:** independent QA review follows every three completed
   Developer tasks. A mandatory Group 1 checkpoint follows the authority,
   activation, Gate Registration, Work Envelope, Claim, and Run foundation and
   must return PASS before trusted execution/evidence work begins.
5. **Checkpoint content:** the Group 1 checkpoint independently verifies the
   accepted JSON/DDL/Event contracts, allow/deny authority states, activation
   source trust, one-Issue limits, registration-before-Envelope ordering,
   Claim/Run bindings, rebuild parity, and fresh standard-library executor
   feasibility for argv, cwd, environment, timeout, output bounds, signal
   observation, redaction, and contained Artifact storage.
6. **Control policy:** keep immutable failed evidence, byte-identical reruns,
   typed recovery, independent re-review, and protected control definitions.
   Do not port report-lint grammar or Markdown workflow aggregates into fid
   product scope.
7. **CF-029/CF-030:** do not mix a workflow-tooling feature into the Phase 21
   product slice. CF-029 remains open with manual QA review-token auditing in
   Phase 21; CF-030 remains trigger-based for the next `report-lint.sh` change.
   Both expire with the temporary V1 delivery workflow rather than becoming V2
   product requirements.

## Phase 21 Calibration

- Operating mode is Heavy because Phase 21 changes public JSON, schema,
  canonical Events, authority, process execution, redaction, and Artifact
  storage.
- One numeric Phase 21 is retained. Two internal task groups and one checkpoint
  preserve ADR-0012's vertical-slice closure while creating a hard feasibility
  boundary before the security-sensitive half.
- Group 1 is complete only as a reviewed foundation, not as a separately
  accepted product phase or capability-adoption point.
- No Phase 21 capability is marked available unless its declared semantic
  contract and limits are complete and independently evidenced. Partial
  implementation remains explicitly unavailable.

## Outcome

The deep retrospective is **COMPLETE**. It changes Phase 21 task sizing,
verifier feasibility, report-lint timing, and review cadence. It does not
weaken a Phase 20 assertion or alter Phase 20's accepted product outcome.

