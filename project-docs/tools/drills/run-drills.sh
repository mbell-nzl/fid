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
#   drill-18  bounded control-conformance routing, scope, and ledger directive
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
sequence_lint="$script_dir/../workflow-sequence-lint.sh"
control_scope_lint="$script_dir/../workflow-control-scope-lint.sh"

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

route_value() {
  "$workflow_route" "$@" | sed 's/^ACTION //'
}

run_sequence() {
  label=$1
  want=$2
  ledger=$3
  output=$("$sequence_lint" "$ledger" 2>&1)
  code=$?
  if [ "$code" -eq 0 ] && [ "$output" = "VALID state=$want" ]; then
    drill_pass "$label"
  else
    drill_fail "$label (exit=$code, output=$output)"
  fi
}

run_invalid_sequence() {
  label=$1
  match=$2
  ledger=$3
  output=$("$sequence_lint" "$ledger" 2>&1)
  code=$?
  if [ "$code" -ne 0 ] && printf '%s\n' "$output" | grep -qF "$match"; then
    drill_pass "$label"
  else
    drill_fail "$label (exit=$code, output=$output)"
  fi
}

run_invalid_append() {
  label=$1
  match=$2
  ledger=$3
  candidate=$4
  before="$workspace/invalid-append.before"
  cp "$ledger" "$before"
  output=$("$sequence_lint" "$ledger" --append-stdin < "$candidate" 2>&1)
  code=$?
  if [ "$code" -ne 0 ] && printf '%s\n' "$output" | grep -qF "$match" &&
     cmp -s "$before" "$ledger"; then
    drill_pass "$label"
  else
    drill_fail "$label (exit=$code, output=$output)"
  fi
}

run_control_scope() {
  label=$1
  want=$2
  match=$3
  shift 3
  output=$("$control_scope_lint" "$@" 2>&1)
  code=$?
  ok=true
  case "$want" in
    zero) [ "$code" -eq 0 ] || ok=false ;;
    nonzero) [ "$code" -ne 0 ] || ok=false ;;
  esac
  if [ -n "$match" ] && ! printf '%s\n' "$output" | grep -qF "$match"; then
    ok=false
  fi
  if $ok; then
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
  | sed '0,/exit_code: 0/s//exit_code: 1/' > "$report"
run_mechanical "drill-09 PASS with nonzero exit detected" nonzero "PASS result did not terminate successfully with exit_code 0"

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
run_route "drill-18 control conformance defect routes to bounded repair" REPAIR control CONFORMANCE_DEFECT
run_route "drill-18 control definition change routes through triage" TRIAGE control DEFINITION_CHANGE
run_route "drill-18 ambiguous control change routes through triage" TRIAGE control AMBIGUOUS

# --- Drill 18: control-scope lint isolation and lifecycle ---
scope_repo="$workspace/control-scope"
mkdir -p "$scope_repo/project-docs/tools" "$scope_repo/cmd"
printf '%s\n' 'old control' > "$scope_repo/project-docs/tools/control.sh"
printf '%s\n' 'old product' > "$scope_repo/cmd/product.go"
git -C "$scope_repo" init -q
git -C "$scope_repo" add -A
git_fixture "$scope_repo" commit "test: establish control scope fixture"
scope_base=$(git -C "$scope_repo" rev-parse HEAD)
printf '%s\n' 'repaired control' > "$scope_repo/project-docs/tools/control.sh"
git -C "$scope_repo" add -A
git_fixture "$scope_repo" commit "fix: repair control implementation"
scope_control_head=$(git -C "$scope_repo" rev-parse HEAD)
run_control_scope "drill-18 isolated control commit satisfies allowlist" zero "PASS control repair scope" \
  "$scope_repo" "$scope_base" "$scope_control_head" project-docs/tools/control.sh
run_control_scope "drill-18 symbolic repair head cannot replace current-head identity" nonzero \
  "current head must be a full canonical commit ID" \
  "$scope_repo" "$scope_base" HEAD project-docs/tools/control.sh
printf '%s\n' 'dirty product' > "$scope_repo/cmd/product.go"
run_control_scope "drill-18 dirty product state invalidates control repair scope" nonzero "worktree is not clean" \
  "$scope_repo" "$scope_base" "$scope_control_head" project-docs/tools/control.sh
git -C "$scope_repo" checkout -q -- cmd/product.go
git -C "$scope_repo" update-index --assume-unchanged cmd/product.go
printf '%s\n' 'hidden product' > "$scope_repo/cmd/product.go"
run_control_scope "drill-18 assume-unchanged cannot hide forbidden product state" nonzero "assume-unchanged or skip-worktree" \
  "$scope_repo" "$scope_base" "$scope_control_head" project-docs/tools/control.sh
git -C "$scope_repo" update-index --no-assume-unchanged cmd/product.go
git -C "$scope_repo" checkout -q -- cmd/product.go
git -C "$scope_repo" update-index --skip-worktree cmd/product.go
printf '%s\n' 'hidden product' > "$scope_repo/cmd/product.go"
run_control_scope "drill-18 skip-worktree cannot hide forbidden product state" nonzero "assume-unchanged or skip-worktree" \
  "$scope_repo" "$scope_base" "$scope_control_head" project-docs/tools/control.sh
git -C "$scope_repo" update-index --no-skip-worktree cmd/product.go
git -C "$scope_repo" checkout -q -- cmd/product.go
printf '%s\n' 'changed product' > "$scope_repo/cmd/product.go"
git -C "$scope_repo" add -A
git_fixture "$scope_repo" commit "test: contaminate control repair"
scope_product_head=$(git -C "$scope_repo" rev-parse HEAD)
newline_allowlist=$(printf 'project-docs/tools/control.sh\ncmd/product.go')
run_control_scope "drill-18 LF allowlist injection cannot authorize product scope" nonzero \
  "allowed path contains LF, CR, or control bytes" \
  "$scope_repo" "$scope_control_head" "$scope_product_head" "$newline_allowlist"
carriage_allowlist=$(printf 'project-docs/tools/control.sh\rcmd/product.go')
run_control_scope "drill-18 CR is rejected in allowlist argv" nonzero \
  "allowed path contains LF, CR, or control bytes" \
  "$scope_repo" "$scope_control_head" "$scope_product_head" "$carriage_allowlist"
tab_allowlist=$(printf 'project-docs/tools/control.sh\tcmd/product.go')
run_control_scope "drill-18 control byte is rejected in allowlist argv" nonzero \
  "allowed path contains LF, CR, or control bytes" \
  "$scope_repo" "$scope_control_head" "$scope_product_head" "$tab_allowlist"
run_control_scope "drill-18 stale claimed repair checkpoint is rejected" nonzero "repair commit is not the checked-out HEAD" \
  "$scope_repo" "$scope_base" "$scope_control_head" project-docs/tools/control.sh
run_control_scope "drill-18 product edit invalidates control repair scope" nonzero "outside allowed control scope: cmd/product.go" \
  "$scope_repo" "$scope_control_head" "$scope_product_head" project-docs/tools/control.sh
run_control_scope "drill-18 protected path cannot be allowlisted" nonzero "allowed path is not a workflow control" \
  "$scope_repo" "$scope_control_head" "$scope_product_head" project-docs/tools/control.sh cmd/product.go
printf '%s\n' 'old product' > "$scope_repo/cmd/product.go"
git -C "$scope_repo" add -A
git_fixture "$scope_repo" commit "test: hide product edit in final diff"
scope_reverted_head=$(git -C "$scope_repo" rev-parse HEAD)
run_control_scope "drill-18 reverted intermediate edit cannot hide a non-isolated repair" nonzero \
  "trusted base is not the sole parent of current HEAD" \
  "$scope_repo" "$scope_base" "$scope_reverted_head" project-docs/tools/control.sh

receipt_repo="$workspace/control-receipt"
mkdir -p "$receipt_repo/project-docs/tools"
printf '%s\n' 'old alpha' > "$receipt_repo/project-docs/tools/alpha.sh"
printf '%s\n' 'old beta' > "$receipt_repo/project-docs/tools/beta.sh"
printf '%s\n' 'unchanged' > "$receipt_repo/project-docs/tools/unused.sh"
git -C "$receipt_repo" init -q
git -C "$receipt_repo" add -A
git_fixture "$receipt_repo" commit "test: establish ordered receipt fixture"
receipt_base=$(git -C "$receipt_repo" rev-parse HEAD)
printf '%s\n' 'new alpha' > "$receipt_repo/project-docs/tools/alpha.sh"
printf '%s\n' 'new beta' > "$receipt_repo/project-docs/tools/beta.sh"
git -C "$receipt_repo" add -A
git_fixture "$receipt_repo" commit "fix: change two allowed controls"
receipt_head=$(git -C "$receipt_repo" rev-parse HEAD)
receipt_alpha=project-docs/tools/alpha.sh
receipt_beta=project-docs/tools/beta.sh
receipt_unused=project-docs/tools/unused.sh
authoritative_receipt="PASS control repair scope: $receipt_base -> $receipt_head; allowed=$receipt_alpha,$receipt_beta"
receipt_output=$("$control_scope_lint" "$receipt_repo" "$receipt_base" \
  "$receipt_head" "$receipt_alpha" "$receipt_beta" 2>&1)
receipt_code=$?
if [ "$receipt_code" -eq 0 ] && [ "$receipt_output" = "$authoritative_receipt" ]; then
  drill_pass "drill-18 receipt binds the exact authoritative ordered allowlist"
else
  drill_fail "drill-18 authoritative allowlist receipt (exit=$receipt_code, output=$receipt_output)"
fi
run_control_scope "drill-18 receipt scope rejects a missing changed path" nonzero \
  "outside allowed control scope: $receipt_beta" \
  "$receipt_repo" "$receipt_base" "$receipt_head" "$receipt_alpha"
additional_receipt="PASS control repair scope: $receipt_base -> $receipt_head; allowed=$receipt_alpha,$receipt_beta,$receipt_unused"
additional_output=$("$control_scope_lint" "$receipt_repo" "$receipt_base" \
  "$receipt_head" "$receipt_alpha" "$receipt_beta" "$receipt_unused" 2>&1)
additional_code=$?
if [ "$additional_code" -eq 0 ] && [ "$additional_output" = "$additional_receipt" ] && \
   [ "$additional_output" != "$authoritative_receipt" ]; then
  drill_pass "drill-18 additional allowlist entry produces a distinct receipt"
else
  drill_fail "drill-18 additional allowlist receipt (exit=$additional_code, output=$additional_output)"
fi
reordered_receipt="PASS control repair scope: $receipt_base -> $receipt_head; allowed=$receipt_beta,$receipt_alpha"
reordered_output=$("$control_scope_lint" "$receipt_repo" "$receipt_base" \
  "$receipt_head" "$receipt_beta" "$receipt_alpha" 2>&1)
reordered_code=$?
if [ "$reordered_code" -eq 0 ] && [ "$reordered_output" = "$reordered_receipt" ] && \
   [ "$reordered_output" != "$authoritative_receipt" ]; then
  drill_pass "drill-18 reordered allowlist produces a distinct receipt"
else
  drill_fail "drill-18 reordered allowlist receipt (exit=$reordered_code, output=$reordered_output)"
fi
run_control_scope "drill-18 comma is rejected in allowlist argv" nonzero \
  "allowed path contains comma" \
  "$receipt_repo" "$receipt_base" "$receipt_head" \
  "$receipt_alpha,$receipt_beta"
run_control_scope "drill-18 duplicate allowlist entry is rejected" nonzero \
  "duplicate allowed path: $receipt_alpha" \
  "$receipt_repo" "$receipt_base" "$receipt_head" \
  "$receipt_alpha" "$receipt_beta" "$receipt_alpha"

symlink_repo="$workspace/control-symlink"
mkdir -p "$symlink_repo/project-docs/tools" "$symlink_repo/cmd"
printf '%s\n' '#!/bin/sh' 'exit 0' > "$symlink_repo/project-docs/tools/control.sh"
printf '%s\n' 'package product' > "$symlink_repo/cmd/product.go"
git -C "$symlink_repo" init -q
git -C "$symlink_repo" add -A
git_fixture "$symlink_repo" commit "test: establish regular control fixture"
symlink_base=$(git -C "$symlink_repo" rev-parse HEAD)
rm "$symlink_repo/project-docs/tools/control.sh"
ln -s ../../cmd/product.go "$symlink_repo/project-docs/tools/control.sh"
git -C "$symlink_repo" add -A
git_fixture "$symlink_repo" commit "test: replace control with product symlink"
symlink_head=$(git -C "$symlink_repo" rev-parse HEAD)
run_control_scope "drill-18 symlink cannot masquerade as an allowed control" nonzero "allowed control is not a regular file" \
  "$symlink_repo" "$symlink_base" "$symlink_head" project-docs/tools/control.sh

rebaseline_repo="$workspace/control-rebaseline"
mkdir -p "$rebaseline_repo/project-docs/tools" "$rebaseline_repo/cmd"
printf '%s\n' 'old control' > "$rebaseline_repo/project-docs/tools/control.sh"
printf '%s\n' 'old product' > "$rebaseline_repo/cmd/product.go"
git -C "$rebaseline_repo" init -q
git -C "$rebaseline_repo" add -A
git_fixture "$rebaseline_repo" commit "test: establish trusted pre-dispatch base"
rebaseline_trusted_base=$(git -C "$rebaseline_repo" rev-parse HEAD)
printf '%s\n' 'unauthorized product change' > "$rebaseline_repo/cmd/product.go"
git -C "$rebaseline_repo" add -A
git_fixture "$rebaseline_repo" commit "test: insert forbidden pre-repair commit"
printf '%s\n' 'repaired control' > "$rebaseline_repo/project-docs/tools/control.sh"
git -C "$rebaseline_repo" add -A
git_fixture "$rebaseline_repo" commit "fix: repair control after forbidden commit"
rebaseline_head=$(git -C "$rebaseline_repo" rev-parse HEAD)
run_control_scope "drill-18 symbolic rebaseline cannot hide a forbidden intermediate commit" nonzero \
  "trusted base must be a full canonical commit ID" \
  "$rebaseline_repo" 'HEAD^' "$rebaseline_head" project-docs/tools/control.sh
run_control_scope "drill-18 authoritative pre-dispatch base detects rebaseline bypass" nonzero \
  "trusted base is not the sole parent of current HEAD" \
  "$rebaseline_repo" "$rebaseline_trusted_base" "$rebaseline_head" project-docs/tools/control.sh

verifier_repo="$workspace/control-trusted-verifier"
mkdir -p "$verifier_repo/project-docs/tools"
cp "$control_scope_lint" \
  "$verifier_repo/project-docs/tools/workflow-control-scope-lint.sh"
cp "$sequence_lint" \
  "$verifier_repo/project-docs/tools/workflow-sequence-lint.sh"
printf '%s\n' 'old control product' > \
  "$verifier_repo/project-docs/tools/product-control.sh"
git -C "$verifier_repo" init -q
git -C "$verifier_repo" add -A
git_fixture "$verifier_repo" commit "test: establish trusted verifier fixture"
verifier_base=$(git -C "$verifier_repo" rev-parse HEAD)
printf '%s\n' '#!/bin/sh' 'echo "PASS weakened HEAD verifier"' > \
  "$verifier_repo/project-docs/tools/workflow-control-scope-lint.sh"
printf '%s\n' '#!/bin/sh' 'echo "VALID state=COMPLETE"' > \
  "$verifier_repo/project-docs/tools/workflow-sequence-lint.sh"
