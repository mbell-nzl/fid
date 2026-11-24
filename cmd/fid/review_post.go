package main

import (
	"fmt"

	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/mbell-nzl/fid/internal/reviewpost"
	"github.com/mbell-nzl/fid/internal/reviewwork"
	fidmigrations "github.com/mbell-nzl/fid/migrations"
	"github.com/spf13/cobra"
)

func init() { rootFlagCapabilities["review post"] = rootFlagCapability{Project: true, JSON: true} }

func newReviewPostCommand() *cobra.Command {
	return &cobra.Command{Use: "post <round>", Short: "Post pending canonical review comments.", Args: canonicalExactArgs("review.post", 1), RunE: func(cmd *cobra.Command, args []string) error {
		root, err := rebuildWorkspaceRoot(cmd)
		if err != nil {
			return reviewFailure(cmd, "review.post", err, "")
		}
		result, err := reviewpost.Round(cmd.Context(), root, fidmigrations.FS, args[0], nil)
		if err != nil {
			return reviewFailure(cmd, "review.post", err, "")
		}
		jsonMode, err := readJSONMode(cmd)
		if err != nil {
			return err
		}
		if jsonMode {
			return cli.WriteSuccess(cmd.OutOrStdout(), "review.post", map[string]any{"posted": result.Posted, "pending": result.Pending}, nil)
		}
		for _, id := range result.Posted {
			if _, err := fmt.Fprintln(cmd.OutOrStdout(), "✓ Posted "+cli.EscapeControlText(id)); err != nil {
				return err
			}
		}
		for _, id := range result.Pending {
			if _, err := fmt.Fprintln(cmd.OutOrStdout(), "⚠ Pending PR comment for "+cli.EscapeControlText(id)); err != nil {
				return err
			}
		}
		return nil
	}}
}

func postFindingResponse(cmd *cobra.Command, root, finding string) {
	conn, err := reviewwork.OpenProjection(root)
	if err != nil {
		return
	}
	var round string
	err = conn.QueryRow(`SELECT round_id FROM round_findings WHERE id=?`, finding).Scan(&round)
	_ = conn.Close()
	if err == nil {
		_, _ = reviewpost.Round(cmd.Context(), root, fidmigrations.FS, round, nil)
	}
}

// Recording has already succeeded. A delivery failure stays visible offline
// in status; it must never turn the canonical write into a failed command.
func postRecordedRound(cmd *cobra.Command, root, round string) {
	_, _ = reviewpost.Round(cmd.Context(), root, fidmigrations.FS, round, nil)
}
