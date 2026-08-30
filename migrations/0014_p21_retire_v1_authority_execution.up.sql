DROP TABLE roles;
DROP TABLE claims;
DROP TABLE gate_runs;

DROP INDEX work_claims_one_active_mutation;
ALTER TABLE work_claims RENAME TO claims;
CREATE UNIQUE INDEX claims_one_active_mutation
  ON claims(issue_id) WHERE status = 'active';

DELETE FROM workflow_transition_requires
WHERE rule IN ('all_acs_verified', 'has_completion_report', 'concerns_addressed')
  AND workflow_id IN (
    SELECT id
    FROM workflows
    WHERE applies_to IN ('epic', 'task', 'subtask')
  );
