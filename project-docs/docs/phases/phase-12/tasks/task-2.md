# Phase 12 - Task 2: Links Domain, Matrix Rows, And `fid link add`

## Objective

Create the links domain layer and the `fid link add` command. Add a typed link
plus its auto-created inverse in one journal event, validate relation/endpoint
compatibility and dedicated-schema rules before append, and append the
`link.add` and `link.remove` rows to the in-code permissions matrix.

## Requirements

- PRD refs: PRD-LINK-001, PRD-LINK-002, PRD-LINK-003, PRD-LINK-004,
  PRD-LINK-005, PRD-LINK-008
- TRD refs: TRD-ARCH-010; `03_TRD_SCHEMA.md` sections 10, 13, 14.2-14.3
- Design refs: `06_DESIGN_PLAN_CLI.md` `fid link add`; GLOSSARY "Link"
- Quality refs: `08_QUALITY_PLAN.md` COV-001/COV-002; CF-007 validate-before-append
- Decisions: DEC-028 (held-lock projection freshness); planning interpretation
  PI-2 (PRD-LINK-008 v0.1 enforcement) in the validation contract
- Depends on: Task 1 (requirement endpoints must exist for issue -> requirement
  links)

## Validation Contract Coverage

Assertions live in `project-docs/docs/phases/phase-12/qa/validation-contract.md`.

- VAL-001 (completes): the `link.add` and `link.remove` matrix rows exist and
  match section 13 (all seven roles); `fid policy show` reflects them.
- VAL-005: valid issue <-> issue link add creates declared + inverse in one
  event with confirmation naming both directions; duplicate add rejects before
  append.
- VAL-006: issue -> requirement `implements` link works and yields the
  `implemented_by` inverse; a link to a missing requirement rejects.
- VAL-007: relation/endpoint-incompatible and unresolvable-endpoint adds reject
  before append naming relation + endpoints.
- VAL-008: relation registry equals the documented seven relations + inverses;
  unknown relation and self-link reject before append.

## Scope

### Files To Create Or Modify

- `internal/links/` (new package) - relation/inverse registry, validation
  (relation/endpoint compatibility, unknown relation, self-link), add service
  building a two-insert (declared + inverse) journal event, types, errors, tests.
- `internal/policy/matrix.go` and `internal/policy/matrix_test.go` - append the
  `link.add` and `link.remove` rows from section 13; update completeness tests.
- `cmd/fid/link.go` (new) and `cmd/fid/link_test.go` - wire `fid link` with the
  `add` subcommand and `--by`; stderr confirmation showing the inverse.
- `cmd/fid/root.go` - register `newLinkCommand()`.

### Files To Read For Context

- `internal/claims/acquire.go` for the lock -> ensure-fresh-with-lock -> build
  event (multiple effects) -> dry-run -> append -> project pattern.
- `internal/journal/project_effects.go` for insert-effect semantics.
- `migrations/0005_v01_ref_triggers.up.sql` `links_*` triggers (mirror their
  rules in the app layer for before-append errors).
- `internal/requirements/` (Task 1) for a requirement-existence helper.
- `cmd/fid/claim.go`, `cmd/fid/concern.go` for command wiring and error render.

### Out Of Scope

- `fid link remove` (Task 3) and `fid link list` (Task 4).
- Decision/escalation endpoints (v0.2, PRD-LINK-007).
- Projecting link add to `activity_log` (link add does not; only remove does,
  per PRD-LINK-006a and the CF-011 precedent).
- Schema migrations. The `links` table and triggers already exist.

## Acceptance Criteria

- [ ] `internal/links` owns the relation/inverse registry and all validation;
  `cmd/fid` contains only flag wiring and confirmation output (VAL-005..008).
- [ ] The in-code matrix contains exactly one `link.add` and one `link.remove`
  row, both true for all seven roles, matching section 13; `fid policy show`
  includes them (VAL-001).
- [ ] The relation/inverse registry equals exactly: blocks/blocked_by,
  relates_to/relates_to, duplicates/duplicated_by, supersedes/superseded_by,
  implements/implemented_by, caused_by/causes, discovered_from/discovered, and
  contains no hierarchy/escalation/resolution relation (VAL-008).
- [ ] `fid link add <from> <relation> <to> --by <role>` for a valid issue <->
  issue relation creates the declared row and its inverse in one journal event,
  projects both `links` rows, and prints a confirmation naming both directions
  (VAL-005).
- [ ] An issue -> requirement `implements` link (requirement synced in Task 1)
  creates the link and its `implemented_by` inverse; an `implements` link to a
  missing requirement ID rejects before append (VAL-006).
- [ ] Relation/endpoint-incompatible adds (e.g. `implements` issue -> issue,
  `blocks` issue -> requirement) and unresolvable-endpoint adds reject before
  append with non-zero exit, zero new journal lines, and an error naming the
  relation and both endpoints (VAL-007).
- [ ] An unknown relation and a self-link (`from == to`) reject before append;
  re-adding an existing link rejects before append naming the existing
  relation/entities (VAL-005, VAL-008).
- [ ] `link.add` authorizes via the matrix; an unknown role is rejected before
  append (VAL-001).

## TDD Plan

- First failing test: `Add` for `blocks` between two issues creates two `links`
  rows (declared + `blocked_by` inverse) and one journal event.
- Additional tests: `relates_to` (self-inverse) issue <-> issue and issue <->
  requirement; `implements` issue -> requirement yields `implemented_by`;
  duplicate add rejection; incompatible relation/endpoint cases; missing
  endpoint; unknown relation; self-link; registry equality; matrix completeness;
  unknown-role rejection. All rejection tests assert zero new journal lines.
- Regression risks: build both rows as effects of one event (not two events);
  validate before append (dry-run project); mirror the relation-compatibility
  trigger so the app error is actionable rather than a raw SQL abort; use
  `EnsureProjectionFreshWithLock` inside `WithLock` (DEC-028).

## Documentation Updates

- `project-docs/docs/phases/phase-12/decisions.md` - record DEC for PI-2 (v0.1
  PRD-LINK-008 structural enforcement: closed relation set + endpoint scope +
  self-link rejection), plus any section-13 vs in-code mismatch.
- `06_DESIGN_PLAN_CLI.md` `fid link add` - reconcile only if implemented output
  differs from the spec; otherwise record verified-no-change.
- `project-docs/docs/phases/phase-12/tasks/task-2-done.md` - completion report
  with command results, coverage, and commit hash.

## Scenario / Parity / Audit Notes

- Held-out scenario impact: None.
- Cross-surface parity impact: completes policy show vs `Authorize` parity for
  the Phase 12 matrix rows (VAL-001).
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
for the full gate at close.

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
- The user declares one direction; the inverse is auto-created. Support both
  directions of each relation as a valid declared input (e.g. a declared
  `blocked_by` produces a `blocks` inverse).
- `relates_to` is its own inverse and is the only relation allowed across mixed
  endpoint types in v0.1.
- Keep `Add` and the validation helpers structured so Task 3 (remove) and Task 4
  (list) can reuse the registry and endpoint-resolution helpers.
- Completion report commit fields must be filled from `git rev-parse HEAD`.
