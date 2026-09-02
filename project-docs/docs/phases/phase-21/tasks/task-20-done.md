# Phase 21 - Task 20: V1 Removal And Migration 0014 - Completion Report

## Workflow outcome

- Task outcome: RECOVERY_REQUIRED
- Category: in-scope-implementation
- Condition: DEC-173 convergence is checkpointed through retained Actor attribution, lock metadata, retired workflow requirement Effects, and the `--by` reader. The attempt-2 package and repository gates remain red because retained V1 test fixtures still compile against removed roles, legacy Claim/Run, and dropped final-schema columns.
- Recovery attempts: 2
- Next action: converge the retained V1 test fixtures and final-schema readers in the amended source scope, preserve their assertion strength, then rerun the complete byte-identical command block.
- Human attention required: no

## Summary of what was implemented

- Added migration 0014: it drops `roles`, V1 `claims`, and `gate_runs`; renames `work_claims` to final `claims`; preserves frozen Phase 23/25 tables; and removes only seeded Issue role/judgment requirements.
- Added red/green migration and command-tree removal guards, removed the V1 dev, completion, policy, legacy Claim, and legacy Run production implementations, and retained generic Claim/Run command leaves.
- Retargeted generic Claim/workflow production queries and Effects to the final `claims` table.
- Removed the last production `--by` reader and retired `workflow_transition_requires` emission. Retained workflow definition requests now fail closed when they contain retired transition requirements; generic Claim still closes its projection handles after the V1 shared helper was removed.

## Files created or modified

- `migrations/0014_p21_retire_v1_authority_execution.up.sql`
- `migrations/0014_p21_retire_v1_authority_execution.down.sql`
- `internal/db/operational_kernel_migration_test.go`
- `cmd/fid/`, `internal/claims/`, `internal/dev/`, `internal/issues/`, `internal/policy/`, `internal/runs/`, and `internal/workflow/`

## Documentation updated

- `project-docs/docs/phases/phase-21/tasks/task-20-done.md`

## Validation Contract Coverage

- VAL-002: PARTIAL - `TestMigration0014RetiresV1AuthorityExecutionWithoutReconstruction` at `internal/db/operational_kernel_migration_test.go:43` first failed because 0014 was absent, then passed after the exact migration files were added; repository-wide migration regressions remain red.
- VAL-074: PARTIAL - the fresh broad table/Event inventories below are classified line by line with `UNCLASSIFIED=0`, but retained frozen/diagnostic dispositions are not yet all converged.
- VAL-075: PARTIAL - `TestTask20V1CommandsAreAbsent` at `cmd/fid/task20_removal_test.go:5` first failed with `fid [dev] resolved command "fid dev"`; source leaves are removed, but the declared absence searches still find V1 residue.

## Tests run

- `go test -race ./internal/db -run '^TestMigration0014RetiresV1AuthorityExecutionWithoutReconstruction$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`; red: `Migrate(14) error = no migration found for version 14: read down for version 14 .: file does not exist`)
- `go test -race ./internal/db -run '^TestMigration0014RetiresV1AuthorityExecutionWithoutReconstruction$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/internal/db	1.677s` (terminated `exit_code: 0`)
- `go test -race ./cmd/fid -run '^TestTask20V1CommandsAreAbsent$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`; red: `fid [dev] resolved command "fid dev" with err <nil>, want unknown command`)
- Real entrypoint exercise: BLOCKED - the exact required package command did not build `cmd/fid`, so no binary could be truthfully exercised after the removal checkpoint.

## Exact quality-gate results

