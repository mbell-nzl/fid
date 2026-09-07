# Phase 22 Remaining-Scope Allocation Audit

Date: 2026-09-07. Author: delegated Tech Lead preflight reviewer.
Scope: accepted tasks 10-18; earlier shared fixtures inspected only where the
remaining tasks depend on them. Read-only source baseline observed at
`7d29403c5e8868085a7c43c89ffe821492e13181`. Task 8 continued concurrently.

This is an allocation audit, not final QA or a validation outcome. No product
code, tests, task specs, ledger, or current-context documents were edited. No
quality gate was run, no prior agent receipt is adopted as this session's
evidence, and no held-out material was read. Line references below describe
the inspected baseline; subsequent task amendments may move them.

## Result And Requested Disposition

No additional closed-domain/DDL contradiction of the DEC-195 kind was found
in the inspected remaining contract rows. Three concrete delivery concerns
needed explicit execution instructions before their owning tasks close:

| ID | Finding | Classification | Narrow proposed owner |
|---|---|---|---|
| RSA-01 | Shared recovery fake returns fixture-shaped strings rather than semantic state; its assertions accept invalid result types. | Accepted conformance work incompletely delivered; no new product decision. | Task 9's already-owned shared runner/fake, with task-6 provenance; extend in 11/15, qualify in 17, independently audit in 18/final QA. |
| RSA-02 | Human recovery rendering is first explicitly checked by task 18, whose scope excludes production implementation. | Acceptance-criterion allocation gap; task 12 already owns the renderer. | Bind Human resume rendering explicitly to task 12; task 18 retains independent parity evidence. |
| RSA-03 | Runnable Human fixtures have no named retained output/generation receipt; present process fixtures disappear with test cleanup. | Delivery/evidence gap in the task-18 preparation contract. | Task 18 test-only generator, retained fixture manifest/instructions, and an exact generation/smoke receipt. |

The delegating root accepted these as bounded implementation/evidence
allocation work and recorded DEC-196 in
`2e99191f16d3efe1c5200b52239d5ebf3e51ae67`. This reviewer inspected that commit:
tasks 9, 12 and 18 now explicitly allocate the repairs described below. The
source allocation gaps are dispositioned; implementation and independent
runtime verification remain outstanding. This audit does not itself authorize
edits or change any validator, assertion, threshold, or accepted behavior.

## RSA-01: The Existing Shared Fake Does Not Model Semantics

Accepted sources: `contracts.md` section 11 requires shared versioned request,
result, error, and negative cases. Task 9's Binding Constraints already require
identical semantic expectations against fake and fid. Tasks 11 and 15 extend
that runner; task 17 conditions availability on those cases; task 18/VAL-073
requires the same inputs and expected semantic outputs.

Observed source:

- `cmd/fid/integration/p22_fake_test.go:19` defines cases as operation/kind,
  result field names, and expected code/exit; there is no shared request object.
- At lines 141-154, the fake returns `field + "-fake"` for each expected field
  and copies the fixture's failure/exit values. It has no stored domain state.
- At lines 110-123, the shared success assertion checks field presence and
  whether Event is an object/null, without checking the fields' types/values.
- `testdata/p22_conformance/operations/run.recover.accepted.json` names
  `recovery`, `run`, `claim`, `idempotency_key`, and `replayed`. Therefore the
  fake supplies `replayed:"replayed-fake"` and `run:"run-fake"`; these pass
  those source-level predicates despite violating the accepted boolean/Run
  types. This is a static counterexample, not a test execution claim.
- The fid-only `p22Advance` check at line 306 checks the prior Run ID, but does
  not turn that into a shared fake/fid semantic expectation.

`git log -1 -S 'result[field] = field + "-fake"'` identifies
`f3be616f0587881ee3f5c3716ffbf6167e3f2a20`, whose body records
`Refs: phase-22/task-6`. `qa/review-2b.md` records historical completion based
on field/channel/Event-presence checks. It records no exception allowing
shape-only evidence to satisfy future semantic qualification. Its historical
verdict need not be rewritten to allocate the causal strengthening now.

Proposed task-9 instruction: explicitly name `p22_fake_test.go`, matching
test-only integration helpers, and `testdata/p22_conformance/` in scope. Replace
the response fabrication with an independent in-memory model driven by shared
versioned requests; the fake must not open fid SQLite, journal, or `.fid/`
files, call production domain operations, or read expected outcomes to select
its response. Use the same independently specified typed values and identity
relationships for both adapters. Retain existing field/channel checks.

Concrete acceptance evidence: recovery creates a distinct linked successor;
first response has boolean `replayed:false`; same-key replay returns the same
identities with `replayed:true` and no additional effects. Invalid requests
leave state unchanged. Task-9 Stop/Attention and budget cases compare their
actual subject/disposition links, deduplicated identity, counters, and committed
denial effects. Demonstrate that a wrong result type, changed predecessor,
wrong counter, or a replay that increments state is rejected by the common
assertions. Existing ordinary/tagged gates execute this scope; no extra tag,
environment switch, production change, or lower acceptance bar is needed.

