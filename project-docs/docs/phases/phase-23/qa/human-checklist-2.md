# Phase 23 Human Verification Checklist — Re-check 2

> Phase: 23 — Review and Human judgment
> Date: 2026-10-03
> Tester: _______________

**Pending Human input: seven blank judgments.** This is the DEC-269 re-check
scope. The [original Human checklist](human-checklist.md) remains immutable,
including its 55 pass, 2 fail and 6 partial results. Deterministic evidence
belongs in [final-review-1.md](final-review-1.md); it does not fill these cells.
That review requires task-16 repair of the full-checkptr gate before phase
acceptance. These fixtures still document the actual reviewed HEAD.

## Fixture identity

- Reviewed source HEAD: `32a55a94efb2a89dbf50b29accb434b9c753418c`.
- Generated at: `2026-10-02T17:17:10.215903366Z`.
- Checkpoint root: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561`.
- Manifest: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/manifest.json`.
- Binary: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/fid`.
- Binary SHA-256: `sha256:9a94cc6405b2c64fa5a6e2eef027b9ab39e81461b1f179b4910e53e0ffec29da`.
- Source dirty flag: `true`; only workflow/QA documents are dirty. Product source is the reviewed HEAD.
- Durable manifest: [human-fixture-manifest.json](final-evidence-1/human-fixture-manifest.json).

Use the [fresh-generation instructions](human-fixtures.md#generate-a-fresh-set)
if the one-hour Claims have expired. After regeneration, substitute the new
manifest's binary, working directories, IDs and input paths.

## Review rework exhaustion — VAL-072

First run the Human denial and inspect resume, Stop and Attention. Record what
the displayed output tells you before reading or executing the resolution step.

### Human rework denial

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/review-exhaustion`. Expected exit: `4`.

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/review-exhaustion claim acquire --input /home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/review-exhaustion/fixture-rework-claim.json --actor actor-recovery --idempotency-key fixture-exhausted-claim
```

### Inspect resume

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/review-exhaustion`. Expected exit: `0`.

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/review-exhaustion resume P22-add6ba --repository /home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/review-exhaustion
```

### Inspect Stop

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/review-exhaustion`. Expected exit: `0`.

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/review-exhaustion stop list --subject envelope:ENV-f39156
```

### Inspect Attention

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/review-exhaustion`. Expected exit: `0`.

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/review-exhaustion attention list --subject envelope:ENV-f39156
```

| Test ID | Covers | Test | Expected Result | Outcome | Notes |
|---|---|---|---|---|---|
| HUMAN-03a | VAL-072 | Explain the cause of the rework block. `(human)` | The cause is clear from the denial and displayed Stop/Attention. | | |
| HUMAN-03b-resolution | VAL-072 | State the exact qualifying action that would resolve the block. `(human)` | The required resolution is clear without additional help. | | |
| HUMAN-03b-internals | VAL-072 | Assess whether the output suggests editing fid internals. `(human)` | It guides you through supported actions without suggesting internal edits. | | |

### Task 21 resolution and admitted retry

Follow these steps in order after recording the initial observations. The approved Human Outcome resolves the matching Stop and Attention atomically; the retry then acquires a Claim.

#### Step 1: exhausted rework denial

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/review-exhaustion`. Expected exit: `4`.

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/review-exhaustion --json claim acquire --input /home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/review-exhaustion/fixture-rework-claim.json --actor actor-recovery --idempotency-key fixture-exhausted-claim
```

#### Step 2: record approved review-rework-extension Outcome

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/review-exhaustion`. Expected exit: `0`.

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/review-exhaustion --json human-outcome record --input /home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/review-exhaustion/fixture-rework-outcome.json --actor actor-human-outcome --idempotency-key fixture-rework-outcome
```

#### Step 3: retry rework Claim

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/review-exhaustion`. Expected exit: `0`.

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/review-exhaustion --json claim acquire --input /home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/review-exhaustion/fixture-rework-claim.json --actor actor-recovery --idempotency-key fixture-exhausted-claim
```


## Requirement get — VAL-073

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/leaf-contract`. Expected exit: `0`.

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/leaf-contract requirement get REQ-leaf
```

| Test ID | Covers | Test | Expected Result | Outcome | Notes |
|---|---|---|---|---|---|
| HUMAN-04a-17 | VAL-073 | Assess the section labels. `(human)` | Labels make the displayed sections clear. | | |
| HUMAN-04b-17 | VAL-073 | Assess the operation purpose line. `(human)` | The view clearly states its purpose. | | |
| HUMAN-04c-17 | VAL-073 | Assess routine-state narration. `(human)` | The view is free of unnecessary routine-state narration. | | |

## Requirement list — VAL-073

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/leaf-contract`. Expected exit: `0`.

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/leaf-contract requirement list
```

| Test ID | Covers | Test | Expected Result | Outcome | Notes |
|---|---|---|---|---|---|
| HUMAN-04a-18 | VAL-073 | Assess the section labels. `(human)` | Labels make the displayed sections clear. | | |

## Summary

| Section | Total | Pass | Fail | Partial | Blank |
|---|---:|---:|---:|---:|---:|
| Review rework exhaustion | 3 | 0 | 0 | 0 | 3 |
| Requirement get | 3 | 0 | 0 | 0 | 3 |
| Requirement list | 1 | 0 | 0 | 0 | 1 |
| **Total** | **7** | **0** | **0** | **0** | **7** |

| | |
|---|---|
| **Overall result** | Pending Human input |
| **Date completed** | |
| **Notes** | ESC-23-012 remains open. A4–A8 and remaining original PARTIAL items stay on the Tech Lead wrap-up agenda; this checklist does not accept a carry-forward. |

Outcome values: pass, fail, partial, blank. Only the Human records judgments.
