# Phase 10 - Task 1: Reconcile role-discovery corrupt-tail handling (CF-003)

## Objective

Reconcile the role-discovery journal scan in `internal/policy` with the journal
corruption contract before the permission matrix builds on it. Today
`collectRoles` silently tolerates a malformed/un-upgradeable journal line once
seeded roles have been found (Phase 1 debt, TRG-002 / CF-003). Now that role
discovery is about to feed permission enforcement, a corrupt journal must be
**reported, not skipped** — matching the replayer (TRD-ARCH-002) and the
journal-canonical ADR-0001.

This is a self-contained, low-blast-radius change (one package) that runs first
so the matrix work in Tasks 2–3 sits on reconciled role handling.

## Requirements

- PRD refs: PRD-ROLE-002 (role validation reads a sound role set)
- TRD refs: TRD-ARCH-002 (malformed line = corruption, reported not skipped)
- Design refs: N/A
- Quality refs: 08_QUALITY_PLAN §4.4 (file/function budgets); COV-001 (`internal/policy` ≥ 80%)
- Carry-forward: CF-003 (role-validation corrupt-tail tolerance)
- Decision: DEC-D (this phase) — align role discovery with the corruption contract

## Validation Contract Coverage

- VAL-001: role discovery reports a malformed/un-upgradeable journal line as
  corruption (non-zero exit, error naming the file) instead of tolerating it
  after seeded roles are found; a clean journal still resolves the role set.

## Scope

### Files To Create Or Modify

- `internal/policy/policy.go` — remove the "tolerate after roles found" branch in
  `collectRoles`; a malformed line and an un-upgradeable event both return the
  same corruption error they return for a malformed first line, naming the file.
- `internal/policy/policy_test.go` — rewrite the
  `TestKnownRolesTRG002CorruptTailTolerance` "tolerated" subtest (decision-backed,
  red/green) to assert reported corruption; keep the "malformed first line
  surfaces an error" subtest; add a positive case proving a clean multi-line
  role journal still resolves.

### Files To Read For Context

- `internal/policy/policy.go` (`collectRoles`, `KnownRoles`)
- `internal/journal/` corruption/upgrade error surface (the error type/shape the
  replayer uses for a malformed line) so role discovery surfaces a consistent
  error
- `project-docs/docs/phases/phase-1/` TRG-002 note (origin of the tolerance)

### Out Of Scope

- The permission matrix, `Authorize`, and `fid policy show` (Tasks 2–3).
- Any change to the replayer/rebuild corruption behaviour itself (already
  reports); this only reconciles the lighter role-discovery scan.
- Changing the malformed-first-line behaviour (already correct).

## Acceptance Criteria

- [ ] A journal with a valid `roles` insert followed by a malformed trailing line
  causes `KnownRoles` to return a corruption error naming the journal file
  (VAL-001).
- [ ] A journal whose first line is malformed still surfaces a corruption error
  naming the file (unchanged) (VAL-001).
- [ ] A clean multi-file/multi-line role journal resolves the active role set
  with no error (VAL-001).
- [ ] A CLI write command run against a corrupt-tail workspace exits non-zero
  with a corruption error naming the file (VAL-001).
- [ ] The error returned for a corrupt role-discovery scan is consistent in
  shape with the replayer's corruption error (no bespoke opaque error).

## TDD Plan

- First failing test: rewrite `TestKnownRolesTRG002CorruptTailTolerance` so the
  "valid role before malformed trailing line" subtest expects a corruption error
  (currently expects tolerance) — red against today's `collectRoles`.
- Additional tests: clean multi-line journal resolves roles; un-upgradeable
  (bad-schema) trailing event also reports corruption; error names the file path.
- Regression risks: `RequireKnownRole`/`KnownRoles` are called by every existing
  write command and `rebuild`; ensure clean journals (the normal path) are
  unaffected. Run the full suite to confirm no existing command regresses on a
  well-formed journal.

## Documentation Updates

- `project-docs/docs/phases/phase-10/decisions.md` - record DEC-D rationale and
  the intentional `TestKnownRolesTRG002CorruptTailTolerance` rewrite (so the
  behaviour change is decision-backed, not a silent test edit).
- No spec doc change expected (behaviour is being aligned *to* the existing
  TRD-ARCH-002 contract); confirm `02_TRD.md`/`03_TRD_SCHEMA.md` need no edit and
  note `verified-no-change` in the completion report.

## Scenario / Parity / Audit Notes

- Held-out scenario impact: None.
- Cross-surface parity impact: role discovery now matches replayer/rebuild
  corruption behaviour — a consistency (parity) improvement across the two
  journal-reading paths.
- Audit or release-readiness impact: closes CF-003; note the disposition in the
  completion report for the wrap audit.

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

- `go test -tags integration -race ./...` — required: a CLI write against a
  corrupt-tail workspace is an integration-level check of the reconciled error.
- `go build -o /tmp/fid ./cmd/fid` — required: confirm the binary still builds.

## Notes For Developer

- This is a deliberate, decision-backed behaviour change (DEC-D). The existing
  "tolerated" subtest encoded acknowledged Phase 1 debt; rewriting it to expect
  reported corruption is a correction, **not** weakening a behaviour test. Do not
  delete the test — change its expectation with a matching red/green cycle and
  cite DEC-D in the diff context.
- Keep the corruption error consistent with the journal package's existing
  corruption/upgrade error so `cmd/fid` error rendering and any callers stay
  uniform; do not invent a new opaque error string.
- `collectRoles` and `KnownRoles` must remain within the §4.4 function/file size
  budgets after the change (the change should shrink, not grow, the logic).
