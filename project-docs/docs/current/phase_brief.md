# Phase Brief - Phase 20

> Default execution brief loaded by every agent session.

## Phase Goal

Deliver one coherent standalone Issue journey on canonical JSON: initialize a
project; create, relate, inspect, update, and transition Issues; explain failure
with typed reasons; rebuild the projection; repeat unchanged. Establish the
`fid-cli.v1` Response Envelope, project resolution/identity, replayable Issue
revision, and the retained Issue contract; delete universal formatters,
`--format`, promised stubs, and Issue fields with no retained meaning.

## Status

**In flight on branch `phase-20`, awaiting Human verification.** Tasks 1-18,
the infrastructure pack, and every review window through `review-6` are
complete (`../phases/phase-20/infra-improvement-pack.md`). The fresh
`../phases/phase-20/qa/final-review-3.md` is **PASS**: MF-01 report conformance
and MF-02 integration placement are resolved, the full phase gate passes, and
45 deterministic validation assertions pass. Only Human-owned VAL-038 remains;
no escalation or protected decision is open. **Next: complete the three blank
rows in `human-checklist.md`, then wrap the phase.**

## Operating Mode And Slice Contract

**Heavy** (DEC-107/118/127): schema/persistence, public-interface, and dependency
changes. Binding process: `project-docs/specs/09_AI_WORKFLOW.md`. Required
artifacts: Phase Runner, QA plan, six mid-phase reviews, final review, Human
checklist, full spec sync, run ledger, triage log.
`../phases/phase-20/contracts.md` defines the accepted CLI behavior, Response
Envelope, leaf schemas, reason codes, DDL/Event semantics, and removals.

## Task Index

| Task | Title | Status | Dependencies |
|---|---|---|---|
| 1-3 | Response contract, project selection, atomic init | Complete | None |
| 3a-3d | Reason-code, strict-config, journey, and gate-record evidence recovery | Complete (DEC-130/131/132) | 1-3 |
| 3e-3l | Infrastructure performance pack | Complete (DEC-133/134) | review-1c |
| 4-8 | Rebuild/authority/locks, field removal, aggregate revision, hierarchy, read contract | Complete `e56b3a0`, `028855b`, `29c5246`, `4950e94`, `4df4d66` (DEC-135..139) | 1-3, 3a-3l |
| 9, 9a | Issue mutation and concurrency contract, then its evidence recovery | Complete `86d927d`, `1f196a9`; review-3b NON-BLOCK (DEC-140) | 6-8 |
| 10-13 | AC/labels, links, workflow, maintenance | Complete `627ca9d`, `49c179a`, `161f89c` (DEC-141); reviews 4b/4c/5b closed | 8-9a |
| 14 | Legacy JSON conversion **and its caller convergence** | Complete (DEC-142) | 13 |
| 15 | Format shim removal and whole-tree convergence | Complete `24f50aef`; review-5d closed MF-02 | 14 |
| 16 | Delete universal formatters and YAML | Complete `20813fc` / `d0f9cea` (DEC-143, DEC-144, DEC-145) | 15, 16a, 16b |
| 16a, 16b | report-lint accepts declared expected-non-zero rows, then parses declaration lines not prose | Complete `d9d6add`, `d8388d5`; Task 16 recovery-support slices (DEC-143/144) | 15 |
| 17 | Delete promised-command stubs | Complete `95ee13d`; control recovery report `0c1094c` (DEC-146) | 16 |
| 18 | Journey acceptance and documentation | Complete; MF-02 resolved by real-process integration journey `90841bb` | 1-17 |

One Developer session owns one task. Review windows close after Tasks 3, 6, 9,
12, 15, and 17; Task 18 followed the closed sixth review. Tasks
3a-3d, 9a, 16a, 16b are recovery sessions and 3e-3l is the DEC-133 infrastructure
pack; none adds a product VAL. `16a`/`16b` ride the ledger as `REPAIR task-16`;
Task 16 and Task 17 now have their single renewed completion events.

## Validation Contract

Full contract: `../phases/phase-20/qa/validation-contract.md`.

- **001-012** transport, selection, authority, locks, channels, field removal, revision, replay (1-11)
- **013-026** Issue reads, mutations, transitions, hierarchy, AC, labels, links, workflows (7-12)
- **027-038** maintenance, legacy boundary, conversion, deletion, inventory, journey, docs (13-18); 038 is Human
- **039-046** config identity/strictness, owner boundary, attribution, init, exact requests, order/bounds, typed endpoints (2-18)

## QA Artifacts

All under `../phases/phase-20/qa/`. `validation-contract.md` is Human-approved
(VAL-002/VAL-043 windows extended to 3a-3c; VAL-032 evidence method revised by
DEC-143/144/145), and `qa-plan.md` is accepted with its unsatisfiable YAML
`go.sum` row replaced by a stronger pair. Bounded re-reviews `review-1b/1c`,
`infra-review-1`, `review-3b`, `review-4b/4c`, `review-5b/5c/5d` closed every
prior finding, and `review-6` is closed. `final-review.md` is REPAIR on MF-01 and
MF-02; `final-review-2.md` closes Task 1's bounded MF-01 share only.
`final-review-3.md` closes MF-01 and MF-02 with a deterministic PASS;
`human-checklist.md` now awaits the three Human VAL-038 judgments.

