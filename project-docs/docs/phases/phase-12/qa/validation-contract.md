# Phase 12 Validation Contract - Requirements + Links

> Phase: 12 - Requirements + links (roadmap v0.1d, slice 3 of 4)
> Author: Tech Lead (draft)
> Last refined: 2026-06-22 by Tech Lead (planning draft)
>
> Location: `project-docs/docs/phases/phase-12/qa/validation-contract.md`.

## Purpose

This file is the finite set of validation assertions that must be true when
Phase 12 is complete. The contract specifies **what must be proven**; the QA
plan specifies **how**.

If a required assertion cannot be satisfied, stop and escalate rather than
weakening it.

## Phase Scoping Notes

Phase 12 is the third v0.1d slice: requirement pointers and typed links. It
ships `fid requirement sync` (minimal), `fid link add`, `fid link remove`, and
`fid link list`. Requirements ship first so issue -> requirement `implements`
links have valid endpoints. All writes are gated by the Phase 10 permissions
matrix; new matrix rows are appended in this phase.

The `links` and `requirements` tables, the `links` endpoint-scope /
ref-integrity / relation-compatibility triggers, and v0.1 issue/requirement
endpoint scope already exist (`migrations/0002`, `migrations/0005`;
`03_TRD_SCHEMA.md` sections 5.2, 10, 14.2-14.3). Phase 12 must not add a schema
migration unless implementation proves a recorded spec gap. The `link.add`,
`link.remove`, and `requirement.sync` rows already appear in
`03_TRD_SCHEMA.md` section 13; Tasks 1-2 append those rows to the in-code matrix
that Phase 10 intentionally kept limited to implemented actions.

Two planning interpretations apply (surfaced to the human for confirmation;
to be recorded as phase decisions during execution):

- **PI-1 (`requirement.sync` parse contract).** The canonical source format is
  the PRD/TRD markdown requirement tables: a row whose first cell is a canonical
  requirement ID (`(PRD|TRD)-[A-Z0-9]+-[0-9]+[a-z]?`) and whose second cell is
  the summary. `type` is derived from the ID prefix and may be filtered by
  `--type`. `--id` restricts the sync to a single requirement. `status` is set
  to `active`. `last_synced_hash` is a deterministic hash over the parsed
  summary so v0.2 drift detection has an anchor. `source_section` is populated
  best-effort from the nearest preceding markdown heading, else null. This is a
  public-interface contract; the exact hashed content is recorded in the Task 1
  completion report.
- **PI-2 (PRD-LINK-008 v0.1 enforcement).** In v0.1 the closed seven-relation
  set plus issue/requirement-only endpoint scope structurally admit no relation
  that re-encodes a dedicated-schema relationship (issue hierarchy via
  `parent_id`, escalation, resolution). The app layer mirrors the
  relation-compatibility trigger so callers get a before-append error, rejects
  unknown relations, and rejects self-links. Decision/escalation
  dedicated-schema duplication activates with link endpoints in v0.2
  (PRD-LINK-007).

Carry-forward review at kickoff:

- CF-008 remains deferred to v0.2.
- CF-011 remains deferred to v0.1e, but PRD-LINK-006a explicitly requires link
  removal to project to `activity_log`; Phase 12 implements that per-command
  projection now (link add does **not** project to `activity_log`, matching the
  `issue create` precedent the roadmap calls out).
- CF-013 and CF-014 remain deferred to v0.1e.
- Lock-ordering lesson (DEC-028): every new write command
  (`requirement sync`, `link add`, `link remove`) refreshes the projection with
  the held-lock freshness helper inside `journal.WithLock`; no nested
  `.fid/write.lock` acquisition.

## Numbering

- IDs use `VAL-NNN`, scoped to Phase 12, sequential from `VAL-001`.
- Carry-forward assertions are re-issued with a new `VAL-NNN` and a source note.
- Retired/superseded assertions stay in the table; IDs are never reused.

## Assertions

