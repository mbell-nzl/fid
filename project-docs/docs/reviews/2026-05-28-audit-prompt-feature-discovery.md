# Audit handoff prompt — product & feature discovery (GPT-5.5)

> This is a **handoff prompt** to run a product-strategy / feature-discovery
> pass with a different model, not an audit result. Paste the block below to the
> auditing model. Companion: `2026-05-28-audit-prompt-watertightness.md`.

---

# fid spec suite — product & feature-discovery audit

You are doing a product-strategy pass on the specification suite for `fid`, a
Go-based CLI issue tracker purpose-built for AGENTIC software-development
workflows. Think like a senior PM who deeply understands both issue-tracking
products and how AI coding agents actually work. Your single job: find useful,
high-value capabilities the suite does NOT yet have — anywhere, including the
roadmap — that would make fid more valuable to its users.

This is NOT a correctness or consistency review. A separate audit covers
contradictions, traceability, and under-specification. Do NOT spend effort on
those. Spend all of it on product value and feature opportunity.

## What fid is, and why it's different (exploit these)

- **Journal-as-source-of-truth.** Every write is an append-only JSONL event; the
  SQLite store is a disposable projection rebuilt by replay. This event log is a
  latent superpower — it enables history, time-travel, audit, undo/correction,
  derived analytics, and git-branchable state. Ask: does the feature set fully
  cash in on having a complete event history?
- **Agent-native.** Output is token-efficient (TOON format), there are per-role
  context packs, an MCP server is planned, and everything is structured for LLM
  consumption. Ask: what would make an AI agent's loop dramatically smoother that
  isn't there yet?
- **Local-first, git-friendly, single binary, no daemon, no network.** Features
  must respect these principles (flag any idea that would break them as an
  explicit tradeoff).
- **Models the whole agentic dev loop**: TL planning → workflow_runner →
  developer → three-tier QA → wrap-up, with acceptance criteria, concerns,
  escalations, decisions, validation contracts, claims, and run/telemetry data.

Users: (1) AI agents executing/QA-ing/planning work, and (2) a human
orchestrator driving the loop. Stated near-term goal: replace ~70% of daily
interaction with "whet" task artifacts. A useful framing question throughout:
**what is the other 30%, and what would push fid past parity into clearly better?**

## What to read (with a product lens)

- `00_PROJECT_CHARTER.md` — vision, value prop, scope, and explicit non-goals (§5.3)
- `01_PRD.md` — current feature/requirement set and version targets
- `05_DESIGN_PLAN.md`, `06_DESIGN_PLAN_CLI.md`, `07_DESIGN_PLAN_FLOWS.md` — the
  command surface, output design, end-to-end user journeys, whet→fid mapping
- `10_IMPLEMENTATION_ROADMAP.md` — what's already planned for v0.1–v0.5
- `02_TRD.md` / `03_TRD_SCHEMA.md` — the data model (tells you what's *capturable*
  today, which often reveals features that are one CLI command away)
- `09_AI_WORKFLOW.md`, `GLOSSARY.md` — how the agentic loop is meant to run
- `scaffold/initial-spec.md` — origin doc; mine it for valuable ideas that were
  dropped (not merely deferred) during formalization

## Discovery method (use several of these lenses)

1. **Jobs-to-be-done.** For each user (developer agent, QA agent, TL/planner
   agent, human orchestrator), enumerate the jobs they do in a real loop and
   find jobs that fid doesn't yet serve, or serves clumsily.
2. **Journey walkthroughs.** Trace `07_DESIGN_PLAN_FLOWS.md` end to end and at
   each step ask "what would the actor reach for next that doesn't exist?"
   (e.g. deciding what to work on, recovering from a mistake, understanding why
   something is blocked, reviewing what changed since yesterday).
3. **Analogous-tool scan.** Compare against Linear, Jira, GitHub Issues, and
   newer agent-oriented trackers. Which genuinely useful primitives are missing
   — dependencies/blockers, saved views/filters, search, bulk operations,
   relations, cycles/sprints, notifications, SLAs? For each, judge whether it's
   valuable *for this agent-native context* (not just cargo-culted).
4. **Event-log leverage.** Given the full journal exists, what history/audit/
   time-travel/undo/diff/analytics features are cheap to offer and high-value,
   but unspecified?
5. **Agent ergonomics.** Idempotency, dry-run/preview, machine-parseable errors
   with correction hints, batch/transactional commands, "what changed" deltas,
   context-window-aware output — what's missing that reduces agent friction or
   token cost?
6. **Anti-features / sharp edges.** Anything currently specced that adds friction
   or risks misuse and should be reconsidered or simplified.

## Discipline (so this is signal, not a wishlist)

- Things already on the roadmap (ANY version) are not findings. You MAY note a
  planned feature that looks under-scoped, but the focus is net-new opportunities.
- Before proposing, check it isn't an explicit non-goal (Charter §5.3). If you
  propose something that conflicts with a stated principle or non-goal, say so
  and frame it as a deliberate tradeoff for the human to weigh.
- Every proposal must be grounded: cite what exists today (`file:section`) and
  state the concrete user need it serves. No generic "add reporting" — say who
  needs what, when, and why fid is the right place for it.
- You are allowed and encouraged to propose net-new scope. Discovery is the
  point. But rank ruthlessly by value.

## Output format

Produce a prioritized list of feature opportunities. For each:
- **Title** + one-line pitch
- **User & job-to-be-done** — who needs this and in what moment of the loop
- **Gap today** — what currently exists/doesn't, with `file:section` citation
- **Proposed capability** — concrete: the command(s)/behavior/data it implies
- **Why fid specifically** — how it exploits the journal / agent-native /
  local-first model (or note the tradeoff if it strains a principle)
- **Value & fit** — rough value (high/med/low), rough effort (S/M/L), and a
  suggested version (v0.x) given the existing roadmap

End with:
- **Top 5 highest-leverage opportunities** (the ones you'd build first), and
- A short **strategic note**: is the current feature set sufficient to hit the
  ~70%-of-whet dogfood goal, and what one capability would most move fid from
  "viable" to "clearly better than the artifacts it replaces"?
