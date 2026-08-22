# Workflow Reliability Overhaul Audit And Architecture Decision

Date: 2026-08-15  
Status: ACCEPTED FOR IMPLEMENTATION by Human authorization on 2026-08-15  
Scope: binding fid V1 delivery workflow and V2 adoption guidance  
Independent inputs: stoppage-history audit, V1 workflow-map audit, and V2
starter-kit adoption audit performed by separate read-only agents.

## Executive Finding

The workflow detects missing evidence and real failures reasonably well, but it
routes almost every unresolved condition as a Human-visible stop. It lacks a
first-class transition for work that is real, local, in scope, and recoverable.

Across Phases 17-20, the independent audit found 31 formal Phase Runner STOPs,
plus Phase 20 ESC-20-004 without a separate STOP:

| Primary classification | Formal STOPs |
|---|---:|
| Protected product or architecture decision | 5 |
| Recoverable implementation state | 6 |
| Evidence or report defect | 15 |
| Orchestration defect | 4 |
| External infrastructure fault | 1 |

Only four formal stops required direct Human judgment: the Phase 17 model
qualification entry, the Phase 17 final Human checklist, Phase 19 kickoff
approval, and the Phase 19 final Human checklist. No implementation failure,
gate failure, report defect, or incomplete task required Human judgment.

The problem is therefore routing, not insufficient strictness. Evidence must
remain strict; ordinary recovery must stop paging the Human.

## Decisive Incident: ESC-20-010

The Task 10 Developer session terminated normally after 616 seconds. It had no
abort, cancellation, interrupt, timeout, context-exhaustion, or kill event. The
Developer correctly described the focused failures as expected compatibility
fallout, made in-scope test migrations, then did not rerun after its final patch
and voluntarily wrote a blocker report.

The remaining work was explicitly owned by Task 10. The focused suite named
legacy channel and flag assertions that needed migration. No product,
architecture, schema, public-interface, dependency, or security question was
open. ESC-20-010 was a recoverable implementation state.

The Phase Runner then treated the child-authored blocker as authoritative. That
exposes the core control defect: an agent's "blocked" response is an observation
to classify, not a lifecycle decision.

Evidence:

- Session:
  `/home/mark/.codex/sessions/2026/08/15/rollout-2026-08-15T08-56-09-01a0020f-8f14-7770-ae5e-d5010bd1e7a9.jsonl`
- Preserved Task 10 tree and report: stash object `779beab1`.
- Task: `project-docs/docs/phases/phase-20/tasks/task-10.md`.

## Root-Cause Analysis

### 1. Binary triage has no recovery transition

TL-Lite currently has only `RESUME` and `BLOCK`:

- `RESUME` ambiguously means either "the trigger is cleared; advance" or "a
  bookkeeping repair is allowed but has not yet been performed."
- `BLOCK` is the only valid result for a genuine unresolved local failure.
- TL-Lite cannot edit product files or run gates.
- The Phase Runner has no mandatory dispatch back to the same task after a
  locally authorized repair.

This either advances incomplete work or, more commonly, forces a STOP.

### 2. Same-task repair is required and obscured simultaneously

`09_AI_WORKFLOW.md` tells the Developer to fix a failed gate and rerun the exact
command. `AGENTS.md` stops only when the failure cannot be fixed inside assigned
scope. However, detailed exception language explicitly permits only authored
harnesses and stale bookkeeping inventories. Agents infer that ordinary product
test failures cannot be repaired after a gate begins.

DEC-138's bookkeeping category mitigated symptoms while reinforcing the false
premise. It is superseded by the general causal-repair rule below.

### 3. QA severity is conflated with authority

A correct QA finding may be "must fix before sign-off" while still having a
bounded, one-task, decision-free repair. The current workflow routes every
Blocking finding to binary triage, where a real unresolved finding must BLOCK.
There is no prescribed Developer recovery and independent re-review loop.

### 4. Report state is inferred from prose

`report-lint.sh` recognizes a blocker only by a "No commit created" phrase.
Partial work with a checkpoint commit can look complete, while a normal
completion and a protected blocker can both exit zero. The linter supplies no
typed outcome for the Phase Runner.

### 5. Command termination is canonical but not mechanically checked

Canonical instructions require an observed exit code or explicit termination
signal, but report-lint's accepted fixture omits it. The Phase Runner must catch
the omission by prose reading and route it through the same broken stop path.

### 6. Durable continuation is underspecified