printf '%s\n' 'changed control product' > \
  "$verifier_repo/project-docs/tools/product-control.sh"
git -C "$verifier_repo" add -A
git_fixture "$verifier_repo" commit "test: weaken verifier with control product"
verifier_head=$(git -C "$verifier_repo" rev-parse HEAD)
weakened_output=$("$verifier_repo/project-docs/tools/workflow-control-scope-lint.sh" 2>&1)
weakened_code=$?
if [ "$weakened_code" -eq 0 ] && [ "$weakened_output" = 'PASS weakened HEAD verifier' ]; then
  drill_pass "drill-18 weakened HEAD verifier demonstrates the self-repair threat"
else
  drill_fail "drill-18 weakened HEAD verifier threat fixture (exit=$weakened_code, output=$weakened_output)"
fi

# Phase Runner materialization is exact: read the verifier blob from the
# authoritative pre-dispatch base, write it to a fresh path, chmod it executable,
# and execute that path rather than the repair candidate's HEAD copy.
trusted_verifier="$workspace/trusted-workflow-control-scope-lint.sh"
git -C "$verifier_repo" show \
  "$verifier_base:project-docs/tools/workflow-control-scope-lint.sh" > \
  "$trusted_verifier"
chmod 700 "$trusted_verifier"
run_control_scope "drill-18 trusted verifier rejects scope-verifier self-repair" nonzero \
  "trusted verifier cannot be allowlisted: project-docs/tools/workflow-control-scope-lint.sh" \
  "$verifier_repo" "$verifier_base" "$verifier_head" \
  project-docs/tools/workflow-control-scope-lint.sh \
  project-docs/tools/workflow-sequence-lint.sh \
  project-docs/tools/product-control.sh
trusted_sequence_output=$("$trusted_verifier" "$verifier_repo" \
  "$verifier_base" "$verifier_head" \
  project-docs/tools/workflow-sequence-lint.sh \
  project-docs/tools/workflow-control-scope-lint.sh \
  project-docs/tools/product-control.sh 2>&1)
trusted_sequence_code=$?
if [ "$trusted_sequence_code" -ne 0 ] && \
   printf '%s\n' "$trusted_sequence_output" | grep -qF \
     'trusted verifier cannot be allowlisted: project-docs/tools/workflow-sequence-lint.sh'; then
  drill_pass "drill-18 trusted verifier rejects sequence-verifier self-repair"
else
  drill_fail "drill-18 trusted sequence-verifier protection (exit=$trusted_sequence_code, output=$trusted_sequence_output)"
fi

# Exercise the actual control-repair topology: the Phase Runner ledger remains
# dirty in the main worktree while a separate clean repair worktree owns the
# isolated commit. Only the verified commit is fast-forwarded back.
lifecycle_repo="$workspace/control-lifecycle-main"
lifecycle_repair="$workspace/control-lifecycle-repair"
mkdir -p "$lifecycle_repo/project-docs/tools" \
  "$lifecycle_repo/project-docs/docs/phases/phase-1"
printf '%s\n' 'old control' > "$lifecycle_repo/project-docs/tools/control.sh"
printf '%s\n' '# Phase 1 Run Ledger' > \
  "$lifecycle_repo/project-docs/docs/phases/phase-1/run-ledger.md"
git -C "$lifecycle_repo" init -q
git -C "$lifecycle_repo" add -A
git_fixture "$lifecycle_repo" commit "test: establish control lifecycle fixture"
lifecycle_base=$(git -C "$lifecycle_repo" rev-parse HEAD)
printf '%s\n' 'TASK-RECOVERY remains runner-owned and uncommitted' >> \
  "$lifecycle_repo/project-docs/docs/phases/phase-1/run-ledger.md"
git -C "$lifecycle_repo" worktree add -q -b control-repair-test \
  "$lifecycle_repair" "$lifecycle_base"
printf '%s\n' 'repaired control' > \
  "$lifecycle_repair/project-docs/tools/control.sh"
git -C "$lifecycle_repair" add -A
git_fixture "$lifecycle_repair" commit "fix: repair control in dedicated worktree"
lifecycle_head=$(git -C "$lifecycle_repair" rev-parse HEAD)
lifecycle_scope_output=$("$control_scope_lint" "$lifecycle_repair" \
  "$lifecycle_base" "$lifecycle_head" project-docs/tools/control.sh 2>&1)
lifecycle_scope_code=$?
if [ "$lifecycle_scope_code" -eq 0 ]; then
  git -C "$lifecycle_repo" merge --ff-only -q control-repair-test
  lifecycle_merge_code=$?
else
  lifecycle_merge_code=1
fi
lifecycle_main_head=$(git -C "$lifecycle_repo" rev-parse HEAD)
lifecycle_main_status=$(git -C "$lifecycle_repo" status --porcelain=v1 \
  --untracked-files=all)
lifecycle_repair_status=$(git -C "$lifecycle_repair" status --porcelain=v1 \
  --untracked-files=all)
if [ "$lifecycle_scope_code" -eq 0 ] && \
   printf '%s\n' "$lifecycle_scope_output" | grep -qF 'PASS control repair scope' && \
   [ "$lifecycle_merge_code" -eq 0 ] && \
   [ "$lifecycle_main_head" = "$lifecycle_head" ] && \
   [ "$lifecycle_main_status" = ' M project-docs/docs/phases/phase-1/run-ledger.md' ] && \
   [ -z "$lifecycle_repair_status" ]; then
  drill_pass "drill-18 dirty main ledger and clean repair worktree integrate safely"
else
  drill_fail "drill-18 dedicated repair worktree lifecycle (scope=$lifecycle_scope_output, main=$lifecycle_main_status, repair=$lifecycle_repair_status)"
fi

# Multi-step conformance: execute a byte-identical command red, repair its
# fixture causally, rerun it green, then route the preserved attempt sequence.
causal="$workspace/causal-repair"
mkdir -p "$causal/src" "$causal/test"
cp "$fixtures/src/greet.sh" "$causal/src/greet.sh"
cp "$fixtures/test/greet_test.sh" "$causal/test/greet_test.sh"
sed 's/hello, /broken, /' "$causal/src/greet.sh" > "$causal/src/greet.sh.tmp"
mv "$causal/src/greet.sh.tmp" "$causal/src/greet.sh"
causal_command='sh test/greet_test.sh'
if (cd "$causal" && sh -c "$causal_command") > "$causal/first.out" 2>&1; then
  first_code=0
else
  first_code=$?
fi
cp "$fixtures/src/greet.sh" "$causal/src/greet.sh"
if (cd "$causal" && sh -c "$causal_command") > "$causal/second.out" 2>&1; then
  second_code=0
else
  second_code=$?
fi
task_sequence="RED"
if [ "$first_code" -ne 0 ] && [ "$(route_value developer RECOVERY_REQUIRED)" = REPAIR ]; then
  task_sequence="$task_sequence>TASK-RECOVERY>REPAIR>VERIFY"
fi
if [ "$second_code" -eq 0 ] && [ "$(route_value recovery PROGRESS)" = REPAIR ] && \
   [ "$(route_value developer COMPLETE)" = ADVANCE ]; then
  task_sequence="$task_sequence>COMPLETE"
fi
if [ "$task_sequence" = "RED>TASK-RECOVERY>REPAIR>VERIFY>COMPLETE" ]; then
  drill_pass "drill-13 lifecycle: red, causal repair, verify, complete"
else
  drill_fail "drill-13 lifecycle sequence ($task_sequence)"
fi

# Report repair is its own loop and does not imply a product gate passed.
report_sequence="REPORT"
if [ "$(route_value report-lint INVALID)" = REPAIR ]; then
  report_sequence="$report_sequence>TASK-RECOVERY>REPORT"
fi
if [ "$(route_value report-lint COMPLETE)" = ADVANCE ]; then
  report_sequence="$report_sequence>COMPLETE"
fi
if [ "$report_sequence" = "REPORT>TASK-RECOVERY>REPORT>COMPLETE" ]; then
  drill_pass "drill-14 lifecycle: report repair then complete"
else
  drill_fail "drill-14 lifecycle sequence ($report_sequence)"
fi

# Exercise a genuinely live, initially silent process. Silence while the PID
# is live maps to WAIT; only observed process completion maps to EVALUATE.
silent_output="$workspace/silent-command.out"
(sleep 1; printf '%s\n' 'terminal receipt') > "$silent_output" &
silent_pid=$!
sleep 0.1
silent_ok=true
kill -0 "$silent_pid" 2>/dev/null || silent_ok=false
[ ! -s "$silent_output" ] || silent_ok=false
[ "$(route_value command LIVE)" = WAIT ] || silent_ok=false
while kill -0 "$silent_pid" 2>/dev/null; do
  sleep 0.1
done
wait "$silent_pid" || silent_ok=false
[ "$(route_value command TERMINATED)" = EVALUATE ] || silent_ok=false
[ "$(sed -n '1p' "$silent_output")" = "terminal receipt" ] || silent_ok=false
if $silent_ok; then
  drill_pass "drill-15 lifecycle: silent live command polled to termination"
else
  drill_fail "drill-15 live-command lifecycle"
fi

# Ledger fixtures exercise the validator the Phase Runner invokes. They prove
# task identity survives a fresh run and QA repair needs an explicitly linked,
# differently identified PASS re-review.
continuation_phase="$workspace/continuation/phase-1"
mkdir -p "$continuation_phase"
valid_continuation="$continuation_phase/valid-continuation.md"
cat > "$valid_continuation" <<'EOF'
### 2026-08-15 10:00 RUN-START
- Sequence contract: v2
- Recovery directive schema: v2
- Branch: phase-1
- ESC baseline (Open at run start): none
- Warm-up/autonomous-resume state: warm-up complete
- Starting task: task-1
### 2026-08-15 10:01 TASK-RECOVERY task-1
- Source: interrupted task-1 session
- Category: session-continuation
- Repair mode: task-scope
- Repair authority: task-1 spec
- Triage disposition: N/A
- Control contract: N/A
- Control mismatch: N/A
- Allowed files: ordinary task scope
- Forbidden changes: protected decisions and out-of-scope files
- Originating report: N/A
- Exact verification: N/A
- Control repair base: N/A
- Trigger signature: task-1:session-continuation:interrupted-session
- Attempt count: 1
- Consecutive no-progress attempts: 0/2
- Checkpoint: committed partial task-1 implementation
- Phase worktree: workspace/phase-1 on phase-1; ledger-only dirty
- Temporary repair worktree: N/A
- Hidden-index flags: none
- Resolution evidence: PENDING - complete and verify task-1
- Next action: complete task-1 and rerun its exact gate
- Human attention required: no
### 2026-08-15 10:02 REPAIR task-1
- Trigger signature: task-1:session-continuation:interrupted-session
- Attempt count: 1
- Starting checkpoint: committed partial task-1 implementation
- Ending checkpoint: verified task-1 implementation
- Control repair base: N/A
- Control repair head: N/A
- Control scope check: N/A for task-scope
- Phase-branch integration: N/A for task-scope
- Outcome: RESOLVED
- Resolution evidence: exact task gate passed
- Next checkpoint: VERIFY
### 2026-08-15 10:03 TASK task-1 COMPLETE
- Report: project-docs/docs/phases/phase-1/tasks/task-1-done.md
- report-lint: PASS - OUTCOME COMPLETE
- Commit: aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa - verified, subject matched
- QA-plan commands beyond the task spec: none required
EOF
run_sequence "drill-16 lifecycle: fresh session retains task identity" ADVANCEABLE "$valid_continuation"

duplicate_completion_report="$continuation_phase/duplicate-completion-report.md"
sed '/^- Report: project-docs\/docs\/phases\/phase-1\/tasks\/task-1-done.md$/a\
- Report: project-docs/docs/phases/phase-1/tasks/task-9-done.md
' "$valid_continuation" > "$duplicate_completion_report"
run_invalid_sequence "drill-18 schema v2 TASK COMPLETE rejects duplicate report authority" \
  "TASK COMPLETE must contain each completion field exactly once" \
  "$duplicate_completion_report"

review_phase="$workspace/review-schema/phase-1"
mkdir -p "$review_phase"
valid_review="$review_phase/run-ledger.md"
cat > "$valid_review" <<'EOF'
### 2026-08-15 09:00 RUN-START
- Sequence contract: v2
- Recovery directive schema: v2
- Branch: phase-1
- ESC baseline (Open at run start): none
- Warm-up/autonomous-resume state: warm-up complete
- Starting task: task-1
### 2026-08-15 09:01 REVIEW review-1
- Path: project-docs/docs/phases/phase-1/qa/review-1.md
- Tasks covered: task-1
- Outcome: PASS
- Owning task: none
- Re-review of: none
- Must-fix: 0 / Improvement: 0 / blocker candidates: 0; QA minted no ESC
EOF
run_sequence "drill-18 schema v2 REVIEW accepts one complete typed record" \
  ADVANCEABLE "$valid_review"
duplicate_review_owner="$review_phase/duplicate-review-owner.md"
sed '/^- Owning task: none$/i\
- Owning task: task-9
' "$valid_review" > "$duplicate_review_owner"
run_invalid_sequence "drill-18 schema v2 REVIEW rejects duplicate ownership" \
  "REVIEW must contain each review field exactly once" "$duplicate_review_owner"

run_end_phase="$workspace/run-end-schema/phase-1"
mkdir -p "$run_end_phase"
valid_run_end="$run_end_phase/run-ledger.md"
cat > "$valid_run_end" <<'EOF'
### 2026-08-15 09:00 RUN-START
- Sequence contract: v2
- Recovery directive schema: v2
- Branch: phase-1
- ESC baseline (Open at run start): none
- Warm-up/autonomous-resume state: warm-up complete
- Starting task: task-1
### 2026-08-15 09:02 RUN-END
- All tasks complete. Final review: project-docs/docs/phases/phase-1/qa/final-review.md
- Human checklist: project-docs/docs/phases/phase-1/qa/human-checklist.md
- Phase ready for Tech Lead wrap-up.
EOF
run_sequence "drill-18 schema v2 RUN-END accepts one complete closure record" \
  ADVANCEABLE "$valid_run_end"
missing_run_end_field="$run_end_phase/missing-run-end-field.md"
sed '/^- Human checklist: /d' "$valid_run_end" > "$missing_run_end_field"
run_invalid_sequence "drill-18 schema v2 RUN-END cannot omit closure fields" \
  "RUN-END must contain each closure field exactly once" "$missing_run_end_field"

control_phase="$workspace/control-fixture/phase-1"
mkdir -p "$control_phase"
cat > "$control_phase/triage-log.md" <<'EOF'
# Phase 1 Triage Log

## TRG-1-001 - byte-exact control mismatch
- **Originating task:** task-1
- **Trigger signature:** task-1:quality-gate-control-defect:literal-command
- **Repair category:** report-repair
- **Repair authority:** TRG-1-001 typed directive
- **Attempt count:** 1
- **Verdict:** REPAIR / control-conformance
- **Control trigger signature:** task-1:quality-gate-control-defect:literal-command
- **Control attempt count:** 1
- **Control contract:** exact commands are compared byte-for-byte
- **Control mismatch:** exact reproducer false-rejects a literal backslash-dot
- **Allowed files:** project-docs/tools/report-lint.sh
- **Forbidden changes:** command, assertion, accepted status, scope, threshold, budget, or tolerance
- **Exact verification:** project-docs/tools/report-lint.sh --self-test
- **Originating report:** project-docs/docs/phases/phase-1/tasks/task-1-done.md
- **Control repair base:** aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa

