# Phase 22 task 20 — Shorten PR feedback without reducing validation

## Authority and objective

Human approved the performance-repair task after reviewing pr-ci-diagnosis.md.
DEC-221 authorizes a bounded serial root repair: measure uncached timings,
optimize redundant fixture preparation, consolidate duplicate CI race execution,
and use exhaustive measured sharding only if needed. ESC-22-005 is resolved as
a decision request by choosing performance repair; CI failure remains open until
verified. The default ten-minute package timeout remains unchanged.

## Scope

- .github/workflows/ci.yml and Makefile only as necessary to run one complete
  race-and-coverage suite; preserve format, lint and build checks. .gitignore may
  ignore only the new controller tests' Python bytecode directory.
- cmd/fid/stop_operation_matrix_test.go and narrow accompanying fixture regression
  helpers/tests; existing command fixture machinery may be reused.
- internal/export/*_test.go for immutable source fixture preparation and distinct
  physical copies. Keep archive validation/import/rebuild paths real.
- Baseline justifies four exhaustive CI shards: cmd/fid 541.652s, workflow
  503.550s, export 368.615s. Own project-docs/tools/ci_tests.py and
  project-docs/tools/test_ci_tests.py for discovery, execution reconciliation,
  coverage union and regressions, plus project-docs/tools/README.md usage.
- Task report, PR diagnosis, QA acceptance addendum, phase decision/progress/ledger,
  current context, quality-plan CI cadence and PR description.

No production Go, migration, dependency, public-interface, archive limit, 100k
performance limit, timeout, assertion meaning, reduced race scope, test deletion
or ordinary-test deferral. Existing historical failed receipts remain intact.
Root owns all changes serially under the Human's delegated repair authority.

## Acceptance criteria

1. Record uncached package/test timings before and after, retaining failures and
   actual process termination. Distinguish local timings from hosted CI evidence.
2. Every original test case and its observable assertions remain. Cached setup is
   immutable; each case gets a distinct writable filesystem/database. Regression
   checks prove sibling/source isolation and marker execution in the correct root.
3. Source fixture reuse must not replace import/rebuild under test or accidentally
   make changed-digest/identity tests compare identical artifacts.
4. CI executes the ordinary suite once with race and coverage; no duplicate full
   test pass. Format/lint/build remain. Any sharding discovers all current tests,
   includes future tests, rejects empty/overlapping/missing selection, preserves
   package-level initialization, and combines coverage without dropping statements.
5. Evidence supports reduced work/wall time. Do not claim the reported hosted CI
   failure resolved until an actual successful hosted rerun is available.

## Validation assertions

Retained VAL-036..039/075 Stop command guards; VAL-054..065 archive validation,
identity and import/fault parity. No assertion definition changes. Fixture isolation
and CI scheduling controls have no independent product VAL; they preserve evidence.

## Exact Commands To Run

From /home/mark/projects/fid, after focused red/green and timing diagnostics:

```bash
gofmt -l cmd internal
golangci-lint run ./...
go test -race -count=1 -json ./...
make coverage
make build
python3 -m unittest discover -s project-docs/tools -p test_ci_tests.py -v
python3 project-docs/tools/ci_tests.py plan --count 4
```

Task-specific exception: no extra tagged integration/benchmark/archive_scale run
is required for changes confined to ordinary-test fixture setup and CI execution;
those sources and production behavior are unchanged. If scope expands to a shared
runtime or tagged fixture, amend this allocation and verification before editing.
The uncached timing pass is comparison evidence; the additional coverage pass is
qualification of the exact CI target, not a proposed duplicate CI step.

## Expected Non-Zero Commands

None for closing gates. Focused regression red failures are preserved as TDD
observations and are never reported as passing gates.

## Development dependency budget

The controller uses Python 3.10+ standard library, available on the existing
Ubuntu CI runner and locally (Python 3.14.4); no pip package, Go module or fid
runtime dependency. A short controller is preferable to fragile shell parsing
of JSON and coverage files. Existing Go remains the actual test/coverage engine.
New official GitHub upload-artifact v7.0.1 and download-artifact v8.0.1 actions
transfer shard manifests/coverage in this workflow. Alternatives: repeating the
suite in a reporting job (defeats this repair), or custom authenticated artifact
API code (more code and token handling). Both actions are MIT-licensed; primary
release/action/license sources are recorded in the completion report.
