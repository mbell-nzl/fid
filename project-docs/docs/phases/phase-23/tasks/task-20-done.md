# Phase 23 - Task 20: Stable requirement identity on every reference - Blocker Report

## Workflow outcome

- Task outcome: BLOCKED
- Category: specification-contradiction
- Condition: The accepted `RequirementRef` shape requires `requirement_id`, but a task-base `requirement.referenced` journal stores an exact `RequirementRef` result without it. The current exact result validator rejects that immutable event after the shape change.
- Recovery attempts: 1
- Next action: Obtain protected authority for a versioned stored-result compatibility contract and its file scope, or amend the requirement. Preserve the source journal and strict validation; then resume task-20.
- Human attention required: yes

## Summary of what was implemented

- Added failing built-binary assertions for stable requirement identity and the Human requirement labels, then made a scoped diagnostic implementation of the new value shape and read result. The implementation is deliberately incomplete because the pre-change replay criterion fails.
- Pre-change reference commit: `d7f960d94252a7fa07df4aabb0798ce160c94763` (`docs(task): record task-19 report-lint recovery`). Before product edits, its binary seeded `/tmp/fid-task20-prechange.NjCRbu/judgment/.fid/journal/2026-10-01.jsonl`. The 22-event journal contains one `requirement.referenced` and one `requirement.linked` event. Its SHA-256 before and after the candidate rebuild was `a8fdbdccfc7fd39e27b80f256811d8b70113a0e30285122a166148455d1ba16e`.
- The pre-change `requirement.referenced` Event is line 18. Its stored `result_json.requirement` keys are `created_at`, `event_id`, `id`, `source`, `status`, and `superseded_by_id`; `requirement_id` is absent. The candidate rebuild stops at that line with `journal-replay-stopped` (exit 6). Adding this field to the closed shape makes the historical result fail exact validation. Task-20 explicitly directs `BLOCKED` if the stored result shape prevents unchanged replay.

## Files created or modified

- `internal/review/requirement.go` (partial diagnostic implementation)
- `internal/journal/current_event_value_catalog.go` (partial diagnostic implementation)
- `internal/journal/current_event_result_projection.go` (partial diagnostic implementation)
- `cmd/fid/integration/p23_requirement_reference_test.go` (failing identity regression and pre-change fixture seed)
- `cmd/fid/integration/p23_human_leaf_pairs_test.go` (failing Human fixture expectation)
- `project-docs/docs/phases/phase-23/tasks/task-20-done.md`

## Documentation updated

- `project-docs/docs/phases/phase-23/tasks/task-20-done.md` only. Accepted contracts and Runner-owned files were left untouched.

## Validation Contract Coverage

- VAL-089: BLOCKED - `TestP23RequirementReferenceAndLinkRealProcess` failed on absent `result.requirement.requirement_id` at `p23_requirement_reference_test.go:35`; `TestP23HumanFixtures` failed on the old reference heading/labels; the task-base journal rebuild stopped at its `requirement.referenced` Event line 18 after the scoped shape change. The unchanged digest is recorded above. The required replay state is not achieved.
- VAL-073: BLOCKED - the regenerated `TestP23HumanFixtures` expectation demonstrates the requirement view label gap, but the protected replay incompatibility prevents completion and the later Human re-check of HUMAN-04a/b/c-17 and HUMAN-04a-18.
- VAL-058: BLOCKED - the Human leaf oracle was updated to derive requirement and link text from JSON fields, but the fixture remains red and the complete no-Human-only-member check cannot pass until the contract conflict is resolved.

## Tests run

- `go test -tags integration -run '^TestP23RequirementReferenceAndLinkRealProcess$' -count=1 ./cmd/fid/integration` - expected RED - `FAIL` (terminated `exit_code: 1`); the failure was `path result.requirement.requirement_id = <nil>, want non-empty string`.
- `FID_TASK20_PRECHANGE_ROOT=/tmp/fid-task20-prechange.NjCRbu go test -tags integration -run '^TestP23Task20SeedPrechangeJournal$' -count=1 -v ./cmd/fid/integration` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	1.911s` (terminated `exit_code: 0`)
- `go test -tags integration -run '^TestP23HumanFixtures$' -count=1 ./cmd/fid/integration` - expected RED - `FAIL` (terminated `exit_code: 1`); the fixture showed `Requirement References and links` and `Requirement RREF-8227f5` where it expected separate stable requirement and reference labels.
- `go build -o /tmp/fid-task20-candidate ./cmd/fid` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `/tmp/fid-task20-candidate --project /tmp/fid-task20-prechange.NjCRbu/judgment --json rebuild` - FAIL - `{"contract_version":"fid-cli.v1","error":{"category":"integrity","code":"journal-replay-stopped","details":{"event_id":null,"journal_file":".fid/journal/2026-10-01.jsonl","line":18},"message":"journal replay stopped","retryable":false,"subject":null},"ok":false,"operation":"project.rebuild"}` (terminated `exit_code: 6`).
- `sha256sum /tmp/fid-task20-prechange.NjCRbu/judgment/.fid/journal/2026-10-01.jsonl` - PASS - `a8fdbdccfc7fd39e27b80f256811d8b70113a0e30285122a166148455d1ba16e  /tmp/fid-task20-prechange.NjCRbu/judgment/.fid/journal/2026-10-01.jsonl` (terminated `exit_code: 0`)

## Exact quality-gate results

- `gofmt -l cmd internal` - BLOCKED - not run after the protected replay contradiction; no termination result.
- `golangci-lint run ./...` - BLOCKED - not run after the protected replay contradiction; no termination result.
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` - BLOCKED - not run after the protected replay contradiction; no termination result. Thus `cmd/fid` and `internal/workflow` wall times were not measured.
- `go tool cover -func=coverage.out` - BLOCKED - not run after the protected replay contradiction; no termination result.
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...` - BLOCKED - not run after the protected replay contradiction; no termination result. The tagged integration wall time was not measured.

## Deviations from spec

- The new identity and Human rendering are unfinished. The task's immutable pre-change replay requirement failed at the stored result shape. No history, validator, request, Event, Effect, or digest was rewritten or weakened.

## Concerns, technical debt, or notes for the Tech Lead

- `internal/journal/current_event_json.go` routes `requirement.referenced` stored results to the now-amended closed `RequirementRef` shape. `internal/journal/current_event_result_projection.go` also compares the stored result exactly with the current projection. A versioned compatibility rule would touch validation semantics and an out-of-scope registry; task-20 provides no authority for it.
- The task-base journal fixture remains at the path above for inspection. The Runner-owned `run-ledger.md` change predates this Developer work and was not edited here.
- No package headroom measurements are available because the required long gates did not run.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `No commit created - protected blocker report`

## Commit subject

- `No commit created - protected blocker report`
