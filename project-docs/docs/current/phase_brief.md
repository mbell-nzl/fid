# Phase Brief - Phase 20

> Default execution brief loaded by every agent session.

## Phase Goal

Deliver one coherent standalone Issue journey on canonical JSON: initialize a
project; create, relate, inspect, update, and transition Issues; explain failure
with typed reasons; rebuild the projection; and repeat unchanged. Establish the
`fid-cli.v1` Response Envelope, project resolution/identity, replayable Issue
revision, and the retained Issue contract. Delete universal formatters,
`--format`, promised stubs, and Issue fields with no retained meaning.

## Status

**In flight on branch `phase-20`. Tasks 1-15 and the infrastructure pack are
complete; every cadence and bounded review through `review-5d` is closed** (infra
record: `../phases/phase-20/infra-improvement-pack.md`). Task 16's deletion work
is committed at `afb262d`; it has stopped twice on control and spec-literal
defects, resolved by DEC-143 and DEC-144. `task-16a` closed at `d9d6add`.
ESC-20-001..013 are resolved; nothing is open.
**Next: `task-16b`, the Task 16 continuation, Task 17, Review 6, Task 18.**

## Operating Mode And Slice Contract

**Heavy** (DEC-107/DEC-118/DEC-127): schema/persistence, public-interface, and
dependency changes. Binding process: `project-docs/specs/09_AI_WORKFLOW.md`.
Required artifacts: Phase Runner, QA plan, six mid-phase reviews, final review,
Human checklist, full spec sync, run ledger, triage log.
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
| 16 | Delete universal formatters and YAML | Deletion committed `afb262d`; continuation owed (DEC-143, DEC-144) | 15, 16a, 16b |
| 16a | report-lint accepts declared expected-non-zero commands | Complete `d9d6add` (DEC-143) | 15 |
| 16b | report-lint parses declaration lines, not prose | **Next** (DEC-144) | 16a |
| 17 | Delete promised-command stubs | Pending | 16 |
| 18 | Journey acceptance and documentation | Pending | 1-17 |

One Developer session owns one task. Review windows close after Tasks 3, 6, 9,
12, 15, and 17; Task 18 follows the sixth review. Tasks 3a-3d, 9a, 16a, and 16b
are DEC-130/131/140/143/144 recovery sessions and 3e-3l is the DEC-133
infrastructure pack; none adds a product VAL. Review 6 is the next cadence
review.

## Validation Contract

Full contract: `../phases/phase-20/qa/validation-contract.md`.

- **001-012** transport, selection, authority, locks, channels, field removal, revision, replay (1-11)
- **013-026** Issue reads, mutations, transitions, hierarchy, AC, labels, links, workflows (7-12)
- **027-038** maintenance, legacy boundary, conversion, deletion, inventory, journey, docs (13-18); 038 is Human
- **039-046** config identity/strictness, owner boundary, attribution, init, exact requests, order/bounds, typed endpoints (2-18)

## QA Artifacts

All under `../phases/phase-20/qa/`. `validation-contract.md` is Human-approved
(VAL-002/VAL-043 windows extended to 3a-3c; VAL-032 evidence method clarified by
DEC-143), and `qa-plan.md` is accepted with its unsatisfiable YAML `go.sum` row
replaced by a stronger pair (DEC-143). Bounded re-reviews `review-1b/1c`,
`infra-review-1`, `review-3b`, `review-4b/4c`, and `review-5b/5c/5d` closed every
prior Blocking and Must-fix finding. Review 6 is due after Task 17 and adds
independent contract/spec consistency. Final review and the Human checklist must
re-run the full gate with real tokens.

## Active Constraints

- **A task acceptance criterion summarizes a validation contract row; it never
  replaces it (DEC-140).** Where a VAL's Evidence or Pass column names a matrix,
  table, or "every X", assert one observable case per cell for every named
  operation and cite the cell-to-test mapping in the completion report.
  Representative or single-operation evidence is a failed criterion.
