# Phase 21 - Task 20: V1 Removal And Migration 0014 - Completion Report

## Workflow outcome

- Task outcome: RECOVERY_REQUIRED
- Category: in-scope-implementation
- Condition: Finalization attempt 6 ran every amended task command to termination. Command 17 remains red because `internal/sqlquery/sqlquery_test.go:113` mutates the removed `roles` table instead of proving read-only rejection against a surviving table, and `internal/workflow/define_test.go:32` requests the retired `has_completion_report` transition requirement instead of asserting its accepted fail-closed rejection.
- Recovery attempts: 6
- Next action: converge those two retained assertions under DEC-173 without weakening their subjects, then rerun the complete byte-identical command block from command 1 and record a real CLI entrypoint exercise.
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
- Attempt 6 current disposition: VAL-002 is PARTIAL because the scoped migration/package gate passes but the repository-wide coverage gate is BLOCKED by a retained final-schema test. VAL-074 is PARTIAL because both required broad inventories are fully classified below with `UNCLASSIFIED=0`, but command 17 is BLOCKED. VAL-075 is PARTIAL because every final V1 path, source, Event, and snake-case absence control passed, while command 17 remains red.

## Tests run

- `go test -race ./internal/db -run '^TestMigration0014RetiresV1AuthorityExecutionWithoutReconstruction$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`; red: `Migrate(14) error = no migration found for version 14: read down for version 14 .: file does not exist`)
- `go test -race ./internal/db -run '^TestMigration0014RetiresV1AuthorityExecutionWithoutReconstruction$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/internal/db	1.677s` (terminated `exit_code: 0`)
- `go test -race ./cmd/fid -run '^TestTask20V1CommandsAreAbsent$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`; red: `fid [dev] resolved command "fid dev" with err <nil>, want unknown command`)
- Real entrypoint exercise: BLOCKED - the exact required package command did not build `cmd/fid`, so no binary could be truthfully exercised after the removal checkpoint.
- Real entrypoint exercise for attempt 6: UNVERIFIED - the dispatch authorized only the 18 task-spec commands; command 1 compiled and tested `cmd/fid/...`, but no additional CLI command was authorized before the required repository-wide gate failed.

## Exact quality-gate results

### Attempt 4 package-gate result (TRG-21-020 / DEC-173)

- `go test -race ./internal/db ./internal/maintenance ./internal/cli ./cmd/fid/...` - BLOCKED - `FAIL` (terminated `exit_code: 1`; completed after 327.878s. Remaining failures are retained V1 DB migration/ref-trigger inventories, V1 closure/bootstrap fixtures, and command-surface expectations. The attempt is preserved as failed evidence.)

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
- `rg -n --glob '*.go' --glob '!**/*_test.go' '(FROM|INTO|UPDATE|JOIN|DELETE FROM) (roles|claims|agent_runs|task_completions|gate_runs|ac_verifications|concerns|workflow_transition_requires)\b|Table:\s*"(roles|claims|agent_runs|task_completions|gate_runs|ac_verifications|concerns|workflow_transition_requires)"' cmd/fid internal` - PASS - `internal/issues/transition_rules.go:125: FROM concerns` (terminated `exit_code: 0`). Classified in the retained broad inventory; `UNCLASSIFIED=0`.
- `rg -n --glob '*.go' --glob '!**/*_test.go' 'EventType:\s*"' cmd/fid internal` - PASS - `internal/workflow/assign.go:111: EventType: "workflow.assigned",` (terminated `exit_code: 0`). Classified in the retained broad inventory; `UNCLASSIFIED=0`.
- `rg -n --glob '*.go' --glob '!**/*_test.go' '(INTO|UPDATE|DELETE FROM) (roles|gate_runs|task_completions|agent_runs|ac_verifications|concerns)\b|Table:\s*"(roles|gate_runs|task_completions|agent_runs|ac_verifications|concerns)"' cmd/fid internal` - BLOCKED - `internal/workflow/authoring_effects.go:97: Op: "insert", Table: "workflow_transition_requires",` (expected `exit_code: 1`; terminated `exit_code: 0`; repaired in checkpoint `556b4e6`, rerun required.)
- `rg -n 'EventType:\s*"(dev\.|issue\.completed|issue\.claim\.|run\.recorded)|"event_type":"(dev\.|issue\.completed|issue\.claim\.|run\.recorded)' cmd/fid internal testdata` - BLOCKED - `internal/issues/transition_test.go:1408: EventType: "issue.claim.seeded",` (expected `exit_code: 1`; terminated `exit_code: 0`.)
- `rg -n --glob '*.go' --glob '!**/*_test.go' '"event"\s*:\s*"[a-z0-9]+_[a-z0-9_]+"|Event:\s*"[a-z0-9]+_[a-z0-9_]"' cmd/fid internal` - PASS - `exit 1 (no output)` (expected `exit_code: 1`; terminated `exit_code: 1`)
- `gofmt -l cmd internal` - BLOCKED - `internal/workflow/workflow_test.go` (terminated `exit_code: 0`; emitted paths require formatting.)
- `golangci-lint run ./...` - BLOCKED - `* typecheck: 1` (terminated `exit_code: 1`; retained policy tests refer to removed V1 authority functions.)
- `go test -race -coverprofile=coverage.out ./...` - BLOCKED - `FAIL` (terminated `exit_code: 1`; package failures include V1 Claim, role, and Run fixtures.)
- `go tool cover -func=coverage.out` - PASS - `total:                                (statements)                                48.4%` (terminated `exit_code: 0`). Generated from the failed coverage test command and not used as completion evidence.

