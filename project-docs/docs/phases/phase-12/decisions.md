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

### DEC-031: One `link_removed` activity event records both link directions for PRD-LINK-006a

- **Date:** 2026-06-27
- **Context:** TRG-002 / Review 2 surfaced that Task 3 objective prose and
  PRD-LINK-006a say "project both removals to the activity log," while
  `fid link remove` writes a single `link_removed` `activity_log` row rather
  than two rows. The validation contract (VAL-009), QA plan, and acceptance
  criteria require activity **evidence** of the removal, not a literal two-row
  projection. Final review (Improvement #2) asked Tech Lead wrap-up to
  disposition whether PRD-LINK-006a wording needs a spec sync note.
- **Decision:** A single `link_removed` activity event satisfies PRD-LINK-006a.
  `fid link remove` deletes the declared link and its auto-inverse as one
  logical write (one journal event) and projects one `activity_log` row whose
  `details` payload names both directions: the declared `from_*`/`relation`/
  `to_*` and the `inverse_relation`/`inverse_from_*`/`inverse_to_*`
  (`internal/links/remove.go:199`). "Project both removals to the activity log"
  is read as "both removals are auditable in the activity log," not "emit two
  activity rows." This matches fid's one-logical-write-one-event model (DEC-028
  era) and the `link add` / `issue create` precedent that one logical action
  yields one activity-relevant record.
- **Consequences:** No schema migration and no spec text change. `01_PRD.md`
  PRD-LINK-006a receives a `verified-no-change` verdict with this
  interpretation recorded; the activity row `details` payload is the durable
  contract for "both directions removed." Future `fid activity list` work
  (CF-011) renders both directions from that one row.
- **Refs:** PRD-LINK-006a; TRG-002; VAL-009; `internal/links/remove.go`;
  phase-12/task-3; `01_PRD.md` PRD-LINK-006a.

## Carry-Forward Review At Kickoff

| ID | Disposition | Rationale |
|---|---|---|
| CF-008 | Re-defer to v0.2 | Resolution-reason field is a product/data-model enhancement unrelated to requirements/links. |
| CF-011 | Re-defer broad reconciliation to v0.1e; deliver link-removal projection now | PRD-LINK-006a requires link removal to project to `activity_log`, so Phase 12 handles that per-command case (DEC-031) without widening the broader timeline policy. Link add does not project, matching the `issue create` precedent. |
| CF-013 | Re-defer to v0.1e | `complete.go` orchestration cleanup belongs with the `fid dev *` composites, not requirement/link primitives. |
| CF-014 | Re-defer to v0.1e | `cmd/fid/issue.go` / `replay.go` file-size split is a behaviour-preserving refactor; Phase 12 touched neither file materially. |

## Wrap-Up Dispositions (Tech Lead, 2026-06-27)

### Validation contract outcome audit

All twelve assertions reach a terminal state.

- VAL-001..011: PASS in `qa/final-review.md` with cited test/CLI evidence. No
  disposition needed.
- VAL-012: was CARRY-FORWARD to the human checklist in final review. The
  generated `qa/human-checklist.md` (HV-12a..HV-12d) is now complete with
  overall **PASS** (completed 2026-06-27). Final outcome: **PASS**. No
  FAIL/BLOCKED/CARRY-FORWARD remains; no defer/cancel/escalation required.

### QA recommendation disposition

Final review had no Blocking or Major findings. Two Improvement findings:

- **VAL-012 human dogfood** — resolved: the human checklist ran and passed
  overall (HV-12a..HV-12d all `pass`). No action remains.
- **PRD-LINK-006a wording** — resolved by **DEC-031**: the single `link_removed`
  activity row whose `details` payload names both the declared and inverse
  directions satisfies "project both removals to the activity log." No spec
  text change; `01_PRD.md` verdict is `verified-no-change`.

### DEC / deviation spec-impact walk

- DEC-029 (`requirement sync` hashes the parsed summary cell only): spec-
  impacting? No. `03_TRD_SCHEMA.md` §§5.2/13 already expose the required
  `requirements` columns and the `requirement.sync` §13 row; matrix matches
  (VAL-001/002/003 PASS) — `verified-no-change`.
- DEC-030 (closed v0.1 relation registry + app-layer prechecks): spec-
  impacting? No. It mirrors the existing `links_*` triggers and §10 relation
  set; in-code `link.add`/`link.remove` rows match §13 (VAL-001/005..008 PASS) —
  `verified-no-change`.
- DEC-031 (one `link_removed` event records both directions): spec-impacting?
  No. PRD-LINK-006a intent (both removals auditable) is MET; interpretation
  recorded here, no PRD edit — `verified-no-change`.
- Developer deviations: none reported across Tasks 1-4.
- Validation-contract deviations: none; all twelve VAL IDs were satisfied as
  written, none weakened.

## Spec Sync Verdict

Record a verdict for every spec at phase wrap-up. `verified-no-change` is a
positive assertion after review, not the default.

