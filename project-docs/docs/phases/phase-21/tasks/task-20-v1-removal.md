# Phase 21 - Task 20: V1 Removal And Migration 0014

## Objective

After all replacements pass, execute the exhaustive convergence inventory: remove owned V1 policy/dev/completion/Claim/Run/Gate surfaces and writers, apply migration 0014, and leave only explicitly frozen Phase 23/25 state.

## Validation Contract Coverage

- VAL-002: owning observable in the Phase 21 validation contract.
- VAL-074: owning observable in the Phase 21 validation contract.
- VAL-075: owning observable in the Phase 21 validation contract.

## Binding Constraints

- Use `v1-convergence-inventory.md` as protected removal authority; classify every live match.
- Migration 0014 drops `roles`, V1 `claims`, and `gate_runs`; renames
  `work_claims` to final `claims`; and deletes the displaced role/judgment
  Issue rows from `workflow_transition_requires`. No data reconstruction or
  compatibility alias is permitted.
- Delete the obsolete `internal/dev/` and `internal/runs/` packages only after
  their trusted Claim/Run/Receipt/Evidence replacements have passed.
- Remove role/judgment workflow rows but do not drop frozen `task_completions`, `agent_runs`, `ac_verifications`, or `concerns`.
- Do not weaken a protected absence command to accommodate residue.

## Scope

### Files To Create Or Modify

- `cmd/fid/`
- `internal/dev/`
- `internal/runs/`
- `internal/policy/`
- `internal/claims/`
- `internal/issues/`
- `internal/workflow/`
- `internal/maintenance/`
- `migrations/0014_p21_retire_v1_authority_execution.up.sql`
- `migrations/0014_p21_retire_v1_authority_execution.down.sql`
- `project-docs/docs/phases/phase-21/tasks/task-20-done.md`

### Files To Read For Context

- `project-docs/docs/phases/phase-21/v1-convergence-inventory.md`
- `project-docs/specs/03_TRD_SCHEMA.md §8.2`
- `project-docs/docs/phases/phase-21/public-interface.md §5`

### Out Of Scope

- Final-schema maintenance/rebuild/export parity and Phase 23/25 frozen-table deletion.

## Acceptance Criteria

- [ ] Migration 0014 is exact and empty-state reversible. (VAL-002)
- [ ] Every retained/fail-closed/frozen caller has its tested final disposition. (VAL-074)
- [ ] Owned V1 paths, writers, Events, tables, commands, and `--by` are absent. (VAL-075)
- [ ] The task-2 nonpersistent V1 compatibility role vocabulary and every
  fallback reader are deleted with `internal/policy`; no alternate fixed-role
  registry remains. (VAL-075)
- [ ] The completion report runs all four `v1-convergence-inventory.md` §7 recipes, classifies every match against the frozen/diagnostic-reader or removal disposition, and records `UNCLASSIFIED=0`; the first two removal searches have zero matches while the broad table-reader and Event-family inventories remain successful, classified inventories. (VAL-074, VAL-075)

## First Failing Test

Add failing schema/source/CLI/Event/writer guards before deleting each listed surface.

## Documentation Updates

- `project-docs/docs/phases/phase-21/tasks/task-20-done.md` - record red/green evidence, deviations, required commands, and spec impact.

## Exact Commands To Run

Run exactly:

