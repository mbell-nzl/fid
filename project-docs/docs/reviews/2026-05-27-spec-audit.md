# Specification Audit Report

**Date:** 2026-05-27
**Scope:** All `project-docs/specs/` files audited against `scaffold/initial-spec.md` for completeness, consistency, and implementation readiness.
**Auditor:** Claude Opus 4.7

---

## Executive Summary

The spec suite is strong: the Charter, PRD, TRD, Quality Plan, AI Workflow, and Glossary are thorough and internally coherent. Two critical gaps will directly degrade agentic implementation quality if not addressed:

1. **The Implementation Roadmap (06) is an empty template.** No phases, no tasks, no PRD/TRD traceability. Phase planning cannot begin.
2. **The Design Plan (03) lacks output design for ~60% of v0.1 CLI commands.** Based on the whet experience, commands without output examples will get crude/basic implementations.

A third structural gap — the initial spec's detailed CLI surface, schema SQL, export layout, and end-to-end flow documentation lives only in `scaffold/initial-spec.md` and isn't incorporated into the formal spec suite — creates a hidden dependency on a "draft" document that agents may not reliably load.

### Resolution Status (updated 2026-05-27)

The spec suite has been restructured. Four new spec documents were created
and existing files renumbered:

- `03_TRD_SCHEMA.md` — full SQL schema, triggers, validation (from initial spec §6)
- `04_TRD_CONFIG.md` — `.fid/` directory structure, config.toml, pricing.toml
- `06_DESIGN_PLAN_CLI.md` — CLI flag specs and output examples for all commands
- `07_DESIGN_PLAN_FLOWS.md` — end-to-end flows, whet mapping, export layout

All cross-references updated across the codebase. See per-issue status
markers below (`DONE` / `OPEN`).

---

## Issues (Prioritised)

### P0 — Blocks phase planning

#### 1. Implementation Roadmap is an empty template — DONE

**File:** `project-docs/specs/10_IMPLEMENTATION_ROADMAP.md`
**Problem:** The entire file is still the scaffolding template with `<PROJECT_NAME>`, `<Item>`, `<Description>` placeholders. No fid-specific content.
**Impact:** Phase planning requires a roadmap. The initial spec §16 has a detailed v0.1–v0.5 breakdown with entity groupings, CLI scope per version, and dependency ordering. None of this is in the formal spec.
**Action:** Fill in `10_IMPLEMENTATION_ROADMAP.md` from initial spec §16. Include:
- Phase 0 (scaffold/spec — essentially complete)
- v0.1 core issue layer items with PRD/TRD refs
- v0.2 workflow engine / governance items
- v0.3 QA / templates / TUI / MCP items
- v0.4 validation contracts / retros / export items
- v0.5 polish items
- Risk register (concurrency edge cases, schema migration complexity, TOON spec drift)

---

### P1 — Will cause crude implementations

#### 2. Design Plan lacks output design for most v0.1 commands — DONE

**File:** `project-docs/specs/05_DESIGN_PLAN.md`
**Problem:** The Design Plan has JSON output examples for only 5 commands: `issue get`, `issue list`, `ready`, `stats cost`, `maintenance check`. Approximately 20+ v0.1 commands have no output shape specified. Missing output design for:

**Read commands (no JSON shape):**
- `issue ac list` — what does the AC + verification grid look like?
- `concern list` — what fields are returned?
- `link list` — how are bidirectional links represented?
- `label list` — flat array or grouped?
- `claim list` — what fields; how is TTL/staleness shown?
- `stats tokens` — parallel to `stats cost` but not shown
- `sql` — how are arbitrary query results shaped?
- `rebuild` — what does rebuild output look like? Progress? Summary?
- `report <name>` — what does a named report return?

**Write commands (no confirmation shape):**
- `init` — what does successful init look like?
- `issue ac add` / `issue ac verify` — how are ACs confirmed?
- `dev start` / `dev gate` / `dev finish` / `dev block` — composite command output
- `concern raise` / `concern address` — confirmation format
- `link add` — does it confirm both directions?
- `run record` — telemetry recording confirmation
- `label add` — confirmation format

**Impact:** From the whet experience: AI developers implement well-designed examples faithfully but produce minimal/crude output for commands without design guidance. Every command that will ship in v0.1 needs at minimum a JSON output example and a write confirmation example.
**Action:** Add output examples for all v0.1-scoped commands. At minimum: one JSON example per read command, one write confirmation per write command. The initial spec §9 has detailed flag documentation per command that should inform these shapes.