| Spec | Verdict | Evidence / Reference |
|---|---|---|
| `00_PROJECT_CHARTER.md` | N/A | No charter-level scope/constraint change in a requirements/links slice. |
| `01_PRD.md` | verified-no-change | PRD-LINK-001..006a, PRD-LINK-008, PRD-GOV-006 all MET (final-review Requirement Compliance); PRD-LINK-006a interpretation recorded in DEC-031. |
| `02_TRD.md` | verified-no-change | TRD-ARCH-010 matrix satisfied (VAL-001/004); one-event-per-write and rebuild parity held (VAL-011); no TRD change. |
| `03_TRD_SCHEMA.md` | verified-no-change | DEC-029/030; no migration added; §13 matrix rows match the in-code matrix (VAL-001 PASS); §§5.2/10/14.2-14.3 satisfied as written. |
| `04_TRD_CONFIG.md` | N/A | No `.fid/` directory or config change. |
| `05_DESIGN_PLAN.md` | verified-no-change | Cross-surface link fields (`direction`, `relation`, `entity_*`, resolved title) consistent (VAL-010 PASS). |
| `06_DESIGN_PLAN_CLI.md` | verified-no-change | `fid requirement sync` and `fid link add/remove/list` flags, counts, confirmation/inverse output, and four-format shape match implementation (VAL-002..011); task reports record no CLI spec change required. |
| `07_DESIGN_PLAN_FLOWS.md` | verified-no-change | End-to-end sync -> issue<->issue link -> issue->requirement link -> list -> remove flow matches documented behaviour (VAL-011 PASS). |
| `08_QUALITY_PLAN.md` | verified-no-change | Phase gate run as defined; COV-001 domain packages >=80% and COV-002 `cmd/fid` 77.9% met; no new gate added. |
| `09_AI_WORKFLOW.md` | N/A | No agent-workflow change. |
| `10_IMPLEMENTATION_ROADMAP.md` | verified-no-change | Phase 12 delivered exactly the v0.1d Phase 12 entry (requirement sync + link add/remove/list with auto-inverse); no re-slice. |
| `GLOSSARY.md` | verified-no-change | "Link" and "Requirement" entries already cover typed links with auto-inverse and PRD/TRD requirement pointers; the closed seven-relation set lives in `03_TRD_SCHEMA.md` §10. Relation registry / auto-inverse are implementation mechanics, not new canonical language. |

## Triage-Log Audit

- [x] Walked every TRG-NNN in `triage-log.md` (TRG-001, TRG-002).
- [x] Confirmed each RESUME verdict was correct in hindsight.
- [x] Flagged any TRG that should have escalated.
- [x] Confirmed autonomous TL-Lite resume remains in effect (DEC-002; no warm-up tally).

Audit note: both triage entries were **RESUME** and both were correct. TRG-001
(Review 1) and TRG-002 (Review 2) each recorded zero blocking findings, no
escalations, and a CONTINUE verdict; their notes were report-quality/parity-
citation improvements. TRG-002's cross-cutting wording observation (one vs two
activity rows) was correctly resumed under the validation contract and is now
durably closed by DEC-031. No TRG was under-escalated; no RESUME should have
been a BLOCK.

## Workflow Notes

- **What worked well?** Sequencing requirements before links (Task 1 before
  Task 2) meant issue->requirement `implements` links always had valid endpoints
  to test against; no endpoint-ordering rework. DEC-028's held-lock freshness
  helper was applied from the first write command, so no nested-lock surprise
  appeared this phase (the Phase 11 lesson held). The closed relation registry
  (DEC-030) let the app reject invalid relations before the SQL trigger, giving
  actionable CLI errors instead of raw aborts.
- **What was inefficient?** The PRD-LINK-006a one-vs-two activity-row ambiguity
  re-surfaced at Review 1, Review 2, final review, and wrap. Pinning the
  interpretation as a phase decision at kickoff (it was visible in the
  validation contract scoping notes) would have saved three rounds of "accepted
  under the contract" commentary; DEC-031 now closes it.
- **Were tasks the right size?** Yes. Four tasks each owned one coherent surface
  (requirement sync; link add + registry; link remove; link list + end-to-end),
  each one session, with clean dependencies. No mid-phase do-now task was needed.
- **Did the QA review cadence feel right?** Yes. Two mid-phase reviews over a
  four-task phase (a 2+2 window) gave per-task signal; both were CONTINUE with
  improvement notes only, and final QA closed every note.
- **Were there token-efficiency opportunities missed?** Minor. Both reviews
  asked completion reports to cite the exact design/CLI sections checked during
  polish/parity review; richer inline citations in the task reports would have
  saved final QA a re-derivation.

### Deep Retrospective Recommendation

No deep retrospective required. Phase 12 closed PASS with no Blocking/Major
findings, all twelve VAL outcomes terminal, and no new carry-forward opened.
Carry forward into Phase 13 planning: pin any genuine spec-wording ambiguity
(like PRD-LINK-006a) as a phase decision at kickoff rather than letting it
recur through every QA round.