## TRG-1-002 - ordinary task-scope repair
- **Originating task:** task-1
- **Trigger signature:** task-1:focused-test:ordinary
- **Repair category:** gate-repair
- **Repair authority:** TRG-1-002
- **Attempt count:** 1
- **Verdict:** REPAIR

## TRG-1-003 - protected product decision
- **Originating task:** task-1
- **Trigger signature:** task-1:product-decision:protected
- **Repair category:** N/A
- **Repair authority:** N/A
- **Attempt count:** 1
- **Verdict:** BLOCK
- **Escalation:** ESC-1-003
EOF
valid_control_repair="$control_phase/valid-control-repair.md"
cat > "$valid_control_repair" <<'EOF'
### 2026-08-15 10:00 RUN-START
- Sequence contract: v2
- Recovery directive schema: v2
- Branch: phase-1
- ESC baseline (Open at run start): none
- Warm-up/autonomous-resume state: warm-up complete
- Starting task: task-1
### 2026-08-15 10:00 TRIAGE TRG-1-001
- Trigger: report-lint violates its byte-exact command contract
- Originating task: task-1
- Trigger signature: task-1:quality-gate-control-defect:literal-command
- Attempt count: 1
- Repair category: report-repair
- Repair authority: TRG-1-001 typed directive
- Checkpoint: exact reproducer failed
- Verdict: REPAIR / control-conformance
- Resolution evidence: bounded tools-only repair
- Triage log: phase triage entry confirmed
### 2026-08-15 10:01 TASK-RECOVERY task-1
- Source: TRG-1-001
- Category: report-repair
- Repair mode: control-conformance
- Repair authority: TRG-1-001 typed directive
- Triage disposition: REPAIR / control-conformance
- Control contract: exact commands are compared byte-for-byte
- Control mismatch: exact reproducer false-rejects a literal backslash-dot
- Allowed files: project-docs/tools/report-lint.sh
- Forbidden changes: command, assertion, accepted status, scope, threshold, budget, or tolerance
- Exact verification: project-docs/tools/report-lint.sh --self-test
- Originating report: project-docs/docs/phases/phase-1/tasks/task-1-done.md
- Control repair base: aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa
- Trigger signature: task-1:quality-gate-control-defect:literal-command
- Attempt count: 1
- Consecutive no-progress attempts: 0/2
- Checkpoint: committed failing control fixture
- Phase worktree: workspace/phase-1 on phase-1; ledger-only dirty
- Temporary repair worktree: workspace/control-repair on control-repair at aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa
- Hidden-index flags: none
- Resolution evidence: PENDING - repair and rerun the exact reproducer
- Next action: repair the control, run workflow-control-scope-lint.sh, and rerun the exact reproducer
- Human attention required: no
### 2026-08-15 10:02 REPAIR task-1
- Trigger signature: task-1:quality-gate-control-defect:literal-command
- Attempt count: 1
- Starting checkpoint: committed failing control fixture
- Ending checkpoint: bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb
- Control repair base: aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa
- Control repair head: bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb
- Control scope check: PASS control repair scope: aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa -> bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb; allowed=project-docs/tools/report-lint.sh (terminated exit_code: 0)
- Phase-branch integration: PASS fast-forwarded to bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb
- Outcome: RESOLVED
- Resolution evidence: exact reproducer passed
- Next checkpoint: REPORT
### 2026-08-15 10:03 TASK-RECOVERY task-1
- Source: TRG-1-001 resolved control and originating report
- Category: report-repair
- Repair mode: task-scope
- Repair authority: TRG-1-001 report continuation
- Triage disposition: N/A
- Control contract: N/A
- Control mismatch: N/A
- Allowed files: project-docs/docs/phases/phase-1/tasks/task-1-done.md
- Forbidden changes: product, task command, evidence, assertion, threshold, or control files
- Exact verification: project-docs/tools/report-lint.sh --self-test
- Originating report: project-docs/docs/phases/phase-1/tasks/task-1-done.md
- Control repair base: aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa
- Trigger signature: task-1:quality-gate-control-defect:literal-command
- Attempt count: 2
- Consecutive no-progress attempts: 0/2
- Checkpoint: control repair integrated at bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb
- Phase worktree: workspace/phase-1 on phase-1; ledger-only dirty
- Temporary repair worktree: N/A
- Hidden-index flags: none
- Resolution evidence: PENDING - rerun exact verification and finalize report
- Next action: rerun exact verification and update only the originating report
- Human attention required: no
### 2026-08-15 10:04 REPAIR task-1
- Trigger signature: task-1:quality-gate-control-defect:literal-command
- Attempt count: 2
- Starting checkpoint: control repair integrated at bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb
- Ending checkpoint: exact report finalized
- Control repair base: N/A
- Control repair head: N/A
- Control scope check: N/A for task-scope report continuation
- Phase-branch integration: N/A for task-scope
- Outcome: RESOLVED
- Resolution evidence: exact verification passed and report finalized
- Next checkpoint: COMPLETE
### 2026-08-15 10:05 TASK task-1 COMPLETE
- Report: project-docs/docs/phases/phase-1/tasks/task-1-done.md
- report-lint: PASS - OUTCOME COMPLETE
- Commit: bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb - verified, subject matched
- QA-plan commands beyond the task spec: none required
EOF
run_sequence "drill-18 scoped control repair retains originating task identity" ADVANCEABLE "$valid_control_repair"

# A protected BLOCK keeps its original verdict. A unique, decision-authorized
# amendment supplies the later control directive without rewriting that record.
amended_phase="$workspace/amended-control/phase-1"
mkdir -p "$amended_phase"
cat > "$amended_phase/decisions.md" <<'EOF'
# Phase 1 Decisions

### DEC-200: Authorize the bounded control repair
EOF
cat > "$amended_phase/triage-log.md" <<'EOF'
# Phase 1 Triage Log

## TRG-1-010 - protected control definition
- **Originating task:** task-1
- **Trigger signature:** task-1:quality-gate-control-defect:blocked-control
- **Repair category:** N/A
- **Repair authority:** N/A
- **Attempt count:** 1
- **Verdict:** BLOCK
- **Escalation:** ESC-1-010

## Control Directive Amendment For TRG-1-010
- **Governing TRG:** TRG-1-010
- **Governing ESC:** ESC-1-010
- **Resolution authority:** DEC-200
- **Verdict:** REPAIR / control-conformance
- **Originating task:** task-1
- **Control trigger signature:** task-1:quality-gate-control-defect:blocked-control
- **Control attempt count:** 2
- **Control contract:** exact commands remain opaque
- **Control mismatch:** exact reproducer rejects accepted output
- **Allowed files:** project-docs/tools/report-lint.sh
- **Forbidden changes:** product, specification, assertion, threshold, budget, or verifier
- **Exact verification:** project-docs/tools/report-lint.sh --self-test
- **Originating report:** project-docs/docs/phases/phase-1/tasks/task-1-done.md
- **Control repair base:** aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa
- **References:** DEC-200 and ESC-1-010
EOF
valid_amended_control="$amended_phase/run-ledger.md"
cat > "$valid_amended_control" <<'EOF'
### 2026-08-15 11:00 RUN-START
- Sequence contract: v2
- Recovery directive schema: v2
- Branch: phase-1
- ESC baseline (Open at run start): none
- Warm-up/autonomous-resume state: warm-up complete
- Starting task: task-1
### 2026-08-15 11:01 TRIAGE TRG-1-010
- Trigger: protected control definition requires authority
- Originating task: task-1
- Trigger signature: task-1:quality-gate-control-defect:blocked-control
- Attempt count: 1
- Repair category: N/A
- Repair authority: N/A
- Checkpoint: exact control reproducer failed
- Verdict: BLOCK (escalated as ESC-1-010)
- Resolution evidence: protected control definition requires a decision
- Triage log: unique phase triage entry and ESC confirmed
### 2026-08-15 11:02 STOP
- Cause: TRG-1-010 / ESC-1-010 requires protected authority
- Triage: TRG-1-010 (BLOCK, escalated as ESC-1-010)
- Protected category: gate/control definition authority
- Next Human or external-authority action: record a phase decision
### 2026-08-15 11:03 RUN-RESUME
- Sequence contract: v2
- Recovery directive schema: v2
- Prior run/session or STOP: TRG-1-010 / ESC-1-010 STOP
- Governing TRG: TRG-1-010
- Governing ESC: ESC-1-010
- Trigger signature: task-1:quality-gate-control-defect:blocked-control
- Attempt count: 2
- Checkpoint: DEC-200 and amendment recorded
- Resolution evidence: DEC-200 authorizes the exact amendment
- Resolution authority: DEC-200
- Continuing task: task-1
- Next recovery: task-1 | gate-repair | control-conformance | task-1:quality-gate-control-defect:blocked-control | 2
### 2026-08-15 11:04 TASK-RECOVERY task-1
- Source: TRG-1-010 / DEC-200
- Category: gate-repair
- Repair mode: control-conformance
- Repair authority: TRG-1-010 / DEC-200
- Triage disposition: REPAIR / control-conformance
- Trigger signature: task-1:quality-gate-control-defect:blocked-control
- Attempt count: 2
- Control contract: exact commands remain opaque
- Control mismatch: exact reproducer rejects accepted output
- Allowed files: project-docs/tools/report-lint.sh
- Forbidden changes: product, specification, assertion, threshold, budget, or verifier
- Originating report: project-docs/docs/phases/phase-1/tasks/task-1-done.md
- Exact verification: project-docs/tools/report-lint.sh --self-test
- Control repair base: aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa
- Consecutive no-progress attempts: 0/2
- Checkpoint: DEC-200 resolution checkpoint
- Phase worktree: workspace/phase-1 on phase-1; ledger-only dirty
- Temporary repair worktree: workspace/control-repair on control-repair at aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa
- Hidden-index flags: none
- Resolution evidence: PENDING - repair and rerun exact verification
- Next action: run trusted-base workflow-control-scope-lint.sh and exact verification
- Human attention required: no
EOF
run_sequence "drill-18 unique resolved-BLOCK amendment authorizes exact control recovery" REPAIR_PENDING "$valid_amended_control"

missing_amendment_phase="$workspace/missing-amendment/phase-1"
mkdir -p "$missing_amendment_phase"
cp "$amended_phase/decisions.md" "$missing_amendment_phase/decisions.md"
sed '/^## Control Directive Amendment For TRG-1-010$/,$d' \
  "$amended_phase/triage-log.md" > "$missing_amendment_phase/triage-log.md"
cp "$valid_amended_control" "$missing_amendment_phase/run-ledger.md"
run_invalid_sequence "drill-18 resolved BLOCK requires its unique control amendment" \
  "lacks one exact decision-authorized control amendment" \
  "$missing_amendment_phase/run-ledger.md"

duplicate_amendment_phase="$workspace/duplicate-amendment/phase-1"
mkdir -p "$duplicate_amendment_phase"
cp "$amended_phase/decisions.md" "$duplicate_amendment_phase/decisions.md"
cp "$amended_phase/triage-log.md" "$duplicate_amendment_phase/triage-log.md"
sed -n '/^## Control Directive Amendment For TRG-1-010$/,$p' \
  "$amended_phase/triage-log.md" >> "$duplicate_amendment_phase/triage-log.md"
cp "$valid_amended_control" "$duplicate_amendment_phase/run-ledger.md"
run_invalid_sequence "drill-18 duplicate control amendments are ambiguous" \
  "lacks one exact decision-authorized control amendment" \
  "$duplicate_amendment_phase/run-ledger.md"

trailing_directive_phase="$workspace/trailing-directive/phase-1"
mkdir -p "$trailing_directive_phase"
sed 's/^- \*\*Control contract:\*\* exact commands are compared byte-for-byte$/&   /' \
  "$control_phase/triage-log.md" > "$trailing_directive_phase/triage-log.md"
cp "$valid_control_repair" "$trailing_directive_phase/run-ledger.md"
run_invalid_sequence "drill-18 authoritative directive whitespace is byte-exact" \
  "does not match its authoritative TRIAGE record" \
  "$trailing_directive_phase/run-ledger.md"

duplicate_trg_phase="$workspace/duplicate-trg/phase-1"
mkdir -p "$duplicate_trg_phase"
cp "$control_phase/triage-log.md" "$duplicate_trg_phase/triage-log.md"
cat >> "$duplicate_trg_phase/triage-log.md" <<'EOF'

## TRG-1-001 - conflicting duplicate
- **Verdict:** BLOCK
EOF
cp "$valid_control_repair" "$duplicate_trg_phase/run-ledger.md"
run_invalid_sequence "drill-18 duplicate external TRG identity is rejected" \
  "TRIAGE ID is not recorded in the phase triage log" \
  "$duplicate_trg_phase/run-ledger.md"

nonexistent_contract_decision_phase="$workspace/nonexistent-contract-decision/phase-1"
mkdir -p "$nonexistent_contract_decision_phase"
sed 's/exact commands are compared byte-for-byte/DEC-999999 says exact commands are compared byte-for-byte/' \
  "$control_phase/triage-log.md" > "$nonexistent_contract_decision_phase/triage-log.md"
sed 's/exact commands are compared byte-for-byte/DEC-999999 says exact commands are compared byte-for-byte/g' \
  "$valid_control_repair" > "$nonexistent_contract_decision_phase/run-ledger.md"
run_sequence "drill-18 non-authority directive prose does not become decision authority" \
  ADVANCEABLE "$nonexistent_contract_decision_phase/run-ledger.md"

nonexistent_repair_authority_phase="$workspace/nonexistent-repair-authority/phase-1"
mkdir -p "$nonexistent_repair_authority_phase"
sed 's/TRG-1-001 typed directive/TRG-1-001 typed directive \/ DEC-999999/g' \
  "$control_phase/triage-log.md" > "$nonexistent_repair_authority_phase/triage-log.md"
sed 's/TRG-1-001 typed directive/TRG-1-001 typed directive \/ DEC-999999/g' \
  "$valid_control_repair" > "$nonexistent_repair_authority_phase/run-ledger.md"
run_invalid_sequence "drill-18 nonexistent DEC cannot enter Repair authority" \
  "cites an unrecorded or malformed phase decision" \
  "$nonexistent_repair_authority_phase/run-ledger.md"

embedded_repair_authority_phase="$workspace/embedded-repair-authority/phase-1"
mkdir -p "$embedded_repair_authority_phase"
sed 's/TRG-1-001 typed directive/TRG-1-001 typed directive \/ NOTDEC-999X/g' \
  "$control_phase/triage-log.md" > "$embedded_repair_authority_phase/triage-log.md"
sed 's/TRG-1-001 typed directive/TRG-1-001 typed directive \/ NOTDEC-999X/g' \
  "$valid_control_repair" > "$embedded_repair_authority_phase/run-ledger.md"
run_invalid_sequence "drill-18 embedded pseudo-ID cannot enter Repair authority" \
  "cites an unrecorded or malformed phase decision" \
  "$embedded_repair_authority_phase/run-ledger.md"

duplicate_report_continuation="$control_phase/duplicate-report-continuation.md"
sed '/^### 2026-08-15 10:03 TASK-RECOVERY task-1$/,/^### 2026-08-15 10:04 REPAIR task-1$/ {
  /- Originating report: project-docs\/docs\/phases\/phase-1\/tasks\/task-1-done.md/a\
