# Phase 21 - Task 2: Event V2, Idempotency, And Role-Free Bootstrap

## Objective

Implement the shared Event-v2 and durable idempotency foundation, then convert `project.initialized` and init/rebuild adapters to the exact role-free bootstrap fixture. This task owns `internal/bootstrap` explicitly because genesis is the first Event-v2 aggregate.

## Validation Contract Coverage

- VAL-003: owning observable in the Phase 21 validation contract.
- VAL-004: owning observable in the Phase 21 validation contract.
- VAL-005: owning observable in the Phase 21 validation contract.
- VAL-006: owning observable in the Phase 21 validation contract.
- VAL-007: owning observable in the Phase 21 validation contract.
- VAL-008: owning observable in the Phase 21 validation contract.

## Binding Constraints

- The genesis Event has the `fid-bootstrap` service Actor, project sequence 1, null repository/idempotency, and empty Effects.
- Initialization remains naturally exclusive and has no idempotency key.
- This task removes role seeding from the new fixture but does not drop the live `roles` table before task 20.
- Preserve retained V1 callers through the exact temporary compatibility rule
  in `contracts.md` section 14: the closed V1 role vocabulary is recognized
  only after the exact role-free Event-v2 bootstrap and only when no legacy
  role Effect exists. The bridge persists nothing, supplements no explicit
  role facts, grants no Phase 21 generic authority, and is removed by task 20.
