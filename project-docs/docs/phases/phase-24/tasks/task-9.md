# Phase 24 - Task 9: Activation descriptors; whole-tree input, help and grammar conformance

## Objective

Add descriptors for the eight `activation` leaves (shapes unchanged),
including the discriminated `component` content. Delete the now-empty
pending-descriptor fence. Make the help, template, typed-template,
exclusivity and structural-validation conformance run over all 39 `--input`
leaves. Assemble the shared fake/fid conformance evidence for every changed
leaf.

Deliver two accepted Phase 21 allocations that were never delivered
(DEC-194, DEC-278). Grammar failures name the offending spelling across the
whole live tree. Unknown commands return `operation-unavailable` instead of
exiting silently.

**Dispatch gate:** this task is not dispatched until the Human decides D5
(contract section 13) and the Tech Lead records it as DEC-279 and marks
VAL-065 `Ready`.

## Validation Contract Coverage

- VAL-011..016 and VAL-060 (completion: all 39 leaves, no fence).
- VAL-035 (assembly across tasks 4-8).
- VAL-064: grammar failures name the offending spelling or value.
- VAL-065: unknown commands return `operation-unavailable` (after DEC-279).
- VAL-052.

Assertions live in `project-docs/docs/phases/phase-24/qa/validation-contract.md`.

## Binding Constraints

- Contract section 8; DEC-276.
- `activation component` `content` is `ActorProfile[]` when
  `component_type` is `authority-profile`, and `GateDefinition[]` when it is
  `gate-definitions`. The descriptor uses discriminated variants. Structural
  errors name the path inside the selected variant. A mismatched variant
  names `input.content[i]`.
- `GateDefinition` literal `environment[].value` and `redact_literals` are
  `sensitive`: their values never appear in failure details.
- Deleting the fence leaves one invariant: registry equals the live
  `--input` set (DEC-192 pattern). No other leaf changes shape.
- **No top-level sensitive field (DEC-278, I-17).** Per-field digests
  (DEC-275) fingerprint each top-level name alone. The conformance test fails
  if any descriptor marks a top-level field `sensitive`.
- **Single structural source.** No family-local `decodeClosed*Input` decoder
  remains after this task. A task that migrates a family's last leaf may
  remove that family's decoder; this task removes the rest. If one must stay,
  structural validation would have two sources: that is a `BLOCKED`
  escalation candidate, not a local decision.
- **Grammar failures (Phase 21 `public-interface.md` section 2; DEC-278).**
  Every live leaf except the Human-only diagnostics `fid journal inspect` and
  `fid sql` returns `invalid-argument`, exit 2, with no read or write. The
  JSON envelope carries the leaf's operation ID (for example `stop.raise`,
  never the bare Cobra name). `message` is the text Human mode printed after
  `error: ` for that failure at `eb34889`, normally the parser's text. Where
  Human mode printed nothing, use the parser's text.

  | Class | `argument` | `value` | `allowed` |
  |---|---|---|---|
  | unknown flag | the spelling as typed | the spelling | the leaf's accepted long flag spellings, sorted, without `--help` |
  | missing required flag | `--<name>` | null | `["--<name>"]` |
  | malformed flag value | `--<name>` | the supplied value | `[<the flag's type or format>]` |
  | wrong positional count | `arguments` | the first extra argument, or null when one is missing | `[<positional spelling from Use>]`, or `["no arguments"]` |

  Human mode keeps printing `error: <message>` until task 10 converges it.
  Leaves that already conform (for example `capabilities` with an extra
  argument) keep their text.
- **Unknown commands (Phase 21 section 4; DEC-278, DEC-279).** Add the
  accepted `operation-unavailable` row to `internal/cli/codes.go`:
  invalid-input, not retryable, exit 2, subject null, details
  `{argument:string, value:string|null, allowed:string[]}`. An unknown word at
  any command level returns it with `argument:"command"`, `value:<word>` and
  `allowed:<the subcommands that command's --help lists, sorted>`. The
  message is Cobra's unknown-command text. The JSON envelope `operation`
  follows DEC-279. Human mode prints `error: <message>` and exits 2. It is
  never silent.

## Scope

### Files To Create Or Modify

- Descriptor files for the activation leaves (new); `cmd/fid/activation.go`
- `cmd/fid/input_descriptor_test.go` (fence removal; whole-tree corpora;
  top-level sensitive check)
- Remaining family-local `decodeClosed*Input` decoders and their callers
  (inventory R5)