#### 3. Design Plan lacks detailed flag specification per command — DONE

**File:** `project-docs/specs/05_DESIGN_PLAN.md`
**Problem:** The command tree (§3.1) shows command names and some flags, but doesn't fully enumerate the flags each command accepts. The initial spec §9 has exhaustive flag examples (e.g., `fid issue create --kind task --title "..." --epic <id> --priority 1 --sequence 10 --spec-md path --label bug --label auth`). Without this, developers will implement minimal flag sets.
**Action:** Either incorporate the initial spec §9 flag-level documentation into the Design Plan, or explicitly reference initial spec §9 as the flag specification source and ensure agents load it.

#### 4. Design Plan missing export directory layout and frontmatter contract — DONE

**File:** `project-docs/specs/05_DESIGN_PLAN.md`
**Problem:** The initial spec §11.2 specifies the default `whet` export layout (directory tree with `docs/current/`, `docs/epics/<epic-id>/tasks/`, `docs/epics/<epic-id>/qa/`, etc.) and §11.3 specifies the YAML frontmatter contract shape. Neither is in the Design Plan. When v0.4 export ships, the developer will have no design guidance.
**Action:** Add the export layout and frontmatter contract to the Design Plan (can be marked as v0.4 scope).

---

### P2 — Significant quality risk

#### 5. Initial spec is the de facto authoritative spec but treated as a "draft" — DONE

**Problem:** The TRD §4.1 says "The full schema is defined in `scaffold/initial-spec.md` §6. This section captures the entity inventory..." The Design Plan's command tree is a summary of initial spec §9. Multiple critical details (SQL schema, trigger list, export layout, frontmatter contract, CLI flag specifications, end-to-end flow, whet mapping) exist only in `scaffold/initial-spec.md`. But the spec suite presents itself as self-contained, and agents may not load the initial spec.
**Impact:** If an agent implements from the spec suite alone (which the doc structure implies it should), it will miss schema-level detail, trigger specifications, and CLI behavior nuances.
**Action:** Choose one of:
- **(A) Absorb critical content.** Move the SQL schema, trigger list, CLI flag specs, and export design into the TRD/Design Plan where they belong. Keep `initial-spec.md` as historical input.
- **(B) Formalize the dependency.** Add explicit references in TRD and Design Plan: "For the full SQL schema, see `scaffold/initial-spec.md` §6" etc., and ensure task specs reference it.

Recommendation: **(A)** for schema and triggers (they're architectural), **(B)** for CLI flags and export layout (they can be pulled on demand per task).

#### 6. End-to-end task flow using fid CLI commands not in specs — DONE

**Problem:** Initial spec §13 documents the complete end-to-end flow (TL planning → Workflow Runner spawn → Developer execution → Runner inspection → mid-epic QA → final QA → TL wrap-up) with specific `fid` CLI commands at each step. This is the primary user journey and it's absent from the spec suite. The AI Workflow (05) documents the *agent workflow for building fid*, not the *fid user workflow*.
**Action:** Add an "End-to-End User Flow" section to the Design Plan showing the fid CLI commands used at each step of the agentic development cycle. This is the primary design artifact for understanding how the commands compose.

#### 7. whet artifact mapping table not in specs — DONE

**Problem:** Initial spec §12 has a detailed table mapping whet artifacts (Phase → epic, Task spec → issues + ACs, Completion report → task_completions + concerns, etc.) to fid entities. This context is important for the v0.1 dogfood goal ("replaces ~70% of daily interaction with whet's task artifacts") and for the v0.4 export layout.
**Action:** Add to the Charter §5 or PRD as a migration/mapping reference.

#### 8. `.fid/` directory structure not consolidated anywhere — DONE

**Problem:** The `.fid/` directory structure is scattered across multiple documents:
- `journal/YYYY-MM-DD.jsonl` — TRD §4.3
- `fid.sqlite` — TRD §4.3
- `config.toml` — initial spec §11.1
- `pricing.toml` — initial spec §6.6 comment

No single document lists the complete `.fid/` directory structure, what each file does, and which files are committed vs gitignored.
**Action:** Add a `.fid/ Directory Structure` section to the TRD.

#### 9. config.toml and pricing.toml formats underspecified — DONE

**Problem:** `config.toml` is described in one small block in the initial spec §11.1 (export section, markdown_root, layout, banner). `pricing.toml` is mentioned in a SQL comment. Neither format is fully specified in the TRD.
**Action:** Add configuration file format specifications to the TRD. Include: what sections exist in config.toml, what the pricing.toml schema looks like (model → in_per_mtok, out_per_mtok), what ships as defaults.

---

### P3 — Inconsistency or moderate gap

#### 10. GLOSSARY "Workflow Runner" vs AI Workflow "Phase Runner" naming collision — OPEN

**Files:** `GLOSSARY.md`, `09_AI_WORKFLOW.md`
**Problem:** The Glossary defines `Workflow Runner` as fid's internal role concept and explicitly says "Avoid: Phase Runner (legacy)." But `09_AI_WORKFLOW.md` uses "Phase Runner" throughout as the name of the agent that orchestrates building fid. These are different concepts (fid's internal role vs the agentic workflow role), but the naming collision will confuse developers implementing fid's `workflow_runner` role.
**Action:** Add a clarifying note to either the Glossary or the AI Workflow doc distinguishing the two: "Phase Runner (the agent building fid) is distinct from Workflow Runner (the fid role concept)." Alternatively, rename the AI Workflow agent to avoid the collision.

