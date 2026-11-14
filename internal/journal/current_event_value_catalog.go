package journal

// The retained Artifact verifier result has the same nine immutable metadata fields.
var currentValueShapes = map[string]string{
	"Artifact": "id:ID digest:Digest media_type:Nonempty size_bytes:Nonnegative retention_class:=evidence availability:=available,missing location_policy:=content-addressed created_at:Timestamp expires_at:?Timestamp",
}