- `go test -race ./internal/db ./internal/maintenance ./internal/cli ./cmd/fid/...` - BLOCKED - `FAIL` (terminated `exit_code: 1`; generic `claims` has exposed retained V1 readers/tests and deleted-package test fixtures.)
- `ls internal/dev` - BLOCKED - `start_test.go` (expected `exit_code: 2`; terminated `exit_code: 0`; source is removed but V1 tests leave the directory present.)
- `ls internal/runs` - BLOCKED - `record_test.go` (expected `exit_code: 2`; terminated `exit_code: 0`; source is removed but the V1 test leaves the directory present.)
- `ls cmd/fid/dev.go` - PASS - `ls: cannot access 'cmd/fid/dev.go': No such file or directory` (expected `exit_code: 2`; terminated `exit_code: 2`)
- `ls cmd/fid/policy.go` - PASS - `ls: cannot access 'cmd/fid/policy.go': No such file or directory` (expected `exit_code: 2`; terminated `exit_code: 2`)
- `ls cmd/fid/issue_complete.go` - PASS - `ls: cannot access 'cmd/fid/issue_complete.go': No such file or directory` (expected `exit_code: 2`; terminated `exit_code: 2`)
- `rg -n --glob '*.go' --glob '!**/*_test.go' 'String(Var)?\([^\n]*"by"|GetString\("by"|--by' cmd/fid internal` - BLOCKED - `internal/issues/ac_verify.go:264:\t\t"--by", actorRole,` (expected `exit_code: 1`; terminated `exit_code: 0`)
- `rg -n --glob '*.go' --glob '!**/*_test.go' 'policy\.Authorize|KnownRoles|ActorRole|Role:' cmd/fid internal` - BLOCKED - `internal/workflow/define.go:56:\t\tActorRole: actor.id,` (expected `exit_code: 1`; terminated `exit_code: 0`)
- `rg -n --glob '*.go' --glob '!**/*_test.go' '(INTO|UPDATE|DELETE FROM) (roles|claims|agent_runs|task_completions|gate_runs|ac_verifications|concerns|workflow_transition_requires)\b|Table:\s*"(roles|claims|agent_runs|task_completions|gate_runs|ac_verifications|concerns|workflow_transition_requires)"' cmd/fid internal` - PASS - `internal/issues/ac_verify.go:160:\t\t\t\tTable:  "ac_verifications",` (terminated `exit_code: 0`). Classified below.
- `rg -n --glob '*.go' --glob '!**/*_test.go' 'EventType:\s*"' cmd/fid internal` - PASS - `internal/labels/add.go:121:\tevent := journal.Event{EventSchema: journal.EventSchemaV2, EventType: "label.added", ActorV2: actor.eventActor(), Aggregate: journal.Aggregate{ID: name, Sequence: sequence, Type: "label"}, CommandName: "label.add", Effects: []journal.Effect{{Op: "insert", Table: "labels", PK: map[string]any{"name": name}, Values: labelAddEffectValues(colour, description)}}}` (terminated `exit_code: 0`). Classified below.
- `rg -n --glob '*.go' --glob '!**/*_test.go' '(INTO|UPDATE|DELETE FROM) (roles|gate_runs|task_completions|agent_runs|ac_verifications|concerns)\b|Table:\s*"(roles|gate_runs|task_completions|agent_runs|ac_verifications|concerns|workflow_transition_requires)"' cmd/fid internal` - BLOCKED - `internal/issues/ac_verify.go:160:\t\t\t\tTable:  "ac_verifications",` (expected `exit_code: 1`; terminated `exit_code: 0`)
- `rg -n 'EventType:\s*"(dev\.|issue\.completed|issue\.claim\.|run\.recorded)|"event_type":"(dev\.|issue\.completed|issue\.claim\.|run\.recorded)' cmd/fid internal testdata` - BLOCKED - `internal/issues/transition_test.go:1408:\t\tEventType:   "issue.claim.seeded",` (expected `exit_code: 1`; terminated `exit_code: 0`)
- `rg -n --glob '*.go' --glob '!**/*_test.go' '"event"\s*:\s*"[a-z0-9]+_[a-z0-9_]+"|Event:\s*"[a-z0-9]+_[a-z0-9_]"' cmd/fid internal` - PASS - `exit 1 (no output)` (expected `exit_code: 1`; terminated `exit_code: 1`)
- `gofmt -l cmd internal` - BLOCKED - `internal/claims/generic.go` (terminated `exit_code: 0`; emitted paths require formatting.)
- `golangci-lint run ./...` - BLOCKED - `* typecheck: 1` (terminated `exit_code: 1`)
- `go test -race -coverprofile=coverage.out ./...` - BLOCKED - `FAIL` (terminated `exit_code: 1`)
- `go tool cover -func=coverage.out` - PASS - `total:\t\t\t\t\t\t\t\t\t\t(statements)\t\t\t\t\t\t48.4%` (terminated `exit_code: 0`). Generated from the failed coverage attempt and not completion evidence.

