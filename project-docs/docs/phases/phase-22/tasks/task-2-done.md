# Phase 22 - Task 2: Typed Stops and Attention domain - Completion Report

## Workflow outcome

- Task outcome: COMPLETE
- Category: none
- Condition: None.
- Recovery attempts: 0
- Next action: None.
- Human attention required: no

## Summary of what was implemented

- Added idempotent Stop raise and resolution operations with immutable original facts, typed conditions, closed codes, subject-bound authority, and qualifying Event checks.
- Added SHA-256 semantic Attention identities, one-open-item reuse with distinct Stop links, exact disposition-Event-set resolution, and atomic Event effects.
- Activated only `stop.raise`, `stop.resolve`, and `attention.resolve` as the approved profile actions while preserving the unavailable `attention.v1` binary capability and absent CLI leaf.
- Reopened the Stop/Attention Event projection family while keeping rejected reads and unavailable recovery branches free of Stop/Attention effects.

## Files created or modified

- `internal/workflow/stop.go`
- `internal/workflow/stop_test.go`
- `internal/workflow/attention_test.go`
- `internal/workflow/activation.go`
- `internal/journal/project_effects.go`
- `cmd/fid/contract_surface_test.go`
- `project-docs/docs/phases/phase-22/tasks/task-2-done.md`

## Documentation updated

- `project-docs/docs/phases/phase-22/tasks/task-2-done.md`

## Validation Contract Coverage

- VAL-004: PASS - `TestAttentionResolutionRequiresExactDispositionEventsAndKeepsStopsImmutable` at `internal/workflow/attention_test.go:12` resolves both linked Stops and confirms their original `operator-stop` code and details remain unchanged at `:60`.
- VAL-005: PASS - `TestEquivalentStopsShareAttentionButKeepDistinctFacts` at `internal/workflow/stop_test.go:11` proves two Stop rows, one open Attention item, and two links at `:31`; `TestRaiseStopUsesOneOpenAttentionAndRejectsRoutineReadsWithoutEffects` at `:46` verifies the real workflow operation returns distinct Stop IDs with the same Attention ID at `:69`.
- VAL-006: PASS - `TestAttentionResolutionRequiresExactDispositionEventsAndKeepsStopsImmutable` rejects a missing disposition Event set with unchanged journal bytes at `internal/workflow/attention_test.go:42`, then accepts the sorted exact set and produces resolved version 2 at `:50`.
- VAL-007: PASS - `TestRaiseStopUsesOneOpenAttentionAndRejectsRoutineReadsWithoutEffects` proves the rejected `GetStop` read leaves journal bytes unchanged at `internal/workflow/stop_test.go:72`; only explicit `RaiseStop` creates Stop/Attention facts.
- VAL-079: PASS - direct in-package positive and rejection branches are `TestEquivalentStopsShareAttentionButKeepDistinctFacts`, `TestRaiseStopUsesOneOpenAttentionAndRejectsRoutineReadsWithoutEffects`, and `TestAttentionResolutionRequiresExactDispositionEventsAndKeepsStopsImmutable`; the required pre-edit and final coverage commands measured 27.4% and 34.7% respectively.

## Tests run

- `go test -race -cover ./internal/workflow` - PASS - `ok  \tgithub.com/mbell-nzl/fid/internal/workflow\t26.550s\tcoverage: 27.4% of statements` (terminated `exit_code: 0`)
  - Pre-edit CF-034 baseline.
- `go test -race ./internal/workflow -run '^TestEquivalentStopsShareAttentionButKeepDistinctFacts$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
  - First red attempt: `Event-v2 "stop.raised" does not permit effect application before its recovery domain is available`.
- `go test -race ./internal/workflow -run '^TestEquivalentStopsShareAttentionButKeepDistinctFacts$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
  - The initial green implementation exposed a test-harness Event-ID error: `parse event id "EVT-stop-1" date`.
