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

- One Developer session handles one task.
- Each task produces one completion report and one traceable commit.
- Run deterministic gates before QA review.
- Phase Runner routes soft blockers through TL-Lite triage.
- Hard failures stop immediately.
- If a task mixes logic/data with UI/wiring, re-slice before execution.

## 5. Review And Close

- Use QA Reviewer for mid-phase reviews every 3 tasks or at named checkpoints.
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
- Keep held-out scenario material outside normal implementation context.

## 8. Audit A Workspace

Run the read-only maintenance audit before handoff or after rebuilding a
projection:

```text
fid maintenance check
fid maintenance check --format markdown
fid maintenance check --warn-as-error --format json
```

The command refreshes a missing or stale SQLite projection before reading but
does not append to the journal. It writes the complete audit result to stdout
in `json` (default), `toon`, `markdown`, or `yaml`. No findings return
`status: pass`; warnings alone return `status: warn` and exit zero unless
`--warn-as-error` promotes the status and exit outcome to `fail`; any structural
error returns `fail` and a non-zero exit.

## 9. Compact the Derived Database

Use VACUUM only when reclaiming local SQLite space is useful:

```text
fid maintenance vacuum --by tl
```

The command is available to `tl` and `human` only. It waits for the shared
write lock, refreshes the derived projection if needed, then compacts only
`.fid/fid.sqlite`. It does not append, compact, or delete journal history, and
it creates no activity event. The stderr confirmation reports the database byte
sizes before and after the operation.
