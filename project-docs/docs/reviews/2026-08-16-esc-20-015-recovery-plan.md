# ESC-20-015 recovery and workflow-control repair plan

**Date:** 2026-08-16  
**Owner:** Recovery super-agent  
**Status:** In progress  
**Baseline:** `phase-20` at `7c74edc`; the append-only Phase 20 ledger is
`STOPPED` at `TRG-20-015` / `ESC-20-015`.

## Objective

Close ESC-20-015 without changing Task 17's command, expected exit, evidence,
VAL-034, threshold, or assertion strength. Repair `report-lint.sh` so literal
commands are compared byte-for-byte, then correct the V1 recovery workflow so
a bounded defect in a canonical control is recoverable without a Human stop
when the repair only restores an already accepted contract.

## Diagnosis

Task 17's expected-nonzero declaration and report row carry the same
backslash-bearing command. `report-lint.sh` transports that command through
`awk -v`; awk interprets `\.` before equality comparison and false-rejects the
truthful evidence.

The Phase Runner correctly retained the failed evidence and entered recovery,
but its only repair dispatch is a same-task Developer continuation. Task 17
does not own the workflow tool. TL-Lite then conflated repair of a control's
implementation with a change to the protected gate definition and emitted
`BLOCK`. The workflow has no bounded cross-scope control-repair route.

## Guardrails

- Preserve the existing STOP, TRIAGE, report failure, and Task 17 evidence as
  immutable history.
- Do not alter Task 17's exact command or its recorded PASS row.
- Do not change a gate command, definition, threshold, accepted status,
  assertion strength, dependency, product behavior, or Phase 20 task boundary.
- A control-conformance repair is eligible only when an accepted contract
  already determines the result and the repair adds no new policy judgment.
- A dedicated recovery session owns only the named canonical control, its
  regression tests, and required workflow documentation. It remains linked to
  the originating task and consumes that trigger's recovery budget.
- Any ambiguity about accepted semantics, any weakening/narrowing, or any
  protected contract change remains `BLOCK`.

## Work Packages And Progress

- [x] **P0 - Preserve and diagnose.** Inspect the stopped ledger, Task 17
  report/spec, parser call sites, current routing contract, and uncommitted
  Phase Runner artifacts. Reproduce the failure mechanism and record this plan.
- [x] **P1 - Byte-preserving report-lint repair.** Added red self-tests for
  backslash-bearing exact commands in expected-nonzero lookup and required
  command coverage. Replaced unsafe literal transport with byte-preserving
  comparison; focused verification and the unchanged Task 17 lint invocation
  pass. Committed as `88b0377`.
- [x] **P2 - Cross-scope recovery contract.** Amended canonical workflow spec,
  Phase Runner, TL-Lite, and Developer instructions to distinguish gate
  definition changes from control-conformance repairs and authorize a dedicated
  bounded control-repair session linked to the originating task.
- [x] **P3 - Mechanical enforcement.** Extended workflow lint and drills so the
  distinction, ownership, budget linkage, no-advance rule, and protected
  counter-cases cannot silently regress. Generated prompts are synchronized;
  workflow commit `92cb0f4` passes lint, drills, and sequence validation.
- [x] **P4 - Phase 20 recovery closure.** Recorded the authoritative resolution,
  archive ESC-20-015, append validated RUN-RESUME / recovery events, finalize
  Task 17's report from unchanged evidence at `0c1094c`, and left the ledger
  `ADVANCEABLE` with Review 6 next.
- [ ] **P5 - Independent adversarial review.** Ask an agent that did not
  implement the repair to probe literal edge cases, unsafe gate weakening,
  cross-scope ownership, recovery budgets, ledger sequencing, and generated
  prompt drift. Initial verdict `REPAIR`: the reviewer reproduced missing
  schema enforcement, self-authorized/unbudgeted recovery, stale/dirty/symlink
  scope-check bypasses, untriaged STOP/resume, and literal-tab command
  false-rejection. Commit `f71abe6` closed that first set, but the reviewer
  then reproduced five deeper bypass classes: self-attested TRIAGE/decision
  authority, non-resolved repair completion and budget reset, incompatibility
  between the dirty runner ledger and a clean repair verifier, STOP/resume
  identity loss, and hidden index flags. The second repair round now requires a
  dedicated clean temporary repair worktree, binds strict TRIAGE to the real
  phase TRG/verdict and STOP/resume identity, verifies cited decisions, rejects
  `assume-unchanged`/`skip-worktree`, and preserves the monotonic two-attempt
  no-progress budget. The mechanical drills exercise those negative cases and
  the actual dirty-main/clean-repair/fast-forward lifecycle. Obtain `PASS` on
  re-review before closing this item. A third adversarial pass against
  `c133cb0` found transactional-append authority loss, delayed schema
  activation, prefix-matched and self-selected TRIAGE authority, premature
  control completion, control-scope escape, unbound STOP/resume and recovery
  budgets, partial DEC validation, allowlist injection, and caller-selected
  repair bases. Commit `173c484` replaces those loose transitions with
  externally bound schema-v2 directives, Runner-owned base/head checkpoints,
  exact report-only continuation, monotonic control retries, full DEC and
  resume identity checks, transactional authority preservation, and expanded
  adversarial drills. Canonical prompt/spec alignment and final independent
  re-review remain in progress.
- [ ] **P6 - Exact verification and handoff.** Run all workflow gates and the
  relevant repository gate, capture exact termination, commit the final review,
  and provide a copy-paste Phase Runner restart prompt with branch, HEAD,
  ledger state, and first required action.

## Required Verification

- `sh -n project-docs/tools/report-lint.sh`
- `project-docs/tools/report-lint.sh --self-test`
- Task 17's unchanged report-lint invocation
- `project-docs/tools/workflow-lint.sh`
- `project-docs/tools/drills/run-drills.sh`
- `project-docs/tools/sync-agents.sh --check`
- `project-docs/tools/workflow-sequence-lint.sh project-docs/docs/phases/phase-20/run-ledger.md`
- `git diff --check`
- Focused or broader product gates only where evidence freshness requires them;
  Task 17's already captured product gate results will not be falsely
  re-attributed to this recovery session.

## Commit Strategy

1. Preserve the STOP artifacts and this plan as the recovery baseline.
2. Commit the report-lint behavior repair with its red/green regression.
3. Commit the workflow authority/routing contract and generated prompt sync.
4. Commit Phase 20 recovery closure and Task 17 report finalization.
5. Commit the independent review and any resulting corrections.

Each checkpoint updates this file before its commit so it remains the durable
anchor across compaction or a fresh recovery session.
