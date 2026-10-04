# Phase 24 - Task 3: Input descriptor framework, help, templates and field-named validation

## Objective

Build the single input-descriptor registry that drives each `--input` leaf's
help input section, its `--example` template and closed structural
validation with field-path errors (contract section 8, DEC-276). Adopt it for
`stop raise`, `stop resolve` and `attention resolve`. Give every leaf in the
live tree a one-line description, and fix the shared "Run JSON input" flag
text. The other 36 leaves sit behind a pending-descriptor fence that later
tasks shrink.

## Validation Contract Coverage

- VAL-010: every live leaf has a one-line description.
- VAL-059: `attention resolve` help describes Attention resolution.
- VAL-011..016 and VAL-060 (framework; the three adopted leaves):
  descriptor/type conformance, help, `--example` form, typed templates,
  exclusivity, field-path validation, and no `--example` side effects.
- VAL-052: gate wall-time budgets.
- Advances VAL-057 (Human authoring, phase-end).

Assertions live in `project-docs/docs/phases/phase-24/qa/validation-contract.md`.

## Binding Constraints

- **Descriptor model.** Each field has:
  - a name and JSON type (object, array, string, integer, boolean);
  - its requirement: required, optional (omittable) or nullable;
  - allowed members, a format (ID, text, digest, OID, relative path,
    timestamp, sorted-unique), and integer or length bounds;
  - a one-line meaning and a typed placeholder;
  - children or an element type;
  - discriminated variants, needed by task 9 for activation `content`;
  - a `sensitive` mark.

  Descriptors live in `internal/cli` (no new package boundary). Per-leaf
  descriptors live in `cmd/fid`. No shared descriptor decoder may be named
  to match `func decodeClosed[A-Za-z]*Input` (inventory R5). Task 9
  declares that pattern `exit_code: 1` (DEC-280, RI-05).
- **Validation order.** First the existing size, UTF-8 and single-object
  rules. Then duplicate keys. Then structure, level by level: unknown keys
  (lexicographic), then members in documented order (missing, null, type,
  enum, format, bounds). Then decoding into the Go type. Domain validation
  runs last and is unchanged. The first structural failure is reported
  deterministically.
- **`invalid-argument` details** (shape unchanged, contract section 6).
  `argument` is `input.<path>`, with `[i]` for array elements. `value` and
  `allowed` follow the failure class:

  | Failure | `value` | `allowed` |
  |---|---|---|
  | unknown key | the key | sorted accepted keys |
  | duplicate key | the key | sorted accepted keys |
  | missing | null | `[name]` |
  | null on non-nullable | `"null"` | `[type]` |
  | wrong type | JSON type found | `[type]` |
  | enum non-member | the value | the members |
  | format violation | the value | `[format]` |

  A `sensitive` field (for example a literal environment value or a redact
  literal) never echoes its value: `value` is null. The message names the
  path.
- **`--help`** is generated from the descriptor: a purpose sentence, the
  input section and the template line. The `--input` flag text names the
  leaf's own input. The shared `runMutationFlags` "Run JSON input" text is
  replaced.
- **`--example`** follows DEC-276 exactly. It is handled before any workspace
  resolution, so it works outside a project. Combined with any other flag or
  argument, it returns `invalid-argument` (exit 2) with no read or write.
  Checks run in this order, and the first conflict is reported: `--json`,
  `--project`, `--input`, `--actor`, `--idempotency-key`, then a positional
  argument. The failure follows the Phase 21 grammar rule (offending spelling,
  accepted alternative):
  - message: `--example cannot be combined with <spelling>`;
  - details: `{argument:"<spelling>", value:<supplied value, or null for
    --json>, allowed:["--example"]}`;
  - a positional argument uses `argument:"arguments"` and the first
    positional as `value`.

  With `--json` present the failure is a `fid-cli.v1` JSON envelope carrying
  the leaf's operation ID. Otherwise the Human form applies until task 10
  converges it.
- **Fence.** A test asserts that registry ∪ fence equals the live `--input`
  set exactly, and that registry ∩ fence is empty. Tasks 4-8 remove their
  leaves; task 9 deletes the empty fence (DEC-192 pattern).
- No input shape changes in this task. JSON results and success output are
  unchanged.

## Scope

### Files To Create Or Modify

- `internal/cli/` input descriptor, validator and help/example renderers
  (new files), with tests
- `cmd/fid/` descriptor files for the stop and attention leaves (new);
  `cmd/fid/stop.go`, `attention.go`, `run.go` (`runMutationFlags`), plus a
  shared `--example` hookup
- `Short` for every leaf starting with `Usage:` today (inventory R15):
  `cmd/fid/{attention,concern,decision,deviation,finding,human_outcome,requirement,review,stop}.go`
