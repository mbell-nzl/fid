# Phase 20 - Task 3a: Reason-code schema and redaction evidence

## Objective

Close Review 1's Task 1 Blocking finding by proving the closed reason-code
registry exhaustively: every registered code's exact detail value types,
nullability, nested object and object-list schemas, allowed subject types, and
closed-details redaction. This is evidence recovery for an accepted assertion,
not a contract change (DEC-130).

## Validation Contract Coverage

- VAL-002: machine failure has the exact failure envelope, code-specific
  subject/details, and redaction behavior; no secret, config, or payload
  disclosure. Accepted pass/fail criteria are unchanged.

## Binding Constraints

- `contracts.md` §2.4 and §2.5 are the only authority for codes, categories,
  retryability, subject rules, and detail objects. All 23 registered codes are
  in scope; no sampling, no representative subset, no family ellipsis.
- Test-first. Every added assertion starts red.
- Production changes are allowed **only** where the matrix proves the registry
  accepts an input that §2.4/§2.5 forbid, and only inside `internal/cli`.
  Enforcing equal-length `duplicate-entity` `identity.fields`/`identity.values`
  is pre-authorized by DEC-130 behind a red test. Any other behavior change,
  or any case where the contract looks wrong, is a hard stop and escalates.
- `internal/cli` still imports no Cobra and no domain package.
- `internal-error` remains a backstop: it carries the empty details object and
  no test may treat it as an expected specified failure of another operation.
- `projection-divergence` details are asserted at the registry schema level
  against `contracts.md` §3.5's result object. Do not implement, import, or
  anticipate Task 13's `maintenance.check` behavior.

## Scope

### Files To Create Or Modify

- `internal/cli/codes_test.go` (primary), `internal/cli/envelope_test.go`.
- `internal/cli/codes.go` only under the constraint above.

### Files To Read For Context

- `project-docs/docs/phases/phase-20/contracts.md` §§2.2-2.5 and §3.5.
- `project-docs/docs/phases/phase-20/qa/review-1.md` Task 1 Blocking finding.
- `project-docs/docs/phases/phase-20/qa/validation-contract.md` VAL-002.

### Out Of Scope

- Command migration, domain error mapping, Cobra wiring, Human-mode text.
- Config, init, bootstrap, and project-resolution work (task-3b owns those).
- Any new reason code, category, subject type, or detail field.

## Acceptance Criteria

- [ ] One table-driven test declares, for all 23 codes, the expected category,
      retryability, allowed subject types, and full ordered detail schema
      including nested object and object-list members; adding, removing, or
      re-specifying any code fails the test. (VAL-002)
- [ ] For every code, a fully populated valid payload is accepted, and every
      nullable field is exercised in both its non-null and null states. (VAL-002)
- [ ] For every code and every declared detail field, wrong-type, missing-field,
      extra-field, and null-for-non-nullable payloads are rejected. (VAL-002)
- [ ] Nested members are asserted to the same depth: `duplicate-entity`
      `identity` (member names, string-array types, field-name order, and
      equal-length parallel arrays, including a rejected length mismatch), and
      `projection-divergence` `errors`/`warnings` finding objects and `info`
      object (exact member names and types, wrong-typed and extra nested
      members rejected). (VAL-002)
- [ ] Subject rules are asserted per code in both directions: each code accepts
      exactly its allowed subject types, rejects every other type, the six
      null-subject codes reject a non-null subject, and subject-bearing codes
      reject a nil subject. (VAL-002)
- [ ] Unknown and empty codes reject, and `internal-error` accepts only the
      empty details object. (VAL-002)
- [ ] Redaction is proven as a closed-details rule: payloads carrying
      environment, config, journal-payload, SQL, file-content, credential, or
      stack-trace members reject, and serialized failures for at least one code
      per category contain exactly the four failure keys, the exact `error`
      member set, the declared detail field names only, and one trailing
      newline. (VAL-002)

## First Failing Test

`internal/cli/codes_test.go` -
`TestReasonCodeDetailSchemaIsExactForEveryCode`.

## Documentation Updates

- None expected. `contracts.md` §2.5's equal-length parallel-array clarification
  is already recorded by DEC-130.

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

- Drive the matrix from one declared expectation table compared against
  `RegisteredCodes()`/`LookupCode`, so a registry edit cannot pass by updating
  only the implementation.
- Do not assert on unexported registry internals in a way that makes the test
  restate `codes.go`; assert observable acceptance/rejection and serialized
  output.
- A payload that the registry accepts but §2.4/§2.5 forbid is a finding to fix
  here; a contract that looks wrong is an escalation, not a local edit.
- Completion report: `task-3a-done.md`. Cite VAL-002 explicitly and state which
  Review 1 Blocking finding each new test closes.
