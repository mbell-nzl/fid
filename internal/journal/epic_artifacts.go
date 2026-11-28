package journal

import (
	"database/sql"
	"fmt"

	"github.com/mbell-nzl/fid/internal/checkrecord"
)

// Every retained row belongs to a declared file in this verification's
// results. Existing content-addressed files may be shared by several checks.
func validateEpicFiles(tx *sql.Tx, event Event, result checkrecord.Result, definition checkrecord.Definition, used map[string]bool) error {
	declared := map[string]bool{}
	for _, path := range definition.Artifacts {
		if path == "" || declared[path] {
			return fmt.Errorf("invalid validation artifact declaration")
		}
		declared[path] = true
	}
	seen := map[string]bool{}
	for _, file := range result.Artifacts {
		if !declared[file.Path] || seen[file.Path] {
			return fmt.Errorf("validation artifact must name a declared path once")
		}
		seen[file.Path] = true
		var digest, availability string
		var size int64
		err := tx.QueryRow(`SELECT digest,size_bytes,availability FROM referenced_artifacts WHERE id=?`, file.ID).Scan(&digest, &size, &availability)
		if err != nil && err != sql.ErrNoRows {
			return err
		}
		for _, e := range event.Effects {
			if e.Table != "referenced_artifacts" || e.PK["id"] != file.ID {
				continue
			}
			if e.Op == "insert" {
				digest, _ = e.Values["digest"].(string)
				if !currentValuesEqual(e.Values["size_bytes"], file.SizeBytes) {
					return fmt.Errorf("validation artifact size differs from retained row")
				}
				size = file.SizeBytes
			}
			availability, _ = e.Values["availability"].(string)
		}
		if digest != file.Digest || size != file.SizeBytes || availability != "available" {
			return fmt.Errorf("validation artifact differs from retained evidence")
		}
		used[file.ID] = true
	}
	for _, file := range result.RefusedArtifacts {
		if !declared[file.Path] || seen[file.Path] || file.Reason == "" {
			return fmt.Errorf("refused validation artifact must name a declared path once with a reason")
		}
		seen[file.Path] = true
	}
	if len(seen) != len(declared) {
		return fmt.Errorf("validation must account for every declared artifact")
	}
	return nil
}
