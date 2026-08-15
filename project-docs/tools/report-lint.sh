#!/bin/sh
# report-lint.sh - deterministic completion-report verification.
#
# Verifies the mechanical evidence rules that the workflow previously asked
# the Phase Runner to check by prose-parsing:
#   1. Required completion-report headings, present and in order, including a
#      typed workflow outcome for deterministic orchestration.
#   2. Non-empty sentinel sections (concerns, dependency-budget note).
#   3. Command coverage: every command in the task spec's
#      `## Exact Commands To Run` required block appears in the report's
#      `## Tests run` or `## Exact quality-gate results` with a PASS result,
#      or BLOCKED inside a blocker report.
#   4. Commit hash resolves via `git rev-parse --verify <hash>^{commit}`.
#   5. Commit subject matches `git log -1 --format=%s <hash>` exactly, which
#      also catches plausible-but-wrong hash prefixes.
#   6. The report file itself is committed with no uncommitted modifications.
#      A report is evidence only once committed; the finalized hash-bearing
#      report lands in a `docs(task): ...` commit immediately after the task
#      commit, because a report cannot carry the hash of the same commit that
#      introduces its final content.
#
# Judgment items stay with agents and humans: QA-plan commands beyond the task
# spec, conditional-command applicability, deviation content, and whether the
# evidence supports a protected blocker. Mechanical lint returns one typed
# outcome: COMPLETE, RECOVERY_REQUIRED, or BLOCKED.
#
# Coverage matching requires each required command to appear verbatim inside
# backticks in the report (the completion-report template's format). A
# narrower substitute - a different command, or the required command with
# narrowing arguments appended - does not match and fails the gate.
set -u

usage() {
  cat <<'EOF'
Usage:
  project-docs/tools/report-lint.sh <completion-report.md> <task-spec.md> [--repo <git-dir>]
  project-docs/tools/report-lint.sh --self-test

Exit code is the number of failed checks (0 = clean).
--repo defaults to the git repository containing the completion report.
EOF
}

failures=0
report_outcome=""

pass() { echo "PASS $1"; }
note() { echo "NOTE $1"; }
fail() {
  echo "FAIL $1" >&2
  failures=$((failures + 1))
}

# Print the body of one `## `-level section: lines after the exact heading
# up to (excluding) the next `## ` heading.
section_body() {
  file=$1
  heading=$2
  awk -v h="$heading" '
    $0 == h { in_section = 1; next }
    in_section && /^## / { exit }
    in_section { print }
  ' "$file"
}

# Print required commands from the first fenced code block after
# `## Exact Commands To Run` in the task spec.
required_commands() {
  awk '
    $0 == "## Exact Commands To Run" { in_section = 1; next }
    in_section && /^## / { exit }
    in_section && /^```/ {
      if (in_block) exit
      in_block = 1
      next
    }
    in_block && NF > 0 { print }
  ' "$1"
}

# Reject evidence tokens that cannot have been printed by the recorded go test
# command. This is intentionally mechanical: it cannot prove the command ran,
# but it can catch a silent-success or unrelated token copied into its record.
validate_go_test_tokens() {
  while IFS= read -r line; do
    case "$line" in
      '- `'*'go test'*'`'*'- '*) ;;
      *) continue ;;
    esac

    case "$line" in
      *BLOCKED*|*UNVERIFIED*) continue ;;
    esac

    token=$(printf '%s\n' "$line" | awk -F '`' 'NF >= 4 { print $4 }')
    case "$token" in
      'exit 0 (no output)')
        fail "go test recorded an impossible silent-success token on line: $line. Long go test runs may print nothing until packages finish; wait for completion. A killed or timed-out run is BLOCKED, never PASS with a silent token."
        ;;
      *'ok '*|*FAIL*|*PASS*|*---*|*'no test files'*)
        pass "go test token looks like output: $token"
        ;;
      *)
        fail "go test token does not look like go test output on line: $line. Long go test runs may print nothing until packages finish; wait for completion. A killed or timed-out run is BLOCKED, never PASS with a silent token."
        ;;
    esac
  done <<EOF
$1
EOF
}

