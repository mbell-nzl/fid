# Phase 22 Human fixture guide

These disposable checkpoints are generated from the built fid binary. Use the
independently verified manifest and commands recorded in `human-checklist.md`.
Human responses and outcomes belong in that checklist and remain blank until
the Human records them. This guide supplies setup and inspection instructions.

## Generate a fresh set

From `/home/mark/projects/fid`, run:

```bash
FID_P22_HUMAN_FIXTURE_ROOT=/home/mark/projects/fid/.cache/phase-22-human go test -tags integration ./cmd/fid/integration/... -run '^TestP22HumanFixtures$' -count=1 -v
```

The successful output prints the absolute `manifest.json` path beneath a fresh
`checkpoints-*` directory. Each invocation creates a separate set; it does not
reuse an earlier fixture or personal project. Preserve the entire directory,
including the retained `fid` executable, source repository, destination
repository, archive and manifest. The ordinary tagged suite exercises the same
generator and post-exit inspection with temporary parent-owned storage.

## Inspect the generated checkpoints

The manifest records the binary path and SHA-256, build source commit and dirty
state, generation time, exact Project/Issue/Run/Claim/Activation identities,
working directories, and argument arrays for Human and JSON inspection.

For each `inspections` entry, set the working directory to its `cwd` and invoke
the manifest's absolute `binary` with the listed `args`, preserving each array
element as one shell argument. The independent QA checklist supplies these as
copyable commands for the delivered set, so no reconstruction from chat is
needed.

Inspect the four checkpoints in checklist order: Recovery checkpoint,
Activation checkpoint A, Activation checkpoint B, and Imported project
checkpoint. Open the specified views and record your interpretation before
consulting deterministic test expectations or later explanations. Inspection
commands do not perform the requested recovery judgments or follow-up mutations.

## Fixture lifetime and verification

The generator creates its projects at their final paths, retains the executable,
exits, and then a separate process invokes every listed inspection command.
Moving a fixture directory changes its repository context; use the recorded
absolute paths. The manifest's binary digest lets QA identify the exact tested
executable independently of a later `/tmp/fid` build.

The recovery checkpoint records `lease_expires_at` and uses the accepted maximum
86,400-second successor lease. If that lease expires before review, regenerate
a fresh set and have QA update the checklist paths. Expiry is a current view of
time, not a reason to edit stored facts or assume executor termination.

Only the Human judges VAL-076, VAL-077 and VAL-078. Generator success proves
fixture readiness and successful inspection; it does not answer those prompts.
