@import AGENTS.md

## Agent Workflow

- Tech Lead (`@tech-lead`): Plans phases and wraps them up. No implementation.
- TL-Lite (`@tl-lite`): Single-use tactical questions and triage.
- Phase Runner (`@phase-runner`): Mechanical dev/QA loop orchestration.
- Developer (`@developer`): Implements one task with red/green TDD.
- QA Reviewer (`@qa-reviewer`): Bounded mid-phase reviews.
- QA Lead (`@qa-lead`): QA plans, phase-end review, human checklist.

See `.claude/agents/` for Claude agent definitions,
`project-docs/specs/09_AI_WORKFLOW.md` for the canonical workflow, and
`project-docs/docs/playbook.md` for the short operator runbook.
Phase planning includes `VAL-NNN` validation assertions in
`project-docs/docs/phases/phase-N/qa/validation-contract.md`; task specs,
completion reports, QA reviews, and final review should cite those IDs.
Canonical project language lives in `project-docs/specs/GLOSSARY.md`; root
`CONTEXT.md` is only a compatibility pointer.

## Docs Layout

- `project-docs/docs/current/` - default-loaded context:
  `phase_brief.md`, `decisions_active.md`, `progress_current.md`,
  `tl_escalations_current.md`, and `tl_handoff.md` when present.
- `project-docs/docs/phases/phase-N/` - per-phase archives:
  `tasks/`, `qa/validation-contract.md`, `qa/qa-plan.md`, `qa/review-M.md`,
  `qa/final-review.md`, `qa/human-checklist.md`, `decisions.md`,
  `progress.md`, `run-ledger.md`, `triage-log.md`, `tl-escalations.md`, and
  handoffs when needed.
- `project-docs/docs/reviews/` - retrospectives, audits, and external reviews.
- `project-docs/templates/` - reusable workflow templates.
- `project-docs/specs/GLOSSARY.md` - canonical language and flagged
  ambiguities.
- `project-docs/skills/` - canonical skill authoring source. Generated Claude
  copies live under `.claude/skills/` when skills are present.
- `docs/adr/` - pre-scaffold design decisions (predate the phase system).
- `scaffold/initial-spec.md` - original design spec; input for filling
  the seven specs progressively.
