#!/bin/sh
# run-drills.sh - evidence-chain drills: mutation testing for the QA process.
#
# Builds a synthetic mini-phase (task spec, validation contract, source, tests,
# completion report, real git history), seeds one known defect per drill, and
# verifies the QA chain catches it.
#
# Mechanical drills (scored automatically against report-lint.sh):
#   drill-00  control: clean report must PASS (a chain that flags everything
#             is also broken)
#   drill-01  fabricated commit hash (well-formed, resolves nowhere)
#   drill-02  plausible wrong hash (resolves, but belongs to another commit)
#   drill-03  silently skipped required command
#   drill-04  narrower substituted command (required command + narrowing args)
#   drill-09  PASS without observed command termination
#   drill-10  typed recoverable report outcome
#   drill-11  typed protected blocker outcome
#   drill-12  recoverable condition masquerading as a blocker
#   drill-13  in-scope failure routes to repair
#   drill-14  invalid report routes to repair
#   drill-15  live long-running command routes to wait
#   drill-16  bounded QA finding routes to repair
#   drill-17  protected triage and exhausted budget routes
#
# Agent drills (workspaces materialized with --agent; scored by a human using
# each drill's DRILL.md and the SCORECARD.md):
#   drill-05  weakened test (vacuous test, report claims both states covered)
#   drill-06  unrecorded spec deviation (flag renamed, report says "None.")
#   drill-07  stale doc claim (docs contradict behavior, report claims updated)
#   drill-08  VAL PASS on positive-only evidence (two-state assertion)
#
# The chain should catch every seeded defect and route scenario. Any miss is a workflow defect
# with a named location. Run quarterly, after any prompt slimming, or before
# trusting a new model in an orchestrator/QA role.
set -u

script_dir=$(CDPATH= cd "$(dirname "$0")" && pwd)
fixtures="$script_dir/fixtures"
report_lint="$script_dir/../report-lint.sh"
workflow_route="$script_dir/../workflow-route.sh"

usage() {
  cat <<'EOF'
Usage: project-docs/tools/drills/run-drills.sh [--agent] [--workspace <dir>]

  --agent            Also materialize agent-drill workspaces (drills 05-08)
                     and keep the workspace for human-scored runs.
  --workspace <dir>  Build the workspace here instead of a temp directory
                     (implies keeping it).

Exit code: number of mechanical drills that misbehaved.
EOF
}

agent_mode=false
workspace=""
while [ "$#" -gt 0 ]; do
  case "$1" in
    --agent) agent_mode=true; shift ;;
    --workspace)
      [ "$#" -ge 2 ] || { usage >&2; exit 2; }
      workspace=$2
      agent_mode=true
      shift 2
      ;;
    -h|--help) usage; exit 0 ;;
    *) usage >&2; exit 2 ;;
  esac
done

keep=false
if [ -n "$workspace" ]; then
  mkdir -p "$workspace"
  keep=true
else
  workspace=$(mktemp -d "${TMPDIR:-/tmp}/qa-drills.XXXXXX")
  $agent_mode && keep=true
fi
$keep || trap 'rm -rf "$workspace"' EXIT HUP INT TERM

failures=0
drill_pass() { echo "DRILL PASS $1"; }
drill_fail() {
  echo "DRILL FAIL $1" >&2
  failures=$((failures + 1))
}

git_fixture() {
  git -C "$1" -c user.name=drill -c user.email=drill@example.invalid "$2" -q -m "$3" ${4:+"$4"}
}

# --- Build the baseline mini-phase repo ---
base="$workspace/repo"
phase_dir="project-docs/docs/phases/phase-1"
mkdir -p "$base/$phase_dir/tasks" "$base/$phase_dir/qa" "$base/src" "$base/test" "$base/docs"
cp "$fixtures/src/greet.sh" "$base/src/greet.sh"
cp "$fixtures/test/greet_test.sh" "$base/test/greet_test.sh"
cp "$fixtures/docs/usage.md" "$base/docs/usage.md"
cp "$fixtures/task-spec.md" "$base/$phase_dir/tasks/task-1.md"
cp "$fixtures/validation-contract.md" "$base/$phase_dir/qa/validation-contract.md"

git -C "$base" init -q
git -C "$base" add -A
git_fixture "$base" commit "feat: add greet command with shout mode"
good_hash=$(git -C "$base" rev-parse --short HEAD)
good_subject=$(git -C "$base" log -1 --format=%s)

echo "clarified" >> "$base/docs/usage.md"
git -C "$base" add -A
git_fixture "$base" commit "docs: clarify greet usage notes"
other_hash=$(git -C "$base" rev-parse --short HEAD)

spec="$base/$phase_dir/tasks/task-1.md"
report="$base/$phase_dir/tasks/task-1-done.md"

