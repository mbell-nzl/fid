#!/bin/sh
# Regression for DEC-257: decision-released repair owners after a BLOCK re-review.
set -eu
script_dir=$(CDPATH= cd "$(dirname "$0")" && pwd)
validator="$script_dir/../workflow-sequence-lint.sh"
scratch=$(mktemp -d "${TMPDIR:-/tmp}/review-block-owners.XXXXXX")
trap 'rm -rf "$scratch"' EXIT HUP INT TERM
phase="$scratch/phase-1"
mkdir -p "$phase"
ledger="$phase/run-ledger.md"
failures=0

# DEC-903 follows DEC-901 so a line from the next section cannot satisfy it.
cat > "$phase/decisions.md" <<'EOF'
### DEC-901: release review-2 repairs to ordered owners
- **Status:** ACCEPTED
- **Ledger binding:**
  - Repair owners: task-6, task-2
### DEC-903: reversed owner order
- **Ledger binding:**
  - Repair owners: task-2, task-6
### DEC-902: duplicate owner list
- **Ledger binding:**
  - Repair owners: task-6, task-6
### DEC-904: re-own the trigger task under its retained signature
- **Ledger binding:**
  - Repair owners: task-6, task-3
EOF
cat > "$phase/triage-log.md" <<'EOF'
## TRG-1-001 - review-2 BLOCK
- **Originating task:** task-3
- **Trigger signature:** task-3:qa-rework:review-2-block
- **Attempt count:** 1
- **Repair category:** N/A
- **Repair authority:** N/A
- **Verdict:** BLOCK
- **If BLOCK:** ESC-1-001
EOF

recovery() {
  cat <<EOF
### 2026-09-06 $2 TASK-RECOVERY $1
- Source: $3
- Category: qa-rework
- Repair mode: task-scope
- Repair authority: $4
- Triage disposition: N/A
- Trigger signature: $5
- Attempt count: $6
- Control contract: N/A
- Control mismatch: N/A
- Allowed files: ordinary task scope
- Forbidden changes: N/A
- Originating report: N/A
- Exact verification: N/A
- Control repair base: N/A
- Consecutive no-progress attempts: 0/2
- Checkpoint: synthetic checkpoint
- Phase worktree: synthetic fixture
- Temporary repair worktree: N/A
- Hidden-index flags: none
- Resolution evidence: PENDING
- Next action: repair
- Human attention required: no
EOF
}
repair() {
  cat <<EOF
### 2026-09-06 $2 REPAIR $1
- Trigger signature: $3
- Attempt count: $4
- Starting checkpoint: synthetic checkpoint
- Ending checkpoint: findings fixed
- Control repair base: N/A
- Control repair head: N/A
- Control scope check: N/A
- Phase-branch integration: N/A
- Outcome: RESOLVED
- Resolution evidence: exact task gate PASS
- Next checkpoint: TASK COMPLETE
EOF
}
complete() {
  cat <<EOF
### 2026-09-06 $2 TASK $1 COMPLETE
- Report: project-docs/docs/phases/phase-1/tasks/$1-done.md
- report-lint: PASS - OUTCOME COMPLETE (terminated exit_code: 0)
- Commit: synthetic - verified
- QA-plan commands beyond the task spec: none required
EOF
}
review() {
  cat <<EOF
### 2026-09-06 $2 REVIEW $1
- Path: project-docs/docs/phases/phase-1/qa/$1.md
- Tasks covered: task-2, task-3, task-6
- Outcome: $3
- Owning task: $4
- Re-review of: $5
- Must-fix: 1 / Improvement: 0 / blocker candidates: 0
EOF
}

