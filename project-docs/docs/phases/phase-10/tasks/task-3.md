# Phase 10 - Task 3: Retrofit permission enforcement onto every write command

## Objective

Wire `policy.Authorize` (Task 2) into every existing write command so each write
is gated by the matrix on its derived action, not by `RequireKnownRole` alone.
This is the cross-cutting slice that makes Phase 10 real: after this task, a
denied (role, action) write is rejected before any journal append with a
role/action-named error, the `--by human` escape valve works without identity
laundering, and `fid policy show` provably matches enforced behaviour.

Enforcing the matrix changes the behaviour of existing commands (DEC-E), so this
task also moves existing test fixtures to permitted role/action combinations.

## Requirements

- PRD refs: PRD-ROLE-001 (`--by` required), PRD-ROLE-002 (unknown role rejected),
  PRD-ROLE-003 (denied combination rejected naming role + action), PRD-ROLE-004
  (human escape valve; no verifier impersonation)
- TRD refs: TRD-ARCH-010 (enforce matrix; every write action gated)
- Design refs: 06_DESIGN_PLAN_CLI §2/§3/§6 (the existing write commands), §11
  (policy show — referenced for parity)
- Quality refs: 08_QUALITY_PLAN §4.4 (budgets); COV-001 (`internal/issues|policy`
  ≥ 80%), COV-002 (`cmd/fid` ≥ 70%); CF-007 (validate-before-append)
- Decisions: DEC-A (action key derived in domain layer), DEC-B (single-source
  matrix), DEC-E (fixtures move to permitted combinations)

## Validation Contract Coverage

- VAL-004: denied (role, action) write rejected before any journal append
  (non-zero exit, role/action-named error, zero new lines); permitted role
  succeeds with one append.
- VAL-005: role-identity validation preserved through `Authorize` — missing
  `--by` and unknown role rejected before append; permitted role proceeds.
- VAL-006: `--by human` permitted where a non-human role is denied; human AC
  verify attributed to `issue.ac.verify.human` (`_role = human`), no
  impersonation.
- VAL-007: action-key derivation — `issue.ac.verify.<role>` (qa_reviewer
  permitted, workflow_runner denied) and `issue.transition.<target-state>`
  (developer→dev_complete permitted, qa_reviewer denied).
- VAL-008: `fid policy show` display matches enforced behaviour for a permitted
  and a denied cell (cross-surface parity).

## Scope

### Files To Create Or Modify

- `internal/issues/create.go` — derive `issue.create.<kind>`; replace
  `RequireKnownRole` with `Authorize`. (Owner-role validation at line ~291 stays
  a known-role check, not a permission action.)
- `internal/issues/update.go` / `update_request.go` — `issue.update`.
- `internal/issues/transition.go` — derive `issue.transition.<target-state>`.
- `internal/issues/complete.go` — `issue.complete`.
- `internal/issues/ac.go` — `issue.ac.add`.
- `internal/issues/ac_verify.go` — derive `issue.ac.verify.<actor role>`.
- `internal/issues/concern_raise.go` — `concern.raise`.
- `internal/issues/concern_address.go` — `concern.address`.
- `cmd/fid/rebuild.go` — `maintenance.rebuild`.
- `cmd/fid/issue.go` (`renderIssueError`) and any other error renderers — add a
  `PermissionDeniedError` case (role/action-named, actionable message).
- Existing `*_test.go` across `internal/issues` and `cmd/fid` — update fixtures
  that used now-denied combinations to permitted roles per DEC-E; add the new
  enforcement/parity tests.

### Files To Read For Context

- `internal/policy/matrix.go`, `internal/policy/authorize.go` (Task 2)
- `03_TRD_SCHEMA.md` §13 (which role is permitted for each action — drives the
  fixture role choices)
- Each write entry point listed above (current `RequireKnownRole` call sites)
- `cmd/fid/issue.go` `renderIssueError` (error-rendering pattern)

### Out Of Scope

- New v0.1d write commands (claim/link/label/requirement.sync) and their action
  rows — Phases 11–13.
- `fid init` (bootstrap; exempt — no role exists yet to check).
- Any schema migration or `activity_log` change.

## Acceptance Criteria

- [ ] Every existing write entry point calls `Authorize` with its derived action
  key; no write retains a bare `RequireKnownRole` where a matrix action applies
  (VAL-004, VAL-007).