- **An absence command records its real non-zero status (DEC-143).** A task spec
  declares the expected status under `## Expected Non-Zero Commands`; the report
  writes ``PASS - `<token>` (expected `exit_code: N`; terminated `exit_code: M`)``
  with `M` equal to the declared `N`. Never rewrite a non-zero absence result as
  `exit 0 (no output)`. VAL-032's YAML absence is a `go.mod` requirement and
  build-graph proof; a `go.sum` checksum retained for an accepted dependency's
  module graph is bookkeeping, not a fid dependency, and the migration
  dependency does not change. That proof passes on the tool's own wording -
  `(main module does not need to vendor module ...)` under `-vendor` (DEC-144);
  an import chain is still a real failure. Only backtick-bullet lines in
  `## Expected Non-Zero Commands` are declarations; other lines are prose.
- **A removal task owns its own invalidated callers (DEC-142).** A task that
  removes a flag, field, handler, or symbol owns every caller that removal
  breaks, including test callers and helpers it orphans. Never defer to a later
  task work that an earlier task needs to pass its own required gate.
- Implement the approved form of `contracts.md` literally; divergence escalates.
- Preserve fail-closed V1 authority until Phase 21 replaces it; Phase 20 maps
  errors and weakens no check.
- Pull no Phase 21 capability forward: no capabilities command, Work Envelope,
  generic authority, Claim/Run reshape, gates/Receipts, idempotency, or
  backend-port conformance.
- The caller never supplies resulting Issue revision; Mutation Events carry the
  absolute result plus prior/precondition values for deterministic replay, and
  revision is aggregate-wide across AC, Issue-label, and link mutations.
- Root `--project`/`--json` never silently apply to legacy leaves; explicit use
  rejects before their handler (DEC-119). Legacy JSON carries tested meaning
  only, and `journal inspect` stays legacy until Phase 25.
- Bare “Envelope” means Work Envelope; JSON transport is Response Envelope.
  `09_AI_WORKFLOW.md` stays the binding V1 process. No open questions; a new
  contract/data-model ambiguity is a hard stop for the Human.
- **Parity matrix:** retained JSON leaves get the closed envelope/request/order
  contract; Human mode keeps concise stdout/stderr; legacy leaves keep their
  unenveloped mode and reject root flags. Evidence: VAL-001..005, VAL-008,
  VAL-029, VAL-044, VAL-045.

## Interfaces / Subsystems In Scope

`cmd/fid/`, `internal/cli/`, `internal/project/` (dispatch/transport);
`internal/config/`, `internal/bootstrap/` (strict config, atomic init);
Issue/label/link/workflow/maintenance packages (journey behavior);
`internal/journal/`, `internal/db/`, `migrations/` (revision, rebuild).
`internal/format/` is deleted; its guards and command tests remain in `cmd/fid`.

## Known Risks

- **Realized in Reviews 1 and 3:** wide tasks produce representative rather than
  exhaustive evidence. Where a VAL names a domain or matrix, assert every
  just-inside/outside pair and registry member (DEC-130) for every operation the
  row names (DEC-140). A green broad command never evidences an unasserted state.
- **Realized in Review 1b:** record each gate with a verbatim token from that
  command; `exit 0 (no output)` is never valid for `go test`, and an uncaptured
  gate is `UNVERIFIED`, never `PASS` (DEC-131/132).
- **Realized in ESC-20-011:** exact VAL-035 inventories include
  `fid issue ac verify`; no family ellipses (DEC-142).
- **Realized in ESC-20-012:** a planned command validated only as "matches
  today" may be unreachable in its finished state, and a control may reject
  evidence canon requires. Before adding an absence gate, state the reachable
  pass condition and check no accepted dependency or rule prevents it (DEC-143).
- **Realized in ESC-20-013:** transcribe an exact-output literal from an observed
  run of that exact command with those exact flags, never from documentation or a
  sibling flag; and self-test any section-parsing control against every real
  shape of that section in the repository (DEC-144).

## Success / Exit Criteria

- [ ] All eighteen tasks plus recovery tasks 3a-3d, 9a, 16a, and 16b complete with exact-command evidence.
- [ ] Every accepted VAL has a phase-end outcome with cited evidence.
- [ ] Complete standalone Issue journey passes and survives rebuild unchanged.
- [ ] Universal formatters/flags/dependency, promised stubs, and retired Issue fields are absent.
- [ ] Fail-closed V1 authority still rejects; no Phase 21 capability is advertised.
- [ ] Final QA and Human checklist have no unresolved blocking outcome.

## Next Planning Checkpoint

After Phase 20 wrap-up, resolve Gate D3.2: the exact one-Issue Routine walking
skeleton boundary for Phase 21.