#### 11. PRD missing `fid merge --journal` requirement — OPEN

**Problem:** Initial spec §16 v0.2 includes `fid merge --journal <path>` (replay foreign journal with conflict detection). The PRD has no requirement for this. The TRD doesn't mention it either.
**Action:** Add a PRD requirement under OPS (e.g., PRD-OPS-011) or explicitly defer it with rationale.

#### 12. PRD missing explicit `fid rebuild` requirement — OPEN

**Problem:** `fid rebuild` is described in the TRD (§4.3) and is mentioned extensively in the initial spec as a foundational v0.1 primitive. But the PRD has no explicit requirement for it — PRD-OPS-003 and PRD-OPS-004 describe the behavior abstractly ("rebuild the query store") without naming the command.
**Impact:** Minor — the TRD and Design Plan both show `rebuild`. But the PRD should be the requirements source.
**Action:** Verify PRD-OPS-004 is sufficient or add explicit `fid rebuild` requirement.

#### 13. PRD doesn't cover `issue update` comprehensively — OPEN

**Problem:** The Design Plan shows `fid issue update <id> --title --priority --labels --owner-role --owner-actor`. PRD-ISSUE-021 covers owner changes but no PRD requirement covers updating title, priority, or labels on existing issues. PRD-ISSUE-015 covers *creation* with these fields.
**Action:** Add a PRD requirement for issue field updates (title, priority, labels, spec_md_path, est_tokens, due_date).

#### 14. Trigger specification is summary-level in TRD — DONE

**Problem:** The initial spec §6.14 has an exhaustive trigger list (kind-invariant rules per table, polymorphic ref integrity, link scope, claim scope, comment scope, kind immutability, mutually-exclusive qa_review columns). The TRD §4.5 summarizes these as bullets. Without the detailed trigger specification, the developer implementing triggers will miss cases.
**Action:** Either expand TRD §4.5 to include the full trigger list from §6.14, or explicitly reference initial spec §6.14 for the implementation-level detail.

#### 15. Open questions not synchronized across documents — OPEN

**Problem:**
- PRD §7 has 3 open questions
- TRD §11 has 3 open questions
- Initial spec §15.1 has 1 open question
- These partially overlap (workflow_runner as deterministic harness appears in both PRD and initial spec) but the TRD questions (TOON subset, rebuild checkpoint, MCP SDK choice) don't appear elsewhere.
**Action:** Consolidate open questions into one location or ensure cross-references.

#### 16. Design Plan missing TOON and markdown output examples — DONE

**Problem:** The Design Plan §5.1 defines four output formats but only shows JSON examples. TOON is specifically called out as token-efficient for agents — showing what TOON output looks like for `issue list` would ensure consistent implementation.
**Action:** Add at least one TOON example (for `issue list` or `stats cost`) and one markdown example.

#### 17. No design for context packs (v0.2) — DONE

