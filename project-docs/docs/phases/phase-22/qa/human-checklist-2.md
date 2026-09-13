# Phase 22 Human Verification — Attempt 2

> Status: Prepared for Human retest after task-19; outcomes are blank.
> Tester:
> Date:

Attempt 1 remains in [human-checklist.md](human-checklist.md). Judge whether the
new views themselves explain the situation. If you still need an explanation,
record that difficulty rather than relying on the discussion of the first attempt.

## Before starting

Use the retained binary and paths below. Read the Human output first and write
your answers. JSON is optional diagnostic material after that; needing it to
understand a Human view is useful feedback. All commands below are inspections;
you do not need to execute a recovery, Gate or other follow-up mutation.

| Fixture record | Value |
|---|---|
| Manifest | [manifest.json](/home/mark/projects/fid/.cache/phase-22-human/checkpoints-386504276/manifest.json) |
| Binary | `/home/mark/projects/fid/.cache/phase-22-human/checkpoints-386504276/fid` |
| SHA-256 | `sha256:e2ca68186dc96297705c519fd9849d0a041862522d8b002c96fecf3e1393fbc9` |
| Source revision | `ee7a5cae5fb7b572be1b8032c0ed17cb6a1afa37` |
| Source dirty | `false` |
| Generated | `2026-09-12T10:18:42.450198226Z` |
| Recovery lease expires (UTC) | `2026-09-13T10:18:47.253596499Z` |
| Recovery lease expires (NZST) | **13 September 2026, 10:18:47 pm** |

Keep the fixtures at these paths. If the lease expires before you test,
ask for fresh fixtures; changing recorded timestamps would invalidate the setup.

Independent review: [task-19-review.md](task-19-review.md).
Technical evidence: [task-19-done.md](../tasks/task-19-done.md).

## Recovery view

Project: `/home/mark/projects/fid/.cache/phase-22-human/checkpoints-386504276/recovery`.

`resume` presents the Issue and its related work. Use it as the primary
view for these questions. `run inspect` provides details about the selected
Run and repository comparison.

```bash
cd /home/mark/projects/fid/.cache/phase-22-human/checkpoints-386504276/recovery
/home/mark/projects/fid/.cache/phase-22-human/checkpoints-386504276/fid --project /home/mark/projects/fid/.cache/phase-22-human/checkpoints-386504276/recovery resume P22-f6fa61 --repository /home/mark/projects/fid/.cache/phase-22-human/checkpoints-386504276/recovery
```

```bash
cd /home/mark/projects/fid/.cache/phase-22-human/checkpoints-386504276/recovery
/home/mark/projects/fid/.cache/phase-22-human/checkpoints-386504276/fid --project /home/mark/projects/fid/.cache/phase-22-human/checkpoints-386504276/recovery run inspect RUN-2aac7e --repository /home/mark/projects/fid/.cache/phase-22-human/checkpoints-386504276/recovery
```

| Test | Covers | Question | Your response / evidence | Outcome | Notes |
|---|---|---|---|---|---|
| HUMAN-01a (human) | VAL-076 | What single recovery action would you take next? State what in the view led you to that choice. | Run the gate. The output shows an active run and claim on an incomplete task. The completion requirements section makes it clear that there is no evidence for OBL-p22-corpus. | PASS | CLI views could be polished, but they are much improved from before. Ok to carry-forward the polish for now. |
| HUMAN-01b (human) | VAL-076 | What does the view tell you about the proof status of the adopted changes? Explain your interpretation. |   The view tells me that the uncommitted changes to base.txt were retained for the successor Run. It explicitly says that adoption does not establish that those changes were tested, reviewed, approved, or proved correct. Qualifying Evidence is still absent. The Run inspection also says there is no liveness observation for the current Run, but that is separate from proof of the changes. | PASS | Much improved from before. |

<details>
<summary>Optional JSON diagnostics (after recording your interpretation)</summary>

```bash
cd /home/mark/projects/fid/.cache/phase-22-human/checkpoints-386504276/recovery
/home/mark/projects/fid/.cache/phase-22-human/checkpoints-386504276/fid --project /home/mark/projects/fid/.cache/phase-22-human/checkpoints-386504276/recovery --json resume P22-f6fa61 --repository /home/mark/projects/fid/.cache/phase-22-human/checkpoints-386504276/recovery
cd /home/mark/projects/fid/.cache/phase-22-human/checkpoints-386504276/recovery
/home/mark/projects/fid/.cache/phase-22-human/checkpoints-386504276/fid --project /home/mark/projects/fid/.cache/phase-22-human/checkpoints-386504276/recovery --json run inspect RUN-2aac7e --repository /home/mark/projects/fid/.cache/phase-22-human/checkpoints-386504276/recovery
```

</details>

## Activation checkpoint A

Project: `/home/mark/projects/fid/.cache/phase-22-human/checkpoints-386504276/prepared`.

`activation status` presents this activation and its resource state.

```bash
cd /home/mark/projects/fid/.cache/phase-22-human/checkpoints-386504276/prepared
/home/mark/projects/fid/.cache/phase-22-human/checkpoints-386504276/fid --project /home/mark/projects/fid/.cache/phase-22-human/checkpoints-386504276/prepared activation status ACT-e869cf
```

| Test | Covers | Question | Your response / evidence | Outcome | Notes |
|---|---|---|---|---|---|
| HUMAN-02a (human) | VAL-077 | In the pre-commit checkpoint, which recovery direction is offered, and what in the view explains it? | `restore-prior`; the view clearly states this at the top. | PASS | Much improved from before.|

