# Phase 20 - Task 18: Journey acceptance and documentation

## Objective

Prove the complete standalone Issue journey, enforce exact retained/legacy
command inventories, and synchronize operator documentation. This task adds
acceptance evidence only; earlier behavior defects are escalated to their owner.

## Validation Contract Coverage

- VAL-035, VAL-036, and VAL-037; close the cross-family portions of VAL-008,
  VAL-017, VAL-041, VAL-044, and VAL-045. VAL-038 is the phase-end Human
  checklist and is supported, not satisfied, by this task.

## Binding Constraints

- The journey runs entirely through `--json`: init, create legal hierarchy,
  add AC/labels/links, list/get/ready, update/transition with revision, provoke a
  typed failure, inspect/check, rebuild, and re-read identical Issue state.
- Drift guard uses exact leaf lists from `../contracts.md` sections 3 and 6; no
  ellipses or inferred command-family membership.
- Rebuild comparison is byte-level on canonical `issue.get` output.
- Do not change product behavior here; a discovered mismatch is a blocker.

## Scope

### Files To Create Or Modify

- `cmd/fid/e2e_p20_issue_journey_test.go` (integration).
- `cmd/fid/contract_surface_test.go` exact command inventory guard.
- `project-docs/docs/playbook.md`, `CHANGELOG.md`; `AGENTS.md` only if Task 16 left drift.

### Files To Read For Context

- `../contracts.md`; roadmap Phase 20; `07_DESIGN_PLAN_FLOWS.md` section 1.

### Out Of Scope

- Product fixes, held-out scenarios, Phase 21 capability advertisement, and live-project migration.

## Acceptance Criteria

- [ ] Exact leaf inventory covers all 42 live executable leaves as 21 retained and 21 legacy. (VAL-035)
- [ ] Every retained leaf accepts only its closed request/`--json`; every legacy leaf rejects root flags. (VAL-035, VAL-044)
- [ ] Journey proves channels/Event identities, attribution boundary, and collection order/limits across families. (VAL-008, VAL-017, VAL-041, VAL-045)
- [ ] Complete journey passes and Issue JSON is byte-identical after rebuild. (VAL-036)
- [ ] Documentation contains no current universal-format or promised-command claim. (VAL-037)
- [ ] Human checklist has one discrete observation for each VAL-038 judgment without pre-answering it. (VAL-038)

## First Failing Test

`cmd/fid/contract_surface_test.go` - `TestLiveLeafCommandsMatchContractInventories`.

## Documentation Updates

- Playbook delivered surface; changelog Phase 20 entry; architecture drift if present.

## Exact Commands To Run

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
go test -tags integration -race ./cmd/fid/integration/...
go test -bench . -benchmem ./...
go build -o /tmp/fid ./cmd/fid
```

## Notes For Developer

- Prove the drift guard fails using a temporary local test-harness edit, revert
  it, and record the permitted harness deviation under the exact-command rule.
- Completion report: `task-18-done.md`.