**Problem:** The initial spec §9.2.1 describes five built-in context packs (developer, qa_reviewer, qa_lead, tl_lite, tl) with their composition. The Design Plan mentions context packs in §4 but shows no output example. The PRD has PRD-OUT-005.
**Action:** Add a context pack output example to the Design Plan, showing the bundled shape for at least the `developer` pack. Even if v0.2, having the design ready prevents crude implementation.

---

### P4 — Minor / forward-looking

#### 18. No design for named reports output — OPEN (v0.2)

**Problem:** Initial spec §9.5 lists 8+ named reports (epic-status, open-escalations, qa-blockers, ac-matrix, etc.) shipping in v0.2. No output design exists for any of them.
**Action:** Add output examples for at least `epic-status` and `ac-matrix` to the Design Plan when v0.2 scope is planned.

#### 19. PRD traceability table has empty Task and Verification columns — OPEN (phase planning)

**File:** `project-docs/specs/01_PRD.md` §6
**Problem:** The traceability table maps PRD IDs to roadmap versions but leaves Task and Verification columns blank. This is expected pre-implementation but should be filled during phase planning.
**Action:** No action needed now — will be filled during phase planning. Flag as a checklist item for the Tech Lead.

#### 20. Quality Plan references commands not yet validated — OPEN (Phase 0)

**File:** `project-docs/specs/08_QUALITY_PLAN.md`
**Problem:** The Quality Plan references specific quality gate commands (`gofmt -l .`, `golangci-lint run ./...`, `go test -race`, etc.) but notes they need to be validated as real/working before Phase 0 exits. This is correct and the exit criteria capture it.
**Action:** No action — Phase 0 handles this.

#### 21. No `.gitignore` specification — DONE

**Problem:** The initial spec says `.fid/fid.sqlite` is gitignored. The TRD says so too. But there's no specification of what `fid init` writes to `.gitignore` (or whether it creates a `.fid/.gitignore`).
**Action:** Add to the TRD or Design Plan: `fid init` creates `.fid/.gitignore` containing `fid.sqlite`.

#### 22. MCP server design is minimal — OPEN (v0.3)

**Problem:** The initial spec §16 v0.3 has detailed MCP design (resources list, prompts list, tools list with rationale for each). The Design Plan just has a row in the surfaces table. The TRD has TRD-ARCH-013 (curated subset, not full surface).
**Action:** No action for v0.1. When v0.3 scope is planned, the MCP design from initial spec §16 should be incorporated into the Design Plan.

#### 23. AI Workflow model guidance references OpenAI models — OPEN (minor)

**File:** `project-docs/specs/09_AI_WORKFLOW.md` §6
**Problem:** The model guidance table lists `gpt-5.5` and `gpt-5.4-mini` under a "Codex" column. This is correct for OpenAI Codex agents but may confuse Claude Code users. The CLAUDE.md and `.claude/agents/` are the Claude-side materialization.
**Action:** Minor — consider adding a note that Claude Code agent model choices are configured separately in `.claude/agents/`.

---

## Cross-Document Consistency Check

| Aspect | Status | Notes |
|---|---|---|
| Entities in initial spec vs TRD | Consistent | TRD §4.1 covers all 30+ entities from initial spec §5–6 |
| ID generation scheme | Consistent | TRD §4.2 matches initial spec §6.0 |
| Persistence model | Consistent | TRD §4.3 matches initial spec §3.1 and ADR-0001 |
| Validation strategy | Consistent (summary-level) | TRD §4.5 summarizes initial spec §6.14; detail gap noted in #14 |
| Workflow states/transitions | Consistent | PRD §3.5 matches initial spec §7; Design Plan doesn't enumerate |
| Role/permission model | Consistent | PRD §3.7, TRD §4.1, ADR-0002 aligned |
| QA three-tier model | Consistent | PRD §3.9–3.10, TRD §4.1, ADR-0003 aligned |
| Output formats | Consistent | Design Plan §5.1, PRD §3.13, initial spec §10 aligned |
| Charter scope vs PRD requirements | Consistent | Charter §5.2 maps cleanly to PRD §3 sections |
| Glossary terms vs spec usage | Mostly consistent | #10 naming collision noted |
| Quality Plan vs AI Workflow | Consistent | Both reference same gate commands, templates, and artifacts |

---

## Completeness Check: Initial Spec Coverage

