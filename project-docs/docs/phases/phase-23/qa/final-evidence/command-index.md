# Phase 23 final command receipts

All commands below were executed in the third phase-end QA session. Each linked JSON contains the exact command, working directory (default /home/mark/projects/fid for older entries), every captured result including session IDs, full available output, and observed final exit_code. Prior child sessions provide no passing command evidence.

No full checkptr rerun was used to replace its failed attempt. Focused and planned named selections are separate evidence. Setup and cleanup are not gates. The truncated historical timing grep is retained as a diagnostic read; complete report receipts and timing excerpts supply that audit.

| Receipt | Result | Verbatim final output token | Observed termination |
|---|---|---|---|
| [absence-fence](absence-fence.json) | PASS (expected exit_code: 1) | `exit 1 (no output)` | `exit_code: 1` |
| [absence-v1-requirements](absence-v1-requirements.json) | PASS (expected exit_code: 1) | `exit 1 (no output)` | `exit_code: 1` |
| [absence-v1-rules](absence-v1-rules.json) | PASS (expected exit_code: 1) | `exit 1 (no output)` | `exit_code: 1` |
| [absence-v1-tables](absence-v1-tables.json) | PASS (expected exit_code: 1) | `exit 1 (no output)` | `exit_code: 1` |
| [all-action-provenance](all-action-provenance.json) | PASS | `PASS provenance recorded for all 30 grant actions; all seven Phase 23 introductions match owning tasks` | `exit_code: 0` |
| [archive-scale](archive-scale.json) | PASS | `ok  	github.com/mbell-nzl/fid/internal/export	142.700s` | `exit_code: 0` |
| [artifact-audit-final](artifact-audit-final.json) | PASS | `PASS final artifacts: 84 VAL outcomes (79 PASS, 5 BLOCKED), 52 QA checks, open full-checkptr failure, OUTCOME BLOCKED` | `exit_code: 0` |
| [artifact-audit](artifact-audit.json) | PASS | `PASS final artifacts: 84 VAL outcomes (79 PASS, 5 BLOCKED), 52 QA checks, open full-checkptr failure, OUTCOME BLOCKED` | `exit_code: 0` |
| [authority-producer-inventory](authority-producer-inventory.json) | PASS | `internal/cli/codes_test.go:71:		{CodeAuthorityDenied, CategoryAuthority, false, nil, authorityDetails()},` | `exit_code: 0` |
| [baseline-migration-diff](baseline-migration-diff.json) | PASS | `   version INTEGER NOT NULL CHECK (version > 0),` | `exit_code: 0` |
| [benchmark-1](benchmark-1.json) | PASS | `ok  	github.com/mbell-nzl/fid/internal/journal	25.864s` | `exit_code: 0` |
| [benchmark-2](benchmark-2.json) | PASS | `ok  	github.com/mbell-nzl/fid/internal/journal	26.751s` | `exit_code: 0` |
| [benchmark-3](benchmark-3.json) | PASS | `ok  	github.com/mbell-nzl/fid/internal/journal	25.203s` | `exit_code: 0` |
| [binary-checklist-identity](binary-checklist-identity.json) | PASS | `3a34fc318991f77b9aba5136e5e6a04538824010489dbae31fd854bae01f5d49  /home/mark/projects/fid/.cache/phase-23-human/checkpoints-3245767079/fid` | `exit_code: 0` |
| [budget-integration](budget-integration.json) | PASS | `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	262.002s` | `exit_code: 0` |
| [budget-ordinary](budget-ordinary.json) | PASS | `ok  	github.com/mbell-nzl/fid/internal/workflow	94.499s	coverage: 62.1% of statements` | `exit_code: 0` |
| [capabilities-final](capabilities-final.json) | PASS | `{"contract_version":"fid-cli.v1","event":null,"ok":true,"operation":"capabilities.get","result":{"items":[{"available":true,"capability_id":"attention.v1","limits":{"human_outcome_resolution":true,"state":true},"schema_max":1,"schema_min":1,"trust_source":"fid-journal"},{"available":true,"capability_id":"evidence.freshness.v1","limits":{"branch_labels":false,"target":"repository-tree"},"schema_max":1,"schema_min":1,"trust_source":"fid-git-observation"},{"available":true,"capability_id":"evidence.gate-definitions.v1","limits":{"executor_policy":"standard-library","network":"allowed","sandbox":"none","write":"allowed"},"schema_max":1,"schema_min":1,"trust_source":"active-kit-activation"},{"available":true,"capability_id":"evidence.receipts.v1","limits":{"attestation":false,"producer":"fid","terminal_receipts":true},"schema_max":1,"schema_min":1,"trust_source":"fid-executed"},{"available":true,"capability_id":"export.portable.v1","limits":{"create":true,"import":true,"import_mode":"empty-target","validate":true},"schema_max":1,"schema_min":1,"trust_source":"fid-journal"},{"available":true,"capability_id":"human.outcomes.v1","limits":{"provenance":"active-kit-activation","record":true},"schema_max":1,"schema_min":1,"trust_source":"active-kit-activation"},{"available":true,"capability_id":"kit.activation.v1","limits":{"abort_before_commit":true,"automatic_respawn":false,"finalize":true,"recover":true},"schema_max":1,"schema_min":1,"trust_source":"active-kit-activation"},{"available":true,"capability_id":"reviews.coverage.v1","limits":{"findings":true,"independence":"distinct-actor","record":true},"schema_max":1,"schema_min":1,"trust_source":"fid-journal"},{"available":false,"capability_id":"telemetry.trajectory.v1","limits":{"ingest":false},"schema_max":1,"schema_min":1,"trust_source":"phase-25-unavailable"},{"available":false,"capability_id":"work.claims.v1","limits":{"max_active_per_issue":1,"parallelism":false,"recover":true},"schema_max":1,"schema_min":1,"trust_source":"phase-24-unavailable"},{"available":true,"capability_id":"workflow.state.v1","limits":{"max_envelope_revisions":1,"max_primary_issues":1,"related_issues":false},"schema_max":1,"schema_min":1,"trust_source":"fid-journal"},{"available":true,"capability_id":"workflow.transitions.v1","limits":{"execution_bound":true,"terminal_evidence_guard":true},"schema_max":1,"schema_min":1,"trust_source":"fid-journal"}]}}` | `exit_code: 0` |
| [capability-staging](capability-staging.json) | PASS | `PASS task-11-close false and final true capability states on fresh initialized workspaces` | `exit_code: 0` |
| [checkptr-causal-file-state](checkptr-causal-file-state.json) | PASS | `Refs: phase-21/task-14` | `exit_code: 0` |
| [checkptr-failure-provenance](checkptr-failure-provenance.json) | PASS | `Refs: phase-20/task-3h` | `exit_code: 0` |
| [checkptr-fixture-diagnosis](checkptr-fixture-diagnosis.json) | PASS | `OBSERVATION: all AppendTerminal observed fields are copied except Signal; the mismatch therefore identifies a non-null process signal. The 5.01s failure meets the fixture 5000ms timeout. Child TestMain bootstraps SQLite before the output helper; contention/startup-time cause remains a hypothesis until isolated by repair.` | `exit_code: 0` |
| [command-index-complete](command-index-complete.json) | PASS | `PASS indexed 111 terminated command receipts with exact commands, final tokens and actual exit codes` | `exit_code: 0` |
| [command-index-final](command-index-final.json) | PASS | `PASS indexed 109 terminated command receipts with exact commands, final tokens and actual exit codes` | `exit_code: 0` |
| [command-index](command-index.json) | PASS | `PASS indexed 107 terminated command receipts with exact commands, final tokens and actual exit codes` | `exit_code: 0` |
| [dependency-diff](dependency-diff.json) | PASS | `exit 0 (no output)` | `exit_code: 0` |
| [diagnose-live-golden](diagnose-live-golden.json) | PASS | `ok  	github.com/mbell-nzl/fid/cmd/fid	7.232s` | `exit_code: 0` |
| [documentation-amendments](documentation-amendments.json) | PASS | `+**Failure code** and **Evaluation reason**)*` | `exit_code: 0` |
| [feasibility-fid](feasibility-fid.json) | FAIL | `AssertionError` | `exit_code: 1` |
| [feasibility-tools](feasibility-tools.json) | PASS | `PASS tool and flag feasibility` | `exit_code: 0` |
| [file-state](file-state.json) | PASS | `PASS final file-state and provenance audit` | `exit_code: 0` |
| [final-state-audit](final-state-audit.json) | PASS | `PASS 106 terminated command receipts; seven failed attempts retained; 52 QA checks cover 80 deterministic assertions` | `exit_code: 0` |
| [fixture-expiry-diff](fixture-expiry-diff.json) | PASS | `CONFIRMED expired retained Claim changes current resume state; fresh fixture generation required` | `exit_code: 0` |
| [full-checkptr](full-checkptr.json) | FAIL | `FAIL` | `exit_code: 1` |
| [gate-01-format](gate-01-format.json) | PASS | `exit 0 (no output)` | `exit_code: 0` |
| [gate-02-lint](gate-02-lint.json) | PASS | `0 issues.` | `exit_code: 0` |
| [gate-03-coverage](gate-03-coverage.json) | PASS | `ok  	github.com/mbell-nzl/fid/internal/workflow	93.800s	coverage: 62.1% of statements` | `exit_code: 0` |
| [gate-04-coverage-summary](gate-04-coverage-summary.json) | PASS | `total:											(statements)					60.5%` | `exit_code: 0` |
| [gate-05-integration](gate-05-integration.json) | PASS | `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	232.284s` | `exit_code: 0` |
| [gate-06-benchmarks](gate-06-benchmarks.json) | PASS | `ok  	github.com/mbell-nzl/fid/internal/workflow	9.995s` | `exit_code: 0` |
| [gate-07-build](gate-07-build.json) | PASS | `exit 0 (no output)` | `exit_code: 0` |
| [human-checklist-refresh-final](human-checklist-refresh-final.json) | PASS | `PASS fixture identities refreshed; all 63 blank Human rows preserved byte-for-byte` | `exit_code: 0` |
| [human-checklist-refresh](human-checklist-refresh.json) | PASS | `PASS fixture identities refreshed; all 63 blank Human rows preserved byte-for-byte` | `exit_code: 0` |
| [human-fixture-generation-final](human-fixture-generation-final.json) | PASS | `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	25.933s` | `exit_code: 0` |
| [human-fixture-generation](human-fixture-generation.json) | PASS | `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	28.411s` | `exit_code: 0` |
| [human-guide-diff](human-guide-diff.json) | PASS | `Guide section line: 749` | `exit_code: 0` |
| [human-live-final](human-live-final.json) | PASS | `PASS 36 retained Human/JSON leaf pairs and all checkpoint reads; 63 judgments remain blank` | `exit_code: 0` |
| [human-live-rerun](human-live-rerun.json) | FAIL | `AssertionError: concern.record` | `exit_code: 1` |
| [human-live](human-live.json) | FAIL | `AssertionError: ('VAL-070 mixed Findings', ['--project', '/home/mark/projects/fid/.cache/phase-23-human/checkpoints-3245767079/mixed-findings', 'resume', 'P22-bb90f3', '--repository', '/home/mark/projects/fid/.cache/phase-23-human/checkpoints-3245767079/mixed-findings'])` | `exit_code: 1` |
| [identity-cpu-count](identity-cpu-count.json) | PASS | `6` | `exit_code: 0` |
| [identity-cpu-model](identity-cpu-model.json) | PASS | `model name	: AMD Ryzen 7 5700G with Radeon Graphics` | `exit_code: 0` |
| [identity-env](identity-env.json) | PASS | `amd64` | `exit_code: 0` |
| [identity-go](identity-go.json) | PASS | `go version go1.26.2 linux/amd64` | `exit_code: 0` |
| [identity-kernel](identity-kernel.json) | PASS | `Linux 6.6.114.1-microsoft-standard-WSL2 x86_64` | `exit_code: 0` |
| [identity-memory](identity-memory.json) | PASS | `MemTotal:        8133236 kB` | `exit_code: 0` |
| [named-audit](named-audit.json) | PASS | `PASS 176 fresh selected tests; no missing/skipped selection or truncated inner result; all corruption receipts verified` | `exit_code: 0` |
| [named-claims](named-claims.json) | PASS | `ok  	github.com/mbell-nzl/fid/internal/claims	1.037s` | `exit_code: 0` |
| [named-cli](named-cli.json) | PASS | `ok  	github.com/mbell-nzl/fid/internal/cli	1.139s` | `exit_code: 0` |
| [named-cmd-fid](named-cmd-fid.json) | PASS | `ok  	github.com/mbell-nzl/fid/cmd/fid	53.488s` | `exit_code: 0` |
| [named-db](named-db.json) | PASS | `ok  	github.com/mbell-nzl/fid/internal/db	10.300s` | `exit_code: 0` |
| [named-evidence](named-evidence.json) | PASS | `ok  	github.com/mbell-nzl/fid/internal/evidence	1.050s` | `exit_code: 0` |
| [named-export](named-export.json) | PASS | `ok  	github.com/mbell-nzl/fid/internal/export	10.451s` | `exit_code: 0` |
| [named-feasibility](named-feasibility.json) | PASS | `PASS 167 named tests, every failure-code producer, and referenced fixtures are executable scopes` | `exit_code: 0` |
| [named-integration](named-integration.json) | PASS | `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	63.320s` | `exit_code: 0` |
| [named-inventory-expansion](named-inventory-expansion.json) | PASS | `PASS 176 current named tests resolve across all task specs, reports and QA-plan references; historical removed fences and temporary planning probes remain historical evidence only` | `exit_code: 0` |
| [named-issues](named-issues.json) | PASS | `ok  	github.com/mbell-nzl/fid/internal/issues	11.026s` | `exit_code: 0` |
| [named-journal](named-journal.json) | PASS | `ok  	github.com/mbell-nzl/fid/internal/journal	20.806s` | `exit_code: 0` |
| [named-links](named-links.json) | PASS | `ok  	github.com/mbell-nzl/fid/internal/links	1.018s` | `exit_code: 0` |
| [named-maintenance](named-maintenance.json) | PASS | `ok  	github.com/mbell-nzl/fid/internal/maintenance	2.836s` | `exit_code: 0` |
| [named-review](named-review.json) | PASS | `ok  	github.com/mbell-nzl/fid/internal/review	1.169s` | `exit_code: 0` |
| [named-workflow](named-workflow.json) | PASS | `ok  	github.com/mbell-nzl/fid/internal/workflow	19.639s` | `exit_code: 0` |
| [preservation-rows](preservation-rows.json) | PASS | `PASS linked-seed row preservation; exactly two added columns are null; foreign_key_check=[]; integrity_check=ok` | `exit_code: 0` |
| [protected-file-hashes](protected-file-hashes.json) | PASS | `f3b0081004bffce43b06755fc9a9ec7653a48ac5a8d3316b373cd491b46005c2  project-docs/docs/phases/phase-23/qa/human-checklist.md` | `exit_code: 0` |
| [rebuild-pragmas](rebuild-pragmas.json) | PASS | `PASS full journal rebuild: foreign_key_check=[]; integrity_check=ok; source journal unchanged` | `exit_code: 0` |
| [reference-binary-identity](reference-binary-identity.json) | PASS | `	build	vcs.modified=false` | `exit_code: 0` |
| [reference-comparison](reference-comparison.json) | PASS | `PASS six canonical bytes and Attention key identical at 74c83b7 and e0893cd; pinned journal identity verified` | `exit_code: 0` |
| [reference-current-probe](reference-current-probe.json) | PASS | `ok  	github.com/mbell-nzl/fid/internal/workflow	3.968s` | `exit_code: 0` |
| [reference-prechange-probe](reference-prechange-probe.json) | PASS | `{"command": ["go", "test", "-race", "-count=1", "-overlay", "/home/mark/projects/fid/project-docs/docs/phases/phase-23/qa/final-evidence/reference-overlay.json", "-run", "^TestQAPhase23ReferenceAttentionKeyAndBytes$", "-v", "./internal/workflow"], "cwd": "/tmp/fid-p23-qa-reference-PgIz9Jst/prechange", "environment": {"GIT_DIR": "/home/mark/projects/fid/.git/worktrees/prechange", "GIT_WORK_TREE": "/tmp/fid-p23-qa-reference-PgIz9Jst/prechange"}, "stdout": "=== RUN   TestQAPhase23ReferenceAttentionKeyAndBytes\n=== PAUSE TestQAPhase23ReferenceAttentionKeyAndBytes\n=== CONT  TestQAPhase23ReferenceAttentionKeyAndBytes\n    qa_reference_test.go:19: canonical executor-terminated: {\"type\":\"executor-terminated\",\"subject\":{\"type\":\"run\",\"id\":\"RUN-pinned\"},\"budget_name\":null,\"capability_id\":null}\n    qa_reference_test.go:19: canonical run-recovered: {\"type\":\"run-recovered\",\"subject\":{\"type\":\"run\",\"id\":\"RUN-pinned\"},\"budget_name\":null,\"capability_id\":null}\n    qa_reference_test.go:19: canonical activation-recovered: {\"type\":\"activation-recovered\",\"subject\":{\"type\":\"activation\",\"id\":\"ACT-pinned\"},\"budget_name\":null,\"capability_id\":null}\n    qa_reference_test.go:19: canonical capability-available: {\"type\":\"capability-available\",\"subject\":{\"type\":\"activation\",\"id\":\"ACT-pinned\"},\"budget_name\":null,\"capability_id\":\"attention.v1\"}\n    qa_reference_test.go:19: canonical budget-observed: {\"type\":\"budget-observed\",\"subject\":{\"type\":\"run\",\"id\":\"RUN-pinned\"},\"budget_name\":\"attempts\",\"capability_id\":null}\n    qa_reference_test.go:19: canonical work-ended: {\"type\":\"work-ended\",\"subject\":{\"type\":\"run\",\"id\":\"RUN-pinned\"},\"budget_name\":null,\"capability_id\":null}\n    qa_reference_test.go:25: subject={type:run,id:RUN-pinned} code=budget-exhausted condition={type:budget-observed,subject:{type:run,id:RUN-pinned},budget_name:attempts,capability_id:null} key=sha256:d4b51bb4a6404d4a2579826ad5c08d802ed83c41814cd5fa74d3ee8c7796c3de\n--- PASS: TestQAPhase23ReferenceAttentionKeyAndBytes (0.00s)\nPASS\nok  \tgithub.com/mbell-nzl/fid/internal/workflow\t3.348s\n", "stderr": "", "exit_code": 0}` | `exit_code: 0` |
| [reference-prechange-worktree](reference-prechange-worktree.json) | PASS | `HEAD is now at 74c83b7 docs(qa): confirm phase 23 task 9 planning contracts` | `exit_code: 0` |
| [reference-task11-build-bound](reference-task11-build-bound.json) | PASS | `{"command": ["go", "build", "-o", "/tmp/fid-p23-qa-reference-PgIz9Jst/fid-task11", "./cmd/fid"], "cwd": "/tmp/fid-p23-qa-reference-PgIz9Jst/task11", "environment": {"GIT_DIR": "/home/mark/projects/fid/.git/worktrees/task111", "GIT_WORK_TREE": "/tmp/fid-p23-qa-reference-PgIz9Jst/task11"}, "stdout": "", "stderr": "", "exit_code": 0}` | `exit_code: 0` |
| [reference-task11-build-diagnostic](reference-task11-build-diagnostic.json) | FAIL | `	Use -buildvcs=false to disable VCS stamping.` | `exit_code: 1` |
| [reference-task11-build](reference-task11-build.json) | FAIL | `	Use -buildvcs=false to disable VCS stamping.` | `exit_code: 1` |
| [reference-task11-git-status](reference-task11-git-status.json) | PASS | `exit 0 (no output)` | `exit_code: 0` |
| [reference-task11-worktree](reference-task11-worktree.json) | PASS | `HEAD is now at 3dff457 perf(test): parallelize serial test suites for ESC-23-005` | `exit_code: 0` |
| [reference-worktree-setup](reference-worktree-setup.json) | FAIL | `AssertionError` | `exit_code: 1` |
| [report-01](report-01.json) | PASS | `PASS report-lint completed with no failures` | `exit_code: 0` |
| [report-02](report-02.json) | PASS | `PASS report-lint completed with no failures` | `exit_code: 0` |
| [report-03](report-03.json) | PASS | `PASS report-lint completed with no failures` | `exit_code: 0` |
| [report-04](report-04.json) | PASS | `PASS report-lint completed with no failures` | `exit_code: 0` |
| [report-05](report-05.json) | PASS | `PASS report-lint completed with no failures` | `exit_code: 0` |
| [report-06](report-06.json) | PASS | `PASS report-lint completed with no failures` | `exit_code: 0` |
| [report-07](report-07.json) | PASS | `PASS report-lint completed with no failures` | `exit_code: 0` |
| [report-08](report-08.json) | PASS | `PASS report-lint completed with no failures` | `exit_code: 0` |
| [report-09](report-09.json) | PASS | `PASS report-lint completed with no failures` | `exit_code: 0` |
| [report-10](report-10.json) | PASS | `PASS report-lint completed with no failures` | `exit_code: 0` |
| [report-11](report-11.json) | PASS | `PASS report-lint completed with no failures` | `exit_code: 0` |
| [report-12](report-12.json) | PASS | `PASS report-lint completed with no failures` | `exit_code: 0` |
| [report-13](report-13.json) | PASS | `PASS report-lint completed with no failures` | `exit_code: 0` |
| [report-14](report-14.json) | PASS | `PASS report-lint completed with no failures` | `exit_code: 0` |
| [report-15](report-15.json) | PASS | `PASS report-lint completed with no failures` | `exit_code: 0` |
| [report-16](report-16.json) | PASS | `PASS report-lint completed with no failures` | `exit_code: 0` |
| [report-17](report-17.json) | PASS | `PASS report-lint completed with no failures` | `exit_code: 0` |
| [report-18](report-18.json) | PASS | `PASS report-lint completed with no failures` | `exit_code: 0` |
| [report-timing-audit](report-timing-audit.json) | PASS | `PASS timing records present for all 18 reports; final uncached controls recorded separately` | `exit_code: 0` |
| [schema-column-contradiction](schema-column-contradiction.json) | PASS | `CONFIRMED contradictory contract count: three required; accepted DDL adds two` | `exit_code: 0` |
| [schema-sql](schema-sql.json) | PASS | `PASS independently applied schema-17 and schema-18 pragma and inventory probes` | `exit_code: 0` |
| [security-source-audit](security-source-audit.json) | PASS | `Security source observations captured; behavioral proof is supplied by named tests and real-process journeys` | `exit_code: 0` |
| [setup-budget-testcache](setup-budget-testcache.json) | SETUP ONLY | `exit 0 (no output)` | `exit_code: 0` |
| [setup-testcache](setup-testcache.json) | SETUP ONLY | `exit 0 (no output)` | `exit_code: 0` |
| [staged-whitespace-after-repair](staged-whitespace-after-repair.json) | PASS | `exit 0 (no output)` | `exit_code: 0` |
| [staged-whitespace](staged-whitespace.json) | FAIL | `project-docs/docs/phases/phase-23/qa/final-evidence/schema-sql.py:30: new blank line at EOF.` | `exit_code: 2` |
| [tmp-after-budget-integration](tmp-after-budget-integration.json) | PASS | `/tmp/fid-integration-binary-894728768` | `exit_code: 0` |
| [tmp-after-integration](tmp-after-integration.json) | PASS | `disk_usage /tmp: usage(total=4164218880, used=1369890816, free=2794328064)` | `exit_code: 0` |
| [tmp-before-budget-integration](tmp-before-budget-integration.json) | PASS | `disk_usage /tmp: usage(total=4164218880, used=1376452608, free=2787766272)` | `exit_code: 0` |
| [tmp-before-integration](tmp-before-integration.json) | PASS | `disk_usage /tmp: usage(total=4164218880, used=1330049024, free=2834169856)` | `exit_code: 0` |

