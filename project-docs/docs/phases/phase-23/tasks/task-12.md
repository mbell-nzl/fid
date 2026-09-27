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
- VAL-053: `fid requirement sync`, its parser, and `requirement.synced` are
  absent from the binary, from both Event catalogs, and from the source tree.
  The `requirements.last_synced_hash` column belongs to the V1 table that
  migration 0018 drops, so its absence is asserted by VAL-061 at task 15.
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
- DEC-260: `requirement.record-reference` and `requirement.link` need a
  trusted Actor resolved from the finalized Kit Activation, but no
  Actor-profile capability and no Envelope grant. Do not add either name to
  either allowlist copy (`internal/workflow/activation.go`,
  `internal/journal/current_event_value_relations.go`); the closed set stays
  at thirty. The real-process journey's Actor profile holds neither name.

## Scope

### Files To Create Or Modify

- `internal/review/` or a dedicated file for requirement reference rules.
- Delete `internal/requirements/`.
- `internal/links/registry.go`, `add.go`, `list.go`, `errors.go`, `types.go`.
- `internal/journal/project_effects.go`, `current_event_catalog.go`.
- `internal/journal/current_event_json.go`,
  `internal/journal/current_event_result_projection.go`, and
  `internal/journal/current_event_value_catalog.go` - the typed-result
  registration for the last two families this task unfences (R-16, DEC-231).
  Deleting the final fence keys routes `requirement.referenced` and
  `requirement.linked` onto the journal's strict result path, so all three
  registries are required here or the first successful
  `requirement record-reference` is rejected at append time. Add: the two
  `fields` grammar entries derived from the `requirement.record-reference` and
  `requirement.link` result rows in `public-interface.md` section 4; the two
  matching `currentExpectedResult` branches; and the `RequirementRef` and
  `RequirementLink` entries in `currentValueShapes` from `public-interface.md`
  section 3. Neither is composite, so neither needs a fill helper. Also, in
  `current_event_json.go` and distinct from the `fields` map, register
  `requirement_refs.source_json` as `Source` in the `effect.Table + "." +
  column` switch at `:14` (DEC-232).
- `internal/journal/current_event_updates.go` - `requirement.referenced`'s
  optional `requirement_refs/update` Effect, which supersedes a predecessor, is
  rejected with `current Event update table not admitted` because
  `validateCurrentUpdate`'s closed `switch effect.Table` at `:14` lists no Phase
  23 table (DEC-232). Admit `requirement_refs` with its exact changed-field set
  - `status` and `superseded_by_id` per `schema.md` section 2 - and its exact
  guard set. Task 8 added the sibling `judgment_concerns` entry; if it is
  missing, that is a task 8 defect.
- This task removes the fence map entirely, so after it no Phase 23 family is
  registered in `currentEventRules` without its entry in **every** registry
  listed in the `execution-scope.md` section 6 inventory - all six, not only the
  three named by DEC-231. That nine-family, six-registry audit belongs to this
  task, and a registry that is not on the inventory and rejects a Phase 23
  family is an escalation, not a local fix.
- `internal/export/archive_semantics_catalog.go` - delete the
  `requirement.synced` semantic rule. The removed Event's archive rule lives
  here, so without this file the declared absence search cannot pass and its
  cause would fall outside task 15 (R-09, DEC-167/173).
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
- `project-docs/specs/03_TRD_SCHEMA.md` section 3.2: record the `links` closed
  endpoint registry as `issue` only, and `requirement_issue_links.relation` as
  the closed set `implements`, `verifies`, `relates-to`. Both are accepted by
  the slice contract; apply them with this implementation.
- `project-docs/docs/phases/phase-23/tasks/task-12-done.md` in the full format,
  including the final fence removal and its parity proof.

## Expected Non-Zero Commands

- `rg -n --glob '!*_test.go' "requirement\.synced|last_synced_hash" cmd internal` - expected `exit_code: 1` - no match is the pass condition.

Record it as
`PASS - exit 1 (no output) (expected exit_code: 1; terminated exit_code: 1)`.
Never convert it to `exit 0 (no output)`. The declaration above is one line, as
the template grammar requires (DEC-260). At the planning base the matches are
  `internal/requirements/sync.go`, `internal/journal/current_event_catalog.go`,
  and `internal/export/archive_semantics_catalog.go`; all three are removed by
  this task's own scope.

## Exact Commands To Run

Working directory: repository root `/home/mark/projects/fid`, or the isolated
same-task worktree root explicitly recorded by the Runner. No command is a
placeholder. Apply gofmt and goimports to edited Go files before the final gate.

After implementation, every command below is required before task close.
The tagged integration gate is mandatory and is never substituted.

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...
go test -run '^$' -bench '^BenchmarkRebuildWorkspace100kEvents$' -benchtime 1x ./internal/journal/
rg -n --glob '!*_test.go' "requirement\.synced|last_synced_hash" cmd internal
```

Two gates run for minutes and are explicitly bounded (DEC-227): the `-timeout
30m` above replaces Go's implicit 10-minute per-binary default, which the tagged
suite already approaches at 560 s. Record in the completion report the wall time
`go test` reports for `cmd/fid/integration`, `cmd/fid`, and `internal/workflow`,
so a regression is attributed to the task that caused it. Budgets: tagged
integration under 900 s and no ordinary package over 600 s (VAL-079). A breach is
a bounded speed repair owned by this task when this task's own new journeys
caused it; otherwise stop and escalate. Never record a timed-out gate as `PASS`.

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
- the declared `rg` absence search, whose pass condition is `exit_code: 1`; see `## Expected Non-Zero Commands`.

## Notes For Developer

- This task belongs to QA Review window 4.
- One task per Developer session; recovery may continue the same task.
- Commit using Conventional Commits with `Refs: phase-23/task-12`.
- Never infer ownership from a historical filename or task label. No-progress
  budgets exclude interruptions and attempts forbidden from repairing the cause.
- Do not read held-out scenario material or activate another agent.

## Repair Amendment - review-9 MF-2 under DEC-261 (2026-09-28)

This runs after task 10's MF-1 repair completes.

- Follow review-9's MF-2 row (VAL-053). Remove `requirement.synced` from
  `p22PortableEventTypes` in `cmd/fid/integration/p22_portable_fake_test.go`.
  The entry was introduced in `5ec69e6`, before Phase 23; task 12's removal
  of the Event invalidated it.
- Add a focused assertion: the fake reader rejects an archive line of type
  `requirement.synced` and still accepts a retained family.
- Do not add Phase 23 families to the fake or re-pin it. That is task 16's
  portable-recovery scope (DEC-261 Part 2).
- In `task-12-done.md`, correct the source-tree absence claim. Record an
  inventory that includes test files, for example
  `rg -n "requirement\.synced" cmd internal`, and classify any remaining hit.
  Keep every earlier result.
- **Gates:** all seven exact commands above, unchanged.
