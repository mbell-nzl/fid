# Build-Readiness Spec Audit

Verdict: **NOT-READY** for first implementation commit.

The roadmap is broadly sliceable, and the prior consistency/coverage findings
are not re-raised here. The remaining problems are narrower: several v0.1
tasks still force an implementer to invent bootstrap, replay, and version-gate
behavior before a correct first test can be written.

## v0.1 Dependency-Ordered Build Probe

This is not a task plan; it is the ordering I would use to prove whether the
spec can be built from an empty repository.

1. **Phase 0 walking skeleton.** Go module, `cmd/fid`, migrations directory,
   testdata, Makefile, lint config, `fid --help`. Buildable from roadmap and
   quality plan.
2. **`.fid/` layout and config parsing.** Create `.fid/`, config, pricing,
   journal directory, and gitignore. **Blocked by BR-001:** init's relationship
   to the canonical journal is not defined.
3. **v0.1 SQLite migration.** Core issue tables, stub governance tables needed
   by v0.1 FKs/views, `projection_meta`, roles, indexes, and active triggers.
   **Blocked by BR-004:** the v0.1 trigger/DDL inventory includes later-version
   surfaces.
4. **Journal envelope, event ID, clock, and file writer.** Append one canonical
   JSONL event, flush/sync, hash line, and expose injectable time/ID providers
   for tests. **Blocked by BR-005:** event IDs and journal file rollover are
   examples, not a contract.
5. **`fid init` bootstrap event or deterministic seed.** Decide whether roles
   and initial project metadata are journal facts or schema/config bootstrap
   facts. **Blocked by BR-001.**
6. **Projection/replay and `projection_meta`.** Apply event effects into a temp
   SQLite DB, update checkpoint, rebuild from scratch, and catch up from tail.
   Buildable once BR-001 and BR-005 are fixed.
7. **Read-time freshness/rebuild boundary.** On every command, detect missing
   or stale DB and catch up before answering. **Blocked by BR-003:** catch-up
   mutates SQLite but its lock behavior is not specified.
8. **ID generation for generated entity rows.** `<PREFIX>-<hash>`, progressive
   length, retry on local insert collision, replay conflict on post-merge
   collision. Buildable; post-merge conflict reporting is sufficiently defined.
9. **Issue CRUD and read formats.** `issue create/get/list/update`, labels at
   creation/update, `issues_enriched`, JSON/TOON/markdown/YAML. Buildable after
   schema and journal foundations.
10. **Hard-coded workflow engine and rules.** Default state graphs, categories,
    `all_acs_verified`, `has_completion_report`, `concerns_addressed`,
    `all_child_tasks_done_or_cancelled`. Buildable and unit-testable.
11. **ACs, completions, gate runs, and concerns.** Append-only AC
    verifications, completion attempts, gate snapshots, concern lifecycle.
    Buildable once workflow and claims exist.
12. **Claims and permissions.** In-code default permissions matrix, unknown
    role rejection, claim acquire/renew/release/force semantics, stale takeover.
    Buildable.
13. **Links.** Issue-to-issue links are buildable. **Blocked by BR-002:**
    v0.1 issue-to-requirement links need requirement rows, but v0.1 has no
    command that can create them.
14. **Telemetry and stats.** `run record`, snapshot pricing, `stats cost`,
    `stats tokens`. Buildable, with pricing-file edge behavior safe to pin in
    task tests.
15. **Developer composites.** `dev start/gate/finish/block` as one event with
    many effects and atomic projection. Buildable once primitive effects,
    workflow, claims, and telemetry are in place.
16. **Maintenance, SQL, activity, and journal inspect.** `rebuild`, `vacuum`,
    read-only SQL, consistency check, activity list, event inspection. Buildable
    after replay and projection freshness are fixed.

No unavoidable dependency cycle appears in the intended roadmap. The blockers
are missing contracts at the first bootstrap/replay layers and one v0.1 feature
that depends on v0.2 requirement ingestion.

## Findings

