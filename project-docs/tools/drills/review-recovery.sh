#!/bin/sh
# Regression for DEC-190: ordinary triage correction and review recovery.
set -eu
script_dir=$(CDPATH= cd "$(dirname "$0")" && pwd)
validator="$script_dir/../workflow-sequence-lint.sh"
scratch=$(mktemp -d "${TMPDIR:-/tmp}/review-recovery.XXXXXX")
trap 'rm -rf "$scratch"' EXIT HUP INT TERM
phase="$scratch/phase-1"
mkdir -p "$phase"
ledger="$phase/run-ledger.md"
failures=0

cat > "$scratch/triage" <<'EOF'
## TRG-1-001 - original misbound triage
- **Originating task:** task-2
- **Trigger signature:** task-2:latent-gate-failure:golden
- **Attempt count:** 2
- **Repair category:** in-scope-implementation
- **Repair authority:** minimum causal scope
- **Verdict:** REPAIR
## TRG-1-002 - same-task correction
- **Originating task:** task-2
- **Trigger signature:** review-1:qa-rework:task-2-stop-attention-evidence
- **Attempt count:** 2
- **Repair category:** in-scope-implementation
- **Repair authority:** DEC-901
- **Verdict:** REPAIR
EOF
echo '### DEC-901: accepted harness repair' > "$phase/decisions.md"
cat > "$scratch/base" <<'EOF'
### 2026-09-06 10:00 RUN-START
- Sequence contract: v2
- Recovery directive schema: v2
- Branch: phase-1
- ESC baseline (Open at run start): none
- Warm-up/autonomous-resume state: autonomous
- Starting task: task-2
### 2026-09-06 10:01 REVIEW review-1
- Path: project-docs/docs/phases/phase-1/qa/review-1.md
- Tasks covered: task-2,task-3
- Outcome: REPAIR
- Owning task: task-2,task-3
- Re-review of: none
- Must-fix: 2 / Improvement: 0 / blocker candidates: 0
### 2026-09-06 10:02 TASK-RECOVERY task-2
- Source: review-1
- Category: qa-rework
- Repair mode: task-scope
- Repair authority: QA review
- Triage disposition: N/A
- Trigger signature: review-1:qa-rework:task-2-stop-attention-evidence
- Attempt count: 1
- Control contract: N/A
- Control mismatch: N/A
- Allowed files: test.sh,task-2-done.md
- Forbidden changes: accepted assertions
- Originating report: N/A
- Exact verification: N/A
- Control repair base: N/A
- Consecutive no-progress attempts: 0/2
- Checkpoint: original repair
- Phase worktree: synthetic fixture
- Temporary repair worktree: N/A
- Hidden-index flags: none
- Resolution evidence: PENDING
- Next action: repair
- Human attention required: no
### 2026-09-06 10:03 REPAIR task-2
- Trigger signature: review-1:qa-rework:task-2-stop-attention-evidence
- Attempt count: 1
- Starting checkpoint: original repair
- Ending checkpoint: findings fixed, retained gate failed
- Control repair base: N/A
- Control repair head: N/A
- Control scope check: N/A
- Phase-branch integration: N/A
- Outcome: MATERIAL-PROGRESS
- Resolution evidence: failed receipt preserved
- Next checkpoint: triage
EOF
cat > "$scratch/old" <<'EOF'
### 2026-09-06 10:04 TRIAGE TRG-1-001
- Trigger: retained gate failure
- Originating task: task-2
- Trigger signature: task-2:latent-gate-failure:golden
- Attempt count: 2
- Repair category: in-scope-implementation
- Repair authority: minimum causal scope
- Checkpoint: material progress
- Verdict: REPAIR
- Resolution evidence: bounded repair
- Triage log: TRG-1-001
EOF
cat > "$scratch/correction" <<'EOF'
### 2026-09-06 10:06 TRIAGE TRG-1-002
- Trigger: correcting same-task directive
- Originating task: task-2
- Trigger signature: review-1:qa-rework:task-2-stop-attention-evidence
- Attempt count: 2
- Repair category: in-scope-implementation
- Repair authority: DEC-901
- Checkpoint: material progress
- Verdict: REPAIR
- Resolution evidence: bounded repair correction
- Triage log: TRG-1-002
EOF
cat > "$scratch/resume" <<'EOF'
### 2026-09-06 10:05 RUN-RESUME
- Sequence contract: v2
- Recovery directive schema: v2
- Prior run/session or STOP: prior runner; no persisted STOP
- Governing TRG: N/A
- Governing ESC: N/A
- Trigger signature: N/A
- Attempt count: N/A
- Checkpoint: pending triage
- Resolution evidence: same continuation preserved
- Resolution authority: N/A
- Continuing task: task-2
- Next recovery: N/A
EOF
cat > "$scratch/recovery" <<'EOF'
### 2026-09-06 10:07 TASK-RECOVERY task-2
- Source: TRG-1-002
- Category: in-scope-implementation
- Repair mode: task-scope
- Repair authority: DEC-901
- Triage disposition: N/A
- Trigger signature: review-1:qa-rework:task-2-stop-attention-evidence
- Attempt count: 2
- Control contract: N/A
- Control mismatch: N/A
- Allowed files: test.sh,task-2-done.md
- Forbidden changes: accepted assertions
- Originating report: N/A
- Exact verification: N/A
- Control repair base: N/A
- Consecutive no-progress attempts: 0/2
- Checkpoint: corrected triage
- Phase worktree: synthetic fixture
- Temporary repair worktree: N/A
- Hidden-index flags: none
- Resolution evidence: PENDING
- Next action: causal repair
- Human attention required: no
EOF

