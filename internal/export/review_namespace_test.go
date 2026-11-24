package export

import (
	"encoding/json"
	"testing"

	"github.com/mbell-nzl/fid/internal/reviewrecord"
)

func TestReviewDeliveryNamespaces(t *testing.T) {
	t.Parallel()
	for _, kind := range []string{"round", "response"} {
		t.Run(kind, func(t *testing.T) {
			t.Parallel()
			v, e := responseArchiveFixture(t, false)
			if kind == "round" {
				v, e = reviewArchiveFixture(t)
				var r reviewrecord.Round
				if err := json.Unmarshal([]byte(e.Effects[0].Values["record_json"].(string)), &r); err != nil {
					t.Fatal(err)
				}
				r.ID, e.Aggregate.ID = "FRS-wrong", "FRS-wrong"
				setReviewArchiveRound(t, &e, r)
			} else {
				e.Aggregate.ID, e.Effects[0].PK["id"] = "RND-wrong", "RND-wrong"
				changeResponse(t, &e, func(r *reviewrecord.Response) { r.ID = "RND-wrong" })
			}
			if err := applyLeanEmptyEvent(v, "replay", e); err == nil {
				t.Fatal("admitted a delivery record in another record's ID namespace")
			}
		})
	}
}
