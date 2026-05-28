# Audit handoff prompt — watertightness + feature gaps (GPT-5.5)

> This is a **handoff prompt** to run an independent second-pass audit with a
> different model, not an audit result. Paste the block below to the auditing
> model. Companion: `2026-05-28-audit-prompt-feature-discovery.md`. Feeds into
> `2026-05-27-spec-audit.md`.

---

# fid spec suite — independent second-pass audit (watertightness + feature gaps)

You are auditing the specification suite for `fid`, a Go-based CLI issue tracker
purpose-built for agentic software-development workflows. Core architecture: an
append-only JSONL **journal** is the durable source of truth; a **SQLite query
cache** is a fully rebuildable projection of that journal (`fid rebuild` replays
it). Single static binary, local-first, no daemon, no network. It models issues
(epics/tasks with a kind hierarchy), acceptance criteria, concerns, links,
labels, claims, comments, governance entities (decisions, escalations), a
three-tier QA model, and run/telemetry data. Seven roles: developer,
qa_reviewer, qa_lead, tl_lite, tl, workflow_runner, human. Output formats: json,
toon, markdown, yaml. Roadmap spans v0.1–v0.5.

## Your objective

This suite has already had one audit round and a follow-up resolution pass. Your
job is a FRESH, INDEPENDENT second pass with two lenses:

1. **Watertightness** — find contradictions, ambiguities, traceability gaps, and
   under-specified behavior that would cause a coding agent to implement the
   wrong thing or diverge across documents.
2. **Missing features** — identify useful, in-scope capabilities the suite does
   NOT specify but a tool of this kind plausibly should, given its users (AI
   agents + a human orchestrator running an agentic dev loop) and its stated
   dogfood goal (replace ~70% of daily interaction with whet's task artifacts).

Bring an adversarial, independent eye. Do not assume prior reviewers were
correct — verify against the files. Where you agree the suite is sound, say so
briefly; spend your effort on genuinely new findings.

## What to read

All files under `project-docs/specs/`:
- `00_PROJECT_CHARTER.md` — vision, scope, non-goals (§5.3 = explicit out-of-scope)
- `01_PRD.md` — product requirements (IDs like PRD-ISSUE-NNN), traceability §6, open questions §7
- `02_TRD.md` — architecture; `03_TRD_SCHEMA.md` — full SQL DDL + triggers; `04_TRD_CONFIG.md` — `.fid/` layout + config formats
- `05_DESIGN_PLAN.md` — UX overview; `06_DESIGN_PLAN_CLI.md` — per-command flag specs + output examples; `07_DESIGN_PLAN_FLOWS.md` — end-to-end flows, whet mapping, export
- `08_QUALITY_PLAN.md`, `09_AI_WORKFLOW.md`, `10_IMPLEMENTATION_ROADMAP.md`, `GLOSSARY.md`

Origin/historical input (NOT the source of truth, but useful to compare against
for dropped detail): `scaffold/initial-spec.md`.

Prior audit + resolution log: `project-docs/docs/reviews/2026-05-27-spec-audit.md`.
Read its per-issue `DONE`/`OPEN` markers and `Resolution` lines first.

## Ground rules

- **Do not re-raise anything marked DONE** in the prior audit. If you think a
  DONE resolution is actually wrong or incomplete, that's a valid NEW finding —
  but say explicitly "prior #N resolved this incorrectly because…" and cite the
  file/section, rather than restating the original issue.
- **Do not re-litigate the intentionally-deferred items**: #18 (named reports,
  v0.2), #19 (PRD traceability table, phase planning), #20 (quality-gate command
  validation, Phase 0), #22 (MCP design, v0.3). These are deferred by design.
- **Verify every claim against the actual file.** Cite `file:section` (and line
  where helpful). A finding without a citation is not actionable.
- **Respect declared scope.** Before proposing a "missing feature," check it
  isn't already an explicit non-goal (Charter §5.3) or deferred to a later
  version. Label each feature idea as: in-scope-gap / later-version / arguably
  out-of-scope. Do not invent scope creep.

## Specific things to probe (not exhaustive — go beyond this list)

Watertightness:
- **Cross-doc contradictions**: command names, flag names, field names, role
  names, enum values, output shapes, and workflow states must agree across PRD,
  TRD, schema, Design Plan/CLI, and Glossary.
- **Traceability both directions**: every CLI command in `06_DESIGN_PLAN_CLI.md`
  should map to a PRD requirement, have schema support in `03_TRD_SCHEMA.md`, and
  appear in the roadmap. Conversely, flag any PRD requirement with no command/
  schema home, and any schema column/table no CLI can read or write.
- **Inverse/lifecycle completeness**: for every additive operation, is the
  reverse specified where it should be? (add/remove label, link/unlink, claim/
  unclaim/steal, verify/unverify AC, block/unblock, open/close, etc.)
- **Failure modes**: concurrency (two agents, same checkout), journal corruption
  / malformed lines, schema migration + event-schema upgrade chain, merge/
  conflict resolution, partial composite writes, stale claims. Are behaviors
  fully specified or only gestured at?
- **Trigger / invariant coverage** in `03_TRD_SCHEMA.md`: do the documented
  triggers actually enforce every stated invariant? Any invariant with no
  enforcement, or trigger with no stated rationale?
- **Ambiguities** that would yield divergent implementations or crude output
  (the kind of under-specification that produces minimal, low-quality command
  output).

Missing features (think about the domain, not just the docs):
- What does a coding agent or human orchestrator predictably need mid-loop that
  isn't expressible? (e.g. querying/filtering, bulk ops, dependency/blocking
  semantics, search, history/audit of a single entity, undo/correction of a bad
  journal event, dry-run/preview, export/import round-trips, watch/notify.)
- Are there obvious gaps in observability, telemetry, or cost/token accounting
  given fid's stated analytics ambitions?
- Anything in `scaffold/initial-spec.md` that was genuinely useful but got
  dropped (not just deferred) during the formalization?

## Output format

Mirror the existing audit so findings drop into the same workflow:
- Group findings by severity: **P0** (blocks implementation) → **P4** (minor/
  forward-looking).
- For each finding: a short title, the **file:section** citation, **Problem**
  (what's wrong/missing and why it matters), and a **Proposed fix** (concrete:
  the requirement ID to add, the text to change, the command to specify).
- Separate "watertightness findings" from "feature-gap findings."
- End with a one-paragraph overall assessment: is the suite implementation-ready,
  and what are the top 3 things to fix before phase planning begins.