## Exact commands

Multiline Python probes retain their exact full source in their receipt and, where applicable, in the adjacent script. Their child-command receipts record actual return codes.

### absence-fence

Working directory: `/home/mark/projects/fid`

```text
rg -n unavailableJudgmentEvents cmd internal
```

### absence-v1-requirements

Working directory: `/home/mark/projects/fid`

```text
rg -n --glob '!*_test.go' "requirement\.synced|last_synced_hash" cmd internal
```

### absence-v1-rules

Working directory: `/home/mark/projects/fid`

```text
rg -n --glob '!*_test.go' "has_completion_report|concerns_addressed" cmd internal
```

### absence-v1-tables

Working directory: `/home/mark/projects/fid`

```text
rg -n --glob '!*_test.go' "task_completions|ac_verifications" cmd internal
```

### all-action-provenance

Working directory: `/home/mark/projects/fid`

Exact multiline invocation: [all-action-provenance.json](all-action-provenance.json).

### archive-scale

Working directory: `/home/mark/projects/fid`

```text
go test -tags archive_scale -race ./internal/export -run '^TestArchive(ExpandedByteBoundaries|CumulativeBoundary|EntryCountBoundaries)$' -count=1
```

### artifact-audit-final

Working directory: `/home/mark/projects/fid`