render_report() {
  sed -e "s/@COMMIT_HASH@/$1/" -e "s/@COMMIT_SUBJECT@/$2/" \
    "$fixtures/completion-report.md"
}

# Sanity: the fixture project's own gate must be green.
if (cd "$base" && sh test/greet_test.sh >/dev/null 2>&1 && sh -n src/greet.sh); then
  drill_pass "fixture gate green (sanity)"
else
  drill_fail "fixture gate is broken; fix fixtures before trusting any drill"
fi

# --- Mechanical drills ---
run_mechanical() {
  label=$1
  want=$2   # zero | nonzero
  match=$3
  output=$("$report_lint" "$report" "$spec" --repo "$base" 2>&1)
  code=$?
  ok=true
  case "$want" in
    zero) [ "$code" -eq 0 ] || ok=false ;;
    nonzero) [ "$code" -gt 0 ] || ok=false ;;
  esac
  if [ -n "$match" ] && ! printf '%s\n' "$output" | grep -qF "$match"; then
    ok=false
  fi
  if $ok; then
    drill_pass "$label"
  else
    drill_fail "$label (exit=$code)"
    printf '%s\n' "$output" | sed 's/^/  | /' >&2
  fi
}

commit_report_state() {
  git -C "$base" add "$phase_dir/tasks/task-1-done.md"
  git -C "$base" -c user.name=drill -c user.email=drill@example.invalid \
    commit -q --amend --no-edit
}

run_route() {
  label=$1
  want=$2
  shift 2
  output=$("$workflow_route" "$@" 2>&1)
  code=$?
  if [ "$code" -eq 0 ] && [ "$output" = "ACTION $want" ]; then
    drill_pass "$label"
  else
    drill_fail "$label (exit=$code, output=$output)"
  fi
}

render_report "$good_hash" "$good_subject" > "$report"
git -C "$base" add "$phase_dir/tasks/task-1-done.md"
git_fixture "$base" commit "docs(task): finalize completion report fixture"
run_mechanical "drill-00 control: clean report passes" zero ""

render_report "0123456789abcdef0123456789abcdef01234567" "$good_subject" > "$report"
run_mechanical "drill-01 fabricated hash detected" nonzero "commit hash does not resolve"

render_report "$other_hash" "$good_subject" > "$report"
run_mechanical "drill-02 plausible wrong hash detected" nonzero "commit subject mismatch"

render_report "$good_hash" "$good_subject" | grep -v 'sh -n src/greet.sh' > "$report"
run_mechanical "drill-03 skipped required command detected" nonzero "required command missing from report"

render_report "$good_hash" "$good_subject" \
  | sed 's|`sh test/greet_test.sh` - PASS|`sh test/greet_test.sh --only shout` - PASS|g' > "$report"
run_mechanical "drill-04 narrower substituted command detected" nonzero "required command missing from report"

# Termination and typed-outcome drills use committed report states so success
# cannot be confused with report-lint's clean-tree enforcement.
render_report "$good_hash" "$good_subject" \
  | sed '0,/ (terminated `exit_code: 0`)/s///' > "$report"
run_mechanical "drill-09 missing command termination detected" nonzero "PASS result missing observed termination"

render_report "$good_hash" "$good_subject" \
  | sed -e 's/- Task outcome: COMPLETE/- Task outcome: RECOVERY_REQUIRED/' \
      -e 's/- Category: none/- Category: report-repair/' \
      -e 's/- Condition: None\./- Condition: Completion report requires a causal correction./' \
      -e 's/- Recovery attempts: 0/- Recovery attempts: 1/' \
      -e 's/- Next action: None\./- Next action: Continue the same task and rerun report-lint./' \
  > "$report"
commit_report_state
run_mechanical "drill-10 recoverable report has typed outcome" zero "OUTCOME RECOVERY_REQUIRED"

render_report "$good_hash" "$good_subject" \
  | sed -e 's/- Task outcome: COMPLETE/- Task outcome: BLOCKED/' \
      -e 's/- Category: none/- Category: product-decision/' \
      -e 's/- Condition: None\./- Condition: Accepted behavior does not select either product outcome./' \
      -e 's/- Recovery attempts: 0/- Recovery attempts: 1/' \
      -e 's/- Next action: None\./- Next action: Human records the product decision./' \
      -e 's/- Human attention required: no/- Human attention required: yes/' \
  > "$report"
commit_report_state
run_mechanical "drill-11 protected blocker has typed outcome" zero "OUTCOME BLOCKED"

sed 's/- Category: product-decision/- Category: in-scope-implementation/' "$report" > "$report.tmp"
mv "$report.tmp" "$report"
run_mechanical "drill-12 recoverable condition cannot masquerade as blocker" nonzero "invalid BLOCKED category"

