# Phase 1 Validation Contract - Drill Fixture Phase

> Phase: 1 - Drill Fixture Phase
> Author: Tech Lead (draft)
> Last refined: 2026-07-04 by Tech Lead (applying QA Lead proposals)

## Purpose

Fixture contract for evidence-chain drills. One assertion, deliberately of the
conditional/binary shape, so positive-only evidence can be seeded and must be
caught.

## Assertions

| ID | Source requirement | Behavior or invariant | Evidence method | Validator | Owning task/window | Pass criterion | Fail criterion | Status | Out-of-scope notes |
|---|---|---|---|---|---|---|---|---|---|
| VAL-001 | fixture-PRD §1 | `--shout` uppercases the greeting; without `--shout` the greeting stays lowercase | `sh test/greet_test.sh` output covering both states, plus one real-entrypoint run | QA Reviewer | task-1 | Both the qualifying (`--shout` uppercase) and non-qualifying (default lowercase) states are proven by test output | Either state unproven, or positive-only evidence offered for a two-state assertion | Ready | Flags other than `--shout` |

## Outcomes

Phase-end final review records every `VAL-NNN` outcome. `PASS` without cited
evidence is invalid; conditional or binary assertions cannot pass on
positive-only evidence.
