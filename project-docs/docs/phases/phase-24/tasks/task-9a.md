# Phase 24 - Task 9a: Grammar failures and unknown commands across the live tree

## Objective

Deliver two accepted Phase 21 allocations that were never delivered
(DEC-194, DEC-278, DEC-280):
- grammar failures name the offending spelling or value across the whole
  live tree, before any workspace read;
- an unknown word under the root or any command group returns
  `operation-unavailable`. Today it either prints help with exit 0 or exits
  silently.

Every message that embeds caller text is free of control characters. This
task was split from task 9 (DEC-280, TL-02) and runs between tasks 9 and 10.

## Validation Contract Coverage

- VAL-064: grammar failures name the offending spelling or value.
- VAL-065: unknown commands return `operation-unavailable` (DEC-279).
- VAL-052.
- Advances VAL-036. Task 10 pairs these failures with Human lines.

Assertions live in `project-docs/docs/phases/phase-24/qa/validation-contract.md`.

## Binding Constraints

- Phase 21 `public-interface.md` sections 2 (grammar rule), 5 (absent
  commands; the Human-only diagnostics) and 7 (code table; "messages are
  never control flow"). Contract section 9. DEC-278 items 3-4, DEC-279, and
  DEC-280 items 2-6 and 8.
- **Scope.** Every live leaf except `fid journal inspect` and `fid sql`
  (Phase 21 section 5), plus the root and every command group (R17).
  Cobra's generated `help` and `completion` commands have no operation ID
  and are out of scope. The one exception: they are members of the root's
  `allowed`.
- **Envelope.** A grammar failure is `invalid-argument`: invalid-input, not
  retryable, exit 2, subject null, and no read or write. The JSON envelope
  carries the leaf's operation ID (`stop.raise`, never the bare Cobra name
  `raise`). At the root or a group it carries the DEC-279 path (`fid`,
  `concern`, `issue.ac`).
- **Classes (DEC-280 item 4).**

  | Class | `argument` | `value` | `allowed` | `message` |
  |---|---|---|---|---|
  | unknown flag (long or shorthand) | the spelling as typed (`--bogus`, `-x`) | the spelling | the leaf's local long flags, sorted, without `--help`, `--json` and `--project`; `[]` at the root or a group | pflag's text (`unknown flag: --bogus`) |
  | missing required flag | `--<name>` | null | `["--<name>"]` | the detecting check's text (`--input is required`) |
  | malformed flag value: a pflag type parse failure, a flag given without its value, or a value outside its declared domain | `--<name>` | the supplied value, or null when it is missing | the declared domain (below) | pflag's text; for the flag-domain check, `invalid argument "<value>" for "--<name>" flag` |
  | wrong positional count | `arguments` | the first extra argument, or null when one is missing | `[<positional spelling from Use>]`, or `["no arguments"]` | `<command path without "fid ">` + ` accepts no arguments` for a leaf that takes none (the Phase 21 section 1 `capabilities` example); otherwise Cobra's text (`accepts 1 arg(s), received 0`) |

  An extra word on a leaf without subcommands is a wrong positional count,
  never `operation-unavailable`. Today a `cobra.NoArgs` leaf (`ready`,
  `concern list`) reports it as `unknown command "<word>" for "<leaf>"`
  and prints nothing. The `message` is never a family fallback text such as
  `invalid … request` or `… input is invalid`. Today `issue update --bogus`
  prints `error: invalid Issue mutation request` (`issueMutationFlagError`),
  and `issue list --bogus` and `issue get` with no ID are silent in both
  modes. All three become class failures. Leaves that already conform
  (`capabilities extra`) keep their text. `--example` exclusivity failures
  stay exactly as task 3 states.
- **Declared flag domains (DEC-280 item 5).** Each flag with a typed or
  closed value declares its domain once in `cmd/fid`. The declarations are
  keyed by leaf and flag spelling, and both product validation and the
  corpus read them.
  - `allowed`: enumerated members, sorted bytewise; an integer range as
    `"<lo>..<hi>"`; the pattern its row writes (`^[A-Z][A-Z0-9]{0,15}$`,
    `<type>:<id>`); otherwise the pflag type name (`int`, `bool`; a
    free-text string flag missing its value gives `["string"]`).
  - The source is the newest accepted contract row for that leaf (Phase
    21 section 4, or the Phase 22-24 contract that amends it).
  - The check runs before workspace resolution, so a grammar failure is
    identical from a non-project directory. Today some domain checks run
    after the project is opened (`claim list --limit 0` and
    `run list --status bogus` reach the family fallback, `argument:"input"`,
    and `receipt list --outcome bogus` reaches `argument:"request"`).
    Domain-package checks stay unchanged (contract section 12).
  - Completeness: every R18 flag has a declaration (a mechanical check
    over the live tree). Every string flag whose row states a domain has
    one; the report lists each with its row citation.
  - Where a row and the live flag set disagree, the live flag set is the
    corpus scope. The report lists the disagreement and this task does not
    repair it (CF-047). For example, the live `activity list` has `--on` and
    no `--limit`, while the Phase 21 row lists `--ref` and
    `--limit <1..10000>`.
