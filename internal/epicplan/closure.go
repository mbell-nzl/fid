package epicplan

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"slices"
)

// CloseBinding preserves the evidence counted by the closing command, including
// ancestry it proved after a pass was recorded. Replay needs no Git objects.
type CloseBinding struct {
	ApprovalID string   `json:"approval_id"`
	Results    []string `json:"results"`
	Merges     []string `json:"merges"`
}

func readCloseBinding(q Queryer, epic string) (CloseBinding, error) {
	var b CloseBinding
	var raw string
	err := q.QueryRow(`SELECT binding_json FROM person_decisions WHERE kind='close-epic' AND subject_id=? AND binding_json IS NOT NULL ORDER BY rowid DESC LIMIT 1`, epic).Scan(&raw)
	if err == sql.ErrNoRows {
		return b, nil
	}
	if err != nil {
		return b, err
	}
	err = json.Unmarshal([]byte(raw), &b)
	return b, err
}

func coversMerges(p Pass, merges []string, closed CloseBinding, roots []string) bool {
	for _, merge := range merges {
		if slices.Contains(p.Merges, merge) || merge == p.Commit || slices.Contains(closed.Results, p.ID) && slices.Contains(closed.Merges, merge) {
			continue
		}
		if len(roots) == 0 || roots[0] == "" {
			return false
		}
		cmd := exec.Command("git", "-C", roots[0], "merge-base", "--is-ancestor", merge, p.Commit)
		cmd.Env = append(os.Environ(), "GIT_NO_LAZY_FETCH=1")
		if err := cmd.Run(); err != nil {
			return false
		}
	}
	return true
}

// CloseEvidence captures immutable result identities after the caller has
// established their currency, using local Git ancestry where necessary.
func CloseEvidence(q Queryer, epic string) (CloseBinding, error) {
	a, err := Current(q, epic)
	if err != nil {
		return CloseBinding{}, err
	}
	merges, err := Merges(q, epic)
	if err != nil {
		return CloseBinding{}, err
	}
	views, err := Views(q, epic)
	if err != nil {
		return CloseBinding{}, err
	}
	b := CloseBinding{ApprovalID: a.ID, Results: []string{}, Merges: merges}
	for _, v := range views {
		if v.Pass == nil {
			return b, fmt.Errorf("epic close lacks %s result", v.ID)
		}
		b.Results = append(b.Results, v.Pass.ID)
	}
	return b, ValidateCloseEvidence(q, epic, b)
}

// ValidateCloseEvidence checks the recorded ancestry attestation against the
// journal state. External Git ancestry is proved by the owning command.
func ValidateCloseEvidence(q Queryer, epic string, b CloseBinding) error {
	a, err := Current(q, epic)
	if err != nil {
		return err
	}
	merges, err := Merges(q, epic)
	if err != nil {
		return err
	}
	views, err := Views(q, epic)
	if err != nil {
		return err
	}
	if a.State != "current" || b.ApprovalID != a.ID || !slices.Equal(b.Merges, merges) || len(b.Results) != len(views) {
		return fmt.Errorf("epic close binding differs from current plan or merges")
	}
	for i, v := range views {
		if v.Pass == nil || v.Pass.ID != b.Results[i] || v.Pass.Generation != v.Generation || v.Pass.Outcome != "passed" {
			return fmt.Errorf("epic close requires a current passing result for %s", v.ID)
		}
	}
	return nil
}
