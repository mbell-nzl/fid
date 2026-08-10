# Phase 20 - Task 14: Legacy JSON renderer conversion

## Objective

Remove runtime formatter use from the eight legacy JSON reads and encode their
existing semantic results directly with the standard library. This bounds the
behavioral conversion before the separate test-caller cleanup and framework
deletion tasks.

## Validation Contract Coverage

- VAL-028, VAL-029, and the runtime portion of VAL-031.

## Binding Constraints

- Convert exactly `policy show`, `claim list`, `activity list`, `concern list`,
  `stats tokens`, `stats cost`, `sql`, and `journal inspect`.
- Results remain unenveloped and carry tested meaning only; do not introduce a
  new shared renderer, byte-compatibility fixture, or forward contract.
- `journal inspect` remains legacy. Its role/revert/Event shape must not appear
  in the retained operation registry or accept root `--json`. Delete its
  `revert` field and the unused Verdict/RevertVerdict helper/test family.
- Leave `internal/format` present until Task 16 so missed imports remain visible.

## Scope

### Files To Create Or Modify

- The eight legacy command handlers and their focused unit/CLI tests.
- `internal/journal/inspect.go` and tests, limited to deletion of revert-verdict code.
- Narrow result structs only where existing domain results cannot be encoded directly.

### Files To Read For Context

- `../contracts.md` sections 5.2 and 6.

### Out Of Scope

- Broad e2e/parity test caller cleanup (Task 15), formatter deletion (Task 16),
  or legacy result redesign.

## Acceptance Criteria

- [ ] Each of the eight reads emits valid unenveloped JSON with its previously tested semantic fields. (VAL-029)
- [ ] `journal inspect` has no revert field/helper code, rejects root `--json`, and is absent from the retained operation registry. (VAL-028)
- [ ] The eight live handlers declare no universal format flag and import no formatter package. (VAL-031)
- [ ] Explicit root flags on every converted legacy leaf still reject before its handler. (VAL-028, VAL-029)

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
