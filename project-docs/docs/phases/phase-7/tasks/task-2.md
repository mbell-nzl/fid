# Phase 7 - Task 2: `fid issue ac verify`

## Objective

Implement `fid issue ac verify <issue-id> <ac-id> --by <role> --status <status>
[--evidence-type --evidence-ref --notes]`: append-only, multi-role verification
recorded as an `ACV-<hash>` row in `ac_verifications`, journaled as one
`issue.ac_verified` event. Preserves full history; latest-per-role reflects the
most recent verification.

## Requirements

- PRD refs: PRD-AC-003 (multi-role independent outcomes), PRD-AC-004 (records role/status/evidence/timestamp), PRD-AC-005 (append-only history)
- TRD refs: TRD_SCHEMA §2.3 `ac_verifications` (append-only; latest-per-role); §1 ID scheme (`ACV-<hash>`)
- Design refs: `06_DESIGN_PLAN_CLI.md` §3 `fid issue ac verify`
- Quality refs: validate-before-append (CF-007); DEC-007 coverage

## Validation Contract Coverage

- VAL-003: Two roles record independent verifications on the same AC; both latest statuses retained.
- VAL-004: Append-only — fail then pass from the same role yields two rows; latest = pass.
- VAL-005: Records verifier_role/status/evidence_type/evidence_ref/notes/verified_at; invalid `--status` rejected before append.

## Scope

### Files To Create Or Modify

- `internal/issues/ac.go` - add `VerifyAcceptanceCriterion` write path.
- `internal/issues/ac_test.go` - verification unit/integration tests.
- `cmd/fid/issue.go` - wire `ac verify` subcommand + stderr confirmation.
- `cmd/fid/issue_ac_test.go` - CLI tests.

### Files To Read For Context

- `internal/issues/detail.go` (`latestVerificationsByRole`, `issueAcceptanceCriteria`).
- `internal/issues/transition_rules.go` (`loadAcceptanceCriterionVerifications` — the read consumer of these rows).
- `internal/issues/issue.go` (write-path pattern).

### Out Of Scope

- `ac list` command (Task 3); transition gating (Task 4).
- `--in-completion` / `ac_verifications.completion_id` linkage (Phase 8); leave `completion_id` NULL.
- Permission matrix enforcement of which roles may verify (v0.1d); only `RequireKnownRole`.
- `activity_log` emission (DEC-020, confirmed; deferred to v0.1e via CF-011).

## Acceptance Criteria

- [ ] `ac verify <issue> <ac> --by developer --status pass` appends one `ACV-<hash>` row and exactly one `issue.ac_verified` event (VAL-005).
- [ ] All recorded fields persist: verifier_role, status, evidence_type, evidence_ref, notes, verified_at (VAL-005).
- [ ] Verifying the same AC as `developer` and `qa_lead` retains both latest statuses independently (VAL-003).
- [ ] Same AC + same role `fail` then `pass` → two retained rows; latest-per-role = pass; after fail-only latest = fail (VAL-004).
- [ ] `--status` is required and must be one of `pass|fail|blocked|not_applicable`; an invalid value is rejected before append (zero new rows/lines) (VAL-005).
- [ ] `--evidence-type` (when given) must be one of `test|source|screenshot|manual|command|commit`; invalid value rejected before append.
- [ ] `<ac-id>` must exist and belong to `<issue-id>`; mismatch/not-found rejected with clear error and zero append.
- [ ] stderr confirmation matches design §3 (status line + optional `evidence:` line).

## TDD Plan

- First failing test: a pass verification projects one `ac_verifications` row with all fields set and an `ACV-` id.
- Additional tests: multi-role independence; fail→pass two-row history + latest-per-role; invalid status; invalid evidence-type; ac/issue mismatch; ac not found; one-event-per-write.
- Regression risks: `latestVerificationsByRole` ordering (verified_at DESC, id DESC) must keep selecting the newest row; do not change the transition-rule loader contract.

## Documentation Updates

- `06_DESIGN_PLAN_CLI.md` §3 - confirm `ac verify` flags/output match implementation; adjust only on divergence.
- Completion report - paste exact transcripts for: pass verify, multi-role, fail→pass history, invalid-status rejection.

## Scenario / Parity / Audit Notes

- Held-out scenario impact: None (AC theme only).
- Cross-surface parity impact: latest-per-role read by `issue get --section acceptance`, `ac list` (Task 3), and the `all_acs_verified` rule loader must agree.
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

- Append-only means **insert** only: never update/overwrite a prior verification row. Two verifications = two rows.
- Generate ids with `journal.GenerateID("ACV", existsFn)`. Event type: `issue.ac_verified`.
- Validate `--status` and `--evidence-type` enums and AC/issue membership in the engine **before** building/appending the event (CF-007).
- `--by human` records a human verification stored verbatim in `verifier_role`; it does not count as developer/QA/TL for rule purposes (design §3).
</content>
