# Phase 21 - Task 22: Final Wire Conformance And Real-Process Journey

## Objective

Prove every final accepted/fail-closed public leaf, all kebab-case closed values, the Gate-before-terminal-transition real-process journey, its negative counterparts, and the mutation-free post-rebuild repeat.

## Validation Contract Coverage

- VAL-072: owning observable in the Phase 21 validation contract.
- VAL-073: owning observable in the Phase 21 validation contract.
- VAL-084: owning observable in the Phase 21 validation contract.
- VAL-085: owning observable in the Phase 21 validation contract.
- VAL-086: owning observable in the Phase 21 validation contract.

## Binding Constraints

- Build and invoke the final binary; do not mock or use in-process substitutes for the journey.
- Mutate once, rebuild durable facts, then repeat only specified reads/evaluation/completion state/resume/export/negative checks.
- Do not read held-out scenarios or create a second Envelope.

## Scope

### Files To Create Or Modify

- `cmd/fid/integration/`
- `internal/cli/`
- `project-docs/docs/phases/phase-21/tasks/task-22-done.md`

### Files To Read For Context

- `project-docs/docs/phases/phase-21/public-interface.md`
- `project-docs/docs/phases/phase-21/contracts.md §16`
- `project-docs/docs/phases/phase-21/gate-d3-2-proposal.md`
- `project-docs/docs/phases/phase-21/qa/validation-contract.md`

### Out Of Scope

- Product repairs outside accepted contracts, import, recovery, Review/Human Outcome, and held-out scenarios.

## Acceptance Criteria

- [ ] Every nonstreaming leaf is exact and unknown fields reject. (VAL-072)
- [ ] Every closed public value is kebab-case. (VAL-073)
- [ ] Happy journey follows Gate/evidence before terminal transition and completes. (VAL-084)
- [ ] All negative counterparts fail closed. (VAL-085)
- [ ] Post-rebuild repeat performs no second Envelope or other mutation. (VAL-086)

## First Failing Test

Add the generated interface corpus, forbidden-enum inventory, ordered binary journey, negative matrix, and post-rebuild Event-count assertion.

## Documentation Updates

- `project-docs/docs/phases/phase-21/tasks/task-22-done.md` - record red/green evidence, deviations, required commands, and spec impact.

## Exact Commands To Run

Run exactly:

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
go test -tags integration -race ./cmd/fid/integration/...
go test -bench . -benchmem ./...
go build -o /tmp/fid ./cmd/fid
project-docs/tools/report-lint.sh project-docs/docs/phases/phase-21/tasks/task-22-done.md
```

Every command is required and must terminate with session-captured evidence. The report-lint command runs before the next dispatch. A Ready QA plan may add a task-specific command; the Tech Lead must amend this section before authorization if it does.

## Notes For Developer

- Start after task 21 and QA Review window 7 close. Final QA follows; this task does not authorize its own contract refinements.
- Completion report: `project-docs/docs/phases/phase-21/tasks/task-22-done.md`.
- Stop for a protected contract change; ordinary in-scope red/green repair remains this task.
- Do not read held-out scenario material.
