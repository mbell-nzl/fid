PRAGMA defer_foreign_keys=ON;
ALTER TABLE issues DROP COLUMN lane;

DELETE FROM workflow_transition_requires WHERE workflow_id IN ('task-default','epic-default');
DELETE FROM workflow_transitions WHERE workflow_id IN ('task-default','epic-default');
DELETE FROM workflow_states WHERE workflow_id IN ('task-default','epic-default');
INSERT INTO workflow_states (workflow_id, state, category, ordinal) VALUES
  ('task-default', 'todo', 'unstarted', 1),
  ('task-default', 'in_progress', 'started', 2),
  ('task-default', 'dev_complete', 'started', 3),
  ('task-default', 'qa_review', 'started', 4),
  ('task-default', 'qa_block', 'started', 5),
  ('task-default', 'failed_gate', 'started', 6),
  ('task-default', 'qa_pass', 'completed', 7),
  ('task-default', 'done', 'completed', 8),
  ('task-default', 'cancelled', 'cancelled', 9),
  ('epic-default', 'todo', 'unstarted', 1),
  ('epic-default', 'in_progress', 'started', 2),
  ('epic-default', 'wrapped', 'started', 3),
  ('epic-default', 'closed', 'completed', 4),
  ('epic-default', 'cancelled', 'cancelled', 5);
INSERT INTO workflow_transitions (workflow_id, from_state, to_state) VALUES
  ('task-default', 'todo', 'in_progress'),
  ('task-default', 'in_progress', 'dev_complete'),
  ('task-default', 'in_progress', 'failed_gate'),
  ('task-default', 'failed_gate', 'in_progress'),
  ('task-default', 'dev_complete', 'qa_review'),
  ('task-default', 'qa_review', 'qa_pass'),
  ('task-default', 'qa_review', 'qa_block'),
  ('task-default', 'qa_block', 'in_progress'),
  ('task-default', 'qa_pass', 'done'),
  ('task-default', 'todo', 'cancelled'),
  ('task-default', 'in_progress', 'cancelled'),
  ('task-default', 'dev_complete', 'cancelled'),
  ('task-default', 'qa_review', 'cancelled'),
  ('task-default', 'qa_block', 'cancelled'),
  ('task-default', 'failed_gate', 'cancelled'),
  ('task-default', 'qa_pass', 'cancelled'),
  ('task-default', 'done', 'cancelled'),
  ('epic-default', 'todo', 'in_progress'),
  ('epic-default', 'in_progress', 'wrapped'),
  ('epic-default', 'wrapped', 'closed'),
  ('epic-default', 'todo', 'cancelled'),
  ('epic-default', 'in_progress', 'cancelled'),
  ('epic-default', 'wrapped', 'cancelled'),
  ('epic-default', 'closed', 'cancelled');
INSERT INTO workflow_transition_requires (workflow_id,from_state,to_state,rule) VALUES ('epic-default','in_progress','wrapped','all_child_tasks_done_or_cancelled');
