# Phase 20 - Task 16: Delete universal formatters and YAML

## Objective

Delete `internal/format`, the format-matrix suite, and `gopkg.in/yaml.v3`, then
add source/module guards that keep universal serialization removed.

## Continuation Status (DEC-143, DEC-144)

The deletion work is **done and committed** at
`afb262ddc3636f6a4c9a05d6b9e77b4aeac276eb`
(`refactor(cli): remove universal formatter framework`). The task has stopped
twice, both times because a control or a spec literal rejected truthful
evidence, never because the work was incomplete. `ESC-20-012` (DEC-143) and
`ESC-20-013` (DEC-144) resolved every cause. This is the **second same-task
continuation**, not a fresh start, and it runs after `task-16b` repairs the
report-lint declaration parser.

Preserved as immutable failed attempts in `task-16-done.md`:

- `rg -n "gopkg.in/yaml.v3" go.mod go.sum` - `BLOCKED`, observed `exit_code: 0`
  on a retained `go.sum` checksum line. **DEC-143 withdrew this command.** Do
  not rerun it. Keep the row, restate it as `WITHDRAWN by DEC-143`, and record
  the two replacement commands below beside it.
- All three report-lint failures (13, 26, and 6 failures). Keep them; the last
  two were control defects that DEC-143 and DEC-144 fixed, not report defects.
- The `BLOCKED` row for `go mod why -vendor -m gopkg.in/yaml.v3` with the
  observed `exit_code: 0` and its vendor-mode line. **DEC-144 corrected the
  literal, not the evidence:** that observed output is a lawful pass, and this
  session re-records the command from its own run.

What changes here: run the amended command list in **this** session, record
this session's own tokens for every command, and finalize the report. A result
captured in an earlier session is `UNVERIFIED`, so the full gate is rerun even
though the product tree is unchanged since `afb262d`.

## Validation Contract Coverage

- VAL-030, VAL-031, VAL-032, VAL-033, and the architecture-doc portion of VAL-037.

## Binding Constraints

- Follow the exact deletion inventory in `../contracts.md` sections 5.1 and 5.3.
- Delete universal TOON/YAML/Markdown rendering; do not move it to a new package.
- Retain only direct Human text in command adapters and canonical Response Envelope JSON.
- Dependency removal requires `go mod tidy`; its intended `go.mod`/`go.sum`
  changes remain part of the task and are validated by the full test/build gate.
- **VAL-032 absence is a module-requirement and build-graph proof (DEC-143).**
  The YAML dependency is absent when `go.mod` requires no yaml module, no yaml
  package is in this module's package graph including its own tests, no
  universal serializer remains in product source, and tidy/tests/build exit 0.
  A `go.sum` checksum line retained for the module graph of an accepted
  retained dependency is bookkeeping, not a fid dependency.
- **No dependency change is approved.** Do not remove, replace, downgrade, or
  pin `github.com/golang-migrate/migrate/v4@v4.19.1`, do not add `exclude` or
  `replace` directives, do not hand-edit `go.sum`, and do not leave the module
  files untidy to satisfy a search.

## Scope

### Files To Create Or Modify

- Delete all `internal/format/` files and `cmd/fid/e2e_v01_format_matrix_test.go`.
- `go.mod`, `go.sum`, source/command guard tests.
- `AGENTS.md` architecture list: replace `internal/format` with `internal/cli` and `internal/project`.
- Continuation: `task-16-done.md` only, unless a rerun command genuinely fails.

### Files To Read For Context

- `../contracts.md` sections 5.1 and 5.3; PRD-FOUND-006.
- `../decisions.md` DEC-143 (VAL-032 authority and the absence-command notation)
  and DEC-144 (the corrected `go mod why -vendor` pass condition).

### Out Of Scope

- Purpose-built future Human briefs, legacy command deletion, and governance stubs.
- `report-lint.sh` itself; `task-16a` and `task-16b` own that control.

## Acceptance Criteria

- [ ] `internal/format` and all product imports are absent. (VAL-030)
- [ ] Universal serializer symbols and format flags are absent from product code. (VAL-031, VAL-033)
- [ ] No yaml module is required by `go.mod` and none is in the main module's
      package graph; `go mod tidy` and the full test/build gate pass. (VAL-032)
