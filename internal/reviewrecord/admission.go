package reviewrecord

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"slices"
	"strings"

	"github.com/mbell-nzl/fid/internal/handoffrecord"
	"github.com/mbell-nzl/fid/internal/pathglob"
	"github.com/mbell-nzl/fid/internal/policy"
)

type Queryer interface {
	Query(string, ...any) (*sql.Rows, error)
	QueryRow(string, ...any) *sql.Row
}

// RequiredProfiles uses only the hand-off's frozen lane, paths and policy.
// Unsupported paths conservatively require test adequacy, as check selection does.
func RequiredProfiles(h handoffrecord.Record, p policy.Effective) ([]string, error) {
	if !p.TestAdequacySet {
		return nil, fmt.Errorf("hand-off has no configured test-adequacy paths")
	}
	profiles := []string{}
	if h.Lane == "standard" || h.Lane == "assured" {
		profiles = append(profiles, "code")
	}
	required := false
	for _, pattern := range p.TestAdequacy {
		g, err := pathglob.Compile(pattern)
		if err != nil {
			return nil, err
		}
		for _, path := range h.Paths() {
			required = required || !pathglob.ValidPath(path) || g.Match(path)
		}
	}
	if required {
		profiles = append(profiles, "test-adequacy")
	}
	return profiles, nil
}

func LatestHandoff(conn Queryer, task string) (handoffrecord.Record, error) {
	var raw string
	if err := conn.QueryRow(`SELECT record_json FROM handoffs WHERE task_id=? ORDER BY rowid DESC LIMIT 1`, task).Scan(&raw); err != nil {
		return handoffrecord.Record{}, fmt.Errorf("review needs a hand-off: %w", err)
	}
	var h handoffrecord.Record
	if err := json.Unmarshal([]byte(raw), &h); err != nil {
		return h, err
	}
	return h, nil
}

// Admit is shared by live writes, replay and import. It never recomputes a
// fingerprint: replay checks its form and its recorded hand-off endpoints.
func Admit(conn Queryer, r Round) error {
	if err := r.Validate(); err != nil {
		return err
	}
	h, err := LatestHandoff(conn, r.TaskID)
	if err != nil {
		return err
	}
	if h.ID != r.HandoffID || h.Base != r.Binding.Base || h.Head != r.Binding.Head || r.At < h.At {
		return fmt.Errorf("review must bind to the latest hand-off's recorded base and head")
	}
	var count int
	if err := conn.QueryRow(`SELECT count(*) FROM holds WHERE task_id=? AND actor=?`, r.TaskID, r.Reviewer.Actor).Scan(&count); err != nil {
		return err
	}
	if count != 0 {
		return fmt.Errorf("a builder cannot review its own task")
	}
	var raw string
	if err := conn.QueryRow(`SELECT effective_json FROM policy_snapshots WHERE sequence=?`, h.PolicySnapshot).Scan(&raw); err != nil {
		return err
	}
	var p policy.Effective
	if err := json.Unmarshal([]byte(raw), &p); err != nil {
		return err
	}
	required, err := RequiredProfiles(h, p)
	if err != nil {
		return err
	}
	for _, profile := range required {
		if !slices.Contains(r.Profiles, profile) {
			return fmt.Errorf("review is missing required profile %s", profile)
		}
	}
	if r.Stream == "separate" && h.Lane != "assured" {
		return fmt.Errorf("separate review requires an Assured hand-off")
	}
	rounds, err := Read(conn, r.TaskID)
	if err != nil {
		return err
	}
	major := false
	for _, prior := range rounds {
		if r.Stream != prior.Stream && r.Reviewer.SessionID == prior.Reviewer.SessionID {
			return fmt.Errorf("separate review must use a different session from every primary round")
		}
	}
	if err := AdmitDispositions(conn, r); err != nil {
		return err
	}
	open, err := OpenMajors(conn, r.TaskID)
	if err != nil {
		return err
	}
	for _, f := range open {
		closed := false
		for _, d := range r.Dispositions {
			closed = closed || d.FindingID == f.ID && d.Closes()
		}
		major = major || !closed
	}
	for _, f := range r.Findings {
		major = major || f.Severity == "major"
	}
	if (r.Verdict == "CHANGES") != major {
		return fmt.Errorf("review verdict disagrees with open major findings")
	}
	return nil
}

func Read(conn Queryer, task string) ([]Round, error) {
	rows, err := conn.Query(`SELECT record_json FROM review_rounds WHERE task_id=? ORDER BY rowid`, task)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	rounds := []Round{}
	for rows.Next() {
		var raw string
		var r Round
		if err := rows.Scan(&raw); err != nil {
			return nil, err
		}
		if err := json.Unmarshal([]byte(raw), &r); err != nil {
			return nil, err
		}
		rounds = append(rounds, r)
	}
	return rounds, rows.Err()
}

// OtherModel compares with every build hold, including ended ones. A historical
// agent without model attribution cannot establish a different model family;
// a person without a model built by hand.
func OtherModel(conn Queryer, task, model string) (bool, error) {
	if model == "" {
		return false, nil
	}
	family, _, ok := strings.Cut(model, "/")
	if !ok {
		return false, nil
	}
	rows, err := conn.Query(`SELECT actor,model FROM holds WHERE task_id=?`, task)
	if err != nil {
		return false, err
	}
	defer func() { _ = rows.Close() }()
	for rows.Next() {
		var actor, built string
		if err := rows.Scan(&actor, &built); err != nil {
			return false, err
		}
		if built == "" && strings.HasPrefix(actor, "agent:") {
			return false, nil
		}
		if len(built) > len(family) && built[:len(family)+1] == family+"/" {
			return false, nil
		}
	}
	return true, rows.Err()
}
