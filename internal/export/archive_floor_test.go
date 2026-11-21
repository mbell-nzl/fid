package export

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/handoffrecord"
	"github.com/mbell-nzl/fid/internal/policy"
)

func TestHandoffFloorAdmissionEveryPath(t *testing.T) {
	t.Parallel()
	for _, mode := range []string{"replay", "import"} {
		for _, name := range []string{"raised", "below floor", "wrong floor", "missing floor", "wrong hard paths", "unnecessary raise", "lowering", "missing raise attribution", "missing lane update", "inconsistent path bytes", "invalid UTF-8"} {
			t.Run(mode+"/"+name, func(t *testing.T) {
				t.Parallel()
				v, e := handoffArchiveFixtureWithLane(t, "standard")
				var p policy.Effective
				var raw string
				if err := v.conn.QueryRow(`SELECT effective_json FROM policy_snapshots WHERE sequence=1`).Scan(&raw); err != nil {
					t.Fatal(err)
				}
				if err := json.Unmarshal([]byte(raw), &p); err != nil {
					t.Fatal(err)
				}
				// The fixture needs both lanes configured to isolate the lane-floor rule.
				l := p.Lanes["assured"]
				l.ChecksSet = true
				p.Lanes["assured"] = l
				encoded, err := p.Canonical()
				if err != nil {
					t.Fatal(err)
				}
				if _, err := v.conn.Exec(`UPDATE policy_snapshots SET effective_json=?,digest=?`, string(encoded), p.Digest()); err != nil {
					t.Fatal(err)
				}
				var r handoffrecord.Record
				if err := json.Unmarshal([]byte(e.Effects[0].Values["record_json"].(string)), &r); err != nil {
					t.Fatal(err)
				}
				r.ChangedPaths = []string{".fid.toml"}
				floor, err := policy.ComputeFloor(p, r.ChangedPaths)
				if err != nil {
					t.Fatal(err)
				}
				r.Floor = &floor
				r.Lane = "assured"
				beforeLane := "standard"
				r.RaisedFrom = &beforeLane
				e.Effects[1].Values["lane"] = "assured"
				switch name {
				case "below floor":
					r.Lane = "standard"
					r.RaisedFrom = nil
					delete(e.Effects[1].Values, "lane")
				case "wrong floor":
					value := "standard"
					r.Floor.Lane = &value
				case "missing floor":
					r.Floor = nil
				case "wrong hard paths":
					r.Floor.HardPaths = []string{}
				case "unnecessary raise":
					r.ChangedPaths = []string{}
					empty, _ := policy.ComputeFloor(p, r.ChangedPaths)
					r.Floor = &empty
				case "missing raise attribution":
					r.RaisedFrom = nil
				case "missing lane update":
					delete(e.Effects[1].Values, "lane")
				case "inconsistent path bytes":
					r.ChangedPathBytes = policy.PathBytes([]string{"other.md"})
				case "invalid UTF-8":
					r.SetPaths([]string{"odd\xff.md"})
					exact, err := policy.ComputeFloor(p, r.Paths())
					if err != nil {
						t.Fatal(err)
					}
					r.Floor = &exact
				case "lowering":
					r.Lane = "routine"
					e.Effects[1].Values["lane"] = "routine"
				}
				e.Effects[0].Values["record_json"] = encodedLean(r)
				before := handoffProjectionSnapshot(t, v)
				err = applyLeanEmptyEvent(v, mode, e)
				if name == "raised" || name == "invalid UTF-8" {
					if err != nil {
						t.Fatalf("valid raise: %v", err)
					}
				} else {
					reason := "floor"
					if name == "missing lane update" {
						reason = "task change"
					}
					if name == "inconsistent path bytes" {
						reason = "path bytes"
					}
					if err == nil || !strings.Contains(err.Error(), reason) {
						t.Fatalf("forged %s = %v", name, err)
					}
					if handoffProjectionSnapshot(t, v) != before {
						t.Fatal("refusal changed projection")
					}
				}
			})
		}
	}
}

func TestHistoricalHandoffRetainsPreFloorPolicy(t *testing.T) {
	t.Parallel()
	for _, mode := range []string{"replay", "import"} {
		t.Run(mode, func(t *testing.T) {
			t.Parallel()
			v, e := handoffArchiveFixture(t)
			var p policy.Effective
			var raw string
			if err := v.conn.QueryRow(`SELECT effective_json FROM policy_snapshots`).Scan(&raw); err != nil {
				t.Fatal(err)
			}
			if err := json.Unmarshal([]byte(raw), &p); err != nil {
				t.Fatal(err)
			}
			p.RulesVersion = 2
			encoded, err := p.Canonical()
			if err != nil {
				t.Fatal(err)
			}
			if _, err := v.conn.Exec(`UPDATE policy_snapshots SET effective_json=?,digest=?`, string(encoded), p.Digest()); err != nil {
				t.Fatal(err)
			}
			var r handoffrecord.Record
			if err := json.Unmarshal([]byte(e.Effects[0].Values["record_json"].(string)), &r); err != nil {
				t.Fatal(err)
			}
			r.Floor = nil
			r.ChangedPaths = []string{".fid.toml"}
			e.Effects[0].Values["record_json"] = encodedLean(r)
			if err := applyLeanEmptyEvent(v, mode, e); err != nil {
				t.Fatalf("historical receipt refused: %v", err)
			}
		})
	}
}
