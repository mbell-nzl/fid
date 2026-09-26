#!/bin/sh
# workflow-sequence-lint.sh - validate or transactionally append a ledger event.
set -eu

if [ "$#" -lt 1 ] || [ "$#" -gt 2 ] || [ ! -f "$1" ]; then
  echo "Usage: workflow-sequence-lint.sh <run-ledger.md> [--append-stdin]" >&2
  exit 2
fi

ledger=$1
mode=validate
if [ "$#" -eq 2 ]; then
  [ "$2" = "--append-stdin" ] || {
    echo "Usage: workflow-sequence-lint.sh <run-ledger.md> [--append-stdin]" >&2
    exit 2
  }
  mode=append
fi

validate_sequence() {
content_path=$1
authority_path=${2:-$1}
phase_dir=$(dirname "$authority_path")
phase_name=$(basename "$phase_dir")
phase_number=""
decisions_file=""
triage_file=""
case "$phase_name" in
  phase-*)
    candidate_phase=${phase_name#phase-}
    case "$candidate_phase" in
      ''|*[!0-9]*) ;;
      *)
        phase_number=$candidate_phase
        decisions_file="$phase_dir/decisions.md"
        triage_file="$phase_dir/triage-log.md"
        ;;
    esac
    ;;
esac

awk -v phase_number="$phase_number" -v decisions_file="$decisions_file" \
    -v triage_file="$triage_file" '
function invalid(message) {
  print "INVALID " message > "/dev/stderr"
  errors++
}

function references_id(text, id) {
  return id != "" && text ~ ("(^|[^[:alnum:]_-])" id "([^[:alnum:]_-]|$)")
}

function references_exact_token(text, token, rest, position, before, after) {
  if (token == "") {
    return 0
  }
  rest = text
  while ((position = index(rest, token)) > 0) {
    before = position > 1 ? substr(rest, position - 1, 1) : ""
    after = substr(rest, position + length(token), 1)
    if ((before == "" || before !~ /[[:alnum:]_.:-]/) &&
        (after == "" || after !~ /[[:alnum:]_.:-]/)) {
      return 1
    }
    rest = substr(rest, position + 1)
  }
  return 0
}

function bounded_id_at(text, start, token_length, before, after) {
  before = start > 1 ? substr(text, start - 1, 1) : ""
  after = substr(text, start + token_length, 1)
  return (before == "" || before !~ /[[:alnum:]_-]/) &&
         (after == "" || after !~ /[[:alnum:]_-]/)
}

function decision_id(text, rest, offset, id) {
  rest = text
  while (match(rest, /DEC-[0-9]+/)) {
    id = substr(rest, RSTART, RLENGTH)
    if (bounded_id_at(rest, RSTART, RLENGTH)) {
      return id
    }
    rest = substr(rest, RSTART + RLENGTH)
  }
  return ""
}

function all_decisions_exist(text, rest, id, found) {
  rest = text
  while (match(rest, /DEC-[0-9]+/)) {
    id = substr(rest, RSTART, RLENGTH)
    if (bounded_id_at(rest, RSTART, RLENGTH)) {
      found = 1
      if (!decision_exists(id)) {
        return 0
      }
    } else {
      return 0
    }
    rest = substr(rest, RSTART + RLENGTH)
  }
  return 1
}

function escalation_id(text) {
  if (match(text, /ESC-[0-9]+-[0-9]+/)) {
    return substr(text, RSTART, RLENGTH)
  }
  return ""
}

function decision_exists(id, line) {
  if (id == "" || decisions_file == "") {
    return 0
  }
  while ((getline line < decisions_file) > 0) {
    if (index(line, "### " id ":") == 1) {
      close(decisions_file)
      return 1
    }
  }
  close(decisions_file)
  return 0
}

function decision_repair_owners_match(id, owners, line, in_section, found) {
  if (id == "" || decisions_file == "") {
    return 0
  }
  while ((getline line < decisions_file) > 0) {
    if (in_section && (index(line, "## ") == 1 || index(line, "### ") == 1)) {
      break
    }
    if (index(line, "### " id ":") == 1) {
      in_section = 1
      continue
    }
    if (in_section && trim_surrounding_space(line) == "- Repair owners: " owners) {
      found = 1
    }
  }
  close(decisions_file)
  return found
}

function first_owner(value, parts) {
  split(value, parts, ",")
  return trim_surrounding_space(parts[1])
}

function is_exact_triage_heading(line, id, prefix, rest) {
  prefix = "## " id
  if (index(line, prefix) != 1) {
    return 0
  }
  rest = substr(line, length(prefix) + 1)
  return rest == "" || index(rest, " - ") == 1
}

function triage_record_count(id, line, count) {
  if (id == "" || triage_file == "") {
    return 0
  }
  while ((getline line < triage_file) > 0) {
    if (is_exact_triage_heading(line, id)) {
      count++
    }
  }
  close(triage_file)
  return count
}

function triage_exists(id) {
  return triage_record_count(id) == 1
}

function amendment_heading(line, id, prefix) {
  prefix = "## Control Directive Amendment For " id
  return line == prefix
}

function amendment_count(id, line, count) {
  if (id == "" || triage_file == "") {
    return 0
  }
  while ((getline line < triage_file) > 0) {
    if (amendment_heading(line, id)) {
      count++
    }
  }
  close(triage_file)
  return count
}

function normalized_triage_verdict(verdict) {
  if (verdict == "REPAIR / control-conformance" || verdict == "REPAIR" || verdict == "RESUME") {
    return verdict
  }
  if (verdict ~ /^BLOCK/) {
    return "BLOCK"
  }
  return ""
}

function triage_verdict_matches(id, verdict, line, in_section,
                                marker, value, esc, found_verdict, verdict_count,
                                found_esc, expected) {
  if (id == "" || triage_file == "") {
    return 0
  }
  esc = escalation_id(verdict)
  found_esc = (esc == "")
  while ((getline line < triage_file) > 0) {
    if (index(line, "## ") == 1) {
      if (in_section) {
        break
      }
      in_section = is_exact_triage_heading(line, id)
      continue
    }
    if (!in_section) {
      continue
    }
    if (!found_esc && references_id(line, esc)) {
      found_esc = 1
    }
    marker = index(line, "- **Final verdict:** ")
    if (marker == 1) {
      value = substr(line, length("- **Final verdict:** ") + 1)
    } else {
      marker = index(line, "- **Verdict:** ")
      if (marker == 1) {
        value = substr(line, length("- **Verdict:** ") + 1)
      } else {
        marker = 0
      }
    }
    if (marker) {
      verdict_count++
      expected = normalized_triage_verdict(verdict)
      found_verdict = (value == expected)
    }
  }
  close(triage_file)
  return verdict_count == 1 && found_verdict && found_esc
}

function section_field_matches(id, section_kind, label, expected, line,
                              in_section, marker, value, count) {
  if (id == "" || triage_file == "") {
    return 0
  }
  marker = "- **" label ":** "
  while ((getline line < triage_file) > 0) {
    if (index(line, "## ") == 1) {
      if (in_section) {
        break
      }
      if (section_kind == "triage") {
        in_section = is_exact_triage_heading(line, id)
      } else {
        in_section = amendment_heading(line, id)
      }
      continue
    }
    if (!in_section || index(line, marker) != 1) {
      continue
    }
    count++
    value = substr(line, length(marker) + 1)
    if (value != expected) {
      close(triage_file)
      return 0
    }
  }
  close(triage_file)
  return count == 1
}

function triage_field_matches(id, label, expected) {
  return section_field_matches(id, "triage", label, expected)
}

function amendment_field_matches(id, label, expected) {
  return section_field_matches(id, "amendment", label, expected)
}

function authority_field_matches(id, verdict, label, expected) {
  if (verdict ~ /^BLOCK/) {
    return amendment_field_matches(id, label, expected)
  }
  return triage_field_matches(id, label, expected)
}

function section_decisions_exist(id, section_kind, line, in_section, value) {
  if (id == "" || triage_file == "") {
    return 0
  }
  while ((getline line < triage_file) > 0) {
    if (index(line, "## ") == 1) {
      if (in_section) {
        break
      }
      if (section_kind == "triage") {
        in_section = is_exact_triage_heading(line, id)
      } else {
        in_section = amendment_heading(line, id)
      }
      continue
    }
    if (in_section) {
      if (index(line, "- **Repair authority:** ") == 1) {
        value = substr(line, length("- **Repair authority:** ") + 1)
      } else if (index(line, "- **Resolution authority:** ") == 1) {
        value = substr(line, length("- **Resolution authority:** ") + 1)
      } else {
        continue
      }
      if (!all_decisions_exist(value)) {
        close(triage_file)
        return 0
      }
    }
  }
  close(triage_file)
  return in_section
}

function authority_decisions_exist(id, verdict) {
  if (verdict ~ /^BLOCK/) {
    return section_decisions_exist(id, "amendment")
  }
  return section_decisions_exist(id, "triage")
}

function valid_resolution_authority(value) {
  if (value ~ /^DEC-[0-9]+$/) {
    return decision_exists(value)
  }
  if (value ~ /^verified-commit:[0-9a-f]+$/) {
    return valid_commit(substr(value, length("verified-commit:") + 1))
  }
  return value ~ /^restored-capability:[A-Za-z0-9._:-]+$/
}

function event_decisions_exist(ledger_authority) {
  ledger_authority = event_repair_authority " " event_triage_authority " " \
    event_resolution_authority
  return all_decisions_exist(ledger_authority)
}

function exact_once(count, value) {
  return count == 1 && value != ""
}

function trim_surrounding_space(value) {
  sub(/^[[:space:]]+/, "", value)
  sub(/[[:space:]]+$/, "", value)
  return value
}

function valid_review_owner_list(value, id, parts, seen,
                                 count, i, task, valid) {
  count = split(value, parts, ",")
  valid = 1
  for (i = 1; i <= count; i++) {
    task = trim_surrounding_space(parts[i])
    if (task == "") {
      invalid("REVIEW REPAIR has empty owning task: " id)
      valid = 0
    } else if (task !~ /^task-[0-9]+[a-z]*$/) {
      invalid("REVIEW REPAIR has invalid owning task " task ": " id)
      valid = 0
    } else if (task in seen) {
      invalid("REVIEW REPAIR has duplicate owning task " task ": " id)
      valid = 0
    } else {
      seen[task] = 1
    }
  }
  return valid
}

