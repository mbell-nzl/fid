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

**In flight on branch `phase-20`. Tasks 1-3l, recovery Review 1c, and the
independent infrastructure review are complete** (infra Blocking finding
corrected, NON-BLOCK at `cd9e6bc`; record anchored in
`../phases/phase-20/infra-improvement-pack.md`). **Task 4's implementation is
committed at `e56b3a0`**; DEC-135 resolved ESC-20-005 (dead lock-holder
formatter) and DEC-136 resolves ESC-20-006 by re-capturing the report's gate
tokens and finalizing it against that commit. No code change remains in Task 4.

## Operating Mode And Slice Contract

**Heavy** (DEC-107/DEC-118/DEC-127): schema/persistence, public-interface, and
dependency changes. The binding process is `project-docs/specs/09_AI_WORKFLOW.md`.
Required: Phase Runner, QA plan, six mid-phase reviews, final review, Human
checklist, full spec sync, run ledger, and triage log.
`../phases/phase-20/contracts.md` defines the accepted CLI behavior, Response
Envelope, leaf schemas, reason codes, DDL/Event semantics, removal inventory,
and legacy expiry ledger.

## Task Index

| Task | Title | Status | Dependencies |
|---|---|---|---|
| 1 | Response contract foundation | Complete; VAL-002 evidence owed by 3a | None |
| 2 | Project selection and config namespace | Complete; VAL-043 evidence owed by 3b | 1 |
| 3 | Atomic project initialization | Complete | 1-2 |
| 3a | Reason-code schema and redaction evidence | Complete (DEC-130) | 1 |
| 3b | Strict-config and init boundary evidence | Complete; VAL-043 gap owed by 3c | 2-3, 3a |
| 3c | Integration-journey convergence and honest gate baseline | Complete (DEC-131) | 3b |
| 3d | Evidence-token control and gate-record correction | Complete (DEC-131/132) | 3c |
| 3e-3l | Infrastructure performance pack | Complete (DEC-133/134) | review-1c |
| 4 | Rebuild, authority, and lock failures | Code committed `e56b3a0`; report repair in progress (DEC-136) | 1-3, 3a-3l |
| 5 | Remove orphaned Issue fields | Implementation done; classifier fix owed (DEC-137) | 4 |
| 6 | Replayable Issue aggregate revision | Pending | 5 |
| 7 | Issue hierarchy invariant | Pending | 5-6 |
| 8 | Issue read contract | Pending | 5-7 |
| 9 | Issue mutation and concurrency contract | Pending | 6-8 |
| 10 | Acceptance-criterion and label contracts | Pending | 8-9 |
| 11 | Link contract and dependency invariants | Pending | 8-10 |
| 12 | Workflow command contract | Pending | 8-11 |
| 13 | Maintenance inspection contract | Pending | 8-12 |
| 14 | Legacy JSON renderer conversion | Pending | 13 |
| 15 | Remaining format caller/test convergence | Pending | 14 |
| 16 | Delete universal formatters and YAML | Pending | 15 |
| 17 | Delete promised-command stubs | Pending | 16 |
| 18 | Journey acceptance and documentation | Pending | 1-17 |

One Developer session owns one task. Review windows close after Tasks 3, 6, 9,
12, 15, and 17; Task 18 follows the sixth review. Tasks 3a-3d are DEC-130/131
recovery sessions. Tasks 3e-3l are the DEC-133 infrastructure pack and add no
product VAL; `review-1b/1c.md` and `infra-review-1.md` are bounded artifacts,
not cadence reviews.

## Validation Contract

Full contract: `../phases/phase-20/qa/validation-contract.md`.

| VAL IDs | Owning work |
|---|---|
| 001-008 | transport, project selection, authority, locks, channels (Tasks 1-4; closes Task 18) |
| 009-012 | field removal, aggregate revision, replay/rebuild (Tasks 5-6, 9-11) |
| 013-020 | Issue reads, mutations, transitions, hierarchy (Tasks 7-9) |
| 021-026 | acceptance criteria, labels, links, dependencies, workflows (Tasks 10-12) |
| 027-034 | maintenance contract, legacy journal boundary, conversion, deletion (Tasks 13-17) |
| 035-038 | inventory, journey, documentation (Task 18); 038 is Human judgment |
| 039-043 | config identity/strictness, owner boundary, attribution, atomic init (Tasks 2-12, 3b-3c) |
| 044-046 | exact requests, order/bounds, opaque typed link endpoints (Tasks 2-18) |

