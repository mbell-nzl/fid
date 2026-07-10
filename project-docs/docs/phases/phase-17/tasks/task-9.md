# Phase 17 - Task 9: Add the Full v0.1 Closure Journey

## Objective

Add the final integration-tagged v0.1 journey over a real temp workspace. Drive
the primitive lifecycle, scheduling, labels, links, telemetry, maintenance,
version guard, git layout, and delete-and-rebuild parity as one release-closure
proof; add no new product behavior.

## Validation Contract Coverage

- VAL-018: primitive lifecycle, permissions, ACs, claims, and concerns.
- VAL-019: sequential and prioritised ready-queue behavior.
- VAL-020: labels, requirement links, inverse links, and relation rejection.
- VAL-021: priced run telemetry and aggregate cost/token read-back.
- VAL-022: maintenance plus delete-and-rebuild structural equality.
- VAL-023: journal trackability/immutability, DB ignore, and seeded-role rebuild.

## Binding Constraints

- Exercise shipped CLI behavior only. If the journey exposes a product defect,
  stop and report rather than implementing a fix inside this audit task.
- Use the primitive roadmap path exactly: create -> claim -> transition ->
  complete -> unclaim, with denied and unmet-rule attempts proven zero-append.
- Also cover AC verify, concern raise/address, both epic scheduling modes,
  label attach, requirement sync, valid/inverse and invalid links, completion-
  linked run telemetry, `maintenance check`, authorized vacuum, and one
  governance version guard.
- Run in a temp git repo: generated `.fid/.gitignore` must ignore only the DB;
  the journal must be addable/trackable. Do not make product code invoke git.
- Snapshot observable/projected state and journal bytes before DB deletion;
  rebuild, then compare structurally without timestamp-string fragility.

## Scope

### Files To Create Or Modify

- `cmd/fid/e2e_v01_closure_test.go` - the integration-tagged journey.
- `cmd/fid/e2e_v01_closure_helpers_test.go` - test-only snapshots/helpers if the
  primary test file would become hard to review.
- `project-docs/specs/07_DESIGN_PLAN_FLOWS.md` - only if a shipped flow has
  drifted; otherwise record verified-no-change.

### Files To Read For Context

- `cmd/fid/e2e_v01b_test.go`
- `cmd/fid/e2e_v01d_closure_test.go`
- `cmd/fid/e2e_v01d_link_test.go`
- `cmd/fid/e2e_v01e_dev_test.go`
- `cmd/fid/e2e_v01e_activity_test.go`
- `project-docs/specs/10_IMPLEMENTATION_ROADMAP.md` v0.1 exit criteria
- `project-docs/specs/08_QUALITY_PLAN.md` sections 5.2 and 9

### Out Of Scope

- Fixing behavior discovered by the smoke.
- Hidden/held-out scenario material or evaluator prompts.
- Public-release threat model, release notes, packaging, or publishing.
- v0.2+ governance CRUD, reports, or workflow configuration.

## Acceptance Criteria

- [ ] Lifecycle and all rejection branches are observable and journal-atomic.
  (VAL-018)
- [ ] Both ready scheduling modes return only the correct next task. (VAL-019)
- [ ] Labels and issue/requirement link inverses read back; invalid relation
  writes nothing. (VAL-020)
- [ ] Completion-linked telemetry produces the expected snapshot prices, USD
  cost, and stats aggregates. (VAL-021)
- [ ] Maintenance passes the coherent fixture, vacuum preserves it, DB deletion
  plus rebuild restores equal state and audit output. (VAL-022)
- [ ] Temp-git checks prove journal trackability, DB ignore, roles after rebuild,
  and byte-identical journal across all non-journal commands. (VAL-023)

## First Failing Test

`TestV01ClosureJourneyAndRebuildParity` is written end-to-end first and fails at
the first missing Phase 17 command or unmet closure assertion.

## Documentation Updates

- `project-docs/specs/07_DESIGN_PLAN_FLOWS.md` - only if the smoke proves drift;
  otherwise record verified-no-change.

## Exact Commands To Run

Run exactly:

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
go test -tags integration -race ./...
go test -bench . -benchmem ./...
go build -o /tmp/fid ./cmd/fid
CGO_ENABLED=0 go build -o /tmp/fid-static ./cmd/fid
stat -c '%s %n' /tmp/fid-static
test "$(stat -c %s /tmp/fid-static)" -lt 31457280
```

Required-command rules live in `AGENTS.md`. Every command is reported `PASS`,
or `BLOCKED` with reason; all commands above are required before task close.

List any conditional command required by this task and why:

- None.

## Notes For Developer

- Depends on tasks 2-8.
- QA Lead owns the separate phase-end v0.1 PRD-to-test traceability audit
  (VAL-024); this task supplies representative journey evidence, not a claim
  that one smoke proves every PRD row.
- Completion report: `project-docs/docs/phases/phase-17/tasks/task-9-done.md`.