function clear_review_owner_state(i, task) {
  for (i = 1; i <= review_owner_total; i++) {
    task = review_owner_tasks[i]
    delete review_owner_tasks[i]
    delete review_owner_position[task]
    delete review_owner_recovery_seen[task]
    delete review_owner_repair_seen[task]
    delete review_owner_complete[task]
  }
  review_owner_total = 0
  review_next_owner = 1
}

function set_review_owner_state(value, parts, count, i, task) {
  clear_review_owner_state()
  count = split(value, parts, ",")
  for (i = 1; i <= count; i++) {
    task = trim_surrounding_space(parts[i])
    review_owner_total++
    review_owner_tasks[review_owner_total] = task
    review_owner_position[task] = review_owner_total
  }
}

function valid_commit(id) {
  return length(id) == 40 && id !~ /[^0-9a-f]/
}

function trigger_matches_task(trigger, task, prefix, suffix) {
  if (task !~ /^task-[0-9]+[a-z]*$/) {
    return 0
  }
  if (index(trigger, task ":") == 1) {
    return 1
  }
  prefix = review_id ":qa-rework:" task
  if (!review_open || !(task in review_owner_position) ||
      index(trigger, prefix) != 1) {
    return 0
  }
  suffix = substr(trigger, length(prefix) + 1)
  return suffix == "" || suffix ~ /^[-:]/
}

function pending_triage_reconciliation() {
  if (!event_directive_v2 || stopped || resume_recovery_pending ||
      latest_triage_verdict != "REPAIR" || triage_session_liveness_pending) {
    return 0
  }
  if (event == "RUN-RESUME") {
    return event_continuing_task == latest_triage_task
  }
  if (event != "TRIAGE" || event_verdict != "REPAIR" ||
      event_originating_task != latest_triage_task ||
      event_attempt != latest_triage_attempt) {
    return 0
  }
  if (recovery_task == "") {
    return 1
  }
  if (event_originating_task != recovery_task || recovery_mode != "task-scope") {
    return 0
  }
  return recovery_stage == "needs-repair" ||
         (recovery_stage == "needs-continuation" &&
          event_trigger_signature == recovery_trigger &&
          event_attempt + 0 == recovery_attempt + 1)
}