- [ ] A denied write (e.g. `issue create --kind epic --by developer`) exits
  non-zero with an error naming the role, the action, and "not permitted", and
  appends zero journal lines/rows; the permitted case (`--by tl`) appends exactly
  one event (VAL-004).
- [ ] Missing `--by` → `--by is required`; unknown role → error listing known
  roles; both before any append (VAL-005).
- [ ] `--by human` succeeds on a write a non-human role is denied; a human AC
  verification reads back with `_role = human`, not another verifier (VAL-006).
- [ ] `issue ac verify --by qa_reviewer` permitted, `--by workflow_runner`
  denied; `issue transition --to dev_complete --by developer` permitted,
  `--by qa_reviewer` denied (VAL-007).
- [ ] For one permitted and one denied cell, `fid policy show`'s boolean matches
  the actual write outcome (VAL-008).
- [ ] `PermissionDeniedError` renders as an actionable stderr line (and a JSON
  error when `--format json` is active on the command).

## TDD Plan

- First failing test: a CLI write with a denied (role, action) is rejected with a
  `PermissionDeniedError` and zero new journal lines (red — command currently
  permits any known role).
- Additional tests: per-command permitted/denied pairs; missing/unknown role
  preserved; human escape valve + human AC-verify attribution; AC-verify
  role-keying and transition state-keying; show-vs-enforce parity.
- Regression risks: **high blast radius.** Many existing tests assume any known
  role can drive any write; those must move to permitted roles (DEC-E) without
  changing what behaviour they assert. Re-run the full suite plus integration to
  confirm no behaviour test is silently weakened. Validate-before-append must hold
  for every new reject path (zero journal lines on denial — assert explicitly).

## Documentation Updates

- `project-docs/docs/phases/phase-10/decisions.md` - record DEC-E with the list of
  fixture role changes made (so the behaviour change is decision-backed and
  auditable at wrap).
- `06_DESIGN_PLAN_CLI.md` §2/§3/§6 — add a note (if not already implied) that each
  write enforces the permission matrix and names the role/action on denial;
  confirm no flag-surface change. Record `verified-no-change` if no edit needed.
- `01_PRD.md` — confirm PRD-ROLE-001..004 are MET as written; note in completion
  report (no requirement reword expected).

## Scenario / Parity / Audit Notes

- Held-out scenario impact: None.
- Cross-surface parity impact: VAL-008 show-vs-enforce parity (the Phase 9
  standing cross-surface parity check, instantiated for the matrix). Assert the
  displayed boolean equals the enforced outcome for both a permitted and a denied
  cell.
- Audit or release-readiness impact: this completes TRD-ARCH-010 enforcement;
  the completion report should list every write action now gated for the wrap
  audit.

## Exact Commands To Run

Run exactly:

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
go test -tags integration -race ./...
go build -o /tmp/fid ./cmd/fid
```

The full gate above (including integration and build) is required before task
close because this task changes behaviour across the whole write surface and is
exercised by integration-level CLI tests. A focused test command may be used
during the red/green loop, but it does not substitute for the full gate at close.
Every required command must be recorded in the completion report as `PASS`, or as
`BLOCKED` with reason in a blocker report.

Conditional commands:

```text
go test -bench . -benchmem ./...
```

- `go test -bench . -benchmem ./...` — required if any benchmark touches a write
  path whose fixture role changed (Phase 9 ESC-006 showed benchmark gates can
  surface fixture/ordering regressions); otherwise record N/A with reason.

## Notes For Developer

- Derive the action key in the **domain layer** from already-known inputs
  (DEC-A): kind for create, target state for transition, actor role for ac verify,
  fixed strings otherwise. Do not push action-key logic into `cmd/fid`.
- Moving a fixture from a denied role to a permitted one (DEC-E) is **not** test
  weakening — the test still asserts the same behaviour, just with a role the
  matrix permits. Pick the smallest permitted role that keeps the test's intent
  (use `tl`/`human` for create/transition-to-done; keep `developer` for
  dev-owned actions). List every fixture role change in the completion report.
- Assert **zero new journal lines** on every permission denial (CF-007), the same
  way Phase 9 asserted validate-before-append for concern rejects.
- The `human` escape valve is matrix data, not a code branch (Task 2 / DEC-B): if
  a human write is denied, the bug is a missing/incorrect matrix cell, not a
  missing special case.
- Keep modified functions within §4.4 budgets; if a write entry point grows past
  60 lines after adding the action-key derivation, extract a small helper.