### Attempt 2 convergence provenance

- `cmd/fid/issue_mutation_response.go:151` and `:236` - invalidating change: removal of the V1 `--by` contract; accepted replacement: no legacy actor reader is used for authority failures.
- `internal/claims/generic.go:90` - invalidating change: deletion of the legacy Claim shared file; accepted replacement: retained generic Claim owns its connection/row close helpers.
- `internal/workflow/authoring.go:61`, `authoring_effects.go:10`, `authoring_effects_test.go:9`, and `authoring_test.go:74` - invalidating change: 0014 drops seeded V1 requirement Effects; accepted replacement: a non-empty retired requirement is explicitly rejected and no dropped-table Effect is emitted.

### Attempt 3 exact command results (TRG-21-020 / DEC-173)

- `go test -race ./internal/db ./internal/maintenance ./internal/cli ./cmd/fid/...` - BLOCKED - `FAIL` (terminated `exit_code: 1`; retained V1 Claim readers in maintenance/issues and `cmd/fid/stats_test.go` still imports removed `internal/runs`.)
- `go test -tags integration -race ./cmd/fid/integration/...` - BLOCKED - `FAIL` (terminated `exit_code: 1`; `cmd/fid/response_errors.go` still names removed V1 policy error types.)
- `ls internal/dev` - BLOCKED - `exit 0 (no output)` (expected `exit_code: 2`; terminated `exit_code: 0`; the removed package directory is empty but remains on disk.)
- `ls internal/runs` - BLOCKED - `exit 0 (no output)` (expected `exit_code: 2`; terminated `exit_code: 0`; the removed package directory is empty but remains on disk.)
- `ls cmd/fid/dev.go` - PASS - `ls: cannot access 'cmd/fid/dev.go': No such file or directory` (expected `exit_code: 2`; terminated `exit_code: 2`)
- `ls cmd/fid/policy.go` - PASS - `ls: cannot access 'cmd/fid/policy.go': No such file or directory` (expected `exit_code: 2`; terminated `exit_code: 2`)
- `ls cmd/fid/issue_complete.go` - PASS - `ls: cannot access 'cmd/fid/issue_complete.go': No such file or directory` (expected `exit_code: 2`; terminated `exit_code: 2`)
- `rg -n --glob '*.go' --glob '!**/*_test.go' 'String(Var)?\([^\n]*"by"|GetString\("by"|--by' cmd/fid internal` - PASS - `exit 1 (no output)` (expected `exit_code: 1`; terminated `exit_code: 1`)
- `rg -n --glob '*.go' --glob '!**/*_test.go' 'policy\.Authorize|KnownRoles|ActorRole|Role:' cmd/fid internal` - PASS - `exit 1 (no output)` (expected `exit_code: 1`; terminated `exit_code: 1`)
- `rg -n --glob '*.go' --glob '!**/*_test.go' '(FROM|INTO|UPDATE|JOIN|DELETE FROM) (roles|claims|agent_runs|task_completions|gate_runs|ac_verifications|concerns|workflow_transition_requires)\b|Table:\s*"(roles|claims|agent_runs|task_completions|gate_runs|ac_verifications|concerns|workflow_transition_requires)"' cmd/fid internal` - BLOCKED - `internal/issues/concern_list.go:116:\t\tFROM concerns` (terminated `exit_code: 0`; broad inventory output changed after the journal fixture conversion and needs refreshed per-line classification.)
- `rg -n --glob '*.go' --glob '!**/*_test.go' 'EventType:\s*"' cmd/fid internal` - BLOCKED - `internal/workflow/assign.go:111: EventType: "workflow.assigned",` (terminated `exit_code: 0`; broad inventory output changed after the journal fixture conversion and needs refreshed per-line classification.)
- `rg -n --glob '*.go' --glob '!**/*_test.go' '(INTO|UPDATE|DELETE FROM) (roles|gate_runs|task_completions|agent_runs|ac_verifications|concerns)\b|Table:\s*"(roles|gate_runs|task_completions|agent_runs|ac_verifications|concerns)"' cmd/fid internal` - PASS - `exit 1 (no output)` (expected `exit_code: 1`; terminated `exit_code: 1`)
- `rg -n 'EventType:\s*"(dev\.|issue\.completed|issue\.claim\.|run\.recorded)|"event_type":"(dev\.|issue\.completed|issue\.claim\.|run\.recorded)' cmd/fid internal testdata` - BLOCKED - `cmd/fid/issue_transition_test.go:875:\t\tEventType:   "issue.claim.seeded",` (expected `exit_code: 1`; terminated `exit_code: 0`)
- `rg -n --glob '*.go' --glob '!**/*_test.go' '"event"\s*:\s*"[a-z0-9]+_[a-z0-9_]+"|Event:\s*"[a-z0-9]+_[a-z0-9_]"' cmd/fid internal` - PASS - `exit 1 (no output)` (expected `exit_code: 1`; terminated `exit_code: 1`)
- `gofmt -l cmd internal` - BLOCKED - `internal/workflow/workflow_test.go` (terminated `exit_code: 0`; emitted paths require formatting.)
- `golangci-lint run ./...` - BLOCKED - `* typecheck: 1` (terminated `exit_code: 1`; retained V1 policy tests refer to deleted authority types.)
- `go test -race -coverprofile=coverage.out ./...` - BLOCKED - `FAIL` (terminated `exit_code: 1`; retained V1 Claim/policy/Run fixtures and final-schema readers remain red.)
- `go tool cover -func=coverage.out` - PASS - `total:\t\t\t\t\t\t\t\t\t\t(statements)\t\t\t\t\t\t52.2%` (terminated `exit_code: 0`). Generated from the failed coverage attempt and not completion evidence.