```text
python3 project-docs/docs/phases/phase-23/qa/final-evidence/audit-artifacts.py
```

### artifact-audit

Working directory: `/home/mark/projects/fid`

```text
python3 project-docs/docs/phases/phase-23/qa/final-evidence/audit-artifacts.py
```

### authority-producer-inventory

Working directory: `/home/mark/projects/fid`

```text
grep -rn 'CodeAuthority\(Missing\|UnknownActor\|Denied\)' --include=*.go cmd internal
```

### baseline-migration-diff

Working directory: `/home/mark/projects/fid`

```text
git diff 0195897a347dfe1872db0f7be3800c5c3ccdefcf -- migrations/0013_p21_operational_kernel.up.sql migrations/0015_p22_recovery.up.sql
```

### benchmark-1

Working directory: `/home/mark/projects/fid`

```text
go test -run '^$' -bench '^BenchmarkRebuildWorkspace100kEvents$' -benchtime 1x ./internal/journal/
```

### benchmark-2

Working directory: `/home/mark/projects/fid`

```text
go test -run '^$' -bench '^BenchmarkRebuildWorkspace100kEvents$' -benchtime 1x ./internal/journal/
```

### benchmark-3

Working directory: `/home/mark/projects/fid`

```text
go test -run '^$' -bench '^BenchmarkRebuildWorkspace100kEvents$' -benchtime 1x ./internal/journal/
```