| id | severity | dimension | evidence | decision or behavior left unspecified | recommended resolution |
|---|---|---|---|---|---|
| BR-001 | BLOCKER | D1, D2, D3 | `docs/adr/0001-journal-canonical-db-derived.md:5`: "source of truth for all project state"; `02_TRD.md:73`: "every mutation flows through it"; `02_TRD.md:212`: "Write path (every CLI write)"; `10_IMPLEMENTATION_ROADMAP.md:63`: "`fid init` ... seed `roles` table ... create initial journal"; `04_TRD_CONFIG.md:113`: "Creates `.fid/journal/` directory (empty)"; `04_TRD_CONFIG.md:116`: "Seeds the `roles` table with the default profile role list." | Is `fid init` a journaled write? If yes, what is the first event and which effects does it contain? If no, are roles/config/pricing deterministic bootstrap facts outside the journal? What does `projection_meta` look like for an empty journal, and does `fid rebuild` re-seed roles without an event? | Add an explicit "bootstrap contract." Either define a `project.initialized` v1 event that captures prefix, role seed, and any DB seed effects, or state that init/config/pricing/default roles are deterministic bootstrap outside the journal and are recreated by migrations/rebuild. Include the expected empty-project `projection_meta` row and first rebuild test. |
| BR-002 | BLOCKER | D1 | `01_PRD.md:166`: "support links from issues to requirements"; `10_IMPLEMENTATION_ROADMAP.md:104`: "v0.1 scope: issue <-> issue, issue -> requirement"; `10_IMPLEMENTATION_ROADMAP.md:56`: "`requirements`" is listed only under "Stub tables"; `03_TRD_SCHEMA.md:972`: triggers validate polymorphic refs on `links`; `03_TRD_SCHEMA.md:979`: "v0.1 accepts `issue` and `requirement`"; `06_DESIGN_PLAN_CLI.md:2188`: "`fid requirement sync` (v0.2+)"; `06_DESIGN_PLAN_CLI.md:2195`: sync "upserts pointer rows". | A valid v0.1 issue-to-requirement link requires a `requirements` row, but the only specified way to create those rows ships in v0.2. An implementer cannot write a v0.1 happy-path test for PRD-LINK-002 using only v0.1 commands. | Choose one: move a minimal `fid requirement sync` or seed/import mechanism into v0.1 before link add; defer issue-to-requirement links and PRD-LINK-002 to v0.2; or explicitly allow natural requirement IDs in links without a `requirements` row and adjust polymorphic integrity/maintenance semantics. |
| BR-003 | MAJOR | D2, D3 | `02_TRD.md:408`: "all CLI writes acquire `.fid/write.lock`"; `02_TRD.md:442`: "On every CLI invocation, the DB is checked against the journal tail"; `02_TRD.md:443`: "caught up synchronously before the command answers"; `06_DESIGN_PLAN_CLI.md:1299`: "catch-up before reads is internal and does not require `--by`"; `docs/adr/0001-journal-canonical-db-derived.md:5`: "no background daemon and no sync window." | Read-side catch-up/rebuild mutates SQLite, but only CLI writes are explicitly inside `.fid/write.lock`. If a writer has appended the journal event and is still projecting, a concurrent reader could see a stale DB and try to catch up against the same event. The spec does not say whether readers wait, acquire the same lock for catch-up, return a lock error, or may serve the old projection. | State that any DB-mutating projection path, including automatic read catch-up and explicit rebuild, acquires the same checkout-scoped projection/write lock before touching SQLite. Specify reader behavior when the lock is held and stale reads are forbidden. Add an integration test that pauses a writer after append and verifies a reader waits or exits with the specified lock error. |
| BR-004 | MAJOR | D1, D4 | `10_IMPLEMENTATION_ROADMAP.md:56`: v0.1 tables exclude `comments` and `qa_reviews`; `10_IMPLEMENTATION_ROADMAP.md:57`: v0.1 triggers include "comment scope" and "mutually-exclusive qa_review columns"; `03_TRD_SCHEMA.md:306`: "QA Layer (v0.3)"; `03_TRD_SCHEMA.md:309`: `CREATE TABLE qa_reviews`; `03_TRD_SCHEMA.md:561`: "Comments (v0.2; qa_review scope added in v0.3)"; `03_TRD_SCHEMA.md:485`: "v0.2 run/git/lifecycle extension"; `03_TRD_SCHEMA.md:486`: `status TEXT NOT NULL DEFAULT 'open'`; `03_TRD_SCHEMA.md:487`: `branch TEXT`; `10_IMPLEMENTATION_ROADMAP.md:185`: run git context is an "Additive migration to `agent_runs`." | The first migration author cannot tell whether to implement the DDL exactly as shown, which creates v0.2 run lifecycle/git columns in the v0.1 `agent_runs` table and asks for v0.1 triggers/checks over v0.2/v0.3 tables, or to mentally filter the schema by roadmap version. This also weakens the v0.1->v0.2 seam because the "additive" migration may already be present. | Add a per-migration DDL inventory or annotated SQL blocks: v0.1 migration includes only v0.1 tables, stubs, and active triggers; v0.2 adds comments plus `agent_runs` lifecycle/git columns; v0.3 adds QA review tables/checks. If the intent is to ship future columns early, update the roadmap and PRD version gates to say so. |
| BR-005 | MAJOR | D2, D3 | `02_TRD.md:223`: appends to `.fid/journal/YYYY-MM-DD.jsonl`; `02_TRD.md:263`: "`event_id` ... Stable generated event ID"; `02_TRD.md:266`: "`occurred_at` ... UTC ISO 8601 timestamp"; `02_TRD.md:318`: example `"event_id": "EVT-20260529-000001"`; `04_TRD_CONFIG.md:17`: "one file per day"; `03_TRD_SCHEMA.md:13`: generated entity IDs use `<PREFIX>-<hash>`, but the prefix table does not define `EVT`. | The journal compatibility surface needs an event ID and day-rollover contract, but the spec only gives examples. The implementer must invent whether `EVT-YYYYMMDD-NNNNNN` is required, whether the sequence is per day/file/project, how to allocate under the write lock, whether the journal filename uses UTC or local date, and what happens if the clock crosses midnight or moves backward during a write. | Define the event ID and journal date contract before the journal writer task. Example resolution: derive both `occurred_at` and filename from one injected UTC clock read; use `EVT-<UTC YYYYMMDD>-<6 digit sequence>` allocated by scanning the target file under `.fid/write.lock`; retry or fail on collision; tests inject clock/sequence providers. A UUID-style event ID is also acceptable if documented and examples are updated. |

