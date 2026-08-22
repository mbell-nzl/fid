# Phase 21 - Task 14: Live Output Framing And Redaction

## Objective

Implement `fid.live.v1` machine framing and Human rendering over sanitized stdout/stderr/control/terminal records. Prove ordering, UTF-8/control handling, bounded chunking, output-limit draining, and redaction before every downstream representation.

## Validation Contract Coverage

- VAL-033: live frames are globally sequenced with one final terminal frame.
- VAL-034: redaction precedes every live/stored/hashed/rendered surface.
- VAL-035: output limits emit one control frame and do not hide termination.
- VAL-036: Human framing is understandable and control-safe.

## Binding Constraints

- Use one shared sanitizer/framer for machine and Human output; Human rendering is not a second semantic contract.
- Do not persist or hash pre-redaction bytes.
- Do not infer pass/fail from truncation.

## Scope

### Files To Create Or Modify

- `internal/evidence/` - executor framing/redaction files.
- `internal/cli/`
- `cmd/fid/`
- `project-docs/docs/phases/phase-21/tasks/task-14-done.md`

### Files To Read For Context

- `project-docs/docs/phases/phase-21/contracts.md §8`
- `project-docs/specs/05_DESIGN_PLAN.md §5`

### Out Of Scope

- Attempt/Receipt persistence and Artifact collection.

## Acceptance Criteria

- [ ] Dual streams, invalid UTF-8, control bytes, and chunk boundaries frame exactly. (VAL-033)
- [ ] Canary values appear on no downstream surface. (VAL-034)
- [ ] High output drains and terminates with one truncation control record. (VAL-035)
- [ ] The QA handoff identifies the exact terminal fixture for the later Human check. (VAL-036)

## First Failing Test

Add a failing deterministic framer test over interleaved streams, invalid UTF-8, controls, split code points, and redaction canaries.

## Documentation Updates

- `project-docs/docs/phases/phase-21/tasks/task-14-done.md` - record red/green evidence, deviations, required command results, and spec impact.

## Exact Commands To Run

Run exactly:

```text
go test -race ./internal/evidence ./internal/cli ./cmd/fid/...
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
project-docs/tools/report-lint.sh project-docs/docs/phases/phase-21/tasks/task-14-done.md
```

Every command above is required and must terminate with captured evidence. The report-lint command runs in this task before the next task is dispatched.

## Notes For Developer

- Start only after task 13 closes with required evidence.
- Completion report: `project-docs/docs/phases/phase-21/tasks/task-14-done.md`.
- Stop for protected contract changes; ordinary in-scope red/green repair remains this task.
- Do not read held-out scenario material.
