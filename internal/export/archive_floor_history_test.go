package export

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/handoffrecord"
	"github.com/mbell-nzl/fid/internal/policy"
)

func TestHistoricalHandoffStillRequiresTheRecordedLane(t *testing.T) {
	t.Parallel()
	for _, mode := range []string{"live", "replay", "import"} {
		for _, lane := range []string{"standard", "assured"} {
			t.Run(mode+"/"+lane, func(t *testing.T) {
				t.Parallel()
				v, event := handoffArchiveFixture(t)
				pathBytesPolicy(t, v, func(p *policy.Effective) {
					p.RulesVersion = 2
					p.Lanes["assured"] = p.Lanes["standard"]
				})
				var record handoffrecord.Record
				if err := json.Unmarshal([]byte(event.Effects[0].Values["record_json"].(string)), &record); err != nil {
					t.Fatal(err)
				}
				record.Lane, record.Floor, record.ChangedPaths = lane, nil, []string{".fid.toml"}
				event.Effects[0].Values["record_json"] = encodedLean(record)
				before := handoffProjectionSnapshot(t, v)
				err := admitPathBytesEvent(t, v, mode, event)
				if lane == "standard" {
					if err != nil {
						t.Fatalf("historical lane refused: %v", err)
					}
				} else if err == nil || !strings.Contains(err.Error(), "requires its started task and lane") || handoffProjectionSnapshot(t, v) != before {
					t.Fatalf("forged historical lane changed projection: %v", err)
				}
			})
		}
	}
}
