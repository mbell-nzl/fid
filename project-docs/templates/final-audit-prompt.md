# Final Consistency Audit Prompt

Copy the prompt below and provide it to a fresh model with repository access.

```text
You are performing a final pre-implementation consistency audit of this
project's specification and agent workflow.

Read these documents in full before analysis:

1. AGENTS.md
2. CLAUDE.md
3. project-docs/specs/00_PROJECT_CHARTER.md
4. project-docs/specs/01_PRD.md
5. project-docs/specs/02_TRD.md
6. project-docs/specs/03_TRD_SCHEMA.md
7. project-docs/specs/04_TRD_CONFIG.md
8. project-docs/specs/05_DESIGN_PLAN.md
9. project-docs/specs/06_DESIGN_PLAN_CLI.md
10. project-docs/specs/07_DESIGN_PLAN_FLOWS.md
11. project-docs/specs/08_QUALITY_PLAN.md
12. project-docs/specs/09_AI_WORKFLOW.md
13. project-docs/specs/10_IMPLEMENTATION_ROADMAP.md
14. .claude/agents/*.md
11. .codex/config.toml
12. .codex/agents/*.toml

Audit dimensions:

1. Cross-document consistency:
   - Project name, scope, stack choices, commands, architecture, requirement IDs,
     quality gates, docs paths, and agent role names.
2. Internal consistency:
   - Section references, examples, tables, command syntax, and terminology.
3. Agent file consistency:
   - Agent responsibilities, tool restrictions, output artifacts, stop
     conditions, context loading, and model/capability mapping.
4. Implementability:
   - Ambiguity, missing definitions, contradictory requirements, untestable
     acceptance criteria, or unresolved placeholders.
5. Staleness and drift:
   - Any agent prompt or docs path that no longer matches the canonical workflow.
6. Product and stack agnosticism:
   - Starter-kit material must not leak a specific product, programming language,
     framework, database, command runner, or UI surface unless it is explicitly a
     placeholder/example to replace.
7. Quality-system completeness:
   - Verify deterministic/human QA split, TL-Lite triage, spec sync,
     workflow notes, audit prompts, held-out scenarios, threat model, and
     pre-release checklist are represented consistently.

Output format:

# Final Consistency Audit - <PROJECT_NAME>

## Summary
- X findings: Y Critical, Z Major, W Minor, V Informational
- Overall assessment: Ready / Needs fixes before implementation

## Findings

### AUDIT-001 [Severity]
**Location:** <document and section>
**Finding:** <what is wrong>
**Recommendation:** <how to fix>

## Consistency Matrix

Include a compact table for key facts and whether each document agrees.
```
