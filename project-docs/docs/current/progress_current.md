# Progress - Current

> In-flight status only. Historical implementation progress remains in phase
> archives.

## Project State

**No implementation phase is active. Phase 20 completed and wrapped on
2026-08-22.** Its final deterministic review is PASS, Human VAL-038 is 3/3
PASS, all 46 validation assertions pass, and the run ledger is closed. No
product, architecture, data-model, public-interface, dependency, workflow, or
Human decision remains open for Phase 20.

Phase 20 delivered the canonical JSON standalone Issue journey and retained
durable kernel behavior while removing universal formatters, `--format`,
promised stubs, and retired Issue fields. The real-process journey passes before
and after rebuild with byte-identical canonical Issue reads.

## Active Direction

ADR-0011 retains the local Issue tracker and durable operational kernel while
removing pre-dogfood V1 compatibility. ADR-0012 delivers convergence through
vertical scenarios. Fail-closed V1 authority remains until Phase 21 proves its
generic Actor/Envelope replacement.

## Carry-Forward

The authoritative register is `decisions_active.md`.

- CF-028: split atomic-init staging/fsync helpers on the next qualifying reopen.
- CF-029: extend evidence-token validation to QA review artifacts.
- CF-030: improve rejected-token diagnostics in `report-lint.sh`.
- CF-031: closed in Phase 20.
- CF-032: add invalid-input tests for `BuildIssueRevisionEffects` on its next
  implementation change.

## Next Work

1. Run the recommended deep retrospective on Phase 20's temporary V1 delivery
   workflow.
2. Resolve Gate D3.2's exact one-Issue Routine walking-skeleton boundary.
3. Plan Phase 21 only after that boundary is Human-accepted.

No implementation task should cite the superseded historical version roadmap
as forward authority.
