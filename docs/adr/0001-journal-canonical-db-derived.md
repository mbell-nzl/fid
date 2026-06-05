# Journal-canonical; DB derived; no daemon

## Decision

`.fid/journal/YYYY-MM-DD.jsonl` is committed to git and is the source of truth for all fid entity state. `.fid/config.toml` and `.fid/pricing.toml` are committed canonical configuration files, not journal-derived state. `.fid/fid.sqlite` is gitignored and materialized synchronously from the journal plus the current schema. The journal append is the write commit point; applying the event to SQLite is a synchronous projection step that can be replayed if interrupted. There is no background daemon and no sync window.

## Context

fid needs a fast SQL query surface for agent reads — joins across 25+ tables, the `issues_enriched` view, named reports, telemetry rollups, role-tuned context packs. It also needs to merge cleanly across branches under heavy agent traffic, where every interaction can produce a write.

Two reference points framed the decision:

- **Beads** (Steve Yegge) — SQLite canonical, JSONL exported via background daemon. Gives full SQL but the daemon has a known failure mode of "syncing the wrong things at the wrong times" because there's an asynchronous window between the DB write and the JSONL export.
- **Beans** (hmans) — markdown files canonical, GraphQL query layer scraped from the file tree. No daemon, no DB, but no SQL either; doesn't scale to fid's append-only completions, multi-role AC verification, or contract-coverage queries.

This ADR takes Beads' query power and Beans' git-mergeability, and removes the daemon by inverting the write order.

## Considered options

- **(a) Commit DB + journal + export (the previous spec position).** Three representations of the same state in git. Binary diff noise on every agent write. Redundant.
- **(b) Journal canonical, DB derived, no daemon.** Chosen.
- **(c) DB canonical, journal exported via daemon (Beads-style).** Rejected — daemon is the specific thing Beads' users complain about, and we'd inherit that.
- **(d) Markdown files canonical, no DB (Beans-style).** Rejected — fid's 25+ tables, multi-role verification, and traceability queries don't fit a files-only model.

## Consequences

- The DB is a cache. `rm -rf .fid/fid.sqlite` is non-destructive; the next CLI invocation rebuilds it.
- `fid rebuild` ships in v0.1, not deferred to v0.2. The replay engine is foundational, not optional.
- `fid init` is the first journaled write. It writes a `project.initialized` event whose effects seed database state such as the default role profile. Config and pricing files are written before the event and remain canonical files.
- The write path is ordered: validate against the current DB, build a fully resolved event, append and flush the JSONL line, then apply that event to SQLite in a DB transaction. If projection fails after the journal append, the next invocation catches the DB up from the journal under the same checkout-scoped projection lock before answering. If the journal append fails, the DB is not modified.
- Events are self-contained pure data: one logical CLI write becomes one semantic event envelope with generic ordered row effects. All generated/external values (hash IDs, timestamps, snapshot pricing, gate output, token counts) are captured at write time. Replay upgrades event schemas, then applies effects verbatim. No "look up X at replay time" and no re-running command/business logic.
- A malformed trailing JSONL line is journal corruption, not a partial commit. Startup / `fid rebuild` stops and reports the corrupt tail so a human can truncate or repair it.
- Schema migrations require versioned events (`event_schema: <n>`) so the replay engine can upgrade old events into the current schema.
- Merges are textual JSONL line-merges, handled by `git merge` directly. The v0.2 `fid merge --journal` story collapses into ordinary git workflow.
- Clone-then-query has a one-time rebuild cost. Append-only linear scan; expected to be sub-second at single-project scale.
- The §3.1 multi-worktree restriction softens: separate worktrees writing on the same day produce mergeable JSONL conflicts instead of an unrecoverable binary blob.
