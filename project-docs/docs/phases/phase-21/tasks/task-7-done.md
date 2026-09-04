# Phase 21 - Task 7: Generic Authority Evaluator - Completion Report

## Workflow outcome

- Task outcome: COMPLETE
- Category: none
- Condition: None.
- Recovery attempts: 0
- Next action: Phase Runner review and dispatch decision.
- Human attention required: no

## Summary of what was implemented

- Added the generic `internal/policy` Envelope Grant evaluator, which fail-closes unless the active Actor capability, action, Issue, Envelope revision, and Grant selector all match exactly.
- Added a thin `internal/workflow` adapter that resolves the Actor from the finalized active Activation and loads the requested Envelope before delegating to the policy evaluator.
- Added the single-dimension deny matrix and project-local `developer` attribution adversarial test without consulting the V1 role matrix.

## Files created or modified

- `internal/policy/envelope_authority.go`
- `internal/workflow/authority.go`
- `cmd/fid/authority_test.go`
- `cmd/fid/envelope_test.go`
- `project-docs/docs/phases/phase-21/tasks/task-7-done.md`

## Documentation updated

- `project-docs/docs/phases/phase-21/tasks/task-7-done.md`

## Validation Contract Coverage

- VAL-023: PASS - `TestEnvelopeAuthorityRequiresEveryBindingAndRejectsRoleAttribution` in `cmd/fid/authority_test.go:12` permits its all-match `claim.acquire` request at `:50`, then rejects independent actor, action, capability, Issue, Envelope, and revision changes at `:54` plus the no-active-Activation case at `:72`.
- VAL-024: PASS - `cmd/fid/authority_test.go:77` invokes every protected operation with the seeded project-local `developer` attribution and `:101` asserts each result is `*policy.EnvelopeAuthorityDeniedError`.

## Tests run

- Red evidence: `go test -race ./cmd/fid -run '^TestEnvelopeAuthorityRequiresEveryBindingAndRejectsRoleAttribution$'` printed `FAIL` before the generic evaluator API existed (terminated `exit_code: 1`).
- Green evidence: `go test -race ./cmd/fid -run '^TestEnvelopeAuthorityRequiresEveryBindingAndRejectsRoleAttribution$'` printed `ok  \tgithub.com/mbell-nzl/fid/cmd/fid\t4.595s` (terminated `exit_code: 0`).
- Real interface exercise: `TestEnvelopeAuthorityRequiresEveryBindingAndRejectsRoleAttribution` invokes the real `fid --json envelope create` and `fid --json envelope activate` entrypoints before evaluating authority; the green focused test printed `ok  \tgithub.com/mbell-nzl/fid/cmd/fid\t4.595s` (terminated `exit_code: 0`). No task-7 CLI command was added because the evaluator is the reusable internal boundary for later protected operations.

## Exact quality-gate results

- `go test -race ./internal/policy ./internal/workflow ./internal/validation ./internal/cli ./cmd/fid/...` - PASS - `ok  \tgithub.com/mbell-nzl/fid/cmd/fid\t152.041s` (terminated `exit_code: 0`)
- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -coverprofile=coverage.out ./...` - PASS - `ok  \tgithub.com/mbell-nzl/fid/internal/workflow\t(cached)\tcoverage: 38.4% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total:\t\t\t\t\t\t\t\t\t\t(statements)\t\t\t\t\t70.5%` (terminated `exit_code: 0`)

## Deviations from spec

- Expected TDD red stage: the focused command first printed `FAIL` (terminated `exit_code: 1`) because the generic evaluator API was absent. The root cause was the deliberately unimplemented task behavior; the repair added the policy evaluator and Activation/Envelope adapter. The byte-identical later command printed `ok  \tgithub.com/mbell-nzl/fid/cmd/fid\t4.595s` (terminated `exit_code: 0`).
- The first implementation attempt printed `FAIL` (terminated `exit_code: 1`) because importing `internal/workflow` into `internal/policy` formed a package cycle. The root cause was the existing workflow-to-policy dependency; the repair kept policy as a pure evaluator and moved Activation/Envelope loading into a workflow adapter. The byte-identical later command printed `ok  \tgithub.com/mbell-nzl/fid/cmd/fid\t4.595s` (terminated `exit_code: 0`).

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `817ad7c50469c51dc12f88a4bc081ecd7a384d7e`

## Commit subject

- `feat(policy): add envelope authority evaluator`
