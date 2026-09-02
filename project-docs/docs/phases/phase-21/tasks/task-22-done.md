# Phase 21 - Task 22: Final Wire Conformance And Real-Process Journey - Completion Report

## Workflow outcome

- Task outcome: BLOCKED
- Category: specification-contradiction
- Condition: The accepted `contracts.md` section 16 mutation prefix fixes the journey at 20 Events with one nonterminal transition and one terminal-success transition. The shipped task-default workflow rejects the required terminal `in_progress -> done` transition with `transition-not-permitted` / `workflow-transition-invalid`; reaching `done` requires three extra transitions and produces 23 Events.
- Recovery attempts: 0
- Next action: Tech Lead/Human must resolve whether the accepted 20-Event journey changes, or whether the shipped task-default workflow is changed under explicit public-interface and workflow-contract authority.
- Human attention required: yes

## Summary of what was implemented

- Added the committed `p21_conformance` corpus and tagged byte-for-byte regeneration check for accepted/fail-closed operation leaves, section 7 failure codes, and section 8 live-frame variants.
- Added the real-binary journey prefix through activation, immutable Gate registration, Envelope creation/activation, Claim acquisition, Run start, nonterminal transition, Gate execution, Receipt Evidence link/evaluation, and the protected terminal-transition assertion.
- Preserved the first real-process failure at the terminal-success transition; no product workflow behavior or accepted contract was changed.

## Files created or modified

- `cmd/fid/integration/p21_conformance_test.go`
- `cmd/fid/integration/testdata/p21_conformance/`
- `project-docs/docs/phases/phase-21/tasks/task-22-done.md`

## Documentation updated

- `project-docs/docs/phases/phase-21/tasks/task-22-done.md`

## Validation Contract Coverage

- VAL-072: PASS - `TestPhase21ConformanceCorpusIsCompleteAndByteExact` at `cmd/fid/integration/p21_conformance_test.go:24` regenerates and byte-compares the committed operation, failure, and live fixture corpus.
- VAL-073: PASS - `TestPhase21ConformanceCorpusIsCompleteAndByteExact` at `cmd/fid/integration/p21_conformance_test.go:55` derives closed members from public-interface sections 3, 7, and 8, requires missing-member count 0, and rejects any non-kebab or underscore value.
- VAL-084: BLOCKED - `TestPhase21RealProcessJourneyStartsOneActivatedEnvelope` at `cmd/fid/integration/p21_conformance_test.go:125` reaches the real Gate and a satisfied Evidence evaluation, then observes the protected terminal workflow rejection before completion.
- VAL-085: BLOCKED - the required complete negative matrix cannot be finalized while the fixed happy journey cannot reach its terminal-success state under the accepted 20-Event bound.
- VAL-086: BLOCKED - exact pre/post-rebuild 20-Event parity depends on the unresolved terminal journey; adding the mandatory shipped workflow transitions produces 23 Events.

## Tests run

- `go test -tags integration ./cmd/fid/integration/...` - FAIL - `FAIL` (terminated `exit_code: 1`)
  Red evidence: the first corpus run reported missing `failures/artifact-not-found.json`; the corpus materialization repair made the unchanged tagged command pass before the later real-journey extension.
- `P21_WRITE_CONFORMANCE_FIXTURES=1 go test -tags integration ./cmd/fid/integration/...` - PASS - `ok  \tgithub.com/mbell-nzl/fid/cmd/fid/integration\t3.627s` (terminated `exit_code: 0`)
- `go test -tags integration ./cmd/fid/integration/...` - PASS - `ok  \tgithub.com/mbell-nzl/fid/cmd/fid/integration\t3.658s` (terminated `exit_code: 0`)
  Green evidence: tagged suite regenerated and byte-compared the committed corpus.
- `go test -tags integration ./cmd/fid/integration/...` - FAIL - `FAIL` (terminated `exit_code: 1`)
  Protected failure: the real final binary returned `transition-not-permitted` with `workflow-transition-invalid` for terminal `in_progress -> done`; its observed process exit was 4.

## Exact quality-gate results

- `gofmt -l cmd internal` - BLOCKED - `not run; protected specification contradiction before required gate` (no observed termination)
- `golangci-lint run ./...` - BLOCKED - `not run; protected specification contradiction before required gate` (no observed termination)
- `go test -race -coverprofile=coverage.out ./...` - BLOCKED - `not run; protected specification contradiction before required gate` (no observed termination)
- `go tool cover -func=coverage.out` - BLOCKED - `not run; protected specification contradiction before required gate` (no observed termination)
- `go test -tags integration -race ./cmd/fid/integration/...` - BLOCKED - `not run; protected specification contradiction before required gate` (no observed termination)
- `ls cmd/fid/integration/testdata/p21_conformance` - BLOCKED - `not run; protected specification contradiction before required gate` (no observed termination)
- `go test -bench . -benchmem ./...` - BLOCKED - `not run; protected specification contradiction before required gate` (no observed termination)
- `go build -o /tmp/fid ./cmd/fid` - BLOCKED - `not run; protected specification contradiction before required gate` (no observed termination)

## Deviations from spec

- No accepted behavior, test assertion, quality-gate command, or public contract was changed. The journey test exposes the contradiction between the fixed section 16 20-Event enumeration and the shipped task-default transition chain.

## Concerns, technical debt, or notes for the Tech Lead

- The active default task workflow requires `todo -> in_progress -> dev_complete -> qa_review -> qa_pass -> done`. After its required one nonterminal transition, the accepted journey cannot move directly to terminal success. Defining/assigning an alternate workflow is also a mutation outside the enumerated 20 Events.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `236540c48a2129d70480b3222ba4122a5471eb48`

## Commit subject

- `test(conformance): expose phase 21 journey guard`
