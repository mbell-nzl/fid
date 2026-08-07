CREATE TABLE workflow_assignments (
  scope       TEXT NOT NULL,
  kind        TEXT,
  workflow_id TEXT NOT NULL REFERENCES workflows(id),
  assigned_at DATETIME,
  CHECK (scope IN ('kind', 'project')),
  CHECK ((scope = 'kind' AND kind IS NOT NULL)
      OR (scope = 'project' AND kind IS NULL))
);

CREATE UNIQUE INDEX idx_workflow_assignments_kind
  ON workflow_assignments(kind) WHERE scope = 'kind';

CREATE UNIQUE INDEX idx_workflow_assignments_project
  ON workflow_assignments(scope) WHERE scope = 'project';
