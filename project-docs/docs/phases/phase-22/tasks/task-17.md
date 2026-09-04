# Phase 22 - Task 17: Maintenance, capability qualification, and rebuild budget

## Objective

Detect recovery-state inconsistencies without repair, enable only completed capabilities, and prove the final schema's rebuild performance. Keep incomplete Claims capability explicitly unavailable.

## Validation Contract Coverage

- VAL-069: Maintenance detects recovery corruption without repairing it.
- VAL-070: Only fully qualified capabilities are advertised available.
- VAL-071: 100k Event rebuild meets the unchanged reference target in three fresh samples.
- VAL-072: Final schema 16 initializes and rebuilds without retired authority.

Full sources, pass/fail criteria, and validator ownership are in
`project-docs/docs/phases/phase-22/qa/validation-contract.md`.

## Binding Constraints

- PROPOSED: execution requires Human acceptance of DEC-180 and a Ready QA plan.
- contracts.md section 11; public-interface.md section 5; CF-024.
- CF-024: retain the exact reference instance, three consecutive fresh processes, and per-sample <30-second self-assertion.
- Do not infer liveness from lease expiry, repair state in maintenance, or advertise work.claims.v1 complete.
- Every available capability needs the same fake/fid subset cases assembled by earlier tasks; task 18 independently audits coverage.
- `execution-scope.md` is incorporated into this scope: it assigns only
  subject-preserving convergence of assertions displaced by this task's accepted
  change. All owned ordinary and tagged code executes under the commands below.

## Scope

### Files To Create Or Modify

- `internal/maintenance/`
- `internal/workflow/capabilities.go`
- `internal/export/`
- `internal/bootstrap/`
- `internal/db/`
- `internal/journal/`
- `cmd/fid/capabilities.go`
- `cmd/fid/maintenance.go`
- `cmd/fid/rebuild.go`
- `cmd/fid/integration/`
- Matching in-package tests and `cmd/fid/integration/` contract fixtures for
  the behavior above, including retained rows invalidated by this change.
- `project-docs/docs/phases/phase-22/tasks/task-17-done.md`.

### Files To Read For Context

- `project-docs/docs/phases/phase-22/contracts.md` - only the sections named above.
- `project-docs/docs/phases/phase-22/public-interface.md` - owned input/result/error types.
- `project-docs/docs/phases/phase-22/schema.md` - owned DDL and Event rows.
- `project-docs/docs/phases/phase-22/execution-scope.md` and `v1-convergence-inventory.md`.
- Canonical source sections cited by this task's validation rows.
- `project-docs/docs/phases/phase-22/tasks/task-16-done.md` - actual predecessor hand-back.

### Out Of Scope

- Threshold change, broader telemetry, parallel Claims, security-tool selection for release, and release readiness.
- Product policy, protected gate definitions, workflow controls, dependencies,
  and unrelated implementation changes beyond this accepted task.

## Acceptance Criteria

- [ ] Maintenance detects recovery corruption without repairing it. (VAL-069)
- [ ] Only fully qualified capabilities are advertised available. (VAL-070)
- [ ] 100k Event rebuild meets the unchanged reference target in three fresh samples. (VAL-071)
- [ ] Final schema 16 initializes and rebuilds without retired authority. (VAL-072)
- [ ] Full required commands terminate; preserve red/green and failed-attempt
  evidence with output tokens and observed exit status.
- [ ] Every converged retained row identifies introducing commit/phase/task,
  invalidating change, and the exact command that executes it.

## First Failing Test

TestMaintenanceReportsBrokenRecoveryLinkWithoutMutation fails before recovery integrity findings exist.
Capture the failure before product edits; then implement the smallest behavior
that satisfies it. Build failures alone do not prove the intended rejection.

## Documentation Updates

- `project-docs/docs/phases/phase-22/tasks/task-17-done.md`: implementation,
  red/green, every exact command, observed termination, assertion evidence,
  convergence provenance, dependencies, and deviations.
- Flag spec impact against the named contracts; only Tech Lead accepts a
  contract/schema change. Record no spec change when implementation conforms.
- Preserve CF-024's original reference identity in the report: Go 1.26.2,
  linux/amd64, Linux 6.6.114.1-microsoft-standard-WSL2 x86_64, six logical CPUs,
  AMD Ryzen 7 5700G with Radeon Graphics, MemTotal 8133236 kB.

## Exact Commands To Run

Working directory: repository root `/home/mark/projects/fid`, or the isolated
same-task worktree root explicitly recorded by the Runner. No command is a
placeholder. Apply gofmt/goimports to edited Go files before the final gate.

After implementation, every command below is required before task close.
The integration gate is mandatory from task 1, including for retained fixtures.

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
go test -tags integration -race ./cmd/fid/integration/...
go version
go env GOOS GOARCH
uname -srm
nproc
grep -m1 '^model name' /proc/cpuinfo
grep -m1 '^MemTotal:' /proc/meminfo
go test -run '^$' -bench '^BenchmarkRebuildWorkspace100kEvents$' -benchmem -benchtime 1x -timeout 30m ./internal/journal
go test -run '^$' -bench '^BenchmarkRebuildWorkspace100kEvents$' -benchmem -benchtime 1x -timeout 30m ./internal/journal
go test -run '^$' -bench '^BenchmarkRebuildWorkspace100kEvents$' -benchmem -benchtime 1x -timeout 30m ./internal/journal
go test -bench . -benchmem ./...
go build -o /tmp/fid ./cmd/fid
```

No exception to the full task gate. Commands whose expected pass is nonzero:
none. Focused red/green selections do not replace these commands. Report each
attempt using the command's own output token and observed termination, as
required by AGENTS.md. No interrupted/hung command is PASS; preserve all failed
attempts and rerun the byte-identical command after authorized causal repair.
Report lint is a Runner control after the report commit, not a Developer gate.

## Notes For Developer

- Start after task 16 closes.
- One task per Developer session; recovery may continue the same task.
- This task belongs to QA Review window 6.
- Commit using Conventional Commits with `Refs: phase-22/task-17`.
- Never infer ownership from a historical filename or task label. No-progress
  budgets exclude interruptions and attempts forbidden from repairing the cause.
- Do not read held-out scenario material or activate another agent.

