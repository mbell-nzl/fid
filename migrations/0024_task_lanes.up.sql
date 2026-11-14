PRAGMA defer_foreign_keys=ON;
ALTER TABLE issues ADD COLUMN lane TEXT CHECK (lane IS NULL OR (lane IN ('routine', 'standard', 'assured') AND kind = 'task'));

-- The reset has no real ledger to convert. Retired states are removed;
-- old data using them is refused by the workflow-state foreign key.
DELETE FROM workflow_transition_requires WHERE workflow_id IN ('task-default','epic-default');
DELETE FROM workflow_transitions WHERE workflow_id IN ('task-default','epic-default');
DELETE FROM workflow_states WHERE workflow_id IN ('task-default','epic-default');
INSERT INTO workflow_states (workflow_id,state,category,ordinal) VALUES
 ('task-default','todo','unstarted',1),
 ('task-default','in_progress','started',2),
 ('task-default','in_review','started',3),
 ('task-default','done','completed',4),
 ('task-default','cancelled','cancelled',5),
 ('epic-default','todo','unstarted',1),
 ('epic-default','in_progress','started',2),
 ('epic-default','done','completed',3),
 ('epic-default','cancelled','cancelled',4);
INSERT INTO workflow_transitions (workflow_id,from_state,to_state) VALUES
 ('task-default','todo','in_progress'),
 ('task-default','in_progress','in_review'),
 ('task-default','in_review','in_progress'),
 ('task-default','in_review','done'),
 ('task-default','todo','cancelled'),
 ('task-default','in_progress','cancelled'),
 ('task-default','in_review','cancelled'),
 ('task-default','done','cancelled'),
 ('epic-default','todo','in_progress'),
 ('epic-default','in_progress','done'),
 ('epic-default','todo','cancelled'),
 ('epic-default','in_progress','cancelled'),
 ('epic-default','done','cancelled');
