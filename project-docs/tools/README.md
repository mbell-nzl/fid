# Repo Tools

Stack-neutral repository commands that do not require project runtime
dependencies.

| Command | Purpose |
|---|---|
| `sync-agents.sh` | Generate provider role prompts from `project-docs/agents/`. Use `--check` to compare without writing. |
| `sync-skills.sh` | Generate provider skill copies from `project-docs/skills/`. Use `--check` to compare without writing. |
| `workflow-lint.sh` | Verify agent workflow artifacts, templates, config, and generated prompt/skill copies have not drifted. |
| `report-lint.sh <report> <task-spec>` | Deterministically verify a completion report: typed workflow outcome, required headings, command coverage, observed termination, commit identity, and that the report is committed and clean. Emits `OUTCOME COMPLETE`, `OUTCOME RECOVERY_REQUIRED`, or `OUTCOME BLOCKED`. `--self-test` proves pass and fail cases. |
| `workflow-route.sh <source> <outcome> [attempts]` | Map typed Developer, report, QA, triage, command/session-liveness, and recovery-budget outcomes to the Phase Runner's next action. |
| `workflow-control-scope-lint.sh <repo> <base> <head> <allowed...>` | Verify a runner-bound control-conformance base-to-HEAD repair is clean, regular-file-only, and changes no path outside its explicit tools allowlist. |
| `workflow-sequence-lint.sh <run-ledger> [--append-stdin]` | Validate schema, triage/STOP authority, typed recovery budgets and ordering, or idempotently validate/append one stdin event at EOF; a byte-identical retry confirms a lost receipt without duplication. Reports state but never chooses the next action. |
| `drills/run-drills.sh` | Evidence-chain drills: seed known defects into a fixture mini-phase and verify the QA chain catches them. Mechanical drills score automatically; agent drills produce workspaces plus a scoring sheet. |

## Control-Conformance Repair

`workflow-control-scope-lint.sh` verifies a bounded repair; it does not grant
repair authority. Under `Recovery directive schema: v2`, authority is the
phase triage-log record with exact `REPAIR / control-conformance` and these
one-line fields: `Originating task`, `Control trigger signature`,
`Control attempt count`, `Control contract`, `Control mismatch`,
`Allowed files`, `Forbidden changes`, `Exact verification`,
`Originating report`, and `Control repair base`.
Each label occurs exactly once, values bind byte-for-byte, and label order or
contiguity does not matter. `Allowed files` is an exact comma-separated ordered
list with no spaces around commas. A standalone historical schema-v1 ledger
remains valid; every newly appended RUN contains exactly one sequence v2 and
schema v2 marker and cannot downgrade. Every mandatory schema-v2 event field
occurs exactly once, using explicit `N/A` where the ledger template permits it.

Keep the uncommitted run ledger in the phase worktree, with no other dirt.
Before triage and dispatch, the Phase Runner records and verifies full phase
HEAD as immutable `Control repair base`. Create a dedicated clean temporary
repair branch/worktree there and supply it to the Developer; the Developer
works only there and cannot rebaseline.
Never use `assume-unchanged`, `skip-worktree`, or equivalent index flags to
conceal dirt. Materialize the scope-linter blob from the trusted directive base
into fresh scratch storage, never from repair HEAD, and execute that blob in the
repair worktree using the directive base, Developer-produced checked-out HEAD,
and individual allowlist paths in CSV order while its index/worktree is clean.
A Developer-selected base is invalid. The PASS receipt must bind the exact
comma-separated ordered `Allowed files` value byte-for-byte as
`PASS control repair scope: <base> -> <head>; allowed=<Allowed files> (terminated exit_code: 0)`. Normal
control-conformance cannot repair `workflow-control-scope-lint.sh` or
`workflow-sequence-lint.sh`; either verifier requires separate protected
decision authority.

After every control scope PASS, fast-forward, append REPAIR, and clean up. For
nonresolved control repair, increment attempt count, use the prior integrated
HEAD as the next runner-owned base, and preserve every semantic directive
value. Control `RESOLVED` permits only the exact linked report-repair path and
verification, retaining the resolved attempt's base, never direct task
completion. The report-repair's later
`RESOLVED` permits completion. Attempt counts stay monotonic; a second
consecutive no-progress result exhausts the budget.

A DEC-resolved BLOCK retains its governing phase-scoped TRG/ESC. Append one
unique separate H2 typed directive amendment naming that TRG and parse it
independently from the original TRG's retained BLOCK verdict. RUN-RESUME then
binds the exact next recovery; do not create a new TRG. STOP resume authority is
exactly `verified-commit:<40hex>`, `restored-capability:<stable-id>`, or
`DEC-<NNN>`, and protected decision changes/resolutions require DEC. Every
boundary-valid DEC citation in an authority field must exist in phase
decisions. A normal resume uses `N/A` governing IDs and `Next recovery`.

Every schema-v2 external TRG and ledger TRIAGE carries exactly one
`Originating task`, `Trigger signature`, `Attempt count`, `Repair category`,
and `Repair authority`; direct REPAIR binds all five through TASK-RECOVERY.

## Child-Session Liveness

A parent polling wait timeout and a child still reported live both route to
`ACTION WAIT`; neither authorizes interruption. A returned child routes to
`ACTION EVALUATE`. An externally interrupted child routes to `ACTION REPAIR`
on the first consecutive occurrence and `ACTION TRIAGE` on the second. Its
ledger REPAIR uses `Outcome: INTERRUPTED`, records interruption cause, elapsed
session duration, and command in flight, and preserves the preceding
no-progress count. Only declared per-command budgets apply to commands; none
bounds the child session lifetime.

Run from the repository root:

```sh
project-docs/tools/sync-agents.sh
project-docs/tools/sync-skills.sh
project-docs/tools/workflow-lint.sh
project-docs/tools/report-lint.sh --self-test
project-docs/tools/drills/run-drills.sh
```