### binary-checklist-identity

Working directory: `/home/mark/projects/fid`

```text
sha256sum /tmp/fid /home/mark/projects/fid/.cache/phase-23-human/checkpoints-3245767079/fid
```

### budget-integration

Working directory: `/home/mark/projects/fid`

```text
go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...
```

### budget-ordinary

Working directory: `/home/mark/projects/fid`

```text
go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...
```

### capabilities-final

Working directory: `/home/mark/projects/fid`

```text
/tmp/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-3245767079/mixed-findings --json capabilities
```

### capability-staging

Working directory: `/home/mark/projects/fid`

```text
python3 project-docs/docs/phases/phase-23/qa/final-evidence/capability-staging.py
```

### checkptr-causal-file-state

Working directory: `/home/mark/projects/fid`

```text
git log -S'TimeoutMS:        5000' --format='%H %ad %s%n%b' --date=iso-strict -- cmd/fid/live_output_test.go
```

### checkptr-failure-provenance

Working directory: `/home/mark/projects/fid`

```text
git log --format='%H %s%n%b' -- cmd/fid/live_output_test.go cmd/fid/workspace_fixture_test.go
```

### checkptr-fixture-diagnosis

Working directory: `/home/mark/projects/fid`

Exact multiline invocation: [checkptr-fixture-diagnosis.json](checkptr-fixture-diagnosis.json).

### command-index-complete

Working directory: `/home/mark/projects/fid`

```text
python3 project-docs/docs/phases/phase-23/qa/final-evidence/index-receipts.py
```

### command-index-final

Working directory: `/home/mark/projects/fid`

```text
python3 project-docs/docs/phases/phase-23/qa/final-evidence/index-receipts.py
```

### command-index

Working directory: `/home/mark/projects/fid`

```text
python3 project-docs/docs/phases/phase-23/qa/final-evidence/index-receipts.py
```

### dependency-diff

Working directory: `/home/mark/projects/fid`

```text
git diff 0195897 -- go.mod go.sum
```

### diagnose-live-golden

Working directory: `/home/mark/projects/fid`

```text
go test -race -run '^TestLiveOutputGoldenFixtures$' -count=1 -v ./cmd/fid
```

### documentation-amendments

Working directory: `/home/mark/projects/fid`

