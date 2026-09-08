# Phase 22 Human Verification Checklist

> Phase: 22 — Recovery and portability
> Prepared: 2026-09-08
> Tester: _______________
> Status: **Runnable snapshot independently verified; initial final QA requires repair. Human judgments pending.**
> Human responses and outcomes have not been recorded.

## Before starting

Use a fresh session without prior chat or an explanation of the intended
answers. Inspect the supplied views and write your own response for each row
before recording an outcome. Record anything unclear in the notes, including
the wording or view that led to your interpretation.

The retained set below was independently inspected after its generator exited.
All twelve actual Human/JSON commands terminated successfully. Binary hash,
clean build revision and four Project identities were independently verified.
This is the inspected pre-repair snapshot. Product repairs require renewed
fixture qualification before final acceptance. No Human judgment is inferred.

| Preparation record | Verified value |
|---|---|
| Fixture guide | [human-fixtures.md](human-fixtures.md) |
| Manifest | [manifest.json](/home/mark/projects/fid/.cache/phase-22-human/checkpoints-4043242815/manifest.json) |
| Binary | `/home/mark/projects/fid/.cache/phase-22-human/checkpoints-4043242815/fid` |
| SHA-256 | `sha256:da111fbdc19a200ff8fe7d263e6fdd7f25915fa3120c83b554649b876dea1d9e` |
| Clean source revision | `e04bfbdd0194393b0bfee200477a14bc6c563fb3`; source_dirty:false independently matched to go version -m |
| Generated | `2026-09-08T10:57:49.722794225Z` |
| Independent evidence | [final-review.md](final-review.md), Command And Feasibility Validation |
| Successor lease expires | `2026-09-09T10:57:54.448190739Z` / **9 September 2026, 22:57:54 NZST** |

Keep these directories at their recorded paths. If the lease expires before
this exercise or product repair changes the behavior under review, generate a
new set with the accepted command below and have QA verify its manifest and
update the exact paths before recording acceptance. Preserve this original set.

```bash
cd /home/mark/projects/fid
FID_P22_HUMAN_FIXTURE_ROOT=/home/mark/projects/fid/.cache/phase-22-human go test -tags integration ./cmd/fid/integration/... -run '^TestP22HumanFixtures$' -count=1 -v
```

## Recovery view

**Recovery checkpoint**

Project: `/home/mark/projects/fid/.cache/phase-22-human/checkpoints-4043242815/recovery`.

- `claim_id`: `CLM-de71f4`
- `issue_id`: `P22-031f58`
- `prior_run_id`: `RUN-ea4d66`
- `project_id`: `project-e2fb86770a8a08b265e48fd79e579254`
- `run_id`: `RUN-853370`

Copy these inspection commands as a block:

```bash
cd /home/mark/projects/fid/.cache/phase-22-human/checkpoints-4043242815/recovery
/home/mark/projects/fid/.cache/phase-22-human/checkpoints-4043242815/fid --project /home/mark/projects/fid/.cache/phase-22-human/checkpoints-4043242815/recovery resume P22-031f58 --repository /home/mark/projects/fid/.cache/phase-22-human/checkpoints-4043242815/recovery
/home/mark/projects/fid/.cache/phase-22-human/checkpoints-4043242815/fid --project /home/mark/projects/fid/.cache/phase-22-human/checkpoints-4043242815/recovery --json resume P22-031f58 --repository /home/mark/projects/fid/.cache/phase-22-human/checkpoints-4043242815/recovery
/home/mark/projects/fid/.cache/phase-22-human/checkpoints-4043242815/fid --project /home/mark/projects/fid/.cache/phase-22-human/checkpoints-4043242815/recovery run inspect RUN-853370 --repository /home/mark/projects/fid/.cache/phase-22-human/checkpoints-4043242815/recovery
/home/mark/projects/fid/.cache/phase-22-human/checkpoints-4043242815/fid --project /home/mark/projects/fid/.cache/phase-22-human/checkpoints-4043242815/recovery --json run inspect RUN-853370 --repository /home/mark/projects/fid/.cache/phase-22-human/checkpoints-4043242815/recovery
```


| Test ID | Covers | Prompt | Your response / evidence | Outcome | Notes |
|---|---|---|---|---|---|
| HUMAN-01a (human) | VAL-076 | What single recovery action would you take next? State what in the view led you to that choice. | | | |
| HUMAN-01b (human) | VAL-076 | What does the view tell you about the proof status of the adopted changes? Explain your interpretation. | | | |

## Activation recovery

**Activation checkpoint A**

Project: `/home/mark/projects/fid/.cache/phase-22-human/checkpoints-4043242815/prepared`.

- `activation_id`: `ACT-5669eb`
- `project_id`: `project-055ee1afea8d14b45fa6f6b7b0eac5d3`

Copy these inspection commands as a block:

