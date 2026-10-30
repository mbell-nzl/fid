ALTER TABLE idempotency_results ADD COLUMN request_field_digests_json TEXT NOT NULL CHECK (json_valid(request_field_digests_json));