"One task per Developer session" was interpreted as "one Developer session per
task." A task should not cross into another task in the same session, but a
task may require multiple continuation sessions. The run ledger lacks explicit
recovery and resolved-resume events.

### 7. The generic stop section contradicts role-specific routing

The detailed Phase Runner section triage-gates report and gate failures while
the generic stop section describes several as immediate Human stops. The most
conservative interpretation wins, producing repeated unnecessary stops.

## Architecture Decision: Three Workflow Dispositions

The binding workflow uses three distinct dispositions:

| Disposition | Meaning | Required Phase Runner action |
|---|---|---|
| `RESUME` | Trigger is cleared, demonstrably misclassified, or already resolved with verified evidence. | Advance from the current checkpoint. |
| `REPAIR` | Trigger is real, local, inside existing authority, and decision-free. | Do not advance. Record recovery, dispatch the same task to a Developer continuation, rerun required verification, and repeat any affected independent review. |
| `BLOCK` | A protected decision, unsafe condition, unrepairable external dependency, missing/contradictory contract, or exhausted bounded recovery prevents progress. | Record triage, mint/reference an ESC, append STOP, and request the named authority or external action. |

`RESUME` never means "advance now and repair later." `REPAIR` never closes the
task or finding. `BLOCK` is not a synonym for red, incomplete, failed, or must
fix.

### Task lifecycle

```text
ACTIVE -> RED -> REPAIR -> VERIFY -> REPORT -> COMPLETE
             ^       |         |
             +-------+---------+

COMPLETE -> REVIEW -> REPAIR -> VERIFY -> REPORT -> REVIEW

Any nonterminal state -> BLOCKED only for a protected condition or an
exhausted no-progress recovery budget.
```

A failed invocation remains failed evidence. A later passing invocation after a
causal repair does not rewrite the first attempt. Both are recorded, and the
required command is rerun byte-identically unless its accepted definition
changed through protected authority.

## Recovery Policy

### Recovery states

These conditions are `REPAIR` unless accompanied by a separate protected
condition:

- Incomplete in-scope implementation or test migration.
- A focused or full gate exposing an in-scope defect.
- A stale test expectation caused by the authorized contract change.
- Completion-report, commit-metadata, inventory, generated-copy, or evidence
  formatting repair.
- A long-running command that remains live and must be polled.
- A session ending while durable in-scope work remains.
- A bounded QA finding owned by an existing task with unchanged acceptance
  criteria.
- A local sandbox/cache/scratch-state failure with a non-destructive recovery.

### Protected blocker categories

`BLOCK` requires one of these closed categories:

- `product-decision`
- `architecture-decision`
- `data-model-decision`
- `public-interface-decision`
- `dependency-approval`
- `security-or-data-loss`
- `external-infrastructure`
- `specification-contradiction`
- `recovery-budget-exhausted`

The report must name the exact condition, evidence, attempted recovery,
resolution condition, and why recovery cannot continue under current task
authority. "Task incomplete," "tests red," "gate failed," "report lint failed,"
or "several files remain" are invalid blocker reasons by themselves.

### Bounded recovery

- Same-session red/green work does not consume an orchestration recovery
  attempt.
- A Phase Runner-mediated fresh continuation consumes one recovery attempt for
  its trigger signature.
- Default budget: two consecutive continuation attempts for the same trigger
  with no material progress.
- A changed tree, newly satisfied obligation, or a new narrower failure is
  progress and must be recorded; merely rerunning unchanged failing work is not.
- Exhaustion routes through TL-Lite as `recovery-budget-exhausted`; it is not
  silently converted into a product decision.
- Context refresh or planned handoff never resets the budget.

### Session continuity

One Developer session handles at most one task. One task may span multiple
Developer continuation sessions. A continuation must read the durable task
spec, report/checkpoint, run-ledger recovery entry, relevant review, and current
tree before editing. It preserves prior failures and continues from existing
work without reverting unrecognized changes.

## QA Recovery Policy

Review outcome and finding severity are separate:

- `PASS`: no must-fix finding remains.
- `REPAIR`: one or more valid must-fix findings have bounded, decision-free
  ownership in existing task scope.
- `BLOCK`: at least one finding needs protected authority or exceeds the review
  window/task scope.

For `REPAIR`, the Reviewer names the owning task, affected assertion, exact
repair boundary, and re-review evidence. The Phase Runner dispatches a Developer
continuation for that task and then spawns a fresh independent Reviewer for the
same review window. The original finding remains in history and is resolved
only by the re-review.

