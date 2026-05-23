# Ship the three-tier QA model in full

## Decision

fid ships the full three-tier quality model: validation contracts, QA plans, and QA reviews. Plans + reviews land in v0.3; contracts + assertions + outcomes + the AC↔assertion linkage + coverage reports land in v0.4. All three layers are schema-optional — a project can use any subset — but they are not optional *from fid*: they ship together as the supported quality model.

## Context

The primary user's agentic coding workflow encodes assertion-level verification (validation contracts → QA plans → QA reviews with per-assertion outcomes). Without contract-layer support in fid, that workflow has to live outside the tracker, defeating the point of using fid for project state.

This sits in tension with the "lightweight, doesn't dictate workflow" framing that motivated ADR-0002 (configurable roles). The resolution: **the lightweight framing applies to roles, permissions, and workflow state machines — surfaces where projects legitimately disagree. It does not apply to the quality model, where the primary user has a specific multi-tier shape they need supported.**

Reference patterns: Cucumber/Gherkin (feature → scenario → step), regulated-industry verification (requirements → procedures → results), Factory Missions (contract → workers → validators). The pattern is well-established; what's distinctive about fid is making it native to an agent-driven tracker.

## Considered options

- **(Q-1) Ship as specced.** Chosen, with selective sharpening so the layering and optionality lead the section instead of trailing it.
- **(Q-2) Defer contracts to v0.5+ or later.** Rejected — primary user needs contracts to start using fid for real work.
- **(Q-3) Drop contracts entirely.** Rejected — same reason; also removes the traceability features the user's workflow depends on.
- **(Q-4) Compress §6.15 exposition aggressively.** Rejected as the primary action — the substance is dense but each subsection earns its keep at implementation time. Light sharpening is folded into (Q-1).

## Consequences

- **Schema cost is real.** Six tables (`validation_contracts`, `validation_assertions`, `qa_plans`, `qa_plan_steps`, `qa_plan_step_assertions`, `val_assertion_outcomes`) plus two column additions. Maintained.
- **Implementation order:** v0.3 ships plans + steps + AC-citing-step backbone; v0.4 ships contracts + assertions + outcomes + linkage + the three coverage/trace named reports.
- **Optionality is per-project, not per-tool.** §6.15 is structured so a reader who only wants ACs + ad-hoc reviews stops after the layering overview; readers who want the full pattern read on. The schema does not require contracts; the docs do not push them as default.
- **Coverage rules only activate when both a contract and a plan are present and active for the same epic.** A project running ACs + reviews-only sees no extra friction.
- **Future readers who land on §6.15 and wonder "why all this for a personal tracker" find their answer in this ADR.**