## QA Artifacts

| Artifact | Path | Status |
|---|---|---|
| Validation contract | `../phases/phase-20/qa/validation-contract.md` | Human-approved; VAL-002/VAL-043 windows extended to tasks 3a-3c (DEC-130/131) |
| QA plan | `../phases/phase-20/qa/qa-plan.md` | Accepted; unchanged by DEC-130/131 |
| Mid-phase reviews | `../phases/phase-20/qa/review-M.md` | After Tasks 3, 6, 9, 12, 15, 17; review 6 includes independent contract/spec consistency |
| Recovery reviews | `../phases/phase-20/qa/review-1b.md`, `review-1c.md`, `infra-review-1.md` | 1b BLOCK corrected; 1c NON-BLOCK; infra BLOCK corrected, re-reviewed NON-BLOCK |
| Final review / Human checklist | `../phases/phase-20/qa/` | Pending; must re-run the full gate with real tokens |

## Active Constraints

- Implement the approved form of `contracts.md` literally; divergence escalates.
- Preserve fail-closed V1 authority until Phase 21 replaces it; Phase 20 maps
  errors and weakens no check.
- Pull no Phase 21 capability forward: no capabilities command, Work Envelope,
  generic authority, Claim/Run reshape, gates/Receipts, idempotency, or
  backend-port conformance.
- The caller never supplies resulting Issue revision; Mutation Events carry the
  absolute result plus prior/precondition values for deterministic replay.
  Revision is aggregate-wide: AC, Issue-label, and link mutations advance every
  affected Issue as well as create/update/transition.
- Root `--project`/`--json` never silently apply to legacy leaves; explicit use
  rejects before their handler. Workflow project scope is `--project-wide`, and
  global `--project` always selects a fid project root (DEC-119).
- Legacy JSON carries tested meaning only, not byte compatibility or a forward
  contract. `journal inspect` stays legacy until Phase 25.
- Bare “Envelope” means Work Envelope; JSON transport is Response Envelope.
  `09_AI_WORKFLOW.md` stays the binding V1 process; product target concepts do
  not replace its roles or artifacts.
- No open questions; a new contract/data-model ambiguity is a hard stop for the Human.

## Parity Matrix

Retained JSON leaves get the closed envelope/request/order contract; Human mode
keeps concise stdout/stderr; legacy leaves keep their unenveloped mode and
reject root flags. Evidence: VAL-001..005, VAL-008, VAL-029, VAL-044, VAL-045.

## Interfaces / Subsystems In Scope

- `cmd/fid/`, `internal/cli/`, `internal/project/` - dispatch and transport;
  `internal/config/`, `internal/bootstrap/` - strict config and atomic init.
- Issue/label/link/workflow/maintenance packages - journey behavior;
  `internal/journal/`, `internal/db/`, `migrations/` - revision and rebuild;
  `internal/format/` and command tests - formatter removal.

## Known Risks

- **Realized in Review 1:** wide tasks produce representative rather than
  exhaustive evidence. Where a VAL names a domain or matrix, every
  just-inside/outside pair and registry member must be asserted (DEC-130).
- **Realized in Review 1b:** record each gate with a verbatim token from that
  command; `exit 0 (no output)` is never valid for `go test`, and an uncaptured
  gate is `UNVERIFIED`, never `PASS` (DEC-131/132).
- Schema/fixture churn is isolated in Tasks 5-6; formatter removal is high-churn
  across Tasks 14-16. Exact VAL-035 inventories include `fid issue ac verify`;
  family ellipses are prohibited.

## Success / Exit Criteria

- [ ] All eighteen tasks plus recovery tasks 3a/3b complete with exact-command evidence.
- [ ] Every accepted VAL has a phase-end outcome with cited evidence.
- [ ] Complete standalone Issue journey passes and survives rebuild unchanged.
- [ ] Universal formatters/flags/dependency, promised stubs, and retired Issue fields are absent.
- [ ] Fail-closed V1 authority still rejects; no Phase 21 capability is advertised.
- [ ] Final QA and Human checklist have no unresolved blocking outcome.

## Next Planning Checkpoint

After Phase 20 wrap-up, resolve Gate D3.2: the exact one-Issue Routine walking
skeleton boundary for Phase 21.