function finish_event() {
  event_strict = strict || event_contract == "v2" || event_directive_schema != ""
  event_directive_strict = directive_strict || event_directive_schema == "v1" || event_directive_schema == "v2"
  event_directive_v2 = directive_v2 || event_directive_schema == "v2"
  if (event_strict && event == "") {
    invalid("unrecognized or malformed event heading")
  }
  if (event_directive_v2 && event != "" && !event_heading_canonical) {
    invalid(event " must use its canonical heading shape")
  }
  if (stop_entry_pending && event != "" && event != "STOP") {
    invalid("TRIAGE BLOCK must be followed by STOP before " event)
  }
  if (recovery_stage == "needs-triage" && event != "" && event != "TRIAGE") {
    invalid("recovery no-progress budget requires TRIAGE before " event)
  }
  if (recovery_stage == "needs-session-triage" && event != "" && event != "TRIAGE") {
    invalid("orchestration session-liveness requires TRIAGE before " event)
  }
  if (resume_recovery_pending && event != "" && event != "TASK-RECOVERY") {
    invalid("RUN-RESUME must be followed by its bound TASK-RECOVERY before " event)
  }
  if (triage_repair_pending && event != "" && event != "TASK-RECOVERY" &&
      !pending_triage_reconciliation()) {
    invalid("TRIAGE REPAIR must be followed by its bound TASK-RECOVERY before " event)
  }
  if (event_strict && event == "REVIEW" && event_outcome !~ /^(PASS|REPAIR|BLOCK)$/) {
    invalid("REVIEW has missing or invalid Outcome: " event_id)
  }
  if (event_strict && event == "REVIEW" && event_id !~ /^(review-[0-9]+[a-z]*|infra-review-[0-9]+[a-z]*)$/) {
    invalid("REVIEW has missing or invalid review ID")
  }
  if (event_strict && event == "REVIEW" && (event_owner == "" || event_rereview == "")) {
    invalid("REVIEW is missing Owning task or Re-review of: " event_id)
  }
  if (event_strict && event == "TRIAGE" && event_verdict !~ /^(RESUME|REPAIR|BLOCK)/) {
    invalid("TRIAGE has missing or invalid Verdict")
  }
  if (event_strict && (event == "TASK-RECOVERY" || event == "REPAIR" || event == "TASK-COMPLETE") && event_task !~ /^task-[0-9]+[a-z]*$/) {
    invalid(event " has missing or invalid task ID")
  }
  if (event_strict && (event == "RUN-START" || event == "RUN-RESUME") && event_contract != "v2") {
    invalid(event " is missing Sequence contract: v2")
  }
  if (event_strict && (event == "RUN-START" || event == "RUN-RESUME") && event_contract_count != 1) {
    invalid(event " must contain exactly one Sequence contract field")
  }
  if (event_directive_schema != "" && event_directive_schema !~ /^v[12]$/) {
    invalid(event " has invalid Recovery directive schema")
  }
  if (directive_strict && (event == "RUN-START" || event == "RUN-RESUME") &&
      event_directive_schema == "") {
    invalid(event " is missing Recovery directive schema")
  }
  if (event_directive_strict && (event == "RUN-START" || event == "RUN-RESUME") &&
      event_directive_schema_count != 1) {
    invalid(event " must contain exactly one Recovery directive schema field")
  }
  if (directive_v2 && (event == "RUN-START" || event == "RUN-RESUME") && event_directive_schema != "v2") {
    invalid(event " cannot downgrade Recovery directive schema from v2")
  }
  if (!directive_strict && event_directive_schema == "v1" && saw_v2 && !stopped) {
    invalid(event " cannot activate legacy Recovery directive schema after an unmarked v2 prefix without a reconciled STOP")
  }
  if (!directive_strict && event_directive_schema == "v2" && saw_v2) {
    invalid(event " cannot activate Recovery directive schema v2 after an unmarked v2 RUN")
  }
  if (event_directive_v2 && event != "" && !event_decisions_exist()) {
    invalid(event " cites an unrecorded or malformed phase decision")
  }
  if (event_directive_v2 && event == "RUN-START") {
    if (!exact_once(event_branch_count, event_branch) ||
        !exact_once(event_esc_baseline_count, event_esc_baseline) ||
        !exact_once(event_warmup_state_count, event_warmup_state) ||
        event_starting_task_count != 1 || event_starting_task !~ /^task-[0-9]+[a-z]*$/) {
      invalid("schema v2 RUN-START must contain each checkpoint field exactly once")
    }
  }
  if (event_directive_v2 && event == "RUN-RESUME") {
    if (!exact_once(event_prior_run_count, event_prior_run) ||
        event_governing_trg_count != 1 || event_governing_esc_count != 1 ||
        event_trigger_count != 1 || event_attempt_count != 1 ||
        !exact_once(event_checkpoint_count, event_checkpoint) ||
        !exact_once(event_resolution_evidence_count, event_resolution_evidence) ||
        event_resolution_authority_count != 1 ||
        event_continuing_task_count != 1 || event_continuing_task !~ /^task-[0-9]+[a-z]*$/ ||
        event_next_recovery_count != 1) {
      invalid("schema v2 RUN-RESUME must contain each reconciliation field exactly once")
    }
    if (!stopped &&
        (event_governing_trg != "N/A" || event_governing_esc != "N/A" ||
         event_trigger_signature != "N/A" || event_attempt != "N/A" ||
         event_resolution_authority != "N/A" || event_next_recovery != "N/A")) {
      invalid("normal schema v2 RUN-RESUME must use N/A for STOP-only fields")
    }
    if (!stopped && event_repair_owners_count) {
      invalid("Repair owners are valid only on a RUN-RESUME after STOP")
    }
  }
  if (event_directive_v2 && event == "TASK-RECOVERY") {
    if (!exact_once(event_source_count, event_source) ||
        !exact_once(event_category_count, event_category) ||
        !exact_once(event_repair_mode_count, event_repair_mode) ||
        !exact_once(event_repair_authority_count, event_repair_authority) ||
        !exact_once(event_triage_disposition_count, event_triage_disposition) ||
        !exact_once(event_trigger_count, event_trigger_signature) ||
        !exact_once(event_attempt_count, event_attempt) ||
        !exact_once(event_control_contract_count, event_control_contract) ||
        !exact_once(event_control_mismatch_count, event_control_mismatch) ||
        !exact_once(event_allowed_files_count, event_allowed_files) ||
        !exact_once(event_forbidden_changes_count, event_forbidden_changes) ||
        !exact_once(event_originating_report_count, event_originating_report) ||
        !exact_once(event_exact_verification_count, event_exact_verification) ||
        !exact_once(event_control_base_count, event_control_base) ||
        !exact_once(event_no_progress_count, event_no_progress) ||
        !exact_once(event_checkpoint_count, event_checkpoint) ||
        !exact_once(event_phase_worktree_count, event_phase_worktree) ||
        !exact_once(event_temporary_worktree_count, event_temporary_worktree) ||
        !exact_once(event_hidden_flags_count, event_hidden_flags) ||
        !exact_once(event_resolution_evidence_count, event_resolution_evidence) ||
        !exact_once(event_next_action_count, event_next_action) ||
        !exact_once(event_human_attention_count, event_human_attention)) {
      invalid("schema v2 TASK-RECOVERY must contain each typed field exactly once")
    }
  }
  if (event_directive_v2 && event == "REPAIR") {
    if (!exact_once(event_trigger_count, event_trigger_signature) ||
        !exact_once(event_attempt_count, event_attempt) ||
        !exact_once(event_starting_checkpoint_count, event_starting_checkpoint) ||
        !exact_once(event_ending_checkpoint_count, event_ending_checkpoint) ||
        !exact_once(event_control_base_count, event_control_base) ||
        !exact_once(event_control_head_count, event_control_head) ||
        !exact_once(event_control_scope_count, event_control_scope_check) ||
        !exact_once(event_phase_integration_count, event_phase_integration) ||
        !exact_once(event_outcome_count, event_outcome) ||
        !exact_once(event_resolution_evidence_count, event_resolution_evidence) ||
        !exact_once(event_next_checkpoint_count, event_next_checkpoint)) {
      invalid("schema v2 REPAIR must contain each checkpoint field exactly once")
    }
  }
  if (event_directive_v2 && event == "TRIAGE") {
    if (!exact_once(event_triage_trigger_count, event_triage_trigger) ||
        !exact_once(event_originating_task_count, event_originating_task) ||
        !exact_once(event_trigger_count, event_trigger_signature) ||
        !exact_once(event_attempt_count, event_attempt) ||
        !exact_once(event_repair_category_count, event_repair_category) ||
        !exact_once(event_triage_authority_count, event_triage_authority) ||
        !exact_once(event_checkpoint_count, event_checkpoint) ||
        !exact_once(event_verdict_count, event_verdict) ||
        !exact_once(event_resolution_evidence_count, event_resolution_evidence) ||
        !exact_once(event_triage_log_count, event_triage_log)) {
      invalid("schema v2 TRIAGE must contain each authority field exactly once")
    }
  }
  if (event_directive_v2 && event == "STOP") {
    if (!exact_once(event_cause_count, event_cause) ||
        !exact_once(event_triage_reference_count, event_triage_reference) ||
        !exact_once(event_protected_category_count, event_protected_category) ||
        !exact_once(event_next_human_action_count, event_next_human_action)) {
      invalid("schema v2 STOP must contain each authority field exactly once")
    }
  }
  if (event_directive_v2 && event == "TASK-COMPLETE") {
    expected_report = "project-docs/docs/phases/phase-" phase_number "/tasks/" event_task "-done.md"
    if (!exact_once(event_report_count, event_report) || event_report != expected_report ||
        !exact_once(event_report_lint_count, event_report_lint) ||
        event_report_lint !~ /PASS/ || event_report_lint !~ /OUTCOME COMPLETE/ ||
        !exact_once(event_commit_count, event_commit) ||
        !exact_once(event_qa_commands_count, event_qa_commands)) {
      invalid("schema v2 TASK COMPLETE must contain each completion field exactly once")
    }
  }
  if (event_directive_v2 && event == "REVIEW") {
    if (!exact_once(event_review_path_count, event_review_path) ||
        !exact_once(event_tasks_covered_count, event_tasks_covered) ||
        !exact_once(event_outcome_count, event_outcome) ||
        !exact_once(event_owner_count, event_owner) ||
        !exact_once(event_rereview_count, event_rereview) ||
        !exact_once(event_must_fix_count, event_must_fix)) {
      invalid("schema v2 REVIEW must contain each review field exactly once")
    }
  }
  if (event_directive_v2 && event == "RUN-END") {
    if (!exact_once(event_run_end_complete_count, event_run_end_complete) ||
        !exact_once(event_human_checklist_count, event_human_checklist) ||
        !exact_once(event_wrapup_ready_count, event_wrapup_ready)) {
      invalid("schema v2 RUN-END must contain each closure field exactly once")
    }
  }
  if (event == "RUN-END" && review_open && !review_task_complete) {
    invalid("RUN-END precedes completion of every QA repair owner: " review_owner)
  }
  if (event_strict && event_directive_strict && event == "TASK-RECOVERY") {
    if (event_source == "") {
      invalid("TASK-RECOVERY is missing Source")
    }
    if (event_category !~ /^(in-scope-implementation|gate-repair|report-repair|qa-rework|session-continuation|local-environment)$/) {
      invalid("TASK-RECOVERY has missing or invalid Category")
    }
    if (event_repair_mode !~ /^(task-scope|control-conformance)$/) {
      invalid("TASK-RECOVERY has missing or invalid Repair mode")
    }
    if (event_directive_v2 &&
        ((event_repair_mode == "control-conformance" &&
          event_triage_disposition != "REPAIR / control-conformance") ||
         (event_repair_mode == "task-scope" && event_triage_disposition != "N/A"))) {
      invalid("TASK-RECOVERY Triage disposition does not match Repair mode")
    }
    if (event_repair_authority == "") {
      invalid("TASK-RECOVERY is missing Repair authority")
    }
    if (event_allowed_files == "") {
      invalid("TASK-RECOVERY is missing Allowed files")
    }
    if (event_forbidden_changes == "") {
      invalid("TASK-RECOVERY is missing Forbidden changes")
    }
    if (event_trigger_signature == "") {
      invalid("TASK-RECOVERY is missing Trigger signature")
    }
    if (event_attempt !~ /^[1-9][0-9]*$/) {
      invalid("TASK-RECOVERY has missing or invalid Attempt count")
    }
    if (event_no_progress !~ /^[01]\/2$/) {
      invalid("TASK-RECOVERY has missing or invalid Consecutive no-progress attempts")
    }
    if (event_checkpoint == "") {
      invalid("TASK-RECOVERY is missing Checkpoint")
    }
    if (event_resolution_evidence == "") {
      invalid("TASK-RECOVERY is missing Resolution evidence")
    }
    if (event_next_action == "") {
      invalid("TASK-RECOVERY is missing Next action")
    }
    if (event_human_attention != "no") {
      invalid("TASK-RECOVERY must record Human attention required: no")
    }
    if (triage_repair_pending) {
      expected_triage_attempt = triage_session_liveness_pending ? latest_triage_attempt + 1 : latest_triage_attempt
      if (event_directive_v2) {
        if (event_task != latest_triage_task ||
            event_trigger_signature != latest_triage_trigger ||
            event_attempt != expected_triage_attempt ||
            event_category != latest_triage_category ||
            event_repair_authority != latest_triage_authority ||
            !references_id(event_source, latest_triage_id)) {
          invalid("TASK-RECOVERY does not preserve its governing schema v2 TRIAGE REPAIR authority")
        }
      } else if (event_trigger_signature != latest_triage_trigger ||
                 index(event_trigger_signature, event_task ":") != 1 ||
                 event_attempt != expected_triage_attempt ||
                 !references_id(event_source, latest_triage_id) ||
                 !references_id(event_repair_authority, latest_triage_id)) {
        invalid("TASK-RECOVERY does not preserve its governing TRIAGE REPAIR authority")
      }
      if (triage_session_liveness_pending && event_no_progress != recovery_no_progress) {
        invalid("session-liveness TRIAGE repair changed Consecutive no-progress attempts")
      }
    }
    if (resume_recovery_pending) {
      if (resume_owner_route) {
        expected_resume_attempt = event_trigger_signature == resume_trigger ? resume_attempt : "1"
        if (event_task != resume_task ||
            !trigger_matches_task(event_trigger_signature, event_task) ||
            event_attempt != expected_resume_attempt ||
            !references_id(event_source, resume_triage) ||
            !references_id(event_repair_authority, resume_triage)) {
          invalid("TASK-RECOVERY does not preserve its RUN-RESUME Repair owners checkpoint")
        }
      } else if (event_task != resume_task || event_trigger_signature != resume_trigger ||
          event_attempt != resume_attempt || !references_id(event_source, resume_triage) ||
          !references_id(event_repair_authority, resume_triage)) {
        invalid("TASK-RECOVERY does not preserve its RUN-RESUME checkpoint")
      }
      if (resume_decision != "" &&
          (!references_id(event_repair_authority, resume_decision) ||
           !references_id(event_source, resume_decision))) {
        invalid("TASK-RECOVERY does not retain its RUN-RESUME decision authority")
      }
      if (event_directive_v2 && resume_resolution_authority != "" &&
          (!references_exact_token(event_repair_authority, resume_resolution_authority) ||
           !references_exact_token(event_source, resume_resolution_authority))) {
        invalid("TASK-RECOVERY does not retain its typed RUN-RESUME resolution authority")
      }
      if (event_no_progress != "0/2") {
        invalid("decision-authorized RUN-RESUME must restart consecutive no-progress at 0/2")
      }
      if (event_directive_v2) {
        recovery_tuple = event_task " | " event_category " | " event_repair_mode " | " event_trigger_signature " | " event_attempt
        if (resume_next_recovery != recovery_tuple) {
          invalid("TASK-RECOVERY does not match the exact schema v2 RUN-RESUME recovery tuple")
        }
      }
    }
    if (recovery_task != "" && recovery_stage == "needs-repair" &&
        !(triage_repair_pending && event_task == recovery_task &&
          (event_directive_v2 || triage_session_liveness_pending))) {
      invalid("overlapping TASK-RECOVERY precedes required REPAIR")
    } else if (recovery_task != "" && recovery_stage == "needs-continuation") {
      if (event_task != recovery_task || event_trigger_signature != recovery_trigger) {
        invalid("continuing recovery must preserve task and Trigger signature")
      }
      if (event_attempt !~ /^[1-9][0-9]*$/ || event_attempt + 0 != recovery_attempt + 1) {
        invalid("continuing recovery Attempt count must increment by one")
      }
      if (recovery_last_outcome == "NO-PROGRESS") {
        expected_no_progress = (substr(recovery_no_progress, 1, 1) + 1) "/2"
      } else if (recovery_last_outcome == "INTERRUPTED") {
        expected_no_progress = recovery_no_progress
      } else {
        expected_no_progress = "0/2"
      }
      if (event_no_progress != expected_no_progress) {
        invalid("continuing recovery has invalid Consecutive no-progress attempts; expected " expected_no_progress)
      }
      if (event_repair_mode != recovery_mode) {
        invalid("continuing recovery cannot change Repair mode")
      }
      if (recovery_mode == "control-conformance" &&
          (event_repair_authority != recovery_authority ||
           event_control_contract != recovery_contract ||
           event_control_mismatch != recovery_mismatch ||
           event_allowed_files != recovery_allowed ||
           event_forbidden_changes != recovery_forbidden ||
           event_originating_report != recovery_report ||
           event_exact_verification != recovery_verification)) {
        invalid("continuing control recovery changed its authoritative directive")
      }
      if (recovery_mode == "control-conformance") {
        expected_retry_base = recovery_head != "" ? recovery_head : recovery_base
        if (event_control_base != expected_retry_base) {
          invalid("continuing control recovery did not use the prior integrated HEAD as its Runner base")
        }
      }
    } else if (recovery_task != "" && recovery_stage == "needs-triage") {
      invalid("recovery budget is exhausted; TASK-RECOVERY cannot precede TRIAGE")
    } else if (recovery_task != "" && recovery_stage == "needs-session-triage") {
      invalid("orchestration session-liveness requires TRIAGE before TASK-RECOVERY")
    } else if (recovery_task != "" && recovery_stage == "needs-report") {
      expected_report = "project-docs/docs/phases/phase-" phase_number "/tasks/" recovery_task "-done.md"
      if (recovery_mode != "control-conformance" || event_repair_mode != "task-scope" ||
          event_task != recovery_task || event_trigger_signature != recovery_trigger ||
          event_attempt !~ /^[1-9][0-9]*$/ || event_attempt + 0 != recovery_attempt + 1 ||
          event_no_progress != "0/2" || event_category != "report-repair" ||
          event_allowed_files != expected_report) {
        invalid("resolved control recovery permits only its linked exact report continuation")
      }
      if (event_directive_v2 &&
          (event_originating_report != recovery_report ||
           event_exact_verification != recovery_verification ||
           event_control_base != recovery_base ||
           !references_id(event_source, latest_triage_id))) {
        invalid("schema v2 report continuation changed its report, verification, Runner base, or governing TRIAGE")
      }
    }
  }
  if (event_strict && event_directive_strict && event == "TASK-RECOVERY" && event_repair_mode == "control-conformance") {
    control_retry = (recovery_task != "" && recovery_stage == "needs-continuation" && recovery_mode == "control-conformance")
    if (event_control_contract == "") {
      invalid("control-conformance TASK-RECOVERY is missing Control contract")
    }
    if (event_control_mismatch == "") {
      invalid("control-conformance TASK-RECOVERY is missing Control mismatch")
    }
    if (event_allowed_files == "") {
      invalid("control-conformance TASK-RECOVERY is missing Allowed files")
    }
    if (event_forbidden_changes == "") {
      invalid("control-conformance TASK-RECOVERY is missing Forbidden changes")
    }
    if (latest_triage_id == "") {
      invalid("control-conformance TASK-RECOVERY has no preceding TRIAGE authority")
    } else {
      if (!triage_exists(latest_triage_id) ||
          !triage_verdict_matches(latest_triage_id, latest_triage_verdict)) {
        invalid("control-conformance authority does not match the phase triage log: " latest_triage_id)
      }
      if (!references_id(event_repair_authority, latest_triage_id)) {
        invalid("control-conformance Repair authority does not reference latest TRIAGE " latest_triage_id)
      }
      if (!all_decisions_exist(event_repair_authority)) {
        invalid("control-conformance Repair authority cites an unrecorded or malformed phase decision")
      }
      if (event_trigger_signature != latest_triage_trigger) {
        invalid("control-conformance Trigger signature does not match latest TRIAGE " latest_triage_id)
      }
      if (index(event_trigger_signature, event_task ":") != 1) {
        invalid("control-conformance Trigger signature does not bind originating task " event_task)
      }
      if (latest_triage_verdict == "REPAIR / control-conformance") {
        # TL-Lite directly authorized this bounded repair.
      } else if (latest_triage_verdict ~ /^BLOCK/) {
        authority_decision = decision_id(event_repair_authority)
        resolved_authority_text = event_directive_v2 ? resolved_block_authority : resolved_block_evidence
        if (resolved_block_triage != latest_triage_id || authority_decision == "" ||
            !decision_exists(authority_decision) ||
            !all_decisions_exist(event_repair_authority) ||
            !all_decisions_exist(resolved_authority_text)) {
          invalid("control-conformance BLOCK authority lacks a bound RUN-RESUME decision for " latest_triage_id)
        }
        if (event_directive_v2) {
          if (amendment_count(latest_triage_id) != 1 ||
              !amendment_field_matches(latest_triage_id, "Governing TRG", latest_triage_id) ||
              !amendment_field_matches(latest_triage_id, "Governing ESC", latest_triage_esc) ||
              !amendment_field_matches(latest_triage_id, "Resolution authority", resolved_block_authority) ||
              !amendment_field_matches(latest_triage_id, "Verdict", "REPAIR / control-conformance") ||
              decision_id(resolved_block_authority) == "" ||
              !section_decisions_exist(latest_triage_id, "triage") ||
              !authority_decisions_exist(latest_triage_id, latest_triage_verdict)) {
            invalid("schema v2 resolved BLOCK lacks one exact decision-authorized control amendment")
          }
        }
      } else {
        invalid("control-conformance latest TRIAGE is not REPAIR / control-conformance or resolved BLOCK: " latest_triage_id)
      }
      if (!control_retry && event_attempt ~ /^[1-9][0-9]*$/ && latest_triage_attempt ~ /^[1-9][0-9]*$/ &&
          (event_attempt + 0 < latest_triage_attempt + 0 || event_attempt + 0 > latest_triage_attempt + 1)) {
        invalid("control-conformance Attempt count exceeds latest TRIAGE authority " latest_triage_id)
      }
    }
    if (event_directive_v2) {
      expected_report = "project-docs/docs/phases/phase-" phase_number "/tasks/" event_task "-done.md"
      if (event_trigger_count != 1 || event_attempt_count != 1 ||
          event_repair_authority_count != 1 || event_control_contract_count != 1 ||
          event_control_mismatch_count != 1 || event_allowed_files_count != 1 ||
          event_forbidden_changes_count != 1 || event_exact_verification_count != 1 ||
          event_originating_report_count != 1 || event_control_base_count != 1) {
        invalid("control-conformance schema v2 directive fields must each occur exactly once")
      }
      if (event_originating_report == "" || event_exact_verification == "" ||
          event_originating_report != expected_report || !valid_commit(event_control_base)) {
        invalid("control-conformance schema v2 directive is missing report, verification, or Runner base")
      }
      if (!authority_decisions_exist(latest_triage_id, latest_triage_verdict) ||
          !authority_field_matches(latest_triage_id, latest_triage_verdict, "Originating task", event_task) ||
          !authority_field_matches(latest_triage_id, latest_triage_verdict, "Control trigger signature", event_trigger_signature) ||
          !authority_field_matches(latest_triage_id, latest_triage_verdict, "Control contract", event_control_contract) ||
          !authority_field_matches(latest_triage_id, latest_triage_verdict, "Control mismatch", event_control_mismatch) ||
          !authority_field_matches(latest_triage_id, latest_triage_verdict, "Allowed files", event_allowed_files) ||
          !authority_field_matches(latest_triage_id, latest_triage_verdict, "Forbidden changes", event_forbidden_changes) ||
          !authority_field_matches(latest_triage_id, latest_triage_verdict, "Exact verification", event_exact_verification) ||
          !authority_field_matches(latest_triage_id, latest_triage_verdict, "Originating report", event_originating_report) ||
          (!control_retry && !authority_field_matches(latest_triage_id, latest_triage_verdict, "Control attempt count", event_attempt)) ||
          (!control_retry && !authority_field_matches(latest_triage_id, latest_triage_verdict, "Control repair base", event_control_base))) {
        invalid("control-conformance schema v2 directive does not match its authoritative TRIAGE record")
      }
      if (index("," event_allowed_files ",", ",project-docs/tools/workflow-control-scope-lint.sh,") ||
          index("," event_allowed_files ",", ",project-docs/tools/workflow-sequence-lint.sh,")) {
        invalid("control-conformance cannot authorize a protected workflow verifier")
      }
    }
    if (event_next_action !~ /workflow-control-scope-lint\.sh/) {
      invalid("control-conformance TASK-RECOVERY Next action is missing exact scope verification")
    }
  }
  if (event_strict && event_directive_strict && event == "REPAIR") {
    if (event_trigger_signature == "" || event_trigger_signature != recovery_trigger) {
      invalid("REPAIR Trigger signature does not match TASK-RECOVERY")
    }
    if (event_attempt !~ /^[1-9][0-9]*$/ || event_attempt != recovery_attempt) {
      invalid("REPAIR Attempt count does not match TASK-RECOVERY")
    }
    if (event_control_scope_check == "") {
      invalid("REPAIR is missing Control scope check")
    }
    if (event_outcome !~ /^(RESOLVED|MATERIAL-PROGRESS|NO-PROGRESS|INTERRUPTED)$/) {
      invalid("REPAIR has missing or invalid Outcome")
    }
    if (event_outcome == "INTERRUPTED" &&
        (!exact_once(event_interruption_cause_count, event_interruption_cause) ||
         !exact_once(event_elapsed_session_count, event_elapsed_session) ||
         !exact_once(event_command_in_flight_count, event_command_in_flight))) {
      invalid("INTERRUPTED REPAIR must record interruption cause, elapsed session duration, and command in flight")
    }
    if (event_outcome != "INTERRUPTED" &&
        (event_interruption_cause_count || event_elapsed_session_count || event_command_in_flight_count)) {
      invalid("non-INTERRUPTED REPAIR cannot carry interruption context")
    }
    if (event_resolution_evidence == "") {
      invalid("REPAIR is missing Resolution evidence")
    }
    if (event_next_checkpoint == "") {
      invalid("REPAIR is missing Next checkpoint")
    }
    if (recovery_mode == "control-conformance" &&
        (event_control_scope_check !~ /^PASS/ || event_control_scope_check !~ /exit_code: 0/)) {
      invalid("control-conformance REPAIR lacks a passing terminated scope check")
    }
    if (recovery_mode == "control-conformance" && event_directive_v2) {
      expected_scope_receipt = "PASS control repair scope: " recovery_base " -> " event_control_head "; allowed=" recovery_allowed " (terminated exit_code: 0)"
      if (event_control_base_count != 1 || event_control_head_count != 1 ||
          event_control_scope_count != 1 || event_phase_integration_count != 1) {
        invalid("control-conformance schema v2 REPAIR checkpoint fields must each occur exactly once")
      }
      if (event_control_base != recovery_base || !valid_commit(event_control_head) ||
          event_control_scope_check != expected_scope_receipt ||
          event_phase_integration !~ /^PASS/ ||
          event_phase_integration !~ event_control_head) {
        invalid("control-conformance REPAIR does not bind the Runner base, repair HEAD, exact allowlist receipt, and phase integration")
      }
    }
  }
  if (event_strict && event_directive_strict && event == "TRIAGE") {
    if (event_id !~ /^TRG-[0-9]+-[0-9]+$/) {
      invalid("TRIAGE has missing or invalid TRG ID")
    } else if (event_id in triage_seen) {
      invalid("TRIAGE reuses an existing TRG ID: " event_id)
    } else if (phase_number != "" && index(event_id, "TRG-" phase_number "-") != 1) {
      invalid("TRIAGE ID does not match ledger phase-" phase_number ": " event_id)
    }
    if (phase_number == "") {
      invalid("strict TRIAGE ledger is not under a phase-N directory")
    } else if (!triage_exists(event_id)) {
      invalid("TRIAGE ID is not recorded in the phase triage log: " event_id)
    } else if (!triage_verdict_matches(event_id, event_verdict)) {
      invalid("TRIAGE verdict does not match the phase triage log: " event_id)
    }
    if (event_directive_v2 &&
        (!section_decisions_exist(event_id, "triage") ||
         !triage_field_matches(event_id, "Originating task", event_originating_task) ||
         !triage_field_matches(event_id, "Trigger signature", event_trigger_signature) ||
         !triage_field_matches(event_id, "Attempt count", event_attempt) ||
         !triage_field_matches(event_id, "Repair category", event_repair_category) ||
         !triage_field_matches(event_id, "Repair authority", event_triage_authority))) {
      invalid("schema v2 TRIAGE identity does not match its unique external record")
    }
    if (event_verdict !~ /^(RESUME|REPAIR( \/ control-conformance)?|BLOCK \(escalated as ESC-[0-9]+-[0-9]+\))$/) {
      invalid("TRIAGE has missing or invalid strict Verdict")
    }
    triage_esc = escalation_id(event_verdict)
    if (triage_esc != "" && phase_number != "" && index(triage_esc, "ESC-" phase_number "-") != 1) {
      invalid("TRIAGE ESC does not match ledger phase-" phase_number ": " triage_esc)
    }
    if (event_trigger_signature == "") {
      invalid("TRIAGE is missing Trigger signature")
    }
    if (event_attempt !~ /^[1-9][0-9]*$/) {
      invalid("TRIAGE has missing or invalid Attempt count")
    }
    if (event_checkpoint == "") {
      invalid("TRIAGE is missing Checkpoint")
    }
    if (event_resolution_evidence == "") {
      invalid("TRIAGE is missing Resolution evidence")
    }
    if (event_triage_log == "") {
      invalid("TRIAGE is missing Triage log")
    }
    if (event_directive_v2) {
      if (!trigger_matches_task(event_trigger_signature, event_originating_task)) {
        invalid("schema v2 TRIAGE Originating task does not match Trigger signature")
      }
      if (event_verdict ~ /^REPAIR/) {
        if (event_repair_category == "N/A" || event_triage_authority == "N/A") {
          invalid("schema v2 TRIAGE REPAIR requires category and repair authority")
        }
      } else if (event_repair_category != "N/A" || event_triage_authority != "N/A") {
        invalid("schema v2 TRIAGE RESUME/BLOCK must use N/A repair fields")
      }
    }
    if (recovery_stage == "needs-triage" && event_verdict !~ /^BLOCK/) {
      invalid("exhausted no-progress recovery TRIAGE must BLOCK")
    }
    if (recovery_stage == "needs-triage" &&
        (event_trigger_signature != recovery_trigger ||
         !trigger_matches_task(event_trigger_signature, recovery_task) ||
         event_attempt != recovery_attempt)) {
      invalid("exhausted no-progress recovery TRIAGE changed task, trigger, or attempt")
    }
    if (recovery_stage == "needs-session-triage" &&
        (event_trigger_signature != recovery_trigger ||
         !trigger_matches_task(event_trigger_signature, recovery_task) ||
         event_attempt != recovery_attempt)) {
      invalid("orchestration session-liveness TRIAGE changed task, trigger, or attempt")
    }
    if (recovery_stage == "needs-session-triage" && event_verdict ~ /^RESUME/) {
      invalid("orchestration session-liveness TRIAGE must REPAIR or BLOCK")
    }
  }
  if (event_strict && event_directive_strict && event == "RUN-RESUME" && stopped) {
    resume_owner_route_ok = 0
    if (event_prior_run == "" || event_trigger_signature == "" ||
        event_attempt !~ /^[1-9][0-9]*$/ || event_checkpoint == "" ||
        event_resolution_evidence == "" || event_continuing_task !~ /^task-[0-9]+[a-z]*$/) {
      invalid("RUN-RESUME after STOP is missing reconciliation authority or continuing task")
    } else {
      if (!triage_exists(latest_triage_id) ||
          !triage_verdict_matches(latest_triage_id, latest_triage_verdict)) {
        invalid("RUN-RESUME governing TRIAGE does not match the phase triage log: " latest_triage_id)
      }
      if (!references_id(event_prior_run, latest_triage_id) ||
          latest_triage_esc == "" || !references_id(event_prior_run, latest_triage_esc)) {
        invalid("RUN-RESUME does not reference its governing TRIAGE and ESC")
      }
      if (event_trigger_signature != latest_triage_trigger) {
        invalid("RUN-RESUME Trigger signature does not match its governing TRIAGE")
      }
      if (event_attempt + 0 != latest_triage_attempt + 1) {
        invalid("RUN-RESUME Attempt count must increment its governing TRIAGE")
      }
      if (event_repair_owners_count) {
        owner_route_errors = errors
        if (event_repair_owners_count != 1 || !event_directive_v2) {
          invalid("RUN-RESUME Repair owners must occur exactly once on a schema v2 RUN-RESUME")
        }
        valid_review_owner_list(event_repair_owners, "RUN-RESUME Repair owners")
        if (block_rereview_id == "" || !review_open || !review_task_complete ||
            latest_triage_verdict !~ /^BLOCK/ || latest_triage_id != block_rereview_triage) {
          invalid("RUN-RESUME Repair owners require a STOP from a BLOCK re-review of a completed open QA repair")
        }
        if (event_resolution_authority !~ /^DEC-[0-9]+$/ ||
            !decision_repair_owners_match(event_resolution_authority, event_repair_owners)) {
          invalid("RUN-RESUME Repair owners do not match a Repair owners line in its Resolution authority decision")
        }
        if (event_continuing_task != first_owner(event_repair_owners)) {
          invalid("RUN-RESUME Continuing task must be the first Repair owner")
        }
        resume_owner_route_ok = (errors == owner_route_errors)
      } else if (!trigger_matches_task(event_trigger_signature, event_continuing_task)) {
        invalid("RUN-RESUME Continuing task does not match its Trigger signature")
      }
      if (event_directive_v2 &&
          (event_governing_trg != latest_triage_id ||
           event_governing_esc != latest_triage_esc ||
           event_next_recovery == "" || event_resolution_authority == "")) {
        invalid("schema v2 RUN-RESUME is missing its governing IDs, resolution authority, or exact next recovery")
      }
      if (event_directive_v2 &&
          (event_governing_trg_count != 1 || event_governing_esc_count != 1 ||
           event_resolution_authority_count != 1 || event_next_recovery_count != 1)) {
        invalid("schema v2 RUN-RESUME authority fields must each occur exactly once")
      }
      decision_text = event_directive_v2 ? event_resolution_authority : event_resolution_evidence
      resume_decision = decision_id(event_resolution_authority)
      if (!all_decisions_exist(decision_text)) {
        invalid("RUN-RESUME cites an unrecorded or malformed phase decision")
      }
      if (event_directive_v2) {
        if (!valid_resolution_authority(event_resolution_authority)) {
          invalid("schema v2 RUN-RESUME has invalid or unverifiable Resolution authority")
        }
      } else if (latest_triage_verdict ~ /^BLOCK/ && decision_id(decision_text) == "") {
        invalid("historical RUN-RESUME after TRIAGE BLOCK requires a recorded phase decision")
      }
    }
  }
  if (event_strict && event_directive_strict && event == "STOP") {
    if (latest_triage_id == "" || latest_triage_verdict !~ /^BLOCK/) {
      invalid("STOP has no preceding TRIAGE BLOCK")
    } else if (!references_id(event_triage_reference, latest_triage_id) ||
               latest_triage_esc == "" || !references_id(event_triage_reference, latest_triage_esc)) {
      invalid("STOP does not reference its TRIAGE BLOCK and concrete ESC")
    }
    if (event_protected_category == "" || event_next_human_action == "") {
      invalid("STOP is missing protected category or next Human action")
    }
  }

  if (event == "TASK-RECOVERY") {
    if (recovery_task == "") {
      recovery_head = ""
      recovery_consecutive_interruptions = 0
    }
    if (review_open && review_task_complete) {
      invalid("QA repair owners of " review_id " are complete; a linked re-review or a " \
              "RUN-RESUME with Repair owners is required before recovery task " event_task)
    } else if (review_open && !(event_task in review_owner_position)) {
      invalid("QA repair owner " review_owner " does not match recovery task " event_task)
    } else if (review_open && event_task != review_owner_tasks[review_next_owner]) {
      invalid("QA repair owner " review_owner " requires recovery task " \
              review_owner_tasks[review_next_owner] " before " event_task)
    }
    if (recovery_task != "" && recovery_task != event_task) {
      invalid("overlapping recovery changes task identity from " recovery_task " to " event_task)
    } else {
      recovery_task = event_task
    }
    recovery_stage = "needs-repair"
    recovery_trigger = event_trigger_signature
    recovery_attempt = event_attempt
    recovery_mode = event_repair_mode
    recovery_no_progress = event_no_progress
    recovery_authority = event_repair_authority
    recovery_contract = event_control_contract
    recovery_mismatch = event_control_mismatch
    recovery_allowed = event_allowed_files
    recovery_forbidden = event_forbidden_changes
    recovery_report = event_originating_report
    recovery_verification = event_exact_verification
    recovery_base = event_control_base
    if (triage_session_liveness_pending) {
      recovery_consecutive_interruptions = 0
    }
    triage_repair_pending = 0
    triage_session_liveness_pending = 0
    resume_recovery_pending = 0
    if (review_open && event_task == review_owner_tasks[review_next_owner]) {
      review_owner_recovery_seen[event_task] = 1
    }
  } else if (event == "REPAIR") {
    if (recovery_task == "") {
      invalid("REPAIR has no preceding TASK-RECOVERY")
    } else if (event_task != recovery_task) {
      invalid("REPAIR task " event_task " does not match recovery task " recovery_task)
    } else {
      if (event_outcome == "RESOLVED") {
        if (recovery_mode == "control-conformance") {
          recovery_stage = "needs-report"
        } else {
          recovery_stage = "needs-complete"
        }
      } else if (event_outcome == "NO-PROGRESS" && recovery_no_progress == "1/2") {
        recovery_stage = "needs-triage"
        triage_needed = 1
        recovery_consecutive_interruptions = 0
      } else if (event_outcome == "INTERRUPTED") {
        recovery_consecutive_interruptions++
        if (recovery_consecutive_interruptions >= 2) {
          recovery_stage = "needs-session-triage"
          triage_needed = 1
        } else {
          recovery_stage = "needs-continuation"
        }
      } else {
        recovery_stage = "needs-continuation"
        recovery_consecutive_interruptions = 0
      }
      recovery_last_outcome = event_outcome
      if (recovery_mode == "control-conformance" && event_control_head != "") {
        recovery_head = event_control_head
      }
      if (review_open && event_task == review_owner_tasks[review_next_owner]) {
        review_owner_repair_seen[event_task] = 1
      }
    }
  } else if (event == "TASK-COMPLETE") {
    if (recovery_task != "") {
      if (event_task != recovery_task) {
        invalid("completed task " event_task " does not match recovery task " recovery_task)
      } else if (recovery_stage != "needs-complete") {
        invalid("task completed before a resolved REPAIR event: " event_task)
      } else {
        recovery_task = ""
        recovery_stage = ""
        recovery_head = ""
      }
    }
    if (review_open && (event_task in review_owner_position)) {
      if (event_task != review_owner_tasks[review_next_owner]) {
        invalid("QA repair owner " review_owner " requires completion of task " \
                review_owner_tasks[review_next_owner] " before " event_task)
      } else if (!review_owner_recovery_seen[event_task] ||
                 !review_owner_repair_seen[event_task]) {
        invalid("QA repair owner completed without TASK-RECOVERY and REPAIR: " event_task)
      } else {
        review_owner_complete[event_task] = 1
        review_next_owner++
        review_task_complete = review_next_owner > review_owner_total
      }
    }
  } else if (event == "REVIEW") {
    block_rereview_id = ""
    block_rereview_triage = ""
    if (event_outcome == "BLOCK" && review_open && review_task_complete &&
        event_rereview == review_id && event_id != review_id) {
      block_rereview_id = event_id
    }
    if (event_outcome == "REPAIR") {
      owner_list_valid = 0
      if (event_owner == "" || event_owner == "none") {
        invalid("REVIEW REPAIR has no owning task: " event_id)
      } else {
        owner_list_valid = valid_review_owner_list(event_owner, event_id)
      }
      if (review_open) {
        if (!review_task_complete) {
          invalid("QA repair was replaced before owning task completed: " review_owner)
        }
        if (event_id == review_id) {
          invalid("QA repair was self-replaced by the same review identity: " event_id)
        }
        if (event_rereview != review_id) {
          invalid("QA REPAIR does not reference required re-review " review_id ": " event_id)
        }
      }
      review_open = 1
      review_id = event_id
      review_owner = event_owner
      review_task_complete = 0
      clear_review_owner_state()
      if (owner_list_valid) {
        set_review_owner_state(event_owner)
      }
    } else if (event_outcome == "PASS" && review_open) {
      if (event_owner != "none") {
        invalid("REVIEW PASS must use Owning task: none: " event_id)
      }
      if (!review_task_complete) {
        invalid("QA re-review passed before owning task completed: " review_owner)
      }
      if (event_id == review_id) {
        invalid("QA repair was self-closed by the same review identity: " event_id)
      }
      if (event_rereview != review_id) {
        invalid("QA PASS does not reference required re-review " review_id ": " event_id)
      }
      review_open = 0
      review_id = ""
      review_owner = ""
      review_task_complete = 0
      clear_review_owner_state()
    } else if (event_outcome == "BLOCK") {
      if (event_owner != "none") {
        invalid("REVIEW BLOCK must use Owning task: none: " event_id)
      }
      triage_needed = 1
    } else if (event_outcome == "PASS" && event_owner != "none") {
      invalid("REVIEW PASS must use Owning task: none: " event_id)
    }
  } else if (event == "TRIAGE") {
    triage_was_exhaustion = (recovery_stage == "needs-triage")
    triage_was_session_liveness = (recovery_stage == "needs-session-triage")
    triage_session_liveness_pending = 0
    latest_triage_id = event_id
    latest_triage_trigger = event_trigger_signature
    latest_triage_verdict = event_verdict
    latest_triage_attempt = event_attempt
    latest_triage_task = event_originating_task
    latest_triage_category = event_repair_category
    latest_triage_authority = event_triage_authority
    latest_triage_esc = escalation_id(event_verdict)
    triage_seen[event_id] = 1
    if (block_rereview_id != "" && block_rereview_triage == "") {
      block_rereview_triage = event_id
    }
    if (event_verdict ~ /^RESUME/) {
      triage_needed = 0
    } else if (event_verdict ~ /^REPAIR/) {
      triage_needed = 0
      triage_repair_pending = 1
      if (triage_was_session_liveness) {
        recovery_stage = "needs-repair"
        triage_session_liveness_pending = 1
      }
    } else if (event_verdict ~ /^BLOCK/) {
      if (triage_was_exhaustion || triage_was_session_liveness) {
        recovery_task = ""
        recovery_stage = ""
        recovery_consecutive_interruptions = 0
        triage_session_liveness_pending = 0
      }
      stop_entry_pending = 1
      triage_needed = 0
    }
  } else if (event == "STOP") {
    stopped = 1
    stop_entry_pending = 0
  } else if (event == "RUN-RESUME") {
    if (stopped && latest_triage_verdict ~ /^BLOCK/) {
      resolved_block_triage = latest_triage_id
      resolved_block_evidence = event_resolution_evidence
      resolved_block_authority = event_resolution_authority
      if (event_directive_strict) {
        resume_recovery_pending = 1
        recovery_task = ""
        recovery_stage = ""
        recovery_trigger = ""
        recovery_attempt = ""
        recovery_last_outcome = ""
        recovery_no_progress = ""
        recovery_head = ""
        recovery_consecutive_interruptions = 0
        resume_triage = latest_triage_id
        resume_task = event_continuing_task
        resume_trigger = event_trigger_signature
        resume_attempt = event_attempt
        resume_decision = decision_id(event_resolution_evidence " " event_resolution_authority)
        resume_resolution_authority = event_resolution_authority
        resume_next_recovery = event_next_recovery
        resume_owner_route = event_repair_owners_count == 1 && resume_owner_route_ok
        resume_owner_route_ok = 0
        if (resume_owner_route) {
          review_id = block_rereview_id
          review_owner = event_repair_owners
          review_task_complete = 0
          set_review_owner_state(event_repair_owners)
        }
        block_rereview_id = ""
        block_rereview_triage = ""
      }
    }
    stopped = 0
    if (event_contract == "v2") {
      strict = 1
      saw_v2 = 1
    }
    if (event_directive_schema == "v1" || event_directive_schema == "v2") {
      directive_strict = 1
      saw_directive_schema = 1
    }
    if (event_directive_schema == "v2") {
      directive_v2 = 1
    }
  } else if (event == "RUN-START") {
    if (event_contract == "v2") {
      strict = 1
      saw_v2 = 1
    }
    if (event_directive_schema == "v1" || event_directive_schema == "v2") {
      directive_strict = 1
      saw_directive_schema = 1
    }
    if (event_directive_schema == "v2") {
      directive_v2 = 1
    }
  }


  event = ""
  event_task = ""
  event_id = ""
  event_outcome = ""
  event_owner = ""
  event_rereview = ""
  event_verdict = ""
  event_contract = ""
  event_directive_schema = ""
  event_source = ""
  event_category = ""
  event_repair_mode = ""
  event_repair_authority = ""
  event_control_contract = ""
  event_control_mismatch = ""
  event_allowed_files = ""
  event_forbidden_changes = ""
  event_trigger_signature = ""
  event_attempt = ""
  event_no_progress = ""
  event_checkpoint = ""
  event_resolution_evidence = ""
  event_next_action = ""
  event_human_attention = ""
  event_control_scope_check = ""
  event_next_checkpoint = ""
  event_triage_reference = ""
  event_protected_category = ""
  event_next_human_action = ""
  event_triage_log = ""
  event_prior_run = ""
  event_continuing_task = ""
  event_originating_report = ""
  event_exact_verification = ""
  event_control_base = ""
  event_control_head = ""
  event_phase_integration = ""
  event_governing_trg = ""
  event_governing_esc = ""
  event_resolution_authority = ""
  event_next_recovery = ""
  event_repair_owners = ""
  event_branch = ""
  event_esc_baseline = ""
  event_warmup_state = ""
  event_starting_task = ""
  event_originating_task = ""
  event_repair_category = ""
  event_triage_authority = ""
  event_triage_trigger = ""
  event_triage_disposition = ""
  event_phase_worktree = ""
  event_temporary_worktree = ""
  event_hidden_flags = ""
  event_interruption_cause = ""
  event_elapsed_session = ""
  event_command_in_flight = ""
  event_cause = ""
  event_starting_checkpoint = ""
  event_ending_checkpoint = ""
  event_report = ""
  event_report_lint = ""
  event_commit = ""
  event_qa_commands = ""
  event_review_path = ""
  event_tasks_covered = ""
  event_must_fix = ""
  event_run_end_complete = ""
  event_human_checklist = ""
  event_wrapup_ready = ""
  event_heading_canonical = 0
  event_contract_count = 0
  event_directive_schema_count = 0
  event_trigger_count = 0
  event_attempt_count = 0
  event_repair_authority_count = 0
  event_control_contract_count = 0
  event_control_mismatch_count = 0
  event_allowed_files_count = 0
  event_forbidden_changes_count = 0
  event_exact_verification_count = 0
  event_originating_report_count = 0
  event_control_base_count = 0
  event_control_head_count = 0
  event_control_scope_count = 0
  event_phase_integration_count = 0
  event_governing_trg_count = 0
  event_governing_esc_count = 0
  event_resolution_authority_count = 0
  event_next_recovery_count = 0
  event_repair_owners_count = 0
  event_source_count = 0
  event_category_count = 0
  event_repair_mode_count = 0
  event_triage_disposition_count = 0
  event_no_progress_count = 0
  event_checkpoint_count = 0
  event_resolution_evidence_count = 0
  event_next_action_count = 0
  event_human_attention_count = 0
  event_phase_worktree_count = 0
  event_temporary_worktree_count = 0
  event_hidden_flags_count = 0
  event_interruption_cause_count = 0
  event_elapsed_session_count = 0
  event_command_in_flight_count = 0
  event_verdict_count = 0
  event_triage_log_count = 0
  event_originating_task_count = 0
  event_repair_category_count = 0
  event_triage_authority_count = 0
  event_triage_trigger_count = 0
  event_branch_count = 0
  event_esc_baseline_count = 0
  event_warmup_state_count = 0
  event_starting_task_count = 0
  event_prior_run_count = 0
  event_continuing_task_count = 0
  event_cause_count = 0
  event_triage_reference_count = 0
  event_protected_category_count = 0
  event_next_human_action_count = 0
  event_starting_checkpoint_count = 0
  event_ending_checkpoint_count = 0
  event_outcome_count = 0
  event_next_checkpoint_count = 0
  event_report_count = 0
  event_report_lint_count = 0
  event_commit_count = 0
  event_qa_commands_count = 0
  event_review_path_count = 0
  event_tasks_covered_count = 0
  event_owner_count = 0
  event_rereview_count = 0
  event_must_fix_count = 0
  event_run_end_complete_count = 0
  event_human_checklist_count = 0
  event_wrapup_ready_count = 0
}

