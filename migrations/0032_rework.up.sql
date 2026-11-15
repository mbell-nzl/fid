-- The build hold still ended at hand-off; rework changes its kept reservation.
ALTER TABLE holds ADD COLUMN rework_at TEXT
 CHECK (rework_at IS NULL OR (end_reason IS 'review' AND rework_at >= ended_at));
