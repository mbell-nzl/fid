# Progress - Current

> In-flight status only. Historical implementation progress remains in phase
> archives.

## Project State

**Phase 21 planning is active; implementation is not approved.** Phase 20
completed and wrapped on 2026-08-22 with deterministic PASS, Human VAL-038 at
3/3 PASS, all 46 assertions passing, and a closed run ledger.

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

1. **Complete:** Phase 20 deep retrospective and DEC-152.
2. **Complete:** Human-accepted Gate D3.2 and DEC-153.
3. Draft and present the exact Phase 21 slice/DDL/Event contract, validation
   contract, ordered two-group task pack, and Heavy-mode phase brief.
4. Stop before implementation for Human execution approval; a Ready QA plan
   with fresh feasibility evidence is also required before Task 1.

No implementation task should cite the superseded historical version roadmap
as forward authority.
