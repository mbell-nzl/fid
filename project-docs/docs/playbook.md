# Workflow Playbook

Use this file as the short operator guide. Canonical workflow rules live in
`project-docs/specs/09_AI_WORKFLOW.md`.

## 1. Pick The Operating Mode

Every phase declares an operating mode, recorded in the phase brief. Entry
criteria and the artifact matrix live in
`project-docs/specs/09_AI_WORKFLOW.md` section 4.

| Question | If yes -> |
|---|---|
| Schema, public interface, dependency, or security/privacy change? | HEAVY |
| Release readiness, first 2-3 phases of a new project, or new domain? | HEAVY |
| Multi-task phase, AFK execution, or product-intent behavior change? | NORMAL (minimum) |
| Docs/tooling/config only, half-day scope, single subsystem, Human present? | LIGHT |

Mode-invariant in all modes: exact-command gates, evidence rules, commit
`Refs:`, glossary discipline, hard stops. If a hard-stop condition appears
mid-phase, promote the phase one mode and re-plan; never push through.

Orchestration: light mode is Human-orchestrated. In normal mode, use **Phase
Runner** when task specs are stable and execution can follow the artifacts
mechanically, or **direct human orchestration** when the phase is exploratory
or design-heavy. Heavy mode uses Phase Runner. All modes use the same
completion-report, QA, and closeout conventions for the artifacts they
produce.

## 2. Prepare The Project

Before the first implementation phase:

1. Fork or copy the starter kit.
2. Fill placeholders in `AGENTS.md`: project name, tech stack, architecture
   boundaries, quality commands, dependency policy, and conventions.
3. Fill the specs under `project-docs/specs/`, especially the charter, PRD,
   TRD, quality plan, and roadmap.
4. Fill `project-docs/specs/GLOSSARY.md` with initial project language.
5. Confirm every command placeholder (`<FORMAT_COMMAND>`, `<TEST_COMMAND>`,
   etc.) is replaced with a real command or marked `N/A` with rationale.

Do not start autonomous implementation while command placeholders remain
unresolved.

## 3. Start A Phase

1. Confirm the relevant specs are current enough for the phase.
2. Pick the operating mode per section 1 and record it in the phase brief.
   Light mode skips the remaining validation-contract and QA-plan artifact
   steps; the Human orchestrates directly.
3. Have Tech Lead draft `project-docs/docs/phases/phase-N/qa/validation-contract.md`
   from the template.
4. Have Tech Lead write task specs in
   `project-docs/docs/phases/phase-N/tasks/`.
5. Ensure every task spec cites covered `VAL-NNN` assertions or states why
   none apply.
6. Have QA Lead write `project-docs/docs/phases/phase-N/qa/qa-plan.md`.
7. Have Tech Lead refresh `project-docs/docs/current/phase_brief.md`.
8. Confirm tasks respect TDD, one-logic slicing, and smallest viable change.

In normal and heavy modes, do not start Developer work before the validation
contract, task specs, and QA plan exist unless the Human records a
phase-specific exception.

## 4. Run The Tasks

- One Developer session handles at most one task. The same task may span
  multiple continuation sessions from a durable checkpoint.
- Each task produces one completion report and one traceable commit.
- Run deterministic gates before QA review.
- Treat a Developer or QA `BLOCKED` response as a candidate to classify, not an
  authoritative lifecycle decision. Those roles never mint `ESC` IDs.
- Phase Runner routes unresolved candidates through TL-Lite using three
  dispositions:

| Disposition | Meaning | Action |
|---|---|---|
| `RESUME` | The trigger is cleared, misclassified, or already resolved with verified evidence. | Advance from the recorded checkpoint. |
| `REPAIR` | The trigger is real, local, in scope, authorized, and decision-free. | Do not advance; record recovery and dispatch a same-task Developer continuation. |
| `BLOCK` | A protected decision, unsafe condition, unrecoverable external dependency, contradictory contract, or exhausted recovery budget prevents progress. | Mint/reference an ESC, append STOP, and request the named authority action. |

`RESUME` never means repair later, and `BLOCK` is not a synonym for incomplete,
red, failed, or must-fix. Human attention is required only for `BLOCK`.