- Originating report: cmd/product.go
}' "$valid_control_repair" > "$duplicate_report_continuation"
run_invalid_sequence "drill-18 report continuation rejects duplicate report authority" \
  "TASK-RECOVERY must contain each typed field exactly once" \
  "$duplicate_report_continuation"

duplicate_stop_triage="$amended_phase/duplicate-stop-triage.md"
sed '/^- Triage: TRG-1-010 (BLOCK, escalated as ESC-1-010)$/i\
- Triage: TRG-1-999 (BLOCK, escalated as ESC-1-999)
' "$valid_amended_control" > "$duplicate_stop_triage"
run_invalid_sequence "drill-18 STOP rejects duplicate governing triage fields" \
  "STOP must contain each authority field exactly once" "$duplicate_stop_triage"

wrong_scope_allowlist="$control_phase/wrong-scope-allowlist.md"
sed 's/; allowed=project-docs\/tools\/report-lint.sh (terminated/; allowed=project-docs\/tools\/workflow-lint.sh (terminated/' \
  "$valid_control_repair" > "$wrong_scope_allowlist"
run_invalid_sequence "drill-18 scope receipt must repeat the authoritative allowlist" \
  "exact allowlist receipt" "$wrong_scope_allowlist"

protected_verifier_phase="$workspace/protected-verifier/phase-1"
mkdir -p "$protected_verifier_phase"
sed 's#project-docs/tools/report-lint.sh#project-docs/tools/workflow-sequence-lint.sh#g' \
  "$control_phase/triage-log.md" > "$protected_verifier_phase/triage-log.md"
sed 's#project-docs/tools/report-lint.sh#project-docs/tools/workflow-sequence-lint.sh#g' \
  "$valid_control_repair" > "$protected_verifier_phase/run-ledger.md"
run_invalid_sequence "drill-18 control repair cannot authorize a protected verifier" \
  "cannot authorize a protected workflow verifier" \
  "$protected_verifier_phase/run-ledger.md"

wrong_control_disposition="$control_phase/wrong-control-disposition.md"
sed '0,/^- Triage disposition: REPAIR \/ control-conformance$/s//- Triage disposition: N\/A/' \
  "$valid_control_repair" > "$wrong_control_disposition"
run_invalid_sequence "drill-18 TASK-RECOVERY disposition must match control mode" \
  "Triage disposition does not match Repair mode" "$wrong_control_disposition"

blank_run_start="$workspace/blank-v2-run-start.md"
cat > "$blank_run_start" <<'EOF'
### 2026-08-15 12:00 RUN-START
- Sequence contract: v2
- Recovery directive schema: v2
EOF
run_invalid_sequence "drill-18 schema v2 RUN-START requires its checkpoint" \
  "RUN-START must contain each checkpoint field exactly once" "$blank_run_start"

normal_resume="$workspace/normal-v2-resume.md"
cat > "$normal_resume" <<'EOF'
### 2026-08-15 12:00 RUN-RESUME
- Sequence contract: v2
- Recovery directive schema: v2
- Prior run/session or STOP: prior Phase Runner session ended cleanly
- Governing TRG: N/A
- Governing ESC: N/A
- Trigger signature: N/A
- Attempt count: N/A
- Checkpoint: task-17 complete; Review 6 next
- Resolution evidence: committed reports and ledger reconcile
- Resolution authority: N/A
- Continuing task: task-17
- Next recovery: N/A
EOF
run_sequence "drill-18 normal schema v2 RUN-RESUME uses explicit N/A authority" \
  ADVANCEABLE "$normal_resume"

normal_resume_with_dec_prose="$workspace/normal-v2-resume-with-dec-prose.md"
sed -e 's/^- Checkpoint: task-17 complete; Review 6 next$/- Checkpoint: retained NOTDEC-999X diagnostic text; task-17 complete; Review 6 next/' \
    -e 's/^- Resolution evidence: committed reports and ledger reconcile$/- Resolution evidence: diagnostic command rg -n "DEC-[0-9]+" exited as expected; DEC-999999 is reproduced prose, not authority/' \
  "$normal_resume" > "$normal_resume_with_dec_prose"
run_sequence "drill-18 checkpoint and command prose do not become decision authority" \
  ADVANCEABLE "$normal_resume_with_dec_prose"

missing_normal_resume_field="$workspace/missing-normal-resume-field.md"
sed '/^- Prior run\/session or STOP: /d' "$normal_resume" > "$missing_normal_resume_field"
run_invalid_sequence "drill-18 normal RUN-RESUME cannot omit reconciliation fields" \
  "RUN-RESUME must contain each reconciliation field exactly once" \
  "$missing_normal_resume_field"

wrong_normal_resume_authority="$workspace/wrong-normal-resume-authority.md"
sed 's/^- Resolution authority: N\/A$/- Resolution authority: restored-capability:invented/' \
  "$normal_resume" > "$wrong_normal_resume_authority"
run_invalid_sequence "drill-18 normal RUN-RESUME cannot imply STOP authority" \
  "normal schema v2 RUN-RESUME must use N/A" "$wrong_normal_resume_authority"

v1_append_ledger="$workspace/v1-run-append.md"
cp "$normal_resume" "$v1_append_ledger"
cp "$v1_append_ledger" "$v1_append_ledger.before"
if v1_append_output=$(cat <<'EOF' | "$sequence_lint" "$v1_append_ledger" --append-stdin 2>&1
### 2026-08-15 12:01 RUN-RESUME
- Sequence contract: v2
- Recovery directive schema: v1
- Prior run/session or STOP: prior Runner
- Governing TRG: N/A
- Governing ESC: N/A
- Trigger signature: N/A
- Attempt count: N/A
- Checkpoint: unchanged
- Resolution evidence: unchanged
- Resolution authority: N/A
- Continuing task: task-17
- Next recovery: N/A
EOF
); then
  v1_append_code=0
else
  v1_append_code=$?
fi
if [ "$v1_append_code" -ne 0 ] && \
   printf '%s\n' "$v1_append_output" | grep -qF 'every newly appended RUN' && \
   cmp -s "$v1_append_ledger.before" "$v1_append_ledger"; then
  drill_pass "drill-18 transactional append forces every new RUN to schema v2"
else
  drill_fail "drill-18 transactional RUN schema adoption (exit=$v1_append_code, output=$v1_append_output)"
fi

canonical_resume_candidate="$workspace/canonical-resume-candidate.md"
cat > "$canonical_resume_candidate" <<'EOF'
### 2026-08-15 12:02 RUN-RESUME
- Sequence contract: v2
- Recovery directive schema: v2
- Prior run/session or STOP: prior Runner
- Governing TRG: N/A
- Governing ESC: N/A
- Trigger signature: N/A
- Attempt count: N/A
- Checkpoint: unchanged
- Resolution evidence: unchanged
- Resolution authority: N/A
- Continuing task: task-17
- Next recovery: N/A
EOF

trailing_v1_run="$workspace/trailing-v1-run-candidate.md"
sed -e '1s/$/ trailing-token/' \
    -e 's/^- Recovery directive schema: v2$/- Recovery directive schema: v1/' \
  "$canonical_resume_candidate" > "$trailing_v1_run"
run_invalid_append "drill-18 trailing RUN token cannot bypass append-time v2 adoption" \
  "every newly appended RUN" "$v1_append_ledger" "$trailing_v1_run"

trailing_v2_run="$workspace/trailing-v2-run-candidate.md"
sed '1s/$/ trailing-token/' "$canonical_resume_candidate" > "$trailing_v2_run"
run_invalid_append "drill-18 schema v2 RUN rejects a trailing heading token" \
  "must use its canonical heading shape" "$v1_append_ledger" "$trailing_v2_run"

whitespace_v2_run="$workspace/whitespace-v2-run-candidate.md"
sed '1s/$/ /' "$canonical_resume_candidate" > "$whitespace_v2_run"
run_invalid_append "drill-18 schema v2 RUN rejects trailing heading whitespace" \
  "must use its canonical heading shape" "$v1_append_ledger" "$whitespace_v2_run"

ordinary_phase="$workspace/ordinary-triage/phase-1"
mkdir -p "$ordinary_phase"
cat > "$ordinary_phase/triage-log.md" <<'EOF'
# Phase 1 Triage Log

## TRG-1-030 - bounded ordinary repair
- **Originating task:** task-1
- **Trigger signature:** task-1:focused-test:legacy-assertion
- **Repair category:** gate-repair
- **Repair authority:** TRG-1-030 task spec
- **Attempt count:** 1
- **Verdict:** REPAIR
EOF
ordinary_repair="$ordinary_phase/run-ledger.md"
cat > "$ordinary_repair" <<'EOF'
### 2026-08-15 12:10 RUN-START
- Sequence contract: v2
- Recovery directive schema: v2
- Branch: phase-1
- ESC baseline (Open at run start): none
- Warm-up/autonomous-resume state: warm-up complete
- Starting task: task-1
### 2026-08-15 12:11 TRIAGE TRG-1-030
- Trigger: focused test failed on a legacy assertion
- Originating task: task-1
- Trigger signature: task-1:focused-test:legacy-assertion
- Attempt count: 1
- Repair category: gate-repair
- Repair authority: TRG-1-030 task spec
- Checkpoint: focused test failed on in-scope legacy assertion
- Verdict: REPAIR
- Resolution evidence: bounded same-task repair is authorized
- Triage log: unique phase triage entry confirmed
### 2026-08-15 12:12 TASK-RECOVERY task-1
- Source: TRG-1-030
- Category: gate-repair
- Repair mode: task-scope
- Repair authority: TRG-1-030 task spec
- Triage disposition: N/A
- Trigger signature: task-1:focused-test:legacy-assertion
- Attempt count: 1
- Control contract: N/A
- Control mismatch: N/A
- Allowed files: ordinary task scope
- Forbidden changes: protected decisions and out-of-scope files
- Originating report: N/A
- Exact verification: N/A
- Control repair base: N/A
- Consecutive no-progress attempts: 0/2
- Checkpoint: focused test failure preserved
- Phase worktree: workspace/phase-1 on phase-1; ledger-only dirty
- Temporary repair worktree: N/A
- Hidden-index flags: none
- Resolution evidence: PENDING - repair and rerun the exact test
- Next action: same-task repair and exact test rerun
- Human attention required: no
EOF
run_sequence "drill-18 schema v2 ordinary TRIAGE binds exact recovery identity" \
  REPAIR_PENDING "$ordinary_repair"

misbound_ordinary_phase="$workspace/misbound-ordinary/phase-1"
mkdir -p "$misbound_ordinary_phase"
sed -e 's/task-1:focused-test:legacy-assertion/task-9:focused-test:other/' \
    -e 's/\*\*Attempt count:\*\* 1/**Attempt count:** 9/' \
  "$ordinary_phase/triage-log.md" > "$misbound_ordinary_phase/triage-log.md"
cp "$ordinary_repair" "$misbound_ordinary_phase/run-ledger.md"
run_invalid_sequence "drill-18 ordinary TRIAGE cannot redirect external task trigger or attempt" \
  "identity does not match its unique external record" \
  "$misbound_ordinary_phase/run-ledger.md"

capability_phase="$workspace/restored-capability/phase-1"
mkdir -p "$capability_phase"
cat > "$capability_phase/triage-log.md" <<'EOF'
# Phase 1 Triage Log

## TRG-1-040 - required local capability unavailable
- **Originating task:** task-1
- **Trigger signature:** task-1:local-environment:required-tool
- **Repair category:** N/A
- **Repair authority:** N/A
- **Attempt count:** 1
- **Verdict:** BLOCK
- **Escalation:** ESC-1-040
EOF
restored_capability="$capability_phase/run-ledger.md"
cat > "$restored_capability" <<'EOF'
### 2026-08-15 12:20 RUN-START
- Sequence contract: v2
- Recovery directive schema: v2
- Branch: phase-1
- ESC baseline (Open at run start): none
- Warm-up/autonomous-resume state: warm-up complete
- Starting task: task-1
### 2026-08-15 12:21 TRIAGE TRG-1-040
- Trigger: required local capability unavailable
- Originating task: task-1
- Trigger signature: task-1:local-environment:required-tool
- Attempt count: 1
- Repair category: N/A
- Repair authority: N/A
- Checkpoint: required tool unavailable
- Verdict: BLOCK (escalated as ESC-1-040)
- Resolution evidence: external capability must be restored
- Triage log: unique phase triage entry and ESC confirmed
### 2026-08-15 12:22 STOP
- Cause: TRG-1-040 / ESC-1-040 required capability unavailable
- Triage: TRG-1-040 (BLOCK, escalated as ESC-1-040)
- Protected category: external capability unavailable
- Next Human or external-authority action: restore the required tool capability
### 2026-08-15 12:23 RUN-RESUME
- Sequence contract: v2
- Recovery directive schema: v2
- Prior run/session or STOP: TRG-1-040 / ESC-1-040 STOP
- Governing TRG: TRG-1-040
- Governing ESC: ESC-1-040
- Trigger signature: task-1:local-environment:required-tool
- Attempt count: 2
- Checkpoint: required tool is available and its probe passed
- Resolution evidence: capability probe terminated exit_code: 0
- Resolution authority: restored-capability:required-tool-v1
- Continuing task: task-1
- Next recovery: task-1 | local-environment | task-scope | task-1:local-environment:required-tool | 2
### 2026-08-15 12:24 TASK-RECOVERY task-1
- Source: TRG-1-040 / restored-capability:required-tool-v1
- Category: local-environment
- Repair mode: task-scope
- Repair authority: TRG-1-040 / restored-capability:required-tool-v1
- Triage disposition: N/A
- Trigger signature: task-1:local-environment:required-tool
- Attempt count: 2
- Control contract: N/A
- Control mismatch: N/A
- Allowed files: ordinary task scope
- Forbidden changes: protected decisions and out-of-scope files
- Originating report: N/A
- Exact verification: N/A
- Control repair base: N/A
- Consecutive no-progress attempts: 0/2
- Checkpoint: restored capability checkpoint
- Phase worktree: workspace/phase-1 on phase-1; ledger-only dirty
- Temporary repair worktree: N/A
- Hidden-index flags: none
- Resolution evidence: PENDING - resume the exact interrupted task gate
- Next action: rerun the exact task gate
- Human attention required: no
EOF
run_sequence "drill-18 restored capability can resume an unchanged protected contract" \
  REPAIR_PENDING "$restored_capability"

invalid_capability_authority="$capability_phase/invalid-capability-authority.md"
sed 's/restored-capability:required-tool-v1/restored by someone/g' \
  "$restored_capability" > "$invalid_capability_authority"
run_invalid_sequence "drill-18 STOP resume requires a typed resolution authority" \
  "invalid or unverifiable Resolution authority" "$invalid_capability_authority"

expanded_capability_authority="$capability_phase/expanded-capability-authority.md"
sed '/^### 2026-08-15 12:24 TASK-RECOVERY task-1$/,$ s/restored-capability:required-tool-v1/restored-capability:required-tool-v1-expanded/g' \
  "$restored_capability" > "$expanded_capability_authority"
run_invalid_sequence "drill-18 restored-capability authority cannot be suffix-expanded" \
  "does not retain its typed RUN-RESUME resolution authority" \
  "$expanded_capability_authority"

verified_commit="$capability_phase/verified-commit.md"
sed 's/restored-capability:required-tool-v1/verified-commit:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa/g' \
  "$restored_capability" > "$verified_commit"
run_sequence "drill-18 verified commit can resume an unchanged protected contract" \
  REPAIR_PENDING "$verified_commit"

