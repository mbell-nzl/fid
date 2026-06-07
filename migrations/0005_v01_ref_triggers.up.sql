CREATE TRIGGER claims_scope_guard_insert
BEFORE INSERT ON claims
FOR EACH ROW
WHEN NEW.ref_type != 'issue'
BEGIN
  SELECT RAISE(ABORT, 'claims: ref_type must be issue in v0.1');
END;

CREATE TRIGGER links_endpoint_scope_guard_insert
BEFORE INSERT ON links
FOR EACH ROW
WHEN NEW.from_type NOT IN ('issue', 'requirement')
  OR NEW.to_type NOT IN ('issue', 'requirement')
BEGIN
  SELECT RAISE(ABORT, 'links: endpoint types must be issue or requirement in v0.1');
END;

CREATE TRIGGER links_from_ref_integrity_guard_insert
BEFORE INSERT ON links
FOR EACH ROW
WHEN (NEW.from_type = 'issue' AND NOT EXISTS (
    SELECT 1
    FROM issues
    WHERE id = NEW.from_id
  ))
  OR (NEW.from_type = 'requirement' AND NOT EXISTS (
    SELECT 1
    FROM requirements
    WHERE id = NEW.from_id
  ))
BEGIN
  SELECT RAISE(ABORT, 'links: from endpoint must resolve to an existing row');
END;

CREATE TRIGGER links_to_ref_integrity_guard_insert
BEFORE INSERT ON links
FOR EACH ROW
WHEN (NEW.to_type = 'issue' AND NOT EXISTS (
    SELECT 1
    FROM issues
    WHERE id = NEW.to_id
  ))
  OR (NEW.to_type = 'requirement' AND NOT EXISTS (
    SELECT 1
    FROM requirements
    WHERE id = NEW.to_id
  ))
BEGIN
  SELECT RAISE(ABORT, 'links: to endpoint must resolve to an existing row');
END;

CREATE TRIGGER links_relation_compatibility_guard_insert
BEFORE INSERT ON links
FOR EACH ROW
WHEN NEW.from_type IN ('issue', 'requirement')
  AND NEW.to_type IN ('issue', 'requirement')
  AND NOT (
  NEW.relation = 'relates_to'
  OR (
    NEW.relation = 'implements'
    AND NEW.from_type = 'issue'
    AND NEW.to_type = 'requirement'
  )
  OR (
    NEW.relation = 'implemented_by'
    AND NEW.from_type = 'requirement'
    AND NEW.to_type = 'issue'
  )
  OR (
    NEW.relation IN (
      'blocks',
      'blocked_by',
      'duplicates',
      'duplicated_by',
      'supersedes',
      'superseded_by',
      'caused_by',
      'causes',
      'discovered_from',
      'discovered'
    )
    AND NEW.from_type = 'issue'
    AND NEW.to_type = 'issue'
  )
)
BEGIN
  SELECT RAISE(ABORT, 'links: relation/endpoint combination is not allowed in v0.1');
END;

CREATE TRIGGER claims_scope_guard_update
BEFORE UPDATE OF ref_type ON claims
FOR EACH ROW
WHEN NEW.ref_type != 'issue'
BEGIN
  SELECT RAISE(ABORT, 'claims: ref_type must be issue in v0.1');
END;

CREATE TRIGGER claims_ref_integrity_guard_insert
BEFORE INSERT ON claims
FOR EACH ROW
WHEN NEW.ref_type = 'issue'
  AND NOT EXISTS (
    SELECT 1
    FROM issues
    WHERE id = NEW.ref_id
  )
BEGIN
  SELECT RAISE(ABORT, 'claims: ref_type/ref_id must resolve to an existing row');
END;

CREATE TRIGGER activity_log_ref_integrity_guard_insert
BEFORE INSERT ON activity_log
FOR EACH ROW
WHEN NEW.ref_type NOT IN ('issue', 'requirement', 'decision', 'escalation')
  OR (NEW.ref_type = 'issue' AND NOT EXISTS (
    SELECT 1
    FROM issues
    WHERE id = NEW.ref_id
  ))
  OR (NEW.ref_type = 'requirement' AND NOT EXISTS (
    SELECT 1
    FROM requirements
    WHERE id = NEW.ref_id
  ))
  OR (NEW.ref_type = 'decision' AND NOT EXISTS (
    SELECT 1
    FROM decisions
    WHERE id = NEW.ref_id
  ))
  OR (NEW.ref_type = 'escalation' AND NOT EXISTS (
    SELECT 1
    FROM escalations
    WHERE id = NEW.ref_id
  ))
BEGIN
  SELECT RAISE(ABORT, 'activity_log: ref_type/ref_id must resolve to an existing row');
END;

CREATE TRIGGER claims_ref_integrity_guard_update
BEFORE UPDATE OF ref_type, ref_id ON claims
FOR EACH ROW
WHEN NEW.ref_type = 'issue'
  AND NOT EXISTS (
    SELECT 1
    FROM issues
    WHERE id = NEW.ref_id
  )
BEGIN
  SELECT RAISE(ABORT, 'claims: ref_type/ref_id must resolve to an existing row');
END;
