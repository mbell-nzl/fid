---
name: qa-evidence
description: Use when planning, collecting, or reviewing deterministic and human evidence for VAL-NNN assertions in QA plans, reviews, completion reports, or final reviews.
---

# QA Evidence

Provenance: Workflow-Extracted Skill.

Use this skill when evidence quality matters more than prose confidence.

## Procedure

1. Start from the phase validation contract and the relevant template:
   `qa-plan-template.md`, `mid-phase-review-template.md`,
   `final-review-template.md`, `completion-report-template.md`, or
   `human-checklist-template.md`.
2. Map each check to a `VAL-NNN` ID or explicitly record `no VAL` when the
   check supports quality but not a validation assertion.
3. Prefer reproducible deterministic evidence: exact command output, test logs,
   file-state proof, API traces, screenshots, rendered artifacts, or commit
   references.
4. Treat required command coverage as evidence, not ceremony. Every command
   required by the task spec or QA plan must appear as `PASS`, or `BLOCKED`
   with a reason in a blocker report. Missing, silently skipped, or narrower
   substituted commands are failed gates.
5. For conditional or binary assertions, collect evidence for both the
   qualifying and non-qualifying states. Positive-only evidence cannot prove
   the assertion.
6. Keep human evidence separate for product intent, perception, dogfooding, or
   other judgment calls. Cite the relevant `VAL-NNN` when a human check proves
   a validation assertion.
7. For commit evidence, copy hashes and subjects from git output; do not
   reconstruct them from memory. Phase orchestration verifies hashes with
   `git rev-parse --verify <hash>^{commit}`.
8. For every claimed pass, record enough detail that another agent or human can
   rerun, inspect, or trace the proof.
9. If evidence is partial, stale, indirect, or contradicted by another artifact,
   record that status plainly in the review artifact rather than upgrading it
   to a pass.

## References

- `project-docs/templates/qa-plan-template.md`
- `project-docs/templates/mid-phase-review-template.md`
- `project-docs/templates/final-review-template.md`
- `project-docs/templates/human-checklist-template.md`
- `project-docs/templates/completion-report-template.md`
- `project-docs/specs/08_QUALITY_PLAN.md`
