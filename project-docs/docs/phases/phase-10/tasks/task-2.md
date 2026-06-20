# Phase 10 - Task 2: In-code permissions matrix + Authorize seam + `fid policy show`

## Objective

Build the v0.1 permissions infrastructure: the in-code permissions matrix
(TRD §13), the `policy.Authorize` enforcement seam, and the read-only
`fid policy show` command. This task creates and tests the matrix and the seam
in isolation; Task 3 wires `Authorize` into the existing write commands.

The matrix is the single source of truth shared by `Authorize` (enforcement) and
`fid policy show` (display) — DEC-B. It lists only the actions for write commands
implemented as of Phase 10 (v0.1a–v0.1c); later v0.1d phases append their action
rows alongside their commands.

## Requirements

- PRD refs: PRD-ROLE-005 (`fid policy show`), PRD-ROLE-003 (denied combinations
  rejected with role/action-named error — seam returns the error here; CLI proof
  in Task 3)
- TRD refs: TRD-ARCH-010 (matrix is a v0.1 code constant; every CLI write action
  has a matrix entry; unknown roles hard-rejected; denied combinations raise
  errors), TRD-ARCH-012 (TOON marshalled in-tree — reuse existing format engine)
- Design refs: 06_DESIGN_PLAN_CLI §11 `fid policy show`
- Quality refs: 08_QUALITY_PLAN §4.4 (budgets); COV-001 (`internal/policy` ≥ 80%),
  COV-002 (`cmd/fid` ≥ 70%)
- Decisions: DEC-A (Authorize seam), DEC-B (single-source matrix, grows per
  phase), DEC-C (`fid policy show` is a read; no `--by`, no journal write)

## Validation Contract Coverage

- VAL-002: `fid policy show` renders the matrix in all four formats and includes
  every implemented v0.1 write action.
- VAL-003: every implemented write action has exactly one matrix entry;
  `Authorize` permits allowed and returns `PermissionDeniedError` for denied.

## Scope

### Files To Create Or Modify

- `internal/policy/matrix.go` (new) — the in-code matrix as a data structure
  keyed by action with a per-role permission boolean; the v0.1a–c action set
  (DEC-B); a lookup/accessor; and an exported view shape for `fid policy show`.
