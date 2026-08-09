# Phase 20 - Task 3b: Strict-config and init boundary evidence

## Objective

Close Review 1's Task 2 Blocking finding by proving the strict config schema and
every documented value domain at its just-inside and just-outside boundary,
through both the parser/writer and `project.init`. Also commit the missing
dangling-symlink collision regression and relieve the config file-size breach.
This is evidence recovery for accepted assertions, not a contract change
(DEC-130).

## Validation Contract Coverage

- VAL-043: config accepts exactly the documented version, sections, keys,
  ID/prefix/path domains, and signed-64-bit executor/maintenance ranges;
  all just-inside values pass and every unknown, duplicate, version, prefix,
  path, relationship, or limit violation rejects.
- VAL-042 and VAL-002 (partial): an existing dangling `.fid` symlink is never
  replaced and fails with the typed collision envelope.
- VAL-039 (supporting): activated config round-trips exactly the target keys.

Accepted pass/fail criteria for all of the above are unchanged.

## Binding Constraints

- `04_TRD_CONFIG.md` §2 and `contracts.md` §2.1/§3.2.1 are the only authority
  for the schema and domains. Boundaries are exact: prefix `^[A-Z][A-Z0-9]{0,15}$`,
  project ID non-empty opaque ASCII at most 128 bytes, `artifact_root` valid
  UTF-8 `1..4096` bytes and contained under `.fid/`, executor byte limits
  `1..9223372036854775807` with excerpt no larger than output, maintenance grace
  `0..9223372036854775807`.
- Every numeric boundary is proven with an explicit just-inside and
  just-outside case, including `math.MaxInt64` accepted and an over-`MaxInt64`
  literal rejected as out of range rather than silently wrapped, truncated, or
  parsed as a float.
- Test-first. Every added assertion starts red.
- Production changes are allowed **only** where the matrix proves the parser
  accepts or silently ignores an input that the specification forbids, and only
  inside `internal/config`, `internal/bootstrap`, or `cmd/fid` as the failing
  case requires. A specification that looks wrong is a hard stop and escalates.
- The strict generated-config grammar stays dependency-free. A perceived need
  for general third-party TOML support is a dependency-budget escalation.
- Do not weaken, rename, or alias any key; `prefix` and
  `started_without_claim_grace_seconds` must still reject.
- The symlink regression uses a disposable fixture and asserts exact filesystem
  state; no timing, sleeps, or process interruption.

## Scope

### Files To Create Or Modify

- `internal/config/config_test.go` (primary boundary matrix).
- `internal/config/project_config.go` (new; pure move of the strict project
  config parser/writer out of `config.go`, no API or behavior change).
- `internal/config/config.go` (only the corresponding removal, plus a repair
  under the constraint above).
- `cmd/fid/init_test.go` and/or `internal/bootstrap/init_test.go` for prefix
  boundary and dangling-symlink collision cases.

### Files To Read For Context

- `project-docs/specs/04_TRD_CONFIG.md` §§1-2 and §5.
- `project-docs/docs/phases/phase-20/contracts.md` §2.1 and §3.2.1.
- `project-docs/docs/phases/phase-20/qa/review-1.md` Task 2 Blocking finding and
  Task 3 Medium improvement.
- `project-docs/docs/phases/phase-20/qa/validation-contract.md` VAL-039/042/043.

### Out Of Scope

- Reason-code registry evidence (task-3a owns it).
- `fid rebuild`, authority mapping, lock typing, and any Task 4 work.
- Legacy `pricing.toml`, which stays separate and expires in Phase 25.
- Splitting `internal/bootstrap/init.go`; that is CF-028, trigger-based.

## Acceptance Criteria

- [ ] A table-driven parser test proves the structural schema in both
      directions: the exact target config parses, and missing/unsupported/
      non-integer `config_version`, unknown top-level key, unknown section,
      unknown key in each of the four sections, duplicate key, duplicate
      section, missing required section, and missing required key each reject
      with no silent ignore. (VAL-043)
- [ ] `[project].id` boundary matrix: 1-byte and 128-byte values accepted;
      empty, 129-byte, non-ASCII, and control-character values rejected. (VAL-043)
- [ ] `[project].issue_prefix` boundary matrix: 1-character and 16-character
      values accepted; empty, 17-character, lowercase, leading-digit,
      hyphen/underscore, and non-ASCII values rejected. Legacy `prefix` still
      rejects with no alias. (VAL-043, VAL-039)
- [ ] `[storage].artifact_root` matrix: 1-byte and 4096-byte relative
      normalized values accepted; empty, 4097-byte, absolute, `.`, lexical `..`
      escape, and invalid-UTF-8 values rejected, and a symlink component that
      resolves outside `.fid/` rejects. (VAL-043)
- [ ] `[executor_limits]` matrix: `1` and `math.MaxInt64` accepted, excerpt
      equal to output accepted; `0`, negative, over-`MaxInt64`, float, quoted,
      and excerpt-exceeds-output-by-one rejected. (VAL-043)
- [ ] `[maintenance].abandoned_attempt_grace_seconds` matrix: `0` and
      `math.MaxInt64` accepted; `-1`, over-`MaxInt64`, and non-integer rejected.
      `started_without_claim_grace_seconds` still rejects. (VAL-043, VAL-027)
- [ ] Writer round-trip is byte-stable for the exact target config and for a
      `MaxInt64`-valued config, proving no precision loss and no extra or
      reordered key. (VAL-043, VAL-039)
- [ ] `fid init --prefix` proves the same prefix domain end to end: 1-character
      and 16-character accepted; 17-character, lowercase, and leading-digit
      rejected with `invalid-argument` before staging, leaving no `.fid/` and no
      staging residue. Derived-prefix cases prove an acceptable basename is
      uppercased and accepted while a non-conforming basename rejects and
      requires an explicit prefix. (VAL-043, VAL-044, VAL-042)
- [ ] The activated config from a successful init round-trips through the
      strict parser with exactly the target keys. (VAL-039, VAL-043)
- [ ] A committed regression proves an existing dangling `.fid` symlink is not
      replaced: the typed `project-already-initialized` failure envelope is
      returned, the symlink is unchanged, and no staging residue remains.
      (VAL-042, VAL-002)
- [ ] `internal/config/config.go` and every file it is split into are at or
      under the 400-line non-test budget, with no exported API change and no
      behavior change in the move. (no VAL; `AGENTS.md` size budget)

## First Failing Test

`internal/config/config_test.go` -
`TestParseProjectConfigBoundaryDomains` (start with the `MaxInt64` executor
limit and one-character prefix rows).

## Documentation Updates

- None expected; `04_TRD_CONFIG.md` §2 and `contracts.md` §3.2.1 already state
  every domain proven here.

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

- None.

## Notes For Developer

- Write the boundary matrix as data: one row per case with the exact input
  literal and expected accept/reject outcome, so a reviewer can read the
  just-inside and just-outside pair for every domain in one place.
- Prove rejection by outcome, not by message text matching, except where the
  contract names the reason code.
- Do the file move last, as a separate mechanical step, so the evidence diff
  stays readable.
- Completion report: `task-3b-done.md`. Cite VAL-043 explicitly, list the
  boundary pairs proven, and state which Review 1 findings each change closes.
