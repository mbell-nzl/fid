# TL Escalations - Current

> **Purpose.** Open items that exceed local-decision authority and need full
> Tech Lead or human input.
>
> **Scope.** Open items only. Resolved phase items move to the phase archive.
>
> **Size target:** under 80 lines.

## Open Escalations

### ESC-009: Phase 17 QA routes lack model-qualification evidence

- **Raised:** 2026-07-11
- **By:** QA Lead / Phase 17 QA-BLK-02; accepted by Tech Lead
- **Trigger:** DEC-053 and `09_AI_WORKFLOW.md` require seeded-defect evidence
  before trusting a new model in a QA role.
- **Context:** The new QA Lead did preserve the real Phase 17 hard stop in
  `qa-plan.md`, but no Human-scored drill artifact qualifies either new QA
  route, and no evidence proves the QA Reviewer escalation obligation.
- **Question / proposed direction:** Preferred: run the full evidence-chain
  drill set against fresh `gpt-5.6-terra` QA Reviewer and `gpt-5.6-sol` QA Lead
  sessions, score every DEC-053 obligation, and commit a model-audit artifact
  under `phase-17/qa/`. Do not infer uncovered obligations; add a scored seed
  for cross-cutting escalation if needed. DEC-055 rejects any routing fallback;
  only passing model-audit evidence can resolve this escalation.
- **Links:** DEC-053; QA-BLK-02; `project-docs/tools/drills/README.md`;
  `project-docs/tools/drills/run-drills.sh`.
- **Status:** Open - blocks all Developer execution; qualification is the only
  accepted resolution.

## Recent Archive Notes

- Phase 16 raised no `ESC-NNN` items. Local unblock and workflow decisions
  DEC-043..048 plus TRG-046 are archived in
  `../phases/phase-16/tl-escalations.md` and `../phases/phase-16/decisions.md`.
- Phase 16 VAL-013 is not an escalation; the completed human checklist records
  PASS for HUMAN-16a..d.

## Standing Escalation Triggers

TL-Lite escalates instead of deciding locally when an issue:

1. Changes a public interface, CLI contract, output shape, journal event schema,
   data model, persistence contract, dependency, tool, quality gate, or workflow.
2. Contradicts `decisions_active.md`, `decisions_principles.md`, `AGENTS.md`, or
   a recorded phase decision.
3. Affects more than one task/subsystem, exposes product-intent ambiguity, or
   changes testing strategy / coverage policy.
4. Resumes context-exhausted or partial work spanning schema/data model,
   dependencies, public interfaces, canonical docs, or multiple Developer
   sessions.

QA findings escalate here if they have architecture implications, ambiguous
acceptance criteria, unclear security/privacy severity, contradict a standing
decision, span tasks/subsystems beyond the review window, plausibly indicate a
material spec deviation, or expose a quality-gate anomaly not anticipated by the
QA plan.
