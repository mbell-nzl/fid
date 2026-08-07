# Phase 19 Human Verification Checklist

> Phase: 19 - Workflow Authoring CLI + Escalation Default Workflow
> Date: 2026-08-02
> Tester: MBEL
>
> This checklist contains only Human product/perception judgments. Run it in a
> fresh scratch `.fid/` workspace using the shipped help text and CLI spec.
> Deterministic checks are recorded in `final-review.md`. VAL-026 may pass only
> when the Human records PASS for all nine rows.

## Workflow Authoring Operator Coherence

One judgment per row. Record `pass`, `fail`, or `partial` and a concise note;
include the exact error where requested.

| Test ID | Covers | Test | Expected Result | Outcome | Notes |
|---|---|---|---|---|---|
| HUMAN-26a | VAL-026 | Judge whether `fid workflow define --help` and the CLI spec are sufficient to define one valid custom workflow without source-reading `(human)` | The Human can define the workflow from shipped operator guidance alone and finds the guidance coherent. | PASS | Defined `hv19-strict` from the shipped help and CLI-spec excerpt without source-reading. Required flags, state/category ordering, transition syntax, requirements, and cancellation behavior were coherent. |
| HUMAN-26b | VAL-026 | Judge whether `fid workflow assign --help` makes kind assignment discoverable and understandable `(human)` | Kind assignment is discoverable and understandable without source-reading. | PASS | Help made the workflow name, `--kind` / `--project` scope choice, and acting role clear; kind assignment was readily discoverable. |
| HUMAN-26c | VAL-026 | Judge whether JSON workflow-list output is coherent for operator inspection `(human)` | The JSON structure is coherent for operator inspection. | PASS | Summary fields, counts, applicability, and assignments were coherent for discovering workflows. Full named-workflow graph inspection is a separate product gap carried forward as CF-026. |
| HUMAN-26d | VAL-026 | Judge whether TOON workflow-list output is coherent for agent-oriented inspection `(human)` | The TOON structure is coherent for agent-oriented inspection. | PASS | TOON output was coherent for agent-oriented discovery and comparison of the listed workflow summaries. |
| HUMAN-26e | VAL-026 | Judge whether Markdown workflow-list output is scannable for a human `(human)` | The Markdown output is readily scannable. | PASS | Markdown workflow summaries were readily scannable. |
| HUMAN-26f | VAL-026 | Judge whether YAML workflow-list output is coherent for configuration-oriented inspection `(human)` | The YAML structure is coherent for configuration-oriented inspection. | PASS | YAML output was coherent for configuration-oriented inspection. |
| HUMAN-26g | VAL-026 | Trigger one invalid definition and judge whether the error supports self-correction `(human)` | The error identifies the problem clearly enough to support self-correction. | PASS | Exact error: `error: invalid workflow slug "Invalid Workflow"`. It names the invalid value and supports correction. |
| HUMAN-26h | VAL-026 | Trigger one invalid assignment and judge whether the error supports self-correction `(human)` | The error identifies the problem clearly enough to support self-correction. | PASS | Exact error: `error: workflow assignment requires a kind or project scope`. It identifies the missing choice and supports correction. |
| HUMAN-26i | VAL-026 | Encounter the DEC-079 custom-state transition limitation and judge whether the shipped message/documentation makes the temporary boundary understandable rather than surprising `(human)` | The temporary limitation is understandable from the message/documentation and is not surprising. | PASS | Exact error: `error: authorize action "issue.transition.backlog": no permissions matrix entry`. The DEC-079 documentation explains this temporary limitation and Phase 20's `fid policy set` path clearly. |

## Summary

| Section | Total | Pass | Fail | Partial | Blank |
|---|---:|---:|---:|---:|---:|
| Workflow Authoring Operator Coherence | 9 | 9 | 0 | 0 | 0 |
| **Total** | **9** | **9** | **0** | **0** | **0** |

| | |
|---|---|
| **Overall result** | PASS — VAL-026 is PASS |
| **Date completed** | 2026-08-07 |
| **Notes** | All nine Human judgments passed. CF-026 carries the distinct workflow-inspection and workflow-lifecycle product gap for Tech Lead review. |

<!--
Outcome values: pass, fail, partial, blank.
Overall result: pass only if all nine rows pass; partial if any partial and no
fail; fail if any fail; blocked while any row remains blank.
-->
