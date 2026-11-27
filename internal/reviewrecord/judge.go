package reviewrecord

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"slices"

	"github.com/mbell-nzl/fid/internal/fingerprint"
	"github.com/mbell-nzl/fid/internal/handoffrecord"
	"github.com/mbell-nzl/fid/internal/policy"
	"github.com/mbell-nzl/fid/internal/waiverrecord"
)

// PersonCall is the first approval, observed merge or hold for this hand-off.
// Its order comes from journal projection, never timestamps across tables.
type PersonCall struct {
	HandoffID string `json:"handoff_id"`
	Verdict   string `json:"verdict"`
	Actor     string `json:"actor"`
	Reason    string `json:"reason"`
	SourceID  string `json:"source_id"`
	At        string `json:"at"`
}

func FirstCall(conn Queryer, handoff string) (*PersonCall, error) {
	var installed int
	if err := conn.QueryRow(`SELECT count(*) FROM sqlite_master WHERE name='merge_calls'`).Scan(&installed); err != nil || installed == 0 {
		return nil, err
	}
	rows, err := conn.Query(`SELECT handoff_id,verdict,actor,reason,source_id,at FROM merge_calls WHERE handoff_id=?`, handoff)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	if !rows.Next() {
		return nil, rows.Err()
	}
	var c PersonCall
	err = rows.Scan(&c.HandoffID, &c.Verdict, &c.Actor, &c.Reason, &c.SourceID, &c.At)
	return &c, err
}

// JudgeDue shares the frozen profile and fingerprint rules with ordinary
// review admission. A shadow round adds no review gate or finding.
func JudgeDue(conn Queryer, h handoffrecord.Record, fp Binding) (bool, error) {
	var state, raw string
	if err := conn.QueryRow(`SELECT status FROM issues WHERE id=?`, h.TaskID).Scan(&state); err != nil {
		return false, err
	}
	if state != "in_review" {
		return false, nil
	}
	if err := conn.QueryRow(`SELECT effective_json FROM policy_snapshots WHERE sequence=?`, h.PolicySnapshot).Scan(&raw); err != nil {
		return false, err
	}
	var p policy.Effective
	if err := json.Unmarshal([]byte(raw), &p); err != nil {
		return false, err
	}
	w, err := waiverrecord.Current(conn, h.TaskID, fingerprint.Result{Algorithm: fp.Algorithm, Fingerprint: fp.Fingerprint})
	if err != nil {
		return false, err
	}
	if w.State == "current" {
		h.Lane = w.Decision.Binding.Lane
	}
	profiles, err := RequiredProfiles(h, p)
	if err != nil {
		return false, err
	}
	rounds, err := Read(conn, h.TaskID)
	if err != nil {
		return false, err
	}
	satisfied := map[string]bool{}
	for _, r := range rounds {
		if r.Stream == "judge" && r.HandoffID == h.ID {
			return false, nil
		}
		if r.Stream == "judge" || r.Binding.Algorithm != fp.Algorithm || r.Binding.Fingerprint != fp.Fingerprint {
			continue
		}
		satisfied[r.Stream] = true
		for _, profile := range profiles {
			satisfied[r.Stream] = satisfied[r.Stream] && slices.Contains(r.Profiles, profile)
		}
	}
	majors, err := OpenMajors(conn, h.TaskID)
	if err != nil {
		return false, err
	}
	return len(majors) == 0 && (len(profiles) == 0 || satisfied["primary"]) && (h.Lane != "assured" || satisfied["separate"]), nil
}

func JudgeIndependent(conn Queryer, task, actor, session string) (bool, error) {
	built, err := BuiltBy(conn, task, actor)
	if err != nil || built {
		return false, err
	}
	rounds, err := Read(conn, task)
	if err != nil {
		return false, err
	}
	for _, r := range rounds {
		if r.Stream != "judge" && r.Reviewer.Actor == actor || session != "" && r.Reviewer.SessionID == session {
			return false, nil
		}
	}
	return true, nil
}

