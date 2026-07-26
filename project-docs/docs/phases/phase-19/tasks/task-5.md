# Phase 19 - Task 5: `fid workflow list` in every output format

## Objective

Ship the read surface: `fid workflow list` returns defined workflows with
their state/transition counts and live assignments, in json, toon, markdown,
and yaml, with `--kind` and `--active` filters. This is the command an agent
uses to discover what workflows exist before assigning or transitioning.

## Validation Contract Coverage

- VAL-018: the JSON payload has exactly the documented field set, with counts
  equal to projected row counts and `null` for `any`.
- VAL-060: all four formats render and carry equivalent data.
- VAL-061: the markdown renderer has byte-exact characterization against a
  **committed** expected artifact (CF-022, partial).
- VAL-019: `--kind <k>` includes exactly `<k>` and `any` workflows and
  excludes other-kind workflows.
- VAL-062: `--active` excludes a test-seeded inactive workflow; the unfiltered
  run includes it.
- VAL-063: `--kind` and `--active` compose to exactly the intersection.
- VAL-020: a kind-assigned workflow renders `{"scope":"kind","kind":"<k>"}`.
- VAL-064: a project-assigned workflow renders `{"scope":"project"}` with the
  `kind` key **absent**.
- VAL-065: an unassigned workflow renders an empty array, never `null`.
- VAL-066: ordering is deterministic and two consecutive runs are
  byte-identical.
- VAL-067: an empty result renders as an empty array with exit 0 in all four
  formats.

Assertions live in
`project-docs/docs/phases/phase-19/qa/validation-contract.md`.

## Binding Constraints

- **Flag surface is exactly as specified** in `06_DESIGN_PLAN_CLI.md`:

  ```
  fid workflow list [--kind <kind>] [--active] [--format json|toon|markdown|yaml]
  ```

  Read-only: no `--by`, no journal append, no write lock. Default format is
  JSON, consistent with every other read command (PRD-OUT-001).
- **JSON payload shape** (PRD-OUT-001; `06_DESIGN_PLAN_CLI.md`): a top-level
  array of objects with `id`, `name`, `applies_to`, `active`, `states`
  (integer count), `transitions` (integer count), and `assigned_to`. Counts
  are actual projected row counts including materialized cancel edges — the
  numbers in the spec's example block are illustrative, not normative.
  `applies_to` renders `null` for `any` workflows, matching the column.
- **`assigned_to` shape** (VAL-020): kind scope renders
  `{"scope":"kind","kind":"<k>"}`; project scope renders `{"scope":"project"}`
  with the `kind` field **omitted** (use `omitempty`); an unassigned workflow
  renders an empty array, never `null`.
- **Filter semantics** (VAL-019): `--kind <k>` returns workflows whose
  `applies_to` is `<k>` **or** NULL (`any`), because an `any` workflow is
  applicable to that kind; it excludes workflows bound to a different kind.
  `--active` restricts to `active = 1`. Without `--active`, inactive
  workflows are included. The two filters compose.
- **Deterministic ordering.** Results are ordered by `workflows.id`;
  `assigned_to` entries are ordered project-scope-last, then by kind. A read
  command must be reproducible for golden comparison.
- **Markdown needs an explicit type case.** `internal/format/writeMarkdown`
  switches on concrete types and errors on unknown ones; add a case for the
  new list type and render it in a new `internal/format` file rather than
  growing an existing one past the 400-line budget. JSON, TOON, and YAML flow
  through the generic path and need no new code — but must still be proven
  (VAL-018).
- **CF-022, partial take.** This task ships **byte-exact golden
  characterization** for the new markdown renderer. It does **not** backfill
  characterization for pre-existing markdown renderers; that stays deferred
  (DEC-082).
- **Empty result is an empty array**, not an error and not `null`, in every
  format.

## Scope

### Files To Create Or Modify

- `cmd/fid/workflow.go` — add the `list` subcommand
- `cmd/fid/workflow_list_test.go` (new) — four-format and filter coverage
- `internal/workflow/list.go` (new) — the read query and list item type
- `internal/workflow/list_test.go` (new)
- `internal/format/markdown_workflow.go` (new) — markdown renderer
- `internal/format/markdown.go` — add the type case
- `internal/format/markdown_workflow_test.go` (new) — byte-exact golden
- `project-docs/specs/06_DESIGN_PLAN_CLI.md` — `fid workflow list` filters,
  ordering, `assigned_to` shapes, and the note that example counts are
  illustrative

