#!/bin/sh
set -eu

script_dir=$(CDPATH= cd "$(dirname "$0")" && pwd)
repo_root=$(CDPATH= cd "$script_dir/../.." && pwd)
cd "$repo_root"

failures=0

pass() {
  echo "PASS $1"
}

fail() {
  echo "FAIL $1" >&2
  failures=$((failures + 1))
}

check_file() {
  path=$1
  label=$2
  if [ -f "$path" ]; then
    pass "$label: $path"
  else
    fail "$label missing: $path"
  fi
}

frontmatter_value() {
  key=$1
  file=$2
  awk -v key="$key" '
    NR == 1 && $0 == "---" {
      in_frontmatter = 1
      next
    }
    in_frontmatter && $0 == "---" {
      exit
    }
    in_frontmatter {
      prefix = key ":"
      if (index($0, prefix) == 1) {
        value = substr($0, length(prefix) + 1)
        gsub(/^[ \t]+|[ \t]+$/, "", value)
        gsub(/^"/, "", value)
        gsub(/"$/, "", value)
        print value
        exit
      }
    }
  ' "$file"
}

check_contains() {
  file=$1
  pattern=$2
  label=$3
  if grep -Eq "$pattern" "$file"; then
    pass "$label: $file"
  else
    fail "$label missing in $file"
  fi
}

check_absent() {
  file=$1
  pattern=$2
  label=$3
  if grep -Eq "$pattern" "$file"; then
    fail "$label still present in $file"
  else
    pass "$label absent: $file"
  fi
}

required_skills="validation-contract qa-evidence phase-wrap-up tdd diagnose zoom-out phase-planning task-pack domain-language-grill architecture-review"

for skill in $required_skills; do
  check_file "project-docs/skills/$skill/SKILL.md" "required canonical skill"
done

for skill_dir in project-docs/skills/*; do
  [ -d "$skill_dir" ] || continue
  skill=$(basename "$skill_dir")
  skill_file="$skill_dir/SKILL.md"
  check_file "$skill_file" "canonical skill"

  if [ -f "$skill_file" ]; then
    name=$(frontmatter_value name "$skill_file")
    description=$(frontmatter_value description "$skill_file")

    if [ -n "$name" ]; then
      pass "frontmatter name present: $skill"
    else
      fail "frontmatter name missing: $skill_file"
    fi

    if [ "$name" = "$skill" ]; then
      pass "frontmatter name matches directory: $skill"
    else
      fail "frontmatter name '$name' does not match directory '$skill'"
    fi

    if [ -n "$description" ]; then
      pass "frontmatter description present: $skill"
    else
      fail "frontmatter description missing: $skill_file"
    fi

    if grep -Eq '^Provenance: (Workflow-Extracted Skill|External-Adapted Skill)(\.|,|;|$)' "$skill_file"; then
      pass "provenance line valid: $skill"
    else
      fail "valid Provenance line missing: $skill_file"
    fi
  fi
done

if output=$(project-docs/tools/sync-skills.sh --check 2>&1); then
  printf '%s\n' "$output"
  pass "generated skill copies match canonical sources"
else
  printf '%s\n' "$output" >&2
  fail "generated skill copies are missing or drifted"
fi

required_agent_sources="tech-lead phase-runner developer qa-lead qa-reviewer tl-lite"

for role in $required_agent_sources; do
  check_file "project-docs/agents/$role.md" "canonical agent source"
done

if output=$(project-docs/tools/sync-agents.sh --check 2>&1); then
  printf '%s\n' "$output"
  pass "generated agent prompts match canonical sources"
else
  printf '%s\n' "$output" >&2
  fail "generated agent prompts are missing or drifted"
fi

agent_files="
AGENTS.md
CLAUDE.md
.codex/config.toml
.codex/agents/tech_lead.toml
.codex/agents/phase_runner.toml
.codex/agents/developer.toml
.codex/agents/qa_lead.toml
.codex/agents/qa_reviewer.toml
.codex/agents/tl_lite.toml
.claude/agents/tech-lead.md
.claude/agents/phase-runner.md
.claude/agents/developer.md
.claude/agents/qa-lead.md
.claude/agents/qa-reviewer.md
.claude/agents/tl-lite.md
"

for path in $agent_files; do
  check_file "$path" "agent prompt/config"
done

# Codex gate-running roles must retain the nested command session handle for
# commands that outlive the initial exec yield. Generated prompts are checked
# alongside canonical sources so provider drift cannot silently restore the
# false "terminated without output" failure mode from ESC-20-003.
codex_gate_session_surfaces="
project-docs/agents/developer.md
project-docs/agents/qa-lead.md
project-docs/agents/qa-reviewer.md
.codex/agents/developer.toml
.codex/agents/qa_lead.toml
.codex/agents/qa_reviewer.toml
"

for file in $codex_gate_session_surfaces; do
  check_contains "$file" 'JSON\.stringify' "Codex full exec-result preservation"
  check_contains "$file" 'session_id' "Codex nested command session handling"
  check_contains "$file" 'write_stdin' "Codex long-command polling"
  check_contains "$file" 'exit_code' "Codex observed command termination"
done

for file in project-docs/agents/phase-runner.md .codex/agents/phase_runner.toml; do
  check_contains "$file" 'exit_code' "Phase Runner command-termination evidence"
  check_contains "$file" 'outer.*execution cell|execution cell.*outer' "Phase Runner wrapper-completion guardrail"
done

check_file "project-docs/docs/playbook.md" "operator playbook"
check_file "project-docs/templates/hooks/claude-settings.example.json" "Claude hook example"
check_file "project-docs/templates/hooks/codex-hooks-notes.md" "Codex hook notes"

template_files="
project-docs/templates/validation-contract-template.md
project-docs/templates/task-spec-template.md
project-docs/templates/completion-report-template.md
project-docs/templates/qa-plan-template.md
project-docs/templates/mid-phase-review-template.md
project-docs/templates/final-review-template.md
project-docs/templates/human-checklist-template.md
project-docs/templates/phase-brief-template.md
project-docs/templates/phase-decisions-template.md
project-docs/templates/phase-runner-prompt.md
project-docs/templates/run-ledger-template.md
project-docs/templates/tl-escalations-template.md
project-docs/templates/triage-log-template.md
project-docs/templates/threat-model-template.md
project-docs/templates/pre-release-checklist.md
project-docs/templates/held-out-scenario-template.yaml
project-docs/templates/audits/audit-prompts.md
"

for path in $template_files; do
  check_file "$path" "template"
done

check_contains "project-docs/templates/task-spec-template.md" 'Validation Contract|VAL-NNN|validation-contract' "task template validation-contract reference"
check_contains "project-docs/templates/completion-report-template.md" 'Validation Contract|VAL-NNN|validation-contract' "completion template validation-contract reference"
check_contains "project-docs/templates/validation-contract-template.md" 'VAL-NNN|validation contract' "validation contract template reference"
check_contains "project-docs/templates/qa-plan-template.md" 'Validation Contract|VAL-NNN|validation-contract' "QA plan template validation-contract reference"
check_contains "project-docs/templates/mid-phase-review-template.md" 'Validation Contract|VAL-NNN|validation-contract' "mid-phase review template validation-contract reference"
check_contains "project-docs/templates/final-review-template.md" 'Validation Contract|VAL-NNN|validation-contract' "final review template validation-contract reference"
check_contains "project-docs/templates/human-checklist-template.md" 'VAL-NNN' "human checklist validation assertion reference"
check_contains "project-docs/templates/phase-brief-template.md" 'Validation Contract|VAL-NNN|validation-contract' "phase brief validation-contract reference"
check_contains "project-docs/templates/phase-runner-prompt.md" 'validation-contract|VAL-NNN' "phase runner prompt validation-contract reference"

# --- Deterministic verification tooling ---

check_file "project-docs/tools/report-lint.sh" "report-lint tool"
if [ -x "project-docs/tools/report-lint.sh" ]; then
  pass "report-lint.sh is executable"
else
  fail "report-lint.sh is not executable"
fi

if output=$(project-docs/tools/report-lint.sh --self-test 2>&1); then
  pass "report-lint self-test green"
else
  printf '%s\n' "$output" >&2
  fail "report-lint self-test failed"
fi

check_file "project-docs/tools/workflow-route.sh" "workflow route tool"
if [ -x "project-docs/tools/workflow-route.sh" ]; then
  pass "workflow-route.sh is executable"
else
  fail "workflow-route.sh is not executable"
fi

check_file "project-docs/tools/drills/run-drills.sh" "evidence-chain drills harness"
if [ -x "project-docs/tools/drills/run-drills.sh" ]; then
  pass "run-drills.sh is executable"
else
  fail "run-drills.sh is not executable"
fi

# Phase Runner surfaces must reference report-lint and the run ledger instead
# of prose-parsing/chat-state equivalents.

runner_surfaces="
project-docs/specs/09_AI_WORKFLOW.md
project-docs/templates/phase-runner-prompt.md
.codex/agents/phase_runner.toml
.claude/agents/phase-runner.md
"

for file in $runner_surfaces; do
  check_contains "$file" 'report-lint\.sh' "report-lint reference"
  check_contains "$file" 'workflow-route\.sh' "workflow route reference"
  check_contains "$file" 'run-ledger' "run-ledger reference"
  check_contains "$file" 'first-step triager' "TL-Lite triage-gated stop invariant"
  check_contains "$file" 'no TRIAGE entry exists' "untriaged-STOP resume rule"
done

report_lint_check_surfaces="
project-docs/specs/09_AI_WORKFLOW.md
project-docs/templates/phase-runner-prompt.md
project-docs/tools/README.md
.codex/agents/phase_runner.toml
.claude/agents/phase-runner.md
"

for file in $report_lint_check_surfaces; do
  check_contains "$file" 'uncommitted modifications|committed and clean' "report-lint committed-report check"
done

check_contains "project-docs/templates/run-ledger-template.md" '^- Triage:' "run-ledger STOP triage field"

codex_runner_surfaces="
project-docs/specs/09_AI_WORKFLOW.md
project-docs/templates/phase-runner-prompt.md
.codex/agents/phase_runner.toml
"

for file in $codex_runner_surfaces; do
  check_contains "$file" 'agent_type' "Codex custom-agent selector guardrail"
  check_contains "$file" 'default.*worker.*explorer|default.*explorer.*worker' "Codex generic-agent guardrail"
done

# --- Spec section references must resolve ---
# Any "§N (Heading)" or "§N.M (Heading)" reference in canonical agent
# sources, generated prompts, or skills must match a real heading in
# 09_AI_WORKFLOW.md. This keeps slim-prompt pointers from dangling when the
# spec is renumbered.

section_ref_files=$(find project-docs/agents project-docs/skills .agents/skills .claude/agents .claude/skills .codex/agents -type f \( -name '*.md' -o -name '*.toml' \) 2>/dev/null | sort)
section_refs=$(grep -hoE "§[0-9]+(\.[0-9]+)* \([^)]+\)" $section_ref_files 2>/dev/null | sort -u || true)

if [ -n "$section_refs" ]; then
  while IFS= read -r ref; do
    [ -n "$ref" ] || continue
    num=${ref%% *}
    num=${num#§}
    heading=${ref#*\(}
    heading=${heading%\)}
    num_esc=$(printf '%s' "$num" | sed 's/\./\\./g')
    case "$num" in
      *.*) pattern="^#{3,4} $num_esc $heading\$" ;;
      *) pattern="^## $num_esc\\. $heading\$" ;;
    esac
    if grep -Eq "$pattern" project-docs/specs/09_AI_WORKFLOW.md; then
      pass "spec section reference resolves: §$num ($heading)"
    else
      fail "spec section reference does not resolve in 09_AI_WORKFLOW.md: §$num ($heading)"
    fi
  done <<EOF
$section_refs
EOF
fi

# --- Skill-declared task headings must exist in the task spec template ---
# Catches heading drift between the task-pack skill and the template.

task_pack_headings=$(awk '
  $0 == "## Required Task Headings" { in_section = 1; next }
  in_section && /^## / { exit }
  in_section {
    line = $0
    while (match(line, /`#[^`]*`/)) {
      print substr(line, RSTART + 1, RLENGTH - 2)
      line = substr(line, RSTART + RLENGTH)
    }
  }
' project-docs/skills/task-pack/SKILL.md)

if [ -z "$task_pack_headings" ]; then
  fail "task-pack skill: no backticked headings under '## Required Task Headings'"
else
  while IFS= read -r heading; do
    [ -n "$heading" ] || continue
    if grep -Fxq "$heading" project-docs/templates/task-spec-template.md; then
      pass "task-pack heading exists in task-spec template: $heading"
    else
      fail "task-pack heading missing from task-spec template: $heading"
    fi
  done <<EOF
$task_pack_headings
EOF
fi

# --- Completion-report heading concordance ---
# The 11 required headings must appear in the template, the workflow spec,
# and both Developer prompts.

completion_heading_files="
project-docs/templates/completion-report-template.md
project-docs/specs/09_AI_WORKFLOW.md
.claude/agents/developer.md
.codex/agents/developer.toml
"

while IFS= read -r heading; do
  for file in $completion_heading_files; do
    if grep -Fq "$heading" "$file"; then
      pass "completion-report heading present: $file: $heading"
    else
      fail "completion-report heading missing: $file: $heading"
    fi
  done
done <<'EOF'
## Workflow outcome
## Summary of what was implemented
## Files created or modified
## Documentation updated
## Validation Contract Coverage
## Tests run
## Exact quality-gate results
## Deviations from spec
## Concerns, technical debt, or notes for the Tech Lead
## Dependency-budget note
## Commit hash
## Commit subject
EOF

# --- Recovery-state concordance ---

repair_verdict_files="
project-docs/specs/09_AI_WORKFLOW.md
project-docs/agents/phase-runner.md
project-docs/agents/tl-lite.md
project-docs/templates/phase-runner-prompt.md
.claude/agents/phase-runner.md
.claude/agents/tl-lite.md
.codex/agents/phase_runner.toml
.codex/agents/tl_lite.toml
"

for file in $repair_verdict_files; do
  check_contains "$file" 'VERDICT: REPAIR' "three-way recovery verdict"
done

same_task_recovery_files="
AGENTS.md
project-docs/specs/09_AI_WORKFLOW.md
project-docs/agents/developer.md
project-docs/agents/phase-runner.md
project-docs/templates/phase-runner-prompt.md
.claude/agents/developer.md
.claude/agents/phase-runner.md
.codex/agents/developer.toml
.codex/agents/phase_runner.toml
"

for file in $same_task_recovery_files; do
  check_contains "$file" 'same-task.*continuation|continuation.*same task' "same-task recovery dispatch"
  check_contains "$file" 'incomplete.*not.*blocker|not.*blocker.*incomplete' "incomplete-work blocker rejection"
done

for file in project-docs/templates/run-ledger-template.md project-docs/specs/09_AI_WORKFLOW.md project-docs/agents/phase-runner.md project-docs/templates/phase-runner-prompt.md; do
  check_contains "$file" 'TASK-RECOVERY' "task recovery ledger event"
  check_contains "$file" 'RUN-RESUME' "resolved run resume event"
done

qa_recovery_files="
project-docs/specs/09_AI_WORKFLOW.md
project-docs/agents/qa-reviewer.md
project-docs/templates/mid-phase-review-template.md
.claude/agents/qa-reviewer.md
.codex/agents/qa_reviewer.toml
"

for file in $qa_recovery_files; do
  check_contains "$file" 'PASS.*REPAIR.*BLOCK' "typed QA review outcome"
  check_contains "$file" 'owning task' "QA recovery ownership"
  check_contains "$file" 're-review' "independent QA re-review"
done

check_contains "project-docs/tools/report-lint.sh" 'OUTCOME.*RECOVERY_REQUIRED' "typed report-lint recovery outcome"
check_contains "project-docs/tools/report-lint.sh" 'PASS result missing observed termination' "report-lint termination enforcement"
check_absent "project-docs/tools/report-lint.sh" 'orchestrator must still stop' "unconditional blocker stop"

# --- Current-state file line budgets (from 09_AI_WORKFLOW.md section 3.5) ---

check_line_budget() {
  file=$1
  tripwire=$2
  label=$3
  if [ ! -f "$file" ]; then
    pass "line budget (file absent, ok in template state): $label"
    return
  fi
  lines=$(wc -l < "$file")
  if [ "$lines" -le "$tripwire" ]; then
    pass "line budget ($lines <= $tripwire): $label"
  else
    fail "line budget exceeded ($lines > $tripwire tripwire): $label"
  fi
}

check_line_budget "project-docs/docs/current/phase_brief.md" 200 "phase_brief.md"
check_line_budget "project-docs/docs/current/decisions_active.md" 200 "decisions_active.md"
check_line_budget "project-docs/docs/current/progress_current.md" 133 "progress_current.md"
check_line_budget "project-docs/docs/current/tl_escalations_current.md" 107 "tl_escalations_current.md"

# --- TL-Lite verdict and output format concordance ---
# These literal strings must appear in both provider prompts and the spec.

tl_lite_verdict_files=".claude/agents/tl-lite.md .codex/agents/tl_lite.toml project-docs/specs/09_AI_WORKFLOW.md"

while IFS= read -r literal; do
  for file in $tl_lite_verdict_files; do
    if grep -Fq "$literal" "$file"; then
      pass "TL-Lite verdict '$literal' present: $file"
    else
      fail "TL-Lite verdict '$literal' missing: $file"
    fi
  done
done <<'EOF'
VERDICT: RESUME
VERDICT: REPAIR
VERDICT: BLOCK
EOF

tl_lite_prompt_files=".claude/agents/tl-lite.md .codex/agents/tl_lite.toml"

while IFS= read -r literal; do
  for file in $tl_lite_prompt_files; do
    if grep -Fq "$literal" "$file"; then
      pass "TL-Lite output format '$literal' present: $file"
    else
      fail "TL-Lite output format '$literal' missing: $file"
    fi
  done
done <<'EOF'
DECISION (local):
ESCALATED: ESC-
EOF

# --- TL-Lite escalation trigger concordance ---
# Each trigger phrase must appear in the spec and both provider prompts.

escalation_files=".claude/agents/tl-lite.md .codex/agents/tl_lite.toml project-docs/specs/09_AI_WORKFLOW.md"

while IFS= read -r trigger; do
  for file in $escalation_files; do
    if grep -Fq "$trigger" "$file"; then
      pass "escalation trigger present: $file: $trigger"
    else
      fail "escalation trigger missing: $file: $trigger"
    fi
  done
done <<'EOF'
Changes a public interface
Changes the data model or persistence contract
Affects more than one task or subsystem
Introduces a dependency, tool, quality gate, or workflow
Exposes ambiguity in product intent or acceptance criteria
EOF

# --- Workflow guardrail concordance ---

exact_command_files="
project-docs/specs/09_AI_WORKFLOW.md
project-docs/templates/task-spec-template.md
project-docs/templates/completion-report-template.md
project-docs/templates/phase-runner-prompt.md
.codex/agents/developer.toml
.codex/agents/phase_runner.toml
.codex/agents/qa_reviewer.toml
.codex/agents/tech_lead.toml
.claude/agents/developer.md
.claude/agents/phase-runner.md
.claude/agents/qa-reviewer.md
.claude/agents/tech-lead.md
"

for file in $exact_command_files; do
  check_contains "$file" 'Exact Commands To Run' "exact-command task gate guardrail"
done

command_coverage_files="
AGENTS.md
project-docs/specs/08_QUALITY_PLAN.md
project-docs/specs/09_AI_WORKFLOW.md
project-docs/templates/task-spec-template.md
project-docs/templates/completion-report-template.md
project-docs/templates/mid-phase-review-template.md
.codex/agents/developer.toml
.codex/agents/phase_runner.toml
.codex/agents/qa_lead.toml
.codex/agents/qa_reviewer.toml
.claude/agents/developer.md
.claude/agents/phase-runner.md
.claude/agents/qa-lead.md
.claude/agents/qa-reviewer.md
"

for file in $command_coverage_files; do
  check_contains "$file" 'silently skipped|skipped gates without' "missing-command guardrail"
  check_contains "$file" 'narrower substituted' "narrower-substitution guardrail"
done

for file in project-docs/specs/09_AI_WORKFLOW.md project-docs/templates/phase-runner-prompt.md .codex/agents/phase_runner.toml .claude/agents/phase-runner.md; do
  check_contains "$file" 'git rev-parse --verify' "commit hash verification guardrail"
done

for file in project-docs/specs/08_QUALITY_PLAN.md project-docs/templates/validation-contract-template.md project-docs/skills/validation-contract/SKILL.md .codex/agents/qa_lead.toml .claude/agents/qa-lead.md; do
  check_contains "$file" 'one observable' "one-observable validation guardrail"
  check_contains "$file" 'positive-only|Positive-only' "conditional evidence guardrail"
done

# --- Operating-mode concordance ---
# The mode policy must exist in the spec, the operator playbook, the phase
# brief template, glossary, Tech Lead prompts, and canonical planning skill.

mode_files="
project-docs/specs/09_AI_WORKFLOW.md
project-docs/docs/playbook.md
project-docs/templates/phase-brief-template.md
project-docs/specs/GLOSSARY.md
project-docs/skills/phase-planning/SKILL.md
.claude/agents/tech-lead.md
.codex/agents/tech_lead.toml
"

for file in $mode_files; do
  check_contains "$file" '[Oo]perating [Mm]ode' "operating-mode reference"
done

for file in project-docs/specs/09_AI_WORKFLOW.md project-docs/docs/playbook.md; do
  check_contains "$file" 'promote the phase one mode' "mode promotion rule"
done

# --- Role capability mapping concordance ---
# Each role name from the spec table must appear in at least one provider
# agent prompt directory. This catches renamed or deleted agent files.

while IFS= read -r role; do
  found_claude=false
  found_codex=false
  for f in .claude/agents/*.md; do
    if grep -Fq "$role" "$f" 2>/dev/null; then
      found_claude=true
      break
    fi
  done
  for f in .codex/agents/*.toml; do
    if grep -Fq "$role" "$f" 2>/dev/null; then
      found_codex=true
      break
    fi
  done
  if $found_claude; then
    pass "role '$role' referenced in Claude agents"
  else
    fail "role '$role' not referenced in any Claude agent prompt"
  fi
  if $found_codex; then
    pass "role '$role' referenced in Codex agents"
  else
    fail "role '$role' not referenced in any Codex agent prompt"
  fi
done <<'EOF'
Tech Lead
QA Lead
Developer
QA Reviewer
TL-Lite
Phase Runner
EOF

# --- Codex project-local config key restrictions ---

unsupported_keys="openai_base_url chatgpt_base_url model_provider model_providers notify profile profiles experimental_realtime_ws_base_url otel"
config_keys=$(awk '
  /^[ \t]*#/ || /^[ \t]*$/ { next }
  /^[ \t]*\[/ {
    line = $0
    gsub(/^[ \t]*\[+/, "", line)
    gsub(/\]+[ \t]*$/, "", line)
    split(line, parts, ".")
    print parts[1]
    next
  }
  /^[ \t]*[A-Za-z0-9_-]+[ \t]*=/ {
    line = $0
    sub(/^[ \t]*/, "", line)
    sub(/[ \t]*=.*/, "", line)
    print line
  }
' .codex/config.toml)

for unsupported in $unsupported_keys; do
  if printf '%s\n' "$config_keys" | grep -Fxq "$unsupported"; then
    fail "unsupported project-local Codex config key present: $unsupported"
  else
    pass "unsupported Codex config key absent: $unsupported"
  fi
done

if [ "$failures" -eq 0 ]; then
  echo "PASS workflow-lint completed with no failures"
else
  echo "FAIL workflow-lint found $failures failure(s)" >&2
fi

exit "$failures"
