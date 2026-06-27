# Phase 12 - Task 1: Requirement Pointers And `fid requirement sync`

## Objective

Create the requirements domain layer and the minimal `fid requirement sync`
command so issue -> requirement links have valid endpoints. Sync parses
canonical PRD/TRD IDs from a source document and upserts pointer rows
(`id`, `type`, `summary`, `source_doc`, `source_section`, `last_synced_hash`,
`status`) through one journal event, reporting created/updated/unchanged counts.
Append the `requirement.sync` row to the in-code permissions matrix.

## Requirements

- PRD refs: PRD-GOV-006, PRD-LINK-002, PRD-ERR-003
- TRD refs: TRD-ARCH-010; `03_TRD_SCHEMA.md` sections 5.2, 13
- Design refs: `06_DESIGN_PLAN_CLI.md` `fid requirement sync`
- Quality refs: `08_QUALITY_PLAN.md` COV-001/COV-002; CF-007 validate-before-append
- Decisions: DEC-028 (held-lock projection freshness); planning interpretation
  PI-1 (requirement.sync parse contract) in the validation contract

## Validation Contract Coverage

Assertions live in `project-docs/docs/phases/phase-12/qa/validation-contract.md`.

- VAL-001 (partial): the `requirement.sync` matrix row exists and matches
  `03_TRD_SCHEMA.md` section 13 (tl + human); `fid policy show` reflects it.
  (Link rows are added in Task 2.)
- VAL-002: sync upserts pointer rows in one event with correct fields and
  created/updated/unchanged counts; unchanged sync appends no event; a changed
  summary updates the row and `last_synced_hash`.
- VAL-003: malformed input (no parseable IDs, duplicate IDs, unknown `--type`,
  missing `--id`) rejects before append with zero new journal lines.
- VAL-004: `requirement.sync` authorized only for tl/human; denied/unknown role
  rejects before append naming role + action.

## Scope

### Files To Create Or Modify

- `internal/requirements/` (new package) - parser, hash, upsert service, types,
  errors, and tests for `fid requirement sync`.
- `internal/policy/matrix.go` and `internal/policy/matrix_test.go` - append the
  `requirement.sync` row from `03_TRD_SCHEMA.md` section 13 and update
  implemented-action completeness tests.
- `cmd/fid/requirement.go` (new) and `cmd/fid/requirement_test.go` - wire
  `fid requirement` with the `sync` subcommand and flags.
- `cmd/fid/root.go` - register `newRequirementCommand()`.

### Files To Read For Context

- `internal/claims/acquire.go` and `internal/claims/shared.go` for the
  lock -> ensure-fresh-with-lock -> build event -> dry-run -> append -> project
  write pattern and ID generation.
- `internal/journal/project_effects.go` for insert/update effect and
  precondition semantics.
- `internal/policy/matrix.go`, `internal/policy/authorize.go`.
- `cmd/fid/claim.go`, `cmd/fid/policy.go` for new-top-level-command wiring and
  `fidformat.Write` / `renderIssueError` usage.
- `migrations/0002_v01_tables.up.sql` (`requirements` DDL).
- `project-docs/specs/01_PRD.md` requirement tables as the canonical parse-format
  reference (do not copy requirement prose into fixtures beyond what a test
  needs).

### Out Of Scope

- `fid requirement get` / `fid requirement list` and drift reporting (v0.2).
- Any `fid link *` command (Tasks 2-4).
- Schema migrations. The `requirements` table already exists.

## Acceptance Criteria

- [ ] `internal/requirements` owns parse + upsert business rules; `cmd/fid`
  contains only flag wiring and output (VAL-002).
- [ ] The in-code matrix contains exactly one `requirement.sync` row matching
  section 13 (tl + human true, all other roles false); `fid policy show`
  includes it (VAL-001 partial).
