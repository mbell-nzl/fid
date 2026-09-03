# Phase 21 - Task 13: Portable Executor Policy

## Objective

Implement the standard-library process adapter for the exact enforceable policy and reject every unsupported restriction at Registration. Establish no OS containment claim.

## Validation Contract Coverage

- VAL-044: owning observable in the Phase 21 validation contract.
- VAL-045: owning observable in the Phase 21 validation contract.
- VAL-046: owning observable in the Phase 21 validation contract.

## Binding Constraints

- Structured argv/no shell, exact cwd/env, limits, credentials, redaction inputs, and Artifact declarations only.
- Accept only `sandbox:none`, `network:allowed`, and `write:allowed`.
- No runtime dependency addition.
- Final-review MF-2 is documentation-only: the `04_TRD_CONFIG.md` section 2
  example must contain exactly the two implemented `[executor_limits]` keys.
  Do not change the rules paragraph, accepted bounds, or product code.

## Scope

### Files To Create Or Modify

- `internal/evidence/`
- `internal/config/`
- `internal/cli/`
- `cmd/fid/`
- `project-docs/specs/04_TRD_CONFIG.md` section 2 example block only
- `project-docs/docs/phases/phase-21/tasks/task-13-done.md`

### Files To Read For Context

- `project-docs/docs/phases/phase-21/contracts.md §5`
- `project-docs/docs/phases/phase-21/public-interface.md §§2,7`

### Out Of Scope

- Live framing, Receipt persistence, OS isolation, and platform-specific sandboxing.

## Acceptance Criteria

- [ ] Accepted policy is observably enforced. (VAL-044)
- [ ] Every stronger/invalid value rejects without persistence. (VAL-045)
- [ ] No containment claim exists. (VAL-046)
- [ ] The section 2 example is accepted verbatim by a config-reading mutation
  and names only `max_output_bytes` and `max_excerpt_bytes` under
  `[executor_limits]`. (no VAL; final-review MF-2 / QA-91)

## First Failing Test

For the final-review recovery, first run the documented-example command below
against the uncorrected example and preserve its exit-2 failure. After the
bounded documentation edit, rerun the byte-identical command and require exit
0. The original task's positive policy observation and one-field unsupported
policy tables remain unchanged.

## Documentation Updates

- `project-docs/docs/phases/phase-21/tasks/task-13-done.md` - record red/green evidence, deviations, required commands, and spec impact.
- `project-docs/specs/04_TRD_CONFIG.md` - remove only the two unsupported
  `[executor_limits]` example keys identified by final-review MF-2.

## Exact Commands To Run

Run exactly:

```text
bash -c 'set -euo pipefail; repo=$(pwd); scratch=$(mktemp -d); go build -o "$scratch/fid" ./cmd/fid; mkdir -p "$scratch/workspace"; cd "$scratch/workspace"; "$scratch/fid" init --prefix FID >/dev/null; awk "/^## 2[.] /{section=1; next} section && /^...toml$/{fence=1; next} fence && /^...$/{exit} fence{print}" "$repo/project-docs/specs/04_TRD_CONFIG.md" > .fid/config.toml; "$scratch/fid" --json issue create --kind task --title T --actor agent:dev --idempotency-key k1'
go test -race ./internal/evidence ./internal/config ./internal/cli ./cmd/fid/...
rg -n -i --glob '*.go' --glob '!**/*_test.go' -e 'sandboxed' -e 'seccomp' -e 'chroot' -e 'namespace isolation' -e 'privilege (drop|separation)' -e 'network (blocked|denied|disabled|isolated)' -e 'filesystem (isolation|containment)' -e 'isolated (process|environment|execution)' -e 'contained (process|execution)' cmd internal
rg -n -i -e 'sandboxed' -e 'seccomp' -e 'chroot' -e 'namespace isolation' -e 'privilege (drop|separation)' -e 'network (blocked|denied|disabled|isolated)' -e 'filesystem (isolation|containment)' -e 'isolated (process|environment|execution)' -e 'contained (process|execution)' project-docs/docs/phases/phase-21/contracts.md project-docs/docs/phases/phase-21/public-interface.md
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

Every command is required and must terminate with session-captured evidence. Report lint is not a Developer task-gate command; after the report commit, the Phase Runner runs the phase brief's two-argument pre-dispatch control. A Ready QA plan may add a task-specific command; the Tech Lead must amend this section before authorization if it does.

## Expected Non-Zero Commands

- `rg -n -i --glob '*.go' --glob '!**/*_test.go' -e 'sandboxed' -e 'seccomp' -e 'chroot' -e 'namespace isolation' -e 'privilege (drop|separation)' -e 'network (blocked|denied|disabled|isolated)' -e 'filesystem (isolation|containment)' -e 'isolated (process|environment|execution)' -e 'contained (process|execution)' cmd internal` - expected `exit_code: 1` - no production containment claim
- `rg -n -i -e 'sandboxed' -e 'seccomp' -e 'chroot' -e 'namespace isolation' -e 'privilege (drop|separation)' -e 'network (blocked|denied|disabled|isolated)' -e 'filesystem (isolation|containment)' -e 'isolated (process|environment|execution)' -e 'contained (process|execution)' project-docs/docs/phases/phase-21/contracts.md project-docs/docs/phases/phase-21/public-interface.md` - expected `exit_code: 1` - no Phase 21 contract containment claim

## Notes For Developer

- Start only after Group 1 checkpoint PASS (VAL-042 and VAL-043) and all
  feasibility probes.
- Completion report: `project-docs/docs/phases/phase-21/tasks/task-13-done.md`.
- Stop for a protected contract change; ordinary in-scope red/green repair remains this task.
- Do not read held-out scenario material.
- Final-review recovery changes only the section 2 example and this task's
  completion report. Preserve the first failed reproduction, record the root
  cause and later byte-identical PASS, rerun every command above, commit the
  report separately as required, and obtain a fresh independent final re-review.