func admitJudge(conn Queryer, r Round, h handoffrecord.Record, _ policy.Effective) error {
	independent, err := JudgeIndependent(conn, r.TaskID, r.Reviewer.Actor, r.Reviewer.SessionID)
	if err != nil {
		return err
	}
	if !independent {
		return fmt.Errorf("judge must be independent of every builder and reviewer")
	}
	due, err := JudgeDue(conn, h, r.Binding)
	if err != nil {
		return err
	}
	if !due {
		return fmt.Errorf("judge requires satisfied reviews and no judge round for this hand-off")
	}
	if r.LaneRaise != nil {
		rank := map[string]int{"routine": 0, "standard": 1, "assured": 2}
		if rank[r.LaneRaise.Lane] <= rank[h.Lane] {
			return fmt.Errorf("judge may only recommend a lane raise")
		}
	}
	return nil
}

func VisibleRound(conn Queryer, r Round) (RoundView, error) {
	v := RoundView{Round: r, Findings: []FindingView{}}
	if r.Stream != "judge" {
		return v, nil
	}
	call, err := FirstCall(conn, r.HandoffID)
	if err != nil {
		return v, err
	}
	v.PersonCall = call
	if call == nil {
		v.Sealed = true
		v.Round = Round{ID: r.ID, TaskID: r.TaskID, HandoffID: r.HandoffID, Stream: "judge", Profiles: []string{"merge"}, Findings: []Finding{}}
	} else {
		agree := call.Verdict == r.Verdict
		v.Agrees = &agree
	}
	return v, nil
}

func (v RoundView) JudgeText() string {
	if v.Sealed {
		return "Judge call recorded and sealed"
	}
	comparison := "disagree"
	if v.Agrees != nil && *v.Agrees {
		comparison = "agree"
	}
	message := fmt.Sprintf("Judge %s: %s; %s: %s (%s); calls %s", v.Verdict, v.Reason, v.PersonCall.Actor, v.PersonCall.Verdict, v.PersonCall.Reason, comparison)
	if v.LaneRaise != nil {
		message += "; recommended lane " + v.LaneRaise.Lane + " (shadow only): " + v.LaneRaise.Reason
	}
	return message
}

// JudgeViews binds shadow currency to the hand-off, without needing git objects.
// Ordinary review currency still follows the change fingerprint.
func JudgeViews(conn Queryer, task string) ([]RoundView, error) {
	rounds, err := Read(conn, task)
	if err != nil {
		return nil, err
	}
	views := []RoundView{}
	for _, r := range rounds {
		if r.Stream != "judge" {
			continue
		}
		v, err := VisibleRound(conn, r)
		if err != nil {
			return nil, err
		}
		h, err := LatestHandoff(conn, task)
		if err != nil {
			return nil, err
		}
		v.Current = r.HandoffID == h.ID
		views = append(views, v)
	}
	return views, nil
}

func RoundSealed(conn Queryer, id string) (bool, error) {
	var raw string
	if err := conn.QueryRow(`SELECT record_json FROM review_rounds WHERE id=?`, id).Scan(&raw); err != nil {
		return false, err
	}
	var r Round
	if err := json.Unmarshal([]byte(raw), &r); err != nil {
		return false, err
	}
	if r.Stream != "judge" {
		return false, nil
	}
	c, err := FirstCall(conn, r.HandoffID)
	return c == nil, err
}

// LatestHold drives waiting display independently of the first-call comparison.
func LatestHold(conn Queryer, handoff string) (*PersonCall, error) {
	var c PersonCall
	err := conn.QueryRow(`SELECT handoff_id,actor,reason,id,at FROM merge_holds WHERE handoff_id=? ORDER BY rowid DESC LIMIT 1`, handoff).Scan(&c.HandoffID, &c.Actor, &c.Reason, &c.SourceID, &c.At)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	c.Verdict = "HOLD"
	return &c, err
}
