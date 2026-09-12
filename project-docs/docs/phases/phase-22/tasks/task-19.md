# Phase 22 - Task 19: Explain recovery and restored work to Humans

## Objective

Repair the Human-facing recovery views after the 2026-09-12 checklist recorded
one PASS, two FAIL and three PARTIAL outcomes. Execute as one serialized root
implementation followed by one independent review under the Human's explicit
approval; this supplements the eighteen completed planned tasks.

## Validation Contract Coverage

- VAL-053: retain Human/JSON factual agreement.
- VAL-076: explain the next recovery action and adoption versus proof.
- VAL-077: explain activation recovery direction and its consequences.
- VAL-078: identify preserved history and destination intervention requirements.
- Human judgments remain Human-owned; automated evidence establishes readiness only.

## Binding Constraints

- DEC-219 allocates this repair. Existing mutation authority, JSON contracts,
  journal Events, schemas, dependencies and gate definitions stay unchanged.
- Presentation uses existing domain projections and verified local import
  metadata. Historical observations cannot be presented as destination proof.
- A ready operation is a candidate whose request and authority are still checked,
  not permission inferred from an active Claim or a blanket recommended mutation.
- Preserve the original checklist and every failed test/gate attempt.

## Scope

### Files To Create Or Modify

- cmd/fid/resume.go, resume_human.go, resume_human_guidance.go,
  run_recovery.go, run_inspect_human.go, activation.go, activation_human.go,
  human_recovery_guidance_test.go.
- internal/workflow/recovery_read_context.go and recovery_read_context_test.go:
  read-only context using existing import binding and observation records.
- cmd/fid/integration/p22_resume_human_test.go,
  p22_human_fixtures_test.go, p22_human_guidance_test.go: behavioral and
  retained-fixture coverage; preserve semantic assertions when wording changes.
- Phase decisions/public-interface, task-19-done.md, QA review/checklist retry
  artifacts and current/phase progress records.

### Files To Read For Context

- contracts.md sections 3, 10; public-interface.md section 6;
  qa/validation-contract.md; qa/human-checklist.md; existing workflow import,
  observation, reconciliation, activation and resume implementations.

### Out Of Scope

- CLI-wide redesign, new source-path persistence, policy changes, archive or
  performance changes, hidden scenarios and completed task report rewrites.

## Acceptance Criteria

- [x] Resume groups immediate actions, recovery context, blockers and completion
  evidence; labels Actor IDs and explains each candidate operation. (VAL-076)
- [x] Adoption explicitly retains existing uncommitted changes for a successor
  without establishing testing, review or correctness. (VAL-076)
- [x] Run inspection explains paths, reconciliation and observed/unknown or
  historical liveness, without claiming it records an observation. (VAL-053/078)
- [x] Imported project/history identity and destination are recognizable; source
  and destination worktrees are distinguished. Fresh destination observations
  and checks are explained separately from completion Evidence. (VAL-078)
- [x] Activation shows resource identities, explains absent prior resources,
  both recovery directions and the activation commit versus Git. (VAL-077)
- [x] Normal, blocked, expired, historical and fresh observation states retain
  factual correctness; JSON shapes and mutation behavior remain unchanged.
- [x] Independent review and fresh neutral fixtures precede Human retest;
  original failed judgments remain intact and acceptance is not inferred.

## First Failing Test

TestHumanRecoveryGuidance explains adoption and separates completion from
immediate actions; capture its behavioral failure before renderer changes.

## Documentation Updates

Record DEC-219 and Human output expectations; preserve attempt 1 and prepare a
separate human-checklist-2.md with command purposes and blank judgments.
Completion report records exact evidence and retained-fixture provenance.

## Exact Commands To Run

Working directory: /home/mark/projects/fid. Apply gofmt/goimports first.

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
go test -tags integration -race ./cmd/fid/integration/...
go build -o /tmp/fid ./cmd/fid
FID_P22_HUMAN_FIXTURE_ROOT=/home/mark/projects/fid/.cache/phase-22-human go test -tags integration ./cmd/fid/integration/... -run '^TestP22HumanFixtures$' -count=1 -v
```

All commands must terminate successfully; no nonzero expectations or task-gate
exceptions. Focused red/green is additional. Prior phase benchmark/archive-scale
qualification remains historical evidence; this task changes neither subsystem.
Run heavy gates serially. Record raw output and observed termination for every
attempt. Independently review the final diff and exercise its public commands.