expanded_verified_commit="$capability_phase/expanded-verified-commit.md"
sed '/^### 2026-08-15 12:24 TASK-RECOVERY task-1$/,$ s/verified-commit:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa/verified-commit:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa-expanded/g' \
  "$verified_commit" > "$expanded_verified_commit"
run_invalid_sequence "drill-18 verified-commit authority cannot be suffix-expanded" \
  "does not retain its typed RUN-RESUME resolution authority" \
  "$expanded_verified_commit"

expanded_decision_authority="$amended_phase/expanded-decision-authority.md"
sed '/^### 2026-08-15 11:04 TASK-RECOVERY task-1$/,$ s/DEC-200/DEC-200-expanded/g' \
  "$valid_amended_control" > "$expanded_decision_authority"
run_invalid_sequence "drill-18 decision authority cannot be suffix-expanded" \
  "cites an unrecorded or malformed phase decision" \
  "$expanded_decision_authority"

duplicate_schema="$control_phase/duplicate-schema.md"
sed '0,/^- Recovery directive schema: v2$/s//&\
- Recovery directive schema: v2/' "$valid_control_repair" > "$duplicate_schema"
run_invalid_sequence "drill-18 schema v2 marker must be singular" "exactly one Recovery directive schema field" "$duplicate_schema"

duplicate_directive_field="$control_phase/duplicate-directive-field.md"
sed '0,/^- Allowed files: project-docs\/tools\/report-lint.sh$/s//&\
- Allowed files: project-docs\/tools\/report-lint.sh/' "$valid_control_repair" > "$duplicate_directive_field"
run_invalid_sequence "drill-18 schema v2 directive fields must be singular" "directive fields must each occur exactly once" "$duplicate_directive_field"

missing_triage_trigger="$control_phase/missing-triage-trigger.md"
sed '/^- Trigger: report-lint violates its byte-exact command contract$/d' \
  "$valid_control_repair" > "$missing_triage_trigger"
run_invalid_sequence "drill-18 schema v2 TRIAGE requires its trigger" \
  "TRIAGE must contain each authority field exactly once" "$missing_triage_trigger"

duplicate_triage_trigger="$control_phase/duplicate-triage-trigger.md"
sed '/^- Trigger: report-lint violates its byte-exact command contract$/a\
- Trigger: conflicting duplicate trigger
' "$valid_control_repair" > "$duplicate_triage_trigger"
run_invalid_sequence "drill-18 schema v2 TRIAGE rejects duplicate triggers" \
  "TRIAGE must contain each authority field exactly once" "$duplicate_triage_trigger"

premature_control_complete="$control_phase/premature-control-complete.md"
sed '/^### 2026-08-15 10:03 TASK-RECOVERY task-1$/,/^### 2026-08-15 10:05 TASK task-1 COMPLETE$/ {
  /^### 2026-08-15 10:05 TASK task-1 COMPLETE$/!d
}' "$valid_control_repair" > "$premature_control_complete"
run_invalid_sequence "drill-18 resolved control cannot complete before report continuation" "task completed before a resolved REPAIR event" "$premature_control_complete"

wrong_report_continuation="$control_phase/wrong-report-continuation.md"
sed '/^### 2026-08-15 10:03 TASK-RECOVERY task-1$/,/^### 2026-08-15 10:04 REPAIR task-1$/ s#phase-1/tasks/task-1-done.md#phase-1/tasks/task-2-done.md#' \
  "$valid_control_repair" > "$wrong_report_continuation"
run_invalid_sequence "drill-18 control continuation cannot redirect its report" "linked exact report continuation" "$wrong_report_continuation"

wrong_report_verification="$control_phase/wrong-report-verification.md"
sed '/^### 2026-08-15 10:03 TASK-RECOVERY task-1$/,/^### 2026-08-15 10:04 REPAIR task-1$/ s#report-lint.sh --self-test#workflow-lint.sh#' \
  "$valid_control_repair" > "$wrong_report_verification"
run_invalid_sequence "drill-18 control continuation preserves exact verification" "schema v2 report continuation changed" "$wrong_report_verification"

control_escape="$control_phase/control-escape.md"
sed '0,/Outcome: RESOLVED/s//Outcome: MATERIAL-PROGRESS/' \
  "$valid_control_repair" > "$control_escape"
run_invalid_sequence "drill-18 non-resolved control cannot escape to task scope" "cannot change Repair mode" "$control_escape"

mutated_control_retry="$control_phase/mutated-control-retry.md"
sed -e '0,/Outcome: RESOLVED/s//Outcome: MATERIAL-PROGRESS/' \
    -e '/^### 2026-08-15 10:03 TASK-RECOVERY task-1$/,/^### 2026-08-15 10:04 REPAIR task-1$/ s/Repair mode: task-scope/Repair mode: control-conformance/' \
    "$valid_control_repair" > "$mutated_control_retry"
run_invalid_sequence "drill-18 control retry cannot replace its directive" "changed its authoritative directive" "$mutated_control_retry"

valid_control_retry="$control_phase/valid-control-retry.md"
sed '/^### 2026-08-15 10:03 TASK-RECOVERY task-1$/,$d; 0,/Outcome: RESOLVED/s//Outcome: MATERIAL-PROGRESS/' \
  "$valid_control_repair" > "$valid_control_retry"
cat >> "$valid_control_retry" <<'EOF'
### 2026-08-15 10:03 TASK-RECOVERY task-1
- Source: TRG-1-001 material-progress control retry
- Category: report-repair
- Repair mode: control-conformance
- Repair authority: TRG-1-001 typed directive
- Triage disposition: REPAIR / control-conformance
- Control contract: exact commands are compared byte-for-byte
- Control mismatch: exact reproducer false-rejects a literal backslash-dot
- Allowed files: project-docs/tools/report-lint.sh
- Forbidden changes: command, assertion, accepted status, scope, threshold, budget, or tolerance
- Exact verification: project-docs/tools/report-lint.sh --self-test
- Originating report: project-docs/docs/phases/phase-1/tasks/task-1-done.md
- Control repair base: bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb
- Trigger signature: task-1:quality-gate-control-defect:literal-command
- Attempt count: 2
- Consecutive no-progress attempts: 0/2
- Checkpoint: first isolated control commit integrated
- Phase worktree: workspace/phase-1 on phase-1; ledger-only dirty
- Temporary repair worktree: workspace/control-retry on control-retry at bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb
- Hidden-index flags: none
- Resolution evidence: PENDING - finish the same bounded control repair
- Next action: repair the same control, run workflow-control-scope-lint.sh, and rerun exact verification
- Human attention required: no
### 2026-08-15 10:04 REPAIR task-1
- Trigger signature: task-1:quality-gate-control-defect:literal-command
- Attempt count: 2
- Starting checkpoint: first isolated control commit integrated
- Ending checkpoint: cccccccccccccccccccccccccccccccccccccccc
- Control repair base: bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb
- Control repair head: cccccccccccccccccccccccccccccccccccccccc
- Control scope check: PASS control repair scope: bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb -> cccccccccccccccccccccccccccccccccccccccc; allowed=project-docs/tools/report-lint.sh (terminated exit_code: 0)
- Phase-branch integration: PASS fast-forwarded to cccccccccccccccccccccccccccccccccccccccc
- Outcome: RESOLVED
- Resolution evidence: exact reproducer passed
- Next checkpoint: REPORT
### 2026-08-15 10:05 TASK-RECOVERY task-1
- Source: TRG-1-001 resolved control and originating report
- Category: report-repair
- Repair mode: task-scope
- Repair authority: TRG-1-001 report continuation
- Triage disposition: N/A
- Control contract: N/A
- Control mismatch: N/A
- Allowed files: project-docs/docs/phases/phase-1/tasks/task-1-done.md
- Forbidden changes: product, task command, evidence, assertion, threshold, or control files
- Exact verification: project-docs/tools/report-lint.sh --self-test
- Originating report: project-docs/docs/phases/phase-1/tasks/task-1-done.md
- Control repair base: bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb
- Trigger signature: task-1:quality-gate-control-defect:literal-command
- Attempt count: 3
- Consecutive no-progress attempts: 0/2
- Checkpoint: control repair integrated at cccccccccccccccccccccccccccccccccccccccc
- Phase worktree: workspace/phase-1 on phase-1; ledger-only dirty
- Temporary repair worktree: N/A
- Hidden-index flags: none
- Resolution evidence: PENDING - rerun exact verification and finalize report
- Next action: rerun exact verification and update only the originating report
- Human attention required: no
### 2026-08-15 10:06 REPAIR task-1
- Trigger signature: task-1:quality-gate-control-defect:literal-command
- Attempt count: 3
- Starting checkpoint: control repair integrated at cccccccccccccccccccccccccccccccccccccccc
- Ending checkpoint: exact report finalized
- Control repair base: N/A
- Control repair head: N/A
- Control scope check: N/A for task-scope report continuation
- Phase-branch integration: N/A for task-scope
- Outcome: RESOLVED
- Resolution evidence: exact verification passed and report finalized
- Next checkpoint: COMPLETE
### 2026-08-15 10:07 TASK task-1 COMPLETE
- Report: project-docs/docs/phases/phase-1/tasks/task-1-done.md
- report-lint: PASS - OUTCOME COMPLETE
- Commit: cccccccccccccccccccccccccccccccccccccccc - verified, subject matched
- QA-plan commands beyond the task spec: none required
EOF
run_sequence "drill-18 material-progress control retry preserves directive and advances base" ADVANCEABLE "$valid_control_retry"

expanded_control_retry="$control_phase/expanded-control-retry.md"
sed '/^### 2026-08-15 10:03 TASK-RECOVERY task-1$/,/^### 2026-08-15 10:04 REPAIR task-1$/ s#Allowed files: project-docs/tools/report-lint.sh#Allowed files: project-docs/tools/report-lint.sh project-docs/tools/workflow-lint.sh#' \
  "$valid_control_retry" > "$expanded_control_retry"
run_invalid_sequence "drill-18 material-progress retry cannot expand control allowlist" "changed its authoritative directive" "$expanded_control_retry"

schema_omitted="$control_phase/schema-omitted.md"
sed '/^- Recovery directive schema: /d' "$valid_continuation" > "$schema_omitted"
run_invalid_sequence "drill-18 fresh v2 ledger cannot omit directive schema" "without a Recovery directive schema" "$schema_omitted"

untriaged_control_repair="$control_phase/untriaged-control-repair.md"
sed '/^### 2026-08-15 10:00 TRIAGE TRG-1-001$/,/^- Triage log: /d' "$valid_control_repair" > "$untriaged_control_repair"
run_invalid_sequence "drill-18 control repair requires preceding triage authority" "has no preceding TRIAGE authority" "$untriaged_control_repair"

self_authorized_control_repair="$control_phase/self-authorized-control-repair.md"
sed 's/Repair authority: TRG-1-001/Repair authority: SELF-AUTHORIZED/' "$valid_control_repair" > "$self_authorized_control_repair"
run_invalid_sequence "drill-18 control repair authority binds the latest triage" "does not reference latest TRIAGE TRG-1-001" "$self_authorized_control_repair"

fabricated_triage_id="$control_phase/fabricated-triage-id.md"
sed 's/TRG-1-001/SELF-AUTHORIZED/g' "$valid_control_repair" > "$fabricated_triage_id"
run_invalid_sequence "drill-18 triage authority requires a concrete TRG ID" "TRIAGE has missing or invalid TRG ID" "$fabricated_triage_id"

incomplete_triage="$control_phase/incomplete-triage.md"
sed '/^- Triage log: /d' "$valid_control_repair" > "$incomplete_triage"
run_invalid_sequence "drill-18 control authority requires a complete triage record" "TRIAGE is missing Triage log" "$incomplete_triage"

unrecorded_phase="$workspace/unrecorded-control/phase-1"
mkdir -p "$unrecorded_phase"
unrecorded_triage="$unrecorded_phase/run-ledger.md"
cp "$valid_control_repair" "$unrecorded_triage"
run_invalid_sequence "drill-18 control authority must exist in the phase triage log" "TRIAGE ID is not recorded in the phase triage log" "$unrecorded_triage"

misbound_phase="$workspace/misbound-control/phase-1"
mkdir -p "$misbound_phase"
cat > "$misbound_phase/triage-log.md" <<'EOF'
# Phase 1 Triage Log

## TRG-1-001 - ordinary repair only
- **Verdict:** REPAIR
EOF
misbound_triage="$misbound_phase/run-ledger.md"
cp "$valid_control_repair" "$misbound_triage"
run_invalid_sequence "drill-18 ledger verdict must match the recorded triage verdict" "TRIAGE verdict does not match the phase triage log" "$misbound_triage"

prefix_phase="$workspace/prefix-control/phase-1"
mkdir -p "$prefix_phase"
sed 's/REPAIR \/ control-conformance$/REPAIR \/ control-conformance REVOKED/' \
  "$control_phase/triage-log.md" > "$prefix_phase/triage-log.md"
prefix_verdict="$prefix_phase/run-ledger.md"
cp "$valid_control_repair" "$prefix_verdict"
run_invalid_sequence "drill-18 triage verdict is exact rather than prefix matched" "TRIAGE verdict does not match the phase triage log" "$prefix_verdict"

directive_phase="$workspace/directive-control/phase-1"
mkdir -p "$directive_phase"
cp "$control_phase/triage-log.md" "$directive_phase/triage-log.md"
directive_mismatch="$directive_phase/run-ledger.md"
sed '0,/Allowed files: project-docs\/tools\/report-lint.sh/s//Allowed files: project-docs\/tools\/workflow-lint.sh/' \
  "$valid_control_repair" > "$directive_mismatch"
run_invalid_sequence "drill-18 control directive must match authoritative triage bytes" "does not match its authoritative TRIAGE record" "$directive_mismatch"

ordinary_triage_control="$control_phase/ordinary-triage-control.md"
sed 's/Verdict: REPAIR \/ control-conformance/Verdict: REPAIR/' "$valid_control_repair" > "$ordinary_triage_control"
run_invalid_sequence "drill-18 ordinary REPAIR cannot authorize cross-scope control edits" "not REPAIR / control-conformance or resolved BLOCK" "$ordinary_triage_control"

unbound_attempt="$control_phase/unbound-control-attempt.md"
sed '0,/Attempt count: 1/! s/Attempt count: 1/Attempt count: 99/g' "$valid_control_repair" > "$unbound_attempt"
run_invalid_sequence "drill-18 control attempt remains within triaged budget" "Attempt count exceeds latest TRIAGE authority" "$unbound_attempt"

unbudgeted_control_repair="$control_phase/unbudgeted-control-repair.md"
sed '/^- Consecutive no-progress attempts: /d' "$valid_control_repair" > "$unbudgeted_control_repair"
run_invalid_sequence "drill-18 control repair cannot discard its recovery budget" "missing or invalid Consecutive no-progress attempts" "$unbudgeted_control_repair"

missing_control_contract="$control_phase/missing-control-contract.md"
sed '/^- Control contract: /d' "$valid_control_repair" > "$missing_control_contract"
run_invalid_sequence "drill-18 control repair requires accepted contract authority" "missing Control contract" "$missing_control_contract"

missing_repair_authority="$control_phase/missing-repair-authority.md"
sed '/^- Repair authority: /d' "$valid_control_repair" > "$missing_repair_authority"
run_invalid_sequence "drill-18 marked recovery requires recorded authority" "missing Repair authority" "$missing_repair_authority"

