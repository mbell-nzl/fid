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
- Added by DEC-173 (bounded to the dispositions below): `internal/bootstrap/`,
  `internal/journal/`, `internal/labels/`, `internal/links/`,
  `internal/requirements/`, `internal/db/`, `testdata/journal/`,
  `cmd/fid/integration/`, and test-convergence only in `internal/activity/`,
  `internal/stats/`, and `internal/sqlquery/`.

### Scope Amendment And Convergence Ownership (DEC-162, DEC-165, DEC-167, DEC-173)

This task owns every path its required absence and inventory controls walk
(`cmd/fid`, `internal`, `testdata`) and every package its gates compile.
Convergence ownership follows the invalidating accepted change: migration 0014,
each V1 leaf/package removal, and the `ActorRole` deletion are task-20 changes,
so task 20 converges every assertion they invalidate wherever it lives.

Production edits in the DEC-173 packages are limited to these dispositions:

- `internal/journal/lock.go`: delete `LockOwner.ActorRole` and the
  `actor_role` lock-metadata field; keep `pid`, `hostname`, `started_at`, and
  `command`; every `LockOwner{...}` call site drops the field.
- `internal/bootstrap/init.go`: drop the `ActorRole: "human"` lock value and
  delete `RoleSeed`, `DefaultRoles`, and `SeededRoleNames`. Do not change
  `BuildProjectInitializedEvent`, the VAL-008 fixture, or init/rebuild
  behaviour.
- Retained families in `internal/labels/`, `internal/links/`,
  `internal/requirements/`, `internal/workflow/`, and `internal/issues/`:
  delete the `ActorRole` option field, the `legacyRole` fallback inside
  `standaloneActorForOptions`, and the `legacy-in-process-N` idempotency-key
  synthesis. `Actor` plus `IdempotencyKey` are the only inputs. The public
  `--actor`/`--idempotency-key` grammar, Event-v2 bytes, digests, replay, error
  envelopes, and every VAL-028..031 assertion keep full strength; in-process
  test callers converge to `Actor: "agent:<role>"` plus an explicit key.

Fixture and test convergence, applying the §1.2 preference order:

- `testdata/journal/p20-issue-revision-effects.golden.jsonl` with
  `internal/journal/replay_test.go`: the three `dev.*` Events become retained
  Issue families carrying the same revision Effects; `wantTypes` and
  `priorByType` change by the same substitution only. The ten-Event revision
  chain is the preserved subject.
- `roles` seeds and assertions in `testdata/journal/replay-v1.golden.jsonl`,
  `testdata/journal/corrupt-tail.jsonl`, `internal/journal/{replay,project,
  phase2_rebuild,rebuild_scale,writer}_test.go`, `internal/db/*_test.go`,
  `internal/sqlquery/sqlquery_test.go`, and `cmd/fid/*_test.go`: converge to a
  surviving table with identical assertion strength; the scale fixture keeps
  its Event count and shape. Adding a projection branch that ignores, skips, or
  tolerates an Effect on a dropped table is a data-model decision and is
  forbidden. If a fixture's subject genuinely requires a legacy `roles` Effect
  against the final schema, stop and escalate.
- `cmd/fid/integration/` closure journey: removed-behaviour steps
  (`issue complete`, `run record`, `policy show`, `has_completion_report`,
  dropped/frozen table snapshots) become their accepted behaviour - unknown
  command, journal unchanged, table absent or empty - and every retained
  assertion stays byte-identical. Delete a step only where a named surviving
  test owns the accepted equivalent, and name both in the report.
- `t.Skip`, build tags, `-run` narrowing, threshold changes, and silent
  assertion loss are never convergence.

Coordination: task 9 is complete and its public contract is frozen. Task 21
retains final-schema init/rebuild/maintenance/export parity and the 100k
benchmark; do not change projection or rebuild semantics, maintenance findings,
export, or the benchmark threshold/shape.

The completion report records per-row provenance - file, line, invalidating
change, accepted replacement - for every converged row outside the original
file list, and `UNCLASSIFIED=0` for both broad inventories.

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
go test -tags integration -race ./cmd/fid/integration/...
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

Every command is required and must terminate with session-captured evidence. Report lint is not a Developer task-gate command; after the report commit, the Phase Runner runs the phase brief's two-argument pre-dispatch control. A Ready QA plan may add a task-specific command; the Tech Lead must amend this section before authorization if it does. The tagged integration command was added by DEC-173 under DEC-167 and runs under the DEC-168 300 s budget; commit durable work before the long commands (DEC-170).

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
- DEC-173 continuation (attempt 2) starts from checkpoint `aef03cf`; the attempt-1 failed rows stay in the report as immutable evidence and every command is rerun byte-identically.
- Stop for a protected contract change; ordinary in-scope red/green repair remains this task.
- Do not read held-out scenario material.