### Attempt 3 convergence provenance

- `internal/journal/project_test.go:18`, `replay_test.go:73`, `rebuild_scale_test.go:192`, `phase2_rebuild_test.go:188`, and `writer_test.go:476` - invalidating change: migration 0014 drops `roles`; accepted replacement: each projection, replay, scale, writer, and rebuild fixture exercises the same Effect behavior against the surviving `labels` table.
- `testdata/journal/replay-v1.golden.jsonl:1` and `corrupt-tail.jsonl:1` - invalidating change: migration 0014 drops `roles`; accepted replacement: the fixture seeds one `labels` row while preserving ordered replay/corrupt-tail subjects.
- `testdata/journal/p20-issue-revision-effects.golden.jsonl:4` - invalidating change: V1 dev Event removal; accepted replacement: three retained Issue Event families carry the same revision Effects and preserve the ten-Event chain.
- `internal/journal/replay_test.go:130` - invalidating change: the three dev fixture values become retained Issue values; accepted replacement: revision assertions use fixture position, preserving every prior/value/precondition check despite repeated Issue family names.
- `internal/journal/phase2_rebuild_test.go:187` - invalidating change: migration 0014 drops V1 Claim/Gate and freezes completion/judgment/telemetry state; accepted replacement: the parity fixture retains its surviving issue/workflow/requirement/link/activity projection subject without Effects on dropped or frozen tables.