A gate/control definition change alters a command, assertion semantics,
accepted status/scope, threshold/budget/tolerance, or policy and remains
protected. A `control-conformance` repair instead restores a canonical
verification/workflow tool to an already accepted contract that uniquely
determines its behavior. TL-Lite records the contract, reproducer, tools-only
allowlist, forbidden changes, originating task/trigger/budget, exact
verification, originating report, and runner-recorded base under the ten typed
schema-v2 directive labels, then returns the exact disposition
`REPAIR / control-conformance`. Phase Runner uses this lifecycle:

1. Confirm the phase worktree has no dirt except its uncommitted run ledger and
   no hidden index flags. Record and verify full phase `HEAD` before triage as
   immutable `Control repair base`.
2. Confirm TL-Lite wrote `Originating task`, `Control trigger signature`,
   `Control attempt count`, `Control contract`, `Control mismatch`,
   `Allowed files`, `Forbidden changes`, `Exact verification`, `Originating
   report`, and `Control repair base`. Append ledger TRIAGE referencing that
   record, then copy the directive byte-for-byte into TASK-RECOVERY.
3. Create a dedicated clean temporary repair branch/worktree at that exact
   base. The Developer works only there and cannot rebaseline.
4. In that repair worktree, require observed PASS from:

```text
project-docs/tools/workflow-control-scope-lint.sh <repo> <base> <head> <allowed...>
```

5. After every scope PASS, fast-forward, verify the phase worktree remains
   ledger-only dirty, append REPAIR, and then clean up.
6. After `MATERIAL-PROGRESS` or `NO-PROGRESS`, increment attempt count, use the
   prior integrated HEAD as the next runner-owned base, and preserve every
   semantic directive value byte-for-byte.
7. Dispatch only the linked task-scope `report-repair`, restricted to exact
   `Originating report` and `Exact verification`, retaining the resolved
   attempt's `Control repair base`. Direct task completion is
   invalid; only that report-repair's later `RESOLVED` permits completion.

Never use `assume-unchanged`, `skip-worktree`, or an equivalent index flag to
conceal dirt. The claimed repair `HEAD` must be checked out in its repair
worktree, that worktree clean, and allowlisted controls regular files; later
commits, dirty state, symlinks, and gitlinks fail. Ambiguity or any definition
change is `BLOCK`; canonical location or ordinary task-file scope alone is not.

Record recovery in the append-only run ledger. Submit each single event to
`project-docs/tools/workflow-sequence-lint.sh <run-ledger> --append-stdin`; it
validates and atomically appends at EOF with an
`APPEND_CONFIRMED state=<state>` receipt. `INVALID` leaves the ledger unchanged;
if the receipt is lost, retry the byte-identical candidate to confirm it without
duplication. Do not mutate the ledger separately.

1. `TASK-RECOVERY` captures the trigger signature, attempt count, durable
   checkpoint, repair mode/authority, allowed files, and next action. New runs
   use `Recovery directive schema: v2` in the same RUN event; historical v1
   entries remain readable and retain general identity/STOP/DEC/state checks.
   Typed control authority binds byte-for-byte to the governing triage record.
2. `REPAIR` captures the same-task continuation outcome and resolution
   evidence without erasing any earlier failed attempt.
3. `RUN-RESUME` captures state reconciliation before a resumed or replacement
   Phase Runner performs another action.

`STOP` requires a preceding TRIAGE BLOCK and concrete ESC; a later resume
cannot legitimize an untriaged STOP. RUN-RESUME retains the governing
phase-scoped TRG/ESC and binds the exact next task, trigger, category, mode, and
attempt. A DEC-resolved BLOCK uses a typed directive amendment under the
existing TRG and then RUN-RESUME, never a new TRG. Verify every boundary-valid
DEC citation exists; one real DEC cannot mask a nonexistent or embedded ID.

These are temporary V1 repository-delivery entries, not forward fid product
Events, Runs, or recovery entities.

The default budget is two consecutive Phase Runner-mediated continuations for
one trigger signature with no material progress. Same-session red/green work
does not consume it. A changed tree, newly satisfied obligation, or narrower
failure is progress. Attempt counts are monotonic across sessions, branches,
and worktrees; material progress resets only consecutive no-progress to `0/2`.
A second consecutive no-progress result exhausts the budget, and a context
refresh or renamed recovery cannot reset it.
- If a task mixes logic/data with UI/wiring, re-slice before execution.

## 5. Review And Close

