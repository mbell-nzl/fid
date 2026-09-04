package evidence

import "testing"

func TestSameGateRepositoryMatchesEqualOptionalValues(t *testing.T) {
	baseOID := "base-oid"
	dirtyDigest := "sha256:dirty"
	left := Repository{
		BaseOID:       &baseOID,
		CaptureSource: "git",
		Dirty:         true,
		DirtyDigest:   &dirtyDigest,
		HeadOID:       "head-oid",
		ProjectID:     "PRJ-test",
		TreeOID:       "tree-oid",
		WorktreeID:    "WT-test",
	}
	baseOIDCopy := string([]byte(baseOID))
	dirtyDigestCopy := string([]byte(dirtyDigest))
	right := left
	right.BaseOID = &baseOIDCopy
	right.DirtyDigest = &dirtyDigestCopy

	if !sameGateRepository(left, right) {
		t.Fatal("sameGateRepository rejected equal optional repository values")
	}
}
