# fid

fid is a local-first work ledger with receipts for agent-written code: a Go
CLI over a canonical journal with a derived SQLite projection. It is being
reset. [ADR-0013](docs/adr/0013-reset-fid-to-a-ledger-with-receipts.md) says
what fid becomes; [RESET.md](RESET.md) is the plan, the rules and the task
list. Read RESET.md before starting work.

## How we work

1. **One task, one PR, one session.** Mark names the task, or take the first
   unfinished agent task in RESET.md; tasks run in table order. If it won't
   fit one reviewable PR, stop and propose a split.
2. **Branch** from an up-to-date `main`, named `reset-<task>-<slug>`.
3. **Set the lane** using RESET.md's rules: Routine, Standard or Assured. If in
   doubt, pick the higher lane. Never lower a lane; only Mark does.
4. **TDD for behaviour changes.** Write the failing test, watch it fail, then
   make it pass. Test observable behaviour, not internals. Never delete or
   weaken a test to get green.
5. **Run the checks**: `make fmt-check lint test build`. For CLI journeys, also
   run `go test -tags integration ./cmd/fid/integration/...`. As in CI, the
   race detector runs only in packages whose non-test code starts goroutines
   (decision G, task 0.7). `make test` finds them by parsing the source; the
   integration tests have none, and run fid as a binary built without race.
   Run the checks one at a time: `make test` takes about three minutes and
   the integration suite about five, and two suites at once slow each other
   enough to have pushed `cmd/fid` past Go's ten-minute test timeout.
6. **Open the PR.** Its description is the completion report:
   - what changed and why;
   - the lane;
   - each command you ran, with its last output line;
   - decisions and deviations;
   - anything left over.

   In the same PR, mark the task done in RESET.md by appending `(done <date>)`
   to its "Done when" cell. It counts as done once merged into `main` with CI
   green on that commit.
7. **Review.** A fresh session on the other model reviews the PR (Codex
   reviews Claude's work, Claude reviews Codex's), covering code and test
   adequacy in one pass. Findings go in a PR comment. Fix each one, or answer
   it with a reason.
8. **Mark merges.** Agents never merge, push to `main` or force-push.

If you have to stop early, push the branch and open a draft PR that says
where you stopped and what comes next. The next session resumes from that
branch and PR; it's still one task and one PR.

## Evidence

- Claim only what you observed in this session. Report failures as failures,
  with their output.
- A command that didn't finish didn't pass. `go test ./...` prints nothing
  while it compiles, sometimes for minutes; wait for it.
- Numbers such as diff size, coverage and function length are signals for
  judgement, not gates. Refuse only on missing evidence.

## Stop and ask Mark when

- a product, architecture, data-model or public-interface decision isn't
  already recorded in ADR-0013, the decisions record or the specs;
- the work would add a dependency (say why it's needed, the alternatives and
  the licence; prefer MIT, BSD, Apache-2.0 or MPL-2.0);
- passing a check would mean changing the check;
- the work falls outside the task. Note discoveries in the PR; don't build
  them.

## Code

- Go, a `spf13/cobra` CLI and pure-Go SQLite (`modernc.org/sqlite`). The
  journal is canonical; SQLite is a projection (ADR-0001).
- `cmd/fid/` is the entry point and `internal/cli/` the command adapters.
  Domain packages live under `internal/` (issues, claims, workflow, review,
  evidence, …). Persistence lives in `internal/journal/`, `internal/db/` and
  `migrations/`. Keep business rules out of `internal/cli/`.
- Return errors; don't panic. Wrap them with `fmt.Errorf("op: %w", err)`.
- Use table-driven tests. Top-level tests call `t.Parallel()` unless they
  touch process-global state; the `paralleltest` linter enforces this, and
  exceptions carry `//nolint:paralleltest // <reason>`.
- Functions around 60 lines and files around 400 lines are guides.
- Use Conventional Commits, with a `Refs: RESET <task>` footer.
- In Codex's sandbox, `.codex/config.toml` points the Go caches at gitignored
  repo directories, so don't ask for unsandboxed runs just to write caches.

## Where things are

- `RESET.md`: plan, rules, tasks, validations and kill criteria.
- `docs/adr/`: decisions. ADR-0013 overrides parts of older ADRs; each
  affected ADR carries a banner saying which parts.
- `docs/reset/2026-10-31-reset-decisions.md`: decision detail and the later
  list. The rest of `docs/reset/` is the superseded reset proposal and its
  research, kept as background.
- `docs/specs/`: specs, with `GLOSSARY.md` for terms. Where they conflict with
  ADR-0013 or RESET.md, those win.
- **The V1 workflow material is no longer in the repo.** Task 0.5 archived it.
  The specs still cite V1 phase documents (`docs/phases/...`),
  `09_AI_WORKFLOW.md` and `10_IMPLEMENTATION_ROADMAP.md`; don't go looking
  for them. If a task depends on one, stop and ask Mark.