missing_control_mismatch="$control_phase/missing-control-mismatch.md"
sed '/^- Control mismatch: /d' "$valid_control_repair" > "$missing_control_mismatch"
run_invalid_sequence "drill-18 control repair requires exact mismatch evidence" "missing Control mismatch" "$missing_control_mismatch"

missing_control_files="$control_phase/missing-control-files.md"
sed '/^- Allowed files: /d' "$valid_control_repair" > "$missing_control_files"
run_invalid_sequence "drill-18 control repair requires bounded file authority" "missing Allowed files" "$missing_control_files"

missing_control_guard="$control_phase/missing-control-guard.md"
sed '/^- Forbidden changes: /d' "$valid_control_repair" > "$missing_control_guard"
run_invalid_sequence "drill-18 control repair requires semantic guardrails" "missing Forbidden changes" "$missing_control_guard"

missing_repair_mode="$control_phase/missing-repair-mode.md"
sed '/^- Repair mode: /d' "$valid_control_repair" > "$missing_repair_mode"
run_invalid_sequence "drill-18 marked recovery requires an explicit repair mode" "missing or invalid Repair mode" "$missing_repair_mode"

invalid_directive_schema="$control_phase/invalid-directive-schema.md"
sed 's/Recovery directive schema: v2/Recovery directive schema: v3/' "$valid_control_repair" > "$invalid_directive_schema"
run_invalid_sequence "drill-18 unknown recovery directive schema is rejected" "invalid Recovery directive schema" "$invalid_directive_schema"

missing_sequence_contract="$control_phase/missing-control-sequence-contract.md"
sed '/^- Sequence contract: /d' "$valid_control_repair" > "$missing_sequence_contract"
run_invalid_sequence "drill-18 directive schema requires sequence contract" "missing Sequence contract: v2" "$missing_sequence_contract"

schema_dropped_on_resume="$control_phase/schema-dropped-on-resume.md"
cat "$valid_control_repair" > "$schema_dropped_on_resume"
cat >> "$schema_dropped_on_resume" <<'EOF'
### 2026-08-15 10:04 RUN-RESUME
- Sequence contract: v2
EOF
run_invalid_sequence "drill-18 marked ledger cannot drop directive schema on resume" "missing Recovery directive schema" "$schema_dropped_on_resume"

delayed_schema="$control_phase/delayed-schema.md"
cat > "$delayed_schema" <<'EOF'
### 2026-08-15 10:00 RUN-START
- Sequence contract: v2
### 2026-08-15 10:01 TASK-RECOVERY task-1
### 2026-08-15 10:02 REPAIR task-1
- Outcome: RESOLVED
### 2026-08-15 10:03 TASK task-1 COMPLETE
### 2026-08-15 10:04 RUN-RESUME
- Sequence contract: v2
- Recovery directive schema: v1
EOF
run_invalid_sequence "drill-18 later schema cannot cure an untyped recovery gap" "cannot activate legacy Recovery directive schema after an unmarked v2 prefix" "$delayed_schema"

fabricated_stop="$workspace/fabricated-stop-resume.md"
cat > "$fabricated_stop" <<'EOF'
### 2026-08-15 10:00 RUN-START
- Sequence contract: v2
- Recovery directive schema: v1
### 2026-08-15 10:01 STOP
- Cause: fabricated stop
- Triage: none
- Protected category: none
- Next Human or external-authority action: none
### 2026-08-15 10:02 RUN-RESUME
- Sequence contract: v2
- Recovery directive schema: v1
- Resolution evidence: fabricated resume
### 2026-08-15 10:03 TASK task-18 COMPLETE
EOF
run_invalid_sequence "drill-18 STOP requires a preceding triaged BLOCK" "STOP has no preceding TRIAGE BLOCK" "$fabricated_stop"

skipped_stop="$control_phase/triage-block-skips-stop.md"
cat > "$skipped_stop" <<'EOF'
### 2026-08-15 10:00 RUN-START
- Sequence contract: v2
- Recovery directive schema: v1
### 2026-08-15 10:01 TRIAGE TRG-1-003
- Trigger: protected product decision
- Trigger signature: task-18:product-decision:ambiguous-output
- Attempt count: 1
- Checkpoint: task-18 blocker candidate
- Verdict: BLOCK (escalated as ESC-1-003)
- Resolution evidence: product authority required
- Triage log: phase triage entry and ESC confirmed
### 2026-08-15 10:02 RUN-RESUME
- Sequence contract: v2
- Recovery directive schema: v1
- Resolution evidence: DEC-200 selected the output
### 2026-08-15 10:03 TASK task-18 COMPLETE
EOF
run_invalid_sequence "drill-18 TRIAGE BLOCK cannot skip its STOP event" "TRIAGE BLOCK must be followed by STOP" "$skipped_stop"

blocked_phase="$workspace/phase-1"
mkdir -p "$blocked_phase"
cat > "$blocked_phase/decisions.md" <<'EOF'
# Phase 1 Decisions

### DEC-200: Authorize the bounded control repair
EOF
cat > "$blocked_phase/triage-log.md" <<'EOF'
# Phase 1 Triage Log

## TRG-1-010 - protected control classification
- **Verdict:** BLOCK
- **Escalation:** ESC-1-010

## TRG-1-020 - recovery budget exhausted
- **Verdict:** BLOCK
- **Escalation:** ESC-1-020
EOF
valid_blocked_control="$blocked_phase/run-ledger.md"
cat > "$valid_blocked_control" <<'EOF'
# Phase 1 Run Ledger

### 2026-08-15 10:00 RUN-START
- Sequence contract: v2
- Recovery directive schema: v1
### 2026-08-15 10:01 TRIAGE TRG-1-010
- Trigger: protected control classification
- Trigger signature: task-1:quality-gate-control-defect:blocked-control
- Attempt count: 1
- Checkpoint: exact control reproducer failed
- Verdict: BLOCK (escalated as ESC-1-010)
- Resolution evidence: protected authority required
- Triage log: phase triage entry and ESC confirmed
### 2026-08-15 10:02 STOP
- Cause: TRG-1-010 / ESC-1-010 requires authority
- Triage: TRG-1-010 (BLOCK, escalated as ESC-1-010)
- Protected category: gate/control definition authority
- Next Human or external-authority action: record a decision
### 2026-08-15 10:03 RUN-RESUME
- Sequence contract: v2
- Recovery directive schema: v1
- Prior run/session or STOP: TRG-1-010 / ESC-1-010 STOP
- Trigger signature: task-1:quality-gate-control-defect:blocked-control
- Attempt count: 2
- Checkpoint: DEC-200 recorded
- Resolution evidence: DEC-200 authorizes only conformance repair
- Continuing task: task-1
### 2026-08-15 10:04 TASK-RECOVERY task-1
- Source: TRG-1-010 / DEC-200
- Category: gate-repair
- Repair mode: control-conformance
- Repair authority: TRG-1-010 and DEC-200
- Control contract: exact commands remain opaque
- Control mismatch: exact reproducer false-rejects opaque text
- Allowed files: project-docs/tools/control.sh
- Forbidden changes: product, specification, assertion, status, scope, threshold, budget, or policy
- Trigger signature: task-1:quality-gate-control-defect:blocked-control
- Attempt count: 2
- Consecutive no-progress attempts: 0/2
- Checkpoint: DEC-200 resolution checkpoint
- Resolution evidence: PENDING - repair and verify exact reproducer
- Next action: use clean repair worktree, run workflow-control-scope-lint.sh, and rerun reproducer
- Human attention required: no
### 2026-08-15 10:05 REPAIR task-1
- Trigger signature: task-1:quality-gate-control-defect:blocked-control
- Attempt count: 2
- Control scope check: PASS - clean scope (terminated exit_code: 0)
- Outcome: RESOLVED
- Resolution evidence: exact reproducer passed
- Next checkpoint: REPORT
### 2026-08-15 10:06 TASK-RECOVERY task-1
- Source: TRG-1-010 / DEC-200 resolved control report continuation
- Category: report-repair
- Repair mode: task-scope
- Repair authority: TRG-1-010 and DEC-200
- Control contract: N/A
- Control mismatch: N/A
- Allowed files: project-docs/docs/phases/phase-1/tasks/task-1-done.md
- Forbidden changes: product, control, command, assertion, threshold, or evidence
- Trigger signature: task-1:quality-gate-control-defect:blocked-control
- Attempt count: 3
- Consecutive no-progress attempts: 0/2
- Checkpoint: resolved control repair
- Resolution evidence: PENDING - rerun verification and finalize report
- Next action: rerun exact verification and update only the report
- Human attention required: no
### 2026-08-15 10:07 REPAIR task-1
- Trigger signature: task-1:quality-gate-control-defect:blocked-control
- Attempt count: 3
- Control scope check: N/A for task-scope report continuation
- Outcome: RESOLVED
- Resolution evidence: exact verification passed and report finalized
- Next checkpoint: COMPLETE
### 2026-08-15 10:08 TASK task-1 COMPLETE
EOF
run_sequence "drill-18 resolved BLOCK binds phase TRG ESC and existing decision" ADVANCEABLE "$valid_blocked_control"

invented_decision="$blocked_phase/invented-decision.md"
sed 's/DEC-200/DEC-999999/g' "$valid_blocked_control" > "$invented_decision"
run_invalid_sequence "drill-18 invented decision cannot authorize blocked control repair" "RUN-RESUME cites an unrecorded or malformed phase decision" "$invented_decision"

second_invented_decision="$blocked_phase/second-invented-decision.md"
sed 's/DEC-200 authorizes only conformance repair/DEC-200 and DEC-999999 authorize only conformance repair/' \
  "$valid_blocked_control" > "$second_invented_decision"
run_invalid_sequence "drill-18 every cited resolution decision must exist" "RUN-RESUME cites an unrecorded or malformed phase decision" "$second_invented_decision"

embedded_decision="$blocked_phase/embedded-decision.md"
sed 's/DEC-200 authorizes only conformance repair/NOTDEC-200X authorizes only conformance repair/' \
  "$valid_blocked_control" > "$embedded_decision"
run_invalid_sequence "drill-18 embedded decision substring is not authority" "RUN-RESUME cites an unrecorded or malformed phase decision" "$embedded_decision"

mismatched_stop_esc="$blocked_phase/mismatched-stop-esc.md"
sed 's/Triage: TRG-1-010 (BLOCK, escalated as ESC-1-010)/Triage: TRG-1-010 (BLOCK, escalated as ESC-1-999)/' "$valid_blocked_control" > "$mismatched_stop_esc"
run_invalid_sequence "drill-18 STOP ESC must match its governing TRIAGE" "does not reference its TRIAGE BLOCK and concrete ESC" "$mismatched_stop_esc"

wrong_phase_authority="$blocked_phase/wrong-phase-authority.md"
sed 's/TRG-1-/TRG-999-/g; s/ESC-1-/ESC-999-/g' "$valid_blocked_control" > "$wrong_phase_authority"
run_invalid_sequence "drill-18 TRG and ESC authority must match ledger phase" "does not match ledger phase-1" "$wrong_phase_authority"

unresolved_resume="$blocked_phase/unresolved-resume.md"
sed '/^- Resolution evidence: DEC-200 authorizes only conformance repair$/d' "$valid_blocked_control" > "$unresolved_resume"
run_invalid_sequence "drill-18 RUN-RESUME after STOP requires resolution evidence" "missing reconciliation authority" "$unresolved_resume"

wrong_resume_trigger="$blocked_phase/wrong-resume-trigger.md"
sed '/^### 2026-08-15 10:03 RUN-RESUME/,/^### 2026-08-15 10:04 TASK-RECOVERY/ s/task-1:quality-gate-control-defect:blocked-control/task-1:quality-gate-control-defect:unrelated/' "$valid_blocked_control" > "$wrong_resume_trigger"
run_invalid_sequence "drill-18 RUN-RESUME retains its governing trigger" "Trigger signature does not match its governing TRIAGE" "$wrong_resume_trigger"

wrong_resume_reference="$blocked_phase/wrong-resume-reference.md"
sed 's/Prior run\/session or STOP: TRG-1-010 \/ ESC-1-010 STOP/Prior run\/session or STOP: TRG-1-999 \/ ESC-1-999 STOP/' "$valid_blocked_control" > "$wrong_resume_reference"
run_invalid_sequence "drill-18 RUN-RESUME names its governing TRIAGE and ESC" "does not reference its governing TRIAGE and ESC" "$wrong_resume_reference"

reused_resume_attempt="$blocked_phase/reused-resume-attempt.md"
sed '/^### 2026-08-15 10:03 RUN-RESUME/,/^### 2026-08-15 10:04 TASK-RECOVERY/ s/Attempt count: 2/Attempt count: 1/' "$valid_blocked_control" > "$reused_resume_attempt"
run_invalid_sequence "drill-18 RUN-RESUME increments the governing attempt" "Attempt count must increment its governing TRIAGE" "$reused_resume_attempt"

wrong_resume_task="$blocked_phase/wrong-resume-task.md"
sed 's/Continuing task: task-1/Continuing task: task-2/' "$valid_blocked_control" > "$wrong_resume_task"
run_invalid_sequence "drill-18 RUN-RESUME binds the continuing task" "Continuing task does not match its Trigger signature" "$wrong_resume_task"

wrong_post_resume_attempt="$blocked_phase/wrong-post-resume-attempt.md"
sed '/^### 2026-08-15 10:04 TASK-RECOVERY task-1$/,/^### 2026-08-15 10:05 REPAIR task-1$/ s/Attempt count: 2/Attempt count: 1/' \
  "$valid_blocked_control" > "$wrong_post_resume_attempt"
run_invalid_sequence "drill-18 first recovery preserves RUN-RESUME attempt" "does not preserve its RUN-RESUME checkpoint" "$wrong_post_resume_attempt"

wrong_post_resume_trigger="$blocked_phase/wrong-post-resume-trigger.md"
sed '/^### 2026-08-15 10:04 TASK-RECOVERY task-1$/,/^### 2026-08-15 10:05 REPAIR task-1$/ s/task-1:quality-gate-control-defect:blocked-control/task-1:quality-gate-control-defect:new-trigger/' \
  "$valid_blocked_control" > "$wrong_post_resume_trigger"
run_invalid_sequence "drill-18 first recovery preserves RUN-RESUME trigger" "does not preserve its RUN-RESUME checkpoint" "$wrong_post_resume_trigger"

missing_first_resume_fields="$blocked_phase/missing-first-resume-fields.md"
sed '/^### 2026-08-15 10:03 RUN-RESUME$/,/^### 2026-08-15 10:04 TASK-RECOVERY task-1$/ {
  /^- Prior run\/session or STOP: /d
  /^- Attempt count: /d
  /^- Continuing task: /d
}' "$valid_blocked_control" > "$missing_first_resume_fields"
run_invalid_sequence "drill-18 first marker RUN-RESUME validates its own STOP reconciliation" "missing reconciliation authority or continuing task" "$missing_first_resume_fields"

append_authority="$blocked_phase/append-authority.md"
cp "$valid_blocked_control" "$append_authority"
if append_output=$(printf '%s\n' \
  '### 2026-08-15 10:09 CORRECTION' \
  '- Correction: append validation retains the real phase authority directory' | \
  "$sequence_lint" "$append_authority" --append-stdin 2>&1); then
  append_code=0
else
  append_code=$?
fi
if [ "$append_code" -eq 0 ] && printf '%s\n' "$append_output" | grep -qF 'APPEND_CONFIRMED state=ADVANCEABLE'; then
  drill_pass "drill-18 transactional append retains external triage and decision authority"
