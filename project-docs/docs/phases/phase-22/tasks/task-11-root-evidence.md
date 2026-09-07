# Task 11 — Root delegate contribution

Root owns the independent activation corpus contribution under the Human's
delegated execution authority. This evidence is attributed to the root session,
not retrospectively attested as Developer execution. Working directory for
every command below: `/home/mark/projects/fid`.

The 26 activation cases run through the existing shared fake/fid adapter
interface and request/replay machinery. The fake executes an independent
in-memory transition model and does not call fid, workflow APIs or SQLite.
The retained 37 recovery cases are unchanged. Both activation starting states
use common request inputs; expected response fixtures do not manufacture fake
results. Assertions check closed response/error shapes, versions, frozen
bindings, direction, resource agreement, Reasons, Event counts, rejected
operation no-effects and immutable same-key replies after later completion.
Explanatory prose is not compared between adapters. These cases qualify the
enumerated operations and states, not every possible activation transition.

## Captured attempts

- Compile-only construction check:
  `go test -tags integration ./cmd/fid/integration/... -run '^$'`
  printed `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	0.716s [no tests to run]`
  (terminated `exit_code: 0`). This is not behavioral evidence.
- Initial fake selector:
  `go test -tags integration ./cmd/fid/integration/... -run '^TestP22DualAdapterActivationConformance$/.*fake' -count=1`
  printed `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	0.444s [no tests to run]`
  (terminated `exit_code: 0`). The slash-aware selector was incorrect; no
  behavioral PASS is claimed for this invocation.
- Corrected fake selector:
  `go test -tags integration ./cmd/fid/integration/... -run '^TestP22DualAdapterActivationConformance$/(prepared|committed)/fake$' -count=1`
  printed `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	0.457s`
  (terminated `exit_code: 0`), executing all 26 fake cases.
- Full pair command:
  `go test -tags integration ./cmd/fid/integration/... -run '^TestP22DualAdapterActivationConformance$' -count=1`
  attempt 1 printed `FAIL	github.com/mbell-nzl/fid/cmd/fid/integration	1.800s`
  (terminated `exit_code: 1`). Premature restore completion returned the wrong
  `activation-invalid` variant with unclosed `expected_statuses` and no subject.
  The Developer repaired the typed prior-resource agreement error under DEC-199.
- After strengthening exact Reason and version-detail assertions, the
  byte-identical full pair command again printed
  `FAIL	github.com/mbell-nzl/fid/cmd/fid/integration	1.800s`
  (terminated `exit_code: 1`). The implementation invented
  `resource-not-observed-old` / `activation-resource` Reason members.
  DEC-200 explicitly allocated the previously omitted activation/Stop subject
  domain, retaining the accepted Reason structure and `activation-incomplete`
  code. The Developer aligned the implementation before further verification.
- The byte-identical full pair command then printed
  `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	1.834s`
  (terminated `exit_code: 0`). Both failed attempts remain above.
- Shared routing regression check:
  `go test -tags integration ./cmd/fid/integration/... -run '^TestP22DualAdapter(Recovery|Activation)Conformance$' -count=1`
  printed `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	4.489s`
  (terminated `exit_code: 0`). Both original recovery and added activation
  corpora executed after the shared adapter routing changes.

The phase-21 activation value is retained because it is the accepted installed
activation phase; using the development phase number 22 was corrected during
fixture construction before a root behavioral run. No product phase policy
changed. Developer final-source gates remain independently required and are
recorded in task-11-done.md after ownership handback.