```bash
cd /home/mark/projects/fid/.cache/phase-22-human/checkpoints-4043242815/prepared
/home/mark/projects/fid/.cache/phase-22-human/checkpoints-4043242815/fid --project /home/mark/projects/fid/.cache/phase-22-human/checkpoints-4043242815/prepared activation status ACT-5669eb
/home/mark/projects/fid/.cache/phase-22-human/checkpoints-4043242815/fid --project /home/mark/projects/fid/.cache/phase-22-human/checkpoints-4043242815/prepared --json activation status ACT-5669eb
```

**Activation checkpoint B**

Project: `/home/mark/projects/fid/.cache/phase-22-human/checkpoints-4043242815/committed`.

- `activation_id`: `ACT-ecdc52`
- `project_id`: `project-64fb8189b9f0e5f1ea91e257bb9e28af`

Copy these inspection commands as a block:

```bash
cd /home/mark/projects/fid/.cache/phase-22-human/checkpoints-4043242815/committed
/home/mark/projects/fid/.cache/phase-22-human/checkpoints-4043242815/fid --project /home/mark/projects/fid/.cache/phase-22-human/checkpoints-4043242815/committed activation status ACT-ecdc52
/home/mark/projects/fid/.cache/phase-22-human/checkpoints-4043242815/fid --project /home/mark/projects/fid/.cache/phase-22-human/checkpoints-4043242815/committed --json activation status ACT-ecdc52
```


| Test ID | Covers | Prompt | Your response / evidence | Outcome | Notes |
|---|---|---|---|---|---|
| HUMAN-02a (human) | VAL-077 | In the pre-commit checkpoint, which recovery direction is offered, and what in the view explains it? | | | |
| HUMAN-02b (human) | VAL-077 | In the post-commit checkpoint, is rollback possible? Explain what you understand from the view. | | | |

## Imported work

**Imported project checkpoint**

Project: `/home/mark/projects/fid/.cache/phase-22-human/checkpoints-4043242815/imported`.

- `archive_path`: `/home/mark/projects/fid/.cache/phase-22-human/checkpoints-4043242815/import.zip`
- `claim_id`: `CLM-a8e739`
- `issue_id`: `P22-5b8426`
- `project_id`: `project-50fca5be4933a95c33e7218071e5c47d`
- `run_id`: `RUN-0d02e8`
- `source_project_path`: `/home/mark/projects/fid/.cache/phase-22-human/checkpoints-4043242815/import-source`

Copy these inspection commands as a block:

```bash
cd /home/mark/projects/fid/.cache/phase-22-human/checkpoints-4043242815/imported
/home/mark/projects/fid/.cache/phase-22-human/checkpoints-4043242815/fid --project /home/mark/projects/fid/.cache/phase-22-human/checkpoints-4043242815/imported resume P22-5b8426 --repository /home/mark/projects/fid/.cache/phase-22-human/checkpoints-4043242815/imported
/home/mark/projects/fid/.cache/phase-22-human/checkpoints-4043242815/fid --project /home/mark/projects/fid/.cache/phase-22-human/checkpoints-4043242815/imported --json resume P22-5b8426 --repository /home/mark/projects/fid/.cache/phase-22-human/checkpoints-4043242815/imported
/home/mark/projects/fid/.cache/phase-22-human/checkpoints-4043242815/fid --project /home/mark/projects/fid/.cache/phase-22-human/checkpoints-4043242815/imported run inspect RUN-0d02e8 --repository /home/mark/projects/fid/.cache/phase-22-human/checkpoints-4043242815/imported
/home/mark/projects/fid/.cache/phase-22-human/checkpoints-4043242815/fid --project /home/mark/projects/fid/.cache/phase-22-human/checkpoints-4043242815/imported --json run inspect RUN-0d02e8 --repository /home/mark/projects/fid/.cache/phase-22-human/checkpoints-4043242815/imported
```


| Test ID | Covers | Prompt | Your response / evidence | Outcome | Notes |
|---|---|---|---|---|---|
| HUMAN-03a (human) | VAL-078 | Which restored Issue, Run and history are you looking at? Record the identities and what helped you recognize the work. | | | |
| HUMAN-03b (human) | VAL-078 | What intervention, if any, remains before this imported work can be mutated? State what in the views led you to that conclusion. | | | |

## Human results

All six outcomes remain blank. Agents may prepare and verify the fixtures;
only the Human records the judgments above. Deterministic QA evidence belongs
in the separate final review.

| Section | Total | Pass | Fail | Partial | Blank |
|---|---:|---:|---:|---:|---:|
| Recovery view | 2 | 0 | 0 | 0 | 2 |
| Activation recovery | 2 | 0 | 0 | 0 | 2 |
| Imported work | 2 | 0 | 0 | 0 | 2 |
| **Total** | **6** | **0** | **0** | **0** | **6** |

| Result record | Human entry |
|---|---|
| Overall result | |
| Date completed | |
| Notes | |

Outcome values are `pass`, `fail` or `partial`; leave an unperformed item
blank. Record the overall result only after completing the six responses.