/^### / {
  finish_event()
  heading_prefix_canonical = $1 == "###" &&
    $2 ~ /^[0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]$/ &&
    $3 ~ /^[0-9][0-9]:[0-9][0-9]$/
  if ($4 == "TASK-RECOVERY") {
    event = "TASK-RECOVERY"
    event_task = $5
    event_heading_canonical = heading_prefix_canonical && NF == 5 &&
      $0 == "### " $2 " " $3 " " $4 " " $5
  } else if ($4 == "REPAIR") {
    event = "REPAIR"
    event_task = $5
    event_heading_canonical = heading_prefix_canonical && NF == 5 &&
      $0 == "### " $2 " " $3 " " $4 " " $5
  } else if ($4 == "TASK" && $6 == "COMPLETE") {
    event = "TASK-COMPLETE"
    event_task = $5
    event_heading_canonical = heading_prefix_canonical && NF == 6 &&
      $0 == "### " $2 " " $3 " TASK " $5 " COMPLETE"
  } else if ($4 == "REVIEW") {
    event = "REVIEW"
    event_id = $5
    event_heading_canonical = heading_prefix_canonical && NF == 5 &&
      $0 == "### " $2 " " $3 " REVIEW " $5
  } else if ($4 == "TRIAGE") {
    event = "TRIAGE"
    event_id = $5
    event_heading_canonical = heading_prefix_canonical && NF == 5 &&
      $0 == "### " $2 " " $3 " TRIAGE " $5
  } else if ($4 == "STOP") {
    event = "STOP"
    event_heading_canonical = heading_prefix_canonical && NF == 4 &&
      $0 == "### " $2 " " $3 " STOP"
  } else if ($4 == "RUN-RESUME") {
    event = "RUN-RESUME"
    event_heading_canonical = heading_prefix_canonical && NF == 4 &&
      $0 == "### " $2 " " $3 " RUN-RESUME"
  } else if ($4 == "RUN-START") {
    event = "RUN-START"
    event_heading_canonical = heading_prefix_canonical && NF == 4 &&
      $0 == "### " $2 " " $3 " RUN-START"
  } else if ($4 == "RUN-END") {
    event = "RUN-END"
    event_heading_canonical = heading_prefix_canonical && NF == 4 &&
      $0 == "### " $2 " " $3 " RUN-END"
  } else if ($4 == "CORRECTION") {
    event = "CORRECTION"
    event_heading_canonical = heading_prefix_canonical && NF == 4 &&
      $0 == "### " $2 " " $3 " CORRECTION"
  } else {
    event = ""
  }
  next
}