- DEC-162: retained V1 tests broken by this task's accepted interface change are
  converged inside this task under `v1-convergence-inventory.md` section 1.2.
  Preserve each test's subject and assertion strength. No `--by` on `fid
  rebuild`, no bootstrap role Effect, no persisted compatibility row, no
  migration 0014, no `t.Skip`, no build tags, no `-run` narrowing, and no
  threshold, timeout, required-command, interface, or migration change.
- DEC-162: `journal.RebuildWorkspace` drops its vestigial `byRole` parameter and
  acquires the lock with the exact owner command `fid rebuild`. Leave
  `LockOwner.ActorRole`, `internal/journal/lock.go` role metadata, and V1
  lock-owner strings in unrelated tests to tasks 20-21.

## Scope

### Files To Create Or Modify

- `internal/journal/`
- `internal/bootstrap/`
- `internal/policy/`
- `internal/cli/`
- `cmd/fid/init.go`
- `cmd/fid/rebuild.go`
- `cmd/fid/init_test.go`
- `cmd/fid/rebuild_test.go`
- Retained Go test files and test fixtures repository-wide, but only where this
  task's accepted change invalidates them (DEC-162). Derive the exact set with:

  ```text
  rg -n --glob '*_test.go' '"rebuild", "--by"' cmd internal
  rg -n --glob '*_test.go' 'RebuildWorkspace\(' cmd internal
  rg -n --glob '*_test.go' '(FROM|INTO) roles' cmd internal
  ```

  then reconcile against the actual gate output; the commands are a starting
  inventory, not the acceptance boundary.
- `project-docs/docs/phases/phase-21/tasks/task-2-done.md`

### Files To Read For Context

- `project-docs/docs/phases/phase-21/contracts.md §§3,14`
- `project-docs/docs/phases/phase-21/public-interface.md §§1,4,6`
- `project-docs/docs/phases/phase-21/v1-convergence-inventory.md`

### Out Of Scope

- Other mutation-family conversions, changes to the V1 permissions matrix,
  and migration 0014.
- Any production change to a retained V1 caller. Test convergence adapts tests
  to the accepted interface; it never changes retained command behaviour to make
  a test pass. `--by` on retained mutations stays until task 8.

## Acceptance Criteria

- [ ] Converted fixtures use exact closed Event v2 and aggregate sequence rules. (VAL-003, VAL-004)
- [ ] Same-key replay, digest conflict, and concurrent serialization match the accepted contract. (VAL-005, VAL-006, VAL-007)
- [ ] Bootstrap emits the exact role-free `project.initialized` fixture. (VAL-008)
- [ ] The temporary V1 compatibility lookup recognizes exactly the seven
  accepted names only for an exact role-free Event-v2 bootstrap, preserves
  explicit legacy role facts as authoritative, rejects unknown/corrupt or
  uninitialized state, and writes no role fact or projection row. (Intermediate
  convergence invariant; its no-write property preserves VAL-008.)
- [ ] The task-20 removal inventory names the compatibility lookup explicitly;
  no new generic Phase 21 package imports or calls it. (No additional VAL: this
  is an intermediate convergence invariant removed before final validation.)
- [ ] Every retained test invalidated by the no-leaf-flags rebuild interface and
  the role-free bootstrap is converged under DEC-162 and
  `v1-convergence-inventory.md` section 1.2, and both previously failed required
  commands now terminate `exit_code: 0` byte-identically. (VAL-003, VAL-008;
  convergence rule, not a new observable.)
- [ ] The completion report names, one line each: every deleted assertion with
  the accepted test that already owns its equivalent, and every fixture that had
  to seed an explicit legacy `roles` Effect, for task-20 removal.

## First Failing Test

For the recovery continuation, first add failing table-driven policy tests for
the exact-bootstrap fallback, explicit-role-fact precedence, and fail-closed
unknown/uninitialized cases. Preserve the existing bootstrap and idempotency
red/green evidence.

The DEC-162 continuation has no new first failing test: the two required
commands are already red and are the red evidence. Preserve their prior failed
output, converge the retained tests, then rerun both commands byte-identically.

## Known Retained-Test Convergence Cases

Non-exhaustive; derived by inspection, not by running the gate.

| Case | Accepted adaptation |
|---|---|
| ~23 `"rebuild", "--by", "<role>"` call sites across `cmd/fid` and `cmd/fid/integration` | drop the argument pair only |
| `journal.RebuildWorkspace(..., "<role>")` call sites in `cmd/fid` and `internal` tests | drop the role argument only |
| `cmd/fid/issue_test.go` unauthorized `rebuild --by intern` rejection | convert to the accepted flag rejection, or delete because `cmd/fid/rebuild_test.go` already owns it; name the choice |
| `cmd/fid/sql_test.go` and `internal/sqlquery/sqlquery_test.go` expecting `developer`/`human` rows from a bootstrapped workspace | repoint the read to migration-seeded state; keep the read-only/structured-rows subject |
| `cmd/fid/integration/e2e_v01_closure_helpers_test.go` roles parity reads | same repointing rule |
| Third subtest of `TestKnownRolesTRG002CorruptTailTolerance` | repair the fixture to a genuinely un-upgradeable shape; keep the error, file, and line assertions |
| `internal/maintenance/vacuum_test.go` before/after roles snapshot | expected to remain valid with an empty snapshot; change nothing unless it fails |

## Documentation Updates

- `project-docs/docs/phases/phase-21/tasks/task-2-done.md` - record red/green evidence, deviations, required commands, and spec impact.

## Exact Commands To Run

Run exactly:

```text
go test -race ./internal/policy -run '^(TestKnownRolesUsesV1CompatibilityVocabularyForRoleFreeV2Bootstrap|TestKnownRolesKeepsExplicitRoleFactsAuthoritative|TestKnownRolesRejectsCompatibilityForUninitializedJournal|TestKnownRolesRejectsCompatibilityForCorruptJournal|TestAuthorizeRejectsUnknownCompatibilityRole)$'
go test -race ./internal/journal ./internal/bootstrap ./internal/cli ./cmd/fid/...
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

Every command is required and must terminate with session-captured evidence. Report lint is not a Developer task-gate command; after the report commit, the Phase Runner runs the phase brief's two-argument pre-dispatch control. A Ready QA plan may add a task-specific command; the Tech Lead must amend this section before authorization if it does.

## Notes For Developer

- Start after task 1 closes.
- This is a same-task recovery under DEC-160/DEC-161/DEC-162. Preserve every
  prior failed required-command attempt in the report and rerun their
  byte-identical commands after the causal repair; do not replace them with the
  focused policy command.
- DEC-162 authorizes the retained-test convergence. Adapting a retained test to
  the accepted interface is in scope; weakening one is not. If a retained test's
  subject cannot survive adaptation and no accepted test owns its equivalent,
  stop and escalate that single case rather than deleting it.
- Completion report: `project-docs/docs/phases/phase-21/tasks/task-2-done.md`.
- Stop for a protected contract change; ordinary in-scope red/green repair remains this task.
- Do not read held-out scenario material.
