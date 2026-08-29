# Phase 21 - Task 9: Retained Inventory Attribution And Judgment Freeze

## Objective

Convert every retained non-Issue inventory mutation to project-local `--actor` plus Event-v2 idempotency, remove the public V1 policy command after generic authority passes, isolate the legacy action registry to enumerated dev/completion callers until task 20, and make Phase 23 judgment mutations fail closed. Rename stats grouping without treating historical attribution as authority.

## Validation Contract Coverage

- VAL-028: owning observable in the Phase 21 validation contract.
- VAL-029: owning observable in the Phase 21 validation contract.
- VAL-030: owning observable in the Phase 21 validation contract.
- VAL-031: owning observable in the Phase 21 validation contract.

## Binding Constraints

- Cover workflow define/assign, AC add, labels, links, requirements, activity
  normalization, concerns, maintenance vacuum, and stats.
- No public role/judgment requirement authoring survives.
- Do not delete the legacy policy implementation while the enumerated
  dev/completion callers remain; prove new generic paths cannot reach it.
- Do not delete replacement-bound dev/completion packages or tables yet.

## Scope

### Files To Create Or Modify

- `internal/policy/`
- `internal/workflow/`
- `internal/issues/`
- `internal/labels/`
- `internal/links/`
- `internal/requirements/`
- `internal/activity/`
- `internal/stats/`
- `internal/maintenance/`
- `internal/cli/`
- `cmd/fid/` (including the `cmd/fid/integration/` package)
- `project-docs/docs/phases/phase-21/tasks/task-9-done.md`

### Files To Read For Context

- `project-docs/docs/phases/phase-21/public-interface.md §§4-5`
- `project-docs/docs/phases/phase-21/v1-convergence-inventory.md`

### Out Of Scope

- Generic Claim/Run implementation, semantics, table, CLI, and contract; final
  V1 removal.
- This does **not** exclude converging an incidental read of a surface this task
  removed or renamed that happens to sit inside a V1 claim, run, dev, or closure
  fixture. See DEC-165.

### Convergence Ownership (DEC-162, DEC-165, DEC-166)

Convergence ownership follows the invalidating accepted change, not the fixture
file's primary subject. This task's public `fid policy` removal, stats `--by` to
`--group-by` rename, and fail-closed vacuum/AC/concern conversions are task-9
changes, so task 9 converges every assertion they invalidated wherever it lives,
including in fixtures whose headline subject belongs to task 10, 11, or 20.

- Repoint each incidental `policy show` row read onto the retained unregistered
  legacy matrix helper, keeping every role/action assertion byte-identical.
  Extend the enumerated legacy test-caller list and name the callers in the
  completion report for task-20 removal. Do not re-register `fid policy`.
- Delete a parity block only where `TestPublicPolicyCommandIsAbsent` plus a
  named surviving denial assertion already own the accepted equivalent, and name
  both in the report.
- Leave every V1 Claim, Run, dev, and closure *behavioural* assertion
  byte-identical: `issue.claim.force_acquire` denial and `tl` force-acquire,
  one-active/expiry/wrong-holder claim paths, the `run.record` allowed-role
  matrix, `dev.start|gate|finish|block` denial and journal-unchanged checks, and
  closure journal-byte equality. Their owning task replaces them later.
- Known remaining sites at the time of this amendment: `cmd/fid/label_test.go`,
  `cmd/fid/link_test.go`, `cmd/fid/issue_label_test.go`,
  `cmd/fid/requirement_test.go`, `cmd/fid/issue_claim_test.go`,
  `cmd/fid/run_test.go`, `cmd/fid/dev_test.go`, and
  `cmd/fid/integration/e2e_v01_closure_helpers_test.go`. This list is a
  navigation aid, not a limit; the gate is the authority.
- `cmd/fid/task10_contract_test.go` is a Phase 20 task-10 artifact, not a Phase
  21 generic-Claims fixture. Rename it and its `TestTask10*` functions to
  subject-based inventory/request-contract names. Convert only Task-9-owned
  mutation rows from `--by` and old missing-role expectations to the accepted
  `--actor` plus idempotency request contract. Missing/malformed Actor requests
  still fail before mutation: JSON uses the exact current closed
  `invalid-argument` envelope, Human output names `--actor` or its accepted
  grammar, and the journal remains byte-identical. Selected-project successes
  use `--actor agent:developer` plus idempotency. Preserve unrelated assertions
  byte-identically.
- Historical task numbering in a path, Go test name, or prose label is not
  current-phase ownership. Use introducing-commit provenance plus the current
  task pack, convergence inventory, and accepted contract. See DEC-166.
- `t.Skip`, build tags, `-run` narrowing, threshold changes, and silent
  assertion loss are never convergence.

### Checkpoint Repair Obligation (DEC-169, MF-2)

Implement the complete `public-interface.md` section 7 process exit contract, not
only the reproduced `capability-unavailable` value:

- a category-to-exit function in `internal/cli` covering `invalid-input`→2,
  `not-found`→3, `conflict`→4, `precondition`→4, `authority`→5, `integrity`→6,
  `internal`→7, and `gate-unsuccessful`→8, keyed on the closed `Category` already
  attached to every code in `internal/cli/codes.go`;
- a typed error that carries the failure's code/category out of the leaf, with
  `cmd/fid/main.go` wired to it and Cobra unknown-command, unknown-flag, and usage
  errors mapped to exit 2;
- every existing `cli.WriteFailure` call site converted to return it, so Human and
  JSON failures share the same exit class.

Nothing about any envelope changes: no JSON field, code, category, message,
subject, details, channel, command grammar, capability availability, or writer
absence moves, and success stays exit 0. Streaming `gate run` terminal-frame exit
semantics are untouched and remain with tasks 14-15. Real-process coverage in
`cmd/fid/integration/` proves the four fail-closed leaves at exit 5, one
representative each of exit 2, 3, and 4, and unknown-command exit 2.

## Acceptance Criteria

- [ ] Every retained family has exact attributed/idempotent behavior. (VAL-028)
- [ ] Role/judgment workflow inputs cannot be authored. (VAL-029)
- [ ] Judgment mutations are unavailable/no-writer. (VAL-030)
- [ ] Stats uses exact grouping/attribution contract. (VAL-031)

## First Failing Test

Add one failing conformance row per retained/fail-closed family before shared conversion.

## Documentation Updates

- `project-docs/docs/phases/phase-21/tasks/task-9-done.md` - record red/green evidence, deviations, required commands, and spec impact.

## Exact Commands To Run

Run exactly:

```text
go test -race ./internal/policy ./internal/workflow ./internal/issues ./internal/labels ./internal/links ./internal/requirements ./internal/activity ./internal/stats ./internal/maintenance ./internal/cli ./cmd/fid/...
go test -tags integration -race ./cmd/fid/integration/...
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

Every command is required and must terminate with session-captured evidence. Report lint is not a Developer task-gate command; after the report commit, the Phase Runner runs the phase brief's two-argument pre-dispatch control. A Ready QA plan may add a task-specific command; the Tech Lead must amend this section before authorization if it does.

## Notes For Developer

- Start after tasks 7-8 close. QA Review window 3 follows tasks 7-9.
- Completion report: `project-docs/docs/phases/phase-21/tasks/task-9-done.md`.
- Stop for a protected contract change; ordinary in-scope red/green repair remains this task.
- Do not read held-out scenario material.
