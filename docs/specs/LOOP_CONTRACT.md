# Loop contract: fid-cli.v1

fid owns this interface; the kit owns the loop. This is the first declaration
of `fid-cli.v1`, already carried by every JSON reply. Use `--json` on every
call. Only the fields, values and process exits declared here are contract.
Other output remains available to people and tools, but a loop must not
rely on it. JSON examples and the fuller command interface live in
[the CLI spec](06_DESIGN_PLAN_CLI.md).

A loop checks `contract_version` at startup and refuses an unsupported
version. It stops loudly on an unknown `next.kind` or `waiting.kind`, malformed
JSON, missing required fields or empty output. Empty output never means no
work. Changing or removing a declared field, kind or exit code requires a
`contract_version` bump. Adding an optional field does not. Informational
queue reasons may grow without a bump; loops branch on kinds, not reasons.

## Transport and exits

Success is one JSON object on stdout: `contract_version` (string), `ok: true`,
`operation` (string) and `result` (object). Failure is one JSON object on
stderr, with the same version and operation, `ok: false`, and `error.code`
(string) and `error.retryable` (boolean). A failure has no result. Success
has no error. A loop reads the process exit as well as the envelope; success
is exit 0. A nonzero exit never means that there is no work.

| Exit | Meaning |
|---|---|
| 0 | Command succeeded, including `next` kind `none` and batch `done` with per-task refusals |
| 2 | Invalid input, including unknown flags and a non-epic scope |
| 3 | Project or issue not found, including an unknown epic scope |
| 4 | Conflict or unmet precondition; consult `error.retryable` |
| 5 | Actor lacks authority |
| 6 | Journal or projection integrity failure |
| 7 | Internal failure |
| 8 | A check ran unsuccessfully (`check-unsuccessful`) |

| Command | Operation | Exits |
|---|---|---|
| `next [--claim] [--epic <epic>]` | `next.get` / `next.claim` | 0, 2, 3, 4, 5 (claim only), 6, 7 |
| `start` | `task.start` | 0, 2, 3, 4, 5, 6, 7 |
| `heartbeat` | `hold.heartbeat` | 0, 2, 3, 4, 5, 6, 7 |
| `park --interrupted` / `park --waiting` | `hold.park` | 0, 2, 3, 4, 5, 6, 7 |
| `answer` (person only) | `task.answer` | 0, 2, 3, 4, 5, 6, 7 |
| `check [name]...` | `check.execute` | 0, 2, 3, 4, 5, 6, 7, 8 |
| `handoff` | `task.handoff` | 0, 2, 3, 4, 5, 6, 7 |
| `review record <task>` | `review.record` | 0, 2, 3, 4, 5, 6, 7 |
| `done [task]` | `task.done` | 0, 2, 3, 4, 5, 6, 7 |

Parser failures for `next` name `next.get`; once parsed, a claiming request
names `next.claim`. Failure codes identify the refusal; their category and
retryability are defined by the CLI's closed error-code registry. In
particular, a stale heartbeat is `hold-not-current` (exit 4, not retryable):
the loop stops that session. Retryable means the same request may succeed
when its external condition changes; it does not authorize blind retries.

## Work selection

`next.get` reads without taking work. `next.claim` takes or returns the
worktree's hold atomically. Both have exactly these kinds:

- `new`: build the task, beginning with `start`.
- `resume`: continue interrupted work, or the worktree's live hold.
- `rework`: fix the task after changes were requested, check and hand off
  again; it is already started.
- `review`: review the declared stream and all required profiles.
- `none`: no work can be served here. `task` and `hold` are null. This is
  the sole successful no-work signal.

`task.id` is the selected task ID. `hold.id` is the session's opaque build
or review hold ID, supplied to subsequent commands. `hold` is null without
`--claim` and for `none`. IDs must be passed back unchanged.

`waiting` is null or an object whose `kind` is one of `waiting-for-you`,
`review`, `rework-cap`. It explains why a worktree remains occupied even
when `kind` is `none`. An answer is null or
`{question_id, question, text, decision_id}`; pass the person's question
and answer text to the resumed agent. It remains on subsequent selections
until a newer question replaces it.

For `review`, the additional `review` object contains
`handoff: {id, pr, head, base}`, `stream` (`primary` or `separate`) and
`required_profiles` (an array containing `code`, `test-adequacy`, or both,
as the frozen hand-off requires). The task is still `task.id`. Review work
uses the same `hold.id` and heartbeat interface. Record all required
profiles from one session. The builder cannot review its own work.

`--epic <epic>` covers **every kind**, including review, resume and rework.
Only that epic's tasks may be served, in the usual order and with the
usual dependency, sequence, reservation and cap rules. Other tasks remain
in the explanation with typed reason `outside-scope`; the scope does not
pull an out-of-epic blocker into scope. Nothing in scope means `none`.
An unknown ID is `issue-not-found` (3); an existing non-epic is
`unsupported-kind` (2). An empty explicit scope is `invalid-argument` (2).
A scope change cannot start another task in an occupied worktree: it
returns `none` until the existing work is resolved or the matching scope
is restored. Unscoped selection is unchanged.