else
  drill_fail "drill-18 transactional append authority (exit=$append_code, output=$append_output)"
fi

unresolved_repair_complete="$workspace/unresolved-repair-complete.md"
sed 's/Outcome: RESOLVED/Outcome: NO-PROGRESS/' "$valid_continuation" > "$unresolved_repair_complete"
run_invalid_sequence "drill-18 NO-PROGRESS repair cannot complete a task" "task completed before a resolved REPAIR event" "$unresolved_repair_complete"

valid_budget_recovery="$workspace/valid-budget-recovery.md"
cat > "$valid_budget_recovery" <<'EOF'
### 2026-08-15 10:00 RUN-START
- Sequence contract: v2
- Recovery directive schema: v1
### 2026-08-15 10:01 TASK-RECOVERY task-1
- Source: failed exact gate
- Category: gate-repair
- Repair mode: task-scope
- Repair authority: task-1 spec
- Control contract: N/A
- Control mismatch: N/A
- Allowed files: ordinary task scope
- Forbidden changes: protected decisions and out-of-scope files
- Trigger signature: task-1:gate-repair:exact-gate
- Attempt count: 1
- Consecutive no-progress attempts: 0/2
- Checkpoint: first failed exact gate
- Resolution evidence: PENDING - repair and rerun
- Next action: repair task-1 and rerun the exact gate
- Human attention required: no
### 2026-08-15 10:02 REPAIR task-1
- Trigger signature: task-1:gate-repair:exact-gate
- Attempt count: 1
- Control scope check: N/A for task-scope
- Outcome: NO-PROGRESS
- Resolution evidence: exact failure is unchanged
- Next checkpoint: another TASK-RECOVERY
### 2026-08-15 10:03 TASK-RECOVERY task-1
- Source: first no-progress repair
- Category: gate-repair
- Repair mode: task-scope
- Repair authority: task-1 spec
- Control contract: N/A
- Control mismatch: N/A
- Allowed files: ordinary task scope
- Forbidden changes: protected decisions and out-of-scope files
- Trigger signature: task-1:gate-repair:exact-gate
- Attempt count: 2
- Consecutive no-progress attempts: 1/2
- Checkpoint: first no-progress outcome
- Resolution evidence: PENDING - causal second repair and exact rerun
- Next action: repair task-1 and rerun the exact gate
- Human attention required: no
### 2026-08-15 10:04 REPAIR task-1
- Trigger signature: task-1:gate-repair:exact-gate
- Attempt count: 2
- Control scope check: N/A for task-scope
- Outcome: RESOLVED
- Resolution evidence: exact gate passed
- Next checkpoint: VERIFY
### 2026-08-15 10:05 TASK task-1 COMPLETE
EOF
run_sequence "drill-18 no-progress continuation preserves and resolves its budget" ADVANCEABLE "$valid_budget_recovery"

changed_budget_trigger="$workspace/changed-budget-trigger.md"
sed '/^### 2026-08-15 10:03 TASK-RECOVERY/,/^### 2026-08-15 10:05 TASK/ s/task-1:gate-repair:exact-gate/task-1:gate-repair:new-trigger/' "$valid_budget_recovery" > "$changed_budget_trigger"
run_invalid_sequence "drill-18 continuation cannot replace its trigger" "must preserve task and Trigger signature" "$changed_budget_trigger"

reused_budget_attempt="$workspace/reused-budget-attempt.md"
sed '/^### 2026-08-15 10:03 TASK-RECOVERY/,/^### 2026-08-15 10:05 TASK/ s/Attempt count: 2/Attempt count: 1/' "$valid_budget_recovery" > "$reused_budget_attempt"
run_invalid_sequence "drill-18 continuation must increment total attempt" "Attempt count must increment by one" "$reused_budget_attempt"

reset_no_progress="$workspace/reset-no-progress.md"
sed '0,/Consecutive no-progress attempts: 1\/2/s//Consecutive no-progress attempts: 0\/2/' "$valid_budget_recovery" > "$reset_no_progress"
run_invalid_sequence "drill-18 continuation cannot reset no-progress count" "expected 1/2" "$reset_no_progress"

exhausted_budget="$workspace/exhausted-budget.md"
sed -e '/^### 2026-08-15 10:05 TASK task-1 COMPLETE$/d' \
    -e '/^### 2026-08-15 10:04 REPAIR task-1$/,$ s/Outcome: RESOLVED/Outcome: NO-PROGRESS/' \
    "$valid_budget_recovery" > "$exhausted_budget"
run_sequence "drill-18 second consecutive NO-PROGRESS requires triage" TRIAGE_PENDING "$exhausted_budget"

exhausted_repair_triage="$blocked_phase/exhausted-repair-triage.md"
cat "$exhausted_budget" > "$exhausted_repair_triage"
cat >> "$exhausted_repair_triage" <<'EOF'
### 2026-08-15 10:06 TRIAGE TRG-1-020
- Trigger: two consecutive no-progress outcomes exhausted recovery
- Trigger signature: task-1:gate-repair:exact-gate
- Attempt count: 2
- Checkpoint: second no-progress outcome
- Verdict: REPAIR
- Resolution evidence: try the same recovery again
- Triage log: phase triage entry confirmed
EOF
run_invalid_sequence "drill-18 exhausted no-progress triage cannot authorize another repair" "must BLOCK" "$exhausted_repair_triage"

exhausted_blocked="$blocked_phase/exhausted-budget-blocked.md"
cat "$exhausted_budget" > "$exhausted_blocked"
cat >> "$exhausted_blocked" <<'EOF'
### 2026-08-15 10:06 TRIAGE TRG-1-020
- Trigger: two consecutive no-progress outcomes exhausted recovery
- Trigger signature: task-1:gate-repair:exact-gate
- Attempt count: 2
- Checkpoint: second no-progress outcome
- Verdict: BLOCK (escalated as ESC-1-020)
- Resolution evidence: recovery-budget-exhausted requires Human authority
- Triage log: phase triage entry and ESC confirmed
### 2026-08-15 10:07 STOP
- Cause: TRG-1-020 / ESC-1-020 exhausted recovery
- Triage: TRG-1-020 (BLOCK, escalated as ESC-1-020)
- Protected category: recovery-budget-exhausted
- Next Human or external-authority action: decide whether to authorize another approach
EOF
run_sequence "drill-18 exhausted no-progress recovery stops through triage" STOPPED "$exhausted_blocked"

wrong_exhaustion_triage="$blocked_phase/wrong-exhaustion-triage.md"
sed '/^### 2026-08-15 10:06 TRIAGE TRG-1-020$/,/^### 2026-08-15 10:07 STOP$/ {
  s/task-1:gate-repair:exact-gate/task-9:gate-repair:other/
  s/Attempt count: 2/Attempt count: 1/
}' "$exhausted_blocked" > "$wrong_exhaustion_triage"
run_invalid_sequence "drill-18 exhausted TRIAGE cannot redirect recovery identity" "changed task, trigger, or attempt" "$wrong_exhaustion_triage"

resumed_budget="$blocked_phase/resumed-budget.md"
cat "$exhausted_blocked" > "$resumed_budget"
cat >> "$resumed_budget" <<'EOF'
### 2026-08-15 10:08 RUN-RESUME
- Sequence contract: v2
- Recovery directive schema: v2
- Prior run/session or STOP: TRG-1-020 / ESC-1-020 STOP
- Governing TRG: TRG-1-020
- Governing ESC: ESC-1-020
- Trigger signature: task-1:gate-repair:exact-gate
- Attempt count: 3
- Checkpoint: DEC-200 authorizes one renewed approach
- Resolution evidence: DEC-200 authorizes one renewed approach
- Resolution authority: DEC-200
- Continuing task: task-1
- Next recovery: task-1 | gate-repair | task-scope | task-1:gate-repair:exact-gate | 3
### 2026-08-15 10:09 TASK-RECOVERY task-1
- Source: TRG-1-020 / DEC-200
- Category: gate-repair
- Repair mode: task-scope
- Repair authority: TRG-1-020 / DEC-200
- Triage disposition: N/A
- Control contract: N/A
- Control mismatch: N/A
- Allowed files: ordinary task scope
- Forbidden changes: protected decisions and out-of-scope files
- Originating report: N/A
- Exact verification: N/A
- Control repair base: N/A
- Trigger signature: task-1:gate-repair:exact-gate
- Attempt count: 3
- Consecutive no-progress attempts: 0/2
- Checkpoint: DEC-200 renewed recovery
- Phase worktree: workspace/phase-1 on phase-1; ledger-only dirty
- Temporary repair worktree: N/A
- Hidden-index flags: none
- Resolution evidence: PENDING - execute the renewed approach
- Next action: repair and rerun exact gate
- Human attention required: no
### 2026-08-15 10:10 REPAIR task-1
- Trigger signature: task-1:gate-repair:exact-gate
- Attempt count: 3
- Starting checkpoint: DEC-200 renewed recovery
- Ending checkpoint: exact gate passed
- Control repair base: N/A
- Control repair head: N/A
- Control scope check: N/A for task-scope
- Phase-branch integration: N/A for task-scope
- Outcome: RESOLVED
- Resolution evidence: exact gate passed
- Next checkpoint: COMPLETE
### 2026-08-15 10:11 TASK task-1 COMPLETE
- Report: project-docs/docs/phases/phase-1/tasks/task-1-done.md
- report-lint: PASS - OUTCOME COMPLETE
- Commit: aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa - verified, subject matched
- QA-plan commands beyond the task spec: none required
EOF
run_sequence "drill-18 decision-authorized exhausted recovery preserves identity" ADVANCEABLE "$resumed_budget"

reset_resumed_budget="$blocked_phase/reset-resumed-budget.md"
sed '/^### 2026-08-15 10:09 TASK-RECOVERY task-1$/,/^### 2026-08-15 10:11 TASK task-1 COMPLETE$/ s/task-1:gate-repair:exact-gate/task-1:gate-repair:new-trigger/' \
  "$resumed_budget" > "$reset_resumed_budget"
run_invalid_sequence "drill-18 resolved budget STOP cannot reset trigger identity" "does not preserve its RUN-RESUME checkpoint" "$reset_resumed_budget"

wrong_resume_tuple="$blocked_phase/wrong-resume-tuple.md"
sed 's/Next recovery: task-1 | gate-repair | task-scope/Next recovery: task-1 | report-repair | task-scope/' \
  "$resumed_budget" > "$wrong_resume_tuple"
run_invalid_sequence "drill-18 schema v2 resume binds exact recovery tuple" "does not match the exact schema v2 RUN-RESUME recovery tuple" "$wrong_resume_tuple"

third_no_progress_recovery="$workspace/third-no-progress-recovery.md"
cat "$exhausted_budget" > "$third_no_progress_recovery"
cat >> "$third_no_progress_recovery" <<'EOF'
### 2026-08-15 10:06 TASK-RECOVERY task-1
- Source: exhausted no-progress budget
- Category: gate-repair
- Repair mode: task-scope
- Repair authority: task-1 spec
- Control contract: N/A
- Control mismatch: N/A
- Allowed files: ordinary task scope
- Forbidden changes: protected decisions and out-of-scope files
- Trigger signature: task-1:gate-repair:exact-gate
- Attempt count: 3
- Consecutive no-progress attempts: 1/2
- Checkpoint: second no-progress outcome
- Resolution evidence: PENDING
- Next action: retry without triage
- Human attention required: no
EOF
run_invalid_sequence "drill-18 exhausted no-progress budget cannot dispatch again" "requires TRIAGE" "$third_no_progress_recovery"

wrong_task="$workspace/wrong-task.md"
cat > "$wrong_task" <<'EOF'
### 2026-08-15 10:00 RUN-START
- Sequence contract: v2
### 2026-08-15 10:01 TASK-RECOVERY task-1
### 2026-08-15 10:02 REPAIR task-2
EOF
run_invalid_sequence "drill-16 changed recovery task identity rejected" "does not match recovery task" "$wrong_task"

valid_qa="$workspace/valid-qa-rereview.md"
cat > "$valid_qa" <<'EOF'
### 2026-08-15 10:00 RUN-START
- Sequence contract: v2
- Recovery directive schema: v1
### 2026-08-15 10:01 REVIEW review-1
- Outcome: REPAIR
- Owning task: task-1
- Re-review of: none
### 2026-08-15 10:02 TASK-RECOVERY task-1
- Source: review-1
- Category: qa-rework
- Repair mode: task-scope
- Repair authority: review-1
- Control contract: N/A
- Control mismatch: N/A
- Allowed files: ordinary task scope
- Forbidden changes: protected decisions and out-of-scope files
- Trigger signature: task-1:qa-rework:review-1
- Attempt count: 1
- Consecutive no-progress attempts: 0/2
- Checkpoint: review-1 repair finding
- Resolution evidence: PENDING - repair and re-review
- Next action: repair task-1 and obtain linked re-review
- Human attention required: no
### 2026-08-15 10:03 REPAIR task-1
- Trigger signature: task-1:qa-rework:review-1
- Attempt count: 1
- Control scope check: N/A for task-scope
- Outcome: RESOLVED
- Resolution evidence: review finding repaired
- Next checkpoint: REVIEW
### 2026-08-15 10:04 TASK task-1 COMPLETE
### 2026-08-15 10:05 REVIEW review-1b
- Outcome: PASS
- Owning task: none
- Re-review of: review-1
EOF
run_sequence "drill-16 lifecycle: QA repair closes on linked fresh re-review" ADVANCEABLE "$valid_qa"

valid_sequential_qa="$workspace/valid-sequential-qa-repair.md"
cat > "$valid_sequential_qa" <<'EOF'
### 2026-08-15 10:00 RUN-START
- Sequence contract: v2
- Recovery directive schema: v1
### 2026-08-15 10:01 REVIEW review-1
- Outcome: REPAIR
- Owning task: task-1
- Re-review of: none
### 2026-08-15 10:02 TASK-RECOVERY task-1
- Source: review-1
- Category: qa-rework
- Repair mode: task-scope
- Repair authority: review-1
- Control contract: N/A
- Control mismatch: N/A
- Allowed files: ordinary task scope
- Forbidden changes: protected decisions and out-of-scope files
- Trigger signature: task-1:qa-rework:review-1
- Attempt count: 1
- Consecutive no-progress attempts: 0/2
- Checkpoint: review-1 repair finding
- Resolution evidence: PENDING - repair and re-review
- Next action: repair task-1 and obtain linked re-review
- Human attention required: no
### 2026-08-15 10:03 REPAIR task-1
- Trigger signature: task-1:qa-rework:review-1
- Attempt count: 1
- Control scope check: N/A for task-scope
- Outcome: RESOLVED
- Resolution evidence: review-1 finding repaired
- Next checkpoint: REVIEW
### 2026-08-15 10:04 TASK task-1 COMPLETE
### 2026-08-15 10:05 REVIEW review-1b
- Outcome: REPAIR
- Owning task: task-2
- Re-review of: review-1
### 2026-08-15 10:06 TASK-RECOVERY task-2
- Source: review-1b
- Category: qa-rework
- Repair mode: task-scope
- Repair authority: review-1b
- Control contract: N/A
- Control mismatch: N/A
- Allowed files: ordinary task scope
- Forbidden changes: protected decisions and out-of-scope files
- Trigger signature: task-2:qa-rework:review-1b
- Attempt count: 1
- Consecutive no-progress attempts: 0/2
- Checkpoint: review-1b repair finding
- Resolution evidence: PENDING - repair and re-review
- Next action: repair task-2 and obtain linked re-review
- Human attention required: no
### 2026-08-15 10:07 REPAIR task-2
- Trigger signature: task-2:qa-rework:review-1b
- Attempt count: 1
- Control scope check: N/A for task-scope
- Outcome: RESOLVED
- Resolution evidence: review-1b finding repaired
- Next checkpoint: REVIEW
### 2026-08-15 10:08 TASK task-2 COMPLETE
### 2026-08-15 10:09 REVIEW review-1c
- Outcome: PASS
- Owning task: none
- Re-review of: review-1b
EOF
run_sequence "drill-16 lifecycle: linked re-review transfers repair ownership" ADVANCEABLE "$valid_sequential_qa"