- Use QA Reviewer for mid-phase reviews every 3 tasks or at named checkpoints.
- QA review outcomes are `PASS`, `REPAIR`, or `BLOCK`. For `REPAIR`, name the
  owning task, repair boundary, assertion, and required evidence; Phase Runner
  dispatches that task and obtains a fresh independent re-review. The original
  finding remains in history until that re-review resolves it.
- QA Lead final review must include deterministic evidence for checklist items
  and verify `VAL-NNN` outcomes.
- Human checklist should contain only human-judgment or dogfooding items.
- Human runs the checklist and records outcomes.
- Tech Lead wraps up by:
  - Auditing every final `VAL-NNN` outcome.
  - Recording spec sync verdicts for all seven specs.
  - Auditing TRG entries in the triage log.
  - Appending workflow notes.
  - Updating current state, phase archives, and carry-forward items.
  - Compressing active context to size budgets.

## 6. Key File Locations

| Artifact | Path |
|---|---|
| Specs | `project-docs/specs/` |
| Glossary | `project-docs/specs/GLOSSARY.md` |
| Templates | `project-docs/templates/` |
| Current context | `project-docs/docs/current/` |
| Phase archives | `project-docs/docs/phases/phase-N/` |
| Validation contract | `project-docs/docs/phases/phase-N/qa/validation-contract.md` |
| QA plan | `project-docs/docs/phases/phase-N/qa/qa-plan.md` |
| Task specs | `project-docs/docs/phases/phase-N/tasks/` |
| Skills (canonical) | `project-docs/skills/` |
| Workflow-lint | `project-docs/tools/workflow-lint.sh` |

## 7. Rules That Do Not Bend

- Files are the communication layer; do not rely on chat summaries.
- Specs and templates are the source of truth, not agent prompts or skills.
- Required artifacts are mandatory unless the human records a specific
  exception.
- Runtime dependency additions require dependency-budget documentation.
- A required validation assertion that cannot be satisfied must be escalated,
  not weakened or silently adapted.
- Deterministic claims require reproducible evidence.
- A required command that did not terminate is never `PASS`; it is `BLOCKED`.
- Failed evidence remains failed after repair; append the byte-identical rerun
  and its observed termination rather than replacing history.
- Incomplete work, a failed gate, report repair, local environment recovery,
  or bounded QA rework is `REPAIR` unless a separate protected condition
  exists.
- A tools-only control implementation defect is also `REPAIR` when an accepted
  contract uniquely determines the fix and the scope verifier proves the
  isolated commit stayed inside its directive.
- Human attention and STOP are reserved for a protected blocker category or
  `recovery-budget-exhausted`.
- Triage and escalation identifiers are phase-scoped (`TRG-<phase>-<NNN>`,
  `ESC-<phase>-<NNN>`), allocated from the phase triage log, and written to disk
  before the verdict is reported. Only TL-Lite or the full Tech Lead allocates
  ESC IDs; Developer and QA candidates never do. Never use a placeholder.
- Run-ledger and triage entries are append-only in wall-clock order. Correct an
  earlier entry by appending a labelled correction, never by rewriting history.
- Keep held-out scenario material outside normal implementation context.

## 8. Audit A Workspace

ADR-0011 makes JSON the only machine representation. The target audit surface
after convergence slice A is:

```text
fid maintenance check
fid maintenance check --warn-as-error --json
```

Human default output is concise terminal text; `--json` emits the canonical
machine envelope. There is no `--format`, TOON, YAML, or generic Markdown
output. The current pre-convergence binary still contains displaced formatter
behavior, but no adopted workspace depends on it and new integrations must not
use it. The audit refreshes a missing/stale projection without appending to the
journal. Warnings exit zero unless `--warn-as-error` promotes the outcome;
structural errors exit non-zero.

## 9. Compact the Derived Database

Use VACUUM only when reclaiming local SQLite space is useful:

```text
fid maintenance vacuum --by tl
```

The current pre-convergence command is restricted to V1 `tl` and `human`
callers. That role check is replacement-bound, not a forward authority
contract. The operation waits for the shared write lock, refreshes the derived
projection if needed, then compacts only
`.fid/fid.sqlite`. It does not append, compact, or delete journal history, and
it creates no activity event. The stderr confirmation reports the database byte
sizes before and after the operation.
