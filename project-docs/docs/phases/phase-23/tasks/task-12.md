# Phase 23 - Task 12: Requirement references and a single traceability model

## Objective

Replace V1 requirement authoring with digest-verified requirement references
and explicit requirement-to-Issue traceability, remove the markdown parser and
its Event, and restrict `links` to Issue endpoints so one fact has one model.

## Validation Contract Coverage

- VAL-051: `requirement.record-reference` stores a digest-verified reference and
  re-recording supersedes while retaining the prior row.
- VAL-052: `requirement.link` records traceability and rejects an unknown
  requirement.
- VAL-053: `fid requirement sync`, its parser, `last_synced_hash`, and
  `requirement.synced` are absent from the binary and latest schema.
- VAL-054: `links` accepts only Issue endpoints.

## Binding Constraints

- `contracts.md` section 12 is the exact authority; `schema.md` section 2 fixes
  `requirement_refs`, `requirement_issue_links`, and the guard trigger.
- Project files remain canonical. fid stores no requirement prose, summary body,
  or parsed section text. Verify the `Source` digest under the write lock and
  reject `source-digest-mismatch` with no Event.
- Supersession creates a new row and marks the prior one `superseded` with a
  pointer. Never edit a prior row in place. Exactly one `active` row exists per
  requirement ID.
- Removing `implements` and `implemented_by` from `links` and restricting
  endpoints to Issues is the point of this task: it eliminates the second
  authoritative model for requirement traceability (roadmap section 7). The
  `links` trigger rebuild itself is executed by migration 0018 in task 15; this
  task owns the registry, validation, adapters, and their tests, and must leave
  the two models reconciled at its own close.
- DEC-192: delete exactly the `"requirement.referenced"` and
  `"requirement.linked"` keys, then remove the now-empty
  `unavailableJudgmentEvents` map and its guard branch. Converge task 1's
  remaining fence assertions by preserving the VAL-003 subject against the
  Phase 22 recovery fence pattern rather than deleting the parity assertion.
- Remove `internal/requirements` entirely, including its parser, attribution,
  errors, and types, and remove the `requirement.synced` Event, its effect rule,
  operation mapping, and fixtures from the replay corpus.
- The V1 `requirements` table itself is dropped by migration 0018 in task 15.
  Until then the table exists but has no writer; do not add a temporary one.
- `execution-scope.md` is incorporated on the usual convergence terms.

## Scope

### Files To Create Or Modify

- `internal/review/` or a dedicated file for requirement reference rules.
- Delete `internal/requirements/`.
- `internal/links/registry.go`, `add.go`, `list.go`, `errors.go`, `types.go`.
- `internal/journal/project_effects.go`, `current_event_catalog.go`.
- `internal/cli/codes.go` - `requirement-not-found` retained, plus any adapter
  change required by the removed leaf.
- `cmd/fid/requirement.go` - replace `sync` with `record-reference`, `link`,
  `get`, `list`; `cmd/fid/link.go` adapters.
- Matching in-package tests and `cmd/fid/integration/` fixtures invalidated by
  this change, including the replay and scale corpora.
- `project-docs/docs/phases/phase-23/tasks/task-12-done.md`.

### Files To Read For Context

- `contracts.md` section 12; `public-interface.md` sections 3, 4, 5;
  `schema.md` sections 2, 4, 6.
- `v1-convergence-inventory.md` rows for requirements and links.
- `task-11-done.md` and the QA Lead checkpoint outcome; `execution-scope.md`.

### Out Of Scope

- Migration 0018 and every other V1 table removal.
- Issue link relations other than the removed requirement endpoints.

## Acceptance Criteria

- [ ] A mismatched digest rejects; re-recording supersedes with a pointer and
  leaves the prior row readable; exactly one active row per requirement ID.
  (VAL-051)
- [ ] A known link is stored with its relation and an unknown requirement
  rejects with no Event. (VAL-052)
- [ ] The `sync` leaf is unknown to the built binary, `last_synced_hash` and
  `requirement.synced` are absent from the latest schema and Event catalog, and
  the source search finds nothing. (VAL-053)
- [ ] Requirement endpoints and both removed relations reject, and the retained
  Issue relations still work. (VAL-054)
- [ ] The judgment fence map and its guard branch are removed with parity
  preserved. (VAL-003)
- [ ] Full required commands terminate with preserved red/green evidence.
- [ ] Every converged retained row records its provenance and executing command.

## First Failing Test

`TestRequirementRecordReferenceSupersedesPriorActiveRow` fails against the
absent requirement reference writer.

## Documentation Updates

- `project-docs/specs/06_DESIGN_PLAN_CLI.md` section 3: replace
  `fid requirement sync` with `fid requirement record-reference|link|get|list`.
- `project-docs/docs/phases/phase-23/tasks/task-12-done.md` in the full format,
  including the final fence removal and its parity proof.
## Exact Commands To Run

Working directory: repository root `/home/mark/projects/fid`, or the isolated
same-task worktree root explicitly recorded by the Runner. No command is a
placeholder. Apply gofmt and goimports to edited Go files before the final gate.

After implementation, every command below is required before task close.
The tagged integration gate is mandatory and is never substituted.

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
go test -tags integration -race ./cmd/fid/integration/...
go test -run '^$' -bench '^BenchmarkRebuildWorkspace100kEvents$' -benchtime 1x ./internal/journal/
```

Required-command rules live in `AGENTS.md` (Commands): every required command
is reported `PASS`, or `BLOCKED` with reason; missing, silently skipped, or
narrower substituted commands are failed gates. No exception to the full task
gate. Focused red/green selections never replace these commands. Report each
attempt with the command's own output token and observed termination. An
interrupted or hung command is never `PASS`; preserve all failed attempts and
rerun the byte-identical command after an authorized causal repair. Report lint
is a Runner control after the report commit, not a Developer gate.

List any conditional command required by this task and why:

- the named fresh-process rebuild benchmark, because this task removes an Event family from the replay corpus and must prove the reference rebuild is unaffected.

## Notes For Developer

- This task belongs to QA Review window 4.
- One task per Developer session; recovery may continue the same task.
- Commit using Conventional Commits with `Refs: phase-23/task-12`.
- Never infer ownership from a historical filename or task label. No-progress
  budgets exclude interruptions and attempts forbidden from repairing the cause.
- Do not read held-out scenario material or activate another agent.
