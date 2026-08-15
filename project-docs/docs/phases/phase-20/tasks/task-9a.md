# Phase 20 - Task 9a: Mutation concurrency, transition, and attribution evidence

## Objective

Close Review 3's three Task 9 Blocking findings by proving the accepted
VAL-018, VAL-019, VAL-040, and VAL-041 assertions exhaustively across
`issue create`, `issue update`, and `issue transition`, plus the transition
rows of VAL-017 and VAL-044. This is evidence recovery for accepted
assertions, not a contract change (DEC-140).

## Validation Contract Coverage

- VAL-018: three-state (absent / matching / stale) expected-revision behavior
  for **both** update and transition, with journal counts.
- VAL-019: the four-way canonical transition table with journal counts and
  exact `transition-not-permitted` details.
- VAL-040: both public owner flags unknown on **both** canonical mutations, no
  owner change reachable through a canonical result, retained internal
  default-owner path unbroken.
- VAL-041: local attribution recorded with no trusted Evidence, Human Outcome,
  or protected-authority entity or effect, across create, update, transition.
- VAL-017 and VAL-044 (transition and owner rows only). Accepted pass/fail
  criteria for every listed row are unchanged.

## Binding Constraints

- `../contracts.md` §2.4, §3.2, §3.2.1, §3.6, and §7 are the only authority for
  codes, details, requests, value domains, revision, and Event effects.
- **No contract change.** No new or removed flag, reason code, detail member,
  result field, Event family, or value domain. No acceptance criterion,
  threshold, or VAL row is weakened.
- **Test-led recovery.** Task 9 implemented this behavior; what is unproven is
  whether it is correct. If a new assertion fails, that is a real defect: fix it
  minimally inside the existing contract using red/green, and record a deviation
  note quoting the verbatim failing output. If the contract itself looks wrong,
  that is a hard stop — escalate, do not edit the contract.
- **Observable evidence only.** Assert decoded canonical JSON from the command
  boundary, journal line counts taken before and after each case, and projection
  reads. Do not assert on unexported internals or source constants, and do not
  parse human message substrings for canonical JSON failures.
- Every rejection case asserts zero journal append **and** unchanged projected
  state, not one or the other.
- Task 9's implementation commit `86d927d` stands and is not reopened. This task
  is additive. Existing legacy Human-mode transition tests stay; do not delete,
  retarget, or weaken any test to make a new one pass.
- If a required cell below cannot be constructed through any supported surface,
  stop and escalate with the reason. Never silently omit a cell.

## Scope

### Files To Create Or Modify

- `cmd/fid/issue_transition_test.go`, `cmd/fid/issue_update_test.go`
  (primary), and focused helpers they need.
- `internal/issues/update_test.go`, `internal/issues/transition_test.go`,
  `internal/issues/create_test.go` where a domain-level assertion is the
  observable one.
- Product files under `cmd/fid/` and `internal/issues/` **only** under the
  defect clause above.

### Files To Read For Context

- `../contracts.md` §2.4, §3.2, §3.2.1, §3.6, §7.
- `../qa/review-3.md` - the three Task 9 Blocking findings.
- `../qa/validation-contract.md` rows VAL-017, VAL-018, VAL-019, VAL-040,
  VAL-041, VAL-044; the Evidence and Pass columns are binding in full.
- `task-9.md` and `task-9-done.md`.

### Out Of Scope

- Acceptance criteria, labels, links, workflows, maintenance (Tasks 10-13).
- Read contract (Task 8), hierarchy invariant (Task 7), formatter removal.
- Idempotency keys, Work Envelope authority, and any Phase 21 capability.

## Acceptance Criteria

- [ ] One table-driven optimistic-concurrency test asserts, at the canonical
      `--json` command boundary, all three `--expected-revision` states for
      **both** `issue update` and `issue transition` (six cells): absent
      succeeds, matching succeeds, and each success appends exactly one journal
      line and returns the incremented aggregate revision; stale fails with
      `state-conflict`, appends zero journal lines, leaves the projected Issue
      unchanged, and reports exact integer `expected_revision` and
      `actual_revision` members. (VAL-018, VAL-044)
- [ ] One canonical `--json` transition table asserts every VAL-019 outcome with
      before/after journal counts: (a) a legal edge with satisfied preconditions
      succeeds and appends exactly one Event; (b) an unknown target state, (c) an
      absent edge between two known states, (d) a stored unsupported or unknown
      precondition, and (e) an unsatisfied precondition each reject with
      `transition-not-permitted`, exact `from_state`, `to_state`, and
      `unsatisfied` members, zero append, and unchanged status. (VAL-019)
- [ ] The successful canonical transition compares its envelope `event` identity
      against the journal tail: the returned Event ID and timestamp equal the
      last appended Event, and the returned revision equals prior plus one.
      (VAL-017, VAL-010, VAL-011)
- [ ] Both `--owner-role` and `--owner-actor` are proven unknown on **both**
      `issue create` and `issue update` (four cells), each with empty stdout, the
      exact canonical failure envelope on stderr, non-zero exit, and zero
      journal append. (VAL-040, VAL-044)
- [ ] A successful canonical `issue update` that changes every mutable field
      returns `changes` containing no owner member, and the projected
      `owner_role`/`owner_actor` are unchanged; the retained internal
      default-owner path still passes its existing regression test, named
      explicitly in the completion report. (VAL-040)
- [ ] Attribution is proven for `issue create`, `issue update`, and
      `issue transition`: each records the supplied local actor on its appended
      Event and projected activity, the appended Event's effects reference only
      the tables that operation is permitted to write, and every other projection
      table's row count is identical before and after the mutation — so no
      trusted Evidence, Human Outcome, or protected-authority entity or effect
      is created. (VAL-041)
- [ ] Every new rejection case asserts empty stdout with the canonical failure
      envelope on stderr; every new success case asserts stdout-only output.
      (VAL-002, VAL-008)

## First Failing Test

`cmd/fid/issue_transition_test.go` -
`TestIssueTransitionExpectedRevisionMatrix`.

This is evidence recovery, so a new assertion may pass on its first run. That is
an expected outcome, not a skipped red step: record the observed first-run result
for each new test with a verbatim token. A failing first run means a real defect
and triggers the defect clause above.

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

List any conditional command required by this task and why:

- None. The QA plan adds no Task 9 command beyond these six.

## Notes For Developer

- Build the concurrency and transition tables from one declared expectation
  table so a cell cannot be dropped silently, and take journal counts inside the
  table body rather than once per file.
- A transition needing a live caller claim must set that up per case; construct
  fixtures so both the absent and matching success cells are reachable.
- Seed a stored unsupported or unknown precondition through the same event-append
  helper the existing rule-failure test uses; if no supported surface can produce
  one, stop and escalate rather than dropping the cell.
- Completion report: `task-9a-done.md`. Cite each VAL explicitly, map every
  required cell above to the test that asserts it, state which Review 3 Blocking
  finding each new test closes, and record every required command with its
  verbatim token and observed termination result.
