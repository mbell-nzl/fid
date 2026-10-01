# Phase 23 Human Verification Checklist

> Phase: 23 — Review and Human judgment
> Date: 2026-10-01
> Tester: _______________

All 63 outcomes are blank pending Human input. Only `(human)` judgments appear here. Deterministic results and feasibility evidence are in [final-review.md](final-review.md). Use the retained final-binary fixtures below. Record your answer and a short output quotation where requested.

## Fixture identity

- Manifest: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-526363314/manifest.json`
- Binary: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-526363314/fid`
- SHA-256: `sha256:3a34fc318991f77b9aba5136e5e6a04538824010489dbae31fd854bae01f5d49`
- Product source: `e0893cdee915013371455198906852081633e051`. The recorded dirty state consists of workflow documents and these QA artifacts; product files were unchanged.
- This binary is byte-identical to `/tmp/fid` built by the full phase gate. Keep it with these workspaces.
- These fixtures contain short-lived Claims. Generate a fresh set immediately before the Human run using the guide linked below, then use that manifest's current paths and IDs. Exact expected output is bound to generation time.
- Run each section from its stated directory. The protected-Decision steps are sequential and consume that fixture; record the before observation first. For a repeat, use the fresh-generation command in [human-fixtures.md](human-fixtures.md#generate-a-fresh-set) and the resulting manifest paths.

## Mixed Findings — VAL-070

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-526363314/mixed-findings`

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-526363314/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-526363314/mixed-findings review get REV-bab716
```

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-526363314/mixed-findings`

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-526363314/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-526363314/mixed-findings review list --issue P22-364fce
```

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-526363314/mixed-findings`

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-526363314/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-526363314/mixed-findings finding list --issue P22-364fce
```

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-526363314/mixed-findings`

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-526363314/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-526363314/mixed-findings resume P22-364fce --repository /home/mark/projects/fid/.cache/phase-23-human/checkpoints-526363314/mixed-findings
```

| Test ID | Covers | Test | Expected result | Outcome | Notes |
|---|---|---|---|---|---|
| HUMAN-01a | VAL-070 | Identify the resolved Findings from the displayed output. `(human)` | Name them and quote what establishes resolution. | | |
| HUMAN-01b | VAL-070 | Identify the Findings that are addressed but remain unresolved. `(human)` | Name them and quote the distinction. | | |
| HUMAN-01c | VAL-070 | Explain what still blocks completion. `(human)` | Give the blocking Findings and the reason visible in fid output. | | |

## Protected Decision — VAL-071

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-526363314/protected-decision`

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-526363314/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-526363314/protected-decision resume P22-f2c375 --repository /home/mark/projects/fid/.cache/phase-23-human/checkpoints-526363314/protected-decision
```

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-526363314/protected-decision`

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-526363314/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-526363314/protected-decision decision list --issue P22-f2c375
```

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-526363314/protected-decision`

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-526363314/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-526363314/protected-decision stop list --subject envelope:ENV-dbf029
```

### Step 1: blocked action before Outcome

Expected command exit: `5`.

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-526363314/protected-decision`

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-526363314/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-526363314/protected-decision --json issue transition P22-f2c375 --to qa_pass --actor agent:actor-recovery --expected-revision 4 --idempotency-key fixture-human-transition
```

| Test ID | Covers | Test | Expected result | Outcome | Notes |
|---|---|---|---|---|---|
| HUMAN-02a-before | VAL-071 | Observe the guarded action before recording the Outcome. `(human)` | The output makes the refusal clear. | | |

### Step 2: record authorized Human Outcome

Expected command exit: `0`.

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-526363314/protected-decision`

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-526363314/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-526363314/protected-decision --json human-outcome record --input /home/mark/projects/fid/.cache/phase-23-human/checkpoints-526363314/protected-decision/fixture-human-outcome.json --actor actor-human-outcome --idempotency-key fixture-human-outcome
```

### Step 3: retry guarded action

Expected command exit: `0`.

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-526363314/protected-decision`

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-526363314/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-526363314/protected-decision --json issue transition P22-f2c375 --to qa_pass --actor agent:actor-recovery --expected-revision 4 --idempotency-key fixture-human-transition
```

| Test ID | Covers | Test | Expected result | Outcome | Notes |
|---|---|---|---|---|---|
| HUMAN-02a-after | VAL-071 | Observe the same guarded action after recording the Outcome. `(human)` | The output makes the permitted action clear. | | |
| HUMAN-02b | VAL-071 | Explain why the action became permitted, using only fid output. `(human)` | State the change in authority or qualifying fact shown by the output. | | |

## Review rework exhaustion — VAL-072

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-526363314/review-exhaustion`

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-526363314/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-526363314/review-exhaustion --json claim acquire --input /home/mark/projects/fid/.cache/phase-23-human/checkpoints-526363314/review-exhaustion/fixture-rework-claim.json --actor actor-recovery --idempotency-key fixture-exhausted-claim
```

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-526363314/review-exhaustion`

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-526363314/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-526363314/review-exhaustion resume P22-30b087 --repository /home/mark/projects/fid/.cache/phase-23-human/checkpoints-526363314/review-exhaustion
```

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-526363314/review-exhaustion`

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-526363314/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-526363314/review-exhaustion stop list --subject envelope:ENV-99f2cc
```

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-526363314/review-exhaustion`

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-526363314/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-526363314/review-exhaustion attention list --subject envelope:ENV-99f2cc
```

| Test ID | Covers | Test | Expected result | Outcome | Notes |
|---|---|---|---|---|---|
| HUMAN-03a | VAL-072 | Explain the cause of the rework block. `(human)` | The cause is clear from the denial and displayed Stop/Attention. | | |
| HUMAN-03b-resolution | VAL-072 | State the exact qualifying action that would resolve the block. `(human)` | The required resolution is clear without additional help. | | |
| HUMAN-03b-internals | VAL-072 | Assess whether the output suggests editing fid internals. `(human)` | It guides you through supported actions without suggesting internal edits. | | |

## Every new view — VAL-073

The following 18 operation views are the scope of contracts §14. Each has three separate judgments. Mutation examples replay existing fixture requests. The JSON counterpart is available in the manifest if needed; judge the displayed Human text.

### requirement.record-reference

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-526363314/leaf-contract`

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-526363314/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-526363314/leaf-contract requirement record-reference --input /home/mark/projects/fid/.cache/phase-23-human/checkpoints-526363314/leaf-contract/leaf-leaf-requirement.json --actor actor-recovery --idempotency-key leaf-requirement
```

| Test ID | Covers | Test | Expected result | Outcome | Notes |
|---|---|---|---|---|---|
| HUMAN-04a-01 | VAL-073 | Assess the section labels. `(human)` | Labels make the displayed sections clear. | | |
| HUMAN-04b-01 | VAL-073 | Assess the operation purpose line. `(human)` | The view clearly states its purpose. | | |
| HUMAN-04c-01 | VAL-073 | Assess routine-state narration. `(human)` | The view is free of unnecessary routine-state narration. | | |

### requirement.link

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-526363314/leaf-contract`

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-526363314/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-526363314/leaf-contract requirement link --input /home/mark/projects/fid/.cache/phase-23-human/checkpoints-526363314/leaf-contract/leaf-leaf-link.json --actor actor-recovery --idempotency-key leaf-link
```

| Test ID | Covers | Test | Expected result | Outcome | Notes |
|---|---|---|---|---|---|
| HUMAN-04a-02 | VAL-073 | Assess the section labels. `(human)` | Labels make the displayed sections clear. | | |
| HUMAN-04b-02 | VAL-073 | Assess the operation purpose line. `(human)` | The view clearly states its purpose. | | |
| HUMAN-04c-02 | VAL-073 | Assess routine-state narration. `(human)` | The view is free of unnecessary routine-state narration. | | |

### concern.record

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-526363314/leaf-contract`

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-526363314/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-526363314/leaf-contract concern record --input /home/mark/projects/fid/.cache/phase-23-human/checkpoints-526363314/leaf-contract/leaf-leaf-concern.json --actor actor-recovery --idempotency-key leaf-concern
```

| Test ID | Covers | Test | Expected result | Outcome | Notes |
|---|---|---|---|---|---|
| HUMAN-04a-03 | VAL-073 | Assess the section labels. `(human)` | Labels make the displayed sections clear. | | |
| HUMAN-04b-03 | VAL-073 | Assess the operation purpose line. `(human)` | The view clearly states its purpose. | | |
| HUMAN-04c-03 | VAL-073 | Assess routine-state narration. `(human)` | The view is free of unnecessary routine-state narration. | | |

### concern.dispose

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-526363314/leaf-contract`

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-526363314/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-526363314/leaf-contract concern dispose --input /home/mark/projects/fid/.cache/phase-23-human/checkpoints-526363314/leaf-contract/leaf-leaf-dispose.json --actor actor-recovery --idempotency-key leaf-dispose
```

| Test ID | Covers | Test | Expected result | Outcome | Notes |
|---|---|---|---|---|---|
| HUMAN-04a-04 | VAL-073 | Assess the section labels. `(human)` | Labels make the displayed sections clear. | | |
| HUMAN-04b-04 | VAL-073 | Assess the operation purpose line. `(human)` | The view clearly states its purpose. | | |
| HUMAN-04c-04 | VAL-073 | Assess routine-state narration. `(human)` | The view is free of unnecessary routine-state narration. | | |

### deviation.record

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-526363314/leaf-contract`

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-526363314/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-526363314/leaf-contract deviation record --input /home/mark/projects/fid/.cache/phase-23-human/checkpoints-526363314/leaf-contract/leaf-leaf-deviation.json --actor actor-recovery --idempotency-key leaf-deviation
```

| Test ID | Covers | Test | Expected result | Outcome | Notes |
|---|---|---|---|---|---|
| HUMAN-04a-05 | VAL-073 | Assess the section labels. `(human)` | Labels make the displayed sections clear. | | |
| HUMAN-04b-05 | VAL-073 | Assess the operation purpose line. `(human)` | The view clearly states its purpose. | | |
| HUMAN-04c-05 | VAL-073 | Assess routine-state narration. `(human)` | The view is free of unnecessary routine-state narration. | | |

### decision.record-reference

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-526363314/leaf-contract`

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-526363314/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-526363314/leaf-contract decision record-reference --input /home/mark/projects/fid/.cache/phase-23-human/checkpoints-526363314/leaf-contract/leaf-leaf-decision.json --actor actor-recovery --idempotency-key leaf-decision
```

| Test ID | Covers | Test | Expected result | Outcome | Notes |
|---|---|---|---|---|---|
| HUMAN-04a-06 | VAL-073 | Assess the section labels. `(human)` | Labels make the displayed sections clear. | | |
| HUMAN-04b-06 | VAL-073 | Assess the operation purpose line. `(human)` | The view clearly states its purpose. | | |
| HUMAN-04c-06 | VAL-073 | Assess routine-state narration. `(human)` | The view is free of unnecessary routine-state narration. | | |

### human-outcome.record

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-526363314/leaf-contract`

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-526363314/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-526363314/leaf-contract human-outcome record --input /home/mark/projects/fid/.cache/phase-23-human/checkpoints-526363314/leaf-contract/leaf-leaf-outcome.json --actor actor-recovery --idempotency-key leaf-outcome
```

| Test ID | Covers | Test | Expected result | Outcome | Notes |
|---|---|---|---|---|---|
| HUMAN-04a-07 | VAL-073 | Assess the section labels. `(human)` | Labels make the displayed sections clear. | | |
| HUMAN-04b-07 | VAL-073 | Assess the operation purpose line. `(human)` | The view clearly states its purpose. | | |
| HUMAN-04c-07 | VAL-073 | Assess routine-state narration. `(human)` | The view is free of unnecessary routine-state narration. | | |

### review.record

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-526363314/leaf-contract`

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-526363314/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-526363314/leaf-contract review record --input /home/mark/projects/fid/.cache/phase-23-human/checkpoints-526363314/leaf-contract/leaf-leaf-review.json --actor actor-denied --idempotency-key leaf-review
```

| Test ID | Covers | Test | Expected result | Outcome | Notes |
|---|---|---|---|---|---|
| HUMAN-04a-08 | VAL-073 | Assess the section labels. `(human)` | Labels make the displayed sections clear. | | |
| HUMAN-04b-08 | VAL-073 | Assess the operation purpose line. `(human)` | The view clearly states its purpose. | | |
| HUMAN-04c-08 | VAL-073 | Assess routine-state narration. `(human)` | The view is free of unnecessary routine-state narration. | | |

### finding.address

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-526363314/leaf-contract`

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-526363314/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-526363314/leaf-contract finding address --input /home/mark/projects/fid/.cache/phase-23-human/checkpoints-526363314/leaf-contract/leaf-leaf-address.json --actor actor-recovery --idempotency-key leaf-address
```

| Test ID | Covers | Test | Expected result | Outcome | Notes |
|---|---|---|---|---|---|
| HUMAN-04a-09 | VAL-073 | Assess the section labels. `(human)` | Labels make the displayed sections clear. | | |
| HUMAN-04b-09 | VAL-073 | Assess the operation purpose line. `(human)` | The view clearly states its purpose. | | |
| HUMAN-04c-09 | VAL-073 | Assess routine-state narration. `(human)` | The view is free of unnecessary routine-state narration. | | |

### review.get

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-526363314/leaf-contract`

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-526363314/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-526363314/leaf-contract review get REV-4f9c55
```

| Test ID | Covers | Test | Expected result | Outcome | Notes |
|---|---|---|---|---|---|
| HUMAN-04a-10 | VAL-073 | Assess the section labels. `(human)` | Labels make the displayed sections clear. | | |
| HUMAN-04b-10 | VAL-073 | Assess the operation purpose line. `(human)` | The view clearly states its purpose. | | |
| HUMAN-04c-10 | VAL-073 | Assess routine-state narration. `(human)` | The view is free of unnecessary routine-state narration. | | |

### review.list

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-526363314/leaf-contract`

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-526363314/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-526363314/leaf-contract review list --issue P22-2df764
```

| Test ID | Covers | Test | Expected result | Outcome | Notes |
|---|---|---|---|---|---|
| HUMAN-04a-11 | VAL-073 | Assess the section labels. `(human)` | Labels make the displayed sections clear. | | |
| HUMAN-04b-11 | VAL-073 | Assess the operation purpose line. `(human)` | The view clearly states its purpose. | | |
| HUMAN-04c-11 | VAL-073 | Assess routine-state narration. `(human)` | The view is free of unnecessary routine-state narration. | | |

### finding.list

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-526363314/leaf-contract`

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-526363314/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-526363314/leaf-contract finding list --issue P22-2df764
```

| Test ID | Covers | Test | Expected result | Outcome | Notes |
|---|---|---|---|---|---|
| HUMAN-04a-12 | VAL-073 | Assess the section labels. `(human)` | Labels make the displayed sections clear. | | |
| HUMAN-04b-12 | VAL-073 | Assess the operation purpose line. `(human)` | The view clearly states its purpose. | | |
| HUMAN-04c-12 | VAL-073 | Assess routine-state narration. `(human)` | The view is free of unnecessary routine-state narration. | | |

### concern.list

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-526363314/leaf-contract`

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-526363314/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-526363314/leaf-contract concern list --subject issue:P22-2df764
```

| Test ID | Covers | Test | Expected result | Outcome | Notes |
|---|---|---|---|---|---|
| HUMAN-04a-13 | VAL-073 | Assess the section labels. `(human)` | Labels make the displayed sections clear. | | |
| HUMAN-04b-13 | VAL-073 | Assess the operation purpose line. `(human)` | The view clearly states its purpose. | | |
| HUMAN-04c-13 | VAL-073 | Assess routine-state narration. `(human)` | The view is free of unnecessary routine-state narration. | | |

### deviation.list

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-526363314/leaf-contract`

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-526363314/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-526363314/leaf-contract deviation list --subject issue:P22-2df764
```

| Test ID | Covers | Test | Expected result | Outcome | Notes |
|---|---|---|---|---|---|
| HUMAN-04a-14 | VAL-073 | Assess the section labels. `(human)` | Labels make the displayed sections clear. | | |
| HUMAN-04b-14 | VAL-073 | Assess the operation purpose line. `(human)` | The view clearly states its purpose. | | |
| HUMAN-04c-14 | VAL-073 | Assess routine-state narration. `(human)` | The view is free of unnecessary routine-state narration. | | |

### decision.list

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-526363314/leaf-contract`

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-526363314/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-526363314/leaf-contract decision list --issue P22-2df764
```

| Test ID | Covers | Test | Expected result | Outcome | Notes |
|---|---|---|---|---|---|
| HUMAN-04a-15 | VAL-073 | Assess the section labels. `(human)` | Labels make the displayed sections clear. | | |
| HUMAN-04b-15 | VAL-073 | Assess the operation purpose line. `(human)` | The view clearly states its purpose. | | |
| HUMAN-04c-15 | VAL-073 | Assess routine-state narration. `(human)` | The view is free of unnecessary routine-state narration. | | |

### human-outcome.list

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-526363314/leaf-contract`

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-526363314/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-526363314/leaf-contract human-outcome list --issue P22-2df764
```

| Test ID | Covers | Test | Expected result | Outcome | Notes |
|---|---|---|---|---|---|
| HUMAN-04a-16 | VAL-073 | Assess the section labels. `(human)` | Labels make the displayed sections clear. | | |
| HUMAN-04b-16 | VAL-073 | Assess the operation purpose line. `(human)` | The view clearly states its purpose. | | |
| HUMAN-04c-16 | VAL-073 | Assess routine-state narration. `(human)` | The view is free of unnecessary routine-state narration. | | |

### requirement.get

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-526363314/leaf-contract`

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-526363314/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-526363314/leaf-contract requirement get REQ-leaf
```

| Test ID | Covers | Test | Expected result | Outcome | Notes |
|---|---|---|---|---|---|
| HUMAN-04a-17 | VAL-073 | Assess the section labels. `(human)` | Labels make the displayed sections clear. | | |
| HUMAN-04b-17 | VAL-073 | Assess the operation purpose line. `(human)` | The view clearly states its purpose. | | |
| HUMAN-04c-17 | VAL-073 | Assess routine-state narration. `(human)` | The view is free of unnecessary routine-state narration. | | |

### requirement.list

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-526363314/leaf-contract`

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-526363314/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-526363314/leaf-contract requirement list
```

| Test ID | Covers | Test | Expected result | Outcome | Notes |
|---|---|---|---|---|---|
| HUMAN-04a-18 | VAL-073 | Assess the section labels. `(human)` | Labels make the displayed sections clear. | | |
| HUMAN-04b-18 | VAL-073 | Assess the operation purpose line. `(human)` | The view clearly states its purpose. | | |
| HUMAN-04c-18 | VAL-073 | Assess routine-state narration. `(human)` | The view is free of unnecessary routine-state narration. | | |

## Summary

| Section | Total | Pass | Fail | Partial | Blank |
|---|---:|---:|---:|---:|---:|
| Mixed Findings | 3 | 0 | 0 | 0 | 3 |
| Protected Decision | 3 | 0 | 0 | 0 | 3 |
| Rework exhaustion | 3 | 0 | 0 | 0 | 3 |
| Eighteen new views | 54 | 0 | 0 | 0 | 54 |
| **Total** | **63** | **0** | **0** | **0** | **63** |

| | |
|---|---|
| **Overall result** | Pending Human input |
| **Date completed** | |
| **Notes** | |

Outcome values: pass, fail, partial, blank. Overall result is pass only when every row passes; fail if any fails; otherwise partial while any judgment is partial or blank.
