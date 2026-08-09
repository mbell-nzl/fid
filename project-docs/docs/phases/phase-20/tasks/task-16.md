# Phase 20 - Task 16: Delete universal formatters and YAML

## Objective

Delete `internal/format`, the format-matrix suite, and `gopkg.in/yaml.v3`, then
add source/module guards that keep universal serialization removed.

## Validation Contract Coverage

- VAL-030, VAL-031, VAL-032, VAL-033, and the architecture-doc portion of VAL-037.

## Binding Constraints

- Follow the exact deletion inventory in `../contracts.md` sections 5.1 and 5.3.
- Delete universal TOON/YAML/Markdown rendering; do not move it to a new package.
- Retain only direct Human text in command adapters and canonical Response Envelope JSON.
- Dependency removal requires `go mod tidy`; its intended `go.mod`/`go.sum`
  changes remain part of the task and are validated by the full test/build gate.

## Scope

### Files To Create Or Modify

- Delete all `internal/format/` files and `cmd/fid/e2e_v01_format_matrix_test.go`.
- `go.mod`, `go.sum`, source/command guard tests.
- `AGENTS.md` architecture list: replace `internal/format` with `internal/cli` and `internal/project`.

### Files To Read For Context

- `../contracts.md` sections 5.1 and 5.3; PRD-FOUND-006.

### Out Of Scope

- Purpose-built future Human briefs, legacy command deletion, and governance stubs.

## Acceptance Criteria

- [ ] `internal/format` and all product imports are absent. (VAL-030)
- [ ] Universal serializer symbols and format flags are absent from product code. (VAL-031, VAL-033)
- [ ] YAML dependency is absent after `go mod tidy`, and the full test/build gate passes. (VAL-032)
- [ ] Architecture documentation names the replacement adapter packages. (VAL-030, VAL-037)

## First Failing Test

`cmd/fid/root_test.go` - extend the surface guard to reject universal serializer registrations.

## Documentation Updates

- `AGENTS.md` architecture package list.

## Exact Commands To Run

```text
go mod tidy
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
go test -tags integration -race ./...
go build -o /tmp/fid ./cmd/fid
rg -n "internal/format" cmd internal
rg -ni "writeTOON|writeYAML|writeMarkdown|yaml.Marshal" cmd internal
```

## Notes For Developer

- Both `rg` commands use no-match (`exit 1 (no output)`) as the pass condition.
- Completion report: `task-16-done.md`.