- Fake/fid conformance assembly in `cmd/fid/integration/` (the Phase 22 fake
  corpus files)
- Grammar and unknown-command conformance: `cmd/fid/canonical_request.go`,
  `cmd/fid/root.go`, `cmd/fid/main.go` if the exit path needs it,
  `internal/cli/codes.go` and its code conformance tests, and every
  `cmd/fid` command file whose leaves bypass `canonical_request.go` for
  grammar failures (for example `issue list --bogus` is silent in JSON mode
  today)
- A whole-tree grammar corpus test (new, in-process over the live tree), and
  retained grammar-failure assertions (convergence with provenance)
- `cmd/fid/integration/testdata/p21_conformance/` if the
  `operation-unavailable` producer changes its manifest or fixtures

### Files To Read For Context

- `project-docs/docs/phases/phase-24/public-interface.md` section 8
- `project-docs/docs/phases/phase-24/input-inventory.md` sections 1, 2 (R5)
- `project-docs/docs/phases/phase-24/decisions.md` DEC-276
- `project-docs/docs/phases/phase-21/public-interface.md` sections 2 (grammar
  rule), 2.2 and 4 (absent commands, `operation-unavailable`) and the code
  table; `project-docs/docs/phases/phase-22/public-interface.md` section 3
  (activation inputs)
- `project-docs/docs/phases/phase-24/decisions.md` DEC-278, DEC-279

### Out Of Scope

- Any input shape change. Human failure lines (task 10): this task fixes the
  JSON failure and keeps Human output nonempty; task 10 sets the line form.

## Acceptance Criteria

- [ ] Reflection conformance covers all 39 leaves. The fence and its test
  are gone, and registry equals the live set. (VAL-011)
- [ ] Help conformance covers all 39 leaves. Full text is pinned for
  `activation component`. (VAL-012)
- [ ] `--example` form, side-effect and typed-template corpora cover all 39
  leaves. (VAL-013, VAL-060, VAL-014)
- [ ] Exclusivity and structural-validation corpora cover all 39 leaves,
  including both `content` variants and a sensitive-field case with
  `value:null`. (VAL-015, VAL-016)
- [ ] The fake/fid conformance table covers the leaf set derived from the
  fake's operation dispatch intersected with the 17 changed leaves (the QA
  plan measured `gate.register`, `envelope.create`, `envelope.complete`,
  `claim.acquire`, `run.start`, `run.close`, `run.abandon`, `gate.run` and
  `run.observe-budget`), with identical acceptance and rejection and no
  one-sided case. (VAL-035)
- [ ] No descriptor marks a top-level field `sensitive`. (I-17)
- [ ] The whole-tree grammar corpus covers every live leaf except
  `journal inspect` and `sql`, in both modes: an unknown flag; a missing
  required flag where the leaf has one; a malformed value where a typed flag
  exists; and a wrong positional count. Each matches the table above, with
  the leaf's operation ID, no effect and nonempty Human stderr. A leaf
  without rows fails the completeness check. (VAL-064)
- [ ] `fid bogus`, `fid concern bogus` and an unknown word under every command
  group return `operation-unavailable` in JSON with the DEC-279 `operation`
  and the details above, and print `error: …` in Human mode, exit 2.
  (VAL-065)
- [ ] Wall times inside budget. (VAL-052)

## First Failing Test

The live `--input` set equals the descriptor registry with no fence. Today 8 activation leaves remain fenced.

## Documentation Updates

- `project-docs/specs/06_DESIGN_PLAN_CLI.md`: the activation input help, the
  grammar-failure details, and `operation-unavailable` for unknown commands.

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
rg -n "func decodeClosed[A-Za-z]*Input" cmd/fid
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

- The `rg` line regenerates R5. No family-local decoder may remain, so it is
  expected to find nothing.

## Expected Non-Zero Commands

- `rg -n "func decodeClosed[A-Za-z]*Input" cmd/fid` - expected `exit_code: 1` - every structural check goes through descriptors

## Notes For Developer

- Report path: `project-docs/docs/phases/phase-24/tasks/task-9-done.md`.
  Commit with `Refs: phase-24/task-9`.
- Red first: record the first failing test before changing product code.
- Retained-test convergence follows DEC-162/165/166. Record each displaced
  assertion with its introducing commit, phase and task. Never infer
  ownership from a filename or task label. Never delete or weaken an
  assertion to pass a gate.
- Any missing registry member, contract ambiguity, or need to change an
  accepted contract shape or assertion is a `BLOCKED` escalation candidate,
  not a local decision (DEC-195).