```text
git diff 0195897 -- project-docs/specs/02_TRD.md project-docs/specs/03_TRD_SCHEMA.md project-docs/specs/06_DESIGN_PLAN_CLI.md project-docs/specs/GLOSSARY.md
```

### feasibility-fid

Working directory: `/home/mark/projects/fid`

Exact multiline invocation: [feasibility-fid.json](feasibility-fid.json).

### feasibility-tools

Working directory: `/home/mark/projects/fid`

Exact multiline invocation: [feasibility-tools.json](feasibility-tools.json).

### file-state

Working directory: `/home/mark/projects/fid`

```text
python3 project-docs/docs/phases/phase-23/qa/final-evidence/file-state.py
```

### final-state-audit

Working directory: `/home/mark/projects/fid`

```text
python3 project-docs/docs/phases/phase-23/qa/final-evidence/audit-final-state.py
```

### fixture-expiry-diff

Working directory: `/home/mark/projects/fid`

Exact multiline invocation: [fixture-expiry-diff.json](fixture-expiry-diff.json).

### full-checkptr

Working directory: `/home/mark/projects/fid`

```text
go test -race -timeout 30m ./...
```

### gate-01-format

Working directory: `/home/mark/projects/fid`

```text
gofmt -l cmd internal
```

### gate-02-lint

Working directory: `/home/mark/projects/fid`

```text
golangci-lint run ./...
```

### gate-03-coverage

Working directory: `/home/mark/projects/fid`

```text
go test -race -gcflags=modernc.org/...=-d=checkptr=0 -coverprofile=coverage.out ./...
```

### gate-04-coverage-summary

Working directory: `/home/mark/projects/fid`

```text
go tool cover -func=coverage.out
```

### gate-05-integration

Working directory: `/home/mark/projects/fid`

```text
go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 ./cmd/fid/integration/...
```

### gate-06-benchmarks

Working directory: `/home/mark/projects/fid`

```text
go test -bench . -benchmem ./...
```

### gate-07-build

Working directory: `/home/mark/projects/fid`

```text
go build -o /tmp/fid ./cmd/fid
```

### human-checklist-refresh-final

Working directory: `/home/mark/projects/fid`

```text
python3 project-docs/docs/phases/phase-23/qa/final-evidence/refresh-human-checklist.py /home/mark/projects/fid/.cache/phase-23-human/checkpoints-3134551499/manifest.json
```

### human-checklist-refresh

Working directory: `/home/mark/projects/fid`

```text
python3 project-docs/docs/phases/phase-23/qa/final-evidence/refresh-human-checklist.py /home/mark/projects/fid/.cache/phase-23-human/checkpoints-526363314/manifest.json
```

### human-fixture-generation-final

Working directory: `/home/mark/projects/fid`

```text
FID_P23_HUMAN_FIXTURE_ROOT=/home/mark/projects/fid/.cache/phase-23-human go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -run '^TestP23HumanFixtures$' -count=1 -v ./cmd/fid/integration
```

### human-fixture-generation

Working directory: `/home/mark/projects/fid`

```text
FID_P23_HUMAN_FIXTURE_ROOT=/home/mark/projects/fid/.cache/phase-23-human go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -run '^TestP23HumanFixtures$' -count=1 -v ./cmd/fid/integration
```

### human-guide-diff

Working directory: `/home/mark/projects/fid`

Exact multiline invocation: [human-guide-diff.json](human-guide-diff.json).

### human-live-final

Working directory: `/home/mark/projects/fid`

```text
python3 project-docs/docs/phases/phase-23/qa/final-evidence/human-live.py
```

### human-live-rerun

Working directory: `/home/mark/projects/fid`

```text
python3 project-docs/docs/phases/phase-23/qa/final-evidence/human-live.py
```

### human-live

Working directory: `/home/mark/projects/fid`

```text
python3 project-docs/docs/phases/phase-23/qa/final-evidence/human-live.py
```

### identity-cpu-count

Working directory: `/home/mark/projects/fid`

```text
nproc
```

### identity-cpu-model

Working directory: `/home/mark/projects/fid`

```text
grep -m1 '^model name' /proc/cpuinfo
```

### identity-env

Working directory: `/home/mark/projects/fid`

```text
go env GOOS GOARCH
```

### identity-go

Working directory: `/home/mark/projects/fid`

```text
go version
```

### identity-kernel

Working directory: `/home/mark/projects/fid`

```text
uname -srm
```

### identity-memory

Working directory: `/home/mark/projects/fid`

```text
grep -m1 '^MemTotal:' /proc/meminfo
```

### named-audit

Working directory: `/home/mark/projects/fid`

```text
python3 project-docs/docs/phases/phase-23/qa/final-evidence/audit-named.py
```

### named-claims

Working directory: `/home/mark/projects/fid`

```text
go test -race -count=1 -run '^(TestClaimAcquireRejectsOpenReviewReworkStop)$' -v ./internal/claims
```

### named-cli

Working directory: `/home/mark/projects/fid`

```text
go test -race -count=1 -run '^(TestReasonCodeDetailSchemaIsExactForEveryCode|TestPhase23FailureCodesAreRegistered)$' -v ./internal/cli
```

### named-cmd-fid

Working directory: `/home/mark/projects/fid`

```text
go test -race -count=1 -run '^(TestPristineWorkspaceClonePreservesRebuildParityAndIsolation|TestMaintenanceCheckUsesFixed300SecondGraceAndRejectsFormat|TestReviewRecordRejectsNonIndependentActor|TestReviewRecordCreatesImmutableReview|TestReviewRecordIndependentReworkAndReceiptBinding|TestConcernListRetainsIssueScopedV1ReadUntilRemoval|TestConcernListGenericReadExcludesV1Rows|TestPhase23HumanLeafTextMatchesRecordedFacts|TestPhase22AttentionReadBoundary|TestLiveLeafCommandsMatchContractInventories|TestEnvelopeCreateAcceptsTypedObligationsAndProtectedDecisions|TestEnvelopeCreateCommandCoversObligationAndDecisionBoundaries|TestEnvelopeCreateCommandRejectsEveryTypedAuthoringNegativeWithoutPersistenceChange|TestEnvelopeCreateCommandRejectsEveryForbiddenTypedFieldCombinationWithoutPersistenceChange|TestLinkAddRejectsRequirementEndpointsAndInvalidInputBeforeAppend|TestResumeReportsRealFindingsAndJudgmentSections|TestRequirementCommandHasReferenceLeaves|TestEvidenceLinkAndEvaluateKeepsTrustedDimensionsIndependent|TestEvidenceEvaluateNullObligationReturnsActiveRevisionInIDOrder|TestClaimAcquireAcceptsReviewMode|TestClaimAcquireKeepsOneActiveClaimAcrossModes|TestDecodeClosedReviewInputRejectsMalformedInapplicability|TestFindingAddressDoesNotResolveOrUnblock|TestReviewRecordMixedTerminalDispositions|TestArtifactVerifyObservesLossAndRestorationWithoutIdentityMutation|TestGateRegisterUnknownActorReturnsContractAuthorityFailure|TestGateRegisterDeniedActorReturnsContractAuthorityFailure|TestCapabilitiesJSONGoldenAndProjectionRows|TestLiveOutputGoldenFixtures)$' -v ./cmd/fid
```

### named-db