| ID | Source requirement | Behavior or invariant | Evidence method | Validator | Owning task/window | Pass criterion | Fail criterion | Status | Out-of-scope notes |
|---|---|---|---|---|---|---|---|---|---|
| VAL-001 | TRD-ARCH-010; `03_TRD_SCHEMA.md` section 13; roadmap v0.1d Phase 12 | The in-code permissions matrix contains exactly one row each for `requirement.sync`, `link.add`, and `link.remove` with the same role booleans as `03_TRD_SCHEMA.md` (requirement.sync: tl + human; link.add / link.remove: all seven roles), and `fid policy show` displays those rows from the same source `Authorize` enforces. | Unit + CLI: matrix completeness test updated for the three actions; `Authorize` allowed/denied pairs; `fid policy show` output includes the rows and matches the matrix object. | QA Lead | tasks-1-2 | All three actions have exactly one in-code row; allowed/denied role checks match section 13; policy show matches the matrix. | A row is missing/duplicated, booleans diverge from section 13, or show/enforce drift. | Draft | Label actions remain the Phase 13 v0.1d slice. |
| VAL-002 | PRD-GOV-006; PRD-LINK-002; `06_DESIGN_PLAN_CLI.md` `fid requirement sync`; PI-1 | `fid requirement sync --source-doc <path> --by <role>` upserts requirement pointer rows in one journal event: new IDs create rows, changed summaries update rows, identical content is unchanged; persisted rows expose `id`, `type`, `summary`, `source_doc`, `source_section`, `last_synced_hash`, `status='active'`; output reports created/updated/unchanged counts; a sync with zero creates/updates appends no journal event. | Domain + CLI: sync a fixture doc; read back rows; re-sync unchanged (unchanged++ / no new event); edit a summary and re-sync (updated++ and hash changes). | QA Lead | task-1 | First sync creates the expected rows with all fields; identical re-sync reports unchanged and appends nothing; a changed summary reports updated and changes `last_synced_hash`. | A field is dropped, counts are wrong, unchanged re-sync appends an event, or a changed summary is reported unchanged. | Draft | Requirement get/list and drift reporting are v0.2 (PRD-GOV-007). |
| VAL-003 | PRD-GOV-006; PRD-ERR-003; CF-007; `06_DESIGN_PLAN_CLI.md` `fid requirement sync` | `fid requirement sync` rejects malformed input before append with non-zero exit and zero new journal lines/rows: a source with no parseable requirement IDs, duplicate IDs within the source, an unknown `--type`, and (with `--id`) an ID absent from the source. | CLI: run each error case against fixtures; compare journal line count / DB state before and after; assert actionable stderr. | QA Lead | task-1 | Each error case exits non-zero, names the problem, and leaves journal/DB unchanged. | Any error case appends/projects, exits zero, or returns an unactionable error. | Draft | Drift reporting is v0.2; v0.1 only stores the hash. |
| VAL-004 | TRD-ARCH-010; `03_TRD_SCHEMA.md` section 13; CF-007 | `requirement.sync` is authorized only for roles permitted in section 13 (tl, human); a denied role and an unknown role are rejected before append with an error naming role + action and zero new journal lines. | CLI: sync as tl (allowed) and as a denied role (e.g. developer) and an unknown role; assert allow path appends one event and deny paths append nothing. | QA Lead | task-1 | tl/human sync succeeds; denied/unknown role rejects before append naming role + `requirement.sync`. | A denied/unknown role succeeds or appends, or an allowed role is rejected. | Draft | Configurable role permissions are v0.2 (PRD-ROLE-006). |
| VAL-005 | PRD-LINK-001; PRD-LINK-003; PRD-LINK-004; `06_DESIGN_PLAN_CLI.md` `fid link add` | `fid link add <from> <relation> <to> --by <role>` for a valid issue <-> issue relation creates the declared link and its auto-inverse in one journal event, projects both `links` rows, and prints a confirmation naming both directions; re-adding an existing link is rejected before append naming the existing relation/entities. | Domain + CLI: add `blocks` between two issues; read back both `links` rows; assert single event and confirmation text; re-add and assert zero append. | QA Lead | task-2 | A valid link add creates exactly the declared + inverse rows in one event with correct inverse relation; duplicate add rejects before append. | Inverse missing/wrong, two events emitted, confirmation omits a direction, or duplicate add appends. | Draft | Decision/escalation endpoints are v0.2 (PRD-LINK-007). |
| VAL-006 | PRD-LINK-002; PRD-LINK-003; `02_TRD.md` section 4.x links | An issue -> requirement `implements` link can be created once the requirement pointer exists (Task 1), and its auto-inverse is the requirement -> issue `implemented_by` row; an `implements` link to a missing requirement ID is rejected before append. | Domain + CLI: sync a requirement, then `fid link add <issue> implements <REQ> --by tl`; read back both rows; attempt `implements` to an unsynced ID and assert rejection. | QA Lead | task-2 | The issue -> requirement link and its `implemented_by` inverse both project; a missing requirement endpoint rejects before append. | The inverse is wrong/absent, or a link to a nonexistent requirement is accepted. | Draft | Requirement reads are v0.2; only sync + endpoint existence are in scope. |
| VAL-007 | PRD-LINK-005; `03_TRD_SCHEMA.md` section 14.3 | `fid link add` rejects a relation/endpoint combination invalid for the endpoint pair before append, naming the relation and both endpoints: e.g. `implements` issue -> issue, `blocks` issue -> requirement, and a link whose endpoint does not resolve to an existing row. | CLI: run each incompatible combination and a missing-endpoint case; assert non-zero exit, zero append, and an error naming relation + endpoints; confirm a compatible combination still succeeds. | QA Lead | task-2 | Each incompatible/unresolvable combination rejects before append with relation + endpoints named; a compatible combination succeeds. | An incompatible combination is accepted/appended, or the error omits relation/endpoints. | Draft | Cross-type `relates_to` remains allowed by spec and is not an error. |
| VAL-008 | PRD-LINK-008; PRD-LINK-004; `03_TRD_SCHEMA.md` section 10; GLOSSARY "Link"; PI-2 | The v0.1 link relation registry contains exactly the seven semantic relations and their inverses, none of which re-encodes a dedicated-schema relationship (issue hierarchy, escalation, resolution); `fid link add` rejects an unknown relation and a self-link (`from == to`) before append. | Unit + CLI: assert the relation/inverse registry equals the documented set; attempt an unknown relation (e.g. `parent_of`) and a self-link and assert before-append rejection with zero new journal lines. | QA Lead | task-2 | Registry matches the documented seven relations + inverses exactly; unknown relation and self-link reject before append. | The registry includes a hierarchy/escalation/resolution relation, or an unknown relation / self-link is accepted. | Draft | Decision/escalation dedicated-schema duplication is v0.2 (PRD-LINK-007). |
| VAL-009 | PRD-LINK-006a; CF-011; `06_DESIGN_PLAN_CLI.md` `fid link remove` | `fid link remove <from> <relation> <to> --by <role>` removes the declared link and its inverse in one journal event, projects the removal to `activity_log`, and prints a confirmation naming both removed directions; removing a nonexistent link is rejected before append naming the missing relation/entities. | Domain + CLI: add then remove a link; assert both `links` rows gone in one event, an `activity_log` removal row exists, and confirmation names both directions; remove a nonexistent link and assert zero append. | QA Lead | task-3 | Removal deletes declared + inverse in one event, writes activity evidence, and confirms both directions; nonexistent removal rejects before append. | A row survives, two events are emitted, activity evidence is missing, or nonexistent removal appends. | Draft | CF-011 broad activity-log reconciliation remains v0.1e; link add does not project to activity_log. |
| VAL-010 | PRD-LINK-006; `06_DESIGN_PLAN_CLI.md` `fid link list`; `05_DESIGN_PLAN.md` cross-surface fields | `fid link list <id>` returns both outgoing and incoming links for the entity with `direction`, `relation`, `entity_type`, `entity_id`, and `entity_title`; json, toon, markdown, and yaml carry the same data. | CLI: build a fixture with outgoing and incoming links (issue and requirement endpoints); run list in all four formats; compare parsed structures to the JSON reference. | QA Lead | task-4 | All four formats expose the same outgoing + incoming links and required fields, including resolved `entity_title`. | A direction or field is missing, titles fail to resolve, or a format diverges from JSON. | Draft | Filtering by relation/direction is not specified for v0.1 list. |
| VAL-011 | Roadmap v0.1d Phase 12 exit; PRD-NFR-008; ADR-0001; DEC-028 | An end-to-end flow through real commands against a temp workspace holds together and stays rebuild-stable: requirement sync -> issue <-> issue link -> issue -> requirement link -> `link list` (both directions) -> `link remove`; each write appends exactly one event under the held lock with no nested-lock failure, and `fid rebuild` reproduces the same link/requirement state. | Integration test: scripted command sequence; assert per-write event counts, no nested-lock error, and rebuild parity of `links`/`requirements`. | QA Lead | task-4 | The full sequence runs, each write appends one event with no nested-lock failure, and rebuild reproduces identical state. | Any write fails under the lock, emits the wrong event count, or rebuild diverges. | Draft | Full v0.1d lifecycle closure smoke (claim/transition/complete) is Phase 13. |
| VAL-012 | `08_QUALITY_PLAN.md` sections 5.2, 5.6; phase goal | A human dogfoods a real `.fid` workspace and can tell from CLI output that a requirement sync recorded pointers, that a link add created an inverse, that `link list` shows both directions, and that `link remove` removed both directions. | Human checklist with one row per observable; QA Lead generates from this VAL after deterministic QA passes. | Human | phase-end | Human confirms sync counts, inverse creation/removal, and bidirectional listing are understandable without reading implementation. | Human cannot identify those outcomes from CLI output, or the checklist is not run. | Draft | QA Lead fans this into per-observable human-checklist rows. |