event == "REVIEW" && /^- Outcome: / {
  event_outcome_count++
  event_outcome = substr($0, length("- Outcome: ") + 1)
}
event == "REVIEW" && /^- Owning task: / {
  event_owner_count++
  event_owner = substr($0, length("- Owning task: ") + 1)
}
event == "REVIEW" && /^- Re-review of: / {
  event_rereview_count++
  event_rereview = substr($0, length("- Re-review of: ") + 1)
}
event == "REVIEW" && /^- Path: / {
  event_review_path_count++
  event_review_path = substr($0, length("- Path: ") + 1)
}
event == "REVIEW" && /^- Tasks covered: / {
  event_tasks_covered_count++
  event_tasks_covered = substr($0, length("- Tasks covered: ") + 1)
}
event == "REVIEW" && /^- Must-fix: / {
  event_must_fix_count++
  event_must_fix = substr($0, length("- Must-fix: ") + 1)
}
event == "TRIAGE" && /^- Verdict: / {
  event_verdict_count++
  event_verdict = substr($0, length("- Verdict: ") + 1)
}
event == "TRIAGE" && /^- Trigger: / {
  event_triage_trigger_count++
  event_triage_trigger = substr($0, length("- Trigger: ") + 1)
}
event == "TRIAGE" && /^- Triage log: / {
  event_triage_log_count++
  event_triage_log = substr($0, length("- Triage log: ") + 1)
}
(event == "TASK-RECOVERY" || event == "REPAIR" || event == "TRIAGE" || event == "RUN-RESUME") && /^- Trigger signature: / {
  event_trigger_count++
  event_trigger_signature = substr($0, length("- Trigger signature: ") + 1)
}
(event == "TASK-RECOVERY" || event == "REPAIR" || event == "TRIAGE" || event == "RUN-RESUME") && /^- Attempt count: / {
  event_attempt_count++
  event_attempt = substr($0, length("- Attempt count: ") + 1)
}
(event == "TASK-RECOVERY" || event == "TRIAGE" || event == "RUN-RESUME") && /^- Checkpoint: / {
  event_checkpoint_count++
  event_checkpoint = substr($0, length("- Checkpoint: ") + 1)
}
(event == "TASK-RECOVERY" || event == "REPAIR" || event == "TRIAGE" || event == "RUN-RESUME") && /^- Resolution evidence: / {
  event_resolution_evidence_count++
  event_resolution_evidence = substr($0, length("- Resolution evidence: ") + 1)
}
(event == "RUN-START" || event == "RUN-RESUME") && /^- Sequence contract: / {
  event_contract_count++
  event_contract = substr($0, length("- Sequence contract: ") + 1)
}
(event == "RUN-START" || event == "RUN-RESUME") && /^- Recovery directive schema: / {
  event_directive_schema_count++
  event_directive_schema = substr($0, length("- Recovery directive schema: ") + 1)
}
event == "RUN-RESUME" && /^- Prior run\/session or STOP: / {
  event_prior_run_count++
  event_prior_run = substr($0, length("- Prior run/session or STOP: ") + 1)
}
event == "RUN-RESUME" && /^- Continuing task: / {
  event_continuing_task_count++
  event_continuing_task = substr($0, length("- Continuing task: ") + 1)
}
event == "RUN-RESUME" && /^- Governing TRG: / {
  event_governing_trg_count++
  event_governing_trg = substr($0, length("- Governing TRG: ") + 1)
}
event == "RUN-RESUME" && /^- Governing ESC: / {
  event_governing_esc_count++
  event_governing_esc = substr($0, length("- Governing ESC: ") + 1)
}
event == "RUN-RESUME" && /^- Resolution authority: / {
  event_resolution_authority_count++
  event_resolution_authority = substr($0, length("- Resolution authority: ") + 1)
}
event == "RUN-RESUME" && /^- Next recovery: / {
  event_next_recovery_count++
  event_next_recovery = substr($0, length("- Next recovery: ") + 1)
}
event == "RUN-RESUME" && /^- Repair owners: / {
  event_repair_owners_count++
  event_repair_owners = substr($0, length("- Repair owners: ") + 1)
}
event == "TASK-RECOVERY" && /^- Repair mode: / {
  event_repair_mode_count++
  event_repair_mode = substr($0, length("- Repair mode: ") + 1)
}
event == "TASK-RECOVERY" && /^- Source: / {
  event_source_count++
  event_source = substr($0, length("- Source: ") + 1)
}
event == "TASK-RECOVERY" && /^- Category: / {
  event_category_count++
  event_category = substr($0, length("- Category: ") + 1)
}
event == "TASK-RECOVERY" && /^- Triage disposition: / {
  event_triage_disposition_count++
  event_triage_disposition = substr($0, length("- Triage disposition: ") + 1)
}
event == "TASK-RECOVERY" && /^- Repair authority: / {
  event_repair_authority_count++
  event_repair_authority = substr($0, length("- Repair authority: ") + 1)
}
event == "TASK-RECOVERY" && /^- Control contract: / {
  event_control_contract_count++
  event_control_contract = substr($0, length("- Control contract: ") + 1)
}
event == "TASK-RECOVERY" && /^- Control mismatch: / {
  event_control_mismatch_count++
  event_control_mismatch = substr($0, length("- Control mismatch: ") + 1)
}
event == "TASK-RECOVERY" && /^- Allowed files: / {
  event_allowed_files_count++
  event_allowed_files = substr($0, length("- Allowed files: ") + 1)
}
event == "TASK-RECOVERY" && /^- Forbidden changes: / {
  event_forbidden_changes_count++
  event_forbidden_changes = substr($0, length("- Forbidden changes: ") + 1)
}
event == "TASK-RECOVERY" && /^- Consecutive no-progress attempts: / {
  event_no_progress_count++
  event_no_progress = substr($0, length("- Consecutive no-progress attempts: ") + 1)
}
event == "TASK-RECOVERY" && /^- Next action: / {
  event_next_action_count++
  event_next_action = substr($0, length("- Next action: ") + 1)
}
event == "TASK-RECOVERY" && /^- Human attention required: / {
  event_human_attention_count++
  event_human_attention = substr($0, length("- Human attention required: ") + 1)
}
event == "TASK-RECOVERY" && /^- Phase worktree: / {
  event_phase_worktree_count++
  event_phase_worktree = substr($0, length("- Phase worktree: ") + 1)
}
event == "TASK-RECOVERY" && /^- Temporary repair worktree: / {
  event_temporary_worktree_count++
  event_temporary_worktree = substr($0, length("- Temporary repair worktree: ") + 1)
}
event == "TASK-RECOVERY" && /^- Hidden-index flags: / {
  event_hidden_flags_count++
  event_hidden_flags = substr($0, length("- Hidden-index flags: ") + 1)
}
(event == "TASK-RECOVERY") && /^- Originating report: / {
  event_originating_report_count++
  event_originating_report = substr($0, length("- Originating report: ") + 1)
}
(event == "TASK-RECOVERY") && /^- Exact verification: / {
  event_exact_verification_count++
  event_exact_verification = substr($0, length("- Exact verification: ") + 1)
}
(event == "TASK-RECOVERY" || event == "REPAIR") && /^- Control repair base: / {
  event_control_base_count++
  event_control_base = substr($0, length("- Control repair base: ") + 1)
}
event == "REPAIR" && /^- Control repair head: / {
  event_control_head_count++
  event_control_head = substr($0, length("- Control repair head: ") + 1)
}
event == "REPAIR" && /^- Phase-branch integration: / {
  event_phase_integration_count++
  event_phase_integration = substr($0, length("- Phase-branch integration: ") + 1)
}
event == "REPAIR" && /^- Control scope check: / {
  event_control_scope_count++
  event_control_scope_check = substr($0, length("- Control scope check: ") + 1)
}
event == "REPAIR" && /^- Outcome: / {
  event_outcome_count++
  event_outcome = substr($0, length("- Outcome: ") + 1)
}
event == "REPAIR" && /^- Interruption cause: / {
  event_interruption_cause_count++
  event_interruption_cause = substr($0, length("- Interruption cause: ") + 1)
}
event == "REPAIR" && /^- Elapsed session duration: / {
  event_elapsed_session_count++
  event_elapsed_session = substr($0, length("- Elapsed session duration: ") + 1)
}
event == "REPAIR" && /^- Command in flight: / {
  event_command_in_flight_count++
  event_command_in_flight = substr($0, length("- Command in flight: ") + 1)
}
event == "REPAIR" && /^- Next checkpoint: / {
  event_next_checkpoint_count++
  event_next_checkpoint = substr($0, length("- Next checkpoint: ") + 1)
}
event == "REPAIR" && /^- Starting checkpoint: / {
  event_starting_checkpoint_count++
  event_starting_checkpoint = substr($0, length("- Starting checkpoint: ") + 1)
}
event == "REPAIR" && /^- Ending checkpoint: / {
  event_ending_checkpoint_count++
  event_ending_checkpoint = substr($0, length("- Ending checkpoint: ") + 1)
}
event == "STOP" && /^- Triage: / {
  event_triage_reference_count++
  event_triage_reference = substr($0, length("- Triage: ") + 1)
}
event == "STOP" && /^- Protected category: / {
  event_protected_category_count++
  event_protected_category = substr($0, length("- Protected category: ") + 1)
}
event == "STOP" && (/^- Next Human or external-authority action: / || /^- Next human action: /) {
  event_next_human_action_count++
  event_next_human_action = substr($0, index($0, ": ") + 2)
}
event == "STOP" && /^- Cause: / {
  event_cause_count++
  event_cause = substr($0, length("- Cause: ") + 1)
}
event == "RUN-START" && /^- Branch: / {
  event_branch_count++
  event_branch = substr($0, length("- Branch: ") + 1)
}
event == "RUN-START" && /^- ESC baseline \(Open at run start\): / {
  event_esc_baseline_count++
  event_esc_baseline = substr($0, length("- ESC baseline (Open at run start): ") + 1)
}
event == "RUN-START" && /^- Warm-up\/autonomous-resume state: / {
  event_warmup_state_count++
  event_warmup_state = substr($0, length("- Warm-up/autonomous-resume state: ") + 1)
}
event == "RUN-START" && /^- Starting task: / {
  event_starting_task_count++
  event_starting_task = substr($0, length("- Starting task: ") + 1)
}
event == "TRIAGE" && /^- Originating task: / {
  event_originating_task_count++
  event_originating_task = substr($0, length("- Originating task: ") + 1)
}
event == "TRIAGE" && /^- Repair category: / {
  event_repair_category_count++
  event_repair_category = substr($0, length("- Repair category: ") + 1)
}
event == "TRIAGE" && /^- Repair authority: / {
  event_triage_authority_count++
  event_triage_authority = substr($0, length("- Repair authority: ") + 1)
}
event == "TASK-COMPLETE" && /^- Report: / {
  event_report_count++
  event_report = substr($0, length("- Report: ") + 1)
}
event == "TASK-COMPLETE" && /^- report-lint: / {
  event_report_lint_count++
  event_report_lint = substr($0, length("- report-lint: ") + 1)
}
event == "TASK-COMPLETE" && /^- Commit: / {
  event_commit_count++
  event_commit = substr($0, length("- Commit: ") + 1)
}
event == "TASK-COMPLETE" && /^- QA-plan commands beyond the task spec: / {
  event_qa_commands_count++
  event_qa_commands = substr($0, length("- QA-plan commands beyond the task spec: ") + 1)
}
event == "RUN-END" && /^- All tasks complete\. Final review: / {
  event_run_end_complete_count++
  event_run_end_complete = substr($0, length("- All tasks complete. Final review: ") + 1)
}
event == "RUN-END" && /^- Human checklist: / {
  event_human_checklist_count++
  event_human_checklist = substr($0, length("- Human checklist: ") + 1)
}
event == "RUN-END" && /^- Phase ready for Tech Lead wrap-up\.$/ {
  event_wrapup_ready_count++
  event_wrapup_ready = "yes"
}

