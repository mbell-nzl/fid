# Phase 20 - Task 3c: Integration-journey convergence and honest gate baseline

## Objective

Close Review 1b's two Blocking findings: add the missing rejected
`max_excerpt_bytes` boundary rows, and converge the stale integration journeys
that still initialize with a non-conforming derived prefix. Then produce one
clean full-gate run on the current tree captured with real terminal tokens
(DEC-131).

## Validation Contract Coverage

- VAL-043: config accepts exactly the documented domains and signed-64-bit
  executor ranges, with qualifying and non-qualifying evidence for **every**
  numeric boundary, including `max_excerpt_bytes` itself.
- VAL-044 (supporting): `project.init` rejects a non-conforming derived prefix
  before any effect; retained journeys supply a valid explicit prefix.

## Binding Constraints

- **The product is correct and does not change.** `contracts.md` §3.2.1 requires
  an omitted `--prefix` to be derived from the directory basename and accepted
  only if it matches `^[A-Z][A-Z0-9]{0,15}$`. `t.TempDir()` roots such as `001`
  must reject. Relaxing that validation is a contract change: hard stop and
  escalate instead.
- Fixtures converge to the contract; assertions never weaken. Supplying an
  explicit valid `--prefix` is convergence. Deleting, skipping, narrowing, or
  loosening any journey assertion is not, and neither is dropping a scenario.
- The exact required command is `go test -tags integration -race ./...` from the
  repository root. It is never narrowed to a package or a `-run` selector for
  gate purposes.
- Record every gate with a **real verbatim token from that command's own
  output**. `go test ./...` prints a line per package, so `exit 0 (no output)`
  is never a valid token for it (`AGENTS.md`, Commands). If a gate fails, it is
  `FAIL`; if it cannot run, it is `BLOCKED` with reason.
- `runInitCommand`'s `--prefix FID` injection stays (DEC-131 accepted it). Do
  not remove it and do not extend that pattern to the binary-driven journeys.

## Scope

### Files To Create Or Modify

- `internal/config/config_test.go` - rejected `max_excerpt_bytes` rows.
- `cmd/fid/e2e_v01_closure_test.go`, `cmd/fid/e2e_v02a_workflow_authoring_test.go`,
  and any other integration/e2e fixture that initializes with a non-conforming
  derived prefix.
- No product file. A required product change is an escalation.

### Files To Read For Context

- `project-docs/docs/phases/phase-20/qa/review-1b.md` (both Blocking findings).
- `project-docs/docs/phases/phase-20/decisions.md` DEC-131.
- `project-docs/docs/phases/phase-20/contracts.md` §3.2.1; `04_TRD_CONFIG.md` §2.

### Out Of Scope

- Re-running historical gates at old commits; DEC-131 replaces that with one
  honest current-tree baseline.
- `report-lint.sh` and completion-report corrections (task-3d owns them).
- Any Task 4 work.

## Acceptance Criteria

- [ ] `max_excerpt_bytes` has its own rejected rows for `0`, a negative integer,
      an over-`MaxInt64` literal, a float, and a quoted value, alongside the
      retained accepted minimum, maximum, and equal-to-output rows. (VAL-043)
- [ ] Provenance for the integration failure is established and recorded: the
      commit that introduced strict derived-prefix rejection is cited from
      `git log -S` or `git log --diff-filter=A`, together with the statement of
      whether it predates task-3a/3b. (no VAL; `AGENTS.md` provenance rule)
- [ ] Every integration/e2e fixture that initializes a workspace supplies a
      prefix that satisfies the contract grammar, and a repository-wide sweep
      confirms no remaining bare `init` call in a non-conforming root. (VAL-044)
- [ ] Each repaired journey keeps every pre-existing behavioral assertion; the
      completion report names each repaired file and states what changed. (VAL-044)
- [ ] `go test -tags integration -race ./...` runs from the repository root and
      is recorded with its real terminal token. (no VAL; required gate)
- [ ] All six gate commands are recorded with real tokens from this session; no
      full-suite `go test` line uses `exit 0 (no output)`. (no VAL; `AGENTS.md`)

## First Failing Test

`internal/config/config_test.go` - the new rejected `max_excerpt_bytes` rows in
`TestParseProjectConfigBoundaryDomains`; then reproduce the failing
`go test -tags integration -race ./...` before repairing the fixtures.

## Documentation Updates

- None expected. DEC-131 already records the disposition.

## Exact Commands To Run

Run exactly:

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
go test -tags integration -race ./...
go build -o /tmp/fid ./cmd/fid
```

Required-command rules live in `AGENTS.md` (Commands): every required command
is reported `PASS`, or `BLOCKED` with reason; missing, silently skipped, or
narrower substituted commands are failed gates.

List any conditional command required by this task and why:

- `git log -S` / `git log --diff-filter=A` for the provenance criterion above;
  a provenance claim without a cited introducing commit is not evidence.

## Notes For Developer

- **This gate is slow: expect several minutes wall clock, roughly three of them
  in `cmd/fid` alone, with no output at all until packages finish.** Go compiles
  everything first and prints one line per package as it completes. Do not kill
  it. A killed or timed-out run is `BLOCKED` with that reason, never `PASS`.
  Capture it with `2>&1 | tee` so the real tail survives.
- The Human-verified failure signature is three `cmd/fid` tests, each
  `.../001/fid init: exit status 1` with `error: invalid project initialization
  argument`, and `ok` for every other package. If you see a different set,
  something has changed since DEC-131 and you should say so.
- Start by reproducing the failure, not by fixing it. Paste the real failing
  output tail into the report before the repair.
- If a journey cannot be converged without changing what it asserts, stop and
  escalate; that is a contract question, not a fixture question.
- This task's whole point is a truthful record. A gate you did not run in this
  session is `UNVERIFIED`, which is not a passing state.
- Completion report: `task-3c-done.md`.
