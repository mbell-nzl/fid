package waiverwork

import (
	"fmt"
	"io"
	"strings"

	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/mbell-nzl/fid/internal/handoffrecord"
	"github.com/mbell-nzl/fid/internal/policy"
)

type Preview struct {
	TaskID, Lane string
	Handoff      handoffrecord.Record
	Floor        policy.Floor
}

func WritePrompt(w io.Writer, p Preview) error {
	var b strings.Builder
	esc := cli.EscapeControlText
	floor, hard := "none", "none"
	if p.Floor.Lane != nil {
		floor = *p.Floor.Lane
	}
	if p.Floor.HardFloor != nil {
		hard = *p.Floor.HardFloor
	}
	floorPaths, hardPaths := strings.Join(p.Floor.Paths, ", "), strings.Join(p.Floor.HardPaths, ", ")
	if floorPaths == "" {
		floorPaths = "none"
	}
	if hardPaths == "" {
		hardPaths = "none"
	}
	fmt.Fprintf(&b, "Waive %s · %s → %s\n", esc(p.TaskID), esc(p.Handoff.Lane), esc(p.Lane))
	fmt.Fprintf(&b, "Floor: %s · paths: %s\n", esc(floor), esc(floorPaths))
	fmt.Fprintf(&b, "Hard floor: %s · paths: %s\n", esc(hard), esc(hardPaths))
	fmt.Fprintf(&b, "New lane: %s\nType %s to waive: ", esc(p.Lane), esc(p.TaskID))
	_, err := io.WriteString(w, b.String())
	return err
}