```text
go test -race ./internal/db ./internal/maintenance ./internal/cli ./cmd/fid/...
ls internal/dev
ls internal/runs
ls cmd/fid/dev.go
ls cmd/fid/policy.go
ls cmd/fid/issue_complete.go
rg -n --glob '*.go' --glob '!**/*_test.go' 'String(Var)?\([^\n]*"by"|GetString\("by"|--by' cmd/fid internal
rg -n --glob '*.go' --glob '!**/*_test.go' 'policy\.Authorize|KnownRoles|ActorRole|Role:' cmd/fid internal
rg -n --glob '*.go' --glob '!**/*_test.go' '(FROM|INTO|UPDATE|JOIN|DELETE FROM) (roles|claims|agent_runs|task_completions|gate_runs|ac_verifications|concerns|workflow_transition_requires)\b|Table:\s*"(roles|claims|agent_runs|task_completions|gate_runs|ac_verifications|concerns|workflow_transition_requires)"' cmd/fid internal
rg -n --glob '*.go' --glob '!**/*_test.go' 'EventType:\s*"' cmd/fid internal
rg -n --glob '*.go' --glob '!**/*_test.go' '(INTO|UPDATE|DELETE FROM) (roles|gate_runs|task_completions|agent_runs|ac_verifications|concerns)\b|Table:\s*"(roles|gate_runs|task_completions|agent_runs|ac_verifications|concerns)"' cmd/fid internal
rg -n 'EventType:\s*"(dev\.|issue\.completed|issue\.claim\.|run\.recorded)|"event_type":"(dev\.|issue\.completed|issue\.claim\.|run\.recorded)' cmd/fid internal testdata
rg -n --glob '*.go' --glob '!**/*_test.go' '"event"\s*:\s*"[a-z0-9]+_[a-z0-9_]+"|Event:\s*"[a-z0-9]+_[a-z0-9_]+"' cmd/fid internal
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

Every command is required and must terminate with session-captured evidence. Report lint is not a Developer task-gate command; after the report commit, the Phase Runner runs the phase brief's two-argument pre-dispatch control. A Ready QA plan may add a task-specific command; the Tech Lead must amend this section before authorization if it does.

## Expected Non-Zero Commands

- `ls internal/dev` - expected `exit_code: 2` - obsolete dev package absent
- `ls internal/runs` - expected `exit_code: 2` - obsolete legacy Runs package absent
- `ls cmd/fid/dev.go` - expected `exit_code: 2` - dev CLI absent
- `ls cmd/fid/policy.go` - expected `exit_code: 2` - V1 policy CLI absent
- `ls cmd/fid/issue_complete.go` - expected `exit_code: 2` - legacy completion CLI absent
- `rg -n --glob '*.go' --glob '!**/*_test.go' 'String(Var)?\([^\n]*"by"|GetString\("by"|--by' cmd/fid internal` - expected `exit_code: 1` - no production mutation/grouping `--by`
- `rg -n --glob '*.go' --glob '!**/*_test.go' 'policy\.Authorize|KnownRoles|ActorRole|Role:' cmd/fid internal` - expected `exit_code: 1` - no V1 policy caller or Actor-role field
- `rg -n --glob '*.go' --glob '!**/*_test.go' '(INTO|UPDATE|DELETE FROM) (roles|gate_runs|task_completions|agent_runs|ac_verifications|concerns)\b|Table:\s*"(roles|gate_runs|task_completions|agent_runs|ac_verifications|concerns)"' cmd/fid internal` - expected `exit_code: 1` - no removed/frozen-state production writer
- `rg -n 'EventType:\s*"(dev\.|issue\.completed|issue\.claim\.|run\.recorded)|"event_type":"(dev\.|issue\.completed|issue\.claim\.|run\.recorded)' cmd/fid internal testdata` - expected `exit_code: 1` - no accepted V1 execution/completion Event fixture
- `rg -n --glob '*.go' --glob '!**/*_test.go' '"event"\s*:\s*"[a-z0-9]+_[a-z0-9_]+"|Event:\s*"[a-z0-9]+_[a-z0-9_]+"' cmd/fid internal` - expected `exit_code: 1` - no production snake-case Activity Event value

## Notes For Developer

- Start only after tasks 1-19 close.
- The two broad §7 inventories are expected to terminate `exit_code: 0`; their matches are evidence only when every line is classified and the completion report records `UNCLASSIFIED=0`. They are not absence controls and must not be narrowed to suppress frozen or diagnostic readers.
- Completion report: `project-docs/docs/phases/phase-21/tasks/task-20-done.md`.
- Stop for a protected contract change; ordinary in-scope red/green repair remains this task.
- Do not read held-out scenario material.
