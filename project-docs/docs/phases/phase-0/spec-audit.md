# Phase 0 Spec Audit

> Task 3 output: classify in-scope placeholder, Proposed, and open-question
> matches with source references.

## Conclusion

No blocking Phase 0 stack, command, workflow, data-model, or public-interface
gap was found. Phase 1 can be planned without unresolved Phase 0 stack choices
or command placeholders.

## Audit Scope

Audited sources:

- `AGENTS.md`
- `project-docs/specs/01_PRD.md` §7
- `project-docs/specs/02_TRD.md` §§1, 9, 11
- `project-docs/docs/current/phase_brief.md`
- `project-docs/docs/current/decisions_active.md`
- `project-docs/docs/current/progress_current.md`
- `project-docs/docs/current/tl_escalations_current.md`

## Findings

| Source | Match | Classification | Why it is or is not blocking |
|---|---|---|---|
| `AGENTS.md:162-168` | `placeholder` command rule | Non-blocking | This is a standing policy that requires placeholder commands to be treated as defects. It does not describe an unresolved project command. |
| `AGENTS.md:264-270` | Required command still a placeholder | Non-blocking | Same policy surface as above. It reinforces the hard-stop rule for future tasks, but it is not itself an open blocker. |
| `project-docs/specs/02_TRD.md:19-26` | Proposed frontmatter and TUI stack decisions | Non-blocking | These are future v0.4/v0.3 implementation choices already marked Proposed and explicitly deferred from Phase 0. They do not block Phase 1 planning. |
| `project-docs/docs/current/phase_brief.md:61-63` | Do not add future Proposed dependencies in Phase 0; stop if a required command or Phase 0 stack choice is not concrete | Non-blocking | This is a phase guardrail, not an unresolved decision. The audited Phase 0 stack and commands are concrete. |
| `project-docs/specs/01_PRD.md:359-370` | Open product questions | Non-blocking | The workflow_runner, markdown sync, and licensing questions are explicitly deferred product decisions. They are not Phase 0 blockers. |
| `project-docs/specs/02_TRD.md:786-799` | Open technical questions | Non-blocking | The TOON breadth, rebuild strategy, and future MCP adapter questions are deferred technical choices. Existing PRD/TRD constraints already keep Phase 0 listener-free and CLI/skill-first. |
| `project-docs/docs/current/phase_brief.md:79-81` | `None blocking Phase 0. PRD §7 future product questions remain deferred.` | Non-blocking | The current phase brief already records the PRD §7 questions as deferred, so they do not block this phase or Phase 1 planning. |

## Cross-Surface Command Check

- `AGENTS.md` lists concrete command gates, including format, lint, tests,
  coverage, build, and phase-level checks.
- `project-docs/specs/08_QUALITY_PLAN.md` and
  `project-docs/specs/09_AI_WORKFLOW.md` align on the same required gate
  surface.
- `.github/workflows/ci.yml` and `Makefile` expose concrete wrappers for the
  phase-2 tooling surface.
- No placeholder-only command remained unresolved in the audited Phase 0
  surfaces.

## Phase Readiness Note

Phase 1 planning can proceed. The only unresolved items found are future-phase
product or technical questions, and they are explicitly deferred by the
current spec set.

