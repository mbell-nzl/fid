# Phase 24 - Human walkthrough after Task 11 (Window D2)

Human: Mark Bell · Recorder: Mark Bell (self) · Date: 2026-10-10
Manifest: .cache/phase-24-human/fixtures-2268063598/manifest.json
Binary: sha256:61fa7a4c17e05de5e662caf1147349460cb93c2a1ca482cdb6e0c23ba29f55b7
Generated at: Task 11 close 93a1f792d27f001b4052113e9056faf97706c3e5
Answers: qa/review-9.md B1-H
Method: Human-mode output only; no JSON, source or specs consulted.

## VAL-054 - refusals

| Row | Fixture | Why refused | Next action | Clear from output alone? |
|---|---|---|---|---|
| HUMAN-02a | idempotency-conflict | The idempotency key conflicts with a different request. | Resend the original input or use a new key. | Yes, and I can see which fields differ from the previous request with the same key: expected_version, observed_at, previous_sample_id, sample_kind, value. |
| HUMAN-02b | repository-mismatch | The repository identity does not match.| Check head and uncommitted changes, then correct expected_repository. | Yes. |
| HUMAN-02c | invalid-argument (input.expected_repository.dirty) | The `input.expected_repository.dirty` value in `input` is invalid. | Correct the `input` so the `input.expected_repository.dirty` value is `true` instead of `yes`. | Yes |
| HUMAN-02d | protected-decision-required | protected Decision requires an approved reference and Human Outcome | record the security Decision Reference and an approving Human Outcome: fid decision record-reference, fid human-outcome record | yes |
| HUMAN-02e | source-digest-mismatch | Source digest does not match the observed project file | correct or omit digest for Task 9 decision | yes |

## VAL-055 - views (all four clean/dirty × same/different states)

| Row | View | Accept/Reject | Reason |
|---|---|---|---|
| HUMAN-03a | stop list | accept | Clearer than before; includes the reason. Each Stop shows labelled subject, code, blocked actions, reason, resolution condition and status, plus a Next line naming the commands that clear it. The same in all four states. |
| HUMAN-03b | attention list | accept | Clearer than before; includes the reason. Each Attention shows status, code and its linked Stops, and each linked Stop shows the same labelled details and Next line as `stop list`. The same in all four states. |
| HUMAN-03c | resume (short HEAD/state, worktree comparison) | accept | Clearer than before; includes the reason. The Repository section shows current and recorded repository as `HEAD <12 hex>; clean` or `dirty`. `Worktree comparison` says `same worktree` for the same root and `different worktree` for the plain copy. |
| HUMAN-03d | review get (repository lines) | accept | Clearer than before; includes the reason. Each Review shows its repository as `HEAD <12 hex>; clean` or `dirty`, matching the fixture state. |

Interim answers only; they don't replace phase-end HUMAN-02/HUMAN-03.

Edit note (2026-10-10): at the Human's instruction, the Tech Lead
corrected HUMAN-02a's wording and HUMAN-02c's boolean spelling (`True`
to `true`). For HUMAN-03a..03d it added what each accepted view shows, as
read from the manifest's Human output. Every accept judgment and the
Human's original reasons are unchanged. The original answers are in
commit a23b9f3325e5981ccc055fa94e0a36a0699b51d1.