cat > "$scratch/start" <<'EOF'
### 2026-09-06 10:00 RUN-START
- Sequence contract: v2
- Recovery directive schema: v2
- Branch: phase-1
- ESC baseline (Open at run start): none
- Warm-up/autonomous-resume state: autonomous
- Starting task: task-3
EOF
review review-1 10:01 REPAIR task-3 none > "$scratch/review1"
{
  recovery task-3 10:02 review-1 'QA review' review-1:qa-rework:task-3-MF-1 1
  repair task-3 10:03 review-1:qa-rework:task-3-MF-1 1
  complete task-3 10:04
} > "$scratch/fix3"
review review-2 10:05 BLOCK none review-1 > "$scratch/review2"
cat > "$scratch/triage" <<'EOF'
### 2026-09-06 10:06 TRIAGE TRG-1-001
- Trigger: QA BLOCK on the review-1 re-review
- Originating task: task-3
- Trigger signature: task-3:qa-rework:review-2-block
- Attempt count: 1
- Repair category: N/A
- Repair authority: N/A
- Checkpoint: review-2 committed
- Verdict: BLOCK (escalated as ESC-1-001)
- Resolution evidence: protected candidates need a decision
- Triage log: TRG-1-001
EOF
cat > "$scratch/stop" <<'EOF'
### 2026-09-06 10:07 STOP
- Cause: review-2 / TRG-1-001 QA BLOCK
- Triage: TRG-1-001 (BLOCK, escalated as ESC-1-001)
- Protected category: product semantics
- Next Human or external-authority action: record the releasing decision
EOF
cat > "$scratch/resume" <<'EOF'
### 2026-09-06 10:08 RUN-RESUME
- Sequence contract: v2
- Recovery directive schema: v2
- Prior run/session or STOP: STOP after TRG-1-001 (ESC-1-001)
- Governing TRG: TRG-1-001
- Governing ESC: ESC-1-001
- Trigger signature: task-3:qa-rework:review-2-block
- Attempt count: 2
- Checkpoint: synthetic checkpoint, clean
- Resolution evidence: DEC-901 releases the review-2 repairs to ordered owners
- Resolution authority: DEC-901
- Continuing task: task-6
- Repair owners: task-6, task-2
- Next recovery: task-6 | qa-rework | task-scope | review-2:qa-rework:task-6-MF-1 | 1
EOF
cat > "$scratch/normal-resume" <<'EOF'
### 2026-09-06 10:08 RUN-RESUME
- Sequence contract: v2
- Recovery directive schema: v2
- Prior run/session or STOP: prior runner; no persisted STOP
- Governing TRG: N/A
- Governing ESC: N/A
- Trigger signature: N/A
- Attempt count: N/A
- Checkpoint: synthetic checkpoint, clean
- Resolution evidence: same run continued
- Resolution authority: N/A
- Continuing task: task-6
- Repair owners: task-6, task-2
- Next recovery: N/A
EOF
recovery task-6 10:09 'TRG-1-001 DEC-901' 'TRG-1-001 DEC-901' \
  review-2:qa-rework:task-6-MF-1 1 > "$scratch/rec6"
repair task-6 10:10 review-2:qa-rework:task-6-MF-1 1 > "$scratch/rep6"
complete task-6 10:11 > "$scratch/done6"
recovery task-2 10:12 review-2 'QA review' review-2:qa-rework:task-2-MF-2 1 > "$scratch/rec2"
repair task-2 10:13 review-2:qa-rework:task-2-MF-2 1 > "$scratch/rep2"
complete task-2 10:14 > "$scratch/done2"
review review-3 10:15 PASS none review-2 > "$scratch/pass3"