validate_pass_termination() {
  while IFS= read -r line; do
    case "$line" in
      '- `'*'` - PASS - '*) ;;
      *) continue ;;
    esac

    if ! printf '%s\n' "$line" | grep -Eq '[(]terminated `exit_code: [0-9][0-9]*`[)]'; then
      fail "PASS result missing observed termination: $line"
    elif ! printf '%s\n' "$line" | grep -Eq '[(]terminated `exit_code: 0`[)][[:space:]]*$'; then
      fail "PASS result did not terminate successfully with exit_code 0: $line"
    fi

    # Backtick fields are: command, token, and termination result. Requiring
    # all three prevents a bare PASS plus exit code from masquerading as
    # command-output evidence.
    token=$(printf '%s\n' "$line" | awk -F '`' 'NF >= 6 { print $4 }')
    if [ -z "$token" ]; then
      fail "PASS result missing verbatim output token: $line"
    fi
  done <<EOF
$1
EOF
}

validate_workflow_outcome() {
  report=$1
  body=$(section_body "$report" '## Workflow outcome')
  report_outcome=$(printf '%s\n' "$body" | sed -n 's/^- Task outcome: //p' | head -n1)
  category=$(printf '%s\n' "$body" | sed -n 's/^- Category: //p' | head -n1)
  condition=$(printf '%s\n' "$body" | sed -n 's/^- Condition: //p' | head -n1)
  attempts=$(printf '%s\n' "$body" | sed -n 's/^- Recovery attempts: //p' | head -n1)
  next_action=$(printf '%s\n' "$body" | sed -n 's/^- Next action: //p' | head -n1)
  human=$(printf '%s\n' "$body" | sed -n 's/^- Human attention required: //p' | head -n1)

  [ -n "$report_outcome" ] || fail "workflow outcome field missing: report_outcome"
  [ -n "$category" ] || fail "workflow outcome field missing: category"
  [ -n "$condition" ] || fail "workflow outcome field missing: condition"
  [ -n "$attempts" ] || fail "workflow outcome field missing: attempts"
  [ -n "$next_action" ] || fail "workflow outcome field missing: next_action"
  [ -n "$human" ] || fail "workflow outcome field missing: human"

  case "$attempts" in
    ''|*[!0-9]*) fail "Recovery attempts must be a non-negative integer: $attempts" ;;
  esac

  case "$report_outcome" in
    COMPLETE)
      [ "$category" = "none" ] || fail "COMPLETE category must be none: $category"
      [ "$human" = "no" ] || fail "COMPLETE cannot require Human attention"
      ;;
    RECOVERY_REQUIRED)
      case "$category" in
        in-scope-implementation|gate-repair|report-repair|qa-rework|session-continuation|local-environment) ;;
        *) fail "invalid RECOVERY_REQUIRED category: $category" ;;
      esac
      [ "$human" = "no" ] || fail "RECOVERY_REQUIRED cannot require Human attention"
      [ "$condition" != "None." ] || fail "RECOVERY_REQUIRED must name its condition"
      [ "$next_action" != "None." ] || fail "RECOVERY_REQUIRED must name its next action"
      ;;
    BLOCKED)
      case "$category" in
        product-decision|architecture-decision|data-model-decision|public-interface-decision|dependency-approval|security-or-data-loss|external-infrastructure|specification-contradiction|recovery-budget-exhausted) ;;
        *) fail "invalid BLOCKED category: $category" ;;
      esac
      [ "$human" = "yes" ] || fail "BLOCKED must require Human attention"
      [ "$condition" != "None." ] || fail "BLOCKED must name its condition"
      [ "$next_action" != "None." ] || fail "BLOCKED must name its resolution action"
      ;;
    *)
      fail "invalid Task outcome: $report_outcome"
      ;;
  esac

  [ -n "$report_outcome" ] && echo "OUTCOME $report_outcome"
}

# DEC-131 requires an additive correction to five already-complete Phase 20
# reports. The exception is deliberately narrower than normal gate coverage:
# only a corrected go test record can remain UNVERIFIED, and only when the
# report preserves the original impossible token in its appended correction.
has_dec131_go_test_correction() {
  report=$1
  cmd=$2
  case "$cmd" in
    *'go test'*) ;;
    *) return 1 ;;
  esac
  grep -Fqx '## Gate-evidence correction (DEC-131)' "$report" &&
    grep -F "\`$cmd\`" "$report" | grep -q 'UNVERIFIED' &&
    grep -Fq 'original recorded token `exit 0 (no output)`' "$report"
}

