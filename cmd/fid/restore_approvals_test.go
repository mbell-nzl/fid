package main

import (
	"testing"

	"github.com/mbell-nzl/fid/internal/approvalrecord"
	"github.com/mbell-nzl/fid/internal/fingerprint"
	"github.com/mbell-nzl/fid/internal/reviewrecord"
	"github.com/mbell-nzl/fid/internal/reviewwork"
	"github.com/mbell-nzl/fid/internal/waiverrecord"
)

func assertV8Approvals(t *testing.T, root string) {
	t.Helper()
	conn, err := reviewwork.OpenProjection(root)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = conn.Close() }()
	var task string
	if err := conn.QueryRow(`SELECT subject_id FROM person_decisions WHERE kind='approve' LIMIT 1`).Scan(&task); err != nil {
		t.Fatal(err)
	}
	approvals, err := approvalrecord.Read(conn, task)
	if err != nil || len(approvals) != 2 {
		t.Fatalf("V8 approvals=%v/%v", approvals, err)
	}
	if approvals[0].Binding.Fingerprint == approvals[1].Binding.Fingerprint {
		t.Fatal("V8 must retain a stale approval")
	}
	h, err := reviewrecord.LatestHandoff(conn, task)
	if err != nil {
		t.Fatal(err)
	}
	waivers, err := waiverrecord.Read(conn, task)
	if err != nil || len(waivers) != 1 {
		t.Fatalf("V8 waivers=%v/%v", waivers, err)
	}
	b := approvals[1].Binding
	waiver, err := waiverrecord.Current(conn, task, fingerprint.Result{Algorithm: b.Algorithm, Fingerprint: b.Fingerprint})
	if err != nil || waiver.State != "stale" || waiver.Decision.ID != waivers[0].ID {
		t.Fatalf("V8 stale waiver=%v/%v", waiver, err)
	}
	if b.HandoffID != h.ID || b.Head != h.Head || b.Base != h.Base {
		t.Fatal("V8 current approval must cover latest hand-off")
	}
	// Restore need not include git objects. These recorded exact endpoints already
	// have a computed binding; the evaluator supplies its own local computation.
	current, err := approvalrecord.Current(conn, task, fingerprint.Result{Algorithm: b.Algorithm, Fingerprint: b.Fingerprint, Base: b.Base, Head: b.Head})
	if err != nil || current.State != "current" || current.Decision.ID != approvals[1].ID {
		t.Fatalf("V8 current approval=%v/%v", current, err)
	}
}
