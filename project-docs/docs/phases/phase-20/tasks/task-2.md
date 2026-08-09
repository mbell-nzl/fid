# Phase 20 - Task 2: Project selection and config namespace

## Objective

Implement project-root discovery, explicit selection, the strict target config
model/parser, transitional bootstrap emission of that strict config, and the
temporary root-flag rejection rule. Resolve the workflow flag collision before
registering global `--project`; Task 3 owns staged init and the `project.init`
Response Envelope.

## Validation Contract Coverage

- VAL-003, VAL-005, VAL-043, the config-model half of VAL-039, the fixed-policy
  wiring half of VAL-027, the flag-collision half of VAL-026, and the
  root-dispatch portion of VAL-044.

## Binding Constraints

- Implement `../contracts.md` sections 2.1 and 3.2 for discovery, strict config,
  root flags, and the closed request registry.
- Root `--project` and `--json` flags must never be silently ignored. A section
  6 legacy leaf rejects either explicitly supplied flag before execution.
- Discovery stops at the containing Git boundary and does not follow symlinks
  outside the selected root.
- The config model contains the exact documented project/storage/executor-limit/
  maintenance keys. Unknown section/key, duplicate, unsupported version, and
  pre-convergence aliases reject instead of being ignored.
- Artifact root is normalized and contained within `.fid/`; executor byte
  ceilings are positive with excerpt not above output, and grace is non-negative.
- `[project].id` is an opaque stored value, never inferred from path. Rename
  `[project].prefix` to `[project].issue_prefix` without an alias.
- Update the existing direct bootstrap path and its fixtures to emit/read the
  strict target config so this task's full gate is green. Task 3 changes how
  those same values are staged and activated; it does not introduce a second
  transitional config model.
- Removing `started_without_claim_grace_seconds` must not make the retained
  maintenance command consume future `abandoned_attempt_grace_seconds`. Wire
  the legacy check to the named fixed 300-second constant; Task 13 contracts
  and fully validates its output.
- Rename the existing local `workflow assign --project` boolean to
  `--project-wide` in this task before registering the persistent root flag;
  Task 12 later migrates the command's transport/result contract.

## Scope

### Files To Create Or Modify

- `internal/project/resolve.go`, `internal/project/resolve_test.go`.
- `internal/config/config.go` and tests; `internal/bootstrap` config emission,
  call sites, and init fixtures only as required to make strict config live.
- `cmd/fid/root.go` and root dispatch/flag tests.
- `cmd/fid/workflow.go` and assignment flag tests, limited to the collision rename.
- `cmd/fid/maintenance.go` and its config-wiring test, limited to removing the
  retired config-field dependency and using the fixed legacy constant.

### Files To Read For Context

- `04_TRD_CONFIG.md` sections 1-2; `../contracts.md` sections 2, 3, and 6.

### Out Of Scope

- Staged/atomic initialization and activation, `project.init` transport, `fid rebuild`,
  authority mapping, lock typing, and workflow transport/result migration beyond
  the required flag-collision rename.

## Acceptance Criteria

- [ ] Nested discovery and explicit project selection follow all boundary rules. (VAL-003)
- [ ] Config model contains exact target keys, uses `issue_prefix` not `prefix`,
      and round-trips a supplied opaque project identity. (VAL-039, VAL-043)
- [ ] The existing direct init path and all bootstrap fixtures emit/read only
      the strict target config before Task 3 changes activation mechanics. (VAL-039, VAL-043)
- [ ] Maintenance no longer reads the retired Claim-grace key or aliases the
      future Attempt-grace value; its existing check receives fixed 300. (VAL-027, VAL-043)
- [ ] Unknown/duplicate/unsupported, escaping-path, and invalid-limit config input rejects instead of being ignored. (VAL-043)
- [ ] `workflow assign --project-wide` retains project-scope behavior and no
      local boolean flag shadows global `--project <path>`. (VAL-026, VAL-044)
- [ ] Explicit `--project` or `--json` on every legacy leaf rejects before its handler runs. (VAL-005)
- [ ] A target operation not yet migrated rejects root flags rather than pretending support. (VAL-005, VAL-044)

## First Failing Test

`internal/project/resolve_test.go` - `TestResolveStopsAtContainingGitRoot`.

## Documentation Updates

- None expected; the accepted contract already records the behavior.

## Exact Commands To Run

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
go test -tags integration -race ./...
go build -o /tmp/fid ./cmd/fid
```

## Notes For Developer

- Mark migrated leaves with one central command annotation/registry consumed by
  root pre-run handling; do not maintain unrelated ad hoc flag checks.
- Keep the strict generated-config grammar dependency-free; a need for general
  third-party TOML support is a dependency-budget escalation, not an implicit addition.
- Completion report: `task-2-done.md`.
