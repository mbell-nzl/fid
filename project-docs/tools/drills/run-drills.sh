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
printf '%s\n' 'dirty product' > "$scope_repo/cmd/product.go"
run_control_scope "drill-18 dirty product state invalidates control repair scope" nonzero "worktree is not clean" \
  "$scope_repo" "$scope_base" "$scope_control_head" project-docs/tools/control.sh
git -C "$scope_repo" checkout -q -- cmd/product.go
printf '%s\n' 'changed product' > "$scope_repo/cmd/product.go"
git -C "$scope_repo" add -A
git_fixture "$scope_repo" commit "test: contaminate control repair"
scope_product_head=$(git -C "$scope_repo" rev-parse HEAD)
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
run_control_scope "drill-18 reverted intermediate edit cannot hide a non-isolated repair" nonzero "exactly one commit" \
  "$scope_repo" "$scope_base" "$scope_reverted_head" project-docs/tools/control.sh

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
valid_continuation="$workspace/valid-continuation.md"
cat > "$valid_continuation" <<'EOF'
### 2026-08-15 10:00 RUN-RESUME
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
- Resolution evidence: PENDING - complete and verify task-1
- Next action: complete task-1 and rerun its exact gate
- Human attention required: no
### 2026-08-15 10:02 REPAIR task-1
- Trigger signature: task-1:session-continuation:interrupted-session
- Attempt count: 1
- Control scope check: N/A for task-scope
- Outcome: RESOLVED
- Resolution evidence: exact task gate passed
- Next checkpoint: VERIFY
### 2026-08-15 10:03 TASK task-1 COMPLETE
EOF
run_sequence "drill-16 lifecycle: fresh session retains task identity" ADVANCEABLE "$valid_continuation"

valid_control_repair="$workspace/valid-control-repair.md"
cat > "$valid_control_repair" <<'EOF'
### 2026-08-15 10:00 RUN-RESUME
- Sequence contract: v2
- Recovery directive schema: v1
### 2026-08-15 10:00 TRIAGE TRG-1-001
- Trigger: report-lint violates its byte-exact command contract
- Trigger signature: task-1:quality-gate-control-defect:literal-command
- Attempt count: 1
- Checkpoint: exact reproducer failed
- Verdict: REPAIR / control-conformance
- Resolution evidence: bounded tools-only repair
- Triage log: phase triage entry confirmed
### 2026-08-15 10:01 TASK-RECOVERY task-1
- Source: TRG-1-001
- Category: report-repair
- Repair mode: control-conformance
- Repair authority: TRG-1-001
- Control contract: exact commands are compared byte-for-byte
- Control mismatch: exact reproducer false-rejects a literal backslash-dot
- Allowed files: project-docs/tools/report-lint.sh
- Forbidden changes: command, assertion, accepted status, scope, threshold, budget, or tolerance
- Trigger signature: task-1:quality-gate-control-defect:literal-command
- Attempt count: 1
- Consecutive no-progress attempts: 0/2
- Checkpoint: committed failing control fixture
- Resolution evidence: PENDING - repair and rerun the exact reproducer
- Next action: repair the control, run workflow-control-scope-lint.sh, and rerun the exact reproducer
- Human attention required: no
### 2026-08-15 10:02 REPAIR task-1
- Trigger signature: task-1:quality-gate-control-defect:literal-command
- Attempt count: 1
- Control scope check: PASS - control repair scope (terminated exit_code: 0)
- Outcome: RESOLVED
- Resolution evidence: exact reproducer passed
- Next checkpoint: VERIFY
### 2026-08-15 10:03 TASK task-1 COMPLETE
EOF
run_sequence "drill-18 scoped control repair retains originating task identity" ADVANCEABLE "$valid_control_repair"

schema_omitted="$workspace/schema-omitted.md"
sed '/^- Recovery directive schema: /d' "$valid_continuation" > "$schema_omitted"
run_invalid_sequence "drill-18 fresh v2 ledger cannot omit directive schema" "without Recovery directive schema: v1" "$schema_omitted"

untriaged_control_repair="$workspace/untriaged-control-repair.md"
sed '/^### 2026-08-15 10:00 TRIAGE TRG-1-001$/,/^- Triage log: /d' "$valid_control_repair" > "$untriaged_control_repair"
run_invalid_sequence "drill-18 control repair requires preceding triage authority" "has no preceding TRIAGE authority" "$untriaged_control_repair"

