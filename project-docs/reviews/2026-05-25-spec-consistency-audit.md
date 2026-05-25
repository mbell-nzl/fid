# Spec Consistency Audit

**Date:** 2026-05-25
**Scope:** 00_PROJECT_CHARTER.md, 01_PRD.md, 02_TRD.md
**Compared against:** scaffold/initial-spec.md, whet golden docs (structural)

## Summary

The three spec documents are well-aligned with each other and faithfully
distil the initial spec into their respective concerns. The charter captures
vision and scope; the PRD captures testable product behaviour; the TRD
captures technical shape. No major structural omissions or contradictions.

Three minor issues were found and fixed in the initial pass. Two decisions
were escalated to the author and resolved; all five fixes applied in-place.

---

## Fixes Applied (trivial, no decision needed)

### 1. TRD: incorrect FK note on `decisions` stub table

**File:** 02_TRD.md line 124
**Issue:** The `decisions` notes column said "FK target for
concerns.escalated_to_esc_id and escalations.resolution_dec." The first FK
is wrong -- `concerns.escalated_to_esc_id` targets `escalations`, not
`decisions`. The initial spec (line 757-762) confirms this.
**Fix:** Removed the incorrect FK reference; the entry now reads "FK target
for escalations.resolution_dec."

### 2. PRD: escalation workflow version mismatch (PRD-WF-005)

**File:** 01_PRD.md line 116
**Issue:** PRD-WF-005 said the system "shall ship default workflows for
epic, task, subtask, and escalation kinds" in v0.1. But escalations are
stub tables in v0.1 with CLI surface in v0.2 (PRD-GOV-002, spec s16).
Shipping the escalation workflow in v0.1 when escalations can't be
created is contradictory.
**Fix:** Split into PRD-WF-005 (epic, task, subtask -- v0.1) and
PRD-WF-005a (escalation -- v0.2).

### 3. PRD: missing YAML output format requirement

**File:** 01_PRD.md
**Issue:** The spec (s10) lists four output formats: JSON, TOON, Markdown,
YAML. The TRD lists `internal/format/` as owning "JSON, TOON, Markdown,
YAML serialisation." But the PRD had no requirement for YAML output -- only
JSON (PRD-OUT-001), TOON (PRD-OUT-002), and Markdown (PRD-OUT-003).
**Fix:** Added PRD-OUT-003a for YAML format support as a Could.

---

## Decisions Resolved

### D1. PRD omits `owner_role` / `owner_actor` fields entirely

**Decision:** Add now (PRD-ISSUE-019, 020, 021).
**Rationale:** The spec documents must contain all initial-spec content.
Ownership (long-lived accountability) is semantically distinct from claims
(short-lived exclusive access) and from `--by` (the acting role on a
command). Three requirements added: default owner role at creation,
optional owner actor field, and explicit update with journaling.

### D2. Charter missing Competitive Landscape and Monetisation sections

**Decision:** Added Landscape section (s9) with positioning against
Linear/Jira, plain markdown, Beads, Beans, Claude Code tasks, and Task
Master. Licence set to MIT. No monetisation section -- not applicable.
**Rationale:** Landscape aids positioning decisions; monetisation is
not planned for a personal tool.

---

## Cross-Document Consistency: Verified Clean

