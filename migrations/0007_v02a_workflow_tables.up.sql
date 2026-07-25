CREATE TABLE workflows (
  id          TEXT PRIMARY KEY,
  name        TEXT NOT NULL,
  applies_to  TEXT,
  active      BOOLEAN DEFAULT 1,
  created_at  DATETIME
);

CREATE TABLE workflow_states (
  workflow_id  TEXT NOT NULL REFERENCES workflows(id),
  state        TEXT NOT NULL,
  category     TEXT NOT NULL,
  ordinal      INTEGER,
  PRIMARY KEY (workflow_id, state)
);

CREATE TABLE workflow_transitions (
  workflow_id  TEXT NOT NULL REFERENCES workflows(id),
  from_state   TEXT NOT NULL,
  to_state     TEXT NOT NULL,
  PRIMARY KEY (workflow_id, from_state, to_state)
);

CREATE TABLE workflow_transition_requires (
  workflow_id  TEXT NOT NULL,
  from_state   TEXT NOT NULL,
  to_state     TEXT NOT NULL,
  rule         TEXT NOT NULL,
  params_json  TEXT,
  FOREIGN KEY (workflow_id, from_state, to_state)
    REFERENCES workflow_transitions(workflow_id, from_state, to_state)
);

INSERT INTO workflows (id, name, applies_to, active) VALUES
  ('epic-default', 'Epic default', 'epic', 1),
  ('task-default', 'Task default', 'task', 1),
  ('subtask-default', 'Subtask default', 'subtask', 1);

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
  ('epic-default', 'cancelled', 'cancelled', 5),
  ('subtask-default', 'todo', 'unstarted', 1),
  ('subtask-default', 'in_progress', 'started', 2),
  ('subtask-default', 'done', 'completed', 3),
  ('subtask-default', 'cancelled', 'cancelled', 4);

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
  ('epic-default', 'closed', 'cancelled'),
  ('subtask-default', 'todo', 'in_progress'),
  ('subtask-default', 'in_progress', 'done'),
  ('subtask-default', 'todo', 'cancelled'),
  ('subtask-default', 'in_progress', 'cancelled'),
  ('subtask-default', 'done', 'cancelled');

INSERT INTO workflow_transition_requires (workflow_id, from_state, to_state, rule, params_json) VALUES
  ('task-default', 'in_progress', 'dev_complete', 'all_acs_verified', '{"roles":["developer"]}'),
  ('task-default', 'in_progress', 'dev_complete', 'has_completion_report', NULL),
  ('task-default', 'qa_review', 'qa_pass', 'all_acs_verified', '{"roles":["qa_lead"]}'),
  ('epic-default', 'in_progress', 'wrapped', 'all_child_tasks_done_or_cancelled', NULL);
