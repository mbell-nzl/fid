CREATE TABLE merge_observations (
 id TEXT PRIMARY KEY,
 task_id TEXT NOT NULL REFERENCES issues(id),
 handoff_id TEXT NOT NULL REFERENCES handoffs(id),
 merge_commit TEXT NOT NULL,
 record_json TEXT NOT NULL,
 UNIQUE (handoff_id, merge_commit)
);
CREATE INDEX merge_observations_by_task ON merge_observations(task_id);
CREATE TABLE ci_results (
 id TEXT PRIMARY KEY,
 task_id TEXT NOT NULL REFERENCES issues(id),
 merge_id TEXT NOT NULL REFERENCES merge_observations(id),
 repository TEXT NOT NULL,
 run_id INTEGER NOT NULL CHECK (run_id > 0),
 attempt INTEGER NOT NULL CHECK (attempt > 0),
 record_json TEXT NOT NULL,
 UNIQUE (repository, run_id, attempt, merge_id)
);
CREATE INDEX ci_results_by_merge ON ci_results(merge_id);
