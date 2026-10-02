# Phase 23 - Task 20: Stable requirement identity on every reference - Completion Report

## Workflow outcome

- Task outcome: COMPLETE
- Category: none
- Condition: None. DEC-270 resolved the attempt-1 specification contradiction.
- Recovery attempts: 2
- Next action: None.
- Human attention required: no

## Summary of what was implemented

- Attempt 1 added failing built-binary assertions for stable requirement identity and Human labels, then made the scoped partial value-shape and read-result changes. Attempt 2 completed the Human requirement views and kept the strict journal validator unchanged. `RequirementRef` now includes `requirement_id` in record-reference results, exact reads, lists, and idempotent retries.
- Pre-change reference commit: `d7f960d94252a7fa07df4aabb0798ce160c94763` (`docs(task): record task-19 report-lint recovery`). Before product edits, its binary seeded `/tmp/fid-task20-prechange.NjCRbu/judgment/.fid/journal/2026-10-01.jsonl`. The 22-event journal contains one `requirement.referenced` and one `requirement.linked` event. Its SHA-256 before and after the candidate rebuild was `a8fdbdccfc7fd39e27b80f256811d8b70113a0e30285122a166148455d1ba16e`.
- The pre-change `requirement.referenced` Event is line 18. Its stored `result_json.requirement` keys are `created_at`, `event_id`, `id`, `source`, `status`, and `superseded_by_id`; `requirement_id` is absent. The attempt-1 candidate rebuild stopped at that line with `journal-replay-stopped` (exit 6). The original criterion required unchanged replay, so attempt 1 correctly returned `BLOCKED`. DEC-270 now requires fail-closed handling of that unsupported history; the final binary produced the same refusal without changing the journal digest.
- `TestP23RequirementNewShapeJournalRoundTrip` records two references and two Issue links, captures five JSON reads, rebuilds, checks the source journal digest and byte-identical reads, then retries the first idempotency key and checks its stable ID without a new Event. The tagged gate ran this committed test.

## Files created or modified

- `internal/review/requirement.go` (read and stored result identity)
- `internal/journal/current_event_value_catalog.go` (`RequirementRef` value shape)
- `internal/journal/current_event_result_projection.go` (retain projected stable identity)
- `cmd/fid/p23_human.go` (separate Requirement Reference and Issue link rendering)
- `cmd/fid/p23_human_test.go` (converged retained Human text assertion)
- `cmd/fid/integration/p23_requirement_reference_test.go` (identity and get/list assertions; removed the diagnostic seed test)
- `cmd/fid/integration/p23_requirement_roundtrip_test.go` (new-shape rebuild and idempotent retry)
- `cmd/fid/integration/p23_requirement_leaf_oracle_test.go` (converged retained result shape assertion)
- `cmd/fid/integration/p23_human_leaf_pairs_test.go` (converged retained Human fixture rows)
- `project-docs/docs/phases/phase-23/tasks/task-20-done.md`

## Documentation updated

- `project-docs/docs/phases/phase-23/tasks/task-20-done.md` only. Accepted contracts and Runner-owned files were left untouched.

## Validation Contract Coverage

- VAL-089: PASS for deterministic task evidence - `TestP23RequirementReferenceAndLinkRealProcess` checks both identities in record results, stable-ID and exact-reference-ID gets, active/all/Issue-filtered lists, supersession, and stable Issue links. `TestP23HumanFixtures` checks the verbatim Requirement heading, reference line, source line, and Issue links in all four requirement leaves. `TestP23RequirementNewShapeJournalRoundTrip` checks rebuild, unchanged digest, five byte-identical JSON reads, and retry identity. The final-binary refusal and before/after SHA-256 lines below prove the DEC-270(b) fail-closed state. Attempt 1 was `BLOCKED` on the original criterion as recorded above.
- VAL-073: Deterministic rendering evidence PASS - `TestP23HumanFixtures` and the direct `fid requirement get REQ-leaf` output below show purpose and labelled sections without routine narration. Human-owned judgment for HUMAN-04a/b/c-17 and HUMAN-04a-18 remains for the phase-end Human re-check.
- VAL-058: PASS for the requirement rows - `TestP23HumanFixtures` compares each Human view with JSON facts and rejects any Human-only fact; `TestPhase23HumanLeafTextMatchesRecordedFacts` checks the same constraint for retained unit rows. Attempt 1 was `BLOCKED` while the fixture was red.

## Tests run