| Initial Spec Section | Covered In Specs? | Gap |
|---|---|---|
| §0 Changes from v0.1 | N/A (historical) | — |
| §1 Overview | Charter §1–2 | OK |
| §2 Goals / Non-Goals | Charter §5, PRD §1–2 | OK |
| §3 Core Principles | Charter §4 | OK |
| §3.1 Durability / concurrency | TRD §4.3 | OK |
| §4 Epic Scheduling | PRD §3.2 (ISSUE-008–014) | OK |
| §5 Data Model overview | TRD §4.1 | OK |
| §6 Schema SQL | `03_TRD_SCHEMA.md` (full DDL) | ~~Gap~~ **Resolved** |
| §6.14 Validation strategy | `03_TRD_SCHEMA.md` §14 (full trigger list) | ~~Gap~~ **Resolved** |
| §6.15 Validation contracts | PRD §3.10, `03_TRD_SCHEMA.md` §12 | OK |
| §6.16 Implementation stack | TRD §1 | OK |
| §7 Workflows | PRD §3.5 | OK (states/transitions not in Design Plan) |
| §8 Templates | PRD §3.13, TRD §4.1 | OK |
| §9 CLI Surface | `06_DESIGN_PLAN_CLI.md` (full flags + output) | ~~Gap~~ **Resolved** |
| §10 Output Formats | Design Plan §5, `06_DESIGN_PLAN_CLI.md` §16 (TOON + markdown) | ~~Gap~~ **Resolved** |
| §11 Export | `07_DESIGN_PLAN_FLOWS.md` §3 (layout + frontmatter) | ~~Gap~~ **Resolved** |
| §12 whet mapping | `07_DESIGN_PLAN_FLOWS.md` §2 | ~~Gap~~ **Resolved** |
| §13 End-to-end flow | `07_DESIGN_PLAN_FLOWS.md` §1 | ~~Gap~~ **Resolved** |
| §14 Linear/Jira borrows | PRD §3.2 (scattered) | OK |
| §14.1 Features not included | Charter §5.3 | OK |
| §15 Open questions | PRD §7, TRD §11 (partial) | **Gap: not synchronized** |
| §16 Implementation roadmap | `10_IMPLEMENTATION_ROADMAP.md` | ~~Gap~~ **Resolved** |
| §17 Out of scope | Charter §5.3 | OK |

---

## Recommendations

### Before phase planning begins (do now):

1. ~~**Fill 10_IMPLEMENTATION_ROADMAP.md** from initial spec §16.~~ **DONE** — filled with v0.1–v0.5 breakdown + risk register.
2. ~~**Add CLI output examples** to the Design Plan for all v0.1 commands.~~ **DONE** — created `06_DESIGN_PLAN_CLI.md` with flag specs and output examples for all v0.1 commands.
3. ~~**Add CLI flag specifications** to the Design Plan.~~ **DONE** — incorporated into `06_DESIGN_PLAN_CLI.md`.

### Before v0.1 implementation begins:

4. ~~**Consolidate `.fid/` directory structure** in the TRD.~~ **DONE** — created `04_TRD_CONFIG.md`.
5. ~~**Add config.toml and pricing.toml format specs** to the TRD.~~ **DONE** — in `04_TRD_CONFIG.md`.
6. ~~**Add end-to-end user flow** to the Design Plan.~~ **DONE** — created `07_DESIGN_PLAN_FLOWS.md` §1.
7. **Resolve the naming collision** between Phase Runner (agent workflow) and Workflow Runner (fid role). **OPEN**
8. **Add `fid issue update` requirement** to the PRD. **OPEN**
9. **Add `fid merge --journal` requirement** to the PRD or explicitly defer. **OPEN**
10. ~~**Add TOON output example** to the Design Plan.~~ **DONE** — TOON example in `06_DESIGN_PLAN_CLI.md` §2 (`issue list`). Markdown example in §16.

### Before v0.2+ phases:

11. Add named report output designs. **OPEN (v0.2)**
12. ~~Add context pack output designs.~~ **DONE** — in `07_DESIGN_PLAN_FLOWS.md` §4.
13. ~~Add export layout and frontmatter contract to Design Plan.~~ **DONE** — in `07_DESIGN_PLAN_FLOWS.md` §3.
14. Add MCP server design. **OPEN (v0.3)**
15. Fill PRD traceability table during phase planning. **OPEN**

---

*End of audit.*