## Definition Of Done And Testability

The acceptance bar is mostly concrete. The v0.1 roadmap has explicit exit
criteria, including rebuild, full issue lifecycle, claims, links, output
formats, maintenance check, telemetry, and "All v0.1 PRD requirements have
passing tests" (`10_IMPLEMENTATION_ROADMAP.md:127`, `10_IMPLEMENTATION_ROADMAP.md:140`).
The task-level gate is concrete (`08_QUALITY_PLAN.md:231` through
`08_QUALITY_PLAN.md:239`), and the phase gate adds integration, benchmarks,
coverage, and review steps (`08_QUALITY_PLAN.md:245` through
`08_QUALITY_PLAN.md:253`).

Most v0.1 behavior can be tested with stdlib/testify, a temp SQLite DB, and a
temp `.fid/journal`. The exceptions are exactly the blockers above:

- init/rebuild role seeding cannot be tested until BR-001 is fixed;
- valid issue-to-requirement links cannot be tested from v0.1 commands until
  BR-002 is fixed;
- read catch-up under concurrent writer projection needs the BR-003 lock rule;
- first journal golden files need BR-005's event ID/date contract.

Golden replay fixtures are conceptually specified: `08_QUALITY_PLAN.md:130`
through `08_QUALITY_PLAN.md:132` says to capture a real event at schema version
N, run it through the upgrader chain, and compare against the current shape.
That is enough for v0.2+ upgraders, but the first v0.1 fixture still depends on
the bootstrap and event ID decisions above.

## v0.1 -> v0.2 Seams

The workflow and permission migration seams are otherwise sound. The TRD says
hard-coded workflow definitions migrate verbatim into workflow tables
(`02_TRD.md:508` through `02_TRD.md:511`), and it requires named rule functions
with v0.2 migration tests (`02_TRD.md:719`). The permissions matrix is a v0.1
code constant migrated to `role_permissions` in v0.2 with tests for every write
action (`02_TRD.md:720`), matching ADR-0002's "No retroactive breakage" rule
(`docs/adr/0002-configurable-roles-default-profile.md:34`).

The unsafe seam is BR-004: `agent_runs.status` and git context columns are
documented as v0.2 additive fields but appear inside the base `CREATE TABLE
agent_runs`, and v0.1 trigger scope names future tables. That needs a precise
per-version migration inventory before implementation.

## Safe To Defer

- Exact numeric exit-code taxonomy beyond zero vs non-zero. The current spec
  gives command messages and success/error polarity; reserve a richer taxonomy
  until there is a programmatic client need.
- Exact retry count for local hash collisions. The spec says retry with a fresh
  UUID and error if collision persists; a bounded implementation choice can be
  made in the ID-generation task.
- Exact stderr wording for fsync/append failure, projection mismatch, and lock
  timeout. The important behavior is already specified: no DB mutation before
  journal commit, stop before reads on projection inconsistency, non-zero on
  lock timeout. Pin wording in the relevant CLI tests.
- Missing or empty `pricing.toml` behavior. v0.1 telemetry can safely treat it
  like a missing model: cost NULL plus warning, as long as the run-record task
  documents the choice and tests it.
- Human-readable markdown polish for every read command. JSON shapes are the
  compatibility surface; markdown can be refined while preserving fields.

## Must Fix Before First Commit

1. **Define the bootstrap contract** for `fid init`, roles, empty journal,
   `projection_meta`, and delete-DB rebuild behavior (BR-001).
2. **Resolve v0.1 requirement links** by moving requirement row creation to
   v0.1, deferring requirement links, or changing link integrity semantics
   (BR-002).
3. **Specify read catch-up locking** so no reader races a writer between journal
   append and SQLite projection (BR-003).
4. **Publish a per-version migration inventory** that removes future trigger and
   `agent_runs` column ambiguity from v0.1 (BR-004).
5. **Pin the event ID and journal day-rollover contract** before writing the
   first journal golden fixture (BR-005).