run_lint() {
  report=$1
  spec=$2
  repo=$3

  [ -f "$report" ] || { fail "completion report not found: $report"; return; }
  [ -f "$spec" ] || { fail "task spec not found: $spec"; return; }

  if [ -z "$repo" ]; then
    repo=$(CDPATH= cd "$(dirname "$report")" && pwd)
  fi
  if ! git -C "$repo" rev-parse --git-dir >/dev/null 2>&1; then
    fail "not a git repository (commit checks impossible): $repo"
    return
  fi

  # --- 1. Required headings, present and in order ---
  required_headings='## Workflow outcome
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
## Commit subject'

  prev_line=0
  order_ok=true
  missing=false
  while IFS= read -r heading; do
    line_no=$(grep -nFx "$heading" "$report" | head -n1 | cut -d: -f1)
    if [ -z "$line_no" ]; then
      fail "required heading missing: $heading"
      missing=true
      continue
    fi
    if [ "$line_no" -lt "$prev_line" ]; then
      order_ok=false
    fi
    prev_line=$line_no
  done <<EOF
$required_headings
EOF
  $missing || pass "all required headings present"
  if $order_ok; then
    pass "headings in required order"
  else
    fail "required headings out of order"
  fi

  # --- 2. Sentinel sections and typed workflow outcome ---
  for heading in \
    '## Concerns, technical debt, or notes for the Tech Lead' \
    '## Dependency-budget note' \
    '## Deviations from spec' \
    '## Validation Contract Coverage'; do
    body=$(section_body "$report" "$heading" | grep -cv '^[[:space:]]*$' || true)
    if [ "$body" -gt 0 ]; then
      pass "section non-empty: $heading"
    else
      fail "section empty (sentinel like '- None.' required): $heading"
    fi
  done

  validate_workflow_outcome "$report"

  # --- 3. Workflow disposition ---
  hash_body=$(section_body "$report" '## Commit hash')
  blocker=false
  recovery=false
  case "$report_outcome" in
    BLOCKED) blocker=true ;;
    RECOVERY_REQUIRED) recovery=true ;;
  esac

  # --- 4. Command coverage vs task spec ---
  cmds=$(required_commands "$spec")
  if [ -z "$cmds" ]; then
    fail "task spec has no required-command block under '## Exact Commands To Run': $spec"
  fi
  results=$(
    { section_body "$report" '## Tests run'
      section_body "$report" '## Exact quality-gate results'; } | grep -v '^[[:space:]]*$' || true
  )
  validate_go_test_tokens "$results"
  validate_pass_termination "$results"
  while IFS= read -r cmd; do
    [ -n "$cmd" ] || continue
    case "$cmd" in
      \<*\>)
        fail "required command is still a placeholder in task spec: $cmd"
        continue
        ;;
    esac
    match=$(printf '%s\n' "$results" | awk -F '`' -v cmd="$cmd" '$1 == "- " && $2 == cmd { print }')
    if [ -z "$match" ]; then
      fail "required command missing from report (skipped or substituted; must appear verbatim in backticks): $cmd"
    else
      final_status=""
      while IFS= read -r result_line; do
        status=$(printf '%s\n' "$result_line" | awk -F '`' '{ value=$3; sub(/^ - /, "", value); sub(/ - .*$/, "", value); print value }')
        case "$status" in
          PASS|FAIL|BLOCKED|UNVERIFIED) ;;
          *) fail "required command has invalid result status '$status': $cmd" ;;
        esac
        final_status=$status
      done <<EOF_MATCHES
$match
EOF_MATCHES

      case "$final_status" in
        PASS)
          pass "required command recorded PASS: $cmd"
          ;;
        BLOCKED)
          if $blocker || $recovery; then
            pass "required command recorded BLOCKED in non-complete report: $cmd"
          else
            fail "required command BLOCKED outside a blocker report: $cmd"
          fi
          ;;
        UNVERIFIED)
          if has_dec131_go_test_correction "$report" "$cmd"; then
            pass "required go test command recorded UNVERIFIED in DEC-131 historical correction: $cmd"
          else
            fail "required command recorded UNVERIFIED without an accepted historical correction: $cmd"
          fi
          ;;
        *)
          fail "required command final result is not PASS/BLOCKED: $cmd (final status: $final_status)"
          ;;
      esac
    fi
  done <<EOF
