# ESC-20-015 final independent workflow review

**Date:** 2026-08-16  
**Reviewer:** Fresh read-only subagent  
**Reviewed checkpoint:** `2b93a041ecfe36ba4cb2437b6dbd1cce3311f709`  
**Verdict:** PASS

## Scope

The reviewer inspected `a50d8e4..2b93a04`, the canonical workflow spec,
prompts, templates, generated copies, sequence validator, control-scope
verifier, mechanical drills, and the Phase 20 ledger. The review was bounded to
workflow infrastructure; it did not reattribute or rerun Task 17 product
evidence.

## Independent adversarial results

- Exact STOP `Resolution authority` retention was independently exercised.
  Suffix-expanded authority in either the next recovery's `Source` or
  `Repair authority` was rejected.
- Schema-v2 TRIAGE rejected a duplicate `Trigger` and requires the field
  exactly once.
- A trailing-token schema-v1 RUN candidate could not bypass append-time v2
  adoption. A schema-v2 RUN heading with trailing whitespace was rejected as
  noncanonical.
- Both rejected append attempts left their ledgers byte-for-byte unchanged;
  the ledgers still validated `ADVANCEABLE` afterward.
- DEC-shaped checkpoint, contract, evidence, and command prose remained valid.
  Nonexistent or embedded pseudo decision IDs in actual `Repair authority` or
  `Resolution authority` fields were rejected.
- A scope receipt with a different allowlist order from the authoritative
  directive was rejected.
- The verifier materialized from the trusted base rejected an attempted repair
  of `workflow-control-scope-lint.sh` itself.

## Required controls

| Control | Result | Terminal token |
|---|---|---|
| Shell syntax for workflow tools | PASS | `exit_code: 0` |
| `project-docs/tools/report-lint.sh --self-test` | PASS | `SELF-TEST PASS report-lint self-test completed with no failures` |
| Exact Task 17 report lint, without `--repo` | PASS | `PASS report-lint completed with no failures` |
| `project-docs/tools/workflow-lint.sh` | PASS | `PASS workflow-lint completed with no failures` |
| `project-docs/tools/drills/run-drills.sh` | PASS | `DRILLS PASS mechanical drills completed with no misses` |
| Agent and skill generated-copy checks | PASS | generated copies match canonical sources |
| Phase 20 sequence validation | PASS | `VALID state=ADVANCEABLE` |
| `git diff --check` | PASS | `exit_code: 0` with no output |
| Final tree state | PASS | clean at the reviewed checkpoint |

## Residual risk

- Transactional append assumes the accepted single-Phase-Runner operating
  model; concurrent writers are not serialized by an inter-process lock.
- Shell and awk behavior was exercised in the current project environment, not
  across a portability matrix.
- Product Go gates were deliberately not rerun because Task 17's committed
  product evidence is immutable and this recovery changed workflow
  infrastructure only.

These are bounded residual risks, not blockers for resuming Phase 20 at Review
6.