overwritten_qa="$workspace/overwritten-qa-repair.md"
cat > "$overwritten_qa" <<'EOF'
### 2026-08-15 10:00 RUN-START
- Sequence contract: v2
### 2026-08-15 10:01 REVIEW review-1
- Outcome: REPAIR
- Owning task: task-1
- Re-review of: none
### 2026-08-15 10:02 REVIEW review-2
- Outcome: REPAIR
- Owning task: task-2
- Re-review of: none
### 2026-08-15 10:03 TASK-RECOVERY task-2
### 2026-08-15 10:04 REPAIR task-2
### 2026-08-15 10:05 TASK task-2 COMPLETE
### 2026-08-15 10:06 REVIEW review-3
- Outcome: PASS
- Owning task: none
- Re-review of: review-2
EOF
run_invalid_sequence "drill-16 unresolved QA repair owner cannot be overwritten" "replaced before owning task completed" "$overwritten_qa"

self_replaced_qa="$workspace/self-replaced-qa-repair.md"
cat > "$self_replaced_qa" <<'EOF'
### 2026-08-15 10:00 RUN-START
- Sequence contract: v2
### 2026-08-15 10:01 REVIEW review-1
- Outcome: REPAIR
- Owning task: task-1
- Re-review of: none
### 2026-08-15 10:02 TASK-RECOVERY task-1
### 2026-08-15 10:03 REPAIR task-1
### 2026-08-15 10:04 TASK task-1 COMPLETE
### 2026-08-15 10:05 REVIEW review-1
- Outcome: REPAIR
- Owning task: task-2
- Re-review of: review-1
EOF
run_invalid_sequence "drill-16 QA repair transfer requires a fresh review ID" "self-replaced by the same review identity" "$self_replaced_qa"

unlinked_transfer_qa="$workspace/unlinked-qa-repair-transfer.md"
cat > "$unlinked_transfer_qa" <<'EOF'
### 2026-08-15 10:00 RUN-START
- Sequence contract: v2
### 2026-08-15 10:01 REVIEW review-1
- Outcome: REPAIR
- Owning task: task-1
- Re-review of: none
### 2026-08-15 10:02 TASK-RECOVERY task-1
### 2026-08-15 10:03 REPAIR task-1
### 2026-08-15 10:04 TASK task-1 COMPLETE
### 2026-08-15 10:05 REVIEW review-1b
- Outcome: REPAIR
- Owning task: task-2
- Re-review of: none
EOF
run_invalid_sequence "drill-16 QA repair transfer requires explicit linkage" "does not reference required re-review" "$unlinked_transfer_qa"

missing_rereview="$workspace/missing-qa-rereview.md"
sed '/^### 2026-08-15 10:05 REVIEW review-1b$/,$d' "$valid_qa" > "$missing_rereview"
run_sequence "drill-16 QA repair cannot advance before re-review" REREVIEW_PENDING "$missing_rereview"

self_rereview="$workspace/self-qa-rereview.md"
sed -e 's/REVIEW review-1b/REVIEW review-1/' "$valid_qa" > "$self_rereview"
run_invalid_sequence "drill-16 same review identity cannot self-close" "self-closed by the same review identity" "$self_rereview"

missing_outcome="$workspace/missing-review-outcome.md"
sed '/^- Outcome: REPAIR$/d' "$valid_qa" > "$missing_outcome"
run_invalid_sequence "drill-16 strict review outcome is mandatory" "missing or invalid Outcome" "$missing_outcome"

missing_task_id="$workspace/missing-task-id.md"
cat > "$missing_task_id" <<'EOF'
### 2026-08-15 10:00 RUN-START
- Sequence contract: v2
### 2026-08-15 10:01 TASK-RECOVERY
EOF
run_invalid_sequence "drill-16 strict recovery task ID is mandatory" "missing or invalid task ID" "$missing_task_id"

missing_review_id="$workspace/missing-review-id.md"
cat > "$missing_review_id" <<'EOF'
### 2026-08-15 10:00 RUN-START
- Sequence contract: v2
### 2026-08-15 10:01 REVIEW
- Outcome: PASS
- Owning task: none
- Re-review of: none
EOF
run_invalid_sequence "drill-16 strict review ID is mandatory" "missing or invalid review ID" "$missing_review_id"

overlapping_recovery="$workspace/overlapping-recovery.md"
cat > "$overlapping_recovery" <<'EOF'
### 2026-08-15 10:00 RUN-START
- Sequence contract: v2
### 2026-08-15 10:01 TASK-RECOVERY task-1
### 2026-08-15 10:02 TASK-RECOVERY task-2
### 2026-08-15 10:03 REPAIR task-2
### 2026-08-15 10:04 TASK task-2 COMPLETE
EOF
run_invalid_sequence "drill-16 overlapping recovery cannot change task" "overlapping recovery changes task identity" "$overlapping_recovery"

qa_without_recovery="$workspace/qa-without-recovery.md"
cat > "$qa_without_recovery" <<'EOF'
### 2026-08-15 10:00 RUN-START
- Sequence contract: v2
### 2026-08-15 10:01 REVIEW review-1
- Outcome: REPAIR
- Owning task: task-1
- Re-review of: none
### 2026-08-15 10:02 TASK task-1 COMPLETE
### 2026-08-15 10:03 REVIEW review-1b
- Outcome: PASS
- Owning task: none
- Re-review of: review-1
EOF
run_invalid_sequence "drill-16 QA repair requires recovery events" "completed without TASK-RECOVERY and REPAIR" "$qa_without_recovery"

triage_repair="$control_phase/triage-repair-pending.md"
cat > "$triage_repair" <<'EOF'
### 2026-08-15 10:00 RUN-START
- Sequence contract: v2
- Recovery directive schema: v1
### 2026-08-15 10:01 TRIAGE TRG-1-002
- Trigger: task-1 needs bounded repair
- Trigger signature: task-1:gate-repair:focused-test
- Attempt count: 1
- Checkpoint: focused test failed
- Verdict: REPAIR
- Resolution evidence: task-scope repair is authorized
- Triage log: phase triage entry confirmed
EOF
run_sequence "drill-16 TRIAGE REPAIR cannot advance before recovery" REPAIR_PENDING "$triage_repair"

triage_repair_complete="$control_phase/triage-repair-complete.md"
cat "$triage_repair" > "$triage_repair_complete"
cat >> "$triage_repair_complete" <<'EOF'
### 2026-08-15 10:02 TASK-RECOVERY task-1
- Source: TRG-1-002
- Category: gate-repair
- Repair mode: task-scope
- Repair authority: TRG-1-002
- Control contract: N/A
- Control mismatch: N/A
- Allowed files: ordinary task scope
- Forbidden changes: protected decisions and out-of-scope files
- Trigger signature: task-1:gate-repair:focused-test
- Attempt count: 1
- Consecutive no-progress attempts: 0/2
- Checkpoint: TRG-1-002 repair authority
- Resolution evidence: PENDING - repair and rerun focused test
- Next action: repair and rerun the exact focused test
- Human attention required: no
### 2026-08-15 10:03 REPAIR task-1
- Trigger signature: task-1:gate-repair:focused-test
- Attempt count: 1
- Control scope check: N/A for task-scope
- Outcome: RESOLVED
- Resolution evidence: focused test passed
- Next checkpoint: COMPLETE
### 2026-08-15 10:04 TASK task-1 COMPLETE
EOF
run_sequence "drill-16 ordinary TRIAGE REPAIR binds its recovery" ADVANCEABLE "$triage_repair_complete"

unrelated_triage_recovery="$control_phase/unrelated-triage-recovery.md"
sed '/^### 2026-08-15 10:02 TASK-RECOVERY task-1$/,$ {
  s/TASK-RECOVERY task-1/TASK-RECOVERY task-2/
  s/REPAIR task-1/REPAIR task-2/
  s/TASK task-1 COMPLETE/TASK task-2 COMPLETE/
  s/task-1:gate-repair:focused-test/task-2:gate-repair:unrelated/
}' "$triage_repair_complete" > "$unrelated_triage_recovery"
run_invalid_sequence "drill-16 TRIAGE REPAIR cannot be cleared by unrelated task" "does not preserve its governing TRIAGE REPAIR authority" "$unrelated_triage_recovery"

# Transactional append rejects a bad candidate without contaminating the
# append-only ledger.
pending_ledger="$workspace/pending-recovery.md"
cat > "$pending_ledger" <<'EOF'
### 2026-08-15 10:00 RUN-START
- Sequence contract: v2
- Recovery directive schema: v1
### 2026-08-15 10:01 TASK-RECOVERY task-1
- Source: interrupted task-1 session
- Category: session-continuation
- Repair mode: task-scope
- Repair authority: task-1 spec
- Control contract: N/A
- Control mismatch: N/A
- Allowed files: ordinary task scope
- Forbidden changes: protected decisions and out-of-scope files
- Trigger signature: task-1:session-continuation:interrupted-session
- Attempt count: 1
- Consecutive no-progress attempts: 0/2
- Checkpoint: committed partial task-1 implementation
- Resolution evidence: PENDING - complete task-1
- Next action: complete task-1 and rerun its exact gate
- Human attention required: no
EOF
cp "$pending_ledger" "$pending_ledger.before"
proposed_output=$(printf '%s\n' '### 2026-08-15 10:02 REPAIR task-2' \
  | "$sequence_lint" "$pending_ledger" --append-stdin 2>&1)
proposed_code=$?
current_output=$("$sequence_lint" "$pending_ledger" 2>&1)
if [ "$proposed_code" -ne 0 ] && \
   printf '%s\n' "$proposed_output" | grep -qF 'does not match recovery task' && \
   cmp -s "$pending_ledger.before" "$pending_ledger" && \
   [ "$current_output" = 'VALID state=REPAIR_PENDING' ]; then
  drill_pass "drill-16 invalid proposed append leaves ledger usable"
else
  drill_fail "drill-16 rejected append is non-mutating (candidate=$proposed_output, current=$current_output)"
fi

# The append command, not a later generic patch, must own the validated write.
# Repeated field text makes an anchor-based insertion land in the wrong event.
atomic_ledger="$workspace/atomic-append.md"
cat > "$atomic_ledger" <<'EOF'
### 2026-08-15 10:00 RUN-START
- Sequence contract: v2
- Recovery directive schema: v1
- Human attention required: no
### 2026-08-15 10:01 TASK task-1 COMPLETE
- Human attention required: no
### 2026-08-15 10:02 TASK-RECOVERY task-2
- Source: incomplete task-2
- Category: session-continuation
- Repair mode: task-scope
- Repair authority: task-2 spec
- Control contract: N/A
- Control mismatch: N/A
- Allowed files: ordinary task scope
- Forbidden changes: protected decisions and out-of-scope files
- Trigger signature: task-2:session-continuation:atomic-append
- Attempt count: 1
- Consecutive no-progress attempts: 0/2
- Checkpoint: committed partial task-2 implementation
- Resolution evidence: PENDING - complete task-2
- Next action: complete task-2 and rerun its exact gate
- Human attention required: no
EOF
atomic_candidate="$workspace/atomic-candidate.md"
cat > "$atomic_candidate" <<'EOF'
### 2026-08-15 10:03 REPAIR task-2
- Trigger signature: task-2:session-continuation:atomic-append
- Attempt count: 1
- Control scope check: N/A for task-scope
- Outcome: RESOLVED
- Resolution evidence: exact task-2 gate passed
- Next checkpoint: VERIFY
EOF
cp "$atomic_ledger" "$atomic_ledger.expected"
printf '\n' >> "$atomic_ledger.expected"
cat "$atomic_candidate" >> "$atomic_ledger.expected"
atomic_output=$("$sequence_lint" "$atomic_ledger" --append-stdin < "$atomic_candidate" 2>&1)
atomic_code=$?
current_output=$("$sequence_lint" "$atomic_ledger" 2>&1)
if [ "$atomic_code" -eq 0 ] && \
   [ "$atomic_output" = 'APPEND_CONFIRMED state=COMPLETION_PENDING' ] && \
   cmp -s "$atomic_ledger.expected" "$atomic_ledger" && \
   [ "$current_output" = 'VALID state=COMPLETION_PENDING' ]; then
  drill_pass "drill-16 valid candidate is atomically appended at EOF"
else
  drill_fail "drill-16 atomic EOF append (append=$atomic_output, current=$current_output)"
fi

# A receipt can be lost after the atomic rename. Retrying the byte-identical
# candidate must confirm the existing EOF event rather than append a duplicate.
ack_ledger="$workspace/acknowledgment-retry.md"
cat > "$ack_ledger" <<'EOF'
### 2026-08-15 10:00 RUN-START
- Sequence contract: v2
- Recovery directive schema: v1
### 2026-08-15 10:01 TASK-RECOVERY task-3
- Source: incomplete task-3
- Category: session-continuation
- Repair mode: task-scope
- Repair authority: task-3 spec
- Control contract: N/A
- Control mismatch: N/A
- Allowed files: ordinary task scope
- Forbidden changes: protected decisions and out-of-scope files
- Trigger signature: task-3:session-continuation:lost-receipt
- Attempt count: 1
- Consecutive no-progress attempts: 0/2
- Checkpoint: committed partial task-3 implementation
- Resolution evidence: PENDING - complete task-3
- Next action: complete task-3 and rerun its exact gate
- Human attention required: no
EOF
ack_candidate="$workspace/acknowledgment-candidate.md"
cat > "$ack_candidate" <<'EOF'
### 2026-08-15 10:02 REPAIR task-3
- Trigger signature: task-3:session-continuation:lost-receipt
- Attempt count: 1
- Control scope check: N/A for task-scope
- Outcome: RESOLVED
- Resolution evidence: exact task-3 gate passed
- Next checkpoint: VERIFY
EOF
cp "$ack_ledger" "$ack_ledger.expected"
printf '\n' >> "$ack_ledger.expected"
cat "$ack_candidate" >> "$ack_ledger.expected"
"$sequence_lint" "$ack_ledger" --append-stdin < "$ack_candidate" \
  1>&- 2> "$workspace/lost-receipt.err"
lost_receipt_code=$?
retry_output=$("$sequence_lint" "$ack_ledger" --append-stdin < "$ack_candidate" 2>&1)
retry_code=$?
if [ "$lost_receipt_code" -ne 0 ] && \
   [ "$retry_code" -eq 0 ] && \
   [ "$retry_output" = 'APPEND_CONFIRMED state=COMPLETION_PENDING' ] && \
   cmp -s "$ack_ledger.expected" "$ack_ledger"; then
  drill_pass "drill-16 lost append receipt is safely confirmed by exact retry"
else
  drill_fail "drill-16 append acknowledgment retry (first=$lost_receipt_code, retry=$retry_output)"
fi

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
