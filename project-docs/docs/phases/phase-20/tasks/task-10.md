# Phase 20 - Task 10: Acceptance-criterion and label contracts

## Objective

Migrate acceptance-criterion add/list, label add/list, and Issue label add/remove
to their exact Response Envelope contracts.

## Validation Contract Coverage

- VAL-008, VAL-010, VAL-011, VAL-017, VAL-021, VAL-022, VAL-041, the family
  request rows of VAL-044, and relevant order/limit cases in VAL-045.

## Binding Constraints

- Exact requests/results are in `../contracts.md` sections 3.1-3.5 and 3.7.
- AC and Issue-label Events use Task 6's aggregate revision helper and return
  the resulting Issue revision; golden fixtures must match section 7.
- Canonical label vocabulary is `--colour`/`colour` across CLI, JSON, domain,
  and SQLite projection.
- `issue ac verify` is legacy and must remain unenveloped for Phase 23; it does
  not inherit the migrated list/add behavior accidentally.
- Missing references and duplicates use their specific closed codes and append nothing.
- AC ordinal/text, label token/description/colour, and repeated label inputs
  obey section 3.2.1 boundary and malformed-input rules.
- Remove `--format` only from the migrated leaves in this task.
- Read the resolved project root from Task 2 command context.

## Scope

### Files To Create Or Modify

- `cmd/fid/issue_ac.go`, `label.go`, `issue_label.go` and focused tests.
- `internal/issues` AC/label operations and `internal/labels` result/error paths.

### Files To Read For Context

- `../contracts.md` sections 2.4, 3, and 6.

### Out Of Scope

- Role-specific AC verification, links, workflows, and concerns.

## Acceptance Criteria

- [ ] Six operations accept only exact requests and emit exact result key/type sets using `colour`. (VAL-021, VAL-022, VAL-044)
- [ ] AC/Issue-label mutations advance/return Issue revision and match golden Event fixtures. (VAL-010, VAL-011)
- [ ] Mutations return matching Event identity; reads return `event:null`. (VAL-017)
- [ ] Missing Issue/label and duplicate label/ordinal return specific codes without append. (VAL-021, VAL-022)
- [ ] AC and label value-domain boundaries/malformed cases reject without append. (VAL-021, VAL-022, VAL-044)
- [ ] AC/label arrays use exact order and bounds. (VAL-045)
- [ ] Attribution is recorded without trusted/protected effects; Human success is stdout-only. (VAL-008, VAL-041)
- [ ] `issue ac verify` remains a guarded legacy leaf. (VAL-044)

## First Failing Test

`cmd/fid/issue_ac_test.go` - `TestAcceptanceCriterionAddJSONMatchesContract`.

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

- Completion report: `task-10-done.md`.
