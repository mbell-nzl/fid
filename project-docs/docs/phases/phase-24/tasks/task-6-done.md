# Phase 24 - Task 6: Gate registration, Envelope and Claim leaves - Recovery Checkpoint

## Workflow outcome

- Task outcome: RECOVERY_REQUIRED
- Category: in-scope-implementation
- Condition: The eight leaves have partial observation and descriptor implementation, but retained Gate/Envelope/Claim fixtures are red and Task 6 parity, refusal, replay, import and Source matrices remain unfinished.
- Recovery attempts: 0
- Next action: Resume Task 6 on this checkpoint. Finish nullable Obligation registration input, converge retained tests and shared integration/fake fixtures with introducing-commit provenance, prove all VAL rows and real CLI behavior, run touched-package race tests, then run all five exact commands in order and complete the report.
- Human attention required: no

OUTCOME RECOVERY_REQUIRED

## Summary of what was implemented

- `gate.register` now computes its recorded project/worktree and binding digest from a repository observation after replay. `claim.acquire` derives its worktree from an observation after replay. `envelope.create` observes after replay, resolves all three Source positions, and checks Registration binding through a shared same-worktree/import predicate. `envelope.complete` drops its caller repository input.
- Adopted descriptors and `--example` for the eight Gate Registration, Envelope and Claim leaves, and added initial focused tests. The implementation and tests are unfinished; this is a checkpoint, not Task 6 acceptance.

## Files created or modified

- `internal/evidence/registration.go`, `internal/workflow/envelope.go`, `internal/workflow/import_binding.go`, `internal/claims/generic.go`
- `cmd/fid/gate.go`, `cmd/fid/envelope.go`, `cmd/fid/claim.go`, `cmd/fid/input_descriptor_command.go`, `cmd/fid/input_descriptors_binding.go`
- `cmd/fid/gate_test.go`, `cmd/fid/envelope_test.go`, `cmd/fid/generic_claim_test.go`, `cmd/fid/input_descriptor_test.go`
- `project-docs/specs/06_DESIGN_PLAN_CLI.md`, `project-docs/specs/07_DESIGN_PLAN_FLOWS.md`, this report

## Documentation updated

- `project-docs/specs/06_DESIGN_PLAN_CLI.md` records the current input shapes and observed bindings.
- `project-docs/specs/07_DESIGN_PLAN_FLOWS.md` records the Gate, Envelope and Claim observation sequence.

## Validation Contract Coverage

- VAL-011..016, VAL-060: PARTIAL - `TestP24GateEnvelopeClaimDescriptors` proves descriptor/type parity and `TestAdoptedExampleProcessExclusivity` passes for adopted leaves; full eight-leaf help, template, structural and side-effect corpus remains.
- VAL-022, VAL-026..028: PARTIAL - `TestGateRegisterObservesRepositoryWithoutCallerIdentity`, `TestEnvelopeCreateObservesScopeAndSourcesWithoutCallerRepository`, and `TestClaimAcquireBindsActiveEnvelopeActorWorktreeScopeAndTTL` pass for original-workspace records. Independent full Repository/Event parity and plain-copy rejection remain.
- VAL-009, VAL-017..019, VAL-023, VAL-029, VAL-035, VAL-066: PENDING - HEAD/replay, removed members, all Source positions, unavailable repositories, imported destinations and fake/fid paired cases still need Task 6 tests and convergence.
- VAL-052: PENDING - final ordinary and tagged gate package timings were not measured on a green tree.

## Tests run