Working directory: `/home/mark/projects/fid`

```text
go test -race -count=1 -run '^(TestMigration0017CreatesExactlyTheAcceptedObjects|TestMigration0017ObligationRegistrationDiscriminator|TestPhase23BaselineMigrationDiffIsExactlySevenAcceptedHunks|TestMigration0017PreservesReceiptEvidenceAndArtifactLinksInPlace|TestMigration0018DeletesOnlyLegacyNonIssueReferences|TestMigration0018CurrentReferenceGuards|TestMigration0018IssuesEnrichedColumnsAndSubjectCount|TestLatestSchemaHasNoV1JudgmentTables|TestMigration0018DownRestoresEmptyV1Tables|TestMigration0018DownMatchesFreshSchema17|TestMigration0018DownPreservesAgentRunAndFullRollback|TestMigration0018RoundTripIntegrity)$' -v ./internal/db
```

### named-evidence

Working directory: `/home/mark/projects/fid`

```text
go test -race -count=1 -run '^(TestHumanOutcomeAdmissionSatisfiesOnlyMatchingClassWithApprovedResult|TestReviewEvidenceStalesOnTargetTreeChange|TestHumanOutcomeEvaluationKeepsFreshnessAndUsesSupersession|TestValidateGateRunReferencesAuthorizesAgainstEveryEnvelopeGrant|TestEvidenceAdmissionTask6PureProductionClassifierMatrix|TestReviewAdmissionReturnsBlockingFindingIDs|TestLoadReviewEvidenceBindingEnforcesVerdictTypeAndLinkTimeIndependence|TestLoadHumanOutcomeEvidenceBindingRejectsClassMismatchAsOutcomeNotPermitted)$' -v ./internal/evidence
```

### named-export

Working directory: `/home/mark/projects/fid`

```text
go test -race -count=1 -run '^(TestCreateIsByteIdenticalWithCompleteDeterministicInventory)$' -v ./internal/export
```

### named-feasibility

Working directory: `/home/mark/projects/fid`

Exact multiline invocation: [named-feasibility.json](named-feasibility.json).

### named-integration

Working directory: `/home/mark/projects/fid`

```text
go test -tags integration -race -count=1 -run '^(TestP23DecisionReferenceProtectedLifecycleRealProcess|TestP23FailurePrecedenceWithProtectedDecisionAndObligation|TestP23FailurePrecedenceWithReadinessStopAndObligation|TestP23CancelledFamiliesAbsentFromProduct|TestP22FinalActivatedActionSetAcceptsProfilesAndGrants|TestP23HumanOutcomeAuthorityAndCompatibilityMatrixRealProcess|TestP23HumanOutcomeRequiredRaisedByEnvelopeActivation|TestP23HumanOutcomeRequiredRaisedByReviewRecord|TestP23HumanOutcomeRequiredRaisedByArtifactAvailability|TestV01ClosureJourneyAndRebuildParity|TestP23AttentionResolutionThroughReviewRecordedAndFindingResolvedRealProcess|TestP23AttentionResolutionHumanOutcomeReadinessRejectsWrongClassAndWrongResultAndLeavesUnrelatedStopsOpen|TestP23ConcernDeviationRealProcessJourney|TestPhase23LeavesRejectUnknownRequestFields|TestPhase23MutationsRejectMalformedObjectsWithoutEvents|TestGrantActionSetIsExactlyThirtyNames|TestProcessExitContract|TestP23ResumeRejectsJudgmentSectionOverflow|TestP22PortableFakeRejectsRemovedRequirementEvent|TestPortableRoundTripCoversEveryPhase23Table|TestPortableAssertCorruptEventsRejected|TestP23AmendedEnvelopeReplayThroughBuiltBinary|TestP23FindingDispositionRealProcessJourney|TestP23FreshnessAndOutcomeSupersessionRealProcessJourney|TestP23StopResolveHumanOutcomeConditionKeysOnConditionNotRaisingCode|TestP23StopResolveDecisionSatisfiedConditionKeysOnConditionNotRaisingCode|TestP23ConcernDispositionReferenceRejections|TestFinalCapabilityMatrixMatchesAcceptedTable|TestP23ReviewClaimRejectsDirtyPathAdoption|TestP23ReviewRunRetainsPhase22Lifecycle|TestP23HumanFixtures|TestP23HumanFixtureStepsRemainRunnable|TestP23GeneratedValueSyntax|TestP23FailureCorpusHasEveryAcceptedCode|TestP23CallerProvidedReviewIDsRemainDistinct|TestP23SharedLeafConformance|TestP23SharedAmendedLeafConformance|TestP23ReviewObligationGuardingBothActionsBlocksBothRealProcess|TestP23UnsatisfiedObligationBlocksOnlyTheActionItGuardsRealProcess|TestP23ConcernListHumanProjection|TestP23GateRunAuthorizesEveryEnvelopeGrantThroughBuiltBinary|TestP23TenNameActivityInventoryFromProducingLeaves|TestP23RequirementReferenceAndLinkRealProcess|TestP23ConditionWireGrammarMatchesJournalAndFixtures|TestP22ClosedValueSchemaCases|TestP23ReviewRecordRealProcessJourney|TestP23ReviewEvidenceRealProcessJourney|TestP23ReviewCyclesRealProcessJourney|TestP23ReviewCycleChargeSurvivesLinkedRecoveryAndRebuild|TestP22ConformanceCorpusHasEveryAcceptedOperationAndFailure)$' -v ./cmd/fid/integration
```

### named-inventory-expansion

Working directory: `/home/mark/projects/fid`

Exact multiline invocation: [named-inventory-expansion.json](named-inventory-expansion.json).

### named-issues

Working directory: `/home/mark/projects/fid`

```text
go test -race -count=1 -run '^(TestV1JudgmentProjectionIsAbsentAfterReplacement|TestLinkedToDecisionRequiresApprovedReferenceInIssueScope|TestAcceptanceCriteriaRemainOrderedWithoutRoleVerifications|TestGetReturnsCentralizedDetailShape|TestListFiltersAndOrdering|TestGetReturnsEnrichedConcernFields)$' -v ./internal/issues
```

### named-journal

Working directory: `/home/mark/projects/fid`

```text
go test -race -count=1 -run '^(TestProjectEventFenceRejectsMalformedPhase23EventsWithoutProjectionChange|TestConcernDeviationFenceConvergencePreservesEveryRemainingEvent|TestTask11ReviewRecordedAndFindingResolvedConditionShapes|TestTask11FindingResolvedConditionDispatchPreservesNullableBinding|TestTask11ConditionShapeRegistryIsClosedAtTenRows|TestConcernDisposedRequiresExactFieldsAndOpenGuard|TestTask8ValueShapesReuseJudgmentSubject|TestTask9ConditionDispatchRetainsPinnedPhase22Bytes|TestTask9ConditionDispatchPreservesNullableBindings|TestCurrentEventColumnDomainsMatchAcceptedDDL|TestCurrentReviewRecordedRegistersItsDistinctJudgmentValues|TestReviewRecordedAloneAdmitsRepeatedFindingDispositionEffects|TestProjectEventPopulatesPhaseTwoRepresentativeTablesViaJournalEffects|TestTask10HumanOutcomeConditionDispatchPreservesNullableBinding|TestTask10HumanOutcomeConditionShapeAndNonStringDiscriminant|TestTask9PrechangeBudgetObservedJournalReplaysWithoutMutation|TestTask9ProtectedDecisionStopDetailsAndConditionAreStrict|TestTask10HumanOutcomeStopDetailsAndConditionAreStrict|TestTask9DecisionConditionShapeAndNonStringDiscriminant|TestCurrentClosedValuesRejectIncompleteAndInvalidFacts|TestCurrentEventRejectsConditionalBindingCorruption|TestExecutionBlockedEventProjectsIdenticallyLiveAndReplay|TestActivationRecoveryEventsProjectIdenticallyLiveAndReplay)$' -v ./internal/journal
```

