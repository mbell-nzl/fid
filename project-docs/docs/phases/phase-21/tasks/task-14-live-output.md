# Phase 21 - Task 14: Live Framing, Redaction, And Controls

## Objective

Implement exact `fid.live.v1` stdout/stderr/control framing, pre-surface redaction, bounded truncation/draining, and Human rendering. Prepare the recorded terminal fixture for Human VAL-036 without self-validating it.

## Validation Contract Coverage

- VAL-036: owning observable in the Phase 21 validation contract.
- VAL-047: owning observable in the Phase 21 validation contract.
- VAL-048: owning observable in the Phase 21 validation contract.
- VAL-049: owning observable in the Phase 21 validation contract.
- VAL-050: owning observable in the Phase 21 validation contract.

## Binding Constraints

- All public values are kebab-case and payload unions are closed.
- Redaction precedes every emitted, stored, hashed, error, Artifact, export, and Human surface.
- The Human alone validates VAL-036.

## Scope

### Files To Create Or Modify

- `internal/evidence/`
- `internal/cli/`
- `cmd/fid/`
- `cmd/fid/testdata/val-036/human-session.txt`
- `cmd/fid/testdata/val-036/live-frames.jsonl`
- `project-docs/docs/phases/phase-21/tasks/task-14-done.md`

### Files To Read For Context

- `project-docs/docs/phases/phase-21/public-interface.md §8`
- `project-docs/docs/phases/phase-21/contracts.md §8`

### Out Of Scope

- Attempt/Receipt persistence and Human verdict recording.

## Acceptance Criteria

- [ ] Machine data frames are exact and globally sequenced. (VAL-047)
- [ ] Every control variant is exact. (VAL-048)
- [ ] Canary redaction covers all surfaces. (VAL-049)
- [ ] Truncation drains without changing outcome. (VAL-050)
- [ ] `human-session.txt` and `live-frames.jsonl` are regenerated and byte-compared by a golden test; together they show interleaved stdout/stderr, all three control kinds (`redaction-applied`, `artifact-rejected`, `output-truncated`), a non-success terminal frame, and child output containing raw ESC, CR, and BEL bytes rendered safely. (VAL-036)

## First Failing Test

Add the failing golden regeneration/byte-comparison test for both VAL-036 fixtures, then the interleaved/control/canary/high-output real-process fixtures before framing.

## Documentation Updates

- `project-docs/docs/phases/phase-21/tasks/task-14-done.md` - record red/green evidence, deviations, required commands, and spec impact.

## Exact Commands To Run

Run exactly:

```text
go test -race ./internal/evidence ./internal/cli ./cmd/fid/...
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

Every command is required and must terminate with session-captured evidence. Report lint is not a Developer task-gate command; after the report commit, the Phase Runner runs the phase brief's two-argument pre-dispatch control. A Ready QA plan may add a task-specific command; the Tech Lead must amend this section before authorization if it does.

## Notes For Developer

- Start after task 13 closes.
- Completion report: `project-docs/docs/phases/phase-21/tasks/task-14-done.md`.
- Stop for a protected contract change; ordinary in-scope red/green repair remains this task.
- Do not read held-out scenario material.
