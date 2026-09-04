# Phase 21 - Task 21: Final-Schema Maintenance, Init, Rebuild, And Export Parity

## Objective

Against the post-0014 binary, prove fresh init and delete-projection/rebuild are role-free, extend deterministic maintenance findings, compare all named live/rebuilt reads and export bytes, and run the 100k Event-v2 check.

## Validation Contract Coverage

- VAL-076: owning observable in the Phase 21 validation contract.
- VAL-077: owning observable in the Phase 21 validation contract.
- VAL-078: owning observable in the Phase 21 validation contract.
- VAL-079: owning observable in the Phase 21 validation contract.
- VAL-080: owning observable in the Phase 21 validation contract.
- VAL-081: owning observable in the Phase 21 validation contract.
- VAL-082: owning observable in the Phase 21 validation contract.
- VAL-083: owning observable in the Phase 21 validation contract.

## Binding Constraints

- This task starts only after destructive removal; it cannot use a pre-0014 bootstrap fixture.
- Maintenance reports but never repairs.
- Record the DEC-158 reference identity exactly: `go1.26.2 linux/amd64`, `Linux 6.6.114.1-microsoft-standard-WSL2 x86_64`, 6 logical processors, `AMD Ryzen 7 5700G with Radeon Graphics`, and `8133236 kB` `MemTotal`.
- Run the bounded 100k benchmark in three consecutive fresh processes. Every sample must preserve parity and terminate under 30 seconds; best/median sampling is forbidden, and the benchmark's per-sample `<30s` `b.Fatalf` remains part of the general phase gate.
- If any sample fails, optimize the implementation within task scope and rerun the byte-identical three-sample sequence; do not weaken, average, or rebaseline the protected threshold.
- DEC-173 coordination: task 20 deleted `LockOwner.ActorRole`/`actor_role`, the bootstrap `RoleSeed`/`DefaultRoles`/`SeededRoleNames` registry, and every `roles` seed in the journal goldens and scale fixture. Do not reintroduce a role lock value, role registry, role fixture, or any projection tolerance for Effects on dropped tables; prove the final state on the fixtures task 20 left.

## Scope

### Files To Create Or Modify

- `internal/bootstrap/`
- `internal/maintenance/`
- `internal/journal/`
- `internal/db/`
- `internal/export/`
- `internal/cli/`
- `cmd/fid/`
- `cmd/fid/integration/`
- `project-docs/docs/phases/phase-21/tasks/task-21-done.md`

### Files To Read For Context

- `project-docs/docs/phases/phase-21/contracts.md §§14-16`
- `project-docs/docs/phases/phase-21/v1-convergence-inventory.md`
- `project-docs/specs/03_TRD_SCHEMA.md §8.2`

### Out Of Scope

- Import, repair, held-out scenarios, and final conformance journey.

## Acceptance Criteria

- [ ] Final-binary fresh init succeeds without roles. (VAL-076)
- [ ] Delete/rebuild reproduces role-free bootstrap/capabilities. (VAL-077)
- [ ] All three atomic maintenance classes are exact. (VAL-078, VAL-079, VAL-080)
- [ ] Named reads and export are live/rebuild identical. (VAL-081, VAL-082)
- [ ] On the recorded DEC-158 reference instance, three consecutive fresh-process 100k Event-v2 rebuild samples each preserve parity and terminate under 30 seconds. (VAL-083)

## First Failing Test

Add the final-binary temp-project init/rebuild test first, then corrupt fixtures, parity snapshot, and bounded benchmark.

## Documentation Updates

- `project-docs/docs/phases/phase-21/tasks/task-21-done.md` - record red/green evidence, deviations, required commands, and spec impact.

## Exact Commands To Run

Run exactly:

```text
go test -race ./internal/bootstrap ./internal/maintenance ./internal/journal ./internal/db ./internal/export ./internal/cli ./cmd/fid/...
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
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

Every command is required and must terminate with session-captured evidence. Report lint is not a Developer task-gate command; after the report commit, the Phase Runner runs the phase brief's two-argument pre-dispatch control. A Ready QA plan may add a task-specific command; the Tech Lead must amend this section before authorization if it does.

## Notes For Developer

- Start after task 20 closes. QA Review window 7 follows tasks 19-21.
- Completion report: `project-docs/docs/phases/phase-21/tasks/task-21-done.md`.
- Stop for a protected contract change; ordinary in-scope red/green repair remains this task.
- Do not read held-out scenario material.