- First red: `go test ./cmd/fid -run '^TestClaimAcquireBindsActiveEnvelopeActorWorktreeScopeAndTTL$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`); omitted `worktree_id` returned the old `invalid claim request` fallback.
- After the Claim implementation, the identical focused command - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	0.634s` (terminated `exit_code: 0`).
- `go test ./cmd/fid -run '^TestGateRegisterObservesRepositoryWithoutCallerIdentity$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`); the old binding validator required caller identity. After the Registration change, the identical command - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	0.475s` (terminated `exit_code: 0`).
- `go test ./cmd/fid -run '^TestP24GateEnvelopeClaimDescriptors$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`); all eight descriptors were behind the pending fence. After adopting them, the identical command - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	0.151s` (terminated `exit_code: 0`).
- `go test ./cmd/fid -run '^TestEnvelopeCreateObservesScopeAndSourcesWithoutCallerRepository$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`); the old Envelope input required `repository_scope`. After the Source/observation change, the identical command - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	0.556s` (terminated `exit_code: 0`).
- `go test ./cmd/fid -run 'TestGateRegister|TestEnvelope|TestClaim' -count=1 -timeout 10m` - FAIL - `FAIL` (terminated `exit_code: 1`); retained fixtures still send removed members, copy path-bound facts, or need nullable `registration`. Preserve this failure and rerun the identical command after repair.
- `go test ./cmd/fid -run '^(TestClaimAcquireBindsActiveEnvelopeActorWorktreeScopeAndTTL|TestGateRegisterObservesRepositoryWithoutCallerIdentity|TestEnvelopeCreateObservesScopeAndSourcesWithoutCallerRepository|TestP24GateEnvelopeClaimDescriptors|TestInputDescriptorRegistryMatchesLiveTreeAndTypes)$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	0.676s` (terminated `exit_code: 0`)
- `go test ./cmd/fid -run '^TestAdoptedExampleProcessExclusivity$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	1.577s` (terminated `exit_code: 0`)
- `go test ./cmd/fid ./internal/workflow ./internal/evidence ./internal/claims ./internal/export -run '^$'` - PASS - `ok  	github.com/mbell-nzl/fid/internal/export	0.005s [no tests to run]` (terminated `exit_code: 0`)
- `go build -o /tmp/fid ./cmd/fid` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `/tmp/fid gate register --example` - PASS - `}` (terminated `exit_code: 0`)

The compile-only command ran no tests. The real `fid --example` invocation ran from `/tmp` and printed the new identity-free input template.

## Exact quality-gate results

Working directory for the required gate: `/home/mark/projects/fid`. None of its commands has been run on a final tree. These are pending commands, not passing evidence.

- `gofmt -l cmd internal` - BLOCKED - not run; Task 6 retained tests and implementation remain unfinished (no termination).
- `golangci-lint run ./...` - BLOCKED - not run; Task 6 retained tests and implementation remain unfinished (no termination).
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` - BLOCKED - not run; retained Gate/Envelope/Claim tests are red (no termination).
- `go tool cover -func=coverage.out` - BLOCKED - not run; no final coverage profile exists (no termination).
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...` - BLOCKED - not run; Task 6 integration convergence is unfinished (no termination).

## Deviations from spec

- No protected contract or gate definition was changed. The retained focused group ended `FAIL` (terminated `exit_code: 1`) because its inputs still contain removed repository/worktree members, its path-bound fixture caches need replacement, and the descriptor rejects a previously valid nullable Obligation Registration. In-scope convergence is pending; the identical command must be rerun after repair. No required gate attempt has been made.
- The first report-lint check ended `FAIL report-lint found 5 failure(s)` (terminated `exit_code: 5`) because five passing test/build rows had prose or punctuation after their termination field. The report correction moved that prose to a separate paragraph and left the observed command tokens and statuses intact.

## Concerns, technical debt, or notes for the Tech Lead

- Retained files already changed: `gate_test.go` was introduced by `0c4c04ea` (Phase 21 Gate task), `envelope_test.go` by `56dad8fa` (Phase 21 Envelope task), `generic_claim_test.go` by `0b67c830` (Phase 21 Claims task), and `input_descriptor_test.go` by `04b6b099` (Phase 24 task 3). The accepted Phase 24 Task 6 input and observation changes, rather than historical names, require convergence. Record introducing commits for any further retained fixture changes before editing.
- Preserve the negative binding matrix while replacing caller `worktree_id` manipulation with a real plain-copy worktree. Do not remove the wrong-worktree assertion. Preserve the original-workspace positive control, and add imported-destination controls and no-Event checks.
- The current descriptor's `registration` field needs the accepted nullable shape for review/Human Obligations. Its retained failure is real; changing those inputs to fake Registration objects would weaken coverage.
- The current worktree has Runner-owned `project-docs/docs/phases/phase-24/triage-log.md` and untracked `run-ledger.md` changes. Leave them untouched during continuation and commits.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `3ed4824b3f9f7d5dee0f2a50799df538a2911953`

## Commit subject

- `feat(binding): checkpoint observed gate envelope claim inputs`
