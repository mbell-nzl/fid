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
| drill-05 | weakened/vacuous test | QA Reviewer |
| drill-06 | unrecorded spec deviation | QA Reviewer |
| drill-07 | stale doc claim | QA Reviewer |
| drill-08 | VAL PASS on positive-only evidence | QA Reviewer / QA Lead |

## Known Uncovered Defect Class

**A required command that ran but never terminated, recorded `PASS`.** Phase 18
produced this defect (DEC-072 / QA-BLK-18-FINAL-01): a `workflow.Load`
connection deadlock left four required gate commands non-terminating, and the
completion report recorded them `PASS` anyway. Only QA Lead's independent rerun
caught it; the mid-phase review cadence structurally cannot, because a report
that claims a result the reviewer does not re-execute is indistinguishable from
one that earned it.

The existing set does not cover this: drill-03 is a *silently skipped* required
command and drill-04 a *narrower substituted* one — both detectable from the
report text alone. A fabricated `PASS` is only detectable if the report is
required to carry evidence the run itself produced.

The completion-report template and Developer prompt now require a verbatim
terminal-evidence token per command (a hung command has no final line to copy).
That is an affordance, not a detector. A `drill-09` for this class needs
`report-lint.sh` to check that each required command's result line carries an
evidence token, and must be demonstrated failing-then-passing per Extending
below. Not yet built — see
`project-docs/docs/reviews/2026-07-26-workflow-retrospective.md` (DEC-087).

## Running

Mechanical drills (00-04), scored automatically; exit code is the miss count:

```sh
project-docs/tools/drills/run-drills.sh
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