reset_fixture() {
  cp "$scratch/triage" "$phase/triage-log.md"
  cat "$scratch/base" "$scratch/old" > "$ledger"
}
pass() { echo "PASS $1"; }
fail() { echo "FAIL $1: $output"; failures=$((failures + 1)); }
accept() {
  label=$1
  candidate=$2
  want=${3:-REPAIR_PENDING}
  code=0
  output=$("$validator" "$ledger" --append-stdin < "$candidate" 2>&1) || code=$?
  if [ "$code" -eq 0 ] && [ "$output" = "APPEND_CONFIRMED state=$want" ]; then
    pass "$label"
  else
    fail "$label (exit=$code)"
  fi
}
reject() {
  label=$1
  candidate=$2
  reason=$3
  cp "$ledger" "$scratch/before"
  code=0
  output=$("$validator" "$ledger" --append-stdin < "$candidate" 2>&1) || code=$?
  if [ "$code" -eq 1 ] && printf '%s\n' "$output" | grep -qF "$reason" &&
      cmp -s "$ledger" "$scratch/before"; then
    pass "$label; rejected append preserved bytes"
  else
    fail "$label (exit=$code)"
  fi
}

# Isolate review identity from the pending-triage ordering defect.
reset_fixture
cp "$scratch/base" "$ledger"
accept 'active review trigger accepted without earlier triage' "$scratch/correction"
reset_fixture
accept 'normal runner resume preserves pending ordinary triage' "$scratch/resume"
reset_fixture
accept 'same-task correction supersedes unconsumed triage' "$scratch/correction"

# Full real-boundary sequence, including an idempotent lost-receipt retry.
reset_fixture
accept 'restart entrypoint' "$scratch/resume"
accept 'correction entrypoint' "$scratch/correction"
cp "$ledger" "$scratch/confirmed"
accept 'byte-identical correction retry' "$scratch/correction"
if cmp -s "$ledger" "$scratch/confirmed"; then
  pass 'retry did not duplicate an event'
else
  output='ledger changed'; fail 'retry did not duplicate an event'
fi
accept 'bound continuation entrypoint' "$scratch/recovery"