QA roles emit blocker candidates. TL-Lite or the full Tech Lead owns ESC
allocation; QA roles do not mint escalation identities.

## Durable And Mechanical Controls

V1 will add:

1. A typed completion-report workflow outcome:
   `COMPLETE | RECOVERY_REQUIRED | BLOCKED`.
2. Closed recovery/blocker categories and structured condition, attempt, next
   action, and Human-attention fields.
3. `report-lint.sh` outcome output and termination-annotation enforcement.
4. Run-ledger events for `TASK-RECOVERY`, `REPAIR`, and `RUN-RESUME`, including
   trigger signature, attempt count, checkpoint, and resolution evidence.
5. Three-verdict concordance and forbidden-legacy-semantics checks in
   workflow-lint.
6. Mechanical drills for report outcomes and state transitions, plus agent
   drills for premature blocker classification and QA recovery.

The V1 Markdown/report mechanisms remain temporary delivery infrastructure.
They are not forward fid product entities.

## Historical Classification Summary

| Incident group | Correct handling under this decision |
|---|---|
| Phase 17 report-lint stops | `REPAIR` report/commit metadata, then lint again. |
| Phase 18 missing QA output | Retry/fresh QA session; `BLOCK` only after no-progress budget or missing contract. |
| Phase 18 performance failure and deadlock | Same-task causal `REPAIR`, exact gates, independent review. |
| Phase 19 evidence ceremony defects | `REPAIR`; convert repeated manual observable to deterministic automation. |
| Phase 20 ESC-001/002/004/009 | QA `REPAIR` with owning task and mandatory re-review. |
| Phase 20 ESC-003 | Orchestration recovery; retain/poll live session, never infer termination. |
| Phase 20 ESC-005/006/007/008 | In-scope code/report/control `REPAIR`; no Human stop. |
| Phase 20 ESC-010 | Continue Task 10; no protected blocker exists. |

This does not declare all historical findings false. Most detection was valid.
It changes the action after a valid, recoverable finding.

## V2 Adoption Guidance

V2 is directionally aligned but not adoption-ready. Its accepted contracts
already describe scoped gate repair, linked-Run interruption recovery, and
same-Issue review rework. Its current harness still implements the superseded
WorkItem model, production roles remain fixture-only, and all required fid
capabilities remain unqualified.

Before V2 adoption:

1. Resolve the contradiction between globally non-actionable open Stops and the
   Builder action required to repair a retryable `gate-failed` condition. A
   gate stop should block completion while permitting an authorized repair Next
   Action on the same Issue/Envelope.
2. Replace `retryable: bool` with a richer recovery disposition: affected
   transition, actor, authority, same-Run eligibility, Human attention, and
   budget class.
3. Make same-Issue causal repair a first-class path with immutable failed and
   passing Attempts/Receipts.
4. Implement actual linked-Run recovery with executor-liveness proof and
   non-destructive repository reconciliation.
5. Implement independent review rework and budget-exhaustion scenarios.
6. Treat a Builder "blocked" response as an observation classified by the
   harness, never authoritative lifecycle state.
7. Complete trusted long-running execution using provider-neutral handles and
   terminal Receipts.
8. Finish V2-2 correction, twelve capability qualifications, V2-3 evidence and
   recovery conformance, V2-4 provider/role qualification, and V2-6 comparative
   dogfood before adoption.

Do not port V1's report-lint grammar, fixed Markdown headings, two-commit report
ceremony, tool-specific `session_id` instructions, Phase Runner/TL-Lite roles,
TRG/ESC routing, or manual stash mechanics. Port the semantic principles:
immutable attempts, decision-free scoped repair, exact registered gates,
silence is not termination, canonical-state recovery, independent finding
resolution, and Human attention only for protected or exhausted conditions.

## Baseline Verification

Before implementation:

- `project-docs/tools/report-lint.sh --self-test` passed.
- `project-docs/tools/drills/run-drills.sh` passed all five mechanical drills.
- `project-docs/tools/workflow-lint.sh` failed only its existing current-state
  line budget: `tl_escalations_current.md` was 134 lines against a 107-line
  tripwire. This confirms current lint checks vocabulary and generated drift,
  not transition coherence.

The implementation must preserve the two passing controls, make the state
semantics mechanically visible, and return workflow-lint to green without
raising the line-budget threshold.