### Files To Read For Context

- `internal/format/format.go`, `internal/format/markdown.go`,
  `internal/format/markdown_lists.go` (renderer conventions)
- `internal/issues/list.go` (read-command query conventions)
- `cmd/fid/policy.go` (simplest `--format` read command)
- `migrations/0009_v02a_workflow_assignments.up.sql` (Task 1)
- `project-docs/specs/06_DESIGN_PLAN_CLI.md` `fid workflow list`

### Out Of Scope

- Any write command or journal event.
- `fid workflow show`/`get` for a single workflow (unscheduled).
- Backfilling markdown characterization for pre-existing renderers (CF-022
  remainder).
- Workflow deactivation (the inactive fixture row is test-seeded).

## Acceptance Criteria

- [ ] The JSON payload has exactly the documented field set (no extra or
  missing keys), the `states` and `transitions` counts equal SQL row counts
  including materialized cancel edges, and `applies_to` is `null` for `any`
  (VAL-018)
- [ ] `fid workflow list` renders in json, toon, markdown, and yaml over one
  fixture, and all four decode to equivalent data (VAL-060)
- [ ] The markdown renderer has a byte-exact golden test whose expected bytes
  are committed, and the completion report says whether they are a `testdata`
  fixture or an in-test literal (VAL-061, CF-022)
- [ ] `--kind <k>` includes `applies_to = <k>` and `any` workflows and
  excludes workflows bound to another kind, with both the included and
  excluded ID sets asserted (VAL-019)
- [ ] `--active` excludes a test-seeded inactive workflow; the unfiltered run
  includes it (VAL-062)
- [ ] `--kind` and `--active` together return exactly the intersection
  (VAL-063)
- [ ] A kind-assigned workflow renders `{"scope":"kind","kind":"<k>"}`
  (VAL-020)
- [ ] A project-assigned workflow renders `{"scope":"project"}` with the
  `kind` key absent, not null (VAL-064)
- [ ] An unassigned workflow renders an empty array, never `null` (VAL-065)
- [ ] Results are ordered by `workflows.id` with `assigned_to` entries
  project-scope-last then by kind, and two consecutive runs are byte-identical
  (VAL-066)
- [ ] An empty result renders as an empty array with exit 0 in all four
  formats (VAL-067)
- [ ] `06_DESIGN_PLAN_CLI.md` updated and the `05_DESIGN_PLAN.md` §3.1
  `workflow list` entry verified in this task

## First Failing Test

`TestWorkflowListRendersDocumentedJSONFieldSet`: run `fid workflow list` over
a fixture and assert the JSON field set and counts — fails until the command
exists.

## Documentation Updates

- `project-docs/specs/06_DESIGN_PLAN_CLI.md` — `fid workflow list` filter
  semantics, deterministic ordering, `assigned_to` shapes, empty-result
  behavior, and a note that the example counts are illustrative.
- `project-docs/specs/05_DESIGN_PLAN.md` §3.1 — verify the `workflow list`
  entry added at planning (DEC-092) matches the shipped flag surface; correct
  it in this task if it diverges.

## Exact Commands To Run

Run exactly:

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
go build -o /tmp/fid ./cmd/fid
```

Required-command rules live in `AGENTS.md` (Commands): every required command
is reported `PASS`, or `BLOCKED` with reason; missing, silently skipped, or
narrower substituted commands are failed gates. A non-terminating command is
never `PASS`.

Conditional command required by this task:

- `go build -o /tmp/fid ./cmd/fid` — VAL-018 needs built-CLI output in all
  four formats.

## Notes For Developer

- Depends on Tasks 1, 3, and 4 (workflows to list and assignments to show).
- `internal/format` non-test files are capped at 400 lines (STRUCT-002) and
  `markdown.go` was split in Phase 17 precisely to stay under it. Put the new
  renderer in its own file.
- Cross-format equivalence is easiest to assert by decoding each format back
  to a generic structure and comparing; do not settle for "all four commands
  exited zero".
- Completion report at
  `project-docs/docs/phases/phase-19/tasks/task-5-done.md`.