### named-links

Working directory: `/home/mark/projects/fid`

```text
go test -race -count=1 -run '^(TestRegistryHelpersCoverEndpointAndCompatibilityRules)$' -v ./internal/links
```

### named-maintenance

Working directory: `/home/mark/projects/fid`

```text
go test -race -count=1 -run '^(TestCheckErrorsReportsAllSeededStructuralViolationFamilies|TestCheckAcceptsPopulatedGenericReferencesWithoutStructuralErrors)$' -v ./internal/maintenance
```

### named-review

Working directory: `/home/mark/projects/fid`

```text
go test -race -count=1 -run '^(TestConcernDisposeCarriedForwardBindsExistingIssue|TestConcernSubjectAndDispositionMatrix|TestConcernSelfDisposalKeepsRaisingAndDisposingActors|TestDeviationAcceptsExplicitPredecessorReference|TestChangesRequiredReviewChargesExactlyOneReviewCycle|TestTerminalFindingDispositionAuthorityMatrix|TestValidateTerminalFindingInputBoundaries|TestValidateVerdictFindingAgreement|TestValidateCoverageRejectsSourceInItemAndExclusion|TestCheckIndependenceRejectsClaimAndRunHistoryDirections|TestRequirementRecordReferenceSupersedesPriorActiveRow|TestDecisionReferenceSourceMustExistAndMatchDigest)$' -v ./internal/review
```

### named-workflow

Working directory: `/home/mark/projects/fid`

```text
go test -race -count=1 -run '^(TestReviewClaimRejectsDirtyPathAdoption|TestEnvelopeCreateAcceptsTypedObligationSet|TestEnvelopeCreateRejectsInvalidTypedObligationAndProtectedDecisionAuthoring|TestEnvelopeCreatePersistsEveryTypedObligationAndProtectedDecision|TestEvaluateEnvelopeObligationsReturnsIndependentOrderedResults|TestEnvelopeActivationEvaluationIncludesEveryObligation|TestEnvelopeAuthoringObligationCardinalityAndTypeFieldMatrix|TestEnvelopeAuthoringAcceptedAndForbiddenTypedFieldCombinations|TestEnvelopeAuthoringReviewCycleBudgetMatrix|TestEnvelopeAuthoringReviewCycleAlongsideRetainedBudgets|TestEnvelopeAuthoringProtectedDecisionMatrix|TestEvaluateEnvelopeObligationsIncludesQualifyingValidFreshAvailableState|TestStopResolutionChecksCurrentActivationAgreementAndCapability|TestHumanOutcomeConditionQualifyingFactRules|TestDecisionSatisfiedConditionQualifyingFactRules|TestGuardedActionsEvaluateOnlyObligationsThatGuardThem|TestCompleteEnvelopeDeniedByGuardingObligationCommitsNothing|TestRecoveryFixturesKeepRepositoryAndProjectionIsolation|TestTask9RetainedAttentionKeyMatchesPrechangeReference|TestHumanOutcomeSupersessionFollowsRecordedEventOrderBothDirections|TestHumanOutcomeEvaluationUsesLaterDifferingResult|TestAbsentHumanOutcomeFreshnessRemainsFresh|TestRegistryEqualsV01Definitions|TestRetiredJudgmentRulesAreNotImplemented|TestRecoverRunRejectsApplicableStopWithoutEffects|TestStopAndAttentionListsFilterLimitAndRemainReadOnly|TestProtectedDecisionBlocksOnlyDeclaredAffectedActions|TestTerminalSuccessEvaluationIncludesEnvelopeStop|TestEquivalentStopsShareAttentionButKeepDistinctFacts)$' -v ./internal/workflow
```

### preservation-rows

Working directory: `/home/mark/projects/fid`

```text
python3 project-docs/docs/phases/phase-23/qa/final-evidence/preservation-rows.py
```

### protected-file-hashes

Working directory: `/home/mark/projects/fid`

```text
sha256sum project-docs/docs/current/tl_escalations_current.md project-docs/docs/phases/phase-23/run-ledger.md project-docs/docs/phases/phase-23/triage-log.md project-docs/docs/phases/phase-23/qa/qa-plan.md project-docs/docs/phases/phase-23/qa/validation-contract.md project-docs/docs/phases/phase-23/qa/human-checklist.md
```

### rebuild-pragmas

Working directory: `/home/mark/projects/fid`

```text
python3 project-docs/docs/phases/phase-23/qa/final-evidence/rebuild-pragmas.py
```

### reference-binary-identity

Working directory: `/home/mark/projects/fid`

```text
go version -m /tmp/fid-p23-qa-reference-PgIz9Jst/fid-task11
```

### reference-comparison

Working directory: `/home/mark/projects/fid`

```text
python3 project-docs/docs/phases/phase-23/qa/final-evidence/compare-reference.py
```

### reference-current-probe

Working directory: `/home/mark/projects/fid`

```text
go test -race -count=1 -overlay /home/mark/projects/fid/project-docs/docs/phases/phase-23/qa/final-evidence/current-overlay.json -run '^TestQAPhase23ReferenceAttentionKeyAndBytes$' -v ./internal/workflow
```

### reference-prechange-probe

Working directory: `/tmp/fid-p23-qa-reference-PgIz9Jst/prechange`

```text
python3 /home/mark/projects/fid/project-docs/docs/phases/phase-23/qa/final-evidence/reference-environment.py prechange go test -race -count=1 -overlay /home/mark/projects/fid/project-docs/docs/phases/phase-23/qa/final-evidence/reference-overlay.json -run '^TestQAPhase23ReferenceAttentionKeyAndBytes$' -v ./internal/workflow
```

### reference-prechange-worktree

Working directory: `/home/mark/projects/fid`

```text
git worktree add --detach /tmp/fid-p23-qa-reference-PgIz9Jst/prechange 74c83b7033729785eb99403ecec9dcde30ac4280
```

### reference-task11-build-bound

Working directory: `/tmp/fid-p23-qa-reference-PgIz9Jst/task11`

```text
python3 /home/mark/projects/fid/project-docs/docs/phases/phase-23/qa/final-evidence/reference-environment.py task11 go build -o /tmp/fid-p23-qa-reference-PgIz9Jst/fid-task11 ./cmd/fid
```

### reference-task11-build-diagnostic

Working directory: `/tmp/fid-p23-qa-reference-PgIz9Jst/task11`

```text
go build -x -o /tmp/fid-p23-qa-reference-PgIz9Jst/fid-task11 ./cmd/fid
```

### reference-task11-build

Working directory: `/tmp/fid-p23-qa-reference-PgIz9Jst/task11`

```text
go build -o /tmp/fid-p23-qa-reference-PgIz9Jst/fid-task11 ./cmd/fid
```