END {
  finish_event()
  if (saw_v2 && !saw_directive_schema) {
    invalid("ledger uses Sequence contract v2 without a Recovery directive schema")
  }
  if (errors) {
    exit 1
  }
  if (stop_entry_pending) {
    print "VALID state=STOP_PENDING"
  } else if (stopped) {
    print "VALID state=STOPPED"
  } else if (resume_recovery_pending) {
    print "VALID state=REPAIR_PENDING"
  } else if (triage_needed) {
    print "VALID state=TRIAGE_PENDING"
  } else if (triage_repair_pending) {
    print "VALID state=REPAIR_PENDING"
  } else if (recovery_task != "") {
    if (recovery_stage != "needs-complete") {
      print "VALID state=REPAIR_PENDING"
    } else {
      print "VALID state=COMPLETION_PENDING"
    }
  } else if (review_open) {
    if (review_task_complete) {
      print "VALID state=REREVIEW_PENDING"
    } else {
      print "VALID state=REPAIR_PENDING"
    }
  } else {
    print "VALID state=ADVANCEABLE"
  }
}
' "$1"
}

print_append_receipt() {
  case "$1" in
    'VALID state='*) state=${1#VALID state=} ;;
    *)
      echo "INVALID validator returned an unexpected result: $1" >&2
      exit 1
      ;;
  esac
  printf 'APPEND_CONFIRMED state=%s\n' "$state"
}

