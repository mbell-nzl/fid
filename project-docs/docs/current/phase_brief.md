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

**In flight on branch `phase-20`. Tasks 1-9a and the infrastructure pack are
complete; Reviews 1c, infra-review-1, 2, and bounded Review 3b are NON-BLOCK**
(record anchored in `../phases/phase-20/infra-improvement-pack.md`). Review 3's
three sustained Task 9 evidence findings were closed by test-led Task 9a at
`1f196a9` and independently re-reviewed in `qa/review-3b.md`. Task 10 is in a
decision-free same-task recovery under DEC-141: its product work and ordinary
suite are complete, while the candidate commit still requires the exact gate
from a normal clone because Go VCS stamping misresolves this linked worktree.
ESC-20-001..010 are resolved; no product or architecture decision is open.
**Next: finish Task 10's exact gate and report, then Task 11.**

## Operating Mode And Slice Contract

**Heavy** (DEC-107/DEC-118/DEC-127): schema/persistence, public-interface, and
dependency changes. The binding process is `project-docs/specs/09_AI_WORKFLOW.md`;
required artifacts are Phase Runner, QA plan, six mid-phase reviews, final
review, Human checklist, full spec sync, run ledger, and triage log.
`../phases/phase-20/contracts.md` defines the accepted CLI behavior, Response
Envelope, leaf schemas, reason codes, DDL/Event semantics, removal inventory,
and legacy expiry ledger.

## Task Index

| Task | Title | Status | Dependencies |
|---|---|---|---|
| 1 | Response contract foundation | Complete; VAL-002 evidence owed by 3a | None |
| 2 | Project selection and config namespace | Complete; VAL-043 evidence owed by 3b | 1 |
| 3 | Atomic project initialization | Complete | 1-2 |
| 3a-3d | Reason-code, strict-config, journey, and gate-record evidence recovery | Complete (DEC-130/131/132) | 1-3 |
| 3e-3l | Infrastructure performance pack | Complete (DEC-133/134) | review-1c |
| 4 | Rebuild, authority, and lock failures | Complete `e56b3a0` (DEC-135/136) | 1-3, 3a-3l |
| 5 | Remove orphaned Issue fields | Complete `028855b` (DEC-137) | 4 |
| 6 | Replayable Issue aggregate revision | Complete `29c5246` | 5 |
| 7 | Issue hierarchy invariant | Complete `4950e94` | 5-6 |
| 8 | Issue read contract | Complete `4df4d66` (DEC-139) | 5-7 |
| 9 | Issue mutation and concurrency contract | Complete `86d927d`; evidence owed by 9a (DEC-140) | 6-8 |
| 9a | Mutation concurrency/transition/attribution evidence | Complete `1f196a9`; Review 3b NON-BLOCK `db36832` (DEC-140) | 9 |
| 10 | Acceptance-criterion and label contracts | Recovery in progress (DEC-141); exact candidate gate/report pending | 8-9a |
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
recovery sessions and 9a is the DEC-140 recovery session. Tasks 3e-3l are the
DEC-133 infrastructure pack and add no product VAL; `review-1b/1c.md`,
`infra-review-1.md`, and `review-3b.md` are bounded artifacts, not cadence
reviews. Review 3b is NON-BLOCK; the next cadence review remains Review 4 after
Task 12.

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

All under `../phases/phase-20/qa/`. `validation-contract.md` is Human-approved
(VAL-002/VAL-043 windows extended to tasks 3a-3c) and `qa-plan.md` is accepted
and unchanged. Cadence reviews `review-M.md` close after Tasks 3, 6, 9, 12, 15,
17; review 3 is BLOCK (DEC-140) and review 6 adds independent contract/spec
consistency. Bounded reviews: `review-1b.md` (BLOCK, corrected), `review-1c.md`
(NON-BLOCK), `infra-review-1.md` (BLOCK, re-reviewed NON-BLOCK), `review-3b.md`
(NON-BLOCK after Task 9a). Final review and Human checklist are pending and must
re-run the full gate with real tokens.

## Active Constraints

- **A task acceptance criterion summarizes a validation contract row; it never
  replaces it (DEC-140).** Where a VAL's Evidence or Pass column names a matrix,
  table, or "every X", assert one observable case per cell for every named
  operation and cite the cell-to-test mapping in the completion report.
  Representative or single-operation evidence is a failed criterion.
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
- **Parity matrix:** retained JSON leaves get the closed envelope/request/order
  contract; Human mode keeps concise stdout/stderr; legacy leaves keep their
  unenveloped mode and reject root flags. Evidence: VAL-001..005, VAL-008,
  VAL-029, VAL-044, VAL-045.

## Interfaces / Subsystems In Scope

`cmd/fid/`, `internal/cli/`, `internal/project/` (dispatch and transport);
`internal/config/`, `internal/bootstrap/` (strict config and atomic init);
Issue/label/link/workflow/maintenance packages (journey behavior);
`internal/journal/`, `internal/db/`, `migrations/` (revision and rebuild);
`internal/format/` and command tests (formatter removal).

## Known Risks

- **Realized in Review 1 and again in Review 3:** wide tasks produce
  representative rather than exhaustive evidence. Where a VAL names a domain or
  matrix, every just-inside/outside pair and registry member must be asserted
  (DEC-130), for every operation the row names (DEC-140). A green broad command
  never evidences a conditional state nobody asserted.
- **Realized in Review 1b:** record each gate with a verbatim token from that
  command; `exit 0 (no output)` is never valid for `go test`, and an uncaptured
  gate is `UNVERIFIED`, never `PASS` (DEC-131/132).
- Schema/fixture churn is isolated in Tasks 5-6; formatter removal is high-churn
  across Tasks 14-16. Exact VAL-035 inventories include `fid issue ac verify`;
  family ellipses are prohibited.

## Success / Exit Criteria

- [ ] All eighteen tasks plus recovery tasks 3a-3d and 9a complete with exact-command evidence.
- [ ] Every accepted VAL has a phase-end outcome with cited evidence.
- [ ] Complete standalone Issue journey passes and survives rebuild unchanged.
- [ ] Universal formatters/flags/dependency, promised stubs, and retired Issue fields are absent.
- [ ] Fail-closed V1 authority still rejects; no Phase 21 capability is advertised.
- [ ] Final QA and Human checklist have no unresolved blocking outcome.

## Next Planning Checkpoint

After Phase 20 wrap-up, resolve Gate D3.2: the exact one-Issue Routine walking
skeleton boundary for Phase 21.
