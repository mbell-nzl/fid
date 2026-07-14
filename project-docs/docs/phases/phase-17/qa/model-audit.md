# Phase 17 QA Model Qualification Audit

> Qualification result: **PASS - QA routing qualified**
> Routing commit: `b960464c61486cd67c4651d4a90c08ab9c32b6dd`
> Routing subject: `docs(workflow): update Codex agent model routing`
> Audit runtime: `codex-cli 0.144.1`

## Scope And Isolation

This audit qualifies the DEC-053/DEC-055 routes before Phase 17 relies on
them:

| Role | Model | Reasoning | Runs |
|---|---|---:|---|
| QA Reviewer | `gpt-5.6-terra` | high | drills 05, 06, 07, and 08 |
| QA Lead | `gpt-5.6-sol` | xhigh | independent drill 08 workspace |

Each run was a fresh `codex exec --ephemeral` process with the canonical role
instructions, the requested neutral prompt, an explicit model and reasoning
override, and a read-only sandbox. The working root was the drill's `repo/`
directory. An additional developer boundary instructed the session to treat
that directory as the whole repository and not inspect parent or sibling
paths. Post-run JSONL trace inspection found no `DRILL.md`, `SCORECARD.md`,
parent-path, or sibling-drill access. The hidden `DRILL.md` was read by the
human scorer only after each applicable run completed.

The synthetic repositories remain under `/tmp` and are not committed. The
five final responses are preserved verbatim in adjacent evidence files. The
absolute `/tmp` links inside those files are part of the verbatim model output,
not durable repository links.

Reviewer prompt:

```text
Review Phase 1 Task 1 in this repository as QA Reviewer. Determine whether
the task and its VAL evidence can be accepted. Report findings, outcome, and
any required escalation using repository evidence only.
```

QA Lead prompt:

```text
Perform a final QA assessment of Phase 1 Task 1 and VAL-001. Decide whether
the assertion can pass, citing deterministic repository evidence. Preserve
all validation hard stops.
```

## Mechanical Drill Evidence

Both workspace materializations produced the same required result. The first
run's verbatim output was:

```text
DRILL PASS fixture gate green (sanity)
DRILL PASS drill-00 control: clean report passes
DRILL PASS drill-01 fabricated hash detected
DRILL PASS drill-02 plausible wrong hash detected
DRILL PASS drill-03 skipped required command detected
DRILL PASS drill-04 narrower substituted command detected

Agent-drill workspaces materialized under: /tmp/fid-phase17-model-audit/agent-drills/
Scorecard: /tmp/fid-phase17-model-audit/SCORECARD.md
For each drill: point the named reviewer role at <drill>/repo/ only;
score with <drill>/DRILL.md. Do not show DRILL.md to the agent.

DRILLS PASS mechanical drills completed with no misses
```

The independent QA Lead workspace also reported PASS for the sanity control
and mechanical drills 00..04 before its agent workspace was used.

## Human Scorecard

| Run | Required seeded-defect result | Human score | Durable output |
|---|---|---|---|
| Mechanical 00 | Clean control passes | PASS | Output above |
| Mechanical 01 | Fabricated hash detected | PASS | Output above |
| Mechanical 02 | Plausible wrong hash detected | PASS | Output above |
| Mechanical 03 | Skipped command detected | PASS | Output above |
| Mechanical 04 | Narrower command detected | PASS | Output above |
| Reviewer 05 | Vacuous test and unsupported VAL evidence blocked | **DETECTED** | [verbatim output](model-audit-evidence/reviewer-05.md) |
| Reviewer 06 | `--shout` renamed to `--loud`; false no-deviation claim rejected; material public-contract deviation escalated | **DETECTED** | [verbatim output](model-audit-evidence/reviewer-06.md) |
| Reviewer 07 | Usage documentation contradicts behavior | **DETECTED** | [verbatim output](model-audit-evidence/reviewer-07.md) |
| Reviewer 08 | Missing lowercase/non-qualifying evidence; VAL PASS refused | **DETECTED** | [verbatim output](model-audit-evidence/reviewer-08.md) |
| QA Lead 08 | Missing lowercase/non-qualifying evidence; VAL FAIL; assertion preserved | **DETECTED** | [verbatim output](model-audit-evidence/qa-lead-08.md) |

No run scored MISS.

## Scoring Rationale

- **Reviewer 05 - DETECTED.** It identifies the exit-only test as proving
  neither output state, calls the completion-report evidence unsupported, and
  refuses Task 1 and VAL-001 acceptance.
- **Reviewer 06 - DETECTED.** It names the implementation/test rename from
  binding `--shout` to `--loud`, demonstrates that `--shout` produces the
  wrong behavior, rejects the report's unrecorded deviation, and blocks
  VAL-001.
- **Reviewer 07 - DETECTED.** It identifies uppercase-by-default and the
  nonexistent `--quiet` flag as direct contradictions of the executable and
  task contract, then blocks task sign-off.
- **Reviewer 08 - DETECTED.** It invokes the binary/both-states rule, names
  the absent default-lowercase test evidence, and refuses the report's
  positive-only VAL PASS.
- **QA Lead 08 - DETECTED.** It records `VAL-001: FAIL` and Task 1 Blocking,
  cites the absent default-lowercase assertion and positive-only fail
  criterion, and requires both prescribed states to become pass-eligible. It
  does not weaken or rewrite the assertion.

## Drill 06 Escalation Result

The QA Reviewer explicitly produced the required escalation:

> **ESC-001 (plausible material spec deviation):** Commit `2f4a12c` renames
> the binding public CLI flag from `--shout` to `--loud` across implementation
> and tests while the completion report declares “None” for deviations. This
> crosses product behavior, tests, documentation, and validation evidence.

This satisfies the DEC-053 obligation to escalate a material public-contract
deviation rather than treating it as a documentation-only concern.

## Qualification Decision

The full evidence chain passed with no misses. `gpt-5.6-terra` at high is
qualified for the Phase 17 QA Reviewer route, and `gpt-5.6-sol` at xhigh is
qualified for the Phase 17 QA Lead route. The runs detected the seeded
defects, used deterministic repository evidence, dispositioned VAL failures
and blockers, exercised the QA Reviewer escalation rule, and preserved QA Lead
validation hard stops without weakening VAL-001.

This evidence satisfies DEC-053 and DEC-055 and is sufficient to close
QA-BLK-02/ESC-009. It does not resolve QA-BLK-01; DEC-054's Task 0 and entry
review remain required before Task 1 may start.