if [ "$mode" = validate ]; then
  validate_sequence "$ledger" "$ledger"
  exit $?
fi

ledger_dir=$(CDPATH= cd "$(dirname "$ledger")" && pwd)
ledger="$ledger_dir/$(basename "$ledger")"
tmp=$(mktemp -d "$ledger_dir/.workflow-sequence.XXXXXX")
trap 'rm -rf "$tmp"' EXIT HUP INT TERM
candidate="$tmp/candidate.md"
framed="$tmp/framed-candidate.md"
combined="$tmp/combined-ledger.md"
cat > "$candidate"

if [ ! -s "$candidate" ]; then
  echo "INVALID append candidate is empty" >&2
  exit 1
fi
if ! sed -n '1p' "$candidate" | grep -q '^### '; then
  echo "INVALID append candidate must start with an event heading" >&2
  exit 1
fi
heading_count=$(awk '/^### / { count++ } END { print count + 0 }' "$candidate")
if [ "$heading_count" -ne 1 ]; then
  echo "INVALID append candidate must contain exactly one event heading" >&2
  exit 1
fi
if awk 'NR == 1 { exit !($4 == "RUN-START" || $4 == "RUN-RESUME") }' "$candidate"; then
  contract_count=$(grep -c '^- Sequence contract: v2$' "$candidate" || true)
  schema_count=$(grep -c '^- Recovery directive schema: v2$' "$candidate" || true)
  if [ "$contract_count" -ne 1 ] || [ "$schema_count" -ne 1 ]; then
    echo "INVALID every newly appended RUN must contain exactly one Sequence contract: v2 and Recovery directive schema: v2" >&2
    exit 1
  fi
fi

printf '\n' > "$framed"
cat "$candidate" >> "$framed"

# An earlier invocation may have committed the event but lost its stdout
# receipt. An exact EOF match makes retrying the same candidate idempotent.
framed_size=$(wc -c < "$framed" | tr -d ' ')
tail -c "$framed_size" "$ledger" > "$tmp/current-suffix.md"
if cmp -s "$framed" "$tmp/current-suffix.md"; then
  if output=$(validate_sequence "$ledger" "$ledger"); then
    print_append_receipt "$output"
    exit $?
  else
    code=$?
    exit "$code"
  fi
fi

# Build beside the ledger so the validated replacement stays on one filesystem.
cp -p "$ledger" "$combined"
cat "$framed" >> "$combined"

if output=$(validate_sequence "$combined" "$ledger"); then
  :
else
  code=$?
  exit "$code"
fi

mv "$combined" "$ledger"
print_append_receipt "$output"