- **Unknown commands (Phase 21 section 5; DEC-278 item 4, DEC-279,
  DEC-280 items 2-3).**
  - Add the accepted `operation-unavailable` row to
    `internal/cli/codes.go`: invalid-input, not retryable, exit 2, subject
    null, details `{argument:string, value:string|null, allowed:string[]}`.
  - One central mechanism in `cmd/fid/root.go` applies to every command with
    subcommands (R17: the root and 28 groups at planning). It replaces the
    five group-local validators in `concern.go`, `export.go`, `issue.go`,
    `issue_ac.go` and `run.go`.
  - An unknown word returns `argument:"command"`, `value:<word>` and
    `allowed:` the subcommands that command's `--help` lists, sorted. The
    root's list includes `completion` and `help`.
  - `message` is exactly `unknown command "<word>" for "<command path>"`
    (Cobra's `CommandPath`, for example `fid issue ac`) on one line, with
    no suggestions. Today Cobra's root text appends `Did you mean this?`
    lines.
  - A command run with no word keeps today's behaviour (Human mode: its help
    on stdout, exit 0).
  - Human mode prints `error: <message>` and exits 2 until task 10 sets the
    line form. Measured baseline: under 23 groups an unknown word prints the
    group help and exits 0 in both modes, so `fid --json activation bogus`
    returns Human text as success. The root and five groups exit 2 silently.
- **Control characters (DEC-280 item 6, PRD-OPS-007).** Add one escape
  function in `internal/cli` in the `visibleLiveText` form
  (`cmd/fid/live_output.go`): `\n`, `\r` and `\t`, and `\uXXXX` for any other
  control character. Apply it to every grammar and unknown-command
  `message`. Text already escaped by Go's `%q` passes through unchanged.
  Task 10's renderer reuses the function. `details.value` keeps the
  supplied text; JSON encoding escapes it, and task 10 escapes the
  `reason:` line.
- **Request-level fallbacks (R19, DEC-280 item 5).** After this task, the
  grammar paths produce class details and no longer use
  `argument:"request"`. They are `canonical_request.go`, the Cobra-error
  fallback in `root.go` `execute`, and `issueMutationFlagError`. The R19
  `rg` then matches exactly the four family fallbacks: `issueMutationFailure`,
  `issueReadFailure`, the `link_response.go` fallback and
  `renderReceiptError`. The report lists each with the causes that still
  reach it, and shows that none is one of the four classes.
- **Flag-combination rules are not a class (CF-047).** These are exclusive
  selectors (`issue create --epic`/`--parent`), exactly-one or at-least-one
  groups (`workflow assign`, `issue update`) and kind-conditional flags.
  They keep today's failures, and the report lists each such rule it finds.

## Scope

### Files To Create Or Modify

- `internal/cli/codes.go` (`operation-unavailable`) and its code
  conformance tests; the `internal/cli` escape function (new), with tests
- `cmd/fid/root.go` (the central unknown-command mechanism; root and group
  unknown flags; the `execute` fallback); `cmd/fid/main.go` if the exit path
  needs it
- `cmd/fid/canonical_request.go`; `cmd/fid/issue_mutation_response.go`
  (`issueMutationFlagError`); `cmd/fid/{concern,export,issue,issue_ac,run}.go`
  (group-local validators)
- Flag-domain declarations (new file in `cmd/fid`), and every leaf
  constructor that registers a declared flag or sets a flag error function
- A whole-tree grammar and unknown-command corpus test (new): in-process
  over the live tree with the R17 and R18 completeness checks, plus
  built-binary rows in both modes. Retained grammar-failure and group-help
  assertions converge with provenance
- `cmd/fid/integration/testdata/p21_conformance/` if the
  `operation-unavailable` producer changes its manifest or fixtures

### Files To Read For Context

- `project-docs/docs/phases/phase-24/public-interface.md` sections 9 and 13
- `project-docs/docs/phases/phase-24/input-inventory.md` section 1, R8,
  R17, R18 and R19