# Keep the candidate and external identity equal, so semantic guards are tested.
for mutation in task attempt review token category; do
  reset_fixture
  case "$mutation" in
    task) edit='s/task-2/task-3/g'; reason='TRIAGE REPAIR must be followed' ;;
    attempt) edit='s/Attempt count: 2/Attempt count: 3/;s/Attempt count:\*\* 2/Attempt count:** 3/'; reason='TRIAGE REPAIR must be followed' ;;
    review) edit='s/review-1:qa-rework/review-9:qa-rework/g'; reason='Originating task does not match Trigger signature' ;;
    token) edit='s/qa-rework:task-2-/qa-rework:task-20-/g'; reason='Originating task does not match Trigger signature' ;;
    category) edit='s/review-1:qa-rework/review-1:gate-repair/g'; reason='Originating task does not match Trigger signature' ;;
  esac
  sed "$edit" "$scratch/correction" > "$scratch/bad"
  # Mutate only the correction record, leaving the old ledger authority intact.
  sed '/^## TRG-1-002/,$d' "$scratch/triage" > "$phase/triage-log.md"
  sed -n '/^## TRG-1-002/,$p' "$scratch/triage" | sed "$edit" >> "$phase/triage-log.md"
  reject "invalid correction $mutation" "$scratch/bad" "$reason"
done

reset_fixture
sed 's/Continuing task: task-2/Continuing task: task-3/' "$scratch/resume" > "$scratch/bad"
reject 'restart cannot redirect pending task' "$scratch/bad" 'TRIAGE REPAIR must be followed'
reset_fixture
sed 's/Resolution authority: N\/A/Resolution authority: DEC-901/' "$scratch/resume" > "$scratch/bad"
reject 'normal restart cannot smuggle STOP authority' "$scratch/bad" 'normal schema v2 RUN-RESUME'

# A valid corrected prefix, used to check binding independently of append success.
reset_fixture
cat "$scratch/correction" >> "$ledger"
for mutation in source authority attempt trigger budget; do
  case "$mutation" in
    source) edit='s/Source: TRG-1-002/Source: TRG-1-001/'; reason='governing schema v2 TRIAGE REPAIR authority' ;;
    authority) edit='s/Repair authority: DEC-901/Repair authority: minimum causal scope/'; reason='governing schema v2 TRIAGE REPAIR authority' ;;
    attempt) edit='s/Attempt count: 2/Attempt count: 1/'; reason='Attempt count must increment by one' ;;
    trigger) edit='s/review-1:qa-rework:task-2-stop-attention-evidence/task-2:latent-gate-failure:golden/'; reason='must preserve task and Trigger signature' ;;
    budget) edit='s/no-progress attempts: 0\/2/no-progress attempts: 1\/2/'; reason='expected 0/2' ;;
  esac
  sed "$edit" "$scratch/recovery" > "$scratch/bad"
  reject "bound recovery rejects $mutation drift" "$scratch/bad" "$reason"
done

# Review identities must survive budget exhaustion and protected STOP resume.
reset_fixture
sed 's/Outcome: MATERIAL-PROGRESS/Outcome: NO-PROGRESS/' "$scratch/base" > "$ledger"
sed 's/no-progress attempts: 0\/2/no-progress attempts: 1\/2/' "$scratch/recovery" >> "$ledger"
sed -n '/^### .* REPAIR task-2/,$p' "$scratch/base" |
  sed 's/10:03/10:08/;s/Attempt count: 1/Attempt count: 2/;s/Outcome: MATERIAL-PROGRESS/Outcome: NO-PROGRESS/' >> "$ledger"
reject 'exhausted review recovery cannot dispatch again' "$scratch/recovery" 'requires TRIAGE'
sed -e 's/Repair category: in-scope-implementation/Repair category: N\/A/' \
    -e 's/Repair authority: DEC-901/Repair authority: N\/A/' \
    -e 's/Verdict: REPAIR/Verdict: BLOCK (escalated as ESC-1-001)/' \
    "$scratch/correction" > "$scratch/block"
