# Repo Tools

Stack-neutral repository commands that do not require project runtime
dependencies.

| Command | Purpose |
|---|---|
| `sync-agents.sh` | Generate provider role prompts from `project-docs/agents/`. Use `--check` to compare without writing. |
| `sync-skills.sh` | Generate provider skill copies from `project-docs/skills/`. Use `--check` to compare without writing. |
| `workflow-lint.sh` | Verify agent workflow artifacts, templates, config, and generated prompt/skill copies have not drifted. |
| `report-lint.sh <report> <task-spec>` | Deterministically verify a completion report: typed workflow outcome, required headings, command coverage, observed termination, commit identity, and that the report is committed and clean. Emits `OUTCOME COMPLETE`, `OUTCOME RECOVERY_REQUIRED`, or `OUTCOME BLOCKED`. `--self-test` proves pass and fail cases. |
| `workflow-route.sh <source> <outcome> [attempts]` | Map typed Developer, report, QA, triage, command-liveness, and recovery-budget outcomes to the Phase Runner's next action. |
| `workflow-sequence-lint.sh <run-ledger> [--append-stdin]` | Validate current or proposed same-task recovery ordering and independently linked QA re-review; reports state but never chooses the next action. |
| `drills/run-drills.sh` | Evidence-chain drills: seed known defects into a fixture mini-phase and verify the QA chain catches them. Mechanical drills score automatically; agent drills produce workspaces plus a scoring sheet. |

Run from the repository root:

```sh
project-docs/tools/sync-agents.sh
project-docs/tools/sync-skills.sh
project-docs/tools/workflow-lint.sh
project-docs/tools/report-lint.sh --self-test
project-docs/tools/drills/run-drills.sh
```
