---
name: validation-contract
description: Use when drafting, refining, or auditing a phase validation contract with VAL-NNN assertions before implementation or during QA/wrap-up.
---

# Validation Contract

Provenance: Workflow-Extracted Skill.

Use this skill to turn phase scope into auditable `VAL-NNN` assertions without
making the skill the source of workflow authority.

## Procedure

1. Load the phase scope plus the relevant requirements from
   `project-docs/specs/`, the glossary, and
   `project-docs/templates/validation-contract-template.md`.
2. Express each in-scope requirement, decision, or invariant as one observable
   validation assertion. Keep each `VAL-NNN` narrow enough that evidence can
   prove or disprove it.
3. For each assertion, fill the source requirement, evidence method, validator
   role, task/window ownership, pass signal, fail or blocked signal, and any
   explicit non-goals.
4. Check the planned task set: every task should cite covered `VAL-NNN` IDs or
   explain why no validation assertion applies.
5. Mark unresolved evidence, ownership, or terminology gaps where the canonical
   workflow expects them: task spec, QA plan, phase decisions, or escalation
   artifacts.
6. During QA or wrap-up, compare final evidence against the contract and keep
   outcome status tied to the original assertion language.

## References

- `project-docs/templates/validation-contract-template.md`
- `project-docs/templates/task-spec-template.md`
- `project-docs/specs/04_QUALITY_PLAN.md`
- `project-docs/specs/05_AI_WORKFLOW.md`
- `project-docs/specs/GLOSSARY.md`
