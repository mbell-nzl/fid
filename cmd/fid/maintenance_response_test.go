package main

import (
	"reflect"
	"testing"

	"github.com/mbell-nzl/fid/internal/maintenance"
)

func TestMaintenanceCheckWarningPromotionPreservesFindingsAndInfo(t *testing.T) {
	t.Parallel()
	result := maintenance.Result{
		Errors:   []maintenance.Finding{},
		Warnings: []maintenance.Finding{{Check: "test-warning", EntityID: "FID-test", Message: "Test warning"}},
		Info:     maintenance.Info{IssueCount: 3, RecentActivity: 2},
	}
	warn := maintenanceCheckResponse(result, false)
	promoted := maintenanceCheckResponse(result, true)
	if warn.Status != "warn" || promoted.Status != "fail" {
		t.Fatalf("statuses = %s, %s", warn.Status, promoted.Status)
	}
	warn.Status = promoted.Status
	if !reflect.DeepEqual(warn, promoted) {
		t.Fatalf("promotion changed findings or info: %#v, %#v", warn, promoted)
	}
}