- `go test -tags integration -run '^TestP23RequirementReferenceAndLinkRealProcess$' -count=1 ./cmd/fid/integration` - expected RED - `FAIL` (terminated `exit_code: 1`); the failure was `path result.requirement.requirement_id = <nil>, want non-empty string`.
- `FID_TASK20_PRECHANGE_ROOT=/tmp/fid-task20-prechange.NjCRbu go test -tags integration -run '^TestP23Task20SeedPrechangeJournal$' -count=1 -v ./cmd/fid/integration` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	1.911s` (terminated `exit_code: 0`)
- `go test -tags integration -run '^TestP23HumanFixtures$' -count=1 ./cmd/fid/integration` - expected RED - `FAIL` (terminated `exit_code: 1`); the fixture showed `Requirement References and links` and `Requirement RREF-8227f5` where it expected separate stable requirement and reference labels.
- `go build -o /tmp/fid-task20-candidate ./cmd/fid` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `/tmp/fid-task20-candidate --project /tmp/fid-task20-prechange.NjCRbu/judgment --json rebuild` - FAIL - `{"contract_version":"fid-cli.v1","error":{"category":"integrity","code":"journal-replay-stopped","details":{"event_id":null,"journal_file":".fid/journal/2026-10-01.jsonl","line":18},"message":"journal replay stopped","retryable":false,"subject":null},"ok":false,"operation":"project.rebuild"}` (terminated `exit_code: 6`).
- `sha256sum /tmp/fid-task20-prechange.NjCRbu/judgment/.fid/journal/2026-10-01.jsonl` - PASS - `a8fdbdccfc7fd39e27b80f256811d8b70113a0e30285122a166148455d1ba16e  /tmp/fid-task20-prechange.NjCRbu/judgment/.fid/journal/2026-10-01.jsonl` (terminated `exit_code: 0`)
- `go test -tags integration -run '^TestP23RequirementNewShapeJournalRoundTrip$' -count=1 ./cmd/fid/integration` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	3.537s` (terminated `exit_code: 0`).
- `go test -tags integration -run '^TestP23HumanFixtures$' -count=1 ./cmd/fid/integration` - FAIL (expected RED) - `FAIL` (terminated `exit_code: 1`); the output showed `Requirement References and links` and `Requirement RREF-3e34b0` instead of separate stable and reference labels.
- `go test -tags integration -run '^TestP23HumanFixtures$' -count=1 ./cmd/fid/integration` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	11.288s` (terminated `exit_code: 0`).
- `go test -run '^TestPhase23HumanLeafTextMatchesRecordedFacts$' -count=1 ./cmd/fid` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	0.154s` (terminated `exit_code: 0`).
- `go test -tags integration -run '^TestP23SharedLeafConformance$' -count=1 ./cmd/fid/integration` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	10.159s` (terminated `exit_code: 0`).
- `go test -tags integration -run '^TestP23RequirementReferenceAndLinkRealProcess$' -count=1 ./cmd/fid/integration` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	1.965s` (terminated `exit_code: 0`).
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 ./cmd/fid ./internal/review ./internal/journal` - PASS - `ok  	github.com/mbell-nzl/fid/internal/journal	101.603s` (terminated `exit_code: 0`); `cmd/fid` 221.604s and `internal/review` 1.101s.
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 -run '^(TestP23HumanFixtures|TestP23SharedLeafConformance|TestP23RequirementReferenceAndLinkRealProcess|TestP23RequirementNewShapeJournalRoundTrip)$' ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	16.419s` (terminated `exit_code: 0`).
- `sha256sum /tmp/fid-task20-prechange.NjCRbu/judgment/.fid/journal/2026-10-01.jsonl` (before final rebuild) - PASS - `a8fdbdccfc7fd39e27b80f256811d8b70113a0e30285122a166148455d1ba16e  /tmp/fid-task20-prechange.NjCRbu/judgment/.fid/journal/2026-10-01.jsonl` (terminated `exit_code: 0`).
- `go build -o /tmp/fid-task20-candidate ./cmd/fid` (final binary) - PASS - `exit 0 (no output)` (terminated `exit_code: 0`).
- `/tmp/fid-task20-candidate --project /tmp/fid-task20-prechange.NjCRbu/judgment --json rebuild` - FAIL (expected DEC-270 refusal) - `{"contract_version":"fid-cli.v1","error":{"category":"integrity","code":"journal-replay-stopped","details":{"event_id":null,"journal_file":".fid/journal/2026-10-01.jsonl","line":18},"message":"journal replay stopped","retryable":false,"subject":null},"ok":false,"operation":"project.rebuild"}` (terminated `exit_code: 6`).
- `sha256sum /tmp/fid-task20-prechange.NjCRbu/judgment/.fid/journal/2026-10-01.jsonl` (after final rebuild) - PASS - `a8fdbdccfc7fd39e27b80f256811d8b70113a0e30285122a166148455d1ba16e  /tmp/fid-task20-prechange.NjCRbu/judgment/.fid/journal/2026-10-01.jsonl` (terminated `exit_code: 0`).
- `FID_P23_HUMAN_FIXTURE_ROOT=/tmp/fid-task20-runtime.I3sG3G go test -tags integration -run '^TestP23HumanFixtures$' -count=1 ./cmd/fid/integration` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	11.921s` (terminated `exit_code: 0`); fixture setup for the direct CLI check.
- `/tmp/fid-task20-candidate --project /tmp/fid-task20-runtime.I3sG3G/checkpoints-1169491080/leaf-contract requirement get REQ-leaf` - PASS - `  REQ-leaf implements P22-0b1cdb` (terminated `exit_code: 0`); output also contained `Requirement REQ-leaf`, `Reference ID: RREF-1334c6; Status: active`, and `Issue links`.

## Exact quality-gate results

- Attempt 1: `gofmt -l cmd internal` - BLOCKED - not run after the protected replay contradiction; no termination result.
- Attempt 1: `golangci-lint run ./...` - BLOCKED - not run after the protected replay contradiction; no termination result.
- Attempt 1: `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` - BLOCKED - not run after the protected replay contradiction; no termination result. Thus `cmd/fid` and `internal/workflow` wall times were not measured.
- Attempt 1: `go tool cover -func=coverage.out` - BLOCKED - not run after the protected replay contradiction; no termination result.
- Attempt 1: `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...` - BLOCKED - not run after the protected replay contradiction; no termination result. The tagged integration wall time was not measured.
- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`).
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`).
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	117.580s	coverage: 62.1% of statements` (terminated `exit_code: 0`); `cmd/fid` 344.249s, `internal/workflow` 117.580s; neither reaches 80% of its 600s ordinary-package budget.
- `go tool cover -func=coverage.out` - PASS - `total:											(statements)			60.5%` (terminated `exit_code: 0`).
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	253.239s` (terminated `exit_code: 0`); tagged package 253.239s, below 80% of its 900s budget.

## Deviations from spec

- DEC-270 amended the original pre-change replay criterion after attempt 1 proved it impossible. Under ADR-0011, pre-task-20 Phase 23 `requirement.referenced` history is unsupported and must fail closed. The replacement requires a committed new-shape round trip plus one-time refusal and unchanged-digest evidence. No history, validator, request shape, Event type, Effect table/key, or digest was rewritten or weakened. The attempt-1 failure token was `journal-replay-stopped` at line 18 (terminated `exit_code: 6`); the final binary returned the same token and termination as the expected DEC-270 result.

## Concerns, technical debt, or notes for the Tech Lead

- Retained assertion provenance: `cmd/fid/p23_human_test.go` requirement row (introduced by `59d2653da7cfe14699dc4d1cd85ef032bec26a5c`, phase-23/task-14) still checks exact Human text and JSON-field sourcing; DEC-269 Part 3 changed the expected stable heading, reference line, and link section. It ran under the focused `go test -run '^TestPhase23HumanLeafTextMatchesRecordedFacts$' -count=1 ./cmd/fid` and ordinary gate.
- Retained assertion provenance: `cmd/fid/integration/p23_human_leaf_pairs_test.go` requirement rows (introduced by `b71ed1ba18a277da22a71c5fdcf18eadb181be26`, phase-23/task-14) still compare verbatim Human output with JSON facts. DEC-269 Part 3 displaced only requirement labels and link grouping. They ran under `TestP23HumanFixtures` and the tagged gate.
- Retained assertion provenance: `cmd/fid/integration/p23_requirement_leaf_oracle_test.go` reference shape row (introduced by `2e7fd2f9fa05bce648510789453c8d2ed8633b1d`, phase-23/task-18) still compares fake and real leaf results with full shape and identity. DEC-269 Part 3 added the stable ID to that row. It ran under `TestP23SharedLeafConformance` and the tagged gate.
- The task-base journal fixture remains at the path above for inspection. The Runner-owned `run-ledger.md` and `triage-log.md` changes were not edited here. No package reached 80% of its time budget.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `0ba005d7c511d7f5acd1517d1c39dd6e58d3e3e0`

## Commit subject

- `feat(requirement): expose stable identity and distinct Human views`