- `cmd/fid/input_descriptor_test.go` (new: registry/fence/live-tree and
  reflection conformance); `cmd/fid/contract_surface_test.go` if help
  assertions converge

### Files To Read For Context

- `project-docs/docs/phases/phase-24/public-interface.md` sections 6, 8
- `project-docs/docs/phases/phase-24/input-inventory.md` sections 1, 2 (R5, R15)
- `project-docs/docs/phases/phase-24/decisions.md` DEC-276
- `cmd/fid/claim.go` `decodeClosedClaimInput`, and `cmd/fid/concern.go`
  `decodeClosedJudgmentInput` (existing closed-input rules)
- `project-docs/docs/phases/phase-22/public-interface.md` section 3 (stop and
  attention input shapes)

### Out Of Scope

- Descriptors for any other leaf (tasks 4-9). Human failure lines (task 10).
- Changing any input field or domain validation message.

## Acceptance Criteria

- [ ] A recursive walk of the live tree (depths 1-3; 92 leaves at planning,
  20 without a `Short`, inventory R15) finds a nonempty one-line `Short` on
  every leaf. (VAL-010)
- [ ] `attention resolve --help` describes Attention resolution input. No leaf
  outside the Run family describes its input as Run input. (VAL-059)
- [ ] Reflection conformance passes for the three adopted leaves, and a
  seeded mismatch fixture fails it. The fence test holds. (VAL-011)
- [ ] Help for the three leaves lists every field with type, requirement,
  bounds/allowed values and meaning, and names the template. Full text is
  pinned for `stop raise`. (VAL-012)
- [ ] `--example` prints one JSON object with required fields only, in
  documented order with 2-space indentation. (VAL-013)
- [ ] `--example` from a non-project temporary directory exits 0 with empty
  stderr and leaves the tree snapshot unchanged. (VAL-060)
- [ ] Each template decodes into its Go type with no unknown, missing or
  wrongly typed member. (VAL-014)
- [ ] The exclusivity corpus (`--json`, `--project <dir>`, `--input`,
  `--actor`, `--idempotency-key`, a positional argument) returns
  `invalid-argument`, exit 2, with the details above and no effect.
  (VAL-015)
- [ ] The structural corpus for the three leaves covers every failure class,
  at top level and nested (for example `input.details.*`,
  `input.resolution_condition.subject.*`), with the exact path. JSON for
  domain failures is unchanged. (VAL-016)
- [ ] Wall times inside budget. (VAL-052)

## First Failing Test

`attention resolve --input` with a key `bogus` returns `argument:"input.bogus"`. Today it returns `argument:"input"`.

## Documentation Updates

- `project-docs/specs/06_DESIGN_PLAN_CLI.md`: the help input-section and
  `--example` conventions, and field-path `invalid-argument`.

## Exact Commands To Run

Working directory: repository root `/home/mark/projects/fid`, or the isolated
same-task worktree root the Runner records. No command is a placeholder.
Apply gofmt and goimports to edited Go files before the final gate. Every
command below is required before task close. The tagged integration gate is
mandatory and is never substituted.

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...
```

The `-timeout 30m` bounds follow DEC-227/259. Record the wall time `go test`
reports for `cmd/fid/integration`, `cmd/fid`, `internal/workflow` and the
slowest other ordinary package. The budgets are tagged integration under
900 s and no ordinary package over 600 s (VAL-052); report any ordinary
package over the 480 s warning. A breach caused by this task's own new journeys is a bounded speed
repair owned by this task; otherwise stop and escalate. Never record a
timed-out gate as `PASS`.

Required-command rules live in `AGENTS.md` (Commands). Every required command
is reported `PASS` or `BLOCKED` with reason, with its own output token and
observed termination. Focused red/green selections never replace these
commands. Report lint is a Runner control, not a Developer gate.

List any conditional command required by this task and why:

- None.

## Notes For Developer

- Report path: `project-docs/docs/phases/phase-24/tasks/task-3-done.md`.
  Commit with `Refs: phase-24/task-3`.
- Red first: record the first failing test before changing product code.
- Retained-test convergence follows DEC-162/165/166. Record each displaced
  assertion with its introducing commit, phase and task. Never infer
  ownership from a filename or task label. Never delete or weaken an
  assertion to pass a gate.
- Any missing registry member, contract ambiguity, or need to change an
  accepted contract shape or assertion is a `BLOCKED` escalation candidate,
  not a local decision (DEC-195).
- Keep `internal/cli` files under 400 lines (`08_QUALITY_PLAN.md`
  budgets). Split renderer and validator files.