### Classified live inventories

- Broad table inventory: `internal/claims/generic.go:151,185,221` are retained final generic Claim writers; `internal/workflow/authoring_effects.go:97` is retained structural workflow authoring; `internal/issues/concern_raise.go:222`, `internal/issues/concern_address.go:196`, and `internal/issues/ac_verify.go:160` are prohibited frozen-state writers awaiting in-scope removal. `UNCLASSIFIED=0`.
- Event inventory: `internal/workflow/{activation,authoring,define,assign,envelope,run}.go`, `internal/evidence/{artifact,evidence,execution,registration}.go`, and `internal/issues/{create_commit,update,transition,ac}.go` are retained Event-v2 families. `internal/bootstrap/init.go:366` is the frozen task-21 bootstrap reader/writer boundary. `internal/labels/{add}.go`, `internal/links/{add,remove}.go`, and `internal/requirements/sync.go` are retained task-9 attribution families outside task-20 scope. `internal/issues/{ac_verify,concern_raise,concern_address}.go` are prohibited V1 judgment Event writers awaiting in-scope removal. `UNCLASSIFIED=0`.

### Attempt 2 exact command results (TRG-21-020 / DEC-173)

- `go test -race ./internal/db ./internal/maintenance ./internal/cli ./cmd/fid/...` - BLOCKED - `FAIL` (terminated `exit_code: 1`; retained V1 db, maintenance, and cmd fixtures still reference removed tables/packages.)
- `go test -tags integration -race ./cmd/fid/integration/...` - BLOCKED - `FAIL` (terminated `exit_code: 1`; generic Claim compilation first exposed the deleted legacy close helpers.)
- `ls internal/dev` - BLOCKED - `exit 0 (no output)` (expected `exit_code: 2`; terminated `exit_code: 0`; the empty removed-package directory remains.)
- `ls internal/runs` - BLOCKED - `exit 0 (no output)` (expected `exit_code: 2`; terminated `exit_code: 0`; the empty removed-package directory remains.)
- `ls cmd/fid/dev.go` - PASS - `ls: cannot access 'cmd/fid/dev.go': No such file or directory` (expected `exit_code: 2`; terminated `exit_code: 2`)
- `ls cmd/fid/policy.go` - PASS - `ls: cannot access 'cmd/fid/policy.go': No such file or directory` (expected `exit_code: 2`; terminated `exit_code: 2`)
- `ls cmd/fid/issue_complete.go` - PASS - `ls: cannot access 'cmd/fid/issue_complete.go': No such file or directory` (expected `exit_code: 2`; terminated `exit_code: 2`)
- `rg -n --glob '*.go' --glob '!**/*_test.go' 'String(Var)?\([^\n]*"by"|GetString\("by"|--by' cmd/fid internal` - BLOCKED - `cmd/fid/issue_mutation_response.go:237: value, err := cmd.Flags().GetString("by")` (expected `exit_code: 1`; terminated `exit_code: 0`; repaired in checkpoint `556b4e6`, rerun required.)
- `rg -n --glob '*.go' --glob '!**/*_test.go' 'policy\.Authorize|KnownRoles|ActorRole|Role:' cmd/fid internal` - PASS - `exit 1 (no output)` (expected `exit_code: 1`; terminated `exit_code: 1`)
- `rg -n --glob '*.go' --glob '!**/*_test.go' '(FROM|INTO|UPDATE|JOIN|DELETE FROM) (roles|claims|agent_runs|task_completions|gate_runs|ac_verifications|concerns|workflow_transition_requires)\b|Table:\s*"(roles|claims|agent_runs|task_completions|gate_runs|ac_verifications|concerns|workflow_transition_requires)"' cmd/fid internal` - PASS - `internal/issues/transition_rules.go:125: FROM concerns` (terminated `exit_code: 0`; classified in the retained broad inventory; `UNCLASSIFIED=0`.)
- `rg -n --glob '*.go' --glob '!**/*_test.go' 'EventType:\s*"' cmd/fid internal` - PASS - `internal/workflow/assign.go:111: EventType: "workflow.assigned",` (terminated `exit_code: 0`; classified in the retained broad inventory; `UNCLASSIFIED=0`.)
- `rg -n --glob '*.go' --glob '!**/*_test.go' '(INTO|UPDATE|DELETE FROM) (roles|gate_runs|task_completions|agent_runs|ac_verifications|concerns)\b|Table:\s*"(roles|gate_runs|task_completions|agent_runs|ac_verifications|concerns)"' cmd/fid internal` - BLOCKED - `internal/workflow/authoring_effects.go:97: Op: "insert", Table: "workflow_transition_requires",` (expected `exit_code: 1`; terminated `exit_code: 0`; repaired in checkpoint `556b4e6`, rerun required.)
- `rg -n 'EventType:\s*"(dev\.|issue\.completed|issue\.claim\.|run\.recorded)|"event_type":"(dev\.|issue\.completed|issue\.claim\.|run\.recorded)' cmd/fid internal testdata` - BLOCKED - `internal/issues/transition_test.go:1408: EventType: "issue.claim.seeded",` (expected `exit_code: 1`; terminated `exit_code: 0`.)
- `rg -n --glob '*.go' --glob '!**/*_test.go' '"event"\s*:\s*"[a-z0-9]+_[a-z0-9_]+"|Event:\s*"[a-z0-9]+_[a-z0-9_]"' cmd/fid internal` - PASS - `exit 1 (no output)` (expected `exit_code: 1`; terminated `exit_code: 1`)
- `gofmt -l cmd internal` - BLOCKED - `internal/workflow/workflow_test.go` (terminated `exit_code: 0`; emitted paths require formatting.)
- `golangci-lint run ./...` - BLOCKED - `* typecheck: 1` (terminated `exit_code: 1`; retained policy tests refer to removed V1 authority functions.)
- `go test -race -coverprofile=coverage.out ./...` - BLOCKED - `FAIL` (terminated `exit_code: 1`; package failures include V1 Claim, role, and Run fixtures.)
- `go tool cover -func=coverage.out` - PASS - `total:                                (statements)                                48.4%` (terminated `exit_code: 0`; generated from the failed coverage test command and not used as completion evidence.)

