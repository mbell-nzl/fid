# Phase 21 - Task 22: Final Wire Conformance And Real-Process Journey - Completion Report

## Workflow outcome

- Task outcome: RECOVERY_REQUIRED
- Category: gate-repair
- Condition: The DEC-174 authored-workflow journey reaches the exact 22-Event contract check, where the real binary emits `activation.finalizeted` instead of the accepted `activation.finalized`. The emitting writer is `internal/workflow/activation.go`, outside this task's named source scope.
- Recovery attempts: 0
- Next action: Obtain bounded authority that permits the exact Event-name repair while preserving the accepted Event schema, public interface, journey assertion, command block, and all failed evidence; then rerun the byte-identical tagged integration gate and complete the remaining post-rebuild assertions.
- Human attention required: no

## Summary of what was implemented

- Preserved the committed `p21_conformance` corpus and attempt-1 journey guard, then extended the same real-binary journey to author the DEC-174 task workflow before Issue creation: `workflow define --applies-to task`, followed by `workflow assign --kind task`.
- The new journey uses exactly one `unstarted`, `started`, `completed`, and `cancelled` state; declares only `unstarted -> started` and `started -> completed`; transitions through Gate execution, Receipt Evidence, terminal completion, Run finish, Claim release, and Envelope completion.
- Added an explicit VAL-085 undeclared-edge negative against non-cancelled states. It requires `transition-not-permitted`, `workflow-transition-invalid`, process exit 4, and a zero Event delta before and after rebuild.
- Added exact ordered 22-Event-family assertions, single `workflow.defined`/`workflow.assigned`/`envelope.created` checks, and the allowed post-rebuild read/evaluation/completion-state/resume/export/negative sequence. The first exact-family assertion exposes the activation Event-name defect before the post-rebuild sequence can run.
- Historical attempt-1 BLOCKED evidence remains immutable at `c964b1ffbfd88bc9d6bf4abc04fa3d5679fb99b3`: its terminal `in_progress -> done` request returned `transition-not-permitted` / `workflow-transition-invalid` at process exit 4 under the seeded task-default workflow. DEC-174 replaced the planning artifact, not that evidence.

## Files created or modified

- `cmd/fid/integration/p21_conformance_test.go`
- `project-docs/docs/phases/phase-21/tasks/task-22-done.md`

## Documentation updated

- `project-docs/docs/phases/phase-21/tasks/task-22-done.md`

## Validation Contract Coverage

- VAL-072: PASS - `TestPhase21ConformanceCorpusIsCompleteAndByteExact` in `cmd/fid/integration/p21_conformance_test.go:25` regenerated and byte-compared the committed operation, failure, and live-fixture corpus in this attempt's tagged integration runs.
- VAL-073: PASS - `TestPhase21ConformanceCorpusIsCompleteAndByteExact` at `cmd/fid/integration/p21_conformance_test.go:55` requires every extracted public closed member, reports no missing member, and rejects non-kebab/underscore values before the independent journey test runs.
- VAL-084: FAIL - `TestPhase21RealProcessJourneyStartsOneActivatedEnvelope` reaches real Gate execution, Receipt linking/evaluation, terminal `started -> completed`, Run finish, Claim release, and Envelope completion, but fails at `cmd/fid/integration/p21_conformance_test.go:139` because Event 9 is `activation.finalizeted`, not the accepted `activation.finalized`.
- VAL-085: FAIL - `p21AssertUndeclaredTransitionRejected` at `cmd/fid/integration/p21_conformance_test.go:185` exercised the required non-cancelled undeclared edge with exit 4, `transition-not-permitted`, `workflow-transition-invalid`, and no Event append, but the complete negative matrix and its post-rebuild counterpart remain pending the exact Event-name repair.
- VAL-086: FAIL - `p21AssertJournalTypes` at `cmd/fid/integration/p21_conformance_test.go:139` observed 22 Events but rejected the wrong ninth Event family. The permitted post-rebuild zero-delta sequence is implemented but was not reached by the red test run.

## Tests run

- `go test -tags integration ./cmd/fid/integration/...` - FAIL - `FAIL` (terminated `exit_code: 1`)
  Red evidence after adding DEC-174 authoring but before updating the stale guard target: `p21_conformance_test.go:119` observed `unstarted -> in_progress` rejecting at exit 4 with `transition-not-permitted` and `workflow-transition-invalid`.