- `go test -race ./internal/workflow -run '^TestEquivalentStopsShareAttentionButKeepDistinctFacts$' -count=1` - PASS - `ok  \tgithub.com/mbell-nzl/fid/internal/workflow\t2.141s` (terminated `exit_code: 0`)
- `go test -race ./internal/workflow -run '^TestAttentionResolutionRequiresExactDispositionEventsAndKeepsStopsImmutable$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
  - Red attempt: `aggregate attention/ATT-814a6f sequence 2, want 1`.
- `go test -race ./internal/workflow -run '^TestAttentionResolutionRequiresExactDispositionEventsAndKeepsStopsImmutable$' -count=1` - PASS - `ok  \tgithub.com/mbell-nzl/fid/internal/workflow\t3.260s` (terminated `exit_code: 0`)
- `/tmp/fid-phase22-task2 --project /tmp/tmp.gBaYkXcrYj --json init --prefix P22` - PASS - `{"contract_version":"fid-cli.v1","event":{"id":"EVT-20260905-000001","occurred_at":"2026-09-05T02:46:37.657134114Z"},"ok":true,"operation":"project.init","result":{"journal_path":".fid/journal/2026-09-05.jsonl","project":{"id":"project-a1b27a80ed0eb60ae3332e573f81d738","issue_prefix":"P22","root":"/tmp/tmp.gBaYkXcrYj"}}}` (terminated `exit_code: 0`)
- `/tmp/fid-phase22-task2 --project /tmp/tmp.gBaYkXcrYj --json capabilities` - PASS - `{"contract_version":"fid-cli.v1","event":null,"ok":true,"operation":"capabilities.get","result":{"items":[{"available":false,"capability_id":"attention.v1"` (terminated `exit_code: 0`)
  - Real CLI entrypoint confirms the accepted unavailable-capability boundary remains intact.

## Exact quality-gate results

- `go test -race -cover ./internal/workflow` - PASS - `ok  \tgithub.com/mbell-nzl/fid/internal/workflow\t(cached)\tcoverage: 34.7% of statements` (terminated `exit_code: 0`)
  - Final CF-034 measurement; earlier after-implementation attempt was `34.6%` (terminated `exit_code: 0`).
- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - FAIL - `* unused: 1` (terminated `exit_code: 1`)
  - Attempt 1: two unchecked `rows.Close` returns and one unused helper in `internal/workflow/stop.go`.
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
  - Attempt 2 after the in-scope lint repair.
- `golangci-lint run ./...` - FAIL - `* typecheck: 1` (terminated `exit_code: 1`)
  - Attempt 3: the converged retained test referenced `matches` before declaration.
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
  - Attempt 4 after correcting the retained test.
- `go test -race -coverprofile=coverage.out ./...` - FAIL - `FAIL` (terminated `exit_code: 1`)
  - Attempt 1: retained `TestPhase21AttentionBoundary` rejected the accepted task-2 Attention implementation surface.
- `go test -race -coverprofile=coverage.out ./...` - PASS - `ok  \tgithub.com/mbell-nzl/fid/internal/workflow\t(cached)\tcoverage: 34.7% of statements` (terminated `exit_code: 0`)
  - Attempt 2 after subject-preserving convergence.
- `go tool cover -func=coverage.out` - PASS - `total:\t\t\t\t\t\t\t\t\t\t(statements)\t\t\t\t61.6%` (terminated `exit_code: 0`)
- `go test -tags integration -race ./cmd/fid/integration/...` - PASS - `ok  \tgithub.com/mbell-nzl/fid/cmd/fid/integration\t8.157s` (terminated `exit_code: 0`)

## Deviations from spec

- No contract, schema, command, tolerance, or policy changed. The initial lint failure was caused by task-2 implementation hygiene and was repaired in scope; the byte-identical later lint attempt passed.
- The first full race/coverage attempt failed because `TestPhase21AttentionBoundary`, introduced by `7277139ee5d5dde4cd3627335d39b14f616f46a1` (`feat(evidence): evaluate trusted receipt evidence`, Refs: phase-21/task-17), prohibited every production Attention reference except the unavailable capability row. Task 2's accepted Stop/Attention domain and public-interface section 2.1 displaced that broad assertion. The converged test still proves the exact unavailable `attention.v1` row, absent `fid attention` leaf, no rejected-read effects, and limits Attention references to `activation.go`, `capabilities.go`, and task-2 `stop.go`. The byte-identical full gate later passed.
- The first intended red test established unavailable Stop Event effects. Its next attempt exposed a malformed test Event ID; correcting the fixture preserved the assertion. The attention-resolution red test then found and repaired the aggregate sequence from 2 to 1.

## Concerns, technical debt, or notes for the Tech Lead

- No CLI wiring was added because task 2 expressly excludes it. The direct workflow API and Event effects are ready for the later public CLI owner; `attention.v1` remains unavailable as required.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending commit

## Commit subject

- Pending commit
