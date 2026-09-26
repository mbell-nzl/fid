# Evidence-Chain Drills

Mutation testing for the QA process itself. The harness builds a synthetic
mini-phase (task spec, validation contract, source, tests, completion report,
real git history), seeds one known defect per drill, and verifies the QA chain
catches it. Any miss is a workflow defect with a named location - a prompt,
template, or spec rule to fix.

All evidence that the QA chain works is otherwise historical. Models drift,
prompts rot, and new models fail differently; the drills detect guardrail
decay before a real defect exploits it.

## Drill Set

| Drill | Defect class | Expected detector |
|---|---|---|
| drill-00 | control: clean report (must pass) | `report-lint.sh` |
| drill-01 | fabricated commit hash | `report-lint.sh` |
| drill-02 | plausible wrong hash (resolves, wrong commit) | `report-lint.sh` (subject match) |
| drill-03 | silently skipped required command | `report-lint.sh` |
| drill-04 | narrower substituted command | `report-lint.sh` |
| drill-09 | malformed PASS: missing or nonzero termination | `report-lint.sh` |
| drill-10 | typed recoverable report | `report-lint.sh` |
| drill-11 | typed protected blocker | `report-lint.sh` |
| drill-12 | recoverable condition presented as BLOCKED | `report-lint.sh` |
| drill-13 | in-scope failure next action | `workflow-route.sh` |
| drill-14 | invalid report next action | `workflow-route.sh` |
| drill-15 | live and missing-termination command handling | `workflow-route.sh` |
| drill-16 | bounded QA finding next action | `workflow-route.sh` |
| drill-17 | protected stop and no-progress budget | `workflow-route.sh` |
| drill-19 | parent wait expiry and externally interrupted child accounting | `workflow-route.sh`, `workflow-sequence-lint.sh` |
| drill-20 | decision-released repair owners after a BLOCK re-review (DEC-257) | `workflow-sequence-lint.sh` |
| drill-05 | weakened/vacuous test | QA Reviewer |
| drill-06 | unrecorded spec deviation | QA Reviewer |
| drill-07 | stale doc claim | QA Reviewer |
| drill-08 | VAL PASS on positive-only evidence | QA Reviewer / QA Lead |

## Recovered Defect Class

**A required command that ran but never terminated, recorded `PASS`.** Phase 18
produced this defect (DEC-072 / QA-BLK-18-FINAL-01): a `workflow.Load`
connection deadlock left four required gate commands non-terminating, and the
completion report recorded them `PASS` anyway. Only QA Lead's independent rerun
caught it; the mid-phase review cadence structurally cannot, because a report
that claims a result the reviewer does not re-execute is indistinguishable from
one that earned it.

Drill 09 now proves that a PASS row without observed termination is rejected.
Drill 15 separately proves that a live command routes to `WAIT`, while missing
termination routes to `REPAIR`. Together they cover both the evidence defect
and the orchestration behavior that caused ESC-20-003. See
`project-docs/docs/reviews/2026-07-26-workflow-retrospective.md` (DEC-087).

Drill 19 proves that a parent polling timeout retains a live child, an external
interruption does not consume or reset the no-progress budget, and two
consecutive external interruptions require orchestration session-liveness
triage with the interruption cause, elapsed duration, and command in flight.

## Running

Mechanical outcome and transition drills, scored automatically; exit code is
the miss count:

```sh
project-docs/tools/drills/run-drills.sh
```

The suite also runs `review-recovery.sh`, the DEC-190 regression for the
Phase 22 restart deadlock. It exercises canonical append, same-task triage
correction, active review identity, STOP/session recovery, rejection without
mutation, and idempotent receipt retry. Its focused command is:

```sh
sh project-docs/tools/drills/review-recovery.sh
```

The suite also runs `review-block-owners.sh`, the DEC-257 regression for the Phase 23
STOP after a BLOCK re-review. It proves that a STOP resume can bind the
ordered `Repair owners` line of its resolving DEC. The owners then recover one
at a time, and a linked PASS closes the review. It also rejects every drift:
a mismatched or duplicated list, a STOP that was not a BLOCK re-review, a
wrong first owner, out-of-order recovery, and an early PASS:

```sh
sh project-docs/tools/drills/review-block-owners.sh
```

Full set including agent drills (05-08):

```sh
project-docs/tools/drills/run-drills.sh --agent
# or pin the workspace location:
project-docs/tools/drills/run-drills.sh --workspace /path/to/drill-run
```

`--agent` materializes one workspace per agent drill under
`<workspace>/agent-drills/<drill>/repo/` plus a `DRILL.md` scoring sheet
outside the repo directory, and a `SCORECARD.md` at the workspace root.

For each agent drill:

1. Spawn the role named in `DRILL.md` (QA Reviewer or QA Lead) against
   `repo/` only. Never show `DRILL.md` to the agent under test - it names the
   seeded defect.
2. Score DETECTED or MISS per the drill's criteria and record it in
   `SCORECARD.md`, along with the model/prompt configuration under test.

## Cadence

- Quarterly, as guardrail-decay detection.
- After slimming or rewriting any agent prompt (regression gate: detection
  must stay 8/8 before the change is trusted).
- Before trusting a new model in an orchestrator or QA role
  (role-qualification: a model that misses drills does not get the role).

## Extending

Add a new drill when a real incident exposes a defect class the chain missed:
seed the same defect here, confirm the drill fails, fix the workflow surface,
confirm the drill passes. Keep one defect per drill so a miss names its
location precisely.
