// Package decisionview presents a fingerprint-bound person's decision. Its
// current/stale text and JSON shape can be shared by approvals and lane waivers.
package decisionview

import (
	"fmt"

	"github.com/mbell-nzl/fid/internal/display"
)

// View leaves absent evidence out of a missing decision. A nil View means the
// decision does not apply or its freshness could not be determined.
type View struct {
	State       string `json:"state"`
	DecisionID  string `json:"decision_id,omitempty"`
	Actor       string `json:"actor,omitempty"`
	At          string `json:"at,omitempty"`
	Fingerprint string `json:"fingerprint,omitempty"`
}

// Text returns uncoloured prose; the renderer escapes and wraps it before
// applying Role. earlier explains what the stale decision covered.
func (v View) Text(earlier string) string {
	if v.State == "missing" {
		return string(display.SymbolUnstarted) + " missing"
	}
	detail := fmt.Sprintf("%s %s · %s", v.DecisionID, v.Actor, v.At)
	if v.State == "stale" {
		return string(display.SymbolWarning) + " stale · " + detail + " · " + earlier
	}
	return string(display.SymbolSuccess) + " " + detail
}

func (v View) Role() display.Role {
	switch v.State {
	case "current":
		return display.RoleSuccess
	case "stale":
		return display.RoleWarning
	default:
		return display.RoleMuted
	}
}