sed -e 's/Repair category:\*\* in-scope-implementation/Repair category:** N\/A/' \
    -e 's/Repair authority:\*\* DEC-901/Repair authority:** N\/A/' \
    -e 's/Verdict:\*\* REPAIR/Verdict:** BLOCK/' \
    "$scratch/triage" > "$phase/triage-log.md"
echo '- **If BLOCK:** ESC-1-001' >> "$phase/triage-log.md"
accept 'exhausted review trigger retains task at BLOCK' "$scratch/block" STOP_PENDING
reject 'normal resume cannot skip required STOP' "$scratch/resume" 'must be followed by STOP'
cat > "$scratch/stop" <<'EOF'
### 2026-09-06 10:09 STOP
- Cause: exhausted review recovery
- Triage: TRG-1-002 BLOCK ESC-1-001
- Protected category: recovery-budget-exhausted
- Next Human or external-authority action: authorize bounded recovery
EOF
accept 'review recovery STOP' "$scratch/stop" STOPPED
cat > "$scratch/stop-resume" <<'EOF'
### 2026-09-06 10:10 RUN-RESUME
- Sequence contract: v2
- Recovery directive schema: v2
- Prior run/session or STOP: TRG-1-002 ESC-1-001
- Governing TRG: TRG-1-002
- Governing ESC: ESC-1-001
- Trigger signature: review-1:qa-rework:task-2-stop-attention-evidence
- Attempt count: 3
- Checkpoint: stopped recovery
- Resolution evidence: protected authority recorded and verified
- Resolution authority: DEC-901
- Continuing task: task-2
- Next recovery: task-2 | in-scope-implementation | task-scope | review-1:qa-rework:task-2-stop-attention-evidence | 3
EOF
sed 's/Continuing task: task-2/Continuing task: task-20/' "$scratch/stop-resume" > "$scratch/bad"
reject 'STOP resume rejects wrong review task' "$scratch/bad" 'Continuing task does not match'
accept 'review trigger survives protected STOP resume' "$scratch/stop-resume"
reject 'normal resume cannot replace STOP-bound recovery' "$scratch/resume" 'RUN-RESUME must be followed'
sed -e 's/Attempt count: 2/Attempt count: 3/' \
    -e 's/Repair authority: DEC-901/Repair authority: DEC-901 TRG-1-002/' \
    -e 's/Source: TRG-1-002/Source: TRG-1-002 DEC-901/' \
    "$scratch/recovery" > "$scratch/stop-recovery"
accept 'STOP recovery retains review trigger and typed authority' "$scratch/stop-recovery"

# Two interruptions keep the same review identity at session-liveness triage.
reset_fixture
sed 's/Outcome: MATERIAL-PROGRESS/Outcome: INTERRUPTED/' "$scratch/base" > "$ledger"
cat > "$scratch/interruption" <<'EOF'
- Interruption cause: external session loss
- Elapsed session duration: 2 minutes
- Command in flight: none
EOF
cat "$scratch/interruption" "$scratch/recovery" >> "$ledger"
sed -n '/^### .* REPAIR task-2/,$p' "$scratch/base" |
  sed 's/10:03/10:08/;s/Attempt count: 1/Attempt count: 2/;s/Outcome: MATERIAL-PROGRESS/Outcome: INTERRUPTED/' >> "$ledger"
cat "$scratch/interruption" >> "$ledger"
accept 'session-liveness triage accepts owning review trigger' "$scratch/correction"
sed 's/Attempt count: 2/Attempt count: 3/' "$scratch/recovery" > "$scratch/session-recovery"
accept 'session-liveness continuation increments review attempt' "$scratch/session-recovery"

if [ "$failures" -ne 0 ]; then
  echo "REVIEW RECOVERY FAIL $failures checks"
  exit 1
fi
echo 'REVIEW RECOVERY PASS all checks'
