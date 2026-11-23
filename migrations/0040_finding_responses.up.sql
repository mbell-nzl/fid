CREATE TABLE finding_responses (
 id TEXT PRIMARY KEY,
 finding_id TEXT NOT NULL REFERENCES round_findings(id),
 deferred_task_id TEXT UNIQUE REFERENCES issues(id),
 record_json TEXT NOT NULL
);
CREATE UNIQUE INDEX finding_one_deferral ON finding_responses(finding_id) WHERE deferred_task_id IS NOT NULL;
