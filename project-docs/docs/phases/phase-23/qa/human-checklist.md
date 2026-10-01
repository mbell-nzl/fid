# Phase 23 Human Verification Checklist

> Phase: 23 — Review and Human judgment
> Date: 2026-10-01
> Tester: MBEL

Human run completed 2026-10-01: 55 pass, 2 fail, 6 partial, 0 blank; overall FAIL (see Summary and Tech Lead analysis at the end). Only `(human)` judgments appear here. Deterministic results and feasibility evidence are in [final-review.md](final-review.md). Use the retained final-binary fixtures below. Record your answer and a short output quotation where requested.

## Fixture identity

- Manifest: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-3134551499/manifest.json`
- Binary: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-3134551499/fid`
- SHA-256: `sha256:3a34fc318991f77b9aba5136e5e6a04538824010489dbae31fd854bae01f5d49`
- Product source: `e0893cdee915013371455198906852081633e051`. The recorded dirty state consists of workflow documents and these QA artifacts; product files were unchanged.
- This binary is byte-identical to `/tmp/fid` built by the full phase gate. Keep it with these workspaces.
- These fixtures contain short-lived Claims. Generate a fresh set immediately before the Human run using the guide linked below, then use that manifest's current paths and IDs. Exact expected output is bound to generation time.
- Run each section from its stated directory. The protected-Decision steps are sequential and consume that fixture; record the before observation first. For a repeat, use the fresh-generation command in [human-fixtures.md](human-fixtures.md#generate-a-fresh-set) and the resulting manifest paths.

## Mixed Findings — VAL-070

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-3134551499/mixed-findings`

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-3134551499/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-3134551499/mixed-findings review get REV-ab5d15
```

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-3134551499/mixed-findings`

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-3134551499/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-3134551499/mixed-findings review list --issue P22-1ce0e7
```

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-3134551499/mixed-findings`

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-3134551499/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-3134551499/mixed-findings finding list --issue P22-1ce0e7
```

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-3134551499/mixed-findings`

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-3134551499/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-3134551499/mixed-findings resume P22-1ce0e7 --repository /home/mark/projects/fid/.cache/phase-23-human/checkpoints-3134551499/mixed-findings
```

| Test ID | Covers | Test | Expected result | Outcome | Notes |
|---|---|---|---|---|---|
| HUMAN-01a | VAL-070 | Identify the resolved Findings from the displayed output. `(human)` | Name them and quote what establishes resolution. | PASS | There is one resolved finding called `FND-fixture-resolved`. Resolution is established by the `Disposition: resolved` text in all of the tested views. |
| HUMAN-01b | VAL-070 | Identify the Findings that are addressed but remain unresolved. `(human)` | Name them and quote the distinction. | PASS | `FND-fixture-addressed` is the one addressed finding, established by `Disposition: addressed`. The distinction is that `addressed` mean a developer/builder has marked this as being addressed but it has not yet been reviewed by a reviewer to mark it as resolved. |
| HUMAN-01c | VAL-070 | Explain what still blocks completion. `(human)` | Give the blocking Findings and the reason visible in fid output. | PASS | `OBL-p22-corpus` is unsatisifed and blocks completion. It comprises three open findings which are all marked as blocking: `FND-fixture-blocking`, `FND-fixture-new` and `FND-fixture-addressed`. All three clearly state `Blocking: true` |

## Protected Decision — VAL-071

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-3134551499/protected-decision`

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-3134551499/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-3134551499/protected-decision resume P22-a6bf3e --repository /home/mark/projects/fid/.cache/phase-23-human/checkpoints-3134551499/protected-decision
```

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-3134551499/protected-decision`

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-3134551499/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-3134551499/protected-decision decision list --issue P22-a6bf3e
```

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-3134551499/protected-decision`

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-3134551499/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-3134551499/protected-decision stop list --subject envelope:ENV-934a66
```

### Step 1: blocked action before Outcome

Expected command exit: `5`.

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-3134551499/protected-decision`

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-3134551499/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-3134551499/protected-decision --json issue transition P22-a6bf3e --to qa_pass --actor agent:actor-recovery --expected-revision 4 --idempotency-key fixture-human-transition
```

| Test ID | Covers | Test | Expected result | Outcome | Notes |
|---|---|---|---|---|---|
| HUMAN-02a-before | VAL-071 | Observe the guarded action before recording the Outcome. `(human)` | The output makes the refusal clear. | PARTIAL | Human-facing output (no `--json` flag) makes it clear that the action is not permitted. It could/should be improved to show the reason. I can, however, see in the `--json` output `"protected-decision-required"` and `"protected Decision requires an approved reference and Human Outcome"`. I also noted from the pre-requisite commands (`fid resume` and `fid stop list`) that the stop details are shown as raw json in the human output. I've marked this as PARTIAL as I'd like these items addressed.|

### Step 2: record authorized Human Outcome

Expected command exit: `0`.

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-3134551499/protected-decision`

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-3134551499/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-3134551499/protected-decision --json human-outcome record --input /home/mark/projects/fid/.cache/phase-23-human/checkpoints-3134551499/protected-decision/fixture-human-outcome.json --actor actor-human-outcome --idempotency-key fixture-human-outcome
```

### Step 3: retry guarded action

Expected command exit: `0`.

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-3134551499/protected-decision`

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-3134551499/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-3134551499/protected-decision --json issue transition P22-a6bf3e --to qa_pass --actor agent:actor-recovery --expected-revision 4 --idempotency-key fixture-human-transition
```

| Test ID | Covers | Test | Expected result | Outcome | Notes |
|---|---|---|---|---|---|
| HUMAN-02a-after | VAL-071 | Observe the same guarded action after recording the Outcome. `(human)` | The output makes the permitted action clear. | PASS | After recording the Human Outcome, the same transition returned exit `0` with `"ok":true`. The action was permitted.|
| HUMAN-02b | VAL-071 | Explain why the action became permitted, using only fid output. `(human)` | State the change in authority or qualifying fact shown by the output. | PASS | The authorized Human Outcome approved the security Decision. Step 2 showed the protected Decision as `satisfied` and its Stop as resolved; the identical transition then succeeded.|

## Review rework exhaustion — VAL-072

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-3134551499/review-exhaustion`

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-3134551499/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-3134551499/review-exhaustion --json claim acquire --input /home/mark/projects/fid/.cache/phase-23-human/checkpoints-3134551499/review-exhaustion/fixture-rework-claim.json --actor actor-recovery --idempotency-key fixture-exhausted-claim
```

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-3134551499/review-exhaustion`

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-3134551499/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-3134551499/review-exhaustion resume P22-db30b8 --repository /home/mark/projects/fid/.cache/phase-23-human/checkpoints-3134551499/review-exhaustion
```

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-3134551499/review-exhaustion`

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-3134551499/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-3134551499/review-exhaustion stop list --subject envelope:ENV-005ebb
```

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-3134551499/review-exhaustion`

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-3134551499/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-3134551499/review-exhaustion attention list --subject envelope:ENV-005ebb
```

| Test ID | Covers | Test | Expected result | Outcome | Notes |
|---|---|---|---|---|---|
| HUMAN-03a | VAL-072 | Explain the cause of the rework block. `(human)` | The cause is clear from the denial and displayed Stop/Attention. | PARTIAL | The cause of the block is that the review-cycle budget has been exhausted. I marked as PARTIAL because `fid attention list`, `fid stop list --subject` displayed output as raw json rather than formatted for human reader. The `fid resume` `Blockers and attention` block also is raw json. The initial command in this sequence (`fid claim acquire`) didn't show anything when I dropped the `--json` flag -- should it?|
| HUMAN-03b-resolution | VAL-072 | State the exact qualifying action that would resolve the block. `(human)` | The required resolution is clear without additional help. | FAIL | `budget-observed for envelope:ENV-005ebb budget review-cycle` - my intepretation is that I would need to observe that the budget has been exhausted and then record the observation to unblock the issue.|
| HUMAN-03b-internals | VAL-072 | Assess whether the output suggests editing fid internals. `(human)` | It guides you through supported actions without suggesting internal edits. | FAIL | Human readable output was insufficient. The json output led me to note `resolver capability stop.resolve` but doesn’t explain the action needed to produce the qualifying fact or where to obtain its `fact_event_id`. I could not determine the resolution procedure from fid output alone. |

## Every new view — VAL-073

The following 18 operation views are the scope of contracts §14. Each has three separate judgments. Mutation examples replay existing fixture requests. The JSON counterpart is available in the manifest if needed; judge the displayed Human text.

### requirement.record-reference

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-3134551499/leaf-contract`

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-3134551499/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-3134551499/leaf-contract requirement record-reference --input /home/mark/projects/fid/.cache/phase-23-human/checkpoints-3134551499/leaf-contract/leaf-leaf-requirement.json --actor actor-recovery --idempotency-key leaf-requirement
```

| Test ID | Covers | Test | Expected result | Outcome | Notes |
|---|---|---|---|---|---|
| HUMAN-04a-01 | VAL-073 | Assess the section labels. `(human)` | Labels make the displayed sections clear. | PASS | n/a |
| HUMAN-04b-01 | VAL-073 | Assess the operation purpose line. `(human)` | The view clearly states its purpose. | PASS | n/a |
| HUMAN-04c-01 | VAL-073 | Assess routine-state narration. `(human)` | The view is free of unnecessary routine-state narration. | PASS | n/a |

### requirement.link

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-3134551499/leaf-contract`

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-3134551499/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-3134551499/leaf-contract requirement link --input /home/mark/projects/fid/.cache/phase-23-human/checkpoints-3134551499/leaf-contract/leaf-leaf-link.json --actor actor-recovery --idempotency-key leaf-link
```

| Test ID | Covers | Test | Expected result | Outcome | Notes |
|---|---|---|---|---|---|
| HUMAN-04a-02 | VAL-073 | Assess the section labels. `(human)` | Labels make the displayed sections clear. | PASS | Like many of the other commands, seems fine for agent use, but I struggle to see how a human can work with many of these commands with `--input` flags that require formatted json (and the ability to know/remember the schema). |
| HUMAN-04b-02 | VAL-073 | Assess the operation purpose line. `(human)` | The view clearly states its purpose. | PASS | n/a |
| HUMAN-04c-02 | VAL-073 | Assess routine-state narration. `(human)` | The view is free of unnecessary routine-state narration. | PASS | n/a |

### concern.record

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-3134551499/leaf-contract`

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-3134551499/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-3134551499/leaf-contract concern record --input /home/mark/projects/fid/.cache/phase-23-human/checkpoints-3134551499/leaf-contract/leaf-leaf-concern.json --actor actor-recovery --idempotency-key leaf-concern
```

| Test ID | Covers | Test | Expected result | Outcome | Notes |
|---|---|---|---|---|---|
| HUMAN-04a-03 | VAL-073 | Assess the section labels. `(human)` | Labels make the displayed sections clear. | PASS | n/a |
| HUMAN-04b-03 | VAL-073 | Assess the operation purpose line. `(human)` | The view clearly states its purpose. | PASS | n/a |
| HUMAN-04c-03 | VAL-073 | Assess routine-state narration. `(human)` | The view is free of unnecessary routine-state narration. | PARTIAL | Applies to other commands, but while the hints `Purpose: Inspect a raised concern and its disposition as judgment.` and `This Concern is recorded judgment, not deterministic proof.` are useful for me now while learning how to use fid, they should eventually be scrubbed (maybe in Phase 25?) as clutter. This sort of info should be explained on the `--help` for each command and in a manual but not get in the way of a learned user. |

### concern.dispose

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-3134551499/leaf-contract`

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-3134551499/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-3134551499/leaf-contract concern dispose --input /home/mark/projects/fid/.cache/phase-23-human/checkpoints-3134551499/leaf-contract/leaf-leaf-dispose.json --actor actor-recovery --idempotency-key leaf-dispose
```

| Test ID | Covers | Test | Expected result | Outcome | Notes |
|---|---|---|---|---|---|
| HUMAN-04a-04 | VAL-073 | Assess the section labels. `(human)` | Labels make the displayed sections clear. |  PASS | n/a |
| HUMAN-04b-04 | VAL-073 | Assess the operation purpose line. `(human)` | The view clearly states its purpose. |  PASS | n/a |
| HUMAN-04c-04 | VAL-073 | Assess routine-state narration. `(human)` | The view is free of unnecessary routine-state narration. |  PASS | n/a |

### deviation.record

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-3134551499/leaf-contract`

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-3134551499/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-3134551499/leaf-contract deviation record --input /home/mark/projects/fid/.cache/phase-23-human/checkpoints-3134551499/leaf-contract/leaf-leaf-deviation.json --actor actor-recovery --idempotency-key leaf-deviation
```

| Test ID | Covers | Test | Expected result | Outcome | Notes |
|---|---|---|---|---|---|
| HUMAN-04a-05 | VAL-073 | Assess the section labels. `(human)` | Labels make the displayed sections clear. | PASS | n/a |
| HUMAN-04b-05 | VAL-073 | Assess the operation purpose line. `(human)` | The view clearly states its purpose. | PASS | n/a |
| HUMAN-04c-05 | VAL-073 | Assess routine-state narration. `(human)` | The view is free of unnecessary routine-state narration. | PASS | n/a |

### decision.record-reference

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-3134551499/leaf-contract`

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-3134551499/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-3134551499/leaf-contract decision record-reference --input /home/mark/projects/fid/.cache/phase-23-human/checkpoints-3134551499/leaf-contract/leaf-leaf-decision.json --actor actor-recovery --idempotency-key leaf-decision
```

| Test ID | Covers | Test | Expected result | Outcome | Notes |
|---|---|---|---|---|---|
| HUMAN-04a-06 | VAL-073 | Assess the section labels. `(human)` | Labels make the displayed sections clear. | PASS | n/a |
| HUMAN-04b-06 | VAL-073 | Assess the operation purpose line. `(human)` | The view clearly states its purpose. | PASS | n/a |
| HUMAN-04c-06 | VAL-073 | Assess routine-state narration. `(human)` | The view is free of unnecessary routine-state narration. | PASS | n/a |

### human-outcome.record

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-3134551499/leaf-contract`

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-3134551499/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-3134551499/leaf-contract human-outcome record --input /home/mark/projects/fid/.cache/phase-23-human/checkpoints-3134551499/leaf-contract/leaf-leaf-outcome.json --actor actor-recovery --idempotency-key leaf-outcome
```

| Test ID | Covers | Test | Expected result | Outcome | Notes |
|---|---|---|---|---|---|
| HUMAN-04a-07 | VAL-073 | Assess the section labels. `(human)` | Labels make the displayed sections clear. | PASS | n/a |
| HUMAN-04b-07 | VAL-073 | Assess the operation purpose line. `(human)` | The view clearly states its purpose. | PASS | n/a |
| HUMAN-04c-07 | VAL-073 | Assess routine-state narration. `(human)` | The view is free of unnecessary routine-state narration. | PASS | n/a |

### review.record

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-3134551499/leaf-contract`

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-3134551499/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-3134551499/leaf-contract review record --input /home/mark/projects/fid/.cache/phase-23-human/checkpoints-3134551499/leaf-contract/leaf-leaf-review.json --actor actor-denied --idempotency-key leaf-review
```

| Test ID | Covers | Test | Expected result | Outcome | Notes |
|---|---|---|---|---|---|
| HUMAN-04a-08 | VAL-073 | Assess the section labels. `(human)` | Labels make the displayed sections clear. | PASS | n/a |
| HUMAN-04b-08 | VAL-073 | Assess the operation purpose line. `(human)` | The view clearly states its purpose. | PASS | n/a |
| HUMAN-04c-08 | VAL-073 | Assess routine-state narration. `(human)` | The view is free of unnecessary routine-state narration. | PASS | n/a |

### finding.address

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-3134551499/leaf-contract`

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-3134551499/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-3134551499/leaf-contract finding address --input /home/mark/projects/fid/.cache/phase-23-human/checkpoints-3134551499/leaf-contract/leaf-leaf-address.json --actor actor-recovery --idempotency-key leaf-address
```

| Test ID | Covers | Test | Expected result | Outcome | Notes |
|---|---|---|---|---|---|
| HUMAN-04a-09 | VAL-073 | Assess the section labels. `(human)` | Labels make the displayed sections clear. | PASS | Should the reason be included in the human output? |
| HUMAN-04b-09 | VAL-073 | Assess the operation purpose line. `(human)` | The view clearly states its purpose. | PASS | n/a |
| HUMAN-04c-09 | VAL-073 | Assess routine-state narration. `(human)` | The view is free of unnecessary routine-state narration. | PASS | n/a |

### review.get

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-3134551499/leaf-contract`

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-3134551499/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-3134551499/leaf-contract review get REV-383b24
```

| Test ID | Covers | Test | Expected result | Outcome | Notes |
|---|---|---|---|---|---|
| HUMAN-04a-10 | VAL-073 | Assess the section labels. `(human)` | Labels make the displayed sections clear. | PASS | n/a |
| HUMAN-04b-10 | VAL-073 | Assess the operation purpose line. `(human)` | The view clearly states its purpose. | PASS | n/a |
| HUMAN-04c-10 | VAL-073 | Assess routine-state narration. `(human)` | The view is free of unnecessary routine-state narration. | PASS | n/a |

### review.list

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-3134551499/leaf-contract`

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-3134551499/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-3134551499/leaf-contract review list --issue P22-90c968
```

| Test ID | Covers | Test | Expected result | Outcome | Notes |
|---|---|---|---|---|---|
| HUMAN-04a-11 | VAL-073 | Assess the section labels. `(human)` | Labels make the displayed sections clear. | PASS | As per existing carry-forward, human-facing output formatting could be much improved (Phase 25 maybe??) to make easier to consume. |
| HUMAN-04b-11 | VAL-073 | Assess the operation purpose line. `(human)` | The view clearly states its purpose. | PASS | n/a |
| HUMAN-04c-11 | VAL-073 | Assess routine-state narration. `(human)` | The view is free of unnecessary routine-state narration. | PASS | n/a |

### finding.list

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-3134551499/leaf-contract`

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-3134551499/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-3134551499/leaf-contract finding list --issue P22-90c968
```

| Test ID | Covers | Test | Expected result | Outcome | Notes |
|---|---|---|---|---|---|
| HUMAN-04a-12 | VAL-073 | Assess the section labels. `(human)` | Labels make the displayed sections clear. | PASS | n/a |
| HUMAN-04b-12 | VAL-073 | Assess the operation purpose line. `(human)` | The view clearly states its purpose. | PASS | n/a |
| HUMAN-04c-12 | VAL-073 | Assess routine-state narration. `(human)` | The view is free of unnecessary routine-state narration. | PASS | n/a |

### concern.list

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-3134551499/leaf-contract`

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-3134551499/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-3134551499/leaf-contract concern list --subject issue:P22-90c968
```

| Test ID | Covers | Test | Expected result | Outcome | Notes |
|---|---|---|---|---|---|
| HUMAN-04a-13 | VAL-073 | Assess the section labels. `(human)` | Labels make the displayed sections clear. | PASS | n/a |
| HUMAN-04b-13 | VAL-073 | Assess the operation purpose line. `(human)` | The view clearly states its purpose. | PASS | n/a |
| HUMAN-04c-13 | VAL-073 | Assess routine-state narration. `(human)` | The view is free of unnecessary routine-state narration. | PASS | n/a |

### deviation.list

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-3134551499/leaf-contract`

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-3134551499/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-3134551499/leaf-contract deviation list --subject issue:P22-90c968
```

| Test ID | Covers | Test | Expected result | Outcome | Notes |
|---|---|---|---|---|---|
| HUMAN-04a-14 | VAL-073 | Assess the section labels. `(human)` | Labels make the displayed sections clear. | PASS | n/a |
| HUMAN-04b-14 | VAL-073 | Assess the operation purpose line. `(human)` | The view clearly states its purpose. | PASS | n/a |
| HUMAN-04c-14 | VAL-073 | Assess routine-state narration. `(human)` | The view is free of unnecessary routine-state narration. | PASS | n/a |

### decision.list

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-3134551499/leaf-contract`

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-3134551499/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-3134551499/leaf-contract decision list --issue P22-90c968
```

| Test ID | Covers | Test | Expected result | Outcome | Notes |
|---|---|---|---|---|---|
| HUMAN-04a-15 | VAL-073 | Assess the section labels. `(human)` | Labels make the displayed sections clear. | PASS | n/a |
| HUMAN-04b-15 | VAL-073 | Assess the operation purpose line. `(human)` | The view clearly states its purpose. | PASS | n/a |
| HUMAN-04c-15 | VAL-073 | Assess routine-state narration. `(human)` | The view is free of unnecessary routine-state narration. | PASS | n/a |

### human-outcome.list

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-3134551499/leaf-contract`

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-3134551499/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-3134551499/leaf-contract human-outcome list --issue P22-90c968
```

| Test ID | Covers | Test | Expected result | Outcome | Notes |
|---|---|---|---|---|---|
| HUMAN-04a-16 | VAL-073 | Assess the section labels. `(human)` | Labels make the displayed sections clear. | PASS | n/a |
| HUMAN-04b-16 | VAL-073 | Assess the operation purpose line. `(human)` | The view clearly states its purpose. | PASS | n/a |
| HUMAN-04c-16 | VAL-073 | Assess routine-state narration. `(human)` | The view is free of unnecessary routine-state narration. | PASS | n/a |

### requirement.get

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-3134551499/leaf-contract`

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-3134551499/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-3134551499/leaf-contract requirement get REQ-leaf
```

| Test ID | Covers | Test | Expected result | Outcome | Notes |
|---|---|---|---|---|---|
| HUMAN-04a-17 | VAL-073 | Assess the section labels. `(human)` | Labels make the displayed sections clear. | PARTIAL | Returned info for `REQ-leaf` but also for `RREF-7a2b10` -- is this expected? I assumed a `get` command fetches the requested ID only. I thought this looked more like a `list` output, but when I tried `fid requirement list` here, it only showed `RREF-7a2b10` and not `REQ-leaf`.|
| HUMAN-04b-17 | VAL-073 | Assess the operation purpose line. `(human)` | The view clearly states its purpose. | PARTIAL | see above |
| HUMAN-04c-17 | VAL-073 | Assess routine-state narration. `(human)` | The view is free of unnecessary routine-state narration. | PARTIAL | see above |

### requirement.list

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-3134551499/leaf-contract`

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-3134551499/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-3134551499/leaf-contract requirement list
```

| Test ID | Covers | Test | Expected result | Outcome | Notes |
|---|---|---|---|---|---|
| HUMAN-04a-18 | VAL-073 | Assess the section labels. `(human)` | Labels make the displayed sections clear. | PASS | Only shows `RREF-7a2b10` -- are more expected? See notes above for `requirement get` |
| HUMAN-04b-18 | VAL-073 | Assess the operation purpose line. `(human)` | The view clearly states its purpose. | PASS | n/a |
| HUMAN-04c-18 | VAL-073 | Assess routine-state narration. `(human)` | The view is free of unnecessary routine-state narration. | PASS | n/a |

## Summary

| Section | Total | Pass | Fail | Partial | Blank |
|---|---:|---:|---:|---:|---:|
| Mixed Findings | 3 | 3 | 0 | 0 | 0 |
| Protected Decision | 3 | 2 | 0 | 1 | 0 |
| Rework exhaustion | 3 | 0 | 2 | 1 | 0 |
| Eighteen new views | 54 | 50 | 0 | 4 | 0 |
| **Total** | **63** | **55** | **2** | **6** | **0** |

| | |
|---|---|
| **Overall result** | FAIL |
| **Date completed** | 2026-10-01 |
| **Notes** | FAIL: HUMAN-03b-resolution and HUMAN-03b-internals. The rework-exhaustion resolution procedure, including how to produce the qualifying fact and its `fact_event_id`, could not be determined from fid output, and Tech Lead analysis finds no supported path in Phase 23 (finding A1 below). PARTIAL: HUMAN-02a-before and HUMAN-03a. Human-mode denials omit the reason. Stop details render as raw JSON in `resume`, `stop list`, and `attention list`. Human-mode `claim acquire` printed nothing on the rework denial. PARTIAL: HUMAN-04c-03. Purpose and judgment hint lines help while learning but become clutter for experienced users; move them to `--help` or a manual (Phase 25?). This note applies to every view that has these lines. PARTIAL: HUMAN-04a/b/c-17. `requirement get REQ-leaf` shows both `RREF-7a2b10` and `REQ-leaf`, each labelled "Requirement"; `requirement list` shows only `RREF-7a2b10`. Cross-cutting, not scored: `--input` JSON is hard for a Human to author without knowing the schema (HUMAN-04a-02); Human view formatting needs improvement under CF-040 (HUMAN-04a-11); `finding address` Human output omits the reason (HUMAN-04a-09). |

Outcome values: pass, fail, partial, blank. Overall result is pass only when every row passes; fail if any fails; otherwise partial while any judgment is partial or blank.

## Tech Lead analysis of Human notes (2026-10-01)

These are verification notes for phase wrap-up, written at the Human's request. They record no decision and change no outcome. Probes ran against scratch copies of the fixtures with the same binary (`sha256:3a34fc31…5d49`). The retained fixtures were not mutated.

| ID | Related rows | Finding | Evidence | Classification for wrap-up |
|---|---|---|---|---|
| A1 | HUMAN-03b-resolution, HUMAN-03b-internals | The `review-rework-exhausted` Stop appears **unresolvable within Phase 23**, not merely under-explained. Its condition is `budget-observed`. Phase 22 `contracts.md` defines that fact as `run.budget-observed` on the same budget showing a non-exhausted total. Phase 23 `contracts.md` §6 says `run.observe-budget` rejects `review-cycles` and that no cycle is reset. The observation also needs an active Run, but `run.start` is one of the actions this Stop blocks. Raising the maximum needs Envelope revision, which is Phase 24 scope. | A scratch-copy `run observe-budget` on the only Run (`RUN-188ace`, status `failed`) returned `run-inactive` (exit 4). Phase 23 `contracts.md` lines 233-245 and 256-257; Phase 22 `contracts.md` line 220. | Contract or design gap affecting VAL-072. Needs a Tech Lead/Human disposition. |
| A2 | HUMAN-03a | Human-mode `claim acquire` is **silent** on the rework denial: exit 4 and nothing on stdout or stderr. This violates Phase 21 `public-interface.md` §1, which requires Human failure to print `error: <message>` on stderr. Other `claim acquire` failures render correctly (for example invalid input prints `error: claim acquire input is invalid` and exits 2), so the defect is specific to the Phase 23 denial path. | Scratch copy: Human mode exit 4 with empty output; `--json` returns `review-rework-exhausted` with `stop_ids`/`attention_ids`. | Product defect introduced in Phase 23. |
| A3 | HUMAN-04a/b/c-17, HUMAN-04a-18 | fid keeps two IDs per requirement. `REQ-leaf` is the project's stable `requirement_id` (schema `requirement_refs.requirement_id`). `RREF-7a2b10` is fid's ID for that version of the reference, and a new one is created on each supersede. The accepted wire shape `RequirementRef = {id,source,status,…}` (Phase 23 `public-interface.md` line 118) **has no `requirement_id` field**, so the stable ID never appears on the reference. `get` accepts either ID. The Human view labels both the reference and the link line "Requirement". `list` returns only references, so `REQ-leaf` never appears. | `requirement get REQ-leaf` and `requirement get RREF-7a2b10` return identical results. The JSON reference has `id:"RREF-7a2b10"` and no `REQ-leaf` member, while the request carried `id:"REQ-leaf"`. | Gap in the accepted interface plus a view labelling defect. |
| A4 | HUMAN-02a-before, HUMAN-03a | Stop details render as raw JSON in `stop list`, `attention list`, and the `resume` "Blockers and attention" block. | Rendered by `cmd/fid/stop_human.go:78`, introduced in `a0a93c7` (Phase 22 task 9). | Existing view, so CF-040 / Phase 25 scope. It still degraded the VAL-072 judgment. |
| A5 | HUMAN-02a-before | Human-mode denials omit the reason: `error: execution-bound transition for issue "P22-a6bf3e" is not permitted`. `--json` carries the `code`, `details.unsatisfied`, and Stop data. Phase 21 `public-interface.md` §1 limits Human failure output to the message plus "at most two non-semantic guidance lines", so showing reasons may need a contract change, not just a rendering fix. | Scratch-copy comparison of Human mode and `--json`. | Contract constraint. Candidate for Phase 25 operating-experience work. |
| A6 | HUMAN-04c-03 (applies across views) | The flagged lines are currently **required**. The `Purpose:` line is what HUMAN-04b assesses (contracts §14, VAL-073). The "recorded judgment, not deterministic proof" line comes from PRD-HUM-001 / VAL-060. The Human's request to move them to `--help` or a manual for experienced users is therefore a change to accepted guidance, not a defect in this phase. | VAL-060 and VAL-073 rows in `validation-contract.md`. | Guidance change request for Phase 25 (CF-040 neighbourhood). |
| A7 | HUMAN-04a-09 | `finding address` Human output omits the request `reason`. The fixture value `PRIVATE_FINDING_REASON` suggests the omission is deliberate, but no explicit exclusion rule was found in the Phase 23 contracts. | `leaf-leaf-address.json` compared with the replayed Human output. | Clarification question. |
| A8 | HUMAN-04a-02 (applies to all `--input` leaves) | This matches the accepted specs, not a deviation. `06_DESIGN_PLAN_CLI.md` §1 and Phase 23 `public-interface.md` §1 require `--input <file\|->` JSON for these mutations, and no Human-friendlier form is specified. `--input -` accepts stdin, so no temp file is needed, but the schema is still required. The deeper issue is that several payloads need values no Human can type: `repository` OIDs, `worktree_id`, and file `sha256` digests. This matters most for leaves that need a Human actor (`human-outcome record`, `attention resolve`, possibly `decision record-reference`), because PRD-AUTH-006/008 make Human provenance meaningful. CF-040 covers output only, so input ergonomics has no owner today. | Leaf `--help` output; `leaf-leaf-outcome.json`. Separately, `attention resolve --help` describes `--input` as "Run JSON input" (apparent copy-paste defect), and the JSON-input leaves have no one-line description in `--help`. | Unowned operating-experience gap. The Human is still considering options; no decision recorded. |