### Attempt 2 convergence provenance

- `cmd/fid/issue_mutation_response.go:151` and `:236` - invalidating change: removal of the V1 `--by` contract; accepted replacement: no legacy actor reader is used for authority failures.
- `internal/claims/generic.go:90` - invalidating change: deletion of the legacy Claim shared file; accepted replacement: retained generic Claim owns its connection/row close helpers.
- `internal/workflow/authoring.go:61`, `authoring_effects.go:10`, `authoring_effects_test.go:9`, and `authoring_test.go:74` - invalidating change: 0014 drops seeded V1 requirement Effects; accepted replacement: a non-empty retired requirement is explicitly rejected and no dropped-table Effect is emitted.

## Deviations from spec

- The required no-role search cannot pass inside the named task-20 scope: its live production matches include `internal/bootstrap/init.go`, `internal/journal/lock.go`, `internal/labels/`, `internal/links/`, and `internal/requirements/`. The task does not authorize changes there; this report preserves the exact broad command rather than narrowing it.
- The first gate failure is immutable evidence. Its root cause is the accepted 0014 final-table rename and V1 source deletion exposing retained readers/tests; the next repair must converge them without changing any command, assertion semantics, or frozen-table policy.
- First report-lint attempt returned `FAIL report-lint found 11 failure(s)` (terminated `exit_code: 11`) because command escaping and post-termination annotations were not byte-exact; this report correction preserves the gate evidence and will rerun the same lint command.

## Concerns, technical debt, or notes for the Tech Lead

- Scope gap: the protected command-8 absence contract requires edits outside task-20's named source scope, notably the task-21-owned bootstrap/journal paths and task-9-owned attribution packages.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `556b4e65e39865a836b955ffa77f85857818d9fc`

## Commit subject

- `fix(convergence): remove retired workflow requirement effects`