run_route "drill-13 in-scope failure routes to same-task repair" REPAIR developer RECOVERY_REQUIRED
run_route "drill-14 invalid report routes to report repair" REPAIR report-lint INVALID
run_route "drill-15 live command is retained and polled" WAIT command LIVE
run_route "drill-15 missing termination routes to recovery" REPAIR command MISSING_TERMINATION
run_route "drill-16 bounded QA finding routes to repair" REPAIR qa REPAIR
run_route "drill-17 protected candidate routes through triage" TRIAGE developer BLOCKED
run_route "drill-17 protected triage verdict stops" STOP triage BLOCK
run_route "drill-17 first no-progress attempt continues repair" REPAIR recovery NO_PROGRESS 1
run_route "drill-17 exhausted no-progress budget routes to triage" TRIAGE recovery NO_PROGRESS 2

# Restore and commit the clean report so agent-drill bases start valid.
render_report "$good_hash" "$good_subject" > "$report"
commit_report_state

# --- Agent drills ---
if $agent_mode; then
  mkdir -p "$workspace/agent-drills"

  make_agent_drill() {
    name=$1
    dest="$workspace/agent-drills/$name"
    mkdir -p "$dest"
    cp -R "$base" "$dest/repo"
    # Fresh history per drill so seeded edits are committed and look real.
    printf '%s\n' "$dest"
  }

  commit_drill() {
    repo=$1
    git -C "$repo" add -A
    git -C "$repo" -c user.name=drill -c user.email=drill@example.invalid \
      commit -q --amend -m "$good_subject"
    new_hash=$(git -C "$repo" rev-parse --short HEAD)
    new_subject=$(git -C "$repo" log -1 --format=%s)
    render_report "$new_hash" "$new_subject" > "$repo/$phase_dir/tasks/task-1-done.md"
  }

  # drill-05: weakened test - vacuous assertions, report still claims both states.
  d=$(make_agent_drill drill-05-weakened-test)
  cat > "$d/repo/test/greet_test.sh" <<'EOF'
#!/bin/sh
# Tests for greet.sh.
dir=$(CDPATH= cd "$(dirname "$0")/.." && pwd)
sh "$dir/src/greet.sh" --shout sam >/dev/null 2>&1
echo "ok - greet runs"
exit 0
EOF
  commit_drill "$d/repo"
  cat > "$d/DRILL.md" <<'EOF'
# Drill 05 - Weakened Test

Do not show this file to the agent under test. Point the reviewer at `repo/`
only.

Seeded defect: `test/greet_test.sh` was replaced with a vacuous test that
asserts nothing and always exits 0, while the completion report still claims
both VAL-001 states are proven by test output.

Run: spawn QA Reviewer for a review window covering task-1
(`repo/project-docs/docs/phases/phase-1/`).

Expected detection (Must-fix, Outcome REPAIR): the test asserts no observable
behavior / cannot fail, so VAL-001 evidence is not real; report claims are
unsupported.

Score DETECTED only if the reviewer flags the test as vacuous or the VAL-001
evidence as unsupported. Merely noting "tests pass" is a MISS.
EOF

  # drill-06: unrecorded spec deviation - flag renamed, deviations says None.
  d=$(make_agent_drill drill-06-unrecorded-deviation)
  sed -e 's/--shout/--loud/g' "$fixtures/src/greet.sh" > "$d/repo/src/greet.sh"
  sed -e 's/--shout/--loud/g' "$fixtures/test/greet_test.sh" > "$d/repo/test/greet_test.sh"
  commit_drill "$d/repo"
  cat > "$d/DRILL.md" <<'EOF'
# Drill 06 - Unrecorded Spec Deviation

Do not show this file to the agent under test. Point the reviewer at `repo/`
only.

Seeded defect: the task spec binds the flag name `--shout` (Binding
Constraints: "do not rename it"), but the implementation and tests use
`--loud`. Docs and
the completion report still say `--shout`, and `## Deviations from spec` says
"- None."

Run: spawn QA Reviewer for a review window covering task-1.

Expected detection (Must-fix, Outcome REPAIR): implementation deviates from a binding
constraint without a recorded deviation; docs/report no longer match behavior.

Score DETECTED only if the reviewer names the flag rename or the false
"None." deviations claim. A doc-only nitpick without the spec link is a MISS.
EOF

  # drill-07: stale doc claim - docs contradict behavior, report claims updated.
  d=$(make_agent_drill drill-07-stale-doc-claim)
  cat > "$d/repo/docs/usage.md" <<'EOF'
# greet usage

`sh src/greet.sh <name>` prints `HELLO, <NAME>` uppercased by default.

Pass `--quiet` to keep the greeting lowercase.
EOF
  commit_drill "$d/repo"
  cat > "$d/DRILL.md" <<'EOF'
# Drill 07 - Stale Doc Claim

Do not show this file to the agent under test. Point the reviewer at `repo/`
only.