- `go test -tags integration ./cmd/fid/integration/...` - FAIL - `FAIL` (terminated `exit_code: 1`)
  Later result after updating the states, adding the negative, and extending the completion/rebuild journey: `p21_conformance_test.go:139` observed Event 9 as `activation.finalizeted`, with the required `activation.finalized` assertion preserved.
- `/tmp/fid --project /tmp/fid-p21-entrypoint.Cc2mai --json init --prefix P21` - PASS - `{"contract_version":"fid-cli.v1","event":{"id":"EVT-20260903-000001","occurred_at":"2026-09-03T04:45:31.177530755Z"},"ok":true,"operation":"project.init","result":{"journal_path":".fid/journal/2026-09-03.jsonl","project":{"id":"project-33ba0ca1de7a529946cc347d9ac990ae","issue_prefix":"P21","root":"/tmp/fid-p21-entrypoint.Cc2mai"}}}` (terminated `exit_code: 0`)
- `/tmp/fid --project /tmp/fid-p21-entrypoint.Cc2mai --json workflow define --name p21-entrypoint --applies-to task --states unstarted,started,completed,cancelled --categories unstarted,started,completed,cancelled --transitions 'unstarted->started,started->completed' --actor agent:developer --idempotency-key p21-entrypoint-define` - PASS - `{"contract_version":"fid-cli.v1","event":{"id":"EVT-20260903-000002","occurred_at":"2026-09-03T04:45:38.717895356Z"},"ok":true,"operation":"workflow.define","result":{"idempotency_key":"p21-entrypoint-define","replayed":false,"workflow":{"applies_to":"task","id":"p21-entrypoint","name":"p21-entrypoint","state_count":4,"transition_count":5}}}` (terminated `exit_code: 0`)
- `/tmp/fid --project /tmp/fid-p21-entrypoint.Cc2mai --json workflow assign --name p21-entrypoint --kind task --actor agent:developer --idempotency-key p21-entrypoint-assign` - PASS - `{"contract_version":"fid-cli.v1","event":{"id":"EVT-20260903-000003","occurred_at":"2026-09-03T04:45:44.062669497Z"},"ok":true,"operation":"workflow.assign","result":{"assignment":{"scope_type":"kind","scope_value":"task","workflow_id":"p21-entrypoint"},"idempotency_key":"p21-entrypoint-assign","replayed":false}}` (terminated `exit_code: 0`)

## Exact quality-gate results

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -coverprofile=coverage.out ./...` - PASS - `ok  \tgithub.com/mbell-nzl/fid/internal/workflow\t(cached)\tcoverage: 27.5% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total:\t\t\t\t\t\t\t\t\t(statements)\t\t\t61.5%` (terminated `exit_code: 0`)
- `go test -tags integration -race ./cmd/fid/integration/...` - BLOCKED - `FAIL` (terminated `exit_code: 1`)
  Failure detail: `p21_conformance_test.go:139: journal Event 9 = "activation.finalizeted", want "activation.finalized"`.
- `ls cmd/fid/integration/testdata/p21_conformance` - PASS - `operations` (terminated `exit_code: 0`)
- `go test -bench . -benchmem ./...` - PASS - `ok  \tgithub.com/mbell-nzl/fid/internal/workflow\t2.334s` (terminated `exit_code: 0`)
- `go build -o /tmp/fid ./cmd/fid` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)

## Deviations from spec

- Required command failure retained: the first focused red result ended `FAIL` (terminated `exit_code: 1`) because the retained guard still requested `unstarted -> in_progress` after the DEC-174 workflow authoring change. The repair changed that test-only target to `unstarted -> started` and added the explicit undeclared-edge rejection; the later tagged result still ended `FAIL` (terminated `exit_code: 1`) at the independent `activation.finalizeted` versus `activation.finalized` mismatch. No command, expected exit, assertion strength, public contract, Event schema, seeded workflow, or corpus was changed.

## Concerns, technical debt, or notes for the Tech Lead

- The exact Event-name repair is uniquely determined by `public-interface.md` section 4.1 and `contracts.md` sections 14 and 16, but its writer is `internal/workflow/activation.go`, which task 22 cannot modify under its current scope. This is a scope gap requiring bounded recovery authority, not a request to accept or reinterpret the misspelled Event.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `a6b3a448aef6fcaab5aed3f3e46dad1df6b2cb1d`

## Commit subject

- `test(conformance): extend phase 21 journey workflow`
