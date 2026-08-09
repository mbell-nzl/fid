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

**In flight on branch `phase-20`. Tasks 1-3 complete; Review 1 returned BLOCK.**
`ESC-20-001` is resolved by DEC-130 without weakening VAL-002 or VAL-043: run
recovery tasks 3a and 3b, then bounded re-review `qa/review-1b.md`, then Task 4
on a non-BLOCK verdict.

## Operating Mode

**Heavy** (DEC-107/DEC-118/DEC-127): schema/persistence, public-interface, and
dependency changes. The binding delivery process is the V1 workflow in
`project-docs/specs/09_AI_WORKFLOW.md`. Required: Phase Runner, QA plan and
feasibility validation, six mid-phase reviews, final review, Human checklist,
full spec sync, run ledger, and triage log.

## Accepted Slice Contract

`../phases/phase-20/contracts.md` defines the accepted CLI behavior, Response
Envelope, leaf schemas, reason codes, DDL/Event semantics, removal inventory,
and legacy expiry ledger.

## Task Index

| Task | Title | Status | Dependencies |
|---|---|---|---|
| 1 | Response contract foundation | Complete; VAL-002 evidence owed by 3a | None |
| 2 | Project selection and config namespace | Complete; VAL-043 evidence owed by 3b | 1 |
| 3 | Atomic project initialization | Complete | 1-2 |
| 3a | Reason-code schema and redaction evidence | Pending (DEC-130) | 1 |
| 3b | Strict-config and init boundary evidence | Pending (DEC-130) | 2-3, 3a |
| 4 | Rebuild, authority, and lock failures | Blocked until review-1b is non-BLOCK | 1-3, 3a-3b |
| 5 | Remove orphaned Issue fields | Pending | 4 |
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
12, 15, and 17; Task 18 follows the sixth review. Tasks 3a/3b are DEC-130
recovery sessions that add no scope and change no accepted window;
`review-1b.md` is an escalation-resolution artifact, not a cadence review.

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
| 039-043 | config identity/strictness, owner boundary, attribution, atomic init (Tasks 2-12, 3b) |
| 044-046 | exact requests, order/bounds, opaque typed link endpoints (Tasks 2-18) |

## QA Artifacts

| Artifact | Path | Status |
|---|---|---|
| Validation contract | `../phases/phase-20/qa/validation-contract.md` | Human-approved; VAL-002/VAL-043 owning windows extended to tasks 3a/3b (DEC-130) |
| QA plan | `../phases/phase-20/qa/qa-plan.md` | Accepted; unchanged by DEC-130 |
| Mid-phase reviews | `../phases/phase-20/qa/review-M.md` | After Tasks 3, 6, 9, 12, 15, 17; review 6 includes independent contract/spec consistency |
| Recovery re-review | `../phases/phase-20/qa/review-1b.md` | Required after tasks 3a-3b, before Task 4; bounded to the recovery commits, the two sustained Blocking findings, and the Medium improvement |
| Final review / Human checklist | `../phases/phase-20/qa/` | Pending |

## Active Constraints

- Implement the approved form of `contracts.md` literally; divergence escalates.
- Preserve fail-closed V1 authority until Phase 21 replaces it; Phase 20 maps
  errors and weakens no check.
- Pull no Phase 21 capability forward: no capabilities command, Work Envelope,
  generic authority, Claim/Run reshape, gates/Receipts/Evidence, idempotency, or
  backend-port conformance.
- The caller never supplies resulting Issue revision. Mutation Events carry the
  absolute result plus prior/precondition values for deterministic replay.
- Issue revision is aggregate-wide: AC, Issue-label, and link mutations advance
  every affected Issue as well as create/update/transition.
- Root `--project`/`--json` never silently apply to legacy leaves; explicit use
  rejects before their handler. Workflow project scope is `--project-wide`, and
  global `--project` always selects a fid project root (DEC-119).
- Legacy JSON carries tested meaning only, not byte compatibility or a forward
  contract. `journal inspect` stays legacy until Phase 25.
- Bare “Envelope” means Work Envelope; JSON transport is Response Envelope.
- `09_AI_WORKFLOW.md` remains the binding V1 process; product target concepts do
  not replace its roles or artifacts.
- No open questions. A newly discovered contract/data-model ambiguity is a hard
  stop that returns to the Human, never a Developer inference.

## Parity Matrix

Retained JSON leaves get the closed envelope/request/order contract; Human mode
keeps concise stdout/stderr; legacy leaves keep their unenveloped mode, promise
no forward contract, and reject explicit root flags. Evidence: VAL-001..005,
VAL-008, VAL-029, VAL-044, VAL-045.

## Interfaces / Subsystems In Scope

- `cmd/fid/`, `internal/cli/`, `internal/project/` - root dispatch and transport;
  `internal/config/`, `internal/bootstrap/` - strict config and atomic init.
- Issue/label/link/workflow/maintenance packages - standalone journey behavior.
- `internal/journal/`, `internal/db/`, `migrations/` - revision and rebuild
  equivalence; `internal/format/` and command tests - formatter removal.

## Known Risks

- **Realized in Review 1:** wide tasks produce representative rather than
  exhaustive evidence. Where a VAL names a domain or matrix, every
  just-inside/just-outside pair and every registry member must be asserted;
  sampling is a Blocking finding (DEC-130).
- Schema/fixture churn is isolated in Tasks 5-6; formatter removal is high-churn
  and split across Tasks 14-16, reviewed after Tasks 6, 15, and 17.
- Exact command inventories in VAL-035 include `fid issue ac verify`; family ellipses are prohibited.

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
