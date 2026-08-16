# Changelog

All notable changes to fid are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Changed

- Phase 20 delivers the standalone Issue journey: project initialization and
  rebuild, legal Issue hierarchy creation, acceptance criteria, labels, typed
  links, ready selection, revisioned updates/transitions, workflow inspection,
  and maintenance inspection now use the `fid-cli.v1` JSON Response Envelope.
- The SQLite projection remains derived from the append-only journal; canonical
  `fid --json issue get` output is preserved byte-for-byte across rebuild.

### Removed

- Universal `--format` output selection and its TOON, YAML, and generic
  Markdown renderers are removed. JSON is the retained machine interface.
- The promised `fid decision`, `fid escalation`, and `fid triage` command stubs
  are removed rather than advertised ahead of their owning phases.

The `0.1.0` notes below are historical release context. In particular, their
formatter and governance-stub descriptions do not describe the current Phase
20 command surface.

## [0.1.0] - 2026-07-18

First release of fid — the Core Issue Layer. Delivered across Phases 0–17 and
closed for early dogfood (DEC-057). The journal is the source of truth and is
committed to git; the SQLite database is a derived, gitignored projection with
no daemon. Static binary ~12.3 MiB (target < 30 MB).

### Added

**Persistence foundation**
- `.fid/` workspace layout with canonical `config.toml` and `pricing.toml`; the
  `fid.sqlite` projection is gitignored.
- SQLite schema via `golang-migrate`, with validation triggers as an
  application-layer backstop and required pragmas (`foreign_keys`,
  `recursive_triggers`).
- Append-only JSONL journal (`event_schema` v1) with a checkout-scoped
  write/projection lock; one logical write is one event and the journal append
  is the commit point.
- Journal replayer with versioned upgraders, corruption detection (malformed
  trailing lines reported, not skipped), logical-conflict precondition errors,
  and startup catch-up when the DB is behind or missing.
- `fid init` (seeds the 7 default roles via `project.initialized`, `--prefix`),
  `fid rebuild` (deterministic sub-second replay into a fresh DB).
- Entity ID generation (`<PREFIX>-<hash>`) with progressive length scaling and
  collision retry.
- Output format engine — JSON (default), TOON, Markdown, YAML — via `--format`
  on read commands.

**Issue lifecycle**
- `fid issue create` / `get` / `list` / `update` / `transition`, with kind
  hierarchy (epic / task / subtask), owner/priority/sequence/label/estimate
  fields, `--section` on-demand markdown, and the enriched read model
  (`issues_enriched`).
- Epic metadata with sequential / prioritised scheduling.
- Hard-coded default workflows (epic / task / subtask) and named transition
  rules (`all_acs_verified`, `has_completion_report`, `concerns_addressed`,
  `all_child_tasks_done_or_cancelled`) with rich rejection messages.
- `fid ready` — next actionable task for sequential and prioritised epics,
  unattached tasks, and cross-project selection, with empty-result explanations.

**Acceptance criteria, completions, concerns**
- `fid issue ac add` / `verify` / `list` — ordered ACs, append-only multi-role
  verifications with evidence type/ref, latest-status-per-role.
- `fid issue complete` — append-only `task_completions` with attempt numbering,
  prior-attempt supersession, gate-result snapshot, and human-intervention flag.
- `fid concern raise` / `address` / `list`, with derived `concerns_raised` /
  `open_concerns_count` on issue reads.

**Claims, permissions, links, labels**
- `fid issue claim` / `unclaim` and `fid claim list` — atomic acquire with TTL,
  `--renew`, `--force`, `--stale`; claim/transition integration (started-category
  transitions require a live claim; transitions out auto-release); `--by human`
  hold signal.
- In-code permissions matrix enforced on every write via mandatory `--by <role>`;
  unknown/denied role-action combinations rejected before append;
  `fid policy show`.
- `fid requirement sync` (minimal pointer upsert) for issue → requirement links.
- `fid link add` / `remove` / `list` with auto-inverse projection and
  relation/endpoint validation.
- `fid label add` / `list` and `fid issue label add` / `remove`.

**Telemetry, composites, maintenance**
- `fid run record` (snapshot pricing + computed `cost_usd`, unknown-model
  warning) and `fid stats cost` / `tokens` aggregated by role / model / epic;
  `pricing.toml` ships with current provider defaults.
- `fid dev start` / `gate` / `finish` / `block` composite developer commands
  (one logical write = one event).
- `fid sql` (read-only escape hatch) and `fid journal inspect <event_id>`.
- `fid maintenance check` (structural + warning audit, `--warn-as-error`) and
  `fid maintenance vacuum`.
- `fid activity list` — projected timeline with entity / event / date filters.
- Stub-table CLI guards for the v0.2 governance families (`decision`,
  `escalation`, `triage`).

### Notes

- v0.1 exit criteria met: init → rebuild reconstructs the DB (roles included) and
  delete-and-rebuild is non-destructive; the full issue lifecycle
  (create → claim → transition → complete → unclaim) works under workflow and
  permission enforcement; every v0.1 PRD requirement maps to passing tests.
- Schema-only stub tables (`decisions`, `escalations`, `triage_events`) ship
  without a write CLI; their commands land in v0.2.

[Unreleased]: https://github.com/mbell-nzl/fid/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/mbell-nzl/fid/releases/tag/v0.1.0
