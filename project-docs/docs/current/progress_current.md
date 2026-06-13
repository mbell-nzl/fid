# Progress - Current

> **Purpose.** In-flight phase status only. Loaded by default in every agent
> session. Historical phase progress lives under `../phases/phase-N/progress.md`.
>
> **Size target:** under 100 lines.

## Project State

**Status:** Phase 4 (v0.1b Issue authoring and discovery) is **CLOSED** on branch
`phase-4`. QA final review VERDICT: PASS (0 blocking / 0 major / 0 improvement);
human checklist 5/5 PASS. Phase 3 and earlier are archived under
`../phases/phase-N/`.

Phase 4 delivered:

- Sequence-vs-scheduling backstop trigger (Task 1, CF-004 closed).
- Issue authoring domain service + `fid issue create` (Tasks 2-3).
- Centralised enriched issue detail read model + `fid issue get --section`
  (Task 4, CF-006 closed); due-date public read-back closed ESC-002.
- `fid issue list` filters/order/formats (Task 5).
- COV-001/COV-002 closeout (Task 6): `internal/issues` 85.5%, `cmd/fid` 79.5%.

Validation contract: VAL-001..VAL-010 all PASS (VAL-010 via human dogfood).
Six tasks, zero spec deviations, zero dependency changes.

**Needs human attention:** During the VAL-010 dogfood, a rejected
`fid issue create` was found to append its event to the canonical journal and
poison the tail, breaking later reads until the bad line is removed. Tracked as
**CF-007** (High, Phase 5) / **DEC-014**. Please confirm CF-007 priority before
Phase 5 planning closes. This did not change the Phase 4 PASS verdict.

Open carry-forward register: CF-003 (role corrupt-tail tolerance, before v0.1d),
CF-005 (`activity_log` projection guard / `activity` effect, Phase 5), CF-007
(rejected-create journal poisoning, Phase 5).

Next action:

- Plan Phase 5 (v0.1b mutation + workflow transitions): `fid issue update`,
  `activity_log` projection discipline / `activity` effects (CF-005), hard-coded
  workflow graphs/rules, `fid issue transition`, and the CF-007 write-path fixup.