$cmds
EOF

  # --- 5 & 6. Commit hash validity and subject match ---
  if $blocker && printf '%s\n' "$hash_body" | grep -qi 'No commit created'; then
    note "commit hash/subject and report-committed checks skipped for no-commit protected blocker"
    return
  fi

  hash=$(printf '%s\n' "$hash_body" | sed -n 's/.*`\([0-9a-f][0-9a-f]*\)`.*/\1/p' | head -n1)
  if [ -z "$hash" ]; then
    fail "no commit hash found in '## Commit hash' section"
    return
  fi
  case "$hash" in
    *[!0-9a-f]*)
      fail "commit hash contains non-hex characters: $hash"
      return
      ;;
  esac
  if [ "${#hash}" -lt 7 ]; then
    fail "commit hash shorter than 7 characters: $hash"
    return
  fi

  if git -C "$repo" rev-parse --verify --quiet "$hash^{commit}" >/dev/null 2>&1; then
    pass "commit hash resolves: $hash"
  else
    fail "commit hash does not resolve (fabricated or unverifiable): $hash"
    return
  fi

  reported_subject=$(section_body "$report" '## Commit subject' \
    | sed -n 's/^- `\(.*\)`$/\1/p' | head -n1)
  if [ -z "$reported_subject" ]; then
    fail "no commit subject found in '## Commit subject' section (expected: - \`<subject>\`)"
    return
  fi
  actual_subject=$(git -C "$repo" log -1 --format=%s "$hash")
  if [ "$reported_subject" = "$actual_subject" ]; then
    pass "commit subject matches: $reported_subject"
  else
    fail "commit subject mismatch (wrong commit or wrong hash prefix): report '$reported_subject' vs git '$actual_subject'"
  fi

  # --- 7. Report file committed and clean ---
  report_abs=$(CDPATH= cd "$(dirname "$report")" && pwd)/$(basename "$report")
  if ! git -C "$repo" ls-files --error-unmatch -- "$report_abs" >/dev/null 2>&1; then
    fail "report file is not committed (a report is evidence only once committed): $report"
  elif [ -n "$(git -C "$repo" status --porcelain -- "$report_abs")" ]; then
    fail "report file has uncommitted modifications (post-commit edits invalidate lint evidence): $report"
  else
    pass "report file committed and clean"
  fi
}

