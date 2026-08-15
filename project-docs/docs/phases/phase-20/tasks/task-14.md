# Phase 20 - Task 14: Legacy JSON renderer conversion and its caller convergence

> **Amended 2026-08-16 by DEC-142 (resolves ESC-20-011).** The original spec
> deferred this task's own broken callers to Task 15, which made its required
> gate unreachable. Caller convergence for the eight converted leaves is now in
> scope. The six exact commands, thresholds, and assertion strength are
> unchanged.

## Objective

Remove runtime formatter use from the eight legacy JSON reads, encode their
existing semantic results directly with the standard library, and converge every
caller those eight conversions invalidate so the tree is green. This bounds the
behavioral conversion before the shim-removal and framework-deletion tasks.

## Validation Contract Coverage

- VAL-028, VAL-029, and the portion of VAL-031 covering the eight converted
  leaves: their handler flag declarations and their test callers. The shim,
  remaining callers, and whole-tree guards close under Task 15; `internal/format`
  deletion closes under Task 16 (VAL-030).

## Binding Constraints

- Convert exactly `policy show`, `claim list`, `activity list`, `concern list`,
  `stats tokens`, `stats cost`, `sql`, and `journal inspect`.
- Results remain unenveloped and carry tested meaning only; do not introduce a
  new shared renderer, byte-compatibility fixture, or forward contract.
- `journal inspect` remains legacy. Its role/revert/Event shape must not appear
  in the retained operation registry or accept root `--json`. Delete its
  `revert` field and the unused Verdict/RevertVerdict helper/test family.
- Leave `internal/format` present until Task 16 so missed imports remain visible.
- **Caller ownership (DEC-142).** This task owns every caller its own eight
  conversions invalidate. Rewrite each remaining `--format` invocation that
  targets `policy show`, `claim list`, `activity list`, `concern list`,
  `stats tokens`, `stats cost`, `sql`, or `journal inspect` to the converted
  surface, and delete the formatter test helpers those callers leave unused.
  Ownership is per invocation, not per file: one file may hold both Task 14 and
  Task 15 invocations, so leave the others alone.
- **Do not shim.** Do not extend `migratedReadTestArgs` in
  `cmd/fid/issue_test.go`, and do not add any equivalent argument-rewriting
  helper, to cover the eight legacy leaves. The flag is gone; callers call the
  real surface.
- Do not delete or weaken a behavioral assertion in order to drop a flag. A
  converted leaf emits the same JSON its `--format json` path emitted, so the
  existing decode-and-assert bodies survive the flag removal.
- Retain tests whose assertion is that a retired format flag is **rejected**.
  Those are VAL-016 and VAL-031 negative evidence, not callers.

## Scope

### Files To Create Or Modify

- The eight legacy command handlers and their focused unit/CLI tests.
- `internal/journal/inspect.go` and tests, limited to deletion of revert-verdict code.
- Narrow result structs only where existing domain results cannot be encoded directly.
- Unit, e2e, and `cmd/fid/integration` tests that invoke the eight converted
  leaves with `--format`, plus the formatter test helpers those callers leave
  unused. `golangci-lint run ./...` enumerates the unused helpers exactly.

### Files To Read For Context

- `../contracts.md` sections 5.2 and 6.
- `../decisions.md` DEC-142 - caller-ownership rule and the Task 14/15 boundary.

### Out Of Scope

- The `migratedReadTestArgs` compatibility shim and the `ready`, `issue get`,
  `issue list`, and legacy `--section` callers routed through it (Task 15).
- `internal/db/issues_workflow_fk_cli_test.go`, retained parity suites, the
  repository source guard, and the command-tree flag guard (Task 15).
- Deleting `internal/format`, the YAML dependency, or
  `cmd/fid/e2e_v01_format_matrix_test.go` (Task 16); legacy result redesign.

## Acceptance Criteria

- [ ] Each of the eight reads emits valid unenveloped JSON with its previously tested semantic fields. (VAL-029)
- [ ] `journal inspect` has no revert field/helper code, rejects root `--json`, and is absent from the retained operation registry. (VAL-028)
- [ ] The eight live handlers declare no universal format flag and import no formatter package. (VAL-031)
- [ ] Explicit root flags on every converted legacy leaf still reject before its handler. (VAL-028, VAL-029)
- [ ] No test in the repository invokes any of the eight converted leaves with `--format`, except tests asserting that the flag is rejected. (VAL-031)
- [ ] All six exact commands below pass in this session with observed termination, with no assertion, threshold, or command weakened. (VAL-029, VAL-031)

## First Failing Test

`cmd/fid/policy_test.go` - add `TestPolicyShowUsesDirectLegacyJSONWithoutFormatFlag`.

## Documentation Updates

- None expected.

## Exact Commands To Run

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
go test -tags integration -race ./cmd/fid/integration/...
go build -o /tmp/fid ./cmd/fid
```

## Notes For Developer

- Use `json.Encoder` with `SetEscapeHTML(false)` and one trailing newline.
- Completion report: `task-14-done.md`.

### Resuming this task after ESC-20-011 (DEC-142)

- This is a **same-task continuation**, not a fresh start. The preserved
  uncommitted handler and test edits from the first session are correct work;
  build on them.
- The three `BLOCKED` gate results already in `task-14-done.md` are immutable
  failed attempts. Do not rewrite or delete them. Rerun the byte-identical six
  commands, record the later passing attempts alongside them, and add the
  mandatory deviation note: each first failure's verbatim token and observed
  termination, root cause, the repair, and the later result.
- Change the report's `Task outcome` from `BLOCKED` to its true final outcome
  and remove the blocker framing that cited the old Task 15 boundary; cite
  DEC-142 instead.
- `migratedReadTestArgs` covers `ready`, `issue get`, and `issue list` only, so
  those callers still compile and are Task 15's. Your red set is the direct
  callers of the eight leaves plus the helpers they orphan.
- If the caller sweep cannot finish in one session, record `RECOVERY_REQUIRED`,
  not `BLOCKED`. Incomplete in-scope work is a continuation, not an escalation.