- `internal/policy/authorize.go` (new) — `Authorize(workspaceRoot, role,
  action string) error`: runs the existing known-role validation
  (`MissingRoleError`/`UnknownRoleError`) then the matrix lookup, returning a new
  `PermissionDeniedError{Role, Action}` (message names role, action, and "not
  permitted", per PRD-ROLE-003). Reuse `RequireKnownRole`/`KnownRoles`.
- `cmd/fid/policy.go` (new) — `newPolicyCommand` → `policy show` subcommand;
  reads the matrix, writes via the format engine (`internal/format`) in
  json/toon/markdown/yaml; `--format` flag, default json; no `--by`.
- `cmd/fid/root.go` — register `newPolicyCommand`.
- Tests: `internal/policy/matrix_test.go`, `internal/policy/authorize_test.go`,
  `cmd/fid/policy_test.go`.

### Files To Read For Context

- `project-docs/specs/03_TRD_SCHEMA.md` §13 (the default permissions matrix — the
  authoritative cell values; copy the v0.1a–c rows exactly)
- `internal/policy/policy.go` (`RequireKnownRole`, `KnownRoles`, error types)
- `internal/format/` (the format engine + how existing list commands dispatch
  json/toon/markdown/yaml) and an existing read command (e.g.
  `cmd/fid/issue.go` `newIssueListCommand`) for the format/output pattern
- `06_DESIGN_PLAN_CLI.md` §11 `fid policy show` JSON example (row shape)

### Out Of Scope

- Wiring `Authorize` into existing write commands and updating their fixtures
  (Task 3).
- Claim/link/label/requirement.sync action rows (added in Phases 11–13).
- The `role_permissions` table, `fid role create`, `fid policy set` (v0.2).
- Any journal write or `activity_log` row for `fid policy show`.

## Acceptance Criteria

- [ ] The in-code matrix contains exactly the implemented v0.1a–c write actions
  (DEC-B list) with TRD §13 cell values (VAL-003).
- [ ] A table-driven test asserts every implemented action key resolves to one
  matrix entry (TRD-ARCH-010 completeness) (VAL-003).
- [ ] `Authorize` returns nil for a permitted (role, action), a
  `PermissionDeniedError` naming role + action for a denied one, and preserves
  `MissingRoleError`/`UnknownRoleError` for empty/unknown roles (VAL-003).
- [ ] `fid policy show` renders one row per action with a per-role boolean in
  json (default), toon, markdown, and yaml, with the same data in each, and
  includes every implemented write action (VAL-002).
- [ ] `fid policy show` takes no `--by`, writes no journal line, and exits 0
  (VAL-002, DEC-C).

## TDD Plan

- First failing test: `Authorize` permits a known allowed combination and returns
  `PermissionDeniedError` for a known denied combination (red — function does not
  exist yet).
- Additional tests: matrix completeness over the implemented action set; missing
  and unknown role still surface their existing errors via `Authorize`;
  `fid policy show` per-format rendering + action presence; show emits no journal
  line.
- Regression risks: format-engine reuse must match existing list-command output
  conventions; ensure `policy show` registration does not collide with other root
  subcommands. No existing write command behaviour changes in this task.

## Documentation Updates

- `06_DESIGN_PLAN_CLI.md` §11 — confirm the `fid policy show` JSON example matches
  the implemented row shape; adjust only if the implementation must diverge
  (record any divergence as a decision).
- `project-docs/docs/phases/phase-10/decisions.md` - note DEC-A/B/C as applied
  (seam name, matrix location, show-is-read).

## Scenario / Parity / Audit Notes

- Held-out scenario impact: None.
- Cross-surface parity impact: the matrix `fid policy show` displays must be the
  same object `Authorize` queries (DEC-B). Full show-vs-enforce parity proof is
  VAL-008 in Task 3 (needs the wired commands); this task must not create a second
  copy of the matrix.
- Audit or release-readiness impact: TRD-ARCH-010 completeness test is the audit
  hook ("every CLI write action has a matrix entry").

## Exact Commands To Run

Run exactly:

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

The full task gate above is required before task close. A focused test command
may be used during the red/green loop, but it does not substitute for the full
gate at close. Every required command must be recorded in the completion report
as `PASS`, or as `BLOCKED` with reason in a blocker report.

Conditional commands:

```text
go test -tags integration -race ./...
go build -o /tmp/fid ./cmd/fid
```

- `go build -o /tmp/fid ./cmd/fid` — required: a new root subcommand is added.
- `go test -tags integration -race ./...` — required if `policy_test.go` uses an
  integration-tagged workspace harness; otherwise record N/A with reason.

## Notes For Developer

- Copy the matrix cell values **verbatim** from `03_TRD_SCHEMA.md` §13 for the
  v0.1a–c action rows only (DEC-B). Do not seed v0.2+ rows or future v0.1d rows.
- The `human` column is `✓` for every write action in §13; encode it as data
  (do not special-case `human` in `Authorize`). The AC-verify role-keying
  (`issue.ac.verify.<role>`) and transition state-keying
  (`issue.transition.<state>`) are just distinct action keys — Task 3 derives the
  key; the matrix only stores rows.
- Keep `Authorize` thin: known-role check (reuse existing code), then matrix
  lookup. Functions/files stay within §4.4 budgets — split matrix data from
  lookup logic if needed.
- `fid policy show` is read-only; reuse the existing `internal/format` engine so
  toon/markdown/yaml stay consistent with other read commands (no bespoke
  marshalling).
