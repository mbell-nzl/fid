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
| `workflow-control-scope-lint.sh <repo> <base> <head> <allowed...>` | Verify that a control-conformance repair is the checked-out clean HEAD, exactly one tools-only commit, regular-file-only, and changes no path outside its explicit allowlist. |
| `workflow-sequence-lint.sh <run-ledger> [--append-stdin]` | Validate schema, triage/STOP authority, typed recovery budgets and ordering, or idempotently validate/append one stdin event at EOF; a byte-identical retry confirms a lost receipt without duplication. Reports state but never chooses the next action. |
| `drills/run-drills.sh` | Evidence-chain drills: seed known defects into a fixture mini-phase and verify the QA chain catches them. Mechanical drills score automatically; agent drills produce workspaces plus a scoring sheet. |

## Control-Conformance Repair

`workflow-control-scope-lint.sh` verifies a bounded repair; it does not grant
repair authority. The authority must be a phase-scoped
`TRG-<phase>-<NNN>` with exact disposition
`REPAIR / control-conformance`, an accepted contract and reproduced mismatch,
an unchanged trigger/budget, exact `project-docs/tools/` allowlist, forbidden
changes, and exact verification. Resolving a prior BLOCK additionally requires
an existing DEC and its phase-scoped `ESC-<phase>-<NNN>`.

Keep the uncommitted run ledger in the phase worktree, with no other dirt.
Create a dedicated clean temporary repair branch/worktree at the recorded base
and supply that worktree to the Developer; the Developer works only there.
Never use `assume-unchanged`, `skip-worktree`, or equivalent index flags to
conceal dirt. Run the scope linter in the repair worktree while its isolated
repair commit is checked-out `HEAD` and its index/worktree is clean. Only after
scope PASS, fast-forward the phase branch to repair `HEAD`, confirm the phase
worktree is still ledger-only dirty, then append `REPAIR`. Remove the temporary
worktree/branch after the append receipt.

`RESOLVED` permits only the ordinary same-task completion/report continuation;
it does not complete the task or authorize phase advance. Attempt counts stay
monotonic per trigger across sessions, branches, and worktrees. Material
progress resets only consecutive no-progress to `0/2`; a second consecutive
no-progress result exhausts the two-attempt budget and returns to triage.

Run from the repository root:

```sh
project-docs/tools/sync-agents.sh
project-docs/tools/sync-skills.sh
project-docs/tools/workflow-lint.sh
project-docs/tools/report-lint.sh --self-test
project-docs/tools/drills/run-drills.sh
```
