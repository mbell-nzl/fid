INSERT INTO workflows (id, name, applies_to, active) VALUES
  ('escalation-default', 'Escalation default', 'escalation', 1);

INSERT INTO workflow_states (workflow_id, state, category, ordinal) VALUES
  ('escalation-default', 'open', 'started', 1),
  ('escalation-default', 'resolved', 'completed', 2);

INSERT INTO workflow_transitions (workflow_id, from_state, to_state) VALUES
  ('escalation-default', 'open', 'resolved');

INSERT INTO workflow_transition_requires (workflow_id, from_state, to_state, rule, params_json) VALUES
  ('escalation-default', 'open', 'resolved', 'linked_to_decision', NULL);