## RSA-02: Assign Human Resume Rendering Before Final Assembly

Accepted `public-interface.md` section 6 requires recovery text to name
prior/successor IDs, adopted paths, authority and remaining budgets, plus
known-live takeover/residual filesystem risk where applicable. VAL-053 requires
IDs, status, reasons, budgets and uncertainty to agree across renderings.

Task 12 owns `cmd/fid/resume.go` and the recovery projection, but its criteria
only name VAL-048..052. Task 18 first explicitly names Human/JSON parity while
its scope consists of integration tests/fixtures and its Binding Constraints
forbid new domain behavior. At the inspected baseline,
`cmd/fid/resume.go:32` renders only Issue ID and completion status. This is
expected unfinished task-12 work, not a newly demonstrated product regression;
the gap is the absence of a directly failing Human-render criterion at its
production owner.

Proposed additional task-12 acceptance criterion:

> From the same bounded Resume projection, the built binary's default Human
> view and `--json` view expose the current Issue/Run identity, predecessor and
> successor links, effective lease and executor uncertainty, adopted paths and
> authority, unresolved Stops/Attention and their qualifying resolution,
> outstanding Obligations/evidence availability, remaining budgets, and next
> actions with blocking reasons. The Human view states known-live takeover and
> its residual filesystem risk when applicable and does not describe adoption
> as proof. Paired deterministic cases cover an ordinary interrupted Run, an
> adopted dirty target, and blocked/force states. (VAL-051/052; prepares the
> already-accepted VAL-053 parity and VAL-076 Human fixture.)

Keep JSON schema and final validator ownership unchanged. Existing task-12
renderer, domain projection, matching command tests and tagged recovery journey
scope suffice. At task 11 likewise exercise activation Human status/recovery
output against section 6, so its pre/post-commit fixture is useful; task 11
already owns `cmd/fid/activation.go` and VAL-045/046. Task 15 already owns
import rendering. Task 18 remains the final completeness/parity audit.

## RSA-03: Make Tomorrow's Fixture Delivery Concrete

Task 18's objective and Binding Constraints assign reproducible Human fixtures
and reserve judgments to Human. QA plan EVID-67..69 and section 9 require a
fresh resume view, paired pre/post-commit activation views, and an imported
project requiring intervention. Its unamended file/command sections name only
test files, corpus files, a completion report, ordinary/tagged gates and build.
They name no fixture generator, retained output contract, or invocation proof.

Existing `p22_recovery_journey_test.go:35` and `p22_fake_test.go:189` use
`t.TempDir()`. `integration/main_test.go:42` deletes the suite binary when
tests end. Their successful execution alone cannot deliver an inspectable
project/binary tomorrow. This does not forbid a test-only generator: it
identifies exactly what the task must add and evidence.

DEC-196's task-18 scope explicitly adds
`cmd/fid/integration/p22_human_fixtures_test.go` and
`project-docs/docs/phases/phase-22/qa/human-fixtures.md`; the existing
`cmd/fid/integration/testdata/p22_conformance/` scope can hold descriptions.
All setup should use public built-binary operations.
Normal tagged tests execute every fixture setup and post-setup read path;
there is no opt-in-only behavioral coverage. One optional destination-parent
variable changes only where delivery output is retained. It creates a fresh
unique child directory and never overwrites an earlier fixture.

Accepted DEC-196 required generation/smoke command, from repository root:

```text
FID_P22_HUMAN_FIXTURE_ROOT=/home/mark/projects/fid/.cache/phase-22-human go test -tags integration ./cmd/fid/integration/... -run '^TestP22HumanFixtures$' -count=1 -v
```

This is a command allocation, not an executed receipt. It supplements the
unchanged full tagged race gate, which must exercise every fixture setup mode.
Its retained mode must preserve a usable built binary beyond TestMain cleanup
(a copied binary or an explicitly durable path) and leave independent projects for
recovery, pre-commit activation, post-commit activation, and import. Validate
readback through separate processes after their setup processes have exited.

Persist a machine-readable manifest with absolute fixture roots, subject IDs,
actual source revision/dirty state, binary digest, and exact Human/JSON read
commands. Emit the final retained root in test output and copy it into the
completion report. Keep binary/projects in the gitignored `.cache/` child;
commit only generator/source instructions and the report. Regenerate from the
final committed source before Human handoff if subsequent code changed the
binary, preserving the earlier receipt rather than relabelling its revision.

The Human-facing instructions must contain neutral prompts and space to record
answers, without supplying the expected next action/direction in advance.
Keep QA's expected outcomes separately inspectable. Setup/readback can be
verified deterministically; VAL-076..078 remain unjudged until Human responds.
The root should also document that QA plan section 10's Human judgment
condition gates final sign-off, while the delegated preparatory final QA review
may explicitly retain those outcomes pending for tomorrow.

## Remaining Task Allocation And Gate Reach

