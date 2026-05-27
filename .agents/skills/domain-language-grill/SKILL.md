---
name: domain-language-grill
description: Grilling session that challenges a plan or design against the project glossary, sharpens terminology, and updates GLOSSARY.md inline as decisions crystallise.
---

<!-- Generated from project-docs/skills/domain-language-grill/SKILL.md by project-docs/tools/sync-skills.sh. -->
<!-- Do not edit directly; update the canonical skill and rerun the sync command. -->

# Domain Language Grill

Provenance: External-Adapted Skill. Based on Matt Pocock's grill-with-docs
skill, adapted for this starter kit's glossary and decision system.

Use this skill to stress-test a plan, design, or feature proposal against
the project's canonical language before implementation begins.

## Procedure

1. Read `project-docs/specs/GLOSSARY.md` and any relevant spec sections.
2. Interview the user about every aspect of the plan, one question at a
   time. For each question, provide a recommended answer. Wait for feedback
   before continuing.
3. If a question can be answered by exploring the codebase or specs, explore
   instead of asking.

## During The Session

### Challenge against the glossary

When the user uses a term that conflicts with the existing language in
`GLOSSARY.md`, call it out immediately. Propose the canonical term or ask
which meaning is intended.

### Sharpen fuzzy language

When the user uses vague or overloaded terms, propose a precise canonical
term. Distinguish between concepts that look similar but have different
roles.

### Discuss concrete scenarios

When domain relationships are being discussed, stress-test them with
specific scenarios that probe edge cases and force precision about
boundaries between concepts.

### Cross-reference with code

When the user states how something works, check whether the code agrees.
Surface contradictions between stated intent and existing implementation.

### Update GLOSSARY.md inline

When a term is resolved, update `project-docs/specs/GLOSSARY.md`
immediately. Do not batch updates. Use the glossary entry format: preferred
term, definition, use-when, avoid, related terms, examples.

The glossary must stay free of implementation details. It is a language
reference, not a spec or architecture document.

### Record decisions through the existing system

When a significant decision crystallises, record it through the project's
`DEC-NNN` decision system in the relevant phase decisions file. Do not
create a parallel ADR directory.

Only record decisions that are hard to reverse, surprising without context,
and the result of genuine trade-offs with alternatives considered.

## References

- `project-docs/specs/GLOSSARY.md`
- `project-docs/templates/glossary-template.md`
- `project-docs/templates/phase-decisions-template.md`
- `project-docs/specs/09_AI_WORKFLOW.md`