# ---------------------------------------------------------------------------
# Self-test: builds a scratch git repo plus fixture reports, then asserts the
# linter accepts the clean case and rejects each seeded defect class.
# ---------------------------------------------------------------------------
self_test() {
  self=$(CDPATH= cd "$(dirname "$0")" && pwd)/$(basename "$0")
  tmp=$(mktemp -d "${TMPDIR:-/tmp}/report-lint-selftest.XXXXXX")
  trap 'rm -rf "$tmp"' EXIT HUP INT TERM
  st_failures=0

  git -C "$tmp" init -q
  git -C "$tmp" -c user.name=fixture -c user.email=fixture@example.invalid \
    commit -q --allow-empty -m "feat: add widget summary line"
  good_hash=$(git -C "$tmp" rev-parse --short HEAD)
  good_subject=$(git -C "$tmp" log -1 --format=%s)
  git -C "$tmp" -c user.name=fixture -c user.email=fixture@example.invalid \
    commit -q --allow-empty -m "chore: unrelated follow-up commit"
  other_hash=$(git -C "$tmp" rev-parse --short HEAD)

  cat > "$tmp/task-spec.md" <<'EOF'
# Phase 1 - Task 1: Fixture task

## Objective

Fixture.

## Exact Commands To Run

Run exactly:

```text
echo format-fixture
echo lint-fixture
echo coverage-fixture
go test ./...
go test -tags integration ./...
```

## Notes For Developer

- Fixture.
EOF

  write_report() {
    out=$1
    r_hash=$2
    r_subject=$3
    cat > "$out" <<EOF
# Phase 1 - Task 1: Fixture task - Completion Report

## Workflow outcome

- Task outcome: COMPLETE
- Category: none
- Condition: None.
- Recovery attempts: 0
- Next action: None.
- Human attention required: no

## Summary of what was implemented

- Fixture summary.

## Files created or modified

- \`fixture.txt\`

## Documentation updated

- None required.

## Validation Contract Coverage

- VAL-001: fixture evidence.

## Tests run


## Exact quality-gate results

- \`echo format-fixture\` - PASS - \`format-fixture\` (terminated \`exit_code: 0\`)
- \`echo lint-fixture\` - PASS - \`lint-fixture\` (terminated \`exit_code: 0\`)
- \`echo coverage-fixture\` - PASS - \`coverage-fixture\` (terminated \`exit_code: 0\`)
- \`go test ./...\` - PASS - \`ok  example.com/fixture  0.001s\` (terminated \`exit_code: 0\`)
- \`go test -tags integration ./...\` - PASS - \`ok  example.com/fixture/integration  0.002s\` (terminated \`exit_code: 0\`)

## Deviations from spec

- None.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- \`$r_hash\`

## Commit subject

- \`$r_subject\`
EOF
  }

  expect() {
    label=$1
    want=$2   # zero | nonzero
    match=$3  # required output substring, empty to skip
    file=$4
    output=$("$self" "$file" "$tmp/task-spec.md" --repo "$tmp" 2>&1)
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
      echo "SELF-TEST PASS $label"
    else
      echo "SELF-TEST FAIL $label (exit=$code)" >&2
      printf '%s\n' "$output" | sed 's/^/  | /' >&2
      st_failures=$((st_failures + 1))
    fi
  }

  commit_fixture() {
    git -C "$tmp" add "$1"
    git -C "$tmp" -c user.name=fixture -c user.email=fixture@example.invalid \
      commit -q -m "test: add report-lint fixture"
  }

  # Clean report must pass. Committed in a follow-up docs commit, mirroring
  # the required two-step pattern (task commit, then finalized report commit).
  write_report "$tmp/clean.md" "$good_hash" "$good_subject"
  git -C "$tmp" add clean.md task-spec.md
  git -C "$tmp" -c user.name=fixture -c user.email=fixture@example.invalid \
    commit -q -m "docs(task): finalize completion report fixture"
  expect "clean report passes with COMPLETE outcome" zero "OUTCOME COMPLETE" "$tmp/clean.md"
  expect "accepted real go test tokens pass" zero "" "$tmp/clean.md"

  # Every PASS must carry observed termination, not merely an output token.
  sed 's/ (terminated `exit_code: 0`)//' "$tmp/clean.md" > "$tmp/missing-termination.md"
  commit_fixture missing-termination.md
  expect "PASS without termination evidence fails" nonzero "PASS result missing observed termination" "$tmp/missing-termination.md"

  # PASS is an exact state and can only accompany observed exit code zero.
  sed '0,/exit_code: 0/s//exit_code: 1/' "$tmp/clean.md" > "$tmp/pass-nonzero.md"
  commit_fixture pass-nonzero.md
  expect "PASS with nonzero exit fails" nonzero "PASS result did not terminate successfully with exit_code 0" "$tmp/pass-nonzero.md"

  # Status parsing is exact; explanatory or negated PASS text is not success.
  sed 's/`echo format-fixture` - PASS -/`echo format-fixture` - NOT PASS -/' "$tmp/clean.md" > "$tmp/not-pass-status.md"
  commit_fixture not-pass-status.md
  expect "NOT PASS does not satisfy command coverage" nonzero "invalid result status 'NOT PASS'" "$tmp/not-pass-status.md"

  # A failed attempt remains visible, while a later byte-identical PASS is the
  # current required-command result. The command text is exact in both rows.
  sed '/^- `echo lint-fixture` - PASS -/i\
- `echo lint-fixture` - FAIL - `FAIL lint-fixture first attempt` (terminated `exit_code: 1`)\
' "$tmp/clean.md" > "$tmp/fail-then-pass.md"
  commit_fixture fail-then-pass.md
  expect "byte-identical FAIL then PASS preserves history and completes" zero "required command recorded PASS: echo lint-fixture" "$tmp/fail-then-pass.md"

  # A full go test run can print nothing until packages finish. A silent token
  # therefore cannot prove it terminated successfully.
  sed 's|`ok  example.com/fixture  0.001s`|`exit 0 (no output)`|' "$tmp/clean.md" > "$tmp/impossible-go-test-token.md"
  commit_fixture impossible-go-test-token.md
  expect "impossible go test token fails" nonzero "go test recorded an impossible silent-success token" "$tmp/impossible-go-test-token.md"

  # Non-go-test output cannot be copied in as evidence for a go test command.
  sed 's|`ok  example.com/fixture  0.001s`|`295 internal/config/project_config.go`|' "$tmp/clean.md" > "$tmp/mismatched-go-test-token.md"
  commit_fixture mismatched-go-test-token.md
  expect "mismatched go test token fails" nonzero "go test token does not look like go test output" "$tmp/mismatched-go-test-token.md"

  # DEC-131 corrections retain the historical record but do not let a
  # fabricated token remain a passing gate result.
  sed \
    -e 's|`go test ./...` - PASS - `ok  example.com/fixture  0.001s`|`go test ./...` - UNVERIFIED - see Gate-evidence correction (DEC-131)|' \
    -e 's|`go test -tags integration ./...` - PASS - `ok  example.com/fixture/integration  0.002s`|`go test -tags integration ./...` - UNVERIFIED - see Gate-evidence correction (DEC-131)|' \
    "$tmp/clean.md" > "$tmp/corrected-historical-go-test.md"
  cat >> "$tmp/corrected-historical-go-test.md" <<'EOF'

## Gate-evidence correction (DEC-131)

- The original recorded token `exit 0 (no output)` is retained as history.
EOF
  commit_fixture corrected-historical-go-test.md
  expect "corrected historical go test records pass" zero "DEC-131 historical correction" "$tmp/corrected-historical-go-test.md"

  # Missing required heading.
  sed '/^## Dependency-budget note$/,/^## Commit hash$/{/^## Dependency-budget note$/d; /^- No dependency changes\.$/d;}' \
    "$tmp/clean.md" > "$tmp/missing-heading.md"
  expect "missing heading fails" nonzero "required heading missing: ## Dependency-budget note" "$tmp/missing-heading.md"

  # Empty concerns section.
  sed '/^- None\.$/d' "$tmp/clean.md" > "$tmp/empty-sections.md"
  expect "empty sentinel section fails" nonzero "section empty" "$tmp/empty-sections.md"

  # Silently skipped required command.
  grep -v 'echo lint-fixture' "$tmp/clean.md" > "$tmp/skipped-command.md"
  expect "skipped required command fails" nonzero "required command missing from report" "$tmp/skipped-command.md"

  # Narrower substituted command (different, narrower command).
  sed 's|echo coverage-fixture|echo coverage-fixt|' "$tmp/clean.md" > "$tmp/narrow-command.md"
  expect "narrower substituted command fails" nonzero "required command missing from report" "$tmp/narrow-command.md"

  # Narrower substitution by appending narrowing arguments.
  sed 's|`echo coverage-fixture`|`echo coverage-fixture --only happy-path`|' "$tmp/clean.md" > "$tmp/narrow-args.md"
  expect "narrowing-arguments substitution fails" nonzero "required command missing from report" "$tmp/narrow-args.md"

  # BLOCKED outside a blocker report.
  sed 's|`echo lint-fixture` - PASS - `lint-fixture` (terminated `exit_code: 0`)|`echo lint-fixture` - BLOCKED - `FAIL` (terminated `exit_code: 1`)|' "$tmp/clean.md" > "$tmp/blocked-non-blocker.md"
  expect "BLOCKED outside blocker report fails" nonzero "BLOCKED outside a blocker report" "$tmp/blocked-non-blocker.md"

  # Fabricated commit hash (well-formed, resolves nowhere).
  write_report "$tmp/fabricated-hash.md" "0123456789abcdef0123456789abcdef01234567" "$good_subject"
  expect "fabricated hash fails" nonzero "commit hash does not resolve" "$tmp/fabricated-hash.md"

  # Plausible wrong hash: resolves, but belongs to a different commit.
  write_report "$tmp/wrong-prefix.md" "$other_hash" "$good_subject"
  expect "plausible wrong hash fails on subject mismatch" nonzero "commit subject mismatch" "$tmp/wrong-prefix.md"

  # Commit subject mismatch with a valid hash.
  write_report "$tmp/subject-mismatch.md" "$good_hash" "feat: add widget summary line polished"
  expect "subject mismatch fails" nonzero "commit subject mismatch" "$tmp/subject-mismatch.md"

  # A durable partial checkpoint is a typed recovery state, not a blocker.
  write_report "$tmp/recovery.md" "$good_hash" "$good_subject"
  sed -e 's|- Task outcome: COMPLETE|- Task outcome: RECOVERY_REQUIRED|' \
      -e 's|- Category: none|- Category: in-scope-implementation|' \
      -e 's|- Condition: None\.|- Condition: Focused tests still expose an in-scope migration gap.|' \
      -e 's|- Recovery attempts: 0|- Recovery attempts: 1|' \
      -e 's|- Next action: None\.|- Next action: Continue the same task and rerun the exact failed command.|' \
      -e 's| - PASS - | - BLOCKED - |g' \
      "$tmp/recovery.md" > "$tmp/recovery-final.md"
  commit_fixture recovery-final.md
  expect "recovery report returns typed recovery outcome" zero "OUTCOME RECOVERY_REQUIRED" "$tmp/recovery-final.md"

  # A well-formed protected blocker passes mechanical lint with a typed outcome.
  write_report "$tmp/blocker.md" "$good_hash" "$good_subject"
  sed -e 's|- Task outcome: COMPLETE|- Task outcome: BLOCKED|' \
      -e 's|- Category: none|- Category: product-decision|' \
      -e 's|- Condition: None\.|- Condition: The accepted behavior is contradictory.|' \
      -e 's|- Next action: None\.|- Next action: Human selects the product behavior.|' \
      -e 's|- Human attention required: no|- Human attention required: yes|' \
      -e 's|- `'"$good_hash"'`|- `No commit created - blocker report`|' \
      -e 's| - PASS - | - BLOCKED - |g' \
      "$tmp/blocker.md" > "$tmp/blocker-final.md"
  expect "protected blocker returns typed blocked outcome" zero "OUTCOME BLOCKED" "$tmp/blocker-final.md"

  # Incomplete in-scope work cannot be presented as a protected blocker.
  sed 's|- Category: product-decision|- Category: in-scope-implementation|' \
    "$tmp/blocker-final.md" > "$tmp/invalid-blocker.md"
  expect "recoverable category cannot masquerade as blocker" nonzero "invalid BLOCKED category" "$tmp/invalid-blocker.md"

  # Finalized report never committed.
  cp "$tmp/clean.md" "$tmp/uncommitted.md"
  expect "uncommitted report fails" nonzero "report file is not committed" "$tmp/uncommitted.md"

  # Committed report modified in the worktree afterwards (post-commit edit).
  printf -- '- Post-commit edit.\n' >> "$tmp/clean.md"
  expect "post-commit edited report fails" nonzero "uncommitted modifications" "$tmp/clean.md"
  git -C "$tmp" checkout -q -- clean.md

  if [ "$st_failures" -eq 0 ]; then
    echo "SELF-TEST PASS report-lint self-test completed with no failures"
    exit 0
  else
    echo "SELF-TEST FAIL report-lint self-test found $st_failures failure(s)" >&2
    exit 1
  fi
}

# --- argument parsing ---
[ "$#" -ge 1 ] || { usage >&2; exit 2; }

case "$1" in
  --self-test) self_test ;;
  -h|--help) usage; exit 0 ;;
esac

[ "$#" -ge 2 ] || { usage >&2; exit 2; }
report_arg=$1
spec_arg=$2
repo_arg=""
shift 2
while [ "$#" -gt 0 ]; do
  case "$1" in
    --repo)
      [ "$#" -ge 2 ] || { usage >&2; exit 2; }
      repo_arg=$2
      shift 2
      ;;
    *)
      usage >&2
      exit 2
      ;;
  esac
done

run_lint "$report_arg" "$spec_arg" "$repo_arg"

if [ "$failures" -eq 0 ]; then
  echo "PASS report-lint completed with no failures"
else
  echo "FAIL report-lint found $failures failure(s)" >&2
fi

[ "$failures" -gt 100 ] && failures=100
exit "$failures"