### reference-task11-git-status

Working directory: `/tmp/fid-p23-qa-reference-PgIz9Jst/task11`

```text
git status --porcelain
```

### reference-task11-worktree

Working directory: `/home/mark/projects/fid`

```text
git worktree add --detach /tmp/fid-p23-qa-reference-PgIz9Jst/task11 3dff4574a71438900100449e36d7044969ba9162
```

### reference-worktree-setup

Working directory: `/home/mark/projects/fid`

Exact multiline invocation: [reference-worktree-setup.json](reference-worktree-setup.json).

### report-01

Working directory: `/home/mark/projects/fid`

```text
project-docs/tools/report-lint.sh project-docs/docs/phases/phase-23/tasks/task-1-done.md project-docs/docs/phases/phase-23/tasks/task-1.md
```

### report-02

Working directory: `/home/mark/projects/fid`

```text
project-docs/tools/report-lint.sh project-docs/docs/phases/phase-23/tasks/task-2-done.md project-docs/docs/phases/phase-23/tasks/task-2.md
```

### report-03

Working directory: `/home/mark/projects/fid`

```text
project-docs/tools/report-lint.sh project-docs/docs/phases/phase-23/tasks/task-3-done.md project-docs/docs/phases/phase-23/tasks/task-3.md
```

### report-04

Working directory: `/home/mark/projects/fid`

```text
project-docs/tools/report-lint.sh project-docs/docs/phases/phase-23/tasks/task-4-done.md project-docs/docs/phases/phase-23/tasks/task-4.md
```

### report-05

Working directory: `/home/mark/projects/fid`

```text
project-docs/tools/report-lint.sh project-docs/docs/phases/phase-23/tasks/task-5-done.md project-docs/docs/phases/phase-23/tasks/task-5.md
```

### report-06

Working directory: `/home/mark/projects/fid`

```text
project-docs/tools/report-lint.sh project-docs/docs/phases/phase-23/tasks/task-6-done.md project-docs/docs/phases/phase-23/tasks/task-6.md
```

### report-07

Working directory: `/home/mark/projects/fid`

```text
project-docs/tools/report-lint.sh project-docs/docs/phases/phase-23/tasks/task-7-done.md project-docs/docs/phases/phase-23/tasks/task-7.md
```

### report-08

Working directory: `/home/mark/projects/fid`

```text
project-docs/tools/report-lint.sh project-docs/docs/phases/phase-23/tasks/task-8-done.md project-docs/docs/phases/phase-23/tasks/task-8.md
```

### report-09

Working directory: `/home/mark/projects/fid`

```text
project-docs/tools/report-lint.sh project-docs/docs/phases/phase-23/tasks/task-9-done.md project-docs/docs/phases/phase-23/tasks/task-9.md
```

### report-10

Working directory: `/home/mark/projects/fid`

```text
project-docs/tools/report-lint.sh project-docs/docs/phases/phase-23/tasks/task-10-done.md project-docs/docs/phases/phase-23/tasks/task-10.md
```

### report-11

Working directory: `/home/mark/projects/fid`

```text
project-docs/tools/report-lint.sh project-docs/docs/phases/phase-23/tasks/task-11-done.md project-docs/docs/phases/phase-23/tasks/task-11.md
```

### report-12

Working directory: `/home/mark/projects/fid`

```text
project-docs/tools/report-lint.sh project-docs/docs/phases/phase-23/tasks/task-12-done.md project-docs/docs/phases/phase-23/tasks/task-12.md
```

### report-13

Working directory: `/home/mark/projects/fid`

```text
project-docs/tools/report-lint.sh project-docs/docs/phases/phase-23/tasks/task-13-done.md project-docs/docs/phases/phase-23/tasks/task-13.md
```

### report-14

Working directory: `/home/mark/projects/fid`

```text
project-docs/tools/report-lint.sh project-docs/docs/phases/phase-23/tasks/task-14-done.md project-docs/docs/phases/phase-23/tasks/task-14.md
```

### report-15

Working directory: `/home/mark/projects/fid`

```text
project-docs/tools/report-lint.sh project-docs/docs/phases/phase-23/tasks/task-15-done.md project-docs/docs/phases/phase-23/tasks/task-15.md
```

### report-16

Working directory: `/home/mark/projects/fid`

```text
project-docs/tools/report-lint.sh project-docs/docs/phases/phase-23/tasks/task-16-done.md project-docs/docs/phases/phase-23/tasks/task-16.md
```

### report-17

Working directory: `/home/mark/projects/fid`

```text
project-docs/tools/report-lint.sh project-docs/docs/phases/phase-23/tasks/task-17-done.md project-docs/docs/phases/phase-23/tasks/task-17.md
```

### report-18

Working directory: `/home/mark/projects/fid`

```text
project-docs/tools/report-lint.sh project-docs/docs/phases/phase-23/tasks/task-18-done.md project-docs/docs/phases/phase-23/tasks/task-18.md
```

### report-timing-audit

Working directory: `/home/mark/projects/fid`

Exact multiline invocation: [report-timing-audit.json](report-timing-audit.json).

### schema-column-contradiction

Working directory: `/home/mark/projects/fid`

Exact multiline invocation: [schema-column-contradiction.json](schema-column-contradiction.json).

### schema-sql

Working directory: `/home/mark/projects/fid`

```text
python3 project-docs/docs/phases/phase-23/qa/final-evidence/schema-sql.py
```

### security-source-audit

Working directory: `/home/mark/projects/fid`

Exact multiline invocation: [security-source-audit.json](security-source-audit.json).

### setup-budget-testcache

Working directory: `/home/mark/projects/fid`

```text
go clean -testcache
```

### setup-testcache

Working directory: `/home/mark/projects/fid`

```text
go clean -testcache
```

### staged-whitespace-after-repair

Working directory: `/home/mark/projects/fid`

```text
git diff --cached --check -- project-docs/docs/phases/phase-23/qa/final-review.md project-docs/docs/phases/phase-23/qa/human-checklist.md project-docs/docs/phases/phase-23/qa/final-evidence
```

### staged-whitespace

Working directory: `/home/mark/projects/fid`

```text
git diff --cached --check -- project-docs/docs/phases/phase-23/qa/final-review.md project-docs/docs/phases/phase-23/qa/human-checklist.md project-docs/docs/phases/phase-23/qa/final-evidence
```

### tmp-after-budget-integration

Working directory: `/home/mark/projects/fid`

Exact multiline invocation: [tmp-after-budget-integration.json](tmp-after-budget-integration.json).

### tmp-after-integration

Working directory: `/home/mark/projects/fid`

```text
python3 -c 'import pathlib,shutil; p=pathlib.Path("/tmp"); d=list(p.glob("fid-integration-binary-*")); print("fid-integration-binary directories:",len(d)); print("disk_usage /tmp:",shutil.disk_usage(p))'
```

### tmp-before-budget-integration

Working directory: `/home/mark/projects/fid`

Exact multiline invocation: [tmp-before-budget-integration.json](tmp-before-budget-integration.json).

### tmp-before-integration

Working directory: `/home/mark/projects/fid`

```text
python3 -c 'import pathlib,shutil; p=pathlib.Path("/tmp"); d=list(p.glob("fid-integration-binary-*")); print("fid-integration-binary directories:",len(d)); print("disk_usage /tmp:",shutil.disk_usage(p))'
```