### Attempt 6 final ordered command results (TRG-21-020 / DEC-173)

- `go test -race ./internal/db ./internal/maintenance ./internal/cli ./cmd/fid/...` - PASS - `ok  ` (terminated `exit_code: 0`)
- `go test -tags integration -race ./cmd/fid/integration/...` - PASS - `ok  ` (terminated `exit_code: 0`)
- `ls internal/dev` - PASS - `ls: cannot access 'internal/dev': No such file or directory` (expected `exit_code: 2`; terminated `exit_code: 2`)
- `ls internal/runs` - PASS - `ls: cannot access 'internal/runs': No such file or directory` (expected `exit_code: 2`; terminated `exit_code: 2`)
- `ls cmd/fid/dev.go` - PASS - `ls: cannot access 'cmd/fid/dev.go': No such file or directory` (expected `exit_code: 2`; terminated `exit_code: 2`)
- `ls cmd/fid/policy.go` - PASS - `ls: cannot access 'cmd/fid/policy.go': No such file or directory` (expected `exit_code: 2`; terminated `exit_code: 2`)
- `ls cmd/fid/issue_complete.go` - PASS - `ls: cannot access 'cmd/fid/issue_complete.go': No such file or directory` (expected `exit_code: 2`; terminated `exit_code: 2`)
- `rg -n --glob '*.go' --glob '!**/*_test.go' 'String(Var)?\([^\n]*"by"|GetString\("by"|--by' cmd/fid internal` - PASS - `exit 1 (no output)` (expected `exit_code: 1`; terminated `exit_code: 1`)
- `rg -n --glob '*.go' --glob '!**/*_test.go' 'policy\.Authorize|KnownRoles|ActorRole|Role:' cmd/fid internal` - PASS - `exit 1 (no output)` (expected `exit_code: 1`; terminated `exit_code: 1`)
- `rg -n --glob '*.go' --glob '!**/*_test.go' '(FROM|INTO|UPDATE|JOIN|DELETE FROM) (roles|claims|agent_runs|task_completions|gate_runs|ac_verifications|concerns|workflow_transition_requires)\b|Table:\s*"(roles|claims|agent_runs|task_completions|gate_runs|ac_verifications|concerns|workflow_transition_requires)"' cmd/fid internal` - PASS - `completion.id` (terminated `exit_code: 0`). Classified below.
- `rg -n --glob '*.go' --glob '!**/*_test.go' 'EventType:\s*"' cmd/fid internal` - PASS - `activation.component-staged` (terminated `exit_code: 0`). Classified below.
- `rg -n --glob '*.go' --glob '!**/*_test.go' '(INTO|UPDATE|DELETE FROM) (roles|gate_runs|task_completions|agent_runs|ac_verifications|concerns)\b|Table:\s*"(roles|gate_runs|task_completions|agent_runs|ac_verifications|concerns)"' cmd/fid internal` - PASS - `exit 1 (no output)` (expected `exit_code: 1`; terminated `exit_code: 1`)
- `rg -n 'EventType:\s*"(dev\.|issue\.completed|issue\.claim\.|run\.recorded)|"event_type":"(dev\.|issue\.completed|issue\.claim\.|run\.recorded)' cmd/fid internal testdata` - PASS - `exit 1 (no output)` (expected `exit_code: 1`; terminated `exit_code: 1`)
- `rg -n --glob '*.go' --glob '!**/*_test.go' '"event"\s*:\s*"[a-z0-9]+_[a-z0-9_]+"|Event:\s*"[a-z0-9]+_[a-z0-9_]+"' cmd/fid internal` - PASS - `exit 1 (no output)` (expected `exit_code: 1`; terminated `exit_code: 1`)
- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -coverprofile=coverage.out ./...` - BLOCKED - `FAIL` (terminated `exit_code: 1`; `TestReadOnlyConnectionRejectsMutationWhenClassifierIsBypassed` at `internal/sqlquery/sqlquery_test.go:113` inserts into removed `roles`, and `TestDefineAppendsOneAtomicWorkflowDefinedEvent` at `internal/workflow/define_test.go:32` requests the retired completion-report transition requirement. This immutable failed attempt requires in-scope assertion convergence.)
- `go tool cover -func=coverage.out` - PASS - `total:` (terminated `exit_code: 0`). Generated from the failed coverage attempt and not completion evidence.

### Attempt 6 broad table inventory provenance

- `internal/workflow/run.go:219` - final `claims` reader; invalidating change: migration 0014 renames `work_claims`; accepted replacement: generic Claim final-table lookup.
- `internal/workflow/evaluation.go:176` - final `claims` reader; invalidating change: migration 0014 renames `work_claims`; accepted replacement: final Claim obligation evaluation.
- `internal/workflow/resume.go:139` - final `claims` reader; invalidating change: migration 0014 renames `work_claims`; accepted replacement: active final Claim resume lookup.
- `internal/workflow/envelope.go:293` - final `claims` reader; invalidating change: migration 0014 renames `work_claims`; accepted replacement: active final Claim count.
- `internal/workflow/registry.go:183` - frozen transition-requirement diagnostic reader; invalidating change: 0014 removes only displaced role/judgment rows; accepted replacement: retained final-schema registry inspection with unavailable requirements rejected before a dropped-table Effect.
- `internal/stats/stats.go:206` - frozen `agent_runs` diagnostic reader; invalidating change: task-20 removes V1 Run writers; accepted replacement: read-only telemetry aggregation.
- `internal/claims/generic.go:151,159,193,229,337,358,457` - final `claims` readers and generic Claim Effects; invalidating change: migration 0014 renames `work_claims`; accepted replacement: final generic Claim lifecycle subject with unchanged assertion strength.
- `internal/issues/detail.go:129,133` - frozen `ac_verifications` diagnostic reader; invalidating change: task-20 removes V1 completion writers; accepted replacement: read-only Issue detail inspection.
- `internal/issues/detail.go:167` - frozen `concerns` diagnostic reader; invalidating change: task-20 removes V1 judgment writers; accepted replacement: read-only Issue detail inspection.
- `internal/issues/detail.go:208` - frozen `task_completions` diagnostic reader; invalidating change: task-20 removes V1 completion writers; accepted replacement: read-only Issue detail inspection.
- `internal/issues/concern_list.go:116` - frozen `concerns` diagnostic reader; invalidating change: task-20 removes V1 judgment writers; accepted replacement: read-only concern listing.
- `internal/issues/transition_rules.go:112,125` - frozen `task_completions` and `concerns` diagnostic readers; invalidating change: task-20 removes V1 completion/judgment writers; accepted replacement: guarded transition diagnosis without state mutation.
- `internal/maintenance/warnings.go:42` - final `claims` reader; invalidating change: migration 0014 renames `work_claims`; accepted replacement: final Claim maintenance warning.
- `internal/maintenance/warnings.go:164,165,191,192` - frozen `ac_verifications`, `task_completions`, and `agent_runs` diagnostic joins; invalidating change: task-20 removes V1 writers; accepted replacement: read-only maintenance findings.
- Broad table inventory: `UNCLASSIFIED=0`.

### Attempt 6 Event inventory provenance

- `internal/links/add.go:112` and `internal/links/remove.go:126` - retained Event-v2 link families; invalidating change: ActorRole deletion; accepted replacement: generic Actor-v2 attribution.
- `internal/labels/add.go:119` - retained Event-v2 label family; invalidating change: ActorRole deletion; accepted replacement: generic Actor-v2 attribution.
- `internal/evidence/evidence.go:152`, `execution.go:297,371`, `artifact.go:434`, and `registration.go:180` - retained Phase 21 Event-v2 evidence families; invalidating change: V1 Run/Gate removal; accepted replacement: Receipt/Evidence execution Events.
- `internal/bootstrap/init.go:345` - frozen VAL-008 `project.initialized` Event; invalidating change: role bootstrap removal; accepted replacement: byte-identical project initialization Event.
- `internal/issues/create_commit.go:133`, `transition.go:230`, `ac.go:135`, and `update.go:145` - retained Event-v2 Issue families; invalidating change: dev/completion/legacy Claim Event removal; accepted replacement: Issue lifecycle Events.
- `internal/requirements/sync.go:94` - retained Event-v2 requirement-sync family; invalidating change: ActorRole deletion; accepted replacement: generic Actor-v2 attribution.
- `internal/workflow/define.go:107`, `assign.go:111`, `envelope.go:234,304,386`, and `activation.go:209,337` - retained Event-v2 workflow families; invalidating change: V1 role/judgment/Run removal; accepted replacement: final work-envelope and activation Events.
- Event inventory: `UNCLASSIFIED=0`.

### Attempt 6 per-row classifications

- `internal/workflow/run.go:219` - retained final `claims` reader; migration 0014 replacement: final Claim lookup.
- `internal/workflow/evaluation.go:176` - retained final `claims` reader; migration 0014 replacement: final Claim obligation evaluation.
- `internal/workflow/resume.go:139` - retained final `claims` reader; migration 0014 replacement: active Claim resume lookup.
- `internal/workflow/envelope.go:293` - retained final `claims` reader; migration 0014 replacement: active Claim count.
- `internal/workflow/registry.go:183` - frozen `workflow_transition_requires` diagnostic reader; 0014 replacement: final-schema inspection with retired requirements rejected before Effects.
- `internal/stats/stats.go:206` - frozen `agent_runs` diagnostic reader; V1 Run removal replacement: read-only telemetry aggregation.
- `internal/claims/generic.go:151` - retained final `claims` reader; migration 0014 replacement: generic Claim uniqueness lookup.
- `internal/claims/generic.go:159` - retained final `claims` Effect writer; migration 0014 replacement: generic Claim acquire Effect.
- `internal/claims/generic.go:193` - retained final `claims` Effect writer; migration 0014 replacement: generic Claim renew Effect.
- `internal/claims/generic.go:229` - retained final `claims` Effect writer; migration 0014 replacement: generic Claim release Effect.
- `internal/claims/generic.go:337` - retained final `claims` reader; migration 0014 replacement: generic Claim get lookup.
- `internal/claims/generic.go:358` - retained final `claims` reader; migration 0014 replacement: active Claim lookup.
- `internal/claims/generic.go:457` - retained final `claims` reader; migration 0014 replacement: generic Claim list query.
- `internal/issues/detail.go:129` - frozen `ac_verifications` diagnostic reader; V1 completion writer removal replacement: read-only Issue detail.
- `internal/issues/detail.go:133` - frozen `ac_verifications` diagnostic reader; V1 completion writer removal replacement: read-only Issue detail.
- `internal/issues/detail.go:167` - frozen `concerns` diagnostic reader; V1 judgment writer removal replacement: read-only Issue detail.
- `internal/issues/detail.go:208` - frozen `task_completions` diagnostic reader; V1 completion writer removal replacement: read-only Issue detail.
- `internal/issues/concern_list.go:116` - frozen `concerns` diagnostic reader; V1 judgment writer removal replacement: read-only concern list.
- `internal/issues/transition_rules.go:112` - frozen `task_completions` diagnostic reader; V1 completion writer removal replacement: guarded transition diagnosis.
- `internal/issues/transition_rules.go:125` - frozen `concerns` diagnostic reader; V1 judgment writer removal replacement: guarded transition diagnosis.
- `internal/maintenance/warnings.go:42` - retained final `claims` reader; migration 0014 replacement: Claim maintenance warning.
- `internal/maintenance/warnings.go:164` - frozen `ac_verifications` diagnostic reader; V1 completion writer removal replacement: read-only maintenance finding.
- `internal/maintenance/warnings.go:165` - frozen `task_completions` diagnostic reader; V1 completion writer removal replacement: read-only maintenance finding.
- `internal/maintenance/warnings.go:191` - frozen `task_completions` diagnostic reader; V1 completion writer removal replacement: read-only maintenance finding.
- `internal/maintenance/warnings.go:192` - frozen `agent_runs` diagnostic reader; V1 Run writer removal replacement: read-only maintenance finding.
- `internal/links/add.go:112` - retained Event-v2 `link.added`; ActorRole deletion replacement: generic Actor-v2 attribution.
- `internal/labels/add.go:119` - retained Event-v2 `label.added`; ActorRole deletion replacement: generic Actor-v2 attribution.
- `internal/links/remove.go:126` - retained Event-v2 `link.removed`; ActorRole deletion replacement: generic Actor-v2 attribution.
- `internal/evidence/evidence.go:152` - retained Event-v2 `evidence.linked`; V1 Run/Gate removal replacement: Receipt/Evidence event family.
- `internal/evidence/execution.go:297` - retained Event-v2 `gate.execution-started`; V1 Run/Gate removal replacement: Receipt/Evidence event family.
- `internal/evidence/execution.go:371` - retained Event-v2 `gate.receipt-recorded`; V1 Run/Gate removal replacement: Receipt/Evidence event family.
- `internal/evidence/artifact.go:434` - retained Event-v2 `artifact.availability-observed`; V1 Run/Gate removal replacement: Receipt/Evidence event family.
- `internal/bootstrap/init.go:345` - frozen VAL-008 Event-v2 `project.initialized`; role bootstrap removal replacement: byte-identical initialization Event.
- `internal/evidence/registration.go:180` - retained Event-v2 `gate.registered`; V1 Run/Gate removal replacement: Gate Registration event.
- `internal/issues/create_commit.go:133` - retained Event-v2 `issue.created`; V1 dev/completion/Claim Event removal replacement: Issue lifecycle Event.
- `internal/issues/transition.go:230` - retained Event-v2 `issue.transitioned`; V1 dev/completion/Claim Event removal replacement: Issue lifecycle Event.
- `internal/requirements/sync.go:94` - retained Event-v2 `requirement.synced`; ActorRole deletion replacement: generic Actor-v2 attribution.
- `internal/workflow/define.go:107` - retained Event-v2 `workflow.defined`; V1 role/judgment removal replacement: final workflow definition Event.
- `internal/issues/ac.go:135` - retained Event-v2 `issue.ac-added`; V1 completion Event removal replacement: Issue lifecycle Event.
- `internal/issues/update.go:145` - retained Event-v2 `issue.updated`; V1 dev/Claim Event removal replacement: Issue lifecycle Event.
- `internal/workflow/envelope.go:234` - retained Event-v2 `envelope.created`; V1 Run removal replacement: work-envelope Event.
- `internal/workflow/envelope.go:304` - retained Event-v2 `envelope.activated`; V1 Run removal replacement: work-envelope Event.
- `internal/workflow/envelope.go:386` - retained Event-v2 `envelope.completed`; V1 Run removal replacement: work-envelope Event.
- `internal/workflow/assign.go:111` - retained Event-v2 `workflow.assigned`; V1 role removal replacement: final workflow assignment Event.
- `internal/workflow/activation.go:209` - retained Event-v2 `activation.prepared`; V1 Run removal replacement: activation Event.
- `internal/workflow/activation.go:337` - retained Event-v2 `activation.component-staged`; V1 Run removal replacement: activation Event.

## Deviations from spec

- The required no-role search cannot pass inside the named task-20 scope: its live production matches include `internal/bootstrap/init.go`, `internal/journal/lock.go`, `internal/labels/`, `internal/links/`, and `internal/requirements/`. The task does not authorize changes there; this report preserves the exact broad command rather than narrowing it.
- The first gate failure is immutable evidence. Its root cause is the accepted 0014 final-table rename and V1 source deletion exposing retained readers/tests; the next repair must converge them without changing any command, assertion semantics, or frozen-table policy.
- Attempt 6 command 17 failed after the preceding 16 commands met their exact pass conditions. Its failures are retained test assertions in authorized DEC-173 packages; no command, expected status, protected contract, or prior failed evidence was changed.
- First report-lint attempt returned `FAIL report-lint found 11 failure(s)` (terminated `exit_code: 11`) because command escaping and post-termination annotations were not byte-exact; this report correction preserves the gate evidence and will rerun the same lint command.
- Second report-lint attempt returned `FAIL report-lint found 9 failure(s)` (terminated `exit_code: 9`); the recovery report remains committed, but its preserved PASS rows need report-format convergence before the control can return a clean recovery receipt.

## Concerns, technical debt, or notes for the Tech Lead

- Scope gap: the protected command-8 absence contract requires edits outside task-20's named source scope, notably the task-21-owned bootstrap/journal paths and task-9-owned attribution packages.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `8cad753ad612ea5fe252afe55b25ba0874039694`

## Commit subject

- `test(convergence): remove unreachable V1 fixture helpers`