## Coverage Rules

Confirmed for Phase 12 execution:

- Phase-scope requirements each map to at least one `VAL-NNN`: PRD-LINK-001 to
  VAL-005; PRD-LINK-002 to VAL-006; PRD-LINK-003 to VAL-005/006; PRD-LINK-004
  to VAL-005/008; PRD-LINK-005 to VAL-007; PRD-LINK-006 to VAL-010;
  PRD-LINK-006a to VAL-009; PRD-LINK-008 to VAL-008; PRD-GOV-006 to
  VAL-002/003; TRD-ARCH-010 to VAL-001/004.
- Every `VAL-NNN` is referenced by a task spec or explicitly marked
  `phase-end` (VAL-012).
- Binary/conditional assertions name both states: VAL-002, VAL-003, VAL-004,
  VAL-005, VAL-006, VAL-007, VAL-008, VAL-009.
- VAL-001, VAL-008, and VAL-010 carry structural/parity conformance proofs
  (matrix-vs-section-13, relation-registry-vs-documented-set, four-format
  parity-vs-JSON).
- Coverage is enforced at the gate, not as a `VAL-NNN`: COV-001 domain package
  coverage and COV-002 `cmd/fid` coverage must stay above thresholds. Each task
  runs the required coverage gate and records package percentages.
- QA Lead maps each `VAL-NNN` to deterministic `(qa-lead)`, mid-phase review,
  or `(human)` verification, and fans VAL-012 into per-observable human rows.
- Standing QA-plan checks carry forward: show-vs-enforce parity for the new
  matrix rows (VAL-001), validate-before-append on every rejection
  (VAL-003/004/005/007/008/009), and no stale-projection nested-lock failure on
  the new write paths (VAL-011).

## Outcomes

Phase-end final review records every `VAL-NNN` outcome: `PASS`, `FAIL`,
`BLOCKED`, or `CARRY-FORWARD`. `PASS` without evidence is invalid.

## Revision Log

| Date | Assertion | Change | Reason | Decision reference |
|---|---|---|---|---|
| 2026-06-22 | VAL-001..012 | created | Phase 12 planning draft | roadmap v0.1d Phase 12; PRD-LINK-001..008; PRD-GOV-006; TRD-ARCH-010 |
