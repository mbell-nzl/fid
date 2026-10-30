package main

import "github.com/mbell-nzl/fid/internal/journal"

func idempotencyConflictDetails(conflict *journal.IdempotencyConflictError) map[string]any {
	return map[string]any{
		"operation_scope":         conflict.OperationScope,
		"idempotency_key":         conflict.IdempotencyKey,
		"stored_request_digest":   conflict.StoredRequestDigest,
		"supplied_request_digest": conflict.SuppliedRequestDigest,
		"differing_fields":        journal.DifferingRequestFields(conflict.StoredFieldDigests, conflict.SuppliedFieldDigests),
	}
}