| Task | Contract allocation inspected | Gate reach / hand-back requirement |
|---|---|---|
| 10 | Activation prepare/resources/version/direction/trust and four owned Event families; `activation.go`, recovery/authority, registration, journal, existing CLI parser. | Ordinary plus tagged gates. Remove the three last fence keys and empty fence; preserve VAL-003 unknown-Event parity. Existing 0013 predecessor/status and 0015 resource/recovery rows admit the declared closed domains. |
| 11 | Activation new leaves/errors and consuming-path config validation; domain materialization stays outside fid. | Both gates. Existing `ReadProjectConfig` already parses the entire closed config; no missing alternate parser was found. Paired malformed-config consuming/independent cases and activation crash/direction cases remain required. |
| 12 | Current repository freshness, Evidence availability, bounded Resume and eligible actions. | Both gates; real canonical dirty observation required. DEC-196 assigns RSA-02's explicit rendering criterion. Independent QA checkpoint still precedes task 13. |
| 13 | Manifest v2/config, full Event semantics, bounded non-mutating archive validation. | Both gates. Every 056-a..k case and below/equal/above bound is already allocated; measured 8 GiB streaming evidence cannot be replaced by ZIP-metadata rejection. Export split explicitly assigned. |
| 14 | Empty-target lock/stage/fsync/publication, preserved identities, response-loss retry and imported authority fence. | Both gates. `internal/export/project/config/journal/db`, authority, Claim and execution scope are named. No new domain import Event is authorized or needed by the declared publication contract. |
| 15 | Import CLI, active/interrupted/abandoned/budget/activation/Artifact corpus and three-way named JSON/ordered Event parity. | Both gates. Same semantic fake/fid expectations must precede capability enablement; preserve intermediate-schema convergence when 0016 lands. |
| 16 | Exact 0016 removal, retained historical migrations, activity inspection, latest-schema/replay convergence. | Both gates plus one fresh reference benchmark. Inspected production escalation matches are in maintenance/rules paths already owned. `triage_events`/Concern lifecycles remain later-owned. |
| 17 | Maintenance corruption detection, exact four changed capability rows and final schema/reference performance. | Both gates; identity capture; three fresh benchmark processes; two full benchmark gates; build. Capability enablement depends on semantic evidence, not RSA-01's field-only runner. |
| 18 | Corpus completeness, final public journey, Human/JSON parity, Human fixture delivery and independent kit consistency. | Both gates plus build; DEC-196 adds RSA-03's concrete generation/smoke receipt. No new production behavior assigned here. Human judgments remain Human-owned. |

Every task already includes the exact tagged integration gate
`go test -tags integration -race ./cmd/fid/integration/...` alongside the
ordinary race/coverage gate. `integration/main_test.go` builds the real CLI.
No excluded suite was discovered among the inspected assigned surfaces. This
is source-level reach analysis, not proof that any pending tests will pass.

## Carry-Forward And Source-Currency Checks For The Root

- CF-024 remains task 17/final QA: exact reference identity and three fresh
  per-sample <30 s results. Task 16's single early sample is supplementary.
  Avoid competing discretionary CPU-heavy work during the reference runs;
  never change the instance/threshold or select successful samples afterward.
- CF-034 belongs to tasks 2/7; final QA audits direct branch tests and actual
  before/after coverage receipts without substituting a new numeric threshold.
- CF-035 fires in task 11. CF-018 fires only if production
  `internal/sqlquery/sqlquery.go` changes; no such planned change was found.
  Test-only convergence does not trigger its split.
- CF-033/036 and deferred report/control tooling remain separate maintenance;
  this audit grants no workflow-control edit. CF-038 fires only if an owning
  causal task opens the live-output fixture; do not silently broaden its wire.
- CF-039 and DEC-195's lesson are directly addressed by RSA-02/03: a named
  scope or prose promise needs a failing criterion and an executable receipt.
- CF-037 and all C-01..C-13 accepted kit dispositions remain external
  qualification obligations where explicitly deferred; do not implement a
  production adapter or claim its qualification in Phase 22.
- Root reported an independent external-kit inspection during this audit:
  clean current HEAD `8fdd7d22c80ec508646ce53d0fae112f5b9216a2`, ten commits
  beyond accepted `c5229617220f97ff9228df261e24718b1d1d6d89`; the scoped
  interface change adds offline ordinary-frame binding in section 7.1,
  recovery decision 0005 is unchanged, and compatibility status/receiving
  obligations changed. This paragraph is attributed coordination context,
  not this reviewer's verification or a gate result. Final QA must verify
  current-source currency itself while preserving the accepted semantic pin
  and DEC-187 dispositions.

## Evidence Limits

Evidence here is task/contract/DDL/source inspection plus the observed Git
provenance query for RSA-01. No absence inventory is labelled a gate result.
The initial optional current-context handoff read reported that
`project-docs/docs/current/tl_handoff.md` was absent; current phase context was
available from the other required documents. No conclusion depends on that
optional file. Pending behavior remains for its Developer and independent QA
sessions to implement and verify.