- [ ] `fid requirement sync --source-doc <path> --by tl` parses canonical IDs
  (`(PRD|TRD)-[A-Z0-9]+-[0-9]+[a-z]?`) from markdown table rows, deriving
  `type` from the ID prefix and `summary` from the row's second cell, and
  upserts rows with `status='active'` and a deterministic `last_synced_hash`
  over the summary (VAL-002).
- [ ] A single sync emits exactly one journal event containing one
  insert/update effect per created/updated requirement; a sync with zero
  creates/updates appends no event (VAL-002).
- [ ] Output reports `created`, `updated`, and `unchanged` counts; re-syncing
  identical content reports unchanged; editing a summary reports updated and
  changes the stored hash (VAL-002).
- [ ] `--type PRD|TRD` filters to that type; `--id <requirement-id>` restricts
  the sync to one ID and errors if the ID is absent from the source (VAL-002,
  VAL-003).
- [ ] Malformed input rejects before append with non-zero exit and zero new
  journal lines: no parseable IDs, duplicate IDs within the source, unknown
  `--type` (VAL-003).
- [ ] `requirement.sync` authorizes via the matrix: tl/human succeed; a denied
  role and an unknown role get `PermissionDeniedError` before append naming
  role + action (VAL-004).

## TDD Plan

- First failing test: syncing a fixture doc with three requirement rows creates
  three `requirements` rows with the expected `id`/`type`/`summary`/`source_doc`/
  `status` and one journal event.
- Additional tests: re-sync unchanged (unchanged=3, no new event); edit one
  summary and re-sync (updated=1, hash changes); `--id` single-requirement sync;
  `--type` filter; no-parseable-IDs, duplicate-IDs, unknown-type errors;
  denied-role and unknown-role rejection with zero append; matrix completeness.
- Regression risks: writes must use journal effects + dry-run before append, not
  direct DB mutation; the unchanged path must append nothing; hash must be
  deterministic across runs; use `EnsureProjectionFreshWithLock` inside
  `WithLock` (no nested lock, DEC-028).

## Documentation Updates

- `project-docs/docs/phases/phase-12/decisions.md` - record DEC for PI-1 (the
  requirement.sync parse contract and exact hashed content), plus any
  section-13 vs in-code mismatch found.
- `06_DESIGN_PLAN_CLI.md` `fid requirement sync` - reconcile only if the
  implemented flags/output differ from the spec; otherwise record
  verified-no-change in the completion report.
- `project-docs/docs/phases/phase-12/tasks/task-1-done.md` - completion report
  with command results, coverage, parse-contract summary, and commit hash.

## Scenario / Parity / Audit Notes

- Held-out scenario impact: None.
- Cross-surface parity impact: policy show vs `Authorize` parity for the
  `requirement.sync` row begins here and is finalized with the link rows in
  Task 2.
- Audit or release-readiness impact: None.

## Exact Commands To Run

Run exactly:

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

The full task gate above is required before task close because this task
introduces a new domain package, a journaled write, and a new command. A focused
test command may be used during the red/green loop, but it does not substitute
for the full gate at close. Every required command must be recorded in the
completion report as `PASS`, or as `BLOCKED` with reason in a blocker report.

Conditional commands:

```text
go test -tags integration -race ./...
go build -o /tmp/fid ./cmd/fid
```

- `go build -o /tmp/fid ./cmd/fid` - required because this task registers a new
  CLI command; record PASS.
- `go test -tags integration -race ./...` - run if any integration fixture is
  touched; otherwise record N/A with reason.

## Notes For Developer

- Use red/green TDD. Do not add a schema migration unless a spec-backed blocker
  is recorded first.
- One logical write = one event: a sync that touches N requirements produces one
  journal event with N effects, not N events.
- Keep parse logic in `internal/requirements` so Task 2's link validation can
  reuse `RequirementExists`-style helpers if needed.
- The source doc is outside `.fid/` (e.g. a PRD path relative to CWD); read it
  as a plain file. Store pointers and summaries only; never copy full
  requirement prose into the row beyond the summary cell.
- Completion report commit fields must be filled from `git rev-parse HEAD`.
