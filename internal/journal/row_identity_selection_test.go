package journal

import (
	"encoding/json"
	"testing"

	"github.com/mbell-nzl/fid/internal/journal/testdata/rowcorpus"
)

type rowCorpus struct {
	Prefix              []Event
	Controls, Malformed []rowProbe
}

func readRowCorpus(t *testing.T) rowCorpus {
	t.Helper()
	raw, err := rowcorpus.JSON()
	if err != nil {
		t.Fatal(err)
	}
	var corpus rowCorpus
	if err := json.Unmarshal(raw, &corpus); err != nil {
		t.Fatal(err)
	}
	return corpus
}

type rowRefusalKey struct {
	schema        int
	table, reason string
}

func rowIdentityRepresentatives(probes []rowProbe) []rowProbe {
	seen := map[rowRefusalKey]bool{}
	var representatives []rowProbe
	for _, probe := range probes {
		key := rowRefusalKey{probe.Event.EventSchema, probe.Table, probe.Reason}
		if !seen[key] {
			seen[key] = true
			representatives = append(representatives, probe)
		}
	}
	return representatives
}

func TestRowIdentityRepresentatives(t *testing.T) {
	t.Parallel()
	corpus := readRowCorpus(t)
	first := map[rowRefusalKey]string{}
	names := map[string]bool{}
	for _, probe := range corpus.Malformed {
		if names[probe.Name] {
			t.Fatalf("duplicate probe name %q", probe.Name)
		}
		names[probe.Name] = true
		key := rowRefusalKey{probe.Event.EventSchema, probe.Table, probe.Reason}
		if _, found := first[key]; !found {
			first[key] = probe.Name
		}
	}
	representatives := rowIdentityRepresentatives(corpus.Malformed)
	counts := map[rowRefusalKey]int{}
	for _, probe := range representatives {
		key := rowRefusalKey{probe.Event.EventSchema, probe.Table, probe.Reason}
		counts[key]++
		if probe.Name != first[key] {
			t.Errorf("representative %q, want first carrier %q", probe.Name, first[key])
		}
	}
	for key := range first {
		if counts[key] != 1 {
			t.Errorf("representative %+v count = %d, want 1", key, counts[key])
		}
	}
	if len(representatives) != len(first) {
		t.Fatalf("representatives = %d, want %d distinct schema/table/reason carriers", len(representatives), len(first))
	}
	t.Logf("%d controls, %d malformed probes, %d malformed path representatives", len(corpus.Controls), len(corpus.Malformed), len(representatives))
}
