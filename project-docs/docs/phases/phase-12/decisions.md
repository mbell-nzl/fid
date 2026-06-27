# Phase 12 Decisions

> Purpose: phase-local decision history, spec sync verdicts, workflow notes, and
> wrap-up decisions. DEC numbering is global and never resets.

## Decisions

### DEC-029: `fid requirement sync` hashes the parsed summary cell only

- **Date:** 2026-06-27
- **Context:** Phase 12 Task 1 needs a public, reusable parse contract for
  requirement pointers before `fid link add ... implements <REQ>` can depend on
  the `requirements` table. The validation contract's PI-1 defines the
  canonical source shape at a high level, but Task 1 must record the exact
  hashed content and the section-capture behavior used by the minimal v0.1
  sync command.
- **Decision:** `fid requirement sync` parses markdown table rows whose first
  cell matches `(PRD|TRD)-[A-Z0-9]+-[0-9]+[a-z]?` and whose second cell is the
  summary. `type` is derived from the ID prefix. `source_section` is the text
  of the nearest preceding markdown heading, or `NULL` when no heading has been
  seen yet. `last_synced_hash` is `sha256:<hex>` over the parsed summary cell
  only after trimming table-cell whitespace; it does not hash the requirement
  ID, notes columns, or surrounding prose. Duplicate canonical IDs in the same
  source reject before append, even if a later `--type` or `--id` filter would
  have excluded one duplicate from the persisted result set.
- **Consequences:** Re-syncing unchanged content appends no event because the
  persisted comparison key is row-local and deterministic. Task 2 can validate
  issue -> requirement endpoints against a stable pointer contract without
  copying full PRD/TRD prose into the read model. No schema migration is needed
  because `migrations/0002_v01_tables.up.sql` already exposes the required
  `requirements` columns. No section-13 mismatch was found for
  `requirement.sync`; the in-code matrix row matches `03_TRD_SCHEMA.md`
  section 13 (`tl` + `human` only).
- **Refs:** PI-1; phase-12/task-1; VAL-001; VAL-002; VAL-003;
  `03_TRD_SCHEMA.md` sections 5.2 and 13; `06_DESIGN_PLAN_CLI.md`
  `fid requirement sync`.

### DEC-030: Phase 12 link enforcement uses a closed v0.1 relation registry with app-layer prechecks

- **Date:** 2026-06-27
- **Context:** Phase 12 Task 2 needs to implement `fid link add` without a new
  schema migration while satisfying PRD-LINK-008 and the existing `links_*`
  triggers. The validation contract's PI-2 says v0.1 should reject
  dedicated-schema duplication structurally via a closed relation set,
  issue/requirement-only endpoint scope, and self-link rejection before append.
- **Decision:** `internal/links` owns a closed registry containing exactly the
  documented seven relations plus inverses:
  `blocks/blocked_by`, `relates_to/relates_to`,
  `duplicates/duplicated_by`, `supersedes/superseded_by`,
  `implements/implemented_by`, `caused_by/causes`,
  `discovered_from/discovered`. `fid link add` accepts either direction as the
  declared relation, resolves endpoint type from the entity ID shape
  (`PRD-*`/`TRD-*` => `requirement`, everything else => `issue` in v0.1),
  rejects unknown relations and `from == to` before append, rejects
  incompatible relation/endpoint combinations in the app layer before the SQL
  trigger fires, and rejects duplicate declared links before append. The write
  path appends the declared row plus auto-created inverse row in one journal
  event under the held-lock freshness helper.
- **Consequences:** v0.1 link behavior stays within the existing schema and
  trigger contract, produces actionable CLI errors instead of raw SQL aborts,
  and does not widen endpoint scope to decisions/escalations before Task 3/4 or
  v0.2 governance work. No section-13 mismatch was found: the in-code
  `link.add` and `link.remove` matrix rows now match `03_TRD_SCHEMA.md`
  section 13 for all seven roles.
- **Refs:** PI-2; phase-12/task-2; VAL-001; VAL-005; VAL-006; VAL-007;
  VAL-008; `03_TRD_SCHEMA.md` sections 10, 13, 14.2-14.3;
  `06_DESIGN_PLAN_CLI.md` `fid link add`.
