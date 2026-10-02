# Phase 23 - Task 20: Stable requirement identity on every reference

> Supplemental task under DEC-268 / DEC-269 Part 3 (Human checklist finding A3).
> Dispatch only after DEC-269's status line reads ACCEPTED.

## Objective

`requirement_refs` stores the project's stable `requirement_id` and a
per-version reference `id`, but the accepted `RequirementRef` wire shape omits
`requirement_id`. The journal result projection deletes it on purpose, and the
Human view labels the reference and the link line both "Requirement". This task
adds `requirement_id` to every `RequirementRef` and makes the Human requirement
views distinguish requirement, reference, and Issue link.

## Validation Contract Coverage

- VAL-089: every `RequirementRef` carries both the stable `requirement_id` and
  the reference `id`, and the Human requirement views distinguish requirement,
  reference, and Issue link.
- VAL-073 (Human, re-checked at phase end as HUMAN-04a/b/c-17 and
  HUMAN-04a-18): `requirement.get` and `requirement.list` views.
- Keep VAL-051, VAL-052, VAL-053, VAL-054, VAL-058, and VAL-060 green.

## Binding Constraints

- `public-interface.md` section 2 as amended by DEC-269:
  `RequirementRef = {id, requirement_id, source, status, superseded_by_id,
  created_at, event_id}`, with `requirement_id` immediately after `id`. It
  applies to the `requirement.record-reference` result, `requirement.get`, and
  `requirement.list` items.
- `contracts.md` section 12 as amended: `get <id>` accepts the stable
  requirement ID (active reference, else newest) or an exact reference ID (that
  row). Its links are always the stable requirement ID's links. This is the
  shipped SQL order. Assert it; do not change it.
- Unchanged: request shapes, DDL, Event types, Effect tables and keys,
  semantic keys, request digests, and `RequirementLink`. `schema.md` already
  has the column. The one changed stored value is the `record-reference`
  result in the `idempotency_results` Effect, which now carries
  `requirement_id` (DEC-270).
- Human view (contracts 14: labelled sections, purpose line, no routine
  narration). Each reference renders:
  - a heading line `Requirement <requirement_id>`;
  - `Reference ID: <id>; Status: <status>`, plus `; Superseded by: <id>` when
    superseded;
  - the existing source line.

  Issue links render under their own section label `Issue links` as
  `<requirement_id> <relation> <issue_id>`. No link line is labelled
  "Requirement". The `requirement.link` view renders only the `Issue links`
  section. Human text adds no member absent from JSON (VAL-058).
- Journal compatibility (amended by DEC-270): a journal written by a
  pre-task-20 Phase 23 binary with `requirement.referenced` is unsupported
  (ADR-0011) and must fail closed. Rebuild refuses it with
  `journal-replay-stopped` and leaves its source digest unchanged. Never
  rewrite journal history, add a legacy or versioned path, or relax result
  validation.
- The closed journal registries this value reaches are owned here
  (`execution-scope.md` section 9): the `RequirementRef` value shape and its
  result-projection fill. Reaching any other registry is an escalation.

## Scope

### Files To Create Or Modify

- `internal/review/requirement.go` (`loadRequirementRef`)
- `internal/journal/current_event_value_catalog.go` (`RequirementRef` shape)
- `internal/journal/current_event_result_projection.go` (remove the
  `requirement_id` deletion for `RequirementRef`)
- `cmd/fid/p23_human.go` (requirement rendering only)
- `internal/review/requirement_test.go`, `cmd/fid/integration/p23_requirement_*_test.go`,
  `cmd/fid/integration/p23_human_leaf_pairs_test.go`,
  `cmd/fid/integration/p23_human_fixtures_test.go` (requirement rows only)
- Retained assertions displaced by this change converge here with per-row
  provenance (`execution-scope.md` section 9)

### Files To Read For Context

- `project-docs/docs/phases/phase-23/decisions.md` DEC-269 Part 3
- `project-docs/docs/phases/phase-23/contracts.md` sections 12 and 14;
  `public-interface.md` section 2
- `project-docs/docs/phases/phase-23/qa/human-checklist.md`, finding A3 and rows
  HUMAN-04a/b/c-17, HUMAN-04a-18

### Out Of Scope

- Input ergonomics and every other view (wrap-up agenda, DEC-268).
- Requirement request shapes and the `requirement link` input.

## Acceptance Criteria

- [ ] Record a reference, supersede it with a new `Source`, and link the
  requirement. Every JSON `RequirementRef` from `record-reference`, `get`
  (by stable ID and by each reference ID), and `list` (with and without
  `--issue`) carries the correct `requirement_id`. (VAL-089)
- [ ] `get` resolves exactly as contracts 12 states, including a superseded
  reference ID returning that row with the stable ID's links. (VAL-089)
- [ ] Human output for those reads names the requirement by its stable ID,
  shows the reference ID on its own labelled line, and lists links under
  `Issue links`. A verbatim Human fixture per leaf, compared with the JSON
  field set, shows no Human-only member. (VAL-089, VAL-058)
- [ ] (DEC-270 a) A journal recorded by this task's binary with a reference,
  a superseding reference, and links rebuilds clean, with its source
  digest unchanged. JSON `requirement get` (stable ID and each reference
  ID) and `requirement list` are byte-identical before and after the
  rebuild. An idempotent `record-reference` retry returns the stored result
  with the correct `requirement_id`. This is a committed test that the
  tagged gate runs. (VAL-089)
- [ ] (DEC-270 b) The final binary's `rebuild` of the preserved pre-change
  journal (task base `d7f960d`) returns `journal-replay-stopped` at the
  `requirement.referenced` line, exit 6, with the source digest unchanged.
  This is one-time report evidence. (VAL-089)
- [ ] No env-gated or default-skipped test is committed;
  `TestP23Task20SeedPrechangeJournal` is removed. (DEC-270)
- [ ] The regenerated requirement Human fixture rows still run. (VAL-073)

## First Failing Test

The built-binary `requirement get REQ-<id> --json` expects
`result.requirement.requirement_id == "REQ-<id>"`. The field is absent today.

## Documentation Updates

- None expected. The Tech Lead applied the contract amendments under DEC-269;
  the canonical spec sync is recorded at wrap-up.

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

- None.
## Planning Amendment - DEC-270 (2026-10-02)

Attempt 1 returned `BLOCKED` correctly (TRG-23-017). The pre-change replay
criterion could not hold, and it protected no real journal:
`requirement.referenced` exists only in unreleased Phase 23 binaries, and
ADR-0011 makes no pre-adoption journal promise. DEC-270 replaces that
criterion with (a) a new-shape round trip and (b) fail-closed evidence for
pre-change history. Nothing else changes: DEC-269 Part 3's shape, exact
validation, scope, commands, and the other criteria stand.

Continue from the preserved, uncommitted partial tree. Its value-catalog,
projection, and `loadRequirementRef` edits are the accepted direction. Then:

- finish `cmd/fid/p23_human.go` requirement rendering;
- add the (a) round-trip test in an owned `p23_requirement_*_test.go` file;
- remove `TestP23Task20SeedPrechangeJournal`;
- converge any displaced retained assertion with per-row provenance;
- run all five commands;
- rerun the exact rebuild for (b) and report it, keeping attempt 1's
  failed evidence in the report.

## Notes For Developer

- Report path: `project-docs/docs/phases/phase-23/tasks/task-20-done.md`.
  Commit with `Refs: phase-23/task-20`.
- Record the pre-change reference commit and source journal digest used for
  the replay criterion.
