DELETE FROM workflow_transition_requires
WHERE workflow_id = 'escalation-default';

DELETE FROM workflow_transitions
WHERE workflow_id = 'escalation-default';

DELETE FROM workflow_states
WHERE workflow_id = 'escalation-default';

DELETE FROM workflows
WHERE id = 'escalation-default';