- `project-docs/docs/phases/phase-24/decisions.md` DEC-278, DEC-279, DEC-280
- `project-docs/docs/phases/phase-21/public-interface.md` sections 1
  (`capabilities` example), 2, 4, 5 and 7; the Phase 22 and 23
  `public-interface.md` rows that amend Phase 21 rows
- `cmd/fid/live_output.go` `visibleLiveText`; `cobra@v1.10.2` `args.go`
  (`legacyArgs`, `NoArgs`)

### Out Of Scope

- The Human failure line form (task 10). This task keeps Human stderr
  `error: <message>` and nonempty.
- Flag-combination rules, and disagreements between rows and live flags
  (CF-047).
- `journal inspect`, `sql`, and Cobra's `help` and `completion`.
- Input shapes and `--input` structural failures (tasks 3-9).

## Acceptance Criteria

- [ ] A whole-tree corpus over every in-scope leaf, in both modes. It covers
  an unknown flag (`allowed` = the sorted local long flags), a missing
  required flag where the leaf has one, and a wrong positional count. Each
  row matches the class table and carries the leaf's operation ID. A leaf
  without rows fails the completeness check. (VAL-064)
- [ ] Malformed-value rows for every declared flag domain, including every
  R18 flag and a flag given without its value. A typed-flag parse failure
  exits 2 (today `issue list --has-concerns=maybe` exits 7). The check fails
  when an R18 flag has no declaration. The report lists each string-domain
  declaration with its row citation, and each disagreement between a row and
  the live flag set. (VAL-064)
- [ ] An unknown flag on the root and on every R17 group returns the
  unknown-flag class with the DEC-279 `operation` and `allowed:[]`.
  (VAL-064)
- [ ] An unknown-flag spelling containing ESC gives a message with no raw
  control byte. The text is identical in the JSON `message` and after the
  Human `error: `. (VAL-064)
- [ ] Every row run again from a non-project directory returns the same
  failure, and no row writes anything. (VAL-064)
- [ ] `fid bogus`, `fid isue` and an unknown word under every R17 command
  return `operation-unavailable` as VAL-065 states, in both modes. A command
  without a row fails the completeness check. A group run with no word
  still prints its help (Human mode, exit 0). (VAL-065)
- [ ] The R19 `rg` matches exactly the four family fallbacks. The report
  accounts for each, and lists every flag-combination rule found. (CF-047)
- [ ] Wall times inside budget. (VAL-052)

## First Failing Test

`fid --json activation bogus` returns `operation-unavailable` with envelope `operation:"activation"`, exit 2. Today it prints the group help on stdout and exits 0.

## Documentation Updates

- `project-docs/specs/06_DESIGN_PLAN_CLI.md`: the grammar-failure class
  table, declared flag domains, and `operation-unavailable` for unknown
  commands.

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
rg -n '"argument": *"request"' cmd/fid -g '!*_test.go'
```

The `-timeout 30m` bounds follow DEC-227/259. Record the wall time `go test`
reports for `cmd/fid/integration`, `cmd/fid`, `internal/workflow` and the
slowest other ordinary package. The budgets are tagged integration under
900 s and no ordinary package over 600 s (VAL-052). Report any ordinary
package over the 480 s warning. A breach caused by this task's own new
corpora is a bounded speed repair owned by this task; otherwise stop and
escalate. Never record a timed-out gate as `PASS`.

Required-command rules live in `AGENTS.md` (Commands). Every required command
is reported `PASS` or `BLOCKED` with reason, with its own output token and
observed termination. Focused red/green selections never replace these
commands. Report lint is a Runner control, not a Developer gate.

List any conditional command required by this task and why:

- The `rg` line regenerates R19. It is expected to terminate `exit_code: 0`
  with exactly the four family-fallback matches. A match in
  `canonical_request.go`, in `root.go`, or in `issueMutationFlagError` fails
  the task.
- R17 and R18 come from the inventory section 1 block. This task owns both
  registries, so it reruns the block before close and records both counts.

## Notes For Developer

- Report path: `project-docs/docs/phases/phase-24/tasks/task-9a-done.md`.
  Commit with `Refs: phase-24/task-9a`.
- Red first: record the first failing test before changing product code.
- Retained-test convergence follows DEC-162/165/166. Record each displaced
  assertion with its introducing commit, phase and task. Never infer
  ownership from a filename or task label. Never delete or weaken an
  assertion to pass a gate.
- Any missing registry member, contract ambiguity, or need to change an
  accepted contract shape or assertion is a `BLOCKED` escalation candidate,
  not a local decision (DEC-195).
