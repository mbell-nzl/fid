CREATE INDEX idx_concerns_issue_status ON concerns(issue_id, status);

CREATE INDEX idx_triage_issue_created ON triage_events(issue_id, created_at);

CREATE VIEW issues_enriched AS
WITH open_concerns AS (
  SELECT
    i.id AS issue_id,
    (
      SELECT COUNT(*)
      FROM concerns c
      WHERE c.issue_id = i.id AND c.status = 'open'
    ) AS n
  FROM issues i
),
latest_triage AS (
  SELECT
    i.id AS issue_id,
    (
      SELECT te.verdict
      FROM triage_events te
      WHERE te.issue_id = i.id
      ORDER BY te.created_at DESC
      LIMIT 1
    ) AS verdict
  FROM issues i
)
SELECT
  i.*,
  (COALESCE(oc.n, 0) > 0) AS concerns_raised,
  COALESCE(oc.n, 0) AS open_concerns_count,
  CASE lt.verdict
    WHEN 'RESUME' THEN 'resumed'
    WHEN 'BLOCK' THEN 'blocked'
    ELSE NULL
  END AS triage_status
FROM issues i
LEFT JOIN open_concerns oc ON oc.issue_id = i.id
LEFT JOIN latest_triage lt ON lt.issue_id = i.id;