## Session commands and completion

- `start`: `task_id`, `hold_id`, `status` (`in_progress`).
- `heartbeat`: `task_id`, `hold_id`, `last_heartbeat_at`, `expires_at`
  (UTC RFC3339 strings) and `heartbeat_interval_seconds` (300). Renew
  build and review sessions at that interval.
- `park`: `task_id`, `hold_id`, `mode` (`interrupted` or `waiting`),
  `kept_for_worktree` (boolean). The interrupted form includes `reason`;
  the waiting form includes `question_id`, `question`. Parking ends the
  session. Started work remains kept for this worktree.
- `answer`: `task_id`, `question_id`, `decision_id`, `answer` (text).
- `check`: `task_id`, `checks[]`, each with `id`, `name`, `outcome`, `exit_code`, `counts`. A counting pass
  has outcome `passed`, exit code 0 and counts true. A failing command
  emits the failure envelope rather than a partial success. Exit 0 with
  an empty checks array means the lane requires no checks.
- `handoff`: `task_id`, `status` (`in_review`), `handoff_id`, `pr`, `head`,
  `base`, `required_checks`, `result_ids`. It freezes the evidence and
  returns the task to review; it does not merge.
- `review record`: `round` with `id`, `task_id`, `handoff_id`, `stream`,
  `profiles`, `verdict` (`APPROVE` or `CHANGES`) and `findings[]` with
  `id`, `profile`, `severity`, `summary`, `location`, `evidence`, `fix`. Pass
  findings to the fixing agent and resolved finding IDs to the next review. CHANGES returns the task
  to the builder. Rework is served at the cap, stopped past it, with
  `waiting.kind: rework-cap`; another worktree can carry on.
- `done <task>`: `status: done`, `completion.task_id`. Pending merge/CI or
  missing evidence is a failure envelope, never a successful waiting task.
- `done` without a task: `outcomes[]` with `task_id` and `status`:
  `completed`, `waiting`, `refused`, `github-unavailable`. Each failed
  outcome includes `refusal.code`; completed outcomes include
  `completion.task_id`. The batch itself exits 0 if it ran, even when
  an individual task waits or is refused; an empty array means no tasks
  needed reconciliation. The loop must inspect every outcome.

## Executable declaration

The integration goldens select only these result paths. `?` denotes an
optional key; a null parent stays null; `[]` selects each array element.
The version, ok, operation, error.code, error.retryable and process exit
are also pinned. IDs, commits and times are checked by shape and replaced
with placeholders in goldens; journeys assert task identity and exclusivity.

```json loop-contract
{
  "contract_version": "fid-cli.v1",
  "next_kinds": ["new", "resume", "rework", "review", "none"],
  "waiting_kinds": ["waiting-for-you", "review", "rework-cap"],
  "fields": {
    "next.get": ["kind", "task.id", "hold.id", "waiting.kind", "answer.question_id", "answer.question", "answer.text", "answer.decision_id", "review?.handoff.id", "review?.handoff.pr", "review?.handoff.head", "review?.handoff.base", "review?.stream", "review?.required_profiles"],
    "next.claim": ["kind", "task.id", "hold.id", "waiting.kind", "answer.question_id", "answer.question", "answer.text", "answer.decision_id", "review?.handoff.id", "review?.handoff.pr", "review?.handoff.head", "review?.handoff.base", "review?.stream", "review?.required_profiles"],
    "task.start": ["task_id", "hold_id", "status"],
    "hold.heartbeat": ["task_id", "hold_id", "last_heartbeat_at", "expires_at", "heartbeat_interval_seconds"],
    "hold.park": ["task_id", "hold_id", "mode", "kept_for_worktree", "reason?", "question_id?", "question?"],
    "task.answer": ["task_id", "question_id", "decision_id", "answer"],
    "check.execute": ["task_id", "checks[].id", "checks[].name", "checks[].outcome", "checks[].exit_code", "checks[].counts"],
    "task.handoff": ["task_id", "status", "handoff_id", "pr", "head", "base", "required_checks", "result_ids"],
    "review.record": ["round.id", "round.task_id", "round.handoff_id", "round.stream", "round.profiles", "round.verdict", "round.findings[].id", "round.findings[].profile", "round.findings[].severity", "round.findings[].summary", "round.findings[].location", "round.findings[].evidence", "round.findings[].fix"],
    "task.done": ["status?", "completion?.task_id", "outcomes?[].task_id", "outcomes?[].status", "outcomes?[].refusal?.code", "outcomes?[].completion?.task_id"]
  }
}
```

RESET 2.4 supplies this declaration and scenario evidence. V9's verifier,
`check:loop-selftest`, lands with the actual loop in RESET 2.5. It tests
startup version refusal and loud failure on an unknown kind using stub agents.