- [ ] Architecture documentation names the replacement adapter packages. (VAL-030, VAL-037)
- [ ] Every command below is recorded with this session's own verbatim token and
      observed termination, and the withdrawn command's failed attempt survives
      in the report with its deviation note.

## First Failing Test

`cmd/fid/root_test.go` - extend the surface guard to reject universal serializer
registrations. Already satisfied at `afb262d`
(`TestProductSourceContainsNoUniversalSerializers`, red first, recorded in the
report). The continuation writes no new test unless a rerun command fails.

## Documentation Updates

- `AGENTS.md` architecture package list.

## Exact Commands To Run

```text
go mod tidy
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
go test -tags integration -race ./cmd/fid/integration/...
go build -o /tmp/fid ./cmd/fid
ls internal/format
rg -n "internal/format" cmd internal
rg -ni "writeTOON|writeYAML|writeMarkdown|yaml.Marshal" cmd internal
rg -ni "yaml" go.mod
go mod why -vendor -m gopkg.in/yaml.v3
```

The last five are the QA plan's Task 16 absence gates; they are required here so
the spec lists the full gate. `-vendor` excludes tests of dependencies, which is
exactly the boundary being asserted.

`go mod why -vendor -m gopkg.in/yaml.v3` passes when it exits 0 and its whole
output is one parenthesized `does not need` line naming the module. Under
`-vendor` the Go tool prints
`(main module does not need to vendor module gopkg.in/yaml.v3)`; the non-vendor
wording `(main module does not need module gopkg.in/yaml.v3)` is equally lawful
(DEC-144). Record the line you actually observe as the token. An import chain,
any extra line, or a non-zero exit is a **real VAL-032 failure**.

List any conditional command required by this task and why:

- None.

## Expected Non-Zero Commands

Each declaration is one line: the exact command in backticks, then
`- expected \`exit_code: <N>\``. Any trailing prose is ignored.

- `ls internal/format` - expected `exit_code: 2` - the deleted formatter
  directory must not exist.
- `rg -n "internal/format" cmd internal` - expected `exit_code: 1` - no product
  import of the deleted package survives.
- `rg -ni "writeTOON|writeYAML|writeMarkdown|yaml.Marshal" cmd internal` - expected `exit_code: 1`
  (no universal serializer survives in product source)
- `rg -ni "yaml" go.mod` - expected `exit_code: 1` - no yaml module is required,
  directly or indirectly.

## Notes For Developer

- Record each of the four commands above in the declared form:
  ``- `<cmd>` - PASS - `exit 1 (no output)` (expected `exit_code: 1`; terminated `exit_code: 1`)``,
  with nothing after the closing parenthesis. Put any explanation on the next
  indented line.
- If `go mod why -vendor` prints an import chain instead of the
  `does not need` line, that is a **real VAL-032 failure**: record it, stop, and
  raise a blocker candidate. Do not substitute another command, and do not touch
  the migration dependency. A `does not need` line whose wording differs from
  both forms above is a notation question, not a failure: record it verbatim as
  `PASS` and add a note.
- Before finalizing, you may run
  `project-docs/tools/report-lint.sh <report> <spec>` against the committed
  report as a self-check. It is not a gate row here and the Phase Runner still
  runs it independently. If it fails on truthful evidence, stop and report the
  failure; never reword, reformat, or delete a real result to satisfy it.
- If `ls` reports a status other than 2 on this machine, record the real status
  and report the mismatch; do not edit the declaration to match.
- Completion report: `task-16-done.md`. Every command row carries **this**
  session's result, so the earlier session's `PASS` rows are replaced by the
  rerun ones and the three absence rows move into the declared form above -
  same command, same real status, lawful notation.
- What is immutable and must survive: the `BLOCKED` row for the withdrawn
  `rg -n "gopkg.in/yaml.v3" go.mod go.sum`, both report-lint failures, the
  earlier red guard attempts, and their deviation notes. Add a deviation note
  covering the withdrawal and the notation change, citing DEC-143.
- Final outcome must be `COMPLETE` with `Category: none` and
  `Human attention required: no`, or a truthful non-complete outcome if a rerun
  command genuinely fails.