| Check | Result |
|---|---|
| Charter design principles vs PRD requirements | Aligned. Each principle maps to enforceable PRD requirements. |
| Charter scope (in/out) vs PRD requirement coverage | Aligned. Every in-scope item has PRD coverage; every out-of-scope item is absent from PRD. |
| Charter interfaces table vs PRD/TRD | Aligned. All seven surfaces present in charter; PRD covers behaviour; TRD covers technical shape. |
| PRD version assignments vs TRD roadmap references | Aligned after fix #2. All PRD version tags match the initial spec's roadmap (s16). |
| PRD entity coverage vs TRD data model | Aligned. Every TRD entity has corresponding PRD requirements. Stub tables (decisions, escalations, triage_events, requirements) correctly marked as v0.1 schema / v0.2 CLI in both docs. |
| PRD transition rules vs TRD validation strategy | Aligned. PRD-WF-007 through PRD-WF-010 name the four rules; TRD-ARCH-009 describes their implementation path. |
| PRD claim semantics vs TRD claim entity | Aligned. TTL default (30 min), force-steal, atomic acquire, polymorphic ref all consistent. |
| PRD role semantics vs TRD role/policy model | Aligned. Seven roles, `--by` mandatory, human-can-do-anything, app-code-v0.1 / table-v0.2 migration path. |
| PRD link relations vs TRD link scope | Aligned. v0.1 scope (issue-issue, issue-requirement), v0.2 expansion (decision, escalation), auto-inverse, relation compatibility validation. |
| PRD output formats vs TRD format module | Aligned after fix #3. JSON default, TOON opt-in, Markdown, YAML. |
| PRD journal/rebuild vs TRD persistence rules | Aligned. One event per write, atomic projection, auto-rebuild, malformed-tail handling. |
| PRD performance targets vs TRD performance requirements | Aligned. PRD-NFR-001/002 map to TRD-PERF-001/002. TRD adds stretch targets (100k events, binary size). |
| TRD architecture vs initial spec project layout | Aligned. Package tree matches s6.16. |
| TRD dependency list vs initial spec s6.16 | Aligned. Same libraries, same rationale. |

## Structural Comparison with whet Golden Docs

| Aspect | whet | fid | Assessment |
|---|---|---|---|
| Charter heading structure | Vision, Problem, Users, Principles, Scope, Interfaces, Monetisation, Success, Constraints, Decisions Log, Competitive, Document Suite | Vision, Problem, Users, Principles, Scope, Interfaces, Success, Constraints, Decisions Log, Document Suite | Missing Monetisation and Competitive -- see D2. Otherwise aligned. |
| PRD structure | Domain-grouped (Person, Notes, Goals, etc.) with inline EARS IDs | Area-grouped tables (INIT, ISSUE, AC, etc.) with tabular EARS IDs | Different but both valid. fid's tabular approach is denser and more scannable for a tool with many cross-cutting concerns. |
| PRD requirement ID scheme | `PRD-<AREA>-NNN` inline | `PRD-<AREA>-NNN` tabular | Consistent scheme. |
| PRD error/empty states | Absent (inline in requirements) | Dedicated section (s5) | fid improves on whet here. |
| PRD traceability | Bottom section with PRD-to-task mapping | Bottom section with PRD-to-version mapping | Consistent pattern. |
| TRD structure | Stack, Architecture, DB Schema (per-table), CLI, Web, MCP, Config, Build, Seed, Performance, Traceability | Stack, Architecture, Module Boundaries, Data Model (grouped), Persistence, Validation, Interfaces, Security, Performance, Observability, Dependencies, Traceability | fid is more detailed. Module boundaries table and validation strategy are valuable additions. |
| TRD schema presentation | Full CREATE TABLE per entity | Entity inventory table (key fields + relationships) | fid defers full DDL to initial-spec, which is fine since the TRD says "full schema defined in scaffold/initial-spec.md s6." |

---

## Minor Observations (no action needed)

1. **TRD references `fid merge --journal` in the interfaces table** (line
   304) as a consumer of the journal event format. This command is in the
   initial spec's v0.2 roadmap but has no PRD requirement. Not a problem
   -- it's a technical operation, not a product behaviour -- but worth
   noting for v0.2 planning.

2. **PRD-OUT-008 (SQL escape hatch) is marked "Should, v0.1"** while the
   initial spec mentions it in both v0.1 CLI (s9.5) and as a v0.3+
   consideration. The v0.1 Should seems right given the spec's explicit
   v0.1 CLI listing.

3. **The `verifier_actor` field** on `ac_verifications` (initial spec
   line 612) and `owner_actor` on `issues` (line 515) appear in the
   initial spec's schema but are not mentioned in either the PRD or TRD
   entity summaries. These are secondary to the `_role` fields and
   correctly omitted from the summary-level documents, but implementation
   should include them.

4. **The `fid query` DSL** (initial spec s9.5, conditional on v0.3) has
   no PRD requirement. The spec marks it as conditional on gaps in named
   reports, so omission from the PRD is correct -- it can be added if
   needed.