<details>
<summary>Optional JSON diagnostics (after recording your interpretation)</summary>

```bash
cd /home/mark/projects/fid/.cache/phase-22-human/checkpoints-386504276/prepared
/home/mark/projects/fid/.cache/phase-22-human/checkpoints-386504276/fid --project /home/mark/projects/fid/.cache/phase-22-human/checkpoints-386504276/prepared --json activation status ACT-e869cf
```

</details>

## Activation checkpoint B

Project: `/home/mark/projects/fid/.cache/phase-22-human/checkpoints-386504276/committed`.

`activation status` presents this activation and its resource state.

```bash
cd /home/mark/projects/fid/.cache/phase-22-human/checkpoints-386504276/committed
/home/mark/projects/fid/.cache/phase-22-human/checkpoints-386504276/fid --project /home/mark/projects/fid/.cache/phase-22-human/checkpoints-386504276/committed activation status ACT-6a252d
```

| Test | Covers | Question | Your response / evidence | Outcome | Notes |
|---|---|---|---|---|---|
| HUMAN-02b (human) | VAL-077 | In the post-commit checkpoint, is rollback possible? Explain what you understand from the view. | Rollback is not possible; the view is very clear about this now. | PASS | Verbose, but much improved from before. |

<details>
<summary>Optional JSON diagnostics (after recording your interpretation)</summary>

```bash
cd /home/mark/projects/fid/.cache/phase-22-human/checkpoints-386504276/committed
/home/mark/projects/fid/.cache/phase-22-human/checkpoints-386504276/fid --project /home/mark/projects/fid/.cache/phase-22-human/checkpoints-386504276/committed --json activation status ACT-6a252d
```

</details>

## Imported project view

Project: `/home/mark/projects/fid/.cache/phase-22-human/checkpoints-386504276/imported`.

`resume` presents the Issue and its related work. Use it as the primary
view for these questions. `run inspect` provides details about the selected
Run and repository comparison.

Source project reference: `/home/mark/projects/fid/.cache/phase-22-human/checkpoints-386504276/import-source`.
Archive reference: `/home/mark/projects/fid/.cache/phase-22-human/checkpoints-386504276/import.zip`.

```bash
cd /home/mark/projects/fid/.cache/phase-22-human/checkpoints-386504276/imported
/home/mark/projects/fid/.cache/phase-22-human/checkpoints-386504276/fid --project /home/mark/projects/fid/.cache/phase-22-human/checkpoints-386504276/imported resume P22-70e7b3 --repository /home/mark/projects/fid/.cache/phase-22-human/checkpoints-386504276/imported
```

```bash
cd /home/mark/projects/fid/.cache/phase-22-human/checkpoints-386504276/imported
/home/mark/projects/fid/.cache/phase-22-human/checkpoints-386504276/fid --project /home/mark/projects/fid/.cache/phase-22-human/checkpoints-386504276/imported run inspect RUN-6fa5b2 --repository /home/mark/projects/fid/.cache/phase-22-human/checkpoints-386504276/imported
```

| Test | Covers | Question | Your response / evidence | Outcome | Notes |
|---|---|---|---|---|---|
| HUMAN-03a (human) | VAL-078 | Which restored Issue, Run and history are you looking at? Record the identities and what helped you recognize the work. | P22-70e7b3, RUN-6fa5b2; view clearly includes this now. | PASS | Very verbose, but the key info is much clearer now. |
| HUMAN-03b (human) | VAL-078 | What intervention, if any, remains before this imported work can be mutated? State what in the views led you to that conclusion. | The import preserved the work history but did not grant mutation authority in the destination. The destination repository must be inspected and reconciled, and the activated harness must record a fresh destination liveness observation for the imported Run. The repository inspection shows no changed paths or conflicting operation, but liveness is still unknown because no observation has been recorded. Recovery must then check that observation, the repository state, and the applicable authority before the work can continue.| PASS | Ok for now. |

<details>
<summary>Optional JSON diagnostics (after recording your interpretation)</summary>

```bash
cd /home/mark/projects/fid/.cache/phase-22-human/checkpoints-386504276/imported
/home/mark/projects/fid/.cache/phase-22-human/checkpoints-386504276/fid --project /home/mark/projects/fid/.cache/phase-22-human/checkpoints-386504276/imported --json resume P22-70e7b3 --repository /home/mark/projects/fid/.cache/phase-22-human/checkpoints-386504276/imported
cd /home/mark/projects/fid/.cache/phase-22-human/checkpoints-386504276/imported
/home/mark/projects/fid/.cache/phase-22-human/checkpoints-386504276/fid --project /home/mark/projects/fid/.cache/phase-22-human/checkpoints-386504276/imported --json run inspect RUN-6fa5b2 --repository /home/mark/projects/fid/.cache/phase-22-human/checkpoints-386504276/imported
```

</details>

## Human results

Use PASS, FAIL or PARTIAL for each performed check. Leave unperformed checks blank.
Explain any remaining uncertainty in your own words; agents do not fill the judgments.

| Total | Pass | Fail | Partial | Blank |
|---|---|---|---|---|
| 6 | 6 | 0 | 0 | 0 |

| Field | Human response |
|---|---|
| Overall result | PASS |
| Remaining concerns | Make a carry-forward to come back and review view polish and verbosity in the future. |