self_authorized_control_repair="$workspace/self-authorized-control-repair.md"
sed 's/Repair authority: TRG-1-001/Repair authority: SELF-AUTHORIZED/' "$valid_control_repair" > "$self_authorized_control_repair"
run_invalid_sequence "drill-18 control repair authority binds the latest triage" "does not reference latest TRIAGE TRG-1-001" "$self_authorized_control_repair"

unbudgeted_control_repair="$workspace/unbudgeted-control-repair.md"
sed '/^- Consecutive no-progress attempts: /d' "$valid_control_repair" > "$unbudgeted_control_repair"
run_invalid_sequence "drill-18 control repair cannot discard its recovery budget" "missing or invalid Consecutive no-progress attempts" "$unbudgeted_control_repair"

missing_control_contract="$workspace/missing-control-contract.md"
sed '/^- Control contract: /d' "$valid_control_repair" > "$missing_control_contract"
run_invalid_sequence "drill-18 control repair requires accepted contract authority" "missing Control contract" "$missing_control_contract"

missing_repair_authority="$workspace/missing-repair-authority.md"
sed '/^- Repair authority: /d' "$valid_control_repair" > "$missing_repair_authority"
run_invalid_sequence "drill-18 marked recovery requires recorded authority" "missing Repair authority" "$missing_repair_authority"

missing_control_mismatch="$workspace/missing-control-mismatch.md"
sed '/^- Control mismatch: /d' "$valid_control_repair" > "$missing_control_mismatch"
run_invalid_sequence "drill-18 control repair requires exact mismatch evidence" "missing Control mismatch" "$missing_control_mismatch"

missing_control_files="$workspace/missing-control-files.md"
sed '/^- Allowed files: /d' "$valid_control_repair" > "$missing_control_files"
run_invalid_sequence "drill-18 control repair requires bounded file authority" "missing Allowed files" "$missing_control_files"

missing_control_guard="$workspace/missing-control-guard.md"
sed '/^- Forbidden changes: /d' "$valid_control_repair" > "$missing_control_guard"
run_invalid_sequence "drill-18 control repair requires semantic guardrails" "missing Forbidden changes" "$missing_control_guard"

missing_repair_mode="$workspace/missing-repair-mode.md"
sed '/^- Repair mode: /d' "$valid_control_repair" > "$missing_repair_mode"
run_invalid_sequence "drill-18 marked recovery requires an explicit repair mode" "missing or invalid Repair mode" "$missing_repair_mode"

invalid_directive_schema="$workspace/invalid-directive-schema.md"
sed 's/Recovery directive schema: v1/Recovery directive schema: v2/' "$valid_control_repair" > "$invalid_directive_schema"
run_invalid_sequence "drill-18 unknown recovery directive schema is rejected" "invalid Recovery directive schema" "$invalid_directive_schema"

missing_sequence_contract="$workspace/missing-control-sequence-contract.md"
sed '/^- Sequence contract: /d' "$valid_control_repair" > "$missing_sequence_contract"
run_invalid_sequence "drill-18 directive schema requires sequence contract" "missing Sequence contract: v2" "$missing_sequence_contract"

schema_dropped_on_resume="$workspace/schema-dropped-on-resume.md"
cat "$valid_control_repair" > "$schema_dropped_on_resume"
cat >> "$schema_dropped_on_resume" <<'EOF'
### 2026-08-15 10:04 RUN-RESUME
- Sequence contract: v2
EOF
run_invalid_sequence "drill-18 marked ledger cannot drop directive schema on resume" "missing Recovery directive schema: v1" "$schema_dropped_on_resume"

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

triage_repair="$workspace/triage-repair-pending.md"
cat > "$triage_repair" <<'EOF'
### 2026-08-15 10:00 RUN-START
- Sequence contract: v2
- Recovery directive schema: v1
### 2026-08-15 10:01 TRIAGE TRG-1-001
- Verdict: REPAIR
EOF
run_sequence "drill-16 TRIAGE REPAIR cannot advance before recovery" REPAIR_PENDING "$triage_repair"

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
