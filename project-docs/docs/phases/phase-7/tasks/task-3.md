# Phase 7 - Task 3: `fid issue ac list`

## Objective

Implement `fid issue ac list <issue-id> [--format json|toon|markdown|yaml]`:
list an issue's acceptance criteria with the latest verification status per
verifier role, dispatched through the existing format engine. Read-only.

## Requirements

- PRD refs: PRD-AC-006 (list ACs with latest verification status per verifier role)
- TRD refs: TRD_SCHEMA §2.3 (latest-per-role query semantics)
- Design refs: `06_DESIGN_PLAN_CLI.md` §3 `fid issue ac list` (JSON shape); PRD-OUT-001–004 format surface
- Quality refs: DEC-007 coverage

## Validation Contract Coverage

- VAL-006: `ac list` returns each AC with `latest_verifications` keyed by role (latest only, not raw history).
- VAL-007: `--format json|toon|markdown|yaml` renders without error; JSON is default.

## Scope

### Files To Create Or Modify

- `internal/issues/ac.go` (or `ac_list.go`) - `ListAcceptanceCriteria` read function returning the design §3 shape.
- `internal/issues/ac_test.go` - read-path tests.
- `cmd/fid/issue.go` - wire `ac list` subcommand + `--format` dispatch.
- `cmd/fid/issue_ac_test.go` - CLI/format tests.

### Files To Read For Context

- `internal/issues/detail.go` (`issueAcceptanceCriteria`, `latestVerificationsByRole` — reuse, do not duplicate).
- `internal/format/format.go` and `cmd/fid/ready.go` / `cmd/fid/issue.go` (existing `--format` dispatch pattern).

### Out Of Scope

- Adding/verifying ACs (Tasks 1–2); transition gating (Task 4).
- Full per-role verification history listing (only latest-per-role is in scope).
- `completion_id` surfacing (Phase 8).

## Acceptance Criteria

- [ ] `fid issue ac list <issue>` returns ACs in ordinal order, each with `id`, `ordinal`, `text`, and `latest_verifications` keyed by verifier role, matching the design §3 JSON shape (VAL-006).
- [ ] After multi-role and repeated verifications, each verifying role appears once with its most recent status; superseded statuses are not listed as separate entries (VAL-006).
- [ ] `--format json|toon|markdown|yaml` each render without error; omitting `--format` yields JSON (VAL-007).
- [ ] Issue-not-found exits non-zero with a clear error.
- [ ] An issue with no ACs returns an empty list (not an error).

## TDD Plan

- First failing test: `ListAcceptanceCriteria` returns the design §3 shape with latest-per-role for a populated issue.
- Additional tests: empty-AC issue → empty list; repeated verification shows latest only; all four formats render; default = JSON; issue-not-found.
- Regression risks: reuse the existing read helpers so `ac list`, `issue get --section acceptance`, and the rule loader stay consistent.

## Documentation Updates

- `06_DESIGN_PLAN_CLI.md` §3 - confirm the `ac list` JSON example matches the implementation output.
- Completion report - paste exact transcripts for one populated `ac list` per format and the empty-list case.

## Scenario / Parity / Audit Notes

- Held-out scenario impact: None (AC theme only).
- Cross-surface parity impact: `ac list` latest-per-role must match `issue get --section acceptance`.
- Audit or release-readiness impact: None.

## Exact Commands To Run

Run exactly:

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

Coverage targets (DEC-007): `internal/issues` ≥ 80%, `cmd/fid` ≥ 70%.

Conditional commands:

```text
go build -o /tmp/fid ./cmd/fid
```

List any conditional command required by this task and why:

- `go build -o /tmp/fid ./cmd/fid` if pasting real CLI transcripts for evidence.

## Notes For Developer

- This is a read command: no journal write, no lock. Use `journal.EnsureProjectionFresh` then open the DB read-only (mirror `Get`).
- Reuse `issueAcceptanceCriteria` / `latestVerificationsByRole`; do not re-implement the latest-per-role query.
- Route output through the shared format engine exactly as `fid issue list` / `fid ready` do; JSON is the default.
</content>