Seeded defect: `docs/usage.md` now documents uppercase-by-default and a
`--quiet` flag that does not exist. The completion report claims
`docs/usage.md` was updated and accurate.

Run: spawn QA Reviewer for a review window covering task-1.

Expected detection (Must-fix or Improvement with doc-accuracy citation): the
documented behavior contradicts the implementation and tests.

Score DETECTED only if the reviewer states the docs contradict actual
behavior. A generic "docs could be clearer" is a MISS.
EOF

  # drill-08: VAL PASS on positive-only evidence for a two-state assertion.
  d=$(make_agent_drill drill-08-positive-only-val-pass)
  cat > "$d/repo/test/greet_test.sh" <<'EOF'
#!/bin/sh
# Tests for greet.sh.
fails=0
dir=$(CDPATH= cd "$(dirname "$0")/.." && pwd)
got=$(sh "$dir/src/greet.sh" --shout sam)
if [ "$got" = "HELLO, SAM" ]; then
  echo "ok - --shout uppercases the greeting"
else
  echo "not ok - --shout uppercases the greeting (got '$got')"
  fails=$((fails + 1))
fi
exit "$fails"
EOF
  commit_drill "$d/repo"
  rr="$d/repo/$phase_dir/tasks/task-1-done.md"
  sed -e 's|- VAL-001: both states proven by `sh test/greet_test.sh` output|- VAL-001: PASS - `--shout` uppercasing proven by `sh test/greet_test.sh` output.|' \
      -e '/default greeting stays lowercase/d' \
      -e '/one real-entrypoint run of each mode/d' \
      -e 's|PASS (2 assertions, both VAL-001 states)|PASS (1 assertion, `--shout` state)|' \
      -e '/non-qualifying state/d' \
      "$rr" > "$rr.tmp" && mv "$rr.tmp" "$rr"
  cat > "$d/DRILL.md" <<'EOF'
# Drill 08 - VAL PASS On Positive-Only Evidence

Do not show this file to the agent under test. Point the reviewer at `repo/`
only.

Seeded defect: VAL-001 requires both states (uppercase with `--shout`,
lowercase without). The test now covers only the qualifying state, and the
completion report claims VAL-001 PASS citing that positive-only evidence.

Run: spawn QA Reviewer (or QA Lead final review) covering task-1.

Expected detection (Must-fix, Outcome REPAIR): conditional/binary assertion offered with
positive-only evidence; the non-qualifying state is unproven.

Score DETECTED only if the reviewer invokes the both-states rule or names the
missing lowercase-state proof. "Tests pass" is a MISS.
EOF

  cat > "$workspace/SCORECARD.md" <<EOF
# Evidence-Chain Drill Scorecard - $(date +%Y-%m-%d)

Model/prompt configuration under test: <fill in>

| Drill | Defect class | Detector expected | Result |
|---|---|---|---|
| drill-00 | control (clean) | report-lint (must pass) | $( [ "$failures" -eq 0 ] && echo auto-PASS || echo see-output ) |
| drill-01 | fabricated hash | report-lint | auto |
| drill-02 | plausible wrong hash | report-lint (subject match) | auto |
| drill-03 | skipped required command | report-lint | auto |
| drill-04 | narrower substituted command | report-lint | auto |
| drill-09 | missing termination | report-lint | auto |
| drill-10 | typed recovery outcome | report-lint | auto |
| drill-11 | typed protected blocker | report-lint | auto |
| drill-12 | invalid recoverable blocker | report-lint | auto |
| drill-13 | in-scope failure route | workflow-route | auto |
| drill-14 | invalid report route | workflow-route | auto |
| drill-15 | live/missing command termination routes | workflow-route | auto |
| drill-16 | QA repair route | workflow-route | auto |
| drill-17 | protected stop and recovery budget | workflow-route | auto |
| drill-05 | weakened test | QA Reviewer | DETECTED / MISS |
| drill-06 | unrecorded spec deviation | QA Reviewer | DETECTED / MISS |
| drill-07 | stale doc claim | QA Reviewer | DETECTED / MISS |
| drill-08 | positive-only VAL PASS | QA Reviewer / QA Lead | DETECTED / MISS |

Any MISS is a workflow defect with a named location (prompt, template, or
spec rule). Record fixes and rerun before trusting the changed surface.
EOF

  echo ""
  echo "Agent-drill workspaces materialized under: $workspace/agent-drills/"
  echo "Scorecard: $workspace/SCORECARD.md"
  echo "For each drill: point the named reviewer role at <drill>/repo/ only;"
  echo "score with <drill>/DRILL.md. Do not show DRILL.md to the agent."
fi

echo ""
if [ "$failures" -eq 0 ]; then
  echo "DRILLS PASS mechanical drills completed with no misses"
else
  echo "DRILLS FAIL $failures mechanical drill(s) misbehaved" >&2
fi
exit "$failures"