## Active Constraints

- **A task acceptance criterion summarizes a validation contract row; it never
  replaces it (DEC-140).** Where a VAL's Evidence or Pass column names a matrix,
  table, or "every X", assert one observable case per cell for every named
  operation and cite the cell-to-test mapping in the report; representative or
  single-operation evidence is a failed criterion.
- **An absence command records its real non-zero status (DEC-143).** A task spec
  declares the expected status under `## Expected Non-Zero Commands` (only
  backtick-bullet lines there are declarations); the report writes
  ``PASS - `<token>` (expected `exit_code: N`; terminated `exit_code: M`)`` with
  `M` equal to `N`, never `exit 0 (no output)`. VAL-032's YAML absence is a
  `go.mod` requirement and build-graph proof: a retained `go.sum` checksum is
  module-graph bookkeeping, the migration dependency does not change, and
  `go mod why -vendor` passes on exit 0 with exactly one stanza -
  `# gopkg.in/yaml.v3` then `(main module does not need to vendor module ...)`.
  An import-path line, any further non-empty line, or a second stanza header is
  a real failure (DEC-144/145).
- **A removal task owns its own invalidated callers (DEC-142).** Removing a flag,
  field, handler, or symbol owns every caller it breaks, test callers and
  orphaned helpers included; never defer work an earlier task's gate needs.
- Implement the approved form of `contracts.md` literally; divergence escalates.
- Preserve fail-closed V1 authority until Phase 21 replaces it; Phase 20 maps
  errors and weakens no check.
- Pull no Phase 21 capability forward: no capabilities command, Work Envelope,
  generic authority, Claim/Run reshape, gates/Receipts, idempotency, backend-port
  conformance.
- The caller never supplies resulting Issue revision; Mutation Events carry the
  absolute result plus prior/precondition values for deterministic replay, and
  revision is aggregate-wide across AC, Issue-label, and link mutations.
- Root `--project`/`--json` never silently apply to legacy leaves; explicit use
  rejects before their handler (DEC-119). Legacy JSON carries tested meaning only;
  `journal inspect` stays legacy until Phase 25.
- Bare “Envelope” means Work Envelope; JSON transport is Response Envelope. No
  open questions; a new contract/data-model ambiguity is a hard stop for the Human.
- **Parity matrix:** retained JSON leaves get the closed envelope/request/order
  contract; Human mode keeps concise stdout/stderr; legacy leaves keep their
  unenveloped mode and reject root flags (VAL-001..005/008/029/044/045).

## Interfaces / Subsystems In Scope

`cmd/fid/`, `internal/cli/`, `internal/project/` (dispatch/transport);
`internal/config/`, `internal/bootstrap/` (strict config, atomic init);
Issue/label/link/workflow/maintenance packages (journey behavior);
`internal/journal/`, `internal/db/`, `migrations/` (revision, rebuild);
`internal/format/` deleted, its guards retained in `cmd/fid`.

## Known Risks

- **Realized in Reviews 1 and 3:** wide tasks produce representative evidence.
  Where a VAL names a domain or matrix, assert every just-inside/outside pair and
  registry member (DEC-130) for every operation the row names (DEC-140); a green
  broad command evidences no unasserted state.
- **Realized in Review 1b:** record each gate with a verbatim token from that
  command; `exit 0 (no output)` is never valid for `go test`, and an uncaptured
  gate is `UNVERIFIED` (DEC-131/132). Exact VAL-035 inventories include
  `fid issue ac verify`; no family ellipses (ESC-20-011, DEC-142).
- **Realized in ESC-20-012/013/014, all on one gate row:** state a required
  command's reachable pass condition and check no dependency or control prevents
  it (DEC-143); transcribe every exact-output *and complete-output* literal from
  an observed run of that command and flags, framing included (DEC-144/145);
  self-test a section-parsing control against every real repo shape (DEC-144).

## Success / Exit Criteria

- [x] All eighteen tasks plus recovery sessions 3a-3d, 9a, 16a, and 16b complete with exact-command evidence.
- [x] Every accepted VAL has a phase-end outcome with cited evidence.
- [x] Complete standalone Issue journey passes and survives rebuild unchanged.
- [x] Universal formatters/flags/dependency, promised stubs, and retired Issue fields are absent.
- [x] Fail-closed V1 authority still rejects; no Phase 21 capability is advertised.
- [ ] Final QA and Human checklist have no unresolved blocking outcome.

## Next Planning Checkpoint

After wrap-up, resolve Gate D3.2: the one-Issue Routine walking-skeleton
boundary for Phase 21.