set_ledger() {
  : > "$ledger"
  for part in "$@"; do
    cat "$scratch/$part" >> "$ledger"
  done
}
stopped_prefix() {
  set_ledger start review1 fix3 review2 triage stop
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

# (a) The full route: resume, ordered owners, linked PASS.
stopped_prefix
code=0
output=$("$validator" "$ledger" 2>&1) || code=$?
if [ "$code" -eq 0 ] && [ "$output" = "VALID state=STOPPED" ]; then
  pass 'BLOCK re-review prefix stops'
else
  fail 'BLOCK re-review prefix stops'
fi
accept 'STOP resume installs DEC-bound repair owners' "$scratch/resume"
accept 'first owner recovers under its review-form signature' "$scratch/rec6"
accept 'first owner repair resolves' "$scratch/rep6" COMPLETION_PENDING
accept 'first owner completes' "$scratch/done6"
accept 'second owner recovers in list order' "$scratch/rec2"
accept 'second owner repair resolves' "$scratch/rep2" COMPLETION_PENDING
accept 'second owner completes' "$scratch/done2" REREVIEW_PENDING
accept 'linked PASS closes the BLOCK re-review' "$scratch/pass3" ADVANCEABLE

# (a) A later owner may be the trigger task under its retained signature.
stopped_prefix
sed -e 's/DEC-901/DEC-904/g' -e 's/Repair owners: task-6, task-2/Repair owners: task-6, task-3/' \
  "$scratch/resume" > "$scratch/resume-retained"
sed 's/DEC-901/DEC-904/g' "$scratch/rec6" > "$scratch/rec6-retained"
accept 'resume binds an owner list that re-owns the trigger task' "$scratch/resume-retained"
accept 'first owner recovers under DEC-904' "$scratch/rec6-retained"
accept 'first owner repair resolves under DEC-904' "$scratch/rep6" COMPLETION_PENDING
accept 'first owner completes under DEC-904' "$scratch/done6"
recovery task-3 10:12 'TRG-1-001 DEC-904' 'DEC-904' task-3:qa-rework:review-2-block 2 > "$scratch/rec3"
accept 'trigger task recovers under its retained signature' "$scratch/rec3"

# (b) The owner list must be byte-identical to the DEC section line.
stopped_prefix
sed 's/Repair owners: task-6, task-2/Repair owners: task-6,task-2/' "$scratch/resume" > "$scratch/bad"
reject 'owner list spacing differs from DEC line' "$scratch/bad" \
  'Repair owners do not match a Repair owners line in its Resolution authority decision'
sed -e 's/Repair owners: task-6, task-2/Repair owners: task-2, task-6/' \
    -e 's/Continuing task: task-6/Continuing task: task-2/' "$scratch/resume" > "$scratch/bad"
reject 'owner list from the next DEC section is not bound' "$scratch/bad" \
  'Repair owners do not match a Repair owners line in its Resolution authority decision'
sed 's/Resolution authority: DEC-901/Resolution authority: verified-commit:0123456789abcdef0123456789abcdef01234567/' \
  "$scratch/resume" > "$scratch/bad"
reject 'owner list requires a DEC resolution authority' "$scratch/bad" \
  'Repair owners do not match a Repair owners line in its Resolution authority decision'

# (c) The STOP must come from a BLOCK re-review of a completed open review.
set_ledger start review1 review2 triage stop
reject 'owner list rejected while the open review owner is incomplete' "$scratch/resume" \
  'Repair owners require a STOP from a BLOCK re-review of a completed open QA repair'
set_ledger start review1 fix3
sed 's/Re-review of: review-1/Re-review of: none/' "$scratch/review2" >> "$ledger"
cat "$scratch/triage" "$scratch/stop" >> "$ledger"
reject 'owner list rejected after an unlinked BLOCK review' "$scratch/resume" \
  'Repair owners require a STOP from a BLOCK re-review of a completed open QA repair'
set_ledger start review1 fix3 triage stop
reject 'owner list rejected after a STOP without a review BLOCK' "$scratch/resume" \
  'Repair owners require a STOP from a BLOCK re-review of a completed open QA repair'

# (d) The Continuing task is the first Repair owner.
stopped_prefix
sed 's/Continuing task: task-6/Continuing task: task-2/' "$scratch/resume" > "$scratch/bad"
reject 'Continuing task must be the first owner' "$scratch/bad" \
  'RUN-RESUME Continuing task must be the first Repair owner'

# (e) Owners recover one at a time, in list order.
stopped_prefix
cat "$scratch/resume" >> "$ledger"
sed 's/Repair authority: QA review/Repair authority: TRG-1-001 DEC-901/;s/Source: review-2/Source: TRG-1-001 DEC-901/' \
  "$scratch/rec2" > "$scratch/bad"
reject 'second owner cannot take the resume recovery' "$scratch/bad" \
  'TASK-RECOVERY does not preserve its RUN-RESUME Repair owners checkpoint'
sed 's/Attempt count: 1/Attempt count: 2/' "$scratch/rec6" > "$scratch/bad"
reject 'first owner review-form recovery starts at attempt 1' "$scratch/bad" \
  'TASK-RECOVERY does not preserve its RUN-RESUME Repair owners checkpoint'
cat "$scratch/rec6" "$scratch/rep6" >> "$ledger"
reject 'second owner cannot recover before the first completes' "$scratch/rec2" \
  'requires recovery task task-6 before task-2'

# (f) A linked PASS waits for every owner.
stopped_prefix
cat "$scratch/resume" "$scratch/rec6" "$scratch/rep6" "$scratch/done6" >> "$ledger"
reject 'PASS before every owner completes' "$scratch/pass3" \
  'QA re-review passed before owning task completed'

# (g) Repair owners belong only to a STOP resume.
set_ledger start review1 fix3
reject 'normal resume cannot carry Repair owners' "$scratch/normal-resume" \
  'Repair owners are valid only on a RUN-RESUME after STOP'

# (h) Owners stay distinct, and the field occurs once.
stopped_prefix
sed -e 's/DEC-901/DEC-902/g' -e 's/Repair owners: task-6, task-2/Repair owners: task-6, task-6/' \
  "$scratch/resume" > "$scratch/bad"
reject 'duplicate Repair owners' "$scratch/bad" 'duplicate owning task task-6'
sed '/^- Repair owners: /p' "$scratch/resume" > "$scratch/bad"
reject 'Repair owners field repeated' "$scratch/bad" \
  'Repair owners must occur exactly once on a schema v2 RUN-RESUME'

# (i) Without a resume, recovery after the BLOCK re-review is rejected legibly.
set_ledger start review1 fix3 review2
reject 'recovery straight after the BLOCK re-review' "$scratch/rec6" \
  'QA repair owners of review-1 are complete; a linked re-review or a RUN-RESUME with Repair owners is required before recovery task task-6'

if [ "$failures" -ne 0 ]; then
  echo "REVIEW BLOCK OWNERS FAIL $failures checks"
  exit 1
fi
echo 'REVIEW BLOCK OWNERS PASS all checks'
